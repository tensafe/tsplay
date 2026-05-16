package tsplay_core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestWorkbenchTaskRunHandlerExposesManualReviewArtifacts(t *testing.T) {
	artifactRoot := t.TempDir()
	handler := NewWorkbenchAPIHandler(artifactRoot)
	flowYAML := `
schema_version: "1"
name: workbench_manual_review_run
steps:
  - action: write_json
    file_path: artifacts/manual-review/evidence.json
    with:
      value:
        score: 0.42
  - action: set_var
    save_as: payload
    with:
      value:
        status: manual_review
        action: manual_review_required
        phase: detect_score
        reason: low detection score
        evidence:
          result_path: artifacts/manual-review/evidence.json
`
	bodyBytes, err := json.Marshal(map[string]any{"flow_yaml": flowYAML})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/workbench/tasks/run", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v\n%s", err, rec.Body.String())
	}
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
	if payload["status"] != FlowRunStatusManualReviewRequired {
		t.Fatalf("status = %#v", payload["status"])
	}
	if payload["requires_manual_review"] != true {
		t.Fatalf("requires_manual_review = %#v", payload["requires_manual_review"])
	}
	review, ok := payload["manual_review"].(map[string]any)
	if !ok {
		t.Fatalf("manual_review = %#v", payload["manual_review"])
	}
	if review["phase"] != "detect_score" {
		t.Fatalf("manual review phase = %#v", review["phase"])
	}
	artifacts, ok := review["artifacts"].([]any)
	if !ok || len(artifacts) != 1 {
		t.Fatalf("manual review artifacts = %#v", review["artifacts"])
	}
	artifact, ok := artifacts[0].(map[string]any)
	if !ok {
		t.Fatalf("artifact = %#v", artifacts[0])
	}
	if artifact["url_path"] != "/workbench-artifacts/artifacts/manual-review/evidence.json" {
		t.Fatalf("artifact url_path = %#v", artifact["url_path"])
	}
	if artifact["exists"] != true {
		t.Fatalf("artifact exists = %#v", artifact["exists"])
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

func TestWorkbenchTaskPlanHandlerUsesProviderForFlowGeneration(t *testing.T) {
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected provider path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"schema_version: \\\"1\\\"\\nname: ai_search_orders\\nsteps:\\n  - action: navigate\\n    url: \\\"https://example.com/admin/orders\\\"\\n  - action: type_text\\n    selector: \\\"[data-testid='order-query']\\\"\\n    text: \\\"{{order_query}}\\\"\\n  - action: click\\n    selector: \\\"text=\\\\\\\"Search\\\\\\\"\\\"\\n  - action: click\\n    selector: \\\"[data-cy='export-link']\\\"\\n\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer providerServer.Close()

	artifactRoot := t.TempDir()
	observationPath := filepath.Join(artifactRoot, "orders-observation.json")
	if err := os.WriteFile(observationPath, []byte(`{
  "url": "https://example.com/admin/orders",
  "title": "订单管理",
  "page_summary": "订单搜索和导出页面",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "text",
      "label": "订单关键词",
      "placeholder": "搜索订单",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["[data-testid='order-query']"]
    },
    {
      "index": 2,
      "tag": "button",
      "type": "button",
      "text": "Search",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["text=\"Search\""]
    },
    {
      "index": 3,
      "tag": "a",
      "type": "link",
      "text": "Export orders",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["[data-cy='export-link']"]
    }
  ]
}`), 0600); err != nil {
		t.Fatalf("write observation: %v", err)
	}

	if _, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:      "demo_admin",
		Name:        "Demo Admin",
		StartURL:    "https://example.com/admin",
		SessionName: "demo_admin_user",
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
	}
	if _, err := UpsertWorkbenchPageCards("demo_admin", artifactRoot, []WorkbenchPageCard{
		{
			ID:              "route:demo_admin:/orders",
			SiteID:          "demo_admin",
			URL:             "https://example.com/admin/orders",
			NormalizedRoute: "/orders",
			Title:           "订单管理",
			Summary:         "支持订单搜索和导出",
			ObservationPath: observationPath,
		},
	}); err != nil {
		t.Fatalf("UpsertWorkbenchPageCards() error = %v", err)
	}
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
		"/api/workbench/tasks/plan",
		strings.NewReader(`{
  "site_id": "demo_admin",
  "provider_id": "codex_main",
  "intent": "搜索订单并导出"
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"generation_mode": "provider"`) {
		t.Fatalf("expected provider generation, body=%s", body)
	}
	if !strings.Contains(body, `"provider_id": "codex_main"`) {
		t.Fatalf("expected provider metadata, body=%s", body)
	}
	if !strings.Contains(body, `ai_search_orders`) {
		t.Fatalf("expected provider-generated flow name, body=%s", body)
	}
	if !strings.Contains(body, `"use_session": "demo_admin_user"`) {
		t.Fatalf("expected generated flow to reuse session, body=%s", body)
	}
}

func TestWorkbenchTaskPlanHandlerFallsBackToLocalDraftWhenProviderFlowIsInvalid(t *testing.T) {
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"schema_version: \\\"1\\\"\\nname: invalid_ai_flow\\nsteps:\\n  - action: wait_for\\n    selector: \\\"#submit\\\"\\n\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer providerServer.Close()

	artifactRoot := t.TempDir()
	observationPath := filepath.Join(artifactRoot, "orders-observation.json")
	if err := os.WriteFile(observationPath, []byte(`{
  "url": "https://example.com/admin/orders",
  "title": "订单管理",
  "page_summary": "订单搜索和导出页面",
  "elements": [
    {
      "index": 1,
      "tag": "input",
      "type": "text",
      "label": "订单关键词",
      "placeholder": "搜索订单",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["[data-testid='order-query']"]
    },
    {
      "index": 2,
      "tag": "button",
      "type": "button",
      "text": "Search",
      "visible": true,
      "enabled": true,
      "selector_candidates": ["text=\"Search\""]
    }
  ]
}`), 0600); err != nil {
		t.Fatalf("write observation: %v", err)
	}

	if _, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:      "demo_admin",
		Name:        "Demo Admin",
		StartURL:    "https://example.com/admin",
		SessionName: "demo_admin_user",
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
	}
	if _, err := UpsertWorkbenchPageCards("demo_admin", artifactRoot, []WorkbenchPageCard{
		{
			ID:              "route:demo_admin:/orders",
			SiteID:          "demo_admin",
			URL:             "https://example.com/admin/orders",
			NormalizedRoute: "/orders",
			Title:           "订单管理",
			Summary:         "支持订单搜索",
			ObservationPath: observationPath,
		},
	}); err != nil {
		t.Fatalf("UpsertWorkbenchPageCards() error = %v", err)
	}
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
		"/api/workbench/tasks/plan",
		strings.NewReader(`{
  "site_id": "demo_admin",
  "provider_id": "codex_main",
  "intent": "搜索订单"
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"generation_mode": "local_fallback"`) {
		t.Fatalf("expected local fallback generation mode, body=%s", body)
	}
	if !strings.Contains(body, `AI-generated flow did not pass validation`) {
		t.Fatalf("expected validation fallback warning, body=%s", body)
	}
	if !strings.Contains(body, `type_text`) {
		t.Fatalf("expected local selector-aware draft to remain available, body=%s", body)
	}
	if strings.Contains(body, `"name": "invalid_ai_flow"`) {
		t.Fatalf("did not expect invalid provider flow to replace local draft, body=%s", body)
	}
}

func TestWorkbenchTaskPlanHandlerNormalizesProviderBooleanUseSessionWithoutSavedSession(t *testing.T) {
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"schema_version: \\\"1\\\"\\nname: ai_news_top10\\nbrowser:\\n  use_session: true\\nsteps:\\n  - action: navigate\\n    url: \\\"https://www.163.com/\\\"\\n  - action: capture_table\\n    selector: \\\"table\\\"\\n    save_as: rows\\n\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer providerServer.Close()

	artifactRoot := t.TempDir()
	observationPath := filepath.Join(artifactRoot, "news-observation.json")
	if err := os.WriteFile(observationPath, []byte(`{
  "url": "https://www.163.com/",
  "title": "网易",
  "page_summary": "首页新闻与内容导航",
  "elements": [
    {
      "index": 1,
      "tag": "table",
      "selector_candidates": ["table"]
    }
  ]
}`), 0600); err != nil {
		t.Fatalf("write observation: %v", err)
	}

	if _, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:   "news_demo",
		Name:     "News Demo",
		StartURL: "https://www.163.com/",
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
	}
	if _, err := UpsertWorkbenchPageCards("news_demo", artifactRoot, []WorkbenchPageCard{
		{
			ID:              "route:news_demo:/",
			SiteID:          "news_demo",
			URL:             "https://www.163.com/",
			NormalizedRoute: "/",
			Title:           "网易",
			Summary:         "首页新闻内容",
			ObservationPath: observationPath,
		},
	}); err != nil {
		t.Fatalf("UpsertWorkbenchPageCards() error = %v", err)
	}
	if _, err := SaveWorkbenchProviderConfig(WorkbenchProviderConfig{
		ProviderID: "asdf",
		Name:       "asdf",
		Type:       WorkbenchProviderTypeOpenAICompatible,
		BaseURL:    providerServer.URL,
		Model:      "gpt-test",
		APIKey:     "sk-test-openai",
		Enabled:    true,
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchProviderConfig() error = %v", err)
	}

	handler := NewWorkbenchAPIHandler(artifactRoot)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/workbench/tasks/plan",
		strings.NewReader(`{
  "site_id": "news_demo",
  "provider_id": "asdf",
  "intent": "帮我获取前10条新闻"
}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"generation_mode": "provider"`) {
		t.Fatalf("expected provider generation, body=%s", body)
	}
	if strings.Contains(body, `"generation_mode": "local_fallback"`) {
		t.Fatalf("did not expect local fallback, body=%s", body)
	}
	if strings.Contains(body, `"validation_error"`) {
		t.Fatalf("did not expect validation error, body=%s", body)
	}
	if strings.Contains(body, `"use_session": "true"`) || strings.Contains(body, `"use_session": true`) {
		t.Fatalf("did not expect boolean or string true use_session to survive, body=%s", body)
	}
	if !strings.Contains(body, `browser.use_session=true 已自动移除`) {
		t.Fatalf("expected normalization warning, body=%s", body)
	}
}

func TestWorkbenchTaskPlanHandlerPassesRealtimeContextToProvider(t *testing.T) {
	var capturedRequestBody string
	providerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read provider request body: %v", err)
		}
		capturedRequestBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"schema_version: \\\"1\\\"\\nname: ai_live_news\\nsteps:\\n  - action: navigate\\n    url: \\\"https://www.163.com/\\\"\\n  - action: wait_for_selector\\n    selector: \\\"main a\\\"\\n    timeout: 10000\\n  - action: extract_text\\n    selector: \\\"main a\\\"\\n    save_as: article_title\\n\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer providerServer.Close()

	artifactRoot := t.TempDir()
	if _, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:   "live_news",
		Name:     "Live News",
		StartURL: "https://www.163.com/",
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
	}
	if _, err := SaveWorkbenchProviderConfig(WorkbenchProviderConfig{
		ProviderID: "live_provider",
		Name:       "live_provider",
		Type:       WorkbenchProviderTypeOpenAICompatible,
		BaseURL:    providerServer.URL,
		Model:      "gpt-test",
		APIKey:     "sk-test-openai",
		Enabled:    true,
	}, artifactRoot); err != nil {
		t.Fatalf("SaveWorkbenchProviderConfig() error = %v", err)
	}

	handler := NewWorkbenchAPIHandler(artifactRoot)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/workbench/tasks/plan",
		strings.NewReader(`{
  "site_id": "live_news",
  "provider_id": "live_provider",
  "intent": "帮我获取前10条最新的新闻",
  "realtime_context": {
    "url": "https://www.163.com/",
    "title": "网易新闻",
    "html": "<html><body><main><h1>网易新闻</h1><a href=\"/news/1\">头条一</a></main></body></html>",
    "observation": {
      "url": "https://www.163.com/",
      "title": "网易新闻",
      "page_summary": "首页有最新新闻列表",
      "content_elements": [
        {
          "index": 1,
          "kind": "headline",
          "text": "头条一",
          "selector": "main a"
        }
      ]
    }
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
	if !strings.Contains(body, `"generation_mode": "provider"`) {
		t.Fatalf("expected provider generation, body=%s", body)
	}
	if !strings.Contains(body, `ai_live_news`) {
		t.Fatalf("expected provider flow name, body=%s", body)
	}
	if !strings.Contains(capturedRequestBody, `realtime_context`) {
		t.Fatalf("expected provider request to include realtime_context, body=%s", capturedRequestBody)
	}
	if !strings.Contains(capturedRequestBody, `头条一`) {
		t.Fatalf("expected provider request to include live observation content, body=%s", capturedRequestBody)
	}
	if !strings.Contains(capturedRequestBody, `网易新闻`) {
		t.Fatalf("expected provider request to include live title/html excerpt, body=%s", capturedRequestBody)
	}
}
