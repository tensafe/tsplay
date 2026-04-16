package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	flowSavedSessionKindStorageState = "storage_state"
	flowSavedSessionKindProfile      = "persistent_profile"
)

var flowSavedSessionNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

type FlowSavedSession struct {
	Name             string `json:"name"`
	Kind             string `json:"kind"`
	StorageStatePath string `json:"storage_state_path,omitempty"`
	Profile          string `json:"profile,omitempty"`
	Session          string `json:"session,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type FlowSavedSessionSaveOptions struct {
	Name             string
	ArtifactRoot     string
	StorageStateJSON string
	StorageStatePath string
	Profile          string
	Session          string
}

func SaveFlowSavedSession(options FlowSavedSessionSaveOptions) (*FlowSavedSession, error) {
	name, err := normalizeFlowSavedSessionName(options.Name)
	if err != nil {
		return nil, err
	}
	root, err := flowSavedSessionRegistryRoot(options.ArtifactRoot)
	if err != nil {
		return nil, err
	}
	existing, err := LoadFlowSavedSession(name, root)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	session := &FlowSavedSession{
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339Nano),
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
	}
	if existing != nil {
		session.CreatedAt = existing.CreatedAt
	}

	storageStateJSON := strings.TrimSpace(options.StorageStateJSON)
	storageStatePath := strings.TrimSpace(options.StorageStatePath)
	profile := strings.TrimSpace(options.Profile)
	profileSession := strings.TrimSpace(options.Session)

	switch {
	case storageStateJSON != "" || storageStatePath != "":
		if profile != "" || profileSession != "" {
			return nil, fmt.Errorf("save_session accepts either storage_state(_path) or profile/session, not both")
		}
		session.Kind = flowSavedSessionKindStorageState
		content, err := flowSavedSessionContent(root, storageStateJSON, storageStatePath)
		if err != nil {
			return nil, err
		}
		targetPath := flowSavedSessionStorageRelativePath(name)
		targetAbs := filepath.Join(root, targetPath)
		if err := os.MkdirAll(filepath.Dir(targetAbs), 0755); err != nil {
			return nil, fmt.Errorf("create session storage directory: %w", err)
		}
		if err := os.WriteFile(targetAbs, content, 0644); err != nil {
			return nil, fmt.Errorf("write session storage state %q: %w", targetAbs, err)
		}
		session.StorageStatePath = filepath.ToSlash(targetPath)
	case profile != "":
		session.Kind = flowSavedSessionKindProfile
		session.Profile = profile
		session.Session = profileSession
	default:
		return nil, fmt.Errorf("save_session requires storage_state, storage_state_path, or profile")
	}

	if err := writeFlowSavedSession(root, session); err != nil {
		return nil, err
	}
	return session, nil
}

func LoadFlowSavedSession(name string, artifactRoot string) (*FlowSavedSession, error) {
	normalizedName, err := normalizeFlowSavedSessionName(name)
	if err != nil {
		return nil, err
	}
	root, err := flowSavedSessionRegistryRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(flowSavedSessionMetadataPath(root, normalizedName))
	if err != nil {
		return nil, err
	}
	var session FlowSavedSession
	if err := json.Unmarshal(content, &session); err != nil {
		return nil, fmt.Errorf("parse saved session %q: %w", normalizedName, err)
	}
	return &session, nil
}

func ListFlowSavedSessions(artifactRoot string) ([]FlowSavedSession, error) {
	root, err := flowSavedSessionRegistryRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	metadataDir := filepath.Join(root, "sessions", "registry")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return nil, fmt.Errorf("create session registry directory: %w", err)
	}
	entries, err := os.ReadDir(metadataDir)
	if err != nil {
		return nil, fmt.Errorf("read session registry directory: %w", err)
	}

	sessions := make([]FlowSavedSession, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(metadataDir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read saved session metadata %q: %w", entry.Name(), err)
		}
		var session FlowSavedSession
		if err := json.Unmarshal(content, &session); err != nil {
			return nil, fmt.Errorf("parse saved session metadata %q: %w", entry.Name(), err)
		}
		sessions = append(sessions, session)
	}
	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].Name == sessions[j].Name {
			return sessions[i].UpdatedAt > sessions[j].UpdatedAt
		}
		return sessions[i].Name < sessions[j].Name
	})
	return sessions, nil
}

func ResolveFlowSavedSessionBrowserConfig(name string, artifactRoot string) (*FlowBrowserConfig, error) {
	session, err := LoadFlowSavedSession(name, artifactRoot)
	if err != nil {
		return nil, err
	}
	switch session.Kind {
	case flowSavedSessionKindStorageState:
		return &FlowBrowserConfig{StorageState: session.StorageStatePath}, nil
	case flowSavedSessionKindProfile:
		return &FlowBrowserConfig{Persistent: true, Profile: session.Profile, Session: session.Session}, nil
	default:
		return nil, fmt.Errorf("saved session %q uses unsupported kind %q", session.Name, session.Kind)
	}
}

func BuildFlowSavedSessionView(session FlowSavedSession, artifactRoot string) map[string]any {
	view := map[string]any{
		"name":       session.Name,
		"kind":       session.Kind,
		"updated_at": session.UpdatedAt,
		"browser": map[string]any{
			"use_session": session.Name,
		},
		"resolved_browser": map[string]any{},
	}
	if session.CreatedAt != "" {
		view["created_at"] = session.CreatedAt
	}
	switch session.Kind {
	case flowSavedSessionKindStorageState:
		view["storage_state_path"] = session.StorageStatePath
		resolved := map[string]any{
			"storage_state": session.StorageStatePath,
		}
		if artifactRoot != "" {
			if root, err := flowSavedSessionRegistryRoot(artifactRoot); err == nil {
				abs := filepath.Join(root, filepath.FromSlash(session.StorageStatePath))
				if rel, err := filepath.Rel(root, abs); err == nil {
					resolved["storage_state_path"] = filepath.ToSlash(rel)
				}
			}
		}
		view["resolved_browser"] = resolved
	case flowSavedSessionKindProfile:
		view["profile"] = session.Profile
		if session.Session != "" {
			view["session"] = session.Session
		}
		resolved := map[string]any{
			"persistent": true,
			"profile":    session.Profile,
		}
		if session.Session != "" {
			resolved["session"] = session.Session
		}
		view["resolved_browser"] = resolved
	}
	return view
}

func normalizeFlowSavedSessionName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("session name is required")
	}
	if !flowSavedSessionNamePattern.MatchString(name) {
		return "", fmt.Errorf("session name %q is invalid; use letters, digits, dot, underscore, or dash", name)
	}
	return name, nil
}

func flowSavedSessionRegistryRoot(artifactRoot string) (string, error) {
	root := strings.TrimSpace(artifactRoot)
	if root == "" {
		root = DefaultFlowArtifactRoot
	}
	return prepareRuntimeFileRoot(root)
}

func flowSavedSessionStorageRelativePath(name string) string {
	return filepath.ToSlash(filepath.Join("sessions", "storage", name+".json"))
}

func flowSavedSessionMetadataPath(root string, name string) string {
	return filepath.Join(root, "sessions", "registry", name+".json")
}

func writeFlowSavedSession(root string, session *FlowSavedSession) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}
	path := flowSavedSessionMetadataPath(root, session.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create session metadata directory: %w", err)
	}
	content, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal saved session metadata: %w", err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("write saved session metadata %q: %w", path, err)
	}
	return nil
}

func flowSavedSessionContent(root string, storageStateJSON string, storageStatePath string) ([]byte, error) {
	switch {
	case storageStateJSON != "":
		content := []byte(storageStateJSON)
		if !json.Valid(content) {
			return nil, fmt.Errorf("storage_state must be valid JSON")
		}
		return content, nil
	case storageStatePath != "":
		policy := &FlowSecurityPolicy{
			AllowBrowserState: true,
			FileInputRoot:     root,
			FileOutputRoot:    root,
		}
		resolved, err := resolveFlowBrowserStatePath(storageStatePath, flowFileInputPath, policy)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(resolved)
		if err != nil {
			return nil, fmt.Errorf("read storage state %q: %w", resolved, err)
		}
		if !json.Valid(content) {
			return nil, fmt.Errorf("storage_state_path %q is not valid JSON", storageStatePath)
		}
		return content, nil
	default:
		return nil, fmt.Errorf("storage_state or storage_state_path is required")
	}
}
