package tsplay_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

const (
	defaultGoddddocrBaseURL                 = "http://127.0.0.1:8088"
	defaultGoddddocrEndpoint                = defaultGoddddocrBaseURL + "/ocr/file"
	defaultGoddddocrReadyEndpoint           = defaultGoddddocrBaseURL + "/ready"
	defaultGoddddocrDetectEndpoint          = defaultGoddddocrBaseURL + "/det/file"
	defaultGoddddocrSlideComparisonEndpoint = defaultGoddddocrBaseURL + "/slide_comparison/file"
	defaultGoddddocrSlideMatchEndpoint      = defaultGoddddocrBaseURL + "/slide_match/file"
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
	for _, name := range []string{"color_filter_colors", "color_filter_custom_ranges"} {
		if value, ok, err := flowStepResolvedParam(ctx, step, name); err != nil {
			return nil, err
		} else if ok {
			fieldValue, err := ocrMultipartFieldString(value)
			if err != nil {
				return nil, fmt.Errorf("ocr_request %s %w", name, err)
			}
			if fieldValue != "" {
				multipartFields[name] = fieldValue
			}
		}
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
	if value, ok, err := flowStepResolvedParam(ctx, step, "probability"); err != nil {
		return nil, err
	} else if ok {
		probability, err := ocrBoolParam(value)
		if err != nil {
			return nil, fmt.Errorf("ocr_request probability %w", err)
		}
		multipartFields["probability"] = probability
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

func runFlowOCRDetectStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	filePath, err := flowStepStringParam(ctx, step, "file_path")
	if err != nil {
		return nil, err
	}

	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return nil, err
	}
	endpoint = normalizeGoddddocrDetectEndpoint(endpoint)

	fieldName, err := flowStepOptionalStringParam(ctx, step, "field_name")
	if err != nil {
		return nil, err
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		fieldName = "file"
	}

	multipartFields := map[string]any{"detailed": true}
	if value, ok, err := flowStepResolvedParam(ctx, step, "detailed"); err != nil {
		return nil, err
	} else if ok {
		detailed, err := ocrBoolParam(value)
		if err != nil {
			return nil, fmt.Errorf("ocr_detect detailed %w", err)
		}
		multipartFields["detailed"] = detailed
	}
	for _, name := range []string{"score_threshold", "nms_threshold"} {
		if value, ok, err := flowStepResolvedParam(ctx, step, name); err != nil {
			return nil, err
		} else if ok {
			threshold, err := ocrFloatParam(value)
			if err != nil {
				return nil, fmt.Errorf("ocr_detect %s %w", name, err)
			}
			multipartFields[name] = threshold
		}
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
		return nil, fmt.Errorf("ocr_detect expected http response object, got %T", responseValue)
	}
	return buildOCRDetectResult(response)
}

func runFlowOCRSlideComparisonStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	return runFlowOCRSlideStep(L, ctx, step, "ocr_slide_comparison", normalizeGoddddocrSlideComparisonEndpoint, false)
}

func runFlowOCRSlideMatchStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	return runFlowOCRSlideStep(L, ctx, step, "ocr_slide_match", normalizeGoddddocrSlideMatchEndpoint, true)
}

func runFlowOCRSlideStep(L *lua.LState, ctx *FlowContext, step FlowStep, actionName string, normalizeEndpoint func(string) string, includeSimpleTarget bool) (any, error) {
	targetFilePath, err := flowStepStringParam(ctx, step, "target_file_path")
	if err != nil {
		return nil, err
	}
	backgroundFilePath, err := flowStepStringParam(ctx, step, "background_file_path")
	if err != nil {
		return nil, err
	}

	endpoint, err := flowStepOptionalStringParam(ctx, step, "url")
	if err != nil {
		return nil, err
	}
	endpoint = normalizeEndpoint(endpoint)

	multipartFields := map[string]any{}
	if includeSimpleTarget {
		if value, ok, err := flowStepResolvedParam(ctx, step, "simple_target"); err != nil {
			return nil, err
		} else if ok {
			simpleTarget, err := ocrBoolParam(value)
			if err != nil {
				return nil, fmt.Errorf("%s simple_target %w", actionName, err)
			}
			multipartFields["simple_target"] = simpleTarget
		}
	}

	with := map[string]any{
		"multipart_files": map[string]any{
			"target_file":     targetFilePath,
			"background_file": backgroundFilePath,
		},
		"response_as": "json",
	}
	if len(multipartFields) > 0 {
		with["multipart_fields"] = multipartFields
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
		return nil, fmt.Errorf("%s expected http response object, got %T", actionName, responseValue)
	}
	return buildOCRSlideResult(response, actionName)
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

func normalizeGoddddocrDetectEndpoint(value string) string {
	return normalizeGoddddocrFeatureEndpoint(value, defaultGoddddocrDetectEndpoint, "/det/file", "det", "detect", "detection")
}

func normalizeGoddddocrSlideComparisonEndpoint(value string) string {
	return normalizeGoddddocrFeatureEndpoint(value, defaultGoddddocrSlideComparisonEndpoint, "/slide_comparison/file", "slide_comparison", "slide-comparison")
}

func normalizeGoddddocrSlideMatchEndpoint(value string) string {
	return normalizeGoddddocrFeatureEndpoint(value, defaultGoddddocrSlideMatchEndpoint, "/slide_match/file", "slide_match", "slide-match")
}

func normalizeGoddddocrFeatureEndpoint(value string, defaultEndpoint string, filePath string, aliases ...string) string {
	endpoint := goddddocrEndpointValue(value)
	if endpoint == "" {
		return defaultEndpoint
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return endpoint
	}
	path := strings.Trim(parsed.Path, "/")
	convertible := map[string]bool{
		"":         true,
		"ready":    true,
		"health":   true,
		"ocr":      true,
		"ocr/file": true,
	}
	for _, alias := range aliases {
		alias = strings.Trim(alias, "/")
		convertible[alias] = true
		convertible[alias+"/file"] = true
	}
	if convertible[path] {
		parsed.Path = filePath
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

func ocrFloatParam(value any) (float64, error) {
	return floatParam(value)
}

func ocrMultipartFieldString(value any) (string, error) {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text), nil
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("must be a string or JSON-serializable value")
	}
	return string(encoded), nil
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
	for _, name := range []string{"model", "time", "request_id", "detection", "slide_comparison", "slide_match"} {
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
	for _, name := range []string{"confidence", "probability", "request_id", "processing_time_ms"} {
		if value, ok := body[name]; ok {
			result[name] = value
		}
	}
	if savePath, ok := response["save_path"]; ok && savePath != nil {
		result["save_path"] = savePath
	}
	return result, nil
}

func buildOCRDetectResult(response map[string]any) (map[string]any, error) {
	if okValue, ok := response["ok"].(bool); ok && !okValue {
		return nil, fmt.Errorf("ocr_detect failed with HTTP status %v: %s", response["status"], describeOCRResponseBody(response["body"]))
	}

	body, ok := response["body"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("ocr_detect expected JSON object response body, got %T", response["body"])
	}
	resultValue, ok := body["result"]
	if !ok {
		return nil, fmt.Errorf("ocr_detect response body missing result")
	}

	result := map[string]any{
		"result":   resultValue,
		"ok":       response["ok"],
		"status":   response["status"],
		"response": response,
	}
	for _, name := range []string{"boxes", "request_id", "processing_time_ms"} {
		if value, ok := body[name]; ok {
			result[name] = value
		}
	}
	if savePath, ok := response["save_path"]; ok && savePath != nil {
		result["save_path"] = savePath
	}
	return result, nil
}

func buildOCRSlideResult(response map[string]any, actionName string) (map[string]any, error) {
	if okValue, ok := response["ok"].(bool); ok && !okValue {
		return nil, fmt.Errorf("%s failed with HTTP status %v: %s", actionName, response["status"], describeOCRResponseBody(response["body"]))
	}

	body, ok := response["body"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s expected JSON object response body, got %T", actionName, response["body"])
	}
	resultValue, ok := body["result"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s response body missing result object", actionName)
	}

	result := map[string]any{
		"result":   resultValue,
		"ok":       response["ok"],
		"status":   response["status"],
		"response": response,
	}
	for _, name := range []string{"target", "target_x", "target_y", "confidence"} {
		if value, ok := resultValue[name]; ok {
			result[name] = value
		}
	}
	for _, name := range []string{"request_id", "processing_time_ms"} {
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
