package tsplay_core

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

const (
	defaultGoddddocrBaseURL       = "http://127.0.0.1:8088"
	defaultGoddddocrEndpoint      = defaultGoddddocrBaseURL + "/ocr/file"
	defaultGoddddocrReadyEndpoint = defaultGoddddocrBaseURL + "/ready"
)

func runFlowOCRReadyStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return nil, err
	}
	endpoint = normalizeGoddddocrReadyEndpoint(endpoint)

	with := map[string]any{"response_as": "json"}
	if value, ok, err := flowStepResolvedParam(ctx, step, "timeout"); err != nil {
		return nil, err
	} else if ok {
		with["timeout"] = value
	}
	if value, ok, err := flowStepResolvedParam(ctx, step, "save_path"); err != nil {
		return nil, err
	} else if ok {
		with["save_path"] = value
	}

	responseValue, err := runFlowHTTPRequestStep(L, ctx, FlowStep{
		Action: "http_request",
		URL:    endpoint,
		With:   with,
	})
	if err != nil {
		return nil, err
	}

	response, ok := responseValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("ocr_ready expected http response object, got %T", responseValue)
	}

	result, ready := buildOCRReadyResult(response)
	strict := true
	if value, ok, err := flowStepResolvedParam(ctx, step, "strict"); err != nil {
		return nil, err
	} else if ok {
		strict, err = ocrBoolParam(value)
		if err != nil {
			return nil, fmt.Errorf("ocr_ready strict %w", err)
		}
	}
	if strict && !ready {
		return nil, fmt.Errorf("ocr_ready failed with HTTP status %v: %s", response["status"], describeOCRResponseBody(response["body"]))
	}
	return result, nil
}

func runFlowOCRRequestStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	filePath, err := flowStepStringParam(ctx, step, "file_path")
	if err != nil {
		return nil, err
	}

	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return nil, err
	}
	endpoint = normalizeGoddddocrEndpoint(endpoint)

	fieldName, err := flowStepOptionalStringParam(ctx, step, "field_name")
	if err != nil {
		return nil, err
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		fieldName = "file"
	}

	multipartFields := map[string]any{"confidence": true}
	if value, ok, err := flowStepResolvedParam(ctx, step, "charset_range"); err != nil {
		return nil, err
	} else if ok && strings.TrimSpace(fmt.Sprint(value)) != "" {
		multipartFields["charset_range"] = value
	}
	if value, ok, err := flowStepResolvedParam(ctx, step, "confidence"); err != nil {
		return nil, err
	} else if ok {
		confidence, err := ocrBoolParam(value)
		if err != nil {
			return nil, fmt.Errorf("ocr_request confidence %w", err)
		}
		multipartFields["confidence"] = confidence
	}

	with := map[string]any{
		"multipart_files":  map[string]any{fieldName: filePath},
		"multipart_fields": multipartFields,
		"response_as":      "json",
	}
	if value, ok, err := flowStepResolvedParam(ctx, step, "timeout"); err != nil {
		return nil, err
	} else if ok {
		with["timeout"] = value
	}
	if value, ok, err := flowStepResolvedParam(ctx, step, "save_path"); err != nil {
		return nil, err
	} else if ok {
		with["save_path"] = value
	}

	responseValue, err := runFlowHTTPRequestStep(L, ctx, FlowStep{
		Action: "http_request",
		URL:    endpoint,
		With:   with,
	})
	if err != nil {
		return nil, err
	}

	response, ok := responseValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("ocr_request expected http response object, got %T", responseValue)
	}
	return buildOCRRequestResult(response)
}

func goddddocrEndpointValue(value string) string {
	endpoint := strings.TrimSpace(value)
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("GODDDDOCR_URL"))
	}
	return endpoint
}

func normalizeGoddddocrEndpoint(value string) string {
	endpoint := goddddocrEndpointValue(value)
	if endpoint == "" {
		return defaultGoddddocrEndpoint
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return endpoint
	}
	switch strings.Trim(parsed.Path, "/") {
	case "", "ocr", "ready":
		parsed.Path = "/ocr/file"
		return parsed.String()
	}
	return endpoint
}

func normalizeGoddddocrReadyEndpoint(value string) string {
	endpoint := goddddocrEndpointValue(value)
	if endpoint == "" {
		return defaultGoddddocrReadyEndpoint
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return endpoint
	}
	switch strings.Trim(parsed.Path, "/") {
	case "", "health", "ocr", "ocr/file":
		parsed.Path = "/ready"
		return parsed.String()
	}
	return endpoint
}

func ocrBoolParam(value any) (bool, error) {
	if text, ok := value.(string); ok {
		parsed, err := strconv.ParseBool(strings.TrimSpace(text))
		if err != nil {
			return false, fmt.Errorf("must be a boolean")
		}
		return parsed, nil
	}
	return boolParam(value)
}

func buildOCRReadyResult(response map[string]any) (map[string]any, bool) {
	httpOK, _ := response["ok"].(bool)
	body, _ := response["body"].(map[string]any)

	serviceStatus := ""
	if body != nil {
		serviceStatus = strings.TrimSpace(fmt.Sprint(body["status"]))
	}
	bodyReady, _ := body["ready"].(bool)
	ready := httpOK && (bodyReady || strings.EqualFold(serviceStatus, "ready"))

	result := map[string]any{
		"ready":    ready,
		"ok":       httpOK,
		"status":   response["status"],
		"url":      response["url"],
		"response": response,
	}
	if serviceStatus != "" {
		result["service_status"] = serviceStatus
	}
	for _, name := range []string{"model", "time", "request_id"} {
		if value, ok := body[name]; ok {
			result[name] = value
		}
	}
	if savePath, ok := response["save_path"]; ok && savePath != nil {
		result["save_path"] = savePath
	}
	return result, ready
}

func buildOCRRequestResult(response map[string]any) (map[string]any, error) {
	if okValue, ok := response["ok"].(bool); ok && !okValue {
		return nil, fmt.Errorf("ocr_request failed with HTTP status %v: %s", response["status"], describeOCRResponseBody(response["body"]))
	}

	body, ok := response["body"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("ocr_request expected JSON object response body, got %T", response["body"])
	}

	resultValue, ok := body["result"]
	if !ok {
		if textValue, hasText := body["text"]; hasText {
			resultValue = textValue
			ok = true
		}
	}
	if !ok {
		return nil, fmt.Errorf("ocr_request response body missing result")
	}
	text := fmt.Sprint(resultValue)

	result := map[string]any{
		"text":     text,
		"result":   text,
		"ok":       response["ok"],
		"status":   response["status"],
		"response": response,
	}
	for _, name := range []string{"confidence", "request_id", "processing_time_ms"} {
		if value, ok := body[name]; ok {
			result[name] = value
		}
	}
	if savePath, ok := response["save_path"]; ok && savePath != nil {
		result["save_path"] = savePath
	}
	return result, nil
}

func describeOCRResponseBody(value any) string {
	body, ok := value.(map[string]any)
	if !ok {
		return fmt.Sprint(value)
	}
	for _, key := range []string{"message", "error", "code"} {
		if item, ok := body[key]; ok && strings.TrimSpace(fmt.Sprint(item)) != "" {
			return fmt.Sprint(item)
		}
	}
	return fmt.Sprint(body)
}
