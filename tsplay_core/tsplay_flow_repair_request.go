package tsplay_core

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type FlowRepairRequestOptions struct {
	Flow        *Flow
	RepairHints []FlowRepairHint
	Context     *FlowRepairContext
}

type FlowRepairRequest struct {
	FlowName            string             `json:"flow_name,omitempty"`
	FlowYAML            string             `json:"flow_yaml,omitempty"`
	RepairHints         []FlowRepairHint   `json:"repair_hints,omitempty"`
	TargetSteps         []string           `json:"target_steps,omitempty"`
	RepairContext       *FlowRepairContext `json:"repair_context,omitempty"`
	AllowedActions      []string           `json:"allowed_actions,omitempty"`
	GenerationRules     []string           `json:"generation_rules,omitempty"`
	SelectorStrategy    []string           `json:"selector_strategy,omitempty"`
	ValidationChecklist []string           `json:"validation_checklist,omitempty"`
	OutputRequirements  []string           `json:"output_requirements,omitempty"`
	Prompt              string             `json:"prompt,omitempty"`
}

func BuildFlowRepairRequest(options FlowRepairRequestOptions) (*FlowRepairRequest, error) {
	if options.Flow == nil {
		return nil, fmt.Errorf("flow is required")
	}

	hints := mergeFlowRepairHints(options.RepairHints, flowRepairHintsFromContext(options.Context))
	if len(hints) == 0 {
		return nil, fmt.Errorf("repair_hints, repair_context, or failed scene is required")
	}

	encoded, err := yaml.Marshal(options.Flow)
	if err != nil {
		return nil, fmt.Errorf("marshal flow yaml: %w", err)
	}

	request := &FlowRepairRequest{
		FlowName:            options.Flow.Name,
		FlowYAML:            string(encoded),
		RepairHints:         hints,
		TargetSteps:         flowRepairTargetSteps(hints),
		RepairContext:       options.Context,
		AllowedActions:      flowRepairAllowedActions(options.Context),
		GenerationRules:     flowRepairGenerationRules(),
		SelectorStrategy:    flowSelectorStrategy(),
		ValidationChecklist: flowRepairValidationChecklist(),
		OutputRequirements:  flowRepairOutputRequirements(),
	}
	request.Prompt = buildFlowRepairRequestPrompt(request)
	return request, nil
}

func ParseFlowRepairHintsInput(text string) ([]FlowRepairHint, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}

	var hints []FlowRepairHint
	if err := json.Unmarshal([]byte(text), &hints); err == nil {
		return hints, nil
	}

	var wrapper struct {
		RepairHints []FlowRepairHint `json:"repair_hints"`
		Draft       *struct {
			RepairHints []FlowRepairHint `json:"repair_hints"`
		} `json:"draft"`
		Context *struct {
			RepairHints []FlowRepairHint `json:"repair_hints"`
		} `json:"context"`
		Repair *struct {
			RepairHints []FlowRepairHint `json:"repair_hints"`
		} `json:"repair"`
	}
	if err := json.Unmarshal([]byte(text), &wrapper); err != nil {
		return nil, fmt.Errorf("repair_hints must be a JSON array or wrapper with repair_hints: %w", err)
	}
	switch {
	case wrapper.RepairHints != nil:
		return wrapper.RepairHints, nil
	case wrapper.Draft != nil:
		return wrapper.Draft.RepairHints, nil
	case wrapper.Context != nil:
		return wrapper.Context.RepairHints, nil
	case wrapper.Repair != nil:
		return wrapper.Repair.RepairHints, nil
	default:
		return nil, fmt.Errorf("repair_hints wrapper does not contain repair_hints")
	}
}

func ParseFlowRepairContextInput(text string) (*FlowRepairContext, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}

	var context FlowRepairContext
	if err := json.Unmarshal([]byte(text), &context); err == nil && flowRepairContextLooksPresent(&context) {
		return &context, nil
	}

	var wrapper struct {
		Context *FlowRepairContext `json:"context"`
		Repair  *struct {
			RepairContext *FlowRepairContext `json:"repair_context"`
		} `json:"repair"`
	}
	if err := json.Unmarshal([]byte(text), &wrapper); err != nil {
		return nil, fmt.Errorf("repair_context must be a JSON FlowRepairContext or wrapper with context: %w", err)
	}
	switch {
	case wrapper.Context != nil:
		return wrapper.Context, nil
	case wrapper.Repair != nil && wrapper.Repair.RepairContext != nil:
		return wrapper.Repair.RepairContext, nil
	default:
		return nil, fmt.Errorf("repair_context wrapper does not contain context")
	}
}

func buildFlowRepairRequestPrompt(request *FlowRepairRequest) string {
	if request == nil {
		return ""
	}

	lines := []string{
		"Repair this TSPlay Flow.",
		"Return only the corrected Flow YAML.",
		fmt.Sprintf("Use schema_version %q.", CurrentFlowSchemaVersion),
		"Prefer structured actions and stable selectors; use lua only as an explicit escape hatch.",
		"Do not paste full HTML into the answer. Use artifact paths and DOM snapshot excerpts as evidence.",
	}
	if request.FlowYAML != "" {
		lines = append(lines, "Original flow:", "```yaml", strings.TrimRight(request.FlowYAML, "\n"), "```")
	}
	lines = append(lines, "Prioritized repair hints:", formatFlowRepairHintsForPrompt(request.RepairHints))
	if summary := formatFlowRepairContextSummaryForPrompt(request.RepairContext); summary != "" {
		lines = append(lines, "Failure context:", summary)
	}
	lines = append(lines, "Generation rules:", formatStringListForPrompt(request.GenerationRules))
	lines = append(lines, "Selector strategy:", formatStringListForPrompt(request.SelectorStrategy))
	lines = append(lines, "Validation checklist:", formatStringListForPrompt(request.ValidationChecklist))
	lines = append(lines, "Output requirements:", formatStringListForPrompt(request.OutputRequirements))
	return strings.Join(lines, "\n")
}

func formatFlowRepairContextSummaryForPrompt(context *FlowRepairContext) string {
	if context == nil {
		return ""
	}
	lines := []string{}
	if context.FailureCategory != "" {
		lines = append(lines, fmt.Sprintf("- failure_category=%s", context.FailureCategory))
	}
	if context.FailureReason != "" {
		lines = append(lines, fmt.Sprintf("- failure_reason=%s", context.FailureReason))
	}
	if context.FailedStepPath != "" {
		lines = append(lines, fmt.Sprintf("- failed_step_path=%s", context.FailedStepPath))
	}
	if len(context.RepairTargets) > 0 {
		lines = append(lines, fmt.Sprintf("- repair_targets=%s", strings.Join(context.RepairTargets, ", ")))
	}
	if len(context.FocusedVariables) > 0 {
		encoded, err := json.Marshal(context.FocusedVariables)
		if err == nil {
			lines = append(lines, fmt.Sprintf("- focused_variables=%s", string(encoded)))
		}
	}
	if context.Artifacts != nil {
		if context.Artifacts.Paths.ScreenshotPath != "" {
			lines = append(lines, fmt.Sprintf("- screenshot_path=%s", context.Artifacts.Paths.ScreenshotPath))
		}
		if context.Artifacts.Paths.HTMLPath != "" {
			lines = append(lines, fmt.Sprintf("- html_path=%s", context.Artifacts.Paths.HTMLPath))
		}
		if context.Artifacts.Paths.DOMSnapshotPath != "" {
			lines = append(lines, fmt.Sprintf("- dom_snapshot_path=%s", context.Artifacts.Paths.DOMSnapshotPath))
		}
		if context.Artifacts.DOMSnapshotExcerpt != "" {
			lines = append(lines, fmt.Sprintf("- dom_snapshot_excerpt=%s", context.Artifacts.DOMSnapshotExcerpt))
		}
	}
	return strings.Join(lines, "\n")
}

func formatStringListForPrompt(items []string) string {
	if len(items) == 0 {
		return "- none"
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item) == "" {
			continue
		}
		lines = append(lines, "- "+item)
	}
	if len(lines) == 0 {
		return "- none"
	}
	return strings.Join(lines, "\n")
}

func flowRepairOutputRequirements() []string {
	return []string{
		"Return only corrected Flow YAML.",
		fmt.Sprintf("Keep schema_version set to %q.", CurrentFlowSchemaVersion),
		"Preserve existing save_as outputs and downstream variable names unless they are the actual bug.",
		"Prefer the smallest repair that matches the evidence in repair_hints and failure context.",
		"Do not inline full HTML or screenshots; refer to artifact paths and DOM snapshot excerpts instead.",
	}
}

func flowRepairHintsFromContext(context *FlowRepairContext) []FlowRepairHint {
	if context == nil {
		return nil
	}
	if len(context.RepairHints) > 0 {
		return append([]FlowRepairHint(nil), context.RepairHints...)
	}
	if context.FailedStep == nil && context.FailedStepPath == "" {
		return nil
	}

	fallback := FlowRepairHint{
		Priority:        0,
		Source:          "runtime_failure",
		StepPath:        context.FailedStepPath,
		Targets:         append([]string(nil), context.RepairTargets...),
		Reason:          firstNonEmpty(context.FailureReason, "Inspect the failed step and apply the smallest repair that matches the failure context."),
		Suggestion:      runtimeFlowRepairSuggestion(context, context.FailedStepPath),
		Error:           context.Error,
		FailureCategory: context.FailureCategory,
	}
	if context.FailedStep != nil {
		fallback.Action = context.FailedStep.Step.Action
		fallback.Name = context.FailedStep.Step.Name
		fallback.Selector = context.FailedStep.Step.Selector
		if context.FailedStep.Trace != nil {
			fallback.PageURL = context.FailedStep.Trace.PageURL
		}
	}
	if context.Artifacts != nil {
		paths := context.Artifacts.Paths
		fallback.Artifacts = &paths
	}
	return []FlowRepairHint{fallback}
}

func mergeFlowRepairHints(groups ...[]FlowRepairHint) []FlowRepairHint {
	merged := make([]FlowRepairHint, 0)
	for _, group := range groups {
		merged = append(merged, group...)
	}
	if len(merged) == 0 {
		return nil
	}
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Priority != merged[j].Priority {
			return merged[i].Priority < merged[j].Priority
		}
		if merged[i].StepPath != merged[j].StepPath {
			return merged[i].StepPath < merged[j].StepPath
		}
		if merged[i].Source != merged[j].Source {
			return merged[i].Source < merged[j].Source
		}
		return merged[i].Reason < merged[j].Reason
	})
	return dedupeFlowRepairHints(merged)
}

func flowRepairTargetSteps(hints []FlowRepairHint) []string {
	seen := map[string]bool{}
	steps := make([]string, 0, len(hints))
	for _, hint := range hints {
		path := strings.TrimSpace(hint.StepPath)
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		steps = append(steps, path)
	}
	if len(steps) == 0 {
		return nil
	}
	return steps
}

func flowRepairAllowedActions(context *FlowRepairContext) []string {
	if context != nil && len(context.AllowedActions) > 0 {
		return append([]string(nil), context.AllowedActions...)
	}
	return FlowActionNames()
}

func flowRepairContextLooksPresent(context *FlowRepairContext) bool {
	if context == nil {
		return false
	}
	return context.FlowName != "" ||
		context.FlowDescription != "" ||
		context.Error != "" ||
		context.FailureCategory != "" ||
		context.FailedStepPath != "" ||
		len(context.RepairHints) > 0 ||
		context.FailedStep != nil
}
