package tsplay_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func BuildWorkbenchTaskPlan(options WorkbenchTaskPlanOptions) (*WorkbenchTaskPlan, error) {
	siteID := normalizeWorkbenchSiteID(options.SiteID)
	if siteID == "" {
		return nil, fmt.Errorf("site_id is required")
	}
	intent := strings.TrimSpace(options.Intent)
	if intent == "" {
		return nil, fmt.Errorf("intent is required")
	}
	site, err := LoadWorkbenchSiteConfig(siteID, options.ArtifactRoot)
	if err != nil {
		return nil, err
	}
	pages, err := ListWorkbenchPageCards(siteID, options.ArtifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	apis, err := ListWorkbenchAPICards(siteID, options.ArtifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	pageCandidates := rankWorkbenchPages(intent, pages)
	apiCandidates := rankWorkbenchAPIs(intent, apis)
	plan := &WorkbenchTaskPlan{
		SiteID:       siteID,
		Intent:       intent,
		MatchedPages: pageCandidates,
		MatchedAPIs:  apiCandidates,
	}

	var flow *Flow
	if len(pageCandidates) > 0 {
		if pageCard := findWorkbenchPageByID(pages, pageCandidates[0].ID); pageCard != nil {
			drafted, err := buildWorkbenchPageFlow(site, *pageCard, intent, options.ArtifactRoot)
			if err == nil && drafted != nil {
				flow = drafted
				plan.Strategy = "ui_first"
				plan.Reason = "Matched a known page card with saved observation data, so TSPlay can draft a selector-aware flow."
			}
		}
	}
	if flow == nil && len(apiCandidates) > 0 {
		if apiCard := findWorkbenchAPIByID(apis, apiCandidates[0].ID); apiCard != nil {
			flow = buildWorkbenchAPIFallbackFlow(*site, *apiCard, intent)
			plan.Strategy = "api_first"
			plan.Reason = "Matched a readable API card, so the planner generated an API-first flow that reuses browser cookies."
		}
	}
	if flow == nil && len(pageCandidates) > 0 {
		if pageCard := findWorkbenchPageByID(pages, pageCandidates[0].ID); pageCard != nil {
			flow = buildWorkbenchPageFallbackFlow(*site, *pageCard, intent)
			plan.Strategy = "ui_first"
			plan.Reason = "Matched a page card, but no saved observation was available for richer drafting, so the planner generated a navigation-first fallback flow."
		}
	}
	if flow == nil {
		plan.Strategy = "needs_input"
		plan.Reason = "No high-confidence page or API candidates were found in the local knowledge store."
		return plan, nil
	}
	plan.Flow = flow
	plan.FlowName = flow.Name
	flowYAML, err := encodeWorkbenchFlowYAML(flow)
	if err != nil {
		return nil, err
	}
	plan.FlowYAML = flowYAML
	return plan, nil
}

func rankWorkbenchPages(intent string, cards []WorkbenchPageCard) []WorkbenchTaskCandidate {
	candidates := make([]WorkbenchTaskCandidate, 0, len(cards))
	for _, card := range cards {
		score := scoreWorkbenchText(intent, flattenWorkbenchPageCard(card))
		if score == 0 {
			continue
		}
		candidates = append(candidates, WorkbenchTaskCandidate{
			Kind:  "page",
			ID:    card.ID,
			Label: firstNonEmpty(card.Title, card.NormalizedRoute),
			URL:   card.URL,
			Score: score,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		return candidates[i].Label < candidates[j].Label
	})
	if len(candidates) > 5 {
		candidates = candidates[:5]
	}
	return candidates
}

func rankWorkbenchAPIs(intent string, cards []WorkbenchAPICard) []WorkbenchTaskCandidate {
	candidates := make([]WorkbenchTaskCandidate, 0, len(cards))
	for _, card := range cards {
		score := scoreWorkbenchText(intent, flattenWorkbenchAPICard(card))
		if score == 0 {
			continue
		}
		candidates = append(candidates, WorkbenchTaskCandidate{
			Kind:   "api",
			ID:     card.ID,
			Label:  firstNonEmpty(card.SemanticName, card.PathTemplate),
			Method: card.Method,
			Path:   card.PathTemplate,
			Score:  score,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		if candidates[i].Method != candidates[j].Method {
			return candidates[i].Method < candidates[j].Method
		}
		return candidates[i].Path < candidates[j].Path
	})
	if len(candidates) > 5 {
		candidates = candidates[:5]
	}
	return candidates
}

func buildWorkbenchPageFlow(site *WorkbenchSiteConfig, card WorkbenchPageCard, intent string, artifactRoot string) (*Flow, error) {
	if strings.TrimSpace(card.ObservationPath) == "" {
		return nil, fmt.Errorf("page %q has no observation_path", card.ID)
	}
	content, err := os.ReadFile(card.ObservationPath)
	if err != nil {
		return nil, err
	}
	var observation PageObservation
	if err := json.Unmarshal(content, &observation); err != nil {
		return nil, err
	}
	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:       intent,
		URL:          card.URL,
		FlowName:     buildDraftFlowName("", intent, &observation),
		ArtifactRoot: artifactRoot,
		Observation:  &observation,
	})
	if err != nil {
		return nil, err
	}
	if draft == nil || draft.Flow == nil {
		return nil, fmt.Errorf("draft flow returned no flow")
	}
	if draft.Flow.Browser == nil {
		draft.Flow.Browser = &FlowBrowserConfig{}
	}
	if site != nil && strings.TrimSpace(site.SessionName) != "" {
		draft.Flow.Browser.UseSession = site.SessionName
	}
	return draft.Flow, nil
}

func buildWorkbenchPageFallbackFlow(site WorkbenchSiteConfig, card WorkbenchPageCard, intent string) *Flow {
	flow := &Flow{
		SchemaVersion: CurrentFlowSchemaVersion,
		Name:          sanitizeArtifactSegment(intent),
		Description:   fmt.Sprintf("Fallback page-first flow for intent %q.", intent),
		Vars: map[string]any{
			"target_url": card.URL,
		},
		Steps: []FlowStep{
			{
				Name:   "open matched page",
				Action: "navigate",
				URL:    "{{target_url}}",
			},
		},
	}
	if site.SessionName != "" {
		flow.Browser = &FlowBrowserConfig{UseSession: site.SessionName}
	}
	if len(card.Tables) > 0 {
		tableSelector := firstNonEmpty(card.Tables[0].Selector, "table")
		flow.Steps = append(flow.Steps, FlowStep{
			Name:     "capture first table",
			Action:   "capture_table",
			Selector: tableSelector,
			SaveAs:   "table_rows",
		})
	}
	return flow
}

func buildWorkbenchAPIFallbackFlow(site WorkbenchSiteConfig, card WorkbenchAPICard, intent string) *Flow {
	requestPayload := map[string]any{}
	if schema, ok := card.RequestSchema.(map[string]any); ok {
		for key := range schema {
			requestPayload[key] = fmt.Sprintf("TODO_%s", sanitizeArtifactSegment(key))
		}
	}
	fullURL := card.URL
	if strings.TrimSpace(fullURL) == "" {
		fullURL = buildWorkbenchAbsoluteURL(site.StartURL, card.PathTemplate)
	}
	flow := &Flow{
		SchemaVersion: CurrentFlowSchemaVersion,
		Name:          sanitizeArtifactSegment(intent),
		Description:   fmt.Sprintf("API-first fallback flow for intent %q using %s %s.", intent, card.Method, card.PathTemplate),
		Vars: map[string]any{
			"target_url": fullURL,
		},
		Steps: []FlowStep{},
	}
	if site.SessionName != "" {
		flow.Browser = &FlowBrowserConfig{UseSession: site.SessionName}
		flow.Steps = append(flow.Steps, FlowStep{
			Name:   "open site for browser cookies",
			Action: "navigate",
			URL:    site.StartURL,
		})
	}
	step := FlowStep{
		Name:   "call matched api",
		Action: "http_request",
		URL:    "{{target_url}}",
		With: map[string]any{
			"method":                 card.Method,
			"use_browser_cookies":    true,
			"use_browser_referer":    true,
			"use_browser_user_agent": true,
		},
		SaveAs: "api_result",
	}
	if len(requestPayload) > 0 && strings.EqualFold(card.Method, "GET") {
		step.With["query"] = requestPayload
	} else if len(requestPayload) > 0 {
		step.With["json"] = requestPayload
	}
	flow.Steps = append(flow.Steps, step)
	return flow
}

func buildWorkbenchAbsoluteURL(startURL string, pathTemplate string) string {
	base, err := urlFromString(startURL)
	if err != nil {
		return pathTemplate
	}
	if strings.HasPrefix(pathTemplate, "http://") || strings.HasPrefix(pathTemplate, "https://") {
		return pathTemplate
	}
	base.Path = pathTemplate
	base.RawPath = pathTemplate
	base.RawQuery = ""
	base.Fragment = ""
	return base.String()
}

func findWorkbenchPageByID(cards []WorkbenchPageCard, id string) *WorkbenchPageCard {
	for i := range cards {
		if cards[i].ID == id {
			return &cards[i]
		}
	}
	return nil
}

func findWorkbenchAPIByID(cards []WorkbenchAPICard, id string) *WorkbenchAPICard {
	for i := range cards {
		if cards[i].ID == id {
			return &cards[i]
		}
	}
	return nil
}

func urlFromString(rawURL string) (*url.URL, error) {
	return url.Parse(strings.TrimSpace(rawURL))
}
