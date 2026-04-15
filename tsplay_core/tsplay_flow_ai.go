package tsplay_core

import "fmt"

func BuildFlowJSONSchema() map[string]any {
	stepProperties := map[string]any{
		"name":              map[string]any{"type": "string", "description": "Optional human-readable step name."},
		"action":            map[string]any{"type": "string", "enum": FlowActionNames(), "description": "TSPlay action name."},
		"args":              map[string]any{"type": "array", "description": "Positional arguments. Prefer named fields for AI-generated flows."},
		"with":              map[string]any{"type": "object", "description": "Extra named parameters. Prefer top-level named fields when available."},
		"save_as":           map[string]any{"type": "string", "pattern": flowIdentifierPattern.String(), "description": "Save the action output as a flow variable."},
		"continue_on_error": map[string]any{"type": "boolean", "description": "Continue to next step when this step fails."},
		"steps":             map[string]any{"type": "array", "minItems": 1, "description": "Nested steps for control actions such as retry.", "items": map[string]any{"$ref": "#/$defs/step"}},
		"url":               map[string]any{"type": "string"},
		"selector":          map[string]any{"type": "string", "description": "Use selector candidates from tsplay.observe_page when possible."},
		"text":              map[string]any{"type": "string"},
		"value":             map[string]any{"type": "string"},
		"timeout":           map[string]any{"type": "integer", "minimum": 1},
		"times":             map[string]any{"type": "integer", "minimum": 1, "default": 3},
		"interval_ms":       map[string]any{"type": "integer", "minimum": 0, "default": 0},
		"seconds":           map[string]any{"type": "number", "exclusiveMinimum": 0},
		"path":              map[string]any{"type": "string"},
		"script":            map[string]any{"type": "string"},
		"code":              map[string]any{"type": "string"},
		"attribute":         map[string]any{"type": "string"},
		"file_path":         map[string]any{"type": "string"},
		"files":             map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		"save_path":         map[string]any{"type": "string"},
		"pattern":           map[string]any{"type": "string"},
		"index":             map[string]any{"type": "integer"},
		"context_index":     map[string]any{"type": "integer"},
	}
	stepSchema := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"action"},
		"properties":           stepProperties,
		"allOf":                buildFlowActionSchemaConstraints(),
	}

	return map[string]any{
		"$schema":              "https://json-schema.org/draft/2020-12/schema",
		"$id":                  "https://tsplay.local/schemas/flow.schema.json",
		"title":                "TSPlay Flow",
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"schema_version", "steps"},
		"properties": map[string]any{
			"schema_version": map[string]any{
				"type":        "string",
				"const":       CurrentFlowSchemaVersion,
				"description": fmt.Sprintf("Required Flow schema version. Current version is %q.", CurrentFlowSchemaVersion),
			},
			"name":        map[string]any{"type": "string", "description": "Stable flow name, for example export_yesterday_orders."},
			"description": map[string]any{"type": "string"},
			"vars": map[string]any{
				"type":                 "object",
				"description":          "Initial variables. Values can be referenced as {{var_name}}.",
				"propertyNames":        map[string]any{"pattern": flowIdentifierPattern.String()},
				"additionalProperties": true,
			},
			"steps": map[string]any{
				"type":        "array",
				"minItems":    1,
				"description": "Linear workflow steps. Use lua only as an escape hatch.",
				"items":       map[string]any{"$ref": "#/$defs/step"},
			},
		},
		"$defs": map[string]any{
			"step":            stepSchema,
			"action_manifest": buildFlowActionManifest(),
			"generation_rules": []string{
				"Always include schema_version.",
				"Prefer named parameters over args for readability and validation.",
				"Use selector_candidates from tsplay.observe_page; prefer data-testid/data-cy/id/placeholder/aria-label/text selectors before XPath.",
				"Use assert_visible/assert_text for business checks and retry for flaky page interactions.",
				"Use save_as for extracted values that later steps need.",
				"Do not use lua, execute_script, evaluate, file actions, or browser state actions unless the user explicitly needs them and MCP allow flags are set.",
			},
		},
	}
}

func buildFlowActionSchemaConstraints() []any {
	constraints := make([]any, 0, len(flowActionSpecs))
	for _, action := range FlowActionNames() {
		spec := flowActionSpecs[action]
		if action == "retry" {
			constraints = append(constraints, map[string]any{
				"if": map[string]any{
					"properties": map[string]any{"action": map[string]any{"const": action}},
					"required":   []string{"action"},
				},
				"then": map[string]any{
					"description": "Constraints for action \"retry\".",
					"required":    []string{"action", "steps"},
					"not":         map[string]any{"required": []string{"args"}},
				},
			})
			continue
		}
		requiredNamedParams := []string{}
		for _, arg := range spec.Args {
			if arg.Required {
				requiredNamedParams = append(requiredNamedParams, arg.Name)
			}
		}
		if spec.VarArgName != "" {
			requiredNamedParams = append(requiredNamedParams, spec.VarArgName)
		}

		namedRequired := append([]string{"action"}, requiredNamedParams...)
		constraints = append(constraints, map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"description": "Named-parameter form.",
						"required":    namedRequired,
						"not":         map[string]any{"required": []string{"args"}},
					},
					buildFlowActionArgsSchema(action, spec),
				},
			},
		})
	}
	return constraints
}

func buildFlowActionArgsSchema(action string, spec flowActionSpec) map[string]any {
	prefixItems := []any{}
	for _, arg := range spec.Args {
		prefixItems = append(prefixItems, flowParamJSONSchema(arg.Name))
	}

	minItems := requiredArgCount(spec)
	maxItems := len(spec.Args)
	argsSchema := map[string]any{
		"description": fmt.Sprintf("Positional args form for action %q.", action),
		"required":    []string{"action", "args"},
		"properties": map[string]any{
			"args": map[string]any{
				"type":        "array",
				"minItems":    minItems,
				"prefixItems": prefixItems,
			},
		},
	}

	args := argsSchema["properties"].(map[string]any)["args"].(map[string]any)
	if spec.VarArgName == "" {
		args["maxItems"] = maxItems
		return argsSchema
	}
	args["minItems"] = len(spec.Args) + 1
	args["items"] = flowParamJSONSchema(spec.VarArgName)
	return argsSchema
}

func flowParamJSONSchema(name string) map[string]any {
	placeholder := map[string]any{
		"type":        "string",
		"pattern":     placeholderPattern.String(),
		"description": "Full variable placeholder, for example {{timeout_ms}}.",
	}
	switch flowParamType(name) {
	case "int":
		return map[string]any{"oneOf": []any{map[string]any{"type": "integer"}, placeholder}}
	case "number":
		return map[string]any{"oneOf": []any{map[string]any{"type": "number"}, placeholder}}
	case "string_list":
		return map[string]any{"oneOf": []any{map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, placeholder}}
	default:
		return map[string]any{"type": "string"}
	}
}

func BuildFlowExamples() []map[string]any {
	return []map[string]any{
		{
			"name":        "search_and_collect_links",
			"description": "Open a search page, type a query, click submit, and save extracted links.",
			"flow": `schema_version: "1"
name: search_and_collect_links
vars:
  query: 山东大学
steps:
  - action: navigate
    url: https://www.baidu.com
  - action: wait_for_selector
    selector: "#kw"
    timeout: 5000
  - action: type_text
    selector: "#kw"
    text: "{{query}}"
  - action: click
    selector: "#su"
  - action: get_all_links
    selector: "xpath=//body"
    save_as: links
`,
		},
		{
			"name":        "use_observation_selector_candidates",
			"description": "Use selector candidates from tsplay.observe_page instead of asking the user for HTML details.",
			"flow": `schema_version: "1"
name: order_search_from_observation
vars:
  orders_url: https://example.com/orders
  keyword: "A10086"
steps:
  - action: navigate
    url: "{{orders_url}}"
  - action: type_text
    selector: '[data-testid="order-query"]'
    text: "{{keyword}}"
  - action: click
    selector: 'text="Search"'
  - action: wait_for_selector
    selector: '[data-testid="order-table"]'
    timeout: 10000
`,
		},
		{
			"name":        "extract_table",
			"description": "Capture a table into a variable for later processing.",
			"flow": `schema_version: "1"
name: capture_orders_table
vars:
  orders_url: https://example.com/orders
steps:
  - action: navigate
    url: "{{orders_url}}"
  - action: wait_for_selector
    selector: "#orders-table"
    timeout: 10000
  - action: capture_table
    selector: "#orders-table"
    save_as: orders
`,
		},
		{
			"name":        "failure_artifact_friendly",
			"description": "Keep steps small and named so trace artifacts make repair easier.",
			"flow": `schema_version: "1"
name: export_orders
vars:
  orders_url: https://example.com/orders
steps:
  - name: open orders page
    action: navigate
    url: "{{orders_url}}"
  - name: wait for export button
    action: wait_for_selector
    selector: 'text="Export orders"'
    timeout: 10000
  - name: click export
    action: click
    selector: 'text="Export orders"'
`,
		},
		{
			"name":        "retry_with_assertions",
			"description": "Retry flaky interactions and assert the business result before continuing.",
			"flow": `schema_version: "1"
name: retry_export_orders
vars:
  orders_url: https://example.com/orders
steps:
  - action: navigate
    url: "{{orders_url}}"
  - action: retry
    times: 3
    interval_ms: 1000
    steps:
      - action: click
        selector: 'text="Export orders"'
      - action: assert_visible
        selector: "#export-result"
        timeout: 5000
      - action: assert_text
        selector: "#export-result"
        text: "Export complete"
        timeout: 5000
`,
		},
		{
			"name":           "lua_escape_hatch",
			"description":    "Use lua only for cases not yet expressible by structured actions.",
			"requires_allow": []string{"allow_lua"},
			"flow": `schema_version: "1"
name: custom_lua_escape_hatch
vars:
  target_url: https://example.com
steps:
  - action: navigate
    url: "{{target_url}}"
  - action: lua
    code: |
      print("Use lua sparingly; prefer structured actions first.")
`,
		},
	}
}
