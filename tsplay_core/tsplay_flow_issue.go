package tsplay_core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type FlowIssue struct {
	Code           string   `json:"code,omitempty"`
	Message        string   `json:"message,omitempty"`
	Path           string   `json:"path,omitempty"`
	StepPath       string   `json:"step_path,omitempty"`
	Action         string   `json:"action,omitempty"`
	Field          string   `json:"field,omitempty"`
	DidYouMean     string   `json:"did_you_mean,omitempty"`
	Suggestion     string   `json:"suggestion,omitempty"`
	AllowedFields  []string `json:"allowed_fields,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
	AllowedParams  []string `json:"allowed_params,omitempty"`
	Line           int      `json:"line,omitempty"`
	Column         int      `json:"column,omitempty"`
}

type FlowParseError struct {
	Issue FlowIssue
}

func (e *FlowParseError) Error() string {
	if e == nil {
		return ""
	}
	return e.Issue.Message
}

var (
	flowUnsupportedActionPattern = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*) uses unsupported action "([^"]+)"`)
	flowUnexpectedParamPattern   = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*) action "([^"]+)" does not accept parameter "([^"]+)"`)
	flowMissingParamPattern      = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*) action "([^"]+)" requires "?([a-z_]+)"?`)
	flowUnknownVarPattern        = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*) action "([^"]+)" parameter "([^"]+)" references unknown variable "([^"]+)"`)
	flowSecurityActionPattern    = regexp.MustCompile(`step ([0-9]+(?:\.[A-Za-z_]+|\.[0-9]+)*) action "([^"]+)"(?: progress checkpoint)? is disabled by security policy(?:; set ([a-z_]+)=true only for trusted flows)?`)
)

type flowDocumentContext string

const (
	flowContextRoot     flowDocumentContext = "flow"
	flowContextBrowser  flowDocumentContext = "browser"
	flowContextViewport flowDocumentContext = "viewport"
	flowContextStep     flowDocumentContext = "step"
)

var flowFieldAliasHints = map[flowDocumentContext]map[string]FlowIssue{
	flowContextRoot: {
		"schemaversion": {
			DidYouMean: "schema_version",
			Suggestion: `Replace "schemaVersion" with "schema_version".`,
		},
	},
	flowContextBrowser: {
		"usesession": {
			DidYouMean: "use_session",
			Suggestion: `Replace "useSession" with "use_session".`,
		},
		"storagestate": {
			DidYouMean: "storage_state",
			Suggestion: `Replace "storageState" with "storage_state".`,
		},
		"savestoragestate": {
			DidYouMean: "save_storage_state",
			Suggestion: `Replace "saveStorageState" with "save_storage_state".`,
		},
	},
	flowContextStep: {
		"result_var": {
			DidYouMean: "save_as",
			Suggestion: `Replace "result_var" with "save_as".`,
		},
		"resultvar": {
			DidYouMean: "save_as",
			Suggestion: `Replace "resultVar" with "save_as".`,
		},
		"saveas": {
			DidYouMean: "save_as",
			Suggestion: `Replace "saveAs" with "save_as".`,
		},
		"with.headers": {
			DidYouMean: "with.headers",
			Suggestion: `Nest "headers" under "with", for example: with: { headers: [...] }.`,
		},
		"headers": {
			DidYouMean: "with.headers",
			Suggestion: `Move "headers" under "with.headers".`,
		},
	},
}

var flowActionAliasHints = map[string]FlowIssue{
	"fill": {
		DidYouMean: "type_text",
		Suggestion: `Replace action "fill" with "type_text".`,
	},
	"type": {
		DidYouMean: "type_text",
		Suggestion: `Replace action "type" with "type_text".`,
	},
	"result_var": {
		DidYouMean: "save_as",
		Suggestion: `Replace "result_var" with "save_as".`,
	},
	"save_file": {
		Suggestion: `TSPlay has no generic "save_file" action; use "write_json", "write_csv", "save_html", or a download action depending on the artifact you need.`,
	},
	"log": {
		Suggestion: `TSPlay has no "log" action; surface completion in the caller, or use "set_var"/"append_var" if the flow needs to keep state.`,
	},
}

func ParseFlow(content []byte, format string) (*Flow, error) {
	if issue, err := validateFlowDocumentFields(content, format); err != nil {
		return nil, err
	} else if issue != nil {
		return nil, &FlowParseError{Issue: *issue}
	}

	var flow Flow
	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(content, &flow); err != nil {
			return nil, err
		}
	default:
		if err := yaml.Unmarshal(content, &flow); err != nil {
			return nil, err
		}
	}
	return &flow, nil
}

func ExtractFlowIssue(err error, flow *Flow) *FlowIssue {
	if err == nil {
		return nil
	}

	var parseErr *FlowParseError
	if errors.As(err, &parseErr) && parseErr != nil {
		issue := parseErr.Issue
		return &issue
	}

	return buildFlowIssueFromValidationError(err.Error(), flow)
}

func buildFlowIssueFromValidationError(message string, flow *Flow) *FlowIssue {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}

	if matches := flowUnsupportedActionPattern.FindStringSubmatch(message); len(matches) == 3 {
		stepPath := matches[1]
		action := matches[2]
		issue := &FlowIssue{
			Code:           "unsupported_action",
			Message:        message,
			StepPath:       stepPath,
			Action:         action,
			AllowedActions: FlowActionNames(),
		}
		if alias, ok := lookupActionAliasHint(action); ok {
			issue.DidYouMean = alias.DidYouMean
			issue.Suggestion = alias.Suggestion
			return issue
		}
		if suggestion := closestString(action, issue.AllowedActions); suggestion != "" {
			issue.DidYouMean = suggestion
			issue.Suggestion = fmt.Sprintf("Replace action %q with %q.", action, suggestion)
		}
		return issue
	}

	if matches := flowUnexpectedParamPattern.FindStringSubmatch(message); len(matches) == 4 {
		stepPath := matches[1]
		action := matches[2]
		field := matches[3]
		issue := &FlowIssue{
			Code:          "unexpected_parameter",
			Message:       message,
			StepPath:      stepPath,
			Action:        action,
			Field:         field,
			AllowedParams: flowAllowedParamsForAction(action),
		}
		issue.Suggestion = buildUnexpectedParamSuggestion(action, field, issue.AllowedParams)
		if suggestion := closestString(field, issue.AllowedParams); suggestion != "" {
			issue.DidYouMean = suggestion
			if issue.Suggestion == "" {
				issue.Suggestion = fmt.Sprintf("Replace parameter %q with %q.", field, suggestion)
			}
		}
		return issue
	}

	if matches := flowMissingParamPattern.FindStringSubmatch(message); len(matches) == 4 {
		return &FlowIssue{
			Code:       "missing_required_parameter",
			Message:    message,
			StepPath:   matches[1],
			Action:     matches[2],
			Field:      matches[3],
			Suggestion: fmt.Sprintf("Fill %q for action %q before validating again.", matches[3], matches[2]),
		}
	}

	if matches := flowUnknownVarPattern.FindStringSubmatch(message); len(matches) == 5 {
		return &FlowIssue{
			Code:       "unknown_variable",
			Message:    message,
			StepPath:   matches[1],
			Action:     matches[2],
			Field:      matches[3],
			DidYouMean: matches[4],
			Suggestion: fmt.Sprintf("Define %q in flow.vars or produce it earlier with save_as/set_var.", matches[4]),
		}
	}

	if matches := flowSecurityActionPattern.FindStringSubmatch(message); len(matches) >= 3 {
		issue := &FlowIssue{
			Code:     "security_policy",
			Message:  message,
			StepPath: matches[1],
			Action:   matches[2],
		}
		if len(matches) >= 4 && strings.TrimSpace(matches[3]) != "" {
			issue.Field = matches[3]
			issue.Suggestion = fmt.Sprintf("Retry with %s=true only if this is a trusted flow.", matches[3])
		}
		return issue
	}

	if strings.Contains(message, "requires allow_browser_state=true") {
		return &FlowIssue{
			Code:       "security_policy",
			Message:    message,
			Field:      "allow_browser_state",
			Suggestion: "Retry with allow_browser_state=true only if this is a trusted flow.",
		}
	}

	return nil
}

func validateFlowDocumentFields(content []byte, format string) (*FlowIssue, error) {
	switch strings.ToLower(format) {
	case "json":
		decoder := json.NewDecoder(bytes.NewReader(content))
		decoder.UseNumber()
		var root any
		if err := decoder.Decode(&root); err != nil {
			return nil, err
		}
		return validateGenericFlowObject(root)
	default:
		var doc yaml.Node
		if err := yaml.Unmarshal(content, &doc); err != nil {
			return nil, err
		}
		if len(doc.Content) == 0 {
			return nil, nil
		}
		return validateYAMLFlowObject(doc.Content[0])
	}
}

func validateGenericFlowObject(root any) (*FlowIssue, error) {
	object, ok := root.(map[string]any)
	if !ok {
		return nil, nil
	}
	return validateGenericObject(flowContextRoot, object, "", "")
}

func validateGenericObject(context flowDocumentContext, object map[string]any, path string, stepPath string) (*FlowIssue, error) {
	currentAction, _ := object["action"].(string)
	allowedFields := flowAllowedFieldsForContext(context)
	for _, key := range sortedMapKeys(object) {
		value := object[key]
		if !allowedFields[key] {
			return buildUnknownFieldIssue(context, key, joinDocPath(path, key), stepPath, currentAction, 0, 0), nil
		}
		switch context {
		case flowContextRoot:
			switch key {
			case "browser":
				child, ok := value.(map[string]any)
				if !ok {
					continue
				}
				if issue, err := validateGenericObject(flowContextBrowser, child, joinDocPath(path, key), ""); issue != nil || err != nil {
					return issue, err
				}
			case "steps":
				if issue, err := validateGenericStepList(value, "", joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			}
		case flowContextBrowser:
			if key == "viewport" {
				child, ok := value.(map[string]any)
				if !ok {
					continue
				}
				if issue, err := validateGenericObject(flowContextViewport, child, joinDocPath(path, key), ""); issue != nil || err != nil {
					return issue, err
				}
			}
		case flowContextStep:
			switch key {
			case "steps":
				if issue, err := validateGenericStepList(value, stepPath, joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			case "condition":
				child, ok := value.(map[string]any)
				if !ok {
					continue
				}
				childStepPath := flowStepPath(stepPath+".condition", 1)
				if issue, err := validateGenericObject(flowContextStep, child, joinDocPath(path, key), childStepPath); issue != nil || err != nil {
					return issue, err
				}
			case "then", "else", "on_error":
				if issue, err := validateGenericStepList(value, stepPath+"."+key, joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			}
		}
	}
	return nil, nil
}

func validateGenericStepList(value any, parentStepPath string, path string) (*FlowIssue, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, nil
	}
	for i, item := range items {
		child, ok := item.(map[string]any)
		if !ok {
			continue
		}
		stepPath := flowStepPath(parentStepPath, i+1)
		if issue, err := validateGenericObject(flowContextStep, child, fmt.Sprintf("%s[%d]", path, i), stepPath); issue != nil || err != nil {
			return issue, err
		}
	}
	return nil, nil
}

func validateYAMLFlowObject(root *yaml.Node) (*FlowIssue, error) {
	if root == nil {
		return nil, nil
	}
	return validateYAMLObject(flowContextRoot, root, "", "")
}

func validateYAMLObject(context flowDocumentContext, node *yaml.Node, path string, stepPath string) (*FlowIssue, error) {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil, nil
	}

	currentAction := yamlObjectStringValue(node, "action")
	allowedFields := flowAllowedFieldsForContext(context)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value
		if !allowedFields[key] {
			return buildUnknownFieldIssue(context, key, joinDocPath(path, key), stepPath, currentAction, keyNode.Line, keyNode.Column), nil
		}
		switch context {
		case flowContextRoot:
			switch key {
			case "browser":
				if issue, err := validateYAMLObject(flowContextBrowser, valueNode, joinDocPath(path, key), ""); issue != nil || err != nil {
					return issue, err
				}
			case "steps":
				if issue, err := validateYAMLStepList(valueNode, "", joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			}
		case flowContextBrowser:
			if key == "viewport" {
				if issue, err := validateYAMLObject(flowContextViewport, valueNode, joinDocPath(path, key), ""); issue != nil || err != nil {
					return issue, err
				}
			}
		case flowContextStep:
			switch key {
			case "steps":
				if issue, err := validateYAMLStepList(valueNode, stepPath, joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			case "condition":
				childStepPath := flowStepPath(stepPath+".condition", 1)
				if issue, err := validateYAMLObject(flowContextStep, valueNode, joinDocPath(path, key), childStepPath); issue != nil || err != nil {
					return issue, err
				}
			case "then", "else", "on_error":
				if issue, err := validateYAMLStepList(valueNode, stepPath+"."+key, joinDocPath(path, key)); issue != nil || err != nil {
					return issue, err
				}
			}
		}
	}
	return nil, nil
}

func validateYAMLStepList(node *yaml.Node, parentStepPath string, path string) (*FlowIssue, error) {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil, nil
	}
	for i, item := range node.Content {
		stepPath := flowStepPath(parentStepPath, i+1)
		if issue, err := validateYAMLObject(flowContextStep, item, fmt.Sprintf("%s[%d]", path, i), stepPath); issue != nil || err != nil {
			return issue, err
		}
	}
	return nil, nil
}

func buildUnknownFieldIssue(context flowDocumentContext, field string, path string, stepPath string, action string, line int, column int) *FlowIssue {
	allowedFields := flowAllowedFieldListForContext(context)
	issue := &FlowIssue{
		Code:          "unknown_field",
		Path:          path,
		StepPath:      stepPath,
		Action:        action,
		Field:         field,
		AllowedFields: allowedFields,
		Line:          line,
		Column:        column,
	}

	if alias, ok := lookupFieldAliasHint(context, field); ok {
		issue.DidYouMean = alias.DidYouMean
		issue.Suggestion = alias.Suggestion
	} else if suggestion := closestString(field, allowedFields); suggestion != "" {
		issue.DidYouMean = suggestion
		issue.Suggestion = fmt.Sprintf("Replace field %q with %q.", field, suggestion)
	}

	issue.Message = buildUnknownFieldMessage(context, issue)
	return issue
}

func buildUnknownFieldMessage(context flowDocumentContext, issue *FlowIssue) string {
	if issue == nil {
		return "unknown flow field"
	}

	prefix := "flow"
	switch context {
	case flowContextBrowser:
		prefix = "flow.browser"
	case flowContextViewport:
		prefix = "flow.browser.viewport"
	case flowContextStep:
		if issue.StepPath != "" {
			prefix = fmt.Sprintf("step %s", issue.StepPath)
		} else {
			prefix = "step"
		}
	}

	message := fmt.Sprintf("%s field %q is unknown", prefix, issue.Field)
	if issue.DidYouMean != "" {
		message += fmt.Sprintf("; did you mean %q?", issue.DidYouMean)
	}
	if issue.Suggestion != "" {
		message += " " + issue.Suggestion
	}
	return message
}

func lookupFieldAliasHint(context flowDocumentContext, field string) (FlowIssue, bool) {
	if exact, ok := flowFieldAliasHints[context][strings.ToLower(field)]; ok {
		return exact, true
	}
	normalized := normalizedFlowKey(field)
	if hint, ok := flowFieldAliasHints[context][normalized]; ok {
		return hint, true
	}
	return FlowIssue{}, false
}

func lookupActionAliasHint(action string) (FlowIssue, bool) {
	if hint, ok := flowActionAliasHints[strings.ToLower(strings.TrimSpace(action))]; ok {
		return hint, true
	}
	return FlowIssue{}, false
}

func flowAllowedFieldsForContext(context flowDocumentContext) map[string]bool {
	fields := map[string]bool{}
	for _, name := range flowAllowedFieldListForContext(context) {
		fields[name] = true
	}
	return fields
}

func flowAllowedFieldListForContext(context flowDocumentContext) []string {
	switch context {
	case flowContextRoot:
		return structFieldNames(reflect.TypeOf(Flow{}))
	case flowContextBrowser:
		return structFieldNames(reflect.TypeOf(FlowBrowserConfig{}))
	case flowContextViewport:
		return structFieldNames(reflect.TypeOf(FlowViewport{}))
	case flowContextStep:
		return structFieldNames(reflect.TypeOf(FlowStep{}))
	default:
		return nil
	}
}

func structFieldNames(typ reflect.Type) []string {
	fields := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		name := strings.TrimSpace(strings.Split(tag, ",")[0])
		if name == "" || name == "-" {
			continue
		}
		fields = append(fields, name)
	}
	sort.Strings(fields)
	return fields
}

func flowAllowedParamsForAction(action string) []string {
	spec, ok := flowActionSpecs[action]
	if !ok {
		return nil
	}
	params := make([]string, 0, len(spec.Args)+4)
	for _, arg := range spec.Args {
		params = append(params, arg.Name)
	}
	if spec.VarArgName != "" {
		params = append(params, spec.VarArgName)
	}
	switch action {
	case "set_var", "append_var":
		params = append(params, "save_as", "value", "with.value")
	case "retry":
		params = []string{"times", "interval_ms", "steps"}
	case "if":
		params = []string{"condition", "then", "else"}
	case "foreach":
		params = []string{"items", "item_var", "index_var", "steps"}
	case "on_error":
		params = []string{"steps", "on_error"}
	case "wait_until":
		params = []string{"condition", "timeout", "interval_ms"}
	case "write_csv":
		params = []string{"file_path", "value", "with.headers"}
	case "read_excel":
		params = []string{"file_path", "sheet", "range", "with.headers", "with.start_row", "with.limit", "with.row_number_field"}
	}
	sort.Strings(params)
	return params
}

func buildUnexpectedParamSuggestion(action string, field string, allowed []string) string {
	switch {
	case action == "navigate" && field == "timeout":
		return `Remove step-level "timeout". Use top-level browser.timeout for page defaults, or the MCP tool timeout when observing/drafting.`
	case field == "headers":
		return `Move "headers" under "with.headers".`
	case field == "with.headers":
		return `Nest "headers" under "with", for example: with: { headers: [...] }.`
	}
	if best := closestString(field, allowed); best != "" {
		return fmt.Sprintf("Replace parameter %q with %q.", field, best)
	}
	if len(allowed) > 0 {
		return fmt.Sprintf("Use one of the supported parameters for %q: %s.", action, strings.Join(allowed, ", "))
	}
	return ""
}

func closestString(value string, candidates []string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(candidates) == 0 {
		return ""
	}

	best := ""
	bestScore := 1 << 30
	normalizedValue := normalizedFlowKey(value)
	for _, candidate := range candidates {
		if normalizedFlowKey(candidate) == normalizedValue {
			return candidate
		}
		score := levenshteinDistance(normalizedValue, normalizedFlowKey(candidate))
		if score < bestScore {
			best = candidate
			bestScore = score
		}
	}
	threshold := 2
	if len(normalizedValue) >= 8 {
		threshold = 3
	}
	if bestScore > threshold {
		return ""
	}
	return best
}

func levenshteinDistance(a string, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len(b)
	}
	if b == "" {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = minInt(
				current[j-1]+1,
				prev[j]+1,
				prev[j-1]+cost,
			)
		}
		prev = current
	}
	return prev[len(b)]
}

func minInt(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}

func normalizedFlowKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer("_", "", "-", "", ".", "", " ", "")
	return replacer.Replace(value)
}

func joinDocPath(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func yamlObjectStringValue(node *yaml.Node, key string) string {
	if node == nil || node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return strings.TrimSpace(node.Content[i+1].Value)
		}
	}
	return ""
}

func sortedMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
