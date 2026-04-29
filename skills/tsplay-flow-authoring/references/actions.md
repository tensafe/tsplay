# TSPlay Flow Actions Cheat Sheet

## How To Use This File

Use this file when you know the business goal but need help choosing the right TSPlay Flow action or remembering the correct YAML shape. These are the most common actions for coder-facing Flow authoring, not the full action catalog.

## 中文动作快查

- 打开页面: `navigate`
- 等页面准备好: `wait_for_selector`
- 输入内容: `type_text`
- 点击按钮或链接: `click`
- 断言元素可见: `assert_visible`
- 断言页面文本: `assert_text`
- 提取文本或数字: `extract_text`
- 抓表格结构化数据: `capture_table`
- 保存一个变量: `set_var`
- 追加结果列表: `append_var`
- 遍历多行数据: `foreach`
- 局部容错继续执行: `on_error`
- 轮询直到满足条件: `wait_until`
- 重试易抖动步骤: `retry`
- 读 JSON / CSV / Excel: `read_json`, `read_csv`, `read_excel`
- 写 JSON 或 CSV: `write_json`, `write_csv`
- 发邮件通知: `send_email`

## Authoring Rules From TSPlay

- Map user intent to page states first, then choose actions.
- Extract page values into variables before branching or looping on them.
- Prefer `save_as` names that describe business meaning, not DOM details.
- Add assertions around the business result, not only around low-level clicks.
- Keep recovery local with `on_error` instead of rewriting the whole Flow.
- Use `type_text`, not `fill`.
- Do not put `timeout` on `navigate`; use `browser.timeout` or downstream waits and assertions.

## 1. Navigation And Readiness / 导航与就绪

### `navigate`

Use to open a page.

```yaml
- action: navigate
  url: "{{page_url}}"
```

Required fields:

- `url`

Notes:

- Do not put `timeout` on this step itself.

### `wait_for_selector`

Use when the page must reach a visible or ready state before the next step.

```yaml
- action: wait_for_selector
  selector: "#import-form"
  timeout: 5000
```

Required fields:

- `selector`

Optional fields:

- `timeout`

### `click`

Use to trigger a button, link, tab, or other clickable control.

```yaml
- action: click
  selector: "#submit"
```

Required fields:

- `selector`

### `type_text`

Use to enter text into an input or textarea.

```yaml
- action: type_text
  selector: "#name"
  text: "{{row.name}}"
```

Required fields:

- `selector`
- `text`

## 2. Assertions And Extraction / 断言与提取

### `assert_visible`

Use when the business signal is that an element appears.

```yaml
- action: assert_visible
  selector: "#export-result"
  timeout: 5000
```

Required fields:

- `selector`

Optional fields:

- `timeout`

### `assert_text`

Use when the business signal is a known success or status text.

```yaml
- action: assert_text
  selector: "#submit-status"
  text: "Imported"
  timeout: 5000
```

Required fields:

- `selector`
- `text`

Optional fields:

- `timeout`

### `extract_text`

Use to read visible text into a variable. Add `pattern` when only part of the text matters.

```yaml
- action: extract_text
  selector: "#summary-count"
  timeout: 5000
  pattern: "([0-9]+)"
  save_as: order_count
```

Required fields:

- `selector`

Common optional fields:

- `timeout`
- `pattern`
- `save_as`

### `capture_table`

Use when the page already has stable table markup and later steps need structured rows.

```yaml
- action: capture_table
  selector: "#orders-table"
  save_as: orders
```

Required fields:

- `selector`

Recommended fields:

- `save_as`

## 3. Variables And Output Shaping / 变量与输出组织

### `set_var`

Use to create one variable from a resolved value. Prefer `with.value` for objects, lists, numbers, or booleans.

```yaml
- action: set_var
  save_as: payload
  with:
    value:
      page_title: "{{page_title}}"
      order_count: "{{order_count}}"
```

Required fields:

- `save_as`
- `with.value` for non-string literals

Notes:

- Use plain `value` when setting a string or placeholder directly.
- Use `with.value` when shaping a JSON-like object.

### `append_var`

Use to build a list ledger such as import results or verification rows.

```yaml
- action: append_var
  save_as: import_results
  with:
    value:
      source_row: "{{row.source_row}}"
      status: success
```

Required fields:

- `save_as`
- `with.value` for objects or lists

Notes:

- The target list is created automatically if it does not exist yet.

## 4. Control Flow / 控制流

### `foreach`

Use to process each item in a list such as CSV rows, Excel rows, order ids, or scraped records.

```yaml
- action: foreach
  items: "{{rows}}"
  item_var: row
  steps:
    - action: type_text
      selector: "#name"
      text: "{{row.name}}"
```

Required fields:

- `items`
- `item_var`
- `steps`

Optional fields:

- `index_var`
- `with.progress_key`

Notes:

- Use `with.progress_key` when resumable progress checkpoints matter.

### `on_error`

Use when one nested task may fail and you want to handle the failure locally instead of aborting the whole Flow.

```yaml
- action: on_error
  steps:
    - action: click
      selector: "#submit"
    - action: assert_text
      selector: "#submit-status"
      text: "Imported"
      timeout: 1000
  on_error:
    - action: append_var
      save_as: import_results
      with:
        value:
          status: failed
          error: "{{last_error}}"
```

Required fields:

- `steps`
- `on_error`

### `wait_until`

Use when a status may become true only after polling.

```yaml
- action: wait_until
  timeout: 10000
  interval_ms: 1000
  condition:
    action: is_visible
    selector: "#ready-badge"
```

Required fields:

- `condition`

Common optional fields:

- `timeout`
- `interval_ms`

### `retry`

Use when one flaky interaction may succeed after a short reattempt window.

```yaml
- action: retry
  times: 3
  interval_ms: 1000
  steps:
    - action: click
      selector: 'text="Export orders"'
    - action: assert_visible
      selector: "#export-result"
      timeout: 5000
```

Required fields:

- `times`
- `steps`

Useful optional fields:

- `interval_ms`

## 5. File Input And Output / 文件输入输出

### `read_json`

Use to read any local JSON file into a Flow variable. This works well for local config, previous run artifacts, handoff payloads, or decoded API results saved earlier.

```yaml
- action: read_json
  file_path: artifacts/import-results.json
  save_as: payload
```

Required fields:

- `file_path`

Recommended fields:

- `save_as`

Notes:

- Returns the decoded JSON value as an object, list, primitive, or null depending on the file content.
- UTF-8 BOM is accepted and ignored when present.
- After `save_as`, later steps can reference nested fields such as `{{payload.meta.status}}` or `{{payload.items[1].name}}`.
- Restricted Flow or MCP contexts require file access permission.

### `read_csv`

Use to load structured CSV rows into a list variable.

```yaml
- action: read_csv
  file_path: demo/data/users.csv
  with:
    row_number_field: source_row
  save_as: rows
```

Required fields:

- `file_path`

Useful optional fields:

- `with.start_row`
- `with.limit`
- `with.row_number_field`
- `save_as`

Note:

- Requires file access permission in MCP mode.

### `read_excel`

Use to load `.xlsx` rows. Add `sheet`, `range`, or explicit `headers` when the workbook is not simple.

```yaml
- action: read_excel
  file_path: "{{input_file}}"
  sheet: Users
  with:
    row_number_field: source_row
  save_as: rows
```

Required fields:

- `file_path`

Useful optional fields:

- `sheet`
- `range`
- `with.headers`
- `with.start_row`
- `with.limit`
- `with.row_number_field`
- `save_as`

Notes:

- Omit `range` to read the whole sheet and use the first non-empty row as headers.
- Use `with.headers` when the chosen range contains data rows but not a header row.
- Requires file access permission in MCP mode.

### `write_json`

Use to write any resolved value to a JSON artifact.

```yaml
- action: write_json
  file_path: artifacts/import-summary.json
  value: "{{payload}}"
```

Required fields:

- `file_path`
- `value`

Note:

- Use `with.value` when writing an object literal directly.

### `write_csv`

Use to write row lists or result ledgers to CSV.

```yaml
- action: write_csv
  file_path: artifacts/import-results.csv
  with:
    value: "{{import_results}}"
    headers:
      - source_row
      - status
      - error
```

Required fields:

- `file_path`
- `value` or `with.value`

Useful optional fields:

- `with.headers`

## 6. Notification And Delivery / 通知与投递

### `send_email`

Use to send a completion summary, failure alert, exported report, or other notification email through SMTP.

```yaml
- action: send_email
  save_as: email_result
  connection: alerts
  with:
    to:
      - "ops@example.com"
    subject: "TSPlay run finished"
    body: "Import completed successfully."
```

Required fields:

- `with.to`
- `with.subject`
- at least one of `with.body` or `with.html`

Useful optional fields:

- `with.cc`
- `with.bcc`
- `with.attachments`
- `with.from_email`
- `with.reply_to`
- `with.headers`
- `connection`
- `with.smtp`
- `save_as`

Notes:

- `with.to`, `with.cc`, and `with.bcc` can be one email string or a list of email strings.
- `connection` loads SMTP settings from `TSPLAY_EMAIL_*` or `TSPLAY_EMAIL_<NAME>_*` environment variables.
- `with.smtp` supports inline SMTP settings such as `host`, `port`, `username`, `password`, `from`, and `tls_mode`.
- Restricted Flow or MCP contexts require `allow_email=true`.
- Attachments read local files and therefore also require `allow_file_access=true` in restricted Flow or MCP contexts.
- Custom headers cannot override reserved headers such as `From`, `To`, `Subject`, or `Content-Type`.

Attachment example:

```yaml
- action: send_email
  save_as: email_result
  with:
    to: "{{recipient_emails}}"
    subject: "{{subject_text}}"
    body: "{{body_text}}"
    attachments:
      - path: artifacts/reports/run.xlsx
        name: run.xlsx
        content_type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
    smtp:
      host: smtp.qq.com
      port: 465
      username: "{{sender_email}}"
      password: "{{sender_password}}"
      from: "{{sender_email}}"
      tls_mode: tls
```

## 7. Top-Level Browser Config / 顶层浏览器配置

### `browser.use_session`

Use at the top of the Flow when login state should be reused instead of replayed.

```yaml
browser:
  use_session: admin
```

### `browser.timeout`

Use at the top of the Flow when many steps share the same readiness budget.

```yaml
browser:
  timeout: 30000
```

## Good Starter Combos

- Search or form submit: `navigate` + `wait_for_selector` + `type_text` + `click` + `assert_text`
- Simple scrape: `navigate` + `wait_for_selector` + `extract_text` + `set_var` + `write_json`
- Local JSON driven flow: `read_json` + `set_var` or direct variable references
- Table flow: `navigate` + `wait_for_selector` + `capture_table`
- Batch import: `read_excel` or `read_csv` + `foreach` + `append_var` + `write_json` + `write_csv`
- Resilient import: `foreach` + `on_error` + `append_var`
- Report notification: `set_var` or `write_excel` + `send_email`

## 中文起手组合

- 搜索或提交表单: `navigate` + `wait_for_selector` + `type_text` + `click` + `assert_text`
- 提取页面摘要: `navigate` + `wait_for_selector` + `extract_text` + `set_var` + `write_json`
- 本地 JSON 驱动流程: `read_json` + `set_var` 或直接变量引用
- 抓取表格: `navigate` + `wait_for_selector` + `capture_table`
- 批量导入: `read_excel` 或 `read_csv` + `foreach` + `append_var` + `write_json` + `write_csv`
- 带容错的导入: `foreach` + `on_error` + `append_var`
- 结果邮件通知: `set_var` 或 `write_excel` + `send_email`
