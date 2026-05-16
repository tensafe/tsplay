package tsplay_core

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	FlowRunStatusSucceeded            = "succeeded"
	FlowRunStatusFailed               = "failed"
	FlowRunStatusManualReviewRequired = "manual_review_required"

	FlowManualReviewStatus = "manual_review"
	FlowManualReviewAction = "manual_review_required"
)

type FlowManualReviewResult struct {
	Required   bool                       `json:"required"`
	Status     string                     `json:"status"`
	Action     string                     `json:"action"`
	Phase      string                     `json:"phase,omitempty"`
	Reason     string                     `json:"reason,omitempty"`
	PayloadVar string                     `json:"payload_var,omitempty"`
	Evidence   map[string]any             `json:"evidence,omitempty"`
	Artifacts  []FlowManualReviewArtifact `json:"artifacts,omitempty"`
	Payload    map[string]any             `json:"payload,omitempty"`
}

type FlowManualReviewArtifact struct {
	Name         string `json:"name"`
	Kind         string `json:"kind,omitempty"`
	Path         string `json:"path"`
	RelativePath string `json:"relative_path,omitempty"`
	ContentType  string `json:"content_type,omitempty"`
	Exists       bool   `json:"exists"`
}

func FinalizeFlowResult(result *FlowResult, runErr error) {
	if result == nil {
		return
	}
	result.ManualReview = ExtractFlowManualReview(result)
	result.Status = FlowResultCompletionStatus(result, runErr)
}

func FlowResultCompletionStatus(result *FlowResult, runErr error) string {
	if runErr != nil {
		return FlowRunStatusFailed
	}
	if result != nil && result.ManualReview != nil && result.ManualReview.Required {
		return FlowRunStatusManualReviewRequired
	}
	return FlowRunStatusSucceeded
}

func ExtractFlowManualReview(result *FlowResult) *FlowManualReviewResult {
	if result == nil {
		return nil
	}
	if result.ManualReview != nil && result.ManualReview.Required {
		return result.ManualReview
	}
	name, payload, ok := findFlowManualReviewPayload(result.Vars)
	if !ok {
		return nil
	}
	status := firstNonEmpty(flowManualReviewString(payload["status"]), FlowManualReviewStatus)
	action := firstNonEmpty(flowManualReviewString(payload["action"]), FlowManualReviewAction)
	review := &FlowManualReviewResult{
		Required:   true,
		Status:     status,
		Action:     action,
		Phase:      flowManualReviewString(payload["phase"]),
		Reason:     flowManualReviewString(payload["reason"]),
		PayloadVar: name,
		Artifacts:  collectFlowManualReviewArtifacts(payload, result.ArtifactRoot),
	}
	if evidence, ok := flowManualReviewMap(payload["evidence"]); ok {
		if compacted, ok := compactTraceValue(evidence, 0).(map[string]any); ok {
			review.Evidence = compacted
		} else {
			review.Evidence = evidence
		}
	}
	if compacted, ok := compactTraceValue(payload, 0).(map[string]any); ok {
		review.Payload = compacted
	} else {
		review.Payload = payload
	}
	return review
}

func findFlowManualReviewPayload(vars map[string]any) (string, map[string]any, bool) {
	if len(vars) == 0 {
		return "", nil, false
	}
	for _, name := range []string{"payload", "manual_review", "manual_review_result", "result"} {
		if payload, ok := flowManualReviewMap(vars[name]); ok && isFlowManualReviewPayload(payload) {
			return name, payload, true
		}
	}
	names := make([]string, 0, len(vars))
	for name := range vars {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		payload, ok := flowManualReviewMap(vars[name])
		if !ok {
			continue
		}
		if isFlowManualReviewPayload(payload) {
			return name, payload, true
		}
	}
	return "", nil, false
}

func isFlowManualReviewPayload(payload map[string]any) bool {
	action := strings.ToLower(strings.TrimSpace(flowManualReviewString(payload["action"])))
	status := strings.ToLower(strings.TrimSpace(flowManualReviewString(payload["status"])))
	return action == FlowManualReviewAction || status == FlowManualReviewStatus
}

func collectFlowManualReviewArtifacts(payload map[string]any, artifactRoot string) []FlowManualReviewArtifact {
	artifacts := []FlowManualReviewArtifact{}
	seen := map[string]bool{}
	appendFromMap := func(values map[string]any) {
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if !flowManualReviewLooksLikeArtifactPath(key) {
				continue
			}
			pathValue := strings.TrimSpace(flowManualReviewString(values[key]))
			if pathValue == "" {
				continue
			}
			artifact := buildFlowManualReviewArtifact(key, pathValue, artifactRoot)
			seenKey := artifact.Name + "\x00" + artifact.Path
			if seen[seenKey] {
				continue
			}
			seen[seenKey] = true
			artifacts = append(artifacts, artifact)
		}
	}
	if evidence, ok := flowManualReviewMap(payload["evidence"]); ok {
		appendFromMap(evidence)
	}
	appendFromMap(payload)
	return artifacts
}

func flowManualReviewLooksLikeArtifactPath(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	return key == "path" ||
		key == "file_path" ||
		key == "artifact_path" ||
		strings.HasSuffix(key, "_path")
}

func buildFlowManualReviewArtifact(name string, pathValue string, artifactRoot string) FlowManualReviewArtifact {
	resolvedPath, relativePath, exists := resolveFlowManualReviewArtifactPath(pathValue, artifactRoot)
	artifact := FlowManualReviewArtifact{
		Name:         strings.TrimSuffix(strings.TrimSpace(name), "_path"),
		Kind:         flowManualReviewArtifactKind(name, pathValue),
		Path:         resolvedPath,
		RelativePath: relativePath,
		ContentType:  mime.TypeByExtension(strings.ToLower(filepath.Ext(pathValue))),
		Exists:       exists,
	}
	if artifact.Name == "" {
		artifact.Name = "artifact"
	}
	return artifact
}

func resolveFlowManualReviewArtifactPath(pathValue string, artifactRoot string) (string, string, bool) {
	cleanPath := filepath.Clean(pathValue)
	resolvedPath := cleanPath
	if !filepath.IsAbs(cleanPath) && strings.TrimSpace(artifactRoot) != "" {
		resolvedPath = filepath.Join(artifactRoot, cleanPath)
	}
	if abs, err := filepath.Abs(resolvedPath); err == nil {
		resolvedPath = abs
	}

	relativePath := ""
	if strings.TrimSpace(artifactRoot) != "" {
		if rootAbs, err := filepath.Abs(artifactRoot); err == nil {
			if rel, relErr := filepath.Rel(rootAbs, resolvedPath); relErr == nil && rel != "." && !strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel) {
				relativePath = filepath.ToSlash(rel)
			}
		}
	} else if !filepath.IsAbs(cleanPath) {
		relativePath = filepath.ToSlash(cleanPath)
	}

	_, err := os.Stat(resolvedPath)
	return resolvedPath, relativePath, err == nil
}

func flowManualReviewArtifactKind(name string, pathValue string) string {
	value := strings.ToLower(name + " " + pathValue)
	switch {
	case strings.Contains(value, "screenshot") || strings.Contains(value, ".png") || strings.Contains(value, ".jpg") || strings.Contains(value, ".jpeg") || strings.Contains(value, ".webp"):
		return "screenshot"
	case strings.Contains(value, "html") || strings.HasSuffix(value, ".htm"):
		return "html"
	case strings.Contains(value, "response") || strings.Contains(value, "result") || strings.HasSuffix(value, ".json"):
		return "json"
	default:
		return "artifact"
	}
}

func flowManualReviewMap(value any) (map[string]any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, true
	case map[string]string:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			result[key] = item
		}
		return result, true
	default:
		return nil, false
	}
}

func flowManualReviewString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}
