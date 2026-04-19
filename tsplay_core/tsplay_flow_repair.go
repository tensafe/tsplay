package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const defaultFlowRepairArtifactExcerpt = 4000

type FlowRepairContextOptions struct {
	Flow               *Flow
	Result             *FlowResult
	Trace              []FlowStepTrace
	Error              string
	ArtifactRoot       string
	MaxArtifactExcerpt int
}

type FlowRepairContext struct {
	FlowName            string                  `json:"flow_name,omitempty"`
	FlowDescription     string                  `json:"flow_description,omitempty"`
	Error               string                  `json:"error,omitempty"`
	FailureCategory     string                  `json:"failure_category,omitempty"`
	FailureReason       string                  `json:"failure_reason,omitempty"`
	FailedStepPath      string                  `json:"failed_step_path,omitempty"`
	RepairHints         []FlowRepairHint        `json:"repair_hints,omitempty"`
	FailedStep          *FlowRepairStepContext  `json:"failed_step,omitempty"`
	NearbySteps         []FlowRepairStepContext `json:"nearby_steps,omitempty"`
	TraceSummary        []FlowRepairTraceItem   `json:"trace_summary,omitempty"`
	Artifacts           *FlowRepairArtifacts    `json:"artifacts,omitempty"`
	Variables           map[string]any          `json:"variables,omitempty"`
	FocusedVariables    map[string]any          `json:"focused_variables,omitempty"`
	RepairTargets       []string                `json:"repair_targets,omitempty"`
	AllowedActions      []string                `json:"allowed_actions,omitempty"`
	GenerationRules     []string                `json:"generation_rules,omitempty"`
	SelectorStrategy    []string                `json:"selector_strategy,omitempty"`
	RepairInstructions  []string                `json:"repair_instructions,omitempty"`
	ValidationChecklist []string                `json:"validation_checklist,omitempty"`
	Prompt              string                  `json:"prompt,omitempty"`
}

type FlowRepairStepContext struct {
	Index    int                  `json:"index"`
	Path     string               `json:"path,omitempty"`
	Relation string               `json:"relation,omitempty"`
	Step     FlowStep             `json:"step"`
	Trace    *FlowRepairTraceItem `json:"trace,omitempty"`
}

type FlowRepairTraceItem struct {
	Index         int                      `json:"index"`
	Path          string                   `json:"path,omitempty"`
	Label         string                   `json:"label,omitempty"`
	Attempt       int                      `json:"attempt,omitempty"`
	Iteration     int                      `json:"iteration,omitempty"`
	Branch        string                   `json:"branch,omitempty"`
	Name          string                   `json:"name,omitempty"`
	Action        string                   `json:"action"`
	Status        string                   `json:"status"`
	SaveAs        string                   `json:"save_as,omitempty"`
	Selector      string                   `json:"selector,omitempty"`
	Text          string                   `json:"text,omitempty"`
	URL           string                   `json:"url,omitempty"`
	ArgsSummary   string                   `json:"args_summary,omitempty"`
	OutputSummary string                   `json:"output_summary,omitempty"`
	Error         string                   `json:"error,omitempty"`
	ErrorStack    string                   `json:"error_stack,omitempty"`
	PageURL       string                   `json:"page_url,omitempty"`
	DurationMS    int64                    `json:"duration_ms,omitempty"`
	Artifacts     *FlowRepairArtifactPaths `json:"artifacts,omitempty"`
	Condition     *FlowRepairTraceItem     `json:"condition,omitempty"`
	Children      []FlowRepairTraceItem    `json:"children,omitempty"`
	Attempts      []FlowRepairTraceItem    `json:"attempts,omitempty"`
}

type FlowRepairArtifacts struct {
	Paths                FlowRepairArtifactPaths `json:"paths"`
	ArtifactSummary      []string                `json:"artifact_summary,omitempty"`
	DOMSnapshotExcerpt   string                  `json:"dom_snapshot_excerpt,omitempty"`
	DOMSnapshotTruncated bool                    `json:"dom_snapshot_truncated,omitempty"`
	RelevantDOM          []string                `json:"relevant_dom,omitempty"`
	RelevantSelectors    []string                `json:"relevant_selectors,omitempty"`
	ReadErrors           []string                `json:"read_errors,omitempty"`
}

type FlowRepairArtifactPaths struct {
	Directory       string `json:"directory,omitempty"`
	ScreenshotPath  string `json:"screenshot_path,omitempty"`
	HTMLPath        string `json:"html_path,omitempty"`
	DOMSnapshotPath string `json:"dom_snapshot_path,omitempty"`
	CaptureError    string `json:"capture_error,omitempty"`
}

func BuildFlowRepairContext(options FlowRepairContextOptions) (*FlowRepairContext, error) {
	if options.Flow == nil {
		return nil, fmt.Errorf("flow is required")
	}
	trace := options.Trace
	if options.Result != nil && len(options.Result.Trace) > 0 {
		trace = options.Result.Trace
	}
	if len(trace) == 0 {
		return nil, fmt.Errorf("failed trace or run_result is required")
	}

	artifactRoot := strings.TrimSpace(options.ArtifactRoot)
	if artifactRoot == "" && options.Result != nil {
		artifactRoot = options.Result.ArtifactRoot
	}
	if artifactRoot == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	excerptLimit := options.MaxArtifactExcerpt
	if excerptLimit <= 0 {
		excerptLimit = defaultFlowRepairArtifactExcerpt
	}

	traceSummary := buildFlowRepairTraceSummary(trace, excerptLimit)
	failedTrace, failedTraceIndex := findFailedFlowStepTrace(trace)
	failedStepNumber := failedTrace.Index
	if failedStepNumber <= 0 {
		failedStepNumber = failedTraceIndex + 1
	}
	var failedFlowStep *FlowStep
	if failedStepNumber >= 1 && failedStepNumber <= len(options.Flow.Steps) {
		failedFlowStep = &options.Flow.Steps[failedStepNumber-1]
	}

	traceByIndex := map[int]FlowRepairTraceItem{}
	traceByPath := map[string]FlowRepairTraceItem{}
	indexFlowRepairTraceItems(traceSummary, traceByIndex, traceByPath)
	failedStepPath := strings.TrimSpace(failedTrace.Path)
	if failedStepPath == "" && failedStepNumber > 0 {
		failedStepPath = fmt.Sprint(failedStepNumber)
	}
	failedLocation := resolveFlowRepairStepLocation(options.Flow, failedStepPath)
	if failedLocation != nil && failedLocation.Step != nil {
		failedFlowStep = failedLocation.Step
		if failedLocation.Position > 0 {
			failedStepNumber = failedLocation.Position
		}
		failedStepPath = failedLocation.Path
	}
	failureCategory, failureReason := classifyFlowRepairFailure(failedTrace, failedFlowStep)

	context := &FlowRepairContext{
		FlowName:            options.Flow.Name,
		FlowDescription:     options.Flow.Description,
		Error:               firstNonEmpty(options.Error, failedTrace.Error),
		FailureCategory:     failureCategory,
		FailureReason:       failureReason,
		FailedStepPath:      failedStepPath,
		TraceSummary:        traceSummary,
		NearbySteps:         buildFlowRepairNearbySteps(options.Flow, failedStepPath, failedLocation, traceByIndex, traceByPath),
		RepairTargets:       buildFlowRepairTargets(failedFlowStep, failedTrace),
		AllowedActions:      FlowActionNames(),
		GenerationRules:     flowRepairGenerationRules(),
		SelectorStrategy:    flowSelectorStrategy(),
		RepairInstructions:  flowRepairInstructions(),
		ValidationChecklist: flowRepairValidationChecklist(),
	}
	if options.Result != nil {
		if vars, ok := compactTraceValue(options.Result.Vars, 0).(map[string]any); ok && len(vars) > 0 {
			context.Variables = vars
		}
		if focused := buildFlowRepairFocusedVariables(options.Flow, failedLocation, failedStepPath, options.Result.Vars); len(focused) > 0 {
			context.FocusedVariables = focused
		}
	}
	if failedLocation != nil && failedLocation.Step != nil {
		failed := FlowRepairStepContext{
			Index:    failedLocation.Position,
			Path:     failedLocation.Path,
			Relation: "failed",
			Step:     *failedLocation.Step,
		}
		if item, ok := traceByPath[failedLocation.Path]; ok {
			failed.Trace = &item
		} else if item, ok := traceByIndex[failedStepNumber]; ok {
			failed.Trace = &item
		}
		context.FailedStep = &failed
	}
	if failedTrace.Artifacts != nil {
		context.Artifacts = buildFlowRepairArtifacts(*failedTrace.Artifacts, artifactRoot, excerptLimit, failedFlowStep, failedTrace)
	}
	context.RepairHints = buildRuntimeFlowRepairHints(context, failedTrace, failedFlowStep)
	context.Prompt = buildFlowRepairPrompt(context)
	return context, nil
}

func ParseFlowRunResultForRepair(runResultText string, traceText string) (*FlowResult, string, error) {
	runResultText = strings.TrimSpace(runResultText)
	traceText = strings.TrimSpace(traceText)
	if runResultText == "" && traceText == "" {
		return nil, "", fmt.Errorf("run_result or trace is required")
	}

	if traceText != "" {
		var trace []FlowStepTrace
		if err := json.Unmarshal([]byte(traceText), &trace); err != nil {
			return nil, "", fmt.Errorf("parse trace: %w", err)
		}
		return &FlowResult{Trace: trace}, "", nil
	}

	var result FlowResult
	if err := json.Unmarshal([]byte(runResultText), &result); err == nil && len(result.Trace) > 0 {
		return &result, "", nil
	}

	var trace []FlowStepTrace
	if err := json.Unmarshal([]byte(runResultText), &trace); err == nil && len(trace) > 0 {
		return &FlowResult{Trace: trace}, "", nil
	}

	var wrapper struct {
		Error  string          `json:"error"`
		Result json.RawMessage `json:"result"`
		Trace  json.RawMessage `json:"trace"`
	}
	if err := json.Unmarshal([]byte(runResultText), &wrapper); err != nil {
		return nil, "", fmt.Errorf("parse run_result: %w", err)
	}
	if len(wrapper.Result) > 0 && string(wrapper.Result) != "null" {
		if err := json.Unmarshal(wrapper.Result, &result); err != nil {
			return nil, wrapper.Error, fmt.Errorf("parse run_result.result: %w", err)
		}
		if len(result.Trace) > 0 {
			return &result, wrapper.Error, nil
		}
	}
	if len(wrapper.Trace) > 0 && string(wrapper.Trace) != "null" {
		if err := json.Unmarshal(wrapper.Trace, &trace); err != nil {
			return nil, wrapper.Error, fmt.Errorf("parse run_result.trace: %w", err)
		}
		if len(trace) > 0 {
			return &FlowResult{Trace: trace}, wrapper.Error, nil
		}
	}
	return nil, wrapper.Error, fmt.Errorf("run_result does not contain trace")
}

func buildFlowRepairTraceSummary(trace []FlowStepTrace, textLimit int) []FlowRepairTraceItem {
	items := make([]FlowRepairTraceItem, 0, len(trace))
	for i, step := range trace {
		index := step.Index
		if index <= 0 {
			index = i + 1
		}
		item := FlowRepairTraceItem{
			Index:         index,
			Path:          step.Path,
			Attempt:       step.Attempt,
			Iteration:     step.Iteration,
			Branch:        step.Branch,
			Name:          step.Name,
			Action:        step.Action,
			Status:        step.Status,
			SaveAs:        step.SaveAs,
			ArgsSummary:   step.ArgsSummary,
			OutputSummary: step.OutputSummary,
			Error:         step.Error,
			PageURL:       step.PageURL,
			DurationMS:    step.DurationMS,
		}
		item.Selector, item.Text, item.URL = flowRepairTraceInterestingArgs(step.Args)
		item.Label = buildFlowRepairTraceLabel(item)
		if step.ErrorStack != "" {
			item.ErrorStack, _ = truncateRepairText(step.ErrorStack, textLimit)
		}
		if step.Artifacts != nil {
			paths := flowRepairArtifactPaths(*step.Artifacts)
			item.Artifacts = &paths
		}
		if len(step.Attempts) > 0 {
			item.Attempts = buildFlowRepairTraceSummary(step.Attempts, textLimit)
		}
		if step.Condition != nil {
			condition := buildFlowRepairTraceSummary([]FlowStepTrace{*step.Condition}, textLimit)
			if len(condition) > 0 {
				item.Condition = &condition[0]
			}
		}
		if len(step.Children) > 0 {
			item.Children = buildFlowRepairTraceSummary(step.Children, textLimit)
		}
		items = append(items, item)
	}
	return items
}

func findFailedFlowStepTrace(trace []FlowStepTrace) (FlowStepTrace, int) {
	for index, step := range trace {
		if failed, ok := findFailedFlowStepTraceInItem(step); ok {
			return failed, index
		}
	}
	return trace[len(trace)-1], len(trace) - 1
}

func findFailedFlowStepTraceInItem(step FlowStepTrace) (FlowStepTrace, bool) {
	if step.Condition != nil {
		if failed, ok := findFailedFlowStepTraceInItem(*step.Condition); ok {
			return failed, true
		}
	}
	for _, attempt := range step.Attempts {
		if failed, ok := findFailedFlowStepTraceInItem(attempt); ok {
			return failed, true
		}
	}
	for _, child := range step.Children {
		if failed, ok := findFailedFlowStepTraceInItem(child); ok {
			return failed, true
		}
	}
	status := strings.ToLower(strings.TrimSpace(step.Status))
	if status == "error" || status == "failed" || status == "failure" {
		return step, true
	}
	return FlowStepTrace{}, false
}

type flowRepairStepLocation struct {
	Step         *FlowStep
	Path         string
	Sequence     []FlowStep
	SequencePath string
	Position     int
}

func indexFlowRepairTraceItems(items []FlowRepairTraceItem, traceByIndex map[int]FlowRepairTraceItem, traceByPath map[string]FlowRepairTraceItem) {
	for _, item := range items {
		if item.Index > 0 {
			if _, ok := traceByIndex[item.Index]; !ok {
				traceByIndex[item.Index] = item
			}
		}
		if path := strings.TrimSpace(item.Path); path != "" {
			traceByPath[path] = item
		}
		if item.Condition != nil {
			indexFlowRepairTraceItems([]FlowRepairTraceItem{*item.Condition}, traceByIndex, traceByPath)
		}
		if len(item.Children) > 0 {
			indexFlowRepairTraceItems(item.Children, traceByIndex, traceByPath)
		}
		if len(item.Attempts) > 0 {
			indexFlowRepairTraceItems(item.Attempts, traceByIndex, traceByPath)
		}
	}
}

func buildFlowRepairTraceLabel(item FlowRepairTraceItem) string {
	label := strings.TrimSpace(item.Path)
	if label == "" && item.Index > 0 {
		label = fmt.Sprint(item.Index)
	}
	switch {
	case item.Name != "" && label != "":
		return fmt.Sprintf("%s %s (%s)", label, item.Name, item.Action)
	case label != "":
		return fmt.Sprintf("%s %s", label, item.Action)
	case item.Name != "":
		return fmt.Sprintf("%s (%s)", item.Name, item.Action)
	default:
		return item.Action
	}
}

func flowRepairTraceInterestingArgs(args any) (selector string, text string, url string) {
	values, ok := args.(map[string]any)
	if !ok {
		return "", "", ""
	}
	selector = strings.TrimSpace(flowRepairTraceStringValue(values["selector"]))
	text = firstNonEmpty(
		strings.TrimSpace(flowRepairTraceStringValue(values["text"])),
		strings.TrimSpace(flowRepairTraceStringValue(values["value"])),
	)
	url = strings.TrimSpace(flowRepairTraceStringValue(values["url"]))
	return selector, text, url
}

func flowRepairTraceStringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return ""
	}
}

func buildFlowRepairNearbySteps(flow *Flow, failedStepPath string, failedLocation *flowRepairStepLocation, traceByIndex map[int]FlowRepairTraceItem, traceByPath map[string]FlowRepairTraceItem) []FlowRepairStepContext {
	if flow == nil || len(flow.Steps) == 0 {
		return nil
	}

	location := failedLocation
	if location == nil {
		location = resolveFlowRepairStepLocation(flow, failedStepPath)
	}
	if location == nil || location.Step == nil {
		return nil
	}
	if location.Sequence != nil && location.Position > 0 {
		return buildFlowRepairNearbyStepsFromSequence(location.Sequence, location.SequencePath, location.Position, location.Path, traceByIndex, traceByPath)
	}

	topLevelIndex := topLevelFlowRepairStepNumber(failedStepPath)
	if topLevelIndex <= 0 || topLevelIndex > len(flow.Steps) {
		return nil
	}
	steps := buildFlowRepairNearbyStepsFromSequence(flow.Steps, "", topLevelIndex, failedStepPath, traceByIndex, traceByPath)
	for i := range steps {
		if steps[i].Relation == "failed" {
			steps[i].Path = location.Path
			steps[i].Step = *location.Step
			if item, ok := traceByPath[location.Path]; ok {
				steps[i].Trace = &item
			}
			break
		}
	}
	return steps
}

func buildFlowRepairNearbyStepsFromSequence(sequence []FlowStep, sequencePath string, failedIndex int, failedPath string, traceByIndex map[int]FlowRepairTraceItem, traceByPath map[string]FlowRepairTraceItem) []FlowRepairStepContext {
	if len(sequence) == 0 || failedIndex <= 0 || failedIndex > len(sequence) {
		return nil
	}
	start := failedIndex - 2
	if start < 1 {
		start = 1
	}
	end := failedIndex + 2
	if end > len(sequence) {
		end = len(sequence)
	}

	steps := make([]FlowRepairStepContext, 0, end-start+1)
	for index := start; index <= end; index++ {
		relation := "next"
		if index < failedIndex {
			relation = "previous"
		} else if index == failedIndex {
			relation = "failed"
		}
		path := flowStepPath(sequencePath, index)
		if relation == "failed" && strings.TrimSpace(failedPath) != "" {
			path = failedPath
		}
		item := FlowRepairStepContext{
			Index:    index,
			Path:     path,
			Relation: relation,
			Step:     sequence[index-1],
		}
		if trace, ok := traceByPath[path]; ok {
			item.Trace = &trace
		} else if sequencePath == "" {
			if trace, ok := traceByIndex[index]; ok {
				item.Trace = &trace
			}
		}
		steps = append(steps, item)
	}
	return steps
}

func resolveFlowRepairStepLocation(flow *Flow, stepPath string) *flowRepairStepLocation {
	if flow == nil || len(flow.Steps) == 0 {
		return nil
	}
	tokens := strings.Split(strings.TrimSpace(stepPath), ".")
	if len(tokens) == 0 {
		return nil
	}
	index, ok := parseFlowRepairStepIndex(tokens[0])
	if !ok {
		return nil
	}
	return resolveFlowRepairStepLocationFromSequence(flow.Steps, "", index, tokens, 1)
}

func resolveFlowRepairStepLocationFromSequence(sequence []FlowStep, sequencePath string, index int, tokens []string, next int) *flowRepairStepLocation {
	if index <= 0 || index > len(sequence) {
		return nil
	}
	step := &sequence[index-1]
	path := flowStepPath(sequencePath, index)
	if next >= len(tokens) {
		return &flowRepairStepLocation{
			Step:         step,
			Path:         path,
			Sequence:     sequence,
			SequencePath: sequencePath,
			Position:     index,
		}
	}
	return resolveFlowRepairStepLocationFromStep(step, path, tokens, next)
}

func resolveFlowRepairStepLocationFromStep(step *FlowStep, path string, tokens []string, next int) *flowRepairStepLocation {
	if step == nil || next >= len(tokens) {
		return &flowRepairStepLocation{Step: step, Path: path}
	}
	field := strings.TrimSpace(tokens[next])
	switch field {
	case "condition":
		if step.Condition == nil {
			return nil
		}
		conditionPath := path + ".condition"
		if next == len(tokens)-1 {
			return &flowRepairStepLocation{Step: step.Condition, Path: conditionPath}
		}
		return resolveFlowRepairStepLocationFromStep(step.Condition, conditionPath, tokens, next+1)
	case "steps":
		if next+1 >= len(tokens) {
			return nil
		}
		index, ok := parseFlowRepairStepIndex(tokens[next+1])
		if !ok {
			return nil
		}
		return resolveFlowRepairStepLocationFromSequence(step.Steps, path+".steps", index, tokens, next+2)
	case "then":
		if next+1 >= len(tokens) {
			return nil
		}
		index, ok := parseFlowRepairStepIndex(tokens[next+1])
		if !ok {
			return nil
		}
		return resolveFlowRepairStepLocationFromSequence(step.Then, path+".then", index, tokens, next+2)
	case "else":
		if next+1 >= len(tokens) {
			return nil
		}
		index, ok := parseFlowRepairStepIndex(tokens[next+1])
		if !ok {
			return nil
		}
		return resolveFlowRepairStepLocationFromSequence(step.Else, path+".else", index, tokens, next+2)
	case "on_error":
		if next+1 >= len(tokens) {
			return nil
		}
		index, ok := parseFlowRepairStepIndex(tokens[next+1])
		if !ok {
			return nil
		}
		return resolveFlowRepairStepLocationFromSequence(step.OnError, path+".on_error", index, tokens, next+2)
	default:
		return nil
	}
}

func parseFlowRepairStepIndex(value string) (int, bool) {
	index := 0
	for _, ch := range strings.TrimSpace(value) {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		index = (index * 10) + int(ch-'0')
	}
	return index, index > 0
}

func topLevelFlowRepairStepNumber(stepPath string) int {
	first := strings.TrimSpace(strings.Split(strings.TrimSpace(stepPath), ".")[0])
	index, ok := parseFlowRepairStepIndex(first)
	if !ok {
		return 0
	}
	return index
}

type flowRepairDOMNode struct {
	Tag                string              `json:"tag,omitempty"`
	XPath              string              `json:"xpath,omitempty"`
	Text               string              `json:"text,omitempty"`
	Href               string              `json:"href,omitempty"`
	SelectorCandidates []string            `json:"selector_candidates,omitempty"`
	Children           []flowRepairDOMNode `json:"children,omitempty"`
}

var flowRepairSelectorTermPattern = regexp.MustCompile(`(?:#([A-Za-z0-9_-]+))|(?:\[(?:name|placeholder|aria-label|data-testid|data-test|data-cy|href)=["']([^"']+)["']\])|(?:text=["']([^"']+)["'])`)
var flowRepairXPathIDPattern = regexp.MustCompile(`\*\[@id="([^"]+)"\]`)

func buildFlowRepairArtifacts(artifacts FlowStepArtifacts, artifactRoot string, excerptLimit int, failedStep *FlowStep, failedTrace FlowStepTrace) *FlowRepairArtifacts {
	context := &FlowRepairArtifacts{
		Paths:           flowRepairArtifactPaths(artifacts),
		ArtifactSummary: buildFlowRepairArtifactSummary(artifacts),
	}
	if artifacts.DOMSnapshotPath == "" {
		return context
	}
	content, err := readFlowRepairArtifactContent(artifacts.DOMSnapshotPath, artifactRoot)
	if err != nil {
		context.ReadErrors = append(context.ReadErrors, err.Error())
		return context
	}
	excerpt, truncated := truncateRepairText(content, excerptLimit)
	context.DOMSnapshotExcerpt = excerpt
	context.DOMSnapshotTruncated = truncated
	context.RelevantDOM, context.RelevantSelectors = extractFlowRepairDOMEvidence(content, failedStep, failedTrace)
	return context
}

func flowRepairArtifactPaths(artifacts FlowStepArtifacts) FlowRepairArtifactPaths {
	return FlowRepairArtifactPaths{
		Directory:       artifacts.Directory,
		ScreenshotPath:  artifacts.ScreenshotPath,
		HTMLPath:        artifacts.HTMLPath,
		DOMSnapshotPath: artifacts.DOMSnapshotPath,
		CaptureError:    artifacts.CaptureError,
	}
}

func buildFlowRepairArtifactSummary(artifacts FlowStepArtifacts) []string {
	summary := []string{}
	if strings.TrimSpace(artifacts.Directory) != "" {
		summary = append(summary, fmt.Sprintf("directory: %s", filepath.Base(artifacts.Directory)))
	}
	if strings.TrimSpace(artifacts.ScreenshotPath) != "" {
		summary = append(summary, fmt.Sprintf("screenshot: %s", filepath.Base(artifacts.ScreenshotPath)))
	}
	if strings.TrimSpace(artifacts.HTMLPath) != "" {
		summary = append(summary, fmt.Sprintf("html: %s", filepath.Base(artifacts.HTMLPath)))
	}
	if strings.TrimSpace(artifacts.DOMSnapshotPath) != "" {
		summary = append(summary, fmt.Sprintf("dom_snapshot: %s", filepath.Base(artifacts.DOMSnapshotPath)))
	}
	if strings.TrimSpace(artifacts.CaptureError) != "" {
		summary = append(summary, fmt.Sprintf("capture_error: %s", artifacts.CaptureError))
	}
	return summary
}

func readFlowRepairArtifactContent(path string, artifactRoot string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}
	if strings.TrimSpace(artifactRoot) == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	rootAbs, err := filepath.Abs(artifactRoot)
	if err != nil {
		return "", fmt.Errorf("resolve artifact root %q: %w", artifactRoot, err)
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", fmt.Errorf("artifact root %q is not accessible: %w", rootAbs, err)
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootReal, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve artifact path %q: %w", path, err)
	}
	candidateReal, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", fmt.Errorf("artifact path %q is not accessible: %w", path, err)
	}
	if err := ensurePathInsideRoot(candidateReal, rootReal); err != nil {
		return "", fmt.Errorf("artifact path %q is outside allowed artifact root %q", path, rootReal)
	}

	content, err := os.ReadFile(candidateReal)
	if err != nil {
		return "", fmt.Errorf("read artifact path %q: %w", candidateReal, err)
	}
	return string(content), nil
}

func readFlowRepairArtifactExcerpt(path string, artifactRoot string, limit int) (string, bool, error) {
	content, err := readFlowRepairArtifactContent(path, artifactRoot)
	if err != nil {
		return "", false, err
	}
	excerpt, truncated := truncateRepairText(content, limit)
	return excerpt, truncated, nil
}

func extractFlowRepairDOMEvidence(content string, failedStep *FlowStep, failedTrace FlowStepTrace) ([]string, []string) {
	terms := flowRepairEvidenceTerms(failedStep, failedTrace)
	if len(terms) == 0 {
		return nil, nil
	}

	var root flowRepairDOMNode
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		return nil, nil
	}

	nodes := flattenFlowRepairDOMNodes(root)
	type evidence struct {
		summary   string
		selectors []string
		score     int
	}
	evidenceItems := []evidence{}
	for _, node := range nodes {
		score := scoreFlowRepairDOMNode(node, terms)
		if score <= 0 {
			continue
		}
		evidenceItems = append(evidenceItems, evidence{
			summary:   summarizeFlowRepairDOMNode(node),
			selectors: inferFlowRepairDOMNodeSelectors(node),
			score:     score,
		})
	}
	sort.SliceStable(evidenceItems, func(i, j int) bool {
		if evidenceItems[i].score != evidenceItems[j].score {
			return evidenceItems[i].score > evidenceItems[j].score
		}
		return evidenceItems[i].summary < evidenceItems[j].summary
	})

	highlights := []string{}
	selectorSet := map[string]bool{}
	selectors := []string{}
	for _, item := range evidenceItems {
		if strings.TrimSpace(item.summary) != "" && len(highlights) < 3 {
			highlights = append(highlights, item.summary)
		}
		for _, selector := range item.selectors {
			if strings.TrimSpace(selector) == "" || selectorSet[selector] {
				continue
			}
			selectorSet[selector] = true
			selectors = append(selectors, selector)
		}
		if len(highlights) >= 3 && len(selectors) >= 5 {
			break
		}
	}

	sort.SliceStable(selectors, func(i, j int) bool {
		leftScore := scoreObservedSelectorCandidate(selectors[i])
		rightScore := scoreObservedSelectorCandidate(selectors[j])
		if leftScore != rightScore {
			return leftScore > rightScore
		}
		return selectors[i] < selectors[j]
	})
	if len(selectors) > 5 {
		selectors = selectors[:5]
	}
	return highlights, selectors
}

func flattenFlowRepairDOMNodes(root flowRepairDOMNode) []flowRepairDOMNode {
	nodes := []flowRepairDOMNode{}
	var visit func(flowRepairDOMNode)
	visit = func(node flowRepairDOMNode) {
		if strings.TrimSpace(node.Tag) != "" || strings.TrimSpace(node.Text) != "" || strings.TrimSpace(node.XPath) != "" {
			nodes = append(nodes, node)
		}
		for _, child := range node.Children {
			visit(child)
		}
	}
	visit(root)
	return nodes
}

func flowRepairEvidenceTerms(failedStep *FlowStep, failedTrace FlowStepTrace) []string {
	terms := []string{}
	seen := map[string]bool{}
	add := func(value string) {
		value = strings.TrimSpace(strings.Join(strings.Fields(value), " "))
		if value == "" {
			return
		}
		lower := strings.ToLower(value)
		if len([]rune(lower)) < 3 && !strings.Contains(lower, "#") {
			return
		}
		if seen[lower] {
			return
		}
		seen[lower] = true
		terms = append(terms, value)
		for _, part := range strings.FieldsFunc(lower, func(r rune) bool {
			return r == ' ' || r == '_' || r == '-' || r == '/' || r == '.'
		}) {
			part = strings.TrimSpace(part)
			if len([]rune(part)) < 4 || seen[part] {
				continue
			}
			seen[part] = true
			terms = append(terms, part)
		}
	}

	if failedStep != nil {
		add(failedStep.Text)
		add(failedStep.Name)
		for _, match := range flowRepairSelectorTermPattern.FindAllStringSubmatch(failedStep.Selector, -1) {
			for _, token := range match[1:] {
				add(token)
			}
		}
		if selector := strings.TrimSpace(failedStep.Selector); selector != "" {
			add(selector)
		}
	}
	traceSelector, traceText, _ := flowRepairTraceInterestingArgs(failedTrace.Args)
	add(traceText)
	if traceSelector != "" {
		add(traceSelector)
		for _, match := range flowRepairSelectorTermPattern.FindAllStringSubmatch(traceSelector, -1) {
			for _, token := range match[1:] {
				add(token)
			}
		}
	}
	return terms
}

func scoreFlowRepairDOMNode(node flowRepairDOMNode, terms []string) int {
	haystack := strings.ToLower(strings.Join([]string{
		node.Tag,
		node.Text,
		node.Href,
		node.XPath,
		strings.Join(node.SelectorCandidates, " "),
	}, " "))
	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		if normalized == "" {
			continue
		}
		if strings.Contains(haystack, normalized) {
			score += len([]rune(normalized)) + 8
		}
	}
	return score
}

func summarizeFlowRepairDOMNode(node flowRepairDOMNode) string {
	parts := []string{}
	if strings.TrimSpace(node.Tag) != "" {
		parts = append(parts, "tag="+strings.ToLower(strings.TrimSpace(node.Tag)))
	}
	if text := strings.TrimSpace(node.Text); text != "" {
		text, _ = truncateRepairText(text, 96)
		parts = append(parts, fmt.Sprintf("text=%q", text))
	}
	if href := strings.TrimSpace(node.Href); href != "" {
		href, _ = truncateRepairText(href, 96)
		parts = append(parts, fmt.Sprintf("href=%q", href))
	}
	if xpath := strings.TrimSpace(node.XPath); xpath != "" {
		parts = append(parts, fmt.Sprintf("xpath=%q", xpath))
	}
	return strings.Join(parts, " ")
}

func inferFlowRepairDOMNodeSelectors(node flowRepairDOMNode) []string {
	selectors := []string{}
	for _, selector := range node.SelectorCandidates {
		selector = strings.TrimSpace(selector)
		if selector != "" {
			selectors = append(selectors, selector)
		}
	}
	if match := flowRepairXPathIDPattern.FindStringSubmatch(node.XPath); len(match) == 2 {
		selectors = append(selectors, "#"+match[1])
	}
	if text := strings.TrimSpace(node.Text); text != "" && len([]rune(text)) <= 80 {
		selectors = append(selectors, fmt.Sprintf("text=%q", text))
	}

	seen := map[string]bool{}
	deduped := make([]string, 0, len(selectors))
	for _, selector := range selectors {
		if seen[selector] {
			continue
		}
		seen[selector] = true
		deduped = append(deduped, selector)
	}
	return deduped
}

func truncateRepairText(value string, limit int) (string, bool) {
	if limit <= 0 {
		limit = defaultFlowRepairArtifactExcerpt
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value, false
	}
	return string(runes[:limit]) + "...(truncated)", true
}

func flowRepairGenerationRules() []string {
	return flowSchemaGenerationRules()
}

func flowRepairInstructions() []string {
	return []string{
		"Identify whether the failure is caused by selector drift, timing, navigation state, hidden elements, changed text, or missing variables.",
		"Patch only the smallest necessary part of the Flow.",
		"Add wait_for_selector or wait_for_text before fragile click/type/extract steps when the page is dynamic.",
		"Prefer extract_text + save_as and set_var over introducing lua just to move values between steps.",
		"Preserve existing save_as variables and downstream variable references unless they are the actual bug.",
		"Return a valid TSPlay Flow YAML/JSON that passes tsplay.validate_flow.",
	}
}

func buildRuntimeFlowRepairHints(context *FlowRepairContext, failedTrace FlowStepTrace, failedStep *FlowStep) []FlowRepairHint {
	if context == nil {
		return nil
	}

	stepPath := strings.TrimSpace(context.FailedStepPath)
	hints := []FlowRepairHint{}
	primary := FlowRepairHint{
		Priority:        0,
		Source:          "runtime_failure",
		StepPath:        stepPath,
		Action:          failedTrace.Action,
		Targets:         append([]string(nil), context.RepairTargets...),
		Reason:          context.FailureReason,
		Suggestion:      runtimeFlowRepairSuggestion(context, stepPath),
		Error:           firstNonEmpty(context.Error, failedTrace.Error, failedTrace.ErrorStack),
		FailureCategory: context.FailureCategory,
		PageURL:         failedTrace.PageURL,
	}
	if failedStep != nil {
		if strings.TrimSpace(failedStep.Action) != "" {
			primary.Action = failedStep.Action
		}
		primary.Name = failedStep.Name
		primary.Selector = failedStep.Selector
	}
	if primary.Name == "" {
		primary.Name = failedTrace.Name
	}
	if primary.Selector == "" && context.FailedStep != nil {
		primary.Selector = context.FailedStep.Step.Selector
	}
	if failedTrace.Artifacts != nil {
		paths := flowRepairArtifactPaths(*failedTrace.Artifacts)
		primary.Artifacts = &paths
	} else if context.Artifacts != nil {
		paths := context.Artifacts.Paths
		primary.Artifacts = &paths
	}
	hints = append(hints, primary)

	if previous := runtimeFlowRepairPreviousStepHint(context); previous != nil {
		hints = append(hints, *previous)
	}
	if artifact := runtimeFlowRepairArtifactSelectorHint(context, failedTrace, failedStep); artifact != nil {
		hints = append(hints, *artifact)
	}

	return dedupeFlowRepairHints(hints)
}

func runtimeFlowRepairSuggestion(context *FlowRepairContext, stepPath string) string {
	switch context.FailureCategory {
	case "selector_or_timing":
		return fmt.Sprintf("Inspect step %s first. Compare its selector with the latest DOM snapshot and screenshot, then either switch to a better selector or add wait_for_selector/retry before interacting.", flowRepairHintStepLabel(stepPath))
	case "text_mismatch":
		return fmt.Sprintf("Inspect step %s first. Compare the expected text with the current page content and update assert_text/extract_text or add wait_for_text if the content now appears later.", flowRepairHintStepLabel(stepPath))
	case "extraction_pattern":
		return fmt.Sprintf("Inspect step %s first. Re-check the extraction regex against the current text and keep save_as names stable while adjusting the pattern.", flowRepairHintStepLabel(stepPath))
	case "polling_timeout":
		return fmt.Sprintf("Inspect step %s first. Confirm the wait_until condition can actually become truthy, then adjust the condition, timeout, interval_ms, or the preceding page-state step.", flowRepairHintStepLabel(stepPath))
	case "navigation":
		return fmt.Sprintf("Inspect step %s and the navigation state around it first. Confirm the flow reached the expected page before the next interaction or assertion ran.", flowRepairHintStepLabel(stepPath))
	case "variable_resolution":
		return fmt.Sprintf("Inspect step %s first. Reconnect the missing variable to an earlier save_as/set_var producer and keep downstream variable names stable.", flowRepairHintStepLabel(stepPath))
	default:
		return fmt.Sprintf("Inspect step %s first, compare it with the failure artifacts, and apply the smallest repair that restores the expected page state.", flowRepairHintStepLabel(stepPath))
	}
}

func runtimeFlowRepairPreviousStepHint(context *FlowRepairContext) *FlowRepairHint {
	if context == nil {
		return nil
	}
	switch context.FailureCategory {
	case "selector_or_timing", "text_mismatch", "polling_timeout", "navigation":
	default:
		return nil
	}
	for _, nearby := range context.NearbySteps {
		if nearby.Relation != "previous" {
			continue
		}
		stepPath := firstNonEmpty(strings.TrimSpace(nearby.Path), fmt.Sprint(nearby.Index))
		return &FlowRepairHint{
			Priority:        2,
			Source:          "runtime_failure",
			StepPath:        stepPath,
			Action:          nearby.Step.Action,
			Name:            nearby.Step.Name,
			Selector:        nearby.Step.Selector,
			Targets:         []string{"page_state", "timing"},
			Reason:          "The failed step may depend on page state prepared by the previous step.",
			Suggestion:      fmt.Sprintf("Also inspect step %s to confirm it still brings the page into the state expected by the failed step.", flowRepairHintStepLabel(stepPath)),
			Error:           context.Error,
			FailureCategory: context.FailureCategory,
			PageURL:         runtimeFlowRepairHintPageURL(nearby),
		}
	}
	return nil
}

func runtimeFlowRepairArtifactSelectorHint(context *FlowRepairContext, failedTrace FlowStepTrace, failedStep *FlowStep) *FlowRepairHint {
	if context == nil || context.Artifacts == nil || len(context.Artifacts.RelevantSelectors) == 0 {
		return nil
	}
	switch context.FailureCategory {
	case "selector_or_timing", "text_mismatch", "step_execution":
	default:
		return nil
	}

	currentSelector := ""
	if failedStep != nil {
		currentSelector = strings.TrimSpace(failedStep.Selector)
	}
	filtered := make([]string, 0, len(context.Artifacts.RelevantSelectors))
	for _, selector := range context.Artifacts.RelevantSelectors {
		selector = strings.TrimSpace(selector)
		if selector == "" || selector == currentSelector {
			continue
		}
		filtered = append(filtered, selector)
	}
	if len(filtered) == 0 {
		return nil
	}
	candidates := filtered
	if len(candidates) > 3 {
		candidates = candidates[:3]
	}

	reason := "Failure artifacts include nearby DOM nodes that overlap with the failed step."
	if len(context.Artifacts.RelevantDOM) > 0 {
		reason = reason + " Closest match: " + context.Artifacts.RelevantDOM[0]
	}
	paths := context.Artifacts.Paths
	return &FlowRepairHint{
		Priority: 1,
		Source:   "artifact_analysis",
		StepPath: context.FailedStepPath,
		Action: firstNonEmpty(failedTrace.Action, func() string {
			if failedStep != nil {
				return failedStep.Action
			}
			return ""
		}()),
		Name: func() string {
			if failedStep != nil {
				return failedStep.Name
			}
			return failedTrace.Name
		}(),
		Selector:        candidates[0],
		Targets:         []string{"selector", "page_state"},
		Reason:          reason,
		Suggestion:      fmt.Sprintf("Compare the failed selector with these artifact-derived candidates: %s. Prefer the strongest non-XPath option that still matches the screenshot and DOM snapshot.", strings.Join(candidates, ", ")),
		Error:           context.Error,
		FailureCategory: context.FailureCategory,
		PageURL:         failedTrace.PageURL,
		Artifacts:       &paths,
	}
}

func runtimeFlowRepairHintPageURL(step FlowRepairStepContext) string {
	if step.Trace == nil {
		return ""
	}
	return step.Trace.PageURL
}

func buildFlowRepairPrompt(context *FlowRepairContext) string {
	if context == nil {
		return ""
	}
	failed := "unknown"
	if context.FailedStep != nil {
		failed = fmt.Sprintf("%s action=%s", flowRepairHintStepLabel(firstNonEmpty(context.FailedStep.Path, fmt.Sprint(context.FailedStep.Index))), context.FailedStep.Step.Action)
		if context.FailedStep.Step.Name != "" {
			failed += fmt.Sprintf(" name=%q", context.FailedStep.Step.Name)
		}
	}
	return fmt.Sprintf(`Repair this TSPlay Flow using the provided context.
Return only the corrected Flow YAML.
Use schema_version %q.
Prefer structured actions and stable selectors; use lua only as an explicit escape hatch.
Do not paste full HTML into the answer. Use html_path, screenshot_path, and dom_snapshot_excerpt as evidence.
Failed step: %s
Failure category: %s
Error: %s
Repair hints:
%s`, CurrentFlowSchemaVersion, failed, context.FailureCategory, context.Error, formatFlowRepairHintsForPrompt(context.RepairHints))
}

func formatFlowRepairHintsForPrompt(hints []FlowRepairHint) string {
	if len(hints) == 0 {
		return "- No structured repair hints available."
	}
	lines := make([]string, 0, len(hints))
	for _, hint := range hints {
		line := fmt.Sprintf("- priority=%d", hint.Priority)
		if strings.TrimSpace(hint.StepPath) != "" {
			line += fmt.Sprintf(" step=%s", hint.StepPath)
		}
		if strings.TrimSpace(hint.Action) != "" {
			line += fmt.Sprintf(" action=%s", hint.Action)
		}
		if strings.TrimSpace(hint.Selector) != "" {
			line += fmt.Sprintf(" selector=%q", hint.Selector)
		}
		if strings.TrimSpace(hint.Reason) != "" {
			line += fmt.Sprintf(" reason=%s", hint.Reason)
		}
		if strings.TrimSpace(hint.Suggestion) != "" {
			line += fmt.Sprintf(" suggestion=%s", hint.Suggestion)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func classifyFlowRepairFailure(trace FlowStepTrace, step *FlowStep) (string, string) {
	action := trace.Action
	if step != nil && strings.TrimSpace(step.Action) != "" {
		action = step.Action
	}
	errorText := strings.ToLower(firstNonEmpty(trace.Error, trace.ErrorStack))
	switch {
	case strings.Contains(errorText, "unknown flow variable"):
		return "variable_resolution", "A referenced flow variable is missing or saved under the wrong name."
	case action == "assert_text" || strings.Contains(errorText, "assert_text failed"):
		return "text_mismatch", "Expected text no longer matches the page content or appears later than expected."
	case action == "extract_text" && strings.Contains(errorText, "pattern"):
		return "extraction_pattern", "The extraction regex no longer matches the text returned by the page."
	case action == "wait_until":
		return "polling_timeout", "The condition never became truthy before timeout."
	case action == "navigate":
		return "navigation", "Navigation did not reach the expected page or failed before the next state appeared."
	case strings.Contains(errorText, "timeout"),
		strings.Contains(errorText, "locator"),
		strings.Contains(errorText, "not visible"),
		strings.Contains(errorText, "assert_visible failed"):
		return "selector_or_timing", "The selector may have drifted, the element may be hidden, or the page was not ready yet."
	default:
		return "step_execution", "The step failed during normal execution and likely needs a small targeted repair."
	}
}

func buildFlowRepairTargets(step *FlowStep, trace FlowStepTrace) []string {
	targets := map[string]bool{}
	if step != nil {
		if strings.TrimSpace(step.Selector) != "" {
			targets["selector"] = true
		}
		if strings.TrimSpace(step.Text) != "" {
			targets["text"] = true
		}
		if strings.TrimSpace(step.Pattern) != "" {
			targets["pattern"] = true
		}
		if strings.TrimSpace(step.SaveAs) != "" {
			targets["save_as"] = true
		}
		if len(flowReferences(step.presentNamedParams())) > 0 {
			targets["variables"] = true
		}
		if step.Timeout != 0 || step.IntervalMS != 0 {
			targets["timing"] = true
		}
	}
	if strings.TrimSpace(trace.PageURL) != "" {
		targets["page_state"] = true
	}
	if trace.Condition != nil {
		targets["condition"] = true
	}
	items := make([]string, 0, len(targets))
	for target := range targets {
		items = append(items, target)
	}
	sort.Strings(items)
	return items
}

func buildFlowRepairFocusedVariables(flow *Flow, failedLocation *flowRepairStepLocation, failedStepPath string, vars map[string]any) map[string]any {
	if flow == nil || len(vars) == 0 {
		return nil
	}

	refs := map[string]bool{}
	location := failedLocation
	if location == nil {
		location = resolveFlowRepairStepLocation(flow, failedStepPath)
	}
	if location != nil && location.Step != nil {
		addFlowRepairStepReferences(refs, *location.Step)
	}
	if location != nil && location.Sequence != nil && location.Position > 0 {
		start := location.Position - 1
		if start < 1 {
			start = 1
		}
		end := location.Position + 1
		if end > len(location.Sequence) {
			end = len(location.Sequence)
		}
		for index := start; index <= end; index++ {
			addFlowRepairStepReferences(refs, location.Sequence[index-1])
		}
	} else {
		topLevelIndex := topLevelFlowRepairStepNumber(failedStepPath)
		if topLevelIndex <= 0 || topLevelIndex > len(flow.Steps) {
			topLevelIndex = 1
		}
		start := topLevelIndex - 1
		if start < 1 {
			start = 1
		}
		end := topLevelIndex + 1
		if end > len(flow.Steps) {
			end = len(flow.Steps)
		}
		for index := start; index <= end; index++ {
			addFlowRepairStepReferences(refs, flow.Steps[index-1])
		}
	}
	if len(refs) == 0 {
		return nil
	}

	focused := map[string]any{}
	for ref := range refs {
		value, ok := vars[ref]
		if !ok {
			continue
		}
		focused[ref] = compactTraceValue(value, 0)
	}
	if len(focused) == 0 {
		return nil
	}
	return focused
}

func addFlowRepairStepReferences(refs map[string]bool, step FlowStep) {
	for _, ref := range flowReferences(step.presentNamedParams()) {
		refs[ref] = true
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
