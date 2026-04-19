package tsplay_core

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	tsplaySecurityPresetReadOnly       = "readonly"
	tsplaySecurityPresetBrowserWrite   = "browser_write"
	tsplaySecurityPresetFullAutomation = "full_automation"
)

type tsplayFlowSecurityResolution struct {
	Preset string             `json:"preset"`
	Policy FlowSecurityPolicy `json:"policy"`
}

func flowSecurityPresetNames() []string {
	return []string{
		tsplaySecurityPresetReadOnly,
		tsplaySecurityPresetBrowserWrite,
		tsplaySecurityPresetFullAutomation,
	}
}

func flowSecurityPolicyResolutionFromToolRequest(
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (tsplayFlowSecurityResolution, error) {
	preset := strings.TrimSpace(request.GetString("security_preset", ""))
	policy, ok := flowSecurityPolicyPreset(preset)
	if preset != "" && !ok {
		return tsplayFlowSecurityResolution{}, fmt.Errorf(
			"security_preset %q is unsupported; use readonly, browser_write, or full_automation",
			preset,
		)
	}

	if _, present := request.GetArguments()["allow_lua"]; present {
		policy.AllowLua = request.GetBool("allow_lua", false)
	}
	if _, present := request.GetArguments()["allow_javascript"]; present {
		policy.AllowJavaScript = request.GetBool("allow_javascript", false)
	}
	if _, present := request.GetArguments()["allow_file_access"]; present {
		policy.AllowFileAccess = request.GetBool("allow_file_access", false)
	}
	if _, present := request.GetArguments()["allow_browser_state"]; present {
		policy.AllowBrowserState = request.GetBool("allow_browser_state", false)
	}
	if _, present := request.GetArguments()["allow_http"]; present {
		policy.AllowHTTP = request.GetBool("allow_http", false)
	}
	if _, present := request.GetArguments()["allow_redis"]; present {
		policy.AllowRedis = request.GetBool("allow_redis", false)
	}
	if _, present := request.GetArguments()["allow_database"]; present {
		policy.AllowDatabase = request.GetBool("allow_database", false)
	}

	policy.FileInputRoot = options.ArtifactRoot
	policy.FileOutputRoot = options.ArtifactRoot

	return tsplayFlowSecurityResolution{
		Preset: preset,
		Policy: policy,
	}, nil
}

func flowSecurityPolicyPreset(name string) (FlowSecurityPolicy, bool) {
	switch strings.TrimSpace(name) {
	case "", tsplaySecurityPresetReadOnly:
		return FlowSecurityPolicy{}, true
	case tsplaySecurityPresetBrowserWrite:
		return FlowSecurityPolicy{
			AllowFileAccess:   true,
			AllowBrowserState: true,
		}, true
	case tsplaySecurityPresetFullAutomation:
		return FlowSecurityPolicy{
			AllowLua:          true,
			AllowJavaScript:   true,
			AllowFileAccess:   true,
			AllowBrowserState: true,
			AllowHTTP:         true,
			AllowRedis:        true,
			AllowDatabase:     true,
		}, true
	default:
		return FlowSecurityPolicy{}, false
	}
}

func newTSPlayToolResult(tool string, payload map[string]any) (*mcp.CallToolResult, error) {
	return newJSONToolResult(enrichTSPlayToolPayload(tool, payload))
}

func enrichTSPlayToolPayload(tool string, payload map[string]any) map[string]any {
	if payload == nil {
		payload = map[string]any{}
	}
	payload["tool"] = tool

	if _, ok := payload["ok"]; !ok {
		payload["ok"] = inferTSPlayToolOK(payload)
	}
	if _, ok := payload["summary"]; !ok {
		payload["summary"] = buildTSPlayToolSummary(tool, payload)
	}
	if _, ok := payload["warnings"]; !ok {
		payload["warnings"] = extractTSPlayToolWarnings(tool, payload)
	}
	if _, ok := payload["artifacts"]; !ok {
		payload["artifacts"] = extractTSPlayToolArtifacts(tool, payload)
	}
	if _, ok := payload["next_action"]; !ok {
		payload["next_action"] = buildTSPlayToolNextAction(tool, payload)
	}

	switch run := payload["run"].(type) {
	case TSPlayBrowserRun:
		payload["run"] = normalizeTSPlayBrowserRun(run)
	case *TSPlayBrowserRun:
		if run != nil {
			payload["run"] = normalizeTSPlayBrowserRun(*run)
		}
	}
	if _, ok := payload["run"]; !ok && (tool == "tsplay.draft_flow" || tool == "tsplay.finalize_flow") {
		payload["run"] = map[string]any{
			"id":     "",
			"tool":   tool,
			"status": "not_started",
			"details": map[string]any{
				"source": "provided_observation",
			},
		}
	}

	return payload
}

func inferTSPlayToolOK(payload map[string]any) bool {
	if ok, okPresent := payload["ok"].(bool); okPresent {
		return ok
	}
	if valid, okPresent := payload["valid"].(bool); okPresent {
		return valid
	}
	if errorText, ok := payload["error"].(string); ok && strings.TrimSpace(errorText) != "" {
		return false
	}
	return true
}

func normalizeTSPlayBrowserRun(run TSPlayBrowserRun) map[string]any {
	value := map[string]any{
		"id":               run.ID,
		"tool":             run.Tool,
		"status":           run.Status,
		"queued_at":        run.QueuedAt,
		"started_at":       run.StartedAt,
		"finished_at":      run.FinishedAt,
		"queue_wait_ms":    run.QueueWaitMS,
		"duration_ms":      run.DurationMS,
		"timeout_ms":       run.TimeoutMS,
		"queue_timeout_ms": run.QueueTimeoutMS,
		"artifact_root":    run.ArtifactRoot,
		"run_root":         run.RunRoot,
		"audit_path":       run.AuditPath,
		"caller":           compactTraceValue(run.Caller, 0),
		"grants":           compactTraceValue(run.Grants, 0),
		"details":          compactTraceValue(run.Details, 0),
	}
	if strings.TrimSpace(run.Error) != "" {
		value["error"] = run.Error
	}
	return value
}

func buildTSPlayToolSummary(tool string, payload map[string]any) string {
	switch tool {
	case "tsplay.list_actions":
		return fmt.Sprintf("Returned %d TSPlay Flow action definitions.", countTSPlayItems(payload["actions"]))
	case "tsplay.list_sessions":
		return fmt.Sprintf("Returned %d reusable browser sessions.", countTSPlayItems(payload["sessions"]))
	case "tsplay.get_session":
		if sessionName := tsplayNestedString(payload["session"], "name"); sessionName != "" {
			return fmt.Sprintf("Loaded browser session %q.", sessionName)
		}
		return "Loaded the requested browser session."
	case "tsplay.export_session_flow_snippet":
		if format := tsplayNestedString(payload["export"], "format"); format != "" {
			return fmt.Sprintf("Exported session snippets in %s format.", format)
		}
		return "Exported reusable browser or Flow snippets."
	case "tsplay.delete_session":
		if deletedName := tsplayNestedString(payload["deleted"], "name"); deletedName != "" {
			return fmt.Sprintf("Deleted browser session %q.", deletedName)
		}
		return "Deleted the requested browser session."
	case "tsplay.save_session":
		if sessionName := tsplayNestedString(payload["session"], "name"); sessionName != "" {
			return fmt.Sprintf("Saved browser session %q.", sessionName)
		}
		return "Saved a reusable browser session."
	case "tsplay.flow_schema":
		return "Returned the TSPlay Flow schema, action manifest, and generation rules."
	case "tsplay.flow_examples":
		return fmt.Sprintf("Returned %d TSPlay Flow examples.", countTSPlayItems(payload["examples"]))
	case "tsplay.observe_page":
		if observation, ok := payload["observation"].(*PageObservation); ok && observation != nil {
			return fmt.Sprintf(
				"Observed %d interactive elements on %s.",
				len(observation.Elements),
				firstNonEmpty(observation.URL, "the page"),
			)
		}
		if payload["ok"] == false {
			return "Page observation failed before a usable observation was produced."
		}
		return "Observed the target page and extracted an AI-friendly page snapshot."
	case "tsplay.draft_flow":
		if draft, ok := payload["draft"].(*FlowDraft); ok && draft != nil {
			return fmt.Sprintf(
				"Drafted flow %q with %d planned actions.",
				firstNonEmpty(draft.FlowName, "draft"),
				len(draft.PlannedActions),
			)
		}
		if payload["ok"] == false {
			return "Could not draft a TSPlay Flow from the provided intent and context."
		}
		return "Drafted a TSPlay Flow from the provided intent."
	case "tsplay.finalize_flow":
		status := firstNonEmpty(stringValue(payload["status"]), "drafted")
		if draft, ok := payload["draft"].(*FlowDraft); ok && draft != nil {
			return fmt.Sprintf("Finalized flow %q with status %s.", firstNonEmpty(draft.FlowName, "draft"), status)
		}
		if payload["ok"] == false {
			return "Could not finalize a TSPlay Flow from the provided intent and context."
		}
		return fmt.Sprintf("Finalized the TSPlay Flow with status %s.", status)
	case "tsplay.validate_flow":
		if valid, ok := payload["valid"].(bool); ok && valid {
			return fmt.Sprintf("Validated flow %q successfully.", firstNonEmpty(stringValue(payload["name"]), "flow"))
		}
		if name := stringValue(payload["name"]); name != "" {
			return fmt.Sprintf("Validation failed for flow %q.", name)
		}
		return "Validation did not succeed for the requested flow."
	case "tsplay.run_flow":
		if result, ok := payload["result"].(*FlowResult); ok && result != nil {
			return fmt.Sprintf(
				"Ran flow %q and recorded %d trace steps.",
				firstNonEmpty(result.Name, "flow"),
				len(result.Trace),
			)
		}
		if payload["ok"] == false {
			return "Flow execution failed before completion."
		}
		return "Executed the requested TSPlay Flow."
	case "tsplay.repair_flow_context":
		if contextPayload, ok := payload["context"].(*FlowRepairContext); ok && contextPayload != nil {
			if contextPayload.FailedStepPath != "" {
				return fmt.Sprintf("Built repair context for failed step %s.", contextPayload.FailedStepPath)
			}
		}
		if payload["ok"] == false {
			return "Could not build repair context from the supplied flow and trace."
		}
		return "Built repair context from the supplied flow and trace."
	case "tsplay.repair_flow":
		if repair, ok := payload["repair"].(*FlowRepairRequest); ok && repair != nil {
			return fmt.Sprintf("Built a repair prompt targeting %d flow step(s).", len(repair.TargetSteps))
		}
		if payload["ok"] == false {
			return "Could not build a repair prompt for the supplied flow."
		}
		return "Built a repair prompt for the supplied flow."
	default:
		if payload["ok"] == false {
			return "The tool returned an error."
		}
		return "The tool completed successfully."
	}
}

func extractTSPlayToolWarnings(tool string, payload map[string]any) []string {
	warnings := []string{}
	switch tool {
	case "tsplay.observe_page":
		if observation, ok := payload["observation"].(*PageObservation); ok && observation != nil && len(observation.Errors) > 0 {
			warnings = append(warnings, observation.Errors...)
		}
	case "tsplay.draft_flow":
		fallthrough
	case "tsplay.finalize_flow":
		if draft, ok := payload["draft"].(*FlowDraft); ok && draft != nil {
			warnings = append(warnings, draft.Warnings...)
			if draft.Validation != nil && !draft.Validation.Valid && strings.TrimSpace(draft.Validation.Error) != "" {
				warnings = append(warnings, draft.Validation.Error)
			}
		}
	case "tsplay.repair_flow_context":
		if contextPayload, ok := payload["context"].(*FlowRepairContext); ok && contextPayload != nil && contextPayload.Artifacts != nil {
			warnings = append(warnings, contextPayload.Artifacts.ReadErrors...)
		}
	}
	if len(warnings) == 0 {
		return []string{}
	}
	return warnings
}

func extractTSPlayToolArtifacts(tool string, payload map[string]any) any {
	switch tool {
	case "tsplay.observe_page":
		if observation, ok := payload["observation"].(*PageObservation); ok && observation != nil {
			return map[string]any{
				"artifact_root":     observation.ArtifactRoot,
				"screenshot_path":   observation.ScreenshotPath,
				"dom_snapshot_path": observation.DOMSnapshotPath,
			}
		}
	case "tsplay.draft_flow":
		if observation, ok := payload["observation"].(*PageObservation); ok && observation != nil {
			return map[string]any{
				"artifact_root":     observation.ArtifactRoot,
				"screenshot_path":   observation.ScreenshotPath,
				"dom_snapshot_path": observation.DOMSnapshotPath,
			}
		}
	case "tsplay.finalize_flow":
		if observation, ok := payload["observation"].(*PageObservation); ok && observation != nil {
			return map[string]any{
				"artifact_root":     observation.ArtifactRoot,
				"screenshot_path":   observation.ScreenshotPath,
				"dom_snapshot_path": observation.DOMSnapshotPath,
			}
		}
	case "tsplay.run_flow":
		if result, ok := payload["result"].(*FlowResult); ok && result != nil {
			artifacts := map[string]any{
				"artifact_root": result.ArtifactRoot,
				"run_root":      result.RunRoot,
			}
			if failedArtifacts := latestTSPlayTraceArtifacts(result.Trace); failedArtifacts != nil {
				artifacts["last_step_artifacts"] = failedArtifacts
			}
			return artifacts
		}
	case "tsplay.repair_flow_context":
		if contextPayload, ok := payload["context"].(*FlowRepairContext); ok && contextPayload != nil && contextPayload.Artifacts != nil {
			return contextPayload.Artifacts
		}
	}
	return nil
}

func buildTSPlayToolNextAction(tool string, payload map[string]any) any {
	ok, _ := payload["ok"].(bool)
	switch tool {
	case "tsplay.observe_page":
		if !ok {
			return map[string]any{
				"tool":   "tsplay.observe_page",
				"reason": "Retry with a valid URL, a longer timeout, or headless=false if the page is timing-sensitive.",
			}
		}
		return map[string]any{
			"tool":   "tsplay.draft_flow",
			"reason": "Draft a flow from this observation by passing the user intent and the returned observation object.",
		}
	case "tsplay.draft_flow":
		if !ok {
			return map[string]any{
				"tool":   "tsplay.observe_page",
				"reason": "Re-observe the page or provide an observation payload before drafting again.",
			}
		}
		if draft, ok := payload["draft"].(*FlowDraft); ok && draft != nil {
			if draft.Validation != nil && draft.Validation.Valid {
				return map[string]any{
					"tool":   "tsplay.run_flow",
					"reason": "Run the returned draft.flow_yaml and inspect the trace.",
				}
			}
			return map[string]any{
				"tool":   "tsplay.validate_flow",
				"reason": "Adjust the drafted flow using the validation output, then validate again before execution.",
			}
		}
	case "tsplay.finalize_flow":
		if !ok {
			return map[string]any{
				"tool":   "tsplay.observe_page",
				"reason": "Re-observe the page or provide an observation payload before finalizing again.",
			}
		}
		switch stringValue(payload["status"]) {
		case "ready":
			return map[string]any{
				"tool":   "tsplay.run_flow",
				"reason": "Run the returned flow_yaml and inspect the execution trace.",
			}
		case "needs_permission":
			return map[string]any{
				"tool":   "tsplay.finalize_flow",
				"reason": "Retry with the matching security_preset or allow_* override once the user approves it.",
			}
		case "needs_input":
			return map[string]any{
				"tool":   "tsplay.finalize_flow",
				"reason": "Fill the remaining TODO variables or unresolved inputs, then finalize again.",
			}
		default:
			return map[string]any{
				"tool":   "tsplay.validate_flow",
				"reason": "Inspect the returned issue or validation output, adjust the flow, then validate again.",
			}
		}
	case "tsplay.validate_flow":
		if ok {
			return map[string]any{
				"tool":   "tsplay.run_flow",
				"reason": "Run the validated flow and inspect the execution trace.",
			}
		}
		if errorText := stringValue(payload["error"]); strings.Contains(errorText, "allow_") {
			return map[string]any{
				"tool":   "tsplay.validate_flow",
				"reason": "Retry validation with a matching security_preset or explicit allow_* override for the blocked action.",
			}
		}
		return map[string]any{
			"tool":   "tsplay.repair_flow",
			"reason": "Build a repair prompt if the flow structure or selectors still need to be updated.",
		}
	case "tsplay.run_flow":
		if !ok {
			return map[string]any{
				"tool":   "tsplay.repair_flow_context",
				"reason": "Build repair context from this run result before asking a model to update the flow.",
			}
		}
	case "tsplay.repair_flow_context":
		if ok {
			return map[string]any{
				"tool":   "tsplay.repair_flow",
				"reason": "Build a repair prompt from this context so a model can propose an updated flow.",
			}
		}
	case "tsplay.repair_flow":
		if ok {
			return map[string]any{
				"tool":   "tsplay.validate_flow",
				"reason": "Validate the updated flow generated from the repair prompt before running it again.",
			}
		}
	}
	return nil
}

func latestTSPlayTraceArtifacts(trace []FlowStepTrace) map[string]any {
	for i := len(trace) - 1; i >= 0; i-- {
		if trace[i].Artifacts == nil {
			continue
		}
		return map[string]any{
			"directory":         trace[i].Artifacts.Directory,
			"screenshot_path":   trace[i].Artifacts.ScreenshotPath,
			"html_path":         trace[i].Artifacts.HTMLPath,
			"dom_snapshot_path": trace[i].Artifacts.DOMSnapshotPath,
			"capture_error":     trace[i].Artifacts.CaptureError,
			"step_index":        trace[i].Index,
			"step_path":         trace[i].Path,
			"action":            trace[i].Action,
			"status":            trace[i].Status,
		}
	}
	return nil
}

func countTSPlayItems(value any) int {
	switch typed := value.(type) {
	case nil:
		return 0
	case []any:
		return len(typed)
	case []string:
		return len(typed)
	case []map[string]any:
		return len(typed)
	default:
		return 0
	}
}

func tsplayNestedString(value any, key string) string {
	typed, ok := value.(map[string]any)
	if !ok {
		return ""
	}
	return stringValue(typed[key])
}

func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}
