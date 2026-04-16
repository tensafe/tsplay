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
		"tsplay.delete_session",
		"tsplay.draft_flow",
		"tsplay.export_session_flow_snippet",
		"tsplay.flow_examples",
		"tsplay.flow_schema",
		"tsplay.get_session",
		"tsplay.list_actions",
		"tsplay.list_sessions",
		"tsplay.observe_page",
		"tsplay.repair_flow",
		"tsplay.repair_flow_context",
		"tsplay.run_flow",
		"tsplay.save_session",
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

func TestHandleSaveSessionToolWithStorageStateJSON(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name":          "admin",
				"storage_state": `{"cookies":[],"origins":[]}`,
			},
		},
	}

	result, err := handleSaveSessionToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("save session: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	session, ok := payload["session"].(map[string]any)
	if !ok {
		t.Fatalf("expected session payload, got %#v", payload["session"])
	}
	if session["name"] != "admin" {
		t.Fatalf("unexpected session name: %#v", session["name"])
	}
	if session["source_type"] != "inline_storage_state" {
		t.Fatalf("expected source_type inline_storage_state, got %#v", session["source_type"])
	}
	if _, ok := session["source"].(string); !ok {
		t.Fatalf("expected source description, got %#v", session["source"])
	}
	browser, ok := session["browser"].(map[string]any)
	if !ok || browser["use_session"] != "admin" {
		t.Fatalf("expected browser use_session, got %#v", session["browser"])
	}
	if _, err := os.Stat(filepath.Join(options.ArtifactRoot, "sessions", "storage", "admin.json")); err != nil {
		t.Fatalf("expected saved storage state file: %v", err)
	}
}

func TestHandleListSessionsTool(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     options.ArtifactRoot,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:         "operator",
		ArtifactRoot: options.ArtifactRoot,
		Profile:      "crm",
		Session:      "operator",
	}); err != nil {
		t.Fatalf("seed operator session: %v", err)
	}

	result, err := handleListSessionsToolWithOptions(context.Background(), mcp.CallToolRequest{}, options)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	sessions, ok := payload["sessions"].([]any)
	if !ok || len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %#v", payload["sessions"])
	}
	first, ok := sessions[0].(map[string]any)
	if !ok {
		t.Fatalf("expected session item, got %#v", sessions[0])
	}
	if _, ok := first["browser"].(map[string]any); !ok {
		t.Fatalf("expected browser snippet, got %#v", first["browser"])
	}
	if _, ok := first["resolved_browser"].(map[string]any); !ok {
		t.Fatalf("expected resolved_browser snippet, got %#v", first["resolved_browser"])
	}
	if _, ok := first["source"].(string); !ok {
		t.Fatalf("expected source description, got %#v", first["source"])
	}
}

func TestHandleGetSessionToolForStorageState(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     options.ArtifactRoot,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "admin",
			},
		},
	}
	result, err := handleGetSessionToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	session, ok := payload["session"].(map[string]any)
	if !ok {
		t.Fatalf("expected session payload, got %#v", payload["session"])
	}
	expanded, ok := session["expanded_browser"].(map[string]any)
	if !ok || expanded["storage_state"] == nil {
		t.Fatalf("expected expanded browser storage_state, got %#v", session["expanded_browser"])
	}
	physical, ok := session["physical_paths"].(map[string]any)
	if !ok {
		t.Fatalf("expected physical_paths, got %#v", session["physical_paths"])
	}
	if _, ok := physical["metadata_path"].(string); !ok {
		t.Fatalf("expected metadata_path, got %#v", physical["metadata_path"])
	}
	if _, ok := physical["storage_state_path"].(string); !ok {
		t.Fatalf("expected storage_state_path, got %#v", physical["storage_state_path"])
	}
}

func TestHandleGetSessionToolForPersistentProfile(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:         "crm-admin",
		ArtifactRoot: options.ArtifactRoot,
		Profile:      "crm",
		Session:      "admin",
	}); err != nil {
		t.Fatalf("seed crm-admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "crm-admin",
			},
		},
	}
	result, err := handleGetSessionToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	session, ok := payload["session"].(map[string]any)
	if !ok {
		t.Fatalf("expected session payload, got %#v", payload["session"])
	}
	expanded, ok := session["expanded_browser"].(map[string]any)
	if !ok || expanded["persistent"] != true {
		t.Fatalf("expected persistent expanded browser, got %#v", session["expanded_browser"])
	}
	physical, ok := session["physical_paths"].(map[string]any)
	if !ok {
		t.Fatalf("expected physical_paths, got %#v", session["physical_paths"])
	}
	if _, ok := physical["profile_dir"].(string); !ok {
		t.Fatalf("expected profile_dir, got %#v", physical["profile_dir"])
	}
}

func TestHandleExportSessionFlowSnippetToolForStorageState(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     options.ArtifactRoot,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "admin",
			},
		},
	}
	result, err := handleExportSessionFlowSnippetToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("export session flow snippet: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	exported, ok := payload["export"].(map[string]any)
	if !ok {
		t.Fatalf("expected export payload, got %#v", payload["export"])
	}
	if exported["format"] != "all" {
		t.Fatalf("expected format all, got %#v", exported["format"])
	}
	snippets, ok := exported["snippets"].(map[string]any)
	if !ok {
		t.Fatalf("expected snippets, got %#v", exported["snippets"])
	}
	browser, ok := snippets["browser"].(map[string]any)
	if !ok || browser["use_session"] != "admin" {
		t.Fatalf("expected browser use_session snippet, got %#v", snippets["browser"])
	}
	browserYAML, ok := snippets["browser_yaml"].(string)
	if !ok || !strings.Contains(browserYAML, "use_session: admin") {
		t.Fatalf("unexpected browser_yaml: %#v", snippets["browser_yaml"])
	}
	expandedBrowserYAML, ok := snippets["expanded_browser_yaml"].(string)
	if !ok || !strings.Contains(expandedBrowserYAML, "storage_state: sessions/storage/admin.json") {
		t.Fatalf("unexpected expanded_browser_yaml: %#v", snippets["expanded_browser_yaml"])
	}
	flowYAML, ok := snippets["flow_yaml"].(string)
	if !ok || !strings.Contains(flowYAML, "schema_version: \"1\"") || !strings.Contains(flowYAML, "steps: []") {
		t.Fatalf("unexpected flow_yaml: %#v", snippets["flow_yaml"])
	}
}

func TestHandleExportSessionFlowSnippetToolForPersistentProfile(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:         "crm-admin",
		ArtifactRoot: options.ArtifactRoot,
		Profile:      "crm",
		Session:      "admin",
	}); err != nil {
		t.Fatalf("seed crm-admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name":   "crm-admin",
				"format": "expanded_flow_json",
			},
		},
	}
	result, err := handleExportSessionFlowSnippetToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("export session flow snippet: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	exported, ok := payload["export"].(map[string]any)
	if !ok {
		t.Fatalf("expected export payload, got %#v", payload["export"])
	}
	if exported["format"] != "expanded_flow_json" {
		t.Fatalf("unexpected export format: %#v", exported["format"])
	}
	if exported["target"] != "flow" || exported["encoding"] != "json" || exported["variant"] != "expanded" {
		t.Fatalf("unexpected export metadata: %#v", exported)
	}
	snippet, ok := exported["snippet"].(string)
	if !ok {
		t.Fatalf("expected snippet string, got %#v", exported["snippet"])
	}
	if !strings.Contains(snippet, "\"persistent\": true") || !strings.Contains(snippet, "\"profile\": \"crm\"") || !strings.Contains(snippet, "\"session\": \"admin\"") {
		t.Fatalf("unexpected json snippet: %#v", snippet)
	}
	snippetData, ok := exported["snippet_data"].(map[string]any)
	if !ok {
		t.Fatalf("expected snippet_data, got %#v", exported["snippet_data"])
	}
	browser, ok := snippetData["browser"].(map[string]any)
	if !ok || browser["persistent"] != true {
		t.Fatalf("unexpected browser in snippet_data: %#v", snippetData["browser"])
	}
}

func TestHandleExportSessionFlowSnippetToolRejectsUnknownFormat(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     options.ArtifactRoot,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name":   "admin",
				"format": "json",
			},
		},
	}
	result, err := handleExportSessionFlowSnippetToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("export session flow snippet: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "unsupported format") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
}

func TestHandleDeleteSessionTool(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     options.ArtifactRoot,
		StorageStateJSON: `{"cookies":[],"origins":[]}`,
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "admin",
			},
		},
	}
	result, err := handleDeleteSessionToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("delete session: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	deleted, ok := payload["deleted"].(map[string]any)
	if !ok {
		t.Fatalf("expected deleted payload, got %#v", payload["deleted"])
	}
	if deleted["deleted_storage_state"] != true {
		t.Fatalf("expected deleted_storage_state=true, got %#v", deleted["deleted_storage_state"])
	}
	if _, err := os.Stat(filepath.Join(options.ArtifactRoot, "sessions", "registry", "admin.json")); !os.IsNotExist(err) {
		t.Fatalf("expected metadata file to be removed, stat err=%v", err)
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
	if _, ok := properties["browser"]; !ok {
		t.Fatalf("expected browser property")
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
	foundBrowserExample := false
	foundNamedSessionExample := false
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
		if item["name"] == "browser_session_with_storage_state" {
			foundBrowserExample = true
		}
		if item["name"] == "reuse_named_session" {
			foundNamedSessionExample = true
		}
	}
	if !foundExtractExample {
		t.Fatalf("missing extract_text_and_set_var example")
	}
	if !foundBrowserExample {
		t.Fatalf("missing browser_session_with_storage_state example")
	}
	if !foundNamedSessionExample {
		t.Fatalf("missing reuse_named_session example")
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
	validation, ok := draft["validation"].(map[string]any)
	if !ok || validation["valid"] != true {
		t.Fatalf("expected validation.valid=true, got %#v", draft["validation"])
	}
}

func TestHandleDraftFlowToolAutoRepairsSelectors(t *testing.T) {
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
      "selector_candidates": ["#query", "[data-testid=\"order-query\"]"],
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
	draft, ok := payload["draft"].(map[string]any)
	if !ok {
		t.Fatalf("expected draft payload, got %#v", payload["draft"])
	}
	if draft["auto_repaired"] != true {
		t.Fatalf("expected auto_repaired=true, got %#v", draft["auto_repaired"])
	}
	repairs, ok := draft["selector_repairs"].([]any)
	if !ok || len(repairs) == 0 {
		t.Fatalf("expected selector repairs, got %#v", draft["selector_repairs"])
	}
	flowYAML, ok := draft["flow_yaml"].(string)
	if !ok || !strings.Contains(flowYAML, `[data-testid="order-query"]`) {
		t.Fatalf("expected repaired selector in flow yaml, got %q", flowYAML)
	}
}

func TestHandleDraftFlowToolReturnsRepairHintsWhenValidationFails(t *testing.T) {
	observation := `{
  "url": "https://example.com/upload",
  "title": "Upload",
  "artifact_root": "/tmp/artifacts",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "file",
      "label": "选择文件",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["#fileInput"]
    },
    {
      "index": 2,
      "tag": "button",
      "type": "submit",
      "text": "上传文件",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["text=\"上传文件\""]
    }
  ]
}`
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent":      "上传文件并提交",
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
	draft, ok := payload["draft"].(map[string]any)
	if !ok {
		t.Fatalf("expected draft payload, got %#v", payload["draft"])
	}
	validation, ok := draft["validation"].(map[string]any)
	if !ok || validation["valid"] != false {
		t.Fatalf("expected validation.valid=false, got %#v", draft["validation"])
	}
	hints, ok := draft["repair_hints"].([]any)
	if !ok || len(hints) == 0 {
		t.Fatalf("expected repair_hints, got %#v", draft["repair_hints"])
	}
	firstHint, ok := hints[0].(map[string]any)
	if !ok {
		t.Fatalf("expected hint object, got %#v", hints[0])
	}
	if firstHint["step_path"] != "3" {
		t.Fatalf("expected hint to target step 3, got %#v", firstHint)
	}
	suggestion, _ := firstHint["suggestion"].(string)
	if !strings.Contains(suggestion, "allow_file_access=true") {
		t.Fatalf("expected allow_file_access suggestion, got %#v", firstHint)
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
	repairHints, ok := contextPayload["repair_hints"].([]any)
	if !ok || len(repairHints) == 0 {
		t.Fatalf("expected repair_hints, got %#v", contextPayload["repair_hints"])
	}
	firstHint, ok := repairHints[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first repair hint object, got %#v", repairHints[0])
	}
	if firstHint["source"] != "runtime_failure" {
		t.Fatalf("expected runtime_failure hint source, got %#v", firstHint)
	}
	if firstHint["step_path"] != "2" {
		t.Fatalf("expected failed step path 2, got %#v", firstHint)
	}
	if firstHint["failure_category"] != "selector_or_timing" {
		t.Fatalf("expected selector_or_timing hint category, got %#v", firstHint)
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

func TestHandleRepairFlowToolWithRepairHints(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: upload_orders
steps:
  - action: navigate
    url: https://example.com/upload
  - action: wait_for_selector
    selector: "#fileInput"
    timeout: 10000
  - action: upload_file
    selector: "#fileInput"
    file_path: "{{upload_file_path}}"
`,
				"repair_hints": `{
  "draft": {
    "repair_hints": [
      {
        "priority": 1,
        "source": "draft_validation",
        "step_path": "3",
        "action": "upload_file",
        "targets": ["security_policy", "action"],
        "reason": "The drafted step is valid structurally, but it is blocked by the current safety flags.",
        "suggestion": "Inspect step 3 first. If this is a trusted automation, rerun draft_flow or validate_flow with allow_file_access=true; otherwise replace the step with a lower-risk action.",
        "error": "step 3 action \"upload_file\" is disabled by security policy; set allow_file_access=true only for trusted flows",
        "failure_category": "validation"
      }
    ]
  }
}`,
			},
		},
	}

	result, err := handleRepairFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("repair flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	repair, ok := payload["repair"].(map[string]any)
	if !ok {
		t.Fatalf("expected repair payload, got %#v", payload["repair"])
	}
	targetSteps, ok := repair["target_steps"].([]any)
	if !ok || len(targetSteps) != 1 || targetSteps[0] != "3" {
		t.Fatalf("expected target step 3, got %#v", repair["target_steps"])
	}
	prompt, ok := repair["prompt"].(string)
	if !ok || !strings.Contains(prompt, "Original flow:") || !strings.Contains(prompt, "step=3 action=upload_file") {
		t.Fatalf("unexpected repair prompt: %q", prompt)
	}
}

func TestHandleRepairFlowToolBuildsContextFromRunResult(t *testing.T) {
	artifactRoot := t.TempDir()
	artifactDir := filepath.Join(artifactRoot, "run-1", "02-click")
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		t.Fatalf("create artifact dir: %v", err)
	}
	domPath := filepath.Join(artifactDir, "dom_snapshot.json")
	htmlPath := filepath.Join(artifactDir, "page.html")
	screenshotPath := filepath.Join(artifactDir, "failure.png")
	if err := os.WriteFile(domPath, []byte(`{"tag":"button","text":"Export orders"}`), 0600); err != nil {
		t.Fatalf("write dom snapshot: %v", err)
	}
	if err := os.WriteFile(htmlPath, []byte(`<html><body>secret full html should stay on disk</body></html>`), 0600); err != nil {
		t.Fatalf("write html: %v", err)
	}
	if err := os.WriteFile(screenshotPath, []byte("png"), 0600); err != nil {
		t.Fatalf("write screenshot: %v", err)
	}

	resultPayload := FlowResult{
		Name:         "repair_me",
		ArtifactRoot: artifactRoot,
		Vars: map[string]any{
			"orders_url": "https://example.com/orders",
		},
		Trace: []FlowStepTrace{
			{
				Index:      1,
				Action:     "navigate",
				Status:     "ok",
				PageURL:    "https://example.com/orders",
				DurationMS: 42,
			},
			{
				Index:      2,
				Name:       "click export",
				Action:     "click",
				Status:     "error",
				Error:      "locator click: timeout",
				PageURL:    "https://example.com/orders",
				DurationMS: 15000,
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
		"result": resultPayload,
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

	toolResult, err := handleRepairFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{
		ArtifactRoot: artifactRoot,
	})
	if err != nil {
		t.Fatalf("repair flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, toolResult, &payload)
	repair, ok := payload["repair"].(map[string]any)
	if !ok {
		t.Fatalf("expected repair payload, got %#v", payload["repair"])
	}
	hints, ok := repair["repair_hints"].([]any)
	if !ok || len(hints) == 0 {
		t.Fatalf("expected repair hints, got %#v", repair["repair_hints"])
	}
	firstHint, ok := hints[0].(map[string]any)
	if !ok || firstHint["source"] != "runtime_failure" || firstHint["step_path"] != "2" {
		t.Fatalf("unexpected runtime repair hint: %#v", hints[0])
	}
	prompt, ok := repair["prompt"].(string)
	if !ok || !strings.Contains(prompt, "Failure context:") || !strings.Contains(prompt, "step=2 action=click") {
		t.Fatalf("unexpected repair prompt: %q", prompt)
	}
	if strings.Contains(prompt, "secret full html") {
		t.Fatalf("repair prompt leaked full html: %s", prompt)
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

func TestHandleValidateFlowToolRejectsHTTPWithoutAllow(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_http_from_mcp
steps:
  - action: http_request
    url: https://example.com/api
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
	if !strings.Contains(payload["error"].(string), "allow_http") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}

	allowedRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_http_from_mcp
steps:
  - action: http_request
    url: https://example.com/api
`,
				"allow_http": true,
			},
		},
	}
	result, err = handleValidateFlowTool(context.Background(), allowedRequest)
	if err != nil {
		t.Fatalf("validate flow with allow_http: %v", err)
	}
	decodeToolText(t, result, &payload)
	if payload["valid"] != true {
		t.Fatalf("expected valid flow, got %#v", payload)
	}
}

func TestHandleValidateFlowToolRejectsRedisWithoutAllow(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_redis_from_mcp
steps:
  - action: redis_get
    key: sessions:admin_cookie
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
	if !strings.Contains(payload["error"].(string), "allow_redis") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}

	allowedRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_redis_from_mcp
steps:
  - action: redis_get
    key: sessions:admin_cookie
`,
				"allow_redis": true,
			},
		},
	}
	result, err = handleValidateFlowTool(context.Background(), allowedRequest)
	if err != nil {
		t.Fatalf("validate flow with allow_redis: %v", err)
	}
	decodeToolText(t, result, &payload)
	if payload["valid"] != true {
		t.Fatalf("expected valid flow, got %#v", payload)
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
