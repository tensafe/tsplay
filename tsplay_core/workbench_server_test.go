package tsplay_core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWorkbenchTaskRunHandlerRunsInlineFlow(t *testing.T) {
	handler := NewWorkbenchAPIHandler(t.TempDir())

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/workbench/tasks/run",
		strings.NewReader(`{
  "flow_yaml": "schema_version: \"1\"\nname: workbench_inline_run\nsteps:\n  - action: set_var\n    save_as: answer\n    value: 42\n"
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"ok": true`) {
		t.Fatalf("expected ok=true, body=%s", body)
	}
	if !strings.Contains(body, `"flow_name": "workbench_inline_run"`) {
		t.Fatalf("expected flow_name, body=%s", body)
	}
	if !strings.Contains(body, `"answer": "42"`) {
		t.Fatalf("expected vars.answer=42, body=%s", body)
	}
}

func TestWorkbenchTaskRepairHandlerBuildsContext(t *testing.T) {
	handler := NewWorkbenchAPIHandler(t.TempDir())

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/workbench/tasks/repair",
		strings.NewReader(`{
  "flow_yaml": "schema_version: \"1\"\nname: workbench_failed_run\nsteps:\n  - action: click\n    selector: \"#submit\"\n",
  "error": "selector not found",
  "run_result": {
    "name": "workbench_failed_run",
    "trace": [
      {
        "index": 1,
        "path": "1",
        "action": "click",
        "status": "error",
        "args_summary": "{\"selector\":\"#submit\"}",
        "error": "selector not found",
        "started_at": "2026-04-26T14:00:00+08:00",
        "finished_at": "2026-04-26T14:00:00+08:00"
      }
    ]
  }
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"ok": true`) {
		t.Fatalf("expected ok=true, body=%s", body)
	}
	if !strings.Contains(body, `"failed_step_path": "1"`) {
		t.Fatalf("expected failed_step_path, body=%s", body)
	}
	if !strings.Contains(body, `"repair_hints"`) {
		t.Fatalf("expected repair_hints, body=%s", body)
	}
	if !strings.Contains(body, `"prompt"`) {
		t.Fatalf("expected repair prompt, body=%s", body)
	}
}

func TestWorkbenchTaskRepairAutoHandlerUsesProvider(t *testing.T) {
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected provider path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"```yaml\\nschema_version: \\\"1\\\"\\nname: repaired_auto_run\\nsteps:\\n  - action: wait_for\\n    selector: \\\"button[type='submit']\\\"\\n```\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer providerServer.Close()

	artifactRoot := t.TempDir()
	if _, err := SaveWorkbenchProviderConfig(WorkbenchProviderConfig{
		ProviderID: "codex_main",
		Name:       "Codex Main",
		Type:       WorkbenchProviderTypeOpenAICompatible,
		BaseURL:    providerServer.URL,
		Model:      "gpt-4.1-mini",
		APIKey:     "sk-test-openai",
		Enabled:    true,
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchProviderConfig() error = %v", err)
	}

	handler := NewWorkbenchAPIHandler(artifactRoot)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/workbench/tasks/repair/auto",
		strings.NewReader(`{
  "provider_id": "codex_main",
  "flow_yaml": "schema_version: \"1\"\nname: workbench_failed_run\nsteps:\n  - action: click\n    selector: \"#submit\"\n",
  "error": "selector not found",
  "run_result": {
    "name": "workbench_failed_run",
    "trace": [
      {
        "index": 1,
        "path": "1",
        "action": "click",
        "status": "error",
        "args_summary": "{\"selector\":\"#submit\"}",
        "error": "selector not found",
        "started_at": "2026-04-26T14:00:00+08:00",
        "finished_at": "2026-04-26T14:00:00+08:00"
      }
    ]
  }
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"ok": true`) {
		t.Fatalf("expected ok=true, body=%s", body)
	}
	if !strings.Contains(body, `"provider_id": "codex_main"`) {
		t.Fatalf("expected provider_id in body=%s", body)
	}
	if !strings.Contains(body, `repaired_auto_run`) {
		t.Fatalf("expected repaired flow in body=%s", body)
	}
}
