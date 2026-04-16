package tsplay_core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestRunFlowIfBranchesOnCondition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "if_branch",
		Vars:          map[string]any{"should_run": true},
		Steps: []FlowStep{
			{
				Action:    "if",
				Condition: &FlowStep{Action: "lua", Code: "return should_run"},
				Then: []FlowStep{
					{Action: "lua", Code: "return 'then'", SaveAs: "branch"},
				},
				Else: []FlowStep{
					{Action: "lua", Code: "return 'else'", SaveAs: "branch"},
				},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["branch"]; got != "then" {
		t.Fatalf("branch = %#v", got)
	}
	if result.Trace[0].Branch != "then" {
		t.Fatalf("trace branch = %q", result.Trace[0].Branch)
	}
	if result.Trace[0].Condition == nil || result.Trace[0].Condition.Status != "ok" {
		t.Fatalf("expected condition trace: %#v", result.Trace[0].Condition)
	}
	if len(result.Trace[0].Children) != 1 {
		t.Fatalf("expected branch child trace: %#v", result.Trace[0].Children)
	}
}

func TestRunFlowForeachIteratesItems(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "foreach_items",
		Vars: map[string]any{
			"numbers": []any{1, 2, 3},
			"total":   0,
		},
		Steps: []FlowStep{
			{
				Action:   "foreach",
				Items:    "{{numbers}}",
				ItemVar:  "number",
				IndexVar: "number_index",
				Steps: []FlowStep{
					{Action: "lua", Code: "return total + number", SaveAs: "total"},
				},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["total"]; got != float64(6) {
		t.Fatalf("total = %#v", got)
	}
	if _, ok := result.Vars["number"]; ok {
		t.Fatalf("item var leaked: %#v", result.Vars["number"])
	}
	if len(result.Trace[0].Children) != 3 {
		t.Fatalf("expected foreach child traces: %#v", result.Trace[0].Children)
	}
	if result.Trace[0].Children[2].Iteration != 3 {
		t.Fatalf("expected iteration marker: %#v", result.Trace[0].Children[2])
	}
}

func TestRunFlowOnErrorHandlesFailure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "on_error_handler",
		Steps: []FlowStep{
			{
				Action: "on_error",
				Steps: []FlowStep{
					{Action: "lua", Code: "error('boom')"},
				},
				OnError: []FlowStep{
					{Action: "lua", Code: "return last_error", SaveAs: "handled_error"},
				},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["handled_error"]; !strings.Contains(fmt.Sprint(got), "boom") {
		t.Fatalf("handled_error = %#v", got)
	}
	if result.Trace[0].Status != "ok" || result.Trace[0].Branch != "on_error" {
		t.Fatalf("unexpected on_error trace: %#v", result.Trace[0])
	}
	if len(result.Trace[0].Children) != 2 || result.Trace[0].Children[0].Status != "error" {
		t.Fatalf("expected failing child plus handler trace: %#v", result.Trace[0].Children)
	}
}

func TestRunFlowWaitUntilPollsCondition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "wait_until_condition",
		Steps: []FlowStep{
			{
				Action:     "wait_until",
				Timeout:    1000,
				IntervalMS: 1,
				Condition:  &FlowStep{Action: "lua", Code: "counter = (counter or 0) + 1; return counter >= 3"},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if len(result.Trace[0].Attempts) != 3 {
		t.Fatalf("wait attempts = %d", len(result.Trace[0].Attempts))
	}
	if result.Trace[0].Attempts[2].Status != "ok" {
		t.Fatalf("unexpected final attempt: %#v", result.Trace[0].Attempts[2])
	}
}

func TestRunFlowExtractTextAndSetVar(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("get_text", L.NewFunction(func(L *lua.LState) int {
		selector := L.CheckString(1)
		switch selector {
		case ".summary .count":
			L.Push(lua.LString("Orders: 12"))
		case ".empty-state":
			L.Push(lua.LString("No orders"))
		default:
			L.Push(lua.LString(""))
		}
		return 1
	}))

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "extract_text_set_var",
		Steps: []FlowStep{
			{Action: "extract_text", Selector: ".summary .count", Pattern: `([0-9]+)`, SaveAs: "order_count"},
			{Action: "set_var", SaveAs: "export_message", Value: "Current orders: {{order_count}}"},
			{
				Action: "if",
				Condition: &FlowStep{
					Action:   "extract_text",
					Selector: ".summary .count",
					Pattern:  `[1-9][0-9]*`,
				},
				Then: []FlowStep{
					{Action: "set_var", SaveAs: "should_export", Value: "yes"},
				},
				Else: []FlowStep{
					{Action: "set_var", SaveAs: "should_export", Value: "no"},
				},
			},
		},
	}

	result, err := RunFlowInState(L, flow)
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["order_count"]; got != "12" {
		t.Fatalf("order_count = %#v", got)
	}
	if got := result.Vars["export_message"]; got != "Current orders: 12" {
		t.Fatalf("export_message = %#v", got)
	}
	if got := result.Vars["should_export"]; got != "yes" {
		t.Fatalf("should_export = %#v", got)
	}
	if result.Trace[0].OutputSummary == "" {
		t.Fatalf("expected extract_text output summary")
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
					{Action: "extract_text", Selector: "#message", Pattern: "done", SaveAs: "message"},
				},
			},
			{Action: "set_var", SaveAs: "message_copy", Value: "{{message}}"},
			{
				Action:    "if",
				Condition: &FlowStep{Action: "is_visible", Selector: "#optional"},
				Then: []FlowStep{
					{Action: "click", Selector: "#optional"},
				},
			},
			{
				Action:  "foreach",
				Items:   []any{"a", "b"},
				ItemVar: "row",
				Steps: []FlowStep{
					{Action: "lua", Code: "return row", SaveAs: "last_row"},
				},
			},
			{
				Action:    "wait_until",
				Condition: &FlowStep{Action: "is_visible", Selector: "#done"},
				Timeout:   1000,
			},
			{
				Action: "on_error",
				Steps: []FlowStep{
					{Action: "click", Selector: "#maybe"},
				},
				OnError: []FlowStep{
					{Action: "lua", Code: "return last_error", SaveAs: "handled"},
				},
			},
			{Action: "lua", Code: "return message_copy", SaveAs: "echo"},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestValidateFlowStrictRejectsControlMissingRequiredFields(t *testing.T) {
	for name, step := range map[string]FlowStep{
		"if":         {Action: "if"},
		"foreach":    {Action: "foreach", Items: []any{"a"}, ItemVar: "item"},
		"on_error":   {Action: "on_error", Steps: []FlowStep{{Action: "lua", Code: "return true"}}},
		"wait_until": {Action: "wait_until"},
	} {
		flow := &Flow{
			SchemaVersion: "1",
			Name:          "bad_" + name,
			Steps:         []FlowStep{step},
		}
		if err := ValidateFlowStrict(flow); err == nil {
			t.Fatalf("expected %s validation error", name)
		}
	}
}

func TestValidateFlowStrictRejectsSetVarWithoutSaveAsOrValue(t *testing.T) {
	for name, step := range map[string]FlowStep{
		"missing_save_as": {Action: "set_var", Value: "hello"},
		"missing_value":   {Action: "set_var", SaveAs: "greeting"},
	} {
		flow := &Flow{
			SchemaVersion: "1",
			Name:          name,
			Steps:         []FlowStep{step},
		}
		if err := ValidateFlowStrict(flow); err == nil {
			t.Fatalf("expected %s validation error", name)
		}
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

func TestValidateFlowStrictAcceptsBrowserConfig(t *testing.T) {
	headless := true
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_config",
		Browser: &FlowBrowserConfig{
			Headless:         &headless,
			StorageState:     "states/admin.json",
			SaveStorageState: "states/admin-latest.json",
			Timeout:          30000,
			UserAgent:        "tsplay-test/1.0",
			Viewport:         &FlowViewport{Width: 1440, Height: 900},
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestValidateFlowStrictAcceptsUseSession(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_use_session",
		Browser: &FlowBrowserConfig{
			UseSession: "admin",
			Timeout:    30000,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestValidateFlowStrictRejectsPersistentWithStorageState(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_config_conflict",
		Browser: &FlowBrowserConfig{
			Persistent:   true,
			StorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected browser config conflict error")
	}
}

func TestValidateFlowStrictRejectsUseSessionWithStorageState(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_config_use_session_conflict",
		Browser: &FlowBrowserConfig{
			UseSession:   "admin",
			StorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlowStrict(flow); err == nil {
		t.Fatalf("expected browser config conflict error")
	}
}

func TestValidateFlowSecurityRejectsFlowBrowserStateByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_config_policy",
		Browser: &FlowBrowserConfig{
			StorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy()); err == nil {
		t.Fatalf("expected browser config security policy error")
	}
}

func TestValidateFlowSecurityRejectsNamedSessionByDefault(t *testing.T) {
	root := t.TempDir()
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     root,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("save session: %v", err)
	}
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_named_session_policy",
		Browser: &FlowBrowserConfig{
			UseSession: "admin",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{
		FileInputRoot:  root,
		FileOutputRoot: root,
	}); err == nil {
		t.Fatalf("expected named session security policy error")
	}
}

func TestValidateFlowSecurityRestrictsBrowserStateRoot(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_state_root",
		Browser: &FlowBrowserConfig{
			StorageState: "../escape.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}
	policy := FlowSecurityPolicy{
		AllowBrowserState: true,
		FileInputRoot:     t.TempDir(),
		FileOutputRoot:    t.TempDir(),
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, policy); err == nil {
		t.Fatalf("expected browser state root error")
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

func TestRunFlowBrowserStorageStateRoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/seed":
			fmt.Fprint(w, `<html><body><script>localStorage.setItem("token","abc123"); window.location.href="/check";</script></body></html>`)
		case "/check":
			fmt.Fprint(w, `<html><body><div id="ready"></div><script>document.getElementById("ready").textContent = localStorage.getItem("token") || "missing";</script></body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	root := t.TempDir()
	security := &FlowSecurityPolicy{
		AllowBrowserState: true,
		FileInputRoot:     root,
		FileOutputRoot:    root,
	}
	headless := true

	saveFlow := &Flow{
		SchemaVersion: "1",
		Name:          "save_browser_state",
		Browser: &FlowBrowserConfig{
			Headless:         &headless,
			SaveStorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL + "/seed"},
			{Action: "assert_text", Selector: "#ready", Text: "abc123", Timeout: 3000},
		},
	}

	if _, err := RunFlow(saveFlow, FlowRunOptions{
		Headless:     true,
		Security:     security,
		ArtifactRoot: root,
	}); err != nil {
		t.Fatalf("save flow: %v", err)
	}
	statePath := filepath.Join(root, "states", "admin.json")
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("expected storage state file: %v", err)
	}

	loadFlow := &Flow{
		SchemaVersion: "1",
		Name:          "load_browser_state",
		Browser: &FlowBrowserConfig{
			Headless:     &headless,
			StorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL + "/check"},
			{Action: "assert_text", Selector: "#ready", Text: "abc123", Timeout: 3000},
		},
	}

	result, err := RunFlow(loadFlow, FlowRunOptions{
		Headless:     true,
		Security:     security,
		ArtifactRoot: root,
	})
	if err != nil {
		t.Fatalf("load flow: %v", err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected load result: %#v", result)
	}
}

func TestRunFlowBrowserUseSessionRoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/seed":
			fmt.Fprint(w, `<html><body><script>localStorage.setItem("token","named-session"); window.location.href="/check";</script></body></html>`)
		case "/check":
			fmt.Fprint(w, `<html><body><div id="ready"></div><script>document.getElementById("ready").textContent = localStorage.getItem("token") || "missing";</script></body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	root := t.TempDir()
	security := &FlowSecurityPolicy{
		AllowBrowserState: true,
		FileInputRoot:     root,
		FileOutputRoot:    root,
	}
	headless := true

	saveFlow := &Flow{
		SchemaVersion: "1",
		Name:          "save_browser_state_for_named_session",
		Browser: &FlowBrowserConfig{
			Headless:         &headless,
			SaveStorageState: "states/admin.json",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL + "/seed"},
			{Action: "assert_text", Selector: "#ready", Text: "named-session", Timeout: 3000},
		},
	}

	if _, err := RunFlow(saveFlow, FlowRunOptions{
		Headless:     true,
		Security:     security,
		ArtifactRoot: root,
	}); err != nil {
		t.Fatalf("save flow: %v", err)
	}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     root,
		StorageStatePath: "states/admin.json",
	}); err != nil {
		t.Fatalf("register named session: %v", err)
	}

	loadFlow := &Flow{
		SchemaVersion: "1",
		Name:          "load_browser_state_from_named_session",
		Browser: &FlowBrowserConfig{
			UseSession: "admin",
			Headless:   &headless,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL + "/check"},
			{Action: "assert_text", Selector: "#ready", Text: "named-session", Timeout: 3000},
		},
	}

	result, err := RunFlow(loadFlow, FlowRunOptions{
		Headless:     true,
		Security:     security,
		ArtifactRoot: root,
	})
	if err != nil {
		t.Fatalf("load flow: %v", err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected load result: %#v", result)
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
