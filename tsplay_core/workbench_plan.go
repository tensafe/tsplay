package tsplay_core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
)

func BuildWorkbenchTaskPlan(options WorkbenchTaskPlanOptions) (*WorkbenchTaskPlan, error) {
	if err := normalizeWorkbenchTaskPlanOptions(&options); err != nil {
		return nil, err
	}
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
	log.Printf(
		"workbench planner knowledge site=%s pages=%d apis=%d realtime_context=%t intent=%q",
		siteID,
		len(pages),
		len(apis),
		options.realtimeContext != nil,
		workbenchCleanText(intent, 120),
	)

	pageCandidates := rankWorkbenchPages(intent, pages)
	pageCandidates = boostWorkbenchPageCandidatesWithRealtimeContext(pageCandidates, pages, options.realtimeContext)
	apiCandidates := rankWorkbenchAPIs(intent, apis)
	log.Printf(
		"workbench planner ranked site=%s top_page=%s top_api=%s",
		siteID,
		workbenchCandidateDebugLabel(firstWorkbenchCandidate(pageCandidates)),
		workbenchCandidateDebugLabel(firstWorkbenchCandidate(apiCandidates)),
	)
	plan := &WorkbenchTaskPlan{
		SiteID:         siteID,
		Intent:         intent,
		MatchedPages:   pageCandidates,
		MatchedAPIs:    apiCandidates,
		GenerationMode: "local",
	}

	var flow *Flow
	if options.realtimeContext != nil && options.realtimeContext.Observation != nil {
		drafted, err := buildWorkbenchRealtimeContextFlow(site, options.realtimeContext, intent, options.ArtifactRoot)
		if err == nil && drafted != nil {
			flow = drafted
			plan.Strategy = "ui_live"
			plan.Reason = "Used real-time page observation passed by caller, so TSPlay can draft directly from the current page context."
			log.Printf("workbench planner chose ui_live site=%s url=%s", siteID, firstNonEmpty(options.realtimeContext.URL, site.StartURL))
		} else if err != nil {
			addWorkbenchPlanWarning(plan, fmt.Sprintf("Real-time page context could not be drafted directly: %s. Falling back to stored site knowledge.", err.Error()))
			log.Printf("workbench planner realtime_draft_failed site=%s err=%v", siteID, err)
		}
	}
	if flow == nil && len(pageCandidates) > 0 {
		if pageCard := findWorkbenchPageByID(pages, pageCandidates[0].ID); pageCard != nil {
			drafted, err := buildWorkbenchPageFlow(site, *pageCard, intent, options.ArtifactRoot)
			if err == nil && drafted != nil {
				flow = drafted
				plan.Strategy = "ui_first"
				plan.Reason = "Matched a known page card with saved observation data, so TSPlay can draft a selector-aware flow."
				log.Printf("workbench planner chose ui_first_draft site=%s page=%s", siteID, pageCard.ID)
			}
		}
	}
	if flow == nil && len(apiCandidates) > 0 {
		if apiCard := findWorkbenchAPIByID(apis, apiCandidates[0].ID); apiCard != nil {
			flow = buildWorkbenchAPIFallbackFlow(*site, *apiCard, intent)
			plan.Strategy = "api_first"
			plan.Reason = "Matched a readable API card, so the planner generated an API-first flow that reuses browser cookies."
			log.Printf("workbench planner chose api_first site=%s api=%s", siteID, apiCard.ID)
		}
	}
	if flow == nil && len(pageCandidates) > 0 {
		if pageCard := findWorkbenchPageByID(pages, pageCandidates[0].ID); pageCard != nil {
			flow = buildWorkbenchPageFallbackFlow(*site, *pageCard, intent)
			plan.Strategy = "ui_first"
			plan.Reason = "Matched a page card, but no saved observation was available for richer drafting, so the planner generated a navigation-first fallback flow."
			log.Printf("workbench planner chose ui_first_fallback site=%s page=%s", siteID, pageCard.ID)
		}
	}
	if flow == nil {
		plan.Strategy = "needs_input"
		plan.Reason = "No high-confidence page or API candidates were found in the local knowledge store."
		log.Printf("workbench planner no_match site=%s", siteID)
		return plan, nil
	}
	plan.Flow = flow
	plan.FlowName = flow.Name
	flowYAML, err := encodeWorkbenchFlowYAML(flow)
	if err != nil {
		return nil, err
	}
	plan.FlowYAML = flowYAML
	log.Printf("workbench planner flow_ready site=%s flow=%s strategy=%s steps=%d", siteID, flow.Name, plan.Strategy, len(flow.Steps))
	return plan, nil
}

func normalizeWorkbenchTaskPlanOptions(options *WorkbenchTaskPlanOptions) error {
	if options == nil || options.realtimeContext != nil {
		return nil
	}
	realtimeContext, err := parseWorkbenchRealtimeContext(options.RealtimeContext)
	if err != nil {
		return err
	}
	options.realtimeContext = realtimeContext
	return nil
}

func parseWorkbenchRealtimeContext(input *WorkbenchRealtimeContextInput) (*workbenchRealtimeContext, error) {
	if input == nil {
		return nil, nil
	}
	context := &workbenchRealtimeContext{
		URL:   strings.TrimSpace(input.URL),
		Title: strings.TrimSpace(input.Title),
		HTML:  strings.TrimSpace(input.HTML),
	}
	if len(bytes.TrimSpace(input.Observation)) > 0 {
		observation, err := parseWorkbenchRealtimeObservation(input.Observation)
		if err != nil {
			return nil, fmt.Errorf("realtime_context.observation: %w", err)
		}
		context.Observation = observation
		if context.URL == "" && observation != nil {
			context.URL = strings.TrimSpace(observation.URL)
		}
		if context.Title == "" && observation != nil {
			context.Title = strings.TrimSpace(observation.Title)
		}
	}
	if context.URL == "" && context.Title == "" && context.HTML == "" && context.Observation == nil {
		return nil, nil
	}
	return context, nil
}

func parseWorkbenchRealtimeObservation(raw json.RawMessage) (*PageObservation, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return nil, err
		}
		return ParseObservationForDraft(text)
	}
	return ParseObservationForDraft(string(trimmed))
}

func firstWorkbenchCandidate(items []WorkbenchTaskCandidate) *WorkbenchTaskCandidate {
	if len(items) == 0 {
		return nil
	}
	return &items[0]
}

func workbenchCandidateDebugLabel(item *WorkbenchTaskCandidate) string {
	if item == nil {
		return "none"
	}
	parts := []string{firstNonEmpty(item.Kind, "candidate"), firstNonEmpty(item.Label, item.ID)}
	if item.Score > 0 {
		parts = append(parts, fmt.Sprintf("score=%d", item.Score))
	}
	return strings.Join(parts, ":")
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

func boostWorkbenchPageCandidatesWithRealtimeContext(candidates []WorkbenchTaskCandidate, cards []WorkbenchPageCard, realtimeContext *workbenchRealtimeContext) []WorkbenchTaskCandidate {
	if realtimeContext == nil {
		return candidates
	}
	merged := map[string]WorkbenchTaskCandidate{}
	for _, item := range candidates {
		merged[item.ID] = item
	}
	for _, card := range cards {
		boost := scoreWorkbenchRealtimePageBoost(realtimeContext, card)
		if boost <= 0 {
			continue
		}
		candidate, ok := merged[card.ID]
		if !ok {
			candidate = WorkbenchTaskCandidate{
				Kind:  "page",
				ID:    card.ID,
				Label: firstNonEmpty(card.Title, card.NormalizedRoute),
				URL:   card.URL,
			}
		}
		candidate.Score += boost
		if candidate.Label == "" {
			candidate.Label = firstNonEmpty(card.Title, card.NormalizedRoute)
		}
		if candidate.URL == "" {
			candidate.URL = card.URL
		}
		merged[card.ID] = candidate
	}
	items := make([]WorkbenchTaskCandidate, 0, len(merged))
	for _, item := range merged {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Score != items[j].Score {
			return items[i].Score > items[j].Score
		}
		return items[i].Label < items[j].Label
	})
	if len(items) > 5 {
		items = items[:5]
	}
	return items
}

func scoreWorkbenchRealtimePageBoost(realtimeContext *workbenchRealtimeContext, card WorkbenchPageCard) int {
	if realtimeContext == nil {
		return 0
	}
	score := 0
	if realtimeURL := strings.TrimSpace(realtimeContext.URL); realtimeURL != "" {
		switch {
		case sameWorkbenchURLTarget(realtimeURL, card.URL):
			score = maxWorkbenchInt(score, 120)
		case sameWorkbenchURLPath(realtimeURL, card.URL):
			score = maxWorkbenchInt(score, 80)
		case sameWorkbenchRoutePath(realtimeURL, card.NormalizedRoute):
			score = maxWorkbenchInt(score, 70)
		}
	}
	if realtimeTitle := strings.ToLower(strings.TrimSpace(realtimeContext.Title)); realtimeTitle != "" {
		cardTitle := strings.ToLower(strings.TrimSpace(card.Title))
		if cardTitle != "" && (strings.Contains(realtimeTitle, cardTitle) || strings.Contains(cardTitle, realtimeTitle)) {
			score = maxWorkbenchInt(score, 30)
		}
	}
	return score
}

func sameWorkbenchURLTarget(left string, right string) bool {
	return normalizeWorkbenchURLForMatch(left) != "" && normalizeWorkbenchURLForMatch(left) == normalizeWorkbenchURLForMatch(right)
}

func sameWorkbenchURLPath(left string, right string) bool {
	return workbenchURLPathForMatch(left) != "" && workbenchURLPathForMatch(left) == workbenchURLPathForMatch(right)
}

func sameWorkbenchRoutePath(rawURL string, route string) bool {
	path := workbenchURLPathForMatch(rawURL)
	route = normalizeWorkbenchRoutePath(route)
	return path != "" && route != "" && path == route
}

func normalizeWorkbenchURLForMatch(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.TrimRight(strings.TrimSpace(raw), "/")
	}
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	host := strings.ToLower(strings.TrimSpace(parsed.Host))
	path := normalizeWorkbenchRoutePath(parsed.Path)
	if scheme == "" && host == "" && path == "" {
		return strings.TrimRight(strings.TrimSpace(raw), "/")
	}
	return scheme + "://" + host + path
}

func workbenchURLPathForMatch(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return normalizeWorkbenchRoutePath(raw)
	}
	return normalizeWorkbenchRoutePath(parsed.Path)
}

func normalizeWorkbenchRoutePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		if parsed, err := url.Parse(value); err == nil {
			value = parsed.Path
		}
	}
	if value == "" {
		value = "/"
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	if len(value) > 1 {
		value = strings.TrimRight(value, "/")
	}
	return value
}

func maxWorkbenchInt(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value > best {
			best = value
		}
	}
	return best
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

func buildWorkbenchRealtimeContextFlow(site *WorkbenchSiteConfig, realtimeContext *workbenchRealtimeContext, intent string, artifactRoot string) (*Flow, error) {
	if realtimeContext == nil || realtimeContext.Observation == nil {
		return nil, fmt.Errorf("realtime context does not include observation")
	}
	observation := *realtimeContext.Observation
	if strings.TrimSpace(observation.URL) == "" {
		observation.URL = firstNonEmpty(strings.TrimSpace(realtimeContext.URL), strings.TrimSpace(site.StartURL))
	}
	if strings.TrimSpace(observation.Title) == "" {
		observation.Title = strings.TrimSpace(realtimeContext.Title)
	}
	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:       intent,
		URL:          observation.URL,
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
