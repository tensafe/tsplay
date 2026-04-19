package tsplay_core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type flowHTTPRequestConfig struct {
	Method              string
	URL                 string
	Headers             map[string]string
	Query               map[string]any
	Body                string
	JSON                any
	Form                map[string]any
	MultipartFiles      map[string]string
	MultipartFields     map[string]any
	TimeoutMS           int
	ResponseAs          string
	UseBrowserCookies   bool
	UseBrowserReferer   bool
	UseBrowserUserAgent bool
	SavePath            string
}

func http_request(L *lua.LState) int {
	values, err := httpRequestValuesFromLua(L)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	config, err := normalizeHTTPRequestConfig(values)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err = applyLuaHTTPRequestRuntimePolicy(L, config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeHTTPRequest(L, config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func json_extract(L *lua.LState) int {
	value := luaValueToGo(L.CheckAny(1))
	path := L.CheckString(2)
	result, err := extractJSONPathValue(value, path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func runFlowHTTPRequestStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	values, err := resolvedHTTPRequestValues(ctx, step)
	if err != nil {
		return nil, err
	}
	config, err := normalizeHTTPRequestConfig(values)
	if err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Security != nil && ctx.Security.AllowFileAccess {
		config, err = rewriteHTTPRequestRuntimePaths(config, *ctx.Security)
		if err != nil {
			return nil, err
		}
	}
	return executeHTTPRequest(L, config)
}

func runFlowJSONExtractStep(ctx *FlowContext, step FlowStep) (any, error) {
	value, ok, err := flowStepResolvedParam(ctx, step, "from")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("action %q requires %q", step.Action, "from")
	}
	path, err := flowStepStringParam(ctx, step, "path")
	if err != nil {
		return nil, err
	}
	return extractJSONPathValue(value, path)
}

func httpRequestValuesFromLua(L *lua.LState) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("http_request requires either a config table or a url")
	}

	first := luaValueToGo(L.CheckAny(1))
	if values, ok := first.(map[string]any); ok {
		return values, nil
	}

	values := map[string]any{"url": first}
	if L.GetTop() >= 2 {
		values["method"] = luaValueToGo(L.CheckAny(2))
	}
	if L.GetTop() >= 3 {
		values["body"] = luaValueToGo(L.CheckAny(3))
	}
	return values, nil
}

func resolvedHTTPRequestValues(ctx *FlowContext, step FlowStep) (map[string]any, error) {
	names := []string{
		"url",
		"method",
		"headers",
		"query",
		"body",
		"json",
		"form",
		"multipart_files",
		"multipart_fields",
		"timeout",
		"response_as",
		"use_browser_cookies",
		"use_browser_referer",
		"use_browser_user_agent",
		"save_path",
	}
	values := map[string]any{}
	for _, name := range names {
		value, ok, err := flowStepResolvedParam(ctx, step, name)
		if err != nil {
			return nil, err
		}
		if ok {
			values[name] = value
		}
	}
	return values, nil
}

func flowStepResolvedParam(ctx *FlowContext, step FlowStep, name string) (any, bool, error) {
	value, ok := step.param(name)
	if !ok {
		return nil, false, nil
	}
	if ctx == nil {
		return value, true, nil
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return nil, false, err
	}
	return resolved, true, nil
}

func normalizeHTTPRequestConfig(values map[string]any) (flowHTTPRequestConfig, error) {
	config := flowHTTPRequestConfig{
		Headers:         map[string]string{},
		Query:           map[string]any{},
		Form:            map[string]any{},
		MultipartFiles:  map[string]string{},
		MultipartFields: map[string]any{},
		TimeoutMS:       30000,
		ResponseAs:      "auto",
	}

	urlValue, ok := values["url"]
	if !ok || strings.TrimSpace(fmt.Sprint(urlValue)) == "" {
		return flowHTTPRequestConfig{}, fmt.Errorf("http_request requires url")
	}
	config.URL = strings.TrimSpace(fmt.Sprint(urlValue))

	if method, ok := values["method"]; ok {
		config.Method = strings.ToUpper(strings.TrimSpace(fmt.Sprint(method)))
	}
	if body, ok := values["body"]; ok {
		bodyText, ok := body.(string)
		if !ok {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request body must be a string")
		}
		config.Body = bodyText
	}
	if jsonValue, ok := values["json"]; ok {
		config.JSON = jsonValue
	}
	if headers, ok := values["headers"]; ok {
		stringMap, err := stringMapValue(headers, "headers")
		if err != nil {
			return flowHTTPRequestConfig{}, err
		}
		config.Headers = stringMap
	}
	if query, ok := values["query"]; ok {
		objectValue, err := objectMapValue(query, "query")
		if err != nil {
			return flowHTTPRequestConfig{}, err
		}
		config.Query = objectValue
	}
	if form, ok := values["form"]; ok {
		objectValue, err := objectMapValue(form, "form")
		if err != nil {
			return flowHTTPRequestConfig{}, err
		}
		config.Form = objectValue
	}
	if multipartFiles, ok := values["multipart_files"]; ok {
		fileMap, err := stringMapValue(multipartFiles, "multipart_files")
		if err != nil {
			return flowHTTPRequestConfig{}, err
		}
		config.MultipartFiles = fileMap
	}
	if multipartFields, ok := values["multipart_fields"]; ok {
		objectValue, err := objectMapValue(multipartFields, "multipart_fields")
		if err != nil {
			return flowHTTPRequestConfig{}, err
		}
		config.MultipartFields = objectValue
	}
	if timeout, ok := values["timeout"]; ok {
		timeoutMS, err := intParam(timeout)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request timeout %w", err)
		}
		if timeoutMS < 1 {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request timeout must be at least 1")
		}
		config.TimeoutMS = timeoutMS
	}
	if responseAs, ok := values["response_as"]; ok {
		config.ResponseAs = strings.ToLower(strings.TrimSpace(fmt.Sprint(responseAs)))
	}
	switch config.ResponseAs {
	case "", "auto":
		config.ResponseAs = "auto"
	case "text", "json":
	default:
		return flowHTTPRequestConfig{}, fmt.Errorf("http_request response_as must be one of auto, text, or json")
	}
	if useCookies, ok := values["use_browser_cookies"]; ok {
		boolValue, err := boolParam(useCookies)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request use_browser_cookies %w", err)
		}
		config.UseBrowserCookies = boolValue
	}
	if useReferer, ok := values["use_browser_referer"]; ok {
		boolValue, err := boolParam(useReferer)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request use_browser_referer %w", err)
		}
		config.UseBrowserReferer = boolValue
	}
	if useUserAgent, ok := values["use_browser_user_agent"]; ok {
		boolValue, err := boolParam(useUserAgent)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request use_browser_user_agent %w", err)
		}
		config.UseBrowserUserAgent = boolValue
	}
	if savePath, ok := values["save_path"]; ok {
		savePathText, ok := savePath.(string)
		if !ok {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request save_path must be a string")
		}
		config.SavePath = savePathText
	}

	bodyModes := 0
	if config.Body != "" {
		bodyModes++
	}
	if config.JSON != nil {
		bodyModes++
	}
	if len(config.Form) > 0 {
		bodyModes++
	}
	if len(config.MultipartFiles) > 0 || len(config.MultipartFields) > 0 {
		bodyModes++
	}
	if bodyModes > 1 {
		return flowHTTPRequestConfig{}, fmt.Errorf("http_request accepts only one of body, json, form, or multipart data")
	}

	if config.Method == "" {
		if bodyModes > 0 {
			config.Method = http.MethodPost
		} else {
			config.Method = http.MethodGet
		}
	}
	return config, nil
}

func applyLuaHTTPRequestRuntimePolicy(L *lua.LState, config flowHTTPRequestConfig) (flowHTTPRequestConfig, error) {
	ctx := flowContextFromState(L)
	if ctx == nil || ctx.Security == nil {
		return config, nil
	}
	if !ctx.Security.AllowHTTP {
		return flowHTTPRequestConfig{}, fmt.Errorf("http_request is disabled by security policy; set allow_http=true only for trusted flows")
	}
	if config.SavePath == "" && len(config.MultipartFiles) == 0 {
		return config, nil
	}
	if !ctx.Security.AllowFileAccess {
		return flowHTTPRequestConfig{}, fmt.Errorf("http_request file access is disabled by security policy; set allow_file_access=true only for trusted flows")
	}
	return rewriteHTTPRequestRuntimePaths(config, *ctx.Security)
}

func rewriteHTTPRequestRuntimePaths(config flowHTTPRequestConfig, policy FlowSecurityPolicy) (flowHTTPRequestConfig, error) {
	if config.SavePath != "" {
		resolved, err := resolveRuntimeFilePath(config.SavePath, flowFileOutputPath, policy)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request save_path %w", err)
		}
		config.SavePath = resolved
	}
	if len(config.MultipartFiles) == 0 {
		return config, nil
	}
	rewritten := map[string]string{}
	for key, path := range config.MultipartFiles {
		resolved, err := resolveRuntimeFilePath(path, flowFileInputPath, policy)
		if err != nil {
			return flowHTTPRequestConfig{}, fmt.Errorf("http_request multipart_files.%s %w", key, err)
		}
		rewritten[key] = resolved
	}
	config.MultipartFiles = rewritten
	return config, nil
}

func executeHTTPRequest(L *lua.LState, config flowHTTPRequestConfig) (map[string]any, error) {
	requestURL, err := addQueryParams(config.URL, config.Query)
	if err != nil {
		return nil, err
	}

	bodyReader, contentType, err := buildHTTPRequestBody(config)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(config.Method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}
	if contentType != "" && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", contentType)
	}
	if err := applyHTTPRequestBrowserHeaders(L, request, config); err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Duration(config.TimeoutMS) * time.Millisecond}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("perform http request: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response body: %w", err)
	}
	if config.SavePath != "" {
		if err := os.MkdirAll(filepath.Dir(config.SavePath), 0755); err != nil {
			return nil, fmt.Errorf("create http response output directory: %w", err)
		}
		if err := os.WriteFile(config.SavePath, bodyBytes, 0600); err != nil {
			return nil, fmt.Errorf("write http response body: %w", err)
		}
	}

	responseBody, err := decodeHTTPResponseBody(bodyBytes, response.Header.Get("Content-Type"), config.ResponseAs)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"method":       config.Method,
		"url":          response.Request.URL.String(),
		"status":       response.StatusCode,
		"ok":           response.StatusCode >= 200 && response.StatusCode < 300,
		"headers":      headerMap(response.Header),
		"content_type": response.Header.Get("Content-Type"),
		"body":         responseBody,
		"save_path":    emptyToNil(config.SavePath),
	}, nil
}

func addQueryParams(rawURL string, params map[string]any) (string, error) {
	if len(params) == 0 {
		return rawURL, nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url %q: %w", rawURL, err)
	}
	values := parsed.Query()
	appendURLValues(values, params)
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func buildHTTPRequestBody(config flowHTTPRequestConfig) (io.Reader, string, error) {
	if len(config.MultipartFiles) > 0 || len(config.MultipartFields) > 0 {
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)
		for key, value := range config.MultipartFields {
			if err := writer.WriteField(key, fmt.Sprint(value)); err != nil {
				return nil, "", fmt.Errorf("write multipart field %q: %w", key, err)
			}
		}
		for key, path := range config.MultipartFiles {
			file, err := os.Open(path)
			if err != nil {
				return nil, "", fmt.Errorf("open multipart file %q: %w", path, err)
			}
			part, err := writer.CreateFormFile(key, filepath.Base(path))
			if err != nil {
				file.Close()
				return nil, "", fmt.Errorf("create multipart file %q: %w", key, err)
			}
			if _, err := io.Copy(part, file); err != nil {
				file.Close()
				return nil, "", fmt.Errorf("copy multipart file %q: %w", key, err)
			}
			file.Close()
		}
		if err := writer.Close(); err != nil {
			return nil, "", fmt.Errorf("close multipart body: %w", err)
		}
		return bytes.NewReader(buffer.Bytes()), writer.FormDataContentType(), nil
	}
	if config.JSON != nil {
		payload, err := json.Marshal(config.JSON)
		if err != nil {
			return nil, "", fmt.Errorf("marshal http_request json body: %w", err)
		}
		return bytes.NewReader(payload), "application/json", nil
	}
	if len(config.Form) > 0 {
		values := url.Values{}
		appendURLValues(values, config.Form)
		encoded := values.Encode()
		return strings.NewReader(encoded), "application/x-www-form-urlencoded", nil
	}
	if config.Body != "" {
		return strings.NewReader(config.Body), "", nil
	}
	return nil, "", nil
}

func applyHTTPRequestBrowserHeaders(L *lua.LState, request *http.Request, config flowHTTPRequestConfig) error {
	if request == nil || L == nil {
		return nil
	}
	context, hasContext := flowBrowserContextFromState(L)
	page, hasPage := flowPageFromState(L)

	if config.UseBrowserCookies {
		if !hasContext {
			return fmt.Errorf("http_request use_browser_cookies requires an active browser context")
		}
		cookies, err := context.Cookies(request.URL.String())
		if err != nil {
			return fmt.Errorf("load browser cookies: %w", err)
		}
		pairs := make([]string, 0, len(cookies))
		for _, cookie := range cookies {
			pairs = append(pairs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
		}
		if len(pairs) > 0 && request.Header.Get("Cookie") == "" {
			request.Header.Set("Cookie", strings.Join(pairs, "; "))
		}
	}
	if config.UseBrowserReferer {
		if !hasPage {
			return fmt.Errorf("http_request use_browser_referer requires an active page")
		}
		if request.Header.Get("Referer") == "" {
			request.Header.Set("Referer", page.URL())
		}
	}
	if config.UseBrowserUserAgent {
		if !hasPage {
			return fmt.Errorf("http_request use_browser_user_agent requires an active page")
		}
		userAgent, err := page.Evaluate("navigator.userAgent")
		if err != nil {
			return fmt.Errorf("load browser user agent: %w", err)
		}
		if request.Header.Get("User-Agent") == "" {
			request.Header.Set("User-Agent", fmt.Sprint(userAgent))
		}
	}
	return nil
}

func decodeHTTPResponseBody(body []byte, contentType string, responseAs string) (any, error) {
	switch responseAs {
	case "text":
		return string(body), nil
	case "json":
		if len(bytes.TrimSpace(body)) == 0 {
			return nil, nil
		}
		var value any
		if err := json.Unmarshal(body, &value); err != nil {
			return nil, fmt.Errorf("decode http response json: %w", err)
		}
		return value, nil
	default:
		if looksLikeJSONContentType(contentType) {
			var value any
			if len(bytes.TrimSpace(body)) == 0 {
				return nil, nil
			}
			if err := json.Unmarshal(body, &value); err == nil {
				return value, nil
			}
		}
		return string(body), nil
	}
}

func looksLikeJSONContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "/json") || strings.Contains(contentType, "+json")
}

func extractJSONPathValue(value any, path string) (any, error) {
	tokens, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}
	current, err := normalizeJSONValue(value)
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		current, err = normalizeJSONValue(current)
		if err != nil {
			return nil, err
		}
		if token.field != nil {
			object, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("json_extract path %q expected object before field %q, got %T", path, *token.field, current)
			}
			next, ok := object[*token.field]
			if !ok {
				return nil, fmt.Errorf("json_extract path %q field %q not found", path, *token.field)
			}
			current = next
			continue
		}
		array, ok := current.([]any)
		if !ok {
			return nil, fmt.Errorf("json_extract path %q expected array before index %d, got %T", path, token.index, current)
		}
		if token.index < 0 || token.index >= len(array) {
			return nil, fmt.Errorf("json_extract path %q index %d out of range", path, token.index)
		}
		current = array[token.index]
	}
	return current, nil
}

type jsonPathToken struct {
	field *string
	index int
}

func parseJSONPath(path string) ([]jsonPathToken, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("json_extract path is required")
	}
	if path == "$" {
		return nil, nil
	}
	if !strings.HasPrefix(path, "$") {
		return nil, fmt.Errorf("json_extract path must start with $")
	}
	tokens := []jsonPathToken{}
	for i := 1; i < len(path); {
		switch path[i] {
		case '.':
			i++
			start := i
			for i < len(path) && path[i] != '.' && path[i] != '[' {
				i++
			}
			if start == i {
				return nil, fmt.Errorf("json_extract path %q has an empty field segment", path)
			}
			field := path[start:i]
			tokens = append(tokens, jsonPathToken{field: &field})
		case '[':
			i++
			if i >= len(path) {
				return nil, fmt.Errorf("json_extract path %q has an unterminated bracket", path)
			}
			if path[i] == '"' || path[i] == '\'' {
				quote := path[i]
				i++
				start := i
				for i < len(path) && path[i] != quote {
					i++
				}
				if i >= len(path) {
					return nil, fmt.Errorf("json_extract path %q has an unterminated quoted field", path)
				}
				field := path[start:i]
				i++
				if i >= len(path) || path[i] != ']' {
					return nil, fmt.Errorf("json_extract path %q has an invalid quoted field segment", path)
				}
				i++
				tokens = append(tokens, jsonPathToken{field: &field})
				continue
			}
			start := i
			for i < len(path) && path[i] != ']' {
				i++
			}
			if i >= len(path) {
				return nil, fmt.Errorf("json_extract path %q has an unterminated index segment", path)
			}
			indexValue := strings.TrimSpace(path[start:i])
			i++
			index, err := strconv.Atoi(indexValue)
			if err != nil {
				return nil, fmt.Errorf("json_extract path %q has an invalid array index %q", path, indexValue)
			}
			tokens = append(tokens, jsonPathToken{index: index})
		default:
			return nil, fmt.Errorf("json_extract path %q has an invalid segment near %q", path, path[i:])
		}
	}
	return tokens, nil
}

func normalizeJSONValue(value any) (any, error) {
	switch typed := value.(type) {
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return typed, nil
		}
		if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
			(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
			var parsed any
			if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
				return parsed, nil
			}
		}
		return typed, nil
	case []byte:
		var parsed any
		if err := json.Unmarshal(typed, &parsed); err != nil {
			return nil, fmt.Errorf("decode json bytes: %w", err)
		}
		return parsed, nil
	default:
		return value, nil
	}
}

func stringMapValue(value any, name string) (map[string]string, error) {
	objectValue, err := objectMapValue(value, name)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for key, item := range objectValue {
		if item == nil {
			continue
		}
		result[key] = fmt.Sprint(item)
	}
	return result, nil
}

func objectMapValue(value any, name string) (map[string]any, error) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, nil
	case map[string]string:
		result := map[string]any{}
		for key, item := range typed {
			result[key] = item
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%s must be an object", name)
	}
}

func appendURLValues(values url.Values, params map[string]any) {
	for key, value := range params {
		switch typed := value.(type) {
		case nil:
			continue
		case []any:
			for _, item := range typed {
				values.Add(key, fmt.Sprint(item))
			}
		case []string:
			for _, item := range typed {
				values.Add(key, item)
			}
		default:
			values.Add(key, fmt.Sprint(typed))
		}
	}
}

func headerMap(header http.Header) map[string]any {
	result := map[string]any{}
	for key, values := range header {
		result[key] = strings.Join(values, ", ")
	}
	return result
}

func emptyToNil(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func boolParam(value any) (bool, error) {
	switch typed := value.(type) {
	case bool:
		return typed, nil
	default:
		return false, fmt.Errorf("must be a boolean")
	}
}
