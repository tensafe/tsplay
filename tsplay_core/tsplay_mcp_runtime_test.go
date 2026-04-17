package tsplay_core

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type runtimeTestSession struct {
	id          string
	clientInfo  mcp.Implementation
	initialized bool
	notify      chan mcp.JSONRPCNotification
}

func (session *runtimeTestSession) Initialize() {
	session.initialized = true
}

func (session *runtimeTestSession) Initialized() bool {
	return session.initialized
}

func (session *runtimeTestSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return session.notify
}

func (session *runtimeTestSession) SessionID() string {
	return session.id
}

func (session *runtimeTestSession) GetClientInfo() mcp.Implementation {
	return session.clientInfo
}

func (session *runtimeTestSession) SetClientInfo(info mcp.Implementation) {
	session.clientInfo = info
}

func TestSaveFlowSavedSessionTracksOwnerAndBlocksCrossSessionUse(t *testing.T) {
	artifactRoot := t.TempDir()
	owner := FlowSavedSessionAccessInfo{
		SessionID:     "session-owner",
		ClientName:    "codex",
		ClientVersion: "1.0.0",
	}

	session, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:               "admin",
		ArtifactRoot:       artifactRoot,
		StorageStateJSON:   `{"cookies":[],"origins":[]}`,
		OwnerSessionID:     owner.SessionID,
		OwnerClientName:    owner.ClientName,
		OwnerClientVersion: owner.ClientVersion,
	})
	if err != nil {
		t.Fatalf("save session: %v", err)
	}
	if session.OwnerSessionID != owner.SessionID {
		t.Fatalf("owner_session_id = %q", session.OwnerSessionID)
	}

	_, err = SaveFlowSavedSession(FlowSavedSessionSaveOptions{
		Name:             "admin",
		ArtifactRoot:     artifactRoot,
		StorageStateJSON: `{"cookies":[{"name":"SESSION"}],"origins":[]}`,
		OwnerSessionID:   "session-other",
	})
	if err == nil || !strings.Contains(err.Error(), "owned by MCP session") {
		t.Fatalf("expected ownership error, got %v", err)
	}

	_, err = ResolveFlowSavedSessionBrowserConfig("admin", artifactRoot, FlowSavedSessionAccessInfo{SessionID: "session-other"})
	if err == nil || !strings.Contains(err.Error(), "owned by MCP session") {
		t.Fatalf("expected resolve ownership error, got %v", err)
	}

	used, err := MarkFlowSavedSessionUsed("admin", artifactRoot, FlowSavedSessionAccessInfo{
		SessionID: owner.SessionID,
		RunID:     "run-123",
	})
	if err != nil {
		t.Fatalf("mark used: %v", err)
	}
	if used.LastUsedBySessionID != owner.SessionID {
		t.Fatalf("last_used_by_session_id = %q", used.LastUsedBySessionID)
	}
	if used.LastUsedByRunID != "run-123" {
		t.Fatalf("last_used_by_run_id = %q", used.LastUsedByRunID)
	}

	_, err = DeleteFlowSavedSession("admin", artifactRoot, FlowSavedSessionAccessInfo{SessionID: "session-other"})
	if err == nil || !strings.Contains(err.Error(), "owned by MCP session") {
		t.Fatalf("expected delete ownership error, got %v", err)
	}

	view := BuildFlowSavedSessionView(*used, artifactRoot)
	ownerView, ok := view["owner"].(map[string]any)
	if !ok || ownerView["session_id"] != owner.SessionID {
		t.Fatalf("owner view = %#v", view["owner"])
	}
	lastUsedBy, ok := view["last_used_by"].(map[string]any)
	if !ok || lastUsedBy["run_id"] != "run-123" {
		t.Fatalf("last_used_by view = %#v", view["last_used_by"])
	}
}

func TestBeginTSPlayBrowserRunWritesAuditAndTracksCaller(t *testing.T) {
	artifactRoot := t.TempDir()
	session := &runtimeTestSession{
		id:          "session-a",
		initialized: true,
		notify:      make(chan mcp.JSONRPCNotification, 1),
		clientInfo: mcp.Implementation{
			Name:    "codex-client",
			Version: "2026.04",
		},
	}
	ctx := server.NewMCPServer("test", "1.0.0").WithContext(context.Background(), session)
	options := normalizeTSPlayMCPServerOptions([]TSPlayMCPServerOptions{{
		ArtifactRoot:        artifactRoot,
		DefaultRunTimeoutMS: 2000,
		QueueTimeoutMS:      200,
	}})
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"url": "https://example.com",
			},
		},
	}

	handle, runCtx, err := beginTSPlayBrowserRun(ctx, request, "tsplay.observe_page", options, nil)
	if err != nil {
		t.Fatalf("begin run: %v", err)
	}
	if runCtx == nil {
		t.Fatalf("expected run context")
	}

	run := handle.finish(nil, map[string]any{
		"url": "https://example.com",
	})
	if run.Caller.SessionID != "session-a" {
		t.Fatalf("caller session = %#v", run.Caller)
	}
	if run.Caller.ClientName != "codex-client" {
		t.Fatalf("caller client = %#v", run.Caller)
	}
	if run.Status != "ok" {
		t.Fatalf("run status = %q", run.Status)
	}
	if run.AuditPath == "" {
		t.Fatalf("expected audit path")
	}

	content, err := os.ReadFile(run.AuditPath)
	if err != nil {
		t.Fatalf("read audit: %v", err)
	}
	var audit struct {
		Run TSPlayBrowserRun `json:"run"`
	}
	if err := json.Unmarshal(content, &audit); err != nil {
		t.Fatalf("decode audit: %v", err)
	}
	if audit.Run.ID != run.ID {
		t.Fatalf("audit run id = %q, want %q", audit.Run.ID, run.ID)
	}
	if audit.Run.Details["url"] != "https://example.com" {
		t.Fatalf("audit details = %#v", audit.Run.Details)
	}
}

func TestBeginTSPlayBrowserRunHonorsSessionConcurrencyLimit(t *testing.T) {
	artifactRoot := t.TempDir()
	session := &runtimeTestSession{
		id:          "session-b",
		initialized: true,
		notify:      make(chan mcp.JSONRPCNotification, 1),
	}
	ctx := server.NewMCPServer("test", "1.0.0").WithContext(context.Background(), session)
	options := normalizeTSPlayMCPServerOptions([]TSPlayMCPServerOptions{{
		ArtifactRoot:                       artifactRoot,
		DefaultRunTimeoutMS:                1000,
		QueueTimeoutMS:                     20,
		MaxConcurrentBrowserRuns:           1,
		MaxConcurrentBrowserRunsPerSession: 1,
	}})

	first, _, err := beginTSPlayBrowserRun(ctx, mcp.CallToolRequest{}, "tsplay.observe_page", options, nil)
	if err != nil {
		t.Fatalf("begin first run: %v", err)
	}
	defer first.finish(nil, nil)

	second, _, err := beginTSPlayBrowserRun(ctx, mcp.CallToolRequest{}, "tsplay.observe_page", options, nil)
	if err == nil || !strings.Contains(err.Error(), "concurrency queue timed out") {
		t.Fatalf("expected queue timeout, got %v", err)
	}
	if second == nil {
		t.Fatalf("expected failed run handle")
	}
	if second.snapshot().Status != "timed_out" {
		t.Fatalf("second run status = %#v", second.snapshot())
	}
}
