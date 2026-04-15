package tsplay_core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLoadFlowYAMLAndValidate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "search.flow.yaml")
	content := []byte(`
schema_version: "1"
name: baidu_search
vars:
  query: 山东大学
steps:
  - action: navigate
    url: https://www.baidu.com
  - action: type_text
    selector: "#kw"
    text: "{{query}}"
  - action: click
    selector: "#su"
`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write flow: %v", err)
	}

	flow, err := LoadFlowFile(path)
	if err != nil {
		t.Fatalf("load flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if flow.Name != "baidu_search" {
		t.Fatalf("unexpected flow name: %s", flow.Name)
	}
	if len(flow.Steps) != 3 {
		t.Fatalf("unexpected step count: %d", len(flow.Steps))
	}
}

func TestBuildActionArgsResolvesVars(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	ctx := &FlowContext{Vars: map[string]any{"query": "山东大学"}}
	step := FlowStep{
		Action:   "type_text",
		Selector: "#kw",
		Text:     "{{query}}",
	}

	args, err := buildActionArgs(L, ctx, step)
	if err != nil {
		t.Fatalf("build args: %v", err)
	}
	if len(args) != 2 {
		t.Fatalf("unexpected arg count: %d", len(args))
	}
	if got := args[0].String(); got != "#kw" {
		t.Fatalf("selector = %q", got)
	}
	if got := args[1].String(); got != "山东大学" {
		t.Fatalf("text = %q", got)
	}
}

func TestRunFlowLuaSteps(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_only",
		Vars:          map[string]any{"prefix": "hello"},
		Steps: []FlowStep{
			{
				Action: "lua",
				Code:   "return prefix .. ' world'",
				SaveAs: "message",
			},
			{
				Action: "lua",
				Code:   "return message .. '!'",
				SaveAs: "final",
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["final"]; got != "hello world!" {
		t.Fatalf("final = %#v", got)
	}
	if len(result.Trace) != 2 {
		t.Fatalf("unexpected trace count: %d", len(result.Trace))
	}
	if result.Trace[0].ArgsSummary == "" {
		t.Fatalf("expected args summary")
	}
	if result.Trace[0].OutputSummary == "" {
		t.Fatalf("expected output summary")
	}
}

func TestRunFlowRetryRetriesNestedSteps(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	attempts := 0
	L.SetGlobal("flaky", L.NewFunction(func(L *lua.LState) int {
		attempts++
		if attempts < 2 {
			L.RaiseError("not yet")
			return 0
		}
		L.Push(lua.LString("ok"))
		return 1
	}))

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "retry_lua",
		Steps: []FlowStep{
			{
				Action:     "retry",
				Times:      3,
				IntervalMS: 1,
				Steps: []FlowStep{
					{Action: "lua", Code: "return flaky()", SaveAs: "flaky_result"},
				},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d", attempts)
	}
	if got := result.Vars["flaky_result"]; got != "ok" {
		t.Fatalf("flaky_result = %#v", got)
	}
	if len(result.Trace) != 1 {
		t.Fatalf("unexpected trace count: %d", len(result.Trace))
	}
	retryTrace := result.Trace[0]
	if retryTrace.Status != "ok" {
		t.Fatalf("retry status = %q", retryTrace.Status)
	}
	if len(retryTrace.Attempts) != 2 {
		t.Fatalf("retry attempts trace count = %d", len(retryTrace.Attempts))
	}
	if retryTrace.Attempts[0].Status != "error" || retryTrace.Attempts[1].Status != "ok" {
		t.Fatalf("unexpected attempt statuses: %#v", retryTrace.Attempts)
	}
	if retryTrace.Attempts[0].Attempt != 1 || retryTrace.Attempts[1].Attempt != 2 {
		t.Fatalf("unexpected attempt numbers: %#v", retryTrace.Attempts)
	}
}

func TestValidateFlowStrictRejectsMissingSchemaVersion(t *testing.T) {
	flow := &Flow{
		Name: "missing_schema",
		Steps: []FlowStep{
			{Action: "lua", Code: "return true"},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected missing schema_version error")
	}
}

func TestValidateFlowStrictRejectsArgCount(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "bad_args",
		Steps: []FlowStep{
			{Action: "type_text", Args: []any{"#kw"}},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected arg count error")
	}
}

func TestValidateFlowStrictRejectsUnknownVariable(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "bad_var",
		Steps: []FlowStep{
			{Action: "type_text", Selector: "#kw", Text: "{{missing}}"},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected unknown variable error")
	}
}

func TestValidateFlowStrictRejectsBadSaveAs(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "bad_save_as",
		Steps: []FlowStep{
			{Action: "lua", Code: "return true", SaveAs: "bad-name"},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected bad save_as error")
	}
}

func TestValidateFlowStrictRejectsTypeMismatch(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "bad_type",
		Steps: []FlowStep{
			{Action: "wait_for_selector", Selector: "#kw", With: map[string]any{"timeout": "5000"}},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected type mismatch error")
	}
}

func TestValidateFlowStrictAcceptsRetryAndAsserts(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "retry_asserts",
		Steps: []FlowStep{
			{
				Action:     "retry",
				Times:      2,
				IntervalMS: 10,
				Steps: []FlowStep{
					{Action: "assert_visible", Selector: "#ready", Timeout: 1000},
					{Action: "assert_text", Selector: "#message", Text: "done"},
					{Action: "get_text", Selector: "#message", SaveAs: "message"},
				},
			},
			{Action: "lua", Code: "return message", SaveAs: "echo"},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestValidateFlowStrictRejectsRetryWithoutSteps(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "bad_retry",
		Steps: []FlowStep{
			{Action: "retry", Times: 2},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected retry nested steps error")
	}
}

func TestValidateFlowSecurityChecksRetryNestedSteps(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "nested_lua_policy",
		Steps: []FlowStep{
			{
				Action: "retry",
				Steps: []FlowStep{
					{Action: "lua", Code: "return true"},
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy()); err == nil {
		t.Fatalf("expected nested lua security policy error")
	}
}

func TestValidateFlowSecurityRejectsLuaByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_policy",
		Steps: []FlowStep{
			{Action: "lua", Code: "return true"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy()); err == nil {
		t.Fatalf("expected lua security policy error")
	}
}

func TestValidateFlowSecurityAllowsTrustedLua(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_policy",
		Steps: []FlowStep{
			{Action: "lua", Code: "return true"},
		},
	}

	if err := ValidateFlowSecurity(flow, TrustedFlowSecurityPolicy()); err != nil {
		t.Fatalf("validate security: %v", err)
	}
}

func TestValidateFlowSecurityRejectsBrowserStateByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_state_policy",
		Steps: []FlowStep{
			{Action: "get_cookies_string"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy()); err == nil {
		t.Fatalf("expected browser state security policy error")
	}
}

func TestValidateFlowSecurityRestrictsFileOutputRoot(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "file_output_policy",
		Steps: []FlowStep{
			{Action: "screenshot", Path: "../escape.png"},
		},
	}
	policy := FlowSecurityPolicy{
		AllowFileAccess: true,
		FileOutputRoot:  t.TempDir(),
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, policy); err == nil {
		t.Fatalf("expected file output root error")
	}
}

func TestRewriteFlowFileAccessArgsUsesOutputRoot(t *testing.T) {
	root := t.TempDir()
	policy := &FlowSecurityPolicy{
		AllowFileAccess: true,
		FileOutputRoot:  root,
	}
	args := []lua.LValue{lua.LString("screens/shot.png")}

	rewritten, err := rewriteFlowFileAccessArgs(FlowStep{Action: "screenshot"}, args, policy)
	if err != nil {
		t.Fatalf("rewrite args: %v", err)
	}
	got := rewritten[0].String()
	if !filepath.IsAbs(got) {
		t.Fatalf("expected absolute output path, got %q", got)
	}
	rootReal, err := prepareRuntimeFileRoot(root)
	if err != nil {
		t.Fatalf("prepare root: %v", err)
	}
	if err := ensurePathInsideRoot(got, rootReal); err != nil {
		t.Fatalf("rewritten path outside root: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(got)); err != nil {
		t.Fatalf("expected output directory to be created: %v", err)
	}
}

func TestRunFlowAssertVisibleAndText(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "assert_browser",
		Steps: []FlowStep{
			{Action: "navigate", URL: `data:text/html,<html><body><div id="ready">Order complete</div></body></html>`},
			{Action: "assert_visible", Selector: "#ready", Timeout: 1000},
			{Action: "assert_text", Selector: "#ready", Text: "complete", Timeout: 1000},
		},
	}

	result, err := RunFlow(flow, FlowRunOptions{Headless: true})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if result == nil || len(result.Trace) != 3 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.Trace[1].Output != true {
		t.Fatalf("assert_visible output = %#v", result.Trace[1].Output)
	}
	if result.Trace[2].OutputSummary == "" {
		t.Fatalf("expected assert_text output summary")
	}
}

func TestRunFlowTraceCapturesLuaFailureDetails(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	artifactRoot := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_failure",
		Steps: []FlowStep{
			{Action: "lua", Code: "error('boom')"},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{ArtifactRoot: artifactRoot})
	if err == nil {
		t.Fatalf("expected run error")
	}
	if result == nil || len(result.Trace) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	trace := result.Trace[0]
	if trace.Status != "error" {
		t.Fatalf("status = %q", trace.Status)
	}
	if trace.ArgsSummary == "" || !strings.Contains(trace.ArgsSummary, "boom") {
		t.Fatalf("unexpected args summary: %q", trace.ArgsSummary)
	}
	if trace.ErrorStack == "" {
		t.Fatalf("expected error stack")
	}
	if trace.Artifacts == nil || trace.Artifacts.Directory == "" {
		t.Fatalf("expected artifact directory")
	}
	if !strings.Contains(trace.Artifacts.CaptureError, "page") {
		t.Fatalf("expected page capture error, got %#v", trace.Artifacts)
	}
}

func TestRunFlowCapturesFailureArtifacts(t *testing.T) {
	artifactRoot := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_failure",
		Steps: []FlowStep{
			{Action: "navigate", URL: "data:text/html,<html><body><h1>Hello</h1></body></html>"},
			{Action: "wait_for_selector", Selector: "#missing", Timeout: 50},
		},
	}

	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: artifactRoot,
	})
	if err == nil {
		t.Fatalf("expected run error")
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	trace := result.Trace[1]
	if trace.Status != "error" {
		t.Fatalf("status = %q", trace.Status)
	}
	if trace.PageURL == "" {
		t.Fatalf("expected page url")
	}
	if trace.Artifacts == nil {
		t.Fatalf("expected artifacts")
	}
	for name, path := range map[string]string{
		"screenshot": trace.Artifacts.ScreenshotPath,
		"html":       trace.Artifacts.HTMLPath,
		"dom":        trace.Artifacts.DOMSnapshotPath,
	} {
		if path == "" {
			t.Fatalf("expected %s artifact path: %#v", name, trace.Artifacts)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s artifact file %s: %v", name, path, err)
		}
	}
}
