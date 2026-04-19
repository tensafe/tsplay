package tsplay_core

import (
	"math"
	"sort"
	"strings"
)

type FlowExampleInputVar struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	Example     string `json:"example,omitempty"`
}

type FlowExamplePitfall struct {
	Wrong     string `json:"wrong"`
	IssueCode string `json:"issue_code,omitempty"`
	Fix       string `json:"fix,omitempty"`
}

type FlowRecommendedExample struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Intent         string                 `json:"intent,omitempty"`
	WhyThisMatches string                 `json:"why_this_matches,omitempty"`
	Description    string                 `json:"description,omitempty"`
	WhenToUse      string                 `json:"when_to_use,omitempty"`
	Relevance      float64                `json:"relevance,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	FocusActions   []string               `json:"focus_actions,omitempty"`
	RequiresAllow  []string               `json:"requires_allow,omitempty"`
	InputVars      []FlowExampleInputVar  `json:"input_vars,omitempty"`
	CommonPitfalls []FlowExamplePitfall   `json:"common_pitfalls,omitempty"`
	FlowYAML       string                 `json:"flow_yaml,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}

type FlowRepairExampleIssue struct {
	Code   string `json:"code,omitempty"`
	Action string `json:"action,omitempty"`
	Field  string `json:"field,omitempty"`
}

type FlowRepairExampleEdit struct {
	Op   string `json:"op"`
	Path string `json:"path,omitempty"`
	From any    `json:"from,omitempty"`
	To   any    `json:"to,omitempty"`
}

type FlowRepairExample struct {
	ID            string                  `json:"id"`
	Title         string                  `json:"title"`
	MatchingIssue FlowRepairExampleIssue  `json:"matching_issue"`
	Why           string                  `json:"why,omitempty"`
	SafeAutofix   bool                    `json:"safe_autofix,omitempty"`
	BeforeYAML    string                  `json:"before_yaml,omitempty"`
	AfterYAML     string                  `json:"after_yaml,omitempty"`
	Edits         []FlowRepairExampleEdit `json:"edits,omitempty"`
	Revalidate    bool                    `json:"revalidate,omitempty"`
	Notes         []string                `json:"notes,omitempty"`
}

type flowExampleDefinition struct {
	Example      FlowRecommendedExample
	MatchTerms   []string
	IssueCodes   []string
	IssueActions []string
	IssueFields  []string
}

func AttachRecommendedExamples(payload map[string]any, intent string, draft *FlowDraft, issue *FlowIssue, limit int) {
	if payload == nil {
		return
	}
	examples := BuildRecommendedFlowExamples(intent, draft, issue, limit)
	if len(examples) == 0 {
		return
	}
	payload["recommended_examples"] = examples
}

func AttachRepairExample(payload map[string]any, issue *FlowIssue) {
	if payload == nil || issue == nil {
		return
	}
	if repair := BuildFlowRepairExample(issue); repair != nil {
		payload["repair_example"] = repair
	}
}

func BuildRecommendedFlowExamples(intent string, draft *FlowDraft, issue *FlowIssue, limit int) []FlowRecommendedExample {
	if limit <= 0 {
		limit = 3
	}

	intentText := strings.ToLower(strings.TrimSpace(intent))
	if intentText == "" && draft != nil {
		intentText = strings.ToLower(strings.TrimSpace(draft.Intent))
	}
	actionSet := collectFlowActionSet(draft)
	definitions := builtinFlowExampleDefinitions()
	type scoredExample struct {
		index   int
		score   float64
		why     string
		example FlowRecommendedExample
	}
	scored := make([]scoredExample, 0, len(definitions))
	for index, definition := range definitions {
		score, why := scoreFlowExampleDefinition(definition, intentText, actionSet, issue)
		example := definition.Example
		if why != "" {
			example.WhyThisMatches = why
		}
		example.Relevance = scoreToRelevance(score)
		scored = append(scored, scoredExample{
			index:   index,
			score:   score,
			why:     why,
			example: example,
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].index < scored[j].index
		}
		return scored[i].score > scored[j].score
	})

	if limit > len(scored) {
		limit = len(scored)
	}
	result := make([]FlowRecommendedExample, 0, limit)
	for _, item := range scored[:limit] {
		result = append(result, item.example)
	}
	return result
}

func BuildFlowRepairExample(issue *FlowIssue) *FlowRepairExample {
	if issue == nil {
		return nil
	}

	switch {
	case issue.Code == "unsupported_action" && strings.EqualFold(issue.Action, "fill"):
		return &FlowRepairExample{
			ID:    "repair_fill_to_type_text",
			Title: "Replace fill with type_text",
			MatchingIssue: FlowRepairExampleIssue{
				Code:   issue.Code,
				Action: issue.Action,
			},
			Why:         "TSPlay uses type_text for text entry. The unsupported fill action usually comes from Playwright-style habits.",
			SafeAutofix: true,
			BeforeYAML: `- action: fill
  selector: "#kw"
  value: "{{query}}"`,
			AfterYAML: `- action: type_text
  selector: "#kw"
  text: "{{query}}"`,
			Edits: []FlowRepairExampleEdit{
				{Op: "replace", Path: "steps[n].action", From: "fill", To: "type_text"},
				{Op: "rename_field", Path: "steps[n].value", To: "text"},
			},
			Revalidate: true,
			Notes: []string{
				"Keep the original selector when it is already valid.",
				"Rename value/text together so the step stays valid.",
			},
		}
	case issue.Code == "unknown_field" && strings.EqualFold(issue.Field, "result_var"):
		return &FlowRepairExample{
			ID:    "repair_result_var_to_save_as",
			Title: "Replace result_var with save_as",
			MatchingIssue: FlowRepairExampleIssue{
				Code:  issue.Code,
				Field: issue.Field,
			},
			Why:         "TSPlay stores step outputs with save_as. result_var is a common cross-framework alias, not a supported Flow field.",
			SafeAutofix: true,
			BeforeYAML: `- action: evaluate
  selector: "body"
  script: |
    return []
  result_var: rows`,
			AfterYAML: `- action: evaluate
  selector: "body"
  script: |
    return []
  save_as: rows`,
			Edits: []FlowRepairExampleEdit{
				{Op: "rename_field", Path: "steps[n].result_var", To: "save_as"},
			},
			Revalidate: true,
			Notes: []string{
				"Use save_as on actions that produce values for later steps.",
			},
		}
	case issue.Code == "unknown_field" && strings.EqualFold(issue.Field, "with.headers"):
		return &FlowRepairExample{
			ID:    "repair_dotted_with_headers",
			Title: "Nest headers under with",
			MatchingIssue: FlowRepairExampleIssue{
				Code:  issue.Code,
				Field: issue.Field,
			},
			Why:         "Dotted fields are not valid YAML keys in the Flow schema. TSPlay expects headers inside a with object.",
			SafeAutofix: true,
			BeforeYAML: `- action: write_csv
  file_path: "reports/out.csv"
  value: "{{rows}}"
  with.headers:
    - title`,
			AfterYAML: `- action: write_csv
  file_path: "reports/out.csv"
  with:
    value: "{{rows}}"
    headers:
      - title`,
			Edits: []FlowRepairExampleEdit{
				{Op: "move_field", Path: "steps[n].with.headers", To: "steps[n].with.headers"},
				{Op: "move_field", Path: "steps[n].value", To: "steps[n].with.value"},
			},
			Revalidate: true,
			Notes: []string{
				"write_csv accepts either value or with.value; keep related nested options under with for clarity.",
			},
		}
	case issue.Code == "unexpected_parameter" && strings.EqualFold(issue.Action, "navigate") && strings.EqualFold(issue.Field, "timeout"):
		return &FlowRepairExample{
			ID:    "repair_navigate_timeout_to_browser_timeout",
			Title: "Move navigate timeout to browser.timeout",
			MatchingIssue: FlowRepairExampleIssue{
				Code:   issue.Code,
				Action: issue.Action,
				Field:  issue.Field,
			},
			Why: "navigate only accepts url at the step level. Use browser.timeout for page-wide defaults or remove the timeout if it was copied from another framework.",
			BeforeYAML: `browser:
  timeout: 30000
steps:
  - action: navigate
    url: "https://example.com"
    timeout: 3000`,
			AfterYAML: `browser:
  timeout: 3000
steps:
  - action: navigate
    url: "https://example.com"`,
			Edits: []FlowRepairExampleEdit{
				{Op: "remove_field", Path: "steps[n].timeout"},
				{Op: "set_field", Path: "browser.timeout", To: 3000},
			},
			Revalidate: true,
			Notes: []string{
				"This is a guided repair, not a guaranteed autofix, because a browser-level timeout changes flow-wide behavior.",
			},
		}
	case issue.Code == "unsupported_action" && strings.EqualFold(issue.Action, "save_file"):
		return &FlowRepairExample{
			ID:    "repair_save_file_to_structured_writer",
			Title: "Replace save_file with a concrete writer action",
			MatchingIssue: FlowRepairExampleIssue{
				Code:   issue.Code,
				Action: issue.Action,
			},
			Why: "TSPlay does not have a generic save_file action. Pick the writer that matches the artifact type, such as write_csv, write_json, save_html, or a download action.",
			BeforeYAML: `- action: save_file
  path: "artifacts/results.csv"
  content: "{{rows}}"`,
			AfterYAML: `- action: write_csv
  file_path: "artifacts/results.csv"
  with:
    value: "{{rows}}"`,
			Edits: []FlowRepairExampleEdit{
				{Op: "replace", Path: "steps[n].action", From: "save_file", To: "write_csv"},
				{Op: "rename_field", Path: "steps[n].path", To: "file_path"},
				{Op: "move_field", Path: "steps[n].content", To: "steps[n].with.value"},
			},
			Revalidate: true,
			Notes: []string{
				"Switch to write_json when the value is an object or list meant to stay as JSON.",
				"Switch to save_html when the payload is page HTML.",
			},
		}
	default:
		return nil
	}
}

func BuiltinFlowRecommendedExamples() []FlowRecommendedExample {
	definitions := builtinFlowExampleDefinitions()
	result := make([]FlowRecommendedExample, 0, len(definitions))
	for _, definition := range definitions {
		result = append(result, definition.Example)
	}
	return result
}

func builtinFlowExampleDefinitions() []flowExampleDefinition {
	return []flowExampleDefinition{
		{
			Example: FlowRecommendedExample{
				ID:             "search_results_to_csv",
				Title:          "Search results to CSV",
				Intent:         "在搜索页面输入关键词，提取结果卡片，并保存为 CSV",
				WhyThisMatches: "用户意图里同时出现了搜索、结果提取和 CSV 导出。",
				Description:    "Open a search page, submit a query, extract result cards, and write structured rows to CSV.",
				WhenToUse:      "搜索页 + 结果列表 + 结构化导出",
				Tags:           []string{"search", "extract", "csv"},
				FocusActions:   []string{"navigate", "wait_for_selector", "type_text", "click", "evaluate", "write_csv"},
				RequiresAllow:  []string{"allow_javascript", "allow_file_access"},
				InputVars: []FlowExampleInputVar{
					{Name: "target_url", Required: true, Description: "Search page URL", Example: "https://www.google.com/"},
					{Name: "search_keyword", Required: true, Description: "Keyword to search for", Example: "人工智能"},
					{Name: "output_file", Required: true, Description: "CSV output path", Example: "artifacts/search_results.csv"},
				},
				CommonPitfalls: []FlowExamplePitfall{
					{Wrong: "fill", IssueCode: "unsupported_action", Fix: "Replace it with type_text."},
					{Wrong: "result_var", IssueCode: "unknown_field", Fix: "Use save_as for extracted results."},
					{Wrong: "save_file", IssueCode: "unsupported_action", Fix: "Use write_csv for CSV output."},
				},
				FlowYAML: `schema_version: "1"
name: search_results_to_csv
vars:
  target_url: "https://www.google.com/"
  search_keyword: "人工智能"
  output_file: "artifacts/search_results.csv"
steps:
  - name: open search page
    action: navigate
    url: "{{target_url}}"
  - name: wait for search box
    action: wait_for_selector
    selector: "textarea[name='q'], textarea#APjFqb, input[name='q']"
    timeout: 15000
  - name: type query
    action: type_text
    selector: "textarea[name='q'], textarea#APjFqb, input[name='q']"
    text: "{{search_keyword}}"
  - name: submit search
    action: click
    selector: "input[name='btnK'], button[type='submit']"
  - name: wait for results
    action: wait_for_selector
    selector: "div.g"
    timeout: 15000
  - name: collect result cards
    action: evaluate
    selector: "body"
    script: |
      const cards = Array.from(document.querySelectorAll('div.g'));
      return cards.map((card, index) => {
        const title = card.querySelector('h3')?.textContent?.trim() || '';
        const link = card.querySelector('a')?.href || '';
        const description = card.querySelector('.VwiC3b, .IsZvec')?.textContent?.trim() || '';
        return { index: index + 1, title, link, description };
      }).filter(item => item.title && item.link);
    save_as: search_results
  - name: write csv
    action: write_csv
    file_path: "{{output_file}}"
    with:
      value: "{{search_results}}"
      headers:
        - index
        - title
        - link
        - description
`,
			},
			MatchTerms:   []string{"search", "搜索", "query", "keyword", "关键词", "google", "谷歌", "results", "结果", "csv", "导出"},
			IssueCodes:   []string{"unsupported_action", "unknown_field"},
			IssueActions: []string{"fill", "save_file", "log"},
			IssueFields:  []string{"result_var"},
		},
		{
			Example: FlowRecommendedExample{
				ID:             "capture_table_to_csv",
				Title:          "Capture table to CSV",
				Intent:         "抓取页面表格并导出成 CSV",
				WhyThisMatches: "用户意图更像稳定表格抓取，不一定需要 JavaScript 自定义提取。",
				Description:    "Capture a stable table into rows and write them to a CSV file.",
				WhenToUse:      "后台列表页 / 管理台表格 / 稳定表头",
				Tags:           []string{"table", "csv", "capture"},
				FocusActions:   []string{"navigate", "wait_for_selector", "capture_table", "write_csv"},
				RequiresAllow:  []string{"allow_file_access"},
				InputVars: []FlowExampleInputVar{
					{Name: "table_url", Required: true, Description: "Page containing the table", Example: "https://example.com/orders"},
					{Name: "output_file", Required: true, Description: "CSV output path", Example: "artifacts/orders.csv"},
				},
				CommonPitfalls: []FlowExamplePitfall{
					{Wrong: "evaluate table DOM first", Fix: "Prefer capture_table when the page already exposes a real table."},
					{Wrong: "with.headers at step root", IssueCode: "unknown_field", Fix: "Nest headers under with."},
				},
				FlowYAML: `schema_version: "1"
name: capture_table_to_csv
vars:
  table_url: "https://example.com/orders"
  output_file: "artifacts/orders.csv"
steps:
  - name: open table page
    action: navigate
    url: "{{table_url}}"
  - name: wait for table
    action: wait_for_selector
    selector: "#orders-table"
    timeout: 10000
  - name: capture rows
    action: capture_table
    selector: "#orders-table"
    save_as: orders
  - name: write csv
    action: write_csv
    file_path: "{{output_file}}"
    with:
      value: "{{orders}}"
      headers:
        - order_id
        - status
        - total
`,
			},
			MatchTerms:   []string{"table", "表格", "grid", "datatable", "rows", "列表", "csv", "导出"},
			IssueCodes:   []string{"unknown_field", "unsupported_action"},
			IssueActions: []string{"save_file"},
			IssueFields:  []string{"with.headers"},
		},
		{
			Example: FlowRecommendedExample{
				ID:             "upload_file_then_submit",
				Title:          "Upload file then submit",
				Intent:         "上传文件后点击提交，并等待成功状态",
				WhyThisMatches: "用户意图涉及文件上传，通常需要 TODO 文件路径和文件权限。",
				Description:    "Open an import page, upload a local file, submit it, and wait for the success state.",
				WhenToUse:      "上传、导入、提交附件",
				Tags:           []string{"upload", "file", "submit"},
				FocusActions:   []string{"navigate", "wait_for_selector", "upload_file", "click"},
				RequiresAllow:  []string{"allow_file_access"},
				InputVars: []FlowExampleInputVar{
					{Name: "upload_url", Required: true, Description: "Upload page URL", Example: "https://example.com/import"},
					{Name: "upload_file_path", Required: true, Description: "Local file path to upload", Example: "TODO"},
				},
				CommonPitfalls: []FlowExamplePitfall{
					{Wrong: "omit upload_file_path", IssueCode: "missing_required_parameter", Fix: "Add a real file path variable before running."},
					{Wrong: "run without file permission", IssueCode: "security_policy", Fix: "Retry with allow_file_access=true only for trusted flows."},
				},
				FlowYAML: `schema_version: "1"
name: upload_file_then_submit
vars:
  upload_url: "https://example.com/import"
  upload_file_path: "TODO"
steps:
  - name: open import page
    action: navigate
    url: "{{upload_url}}"
  - name: wait for file input
    action: wait_for_selector
    selector: "input[type='file']"
    timeout: 10000
  - name: choose file
    action: upload_file
    selector: "input[type='file']"
    file_path: "{{upload_file_path}}"
  - name: submit import
    action: click
    selector: "button[type='submit']"
  - name: wait for success
    action: wait_for_selector
    selector: ".upload-success, .alert-success"
    timeout: 15000
`,
			},
			MatchTerms:   []string{"upload", "上传", "import", "导入", "file", "附件"},
			IssueCodes:   []string{"security_policy", "missing_required_parameter"},
			IssueActions: []string{"upload_file"},
		},
		{
			Example: FlowRecommendedExample{
				ID:             "login_and_download_report",
				Title:          "Login and download report",
				Intent:         "登录后台后导出报表到本地文件",
				WhyThisMatches: "用户意图同时包含登录、报表、下载，适合用下载动作而不是泛化保存动作。",
				Description:    "Log in, wait for the export control, and download a report file.",
				WhenToUse:      "登录后导出 / 下载报表 / 需要本地落盘",
				Tags:           []string{"login", "download", "report"},
				FocusActions:   []string{"navigate", "type_text", "click", "wait_for_selector", "download_file"},
				RequiresAllow:  []string{"allow_file_access"},
				InputVars: []FlowExampleInputVar{
					{Name: "login_url", Required: true, Description: "Login page URL", Example: "https://example.com/login"},
					{Name: "username", Required: true, Description: "Account username", Example: "demo"},
					{Name: "password", Required: true, Description: "Account password", Example: "demo-pass"},
					{Name: "save_path", Required: true, Description: "Local download path", Example: "artifacts/monthly-report.csv"},
				},
				CommonPitfalls: []FlowExamplePitfall{
					{Wrong: "save_file", IssueCode: "unsupported_action", Fix: "Use download_file or download_url for actual downloads."},
					{Wrong: "navigate.timeout", IssueCode: "unexpected_parameter", Fix: "Move timeout to browser.timeout."},
				},
				FlowYAML: `schema_version: "1"
name: login_and_download_report
vars:
  login_url: "https://example.com/login"
  username: "demo"
  password: "demo-pass"
  save_path: "artifacts/monthly-report.csv"
steps:
  - name: open login page
    action: navigate
    url: "{{login_url}}"
  - name: wait for username
    action: wait_for_selector
    selector: "#username"
    timeout: 10000
  - name: fill username
    action: type_text
    selector: "#username"
    text: "{{username}}"
  - name: fill password
    action: type_text
    selector: "#password"
    text: "{{password}}"
  - name: submit login
    action: click
    selector: "button[type='submit']"
  - name: wait for export button
    action: wait_for_selector
    selector: "text=\"Export report\""
    timeout: 15000
  - name: download report
    action: download_file
    selector: "text=\"Export report\""
    save_path: "{{save_path}}"
`,
			},
			MatchTerms:   []string{"login", "登录", "download", "下载", "report", "报表", "export", "导出"},
			IssueCodes:   []string{"unsupported_action", "unexpected_parameter", "security_policy"},
			IssueActions: []string{"save_file", "download_file"},
			IssueFields:  []string{"timeout"},
		},
		{
			Example: FlowRecommendedExample{
				ID:             "extract_text_then_branch",
				Title:          "Extract text then branch",
				Intent:         "提取页面文本到变量，再按条件分支",
				WhyThisMatches: "用户意图更像页面状态判断，适合先提取文本再 set_var / if，而不是直接写 lua 或 log。",
				Description:    "Extract business text into variables, derive a message, and branch with if based on visible page state.",
				WhenToUse:      "计数判断 / 状态分支 / 空态断言",
				Tags:           []string{"extract", "branch", "if"},
				FocusActions:   []string{"navigate", "wait_for_selector", "extract_text", "set_var", "if", "assert_text"},
				InputVars: []FlowExampleInputVar{
					{Name: "orders_url", Required: true, Description: "Page URL to inspect", Example: "https://example.com/orders"},
				},
				CommonPitfalls: []FlowExamplePitfall{
					{Wrong: "log", IssueCode: "unsupported_action", Fix: "Use set_var for derived messages and let the caller report completion."},
					{Wrong: "lua for simple branching", Fix: "Prefer extract_text + set_var + if when the logic is already expressible in Flow."},
				},
				FlowYAML: `schema_version: "1"
name: extract_text_then_branch
vars:
  orders_url: "https://example.com/orders"
steps:
  - name: open orders page
    action: navigate
    url: "{{orders_url}}"
  - name: wait for summary count
    action: wait_for_selector
    selector: ".summary .count"
    timeout: 10000
  - name: extract order count
    action: extract_text
    selector: ".summary .count"
    pattern: "([0-9]+)"
    save_as: order_count
  - name: build message
    action: set_var
    save_as: export_message
    value: "Current orders: {{order_count}}"
  - name: branch on count
    action: if
    condition:
      action: extract_text
      selector: ".summary .count"
      pattern: "[1-9][0-9]*"
    then:
      - action: click
        selector: "text=\"Export orders\""
    else:
      - action: assert_text
        selector: ".empty-state"
        text: "No orders"
`,
			},
			MatchTerms:   []string{"extract", "提取", "text", "文本", "count", "计数", "branch", "分支", "if", "condition", "判断", "状态"},
			IssueCodes:   []string{"unsupported_action"},
			IssueActions: []string{"log"},
		},
	}
}

func scoreFlowExampleDefinition(definition flowExampleDefinition, intentText string, actionSet map[string]struct{}, issue *FlowIssue) (float64, string) {
	score := 0.0
	reasons := []string{}

	if containsAnyToken(intentText, definition.MatchTerms) {
		score += 3
		reasons = append(reasons, definition.Example.WhyThisMatches)
	}
	for _, action := range definition.Example.FocusActions {
		if _, ok := actionSet[strings.ToLower(strings.TrimSpace(action))]; ok {
			score += 0.6
		}
	}
	if issue != nil {
		if stringSliceContainsFold(definition.IssueCodes, issue.Code) {
			score += 1.2
		}
		if stringSliceContainsFold(definition.IssueActions, issue.Action) {
			score += 2.2
			reasons = append(reasons, "当前 issue 的 action 和这个例子直接相关。")
		}
		if stringSliceContainsFold(definition.IssueFields, issue.Field) {
			score += 2.2
			reasons = append(reasons, "当前 issue 的字段问题和这个例子直接相关。")
		}
		for _, required := range definition.Example.RequiresAllow {
			if strings.Contains(strings.ToLower(issue.Message), strings.ToLower(required)) || strings.Contains(strings.ToLower(issue.Suggestion), strings.ToLower(required)) {
				score += 0.8
			}
		}
	}
	if len(reasons) == 0 {
		reasons = append(reasons, definition.Example.WhyThisMatches)
	}
	return score, strings.TrimSpace(reasons[0])
}

func collectFlowActionSet(draft *FlowDraft) map[string]struct{} {
	set := map[string]struct{}{}
	if draft == nil {
		return set
	}
	for _, action := range draft.PlannedActions {
		action = strings.ToLower(strings.TrimSpace(action))
		if action != "" {
			set[action] = struct{}{}
		}
	}
	if draft.Flow != nil {
		collectFlowStepActions(draft.Flow.Steps, set)
	}
	return set
}

func collectFlowStepActions(steps []FlowStep, set map[string]struct{}) {
	for _, step := range steps {
		action := strings.ToLower(strings.TrimSpace(step.Action))
		if action != "" {
			set[action] = struct{}{}
		}
		if len(step.Steps) > 0 {
			collectFlowStepActions(step.Steps, set)
		}
		if step.Condition != nil {
			collectFlowStepActions([]FlowStep{*step.Condition}, set)
		}
		if len(step.Then) > 0 {
			collectFlowStepActions(step.Then, set)
		}
		if len(step.Else) > 0 {
			collectFlowStepActions(step.Else, set)
		}
		if len(step.OnError) > 0 {
			collectFlowStepActions(step.OnError, set)
		}
	}
}

func containsAnyToken(text string, tokens []string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token != "" && strings.Contains(text, strings.ToLower(token)) {
			return true
		}
	}
	return false
}

func stringSliceContainsFold(values []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			return true
		}
	}
	return false
}

func scoreToRelevance(score float64) float64 {
	relevance := 0.45 + score*0.08
	if relevance > 0.99 {
		relevance = 0.99
	}
	if relevance < 0.45 {
		relevance = 0.45
	}
	return math.Round(relevance*100) / 100
}
