package tsplay_core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	WorkbenchProviderTypeCodexOpenAI      = "codex_openai"
	WorkbenchProviderTypeOpenAICompatible = "openai_compatible"
	WorkbenchProviderTypeOllama           = "ollama"

	workbenchAutoProviderID             = "codex_auto"
	defaultWorkbenchOpenAIBaseURL       = "https://api.openai.com/v1"
	defaultWorkbenchOllamaBaseURL       = "http://127.0.0.1:11434"
	defaultWorkbenchOpenAIModel         = "gpt-4.1-mini"
	defaultWorkbenchRepairSystemPrompt  = "You repair TSPlay Flow YAML. Return only valid corrected YAML and do not add prose."
	defaultWorkbenchProviderHTTPTimeout = 90 * time.Second
	defaultWorkbenchProviderTemperature = 0.1
	defaultWorkbenchProviderMaxTokens   = 4000
)

var workbenchYAMLFencePattern = regexp.MustCompile("(?s)```(?:yaml|yml)?\\s*(.*?)```")

type workbenchResolvedProvider struct {
	Config       WorkbenchProviderConfig
	BaseURL      string
	Model        string
	APIKey       string
	APIKeySource string
	SystemPrompt string
}

type workbenchProviderMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func normalizeWorkbenchProviderType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case WorkbenchProviderTypeCodexOpenAI:
		return WorkbenchProviderTypeCodexOpenAI
	case WorkbenchProviderTypeOpenAICompatible:
		return WorkbenchProviderTypeOpenAICompatible
	case WorkbenchProviderTypeOllama:
		return WorkbenchProviderTypeOllama
	default:
		return ""
	}
}

func normalizeWorkbenchProviderBaseURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "://") {
		value = "http://" + value
	}
	return strings.TrimRight(value, "/")
}

func defaultWorkbenchProviderBaseURL(providerType string) string {
	switch normalizeWorkbenchProviderType(providerType) {
	case WorkbenchProviderTypeCodexOpenAI, WorkbenchProviderTypeOpenAICompatible:
		return normalizeWorkbenchProviderBaseURL(firstNonEmpty(os.Getenv("OPENAI_BASE_URL"), defaultWorkbenchOpenAIBaseURL))
	case WorkbenchProviderTypeOllama:
		return normalizeWorkbenchProviderBaseURL(firstNonEmpty(os.Getenv("OLLAMA_HOST"), defaultWorkbenchOllamaBaseURL))
	default:
		return ""
	}
}

func defaultWorkbenchProviderAPIKeyEnv(providerType string) string {
	switch normalizeWorkbenchProviderType(providerType) {
	case WorkbenchProviderTypeCodexOpenAI, WorkbenchProviderTypeOpenAICompatible:
		return "OPENAI_API_KEY"
	default:
		return ""
	}
}

func defaultWorkbenchProviderModel(providerType string) string {
	switch normalizeWorkbenchProviderType(providerType) {
	case WorkbenchProviderTypeCodexOpenAI, WorkbenchProviderTypeOpenAICompatible:
		return firstNonEmpty(strings.TrimSpace(os.Getenv("OPENAI_MODEL")), defaultWorkbenchOpenAIModel)
	case WorkbenchProviderTypeOllama:
		return strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
	default:
		return ""
	}
}

func detectWorkbenchAutoProviderConfig() *WorkbenchProviderConfig {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return nil
	}
	now := time.Now().Format(time.RFC3339Nano)
	return &WorkbenchProviderConfig{
		ProviderID: workbenchAutoProviderID,
		Name:       "Codex OpenAI (Auto)",
		Type:       WorkbenchProviderTypeCodexOpenAI,
		BaseURL:    normalizeWorkbenchProviderBaseURL(os.Getenv("OPENAI_BASE_URL")),
		Model:      strings.TrimSpace(os.Getenv("OPENAI_MODEL")),
		APIKeyEnv:  "OPENAI_API_KEY",
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func resolveWorkbenchProviderConfig(config WorkbenchProviderConfig) (*workbenchResolvedProvider, error) {
	providerType := normalizeWorkbenchProviderType(config.Type)
	if providerType == "" {
		return nil, fmt.Errorf("provider type %q is not supported", config.Type)
	}
	if !config.Enabled {
		return nil, fmt.Errorf("provider %q is disabled", config.ProviderID)
	}

	baseURL := normalizeWorkbenchProviderBaseURL(config.BaseURL)
	if baseURL == "" {
		baseURL = defaultWorkbenchProviderBaseURL(providerType)
	}
	if baseURL == "" {
		return nil, fmt.Errorf("provider %q does not have a base_url", config.ProviderID)
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("provider %q base_url %q is invalid", config.ProviderID, baseURL)
	}

	model := strings.TrimSpace(config.Model)
	if model == "" {
		model = defaultWorkbenchProviderModel(providerType)
	}
	if model == "" {
		return nil, fmt.Errorf("provider %q does not have a model", config.ProviderID)
	}

	apiKey := strings.TrimSpace(config.APIKey)
	apiKeySource := ""
	if apiKey != "" {
		apiKeySource = "config"
	}
	if apiKey == "" {
		for _, envName := range []string{strings.TrimSpace(config.APIKeyEnv), defaultWorkbenchProviderAPIKeyEnv(providerType)} {
			if envName == "" {
				continue
			}
			if value := strings.TrimSpace(os.Getenv(envName)); value != "" {
				apiKey = value
				apiKeySource = "env:" + envName
				break
			}
		}
	}
	if apiKey == "" && providerType != WorkbenchProviderTypeOllama {
		return nil, fmt.Errorf("provider %q requires an API key", config.ProviderID)
	}

	systemPrompt := strings.TrimSpace(config.SystemPrompt)
	if systemPrompt == "" {
		systemPrompt = defaultWorkbenchRepairSystemPrompt
	}

	resolved := &workbenchResolvedProvider{
		Config: WorkbenchProviderConfig{
			ProviderID:   normalizeWorkbenchSiteID(config.ProviderID),
			Name:         firstNonEmpty(strings.TrimSpace(config.Name), normalizeWorkbenchSiteID(config.ProviderID)),
			Type:         providerType,
			BaseURL:      baseURL,
			Model:        model,
			APIKeyEnv:    strings.TrimSpace(config.APIKeyEnv),
			SystemPrompt: systemPrompt,
			Enabled:      true,
			CreatedAt:    config.CreatedAt,
			UpdatedAt:    config.UpdatedAt,
		},
		BaseURL:      baseURL,
		Model:        model,
		APIKey:       apiKey,
		APIKeySource: apiKeySource,
		SystemPrompt: systemPrompt,
	}
	return resolved, nil
}

func BuildWorkbenchProviderView(config WorkbenchProviderConfig) WorkbenchProviderView {
	view := WorkbenchProviderView{
		ProviderID:   normalizeWorkbenchSiteID(config.ProviderID),
		Name:         strings.TrimSpace(config.Name),
		Type:         normalizeWorkbenchProviderType(config.Type),
		BaseURL:      normalizeWorkbenchProviderBaseURL(config.BaseURL),
		Model:        strings.TrimSpace(config.Model),
		APIKeyEnv:    strings.TrimSpace(config.APIKeyEnv),
		SystemPrompt: strings.TrimSpace(config.SystemPrompt),
		Enabled:      config.Enabled,
		CreatedAt:    config.CreatedAt,
		UpdatedAt:    config.UpdatedAt,
		Detected:     normalizeWorkbenchSiteID(config.ProviderID) == workbenchAutoProviderID,
	}
	if view.Name == "" {
		view.Name = view.ProviderID
	}
	if apiKey := strings.TrimSpace(config.APIKey); apiKey != "" {
		view.HasAPIKey = true
		view.APIKeyMasked = maskWorkbenchSecret(apiKey)
	}
	if resolved, err := resolveWorkbenchProviderConfig(config); err == nil {
		view.ResolvedBaseURL = resolved.BaseURL
		view.ResolvedModel = resolved.Model
		view.ResolvedAPIKeySource = resolved.APIKeySource
		view.Ready = true
		view.Status = "ready"
		if view.APIKeyMasked == "" && resolved.APIKey != "" {
			view.HasAPIKey = true
			view.APIKeyMasked = maskWorkbenchSecret(resolved.APIKey)
		}
		return view
	} else {
		view.Error = err.Error()
		view.Status = "needs_config"
		if view.HasAPIKey {
			return view
		}
		for _, envName := range []string{strings.TrimSpace(config.APIKeyEnv), defaultWorkbenchProviderAPIKeyEnv(config.Type)} {
			if envName == "" {
				continue
			}
			if value := strings.TrimSpace(os.Getenv(envName)); value != "" {
				view.HasAPIKey = true
				view.APIKeyMasked = maskWorkbenchSecret(value)
				break
			}
		}
		return view
	}
}

func RunWorkbenchProviderPrompt(config WorkbenchProviderConfig, systemPrompt string, userPrompt string) (string, WorkbenchProviderView, error) {
	resolved, err := resolveWorkbenchProviderConfig(config)
	if err != nil {
		return "", BuildWorkbenchProviderView(config), err
	}
	if strings.TrimSpace(systemPrompt) == "" {
		systemPrompt = resolved.SystemPrompt
	}
	messages := []workbenchProviderMessage{}
	if strings.TrimSpace(systemPrompt) != "" {
		messages = append(messages, workbenchProviderMessage{
			Role:    "system",
			Content: strings.TrimSpace(systemPrompt),
		})
	}
	messages = append(messages, workbenchProviderMessage{
		Role:    "user",
		Content: strings.TrimSpace(userPrompt),
	})

	var output string
	switch resolved.Config.Type {
	case WorkbenchProviderTypeOllama:
		output, err = runWorkbenchOllamaPrompt(*resolved, messages)
	default:
		output, err = runWorkbenchOpenAICompatiblePrompt(*resolved, messages)
	}
	viewConfig := resolved.Config
	viewConfig.APIKey = firstNonEmpty(strings.TrimSpace(config.APIKey), strings.TrimSpace(resolved.APIKey))
	view := BuildWorkbenchProviderView(viewConfig)
	if err != nil {
		return "", view, err
	}
	return strings.TrimSpace(output), view, nil
}

func ExtractWorkbenchFlowYAML(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	matches := workbenchYAMLFencePattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		candidate := strings.TrimSpace(match[1])
		if candidate == "" {
			continue
		}
		if _, err := ParseFlow([]byte(candidate), "yaml"); err == nil {
			return candidate
		}
	}
	if _, err := ParseFlow([]byte(text), "yaml"); err == nil {
		return text
	}
	for _, match := range matches {
		candidate := strings.TrimSpace(match[1])
		if candidate != "" {
			return candidate
		}
	}
	return text
}

func runWorkbenchOpenAICompatiblePrompt(resolved workbenchResolvedProvider, messages []workbenchProviderMessage) (string, error) {
	requestBody := map[string]any{
		"model":       resolved.Model,
		"messages":    messages,
		"temperature": defaultWorkbenchProviderTemperature,
		"max_tokens":  defaultWorkbenchProviderMaxTokens,
	}
	urlValue := workbenchOpenAIChatEndpoint(resolved.BaseURL)
	body, err := workbenchProviderPOST(urlValue, requestBody, func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+resolved.APIKey)
	})
	if err != nil {
		return "", err
	}
	var response struct {
		Choices []struct {
			Message struct {
				Content any `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decode openai-compatible response: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("openai-compatible response does not contain choices")
	}
	content := workbenchMessageContentString(response.Choices[0].Message.Content)
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("openai-compatible response does not contain message content")
	}
	return content, nil
}

func runWorkbenchOllamaPrompt(resolved workbenchResolvedProvider, messages []workbenchProviderMessage) (string, error) {
	requestBody := map[string]any{
		"model":    resolved.Model,
		"messages": messages,
		"stream":   false,
		"options": map[string]any{
			"temperature": defaultWorkbenchProviderTemperature,
			"max_tokens":  defaultWorkbenchProviderMaxTokens,
		},
	}
	urlValue := workbenchOllamaChatEndpoint(resolved.BaseURL)
	body, err := workbenchProviderPOST(urlValue, requestBody, nil)
	if err == nil {
		var response any
		if err := json.Unmarshal(body, &response); err != nil {
			return "", fmt.Errorf("decode ollama response: %w", err)
		}
		content := workbenchProviderResponseContent(response)
		if strings.TrimSpace(content) != "" {
			return strings.TrimSpace(content), nil
		}
	}

	content, generateErr := runWorkbenchOllamaGeneratePrompt(resolved, messages)
	if strings.TrimSpace(content) != "" && generateErr == nil {
		return strings.TrimSpace(content), nil
	}
	if err != nil {
		if generateErr != nil {
			return "", fmt.Errorf("ollama chat request failed: %v; ollama generate fallback failed: %v", err, generateErr)
		}
		return "", err
	}
	if generateErr != nil {
		return "", fmt.Errorf("ollama response does not contain message content; generate fallback failed: %v", generateErr)
	}
	return "", fmt.Errorf("ollama response does not contain message content")
}

func runWorkbenchOllamaGeneratePrompt(resolved workbenchResolvedProvider, messages []workbenchProviderMessage) (string, error) {
	requestBody := map[string]any{
		"model":  resolved.Model,
		"prompt": workbenchMessagesToPrompt(messages),
		"stream": false,
		"options": map[string]any{
			"temperature": defaultWorkbenchProviderTemperature,
			"max_tokens":  defaultWorkbenchProviderMaxTokens,
		},
	}
	urlValue := workbenchOllamaGenerateEndpoint(resolved.BaseURL)
	body, err := workbenchProviderPOST(urlValue, requestBody, nil)
	if err != nil {
		return "", err
	}
	var response any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decode ollama generate response: %w", err)
	}
	content := workbenchProviderResponseContent(response)
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("ollama generate response does not contain content")
	}
	return strings.TrimSpace(content), nil
}

func workbenchProviderPOST(urlValue string, bodyValue any, customize func(*http.Request)) ([]byte, error) {
	encoded, err := json.Marshal(bodyValue)
	if err != nil {
		return nil, fmt.Errorf("encode provider request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, urlValue, bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("create provider request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if customize != nil {
		customize(req)
	}
	client := &http.Client{Timeout: defaultWorkbenchProviderHTTPTimeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider request failed: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("read provider response: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("provider request failed with status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func workbenchOpenAIChatEndpoint(baseURL string) string {
	baseURL = normalizeWorkbenchProviderBaseURL(baseURL)
	switch {
	case strings.HasSuffix(baseURL, "/chat/completions"):
		return baseURL
	case strings.HasSuffix(baseURL, "/v1"):
		return baseURL + "/chat/completions"
	default:
		return baseURL + "/v1/chat/completions"
	}
}

func workbenchOllamaChatEndpoint(baseURL string) string {
	baseURL = normalizeWorkbenchProviderBaseURL(baseURL)
	if strings.HasSuffix(baseURL, "/api/chat") {
		return baseURL
	}
	return baseURL + "/api/chat"
}

func workbenchOllamaGenerateEndpoint(baseURL string) string {
	baseURL = normalizeWorkbenchProviderBaseURL(baseURL)
	switch {
	case strings.HasSuffix(baseURL, "/api/generate"):
		return baseURL
	case strings.HasSuffix(baseURL, "/api/chat"):
		return strings.TrimSuffix(baseURL, "/api/chat") + "/api/generate"
	default:
		return baseURL + "/api/generate"
	}
}

func workbenchMessageContentString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			if content := workbenchMessageContentString(item); content != "" {
				parts = append(parts, content)
			}
		}
		return strings.TrimSpace(strings.Join(parts, "\n"))
	case map[string]any:
		if text := strings.TrimSpace(fmt.Sprint(typed["text"])); text != "" && text != "<nil>" {
			return text
		}
		if nested, ok := typed["content"]; ok {
			return workbenchMessageContentString(nested)
		}
		return ""
	default:
		return strings.TrimSpace(fmt.Sprint(value))
	}
}

func workbenchProviderResponseContent(value any) string {
	switch typed := value.(type) {
	case map[string]any:
		for _, key := range []string{"message", "data", "response", "content", "output_text", "text", "completion", "result"} {
			if nested, ok := typed[key]; ok {
				if content := workbenchProviderResponseContent(nested); content != "" {
					return content
				}
			}
		}
		if choices, ok := typed["choices"]; ok {
			if content := workbenchProviderResponseContent(choices); content != "" {
				return content
			}
		}
		return ""
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			if content := workbenchProviderResponseContent(item); content != "" {
				parts = append(parts, content)
			}
		}
		return strings.TrimSpace(strings.Join(parts, "\n"))
	default:
		return workbenchMessageContentString(value)
	}
}

func workbenchMessagesToPrompt(messages []workbenchProviderMessage) string {
	parts := make([]string, 0, len(messages))
	for _, message := range messages {
		content := strings.TrimSpace(message.Content)
		if content == "" {
			continue
		}
		role := strings.TrimSpace(message.Role)
		if role == "" {
			role = "user"
		}
		parts = append(parts, fmt.Sprintf("%s:\n%s", strings.ToUpper(role), content))
	}
	return strings.TrimSpace(strings.Join(parts, "\n\n"))
}

func maskWorkbenchSecret(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}
