package tsplay_core

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type FlowActionCapabilities struct {
	NeedsRuntime           bool
	NeedsPage              bool
	NeedsContext           bool
	NeedsBrowserState      bool
	ConditionalRuntimeArgs []string
}

func (capabilities FlowActionCapabilities) RequiresPlaywright() bool {
	return capabilities.NeedsRuntime || capabilities.NeedsPage || capabilities.NeedsContext || capabilities.NeedsBrowserState
}

func (capabilities FlowActionCapabilities) manifestValue() map[string]any {
	value := map[string]any{
		"needs_playwright":    capabilities.RequiresPlaywright(),
		"needs_runtime":       capabilities.RequiresPlaywright(),
		"needs_page":          capabilities.NeedsPage,
		"needs_context":       capabilities.NeedsContext,
		"needs_browser_state": capabilities.NeedsBrowserState,
	}
	if len(capabilities.ConditionalRuntimeArgs) > 0 {
		value["conditional_playwright_args"] = append([]string(nil), capabilities.ConditionalRuntimeArgs...)
	}
	return value
}

type PlaywrightUsageReason struct {
	Path        string `json:"path,omitempty"`
	Action      string `json:"action,omitempty"`
	Requirement string `json:"requirement,omitempty"`
	Detail      string `json:"detail,omitempty"`
}

type PlaywrightUsage struct {
	NeedsPlaywright   bool                    `json:"needs_playwright,omitempty"`
	NeedsRuntime      bool                    `json:"needs_runtime,omitempty"`
	NeedsPage         bool                    `json:"needs_page,omitempty"`
	NeedsContext      bool                    `json:"needs_context,omitempty"`
	NeedsBrowserState bool                    `json:"needs_browser_state,omitempty"`
	Reasons           []PlaywrightUsageReason `json:"reasons,omitempty"`
}

func (usage PlaywrightUsage) NeedsBrowser() bool {
	return usage.NeedsPage || usage.NeedsContext || usage.NeedsBrowserState
}

func (usage PlaywrightUsage) Summary(limit int) string {
	if len(usage.Reasons) == 0 {
		switch {
		case usage.NeedsBrowserState:
			return "browser state APIs are required"
		case usage.NeedsContext:
			return "browser context APIs are required"
		case usage.NeedsPage:
			return "page APIs are required"
		case usage.NeedsRuntime:
			return "Playwright runtime is required"
		default:
			return ""
		}
	}

	parts := make([]string, 0, len(usage.Reasons))
	seen := map[string]bool{}
	for _, reason := range usage.Reasons {
		label := strings.TrimSpace(reason.Path)
		if label == "" {
			label = strings.TrimSpace(reason.Detail)
		}
		if label == "" {
			label = strings.TrimSpace(reason.Action)
		}
		if label == "" {
			label = strings.TrimSpace(reason.Requirement)
		}
		if label == "" || seen[label] {
			continue
		}
		seen[label] = true
		parts = append(parts, label)
	}
	if len(parts) == 0 {
		return ""
	}
	if limit <= 0 || limit >= len(parts) {
		return strings.Join(parts, ", ")
	}
	return fmt.Sprintf("%s, +%d more", strings.Join(parts[:limit], ", "), len(parts)-limit)
}

func (usage *PlaywrightUsage) normalize() {
	if usage.NeedsPage || usage.NeedsContext || usage.NeedsBrowserState {
		usage.NeedsRuntime = true
	}
	usage.NeedsPlaywright = usage.NeedsRuntime || usage.NeedsPage || usage.NeedsContext || usage.NeedsBrowserState
}

func (usage *PlaywrightUsage) merge(other PlaywrightUsage) {
	usage.NeedsRuntime = usage.NeedsRuntime || other.NeedsRuntime
	usage.NeedsPage = usage.NeedsPage || other.NeedsPage
	usage.NeedsContext = usage.NeedsContext || other.NeedsContext
	usage.NeedsBrowserState = usage.NeedsBrowserState || other.NeedsBrowserState
	for _, reason := range other.Reasons {
		usage.addReason(reason)
	}
	usage.normalize()
}

func (usage *PlaywrightUsage) addReason(reason PlaywrightUsageReason) {
	switch reason.Requirement {
	case "browser_state":
		usage.NeedsBrowserState = true
		usage.NeedsContext = true
	case "context":
		usage.NeedsContext = true
	case "page":
		usage.NeedsPage = true
	case "", "runtime":
		usage.NeedsRuntime = true
		if reason.Requirement == "" {
			reason.Requirement = "runtime"
		}
	default:
		usage.NeedsRuntime = true
	}
	usage.normalize()

	key := fmt.Sprintf("%s|%s|%s|%s", reason.Path, reason.Action, reason.Requirement, reason.Detail)
	for _, existing := range usage.Reasons {
		existingKey := fmt.Sprintf("%s|%s|%s|%s", existing.Path, existing.Action, existing.Requirement, existing.Detail)
		if existingKey == key {
			return
		}
	}
	usage.Reasons = append(usage.Reasons, reason)
}

func (usage *PlaywrightUsage) addCapabilityReason(capabilities FlowActionCapabilities, path string, action string, detail string) {
	requirement := "runtime"
	switch {
	case capabilities.NeedsBrowserState:
		requirement = "browser_state"
	case capabilities.NeedsContext:
		requirement = "context"
	case capabilities.NeedsPage:
		requirement = "page"
	case capabilities.RequiresPlaywright():
		requirement = "runtime"
	}
	usage.addReason(PlaywrightUsageReason{
		Path:        path,
		Action:      action,
		Requirement: requirement,
		Detail:      detail,
	})
}

var (
	flowPlaywrightLuaGlobalPattern = regexp.MustCompile(`\b(?:page|browser|context)\b`)
	flowHTTPRequestBrowserArgs     = []string{"use_browser_cookies", "use_browser_referer", "use_browser_user_agent"}
	flowActionCapabilitiesRegistry = buildFlowActionCapabilitiesRegistry()
)

func flowActionCapabilitiesFor(action string) (FlowActionCapabilities, bool) {
	capabilities, ok := flowActionCapabilitiesRegistry[action]
	return capabilities, ok
}

func luaActionCapabilitiesFor(action string) (FlowActionCapabilities, bool) {
	if capabilities, ok := flowActionCapabilitiesFor(action); ok {
		return capabilities, true
	}
	switch action {
	case "intercept_request":
		return FlowActionCapabilities{
			NeedsRuntime: true,
			NeedsPage:    true,
		}, true
	default:
		return FlowActionCapabilities{}, false
	}
}

func buildFlowActionCapabilitiesRegistry() map[string]FlowActionCapabilities {
	registry := map[string]FlowActionCapabilities{}
	register := func(capabilities FlowActionCapabilities, names ...string) {
		for _, name := range names {
			registry[name] = capabilities
		}
	}

	pageCapabilities := FlowActionCapabilities{
		NeedsRuntime: true,
		NeedsPage:    true,
	}
	register(pageCapabilities,
		"navigate",
		"click",
		"reload",
		"go_back",
		"go_forward",
		"type_text",
		"get_text",
		"extract_text",
		"set_value",
		"select_option",
		"hover",
		"scroll_to",
		"wait_for_network_idle",
		"wait_for_selector",
		"wait_for_text",
		"assert_visible",
		"assert_text",
		"screenshot",
		"screenshot_element",
		"save_html",
		"accept_alert",
		"dismiss_alert",
		"set_alert_text",
		"execute_script",
		"evaluate",
		"upload_file",
		"upload_multiple_files",
		"download_file",
		"download_url",
		"get_attribute",
		"get_html",
		"get_all_links",
		"capture_table",
		"find_element",
		"find_elements",
		"is_visible",
		"is_enabled",
		"is_checked",
		"is_selected",
		"is_aria_selected",
		"new_tab",
		"close_tab",
		"switch_to_tab",
		"block_request",
		"get_response",
	)

	register(FlowActionCapabilities{
		NeedsRuntime:      true,
		NeedsContext:      true,
		NeedsBrowserState: true,
	}, "get_storage_state", "get_cookies_string")

	register(FlowActionCapabilities{}, "sleep",
		"set_var",
		"append_var",
		"retry",
		"if",
		"foreach",
		"on_error",
		"wait_until",
		"db_transaction",
		"read_csv",
		"read_excel",
		"write_json",
		"write_csv",
		"json_extract",
		"redis_get",
		"redis_set",
		"redis_del",
		"redis_incr",
		"db_insert",
		"db_insert_many",
		"db_upsert",
		"db_query",
		"db_query_one",
		"db_execute",
	)

	register(FlowActionCapabilities{
		ConditionalRuntimeArgs: append([]string(nil), flowHTTPRequestBrowserArgs...),
	}, "http_request")

	register(FlowActionCapabilities{
		ConditionalRuntimeArgs: []string{"code"},
	}, "lua")

	return registry
}

func AnalyzeFlowPlaywrightUsage(flow *Flow) PlaywrightUsage {
	usage := PlaywrightUsage{}
	if flow == nil {
		return usage
	}
	staticCtx := newStaticFlowAnalysisContext(flow)
	usage.merge(analyzeFlowBrowserConfigPlaywrightUsage(flow.Browser))
	usage.merge(analyzeFlowStepListPlaywrightUsage(flow.Steps, "steps", staticCtx))
	usage.normalize()
	return usage
}

func flowUsesPlaywright(flow *Flow) bool {
	return AnalyzeFlowPlaywrightUsage(flow).NeedsPlaywright
}

func analyzeFlowBrowserConfigPlaywrightUsage(config *FlowBrowserConfig) PlaywrightUsage {
	usage := PlaywrightUsage{}
	if config == nil {
		return usage
	}

	add := func(path string, detail string) {
		usage.addReason(PlaywrightUsageReason{
			Path:        path,
			Action:      "browser",
			Requirement: "browser_state",
			Detail:      detail,
		})
	}

	if strings.TrimSpace(config.UseSession) != "" {
		add("browser.use_session", "browser.use_session reuses a saved Playwright session")
	}
	loadPath, _ := config.loadStorageStatePath()
	if loadPath != "" {
		add("browser.storage_state", "browser.storage_state loads saved browser state before the flow runs")
	}
	if strings.TrimSpace(config.SaveStorageState) != "" {
		add("browser.save_storage_state", "browser.save_storage_state saves browser state after the flow finishes")
	}
	if config.Persistent {
		add("browser.persistent", "browser.persistent requires a persistent browser context")
	}
	if strings.TrimSpace(config.Profile) != "" {
		add("browser.profile", "browser.profile requires a persistent browser profile")
	}
	if strings.TrimSpace(config.Session) != "" {
		add("browser.session", "browser.session scopes a persistent browser session")
	}
	usage.normalize()
	return usage
}

func analyzeFlowStepListPlaywrightUsage(steps []FlowStep, listPath string, ctx *FlowContext) PlaywrightUsage {
	usage := PlaywrightUsage{}
	for index, step := range steps {
		stepPath := fmt.Sprintf("%s[%d]", listPath, index+1)
		usage.merge(analyzeFlowStepPlaywrightUsage(step, stepPath, ctx))
	}
	usage.normalize()
	return usage
}

func analyzeFlowStepPlaywrightUsage(step FlowStep, stepPath string, ctx *FlowContext) PlaywrightUsage {
	switch step.Action {
	case "retry", "foreach", "db_transaction":
		return analyzeFlowStepListPlaywrightUsage(step.Steps, stepPath+".steps", ctx)
	case "if":
		usage := PlaywrightUsage{}
		if step.Condition != nil {
			usage.merge(analyzeFlowStepPlaywrightUsage(*step.Condition, stepPath+".condition", ctx))
		}
		usage.merge(analyzeFlowStepListPlaywrightUsage(step.Then, stepPath+".then", ctx))
		usage.merge(analyzeFlowStepListPlaywrightUsage(step.Else, stepPath+".else", ctx))
		return usage
	case "on_error":
		usage := PlaywrightUsage{}
		usage.merge(analyzeFlowStepListPlaywrightUsage(step.Steps, stepPath+".steps", ctx))
		usage.merge(analyzeFlowStepListPlaywrightUsage(step.OnError, stepPath+".on_error", ctx))
		return usage
	case "wait_until":
		if step.Condition == nil {
			return PlaywrightUsage{}
		}
		return analyzeFlowStepPlaywrightUsage(*step.Condition, stepPath+".condition", ctx)
	}

	switch step.Action {
	case "http_request":
		return analyzeFlowHTTPRequestPlaywrightUsage(step, stepPath, ctx)
	case "lua":
		return analyzeFlowLuaStepPlaywrightUsage(step, stepPath, ctx)
	}

	capabilities, ok := flowActionCapabilitiesFor(step.Action)
	if !ok || !capabilities.RequiresPlaywright() {
		return PlaywrightUsage{}
	}
	usage := PlaywrightUsage{}
	usage.addCapabilityReason(capabilities, stepPath+"."+step.Action, step.Action, describeFixedPlaywrightRequirement(step.Action, capabilities))
	usage.normalize()
	return usage
}

func analyzeFlowHTTPRequestPlaywrightUsage(step FlowStep, stepPath string, ctx *FlowContext) PlaywrightUsage {
	usage := PlaywrightUsage{}

	addReason := func(param string, requirement PlaywrightUsageReason) {
		requirement.Path = fmt.Sprintf("%s.http_request.%s", stepPath, param)
		requirement.Action = step.Action
		usage.addReason(requirement)
	}

	if needed, dynamic := flowStepBoolParamMayRequirePlaywright(step, "use_browser_cookies", ctx); needed {
		detail := "http_request reuses cookies from the active browser context"
		if dynamic {
			detail = "http_request may reuse cookies from the active browser context after variable resolution"
		}
		addReason("use_browser_cookies", PlaywrightUsageReason{
			Requirement: "browser_state",
			Detail:      detail,
		})
	}
	if needed, dynamic := flowStepBoolParamMayRequirePlaywright(step, "use_browser_referer", ctx); needed {
		detail := "http_request copies the current page URL into the Referer header"
		if dynamic {
			detail = "http_request may copy the current page URL into the Referer header after variable resolution"
		}
		addReason("use_browser_referer", PlaywrightUsageReason{
			Requirement: "page",
			Detail:      detail,
		})
	}
	if needed, dynamic := flowStepBoolParamMayRequirePlaywright(step, "use_browser_user_agent", ctx); needed {
		detail := "http_request reads the current page user agent"
		if dynamic {
			detail = "http_request may read the current page user agent after variable resolution"
		}
		addReason("use_browser_user_agent", PlaywrightUsageReason{
			Requirement: "page",
			Detail:      detail,
		})
	}
	usage.normalize()
	return usage
}

func analyzeFlowLuaStepPlaywrightUsage(step FlowStep, stepPath string, ctx *FlowContext) PlaywrightUsage {
	rawCode := any(step.Code)
	if rawCode == "" {
		if value, ok := step.param("code"); ok {
			rawCode = value
		}
	}

	code := strings.TrimSpace(step.Code)
	if resolvedCtx := ctx; resolvedCtx != nil {
		if resolvedCode, err := step.luaCode(resolvedCtx); err == nil {
			code = strings.TrimSpace(resolvedCode)
		}
	}
	if code == "" {
		code = strings.TrimSpace(fmt.Sprint(rawCode))
	}
	if code == "" {
		return PlaywrightUsage{}
	}

	if len(flowReferences(rawCode)) > 0 && len(flowReferences(code)) > 0 {
		usage := PlaywrightUsage{}
		usage.addReason(PlaywrightUsageReason{
			Path:        stepPath + ".lua.code",
			Action:      step.Action,
			Requirement: "runtime",
			Detail:      "dynamic Lua code may call Playwright APIs after variable resolution",
		})
		return usage
	}

	return analyzeLuaScriptPlaywrightUsageAtPath(code, stepPath+".lua")
}

func newStaticFlowAnalysisContext(flow *Flow) *FlowContext {
	vars := map[string]any{}
	if flow != nil {
		for key, value := range flow.Vars {
			vars[key] = value
		}
	}
	ctx := &FlowContext{Vars: vars}
	if len(vars) == 0 {
		return ctx
	}
	for iteration := 0; iteration < len(vars)+2; iteration++ {
		changed := false
		for key, value := range vars {
			resolved, err := resolveValue(value, ctx)
			if err != nil {
				continue
			}
			if reflect.DeepEqual(resolved, value) {
				continue
			}
			vars[key] = resolved
			changed = true
		}
		if !changed {
			break
		}
	}
	return ctx
}

func flowStepBoolParamMayRequirePlaywright(step FlowStep, name string, ctx *FlowContext) (bool, bool) {
	value, ok := step.param(name)
	if !ok {
		return false, false
	}

	resolved := value
	if ctx != nil {
		if candidate, err := resolveValue(value, ctx); err == nil {
			resolved = candidate
		}
	}

	if len(flowReferences(value)) > 0 && len(flowReferences(resolved)) > 0 {
		return true, true
	}
	return playwrightTruthyValue(resolved), false
}

func playwrightTruthyValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		trimmed := strings.TrimSpace(strings.ToLower(typed))
		return trimmed != "" && trimmed != "false" && trimmed != "0" && trimmed != "no"
	case int:
		return typed != 0
	case int8:
		return typed != 0
	case int16:
		return typed != 0
	case int32:
		return typed != 0
	case int64:
		return typed != 0
	case uint:
		return typed != 0
	case uint8:
		return typed != 0
	case uint16:
		return typed != 0
	case uint32:
		return typed != 0
	case uint64:
		return typed != 0
	case float32:
		return typed != 0
	case float64:
		return typed != 0
	default:
		return true
	}
}

func describeFixedPlaywrightRequirement(action string, capabilities FlowActionCapabilities) string {
	switch {
	case capabilities.NeedsBrowserState:
		return fmt.Sprintf("action %q uses browser state from the active context", action)
	case capabilities.NeedsContext:
		return fmt.Sprintf("action %q uses the active browser context", action)
	case capabilities.NeedsPage:
		return fmt.Sprintf("action %q uses the active page", action)
	default:
		return fmt.Sprintf("action %q requires the Playwright runtime", action)
	}
}

func AnalyzeLuaScriptPlaywrightUsage(script string) PlaywrightUsage {
	return analyzeLuaScriptPlaywrightUsageAtPath(script, "lua")
}

func analyzeLuaScriptPlaywrightUsageAtPath(script string, basePath string) PlaywrightUsage {
	usage := PlaywrightUsage{}
	sanitized := stripLuaCommentsAndStrings(script)
	if strings.TrimSpace(sanitized) == "" {
		return usage
	}

	calls := luaCalledIdentifiers(sanitized)
	sort.Strings(calls)
	for _, call := range calls {
		if call == "http_request" {
			for _, arg := range flowHTTPRequestBrowserArgs {
				if !luaScriptContainsConfigKey(sanitized, arg) {
					continue
				}
				reason := PlaywrightUsageReason{
					Path:   fmt.Sprintf("%s.http_request.%s", basePath, arg),
					Action: call,
				}
				switch arg {
				case "use_browser_cookies":
					reason.Requirement = "browser_state"
					reason.Detail = "Lua http_request reuses cookies from the active browser context"
				case "use_browser_referer":
					reason.Requirement = "page"
					reason.Detail = "Lua http_request copies the current page URL into the Referer header"
				case "use_browser_user_agent":
					reason.Requirement = "page"
					reason.Detail = "Lua http_request reads the current page user agent"
				}
				usage.addReason(reason)
			}
			continue
		}

		capabilities, ok := luaActionCapabilitiesFor(call)
		if !ok || !capabilities.RequiresPlaywright() {
			continue
		}
		usage.addCapabilityReason(
			capabilities,
			fmt.Sprintf("%s.%s()", basePath, call),
			call,
			describeFixedPlaywrightRequirement(call, capabilities),
		)
	}

	if luaUsesObject(sanitized, "page") {
		usage.addReason(PlaywrightUsageReason{
			Path:        basePath + ".page",
			Action:      "lua",
			Requirement: "page",
			Detail:      "Lua code accesses the active page object",
		})
	}
	if luaUsesObject(sanitized, "browser") {
		usage.addReason(PlaywrightUsageReason{
			Path:        basePath + ".browser",
			Action:      "lua",
			Requirement: "page",
			Detail:      "Lua code accesses the browser object created by Playwright",
		})
	}
	if luaUsesObject(sanitized, "context") || flowPlaywrightLuaGlobalPattern.MatchString(sanitized) && strings.Contains(sanitized, "context") {
		usage.addReason(PlaywrightUsageReason{
			Path:        basePath + ".context",
			Action:      "lua",
			Requirement: "context",
			Detail:      "Lua code accesses the active browser context",
		})
	}
	usage.normalize()
	return usage
}

func stripLuaCommentsAndStrings(source string) string {
	var out strings.Builder
	out.Grow(len(source))

	writeBlank := func(segment string) {
		for _, r := range segment {
			switch r {
			case '\n', '\r', '\t':
				out.WriteRune(r)
			default:
				out.WriteByte(' ')
			}
		}
	}

	for index := 0; index < len(source); {
		if strings.HasPrefix(source[index:], "--") {
			if end, ok := luaLongBracketRange(source, index+2); ok {
				writeBlank(source[index:end])
				index = end
				continue
			}
			lineEnd := strings.IndexByte(source[index:], '\n')
			if lineEnd < 0 {
				writeBlank(source[index:])
				break
			}
			writeBlank(source[index : index+lineEnd])
			index += lineEnd
			continue
		}
		if source[index] == '"' || source[index] == '\'' {
			end := luaQuotedStringEnd(source, index)
			writeBlank(source[index:end])
			index = end
			continue
		}
		if end, ok := luaLongBracketRange(source, index); ok {
			writeBlank(source[index:end])
			index = end
			continue
		}
		out.WriteByte(source[index])
		index++
	}
	return out.String()
}

func luaLongBracketRange(source string, start int) (int, bool) {
	if start >= len(source) || source[start] != '[' {
		return 0, false
	}
	index := start + 1
	for index < len(source) && source[index] == '=' {
		index++
	}
	if index >= len(source) || source[index] != '[' {
		return 0, false
	}
	closing := "]" + strings.Repeat("=", index-start-1) + "]"
	closeIndex := strings.Index(source[index+1:], closing)
	if closeIndex < 0 {
		return len(source), true
	}
	return index + 1 + closeIndex + len(closing), true
}

func luaQuotedStringEnd(source string, start int) int {
	quote := source[start]
	index := start + 1
	for index < len(source) {
		switch source[index] {
		case '\\':
			index += 2
			continue
		case quote:
			return index + 1
		default:
			index++
		}
	}
	return len(source)
}

func luaCalledIdentifiers(source string) []string {
	seen := map[string]bool{}
	calls := []string{}
	for index := 0; index < len(source); {
		if !isLuaIdentifierStart(source[index]) {
			index++
			continue
		}
		start := index
		index++
		for index < len(source) && isLuaIdentifierPart(source[index]) {
			index++
		}
		identifier := source[start:index]
		next := skipLuaWhitespace(source, index)
		if next >= len(source) || source[next] != '(' {
			continue
		}
		prev := previousNonSpaceIndex(source, start)
		if prev >= 0 && (source[prev] == '.' || source[prev] == ':') {
			continue
		}
		if previousLuaIdentifier(source, start) == "function" {
			continue
		}
		if seen[identifier] {
			continue
		}
		seen[identifier] = true
		calls = append(calls, identifier)
	}
	return calls
}

func luaUsesObject(source string, name string) bool {
	for index := 0; index < len(source); {
		if !isLuaIdentifierStart(source[index]) {
			index++
			continue
		}
		start := index
		index++
		for index < len(source) && isLuaIdentifierPart(source[index]) {
			index++
		}
		if source[start:index] != name {
			continue
		}
		next := skipLuaWhitespace(source, index)
		if next < len(source) {
			switch source[next] {
			case '.', ':', '[':
				return true
			}
		}
	}
	return false
}

func luaScriptContainsConfigKey(source string, key string) bool {
	pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(key) + `\s*=`)
	return pattern.MatchString(source)
}

func isLuaIdentifierStart(value byte) bool {
	return (value >= 'A' && value <= 'Z') || (value >= 'a' && value <= 'z') || value == '_'
}

func isLuaIdentifierPart(value byte) bool {
	return isLuaIdentifierStart(value) || (value >= '0' && value <= '9')
}

func skipLuaWhitespace(source string, index int) int {
	for index < len(source) {
		switch source[index] {
		case ' ', '\t', '\n', '\r':
			index++
		default:
			return index
		}
	}
	return index
}

func previousNonSpaceIndex(source string, index int) int {
	for index--; index >= 0; index-- {
		switch source[index] {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return index
		}
	}
	return -1
}

func previousLuaIdentifier(source string, index int) string {
	end := previousNonSpaceIndex(source, index)
	if end < 0 || !isLuaIdentifierPart(source[end]) {
		return ""
	}
	start := end
	for start >= 0 && isLuaIdentifierPart(source[start]) {
		start--
	}
	return source[start+1 : end+1]
}
