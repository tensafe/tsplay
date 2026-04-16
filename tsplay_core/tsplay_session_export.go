package tsplay_core

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type flowSavedSessionSnippetBrowser struct {
	UseSession   string `yaml:"use_session,omitempty"`
	StorageState string `yaml:"storage_state,omitempty"`
	Persistent   bool   `yaml:"persistent,omitempty"`
	Profile      string `yaml:"profile,omitempty"`
	Session      string `yaml:"session,omitempty"`
}

type flowSavedSessionSnippetFlow struct {
	SchemaVersion string                         `yaml:"schema_version"`
	Name          string                         `yaml:"name"`
	Browser       flowSavedSessionSnippetBrowser `yaml:"browser"`
	Steps         []map[string]any               `yaml:"steps"`
}

func BuildFlowSavedSessionFlowSnippet(session FlowSavedSession, artifactRoot string) map[string]any {
	recommendedBrowser := flowSavedSessionSnippetBrowser{
		UseSession: session.Name,
	}
	expandedBrowser := buildFlowSavedSessionSnippetBrowser(session, artifactRoot)

	return map[string]any{
		"browser":          map[string]any{"use_session": session.Name},
		"expanded_browser": buildFlowSavedSessionResolvedBrowser(session, artifactRoot),
		"browser_yaml": mustMarshalFlowSavedSessionSnippetYAML(struct {
			Browser flowSavedSessionSnippetBrowser `yaml:"browser"`
		}{Browser: recommendedBrowser}),
		"expanded_browser_yaml": mustMarshalFlowSavedSessionSnippetYAML(struct {
			Browser flowSavedSessionSnippetBrowser `yaml:"browser"`
		}{Browser: expandedBrowser}),
		"flow_yaml": mustMarshalFlowSavedSessionSnippetYAML(flowSavedSessionSnippetFlow{
			SchemaVersion: "1",
			Name:          fmt.Sprintf("reuse_%s_session", session.Name),
			Browser:       recommendedBrowser,
			Steps:         []map[string]any{},
		}),
		"expanded_flow_yaml": mustMarshalFlowSavedSessionSnippetYAML(flowSavedSessionSnippetFlow{
			SchemaVersion: "1",
			Name:          fmt.Sprintf("reuse_%s_session_expanded", session.Name),
			Browser:       expandedBrowser,
			Steps:         []map[string]any{},
		}),
		"recommended_usage": "Prefer browser.use_session for reusable flows; use the expanded browser block only when you need the resolved runtime config inline.",
	}
}

func buildFlowSavedSessionSnippetBrowser(session FlowSavedSession, artifactRoot string) flowSavedSessionSnippetBrowser {
	config, err := ResolveFlowSavedSessionBrowserConfig(session.Name, artifactRoot)
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

func mustMarshalFlowSavedSessionSnippetYAML(value any) string {
	encoded, err := yaml.Marshal(value)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(encoded))
}
