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
	want := []string{
		"tsplay.draft_flow",
		"tsplay.flow_examples",
		"tsplay.flow_schema",
		"tsplay.list_actions",
		"tsplay.observe_page",
		"tsplay.repair_flow_context",
		"tsplay.run_flow",
		"tsplay.validate_flow",
	}
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

func TestHandleFlowSchemaTool(t *testing.T) {
	result, err := handleFlowSchemaTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("flow schema: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	schema, ok := payload["schema"].(map[string]any)
	if !ok {
		t.Fatalf("expected schema, got %#v", payload["schema"])
	}
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("expected properties, got %#v", schema["properties"])
	}
	if _, ok := properties["steps"]; !ok {
		t.Fatalf("expected steps property")
	}
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatalf("expected $defs, got %#v", schema["$defs"])
	}
	manifest, ok := defs["action_manifest"].([]any)
	if !ok || len(manifest) == 0 {
		t.Fatalf("expected action manifest, got %#v", defs["action_manifest"])
	}
	if rules, ok := payload["generation_rules"].([]any); !ok || len(rules) == 0 {
		t.Fatalf("expected generation rules, got %#v", payload["generation_rules"])
	}
	if selectors, ok := payload["selector_strategy"].([]any); !ok || len(selectors) == 0 {
		t.Fatalf("expected selector strategy, got %#v", payload["selector_strategy"])
	}
}

func TestHandleFlowExamplesTool(t *testing.T) {
	result, err := handleFlowExamplesTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("flow examples: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	examples, ok := payload["examples"].([]any)
	if !ok || len(examples) < 3 {
		t.Fatalf("expected examples, got %#v", payload["examples"])
	}
	if hints, ok := payload["example_selection_hints"].([]any); !ok || len(hints) == 0 {
		t.Fatalf("expected example selection hints, got %#v", payload["example_selection_hints"])
	}

	foundExtractExample := false
	for _, example := range examples {
		item, ok := example.(map[string]any)
		if !ok {
			t.Fatalf("example is %T", example)
		}
		content, ok := item["flow"].(string)
		if !ok || strings.TrimSpace(content) == "" {
			t.Fatalf("example missing flow: %#v", item)
		}
		flow, err := ParseFlow([]byte(content), "yaml")
		if err != nil {
			t.Fatalf("parse example %q: %v", item["name"], err)
		}
		if err := ValidateFlow(flow); err != nil {
			t.Fatalf("validate example %q: %v", item["name"], err)
		}
		if item["name"] == "extract_text_and_set_var" {
			foundExtractExample = true
		}
	}
	if !foundExtractExample {
		t.Fatalf("missing extract_text_and_set_var example")
	}
}

func TestHandleDraftFlowToolWithObservation(t *testing.T) {
	observation := `{
  "url": "https://example.com/orders",
  "title": "Orders",
  "artifact_root": "/tmp/artifacts",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "text",
      "label": "Order keyword",
      "placeholder": "Search orders",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["[data-testid=\"order-query\"]", "#query"],
      "attributes": {"data-testid": "order-query"}
    },
    {
      "index": 2,
      "tag": "button",
      "type": "button",
      "text": "Search",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["#search-button", "text=\"Search\""]
    }
  ]
}`
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent":      "搜索订单",
				"observation": observation,
			},
		},
	}

	result, err := handleDraftFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("draft flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	draft, ok := payload["draft"].(map[string]any)
	if !ok {
		t.Fatalf("expected draft payload, got %#v", payload["draft"])
	}
	flowYAML, ok := draft["flow_yaml"].(string)
	if !ok || strings.TrimSpace(flowYAML) == "" {
		t.Fatalf("expected flow_yaml, got %#v", draft["flow_yaml"])
	}
	flow, err := ParseFlow([]byte(flowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse draft flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate draft flow: %v", err)
	}
}

func TestHandleDraftFlowToolRequiresIntentAndObservationOrURL(t *testing.T) {
	result, err := handleDraftFlowTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("draft flow missing intent: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "intent") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent": "搜索订单",
			},
		},
	}
	result, err = handleDraftFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("draft flow missing source: %v", err)
	}
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "url or observation") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func TestHandleRepairFlowContextTool(t *testing.T) {
	artifactRoot := t.TempDir()
	artifactDir := filepath.Join(artifactRoot, "run-1", "02-click")
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		t.Fatalf("create artifact dir: %v", err)
	}
	domPath := filepath.Join(artifactDir, "dom_snapshot.json")
	htmlPath := filepath.Join(artifactDir, "page.html")
	screenshotPath := filepath.Join(artifactDir, "failure.png")
	if err := os.WriteFile(domPath, []byte(`{"tag":"button","text":"Export orders","selector_candidates":["text=\"Export orders\"","#export"]}`), 0600); err != nil {
		t.Fatalf("write dom snapshot: %v", err)
	}
	if err := os.WriteFile(htmlPath, []byte(`<html><body>secret full html should stay on disk</body></html>`), 0600); err != nil {
		t.Fatalf("write html: %v", err)
	}
	if err := os.WriteFile(screenshotPath, []byte("png"), 0600); err != nil {
		t.Fatalf("write screenshot: %v", err)
	}

	result := FlowResult{
		Name:         "repair_me",
		ArtifactRoot: artifactRoot,
		Vars: map[string]any{
			"orders_url": "https://example.com/orders",
		},
		Trace: []FlowStepTrace{
			{
				Index:       1,
				Action:      "navigate",
				Status:      "ok",
				PageURL:     "https://example.com/orders",
				DurationMS:  42,
				ArgsSummary: `{"url":"https://example.com/orders"}`,
			},
			{
				Index:       2,
				Name:        "click export",
				Action:      "click",
				Status:      "error",
				ArgsSummary: `{"selector":"text=\"Old export\""}`,
				Error:       "locator click: timeout",
				ErrorStack:  strings.Repeat("stack frame\n", 20),
				PageURL:     "https://example.com/orders",
				Artifacts: &FlowStepArtifacts{
					Directory:       artifactDir,
					ScreenshotPath:  screenshotPath,
					HTMLPath:        htmlPath,
					DOMSnapshotPath: domPath,
				},
			},
		},
	}
	wrappedRunResult, err := json.Marshal(map[string]any{
		"ok":     false,
		"error":  "flow failed",
		"result": result,
	})
	if err != nil {
		t.Fatalf("marshal run result: %v", err)
	}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: repair_me
vars:
  orders_url: https://example.com/orders
steps:
  - action: navigate
    url: "{{orders_url}}"
  - name: click export
    action: click
    selector: 'text="Old export"'
`,
				"run_result":           string(wrappedRunResult),
				"max_artifact_excerpt": 40,
			},
		},
	}

	toolResult, err := handleRepairFlowContextToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{
		ArtifactRoot: artifactRoot,
	})
	if err != nil {
		t.Fatalf("repair flow context: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, toolResult, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	contextPayload, ok := payload["context"].(map[string]any)
	if !ok {
		t.Fatalf("expected context, got %#v", payload["context"])
	}
	failedStep, ok := contextPayload["failed_step"].(map[string]any)
	if !ok {
		t.Fatalf("expected failed step, got %#v", contextPayload["failed_step"])
	}
	if failedStep["index"] != float64(2) {
		t.Fatalf("unexpected failed step index: %#v", failedStep["index"])
	}
	artifacts, ok := contextPayload["artifacts"].(map[string]any)
	if !ok {
		t.Fatalf("expected artifacts, got %#v", contextPayload["artifacts"])
	}
	excerpt, ok := artifacts["dom_snapshot_excerpt"].(string)
	if !ok || !strings.Contains(excerpt, "Export orders") {
		t.Fatalf("expected dom excerpt, got %#v", artifacts["dom_snapshot_excerpt"])
	}
	encoded, err := json.Marshal(contextPayload)
	if err != nil {
		t.Fatalf("marshal context: %v", err)
	}
	if strings.Contains(string(encoded), "secret full html") {
		t.Fatalf("context leaked full html content: %s", encoded)
	}
	if !strings.Contains(string(encoded), htmlPath) {
		t.Fatalf("expected html path in context: %s", encoded)
	}
	if contextPayload["failure_category"] != "selector_or_timing" {
		t.Fatalf("unexpected failure category: %#v", contextPayload["failure_category"])
	}
	if _, ok := contextPayload["validation_checklist"].([]any); !ok {
		t.Fatalf("expected validation checklist, got %#v", contextPayload["validation_checklist"])
	}
	focusedVars, ok := contextPayload["focused_variables"].(map[string]any)
	if !ok || focusedVars["orders_url"] != "https://example.com/orders" {
		t.Fatalf("unexpected focused variables: %#v", contextPayload["focused_variables"])
	}
}

func TestHandleRepairFlowContextToolMissingTrace(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: repair_missing_trace
steps:
  - action: navigate
    url: https://example.com
`,
			},
		},
	}

	result, err := handleRepairFlowContextTool(context.Background(), request)
	if err != nil {
		t.Fatalf("repair flow context missing trace: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "run_result or trace") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func TestBuildFlowRepairContextRestrictsArtifactRoot(t *testing.T) {
	artifactRoot := t.TempDir()
	outsideDir := t.TempDir()
	outsideDOMPath := filepath.Join(outsideDir, "dom_snapshot.json")
	if err := os.WriteFile(outsideDOMPath, []byte("outside secret"), 0600); err != nil {
		t.Fatalf("write outside dom: %v", err)
	}

	contextPayload, err := BuildFlowRepairContext(FlowRepairContextOptions{
		Flow: &Flow{
			SchemaVersion: CurrentFlowSchemaVersion,
			Name:          "restrict_artifact_read",
			Steps: []FlowStep{
				{Action: "click", Selector: "#missing"},
			},
		},
		Result: &FlowResult{
			Trace: []FlowStepTrace{
				{
					Index:  1,
					Action: "click",
					Status: "error",
					Error:  "timeout",
					Artifacts: &FlowStepArtifacts{
						DOMSnapshotPath: outsideDOMPath,
					},
				},
			},
		},
		ArtifactRoot: artifactRoot,
	})
	if err != nil {
		t.Fatalf("build repair context: %v", err)
	}
	if contextPayload.Artifacts == nil {
		t.Fatalf("expected artifact context")
	}
	if contextPayload.Artifacts.DOMSnapshotExcerpt != "" {
		t.Fatalf("unexpected outside dom excerpt: %q", contextPayload.Artifacts.DOMSnapshotExcerpt)
	}
	if len(contextPayload.Artifacts.ReadErrors) == 0 {
		t.Fatalf("expected read error for outside artifact")
	}
	encoded, err := json.Marshal(contextPayload)
	if err != nil {
		t.Fatalf("marshal context: %v", err)
	}
	if strings.Contains(string(encoded), "outside secret") {
		t.Fatalf("context leaked outside artifact content: %s", encoded)
	}
}

func TestHandleObservePageToolMissingURL(t *testing.T) {
	result, err := handleObservePageTool(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("observe page missing url: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "url") {
		t.Fatalf("unexpected error: %#v", payload["error"])
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

func TestFlowResultForToolCompactsVars(t *testing.T) {
	result := flowResultForTool(&FlowResult{
		Name: "large_vars",
		Vars: map[string]any{
			"html": strings.Repeat("<div>content</div>", 200),
		},
	})

	html, ok := result.Vars["html"].(string)
	if !ok {
		t.Fatalf("expected html string, got %#v", result.Vars["html"])
	}
	if len(html) > 1200 {
		t.Fatalf("expected compacted html, got length %d", len(html))
	}
	if !strings.Contains(html, "truncated") {
		t.Fatalf("expected truncation marker, got %q", html)
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
