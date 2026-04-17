package tsplay_core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/playwright-community/playwright-go"
)

func TestFlowUsesPlaywright(t *testing.T) {
	tests := []struct {
		name string
		flow *Flow
		want bool
	}{
		{
			name: "data_only_flow",
			flow: &Flow{
				Steps: []FlowStep{
					{Action: "set_var", SaveAs: "answer", Value: "ok"},
					{Action: "http_request", URL: "https://example.com"},
				},
			},
			want: false,
		},
		{
			name: "http_request_false_static_var",
			flow: &Flow{
				Vars: map[string]any{
					"use_browser": false,
				},
				Steps: []FlowStep{
					{
						Action: "http_request",
						URL:    "https://example.com",
						With: map[string]any{
							"use_browser_cookies": "{{use_browser}}",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "nested_browser_step",
			flow: &Flow{
				Steps: []FlowStep{
					{
						Action: "retry",
						Steps: []FlowStep{
							{Action: "navigate", URL: "https://example.com"},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "http_request_uses_browser_state",
			flow: &Flow{
				Steps: []FlowStep{
					{
						Action: "http_request",
						URL:    "https://example.com",
						With: map[string]any{
							"use_browser_cookies": true,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "lua_without_browser_calls",
			flow: &Flow{
				Steps: []FlowStep{
					{Action: "lua", Code: "local total = 1 + 1\nreturn total"},
				},
			},
			want: false,
		},
		{
			name: "lua_with_browser_calls",
			flow: &Flow{
				Steps: []FlowStep{
					{Action: "lua", Code: "navigate('https://example.com')"},
				},
			},
			want: true,
		},
		{
			name: "browser_config_requires_runtime",
			flow: &Flow{
				Browser: &FlowBrowserConfig{SaveStorageState: "state.json"},
				Steps: []FlowStep{
					{Action: "set_var", SaveAs: "answer", Value: "ok"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flowUsesPlaywright(tt.flow); got != tt.want {
				t.Fatalf("flowUsesPlaywright() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlowActionCapabilitiesCoverAllActions(t *testing.T) {
	for action := range flowActionSpecs {
		capabilities, ok := flowActionCapabilitiesFor(action)
		if !ok {
			t.Fatalf("missing capabilities for action %q", action)
		}
		if capabilities.NeedsPage && !capabilities.RequiresPlaywright() {
			t.Fatalf("action %q needs a page but is not marked as requiring Playwright", action)
		}
		if capabilities.NeedsContext && !capabilities.RequiresPlaywright() {
			t.Fatalf("action %q needs a context but is not marked as requiring Playwright", action)
		}
		if capabilities.NeedsBrowserState && !capabilities.RequiresPlaywright() {
			t.Fatalf("action %q needs browser state but is not marked as requiring Playwright", action)
		}
	}

	for action := range flowActionCapabilitiesRegistry {
		if _, ok := flowActionSpecs[action]; !ok {
			t.Fatalf("capabilities registered for unknown action %q", action)
		}
	}
}

func TestAnalyzeFlowPlaywrightUsageReportsReasons(t *testing.T) {
	usage := AnalyzeFlowPlaywrightUsage(&Flow{
		Browser: &FlowBrowserConfig{
			SaveStorageState: "state.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
			{
				Action: "http_request",
				URL:    "https://example.com/api",
				With: map[string]any{
					"use_browser_cookies": true,
				},
			},
		},
	})

	if !usage.NeedsPlaywright || !usage.NeedsRuntime {
		t.Fatalf("usage = %#v", usage)
	}
	if !usage.NeedsPage {
		t.Fatalf("expected page requirement: %#v", usage)
	}
	if !usage.NeedsContext || !usage.NeedsBrowserState {
		t.Fatalf("expected context/browser state requirement: %#v", usage)
	}
	summary := usage.Summary(10)
	for _, want := range []string{"browser.save_storage_state", "steps[1].navigate", "steps[2].http_request.use_browser_cookies"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary %q does not contain %q", summary, want)
		}
	}
}

func TestAnalyzeLuaScriptPlaywrightUsageIgnoresCommentsAndStrings(t *testing.T) {
	usage := AnalyzeLuaScriptPlaywrightUsage(`
-- navigate("https://example.com")
local snippet = "page:click('#submit')"
return http_request({url = "https://example.com/api"})
`)
	if usage.NeedsPlaywright {
		t.Fatalf("usage = %#v", usage)
	}
}

func TestAnalyzeLuaScriptPlaywrightUsageDetectsBrowserCalls(t *testing.T) {
	usage := AnalyzeLuaScriptPlaywrightUsage(`
http_request({
  url = "https://example.com/api",
  use_browser_cookies = true,
})
page:click("#submit")
`)
	if !usage.NeedsPlaywright || !usage.NeedsRuntime {
		t.Fatalf("usage = %#v", usage)
	}
	if !usage.NeedsPage {
		t.Fatalf("expected page requirement: %#v", usage)
	}
	if !usage.NeedsContext || !usage.NeedsBrowserState {
		t.Fatalf("expected context/browser state requirement: %#v", usage)
	}
	summary := usage.Summary(10)
	for _, want := range []string{"lua.http_request.use_browser_cookies", "lua.page"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary %q does not contain %q", summary, want)
		}
	}
}

func TestRunFlowSkipsPlaywrightForDataOnlyFlows(t *testing.T) {
	csvPath := filepath.Join(t.TempDir(), "users.csv")
	if err := os.WriteFile(csvPath, []byte("name\nalice\n"), 0600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	flows := []struct {
		name  string
		flow  *Flow
		check func(*testing.T, *FlowResult)
	}{
		{
			name: "set_var_only",
			flow: &Flow{
				SchemaVersion: CurrentFlowSchemaVersion,
				Name:          "data-only-set-var",
				Steps: []FlowStep{
					{Action: "set_var", SaveAs: "answer", Value: "ok"},
				},
			},
			check: func(t *testing.T, result *FlowResult) {
				if got := result.Vars["answer"]; got != "ok" {
					t.Fatalf("result.Vars[answer] = %#v", got)
				}
			},
		},
		{
			name: "lua_without_browser",
			flow: &Flow{
				SchemaVersion: CurrentFlowSchemaVersion,
				Name:          "data-only-lua",
				Steps: []FlowStep{
					{Action: "lua", Code: `local total = 1 + 1 return total`},
				},
			},
			check: func(t *testing.T, result *FlowResult) {
				if result.Trace[0].Output != float64(2) {
					t.Fatalf("trace output = %#v", result.Trace[0].Output)
				}
			},
		},
		{
			name: "read_csv",
			flow: &Flow{
				SchemaVersion: CurrentFlowSchemaVersion,
				Name:          "data-only-read-csv",
				Steps: []FlowStep{
					{Action: "read_csv", FilePath: csvPath, SaveAs: "rows"},
				},
			},
			check: func(t *testing.T, result *FlowResult) {
				rows, ok := result.Vars["rows"].([]any)
				if !ok || len(rows) != 1 {
					t.Fatalf("rows = %#v", result.Vars["rows"])
				}
				row, ok := rows[0].(map[string]any)
				if !ok || row["name"] != "alice" {
					t.Fatalf("row = %#v", rows[0])
				}
			},
		},
	}

	restore := stubPlaywrightRuntime(t,
		func() error {
			t.Fatalf("unexpected Playwright install")
			return nil
		},
		func() (*playwright.Playwright, error) {
			t.Fatalf("unexpected Playwright startup")
			return nil, nil
		},
	)
	defer restore()

	for _, tt := range flows {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunFlow(tt.flow, FlowRunOptions{})
			if err != nil {
				t.Fatalf("RunFlow() error = %v", err)
			}
			if result.Playwright != nil {
				t.Fatalf("expected no playwright analysis for data-only flow: %#v", result.Playwright)
			}
			tt.check(t, result)
		})
	}
}

func TestRunFlowStartsPlaywrightWhenBrowserFlowNeedsIt(t *testing.T) {
	installCalls := 0
	runCalls := 0
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalls++
			return nil
		},
		func() (*playwright.Playwright, error) {
			runCalls++
			return nil, errors.New("boom")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: CurrentFlowSchemaVersion,
		Name:          "browser-flow",
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{})
	if err == nil {
		t.Fatalf("expected RunFlow() to fail when Playwright startup fails")
	}
	if installCalls != 1 {
		t.Fatalf("installCalls = %d, want 1", installCalls)
	}
	if runCalls != 1 {
		t.Fatalf("runCalls = %d, want 1", runCalls)
	}
	if !strings.Contains(err.Error(), "steps[1].navigate") {
		t.Fatalf("expected reason in error: %v", err)
	}
	if !strings.Contains(err.Error(), "could not start Playwright: boom") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func stubPlaywrightRuntime(t *testing.T, install func() error, run func() (*playwright.Playwright, error)) func() {
	t.Helper()

	oldInstall := playwrightInstallFunc
	oldRun := playwrightRunFunc

	playwrightInstallMu.Lock()
	oldDone := playwrightInstallDone
	playwrightInstallDone = false
	playwrightInstallMu.Unlock()

	playwrightInstallFunc = install
	playwrightRunFunc = run

	return func() {
		playwrightInstallFunc = oldInstall
		playwrightRunFunc = oldRun
		playwrightInstallMu.Lock()
		playwrightInstallDone = oldDone
		playwrightInstallMu.Unlock()
	}
}
