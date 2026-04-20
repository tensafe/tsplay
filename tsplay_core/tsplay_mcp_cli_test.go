package tsplay_core

import (
	"context"
	"testing"
)

func TestInvokeTSPlayToolListActions(t *testing.T) {
	payload, err := InvokeTSPlayTool(context.Background(), "tsplay.list_actions", nil)
	if err != nil {
		t.Fatalf("invoke list_actions: %v", err)
	}

	if payload["tool"] != "tsplay.list_actions" {
		t.Fatalf("unexpected tool payload: %#v", payload["tool"])
	}
	actions, ok := payload["actions"].([]any)
	if !ok || len(actions) == 0 {
		t.Fatalf("expected non-empty actions, got %#v", payload["actions"])
	}
}

func TestInvokeTSPlayToolValidateFlow(t *testing.T) {
	payload, err := InvokeTSPlayTool(context.Background(), "tsplay.validate_flow", map[string]any{
		"flow": `schema_version: "1"
name: validate_from_cli_helper
steps:
  - action: set_var
    save_as: payload
    with:
      value:
        lesson: "mcp"
`,
	})
	if err != nil {
		t.Fatalf("invoke validate_flow: %v", err)
	}

	valid, ok := payload["valid"].(bool)
	if !ok || !valid {
		t.Fatalf("expected valid=true, got %#v", payload)
	}
	if payload["tool"] != "tsplay.validate_flow" {
		t.Fatalf("unexpected tool payload: %#v", payload["tool"])
	}
}

func TestInvokeTSPlayToolUnknownTool(t *testing.T) {
	if _, err := InvokeTSPlayTool(context.Background(), "tsplay.not_real", nil); err == nil {
		t.Fatal("expected unsupported tool error")
	}
}
