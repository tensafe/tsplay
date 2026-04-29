# TSPlay Flow Prompt Examples

## How To Use These

Use these templates when the user wants Codex to write or repair a TSPlay Flow but the request is still underspecified. Favor prompts that include the page URL, goal, inputs, expected output, and authorization boundary.

For repo-backed starting points grouped by category, read `example-index.md`.

## 中文优先用法

如果用户是中文提问，优先把需求整理成这 5 个字段：

- 页面
- 目标
- 输入
- 输出
- 授权

## Template 1: Write A New Flow

```text
帮我写一条 TSPlay Flow。
- 页面: <URL>
- 目标: <要完成的业务动作>
- 输入: <关键词、文件、筛选条件，若无则写无>
- 输出: <希望 save_as 什么，或写到哪个 artifacts 路径>
- 授权: <readonly / browser_write / full_automation，或 allow_* 说明>
- 交付要求: 生成一条可 review 的 `.flow.yaml`，变量名和步骤名要清楚
```

## Template 2: Repair An Existing Flow

```text
帮我修这条 TSPlay Flow。
- 文件: <flow 文件路径>
- 问题: <报错、超时、selector 失效、变量不对等>
- 预期: <修完后应该看到什么结果>
- 限制: <不要改业务意图 / 保持 artifact 路径 / 不要转成 Lua>
```

## Template 3: Turn A Requirement Into Flow

```text
把下面需求转成 TSPlay Flow。
- 页面: <URL>
- 需求: <自然语言业务需求>
- 输入: <运行时变量>
- 输出: <json/csv/excel/save_as>
- 授权边界: <最小权限>
- 风格: 优先复用仓库里已有教程模式，保持步骤可 review
```

## Template 4: Write A Tutorial-Style Flow

```text
帮我按 TSPlay 教程风格写一条 flow。
- lesson 编号: <例如 161>
- 主题: <断言、上传、下载、session、db、redis、mcp 等>
- 页面: <URL 或 demo 页面>
- 产物: <artifacts/tutorials/... 下的目标文件>
- 要求: name、save_as、artifact 路径都按教程风格写清楚
```

## Template 5: Add Session Reuse

```text
帮我把这条 TSPlay Flow 改成复用登录态。
- 文件: <flow 文件路径>
- session 名称: <例如 admin>
- 页面: <目标页面>
- 要求: 把会话配置放到顶层 browser，避免把登录步骤散在各个 step 里
```

## Template 6: Review A Flow For Readability

```text
帮我 review 这条 TSPlay Flow，重点不是能不能跑，而是好不好维护。
- 文件: <flow 文件路径>
- 检查点:
  - name 和 description 是否清楚
  - save_as 是否表达业务含义
  - artifact 路径是否稳定
  - 是否优先用了 Flow 而不是不必要的 Lua 绕路
```

## Template 7: Extract A Table

```text
帮我写一条 TSPlay Flow 来提取表格。
- 页面: <URL>
- 表格位置: <已知 selector，或写未知>
- 目标: 抓取表头和所有行
- 输出: 保存到变量、JSON、CSV，或三者都要
- 授权: <最小权限>
- 要求: 如果 selector 不确定，优先说明应该先走 observe_page 还是 capture_table
```

## Template 8: Build Through MCP Finalize

```text
帮我按 TSPlay MCP 的方式把需求收敛成 Flow。
- 页面: <URL>
- 意图: <自然语言需求>
- 输入: <变量>
- 授权: <readonly / browser_write / full_automation>
- 要求: 优先走 finalize_flow 思路，必要时再拆成 observe / draft / validate / repair
```

## Template 9: Send Email Notification

```text
帮我写一条 TSPlay Flow 发送邮件通知。
- 目标: <发送成功通知 / 失败告警 / 导出结果邮件 / 附件邮件>
- 收件人: <邮箱字符串或邮箱列表>
- 主题: <邮件主题>
- 内容: <正文文本或 HTML>
- 附件: <无，或附件文件路径 / 生成后的 artifact 路径>
- SMTP: <用 connection 走环境变量，或直接给 inline smtp 配置>
- 授权: 需要说明 allow_email；如果有附件，也说明文件权限
```

## Template 10: Read Local JSON

```text
帮我写一条 TSPlay Flow 读取本地 JSON。
- 文件: <json 文件路径>
- 目标: <提取状态 / 提取列表 / 驱动后续步骤 / 继续写出新结果>
- 需要的字段: <例如 payload.meta.status, payload.items[1].name>
- 输出: <save_as / JSON / CSV / 下游变量>
- 授权: 需要说明文件权限
```

## Good Input Shape

The highest-signal Flow requests usually include:

- Page URL or exact local demo page
- Business goal in one sentence
- Runtime inputs
- Expected output variables or artifact files
- Minimum allowed security boundary
- Whether this should become a tutorial, a one-off local Flow, or an MCP-facing Flow

## Good Starting Patterns In This Repo

- Minimal browser assert flow: `script/tutorials/10_assert_page_state.flow.yaml`
- Session reuse flow: `script/tutorials/50_use_session_batch_import_excel.flow.yaml`
- On-error recovery flow: `script/tutorials/27_on_error_import_excel_writeback.flow.yaml`
- Review readability example: `script/tutorials/131_review_readability_after.flow.yaml`
- MCP schema and example discovery: `docs/tutorials/112-mcp-flow-schema-and-examples.md`
- Local JSON read pattern: `tsplay_core/tsplay_table_test.go` in `TestRunFlowReadJSONValue`
- Email notification flow: `script/tutorials/send_email_qq_test.flow.yaml`
- Email attachment flow: `script/tutorials/send_email_qq_attachment_test.flow.yaml`

## Ready-Made Flow Snippets

### Snippet 1: Search And Assert

```yaml
schema_version: "1"
name: search_and_assert_flow
vars:
  page_url: https://example.com/search
  keyword: hello
steps:
  - action: navigate
    url: "{{page_url}}"

  - action: wait_for_selector
    selector: "#kw"
    timeout: 5000

  - action: type_text
    selector: "#kw"
    text: "{{keyword}}"

  - action: click
    selector: "#su"

  - action: assert_visible
    selector: "#results"
    timeout: 5000
```

### Snippet 2: Extract Text And Write JSON

```yaml
schema_version: "1"
name: extract_summary_flow
vars:
  page_url: http://127.0.0.1:8000/demo/extract.html
steps:
  - action: navigate
    url: "{{page_url}}"

  - action: wait_for_selector
    selector: "#notice"
    timeout: 5000

  - action: extract_text
    selector: "#page-title"
    timeout: 5000
    save_as: page_title

  - action: set_var
    save_as: payload
    with:
      value:
        page_title: "{{page_title}}"

  - action: write_json
    file_path: artifacts/extract-summary.json
    value: "{{payload}}"
```

### Snippet 3: Read Excel, Loop, And Recover Locally

```yaml
schema_version: "1"
name: import_with_recovery_flow
vars:
  page_url: http://127.0.0.1:8000/demo/import_workflow.html
  input_file: demo/data/import_users.xlsx
  import_results: []
steps:
  - action: read_excel
    file_path: "{{input_file}}"
    sheet: Users
    with:
      row_number_field: source_row
    save_as: rows

  - action: navigate
    url: "{{page_url}}"

  - action: foreach
    items: "{{rows}}"
    item_var: row
    steps:
      - action: on_error
        steps:
          - action: type_text
            selector: "#name"
            text: "{{row.name}}"
          - action: type_text
            selector: "#phone"
            text: "{{row.phone}}"
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
```

### Snippet 4: Send A Completion Email

```yaml
schema_version: "1"
name: send_completion_email_flow
vars:
  recipient_emails:
    - "ops@example.com"
  subject_text: "TSPlay run finished"
  body_text: "Import completed successfully."
steps:
  - action: send_email
    save_as: email_result
    connection: alerts
    with:
      to: "{{recipient_emails}}"
      subject: "{{subject_text}}"
      body: "{{body_text}}"
```

### Snippet 5: Read Local JSON And Reuse Fields

```yaml
schema_version: "1"
name: read_local_json_flow
steps:
  - action: read_json
    file_path: artifacts/payload.json
    save_as: payload

  - action: set_var
    save_as: status
    value: "{{payload.meta.status}}"

  - action: set_var
    save_as: second_name
    value: "{{payload.items[1].name}}"
```

## 中文使用建议

- 如果目标是“能跑一次”，先写最小 Flow，不要一上来就做很重的抽象。
- 如果目标是“给团队复用”，优先把 `name`、`description`、`save_as`、artifact 路径写清楚。
- 如果用户不知道 selector，先考虑 MCP 的 `observe_page` 路线，不要逼用户自己贴 HTML。
- 如果目标是“读本地 JSON 再继续处理”，优先先确认文件路径和要用的字段路径。
- 如果目标是“发通知”，优先先决定 SMTP 来源是 `connection` 还是 `with.smtp`，再决定是否带附件。
