package tsplay_core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestNewTSPlayMCPServerOnlyRegistersTSPlayTools(t *testing.T) {
	names := toolNamesForTest(NewTSPlayMCPServer())
	want := []string{"tsplay.list_actions", "tsplay.run_flow", "tsplay.validate_flow"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("tool names = %#v, want %#v", names, want)
	}
}

func TestHandleFlowListActionsTool(t *testing.T) {
	result, err := handleFlowListActionsTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("list actions: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	actions, ok := payload["actions"].([]any)
	if !ok || len(actions) == 0 {
		t.Fatalf("expected actions, got %#v", payload["actions"])
	}

	foundNavigate := false
	for _, action := range actions {
		item, ok := action.(map[string]any)
		if ok && item["name"] == "navigate" {
			foundNavigate = true
			break
		}
	}
	if !foundNavigate {
		t.Fatalf("navigate action not found in manifest")
	}
}

func TestHandleValidateFlowTool(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_from_mcp
steps:
  - action: lua
    code: return "ok"
`,
				"allow_lua": true,
			},
		},
	}

	result, err := handleValidateFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("validate flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["valid"] != true {
		t.Fatalf("expected valid flow, got %#v", payload)
	}
	if payload["name"] != "validate_from_mcp" {
		t.Fatalf("unexpected flow name: %#v", payload["name"])
	}
}

func TestHandleValidateFlowToolRejectsLuaWithoutAllow(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_from_mcp
steps:
  - action: lua
    code: return "ok"
`,
			},
		},
	}

	result, err := handleValidateFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("validate flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["valid"] != false {
		t.Fatalf("expected invalid flow, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "allow_lua") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func TestHandleValidateFlowToolRestrictsFlowPathRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	insidePath := filepath.Join(root, "ok.flow.yaml")
	outsidePath := filepath.Join(outside, "outside.flow.yaml")
	content := []byte(`
schema_version: "1"
name: path_flow
steps:
  - action: navigate
    url: https://example.com
`)
	if err := os.WriteFile(insidePath, content, 0600); err != nil {
		t.Fatalf("write inside flow: %v", err)
	}
	if err := os.WriteFile(outsidePath, content, 0600); err != nil {
		t.Fatalf("write outside flow: %v", err)
	}

	options := TSPlayMCPServerOptions{FlowPathRoot: root}
	insideRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow_path": "ok.flow.yaml",
			},
		},
	}
	insideResult, err := handleValidateFlowToolWithOptions(context.Background(), insideRequest, options)
	if err != nil {
		t.Fatalf("validate inside flow path: %v", err)
	}
	var insidePayload map[string]any
	decodeToolText(t, insideResult, &insidePayload)
	if insidePayload["valid"] != true {
		t.Fatalf("expected inside flow path to validate, got %#v", insidePayload)
	}

	outsideRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow_path": outsidePath,
			},
		},
	}
	outsideResult, err := handleValidateFlowToolWithOptions(context.Background(), outsideRequest, options)
	if err != nil {
		t.Fatalf("validate outside flow path: %v", err)
	}
	var outsidePayload map[string]any
	decodeToolText(t, outsideResult, &outsidePayload)
	if outsidePayload["valid"] != false {
		t.Fatalf("expected outside flow path to fail, got %#v", outsidePayload)
	}
	if !strings.Contains(outsidePayload["error"].(string), "outside allowed flow root") {
		t.Fatalf("unexpected error: %#v", outsidePayload["error"])
	}
}

func TestHandleValidateFlowToolRestrictsArtifactRoot(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: artifact_root
steps:
  - action: screenshot
    path: ../escape.png
`,
				"allow_file_access": true,
			},
		},
	}

	result, err := handleValidateFlowToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("validate flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["valid"] != false {
		t.Fatalf("expected invalid flow, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "file output root") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func toolNamesForTest(mcpServer *server.MCPServer) []string {
	value := reflect.ValueOf(mcpServer).Elem().FieldByName("tools")
	names := make([]string, 0, value.Len())
	for _, key := range value.MapKeys() {
		names = append(names, key.String())
	}
	sort.Strings(names)
	return names
}

func TestHandleRunFlowToolMissingFlow(t *testing.T) {
	result, err := handleRunFlowTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("run flow missing flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "flow") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func decodeToolText(t *testing.T, result *mcp.CallToolResult, target any) {
	t.Helper()
	if result == nil || len(result.Content) == 0 {
		t.Fatalf("empty tool result")
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("first content is %T, expected TextContent", result.Content[0])
	}
	if err := json.Unmarshal([]byte(text.Text), target); err != nil {
		t.Fatalf("decode tool text: %v\n%s", err, text.Text)
	}
}
