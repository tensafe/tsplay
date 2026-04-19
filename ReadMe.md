# TSPlay

> 面向 AI Agent 和交付团队的浏览器自动化执行引擎。  
> 把 `observe -> draft -> validate -> run -> repair` 和会话复用收敛到同一条链路。

TSPlay 基于 Go + Playwright，提供 `Lua CLI / Lua Script`、`Flow DSL`、`MCP Server` 三层能力。

它不是单纯把浏览器动作封成一堆函数，而是把这些关键环节统一起来：

- 页面自动化执行
- 结构化 Flow 生成与校验
- 失败 trace 与现场留痕
- 登录态 / 会话复用
- Agent 调用与安全授权

如果你想做的是“可长期维护、可被 AI 生成、可被团队交付”的浏览器自动化，TSPlay 更接近这个方向。

## 适用场景

- Web RPA：登录、检索、点击、上传、下载、导出
- 页面数据提取：文本、属性、链接、表格、HTML、Cookie、Storage State
- 业务流程自动化：带变量、控制流、断言、失败恢复和断点续跑
- Agent Browser Tool：让 Codex、OpenClaw 等模型先观察页面，再草拟、执行和修复 Flow
- 浏览器 + 数据源联动：HTTP API、Redis、CSV/Excel、数据库

## 三层能力

| 能力层 | 适合场景 | 入口 |
| --- | --- | --- |
| `Lua CLI / Script` | 临时调试、手工探索、一次性任务 | `go run . -action cli` / `go run . -script ...` |
| `Flow DSL` | 版本化、可审查、可复用、可由 AI 生成的业务流程 | `go run . -flow ...` |
| `MCP Server` | 给 Agent 暴露观察、生成、执行、修复和会话能力 | `go run . -action srv` |

## 核心能力

- 基于 Playwright 驱动 Chromium
- 支持 Lua 直接控制浏览器
- 支持结构化 Flow YAML / JSON
- Flow 支持变量替换、`save_as`、`retry`、`if`、`foreach`、`on_error`、`wait_until`
- 支持 `http_request`、`json_extract`、`read_csv`、`read_excel`、`write_json`、`write_csv`
- 支持 `redis_get/set/del/incr` 和 `db_insert/db_query/db_transaction` 等数据动作
- 失败时自动落盘截图、HTML、DOM snapshot
- 支持从“用户意图 + 页面观察”自动草拟 Flow
- 支持为失败 Flow 生成统一 repair context 和 repair request
- 支持命名浏览器会话保存、复用和导出
- MCP 模式下带显式安全边界和能力授权

## 快速开始

### 环境要求

- Go `1.23.6+`
- 能运行 Playwright Chromium 的系统环境
- 首次执行需要浏览器的能力时，程序会自动执行 `playwright.Install()` 下载浏览器

### 1. 拉取依赖

```bash
go mod download
```

### 2. 选一种方式跑起来

| 想做什么 | 命令 |
| --- | --- |
| 启动交互式 CLI | `go run . -action cli` |
| 运行 Lua 脚本 | `go run . -script script/open_url.lua` |
| 运行 Flow | `go run . -flow script/demo_baidu.flow.yaml` |
| 启动 MCP Server | `go run . -action srv` |

如果想隐藏浏览器窗口，可以追加 `-headless`。

### 3. CLI 试跑

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

### 4. 运行第一个 Flow

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

执行：

```bash
go run . -flow script/demo_baidu.flow.yaml
```

Flow 执行后会输出结构化 JSON，包含变量、trace、耗时和失败现场路径。

## 为什么 Flow 是主线

相比裸 Lua，Flow 更适合做长期维护的业务资产：

- 更容易被 AI 严格生成
- 更容易做人工 review 和版本 diff
- 更容易记录失败上下文和 trace
- 更容易复用登录态、数据源和控制流
- 更适合通过 MCP 暴露给 Agent

常见 Flow 能力包括：

- 变量：`vars`、`save_as`、`set_var`、`append_var`
- 控制流：`retry`、`if`、`foreach`、`on_error`、`wait_until`
- 页面动作：点击、输入、等待、断言、截图、上传、下载
- 数据动作：HTTP、Redis、CSV、Excel、数据库
- 浏览器状态：`use_session`、`storage_state`、`save_storage_state`

在 MCP 模式下，推荐通过这些工具获取完整约束而不是手写猜测：

- `tsplay.list_actions`
- `tsplay.flow_schema`
- `tsplay.flow_examples`

## 执行结果与失败现场

每次运行 Flow 都会返回 step trace，通常包含：

- `action`
- 参数摘要
- `status`
- `duration_ms`
- 输出摘要
- 当前 `page_url`

当步骤失败时，TSPlay 会把现场资料写到 artifact root，默认目录是 `artifacts/`。常见产物包括：

- `failure.png`
- `page.html`
- `dom_snapshot.json`

可以通过命令行参数调整输出目录：

```bash
go run . -flow script/demo_baidu.flow.yaml -artifact-root artifacts
```

## MCP / Agent 集成

TSPlay 可以作为 MCP Server 启动，让 Agent 不必直接读整页 HTML，也不必手写 selector。

如果你想先从“用户通过大模型驱动 TSPlay”这条路线入门，建议配合阅读：
[docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md)

### 启动方式

```bash
go run . -action srv
go run . -action srv -addr :8081
go run . -action srv -flow-root script -artifact-root artifacts
go run . -action mcp-stdio -flow-root script -artifact-root artifacts
```

默认约束：

- `flow_path` 只允许读取 `script/` 下的文件，可用 `-flow-root` 调整
- 文件类输入输出默认限制在 artifact root 下
- `run_flow` 默认 `headless=true`

### MCP 工具分组

| 分组 | 工具 |
| --- | --- |
| Flow 认知 | `tsplay.list_actions`、`tsplay.flow_schema`、`tsplay.flow_examples` |
| 页面观察与草拟 | `tsplay.observe_page`、`tsplay.draft_flow`、`tsplay.finalize_flow` |
| 校验、执行与修复 | `tsplay.validate_flow`、`tsplay.run_flow`、`tsplay.repair_flow_context`、`tsplay.repair_flow` |
| 会话管理 | `tsplay.save_session`、`tsplay.list_sessions`、`tsplay.get_session`、`tsplay.export_session_flow_snippet`、`tsplay.delete_session` |

### 推荐调用顺序

如果你想给小模型或产品化接入一条更短的默认路径，优先用：

1. `tsplay.finalize_flow`
2. 如果 `status=ready`，直接 `tsplay.run_flow`
3. 如果 `status=needs_permission`，补授权后再次 `tsplay.finalize_flow`
4. 如果 `status=needs_input`，补输入后再次 `tsplay.finalize_flow`
5. 如果 `status=needs_repair`，再进入 `tsplay.validate_flow` / `tsplay.repair_flow_context` / `tsplay.repair_flow`

如果你需要更细粒度控制，再走完整链路：

1. 先调 `tsplay.flow_schema`，拿到严格约束
2. 需要参考模板时调 `tsplay.flow_examples`
3. 有明确 URL 和意图时优先调 `tsplay.draft_flow`
4. 需要更细粒度观察时先调 `tsplay.observe_page`
5. 如需单独校验，再调 `tsplay.validate_flow`
6. 成功后用 `tsplay.run_flow` 执行
7. 失败后用 `tsplay.repair_flow_context` / `tsplay.repair_flow` 收敛修复
8. 如果流程要长期复用登录态，再调 `tsplay.save_session`

`tsplay.draft_flow`、`tsplay.finalize_flow` 和 `tsplay.repair_flow_context` 都会返回统一结构的 `repair_hints` 或修复线索，方便继续交给模型修正。

黄金路径工具现在都会返回统一 envelope，顶层至少包含这些字段：

- `ok`
- `tool`
- `summary`
- `artifacts`
- `next_action`
- `warnings`
- `run`

其中浏览器类工具会稳定返回 `run.id`、`status`、`queue_wait_ms`、`duration_ms`、`timeout_ms`、`audit_path`、`run_root` 等运行元信息。

另外，和 Flow 结构相关的工具现在会尽量返回结构化诊断：

- `draft.validation.issue`
- `validate_flow.issue`
- `finalize_flow.issue`

这些字段适合给模型直接消费，尤其是：

- 未知字段
- 不支持的 action
- 参数名写错
- 安全策略阻塞

例如常见的 `fill -> type_text`、`result_var -> save_as`、`with.headers` 写法错误、`navigate.timeout` 误用，都会尽量返回更明确的修复提示。

### 从用户意图草拟 Flow

`tsplay.draft_flow` 适合把“用户想做什么”先落成一份可审阅草稿：

```json
{
  "intent": "搜索订单并导出",
  "url": "https://example.com/orders"
}
```

如果你已经先调过 `tsplay.observe_page`，也可以直接把返回的 `observation` 继续传给 `draft_flow`，减少模型自行猜页面结构的空间。

它会尽量完成这些事情：

- 自动打开页面并观察可交互元素
- 匹配搜索、导出、上传、登录、下拉选择等常见动作
- 生成结构化 Flow YAML
- 自动执行一轮与 `validate_flow` 对齐的校验
- 在 observation 里存在更稳 selector 时自动修正一轮
- 返回 `validation`、`selector_repairs`、`repair_hints`、建议变量和未解决项

### 一步式收敛：`tsplay.finalize_flow`

`tsplay.finalize_flow` 适合做“小模型友好”的默认入口：它会复用 `draft_flow` 的观察和草拟能力，但直接返回是否已经可以执行。

典型输入：

```json
{
  "intent": "搜索订单并导出",
  "url": "https://example.com/orders",
  "security_preset": "readonly"
}
```

重点看这些字段：

- `status`：`ready` / `needs_input` / `needs_permission` / `needs_repair`
- `flow_yaml`
- `validation`
- `issue`
- `suggested_vars`
- `next_action`

## 浏览器会话与 Flow 级配置

如果业务流程依赖登录态，推荐把浏览器配置放到 Flow 顶层，而不是散落在步骤里：

```yaml
schema_version: "1"
name: admin_orders
browser:
  headless: true
  use_session: admin
  save_storage_state: states/admin-latest.json
  timeout: 30000
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

常用 `browser` 字段：

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

- `use_session` 会从 `tsplay.save_session` 保存的命名会话里自动展开
- `save_storage_state` 会在 Flow 结束后保存当前登录态
- `profile` / `session` 会启用 persistent browser context
- `persistent profile/session` 不能和 `storage_state` 或 `use_session` 同时使用

如果希望业务方只记一个会话名，可以先保存：

```json
{
  "name": "admin",
  "storage_state_path": "states/admin.json"
}
```

后续 Flow 里直接写：

```yaml
browser:
  use_session: admin
```

## 安全边界

MCP 模式默认不是全放开。高风险能力需要按请求显式授权：

| 授权参数 | 放行动作 |
| --- | --- |
| `allow_lua=true` | `lua` |
| `allow_javascript=true` | `execute_script`、`evaluate` |
| `allow_file_access=true` | `screenshot`、`save_html`、`read_csv`、`read_excel`、上传下载、`write_json`、`write_csv` |
| `allow_browser_state=true` | Cookie / Storage State / `browser.use_session` / persistent profile |
| `allow_http=true` | `http_request` |
| `allow_redis=true` | `redis_get`、`redis_set`、`redis_del`、`redis_incr`、`foreach.with.progress_key` |
| `allow_database=true` | `db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction` |

如果不想逐个拼 `allow_*`，也可以用 `security_preset`：

- `readonly`：默认最小权限
- `browser_write`：开启文件读写和浏览器状态能力，适合上传、下载、截图、Storage State 复用
- `full_automation`：开启全部 MCP 安全能力

显式传入的 `allow_*` 会覆盖 `security_preset` 中对应字段。

补充说明：

- 文件类动作即使被授权，也只能在 artifact root 范围内读写
- Flow 顶层 `browser` 里的相对路径也会解析到 artifact root 下
- 本地命令行运行 `go run . -flow ...` 仍保持更灵活的本地使用方式

## 外部系统集成

除了页面动作，TSPlay 也支持把浏览器流程和数据动作放进同一条 Flow 里。

### HTTP

可直接用 `http_request` 调外部 API，再用 `json_extract` 继续编排，适合：

- OCR 验证码识别
- 内部查单 / 补数接口
- webhook / 通知接口

### Redis

适合共享 cookie、游标、去重标记、断点续跑 checkpoint。

环境变量约定：

- 默认连接：`TSPLAY_REDIS_*`
- 命名连接：`TSPLAY_REDIS_<NAME>_*`
- 也支持 URL 形式：`TSPLAY_REDIS_URL`、`TSPLAY_REDIS_<NAME>_URL`

### CSV / Excel

适合批量导入、分段执行、结果台账回写。

- `read_csv` 默认把第一行非空行当表头
- `read_excel` 当前支持 `.xlsx`
- `read_excel.range` 支持 `A2:B20` 这类矩形范围
- 可配 `with.start_row`、`with.limit`、`with.row_number_field` 做断点续跑

### 数据库

适合把结构化结果直接落表，或在流程中查询业务数据。

环境变量约定：

- 默认连接：`TSPLAY_DB_*`
- 命名连接：`TSPLAY_DB_<NAME>_*`

支持的常见 driver：

- `mysql`
- `pgsql`
- `sqlserver`
- `oracle`

说明：

- `db_*` 动作在 MCP 模式下需要 `allow_database=true`
- SQL Server / Oracle 需要带对应 build tags 构建带驱动的二进制

## 文档入口

README 负责项目定位和快速上手，更多培训、落地和 Enablement 材料放在 `docs/`。

| 内容 | 说明 | 入口 |
| --- | --- | --- |
| 文档索引 | 仓库文档地图和推荐阅读顺序 | [docs/README.md](docs/README.md) |
| 培训体系总览 | 面向实施、测试、开发和讲师的统一入口 | [docs/training/README.md](docs/training/README.md) |
| 学习路径 | 从新手到 MCP Integrator / Trainer 的路线图 | [docs/training/learning-path.md](docs/training/learning-path.md) |
| Bootcamp 课程表 | 2 天训练营和 4 周落地节奏 | [docs/training/bootcamp-plan.md](docs/training/bootcamp-plan.md) |
| 实训实验 | 基于 `demo/` 和 `script/` 的练习清单 | [docs/training/labs.md](docs/training/labs.md) |
| 考核与认证 | 评分维度、证据标准和结业要求 | [docs/training/assessment.md](docs/training/assessment.md) |
| 讲师手册 | 讲师备课、授课和复盘指南 | [docs/training/trainer-playbook.md](docs/training/trainer-playbook.md) |

## 项目结构

```text
.
├── main.go
├── docs/             # 文档索引、培训体系、实验与讲师材料
├── tsplay_core/      # 核心执行引擎、Flow、MCP、观察与修复能力
├── script/           # Lua 与 Flow 示例
├── demo/             # 本地演示页面
├── tsplay_test/      # 测试与演示资源
└── mcp_test/         # MCP 相关试验代码
```

## 开发与测试

```bash
go test ./...
```

如果你是第一次接触 TSPlay，推荐按这个顺序阅读：

1. 先看本页，理解三层能力和快速开始
2. 再看 [docs/README.md](docs/README.md)，定位后续材料
3. 想学 Flow 交付时看 [docs/training/learning-path.md](docs/training/learning-path.md)
4. 想接 Agent / MCP 时重点看本页的 MCP 章节和 `tsplay_core/tsplay_mcpserver.go`
