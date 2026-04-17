# TSPlay

TSPlay 是一个基于 Go + Playwright 的浏览器自动化执行引擎，提供三层互相配合的能力：

- `Lua CLI / Lua Script`：适合临时调试、手工探索和快速验证
- `Flow DSL`：适合版本化、可审查、可被 AI 生成的结构化流程
- `MCP Server`：适合接入 OpenClaw、Codex 等 Agent，让模型先观察页面，再生成、校验、执行和修复流程

它的重点不是只把浏览器动作包成一堆函数，而是把“执行”“校验”“失败留痕”“Agent 集成”放进同一条链路里。

## 适用场景

- Web RPA：登录、表单填写、点击、下载、上传
- 页面数据提取：链接、属性、表格、HTML、Cookie、Storage State
- 带断言和重试的业务流程自动化
- 给大模型提供页面观察、Flow 生成、执行和修复能力

## 核心能力

- 基于 Playwright 驱动 Chromium
- 支持 Lua 脚本直接控制浏览器
- 支持结构化 Flow YAML/JSON
- Flow 支持变量替换、`save_as`、`extract_text`、`set_var`、`append_var`、`retry`、`if`、`foreach`、`on_error`、`wait_until`、`http_request`、`json_extract`、`read_csv`、`read_excel`、`write_json`、`write_csv`、`redis_get/set/del/incr`、`db_insert`、断言和失败 trace
- 失败时自动落盘现场资料：截图、HTML、DOM snapshot
- 支持基于“用户意图 + 页面观察”自动草拟 Flow
- 可作为 MCP Server 暴露给 Agent 调用
- MCP 模式带安全边界，可按能力显式授权
- 附带面向实施、测试、自动化开发和讲师的培训体系文档

## 运行模式

| 模式 | 入口 | 适合场景 |
| --- | --- | --- |
| 交互式 CLI | `go run . -action cli` | 手动调试、边试边写 |
| Lua 脚本 | `go run . -script script/open_url.lua` | 自定义逻辑、一次性任务 |
| Flow DSL | `go run . -flow script/demo_baidu.flow.yaml` | 结构化流程、版本管理、AI 生成 |
| MCP Server | `go run . -action srv` | 接入 OpenClaw / Codex / 其他 Agent |

## 文档与培训入口

| 内容 | 说明 | 入口 |
| --- | --- | --- |
| 项目总览 | 快速理解 TSPlay 的模式、能力和运行方式 | [ReadMe.md](ReadMe.md) |
| 文档索引 | 查看仓库内的文档地图和推荐阅读顺序 | [docs/README.md](docs/README.md) |
| 培训体系总览 | 面向实施、测试、开发和讲师的统一培训入口 | [docs/training/README.md](docs/training/README.md) |
| 学习路径 | 按角色和能力等级安排学习与晋级 | [docs/training/learning-path.md](docs/training/learning-path.md) |
| Bootcamp 课程表 | 2 天训练营和 4 周落地节奏 | [docs/training/bootcamp-plan.md](docs/training/bootcamp-plan.md) |
| 实训实验 | 基于 `demo/` 和 `script/` 的练习清单 | [docs/training/labs.md](docs/training/labs.md) |
| 考核与认证 | 评分维度、晋级门槛和结业标准 | [docs/training/assessment.md](docs/training/assessment.md) |
| 讲师手册 | 讲师备课、授课、复盘和版本维护指南 | [docs/training/trainer-playbook.md](docs/training/trainer-playbook.md) |

## 快速开始

### 环境要求

- Go `1.23.6+`
- 能运行 Playwright Chromium 的系统环境
- 首次启动时程序会尝试自动执行 `playwright.Install()` 下载浏览器

### 1. 拉取依赖

```bash
go mod download
```

### 2. 启动交互式 CLI

```bash
go run . -action cli
```

启动后先输入：

```text
start
```

然后可以直接执行 Lua 风格命令：

```lua
navigate("https://www.baidu.com")
wait_for_network_idle()
type_text("#kw", "山东大学")
click("#su")
```

如果想隐藏浏览器窗口：

```bash
go run . -action cli -headless
```

### 3. 运行 Lua 脚本

```bash
go run . -script script/open_url.lua
go run . -script script/open_url.lua -headless
```

### 4. 运行 Flow

```bash
go run . -flow script/demo_baidu.flow.yaml
go run . -flow script/demo_baidu.flow.yaml -headless
```

命令行执行 Flow 后会输出结构化 JSON 结果，其中包含变量、每一步 trace、执行耗时和失败现场路径。

## Lua 示例

```lua
print("hello!")
navigate("https://www.baidu.com")
wait_for_network_idle()
type_text("#kw", "山东大学")
click("#su")
wait_for_network_idle()

links = get_all_links("xpath=//body")
print("links count:", #links)
```

## Flow DSL

Flow 更适合作为长期维护的业务资产。它比裸 Lua 更容易：

- 让 AI 严格生成
- 被人工审查
- 做版本 diff
- 记录失败上下文
- 配合 MCP 修复流程

### 最小示例

```yaml
schema_version: "1"
name: baidu_search
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

  - action: wait_for_network_idle

  - action: get_all_links
    selector: "xpath=//body"
    save_as: links
```

### 关键字段

| 字段 | 说明 |
| --- | --- |
| `schema_version` | 必填，当前版本固定为 `"1"` |
| `name` | Flow 名称 |
| `description` | 可选描述 |
| `vars` | 初始变量，支持在步骤中用 `{{var_name}}` 引用 |
| `steps` | 按顺序执行的步骤列表 |
| `save_as` | 把动作返回值保存为变量，供后续步骤复用 |
| `continue_on_error` | 当前步骤失败后继续往下执行 |

### 参数写法

Flow 同时支持命名参数和 `args` 两种写法：

```yaml
- action: type_text
  selector: "#kw"
  text: "{{query}}"

- action: type_text
  args: ["#kw", "{{query}}"]
```

### 常用控制能力

先用 `extract_text + set_var` 把页面上的文本或数字变成业务变量：

```yaml
- action: extract_text
  selector: ".summary .count"
  pattern: '([0-9]+)'
  save_as: order_count

- action: set_var
  save_as: export_message
  value: "Current orders: {{order_count}}"
```

使用 `retry` 包住不稳定步骤：

```yaml
- action: retry
  times: 3
  interval_ms: 1000
  steps:
    - action: click
      selector: 'text="导出"'
    - action: assert_visible
      selector: "#export-result"
      timeout: 5000
```

使用 `if` 处理可选弹窗或不同页面状态：

```yaml
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
      selector: "#main"
      timeout: 5000
```

使用 `foreach` 处理列表数据，循环变量只在循环期间可见：

```yaml
- action: foreach
  items: "{{order_ids}}"
  item_var: order_id
  index_var: order_index
  steps:
    - action: type_text
      selector: "#order-id"
      text: "{{order_id}}"
    - action: click
      selector: "#search"
```

使用 `on_error` 做局部失败恢复，handler 中可以引用 `{{last_error}}`：

```yaml
- action: on_error
  steps:
    - action: click
      selector: 'text="导出"'
  on_error:
    - action: reload
    - action: wait_for_selector
      selector: "#main"
      timeout: 10000
```

使用 `wait_until` 轮询条件：

```yaml
- action: wait_until
  timeout: 30000
  interval_ms: 500
  condition:
    action: is_visible
    selector: "#export-result"
```

### 调三方 HTTP / OCR

可以直接在 Flow 里发起 HTTP 请求，再把 JSON 结果继续编排进页面流程：

```yaml
- action: screenshot_element
  selector: "#captcha-img"
  path: "captcha/login.png"

- action: http_request
  url: "https://ocr.example.com/recognize"
  save_as: ocr_result
  with:
    method: POST
    headers:
      Authorization: "Bearer {{ocr_token}}"
    multipart_files:
      image: "captcha/login.png"
    response_as: json

- action: json_extract
  from: "{{ocr_result}}"
  path: "$.body.text"
  save_as: captcha_text

- action: type_text
  selector: "#captcha-input"
  text: "{{captcha_text}}"
```

这组能力适合：

- OCR 验证码识别
- 调内部 HTTP API 查单、补数、提交工单
- 调 webhook / 通知接口
- 浏览器取数后，再用接口补充结构化数据

### 用 Redis 共享会话 / 游标

如果 cookie、session、游标或去重标记是提前写进 Redis 的，可以直接在 Flow 里读取：

```yaml
- action: redis_get
  key: "sessions:admin_cookie"
  save_as: cookie_header

- action: redis_get
  key: "sessions:admin_payload"
  save_as: session_payload

- action: json_extract
  from: "{{session_payload}}"
  path: "$.cookie"
  save_as: cookie_value

- action: redis_set
  key: "orders:last_processed_id"
  value: "{{current_order_id}}"
  ttl_seconds: 86400
```

如果要写结构化对象，推荐把对象放进 `with.value`：

```yaml
- action: redis_set
  key: "sessions:admin_payload"
  with:
    value:
      cookie: "SESSION=abc"
      user: "admin"
```

Redis 连接默认从环境变量读取：

- 默认连接：
  - `TSPLAY_REDIS_ADDR`
  - `TSPLAY_REDIS_USERNAME`
  - `TSPLAY_REDIS_PASSWORD`
  - `TSPLAY_REDIS_DB`
- 命名连接，例如 `connection: "sessions"`：
  - `TSPLAY_REDIS_SESSIONS_ADDR`
  - `TSPLAY_REDIS_SESSIONS_USERNAME`
  - `TSPLAY_REDIS_SESSIONS_PASSWORD`
  - `TSPLAY_REDIS_SESSIONS_DB`

如果更习惯 URL，也可以用：

- `TSPLAY_REDIS_URL`
- `TSPLAY_REDIS_SESSIONS_URL`

### 从本地 CSV / Excel 批量录入

如果需要把本地表格按行录入网页，可以先把文件读成列表，再配合 `foreach`：

```yaml
- action: read_csv
  file_path: imports/users.csv
  with:
    row_number_field: source_row
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
```

Excel 用法类似：

```yaml
- action: read_excel
  file_path: imports/users.xlsx
  sheet: Users
  with:
    row_number_field: source_row
  save_as: rows
```

如果 Excel 里前面有说明、标题或多块表，可以只读指定范围：

```yaml
- action: read_excel
  file_path: imports/users.xlsx
  sheet: Users
  range: A2:B20
  with:
    headers:
      - name
      - phone
  save_as: rows
```

如果要做分段导入或断点续跑，可以按源文件行号切片：

```yaml
- action: read_excel
  file_path: imports/users.xlsx
  sheet: Users
  with:
    start_row: 102
    limit: 50
    row_number_field: source_row
  save_as: rows
```

如果希望循环过程中自动记住下次该从哪一行继续，也可以让 `foreach` 在每次成功后顺手写一个 Redis checkpoint。Redis 没配置时，这一步会自动跳过，不影响主流程：

```yaml
- action: foreach
  items: "{{rows}}"
  item_var: row
  with:
    progress_key: imports:users:resume_row
  steps:
    - action: type_text
      selector: "#name"
      text: "{{row.name}}"
    - action: click
      selector: "#submit"
```

导入结果也可以边跑边累计，最后统一写回 `json/csv`：

```yaml
- action: foreach
  items: "{{rows}}"
  item_var: row
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
```

规则说明：

- `read_csv` 会把第一行非空行当成表头
- `read_excel` 在不传 `range` 时，也会把第一行非空行当成表头
- `read_excel.range` 支持 `A2:B20` 这类矩形区域
- 如果 `range` 本身包含表头，就直接读取该范围
- 如果 `range` 只有数据行，可通过 `with.headers` 显式指定列名
- `with.start_row` 按源文件中的真实行号恢复，例如 CSV 第 3 行或 Excel 第 102 行
- `with.limit` 可以把一次导入切成固定批次
- `with.row_number_field` 会把源文件行号写回到每条记录里，方便结果台账和断点续跑
- `foreach.with.progress_key` 会在每次成功迭代后写入“下一条应继续的源行号”；如果记录里没有 `source_row/row_number/row`，就退回到下一次迭代序号
- `foreach.with.progress_connection` 可以指定 Redis 连接别名；不写则使用默认连接
- `foreach.with.progress_value` 可以覆盖默认 checkpoint 值，例如手工写业务游标
- `foreach` checkpoint 需要 `allow_redis=true`，但如果环境里没有配置对应 Redis 连接，会自动跳过，不会让导入失败
- `append_var` 用来累积每条导入记录的结果
- `write_json` 和 `write_csv` 用来把结果台账落到本地文件
- 返回结果是 `foreach` 可直接遍历的行对象列表
- 表头是 `name`、`phone` 这种简单字段时，可直接写 `{{row.name}}`
- 如果表头里有空格或特殊字符，可写成 `{{row["User Name"]}}`
- `read_excel` 当前支持 `.xlsx`，暂不支持老式 `.xls`

Flow 会在执行前做严格校验，包括：

- action 是否存在
- 参数数量和类型是否匹配
- 变量引用是否有效
- `save_as` 名称是否合法
- 高风险能力是否被授权

## 执行结果与失败现场

每次执行 Flow 都会返回 step trace，通常包含：

- `action`
- 参数摘要
- `status`
- `duration_ms`
- 输出摘要
- 当前 `page_url`

步骤失败时，TSPlay 会把现场资料写入 artifact root，默认目录是 `artifacts/`。常见产物包括：

- `failure.png`
- `page.html`
- `dom_snapshot.json`

可以通过命令行参数修改目录：

```bash
go run . -flow script/demo_baidu.flow.yaml -artifact-root artifacts
```

## MCP / Agent 集成

TSPlay 可以作为 MCP Server 启动，让 Agent 不必直接读整页 HTML，也不必手写 selector。

### 启动方式

```bash
go run . -action srv
go run . -action srv -addr :8081
go run . -action srv -flow-root script -artifact-root artifacts
```

### 当前暴露的 MCP 工具

| 工具名 | 说明 |
| --- | --- |
| `tsplay.list_actions` | 返回 Flow 可用 action 及参数 schema |
| `tsplay.flow_schema` | 返回 Flow JSON Schema、生成规则、selector 策略、authoring checklist 和 action manifest |
| `tsplay.flow_examples` | 返回带 focus_actions 的参考示例和示例选择提示 |
| `tsplay.draft_flow` | 输入用户意图和页面 URL / observation，自动观察页面、草拟 Flow、自动校验，并在必要时做一轮 selector 修正 |
| `tsplay.delete_session` | 删除命名浏览器会话注册；storage-state 会话会顺手删除复制出的 state 文件，persistent profile 数据会保留 |
| `tsplay.export_session_flow_snippet` | 导出单个命名会话对应的可复制 `browser` / Flow 片段；支持 `browser` / `flow`、展开版，以及 YAML / JSON |
| `tsplay.get_session` | 读取单个命名浏览器会话的详情，返回展开后的 browser 配置和物理路径 |
| `tsplay.list_sessions` | 列出已保存的命名浏览器会话，并返回可直接写进 Flow 的 `browser.use_session` 片段 |
| `tsplay.observe_page` | 打开页面并返回截图路径、DOM snapshot、可交互元素和候选 selector |
| `tsplay.repair_flow` | 输入原始 Flow + repair_hints + 可选失败现场，直接返回统一的 repair_request 和可喂给模型的 prompt |
| `tsplay.repair_flow_context` | 根据失败 Flow 和 run result 组织修复上下文，附带失败分类、修复焦点、统一结构的 repair_hints 和校验清单 |
| `tsplay.validate_flow` | 只校验 Flow，不启动浏览器 |
| `tsplay.run_flow` | 启动 Playwright 执行 Flow，并返回 trace |
| `tsplay.save_session` | 把 storage state 或 persistent profile 注册成命名会话，后续 Flow 可直接按名字复用 |

### 推荐给 Agent 的调用顺序

1. 调 `tsplay.flow_schema`，拿到严格约束
2. 调 `tsplay.flow_examples`，拿到参考模板
3. 如果已经有明确 URL 和用户意图，优先调 `tsplay.draft_flow`
4. 如果需要更细粒度控制，也可以先调 `tsplay.observe_page` 再把 observation 传给 `tsplay.draft_flow`
5. 查看 `tsplay.draft_flow` 返回里的 `validation`、`selector_repairs` 和 `repair_hints`
6. 如需单独二次校验，再调 `tsplay.validate_flow`
7. 如果业务要长期复用登录态，可在成功后用 `tsplay.save_session`
8. 调 `tsplay.list_sessions` 看命名会话、最近使用时间和来源说明
9. 如需查看某个会话的展开配置和物理路径，再调 `tsplay.get_session`
10. 如果想直接拿可粘贴的 YAML，再调 `tsplay.export_session_flow_snippet`
11. 后续 Flow 直接用 `browser.use_session`
12. 不再需要时可调 `tsplay.delete_session`
13. 成功后再调用 `tsplay.run_flow`
14. 失败时可直接调 `tsplay.repair_flow`，或者先调 `tsplay.repair_flow_context` 再交给 `tsplay.repair_flow`

`tsplay.draft_flow` 和 `tsplay.repair_flow_context` 现在都会返回同一种 `repair_hints` 结构，而 `tsplay.repair_flow` 会把原始 Flow、repair_hints、失败现场、规则和 prompt 统一收口成一个 repair_request。

### 从用户意图草拟 Flow

`tsplay.draft_flow` 适合先把“用户想做什么”落成一份可审阅的草稿：

```json
{
  "intent": "搜索订单并导出",
  "url": "https://example.com/orders"
}
```

工具会：

- 自动打开页面并观察可交互元素
- 尽量匹配搜索、导出、上传、登录、下拉选择等常见动作
- 生成一份结构化 Flow YAML
- 自动跑一遍与 `tsplay.validate_flow` 对齐的校验
- 如果 observation 里存在更稳的候选 selector，会自动修正一轮
- 如果自动校验还没过，会直接返回按 step 排序的 `repair_hints`
- 标出匹配到的元素、建议变量、假设项和未解决部分

如果希望校验放开高风险动作，也可以和 `tsplay.validate_flow` 一样传：

- `allow_lua`
- `allow_javascript`
- `allow_file_access`
- `allow_browser_state`
- `allow_http`
- `allow_redis`

如果已经提前调过 `tsplay.observe_page`，也可以把 observation JSON 直接传给 `tsplay.draft_flow`，避免重复打开页面。

### 安全边界

MCP 模式默认不是全放开。高风险能力需要在单次请求里显式授权：

| 授权参数 | 放行动作 |
| --- | --- |
| `allow_lua=true` | `lua` |
| `allow_javascript=true` | `execute_script`、`evaluate` |
| `allow_file_access=true` | `screenshot`、`screenshot_element`、`save_html`、`read_csv`、`read_excel`、`upload_file`、`upload_multiple_files`、`download_file`、`download_url` |
| `allow_browser_state=true` | `get_storage_state`、`get_cookies_string`、Flow 顶层 `browser.use_session` / `storage_state` / `save_storage_state` / `persistent profile` |
| `allow_http=true` | `http_request` |
| `allow_redis=true` | `redis_get`、`redis_set`、`redis_del`、`redis_incr`、`foreach.with.progress_key` |
| `allow_database=true` | `db_insert` |

补充说明：

- `flow_path` 默认只允许读取 `script/` 目录内的文件，可用 `-flow-root` 调整
- 文件类动作即使被授权，也只允许在 artifact root 内读写
- Flow 顶层 `browser` 里的相对路径也会解析到 artifact root 下
- `run_flow` 默认 `headless=true`
- 直接命令行执行 `go run . -flow ...` 仍保持本地使用的灵活性

### Flow 级浏览器配置

真实业务如果依赖登录态，推荐把浏览器会话配置放到 Flow 顶层，而不是散落在步骤里：

```yaml
schema_version: "1"
name: admin_orders
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
```

支持的 `browser` 字段：

- `headless`
- `use_session`
- `storage_state` / `storage_state_path` / `load_storage_state`
- `save_storage_state`
- `persistent`
- `profile`
- `session`
- `timeout`
- `user_agent`
- `viewport.width` / `viewport.height`

说明：

- `storage_state`、`storage_state_path`、`load_storage_state` 是同义入口，表示运行前加载登录态
- `use_session` 会从 `tsplay.save_session` 保存的命名会话里自动展开为 storage state 或 persistent profile 配置
- `save_storage_state` 会在 Flow 结束后把当前登录态保存回文件
- `profile` / `session` 会启用 persistent browser context，目录放在 artifact root 下
- `persistent profile/session` 目前不能和 `storage_state` 或 `use_session` 同时使用

### 显式会话工具

如果希望业务方只记“会话名”，可以先保存一个命名会话：

```json
{
  "name": "admin",
  "storage_state_path": "states/admin.json"
}
```

或者把 persistent profile 注册成命名会话：

```json
{
  "name": "crm-admin",
  "profile": "crm",
  "session": "admin"
}
```

之后调用 `tsplay.list_sessions` 会看到类似返回：

- `browser: { "use_session": "admin" }`
- `resolved_browser: { "storage_state": "sessions/storage/admin.json" }`
- `last_used_at: "2026-04-16T08:30:00Z"`
- `source: "copied from storage_state_path states/admin.json"`

如果想查看单个会话的展开结果，可以再调 `tsplay.get_session`：

```json
{
  "name": "admin"
}
```

返回里会包含：

- `expanded_browser`
- `physical_paths.metadata_path`
- `physical_paths.storage_state_path` 或 `physical_paths.profile_dir`

如果想直接拿可粘贴进 Flow 文件的片段，可以再调 `tsplay.export_session_flow_snippet`：

```json
{
  "name": "admin",
  "format": "browser"
}
```

返回里会包含：

- `export.format`
- `export.target`
- `export.encoding`
- `export.snippet`
- `export.snippet_data`

`format` 默认是 `all`，会把常用片段一次性都返回出来。也支持指定导出一种：

- `browser` / `browser_yaml`
- `expanded_browser` / `expanded_browser_yaml`
- `flow` / `flow_yaml`
- `expanded_flow` / `expanded_flow_yaml`
- `browser_json`
- `expanded_browser_json`
- `flow_json`
- `expanded_flow_json`

比如：

- 想只拿推荐的 `browser:` YAML，就传 `format: "browser"`
- 想只拿完整 Flow YAML，就传 `format: "flow"`
- 想给上层 AI 或程序直接消费 JSON，就传 `format: "flow_json"` 或 `expanded_flow_json`

这样后续 Flow 里只需要写：

```yaml
browser:
  use_session: admin
  headless: true
```

如果某个命名会话已经不用了，可以删除：

```json
{
  "name": "admin"
}
```

说明：

- `list_sessions` 会返回 `last_used_at`，只有 Flow 真正通过 `browser.use_session` 运行时才会更新
- `list_sessions` 也会返回 `source_type` 和 `source`，方便知道这个会话是从 JSON、文件还是 persistent profile 注册来的
- `delete_session` 会删除命名注册；如果这个会话是 storage-state 副本，也会删除对应的 session 文件
- `delete_session` 不会自动删除 persistent profile 目录，避免误删浏览器资料

## Action 速查

下面的 `模式` 列含义：

- `Lua / Flow`：Lua 脚本和 Flow 都可用
- `Lua`：仅 Lua CLI / Lua 脚本可用
- `Flow`：仅 Flow 步骤可用

### 导航与窗口

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `navigate` | 打开指定 URL | `url` | Lua / Flow |
| `reload` | 刷新当前页面 | - | Lua / Flow |
| `go_back` | 返回上一页 | - | Lua / Flow |
| `go_forward` | 前进到下一页 | - | Lua / Flow |
| `new_tab` | 新开标签页 | `url` | Lua / Flow |
| `close_tab` | 关闭当前标签页 | - | Lua / Flow |
| `switch_to_tab` | 切换标签页 | `index` | Lua / Flow |

### 输入与交互

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `click` | 点击元素 | `selector` | Lua / Flow |
| `type_text` | 输入文本 | `selector`, `text` | Lua / Flow |
| `get_text` | 读取元素文本 | `selector` | Lua / Flow |
| `extract_text` | 提取元素文本，可选等待并用正则提取首个匹配 | `selector`, `timeout`, `pattern` | Flow |
| `set_value` | 设置输入框值 | `selector`, `value` | Lua / Flow |
| `select_option` | 选择下拉项 | `selector`, `value` | Lua / Flow |
| `hover` | 悬停元素 | `selector` | Lua / Flow |
| `scroll_to` | 滚动到元素位置 | `selector` | Lua / Flow |
| `accept_alert` | 接受弹窗 | - | Lua / Flow |
| `dismiss_alert` | 取消弹窗 | - | Lua / Flow |
| `set_alert_text` | 给弹窗输入文本 | `text` | Lua / Flow |

### 等待与流程控制

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `wait_for_network_idle` | 等待网络空闲 | - | Lua / Flow |
| `wait_for_selector` | 等待元素出现 | `selector`, `timeout` | Lua / Flow |
| `wait_for_text` | 等待元素出现指定文本 | `selector`, `text`, `timeout` | Lua / Flow |
| `sleep` | 暂停执行 | `seconds` | Lua / Flow |
| `set_var` | 把值显式写入 Flow 变量 | `save_as`, `value` | Flow |
| `append_var` | 把值追加到列表变量 | `save_as`, `with.value` | Flow |
| `retry` | 重试一组嵌套步骤 | `times`, `interval_ms`, `steps` | Flow |
| `if` | 条件分支 | `condition`, `then`, `else` | Flow |
| `foreach` | 遍历列表并执行嵌套步骤，可选写 Redis checkpoint | `items`, `item_var`, `index_var`, `steps`, `with.progress_key?`, `with.progress_connection?`, `with.progress_value?` | Flow |
| `on_error` | 局部错误处理 | `steps`, `on_error` | Flow |
| `wait_until` | 轮询条件直到满足 | `condition`, `timeout`, `interval_ms` | Flow |
| `http_request` | 发起 HTTP 请求并返回结构化响应 | `url`, `with.method`, `with.headers`, `with.query`, `with.json`, `with.form`, `with.multipart_files`, `with.response_as` | Lua / Flow |
| `json_extract` | 用 JSON 路径提取字段 | `from`, `path` | Lua / Flow |
| `read_csv` | 读取本地 CSV 并返回行对象列表 | `file_path`, `with.start_row?`, `with.limit?`, `with.row_number_field?` | Lua / Flow |
| `read_excel` | 读取本地 Excel `.xlsx` 并返回行对象列表 | `file_path`, `sheet?`, `range?`, `with.headers?`, `with.start_row?`, `with.limit?`, `with.row_number_field?` | Lua / Flow |
| `write_json` | 把任意值写入本地 JSON 文件 | `file_path`, `with.value` | Lua / Flow |
| `write_csv` | 把行数据写入本地 CSV 文件 | `file_path`, `with.value`, `with.headers?` | Lua / Flow |
| `redis_get` | 从 Redis 读取一个键 | `key`, `connection?` | Lua / Flow |
| `redis_set` | 写入一个 Redis 键，可选 TTL | `key`, `value`, `ttl_seconds?`, `connection?` | Lua / Flow |
| `redis_del` | 删除一个 Redis 键 | `key`, `connection?` | Lua / Flow |
| `redis_incr` | 递增 Redis 计数 | `key`, `delta?`, `connection?` | Lua / Flow |
| `db_insert` | 把一条结构化结果写入数据库表 | `with.table`, `with.row`, `with.columns?`, `connection?`, `with.driver?` | Lua / Flow |
| `assert_visible` | 断言元素可见 | `selector`, `timeout` | Flow |
| `assert_text` | 断言文本匹配 | `selector`, `text`, `timeout` | Flow |
| `lua` | 在 Flow 中执行 Lua 代码 | `code` | Flow |

### 截图、文件与脚本执行

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `screenshot` | 整页截图 | `path` | Lua / Flow |
| `screenshot_element` | 元素截图 | `selector`, `path` | Lua / Flow |
| `save_html` | 保存页面 HTML | `path` | Lua / Flow |
| `execute_script` | 执行 JavaScript | `script` | Lua / Flow |
| `evaluate` | 对元素执行 JS 表达式 | `selector`, `script` | Lua / Flow |
| `upload_file` | 上传单个文件 | `selector`, `file_path` | Lua / Flow |
| `upload_multiple_files` | 上传多个文件 | `selector`, `files` | Lua / Flow |
| `download_file` | 通过页面元素触发下载 | `selector`, `save_path` | Lua / Flow |
| `download_url` | 直接下载 URL 到本地 | `url`, `save_path` | Lua / Flow |

### 数据提取与状态检查

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `get_attribute` | 获取元素属性 | `selector`, `attribute` | Lua / Flow |
| `get_html` | 获取元素或整页 HTML | `selector?` | Lua / Flow |
| `get_all_links` | 提取页面或区域中的链接 | `selector?` | Lua / Flow |
| `capture_table` | 提取表格数据 | `selector` | Lua / Flow |
| `find_element` | 获取单个元素信息 | `selector` | Lua / Flow |
| `find_elements` | 获取多个元素信息 | `selector` | Lua / Flow |
| `is_visible` | 判断元素是否可见 | `selector` | Lua / Flow |
| `is_enabled` | 判断元素是否可用 | `selector` | Lua / Flow |
| `is_checked` | 判断复选框/单选框是否勾选 | `selector` | Lua / Flow |
| `is_selected` | 判断下拉项是否被选中 | `selector` | Lua / Flow |
| `is_aria_selected` | 判断 ARIA selected 状态 | `selector` | Lua / Flow |

### 网络与浏览器状态

| Action | 说明 | 常用参数 | 模式 |
| --- | --- | --- | --- |
| `block_request` | 阻止指定网络请求 | `pattern` | Lua / Flow |
| `get_response` | 获取某个请求的响应 | `url` | Lua / Flow |
| `intercept_request` | 用 Lua 回调拦截请求 | `callback` | Lua |
| `get_storage_state` | 获取浏览器状态 | `context_index?` | Lua / Flow |
| `get_cookies_string` | 获取 Cookie 字符串 | `context_index?` | Lua / Flow |

说明：

- `intercept_request` 当前只在 Lua 中提供，Flow 没有对应的回调式写法
- `get_storage_state` 和 `get_cookies_string` 在 Flow 里支持可选的 `context_index`
- `http_request` 在 MCP 模式下需要显式传 `allow_http=true`
- `http_request` 如果要上传本地文件或保存响应到文件，还需要 `allow_file_access=true`
- `read_csv` 和 `read_excel` 在 MCP 模式下需要显式传 `allow_file_access=true`
- `write_json` 和 `write_csv` 在 MCP 模式下也需要显式传 `allow_file_access=true`
- `read_excel` 当前支持 `.xlsx`，不支持 `.xls`
- `read_excel.range` 支持 `A2:B20` 这类范围；数据区不含表头时可配 `with.headers`
- `read_csv` / `read_excel` 的 `with.start_row` 使用源文件真实行号，适合做断点续跑
- `redis_*` 在 MCP 模式下需要显式传 `allow_redis=true`
- `foreach.with.progress_key` 在 MCP 模式下同样需要显式传 `allow_redis=true`
- `db_insert` 在 MCP 模式下需要显式传 `allow_database=true`
- 新连接方式推荐用 `TSPLAY_DB_*` 或 `TSPLAY_DB_<NAME>_*`
- `db_insert` 默认内置 `mysql`、`postgres`
- `db_insert` 也支持 `sqlserver` 和 Oracle，但分别需要用 `-tags tsplay_sqlserver`、`-tags tsplay_oracle` 构建带对应驱动的二进制
- Oracle 的可选驱动使用 `github.com/sijms/go-ora/v2`，是 pure Go，不依赖本机 Oracle client
- `db_insert` 在 `driver=mysql` 时，仍兼容 `TSPLAY_MYSQL_*` 或 `TSPLAY_MYSQL_<NAME>_*` 连接配置

## 项目结构

```text
.
├── main.go
├── docs/             # 文档索引、培训体系、实验与讲师材料
├── tsplay_core/      # 核心执行引擎、Flow、MCP、页面观察与修复能力
├── script/           # Lua 与 Flow 示例
├── demo/             # 本地演示页面
├── tsplay_test/      # 测试与演示资源
└── mcp_test/         # MCP 相关试验代码
```

## 开发与测试

```bash
go test ./...
```
