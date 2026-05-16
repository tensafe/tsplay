package tsplay_core

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
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

func TestParseFlowRejectsUnknownStepFieldWithSuggestion(t *testing.T) {
	_, err := ParseFlow([]byte(`
schema_version: "1"
name: bad_alias
steps:
  - action: evaluate
    selector: "div.g"
    script: return []
    result_var: rows
`), "yaml")
	if err == nil {
		t.Fatalf("expected parse error")
	}

	var parseErr *FlowParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FlowParseError, got %T", err)
	}
	if parseErr.Issue.Code != "unknown_field" {
		t.Fatalf("unexpected issue: %#v", parseErr.Issue)
	}
	if parseErr.Issue.StepPath != "1" {
		t.Fatalf("unexpected step path: %#v", parseErr.Issue.StepPath)
	}
	if parseErr.Issue.Field != "result_var" || parseErr.Issue.DidYouMean != "save_as" {
		t.Fatalf("unexpected suggestion: %#v", parseErr.Issue)
	}
}

func TestParseFlowRejectsDottedStepFieldWithSuggestion(t *testing.T) {
	_, err := ParseFlow([]byte(`
schema_version: "1"
name: bad_nested
steps:
  - action: write_csv
    file_path: reports/out.csv
    value: "{{rows}}"
    with.headers:
      - title
`), "yaml")
	if err == nil {
		t.Fatalf("expected parse error")
	}

	var parseErr *FlowParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FlowParseError, got %T", err)
	}
	if parseErr.Issue.Field != "with.headers" || parseErr.Issue.DidYouMean != "with.headers" {
		t.Fatalf("unexpected dotted-field issue: %#v", parseErr.Issue)
	}
}

func TestParseFlowRejectsBooleanBrowserUseSessionYAML(t *testing.T) {
	_, err := ParseFlow([]byte(`
schema_version: "1"
name: bad_browser_session_type
browser:
  use_session: true
steps:
  - action: navigate
    url: https://example.com
`), "yaml")
	if err == nil {
		t.Fatalf("expected parse error")
	}

	var parseErr *FlowParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FlowParseError, got %T", err)
	}
	if parseErr.Issue.Code != "invalid_type" {
		t.Fatalf("unexpected issue: %#v", parseErr.Issue)
	}
	if parseErr.Issue.Field != "use_session" {
		t.Fatalf("unexpected field: %#v", parseErr.Issue)
	}
	if !strings.Contains(parseErr.Issue.Message, "must be a string, got boolean") {
		t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
	}
}

func TestParseFlowRejectsBooleanBrowserUseSessionJSON(t *testing.T) {
	_, err := ParseFlow([]byte(`{
  "schema_version": "1",
  "name": "bad_browser_session_type_json",
  "browser": {
    "use_session": true
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`), "json")
	if err == nil {
		t.Fatalf("expected parse error")
	}

	var parseErr *FlowParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FlowParseError, got %T", err)
	}
	if parseErr.Issue.Code != "invalid_type" {
		t.Fatalf("unexpected issue: %#v", parseErr.Issue)
	}
	if parseErr.Issue.Field != "use_session" {
		t.Fatalf("unexpected field: %#v", parseErr.Issue)
	}
}

func TestParseFlowRejectsInvalidBrowserCDPPortYAML(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{name: "zero", value: "0", wantErr: "between 1 and 65535"},
		{name: "negative", value: "-1", wantErr: "between 1 and 65535"},
		{name: "overflow", value: "65536", wantErr: "between 1 and 65535"},
		{name: "float", value: "9222.5", wantErr: "must be an integer"},
		{name: "string", value: `"9222"`, wantErr: "must be an integer"},
		{name: "boolean", value: "true", wantErr: "must be an integer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(fmt.Sprintf(`
schema_version: "1"
name: bad_browser_cdp_port
browser:
  cdp_port: %s
steps:
  - action: navigate
    url: https://example.com
`, tt.value)), "yaml")
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != "cdp_port" {
				t.Fatalf("unexpected field: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsInvalidBrowserCDPPortJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{name: "zero", value: "0", wantErr: "between 1 and 65535"},
		{name: "negative", value: "-1", wantErr: "between 1 and 65535"},
		{name: "overflow", value: "65536", wantErr: "between 1 and 65535"},
		{name: "float", value: "9222.5", wantErr: "must be an integer"},
		{name: "string", value: `"9222"`, wantErr: "must be an integer"},
		{name: "boolean", value: "true", wantErr: "must be an integer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(fmt.Sprintf(`{
  "schema_version": "1",
  "name": "bad_browser_cdp_port_json",
  "browser": {
    "cdp_port": %s
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`, tt.value)), "json")
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != "cdp_port" {
				t.Fatalf("unexpected field: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsBlankBrowserCDPEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
	}{
		{
			name:   "yaml_empty_string",
			format: "yaml",
			body: `
schema_version: "1"
name: blank_browser_cdp_endpoint
browser:
  cdp_endpoint: ""
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "yaml_whitespace_string",
			format: "yaml",
			body: `
schema_version: "1"
name: blank_browser_cdp_endpoint
browser:
  cdp_endpoint: "   "
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "json_empty_string",
			format: "json",
			body: `{
  "schema_version": "1",
  "name": "blank_browser_cdp_endpoint_json",
  "browser": {
    "cdp_endpoint": ""
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(tt.body), tt.format)
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != "cdp_endpoint" || parseErr.Issue.Code != "invalid_value" {
				t.Fatalf("unexpected issue: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, "cannot be blank") {
				t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsNonStringBrowserCDPEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
	}{
		{
			name:   "yaml_number",
			format: "yaml",
			body: `
schema_version: "1"
name: bad_browser_cdp_endpoint
browser:
  cdp_endpoint: 9222
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "json_boolean",
			format: "json",
			body: `{
  "schema_version": "1",
  "name": "bad_browser_cdp_endpoint_json",
  "browser": {
    "cdp_endpoint": true
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(tt.body), tt.format)
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != "cdp_endpoint" || parseErr.Issue.Code != "invalid_type" {
				t.Fatalf("unexpected issue: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, "must be a string") {
				t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsNonStringBrowserCDPPathFields(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
		field  string
	}{
		{
			name:   "yaml_cdp_executable_number",
			format: "yaml",
			field:  "cdp_executable",
			body: `
schema_version: "1"
name: bad_browser_cdp_executable
browser:
  cdp_executable: 123
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "yaml_cdp_user_data_dir_boolean",
			format: "yaml",
			field:  "cdp_user_data_dir",
			body: `
schema_version: "1"
name: bad_browser_cdp_user_data_dir
browser:
  cdp_user_data_dir: true
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "json_cdp_executable_number",
			format: "json",
			field:  "cdp_executable",
			body: `{
  "schema_version": "1",
  "name": "bad_browser_cdp_executable_json",
  "browser": {
    "cdp_executable": 123
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
		{
			name:   "json_cdp_user_data_dir_boolean",
			format: "json",
			field:  "cdp_user_data_dir",
			body: `{
  "schema_version": "1",
  "name": "bad_browser_cdp_user_data_dir_json",
  "browser": {
    "cdp_user_data_dir": true
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(tt.body), tt.format)
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != tt.field || parseErr.Issue.Code != "invalid_type" {
				t.Fatalf("unexpected issue: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, "must be a string") {
				t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsBlankBrowserCDPPathFields(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
		field  string
	}{
		{
			name:   "yaml_cdp_executable_empty",
			format: "yaml",
			field:  "cdp_executable",
			body: `
schema_version: "1"
name: blank_browser_cdp_executable
browser:
  cdp_executable: ""
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "yaml_cdp_user_data_dir_whitespace",
			format: "yaml",
			field:  "cdp_user_data_dir",
			body: `
schema_version: "1"
name: blank_browser_cdp_user_data_dir
browser:
  cdp_user_data_dir: "   "
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "json_cdp_executable_empty",
			format: "json",
			field:  "cdp_executable",
			body: `{
  "schema_version": "1",
  "name": "blank_browser_cdp_executable_json",
  "browser": {
    "cdp_executable": ""
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
		{
			name:   "json_cdp_user_data_dir_empty",
			format: "json",
			field:  "cdp_user_data_dir",
			body: `{
  "schema_version": "1",
  "name": "blank_browser_cdp_user_data_dir_json",
  "browser": {
    "cdp_user_data_dir": ""
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(tt.body), tt.format)
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != tt.field || parseErr.Issue.Code != "invalid_value" {
				t.Fatalf("unexpected issue: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, "cannot be blank") {
				t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
			}
		})
	}
}

func TestParseFlowRejectsNonBooleanBrowserCDPLaunch(t *testing.T) {
	tests := []struct {
		name   string
		format string
		body   string
	}{
		{
			name:   "yaml_quoted_true",
			format: "yaml",
			body: `
schema_version: "1"
name: bad_browser_cdp_launch
browser:
  cdp_launch: "true"
steps:
  - action: navigate
    url: https://example.com
`,
		},
		{
			name:   "json_string_true",
			format: "json",
			body: `{
  "schema_version": "1",
  "name": "bad_browser_cdp_launch_json",
  "browser": {
    "cdp_launch": "true"
  },
  "steps": [
    {
      "action": "navigate",
      "url": "https://example.com"
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlow([]byte(tt.body), tt.format)
			if err == nil {
				t.Fatalf("expected parse error")
			}
			var parseErr *FlowParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("expected FlowParseError, got %T", err)
			}
			if parseErr.Issue.Field != "cdp_launch" || parseErr.Issue.Code != "invalid_type" {
				t.Fatalf("unexpected issue: %#v", parseErr.Issue)
			}
			if !strings.Contains(parseErr.Issue.Message, "must be a boolean") {
				t.Fatalf("unexpected message: %q", parseErr.Issue.Message)
			}
		})
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

func TestRunFlowAppendVarAndWriteResults(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "append_and_write_results",
		Steps: []FlowStep{
			{
				Action: "append_var",
				SaveAs: "import_results",
				With: map[string]any{
					"value": map[string]any{
						"source_row": 2,
						"status":     "success",
					},
				},
			},
			{
				Action: "append_var",
				SaveAs: "import_results",
				With: map[string]any{
					"value": map[string]any{
						"source_row": 3,
						"status":     "failed",
						"error":      "boom",
					},
				},
			},
			{
				Action:   "write_json",
				FilePath: "reports/import-results.json",
				With: map[string]any{
					"value": "{{import_results}}",
				},
			},
			{
				Action:   "write_csv",
				FilePath: "reports/import-results.csv",
				With: map[string]any{
					"value": "{{import_results}}",
					"headers": []any{
						"source_row",
						"status",
						"error",
					},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	items, ok := result.Vars["import_results"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("import_results = %#v", result.Vars["import_results"])
	}

	jsonContent, err := os.ReadFile(filepath.Join(root, "reports", "import-results.json"))
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	if !strings.Contains(string(jsonContent), "\"status\": \"failed\"") {
		t.Fatalf("unexpected json content: %s", string(jsonContent))
	}

	csvContent, err := os.ReadFile(filepath.Join(root, "reports", "import-results.csv"))
	if err != nil {
		t.Fatalf("read csv: %v", err)
	}
	if !strings.Contains(string(csvContent), "source_row,status,error") {
		t.Fatalf("missing csv header: %s", string(csvContent))
	}
	if !strings.Contains(string(csvContent), "3,failed,boom") {
		t.Fatalf("missing csv row: %s", string(csvContent))
	}
}

func TestRunFlowWriteExcelAndReadBack(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "write_excel_round_trip",
		Steps: []FlowStep{
			{
				Action: "append_var",
				SaveAs: "import_results",
				With: map[string]any{
					"value": map[string]any{
						"source_row": 2,
						"status":     "success",
					},
				},
			},
			{
				Action: "append_var",
				SaveAs: "import_results",
				With: map[string]any{
					"value": map[string]any{
						"source_row": 3,
						"status":     "failed",
						"error":      "boom",
					},
				},
			},
			{
				Action: "write_excel",
				Args: []any{
					"reports/import-results.xlsx",
					"{{import_results}}",
					[]any{"source_row", "status", "error"},
					"Results",
				},
			},
			{
				Action:   "read_excel",
				FilePath: "reports/import-results.xlsx",
				Sheet:    "Results",
				SaveAs:   "reloaded",
			},
			{
				Action: "set_var",
				SaveAs: "reloaded_error",
				Value:  "{{reloaded[1].error}}",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	xlsxPath := filepath.Join(root, "reports", "import-results.xlsx")
	if _, err := os.Stat(xlsxPath); err != nil {
		t.Fatalf("stat xlsx: %v", err)
	}
	if got := result.Vars["reloaded_error"]; got != "boom" {
		t.Fatalf("reloaded_error = %#v", got)
	}
}

func TestRunFlowWriteExcelWorkbookAndReadNamedSheet(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "write_excel_workbook_round_trip",
		Steps: []FlowStep{
			{
				Action:   "write_excel",
				FilePath: "reports/workbook.xlsx",
				With: map[string]any{
					"value": map[string]any{
						"sheets": []any{
							map[string]any{
								"name": "Summary",
								"value": []any{
									map[string]any{
										"count":  2,
										"active": true,
									},
								},
							},
							map[string]any{
								"name":    "Errors",
								"headers": []any{"source_row", "error"},
								"value": []any{
									[]any{3, "boom"},
								},
							},
						},
					},
				},
			},
			{
				Action:   "read_excel",
				FilePath: "reports/workbook.xlsx",
				Sheet:    "Errors",
				SaveAs:   "errors",
			},
			{
				Action: "set_var",
				SaveAs: "error_message",
				Value:  "{{errors[0].error}}",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	if got := result.Vars["error_message"]; got != "boom" {
		t.Fatalf("error_message = %#v", got)
	}
}

func TestValidateFlowSecurityRejectsForeachProgressCheckpointWithoutAllow(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "foreach_checkpoint_policy",
		Vars:          map[string]any{"rows": []any{map[string]any{"source_row": 2}}},
		Steps: []FlowStep{
			{
				Action:  "foreach",
				Items:   "{{rows}}",
				ItemVar: "row",
				With: map[string]any{
					"progress_key": "imports:users:resume_row",
				},
				Steps: []FlowStep{
					{Action: "set_var", SaveAs: "last_row", Value: "{{row.source_row}}"},
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected redis security policy error")
	}
	if !strings.Contains(err.Error(), "allow_redis") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowForeachProgressCheckpointSkipsWithoutRedisConfig(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "foreach_checkpoint_skip",
		Vars: map[string]any{
			"rows": []any{
				map[string]any{"source_row": 2, "name": "Alice"},
				map[string]any{"source_row": 3, "name": "Bob"},
			},
		},
		Steps: []FlowStep{
			{
				Action:  "foreach",
				Items:   "{{rows}}",
				ItemVar: "row",
				With: map[string]any{
					"progress_key": "imports:users:resume_row",
				},
				Steps: []FlowStep{
					{Action: "append_var", SaveAs: "processed_rows", Value: "{{row.source_row}}"},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowRedis: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	output, ok := result.Trace[0].Output.(map[string]any)
	if !ok {
		t.Fatalf("foreach output = %#v", result.Trace[0].Output)
	}
	checkpoint, ok := output["checkpoint"].(map[string]any)
	if !ok {
		t.Fatalf("checkpoint summary = %#v", result.Trace[0].Output)
	}
	if checkpoint["status"] != "skipped" {
		t.Fatalf("checkpoint status = %#v", checkpoint["status"])
	}
	if checkpoint["reason"] != "redis connection not configured" {
		t.Fatalf("checkpoint reason = %#v", checkpoint["reason"])
	}
	if checkpoint["writes"] != 0 {
		t.Fatalf("checkpoint writes = %#v", checkpoint["writes"])
	}
	if got := fmt.Sprint(result.Vars["processed_rows"]); !strings.Contains(got, "2") || !strings.Contains(got, "3") {
		t.Fatalf("processed_rows = %#v", result.Vars["processed_rows"])
	}
}

func TestRunFlowForeachProgressCheckpointWritesNextSourceRow(t *testing.T) {
	server := newRedisTestServer(t)
	defer server.Close()

	t.Setenv("TSPLAY_REDIS_ADDR", server.Addr())

	root := t.TempDir()
	csvPath := filepath.Join(root, "users.csv")
	if err := os.WriteFile(csvPath, []byte("name,phone\nAlice,13800000000\nBob,13900000000\nCarol,13700000000\n"), 0644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "foreach_checkpoint_write",
		Steps: []FlowStep{
			{
				Action:   "read_csv",
				FilePath: "users.csv",
				SaveAs:   "rows",
				With: map[string]any{
					"row_number_field": "source_row",
				},
			},
			{
				Action:  "foreach",
				Items:   "{{rows}}",
				ItemVar: "row",
				With: map[string]any{
					"progress_key": "imports:users:resume_row",
				},
				Steps: []FlowStep{
					{Action: "append_var", SaveAs: "processed_rows", Value: "{{row.source_row}}"},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			AllowRedis:      true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	stored, err := redisGet("imports:users:resume_row", "")
	if err != nil {
		t.Fatalf("redis get checkpoint: %v", err)
	}
	if stored != "5" {
		t.Fatalf("stored checkpoint = %#v", stored)
	}

	output, ok := result.Trace[1].Output.(map[string]any)
	if !ok {
		t.Fatalf("foreach output = %#v", result.Trace[1].Output)
	}
	checkpoint, ok := output["checkpoint"].(map[string]any)
	if !ok {
		t.Fatalf("checkpoint summary = %#v", result.Trace[1].Output)
	}
	if checkpoint["status"] != "ok" {
		t.Fatalf("checkpoint status = %#v", checkpoint["status"])
	}
	if checkpoint["writes"] != 3 {
		t.Fatalf("checkpoint writes = %#v", checkpoint["writes"])
	}
	if checkpoint["last_value"] != 5 {
		t.Fatalf("checkpoint last_value = %#v", checkpoint["last_value"])
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

func TestValidateFlowStrictAcceptsDrag(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "drag_slider",
		Steps: []FlowStep{
			{
				Action:   "drag",
				Selector: "#slider",
				With: map[string]any{
					"delta_x":    120,
					"delta_y":    0,
					"move_steps": 24,
					"timeout":    5000,
				},
			},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate drag flow: %v", err)
	}
}

func TestValidateFlowStrictAcceptsClickAtAndClickBox(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "click_detection_box",
		Steps: []FlowStep{
			{
				Action:   "click_at",
				Selector: "#captcha",
				With: map[string]any{
					"x":       52,
					"y":       52,
					"timeout": 5000,
				},
			},
			{
				Action:   "click_box",
				Selector: "#captcha",
				With: map[string]any{
					"box":        map[string]any{"x1": 28, "y1": 28, "x2": 76, "y2": 76},
					"image_path": "artifacts/captcha/det-source.png",
					"auto_scale": true,
					"scale_x":    1,
					"scale_y":    1,
					"timeout":    5000,
				},
			},
		},
	}

	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate click box flow: %v", err)
	}
}

func TestRunFlowClickBoxAutoScalesFromImagePath(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "captcha-2x.png")
	img := image.NewRGBA(image.Rect(0, 0, 400, 200))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.RGBA{R: 240, G: 246, B: 252, A: 255})
		}
	}
	file, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("create scaled captcha image: %v", err)
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode scaled captcha image: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close scaled captcha image: %v", err)
	}

	page := `<!doctype html>
<html>
<head>
  <style>
    #captcha { position: relative; width: 200px; height: 100px; background: #eef2f7; }
    #target { position: absolute; left: 20px; top: 30px; width: 20px; height: 20px; background: #2da44e; }
  </style>
</head>
<body>
  <div id="captcha"><div id="target"></div></div>
  <div id="success" hidden>clicked</div>
  <script>
    const captcha = document.getElementById("captcha");
    const target = document.getElementById("target");
    const success = document.getElementById("success");
    captcha.addEventListener("click", function(event) {
      const box = target.getBoundingClientRect();
      if (event.clientX >= box.left && event.clientX <= box.right &&
          event.clientY >= box.top && event.clientY <= box.bottom) {
        success.hidden = false;
      }
    });
  </script>
</body>
</html>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page)
	}))
	defer server.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "click_box_autoscale",
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "wait_for_selector", Selector: "#captcha", Timeout: 5000},
			{
				Action:   "click_box",
				Selector: "#captcha",
				SaveAs:   "click_result",
				With: map[string]any{
					"box":        []any{40, 60, 80, 100},
					"image_path": "captcha-2x.png",
					"timeout":    5000,
				},
			},
			{Action: "assert_visible", Selector: "#success", Timeout: 5000},
		},
	}

	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: root,
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
		},
	})
	if err != nil {
		t.Fatalf("run click_box autoscale flow: %v", err)
	}
	clickResult, ok := result.Vars["click_result"].(map[string]any)
	if !ok {
		t.Fatalf("click_result = %#v", result.Vars["click_result"])
	}
	if got := clickResult["auto_scale"]; got != true {
		t.Fatalf("auto_scale = %#v", got)
	}
	if got := clickResult["scale_x"]; got != 2.0 {
		t.Fatalf("scale_x = %#v", got)
	}
	if got := clickResult["scale_y"]; got != 2.0 {
		t.Fatalf("scale_y = %#v", got)
	}
	if got := clickResult["x"]; got != 30.0 {
		t.Fatalf("click x = %#v", got)
	}
	if got := clickResult["y"]; got != 40.0 {
		t.Fatalf("click y = %#v", got)
	}
}

func TestValidateFlowStrictAcceptsRetryAndAsserts(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "retry_asserts",
		Vars: map[string]any{
			"score": 0.91,
		},
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
				Action: "assert_number",
				SaveAs: "score_gate",
				With: map[string]any{
					"value":    "{{score}}",
					"op":       ">=",
					"expected": 0.8,
					"label":    "score",
				},
			},
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

func TestRunFlowAssertNumberFailsWithLabel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "number_gate",
		Vars: map[string]any{
			"confidence": 0.42,
		},
		Steps: []FlowStep{
			{
				Action: "assert_number",
				With: map[string]any{
					"value":    "{{confidence}}",
					"op":       ">=",
					"expected": 0.8,
					"label":    "OCR confidence",
				},
			},
		},
	}

	_, err := RunFlowInState(L, flow)
	if err == nil {
		t.Fatalf("expected assert_number failure")
	}
	if !strings.Contains(err.Error(), "OCR confidence") || !strings.Contains(err.Error(), "0.42 >= 0.8") {
		t.Fatalf("unexpected assert_number error: %v", err)
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

func TestValidateFlowSecurityRejectsHTTPByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "http_policy",
		Steps: []FlowStep{
			{Action: "http_request", URL: "https://example.com/api"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected http security policy error")
	}
	if !strings.Contains(err.Error(), "allow_http") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFlowSecurityRejectsHTTPRequestFileAccessWithoutAllow(t *testing.T) {
	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "http_file_policy",
		Steps: []FlowStep{
			{
				Action: "http_request",
				URL:    "https://example.com/api",
				With: map[string]any{
					"multipart_files": map[string]any{"image": "captcha.png"},
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, FlowSecurityPolicy{
		AllowHTTP:      true,
		FileInputRoot:  root,
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected http file access security policy error")
	}
	if !strings.Contains(err.Error(), "allow_file_access") {
		t.Fatalf("unexpected error: %v", err)
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

func TestParseFlowAcceptsBrowserCDPConfig(t *testing.T) {
	flow, err := ParseFlow([]byte(`
schema_version: "1"
name: browser_cdp
browser:
  cdp_launch: true
  cdp_endpoint: "http://127.0.0.1:9222"
  cdp_executable: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  cdp_user_data_dir: "profiles/cdp-smoke"
  timeout: 30000
  viewport:
    width: 1280
    height: 720
steps:
  - action: navigate
    url: https://example.com
`), "yaml")
	if err != nil {
		t.Fatalf("parse flow: %v", err)
	}
	if flow.Browser == nil || !flow.Browser.CDPLaunch || flow.Browser.CDPEndpoint != "http://127.0.0.1:9222" || flow.Browser.CDPExecutable == "" || flow.Browser.CDPUserDataDir != "profiles/cdp-smoke" {
		t.Fatalf("unexpected browser CDP config: %#v", flow.Browser)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestParseFlowAcceptsBrowserCDPPortIntegerVariantsYAML(t *testing.T) {
	tests := map[string]int{
		"9222":   9222,
		"+9222":  9222,
		"0x2406": 9222,
		"9_222":  9222,
	}
	for value, want := range tests {
		t.Run(value, func(t *testing.T) {
			flow, err := ParseFlow([]byte(fmt.Sprintf(`
schema_version: "1"
name: browser_cdp_port_variant
browser:
  cdp_port: %s
steps:
  - action: navigate
    url: https://example.com
`, value)), "yaml")
			if err != nil {
				t.Fatalf("parse flow: %v", err)
			}
			if flow.Browser == nil || flow.Browser.CDPPort != want {
				t.Fatalf("expected cdp_port %d, got %#v", want, flow.Browser)
			}
			if err := ValidateFlowStrict(flow); err != nil {
				t.Fatalf("validate flow: %v", err)
			}
		})
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

func TestValidateFlowStrictRejectsConflictingBrowserCDPConfig(t *testing.T) {
	tests := []struct {
		name    string
		browser FlowBrowserConfig
		want    string
	}{
		{
			name:    "endpoint_and_port",
			browser: FlowBrowserConfig{CDPEndpoint: "http://127.0.0.1:9222", CDPPort: 9222},
			want:    "cannot be combined",
		},
		{
			name:    "invalid_port",
			browser: FlowBrowserConfig{CDPPort: 70000},
			want:    "cdp_port",
		},
		{
			name:    "remote_launch_endpoint",
			browser: FlowBrowserConfig{CDPLaunch: true, CDPEndpoint: "http://192.0.2.1:9222"},
			want:    "only start or reuse a local browser",
		},
		{
			name:    "local_launch_endpoint_without_port",
			browser: FlowBrowserConfig{CDPLaunch: true, CDPEndpoint: "http://127.0.0.1"},
			want:    "explicit port",
		},
		{
			name:    "storage_state",
			browser: FlowBrowserConfig{CDPPort: 9222, StorageState: "states/admin.json"},
			want:    "storage_state",
		},
		{
			name:    "launch_storage_state",
			browser: FlowBrowserConfig{CDPLaunch: true, StorageState: "states/admin.json"},
			want:    "storage_state",
		},
		{
			name:    "use_session",
			browser: FlowBrowserConfig{CDPPort: 9222, UseSession: "admin"},
			want:    "use_session",
		},
		{
			name:    "user_agent",
			browser: FlowBrowserConfig{CDPPort: 9222, UserAgent: "tsplay-test"},
			want:    "user_agent",
		},
		{
			name:    "executable_user_agent",
			browser: FlowBrowserConfig{CDPExecutable: "/tmp/chrome", UserAgent: "tsplay-test"},
			want:    "user_agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flow := &Flow{
				SchemaVersion: "1",
				Name:          "browser_cdp_conflict",
				Browser:       &tt.browser,
				Steps: []FlowStep{
					{Action: "navigate", URL: "https://example.com"},
				},
			}
			err := ValidateFlowStrict(flow)
			if err == nil {
				t.Fatalf("expected browser CDP config error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.want)
			}
		})
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

func TestValidateFlowSecurityRejectsBrowserCDPByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_cdp_policy",
		Browser: &FlowBrowserConfig{
			CDPPort: 9222,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected browser CDP security policy error")
	}
	if !strings.Contains(err.Error(), "allow_browser_state") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFlowSecurityRestrictsCDPUserDataDirRoot(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "browser_cdp_profile_policy",
		Browser: &FlowBrowserConfig{
			CDPLaunch:      true,
			CDPUserDataDir: "../escape-profile",
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
		t.Fatalf("expected CDP user data dir root policy error")
	}
}

func TestRunFlowRejectsCDPOptionWithoutBrowserStateGrant(t *testing.T) {
	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_option_policy",
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		BrowserCDPPort: 9222,
		Security:       &FlowSecurityPolicy{},
	})
	if err == nil {
		t.Fatalf("expected browser CDP option security policy error")
	}
	if !strings.Contains(err.Error(), "allow_browser_state") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowRejectsInvalidCDPOptionsBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	tests := []struct {
		name    string
		options FlowRunOptions
		want    string
	}{
		{
			name: "negative_port",
			options: FlowRunOptions{
				BrowserCDPPort: -1,
			},
			want: "between 1 and 65535",
		},
		{
			name: "overflow_port",
			options: FlowRunOptions{
				BrowserCDPPort: 65536,
			},
			want: "between 1 and 65535",
		},
		{
			name: "endpoint_and_port",
			options: FlowRunOptions{
				BrowserCDPEndpoint: "http://127.0.0.1:9222",
				BrowserCDPPort:     9222,
			},
			want: "cannot both be set",
		},
		{
			name: "invalid_endpoint",
			options: FlowRunOptions{
				BrowserCDPEndpoint: "127.0.0.1:70000/json/version",
			},
			want: "invalid port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := tt.options
			options.Security = &FlowSecurityPolicy{AllowBrowserState: true}
			_, err := RunFlow(&Flow{
				SchemaVersion: "1",
				Name:          "invalid_cdp_options",
				Steps: []FlowStep{
					{Action: "navigate", URL: "https://example.com"},
				},
			}, options)
			if err == nil {
				t.Fatalf("expected invalid CDP option error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.want)
			}
		})
	}

	if installCalled || runCalled {
		t.Fatalf("invalid CDP options should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func TestRunFlowCDPOptionsTriggerPlaywrightForNonBrowserFlow(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return nil
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("expected playwright startup")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_option_non_browser_flow",
		Steps: []FlowStep{
			{Action: "set_var", SaveAs: "answer", Value: "ok"},
		},
	}, FlowRunOptions{
		BrowserCDPPort: 9222,
		Security:       &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected Playwright startup error")
	}
	if !strings.Contains(err.Error(), "browser.cdp_port") {
		t.Fatalf("expected CDP reason in error, got %v", err)
	}
	if !installCalled || !runCalled {
		t.Fatalf("CDP options should force Playwright startup, install=%v run=%v", installCalled, runCalled)
	}
}

func TestRunFlowRejectsCDPOverrideWithUseSessionBeforeMarkingSessionOrPlaywrightStart(t *testing.T) {
	root := t.TempDir()
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     root,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("save session: %v", err)
	}

	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_override_use_session_conflict",
		Browser: &FlowBrowserConfig{
			UseSession: "admin",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		ArtifactRoot:   root,
		BrowserCDPPort: 9222,
		RunID:          "run-should-not-mark-session",
		Security:       &FlowSecurityPolicy{AllowBrowserState: true, FileInputRoot: root, FileOutputRoot: root},
	})
	if err == nil {
		t.Fatalf("expected CDP override and use_session conflict")
	}
	if !strings.Contains(err.Error(), "use_session") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("CDP override use_session conflict should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
	session, loadErr := LoadFlowSavedSession("admin", root)
	if loadErr != nil {
		t.Fatalf("load session: %v", loadErr)
	}
	if session.LastUsedAt != "" || session.LastUsedByRunID != "" {
		t.Fatalf("session should not be marked used on preflight conflict: %#v", session)
	}
}

func TestRunFlowRejectsCDPOverrideWithMissingUseSessionBeforeSessionResolution(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_override_missing_use_session_conflict",
		Browser: &FlowBrowserConfig{
			UseSession: "missing-admin",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		ArtifactRoot:   t.TempDir(),
		BrowserCDPPort: 9222,
		Security:       &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected CDP override and use_session conflict")
	}
	if !strings.Contains(err.Error(), "use_session") || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected use_session conflict before session resolution, got %v", err)
	}
	if strings.Contains(err.Error(), "missing-admin") {
		t.Fatalf("should not try to resolve missing session before CDP conflict, got %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("CDP override missing use_session conflict should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func TestRunFlowRejectsRemoteCDPLaunchOptionBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "remote_cdp_launch_option",
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		BrowserCDPLaunch:   true,
		BrowserCDPEndpoint: "http://192.0.2.1:9222",
		Security:           &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected remote CDP launch option error")
	}
	if !strings.Contains(err.Error(), "only start or reuse a local browser") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("remote CDP launch option should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func TestRunFlowRejectsLocalCDPLaunchEndpointWithoutPortBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "local_cdp_launch_endpoint_without_port",
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		BrowserCDPLaunch:   true,
		BrowserCDPEndpoint: "http://127.0.0.1",
		Security:           &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected local CDP launch endpoint without port error")
	}
	if !strings.Contains(err.Error(), "explicit port") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("local CDP launch endpoint without port should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func TestRunFlowRejectsCDPWithBrowserVideoBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t,
		func() error {
			installCalled = true
			return fmt.Errorf("unexpected playwright install")
		},
		func() (*playwright.Playwright, error) {
			runCalled = true
			return nil, fmt.Errorf("unexpected playwright startup")
		},
	)
	defer restore()

	videoPath := filepath.Join(t.TempDir(), "nested", "browser.webm")
	_, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_with_browser_video",
		Browser: &FlowBrowserConfig{
			CDPPort: 9222,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}, FlowRunOptions{
		BrowserVideoOutputPath: videoPath,
		Security:               &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected CDP browser video conflict error")
	}
	if !strings.Contains(err.Error(), "browser video recording is not supported") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("CDP browser video conflict should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
	if _, statErr := os.Stat(filepath.Dir(videoPath)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("browser video directory should not be created before conflict validation, stat err=%v", statErr)
	}
}

func TestNormalizeCDPEndpointAcceptsCommonForms(t *testing.T) {
	tests := map[string]string{
		"127.0.0.1:9222":                               "http://127.0.0.1:9222",
		"127.0.0.1:9222/json/version":                  "http://127.0.0.1:9222",
		"127.0.0.1:9222/json/list?x=1#frag":            "http://127.0.0.1:9222",
		"localhost:9222":                               "http://localhost:9222",
		"localhost:9222/devtools/browser/a":            "http://localhost:9222",
		"[::1]:9222/json/version":                      "http://[::1]:9222",
		"http://127.0.0.1:9222?x=1#frag":               "http://127.0.0.1:9222",
		"http://127.0.0.1:9222/json/new":               "http://127.0.0.1:9222",
		"http://127.0.0.1:9222/json/protocol":          "http://127.0.0.1:9222",
		"http://127.0.0.1:9222/devtools/page/a":        "http://127.0.0.1:9222",
		"http://127.0.0.1:9222/devtools/browser/a":     "http://127.0.0.1:9222",
		"http://127.0.0.1:9222/json/version":           "http://127.0.0.1:9222",
		"https://localhost:9222/json/version?x=1#frag": "https://localhost:9222",
		"ws://127.0.0.1:9222/devtools/browser/a":       "ws://127.0.0.1:9222/devtools/browser/a",
		"ws://[::1]:9222/devtools/browser/a":           "ws://[::1]:9222/devtools/browser/a",
		" wss://localhost:9222/devtools/browser/a \t":  "wss://localhost:9222/devtools/browser/a",
	}
	for input, want := range tests {
		got, err := normalizeCDPEndpoint(input)
		if err != nil {
			t.Fatalf("normalize %q: %v", input, err)
		}
		if got != want {
			t.Fatalf("normalize %q = %q, want %q", input, got, want)
		}
	}
}

func TestCDPHTTPBaseNormalizesWebSocketEndpoints(t *testing.T) {
	tests := map[string]string{
		"ws://127.0.0.1:9222/devtools/browser/a":  "http://127.0.0.1:9222",
		"wss://localhost:9222/devtools/browser/a": "https://localhost:9222",
		"ws://[::1]:9222/devtools/browser/a":      "http://[::1]:9222",
	}
	for input, want := range tests {
		got, err := cdpHTTPBase(input)
		if err != nil {
			t.Fatalf("http base %q: %v", input, err)
		}
		if got != want {
			t.Fatalf("http base %q = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeCDPEndpointRejectsInvalidForms(t *testing.T) {
	tests := []string{
		"127.0.0.1:0",
		"127.0.0.1:70000/json/version",
		"http://[::1]:70000/json/version",
		"localhost:notaport/json/version",
		"not a url",
		"ftp://127.0.0.1:9222",
	}
	for _, input := range tests {
		if got, err := normalizeCDPEndpoint(input); err == nil {
			t.Fatalf("normalize %q unexpectedly succeeded with %q", input, got)
		}
	}
}

func TestCDPEndpointIsLocalRecognizesLoopbackHosts(t *testing.T) {
	localEndpoints := []string{
		"http://localhost:9222",
		"http://LOCALHOST:9222/json/version",
		"http://127.0.0.1:9222",
		"http://127.0.1.1:9222",
		"http://[::1]:9222/json/version",
		"ws://[::1]:9222/devtools/browser/a",
	}
	for _, endpoint := range localEndpoints {
		if !cdpEndpointIsLocal(endpoint) {
			t.Fatalf("expected local CDP endpoint: %s", endpoint)
		}
	}

	remoteEndpoints := []string{
		"http://192.0.2.1:9222",
		"http://example.com:9222",
		"ws://10.0.0.8:9222/devtools/browser/a",
	}
	for _, endpoint := range remoteEndpoints {
		if cdpEndpointIsLocal(endpoint) {
			t.Fatalf("expected remote CDP endpoint: %s", endpoint)
		}
	}
}

func TestEnsureLocalCDPBrowserRejectsNonLocalLaunchEndpointBeforeReachability(t *testing.T) {
	_, _, err := ensureLocalCDPBrowser(FlowBrowserConfig{
		CDPLaunch:   true,
		CDPEndpoint: "http://192.0.2.1:9222",
		Timeout:     1,
	}, FlowRunOptions{})
	if err == nil {
		t.Fatalf("expected non-local CDP launch endpoint to be rejected")
	}
	if !strings.Contains(err.Error(), "local browser") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveLocalCDPBrowserExecutableAcceptsExplicitPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "chrome-test")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	resolved, err := resolveLocalCDPBrowserExecutable(path)
	if err != nil {
		t.Fatalf("resolve executable: %v", err)
	}
	if resolved != path {
		t.Fatalf("resolved = %q, want %q", resolved, path)
	}
}

func TestResolveLocalCDPBrowserExecutableUsesEnvironmentCandidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "chrome-env-test")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	t.Setenv("TSPLAY_BROWSER_EXECUTABLE", path)
	t.Setenv("CHROME_EXECUTABLE", "")
	t.Setenv("CHROME_PATH", "")

	resolved, err := resolveLocalCDPBrowserExecutable("")
	if err != nil {
		t.Fatalf("resolve executable from env: %v", err)
	}
	if resolved != path {
		t.Fatalf("resolved = %q, want %q", resolved, path)
	}
}

func TestCDPLaunchDefaultUserDataDirIsUnique(t *testing.T) {
	root := t.TempDir()
	rootReal, err := prepareRuntimeFileRoot(root)
	if err != nil {
		t.Fatalf("prepare root: %v", err)
	}
	config := FlowBrowserConfig{CDPLaunch: true}
	first, err := config.cdpLaunchUserDataDir(FlowRunOptions{ArtifactRoot: root})
	if err != nil {
		t.Fatalf("first dir: %v", err)
	}
	second, err := config.cdpLaunchUserDataDir(FlowRunOptions{ArtifactRoot: root})
	if err != nil {
		t.Fatalf("second dir: %v", err)
	}
	if first == second {
		t.Fatalf("default CDP launch profile dir should be unique, both were %q", first)
	}
	for _, dir := range []string{first, second} {
		if err := ensurePathInsideRoot(dir, rootReal); err != nil {
			t.Fatalf("dir %q outside root: %v", dir, err)
		}
		if _, err := os.Stat(dir); err != nil {
			t.Fatalf("expected profile dir %q: %v", dir, err)
		}
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

func TestResolveRuntimeFilePathAllowsAliasInsideRoot(t *testing.T) {
	realRoot := t.TempDir()
	realRootCanonical, err := filepath.EvalSymlinks(realRoot)
	if err != nil {
		t.Fatalf("resolve real root: %v", err)
	}
	aliasParent := t.TempDir()
	aliasRoot := filepath.Join(aliasParent, "alias-root")
	if err := os.Symlink(realRoot, aliasRoot); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	inputPath := filepath.Join(aliasRoot, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"ok":true}`), 0600); err != nil {
		t.Fatalf("write input via alias: %v", err)
	}
	input, err := resolveRuntimeFilePath(inputPath, flowFileInputPath, FlowSecurityPolicy{
		FileInputRoot: aliasRoot,
	})
	if err != nil {
		t.Fatalf("resolve input alias path: %v", err)
	}
	wantInput := filepath.Join(realRootCanonical, "input.json")
	if input != wantInput {
		t.Fatalf("input path = %q, want %q", input, wantInput)
	}

	outputPath := filepath.Join(aliasRoot, "profiles", "cdp")
	output, err := resolveRuntimeFilePath(outputPath, flowFileOutputPath, FlowSecurityPolicy{
		FileOutputRoot: aliasRoot,
	})
	if err != nil {
		t.Fatalf("resolve output alias path: %v", err)
	}
	wantOutput := filepath.Join(realRootCanonical, "profiles", "cdp")
	if output != wantOutput {
		t.Fatalf("output path = %q, want %q", output, wantOutput)
	}
	if _, err := os.Stat(filepath.Join(realRootCanonical, "profiles")); err != nil {
		t.Fatalf("expected output parent to be created under real root: %v", err)
	}
	if err := validatePathWithinRoot(wantOutput, aliasRoot); err != nil {
		t.Fatalf("static validation should allow canonical path inside alias root: %v", err)
	}
	if err := validatePathWithinRoot(filepath.Join(aliasRoot, "profiles", "cdp"), realRootCanonical); err != nil {
		t.Fatalf("static validation should allow alias path inside canonical root: %v", err)
	}
}

func TestResolveRuntimeFilePathRejectsOutputTraversalBeforeCreatingParent(t *testing.T) {
	root := t.TempDir()
	outsideParent := filepath.Join(filepath.Dir(root), "tsplay-outside-parent")
	_ = os.RemoveAll(outsideParent)
	t.Cleanup(func() {
		_ = os.RemoveAll(outsideParent)
	})

	_, err := resolveRuntimeFilePath(filepath.Join("..", filepath.Base(outsideParent), "nested", "file.json"), flowFileOutputPath, FlowSecurityPolicy{
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected traversal path to be rejected")
	}
	if _, statErr := os.Stat(outsideParent); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("outside parent should not be created, stat err=%v", statErr)
	}
}

func TestResolveRuntimeFilePathRejectsOutputSymlinkAncestorBeforeCreatingParent(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	_, err := resolveRuntimeFilePath(filepath.Join("link", "nested", "file.json"), flowFileOutputPath, FlowSecurityPolicy{
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected symlink ancestor path to be rejected")
	}
	if _, statErr := os.Stat(filepath.Join(outside, "nested")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("outside symlink target should not receive created parent, stat err=%v", statErr)
	}
}

func TestResolveRuntimeFilePathRejectsOutputSymlinkTargetOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	target := filepath.Join(root, "profile")
	outsideTarget := filepath.Join(outside, "profile")
	if err := os.Symlink(outsideTarget, target); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	_, err := resolveRuntimeFilePath(target, flowFileOutputPath, FlowSecurityPolicy{
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected output symlink target outside root to be rejected")
	}
	if _, statErr := os.Stat(outsideTarget); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("outside symlink target should not be created, stat err=%v", statErr)
	}
}

func TestResolveRuntimeFilePathRejectsExistingOutputSymlinkTargetOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	outsideTarget := filepath.Join(outside, "profile")
	if err := os.MkdirAll(outsideTarget, 0755); err != nil {
		t.Fatalf("create outside target: %v", err)
	}
	target := filepath.Join(root, "profile")
	if err := os.Symlink(outsideTarget, target); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	_, err := resolveRuntimeFilePath(target, flowFileOutputPath, FlowSecurityPolicy{
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected existing output symlink target outside root to be rejected")
	}
	if !strings.Contains(err.Error(), "outside allowed file output root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFlowSecurityRejectsOutputSymlinkTargetOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	outsideTarget := filepath.Join(outside, "profile")
	if err := os.MkdirAll(outsideTarget, 0755); err != nil {
		t.Fatalf("create outside target: %v", err)
	}
	target := filepath.Join(root, "profile")
	if err := os.Symlink(outsideTarget, target); err != nil {
		t.Skipf("symlink not available: %v", err)
	}
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "cdp_profile_symlink_static",
		Browser: &FlowBrowserConfig{
			CDPLaunch:      true,
			CDPUserDataDir: "profile",
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: "https://example.com"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, FlowSecurityPolicy{
		AllowBrowserState: true,
		FileOutputRoot:    root,
	})
	if err == nil {
		t.Fatalf("expected static browser cdp_user_data_dir symlink target to be rejected")
	}
	if !strings.Contains(err.Error(), "cdp_user_data_dir") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCDPLaunchUserDataDirAllowsCanonicalPathInsideAliasRoot(t *testing.T) {
	realRoot := t.TempDir()
	realRootCanonical, err := filepath.EvalSymlinks(realRoot)
	if err != nil {
		t.Fatalf("resolve real root: %v", err)
	}
	aliasParent := t.TempDir()
	aliasRoot := filepath.Join(aliasParent, "artifact-root")
	if err := os.Symlink(realRoot, aliasRoot); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	dir := filepath.Join(realRootCanonical, "profiles", "cdp")
	resolved, err := (FlowBrowserConfig{
		CDPLaunch:      true,
		CDPUserDataDir: dir,
	}).cdpLaunchUserDataDir(FlowRunOptions{
		ArtifactRoot: aliasRoot,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
			FileOutputRoot:    aliasRoot,
		},
	})
	if err != nil {
		t.Fatalf("resolve cdp user data dir inside alias root: %v", err)
	}
	if resolved != dir {
		t.Fatalf("resolved = %q, want %q", resolved, dir)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected cdp user data dir to be created: %v", err)
	}
}

func TestEnsureLocalCDPBrowserDoesNotCreateProfileWhenExecutableMissing(t *testing.T) {
	root := t.TempDir()
	profileDir := filepath.Join(root, "profile")

	_, _, err := ensureLocalCDPBrowser(FlowBrowserConfig{
		CDPLaunch:      true,
		CDPExecutable:  filepath.Join(root, "missing-browser"),
		CDPUserDataDir: profileDir,
		Timeout:        1,
	}, FlowRunOptions{
		ArtifactRoot: root,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
			FileOutputRoot:    root,
		},
	})
	if err == nil {
		t.Fatalf("expected missing executable error")
	}
	if _, statErr := os.Stat(profileDir); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("profile dir should not be created when executable is missing, stat err=%v", statErr)
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

func TestRunFlowLuaExtractAndAssertHelpers(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_extract_and_assert",
		Steps: []FlowStep{
			{
				Action: "navigate",
				URL:    `data:text/html,<html><body><div id="ready">Order complete</div><div id="status">Status: shipped</div></body></html>`,
			},
			{
				Action: "lua",
				Code:   `return extract_text("#status", "Status: (.*)")`,
				SaveAs: "status_text",
			},
			{
				Action: "lua",
				Code:   `return assert_visible("#ready", 1000)`,
				SaveAs: "visible_ok",
			},
			{
				Action: "lua",
				Code:   `return assert_text("#ready", "complete", 1000)`,
				SaveAs: "text_assert",
			},
		},
	}

	result, err := RunFlow(flow, FlowRunOptions{Headless: true})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["status_text"]; got != "shipped" {
		t.Fatalf("status_text = %#v", got)
	}
	if got := result.Vars["visible_ok"]; got != true {
		t.Fatalf("visible_ok = %#v", got)
	}
	assertResult, ok := result.Vars["text_assert"].(map[string]any)
	if !ok {
		t.Fatalf("text_assert = %#v", result.Vars["text_assert"])
	}
	if got := assertResult["text"]; got != "complete" {
		t.Fatalf("assert text = %#v", got)
	}
}

func TestLuaAssertHelpersWorkWithoutFlowContext(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	for _, fn := range GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()
	runtime, err := OpenFlowBrowser(pw, FlowBrowserConfig{Headless: playwright.Bool(true)}, FlowRunOptions{}, nil)
	if err != nil {
		t.Fatalf("open browser: %v", err)
	}
	defer runtime.Browser.Close()
	setFlowBrowserGlobals(L, runtime.Browser, runtime.Context, runtime.Page)
	if _, err := runtime.Page.Goto(`data:text/html,<html><body><div id="ready">lua assertion ok</div></body></html>`); err != nil {
		t.Fatalf("goto page: %v", err)
	}

	if err := L.DoString(`assert_visible("#ready", 1000); assert_text("#ready", "lua assertion ok", 1000)`); err != nil {
		t.Fatalf("lua assert helpers without flow context: %v", err)
	}
}

func TestRunFlowLuaSetAndAppendVarHelpers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_set_and_append_vars",
		Vars: map[string]any{
			"message": "hello",
		},
		Steps: []FlowStep{
			{
				Action: "lua",
				Code:   `return set_var("copied_message", "{{message}}")`,
				SaveAs: "set_result",
			},
			{
				Action: "lua",
				Code:   `return append_var("items", copied_message)`,
				SaveAs: "append_one",
			},
			{
				Action: "lua",
				Code:   `return append_var("items", "tail")`,
				SaveAs: "append_two",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowLua: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["copied_message"]; got != "hello" {
		t.Fatalf("copied_message = %#v", got)
	}
	if got := result.Vars["set_result"]; got != "hello" {
		t.Fatalf("set_result = %#v", got)
	}

	appendOne, ok := result.Vars["append_one"].([]any)
	if !ok || len(appendOne) != 1 || appendOne[0] != "hello" {
		t.Fatalf("append_one = %#v", result.Vars["append_one"])
	}
	appendTwo, ok := result.Vars["append_two"].([]any)
	if !ok || len(appendTwo) != 2 || appendTwo[0] != "hello" || appendTwo[1] != "tail" {
		t.Fatalf("append_two = %#v", result.Vars["append_two"])
	}
	items, ok := result.Vars["items"].([]any)
	if !ok || len(items) != 2 || items[0] != "hello" || items[1] != "tail" {
		t.Fatalf("items = %#v", result.Vars["items"])
	}
}

func TestLuaSetAndAppendVarHelpersWithoutFlowContext(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	ensureFlowActionGlobals(L)

	if err := L.DoString(`
set_var("message", "hello")
append_var("items", message)
append_var("items", "tail")
`); err != nil {
		t.Fatalf("run lua: %v", err)
	}

	if got := luaValueToGo(L.GetGlobal("message")); got != "hello" {
		t.Fatalf("message = %#v", got)
	}
	items, ok := luaValueToGo(L.GetGlobal("items")).([]any)
	if !ok || len(items) != 2 || items[0] != "hello" || items[1] != "tail" {
		t.Fatalf("items = %#v", luaValueToGo(L.GetGlobal("items")))
	}
}

func TestRunFlowHTTPRequestAndJSONExtract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "山东" {
			t.Fatalf("unexpected query: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"answer":"山东","items":["山东","济南"]}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "http_json_extract",
		Steps: []FlowStep{
			{
				Action: "http_request",
				URL:    server.URL,
				SaveAs: "api_result",
				With: map[string]any{
					"query":       map[string]any{"q": "山东"},
					"response_as": "json",
				},
			},
			{
				Action: "json_extract",
				From:   "{{api_result}}",
				Path:   "$.body.answer",
				SaveAs: "answer",
			},
			{
				Action: "json_extract",
				From:   "{{api_result}}",
				Path:   "$.body.items[1]",
				SaveAs: "city",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowHTTP: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["answer"]; got != "山东" {
		t.Fatalf("answer = %#v", got)
	}
	if got := result.Vars["city"]; got != "济南" {
		t.Fatalf("city = %#v", got)
	}
	apiResult, ok := result.Vars["api_result"].(map[string]any)
	if !ok {
		t.Fatalf("api_result = %#v", result.Vars["api_result"])
	}
	if got := apiResult["status"]; got != 200 {
		t.Fatalf("status = %#v", got)
	}
}

func TestRunFlowHTTPRequestMultipartAndSavePath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "captcha.txt"), []byte("captcha-image"), 0600); err != nil {
		t.Fatalf("write captcha file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("scene"); got != "login" {
			t.Fatalf("unexpected scene: %q", got)
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			t.Fatalf("read multipart file: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read multipart content: %v", err)
		}
		if string(content) != "captcha-image" {
			t.Fatalf("unexpected multipart content: %q", string(content))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"text":"ABCD"}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "http_multipart",
		Steps: []FlowStep{
			{
				Action:   "http_request",
				URL:      server.URL,
				SaveAs:   "ocr_result",
				SavePath: "responses/ocr.json",
				With: map[string]any{
					"multipart_files":  map[string]any{"image": "captcha.txt"},
					"multipart_fields": map[string]any{"scene": "login"},
					"response_as":      "json",
				},
			},
			{
				Action: "json_extract",
				From:   "{{ocr_result}}",
				Path:   "$.body.text",
				SaveAs: "captcha_text",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["captcha_text"]; got != "ABCD" {
		t.Fatalf("captcha_text = %#v", got)
	}
	savedBody, err := os.ReadFile(filepath.Join(root, "responses", "ocr.json"))
	if err != nil {
		t.Fatalf("read saved response: %v", err)
	}
	if !strings.Contains(string(savedBody), `"text":"ABCD"`) {
		t.Fatalf("unexpected saved response: %s", string(savedBody))
	}
}

func TestRunFlowOCRReady(t *testing.T) {
	root := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ready" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ready","model":"old","detection":true,"slide_comparison":true,"slide_match":true}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "goddddocr_ready",
		Steps: []FlowStep{
			{
				Action:   "ocr_ready",
				URL:      server.URL + "/ocr/file",
				SaveAs:   "ocr_service",
				SavePath: "responses/ready.json",
			},
			{
				Action: "json_extract",
				From:   "{{ocr_service}}",
				Path:   "$.model",
				SaveAs: "ocr_model",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	readyResult, ok := result.Vars["ocr_service"].(map[string]any)
	if !ok {
		t.Fatalf("ocr_service = %#v", result.Vars["ocr_service"])
	}
	if got := readyResult["ready"]; got != true {
		t.Fatalf("ready = %#v", got)
	}
	if got := result.Vars["ocr_model"]; got != "old" {
		t.Fatalf("ocr_model = %#v", got)
	}
	if got := readyResult["detection"]; got != true {
		t.Fatalf("detection = %#v", got)
	}
	savedBody, err := os.ReadFile(filepath.Join(root, "responses", "ready.json"))
	if err != nil {
		t.Fatalf("read saved response: %v", err)
	}
	if !strings.Contains(string(savedBody), `"status":"ready"`) {
		t.Fatalf("unexpected saved response: %s", string(savedBody))
	}
}

func TestRunFlowOCRRequest(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "captcha.png"), []byte("captcha-image"), 0600); err != nil {
		t.Fatalf("write captcha file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ocr/file" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("charset_range"); got != "0123456789abcdef" {
			t.Fatalf("unexpected charset_range: %q", got)
		}
		if got := r.FormValue("confidence"); got != "true" {
			t.Fatalf("unexpected confidence: %q", got)
		}
		if got := r.FormValue("probability"); got != "true" {
			t.Fatalf("unexpected probability: %q", got)
		}
		if got := r.FormValue("color_filter_colors"); got != `["red","blue"]` {
			t.Fatalf("unexpected color_filter_colors: %q", got)
		}
		if got := r.FormValue("color_filter_custom_ranges"); got != `[[[90,30,30],[110,255,255]]]` {
			t.Fatalf("unexpected color_filter_custom_ranges: %q", got)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("read multipart file: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read multipart content: %v", err)
		}
		if string(content) != "captcha-image" {
			t.Fatalf("unexpected multipart content: %q", string(content))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"result":"3n3d","confidence":0.99,"probability":{"text":"3n3d","charsets":["","3","n","d"],"probability":[[0.01,0.97,0.01,0.01]],"confidence":0.97},"request_id":"req-1","processing_time_ms":12.5}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "goddddocr_ocr",
		Steps: []FlowStep{
			{
				Action:   "ocr_request",
				FilePath: "captcha.png",
				URL:      server.URL,
				SaveAs:   "ocr_result",
				SavePath: "responses/ocr.json",
				With: map[string]any{
					"charset_range":              "0123456789abcdef",
					"color_filter_colors":        []any{"red", "blue"},
					"color_filter_custom_ranges": []any{[]any{[]any{90, 30, 30}, []any{110, 255, 255}}},
					"probability":                true,
				},
			},
			{
				Action: "json_extract",
				From:   "{{ocr_result}}",
				Path:   "$.text",
				SaveAs: "captcha_text",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["captcha_text"]; got != "3n3d" {
		t.Fatalf("captcha_text = %#v", got)
	}
	ocrResult, ok := result.Vars["ocr_result"].(map[string]any)
	if !ok {
		t.Fatalf("ocr_result = %#v", result.Vars["ocr_result"])
	}
	if got := ocrResult["confidence"]; got != 0.99 {
		t.Fatalf("confidence = %#v", got)
	}
	probability, ok := ocrResult["probability"].(map[string]any)
	if !ok {
		t.Fatalf("probability = %#v", ocrResult["probability"])
	}
	if got := probability["text"]; got != "3n3d" {
		t.Fatalf("probability.text = %#v", got)
	}
	if got := ocrResult["request_id"]; got != "req-1" {
		t.Fatalf("request_id = %#v", got)
	}
	savedBody, err := os.ReadFile(filepath.Join(root, "responses", "ocr.json"))
	if err != nil {
		t.Fatalf("read saved response: %v", err)
	}
	if !strings.Contains(string(savedBody), `"result":"3n3d"`) {
		t.Fatalf("unexpected saved response: %s", string(savedBody))
	}
}

func TestRunFlowOCRDetect(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "captcha.png"), []byte("captcha-image"), 0600); err != nil {
		t.Fatalf("write captcha file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/det/file" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("detailed"); got != "true" {
			t.Fatalf("unexpected detailed: %q", got)
		}
		if got := r.FormValue("score_threshold"); got != "0.2" {
			t.Fatalf("unexpected score_threshold: %q", got)
		}
		if got := r.FormValue("nms_threshold"); got != "0.35" {
			t.Fatalf("unexpected nms_threshold: %q", got)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("read multipart file: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read multipart content: %v", err)
		}
		if string(content) != "captcha-image" {
			t.Fatalf("unexpected multipart content: %q", string(content))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"result":[[1,2,3,4]],"boxes":[{"x1":1,"y1":2,"x2":3,"y2":4,"score":0.92,"label":0}],"request_id":"det-1","processing_time_ms":4.2}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "goddddocr_detect",
		Steps: []FlowStep{
			{
				Action:   "ocr_detect",
				FilePath: "captcha.png",
				URL:      server.URL,
				SaveAs:   "det_result",
				SavePath: "responses/det.json",
				With: map[string]any{
					"score_threshold": 0.2,
					"nms_threshold":   0.35,
				},
			},
			{
				Action: "json_extract",
				From:   "{{det_result}}",
				Path:   "$.boxes[0].score",
				SaveAs: "first_score",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["first_score"]; got != 0.92 {
		t.Fatalf("first_score = %#v", got)
	}
	detResult, ok := result.Vars["det_result"].(map[string]any)
	if !ok {
		t.Fatalf("det_result = %#v", result.Vars["det_result"])
	}
	if got := detResult["request_id"]; got != "det-1" {
		t.Fatalf("request_id = %#v", got)
	}
	savedBody, err := os.ReadFile(filepath.Join(root, "responses", "det.json"))
	if err != nil {
		t.Fatalf("read saved response: %v", err)
	}
	if !strings.Contains(string(savedBody), `"request_id":"det-1"`) {
		t.Fatalf("unexpected saved response: %s", string(savedBody))
	}
}

func TestRunFlowOCRSlideComparison(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "target.png"), []byte("target-image"), 0600); err != nil {
		t.Fatalf("write target file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "background.png"), []byte("background-image"), 0600); err != nil {
		t.Fatalf("write background file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slide_comparison/file" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		for name, want := range map[string]string{"target_file": "target-image", "background_file": "background-image"} {
			file, _, err := r.FormFile(name)
			if err != nil {
				t.Fatalf("read multipart file %s: %v", name, err)
			}
			content, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				t.Fatalf("read multipart content %s: %v", name, err)
			}
			if string(content) != want {
				t.Fatalf("unexpected multipart content %s: %q", name, string(content))
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"result":{"target":[40,20],"target_x":40,"target_y":20},"request_id":"slide-1","processing_time_ms":5.5}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "goddddocr_slide_comparison",
		Steps: []FlowStep{
			{
				Action:             "ocr_slide_comparison",
				URL:                server.URL + "/slide_comparison",
				TargetFilePath:     "target.png",
				BackgroundFilePath: "background.png",
				SaveAs:             "slide_result",
			},
			{
				Action: "json_extract",
				From:   "{{slide_result}}",
				Path:   "$.target_x",
				SaveAs: "gap_x",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["gap_x"]; got != float64(40) {
		t.Fatalf("gap_x = %#v", got)
	}
}

func TestRunFlowOCRSlideMatch(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "target.png"), []byte("target-image"), 0600); err != nil {
		t.Fatalf("write target file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "background.png"), []byte("background-image"), 0600); err != nil {
		t.Fatalf("write background file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slide_match/file" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("simple_target"); got != "true" {
			t.Fatalf("unexpected simple_target: %q", got)
		}
		for name := range map[string]bool{"target_file": true, "background_file": true} {
			file, _, err := r.FormFile(name)
			if err != nil {
				t.Fatalf("read multipart file %s: %v", name, err)
			}
			file.Close()
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"result":{"target":[48,16],"target_x":48,"target_y":16,"confidence":0.88},"request_id":"match-1","processing_time_ms":6.1}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "goddddocr_slide_match",
		Steps: []FlowStep{
			{
				Action:             "ocr_slide_match",
				URL:                server.URL,
				TargetFilePath:     "target.png",
				BackgroundFilePath: "background.png",
				SaveAs:             "match_result",
				With: map[string]any{
					"simple_target": true,
				},
			},
			{
				Action: "json_extract",
				From:   "{{match_result}}",
				Path:   "$.confidence",
				SaveAs: "match_confidence",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["match_confidence"]; got != 0.88 {
		t.Fatalf("match_confidence = %#v", got)
	}
}

func TestGoddddocrOCRTutorialFlowValidates(t *testing.T) {
	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_ocr.flow.yaml"))
	if err != nil {
		t.Fatalf("load tutorial flow: %v", err)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate tutorial flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowHTTP: true, AllowFileAccess: true}); err != nil {
		t.Fatalf("validate tutorial flow security: %v", err)
	}
}

func TestGoddddocrLoginTutorialFlowValidates(t *testing.T) {
	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_login.flow.yaml"))
	if err != nil {
		t.Fatalf("load login tutorial flow: %v", err)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate login tutorial flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowHTTP: true, AllowFileAccess: true}); err != nil {
		t.Fatalf("validate login tutorial flow security: %v", err)
	}
}

func TestGoddddocrDetSlideTutorialFlowValidates(t *testing.T) {
	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide tutorial flow: %v", err)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate det slide tutorial flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowHTTP: true, AllowFileAccess: true}); err != nil {
		t.Fatalf("validate det slide tutorial flow security: %v", err)
	}
}

func TestGoddddocrDetSlideRecoveryTutorialFlowValidates(t *testing.T) {
	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide_recovery.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide recovery tutorial flow: %v", err)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate det slide recovery tutorial flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowHTTP: true, AllowFileAccess: true}); err != nil {
		t.Fatalf("validate det slide recovery tutorial flow security: %v", err)
	}
}

func TestGoddddocrDetSlideManualReviewTutorialFlowValidates(t *testing.T) {
	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide_manual_review.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide manual review tutorial flow: %v", err)
	}
	if err := ValidateFlowStrict(flow); err != nil {
		t.Fatalf("validate det slide manual review tutorial flow: %v", err)
	}
	if err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowHTTP: true, AllowFileAccess: true}); err != nil {
		t.Fatalf("validate det slide manual review tutorial flow security: %v", err)
	}
}

func TestRunGoddddocrDetSlideTutorialFlowWithDemo(t *testing.T) {
	root := t.TempDir()

	demoServer := httptest.NewServer(http.FileServer(http.Dir("..")))
	defer demoServer.Close()

	var detCalled bool
	var slideCalled bool
	ocrServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/ready":
			fmt.Fprint(w, `{"status":"ready","model":"old","detection":true,"slide_match":true,"slide_comparison":true}`)
		case "/det/file":
			detCalled = true
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse det multipart form: %v", err)
			}
			if got := r.FormValue("detailed"); got != "true" {
				t.Fatalf("unexpected detailed: %q", got)
			}
			if got := r.FormValue("score_threshold"); got != "0.2" {
				t.Fatalf("unexpected score_threshold: %q", got)
			}
			if got := r.FormValue("nms_threshold"); got != "0.45" {
				t.Fatalf("unexpected nms_threshold: %q", got)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read det multipart file: %v", err)
			}
			content, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				t.Fatalf("read det multipart content: %v", err)
			}
			if len(content) == 0 {
				t.Fatalf("expected det screenshot content")
			}
			fmt.Fprint(w, `{"result":[[28,28,76,76]],"boxes":[{"x1":28,"y1":28,"x2":76,"y2":76,"score":0.91,"label":0}],"request_id":"det-demo","processing_time_ms":3.2}`)
		case "/slide_match/file":
			slideCalled = true
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse slide multipart form: %v", err)
			}
			if got := r.FormValue("simple_target"); got != "true" {
				t.Fatalf("unexpected simple_target: %q", got)
			}
			for _, field := range []string{"target_file", "background_file"} {
				file, _, err := r.FormFile(field)
				if err != nil {
					t.Fatalf("read slide multipart file %s: %v", field, err)
				}
				content, err := io.ReadAll(file)
				file.Close()
				if err != nil {
					t.Fatalf("read slide multipart content %s: %v", field, err)
				}
				if len(content) == 0 {
					t.Fatalf("expected slide screenshot content for %s", field)
				}
			}
			fmt.Fprint(w, `{"result":{"target":[126,82],"target_x":126,"target_y":82,"confidence":0.96},"request_id":"slide-demo","processing_time_ms":4.8}`)
		default:
			t.Fatalf("unexpected OCR path: %s", r.URL.Path)
		}
	}))
	defer ocrServer.Close()

	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide tutorial flow: %v", err)
	}
	flow.Vars["page_url"] = demoServer.URL + "/demo/slider_login.html"
	flow.Vars["goddddocr_url"] = ocrServer.URL

	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: root,
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run det slide tutorial flow: %v", err)
	}
	if !detCalled || !slideCalled {
		t.Fatalf("expected det and slide endpoints to be called, det=%v slide=%v", detCalled, slideCalled)
	}
	payload, ok := result.Vars["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload = %#v", result.Vars["payload"])
	}
	if got := payload["slide_target_x"]; got != 126.0 {
		t.Fatalf("slide_target_x = %#v", got)
	}
	if got := payload["service_detection"]; got != true {
		t.Fatalf("service_detection = %#v", got)
	}
	clickResult, ok := payload["detect_click_result"].(map[string]any)
	if !ok {
		t.Fatalf("detect_click_result = %#v", payload["detect_click_result"])
	}
	if got, ok := clickResult["x"].(float64); !ok || got < 51 || got > 53 {
		t.Fatalf("detect click x = %#v", clickResult["x"])
	}
	if got, ok := clickResult["y"].(float64); !ok || got < 51 || got > 53 {
		t.Fatalf("detect click y = %#v", clickResult["y"])
	}
	if _, err := os.Stat(filepath.Join(root, "artifacts", "goddddocr", "det-slide-flow-result.json")); err != nil {
		t.Fatalf("expected det slide result artifact: %v", err)
	}
}

func TestRunGoddddocrDetSlideRecoveryTutorialFlowRetriesAndWritesDiagnostics(t *testing.T) {
	root := t.TempDir()

	demoServer := httptest.NewServer(http.FileServer(http.Dir("..")))
	defer demoServer.Close()

	var detCalls int
	var slideCalls int
	ocrServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/ready":
			fmt.Fprint(w, `{"status":"ready","model":"old","detection":true,"slide_match":true,"slide_comparison":true}`)
		case "/det/file":
			detCalls++
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse det multipart form: %v", err)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read det multipart file: %v", err)
			}
			content, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				t.Fatalf("read det multipart content: %v", err)
			}
			if len(content) == 0 {
				t.Fatalf("expected det screenshot content")
			}
			switch detCalls {
			case 1:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"error":"temporary detection failure"}`)
			case 2:
				fmt.Fprint(w, `{"result":[[180,28,220,76]],"boxes":[{"x1":180,"y1":28,"x2":220,"y2":76,"score":0.61,"label":0}],"request_id":"det-offset","processing_time_ms":3.4}`)
			default:
				fmt.Fprint(w, `{"result":[[28,28,76,76]],"boxes":[{"x1":28,"y1":28,"x2":76,"y2":76,"score":0.93,"label":0}],"request_id":"det-ok","processing_time_ms":3.1}`)
			}
		case "/slide_match/file":
			slideCalls++
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse slide multipart form: %v", err)
			}
			for _, field := range []string{"target_file", "background_file"} {
				file, _, err := r.FormFile(field)
				if err != nil {
					t.Fatalf("read slide multipart file %s: %v", field, err)
				}
				content, err := io.ReadAll(file)
				file.Close()
				if err != nil {
					t.Fatalf("read slide multipart content %s: %v", field, err)
				}
				if len(content) == 0 {
					t.Fatalf("expected slide screenshot content for %s", field)
				}
			}
			if slideCalls == 1 {
				fmt.Fprint(w, `{"result":{"target":[40,82],"target_x":40,"target_y":82,"confidence":0.44},"request_id":"slide-short","processing_time_ms":5.2}`)
				return
			}
			fmt.Fprint(w, `{"result":{"target":[126,82],"target_x":126,"target_y":82,"confidence":0.96},"request_id":"slide-ok","processing_time_ms":4.8}`)
		default:
			t.Fatalf("unexpected OCR path: %s", r.URL.Path)
		}
	}))
	defer ocrServer.Close()

	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide_recovery.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide recovery tutorial flow: %v", err)
	}
	flow.Vars["page_url"] = demoServer.URL + "/demo/slider_login.html"
	flow.Vars["goddddocr_url"] = ocrServer.URL

	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: root,
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run det slide recovery tutorial flow: %v", err)
	}
	if detCalls != 4 {
		t.Fatalf("det calls = %d", detCalls)
	}
	if slideCalls != 2 {
		t.Fatalf("slide calls = %d", slideCalls)
	}

	payload, ok := result.Vars["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload = %#v", result.Vars["payload"])
	}
	if got := payload["status"]; got != "success" {
		t.Fatalf("payload status = %#v", got)
	}
	retryResult, ok := payload["retry_result"].(map[string]any)
	if !ok {
		t.Fatalf("retry_result = %#v", payload["retry_result"])
	}
	if got := retryResult["attempts"]; got != 4 {
		t.Fatalf("retry attempts = %#v", got)
	}
	if got := retryResult["status"]; got != "succeeded" {
		t.Fatalf("retry status = %#v", got)
	}

	diagnosticPath := filepath.Join(root, "artifacts", "goddddocr", "det-slide-recovery-diagnostic.json")
	diagnosticBytes, err := os.ReadFile(diagnosticPath)
	if err != nil {
		t.Fatalf("read recovery diagnostic: %v", err)
	}
	var diagnostic map[string]any
	if err := json.Unmarshal(diagnosticBytes, &diagnostic); err != nil {
		t.Fatalf("parse recovery diagnostic: %v", err)
	}
	if got := diagnostic["status"]; got != "retrying" {
		t.Fatalf("diagnostic status = %#v", got)
	}
	if got := diagnostic["phase"]; got != "slide" {
		t.Fatalf("diagnostic phase = %#v", got)
	}
	if got := fmt.Sprint(diagnostic["error"]); !strings.Contains(got, "assert_visible") {
		t.Fatalf("diagnostic error = %#v", got)
	}
	for _, path := range []string{
		filepath.Join(root, "artifacts", "goddddocr", "det-slide-recovery-result.json"),
		filepath.Join(root, "artifacts", "goddddocr", "det-slide-recovery-failure.png"),
		filepath.Join(root, "artifacts", "goddddocr", "det-slide-recovery-failure.html"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected recovery artifact %s: %v", path, err)
		}
	}
}

func TestRunGoddddocrDetSlideManualReviewStopsOnLowDetectionScore(t *testing.T) {
	result, root, detCalls, slideCalls := runGoddddocrDetSlideManualReviewTutorial(t, 0.42, 0.96)
	if detCalls != 1 {
		t.Fatalf("det calls = %d", detCalls)
	}
	if slideCalls != 0 {
		t.Fatalf("slide calls = %d", slideCalls)
	}

	payload := requireFlowPayload(t, result)
	if got := payload["status"]; got != "manual_review" {
		t.Fatalf("payload status = %#v", got)
	}
	if got := payload["action"]; got != "manual_review_required" {
		t.Fatalf("payload action = %#v", got)
	}
	if got := payload["phase"]; got != "detect_score" {
		t.Fatalf("payload phase = %#v", got)
	}
	if got := fmt.Sprint(payload["reason"]); !strings.Contains(got, "detection score") || !strings.Contains(got, "0.42 >= 0.8") {
		t.Fatalf("payload reason = %#v", got)
	}
	assertManualReviewArtifacts(t, root)
}

func TestRunGoddddocrDetSlideManualReviewStopsOnLowSlideConfidence(t *testing.T) {
	result, root, detCalls, slideCalls := runGoddddocrDetSlideManualReviewTutorial(t, 0.93, 0.44)
	if detCalls != 1 {
		t.Fatalf("det calls = %d", detCalls)
	}
	if slideCalls != 1 {
		t.Fatalf("slide calls = %d", slideCalls)
	}

	payload := requireFlowPayload(t, result)
	if got := payload["status"]; got != "manual_review" {
		t.Fatalf("payload status = %#v", got)
	}
	if got := payload["phase"]; got != "slide_confidence" {
		t.Fatalf("payload phase = %#v", got)
	}
	if got := fmt.Sprint(payload["reason"]); !strings.Contains(got, "slide confidence") || !strings.Contains(got, "0.44 >= 0.8") {
		t.Fatalf("payload reason = %#v", got)
	}
	if got := payload["drag_result"]; fmt.Sprint(got) != "map[]" {
		t.Fatalf("drag_result should remain empty when slide confidence is low, got %#v", got)
	}
	assertManualReviewArtifacts(t, root)
}

func runGoddddocrDetSlideManualReviewTutorial(t *testing.T, detScore float64, slideConfidence float64) (*FlowResult, string, int, int) {
	t.Helper()
	root := t.TempDir()

	demoServer := httptest.NewServer(http.FileServer(http.Dir("..")))
	t.Cleanup(demoServer.Close)

	var detCalls int
	var slideCalls int
	ocrServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/ready":
			fmt.Fprint(w, `{"status":"ready","model":"old","detection":true,"slide_match":true,"slide_comparison":true}`)
		case "/det/file":
			detCalls++
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse det multipart form: %v", err)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read det multipart file: %v", err)
			}
			content, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				t.Fatalf("read det multipart content: %v", err)
			}
			if len(content) == 0 {
				t.Fatalf("expected det screenshot content")
			}
			fmt.Fprintf(w, `{"result":[[28,28,76,76]],"boxes":[{"x1":28,"y1":28,"x2":76,"y2":76,"score":%.2f,"label":0}],"request_id":"det-gate","processing_time_ms":3.2}`, detScore)
		case "/slide_match/file":
			slideCalls++
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				t.Fatalf("parse slide multipart form: %v", err)
			}
			for _, field := range []string{"target_file", "background_file"} {
				file, _, err := r.FormFile(field)
				if err != nil {
					t.Fatalf("read slide multipart file %s: %v", field, err)
				}
				content, err := io.ReadAll(file)
				file.Close()
				if err != nil {
					t.Fatalf("read slide multipart content %s: %v", field, err)
				}
				if len(content) == 0 {
					t.Fatalf("expected slide screenshot content for %s", field)
				}
			}
			fmt.Fprintf(w, `{"result":{"target":[126,82],"target_x":126,"target_y":82,"confidence":%.2f},"request_id":"slide-gate","processing_time_ms":4.8}`, slideConfidence)
		default:
			t.Fatalf("unexpected OCR path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(ocrServer.Close)

	flow, err := LoadFlowFile(filepath.Join("..", "script", "tutorials", "goddddocr_det_slide_manual_review.flow.yaml"))
	if err != nil {
		t.Fatalf("load det slide manual review tutorial flow: %v", err)
	}
	flow.Vars["page_url"] = demoServer.URL + "/demo/slider_login.html"
	flow.Vars["goddddocr_url"] = ocrServer.URL

	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: root,
		Security: &FlowSecurityPolicy{
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run det slide manual review tutorial flow: %v", err)
	}
	return result, root, detCalls, slideCalls
}

func requireFlowPayload(t *testing.T, result *FlowResult) map[string]any {
	t.Helper()
	payload, ok := result.Vars["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload = %#v", result.Vars["payload"])
	}
	return payload
}

func assertManualReviewArtifacts(t *testing.T, root string) {
	t.Helper()
	for _, path := range []string{
		filepath.Join(root, "artifacts", "goddddocr", "manual-review-result.json"),
		filepath.Join(root, "artifacts", "goddddocr", "manual-review-evidence.png"),
		filepath.Join(root, "artifacts", "goddddocr", "manual-review-evidence.html"),
		filepath.Join(root, "artifacts", "goddddocr", "manual-review-det-response.json"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected manual review artifact %s: %v", path, err)
		}
	}
}

func TestRunFlowLuaHTTPRequestHonorsAllowHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("request should not have been sent: %s", r.URL.String())
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_http_policy",
		Steps: []FlowStep{
			{
				Action: "lua",
				Code: fmt.Sprintf(`return http_request({
  url = %q,
  response_as = "json"
})`, server.URL),
			},
		},
	}

	_, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowLua: true},
	})
	if err == nil {
		t.Fatalf("expected allow_http runtime error")
	}
	if !strings.Contains(err.Error(), "allow_http") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowLuaHTTPRequestRequiresFileAccessForMultipartAndSavePath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "captcha.txt"), []byte("captcha-image"), 0600); err != nil {
		t.Fatalf("write captcha file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_http_file_policy",
		Steps: []FlowStep{
			{
				Action: "lua",
				Code: `return http_request({
  url = "https://example.com/ocr",
  save_path = "responses/ocr.json",
  multipart_files = {image = "captcha.txt"}
})`,
			},
		},
	}

	_, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowLua:  true,
			AllowHTTP: true,
		},
	})
	if err == nil {
		t.Fatalf("expected allow_file_access runtime error")
	}
	if !strings.Contains(err.Error(), "allow_file_access") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "responses", "ocr.json")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no saved response, stat err=%v", statErr)
	}
}

func TestRunFlowLuaHTTPRequestUsesFlowFileRoots(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "captcha.txt"), []byte("captcha-image"), 0600); err != nil {
		t.Fatalf("write captcha file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("scene"); got != "login" {
			t.Fatalf("unexpected scene: %q", got)
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			t.Fatalf("read multipart file: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read multipart content: %v", err)
		}
		if string(content) != "captcha-image" {
			t.Fatalf("unexpected multipart content: %q", string(content))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"text":"ABCD"}`)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_http_file_roots",
		Steps: []FlowStep{
			{
				Action: "lua",
				SaveAs: "ocr_result",
				Code: fmt.Sprintf(`return http_request({
  url = %q,
  save_path = "responses/ocr.json",
  multipart_files = {image = "captcha.txt"},
  multipart_fields = {scene = "login"},
  response_as = "json"
})`, server.URL),
			},
			{
				Action: "lua",
				SaveAs: "captcha_text",
				Code:   `return json_extract(ocr_result, "$.body.text")`,
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowLua:        true,
			AllowHTTP:       true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["captcha_text"]; got != "ABCD" {
		t.Fatalf("captcha_text = %#v", got)
	}
	savedBody, err := os.ReadFile(filepath.Join(root, "responses", "ocr.json"))
	if err != nil {
		t.Fatalf("read saved response: %v", err)
	}
	if !strings.Contains(string(savedBody), `"text":"ABCD"`) {
		t.Fatalf("unexpected saved response: %s", string(savedBody))
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

func TestRunFlowConnectsOverCDP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">cdp attached</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_attach_flow",
		Browser: &FlowBrowserConfig{
			CDPPort: port,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "cdp attached", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless: true,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err != nil {
		t.Fatalf("run flow over CDP: %v", err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected after CDP flow")
	}
}

func TestRunFlowCDPOptionReportsPlaywrightUsageForNonBrowserFlow(t *testing.T) {
	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_option_non_browser_result",
		Steps: []FlowStep{
			{Action: "set_var", SaveAs: "answer", Value: "ok"},
		},
	}, FlowRunOptions{
		Headless:       true,
		BrowserCDPPort: port,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err != nil {
		t.Fatalf("run non-browser flow over CDP option: %v", err)
	}
	if result == nil || len(result.Trace) != 1 || result.Vars["answer"] != "ok" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.Playwright == nil || !result.Playwright.NeedsBrowserState {
		t.Fatalf("expected CDP Playwright usage in result: %#v", result.Playwright)
	}
	if summary := result.Playwright.Summary(10); !strings.Contains(summary, "browser.cdp_port") {
		t.Fatalf("expected browser.cdp_port in result Playwright summary, got %q", summary)
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected after CDP option flow")
	}
}

func TestRunFlowConnectsOverCDPWebSocketEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">cdp websocket attached</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	endpoint := cdpWebSocketEndpoint(t, port)
	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_ws_attach_flow",
		Browser: &FlowBrowserConfig{
			CDPEndpoint: endpoint,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "cdp websocket attached", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless: true,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err != nil {
		t.Fatalf("run flow over CDP websocket endpoint %q: %v", endpoint, err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected after websocket CDP flow")
	}
}

func TestRunFlowConnectsOverCDPBareJSONVersionEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">cdp bare endpoint attached</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	endpoint := fmt.Sprintf("127.0.0.1:%d/json/version", port)
	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_bare_json_version_attach_flow",
		Browser: &FlowBrowserConfig{
			CDPEndpoint: endpoint,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "cdp bare endpoint attached", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless: true,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err != nil {
		t.Fatalf("run flow over bare CDP endpoint %q: %v", endpoint, err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected after bare endpoint CDP flow")
	}
}

func TestRunFlowCDPLaunchUsesExistingLocalEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">reused existing cdp</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_launch_reuses_existing_flow",
		Browser: &FlowBrowserConfig{
			CDPLaunch:     true,
			CDPPort:       port,
			CDPExecutable: filepath.Join(t.TempDir(), "missing-browser"),
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "reused existing cdp", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless: true,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err != nil {
		t.Fatalf("run flow over existing CDP endpoint with cdp_launch: %v", err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected when cdp_launch reuses an existing endpoint")
	}
}

func TestRunFlowCDPFailureKeepsExternalBrowserConnected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">actual text</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	defer pw.Stop()

	port := freeTCPPort(t)
	remoteBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	if err != nil {
		t.Fatalf("launch remote-debugging browser: %v", err)
	}
	defer remoteBrowser.Close()

	_, err = RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_attach_failed_flow",
		Browser: &FlowBrowserConfig{
			CDPPort: port,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "expected mismatch", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless: true,
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
		},
	})
	if err == nil {
		t.Fatalf("expected flow assertion failure")
	}
	if !remoteBrowser.IsConnected() {
		t.Fatalf("external browser should remain connected after failed CDP flow")
	}
}

func TestRunFlowLaunchesLocalBrowserOverCDP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><body><div id="ready">cdp launched</div></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	executable := pw.Chromium.ExecutablePath()
	if err := pw.Stop(); err != nil {
		t.Fatalf("stop playwright: %v", err)
	}
	if executable == "" {
		t.Skip("playwright chromium executable path is empty")
	}

	profileRoot, err := prepareRuntimeFileRoot(t.TempDir())
	if err != nil {
		t.Fatalf("prepare profile root: %v", err)
	}
	profileDir := filepath.Join(profileRoot, "cdp-profile")
	result, err := RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_launch_flow",
		Browser: &FlowBrowserConfig{
			CDPLaunch:      true,
			CDPExecutable:  executable,
			CDPUserDataDir: profileDir,
			Timeout:        15000,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL},
			{Action: "assert_text", Selector: "#ready", Text: "cdp launched", Timeout: 3000},
		},
	}, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: t.TempDir(),
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
			FileOutputRoot:    profileRoot,
		},
	})
	if err != nil {
		t.Fatalf("run flow over launched CDP browser: %v", err)
	}
	if result == nil || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if _, err := os.Stat(profileDir); err != nil {
		t.Fatalf("expected CDP launch profile dir: %v", err)
	}
}

func TestRunFlowCleansUpLaunchedCDPBrowserWhenConnectFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("process liveness check uses Unix signal semantics")
	}

	helperDir := t.TempDir()
	helperSource := filepath.Join(helperDir, "fake_cdp_browser.go")
	helperBinary := filepath.Join(helperDir, "fake-cdp-browser")
	helperCode := `package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := ""
	userDataDir := ""
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--remote-debugging-port=") {
			port = strings.TrimPrefix(arg, "--remote-debugging-port=")
		}
		if strings.HasPrefix(arg, "--user-data-dir=") {
			userDataDir = strings.TrimPrefix(arg, "--user-data-dir=")
		}
	}
	if port == "" {
		os.Exit(2)
	}
	if userDataDir != "" {
		_ = os.MkdirAll(userDataDir, 0755)
		_ = os.WriteFile(filepath.Join(userDataDir, "fake-browser.pid"), []byte(fmt.Sprint(os.Getpid())), 0644)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/json/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, ` + "`" + `{"Browser":"FakeChrome/1.0"}` + "`" + `)
	})
	if err := http.ListenAndServe("127.0.0.1:"+port, mux); err != nil {
		os.Exit(3)
	}
}
`
	if err := os.WriteFile(helperSource, []byte(helperCode), 0644); err != nil {
		t.Fatalf("write helper source: %v", err)
	}
	if output, err := exec.Command("go", "build", "-o", helperBinary, helperSource).CombinedOutput(); err != nil {
		t.Fatalf("build fake CDP browser: %v\n%s", err, output)
	}

	profileRoot, err := prepareRuntimeFileRoot(t.TempDir())
	if err != nil {
		t.Fatalf("prepare profile root: %v", err)
	}
	profileDir := filepath.Join(profileRoot, "fake-profile")
	port := freeTCPPort(t)
	_, err = RunFlow(&Flow{
		SchemaVersion: "1",
		Name:          "cdp_launch_connect_failure_cleanup",
		Browser: &FlowBrowserConfig{
			CDPLaunch:      true,
			CDPExecutable:  helperBinary,
			CDPUserDataDir: profileDir,
			CDPPort:        port,
			Timeout:        2000,
		},
		Steps: []FlowStep{
			{Action: "navigate", URL: `data:text/html,<html></html>`},
		},
	}, FlowRunOptions{
		Headless:     true,
		ArtifactRoot: t.TempDir(),
		Security: &FlowSecurityPolicy{
			AllowBrowserState: true,
			FileOutputRoot:    profileRoot,
		},
	})
	if err == nil {
		t.Fatalf("expected ConnectOverCDP failure for fake browser")
	}

	pidBytes, readErr := os.ReadFile(filepath.Join(profileDir, "fake-browser.pid"))
	if readErr != nil {
		t.Fatalf("read fake browser pid: %v", readErr)
	}
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if parseErr != nil {
		t.Fatalf("parse pid %q: %v", string(pidBytes), parseErr)
	}
	defer func() {
		if processStillRunning(pid) {
			if process, findErr := os.FindProcess(pid); findErr == nil {
				_ = process.Kill()
			}
		}
	}()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if !processStillRunning(pid) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("fake CDP browser process %d was not cleaned up after ConnectOverCDP failure: %v", pid, err)
}

func processStillRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}

func cdpWebSocketEndpoint(t *testing.T, port int) string {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	client := http.Client{Timeout: 500 * time.Millisecond}
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/json/version", port))
		if err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("unexpected /json/version status %d", resp.StatusCode)
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			time.Sleep(100 * time.Millisecond)
			continue
		}
		var payload struct {
			WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
		}
		err = json.NewDecoder(resp.Body).Decode(&payload)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if strings.TrimSpace(payload.WebSocketDebuggerURL) != "" {
			return payload.WebSocketDebuggerURL
		}
		lastErr = fmt.Errorf("missing webSocketDebuggerUrl in /json/version response")
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("fetch CDP websocket endpoint on port %d: %v", port, lastErr)
	return ""
}

func TestRunFlowLuaSaveStorageStateCreatesParentDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/seed":
			fmt.Fprint(w, `<html><body><script>localStorage.setItem("token","lua-parent-dir"); window.location.href="/check";</script></body></html>`)
		case "/check":
			fmt.Fprint(w, `<html><body><div id="ready"></div><script>document.getElementById("ready").textContent = localStorage.getItem("token") || "missing";</script></body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	root := t.TempDir()
	statePath := filepath.Join(root, "artifacts", "states", "admin.json")
	if _, err := os.Stat(filepath.Dir(statePath)); !os.IsNotExist(err) {
		t.Fatalf("expected storage state directory to start missing, stat err=%v", err)
	}

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_save_storage_state_creates_parent",
		Steps: []FlowStep{
			{Action: "navigate", URL: server.URL + "/seed"},
			{Action: "assert_text", Selector: "#ready", Text: "lua-parent-dir", Timeout: 3000},
			{
				Action: "lua",
				Code:   fmt.Sprintf("return save_storage_state(%q)", filepath.ToSlash(statePath)),
				SaveAs: "saved_state_path",
			},
		},
	}

	result, err := RunFlow(flow, FlowRunOptions{Headless: true})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["saved_state_path"]; got != filepath.ToSlash(statePath) {
		t.Fatalf("saved_state_path = %#v", got)
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("expected storage state file: %v", err)
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
	session, err := LoadFlowSavedSession("admin", root)
	if err != nil {
		t.Fatalf("load named session: %v", err)
	}
	if session.LastUsedAt == "" {
		t.Fatalf("expected last_used_at to be set")
	}
	if session.SourceType != "storage_state_path" {
		t.Fatalf("expected source_type storage_state_path, got %q", session.SourceType)
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

func freeTCPPort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for free port: %v", err)
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected TCP addr: %#v", listener.Addr())
	}
	return addr.Port
}
