package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type FlowDraftOptions struct {
	Intent       string
	URL          string
	FlowName     string
	ArtifactRoot string
	Observation  *PageObservation
	Security     *FlowSecurityPolicy
}

type FlowDraft struct {
	Intent            string                    `json:"intent"`
	FlowName          string                    `json:"flow_name"`
	Flow              *Flow                     `json:"flow,omitempty"`
	FlowYAML          string                    `json:"flow_yaml,omitempty"`
	InitialValidation *FlowDraftValidation      `json:"initial_validation,omitempty"`
	Validation        *FlowDraftValidation      `json:"validation,omitempty"`
	AutoRepaired      bool                      `json:"auto_repaired,omitempty"`
	SelectorRepairs   []FlowDraftSelectorRepair `json:"selector_repairs,omitempty"`
	RepairHints       []FlowRepairHint          `json:"repair_hints,omitempty"`
	PlannedActions    []string                  `json:"planned_actions,omitempty"`
	SuggestedVars     map[string]any            `json:"suggested_vars,omitempty"`
	MatchedElements   []FlowDraftMatch          `json:"matched_elements,omitempty"`
	Assumptions       []string                  `json:"assumptions,omitempty"`
	Warnings          []string                  `json:"warnings,omitempty"`
	Unresolved        []string                  `json:"unresolved,omitempty"`
	NextSteps         []string                  `json:"next_steps,omitempty"`
}

type FlowDraftMatch struct {
	Purpose  string `json:"purpose"`
	Source   string `json:"source"`
	Index    int    `json:"index,omitempty"`
	Selector string `json:"selector,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Text     string `json:"text,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type FlowDraftValidation struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name,omitempty"`
	Steps int    `json:"steps,omitempty"`
	Error string `json:"error,omitempty"`
}

type FlowDraftSelectorRepair struct {
	StepPath      string `json:"step_path"`
	Action        string `json:"action,omitempty"`
	Name          string `json:"name,omitempty"`
	Purpose       string `json:"purpose,omitempty"`
	ObservedIndex int    `json:"observed_index,omitempty"`
	FromSelector  string `json:"from_selector,omitempty"`
	ToSelector    string `json:"to_selector,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

type draftIntentAction struct {
	Kind     string
	Position int
}

type draftObservedNode struct {
	Tag      string              `json:"tag"`
	XPath    string              `json:"xpath"`
	Text     string              `json:"text"`
	Children []draftObservedNode `json:"children"`
}

type flattenedDraftNode struct {
	Tag   string
	XPath string
	Text  string
}

type flowDraftBuilder struct {
	intent            string
	intentLower       string
	url               string
	observation       *PageObservation
	artifactRoot      string
	flow              *Flow
	draft             *FlowDraft
	domNodes          []flattenedDraftNode
	quotedLiterals    []string
	searchLiteralUsed bool
}

var draftIDXPathPattern = regexp.MustCompile(`^\s*//\*\[@id="([^"]+)"\]\s*$`)
var draftQuotedLiteralPattern = regexp.MustCompile(`["'“”‘’]([^"'“”‘’]{1,160})["'“”‘’]`)
var draftFilePathPattern = regexp.MustCompile(`([~/A-Za-z]:)?[/\\][^"'“”‘’\s]+\.[A-Za-z0-9]{1,8}`)
var draftValidationStepPattern = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*)`)
var draftValidationAllowOptionPattern = regexp.MustCompile(`set ([a-z_]+)=true`)
var draftValidationUnknownVariablePattern = regexp.MustCompile(`unknown variable "([^"]+)"`)
var draftValidationMissingFieldPattern = regexp.MustCompile(`requires "?([a-z_]+)"?`)

func BuildDraftFlow(options FlowDraftOptions) (*FlowDraft, error) {
	if strings.TrimSpace(options.Intent) == "" {
		return nil, fmt.Errorf("intent is required")
	}
	if options.Observation == nil {
		return nil, fmt.Errorf("observation is required")
	}

	artifactRoot := strings.TrimSpace(options.ArtifactRoot)
	if artifactRoot == "" {
		artifactRoot = firstNonEmpty(options.Observation.ArtifactRoot, DefaultFlowArtifactRoot)
	}

	builder := &flowDraftBuilder{
		intent:       strings.TrimSpace(options.Intent),
		intentLower:  strings.ToLower(strings.TrimSpace(options.Intent)),
		url:          strings.TrimSpace(options.URL),
		observation:  options.Observation,
		artifactRoot: artifactRoot,
		draft: &FlowDraft{
			Intent:          strings.TrimSpace(options.Intent),
			FlowName:        buildDraftFlowName(options.FlowName, options.Intent, options.Observation),
			SuggestedVars:   map[string]any{},
			Assumptions:     []string{},
			Warnings:        []string{},
			Unresolved:      []string{},
			NextSteps:       []string{"Review the drafted selectors and TODO variables.", "Review the auto validation result and any selector repairs.", "Run tsplay.run_flow and iterate with tsplay.repair_flow_context if needed."},
			MatchedElements: []FlowDraftMatch{},
		},
	}
	builder.flow = &Flow{
		SchemaVersion: CurrentFlowSchemaVersion,
		Name:          builder.draft.FlowName,
		Description:   fmt.Sprintf("Auto-drafted from user intent and page observation. Intent: %s", builder.intent),
		Vars:          map[string]any{},
		Steps:         []FlowStep{},
	}
	builder.quotedLiterals = extractDraftQuotedLiterals(builder.intent)
	builder.loadDOMNodes()
	builder.bootstrapNavigateStep()
	builder.buildIntentSteps()

	if len(builder.flow.Steps) == 0 {
		return nil, fmt.Errorf("could not draft any steps from the provided observation")
	}
	if len(builder.draft.PlannedActions) == 0 {
		builder.draft.PlannedActions = []string{"navigate"}
	}
	if len(builder.draft.SuggestedVars) == 0 {
		builder.draft.SuggestedVars = nil
	}
	if len(builder.draft.MatchedElements) == 0 {
		builder.draft.MatchedElements = nil
	}
	if len(builder.draft.Assumptions) == 0 {
		builder.draft.Assumptions = nil
	}
	if len(builder.draft.Warnings) == 0 {
		builder.draft.Warnings = nil
	}
	if len(builder.draft.Unresolved) == 0 {
		builder.draft.Unresolved = nil
	}
	builder.draft.InitialValidation = validateDraftFlow(builder.flow, options.Security)
	repairs := repairDraftFlowSelectors(builder.flow, builder.observation, builder.draft.MatchedElements)
	if len(repairs) > 0 {
		builder.draft.AutoRepaired = true
		builder.draft.SelectorRepairs = repairs
		applyDraftSelectorRepairsToMatches(builder.draft.MatchedElements, repairs)
	}
	builder.draft.Validation = validateDraftFlow(builder.flow, options.Security)
	if len(builder.draft.SelectorRepairs) == 0 {
		builder.draft.SelectorRepairs = nil
	}
	builder.draft.RepairHints = buildDraftFlowRepairHints(builder.flow, builder.draft)
	if len(builder.draft.RepairHints) == 0 {
		builder.draft.RepairHints = nil
	}

	encoded, err := yaml.Marshal(builder.flow)
	if err != nil {
		return nil, fmt.Errorf("marshal drafted flow: %w", err)
	}
	builder.draft.Flow = builder.flow
	builder.draft.FlowYAML = string(encoded)
	return builder.draft, nil
}

func ParseObservationForDraft(observationText string) (*PageObservation, error) {
	observationText = strings.TrimSpace(observationText)
	if observationText == "" {
		return nil, nil
	}

	var observation PageObservation
	if err := json.Unmarshal([]byte(observationText), &observation); err == nil && observation.URL != "" {
		return &observation, nil
	}

	var wrapper struct {
		Observation *PageObservation `json:"observation"`
	}
	if err := json.Unmarshal([]byte(observationText), &wrapper); err == nil && wrapper.Observation != nil && wrapper.Observation.URL != "" {
		return wrapper.Observation, nil
	}
	return nil, fmt.Errorf("observation must be a JSON PageObservation or a wrapper with an observation field")
}

func (b *flowDraftBuilder) bootstrapNavigateStep() {
	url := firstNonEmpty(b.observation.URL, b.flowStepURLFallback())
	if strings.TrimSpace(url) == "" {
		b.draft.Assumptions = append(b.draft.Assumptions, "The observation does not include a page URL, so the draft starts after navigation.")
		return
	}
	b.flow.Vars["target_url"] = url
	b.draft.SuggestedVars["target_url"] = url
	b.flow.Steps = append(b.flow.Steps, FlowStep{
		Name:   "open target page",
		Action: "navigate",
		URL:    "{{target_url}}",
	})
	b.draft.PlannedActions = append(b.draft.PlannedActions, "navigate")
}

func (b *flowDraftBuilder) flowStepURLFallback() string {
	return b.url
}

func (b *flowDraftBuilder) buildIntentSteps() {
	actions := detectDraftIntentActions(b.intentLower)
	if len(actions) == 0 {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent did not map cleanly to a known action pattern, so only a minimal navigation step was drafted.")
		b.draft.Unresolved = append(b.draft.Unresolved, "intent_to_actions")
		return
	}

	for _, action := range actions {
		if action.Kind == "submit" && draftHasAnyPlannedAction(b.draft.PlannedActions, "upload", "login", "search", "export") {
			continue
		}
		switch action.Kind {
		case "login":
			if !b.addLoginSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "search":
			if !b.addSearchSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "select":
			if !b.addSelectSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "upload":
			if !b.addUploadSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "extract_title":
			if !b.addTitleExtractionStep() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "extract_count":
			if !b.addCountExtractionStep() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "capture_table":
			if !b.addCaptureTableStep() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "export":
			if !b.addExportSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		case "submit":
			if !b.addSubmitSteps() {
				b.draft.Unresolved = append(b.draft.Unresolved, action.Kind)
			}
		}
	}
}

func detectDraftIntentActions(intentLower string) []draftIntentAction {
	type actionKeywords struct {
		kind     string
		keywords []string
	}
	keywordSets := []actionKeywords{
		{kind: "login", keywords: []string{"login", "log in", "sign in", "登录"}},
		{kind: "search", keywords: []string{"search", "query", "find", "lookup", "filter", "搜索", "查询", "检索", "查找", "筛选"}},
		{kind: "select", keywords: []string{"select", "choose", "pick", "选择", "选中"}},
		{kind: "upload", keywords: []string{"upload", "上传"}},
		{kind: "extract_title", keywords: []string{"title", "heading", "标题"}},
		{kind: "extract_count", keywords: []string{"count", "total", "summary", "数量", "总数", "统计"}},
		{kind: "capture_table", keywords: []string{"table", "grid", "表格", "列表", "清单"}},
		{kind: "export", keywords: []string{"export", "download", "导出", "下载"}},
		{kind: "submit", keywords: []string{"submit", "confirm", "save", "提交", "确认", "保存", "继续"}},
	}

	actions := make([]draftIntentAction, 0, len(keywordSets))
	seen := map[string]bool{}
	for _, item := range keywordSets {
		position := findFirstDraftKeywordPosition(intentLower, item.keywords)
		if position < 0 || seen[item.kind] {
			continue
		}
		actions = append(actions, draftIntentAction{Kind: item.kind, Position: position})
		seen[item.kind] = true
	}
	sort.SliceStable(actions, func(i, j int) bool {
		if actions[i].Position == actions[j].Position {
			return actions[i].Kind < actions[j].Kind
		}
		return actions[i].Position < actions[j].Position
	})
	return actions
}

func findFirstDraftKeywordPosition(intentLower string, keywords []string) int {
	best := -1
	for _, keyword := range keywords {
		index := strings.Index(intentLower, strings.ToLower(keyword))
		if index < 0 {
			continue
		}
		if best < 0 || index < best {
			best = index
		}
	}
	return best
}

func (b *flowDraftBuilder) addSearchSteps() bool {
	input, reason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftTextInput(element) {
			return -1, ""
		}
		metadata := observedElementMetadata(element)
		score := 10
		score += scoreDraftKeywords(metadata, []string{"search", "query", "keyword", "filter", "搜索", "查询", "关键", "筛选"})
		score += scoreDraftKeywords(metadata, []string{"order", "订单"})
		return score, "matched a text input whose label, placeholder, or metadata looks like a search field"
	})
	if input == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions search/query behavior, but no reliable text input was found.")
		return false
	}

	selector := bestObservedSelector(*input)
	if selector == "" {
		b.draft.Assumptions = append(b.draft.Assumptions, "A likely search input was found, but it has no selector candidates.")
		return false
	}

	queryVar := "query_text"
	if strings.Contains(observedElementMetadata(*input), "order") || strings.Contains(observedElementMetadata(*input), "订单") {
		queryVar = "order_query"
	}
	queryValue := b.consumeQuotedLiteral("TODO")
	b.flow.Vars[queryVar] = queryValue
	b.draft.SuggestedVars[queryVar] = queryValue
	if queryValue == "TODO" {
		b.draft.Assumptions = append(b.draft.Assumptions, fmt.Sprintf("Fill %q with the actual search content before execution.", queryVar))
	}

	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "wait for search input", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
		FlowStep{Name: "fill search query", Action: "type_text", Selector: selector, Text: fmt.Sprintf("{{%s}}", queryVar)},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "search")
	b.noteElementMatch("search_input", *input, selector, reason)

	button, buttonReason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftButtonLike(element) {
			return -1, ""
		}
		metadata := observedElementMetadata(element)
		score := 1 + scoreDraftKeywords(metadata, []string{"search", "query", "find", "filter", "搜索", "查询", "查找", "筛选"})
		return score, "matched a visible button or link whose text looks like a search trigger"
	})
	if button == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "No explicit search button was found, so the draft stops after filling the query.")
		return true
	}

	buttonSelector := bestObservedSelector(*button)
	if buttonSelector == "" {
		return true
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "submit search", Action: "click", Selector: buttonSelector},
		FlowStep{Name: "wait for search result update", Action: "wait_for_network_idle"},
	)
	b.noteElementMatch("search_submit", *button, buttonSelector, buttonReason)
	return true
}

func (b *flowDraftBuilder) addExportSteps() bool {
	button, reason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftButtonLike(element) {
			return -1, ""
		}
		metadata := observedElementMetadata(element)
		score := scoreDraftKeywords(metadata, []string{"export", "download", "导出", "下载"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a visible control whose text looks like export/download"
	})
	if button == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions export/download, but no matching control was found.")
		return false
	}

	selector := bestObservedSelector(*button)
	if selector == "" {
		return false
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "wait for export control", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
		FlowStep{Name: "click export", Action: "click", Selector: selector},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "export")
	b.noteElementMatch("export", *button, selector, reason)
	return true
}

func (b *flowDraftBuilder) addUploadSteps() bool {
	fileInput, reason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if strings.EqualFold(element.Type, "file") {
			return 100, "matched a file input"
		}
		return -1, ""
	})
	if fileInput == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions upload, but no file input was found.")
		return false
	}
	selector := bestObservedSelector(*fileInput)
	if selector == "" {
		return false
	}

	filePath := extractDraftFilePath(b.intent)
	if filePath == "" {
		filePath = "TODO"
		b.draft.Assumptions = append(b.draft.Assumptions, `Fill "upload_file_path" with a real local file path before execution.`)
	}
	b.flow.Vars["upload_file_path"] = filePath
	b.draft.SuggestedVars["upload_file_path"] = filePath
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "wait for upload input", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
		FlowStep{Name: "choose upload file", Action: "upload_file", Selector: selector, FilePath: "{{upload_file_path}}"},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "upload")
	b.noteElementMatch("upload_input", *fileInput, selector, reason)

	button, buttonReason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftButtonLike(element) {
			return -1, ""
		}
		metadata := observedElementMetadata(element)
		score := scoreDraftKeywords(metadata, []string{"upload", "submit", "发送", "上传", "提交"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a visible upload/submit button"
	})
	if button == nil {
		return true
	}
	buttonSelector := bestObservedSelector(*button)
	if buttonSelector == "" {
		return true
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "submit upload", Action: "click", Selector: buttonSelector},
		FlowStep{Name: "wait for upload request", Action: "wait_for_network_idle"},
	)
	b.noteElementMatch("upload_submit", *button, buttonSelector, buttonReason)
	return true
}

func (b *flowDraftBuilder) addLoginSteps() bool {
	usernameInput, usernameReason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftTextInput(element) {
			return -1, ""
		}
		metadata := observedElementMetadata(element)
		score := scoreDraftKeywords(metadata, []string{"user", "username", "account", "email", "phone", "用户", "用户名", "账号", "邮箱", "手机号"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a username/account input"
	})
	passwordInput, passwordReason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		metadata := observedElementMetadata(element)
		if strings.EqualFold(element.Type, "password") {
			return 100, "matched a password input"
		}
		score := scoreDraftKeywords(metadata, []string{"password", "pwd", "密码"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a password input"
	})
	if usernameInput == nil && passwordInput == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions login, but neither username nor password inputs were found.")
		return false
	}

	literals := extractDraftQuotedLiterals(b.intent)
	usernameValue := "TODO"
	passwordValue := "TODO"
	if len(literals) >= 1 {
		usernameValue = literals[0]
	}
	if len(literals) >= 2 {
		passwordValue = literals[1]
	}
	b.flow.Vars["username"] = usernameValue
	b.flow.Vars["password"] = passwordValue
	b.draft.SuggestedVars["username"] = usernameValue
	b.draft.SuggestedVars["password"] = passwordValue
	if usernameValue == "TODO" || passwordValue == "TODO" {
		b.draft.Assumptions = append(b.draft.Assumptions, `Fill "username" and "password" before execution.`)
	}

	if usernameInput != nil {
		selector := bestObservedSelector(*usernameInput)
		if selector != "" {
			b.flow.Steps = append(b.flow.Steps,
				FlowStep{Name: "wait for username input", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
				FlowStep{Name: "fill username", Action: "type_text", Selector: selector, Text: "{{username}}"},
			)
			b.noteElementMatch("login_username", *usernameInput, selector, usernameReason)
		}
	}
	if passwordInput != nil {
		selector := bestObservedSelector(*passwordInput)
		if selector != "" {
			b.flow.Steps = append(b.flow.Steps,
				FlowStep{Name: "fill password", Action: "type_text", Selector: selector, Text: "{{password}}"},
			)
			b.noteElementMatch("login_password", *passwordInput, selector, passwordReason)
		}
	}

	button, buttonReason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftButtonLike(element) {
			return -1, ""
		}
		score := scoreDraftKeywords(observedElementMetadata(element), []string{"login", "sign in", "submit", "登录", "进入", "提交"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a login button"
	})
	if button != nil {
		selector := bestObservedSelector(*button)
		if selector != "" {
			b.flow.Steps = append(b.flow.Steps,
				FlowStep{Name: "submit login", Action: "click", Selector: selector},
				FlowStep{Name: "wait for login result", Action: "wait_for_network_idle"},
			)
			b.noteElementMatch("login_submit", *button, selector, buttonReason)
		}
	}
	b.draft.PlannedActions = append(b.draft.PlannedActions, "login")
	return true
}

func (b *flowDraftBuilder) addSelectSteps() bool {
	selectElement, reason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !(strings.EqualFold(element.Tag, "select") || strings.EqualFold(element.Role, "combobox")) {
			return -1, ""
		}
		return 100 + scoreDraftKeywords(observedElementMetadata(element), []string{"select", "choose", "option", "选择", "选项"}), "matched a dropdown/select control"
	})
	if selectElement == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions selecting an option, but no dropdown-like control was found.")
		return false
	}
	selector := bestObservedSelector(*selectElement)
	if selector == "" {
		return false
	}

	value := b.consumeQuotedLiteral("TODO")
	b.flow.Vars["selected_value"] = value
	b.draft.SuggestedVars["selected_value"] = value
	if value == "TODO" {
		b.draft.Assumptions = append(b.draft.Assumptions, `Fill "selected_value" with the target option before execution.`)
	}

	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "wait for select control", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
		FlowStep{Name: "select option", Action: "select_option", Selector: selector, Value: "{{selected_value}}"},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "select")
	b.noteElementMatch("select_control", *selectElement, selector, reason)
	return true
}

func (b *flowDraftBuilder) addSubmitSteps() bool {
	button, reason := findBestObservedElement(b.observation.Elements, func(element PageObservationElement) (int, string) {
		if !isDraftButtonLike(element) {
			return -1, ""
		}
		score := scoreDraftKeywords(observedElementMetadata(element), []string{"submit", "confirm", "save", "continue", "提交", "确认", "保存", "继续"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a generic submit/confirm control"
	})
	if button == nil {
		return false
	}
	selector := bestObservedSelector(*button)
	if selector == "" {
		return false
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "click submit control", Action: "click", Selector: selector},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "submit")
	b.noteElementMatch("submit", *button, selector, reason)
	return true
}

func (b *flowDraftBuilder) addTitleExtractionStep() bool {
	node, reason := findBestDraftNode(b.domNodes, func(node flattenedDraftNode) (int, string) {
		tag := strings.ToUpper(strings.TrimSpace(node.Tag))
		if tag != "H1" && tag != "H2" && tag != "H3" {
			return -1, ""
		}
		score := 100 + len(strings.TrimSpace(node.Text))
		return score, "matched a heading node from the DOM snapshot"
	})
	if node == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions a title/heading, but no heading node was found in the DOM snapshot.")
		return false
	}
	selector := selectorFromDraftXPath(node.XPath)
	if selector == "" {
		return false
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "extract page title", Action: "extract_text", Selector: selector, SaveAs: "page_title"},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "extract_title")
	b.draft.MatchedElements = append(b.draft.MatchedElements, FlowDraftMatch{
		Purpose:  "extract_title",
		Source:   "dom_snapshot",
		Selector: selector,
		Tag:      node.Tag,
		Text:     node.Text,
		Reason:   reason,
	})
	return true
}

func (b *flowDraftBuilder) addCountExtractionStep() bool {
	node, reason := findBestDraftNode(b.domNodes, func(node flattenedDraftNode) (int, string) {
		textLower := strings.ToLower(node.Text)
		score := 0
		if draftContainsDigits(node.Text) {
			score += 50
		}
		score += scoreDraftKeywords(textLower, []string{"count", "total", "summary", "数量", "总数", "统计"})
		if score <= 0 {
			return -1, ""
		}
		return score, "matched a DOM snapshot node whose text looks like a count or summary"
	})
	if node == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions a count/summary, but no matching text node was found in the DOM snapshot.")
		return false
	}
	selector := selectorFromDraftXPath(node.XPath)
	if selector == "" {
		return false
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "extract summary count", Action: "extract_text", Selector: selector, Pattern: "([0-9]+)", SaveAs: "summary_count"},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "extract_count")
	b.draft.MatchedElements = append(b.draft.MatchedElements, FlowDraftMatch{
		Purpose:  "extract_count",
		Source:   "dom_snapshot",
		Selector: selector,
		Tag:      node.Tag,
		Text:     node.Text,
		Reason:   reason,
	})
	return true
}

func (b *flowDraftBuilder) addCaptureTableStep() bool {
	node, reason := findBestDraftNode(b.domNodes, func(node flattenedDraftNode) (int, string) {
		if strings.ToUpper(strings.TrimSpace(node.Tag)) != "TABLE" {
			return -1, ""
		}
		return 100, "matched a table node from the DOM snapshot"
	})
	if node == nil {
		b.draft.Assumptions = append(b.draft.Assumptions, "The intent mentions a table/list, but no table node was found in the DOM snapshot.")
		return false
	}
	selector := selectorFromDraftXPath(node.XPath)
	if selector == "" {
		return false
	}
	b.flow.Steps = append(b.flow.Steps,
		FlowStep{Name: "wait for table", Action: "wait_for_selector", Selector: selector, Timeout: 10000},
		FlowStep{Name: "capture table", Action: "capture_table", Selector: selector, SaveAs: "table_rows"},
	)
	b.draft.PlannedActions = append(b.draft.PlannedActions, "capture_table")
	b.draft.MatchedElements = append(b.draft.MatchedElements, FlowDraftMatch{
		Purpose:  "capture_table",
		Source:   "dom_snapshot",
		Selector: selector,
		Tag:      node.Tag,
		Text:     node.Text,
		Reason:   reason,
	})
	return true
}

func (b *flowDraftBuilder) loadDOMNodes() {
	if b.observation == nil || strings.TrimSpace(b.observation.DOMSnapshotPath) == "" {
		return
	}
	nodes, err := loadDraftDOMNodes(b.observation.DOMSnapshotPath, b.artifactRoot)
	if err != nil {
		b.draft.Warnings = append(b.draft.Warnings, err.Error())
		return
	}
	b.domNodes = nodes
}

func loadDraftDOMNodes(path string, artifactRoot string) ([]flattenedDraftNode, error) {
	content, err := readDraftArtifactFile(path, artifactRoot)
	if err != nil {
		return nil, err
	}
	var root draftObservedNode
	if err := json.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("parse dom snapshot %q: %w", path, err)
	}
	nodes := []flattenedDraftNode{}
	flattenDraftDOMNodes(root, &nodes)
	return nodes, nil
}

func flattenDraftDOMNodes(node draftObservedNode, out *[]flattenedDraftNode) {
	if strings.TrimSpace(node.Tag) != "" && strings.TrimSpace(node.XPath) != "" {
		*out = append(*out, flattenedDraftNode{
			Tag:   strings.TrimSpace(node.Tag),
			XPath: strings.TrimSpace(node.XPath),
			Text:  strings.TrimSpace(node.Text),
		})
	}
	for _, child := range node.Children {
		flattenDraftDOMNodes(child, out)
	}
}

func readDraftArtifactFile(path string, artifactRoot string) ([]byte, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("artifact path is empty")
	}
	root := strings.TrimSpace(artifactRoot)
	if root == "" {
		root = DefaultFlowArtifactRoot
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve artifact root %q: %w", root, err)
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return nil, fmt.Errorf("artifact root %q is not accessible: %w", rootAbs, err)
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootReal, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return nil, fmt.Errorf("resolve artifact path %q: %w", path, err)
	}
	candidateReal, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return nil, fmt.Errorf("artifact path %q is not accessible: %w", path, err)
	}
	if err := ensurePathInsideRoot(candidateReal, rootReal); err != nil {
		return nil, fmt.Errorf("artifact path %q is outside allowed artifact root %q", path, rootReal)
	}
	content, err := os.ReadFile(candidateReal)
	if err != nil {
		return nil, fmt.Errorf("read artifact path %q: %w", candidateReal, err)
	}
	return content, nil
}

func findBestObservedElement(elements []PageObservationElement, score func(PageObservationElement) (int, string)) (*PageObservationElement, string) {
	bestIndex := -1
	bestScore := -1
	bestReason := ""
	for index := range elements {
		currentScore, reason := score(elements[index])
		if currentScore > bestScore {
			bestIndex = index
			bestScore = currentScore
			bestReason = reason
		}
	}
	if bestIndex < 0 || bestScore < 0 {
		return nil, ""
	}
	return &elements[bestIndex], bestReason
}

func findBestDraftNode(nodes []flattenedDraftNode, score func(flattenedDraftNode) (int, string)) (*flattenedDraftNode, string) {
	bestIndex := -1
	bestScore := -1
	bestReason := ""
	for index := range nodes {
		currentScore, reason := score(nodes[index])
		if currentScore > bestScore {
			bestIndex = index
			bestScore = currentScore
			bestReason = reason
		}
	}
	if bestIndex < 0 || bestScore < 0 {
		return nil, ""
	}
	return &nodes[bestIndex], bestReason
}

func isDraftTextInput(element PageObservationElement) bool {
	if !element.Visible || !element.Enabled {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(element.Type)) {
	case "text", "search", "email", "number", "tel", "url", "textarea":
		return true
	}
	switch strings.ToLower(strings.TrimSpace(element.Tag)) {
	case "input", "textarea":
		return strings.ToLower(strings.TrimSpace(element.Type)) != "file" && strings.ToLower(strings.TrimSpace(element.Type)) != "password"
	}
	return strings.EqualFold(element.Role, "textbox")
}

func isDraftButtonLike(element PageObservationElement) bool {
	if !element.Visible || !element.Enabled {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(element.Tag)) {
	case "button", "a":
		return true
	}
	switch strings.ToLower(strings.TrimSpace(element.Role)) {
	case "button", "link":
		return true
	}
	return false
}

func observedElementMetadata(element PageObservationElement) string {
	parts := []string{
		element.Tag,
		element.Type,
		element.Role,
		element.ID,
		element.Name,
		element.Text,
		element.Label,
		element.Placeholder,
		element.AriaLabel,
		element.Href,
		element.NearText,
	}
	keys := make([]string, 0, len(element.Attributes))
	for key := range element.Attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, key, element.Attributes[key])
	}
	return strings.ToLower(strings.Join(parts, " "))
}

func scoreDraftKeywords(haystack string, keywords []string) int {
	score := 0
	for _, keyword := range keywords {
		if strings.Contains(haystack, strings.ToLower(keyword)) {
			score += len([]rune(keyword)) + 10
		}
	}
	return score
}

func bestObservedSelector(element PageObservationElement) string {
	if selector := preferredObservedSelector(element); strings.TrimSpace(selector) != "" {
		return selector
	}
	for _, selector := range element.SelectorCandidates {
		if strings.TrimSpace(selector) != "" {
			return selector
		}
	}
	return ""
}

func selectorFromDraftXPath(xpath string) string {
	xpath = strings.TrimSpace(xpath)
	if xpath == "" {
		return ""
	}
	if matches := draftIDXPathPattern.FindStringSubmatch(xpath); len(matches) == 2 {
		return "#" + matches[1]
	}
	return "xpath=" + xpath
}

func (b *flowDraftBuilder) noteElementMatch(purpose string, element PageObservationElement, selector string, reason string) {
	b.draft.MatchedElements = append(b.draft.MatchedElements, FlowDraftMatch{
		Purpose:  purpose,
		Source:   "observed_element",
		Index:    element.Index,
		Selector: selector,
		Tag:      element.Tag,
		Text:     firstNonEmpty(element.Text, element.Label, element.Placeholder),
		Reason:   reason,
	})
}

func extractDraftQuotedLiterals(intent string) []string {
	matches := draftQuotedLiteralPattern.FindAllStringSubmatch(intent, -1)
	values := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		value := strings.TrimSpace(match[1])
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

func (b *flowDraftBuilder) consumeQuotedLiteral(fallback string) string {
	if !b.searchLiteralUsed && len(b.quotedLiterals) > 0 {
		b.searchLiteralUsed = true
		return b.quotedLiterals[0]
	}
	return fallback
}

func extractDraftFilePath(intent string) string {
	match := draftFilePathPattern.FindString(intent)
	return strings.TrimSpace(match)
}

func draftContainsDigits(value string) bool {
	for _, char := range value {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

func buildDraftFlowName(explicit string, intent string, observation *PageObservation) string {
	if candidate := sanitizeDraftIdentifier(explicit); candidate != "" {
		return candidate
	}
	actions := detectDraftIntentActions(strings.ToLower(strings.TrimSpace(intent)))
	parts := []string{"draft"}
	for _, action := range actions {
		parts = append(parts, action.Kind)
	}
	if observation != nil {
		if titlePart := sanitizeDraftIdentifier(observation.Title); titlePart != "" {
			parts = append(parts, titlePart)
		}
	}
	name := sanitizeDraftIdentifier(strings.Join(parts, "_"))
	if name == "" {
		return "draft_observed_page"
	}
	return name
}

func sanitizeDraftIdentifier(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, " ", "_")
	builder := strings.Builder{}
	lastUnderscore := false
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastUnderscore = false
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastUnderscore = false
		default:
			if !lastUnderscore {
				builder.WriteRune('_')
				lastUnderscore = true
			}
		}
	}
	result := strings.Trim(builder.String(), "_")
	if result == "" {
		return ""
	}
	if result[0] >= '0' && result[0] <= '9' {
		result = "draft_" + result
	}
	return result
}

func draftHasAnyPlannedAction(actions []string, wants ...string) bool {
	for _, action := range actions {
		for _, want := range wants {
			if action == want {
				return true
			}
		}
	}
	return false
}

func validateDraftFlow(flow *Flow, security *FlowSecurityPolicy) *FlowDraftValidation {
	validation := &FlowDraftValidation{
		Valid: false,
	}
	if flow == nil {
		validation.Error = "flow is nil"
		return validation
	}
	validation.Name = flow.Name
	validation.Steps = len(flow.Steps)
	if err := ValidateFlow(flow); err != nil {
		validation.Error = err.Error()
		return validation
	}
	if security != nil {
		if err := ValidateFlowSecurity(flow, *security); err != nil {
			validation.Error = err.Error()
			return validation
		}
	}
	validation.Valid = true
	return validation
}

type draftSelectorRepairTarget struct {
	Preferred     string
	ObservedIndex int
	Purpose       string
}

func repairDraftFlowSelectors(flow *Flow, observation *PageObservation, matches []FlowDraftMatch) []FlowDraftSelectorRepair {
	if flow == nil || observation == nil || len(observation.Elements) == 0 {
		return nil
	}
	targets := buildDraftSelectorRepairTargets(observation.Elements, matches)
	if len(targets) == 0 {
		return nil
	}
	repairs := []FlowDraftSelectorRepair{}
	repairDraftSelectorStepSequence(flow.Steps, "", targets, &repairs)
	return repairs
}

func buildDraftSelectorRepairTargets(elements []PageObservationElement, matches []FlowDraftMatch) map[string]draftSelectorRepairTarget {
	purposeByIndex := map[int]string{}
	for _, match := range matches {
		if match.Source != "observed_element" || match.Index == 0 || strings.TrimSpace(match.Purpose) == "" {
			continue
		}
		if _, ok := purposeByIndex[match.Index]; !ok {
			purposeByIndex[match.Index] = match.Purpose
		}
	}

	targets := map[string]draftSelectorRepairTarget{}
	for _, element := range elements {
		preferred := preferredObservedSelector(element)
		if preferred == "" {
			continue
		}
		target := draftSelectorRepairTarget{
			Preferred:     preferred,
			ObservedIndex: element.Index,
			Purpose:       purposeByIndex[element.Index],
		}
		for _, selector := range element.SelectorCandidates {
			selector = strings.TrimSpace(selector)
			if selector == "" {
				continue
			}
			targets[selector] = target
		}
	}
	return targets
}

func repairDraftSelectorStepSequence(steps []FlowStep, parentPath string, targets map[string]draftSelectorRepairTarget, repairs *[]FlowDraftSelectorRepair) {
	for index := range steps {
		stepPath := flowStepPath(parentPath, index+1)
		repairDraftSelectorStep(&steps[index], stepPath, targets, repairs)
	}
}

func repairDraftSelectorStep(step *FlowStep, stepPath string, targets map[string]draftSelectorRepairTarget, repairs *[]FlowDraftSelectorRepair) {
	if step == nil {
		return
	}
	current := strings.TrimSpace(step.Selector)
	if current != "" && len(flowReferences(current)) == 0 {
		if target, ok := targets[current]; ok && target.Preferred != "" && target.Preferred != current {
			*repairs = append(*repairs, FlowDraftSelectorRepair{
				StepPath:      stepPath,
				Action:        step.Action,
				Name:          step.Name,
				Purpose:       target.Purpose,
				ObservedIndex: target.ObservedIndex,
				FromSelector:  current,
				ToSelector:    target.Preferred,
				Reason:        "Replaced a weaker selector candidate with the preferred selector from the page observation.",
			})
			step.Selector = target.Preferred
		}
	}
	if step.Condition != nil {
		repairDraftSelectorStep(step.Condition, stepPath+".condition", targets, repairs)
	}
	if len(step.Steps) > 0 {
		nestedPath := stepPath
		if step.Action == "on_error" {
			nestedPath = stepPath + ".try"
		}
		repairDraftSelectorStepSequence(step.Steps, nestedPath, targets, repairs)
	}
	if len(step.Then) > 0 {
		repairDraftSelectorStepSequence(step.Then, stepPath+".then", targets, repairs)
	}
	if len(step.Else) > 0 {
		repairDraftSelectorStepSequence(step.Else, stepPath+".else", targets, repairs)
	}
	if len(step.OnError) > 0 {
		repairDraftSelectorStepSequence(step.OnError, stepPath+".on_error", targets, repairs)
	}
}

func applyDraftSelectorRepairsToMatches(matches []FlowDraftMatch, repairs []FlowDraftSelectorRepair) {
	if len(matches) == 0 || len(repairs) == 0 {
		return
	}
	selectorByIndex := map[int]string{}
	for _, repair := range repairs {
		if repair.ObservedIndex == 0 || strings.TrimSpace(repair.ToSelector) == "" {
			continue
		}
		selectorByIndex[repair.ObservedIndex] = repair.ToSelector
	}
	if len(selectorByIndex) == 0 {
		return
	}
	for index := range matches {
		if matches[index].Source != "observed_element" || matches[index].Index == 0 {
			continue
		}
		if selector, ok := selectorByIndex[matches[index].Index]; ok {
			matches[index].Selector = selector
		}
	}
}

func buildDraftFlowRepairHints(flow *Flow, draft *FlowDraft) []FlowRepairHint {
	if flow == nil || draft == nil || draft.Validation == nil || draft.Validation.Valid {
		return nil
	}

	hints := buildDraftValidationRepairHints(flow, draft.Validation.Error)
	if len(hints) == 0 {
		hints = append(hints, FlowRepairHint{
			Priority:        1,
			Source:          "draft_validation",
			Reason:          "The drafted Flow still does not pass validation after the selector repair pass.",
			Suggestion:      "Start from the validation error, compare the affected step with tsplay.flow_schema and the current page observation, then apply the smallest possible fix.",
			Error:           draft.Validation.Error,
			FailureCategory: "validation",
		})
	}

	if len(draft.SelectorRepairs) > 0 {
		repairedSteps := map[string]bool{}
		for _, repair := range draft.SelectorRepairs {
			if strings.TrimSpace(repair.StepPath) != "" {
				repairedSteps[repair.StepPath] = true
			}
		}
		for index := range hints {
			if repairedSteps[hints[index].StepPath] {
				hints[index].Reason += " This step was already selector-repaired once, so inspect the action fields, variables, or policy flags next."
			}
		}
	}

	sort.SliceStable(hints, func(i, j int) bool {
		if hints[i].Priority != hints[j].Priority {
			return hints[i].Priority < hints[j].Priority
		}
		if hints[i].StepPath == hints[j].StepPath {
			return hints[i].Reason < hints[j].Reason
		}
		return hints[i].StepPath < hints[j].StepPath
	})
	return dedupeFlowRepairHints(hints)
}

func buildDraftValidationRepairHints(flow *Flow, validationError string) []FlowRepairHint {
	validationError = strings.TrimSpace(validationError)
	if validationError == "" {
		return nil
	}

	stepPath := parseDraftValidationStepPath(validationError)
	step, found := findFlowStepByPath(flow, stepPath)
	hint := FlowRepairHint{
		Priority:        1,
		Source:          "draft_validation",
		StepPath:        stepPath,
		Error:           validationError,
		FailureCategory: "validation",
	}
	if found {
		hint.Action = step.Action
		hint.Name = step.Name
		hint.Selector = step.Selector
	}

	switch {
	case strings.Contains(validationError, "disabled by security policy"):
		option := parseDraftValidationAllowOption(validationError)
		hint.Targets = []string{"security_policy", "action"}
		hint.Reason = "The drafted step is valid structurally, but it is blocked by the current safety flags."
		if option != "" {
			hint.Suggestion = fmt.Sprintf("Inspect step %s first. If this is a trusted automation, rerun draft_flow or validate_flow with %s=true; otherwise replace the step with a lower-risk action.", flowRepairHintStepLabel(stepPath), option)
		} else {
			hint.Suggestion = "Inspect the blocked action first and decide whether to enable the matching allow_* flag or replace it with a lower-risk action."
		}
	case strings.Contains(validationError, "references unknown variable"):
		unknownVar := parseDraftValidationUnknownVariable(validationError)
		hint.Targets = []string{"variables", "save_as"}
		hint.Reason = "This step depends on a variable that is not defined yet."
		if unknownVar != "" {
			hint.Suggestion = fmt.Sprintf("Inspect step %s and the steps right before it. Make sure %q exists in flow.vars or is produced earlier with save_as/set_var.", flowRepairHintStepLabel(stepPath), unknownVar)
		} else {
			hint.Suggestion = "Inspect the affected step and the previous producing steps. Make sure every {{var}} is defined in vars or produced earlier with save_as/set_var."
		}
	case strings.Contains(validationError, "requires"):
		field := parseDraftValidationMissingField(validationError)
		hint.Targets = []string{"required_fields"}
		if field != "" {
			hint.Targets = append(hint.Targets, field)
		}
		hint.Reason = "The drafted step is missing a required field for its action."
		if field != "" {
			hint.Suggestion = fmt.Sprintf("Inspect step %s first and fill %q using the page observation, extracted variables, or a safer default.", flowRepairHintStepLabel(stepPath), field)
		} else {
			hint.Suggestion = "Inspect the affected step first and fill the missing required field using the action schema and page observation."
		}
	case strings.Contains(validationError, "uses unsupported action"):
		hint.Targets = []string{"action"}
		hint.Reason = "The action name is not part of the current Flow DSL."
		hint.Suggestion = "Inspect the action name first. Replace it with a supported action from tsplay.list_actions, or move the special logic into a lua escape hatch only if necessary."
	case strings.Contains(validationError, "is not a valid variable name"):
		hint.Targets = []string{"variables", "save_as"}
		hint.Reason = "A variable or save_as name does not match the Flow identifier rules."
		hint.Suggestion = "Rename the offending variable to letters, numbers, and underscores only, and keep the first character non-numeric."
	case strings.Contains(validationError, "schema_version"):
		hint.Priority = 0
		hint.Targets = []string{"schema_version"}
		hint.Reason = "The Flow schema version is missing or unsupported."
		hint.Suggestion = fmt.Sprintf("Set schema_version to %q before validating again.", CurrentFlowSchemaVersion)
	case strings.Contains(validationError, "flow must contain at least one step"):
		hint.Priority = 0
		hint.Targets = []string{"steps"}
		hint.Reason = "The draft did not produce any executable steps."
		hint.Suggestion = "Start from the user intent and page observation again, then add at least one concrete Flow step."
	default:
		hint.Targets = []string{"validation"}
		hint.Reason = "The draft failed validation and needs a small targeted repair."
		if stepPath != "" {
			hint.Suggestion = fmt.Sprintf("Inspect step %s first and compare its fields against tsplay.flow_schema plus the observed selector candidates.", stepPath)
		} else {
			hint.Suggestion = "Inspect the validation error and compare the draft against tsplay.flow_schema before retrying."
		}
	}

	return []FlowRepairHint{hint}
}

func parseDraftValidationStepPath(validationError string) string {
	matches := draftValidationStepPattern.FindStringSubmatch(validationError)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func parseDraftValidationAllowOption(validationError string) string {
	matches := draftValidationAllowOptionPattern.FindStringSubmatch(validationError)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func parseDraftValidationUnknownVariable(validationError string) string {
	matches := draftValidationUnknownVariablePattern.FindStringSubmatch(validationError)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func parseDraftValidationMissingField(validationError string) string {
	matches := draftValidationMissingFieldPattern.FindStringSubmatch(validationError)
	if len(matches) == 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func findFlowStepByPath(flow *Flow, stepPath string) (FlowStep, bool) {
	if flow == nil || strings.TrimSpace(stepPath) == "" {
		return FlowStep{}, false
	}
	return findFlowStepInSequence(flow.Steps, strings.Split(stepPath, "."))
}

func findFlowStepInSequence(steps []FlowStep, tokens []string) (FlowStep, bool) {
	if len(tokens) == 0 {
		return FlowStep{}, false
	}
	index, err := strconv.Atoi(tokens[0])
	if err != nil || index < 1 || index > len(steps) {
		return FlowStep{}, false
	}
	step := steps[index-1]
	if len(tokens) == 1 {
		return step, true
	}
	if _, err := strconv.Atoi(tokens[1]); err == nil {
		return findFlowStepInSequence(step.Steps, tokens[1:])
	}
	if len(tokens) < 3 {
		return FlowStep{}, false
	}
	branch := tokens[1]
	rest := tokens[2:]
	switch branch {
	case "condition":
		if step.Condition == nil {
			return FlowStep{}, false
		}
		return findFlowStepInSequence([]FlowStep{*step.Condition}, rest)
	case "then":
		return findFlowStepInSequence(step.Then, rest)
	case "else":
		return findFlowStepInSequence(step.Else, rest)
	case "on_error":
		return findFlowStepInSequence(step.OnError, rest)
	case "try":
		return findFlowStepInSequence(step.Steps, rest)
	default:
		return findFlowStepInSequence(step.Steps, tokens[1:])
	}
}
