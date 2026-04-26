package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	defaultWorkbenchPlanSystemPrompt            = "You generate TSPlay Flow YAML for TSPlay Workbench. Return only one valid YAML document with no markdown fences or prose. Use schema_version \"1\", only supported TSPlay actions and fields, prefer structured actions over lua, prefer stable selectors from the provided observation context, preserve or reuse browser.use_session when available, and place unresolved user inputs in flow.vars as TODO values referenced through {{var_name}}. browser.use_session must be a saved session name string when present; never output boolean true/false for browser.use_session."
	defaultWorkbenchPlanPageContextLimit        = 3
	defaultWorkbenchPlanAPIContextLimit         = 3
	defaultWorkbenchPlanEntityContextLimit      = 5
	defaultWorkbenchPlanObservationElementLimit = 16
	defaultWorkbenchPlanObservationContentLimit = 12
	defaultWorkbenchPlanExampleLimit            = 2
	defaultWorkbenchPlanTextListLimit           = 6
	defaultWorkbenchPlanHTMLExcerptLimit        = 4000
)

type workbenchTaskPlanningKnowledge struct {
	Site     *WorkbenchSiteConfig
	Pages    []WorkbenchPageCard
	APIs     []WorkbenchAPICard
	Entities []WorkbenchEntityCard
}

func BuildWorkbenchProviderTaskPlan(
	options WorkbenchTaskPlanOptions,
	basePlan *WorkbenchTaskPlan,
	providerConfig WorkbenchProviderConfig,
) (*WorkbenchTaskPlan, WorkbenchProviderView, string, error) {
	if err := normalizeWorkbenchTaskPlanOptions(&options); err != nil {
		return nil, WorkbenchProviderView{}, "", err
	}
	knowledge, err := loadWorkbenchTaskPlanningKnowledge(options.SiteID, options.ArtifactRoot)
	if err != nil {
		return nil, WorkbenchProviderView{}, "", err
	}

	prompt, err := buildWorkbenchProviderTaskPlanPrompt(options, basePlan, knowledge)
	if err != nil {
		return nil, WorkbenchProviderView{}, "", err
	}

	modelOutput, providerView, err := RunWorkbenchProviderPrompt(providerConfig, defaultWorkbenchPlanSystemPrompt, prompt)
	if err != nil {
		return nil, providerView, modelOutput, err
	}

	flowYAML := strings.TrimSpace(ExtractWorkbenchFlowYAML(modelOutput))
	if flowYAML == "" {
		return nil, providerView, modelOutput, fmt.Errorf("provider did not return flow yaml")
	}

	normalizedFlowYAML, normalizationWarnings := normalizeWorkbenchProviderFlowYAML(flowYAML, basePlan, knowledge.Site)
	if strings.TrimSpace(normalizedFlowYAML) != "" {
		flowYAML = normalizedFlowYAML
	}

	flow, err := ParseFlow([]byte(flowYAML), "yaml")
	if err != nil {
		return nil, providerView, modelOutput, err
	}
	applyWorkbenchProviderFlowDefaults(flow, basePlan, knowledge.Site, options.Intent)
	if err := ValidateFlow(flow); err != nil {
		return nil, providerView, modelOutput, err
	}

	encodedYAML, err := encodeWorkbenchFlowYAML(flow)
	if err != nil {
		return nil, providerView, modelOutput, err
	}

	plan := cloneWorkbenchTaskPlan(basePlan)
	plan.SiteID = normalizeWorkbenchSiteID(options.SiteID)
	plan.Intent = strings.TrimSpace(options.Intent)
	plan.Flow = flow
	plan.FlowName = flow.Name
	plan.FlowYAML = encodedYAML
	plan.GenerationMode = "provider"
	plan.Provider = &providerView
	plan.ModelOutput = modelOutput
	plan.ValidationError = ""
	plan.Strategy = firstNonEmpty(strings.TrimSpace(plan.Strategy), deriveWorkbenchProviderPlanStrategy(plan, knowledge))
	plan.Reason = buildWorkbenchProviderPlanReason(plan.Reason, providerView)
	for _, warning := range normalizationWarnings {
		addWorkbenchPlanWarning(plan, warning)
	}
	return plan, providerView, modelOutput, nil
}

func loadWorkbenchTaskPlanningKnowledge(siteID string, artifactRoot string) (*workbenchTaskPlanningKnowledge, error) {
	site, err := LoadWorkbenchSiteConfig(siteID, artifactRoot)
	if err != nil {
		return nil, err
	}

	pages, err := ListWorkbenchPageCards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	apis, err := ListWorkbenchAPICards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	entities, err := ListWorkbenchEntityCards(siteID, artifactRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &workbenchTaskPlanningKnowledge{
		Site:     site,
		Pages:    pages,
		APIs:     apis,
		Entities: entities,
	}, nil
}

func buildWorkbenchProviderTaskPlanPrompt(
	options WorkbenchTaskPlanOptions,
	basePlan *WorkbenchTaskPlan,
	knowledge *workbenchTaskPlanningKnowledge,
) (string, error) {
	if knowledge == nil || knowledge.Site == nil {
		return "", fmt.Errorf("site context is required")
	}
	realtimeContext := compactWorkbenchRealtimeContextForAI(options.realtimeContext)
	if len(knowledge.Pages) == 0 && len(knowledge.APIs) == 0 && strings.TrimSpace(flowYAMLFromPlan(basePlan)) == "" && realtimeContext == nil {
		return "", fmt.Errorf("no explored site context is available for provider planning")
	}

	contextPayload := map[string]any{
		"intent": strings.TrimSpace(options.Intent),
		"site": map[string]any{
			"site_id":         knowledge.Site.SiteID,
			"name":            knowledge.Site.Name,
			"start_url":       knowledge.Site.StartURL,
			"allowed_domains": knowledge.Site.AllowedDomains,
			"session_name":    knowledge.Site.SessionName,
		},
		"generation_rules":  flowSchemaGenerationRules(),
		"selector_strategy": flowSelectorStrategy(),
		"authoring_checklist": []string{
			"Return exactly one TSPlay Flow YAML document.",
			"Set schema_version to \"1\".",
			"Use only supported TSPlay actions and fields from the manifest below.",
			"Prefer stable selectors from selector_candidates, data-testid, data-cy, id, aria-label, name, placeholder, or clear visible text.",
			"Reuse the saved browser session when one is available.",
			"If browser.use_session is present, it must be a saved session name string. Never emit true/false for browser.use_session.",
			"If required user input is unknown, define it in flow.vars as TODO and reference it with {{var_name}}.",
			"Avoid lua unless no structured action can express the task.",
		},
		"supported_actions":    compactWorkbenchFlowActionManifest(),
		"matched_pages":        buildWorkbenchPageContextsForAI(basePlan, knowledge.Pages),
		"matched_apis":         buildWorkbenchAPIContextsForAI(basePlan, knowledge.APIs),
		"known_entities":       buildWorkbenchEntityContextsForAI(knowledge.Entities),
		"recommended_examples": buildWorkbenchRecommendedExamplesForAI(strings.TrimSpace(options.Intent)),
	}
	if baselineYAML := flowYAMLFromPlan(basePlan); baselineYAML != "" {
		contextPayload["baseline_plan"] = map[string]any{
			"strategy":        basePlan.Strategy,
			"reason":          basePlan.Reason,
			"generation_mode": basePlan.GenerationMode,
			"flow_yaml":       baselineYAML,
		}
	}
	if realtimeContext != nil {
		contextPayload["realtime_context"] = realtimeContext
	}

	encoded, err := json.MarshalIndent(contextPayload, "", "  ")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.Join([]string{
		"Generate one TSPlay Flow YAML document for the user intent using the explored website context below.",
		"Return only YAML with no markdown fences and no explanation.",
		"",
		"Context JSON:",
		string(encoded),
	}, "\n")), nil
}

func compactWorkbenchFlowActionManifest() []map[string]any {
	manifest := buildFlowActionManifest()
	compacted := make([]map[string]any, 0, len(manifest))
	for _, item := range manifest {
		compactItem := map[string]any{
			"name":        item["name"],
			"description": item["description"],
		}
		if aliases, ok := item["common_aliases"]; ok {
			compactItem["common_aliases"] = aliases
		}
		if args, ok := item["args"].([]map[string]any); ok && len(args) > 0 {
			compactItem["args"] = args
		}
		compacted = append(compacted, compactItem)
	}
	return compacted
}

func buildWorkbenchPageContextsForAI(basePlan *WorkbenchTaskPlan, cards []WorkbenchPageCard) []map[string]any {
	contexts := []map[string]any{}
	seen := map[string]bool{}
	appendCard := func(card *WorkbenchPageCard, score int) {
		if card == nil || seen[card.ID] {
			return
		}
		seen[card.ID] = true
		context := compactWorkbenchPageCardForAI(*card)
		if score > 0 {
			context["score"] = score
		}
		if observation := loadWorkbenchObservationForAI(card.ObservationPath); observation != nil {
			context["observation"] = compactWorkbenchObservationForAI(observation)
		}
		contexts = append(contexts, context)
	}

	for _, candidate := range planMatchedPages(basePlan) {
		appendCard(findWorkbenchPageByID(cards, candidate.ID), candidate.Score)
		if len(contexts) >= defaultWorkbenchPlanPageContextLimit {
			return contexts
		}
	}
	for i := range cards {
		appendCard(&cards[i], 0)
		if len(contexts) >= defaultWorkbenchPlanPageContextLimit {
			break
		}
	}
	return contexts
}

func buildWorkbenchAPIContextsForAI(basePlan *WorkbenchTaskPlan, cards []WorkbenchAPICard) []map[string]any {
	contexts := []map[string]any{}
	seen := map[string]bool{}
	appendCard := func(card *WorkbenchAPICard, score int) {
		if card == nil || seen[card.ID] {
			return
		}
		seen[card.ID] = true
		context := compactWorkbenchAPICardForAI(*card)
		if score > 0 {
			context["score"] = score
		}
		contexts = append(contexts, context)
	}

	for _, candidate := range planMatchedAPIs(basePlan) {
		appendCard(findWorkbenchAPIByID(cards, candidate.ID), candidate.Score)
		if len(contexts) >= defaultWorkbenchPlanAPIContextLimit {
			return contexts
		}
	}
	for i := range cards {
		appendCard(&cards[i], 0)
		if len(contexts) >= defaultWorkbenchPlanAPIContextLimit {
			break
		}
	}
	return contexts
}

func buildWorkbenchEntityContextsForAI(cards []WorkbenchEntityCard) []map[string]any {
	contexts := make([]map[string]any, 0, minWorkbenchInt(len(cards), defaultWorkbenchPlanEntityContextLimit))
	for i := 0; i < len(cards) && i < defaultWorkbenchPlanEntityContextLimit; i++ {
		card := cards[i]
		fields := make([]map[string]any, 0, minWorkbenchInt(len(card.Fields), defaultWorkbenchPlanTextListLimit))
		for j := 0; j < len(card.Fields) && j < defaultWorkbenchPlanTextListLimit; j++ {
			fields = append(fields, map[string]any{
				"name":  card.Fields[j].Name,
				"label": card.Fields[j].Label,
				"type":  card.Fields[j].Type,
			})
		}
		contexts = append(contexts, map[string]any{
			"id":     card.ID,
			"name":   card.Name,
			"label":  card.Label,
			"fields": fields,
		})
	}
	return contexts
}

func compactWorkbenchPageCardForAI(card WorkbenchPageCard) map[string]any {
	value := map[string]any{
		"id":               card.ID,
		"title":            firstNonEmpty(card.Title, card.NormalizedRoute, card.URL),
		"url":              card.URL,
		"normalized_route": card.NormalizedRoute,
		"summary":          card.Summary,
		"menu_path":        trimWorkbenchStringList(card.MenuPath, defaultWorkbenchPlanTextListLimit),
		"breadcrumbs":      trimWorkbenchStringList(card.Breadcrumbs, defaultWorkbenchPlanTextListLimit),
		"text_snippets":    trimWorkbenchStringList(card.TextSnippets, defaultWorkbenchPlanTextListLimit),
		"key_elements":     trimWorkbenchStringList(card.KeyElements, defaultWorkbenchPlanTextListLimit),
		"flow_hints":       card.FlowHints,
	}
	if len(card.InputFields) > 0 {
		fields := make([]map[string]any, 0, minWorkbenchInt(len(card.InputFields), defaultWorkbenchPlanTextListLimit))
		for i := 0; i < len(card.InputFields) && i < defaultWorkbenchPlanTextListLimit; i++ {
			fields = append(fields, map[string]any{
				"name":     card.InputFields[i].Name,
				"label":    card.InputFields[i].Label,
				"selector": card.InputFields[i].Selector,
			})
		}
		value["input_fields"] = fields
	}
	if len(card.Actions) > 0 {
		actions := make([]map[string]any, 0, minWorkbenchInt(len(card.Actions), defaultWorkbenchPlanTextListLimit))
		for i := 0; i < len(card.Actions) && i < defaultWorkbenchPlanTextListLimit; i++ {
			actions = append(actions, map[string]any{
				"label":    card.Actions[i].Label,
				"kind":     card.Actions[i].Kind,
				"selector": card.Actions[i].Selector,
				"risk":     card.Actions[i].Risk,
			})
		}
		value["actions"] = actions
	}
	if len(card.Tables) > 0 {
		tables := make([]map[string]any, 0, minWorkbenchInt(len(card.Tables), 2))
		for i := 0; i < len(card.Tables) && i < 2; i++ {
			tables = append(tables, map[string]any{
				"name":     card.Tables[i].Name,
				"selector": card.Tables[i].Selector,
				"columns":  trimWorkbenchStringList(card.Tables[i].Columns, defaultWorkbenchPlanTextListLimit),
			})
		}
		value["tables"] = tables
	}
	return value
}

func compactWorkbenchRealtimeContextForAI(realtimeContext *workbenchRealtimeContext) map[string]any {
	if realtimeContext == nil {
		return nil
	}
	value := map[string]any{}
	if urlValue := strings.TrimSpace(realtimeContext.URL); urlValue != "" {
		value["url"] = urlValue
	}
	if title := strings.TrimSpace(realtimeContext.Title); title != "" {
		value["title"] = title
	}
	if html := strings.TrimSpace(realtimeContext.HTML); html != "" {
		value["html_excerpt"] = workbenchCleanText(html, defaultWorkbenchPlanHTMLExcerptLimit)
	}
	if realtimeContext.Observation != nil {
		value["observation"] = compactWorkbenchObservationForAI(realtimeContext.Observation)
	}
	if len(value) == 0 {
		return nil
	}
	return value
}

func normalizeWorkbenchProviderFlowYAML(flowYAML string, basePlan *WorkbenchTaskPlan, site *WorkbenchSiteConfig) (string, []string) {
	flowYAML = strings.TrimSpace(flowYAML)
	if flowYAML == "" {
		return "", nil
	}

	var doc map[string]any
	if err := yaml.Unmarshal([]byte(flowYAML), &doc); err != nil {
		return flowYAML, nil
	}

	browser, ok := doc["browser"].(map[string]any)
	if !ok || len(browser) == 0 {
		return flowYAML, nil
	}

	useSessionValue, exists := browser["use_session"]
	if !exists {
		return flowYAML, nil
	}

	preferredSession := resolveWorkbenchProviderUseSessionFallback(basePlan, site)
	normalizedValue, changed, warning := normalizeWorkbenchProviderUseSessionValue(useSessionValue, preferredSession)
	if !changed {
		return flowYAML, nil
	}

	if normalizedValue == nil {
		delete(browser, "use_session")
	} else {
		browser["use_session"] = normalizedValue
	}
	if len(browser) == 0 {
		delete(doc, "browser")
	}

	encoded, err := yaml.Marshal(doc)
	if err != nil {
		return flowYAML, nil
	}
	if strings.TrimSpace(warning) == "" {
		return string(encoded), nil
	}
	return string(encoded), []string{warning}
}

func normalizeWorkbenchProviderUseSessionValue(value any, preferredSession string) (any, bool, string) {
	switch typed := value.(type) {
	case bool:
		if typed {
			if preferredSession != "" {
				return preferredSession, true, fmt.Sprintf("AI 生成的 browser.use_session=true 已自动改写为已保存 Session %q。", preferredSession)
			}
			return nil, true, "AI 生成的 browser.use_session=true 已自动移除，因为当前站点没有配置可复用的已保存 Session。"
		}
		return nil, true, "AI 生成的 browser.use_session=false 已自动移除。"
	case string:
		normalized := strings.ToLower(strings.TrimSpace(typed))
		switch normalized {
		case "":
			return nil, true, "AI 生成的 browser.use_session 为空，已自动移除。"
		case "true":
			if preferredSession != "" {
				return preferredSession, true, fmt.Sprintf("AI 生成的 browser.use_session=\"true\" 已自动改写为已保存 Session %q。", preferredSession)
			}
			return nil, true, "AI 生成的 browser.use_session=\"true\" 已自动移除，因为当前站点没有配置可复用的已保存 Session。"
		case "false":
			return nil, true, "AI 生成的 browser.use_session=\"false\" 已自动移除。"
		}
	}
	return nil, false, ""
}

func compactWorkbenchAPICardForAI(card WorkbenchAPICard) map[string]any {
	return map[string]any{
		"id":              card.ID,
		"method":          card.Method,
		"path_template":   card.PathTemplate,
		"url":             card.URL,
		"semantic_name":   firstNonEmpty(card.SemanticName, card.PathTemplate),
		"trigger_route":   card.TriggerRoute,
		"trigger_action":  card.TriggerAction,
		"operation_type":  card.OperationType,
		"risk":            card.Risk,
		"request_schema":  card.RequestSchema,
		"response_schema": card.ResponseSchema,
		"resource_type":   card.ResourceType,
	}
}

func compactWorkbenchObservationForAI(observation *PageObservation) map[string]any {
	if observation == nil {
		return nil
	}
	value := map[string]any{
		"url":                  observation.URL,
		"title":                observation.Title,
		"page_summary":         observation.PageSummary,
		"dom_snapshot_excerpt": observation.DOMSnapshotExcerpt,
	}
	if len(observation.ContentElements) > 0 {
		content := make([]map[string]any, 0, minWorkbenchInt(len(observation.ContentElements), defaultWorkbenchPlanObservationContentLimit))
		for i := 0; i < len(observation.ContentElements) && i < defaultWorkbenchPlanObservationContentLimit; i++ {
			item := observation.ContentElements[i]
			content = append(content, map[string]any{
				"index":    item.Index,
				"kind":     item.Kind,
				"tag":      item.Tag,
				"text":     item.Text,
				"href":     item.Href,
				"selector": item.Selector,
				"xpath":    item.XPath,
			})
		}
		value["content_elements"] = content
	}
	if len(observation.Elements) > 0 {
		elements := make([]map[string]any, 0, minWorkbenchInt(len(observation.Elements), defaultWorkbenchPlanObservationElementLimit))
		for i := 0; i < len(observation.Elements) && i < defaultWorkbenchPlanObservationElementLimit; i++ {
			item := observation.Elements[i]
			elements = append(elements, map[string]any{
				"index":               item.Index,
				"tag":                 item.Tag,
				"type":                item.Type,
				"role":                item.Role,
				"id":                  item.ID,
				"name":                item.Name,
				"text":                item.Text,
				"label":               item.Label,
				"placeholder":         item.Placeholder,
				"aria_label":          item.AriaLabel,
				"near_text":           item.NearText,
				"primary_selector":    item.PrimarySelector,
				"selector_candidates": trimWorkbenchStringList(item.SelectorCandidates, 4),
			})
		}
		value["interactive_elements"] = elements
	}
	return value
}

func buildWorkbenchRecommendedExamplesForAI(intent string) []map[string]any {
	examples := BuildRecommendedFlowExamples(intent, nil, nil, defaultWorkbenchPlanExampleLimit)
	result := make([]map[string]any, 0, len(examples))
	for _, example := range examples {
		result = append(result, map[string]any{
			"id":            example.ID,
			"title":         example.Title,
			"description":   example.Description,
			"when_to_use":   example.WhenToUse,
			"focus_actions": example.FocusActions,
			"flow_yaml":     example.FlowYAML,
		})
	}
	return result
}

func loadWorkbenchObservationForAI(path string) *PageObservation {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var observation PageObservation
	if err := json.Unmarshal(content, &observation); err != nil {
		return nil
	}
	return &observation
}

func applyWorkbenchProviderFlowDefaults(flow *Flow, basePlan *WorkbenchTaskPlan, site *WorkbenchSiteConfig, intent string) {
	if flow == nil {
		return
	}
	flow.SchemaVersion = CurrentFlowSchemaVersion
	if strings.TrimSpace(flow.Name) == "" {
		flow.Name = buildDraftFlowName("", intent, nil)
	}
	if basePlan != nil && basePlan.Flow != nil && len(basePlan.Flow.Vars) > 0 {
		if flow.Vars == nil {
			flow.Vars = map[string]any{}
		}
		for key, value := range basePlan.Flow.Vars {
			if _, exists := flow.Vars[key]; exists {
				continue
			}
			flow.Vars[key] = value
		}
	}
	useSession := resolveWorkbenchProviderUseSessionFallback(basePlan, site)
	if useSession != "" {
		if flow.Browser == nil {
			flow.Browser = &FlowBrowserConfig{}
		}
		if flow.Browser.UseSession == "" {
			flow.Browser.UseSession = useSession
		}
	}
	if flow.Description == "" {
		flow.Description = fmt.Sprintf("AI-generated from explored site context. Intent: %s", strings.TrimSpace(intent))
	}
}

func resolveWorkbenchProviderUseSessionFallback(basePlan *WorkbenchTaskPlan, site *WorkbenchSiteConfig) string {
	useSession := ""
	if basePlan != nil && basePlan.Flow != nil && basePlan.Flow.Browser != nil && strings.TrimSpace(basePlan.Flow.Browser.UseSession) != "" {
		useSession = strings.TrimSpace(basePlan.Flow.Browser.UseSession)
	}
	if useSession == "" && site != nil && strings.TrimSpace(site.SessionName) != "" {
		useSession = strings.TrimSpace(site.SessionName)
	}
	return useSession
}

func cloneWorkbenchTaskPlan(plan *WorkbenchTaskPlan) *WorkbenchTaskPlan {
	if plan == nil {
		return &WorkbenchTaskPlan{}
	}
	cloned := *plan
	cloned.MatchedPages = append([]WorkbenchTaskCandidate(nil), plan.MatchedPages...)
	cloned.MatchedAPIs = append([]WorkbenchTaskCandidate(nil), plan.MatchedAPIs...)
	cloned.Warnings = append([]string(nil), plan.Warnings...)
	if plan.Provider != nil {
		provider := *plan.Provider
		cloned.Provider = &provider
	}
	return &cloned
}

func deriveWorkbenchProviderPlanStrategy(plan *WorkbenchTaskPlan, knowledge *workbenchTaskPlanningKnowledge) string {
	if plan != nil && strings.TrimSpace(plan.Strategy) != "" {
		return strings.TrimSpace(plan.Strategy)
	}
	if plan != nil && len(plan.MatchedPages) > 0 {
		return "ui_first"
	}
	if plan != nil && len(plan.MatchedAPIs) > 0 {
		return "api_first"
	}
	if knowledge != nil && len(knowledge.Pages) > 0 {
		return "ui_first"
	}
	if knowledge != nil && len(knowledge.APIs) > 0 {
		return "api_first"
	}
	return "context_ai"
}

func buildWorkbenchProviderPlanReason(baseReason string, providerView WorkbenchProviderView) string {
	providerLabel := firstNonEmpty(strings.TrimSpace(providerView.Name), strings.TrimSpace(providerView.ProviderID), "AI provider")
	modelLabel := firstNonEmpty(strings.TrimSpace(providerView.ResolvedModel), strings.TrimSpace(providerView.Model))
	if modelLabel != "" {
		providerLabel = fmt.Sprintf("%s (%s)", providerLabel, modelLabel)
	}
	reason := fmt.Sprintf("Used %s with explored site context to generate the Flow.", providerLabel)
	if strings.TrimSpace(baseReason) != "" {
		reason += " 本地匹配结果也已作为 baseline 提供给模型。"
	}
	return reason
}

func planMatchedPages(plan *WorkbenchTaskPlan) []WorkbenchTaskCandidate {
	if plan == nil {
		return nil
	}
	return plan.MatchedPages
}

func planMatchedAPIs(plan *WorkbenchTaskPlan) []WorkbenchTaskCandidate {
	if plan == nil {
		return nil
	}
	return plan.MatchedAPIs
}

func flowYAMLFromPlan(plan *WorkbenchTaskPlan) string {
	if plan == nil {
		return ""
	}
	return strings.TrimSpace(plan.FlowYAML)
}

func trimWorkbenchStringList(items []string, limit int) []string {
	if limit <= 0 || len(items) <= limit {
		return append([]string(nil), items...)
	}
	return append([]string(nil), items[:limit]...)
}
