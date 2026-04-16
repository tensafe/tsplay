package tsplay_core

import (
	"strconv"
	"strings"
)

type FlowRepairHint struct {
	Priority        int                      `json:"priority"`
	Source          string                   `json:"source,omitempty"`
	StepPath        string                   `json:"step_path,omitempty"`
	Action          string                   `json:"action,omitempty"`
	Name            string                   `json:"name,omitempty"`
	Selector        string                   `json:"selector,omitempty"`
	Targets         []string                 `json:"targets,omitempty"`
	Reason          string                   `json:"reason"`
	Suggestion      string                   `json:"suggestion,omitempty"`
	Error           string                   `json:"error,omitempty"`
	FailureCategory string                   `json:"failure_category,omitempty"`
	PageURL         string                   `json:"page_url,omitempty"`
	Artifacts       *FlowRepairArtifactPaths `json:"artifacts,omitempty"`
}

func flowRepairHintStepLabel(stepPath string) string {
	if strings.TrimSpace(stepPath) == "" {
		return "the reported step"
	}
	return stepPath
}

func dedupeFlowRepairHints(hints []FlowRepairHint) []FlowRepairHint {
	seen := map[string]bool{}
	deduped := make([]FlowRepairHint, 0, len(hints))
	for _, hint := range hints {
		key := strings.Join([]string{
			strconv.Itoa(hint.Priority),
			hint.Source,
			hint.StepPath,
			hint.Action,
			hint.Reason,
			hint.Suggestion,
			hint.Error,
			hint.FailureCategory,
		}, "|")
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, hint)
	}
	return deduped
}
