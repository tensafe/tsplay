package tsplay_core

import (
	"fmt"
	"regexp"
	"strings"
)

type FlowActionCapabilities struct {
	NeedsPlaywright            bool
	NeedsPage                  bool
	NeedsBrowserState          bool
	ConditionalPlaywrightArgs  []string
	dynamicPlaywrightEvaluator func(FlowStep) bool
}

func (capabilities FlowActionCapabilities) RequiresPlaywright() bool {
	return capabilities.NeedsPlaywright || capabilities.NeedsPage || capabilities.NeedsBrowserState
}

func (capabilities FlowActionCapabilities) manifestValue() map[string]any {
	value := map[string]any{
		"needs_playwright":    capabilities.RequiresPlaywright(),
		"needs_page":          capabilities.NeedsPage,
		"needs_browser_state": capabilities.NeedsBrowserState,
	}
	if len(capabilities.ConditionalPlaywrightArgs) > 0 {
		value["conditional_playwright_args"] = append([]string(nil), capabilities.ConditionalPlaywrightArgs...)
	}
	return value
}

var (
	flowPlaywrightLuaCallPattern   = regexp.MustCompile(`\b(?:navigate|click|reload|go_back|go_forward|type_text|get_text|set_value|select_option|hover|scroll_to|wait_for_network_idle|wait_for_selector|wait_for_text|screenshot|screenshot_element|save_html|accept_alert|dismiss_alert|set_alert_text|execute_script|evaluate|upload_file|upload_multiple_files|download_file|download_url|get_attribute|get_html|get_all_links|capture_table|find_element|find_elements|is_visible|is_enabled|is_checked|is_selected|is_aria_selected|new_tab|close_tab|switch_to_tab|intercept_request|block_request|get_response|get_storage_state|get_cookies_string)\s*\(`)
	flowPlaywrightLuaGlobalPattern = regexp.MustCompile(`\b(?:page|browser|context)\b`)
	flowHTTPRequestBrowserArgs     = []string{"use_browser_cookies", "use_browser_referer", "use_browser_user_agent"}
	flowActionCapabilitiesRegistry = buildFlowActionCapabilitiesRegistry()
)

func flowActionCapabilitiesFor(action string) (FlowActionCapabilities, bool) {
	capabilities, ok := flowActionCapabilitiesRegistry[action]
	return capabilities, ok
}

func buildFlowActionCapabilitiesRegistry() map[string]FlowActionCapabilities {
	registry := map[string]FlowActionCapabilities{}
	register := func(capabilities FlowActionCapabilities, names ...string) {
		for _, name := range names {
			registry[name] = capabilities
		}
	}

	pageCapabilities := FlowActionCapabilities{
		NeedsPlaywright: true,
		NeedsPage:       true,
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
		NeedsPlaywright:   true,
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
		ConditionalPlaywrightArgs:  append([]string(nil), flowHTTPRequestBrowserArgs...),
		dynamicPlaywrightEvaluator: flowHTTPRequestUsesPlaywright,
	}, "http_request")

	register(FlowActionCapabilities{
		ConditionalPlaywrightArgs:  []string{"code"},
		dynamicPlaywrightEvaluator: flowLuaStepUsesPlaywright,
	}, "lua")

	return registry
}

func flowHTTPRequestUsesPlaywright(step FlowStep) bool {
	for _, name := range flowHTTPRequestBrowserArgs {
		if flowStepBoolParamMayRequirePlaywright(step, name) {
			return true
		}
	}
	return false
}

func flowStepBoolParamMayRequirePlaywright(step FlowStep, name string) bool {
	value, ok := step.param(name)
	if !ok {
		return false
	}
	if len(flowReferences(value)) > 0 {
		return true
	}
	switch typed := value.(type) {
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
		return value != nil
	}
}

func flowLuaStepUsesPlaywright(step FlowStep) bool {
	rawCode := any(step.Code)
	code := strings.TrimSpace(step.Code)
	if code == "" {
		if value, ok := step.param("code"); ok {
			rawCode = value
			code = strings.TrimSpace(fmt.Sprint(value))
		}
	}
	if code == "" {
		return false
	}
	if len(flowReferences(rawCode)) > 0 {
		return true
	}
	return flowPlaywrightLuaCallPattern.MatchString(code) || flowPlaywrightLuaGlobalPattern.MatchString(code)
}
