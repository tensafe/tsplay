package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	FlowName           string                  `json:"flow_name,omitempty"`
	FlowDescription    string                  `json:"flow_description,omitempty"`
	Error              string                  `json:"error,omitempty"`
	FailedStep         *FlowRepairStepContext  `json:"failed_step,omitempty"`
	NearbySteps        []FlowRepairStepContext `json:"nearby_steps,omitempty"`
	TraceSummary       []FlowRepairTraceItem   `json:"trace_summary,omitempty"`
	Artifacts          *FlowRepairArtifacts    `json:"artifacts,omitempty"`
	Variables          map[string]any          `json:"variables,omitempty"`
	AllowedActions     []string                `json:"allowed_actions,omitempty"`
	GenerationRules    []string                `json:"generation_rules,omitempty"`
	RepairInstructions []string                `json:"repair_instructions,omitempty"`
	Prompt             string                  `json:"prompt,omitempty"`
}

type FlowRepairStepContext struct {
	Index    int                  `json:"index"`
	Relation string               `json:"relation,omitempty"`
	Step     FlowStep             `json:"step"`
	Trace    *FlowRepairTraceItem `json:"trace,omitempty"`
}

type FlowRepairTraceItem struct {
	Index         int                      `json:"index"`
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

	traceByIndex := map[int]FlowRepairTraceItem{}
	for _, item := range traceSummary {
		traceByIndex[item.Index] = item
	}

	context := &FlowRepairContext{
		FlowName:           options.Flow.Name,
		FlowDescription:    options.Flow.Description,
		Error:              firstNonEmpty(options.Error, failedTrace.Error),
		TraceSummary:       traceSummary,
		NearbySteps:        buildFlowRepairNearbySteps(options.Flow, failedStepNumber, traceByIndex),
		AllowedActions:     FlowActionNames(),
		GenerationRules:    flowRepairGenerationRules(),
		RepairInstructions: flowRepairInstructions(),
	}
	if options.Result != nil {
		if vars, ok := compactTraceValue(options.Result.Vars, 0).(map[string]any); ok && len(vars) > 0 {
			context.Variables = vars
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
	return []string{
		`Keep schema_version: "1".`,
		"Prefer structured actions over lua for navigation, clicking, typing, waiting, and extraction.",
		"Repair selectors using DOM snapshot evidence and selector candidates when available.",
		"Prefer stable selectors: data-testid, data-cy, id, placeholder, aria-label, role/text; use XPath only when necessary.",
		"Keep steps small and named so future trace artifacts identify the failure location.",
		"Do not inline page.html contents into the repaired Flow; use artifact paths only as evidence.",
	}
}

func flowRepairInstructions() []string {
	return []string{
		"Identify whether the failure is caused by selector drift, timing, navigation state, hidden elements, changed text, or missing variables.",
		"Patch only the smallest necessary part of the Flow.",
		"Add wait_for_selector or wait_for_text before fragile click/type/extract steps when the page is dynamic.",
		"Preserve existing save_as variables and downstream variable references unless they are the actual bug.",
		"Return a valid TSPlay Flow YAML/JSON that passes tsplay.validate_flow.",
	}
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
Error: %s`, CurrentFlowSchemaVersion, failed, context.Error)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
