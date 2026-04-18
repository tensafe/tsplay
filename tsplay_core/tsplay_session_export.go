package tsplay_core

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	flowSavedSessionSnippetFormatAll                 = "all"
	flowSavedSessionSnippetFormatBrowserYAML         = "browser_yaml"
	flowSavedSessionSnippetFormatExpandedBrowserYAML = "expanded_browser_yaml"
	flowSavedSessionSnippetFormatFlowYAML            = "flow_yaml"
	flowSavedSessionSnippetFormatExpandedFlowYAML    = "expanded_flow_yaml"
	flowSavedSessionSnippetFormatBrowserJSON         = "browser_json"
	flowSavedSessionSnippetFormatExpandedBrowserJSON = "expanded_browser_json"
	flowSavedSessionSnippetFormatFlowJSON            = "flow_json"
	flowSavedSessionSnippetFormatExpandedFlowJSON    = "expanded_flow_json"
)

type flowSavedSessionSnippetBrowser struct {
	UseSession   string `json:"use_session,omitempty" yaml:"use_session,omitempty"`
	StorageState string `json:"storage_state,omitempty" yaml:"storage_state,omitempty"`
	Persistent   bool   `json:"persistent,omitempty" yaml:"persistent,omitempty"`
	Profile      string `json:"profile,omitempty" yaml:"profile,omitempty"`
	Session      string `json:"session,omitempty" yaml:"session,omitempty"`
}

type flowSavedSessionSnippetFlow struct {
	SchemaVersion string                         `json:"schema_version" yaml:"schema_version"`
	Name          string                         `json:"name" yaml:"name"`
	Browser       flowSavedSessionSnippetBrowser `json:"browser" yaml:"browser"`
	Steps         []map[string]any               `json:"steps" yaml:"steps"`
}

func BuildFlowSavedSessionFlowSnippet(session FlowSavedSession, artifactRoot string) map[string]any {
	return BuildFlowSavedSessionFlowSnippetForActor(session, artifactRoot, FlowSavedSessionAccessInfo{})
}

func BuildFlowSavedSessionFlowSnippetForActor(session FlowSavedSession, artifactRoot string, actor FlowSavedSessionAccessInfo) map[string]any {
	recommendedBrowserDoc := buildFlowSavedSessionRecommendedBrowserDoc(session)
	expandedBrowserDoc := buildFlowSavedSessionExpandedBrowserDoc(session, artifactRoot, actor)
	recommendedFlowDoc := buildFlowSavedSessionRecommendedFlowDoc(session)
	expandedFlowDoc := buildFlowSavedSessionExpandedFlowDoc(session, artifactRoot, actor)

	return map[string]any{
		"browser":               recommendedBrowserDoc["browser"],
		"expanded_browser":      buildFlowSavedSessionSnippetBrowser(session, artifactRoot, actor),
		"flow":                  recommendedFlowDoc,
		"expanded_flow":         expandedFlowDoc,
		"browser_yaml":          mustMarshalFlowSavedSessionSnippetYAML(recommendedBrowserDoc),
		"expanded_browser_yaml": mustMarshalFlowSavedSessionSnippetYAML(expandedBrowserDoc),
		"flow_yaml":             mustMarshalFlowSavedSessionSnippetYAML(recommendedFlowDoc),
		"expanded_flow_yaml":    mustMarshalFlowSavedSessionSnippetYAML(expandedFlowDoc),
		"browser_json":          mustMarshalFlowSavedSessionSnippetJSON(recommendedBrowserDoc),
		"expanded_browser_json": mustMarshalFlowSavedSessionSnippetJSON(expandedBrowserDoc),
		"flow_json":             mustMarshalFlowSavedSessionSnippetJSON(recommendedFlowDoc),
		"expanded_flow_json":    mustMarshalFlowSavedSessionSnippetJSON(expandedFlowDoc),
		"recommended_usage":     "Prefer browser.use_session for reusable flows; use the expanded browser block only when you need the resolved runtime config inline.",
	}
}

func ExportFlowSavedSessionFlowSnippet(session FlowSavedSession, artifactRoot string, format string) (map[string]any, error) {
	return ExportFlowSavedSessionFlowSnippetForActor(session, artifactRoot, format, FlowSavedSessionAccessInfo{})
}

func ExportFlowSavedSessionFlowSnippetForActor(session FlowSavedSession, artifactRoot string, format string, actor FlowSavedSessionAccessInfo) (map[string]any, error) {
	bundle := BuildFlowSavedSessionFlowSnippetForActor(session, artifactRoot, actor)
	normalizedFormat, err := normalizeFlowSavedSessionSnippetFormat(format)
	if err != nil {
		return nil, err
	}
	if normalizedFormat == flowSavedSessionSnippetFormatAll {
		return map[string]any{
			"format":            normalizedFormat,
			"snippets":          bundle,
			"recommended_usage": bundle["recommended_usage"],
		}, nil
	}

	target, variant, encoding, contentType, err := describeFlowSavedSessionSnippetFormat(normalizedFormat)
	if err != nil {
		return nil, err
	}
	snippet, ok := bundle[normalizedFormat].(string)
	if !ok {
		return nil, fmt.Errorf("snippet %q is unavailable", normalizedFormat)
	}
	dataKey := strings.TrimSuffix(normalizedFormat, "_yaml")
	dataKey = strings.TrimSuffix(dataKey, "_json")
	snippetData, ok := bundle[dataKey]
	if !ok {
		return nil, fmt.Errorf("snippet data %q is unavailable", dataKey)
	}

	return map[string]any{
		"format":            normalizedFormat,
		"target":            target,
		"variant":           variant,
		"encoding":          encoding,
		"content_type":      contentType,
		"snippet":           snippet,
		"snippet_data":      snippetData,
		"recommended_usage": bundle["recommended_usage"],
	}, nil
}

func buildFlowSavedSessionRecommendedBrowserDoc(session FlowSavedSession) map[string]flowSavedSessionSnippetBrowser {
	return map[string]flowSavedSessionSnippetBrowser{
		"browser": {
			UseSession: session.Name,
		},
	}
}

func buildFlowSavedSessionExpandedBrowserDoc(session FlowSavedSession, artifactRoot string, actor FlowSavedSessionAccessInfo) map[string]flowSavedSessionSnippetBrowser {
	return map[string]flowSavedSessionSnippetBrowser{
		"browser": buildFlowSavedSessionSnippetBrowser(session, artifactRoot, actor),
	}
}

func buildFlowSavedSessionRecommendedFlowDoc(session FlowSavedSession) flowSavedSessionSnippetFlow {
	return flowSavedSessionSnippetFlow{
		SchemaVersion: "1",
		Name:          fmt.Sprintf("reuse_%s_session", session.Name),
		Browser: flowSavedSessionSnippetBrowser{
			UseSession: session.Name,
		},
		Steps: []map[string]any{},
	}
}

func buildFlowSavedSessionExpandedFlowDoc(session FlowSavedSession, artifactRoot string, actor FlowSavedSessionAccessInfo) flowSavedSessionSnippetFlow {
	return flowSavedSessionSnippetFlow{
		SchemaVersion: "1",
		Name:          fmt.Sprintf("reuse_%s_session_expanded", session.Name),
		Browser:       buildFlowSavedSessionSnippetBrowser(session, artifactRoot, actor),
		Steps:         []map[string]any{},
	}
}

func buildFlowSavedSessionSnippetBrowser(session FlowSavedSession, artifactRoot string, actor FlowSavedSessionAccessInfo) flowSavedSessionSnippetBrowser {
	config, err := ResolveFlowSavedSessionBrowserConfig(session.Name, artifactRoot, actor)
	if err != nil || config == nil {
		return flowSavedSessionSnippetBrowser{}
	}
	return flowSavedSessionSnippetBrowser{
		StorageState: config.StorageState,
		Persistent:   config.Persistent,
		Profile:      config.Profile,
		Session:      config.Session,
	}
}

func normalizeFlowSavedSessionSnippetFormat(format string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	switch normalized {
	case "", flowSavedSessionSnippetFormatAll:
		return flowSavedSessionSnippetFormatAll, nil
	case "browser":
		return flowSavedSessionSnippetFormatBrowserYAML, nil
	case "expanded_browser":
		return flowSavedSessionSnippetFormatExpandedBrowserYAML, nil
	case "flow":
		return flowSavedSessionSnippetFormatFlowYAML, nil
	case "expanded_flow":
		return flowSavedSessionSnippetFormatExpandedFlowYAML, nil
	case flowSavedSessionSnippetFormatBrowserYAML,
		flowSavedSessionSnippetFormatExpandedBrowserYAML,
		flowSavedSessionSnippetFormatFlowYAML,
		flowSavedSessionSnippetFormatExpandedFlowYAML,
		flowSavedSessionSnippetFormatBrowserJSON,
		flowSavedSessionSnippetFormatExpandedBrowserJSON,
		flowSavedSessionSnippetFormatFlowJSON,
		flowSavedSessionSnippetFormatExpandedFlowJSON:
		return normalized, nil
	default:
		return "", fmt.Errorf("unsupported format %q; use one of all, browser, expanded_browser, flow, expanded_flow, browser_json, expanded_browser_json, flow_json, or expanded_flow_json", format)
	}
}

func describeFlowSavedSessionSnippetFormat(format string) (target string, variant string, encoding string, contentType string, err error) {
	switch format {
	case flowSavedSessionSnippetFormatBrowserYAML:
		return "browser", "recommended", "yaml", "application/yaml", nil
	case flowSavedSessionSnippetFormatExpandedBrowserYAML:
		return "browser", "expanded", "yaml", "application/yaml", nil
	case flowSavedSessionSnippetFormatFlowYAML:
		return "flow", "recommended", "yaml", "application/yaml", nil
	case flowSavedSessionSnippetFormatExpandedFlowYAML:
		return "flow", "expanded", "yaml", "application/yaml", nil
	case flowSavedSessionSnippetFormatBrowserJSON:
		return "browser", "recommended", "json", "application/json", nil
	case flowSavedSessionSnippetFormatExpandedBrowserJSON:
		return "browser", "expanded", "json", "application/json", nil
	case flowSavedSessionSnippetFormatFlowJSON:
		return "flow", "recommended", "json", "application/json", nil
	case flowSavedSessionSnippetFormatExpandedFlowJSON:
		return "flow", "expanded", "json", "application/json", nil
	default:
		return "", "", "", "", fmt.Errorf("unsupported format %q", format)
	}
}

func mustMarshalFlowSavedSessionSnippetYAML(value any) string {
	encoded, err := yaml.Marshal(value)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(encoded))
}

func mustMarshalFlowSavedSessionSnippetJSON(value any) string {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return ""
	}
	return string(encoded)
}
