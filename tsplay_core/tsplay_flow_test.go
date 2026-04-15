package tsplay_core

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLoadFlowYAMLAndValidate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "search.flow.yaml")
	content := []byte(`
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
		Name: "lua_only",
		Vars: map[string]any{"prefix": "hello"},
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
}
