package tsplay_core

import (
	"errors"
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

func TestRunFlowSkipsPlaywrightForNonBrowserFlow(t *testing.T) {
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

	result, err := RunFlow(&Flow{
		SchemaVersion: CurrentFlowSchemaVersion,
		Name:          "data-only",
		Steps: []FlowStep{
			{Action: "set_var", SaveAs: "answer", Value: "ok"},
		},
	}, FlowRunOptions{})
	if err != nil {
		t.Fatalf("RunFlow() error = %v", err)
	}
	if got := result.Vars["answer"]; got != "ok" {
		t.Fatalf("result.Vars[answer] = %#v", got)
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
