package tsplay_core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunWorkbenchProviderPromptWithOpenAICompatible(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test-openai" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"```yaml\\nschema_version: \\\"1\\\"\\nname: repaired_by_openai\\nsteps:\\n  - action: wait_for\\n    selector: \\\"button[type='submit']\\\"\\n```\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer server.Close()

	output, view, err := RunWorkbenchProviderPrompt(WorkbenchProviderConfig{
		ProviderID: "codex_main",
		Name:       "Codex Main",
		Type:       WorkbenchProviderTypeOpenAICompatible,
		BaseURL:    server.URL,
		Model:      "gpt-4.1-mini",
		APIKey:     "sk-test-openai",
		Enabled:    true,
	}, "", "repair this flow")
	if err != nil {
		t.Fatalf("RunWorkbenchProviderPrompt(openai) error = %v", err)
	}
	if !view.Ready {
		t.Fatalf("expected ready provider view: %#v", view)
	}
	if !strings.Contains(output, "repaired_by_openai") {
		t.Fatalf("unexpected output: %q", output)
	}
	if yamlText := ExtractWorkbenchFlowYAML(output); !strings.Contains(yamlText, "schema_version") {
		t.Fatalf("expected yaml output, got %q", yamlText)
	}
}

func TestRunWorkbenchProviderPromptWithOllama(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"message\": {\n" +
			"    \"content\": \"schema_version: \\\"1\\\"\\nname: repaired_by_ollama\\nsteps:\\n  - action: wait_for\\n    selector: \\\".submit-button\\\"\\n\"\n" +
			"  }\n" +
			"}"))
	}))
	defer server.Close()

	output, view, err := RunWorkbenchProviderPrompt(WorkbenchProviderConfig{
		ProviderID: "ollama_local",
		Name:       "Ollama Local",
		Type:       WorkbenchProviderTypeOllama,
		BaseURL:    server.URL,
		Model:      "qwen2.5-coder:7b",
		Enabled:    true,
	}, "", "repair this flow")
	if err != nil {
		t.Fatalf("RunWorkbenchProviderPrompt(ollama) error = %v", err)
	}
	if !view.Ready {
		t.Fatalf("expected ready provider view: %#v", view)
	}
	if !strings.Contains(output, "repaired_by_ollama") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRunWorkbenchProviderPromptWithOllamaCompatibleChoicesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\n" +
			"  \"choices\": [\n" +
			"    {\n" +
			"      \"message\": {\n" +
			"        \"content\": \"schema_version: \\\"1\\\"\\nname: repaired_from_choices\\nsteps:\\n  - action: wait_for_selector\\n    selector: \\\".submit-button\\\"\\n\"\n" +
			"      }\n" +
			"    }\n" +
			"  ]\n" +
			"}"))
	}))
	defer server.Close()

	output, view, err := RunWorkbenchProviderPrompt(WorkbenchProviderConfig{
		ProviderID: "ollama_gateway",
		Name:       "Ollama Gateway",
		Type:       WorkbenchProviderTypeOllama,
		BaseURL:    server.URL,
		Model:      "deepseek",
		Enabled:    true,
	}, "", "repair this flow")
	if err != nil {
		t.Fatalf("RunWorkbenchProviderPrompt(ollama choices) error = %v", err)
	}
	if !view.Ready {
		t.Fatalf("expected ready provider view: %#v", view)
	}
	if !strings.Contains(output, "repaired_from_choices") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRunWorkbenchProviderPromptWithOllamaGenerateFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/chat":
			_, _ = w.Write([]byte("{\n" +
				"  \"message\": {\n" +
				"    \"role\": \"assistant\",\n" +
				"    \"content\": \"\"\n" +
				"  },\n" +
				"  \"done\": true\n" +
				"}"))
		case "/api/generate":
			_, _ = w.Write([]byte("{\n" +
				"  \"data\": {\n" +
				"    \"response\": \"schema_version: \\\"1\\\"\\nname: repaired_from_generate\\nsteps:\\n  - action: wait_for_selector\\n    selector: \\\".submit-button\\\"\\n\"\n" +
				"  },\n" +
				"  \"status\": \"success\"\n" +
				"}"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	output, view, err := RunWorkbenchProviderPrompt(WorkbenchProviderConfig{
		ProviderID: "ollama_generate",
		Name:       "Ollama Generate",
		Type:       WorkbenchProviderTypeOllama,
		BaseURL:    server.URL,
		Model:      "deepseek",
		Enabled:    true,
	}, "", "repair this flow")
	if err != nil {
		t.Fatalf("RunWorkbenchProviderPrompt(ollama generate fallback) error = %v", err)
	}
	if !view.Ready {
		t.Fatalf("expected ready provider view: %#v", view)
	}
	if !strings.Contains(output, "repaired_from_generate") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestBuildWorkbenchProviderViewIncludesAutoProvider(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "sk-env-12345678")
	t.Setenv("OPENAI_MODEL", "gpt-4.1-mini")

	auto := detectWorkbenchAutoProviderConfig()
	if auto == nil {
		t.Fatalf("expected auto provider")
	}
	view := BuildWorkbenchProviderView(*auto)
	if view.ProviderID != workbenchAutoProviderID {
		t.Fatalf("unexpected provider id: %q", view.ProviderID)
	}
	if !view.Ready {
		t.Fatalf("expected auto provider to be ready: %#v", view)
	}
	if !view.Detected {
		t.Fatalf("expected detected=true: %#v", view)
	}
}
