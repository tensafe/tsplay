package tsplay_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type workbenchKnowledgeBundle[T any] struct {
	Items []T `json:"items"`
}

func SaveWorkbenchSiteConfig(config WorkbenchSiteConfig, artifactRoot string) (*WorkbenchSiteConfig, error) {
	siteID := normalizeWorkbenchSiteID(config.SiteID)
	if siteID == "" {
		return nil, fmt.Errorf("site_id is required")
	}
	startURL := strings.TrimSpace(config.StartURL)
	if startURL == "" {
		return nil, fmt.Errorf("start_url is required")
	}
	parsed, err := url.Parse(startURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("start_url %q must be a valid absolute URL", startURL)
	}

	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339Nano)
	existing, _ := LoadWorkbenchSiteConfig(siteID, artifactRoot)
	saved := &WorkbenchSiteConfig{
		SiteID:         siteID,
		Name:           strings.TrimSpace(config.Name),
		StartURL:       startURL,
		AllowedDomains: normalizeAllowedDomains(config.AllowedDomains, parsed.Hostname()),
		SessionName:    strings.TrimSpace(config.SessionName),
		ProviderID:     strings.TrimSpace(config.ProviderID),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if existing != nil {
		saved.CreatedAt = firstNonEmpty(existing.CreatedAt, now)
	}
	if saved.Name == "" {
		saved.Name = siteID
	}
	if err := os.MkdirAll(filepath.Join(root, "sites"), 0755); err != nil {
		return nil, fmt.Errorf("create workbench sites directory: %w", err)
	}
	if err := writeWorkbenchJSON(filepath.Join(root, "sites", siteID+".json"), saved); err != nil {
		return nil, err
	}
	return saved, nil
}

func LoadWorkbenchSiteConfig(siteID string, artifactRoot string) (*WorkbenchSiteConfig, error) {
	normalized := normalizeWorkbenchSiteID(siteID)
	if normalized == "" {
		return nil, fmt.Errorf("site_id is required")
	}
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(filepath.Join(root, "sites", normalized+".json"))
	if err != nil {
		return nil, err
	}
	var config WorkbenchSiteConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("parse site config %q: %w", normalized, err)
	}
	return &config, nil
}

func ListWorkbenchSiteConfigs(artifactRoot string) ([]WorkbenchSiteConfig, error) {
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "sites")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create workbench sites directory: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read workbench sites directory: %w", err)
	}
	items := make([]WorkbenchSiteConfig, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read site config %q: %w", entry.Name(), err)
		}
		var config WorkbenchSiteConfig
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("parse site config %q: %w", entry.Name(), err)
		}
		items = append(items, config)
	}
	sort.Slice(items, func(i, j int) bool {
		left := firstNonEmpty(items[i].UpdatedAt, items[i].CreatedAt)
		right := firstNonEmpty(items[j].UpdatedAt, items[j].CreatedAt)
		if left != right {
			return left > right
		}
		return items[i].SiteID < items[j].SiteID
	})
	return items, nil
}

func SaveWorkbenchProviderConfig(config WorkbenchProviderConfig, artifactRoot string) (*WorkbenchProviderConfig, error) {
	providerID := normalizeWorkbenchSiteID(config.ProviderID)
	if providerID == "" {
		return nil, fmt.Errorf("provider_id is required")
	}
	if providerID == workbenchAutoProviderID {
		return nil, fmt.Errorf("provider_id %q is reserved", workbenchAutoProviderID)
	}
	providerType := normalizeWorkbenchProviderType(config.Type)
	if providerType == "" {
		return nil, fmt.Errorf("provider type %q is not supported", config.Type)
	}
	baseURL := normalizeWorkbenchProviderBaseURL(config.BaseURL)
	if baseURL != "" {
		parsed, err := url.Parse(baseURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return nil, fmt.Errorf("base_url %q must be a valid absolute URL", config.BaseURL)
		}
	}

	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339Nano)
	existing, _ := LoadWorkbenchProviderConfig(providerID, artifactRoot)
	saved := &WorkbenchProviderConfig{
		ProviderID:   providerID,
		Name:         strings.TrimSpace(config.Name),
		Type:         providerType,
		BaseURL:      baseURL,
		Model:        strings.TrimSpace(config.Model),
		APIKey:       strings.TrimSpace(config.APIKey),
		APIKeyEnv:    strings.TrimSpace(config.APIKeyEnv),
		SystemPrompt: strings.TrimSpace(config.SystemPrompt),
		Enabled:      config.Enabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if existing != nil {
		saved.CreatedAt = firstNonEmpty(existing.CreatedAt, now)
		if saved.APIKey == "" {
			saved.APIKey = existing.APIKey
		}
	}
	if saved.Name == "" {
		saved.Name = providerID
	}
	if err := os.MkdirAll(filepath.Join(root, "providers"), 0755); err != nil {
		return nil, fmt.Errorf("create workbench providers directory: %w", err)
	}
	if err := writeWorkbenchJSON(filepath.Join(root, "providers", providerID+".json"), saved); err != nil {
		return nil, err
	}
	return saved, nil
}

func LoadWorkbenchProviderConfig(providerID string, artifactRoot string) (*WorkbenchProviderConfig, error) {
	normalized := normalizeWorkbenchSiteID(providerID)
	if normalized == "" {
		return nil, fmt.Errorf("provider_id is required")
	}
	if normalized == workbenchAutoProviderID {
		auto := detectWorkbenchAutoProviderConfig()
		if auto == nil {
			return nil, fmt.Errorf("provider %q is not available", providerID)
		}
		return auto, nil
	}
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(filepath.Join(root, "providers", normalized+".json"))
	if err != nil {
		return nil, err
	}
	var config WorkbenchProviderConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("parse provider config %q: %w", normalized, err)
	}
	return &config, nil
}

func ListWorkbenchProviderConfigs(artifactRoot string) ([]WorkbenchProviderConfig, error) {
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "providers")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create workbench providers directory: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read workbench providers directory: %w", err)
	}
	items := make([]WorkbenchProviderConfig, 0, len(entries)+1)
	if auto := detectWorkbenchAutoProviderConfig(); auto != nil {
		items = append(items, *auto)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read provider config %q: %w", entry.Name(), err)
		}
		var config WorkbenchProviderConfig
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("parse provider config %q: %w", entry.Name(), err)
		}
		items = append(items, config)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].ProviderID == workbenchAutoProviderID || items[j].ProviderID == workbenchAutoProviderID {
			return items[i].ProviderID == workbenchAutoProviderID
		}
		left := firstNonEmpty(items[i].UpdatedAt, items[i].CreatedAt)
		right := firstNonEmpty(items[j].UpdatedAt, items[j].CreatedAt)
		if left != right {
			return left > right
		}
		return items[i].ProviderID < items[j].ProviderID
	})
	return items, nil
}

func UpsertWorkbenchPageCards(siteID string, artifactRoot string, cards []WorkbenchPageCard) ([]WorkbenchPageCard, error) {
	current, err := ListWorkbenchPageCards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	merged := map[string]WorkbenchPageCard{}
	for _, item := range current {
		merged[item.ID] = item
	}
	for _, item := range cards {
		if strings.TrimSpace(item.ID) == "" {
			continue
		}
		merged[item.ID] = item
	}
	items := make([]WorkbenchPageCard, 0, len(merged))
	for _, item := range merged {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].NormalizedRoute != items[j].NormalizedRoute {
			return items[i].NormalizedRoute < items[j].NormalizedRoute
		}
		return items[i].URL < items[j].URL
	})
	if err := saveWorkbenchKnowledgeBundle(filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "pages.json"), items); err != nil {
		return nil, err
	}
	return items, nil
}

func UpsertWorkbenchAPICards(siteID string, artifactRoot string, cards []WorkbenchAPICard) ([]WorkbenchAPICard, error) {
	current, err := ListWorkbenchAPICards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	merged := map[string]WorkbenchAPICard{}
	for _, item := range current {
		merged[item.ID] = item
	}
	for _, item := range cards {
		if strings.TrimSpace(item.ID) == "" {
			continue
		}
		existing, ok := merged[item.ID]
		if ok {
			if existing.RequestSchema == nil && item.RequestSchema != nil {
				existing.RequestSchema = item.RequestSchema
			}
			if existing.ResponseSchema == nil && item.ResponseSchema != nil {
				existing.ResponseSchema = item.ResponseSchema
			}
			if existing.SemanticName == "" {
				existing.SemanticName = item.SemanticName
			}
			if existing.TriggerRoute == "" {
				existing.TriggerRoute = item.TriggerRoute
			}
			if existing.TriggerAction == "" {
				existing.TriggerAction = item.TriggerAction
			}
			if existing.Risk == "" {
				existing.Risk = item.Risk
			}
			if existing.OperationType == "" {
				existing.OperationType = item.OperationType
			}
			if existing.Status == 0 {
				existing.Status = item.Status
			}
			if existing.ContentType == "" {
				existing.ContentType = item.ContentType
			}
			if existing.URL == "" {
				existing.URL = item.URL
			}
			if item.UpdatedAt != "" {
				existing.UpdatedAt = item.UpdatedAt
			}
			merged[item.ID] = existing
			continue
		}
		merged[item.ID] = item
	}
	items := make([]WorkbenchAPICard, 0, len(merged))
	for _, item := range merged {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Method != items[j].Method {
			return items[i].Method < items[j].Method
		}
		return items[i].PathTemplate < items[j].PathTemplate
	})
	if err := saveWorkbenchKnowledgeBundle(filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "apis.json"), items); err != nil {
		return nil, err
	}
	return items, nil
}

func UpsertWorkbenchEntityCards(siteID string, artifactRoot string, cards []WorkbenchEntityCard) ([]WorkbenchEntityCard, error) {
	current, err := ListWorkbenchEntityCards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	merged := map[string]WorkbenchEntityCard{}
	for _, item := range current {
		merged[item.ID] = item
	}
	for _, item := range cards {
		if strings.TrimSpace(item.ID) == "" {
			continue
		}
		merged[item.ID] = item
	}
	items := make([]WorkbenchEntityCard, 0, len(merged))
	for _, item := range merged {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	if err := saveWorkbenchKnowledgeBundle(filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "entities.json"), items); err != nil {
		return nil, err
	}
	return items, nil
}

func ListWorkbenchPageCards(siteID string, artifactRoot string) ([]WorkbenchPageCard, error) {
	return loadWorkbenchKnowledgeBundle[WorkbenchPageCard](filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "pages.json"))
}

func ListWorkbenchAPICards(siteID string, artifactRoot string) ([]WorkbenchAPICard, error) {
	return loadWorkbenchKnowledgeBundle[WorkbenchAPICard](filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "apis.json"))
}

func ListWorkbenchEntityCards(siteID string, artifactRoot string) ([]WorkbenchEntityCard, error) {
	return loadWorkbenchKnowledgeBundle[WorkbenchEntityCard](filepath.Join(workbenchKnowledgeSiteDir(siteID, artifactRoot), "entities.json"))
}

func SaveWorkbenchExploreResult(result WorkbenchExploreResult, artifactRoot string) (*WorkbenchExploreResult, error) {
	siteID := normalizeWorkbenchSiteID(result.Site.SiteID)
	if siteID == "" {
		return nil, fmt.Errorf("site_id is required")
	}
	if strings.TrimSpace(result.RunID) == "" || strings.TrimSpace(result.RunRoot) == "" {
		return nil, fmt.Errorf("run_id and run_root are required")
	}
	if err := os.MkdirAll(result.RunRoot, 0755); err != nil {
		return nil, fmt.Errorf("create explore run directory: %w", err)
	}
	copyResult := result
	copyResult.Site.SiteID = siteID
	if err := writeWorkbenchJSON(filepath.Join(result.RunRoot, "result.json"), copyResult); err != nil {
		return nil, err
	}
	if _, err := UpsertWorkbenchPageCards(siteID, artifactRoot, result.Pages); err != nil {
		return nil, err
	}
	if _, err := UpsertWorkbenchAPICards(siteID, artifactRoot, result.APIs); err != nil {
		return nil, err
	}
	if _, err := UpsertWorkbenchEntityCards(siteID, artifactRoot, result.Entities); err != nil {
		return nil, err
	}
	return &copyResult, nil
}

func workbenchRoot(artifactRoot string) (string, error) {
	root, err := prepareRuntimeFileRoot(firstNonEmpty(strings.TrimSpace(artifactRoot), DefaultFlowArtifactRoot))
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "workbench"), nil
}

func workbenchSiteRunRoot(siteID string, artifactRoot string, runID string) (string, error) {
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "sites_data", normalizeWorkbenchSiteID(siteID), "runs", sanitizeArtifactSegment(runID)), nil
}

func workbenchKnowledgeSiteDir(siteID string, artifactRoot string) string {
	root, err := workbenchRoot(artifactRoot)
	if err != nil {
		return filepath.Join(DefaultFlowArtifactRoot, "workbench", "knowledge", normalizeWorkbenchSiteID(siteID))
	}
	return filepath.Join(root, "knowledge", normalizeWorkbenchSiteID(siteID))
}

func normalizeWorkbenchSiteID(value string) string {
	value = sanitizeArtifactSegment(value)
	if value == "flow" {
		return ""
	}
	return value
}

func normalizeAllowedDomains(domains []string, fallback string) []string {
	normalized := make([]string, 0, len(domains)+1)
	seen := map[string]struct{}{}
	appendDomain := func(value string) {
		value = strings.ToLower(strings.TrimSpace(value))
		value = strings.TrimPrefix(value, "https://")
		value = strings.TrimPrefix(value, "http://")
		value = strings.TrimSuffix(value, "/")
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	for _, domain := range domains {
		appendDomain(domain)
	}
	appendDomain(fallback)
	sort.Strings(normalized)
	return normalized
}

func writeWorkbenchJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create workbench directory %q: %w", filepath.Dir(path), err)
	}
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal workbench json %q: %w", path, err)
	}
	if err := os.WriteFile(path, encoded, 0644); err != nil {
		return fmt.Errorf("write workbench json %q: %w", path, err)
	}
	return nil
}

func saveWorkbenchKnowledgeBundle[T any](path string, items []T) error {
	return writeWorkbenchJSON(path, workbenchKnowledgeBundle[T]{Items: items})
}

func loadWorkbenchKnowledgeBundle[T any](path string) ([]T, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bundle workbenchKnowledgeBundle[T]
	if err := json.Unmarshal(content, &bundle); err != nil {
		return nil, fmt.Errorf("parse knowledge bundle %q: %w", path, err)
	}
	if bundle.Items == nil {
		return []T{}, nil
	}
	return bundle.Items, nil
}
