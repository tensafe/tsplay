package tsplay_core

import "fmt"

func BuildFlowJSONSchema() map[string]any {
	stepProperties := map[string]any{
		"name":              map[string]any{"type": "string", "description": "Optional human-readable step name."},
		"action":            map[string]any{"type": "string", "enum": FlowActionNames(), "description": "TSPlay action name."},
		"args":              map[string]any{"type": "array", "description": "Positional arguments. Prefer named fields for AI-generated flows."},
		"with":              map[string]any{"type": "object", "description": "Extra named parameters. Prefer top-level named fields when available. For set_var non-string literals, use with.value."},
		"save_as":           map[string]any{"type": "string", "pattern": flowIdentifierPattern.String(), "description": "Save the action output as a flow variable."},
		"continue_on_error": map[string]any{"type": "boolean", "description": "Continue to next step when this step fails."},
		"steps":             map[string]any{"type": "array", "minItems": 1, "description": "Nested steps for control actions such as retry.", "items": map[string]any{"$ref": "#/$defs/step"}},
		"condition":         map[string]any{"$ref": "#/$defs/step", "description": "Condition step for if and wait_until. Its output truthiness controls the branch."},
		"then":              map[string]any{"type": "array", "minItems": 1, "items": map[string]any{"$ref": "#/$defs/step"}},
		"else":              map[string]any{"type": "array", "minItems": 1, "items": map[string]any{"$ref": "#/$defs/step"}},
		"on_error":          map[string]any{"type": "array", "minItems": 1, "items": map[string]any{"$ref": "#/$defs/step"}},
		"url":               map[string]any{"type": "string"},
		"selector":          map[string]any{"type": "string", "description": "Use selector candidates from tsplay.observe_page when possible."},
		"text":              map[string]any{"type": "string"},
		"value":             map[string]any{"type": "string", "description": "String value. For set_var non-string literals, put the literal in with.value."},
		"timeout":           map[string]any{"type": "integer", "minimum": 1},
		"times":             map[string]any{"type": "integer", "minimum": 1, "default": 3},
		"interval_ms":       map[string]any{"type": "integer", "minimum": 0, "default": 0},
		"items":             map[string]any{"description": "List value or variable placeholder for foreach."},
		"item_var":          map[string]any{"type": "string", "pattern": flowIdentifierPattern.String()},
		"index_var":         map[string]any{"type": "string", "pattern": flowIdentifierPattern.String()},
		"seconds":           map[string]any{"type": "number", "exclusiveMinimum": 0},
		"path":              map[string]any{"type": "string"},
		"range":             map[string]any{"type": "string", "description": "Optional Excel cell range such as A2:B20 for read_excel."},
		"script":            map[string]any{"type": "string"},
		"code":              map[string]any{"type": "string"},
		"attribute":         map[string]any{"type": "string"},
		"sheet":             map[string]any{"type": "string", "description": "Optional Excel sheet name for read_excel."},
		"key":               map[string]any{"type": "string"},
		"connection":        map[string]any{"type": "string", "description": "Optional named external connection such as a Redis alias."},
		"file_path":         map[string]any{"type": "string"},
		"files":             map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		"save_path":         map[string]any{"type": "string"},
		"pattern":           map[string]any{"type": "string"},
		"from":              map[string]any{"description": "Input value for json_extract. Can be a variable placeholder or structured data."},
		"delta":             map[string]any{"type": "integer"},
		"ttl_seconds":       map[string]any{"type": "integer", "minimum": 1},
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
			"browser": map[string]any{
				"type":                 "object",
				"description":          "Optional browser launch/session config applied to the whole flow.",
				"additionalProperties": false,
				"properties": map[string]any{
					"headless":           map[string]any{"type": "boolean"},
					"use_session":        map[string]any{"type": "string", "description": "Reuse a named saved session created by tsplay.save_session."},
					"storage_state":      map[string]any{"type": "string", "description": "Load browser storage state from a file relative to the artifact root."},
					"storage_state_path": map[string]any{"type": "string", "description": "Alias of storage_state."},
					"load_storage_state": map[string]any{"type": "string", "description": "Alias of storage_state."},
					"save_storage_state": map[string]any{"type": "string", "description": "Save browser storage state to a file relative to the artifact root after the flow finishes."},
					"persistent":         map[string]any{"type": "boolean", "description": "Use a persistent browser profile stored under the artifact root."},
					"profile":            map[string]any{"type": "string", "description": "Persistent browser profile name."},
					"session":            map[string]any{"type": "string", "description": "Optional session name inside the profile."},
					"timeout":            map[string]any{"type": "integer", "minimum": 0, "description": "Default browser/page timeout in milliseconds."},
					"user_agent":         map[string]any{"type": "string"},
					"viewport": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"required":             []string{"width", "height"},
						"properties": map[string]any{
							"width":  map[string]any{"type": "integer", "minimum": 1},
							"height": map[string]any{"type": "integer", "minimum": 1},
						},
					},
				},
			},
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
			"step":                stepSchema,
			"action_manifest":     buildFlowActionManifest(),
			"generation_rules":    flowSchemaGenerationRules(),
			"selector_strategy":   flowSelectorStrategy(),
			"authoring_checklist": flowAuthoringChecklist(),
			"repair_checklist":    flowRepairValidationChecklist(),
		},
	}
}

func buildFlowActionSchemaConstraints() []any {
	constraints := make([]any, 0, len(flowActionSpecs))
	for _, action := range FlowActionNames() {
		spec := flowActionSpecs[action]
		if special := flowSpecialActionSchemaConstraint(action); special != nil {
			constraints = append(constraints, special)
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

func flowSpecialActionSchemaConstraint(action string) map[string]any {
	if action == "set_var" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "save_as", "value"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "save_as", "with"},
						"not":      map[string]any{"required": []string{"args"}},
						"properties": map[string]any{
							"with": map[string]any{
								"type":     "object",
								"required": []string{"value"},
							},
						},
					},
				},
			},
		}
	}
	if action == "append_var" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "save_as", "value"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "save_as", "with"},
						"not":      map[string]any{"required": []string{"args"}},
						"properties": map[string]any{
							"with": map[string]any{
								"type":     "object",
								"required": []string{"value"},
							},
						},
					},
				},
			},
		}
	}
	if action == "redis_set" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "key", "value"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "key", "with"},
						"not":      map[string]any{"required": []string{"args"}},
						"properties": map[string]any{
							"with": map[string]any{
								"type":     "object",
								"required": []string{"value"},
							},
						},
					},
					buildFlowActionArgsSchema(action, flowActionSpecs[action]),
				},
			},
		}
	}
	if action == "write_json" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "file_path", "value"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "file_path", "with"},
						"not":      map[string]any{"required": []string{"args"}},
						"properties": map[string]any{
							"with": map[string]any{
								"type":     "object",
								"required": []string{"value"},
							},
						},
					},
					map[string]any{
						"required": []string{"action", "args"},
						"properties": map[string]any{
							"args": map[string]any{
								"type":     "array",
								"minItems": 2,
								"maxItems": 2,
								"prefixItems": []any{
									flowParamJSONSchema("file_path"),
									map[string]any{},
								},
							},
						},
					},
				},
			},
		}
	}
	if action == "write_csv" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "file_path", "value"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "file_path", "with"},
						"not":      map[string]any{"required": []string{"args"}},
						"properties": map[string]any{
							"with": map[string]any{
								"type":     "object",
								"required": []string{"value"},
							},
						},
					},
					map[string]any{
						"required": []string{"action", "args"},
						"properties": map[string]any{
							"args": map[string]any{
								"type":     "array",
								"minItems": 2,
								"maxItems": 3,
								"prefixItems": []any{
									flowParamJSONSchema("file_path"),
									map[string]any{},
									map[string]any{"oneOf": []any{map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, flowPlaceholderJSONSchema()}},
								},
							},
						},
					},
				},
			},
		}
	}
	if action == "read_excel" {
		return map[string]any{
			"if": map[string]any{
				"properties": map[string]any{"action": map[string]any{"const": action}},
				"required":   []string{"action"},
			},
			"then": map[string]any{
				"description": fmt.Sprintf("Constraints for action %q.", action),
				"anyOf": []any{
					map[string]any{
						"required": []string{"action", "file_path"},
						"not":      map[string]any{"required": []string{"args"}},
					},
					map[string]any{
						"required": []string{"action", "args"},
						"properties": map[string]any{
							"args": map[string]any{
								"type":     "array",
								"minItems": 1,
								"maxItems": 3,
								"prefixItems": []any{
									flowParamJSONSchema("file_path"),
									flowParamJSONSchema("sheet"),
									flowParamJSONSchema("range"),
								},
							},
						},
					},
				},
			},
		}
	}
	required := []string{}
	switch action {
	case "retry":
		required = []string{"steps"}
	case "if":
		required = []string{"condition"}
	case "foreach":
		required = []string{"items", "item_var", "steps"}
	case "on_error":
		required = []string{"steps", "on_error"}
	case "wait_until":
		required = []string{"condition"}
	default:
		return nil
	}
	return map[string]any{
		"if": map[string]any{
			"properties": map[string]any{"action": map[string]any{"const": action}},
			"required":   []string{"action"},
		},
		"then": map[string]any{
			"description": fmt.Sprintf("Constraints for action %q.", action),
			"required":    append([]string{"action"}, required...),
			"not":         map[string]any{"required": []string{"args"}},
		},
	}
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

func flowPlaceholderJSONSchema() map[string]any {
	return map[string]any{
		"type":        "string",
		"pattern":     placeholderPattern.String(),
		"description": "Full variable placeholder, for example {{timeout_ms}}.",
	}
}

func flowParamJSONSchema(name string) map[string]any {
	placeholder := flowPlaceholderJSONSchema()
	switch flowParamType(name) {
	case "any":
		return map[string]any{}
	case "bool":
		return map[string]any{"oneOf": []any{map[string]any{"type": "boolean"}, placeholder}}
	case "int":
		return map[string]any{"oneOf": []any{map[string]any{"type": "integer"}, placeholder}}
	case "number":
		return map[string]any{"oneOf": []any{map[string]any{"type": "number"}, placeholder}}
	case "object":
		return map[string]any{"oneOf": []any{map[string]any{"type": "object"}, placeholder}}
	case "string_list":
		return map[string]any{"oneOf": []any{map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, placeholder}}
	default:
		return map[string]any{"type": "string"}
	}
}

func BuildFlowExamples() []map[string]any {
	return []map[string]any{
		{
			"name":          "import_rows_from_csv",
			"description":   "Load local CSV rows and iterate through them to fill repeated form fields.",
			"focus_actions": []string{"read_csv", "foreach", "type_text", "click"},
			"when_to_use":   "A trusted local automation needs to read structured rows from a CSV file and submit them one by one.",
			"flow": `schema_version: "1"
name: import_rows_from_csv
steps:
  - action: read_csv
    file_path: imports/users.csv
    save_as: rows
  - action: foreach
    items: "{{rows}}"
    item_var: row
    steps:
      - action: type_text
        selector: "#name"
        text: "{{row.name}}"
      - action: type_text
        selector: "#phone"
        text: "{{row.phone}}"
      - action: click
        selector: "#submit"
`,
		},
		{
			"name":          "import_rows_from_excel_range",
			"description":   "Read a bounded Excel range, assign explicit headers, and iterate through the rows.",
			"focus_actions": []string{"read_excel", "foreach", "type_text", "click"},
			"when_to_use":   "The sheet contains titles, notes, or multiple tables, so only one rectangular range should be imported.",
			"flow": `schema_version: "1"
name: import_rows_from_excel_range
steps:
  - action: read_excel
    file_path: imports/users.xlsx
    sheet: Users
    range: A2:B20
    with:
      headers:
        - name
        - phone
    save_as: rows
  - action: foreach
    items: "{{rows}}"
    item_var: row
    steps:
      - action: type_text
        selector: "#name"
        text: "{{row.name}}"
      - action: type_text
        selector: "#phone"
        text: "{{row.phone}}"
      - action: click
        selector: "#submit"
`,
		},
		{
			"name":          "resume_import_with_writeback",
			"description":   "Resume a batch import from a source row, checkpoint the next row in Redis when available, and write a result ledger to JSON or CSV.",
			"focus_actions": []string{"read_excel", "foreach", "on_error", "append_var", "write_json", "write_csv"},
			"when_to_use":   "The import may be resumed in chunks, each source row needs a durable success or failure record, and a Redis-backed checkpoint is helpful but optional.",
			"flow": `schema_version: "1"
name: resume_import_with_writeback
vars:
  import_results: []
  resume_from_row: 2
steps:
  - action: read_excel
    file_path: imports/users.xlsx
    sheet: Users
    with:
      start_row: "{{resume_from_row}}"
      limit: 100
      row_number_field: source_row
    save_as: rows
  - action: foreach
    items: "{{rows}}"
    item_var: row
    with:
      progress_key: imports:users:resume_row
    steps:
      - action: on_error
        steps:
          - action: type_text
            selector: "#name"
            text: "{{row.name}}"
          - action: click
            selector: "#submit"
          - action: append_var
            save_as: import_results
            with:
              value:
                source_row: "{{row.source_row}}"
                status: success
        on_error:
          - action: append_var
            save_as: import_results
            with:
              value:
                source_row: "{{row.source_row}}"
                status: failed
                error: "{{last_error}}"
  - action: write_json
    file_path: reports/import-results.json
    with:
      value: "{{import_results}}"
  - action: write_csv
    file_path: reports/import-results.csv
    with:
      value: "{{import_results}}"
      headers:
        - source_row
        - status
        - error
`,
		},
		{
			"name":          "write_scraped_rows_to_mysql",
			"description":   "Iterate through structured scrape results and insert each row into a MySQL table directly from the Flow.",
			"focus_actions": []string{"foreach", "db_insert"},
			"when_to_use":   "Scraped data already exists as Flow variables and needs to be reported to a MySQL table with explicit column mapping.",
			"flow": `schema_version: "1"
name: write_scraped_rows_to_mysql
vars:
  query: 山东大学
  results:
    - title: 山东大学
      url: https://www.sdu.edu.cn/
      rank: 1
steps:
  - action: foreach
    items: "{{results}}"
    item_var: item
    steps:
      - action: db_insert
        connection: reporting
        with:
          driver: postgres
          table: crawl_results
          columns:
            - keyword
            - title
            - url
            - rank
          row:
            keyword: "{{query}}"
            title: "{{item.title}}"
            url: "{{item.url}}"
            rank: "{{item.rank}}"
`,
		},
		{
			"name":          "search_and_collect_links",
			"description":   "Open a search page, type a query, click submit, and save extracted links.",
			"focus_actions": []string{"navigate", "wait_for_selector", "type_text", "click", "get_all_links"},
			"when_to_use":   "Simple page navigation plus a single extraction result.",
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
			"name":          "use_observation_selector_candidates",
			"description":   "Use selector candidates from tsplay.observe_page instead of asking the user for HTML details.",
			"focus_actions": []string{"navigate", "type_text", "click", "wait_for_selector"},
			"when_to_use":   "The user describes intent but does not know HTML details.",
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
			"name":          "extract_table",
			"description":   "Capture a table into a variable for later processing.",
			"focus_actions": []string{"navigate", "wait_for_selector", "capture_table"},
			"when_to_use":   "The page already has stable table markup and later steps need structured rows.",
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
			"name":          "failure_artifact_friendly",
			"description":   "Keep steps small and named so trace artifacts make repair easier.",
			"focus_actions": []string{"navigate", "wait_for_selector", "click"},
			"when_to_use":   "The flow is likely to be repaired automatically later.",
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
			"name":          "retry_with_assertions",
			"description":   "Retry flaky interactions and assert the business result before continuing.",
			"focus_actions": []string{"retry", "click", "assert_visible", "assert_text"},
			"when_to_use":   "The page is dynamic and a click may succeed only after a short delay.",
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
			"name":          "control_flow_branch_loop_recover",
			"description":   "Use if, foreach, on_error, and wait_until for common business control flow without lua.",
			"focus_actions": []string{"if", "foreach", "on_error", "wait_until"},
			"when_to_use":   "The flow needs login branching, list traversal, and local failure recovery.",
			"flow": `schema_version: "1"
name: batch_export_orders
vars:
  orders_url: https://example.com/orders
  username: demo
  order_ids:
    - A1001
    - A1002
steps:
  - action: navigate
    url: "{{orders_url}}"
  - action: if
    condition:
      action: is_visible
      selector: ".login-dialog"
    then:
      - action: type_text
        selector: "#username"
        text: "{{username}}"
      - action: click
        selector: "#login"
    else:
      - action: wait_for_selector
        selector: "#orders-table"
        timeout: 10000
  - action: foreach
    items: "{{order_ids}}"
    item_var: order_id
    index_var: order_index
    steps:
      - action: type_text
        selector: "#order-id"
        text: "{{order_id}}"
      - action: on_error
        steps:
          - action: click
            selector: 'text="Export"'
          - action: wait_until
            timeout: 30000
            interval_ms: 500
            condition:
              action: is_visible
              selector: "#export-result"
        on_error:
          - action: reload
          - action: wait_for_selector
            selector: "#orders-table"
            timeout: 10000
`,
		},
		{
			"name":          "extract_text_and_set_var",
			"description":   "Extract text into variables and build a branch-friendly Flow without asking users for HTML internals.",
			"focus_actions": []string{"extract_text", "set_var", "if", "assert_text"},
			"when_to_use":   "The next action depends on page text, counts, or labels that should become Flow variables first.",
			"flow": `schema_version: "1"
name: extract_summary_and_branch
vars:
  orders_url: https://example.com/orders
steps:
  - action: navigate
    url: "{{orders_url}}"
  - action: extract_text
    selector: ".summary .count"
    pattern: '([0-9]+)'
    save_as: order_count
  - action: set_var
    save_as: export_message
    value: "Current orders: {{order_count}}"
  - action: if
    condition:
      action: extract_text
      selector: ".summary .count"
      pattern: '[1-9][0-9]*'
    then:
      - action: click
        selector: 'text="Export orders"'
    else:
      - action: assert_text
        selector: ".empty-state"
        text: "No orders"
`,
		},
		{
			"name":           "browser_session_with_storage_state",
			"description":    "Reuse login state and browser profile settings at the flow level instead of pushing that logic into individual steps.",
			"focus_actions":  []string{"navigate", "assert_visible"},
			"when_to_use":    "The business flow depends on a stable login session, custom timeout, viewport, or user agent.",
			"requires_allow": []string{"allow_browser_state"},
			"flow": `schema_version: "1"
name: admin_orders_with_saved_session
browser:
  headless: true
  storage_state: states/admin.json
  save_storage_state: states/admin-latest.json
  timeout: 30000
  user_agent: tsplay-bot/1.0
  viewport:
    width: 1440
    height: 900
steps:
  - action: navigate
    url: https://example.com/admin/orders
  - action: assert_visible
    selector: "#orders-table"
    timeout: 10000
`,
		},
		{
			"name":          "reuse_named_session",
			"description":   "Reference a named reusable session alias instead of hardcoding storage_state paths in every Flow.",
			"focus_actions": []string{"navigate", "assert_visible"},
			"when_to_use":   "A business team already saved a login session with tsplay.save_session and wants future flows to reuse it by name.",
			"flow": `schema_version: "1"
name: admin_orders_from_named_session
browser:
  use_session: admin
  headless: true
  timeout: 30000
steps:
  - action: navigate
    url: https://example.com/admin/orders
  - action: assert_visible
    selector: "#orders-table"
    timeout: 10000
`,
		},
		{
			"name":           "lua_escape_hatch",
			"description":    "Use lua only for cases not yet expressible by structured actions.",
			"focus_actions":  []string{"lua"},
			"when_to_use":    "Only when structured Flow actions cannot express the needed business rule yet.",
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

func flowSchemaGenerationRules() []string {
	return []string{
		`Always include schema_version: "1".`,
		"Prefer named parameters over args for readability and validation.",
		"Generate structured Flow first. Use lua only as an explicit escape hatch.",
		"Put session-wide browser concerns such as headless, use_session, storage_state, timeout, viewport, user_agent, or persistent profile naming in the top-level browser block.",
		"Turn visible page facts into variables with extract_text + save_as, then use set_var when later steps need a stable derived value.",
		"Use assert_visible/assert_text for business checks, retry for flaky page interactions, if for optional page states, foreach for lists, on_error for local recovery, and wait_until for polling conditions.",
		"Keep steps small and named so trace artifacts and repair context point to an exact failure location.",
		"Do not use execute_script, evaluate, file actions, or browser state actions unless the request explicitly needs them and the MCP allow flags are set.",
	}
}

func flowSelectorStrategy() []string {
	return []string{
		"Prefer selectors from tsplay.observe_page selector_candidates when available.",
		"Selector priority: data-testid, data-cy, id, placeholder, aria-label, role/text, stable class combinations, XPath only as a last resort.",
		"Prefer selectors tied to user intent such as button text, input placeholder, or table/test ids over brittle DOM depth.",
		"When a selector may appear late, pair it with wait_for_selector, retry, or wait_until instead of using sleep first.",
	}
}

func flowAuthoringChecklist() []string {
	return []string{
		"Map the user intent to page states first, then choose actions.",
		"Extract page values into variables before branching or looping on them.",
		"Prefer save_as variables that describe business meaning, not DOM details.",
		"Add assertions around the business result, not only around low-level clicks.",
		"When adding recovery logic, keep it local with on_error instead of rewriting the whole Flow.",
	}
}

func flowRepairValidationChecklist() []string {
	return []string{
		"Validate the repaired Flow before running it.",
		"Check whether selector, text, timeout, pattern, or variable references are the actual failure point.",
		"Prefer the smallest repair that keeps existing save_as outputs and downstream variable names stable.",
		"Use failure artifact paths and DOM snapshot excerpts as evidence; do not inline full HTML into the repaired Flow.",
	}
}

func flowExampleSelectionHints() []string {
	return []string{
		"Start from the example whose focus_actions are closest to the requested business flow.",
		"When the user does not know selectors, combine tsplay.observe_page with the use_observation_selector_candidates example.",
		"When the next step depends on page text or counts, start from extract_text_and_set_var.",
		"When the page is flaky, start from retry_with_assertions before reaching for lua.",
	}
}
