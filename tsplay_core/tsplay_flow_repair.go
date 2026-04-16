package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	Relation string               `json:"relation,omitempty"`
	Step     FlowStep             `json:"step"`
	Trace    *FlowRepairTraceItem `json:"trace,omitempty"`
}

type FlowRepairTraceItem struct {
	Index         int                      `json:"index"`
	Path          string                   `json:"path,omitempty"`
	Attempt       int                      `json:"attempt,omitempty"`
	Iteration     int                      `json:"iteration,omitempty"`
	Branch        string                   `json:"branch,omitempty"`
	Name          string                   `json:"name,omitempty"`
	Action        string                   `json:"action"`
	Status        string                   `json:"status"`
	SaveAs        string                   `json:"save_as,omitempty"`
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
	DOMSnapshotExcerpt   string                  `json:"dom_snapshot_excerpt,omitempty"`
	DOMSnapshotTruncated bool                    `json:"dom_snapshot_truncated,omitempty"`
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
	failedTraceIndex := findFailedTraceIndex(trace)
	failedTrace := trace[failedTraceIndex]
	failedStepNumber := failedTrace.Index
	if failedStepNumber <= 0 {
		failedStepNumber = failedTraceIndex + 1
	}
	var failedFlowStep *FlowStep
	if failedStepNumber >= 1 && failedStepNumber <= len(options.Flow.Steps) {
		failedFlowStep = &options.Flow.Steps[failedStepNumber-1]
	}

	traceByIndex := map[int]FlowRepairTraceItem{}
	for _, item := range traceSummary {
		traceByIndex[item.Index] = item
	}
	failureCategory, failureReason := classifyFlowRepairFailure(failedTrace, failedFlowStep)
	failedStepPath := strings.TrimSpace(failedTrace.Path)
	if failedStepPath == "" && failedStepNumber > 0 {
		failedStepPath = fmt.Sprint(failedStepNumber)
	}

	context := &FlowRepairContext{
		FlowName:            options.Flow.Name,
		FlowDescription:     options.Flow.Description,
		Error:               firstNonEmpty(options.Error, failedTrace.Error),
		FailureCategory:     failureCategory,
		FailureReason:       failureReason,
		FailedStepPath:      failedStepPath,
		TraceSummary:        traceSummary,
		NearbySteps:         buildFlowRepairNearbySteps(options.Flow, failedStepNumber, traceByIndex),
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
		if focused := buildFlowRepairFocusedVariables(options.Flow, failedStepNumber, options.Result.Vars); len(focused) > 0 {
			context.FocusedVariables = focused
		}
	}
	if failedStepNumber >= 1 && failedStepNumber <= len(options.Flow.Steps) {
		failed := FlowRepairStepContext{
			Index:    failedStepNumber,
			Relation: "failed",
			Step:     options.Flow.Steps[failedStepNumber-1],
		}
		if item, ok := traceByIndex[failedStepNumber]; ok {
			failed.Trace = &item
		}
		context.FailedStep = &failed
	}
	if failedTrace.Artifacts != nil {
		context.Artifacts = buildFlowRepairArtifacts(*failedTrace.Artifacts, artifactRoot, excerptLimit)
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

func findFailedTraceIndex(trace []FlowStepTrace) int {
	for index, step := range trace {
		status := strings.ToLower(strings.TrimSpace(step.Status))
		if status == "error" || status == "failed" || status == "failure" {
			return index
		}
	}
	return len(trace) - 1
}

func buildFlowRepairNearbySteps(flow *Flow, failedStepNumber int, traceByIndex map[int]FlowRepairTraceItem) []FlowRepairStepContext {
	if flow == nil || len(flow.Steps) == 0 || failedStepNumber <= 0 {
		return nil
	}
	start := failedStepNumber - 2
	if start < 1 {
		start = 1
	}
	end := failedStepNumber + 2
	if end > len(flow.Steps) {
		end = len(flow.Steps)
	}

	steps := make([]FlowRepairStepContext, 0, end-start+1)
	for index := start; index <= end; index++ {
		relation := "next"
		if index < failedStepNumber {
			relation = "previous"
		} else if index == failedStepNumber {
			relation = "failed"
		}
		item := FlowRepairStepContext{
			Index:    index,
			Relation: relation,
			Step:     flow.Steps[index-1],
		}
		if trace, ok := traceByIndex[index]; ok {
			item.Trace = &trace
		}
		steps = append(steps, item)
	}
	return steps
}

func buildFlowRepairArtifacts(artifacts FlowStepArtifacts, artifactRoot string, excerptLimit int) *FlowRepairArtifacts {
	context := &FlowRepairArtifacts{
		Paths: flowRepairArtifactPaths(artifacts),
	}
	if artifacts.DOMSnapshotPath == "" {
		return context
	}
	excerpt, truncated, err := readFlowRepairArtifactExcerpt(artifacts.DOMSnapshotPath, artifactRoot, excerptLimit)
	if err != nil {
		context.ReadErrors = append(context.ReadErrors, err.Error())
		return context
	}
	context.DOMSnapshotExcerpt = excerpt
	context.DOMSnapshotTruncated = truncated
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

func readFlowRepairArtifactExcerpt(path string, artifactRoot string, limit int) (string, bool, error) {
	if strings.TrimSpace(path) == "" {
		return "", false, nil
	}
	if strings.TrimSpace(artifactRoot) == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	rootAbs, err := filepath.Abs(artifactRoot)
	if err != nil {
		return "", false, fmt.Errorf("resolve artifact root %q: %w", artifactRoot, err)
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", false, fmt.Errorf("artifact root %q is not accessible: %w", rootAbs, err)
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootReal, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return "", false, fmt.Errorf("resolve artifact path %q: %w", path, err)
	}
	candidateReal, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", false, fmt.Errorf("artifact path %q is not accessible: %w", path, err)
	}
	if err := ensurePathInsideRoot(candidateReal, rootReal); err != nil {
		return "", false, fmt.Errorf("artifact path %q is outside allowed artifact root %q", path, rootReal)
	}

	content, err := os.ReadFile(candidateReal)
	if err != nil {
		return "", false, fmt.Errorf("read artifact path %q: %w", candidateReal, err)
	}
	excerpt, truncated := truncateRepairText(string(content), limit)
	return excerpt, truncated, nil
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
		return &FlowRepairHint{
			Priority:        2,
			Source:          "runtime_failure",
			StepPath:        fmt.Sprint(nearby.Index),
			Action:          nearby.Step.Action,
			Name:            nearby.Step.Name,
			Selector:        nearby.Step.Selector,
			Targets:         []string{"page_state", "timing"},
			Reason:          "The failed step may depend on page state prepared by the previous step.",
			Suggestion:      fmt.Sprintf("Also inspect step %d to confirm it still brings the page into the state expected by the failed step.", nearby.Index),
			Error:           context.Error,
			FailureCategory: context.FailureCategory,
			PageURL:         runtimeFlowRepairHintPageURL(nearby),
		}
	}
	return nil
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
		failed = fmt.Sprintf("#%d action=%s", context.FailedStep.Index, context.FailedStep.Step.Action)
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

func buildFlowRepairFocusedVariables(flow *Flow, failedStepNumber int, vars map[string]any) map[string]any {
	if flow == nil || failedStepNumber <= 0 || len(vars) == 0 {
		return nil
	}
	start := failedStepNumber - 1
	if start < 1 {
		start = 1
	}
	end := failedStepNumber + 1
	if end > len(flow.Steps) {
		end = len(flow.Steps)
	}

	refs := map[string]bool{}
	for index := start; index <= end; index++ {
		for _, ref := range flowReferences(flow.Steps[index-1].presentNamedParams()) {
			refs[ref] = true
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
