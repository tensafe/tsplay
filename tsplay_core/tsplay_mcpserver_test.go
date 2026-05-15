package tsplay_core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		"tsplay.finalize_flow",
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

func TestWithDraftFlowObservationInputAcceptsStringOrObject(t *testing.T) {
	tool := mcp.NewTool("test.tool", withDraftFlowObservationInput())
	observationSchema, ok := tool.InputSchema.Properties["observation"].(map[string]any)
	if !ok {
		t.Fatalf("expected observation schema, got %#v", tool.InputSchema.Properties["observation"])
	}
	oneOf, ok := observationSchema["oneOf"].([]any)
	if !ok || len(oneOf) != 2 {
		t.Fatalf("expected observation oneOf schema, got %#v", observationSchema["oneOf"])
	}
	stringSchema, _ := oneOf[0].(map[string]any)
	objectSchema, _ := oneOf[1].(map[string]any)
	if stringSchema["type"] != "string" {
		t.Fatalf("expected string schema first, got %#v", stringSchema)
	}
	if objectSchema["type"] != "object" || objectSchema["additionalProperties"] != true {
		t.Fatalf("expected flexible object schema, got %#v", objectSchema)
	}
}

func TestWithBrowserCDPPortInputDeclaresRange(t *testing.T) {
	tool := mcp.NewTool("test.tool", withBrowserCDPPortInput("desc"))
	portSchema, ok := tool.InputSchema.Properties["browser_cdp_port"].(map[string]any)
	if !ok {
		t.Fatalf("expected browser_cdp_port schema, got %#v", tool.InputSchema.Properties["browser_cdp_port"])
	}
	if portSchema["type"] != "number" {
		t.Fatalf("expected number schema, got %#v", portSchema)
	}
	if portSchema["minimum"] != float64(1) || portSchema["maximum"] != float64(65535) {
		t.Fatalf("expected port range 1..65535, got %#v", portSchema)
	}
}

func TestBrowserCDPPortFromToolRequestValidatesIntegerRange(t *testing.T) {
	valid := []any{9222, int64(9222), float64(9222), json.Number("9222"), "9222"}
	for _, value := range valid {
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]any{"browser_cdp_port": value},
			},
		}
		port, present, err := browserCDPPortFromToolRequest(request)
		if err != nil {
			t.Fatalf("expected valid port for %#v, got %v", value, err)
		}
		if !present || port != 9222 {
			t.Fatalf("expected port 9222 for %#v, got present=%v port=%d", value, present, port)
		}
	}

	invalid := []any{0, -1, 65536, 9222.5, json.Number("9222.5"), "abc", true}
	for _, value := range invalid {
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]any{"browser_cdp_port": value},
			},
		}
		if _, _, err := browserCDPPortFromToolRequest(request); err == nil {
			t.Fatalf("expected invalid port error for %#v", value)
		}
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
	foundWriteExcel := false
	foundSendEmail := false
	foundZipCompress := false
	foundZipExtract := false
	for _, action := range actions {
		item, ok := action.(map[string]any)
		if !ok {
			continue
		}
		if item["name"] == "navigate" {
			foundNavigate = true
		}
		if item["name"] == "write_excel" {
			foundWriteExcel = true
		}
		if item["name"] == "send_email" {
			foundSendEmail = true
		}
		if item["name"] == "zip_compress" {
			foundZipCompress = true
		}
		if item["name"] == "zip_extract" {
			foundZipExtract = true
		}
	}
	if !foundNavigate {
		t.Fatalf("navigate action not found in manifest")
	}
	if !foundWriteExcel {
		t.Fatalf("write_excel action not found in manifest")
	}
	if !foundSendEmail {
		t.Fatalf("send_email action not found in manifest")
	}
	if !foundZipCompress {
		t.Fatalf("zip_compress action not found in manifest")
	}
	if !foundZipExtract {
		t.Fatalf("zip_extract action not found in manifest")
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

func TestSessionToolsHideOwnedSessionsFromOtherMCPCallers(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	if _, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:               "admin",
		ArtifactRoot:       options.ArtifactRoot,
		StorageStateJSON:   `{"cookies":[],"origins":[]}`,
		OwnerSessionID:     "session-owner",
		OwnerClientName:    "codex",
		OwnerClientVersion: "1.0.0",
	}); err != nil {
		t.Fatalf("seed admin session: %v", err)
	}

	makeCtx := func(id string) context.Context {
		session := &runtimeTestSession{
			id:          id,
			initialized: true,
			notify:      make(chan mcp.JSONRPCNotification, 1),
			clientInfo:  mcp.Implementation{Name: "codex", Version: "1.0.0"},
		}
		return server.NewMCPServer("test", "1.0.0").WithContext(context.Background(), session)
	}

	otherCtx := makeCtx("session-other")
	listResult, err := handleListSessionsToolWithOptions(otherCtx, mcp.CallToolRequest{}, options)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	var listPayload map[string]any
	decodeToolText(t, listResult, &listPayload)
	sessions, ok := listPayload["sessions"].([]any)
	if !ok || len(sessions) != 0 {
		t.Fatalf("expected no visible sessions for other caller, got %#v", listPayload["sessions"])
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{"name": "admin"},
		},
	}
	getResult, err := handleGetSessionToolWithOptions(otherCtx, request, options)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	var getPayload map[string]any
	decodeToolText(t, getResult, &getPayload)
	if getPayload["ok"] != false || !strings.Contains(getPayload["error"].(string), "owned by MCP session") {
		t.Fatalf("expected ownership error, got %#v", getPayload)
	}

	exportResult, err := handleExportSessionFlowSnippetToolWithOptions(otherCtx, request, options)
	if err != nil {
		t.Fatalf("export session: %v", err)
	}
	var exportPayload map[string]any
	decodeToolText(t, exportResult, &exportPayload)
	if exportPayload["ok"] != false || !strings.Contains(exportPayload["error"].(string), "owned by MCP session") {
		t.Fatalf("expected ownership error, got %#v", exportPayload)
	}

	ownerCtx := makeCtx("session-owner")
	ownerGetResult, err := handleGetSessionToolWithOptions(ownerCtx, request, options)
	if err != nil {
		t.Fatalf("owner get session: %v", err)
	}
	var ownerPayload map[string]any
	decodeToolText(t, ownerGetResult, &ownerPayload)
	if ownerPayload["ok"] != true {
		t.Fatalf("expected owner access, got %#v", ownerPayload)
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
	foundNavigate := false
	foundHTTPRequest := false
	for _, raw := range manifest {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("expected manifest item object, got %#v", raw)
		}
		capabilities, ok := item["capabilities"].(map[string]any)
		if !ok {
			t.Fatalf("expected capabilities on manifest item, got %#v", item["capabilities"])
		}
		switch item["name"] {
		case "navigate":
			foundNavigate = true
			if capabilities["needs_playwright"] != true || capabilities["needs_runtime"] != true || capabilities["needs_page"] != true {
				t.Fatalf("navigate capabilities = %#v", capabilities)
			}
		case "http_request":
			foundHTTPRequest = true
			args, ok := capabilities["conditional_playwright_args"].([]any)
			if !ok || len(args) == 0 {
				t.Fatalf("http_request conditional capabilities = %#v", capabilities)
			}
		}
	}
	if !foundNavigate {
		t.Fatalf("expected navigate manifest item")
	}
	if !foundHTTPRequest {
		t.Fatalf("expected http_request manifest item")
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
  "page_summary": "Observed \"Orders\" with 2 interactive elements and 2 content elements. Top content: Order Center; Export orders.",
  "dom_snapshot_excerpt": "h1: Order Center\na: Export orders -> /export",
  "content_elements": [
    {
      "index": 1,
      "kind": "headline",
      "tag": "h1",
      "text": "Order Center",
      "xpath": "/html/body/main/h1[1]",
      "selector": "xpath=/html/body/main/h1[1]"
    },
    {
      "index": 2,
      "kind": "article_link",
      "tag": "a",
      "text": "Export orders",
      "href": "/export",
      "xpath": "/html/body/main/a[1]",
      "selector": "xpath=/html/body/main/a[1]"
    }
  ],
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
	if payload["tool"] != "tsplay.draft_flow" {
		t.Fatalf("expected tool metadata, got %#v", payload["tool"])
	}
	if _, ok := payload["summary"].(string); !ok {
		t.Fatalf("expected summary string, got %#v", payload["summary"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.run_flow" {
		t.Fatalf("expected next_action tsplay.run_flow, got %#v", payload["next_action"])
	}
	run, ok := payload["run"].(map[string]any)
	if !ok || run["status"] != "not_started" {
		t.Fatalf("expected synthetic run metadata for observation-only draft, got %#v", payload["run"])
	}
	runDetails, ok := run["details"].(map[string]any)
	if !ok || runDetails["source"] != "provided_observation" {
		t.Fatalf("expected provided_observation synthetic run source, got %#v", run["details"])
	}
	security, ok := payload["security"].(map[string]any)
	if !ok {
		t.Fatalf("expected security metadata, got %#v", payload["security"])
	}
	if preset, ok := security["preset"].(string); !ok || preset != "" {
		t.Fatalf("expected empty preset by default, got %#v", security["preset"])
	}
	if summary, ok := payload["page_summary"].(string); !ok || !strings.Contains(summary, "Order Center") {
		t.Fatalf("expected page_summary passthrough, got %#v", payload["page_summary"])
	}
	if excerpt, ok := payload["dom_snapshot_excerpt"].(string); !ok || !strings.Contains(excerpt, "Export orders") {
		t.Fatalf("expected dom_snapshot_excerpt passthrough, got %#v", payload["dom_snapshot_excerpt"])
	}
	content, ok := payload["content_elements"].([]any)
	if !ok || len(content) != 2 {
		t.Fatalf("expected content_elements passthrough, got %#v", payload["content_elements"])
	}
}

func TestHandleDraftFlowToolExposesTopLevelIssue(t *testing.T) {
	observation := `{
  "url": "https://example.com/upload",
  "title": "Upload",
  "artifact_root": "/tmp/artifacts",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "file",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["#file-input"]
    },
    {
      "index": 2,
      "tag": "button",
      "type": "button",
      "text": "Upload",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["text=\"Upload\""]
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
	issue, ok := payload["issue"].(map[string]any)
	if !ok {
		t.Fatalf("expected top-level issue payload, got %#v", payload["issue"])
	}
	if issue["code"] != "security_policy" {
		t.Fatalf("unexpected issue payload: %#v", issue)
	}
}

func TestHandleDraftFlowToolUsesContentElementsForContentIntent(t *testing.T) {
	observation := `{
  "url": "https://money.163.com/",
  "title": "网易财经",
  "artifact_root": "/tmp/artifacts",
  "page_summary": "Observed \"网易财经\" with 1 interactive elements and 3 content elements. Top content: 财经要闻; 头条新闻一; 头条新闻二.",
  "dom_snapshot_excerpt": "h2: 财经要闻\na: 头条新闻一 -> https://money.163.com/story/1",
  "content_elements": [
    {
      "index": 1,
      "kind": "headline",
      "tag": "h2",
      "text": "财经要闻",
      "xpath": "/html/body/main/section[1]/h2[1]",
      "selector": "xpath=/html/body/main/section[1]/h2[1]"
    },
    {
      "index": 2,
      "kind": "article_link",
      "tag": "a",
      "text": "头条新闻一",
      "href": "https://money.163.com/story/1",
      "xpath": "/html/body/main/section[1]/ul[1]/li[1]/a[1]",
      "selector": "xpath=/html/body/main/section[1]/ul[1]/li[1]/a[1]"
    },
    {
      "index": 3,
      "kind": "article_link",
      "tag": "a",
      "text": "头条新闻二",
      "href": "https://money.163.com/story/2",
      "xpath": "/html/body/main/section[1]/ul[1]/li[2]/a[1]",
      "selector": "xpath=/html/body/main/section[1]/ul[1]/li[2]/a[1]"
    }
  ],
  "elements": [
    {
      "index": 1,
      "tag": "a",
      "type": "link",
      "text": "头条新闻一",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["xpath=/html/body/main/section[1]/ul[1]/li[1]/a[1]"]
    }
  ]
}`
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent":      "查看财经要闻内容",
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
	flowYAML, ok := draft["flow_yaml"].(string)
	if !ok || !strings.Contains(flowYAML, `action: get_all_links`) {
		t.Fatalf("expected content-oriented draft flow, got %#v", draft["flow_yaml"])
	}
}

func TestHandleDraftFlowToolAcceptsStructuredObservationObject(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent": "查看财经要闻内容",
				"observation": map[string]any{
					"title": "网易财经",
					"content_elements\\": []any{
						map[string]any{
							"index":    1,
							"kind":     "headline",
							"tag":      "h2",
							"text":     "财经要闻",
							"xpath":    "/html/body/main/section[1]/h2[1]",
							"selector": "xpath=/html/body/main/section[1]/h2[1]",
						},
						map[string]any{
							"index":    2,
							"kind":     "article_link",
							"tag":      "a",
							"text":     "头条新闻一",
							"href":     "https://money.163.com/story/1",
							"xpath":    "/html/body/main/section[1]/ul[1]/li[1]/a[1]",
							"selector": "xpath=/html/body/main/section[1]/ul[1]/li[1]/a[1]",
						},
						map[string]any{
							"index":    3,
							"kind":     "article_link",
							"tag":      "a",
							"text":     "头条新闻二",
							"href":     "https://money.163.com/story/2",
							"xpath":    "/html/body/main/section[1]/ul[1]/li[2]/a[1]",
							"selector": "xpath=/html/body/main/section[1]/ul[1]/li[2]/a[1]",
						},
					},
				},
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
	flowYAML, ok := draft["flow_yaml"].(string)
	if !ok || !strings.Contains(flowYAML, `action: get_all_links`) {
		t.Fatalf("expected content-oriented flow from structured observation, got %#v", draft["flow_yaml"])
	}
}

func TestHandleDraftFlowToolIncludesRecommendedExamples(t *testing.T) {
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
				"intent":      "搜索订单并导出 CSV",
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
	recommended, ok := payload["recommended_examples"].([]any)
	if !ok || len(recommended) == 0 {
		t.Fatalf("expected recommended_examples, got %#v", payload["recommended_examples"])
	}
	first, ok := recommended[0].(map[string]any)
	if !ok || first["id"] != "search_results_to_csv" {
		t.Fatalf("expected search_results_to_csv first, got %#v", payload["recommended_examples"])
	}
	flowYAML, ok := first["flow_yaml"].(string)
	if !ok || strings.TrimSpace(flowYAML) == "" {
		t.Fatalf("expected flow_yaml in recommended example, got %#v", first["flow_yaml"])
	}
	flow, err := ParseFlow([]byte(flowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse recommended example: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate recommended example: %v", err)
	}
}

func TestHandleFinalizeFlowToolWithObservationReady(t *testing.T) {
	observation := `{
  "url": "https://example.com/orders",
  "title": "Orders",
  "artifact_root": "/tmp/artifacts",
  "page_summary": "Observed \"Orders\" with 2 interactive elements and 2 content elements. Top content: Order Center; Export orders.",
  "dom_snapshot_excerpt": "h1: Order Center\na: Export orders -> /export",
  "content_elements": [
    {
      "index": 1,
      "kind": "headline",
      "tag": "h1",
      "text": "Order Center",
      "xpath": "/html/body/main/h1[1]",
      "selector": "xpath=/html/body/main/h1[1]"
    }
  ],
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
				"intent":      `搜索订单 "A10086"`,
				"observation": observation,
			},
		},
	}

	result, err := handleFinalizeFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("finalize flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true || payload["status"] != "ready" {
		t.Fatalf("expected ready finalize payload, got %#v", payload)
	}
	if _, ok := payload["flow_yaml"].(string); !ok {
		t.Fatalf("expected flow_yaml, got %#v", payload["flow_yaml"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.run_flow" {
		t.Fatalf("expected next_action tsplay.run_flow, got %#v", payload["next_action"])
	}
	recommended, ok := payload["recommended_examples"].([]any)
	if !ok || len(recommended) == 0 {
		t.Fatalf("expected recommended_examples, got %#v", payload["recommended_examples"])
	}
	first, ok := recommended[0].(map[string]any)
	if !ok || first["id"] != "search_results_to_csv" {
		t.Fatalf("expected search_results_to_csv first, got %#v", payload["recommended_examples"])
	}
	if summary, ok := payload["page_summary"].(string); !ok || !strings.Contains(summary, "Order Center") {
		t.Fatalf("expected page_summary passthrough, got %#v", payload["page_summary"])
	}
	if excerpt, ok := payload["dom_snapshot_excerpt"].(string); !ok || !strings.Contains(excerpt, "Export orders") {
		t.Fatalf("expected dom_snapshot_excerpt passthrough, got %#v", payload["dom_snapshot_excerpt"])
	}
	content, ok := payload["content_elements"].([]any)
	if !ok || len(content) != 1 {
		t.Fatalf("expected content_elements passthrough, got %#v", payload["content_elements"])
	}
}

func TestHandleDraftFlowToolPreflightCDPErrorUsesPreflightRunSource(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"intent":               "读取页面标题",
				"url":                  "https://example.com",
				"browser_cdp_launch":   true,
				"browser_cdp_endpoint": "http://192.0.2.1:9222",
				"allow_browser_state":  true,
			},
		},
	}

	result, err := handleDraftFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("draft flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "only start or reuse a local browser") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	run, ok := payload["run"].(map[string]any)
	if !ok || run["status"] != "not_started" {
		t.Fatalf("expected synthetic run metadata, got %#v", payload["run"])
	}
	details, ok := run["details"].(map[string]any)
	if !ok || details["source"] != "preflight" {
		t.Fatalf("expected preflight synthetic run source, got %#v", run["details"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleFinalizeFlowToolWithObservationNeedsInput(t *testing.T) {
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

	result, err := handleFinalizeFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("finalize flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true || payload["status"] != "needs_input" {
		t.Fatalf("expected needs_input finalize payload, got %#v", payload)
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.finalize_flow" {
		t.Fatalf("expected next_action tsplay.finalize_flow, got %#v", payload["next_action"])
	}
}

func TestHandleFinalizeFlowToolWithObservationNeedsPermission(t *testing.T) {
	observation := `{
  "url": "https://example.com/upload",
  "title": "Upload",
  "artifact_root": "/tmp/artifacts",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "file",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["#file-input"]
    },
    {
      "index": 2,
      "tag": "button",
      "type": "button",
      "text": "Upload",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["text=\"Upload\""]
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

	result, err := handleFinalizeFlowTool(context.Background(), request)
	if err != nil {
		t.Fatalf("finalize flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true || payload["status"] != "needs_permission" {
		t.Fatalf("expected needs_permission finalize payload, got %#v", payload)
	}
	issue, ok := payload["issue"].(map[string]any)
	if !ok || issue["code"] != "security_policy" {
		t.Fatalf("expected security issue, got %#v", payload["issue"])
	}
	recommended, ok := payload["recommended_examples"].([]any)
	if !ok || len(recommended) == 0 {
		t.Fatalf("expected recommended_examples, got %#v", payload["recommended_examples"])
	}
	first, ok := recommended[0].(map[string]any)
	if !ok || first["id"] != "upload_file_then_submit" {
		t.Fatalf("expected upload_file_then_submit first, got %#v", payload["recommended_examples"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.finalize_flow" {
		t.Fatalf("expected next_action tsplay.finalize_flow, got %#v", payload["next_action"])
	}
}

func TestFinalizeFlowStatusNeedsRepair(t *testing.T) {
	draft := &FlowDraft{
		Validation: &FlowDraftValidation{
			Valid: false,
			Error: `step 2 uses unsupported action "fill"`,
			Issue: &FlowIssue{
				Code:       "unsupported_action",
				DidYouMean: "type_text",
			},
		},
	}

	status, reason := finalizeFlowStatus(draft)
	if status != "needs_repair" {
		t.Fatalf("expected needs_repair, got %q", status)
	}
	if !strings.Contains(reason, `unsupported action "fill"`) {
		t.Fatalf("unexpected reason: %q", reason)
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
	if _, ok := draft["auto_repaired"]; ok {
		t.Fatalf("expected auto_repaired to be omitted when the best selector is chosen up front, got %#v", draft["auto_repaired"])
	}
	if repairs, ok := draft["selector_repairs"]; ok && repairs != nil {
		t.Fatalf("expected selector_repairs to be omitted, got %#v", draft["selector_repairs"])
	}
	flowYAML, ok := draft["flow_yaml"].(string)
	if !ok || !strings.Contains(flowYAML, `[data-testid="order-query"]`) {
		t.Fatalf("expected stable selector in flow yaml, got %q", flowYAML)
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
	if failedStep["path"] != "2" {
		t.Fatalf("unexpected failed step path: %#v", failedStep["path"])
	}
	artifacts, ok := contextPayload["artifacts"].(map[string]any)
	if !ok {
		t.Fatalf("expected artifacts, got %#v", contextPayload["artifacts"])
	}
	if summary, ok := artifacts["artifact_summary"].([]any); !ok || len(summary) == 0 {
		t.Fatalf("expected artifact summary, got %#v", artifacts["artifact_summary"])
	}
	excerpt, ok := artifacts["dom_snapshot_excerpt"].(string)
	if !ok || !strings.Contains(excerpt, "Export orders") {
		t.Fatalf("expected dom excerpt, got %#v", artifacts["dom_snapshot_excerpt"])
	}
	relevantSelectors, ok := artifacts["relevant_selectors"].([]any)
	if !ok || len(relevantSelectors) == 0 {
		t.Fatalf("expected relevant selectors, got %#v", artifacts["relevant_selectors"])
	}
	if relevantSelectors[0] != "#export" && relevantSelectors[0] != `text="Export orders"` {
		t.Fatalf("unexpected relevant selectors: %#v", relevantSelectors)
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
	traceSummary, ok := contextPayload["trace_summary"].([]any)
	if !ok || len(traceSummary) < 2 {
		t.Fatalf("expected trace summary, got %#v", contextPayload["trace_summary"])
	}
	secondTrace, ok := traceSummary[1].(map[string]any)
	if !ok || secondTrace["label"] != "2 click export (click)" {
		t.Fatalf("unexpected trace label: %#v", traceSummary[1])
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
	if payload["tool"] != "tsplay.observe_page" {
		t.Fatalf("expected tool metadata, got %#v", payload["tool"])
	}
	if _, ok := payload["summary"].(string); !ok {
		t.Fatalf("expected summary string, got %#v", payload["summary"])
	}
	warnings, ok := payload["warnings"].([]any)
	if !ok || len(warnings) != 0 {
		t.Fatalf("expected warnings array, got %#v", payload["warnings"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.observe_page" {
		t.Fatalf("expected retry next_action, got %#v", payload["next_action"])
	}
	run, ok := payload["run"].(map[string]any)
	if !ok {
		t.Fatalf("expected run metadata, got %#v", payload["run"])
	}
	for _, key := range []string{"id", "status", "queue_wait_ms", "duration_ms", "timeout_ms", "audit_path", "run_root"} {
		if _, ok := run[key]; !ok {
			t.Fatalf("expected run metadata key %q in %#v", key, run)
		}
	}
}

func TestHandleObservePageToolIncludesObservationSummaryFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html>
<head><title>Finance</title></head>
<body>
  <main>
    <h1>财经要闻</h1>
    <p>今日市场聚焦宏观经济数据和行业财报。</p>
    <a href="/story/1">头条新闻一</a>
    <a href="/story/2">头条新闻二</a>
  </main>
</body>
</html>`)
	}))
	defer server.Close()

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url": server.URL,
			},
		},
	}

	result, err := handleObservePageToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	if summary, ok := payload["page_summary"].(string); !ok || !strings.Contains(summary, "财经要闻") {
		t.Fatalf("expected page_summary, got %#v", payload["page_summary"])
	}
	if excerpt, ok := payload["dom_snapshot_excerpt"].(string); !ok || !strings.Contains(excerpt, "头条新闻一") {
		t.Fatalf("expected dom_snapshot_excerpt, got %#v", payload["dom_snapshot_excerpt"])
	}
	content, ok := payload["content_elements"].([]any)
	if !ok || len(content) == 0 {
		t.Fatalf("expected content_elements, got %#v", payload["content_elements"])
	}
	observation, ok := payload["observation"].(map[string]any)
	if !ok {
		t.Fatalf("expected observation payload, got %#v", payload["observation"])
	}
	if observation["page_summary"] != payload["page_summary"] {
		t.Fatalf("expected observation.page_summary to match top-level summary, got %#v", observation["page_summary"])
	}
}

func TestHandleObservePageToolRestrictsCDPUserDataDirRoot(t *testing.T) {
	installCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, nil)
	defer restore()

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url":                       "https://example.com",
				"browser_cdp_launch":        true,
				"browser_cdp_user_data_dir": "../escape-profile",
				"allow_browser_state":       true,
			},
		},
	}

	result, err := handleObservePageToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "file output root") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if installCalled {
		t.Fatalf("invalid CDP user data dir should be rejected before starting Playwright")
	}
}

func TestHandleObservePageToolRejectsInvalidCDPPortBeforeRun(t *testing.T) {
	installCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, nil)
	defer restore()

	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url":                 "https://example.com",
				"browser_cdp_port":    0,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleObservePageToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "between 1 and 65535") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP port validation: %#v", payload["run"])
	}
	if installCalled {
		t.Fatalf("invalid CDP port should be rejected before starting Playwright")
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleObservePageToolRejectsRemoteCDPLaunchBeforeRun(t *testing.T) {
	installCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, nil)
	defer restore()

	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url":                  "https://example.com",
				"browser_cdp_launch":   true,
				"browser_cdp_endpoint": "http://192.0.2.1:9222",
				"allow_browser_state":  true,
			},
		},
	}

	result, err := handleObservePageToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "only start or reuse a local browser") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP launch endpoint validation: %#v", payload["run"])
	}
	if installCalled {
		t.Fatalf("invalid CDP launch endpoint should be rejected before starting Playwright")
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleObservePageToolRejectsLocalCDPLaunchEndpointWithoutPortBeforeRun(t *testing.T) {
	installCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, nil)
	defer restore()

	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url":                  "https://example.com",
				"browser_cdp_launch":   true,
				"browser_cdp_endpoint": "http://127.0.0.1",
				"allow_browser_state":  true,
			},
		},
	}

	result, err := handleObservePageToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "explicit port") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP launch endpoint validation: %#v", payload["run"])
	}
	if installCalled {
		t.Fatalf("invalid CDP launch endpoint should be rejected before starting Playwright")
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
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
	if payload["tool"] != "tsplay.validate_flow" {
		t.Fatalf("expected tool metadata, got %#v", payload["tool"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.run_flow" {
		t.Fatalf("expected next_action tsplay.run_flow, got %#v", payload["next_action"])
	}
	security, ok := payload["security"].(map[string]any)
	if !ok {
		t.Fatalf("expected security metadata, got %#v", payload["security"])
	}
	policy, ok := security["policy"].(map[string]any)
	if !ok || policy["allow_lua"] != true {
		t.Fatalf("expected allow_lua in security policy, got %#v", payload["security"])
	}
}

func TestHandleValidateFlowToolAcceptsFlowYAMLAlias(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow_yaml": `
schema_version: "1"
name: validate_from_flow_yaml
steps:
  - action: click
    selector: "#submit"
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
	if payload["valid"] != true || payload["name"] != "validate_from_flow_yaml" {
		t.Fatalf("expected flow_yaml alias to validate, got %#v", payload)
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
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.validate_flow" {
		t.Fatalf("expected validate retry next_action, got %#v", payload["next_action"])
	}
}

func TestHandleValidateFlowToolReturnsIssueForUnsupportedActionAlias(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: bad_alias
steps:
  - action: fill
    selector: "#kw"
    text: "hello"
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
	issue, ok := payload["issue"].(map[string]any)
	if !ok {
		t.Fatalf("expected issue payload, got %#v", payload["issue"])
	}
	if issue["code"] != "unsupported_action" || issue["did_you_mean"] != "type_text" {
		t.Fatalf("unexpected issue payload: %#v", issue)
	}
	repairExample, ok := payload["repair_example"].(map[string]any)
	if !ok || repairExample["id"] != "repair_fill_to_type_text" {
		t.Fatalf("expected fill repair example, got %#v", payload["repair_example"])
	}
}

func TestHandleValidateFlowToolReturnsRepairExampleForUnknownFieldAlias(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: bad_result_var
steps:
  - action: evaluate
    selector: "body"
    script: |
      return []
    result_var: rows
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
	repairExample, ok := payload["repair_example"].(map[string]any)
	if !ok || repairExample["id"] != "repair_result_var_to_save_as" {
		t.Fatalf("expected result_var repair example, got %#v", payload["repair_example"])
	}
}

func TestHandleValidateFlowToolReturnsIssueForUnexpectedParameter(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: bad_param
steps:
  - action: navigate
    url: "https://example.com"
    timeout: 3000
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
	issue, ok := payload["issue"].(map[string]any)
	if !ok {
		t.Fatalf("expected issue payload, got %#v", payload["issue"])
	}
	if issue["code"] != "unexpected_parameter" || issue["field"] != "timeout" {
		t.Fatalf("unexpected issue payload: %#v", issue)
	}
	if suggestion, _ := issue["suggestion"].(string); !strings.Contains(suggestion, "browser.timeout") {
		t.Fatalf("expected browser.timeout hint, got %#v", issue)
	}
	repairExample, ok := payload["repair_example"].(map[string]any)
	if !ok || repairExample["id"] != "repair_navigate_timeout_to_browser_timeout" {
		t.Fatalf("expected navigate timeout repair example, got %#v", payload["repair_example"])
	}
}

func TestFlowSecurityPolicyResolutionFromToolRequestSupportsPresetOverrides(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"security_preset": "full_automation",
				"allow_http":      false,
				"allow_email":     false,
			},
		},
	}

	resolution, err := flowSecurityPolicyResolutionFromToolRequest(request, options)
	if err != nil {
		t.Fatalf("resolve security preset: %v", err)
	}
	if resolution.Preset != "full_automation" {
		t.Fatalf("unexpected preset: %#v", resolution.Preset)
	}
	if resolution.Policy.AllowLua != true || resolution.Policy.AllowJavaScript != true || resolution.Policy.AllowDatabase != true {
		t.Fatalf("expected full automation grants, got %#v", resolution.Policy)
	}
	if resolution.Policy.AllowHTTP != false {
		t.Fatalf("expected explicit allow_http override to win, got %#v", resolution.Policy)
	}
	if resolution.Policy.AllowEmail != false {
		t.Fatalf("expected explicit allow_email override to win, got %#v", resolution.Policy)
	}
	if resolution.Policy.FileInputRoot != options.ArtifactRoot || resolution.Policy.FileOutputRoot != options.ArtifactRoot {
		t.Fatalf("expected artifact roots in security policy, got %#v", resolution.Policy)
	}
}

func TestHandleValidateFlowToolSupportsSecurityPreset(t *testing.T) {
	options := TSPlayMCPServerOptions{ArtifactRoot: t.TempDir()}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: validate_browser_write
steps:
  - action: screenshot
    path: capture.png
`,
				"security_preset": "browser_write",
			},
		},
	}

	result, err := handleValidateFlowToolWithOptions(context.Background(), request, options)
	if err != nil {
		t.Fatalf("validate flow with browser_write: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["valid"] != true {
		t.Fatalf("expected valid flow with browser_write preset, got %#v", payload)
	}
	security, ok := payload["security"].(map[string]any)
	if !ok || security["preset"] != "browser_write" {
		t.Fatalf("expected browser_write preset metadata, got %#v", payload["security"])
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

func TestHandleValidateFlowToolRejectsInvalidCDPEndpoint(t *testing.T) {
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: invalid_cdp_endpoint_from_mcp
browser:
  cdp_endpoint: "127.0.0.1:70000/json/version"
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
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
	if !strings.Contains(payload["error"].(string), "invalid port") {
		t.Fatalf("unexpected error: %#v", payload["error"])
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
	if payload["tool"] != "tsplay.run_flow" {
		t.Fatalf("expected tool metadata, got %#v", payload["tool"])
	}
	nextAction, ok := payload["next_action"].(map[string]any)
	if !ok || nextAction["tool"] != "tsplay.repair_flow_context" {
		t.Fatalf("expected repair_flow_context next_action, got %#v", payload["next_action"])
	}
}

func TestHandleRunFlowToolRejectsCDPWithoutBrowserStateBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: cdp_override_without_grant
steps:
  - action: navigate
    url: https://example.com
`,
				"browser_cdp_port": 9222,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "allow_browser_state") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP grant check: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsCDPPathOverridesWithoutBrowserStateBeforeRun(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantArg string
	}{
		{
			name:    "executable",
			args:    map[string]any{"browser_cdp_executable": "/definitely/missing/chrome"},
			wantArg: "browser_cdp_executable",
		},
		{
			name:    "user_data_dir",
			args:    map[string]any{"browser_cdp_user_data_dir": "profiles/cdp"},
			wantArg: "browser_cdp_user_data_dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artifactRoot := t.TempDir()
			args := map[string]any{
				"flow": `
schema_version: "1"
name: cdp_path_override_without_grant
steps:
  - action: navigate
    url: https://example.com
`,
			}
			for key, value := range tt.args {
				args[key] = value
			}
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: args,
				},
			}

			result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
			if err != nil {
				t.Fatalf("run flow: %v", err)
			}

			var payload map[string]any
			decodeToolText(t, result, &payload)
			if payload["ok"] != false {
				t.Fatalf("expected ok=false, got %#v", payload)
			}
			errorText, _ := payload["error"].(string)
			if !strings.Contains(errorText, "allow_browser_state") || !strings.Contains(errorText, tt.wantArg) {
				t.Fatalf("expected browser state grant error mentioning %q, got %#v", tt.wantArg, payload["error"])
			}
			if _, ok := payload["run"]; ok {
				t.Fatalf("run metadata should not be created before CDP grant check: %#v", payload["run"])
			}
			if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
				t.Fatalf("expected no browser run directory, got %#v", entries)
			}
		})
	}
}

func TestHandleRunFlowToolRejectsCDPOverrideWithUseSessionBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: cdp_override_use_session_conflict
browser:
  use_session: missing-admin
steps:
  - action: navigate
    url: https://example.com
`,
				"browser_cdp_port":    9222,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	errorText, _ := payload["error"].(string)
	if !strings.Contains(errorText, "use_session") || !strings.Contains(errorText, "cannot be combined") {
		t.Fatalf("expected use_session conflict, got %#v", payload["error"])
	}
	if strings.Contains(errorText, "missing-admin") {
		t.Fatalf("should reject CDP/use_session conflict before resolving missing session, got %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP use_session conflict validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsFlowCDPWithoutBrowserStateBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: flow_cdp_without_grant
browser:
  cdp_port: 9222
steps:
  - action: navigate
    url: https://example.com
`,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "allow_browser_state") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before flow CDP grant check: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsFlowCDPPortZeroBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: flow_cdp_port_zero
browser:
  cdp_port: 0
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	errorText, _ := payload["error"].(string)
	if !strings.Contains(errorText, "cdp_port") || !strings.Contains(errorText, "between 1 and 65535") {
		t.Fatalf("expected cdp_port range error, got %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before flow cdp_port validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsBlankFlowCDPEndpointBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: blank_flow_cdp_endpoint
browser:
  cdp_endpoint: ""
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	errorText, _ := payload["error"].(string)
	if !strings.Contains(errorText, "cdp_endpoint") || !strings.Contains(errorText, "cannot be blank") {
		t.Fatalf("expected cdp_endpoint blank error, got %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before flow cdp_endpoint validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsBlankFlowCDPPathFieldBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: blank_flow_cdp_executable
browser:
  cdp_executable: ""
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	errorText, _ := payload["error"].(string)
	if !strings.Contains(errorText, "cdp_executable") || !strings.Contains(errorText, "cannot be blank") {
		t.Fatalf("expected cdp_executable blank error, got %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before flow cdp_executable validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsRemoteFlowCDPLaunchBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: remote_flow_cdp_launch
browser:
  cdp_launch: true
  cdp_endpoint: "http://192.0.2.1:9222"
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "only start or reuse a local browser") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before flow CDP launch validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsInvalidCDPPortBeforeRun(t *testing.T) {
	artifactRoot := t.TempDir()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"flow": `
schema_version: "1"
name: invalid_cdp_port
steps:
  - action: navigate
    url: https://example.com
`,
				"browser_cdp_port":    0,
				"allow_browser_state": true,
			},
		},
	}

	result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}

	var payload map[string]any
	decodeToolText(t, result, &payload)
	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %#v", payload)
	}
	if !strings.Contains(payload["error"].(string), "between 1 and 65535") {
		t.Fatalf("unexpected error: %#v", payload["error"])
	}
	if _, ok := payload["run"]; ok {
		t.Fatalf("run metadata should not be created before CDP port validation: %#v", payload["run"])
	}
	if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
		t.Fatalf("expected no browser run directory, got %#v", entries)
	}
}

func TestHandleRunFlowToolRejectsInvalidCDPToolArgsBeforeRun(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr string
	}{
		{
			name:    "blank_endpoint",
			args:    map[string]any{"browser_cdp_endpoint": ""},
			wantErr: "browser_cdp_endpoint cannot be blank",
		},
		{
			name:    "whitespace_endpoint",
			args:    map[string]any{"browser_cdp_endpoint": " \t "},
			wantErr: "browser_cdp_endpoint cannot be blank",
		},
		{
			name:    "non_string_endpoint",
			args:    map[string]any{"browser_cdp_endpoint": 9222},
			wantErr: "browser_cdp_endpoint must be a string",
		},
		{
			name:    "blank_executable",
			args:    map[string]any{"browser_cdp_executable": ""},
			wantErr: "browser_cdp_executable cannot be blank",
		},
		{
			name:    "whitespace_executable",
			args:    map[string]any{"browser_cdp_executable": "\n\t"},
			wantErr: "browser_cdp_executable cannot be blank",
		},
		{
			name:    "non_string_executable",
			args:    map[string]any{"browser_cdp_executable": false},
			wantErr: "browser_cdp_executable must be a string",
		},
		{
			name:    "blank_user_data_dir",
			args:    map[string]any{"browser_cdp_user_data_dir": ""},
			wantErr: "browser_cdp_user_data_dir cannot be blank",
		},
		{
			name:    "whitespace_user_data_dir",
			args:    map[string]any{"browser_cdp_user_data_dir": " \r\n "},
			wantErr: "browser_cdp_user_data_dir cannot be blank",
		},
		{
			name:    "non_string_user_data_dir",
			args:    map[string]any{"browser_cdp_user_data_dir": 1},
			wantErr: "browser_cdp_user_data_dir must be a string",
		},
		{
			name:    "blank_launch",
			args:    map[string]any{"browser_cdp_launch": ""},
			wantErr: "browser_cdp_launch must be a boolean",
		},
		{
			name:    "invalid_launch",
			args:    map[string]any{"browser_cdp_launch": "yes"},
			wantErr: "browser_cdp_launch must be a boolean",
		},
		{
			name:    "non_boolean_launch",
			args:    map[string]any{"browser_cdp_launch": 1},
			wantErr: "browser_cdp_launch must be a boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artifactRoot := t.TempDir()
			args := map[string]any{
				"flow": `
schema_version: "1"
name: invalid_cdp_tool_arg
steps:
  - action: navigate
    url: https://example.com
`,
				"allow_browser_state": true,
			}
			for key, value := range tt.args {
				args[key] = value
			}
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: args,
				},
			}

			result, err := handleRunFlowToolWithOptions(context.Background(), request, TSPlayMCPServerOptions{ArtifactRoot: artifactRoot})
			if err != nil {
				t.Fatalf("run flow: %v", err)
			}

			var payload map[string]any
			decodeToolText(t, result, &payload)
			if payload["ok"] != false {
				t.Fatalf("expected ok=false, got %#v", payload)
			}
			if !strings.Contains(payload["error"].(string), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %#v", tt.wantErr, payload["error"])
			}
			if _, ok := payload["run"]; ok {
				t.Fatalf("run metadata should not be created before CDP tool arg validation: %#v", payload["run"])
			}
			if entries, err := os.ReadDir(filepath.Join(artifactRoot, defaultTSPlayBrowserRunFolderName)); err == nil && len(entries) > 0 {
				t.Fatalf("expected no browser run directory, got %#v", entries)
			}
		})
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
