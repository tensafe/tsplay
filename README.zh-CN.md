# TSPlay

[English](ReadMe.md) | 简体中文

> 面向 AI Agent 和交付团队的浏览器自动化执行引擎。  
> 用一条统一链路收敛 `observe -> draft -> validate -> run -> repair`、会话复用和安全授权。

TSPlay 基于 Go + Playwright，提供三层入口：`Lua CLI / Script`、`Flow DSL`、`MCP Server`。

它不是把浏览器动作简单封成一组函数，而是把这些在真实交付里最容易散掉的能力组织到一起：

- 页面观察与可交互元素提取
- 结构化 Flow 草拟、校验和执行
- 失败 trace、截图、HTML、DOM snapshot 留痕
- 登录态 / 浏览器会话复用
- 面向 Agent 的 MCP 工具和安全授权边界

如果你想做的是“可长期维护、可被 AI 生成、可被团队审查和交付”的浏览器自动化，TSPlay 更接近这个方向。

## 适用场景

- Web RPA：登录、检索、点击、上传、下载、导出
- 页面数据提取：文本、属性、链接、表格、HTML、Cookie、Storage State
- 业务流程自动化：变量、控制流、断言、失败恢复、断点续跑
- Agent Browser Tool：让 Codex、OpenClaw 等模型先观察页面，再草拟、执行和修复 Flow
- 浏览器和外部系统联动：HTTP API、Redis、CSV/Excel、数据库

## 为什么 TSPlay 不是“又一个 Playwright 封装”

- `Flow` 可版本化、可评审、可复用，也更适合 AI 严格生成
- 失败时自动沉淀 trace、截图、HTML 和 DOM snapshot，方便排障和修复
- 会话可以命名保存和复用，适合长期交付而不是一次性脚本
- MCP 模式自带安全边界，适合把能力暴露给 Agent，而不是直接放开浏览器

## 三种使用方式

| 方式 | 适合什么情况 | 入口 |
| --- | --- | --- |
| `Lua CLI / Script` | 临时调试、页面探索、一次性任务 | `go run . -action cli` / `go run . -script ...` |
| `Flow DSL` | 版本化、可审查、可复用、可由 AI 生成的业务流程 | `go run . -flow ...` |
| `MCP Server` | 给 Agent 暴露观察、生成、执行、修复和会话能力 | `go run . -action srv` |

日常交付推荐以 `Flow` 为主线。  
CLI 适合探索页面，MCP 适合接入 AI 产品或 Agent 工作流。

## 能力对照矩阵

这张表的重点不是追求所有能力都机械地一比一复制，而是区分哪些能力应该强同步，哪些更适合留在 `Flow` 这一层。

| 能力类别 | 典型能力 | Flow | Lua | MCP | 建议 |
| --- | --- | --- | --- | --- | --- |
| 页面原子动作 | `navigate`、`click`、`type_text`、`select_option` | 是 | 是 | 是 | 应保持同步 |
| 文件与表格 I/O | `screenshot`、`save_html`、`read_csv`、`read_excel`、`write_json`、`write_csv` | 是 | 是 | 是 | 应保持同步，MCP 下受 `allow_file_access` 约束 |
| HTTP 请求 | `http_request`、`json_extract` | 是 | 是 | 是 | 应保持同步；Lua 在 Flow / MCP 安全上下文中也遵守 `allow_http`、`allow_file_access` 和文件根目录 |
| Redis 操作 | `redis_get`、`redis_set`、`redis_del`、`redis_incr` | 是 | 是 | 是 | 应保持同步；Lua 在 Flow / MCP 安全上下文中也遵守 `allow_redis` |
| 数据库操作 | `db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction` | 是 | 是 | 是 | 应保持同步；Lua 在 Flow / MCP 安全上下文中也遵守 `allow_database`，`db_transaction` 会自动提交或回滚 |
| 浏览器状态 | `get_storage_state`、`get_cookies_string`、`browser.use_session` | 是 | 是 | 是 | 应保持同步，MCP 下受 `allow_browser_state` 约束 |
| Flow 便捷动作 | `extract_text`、`assert_visible`、`assert_text`、`set_var`、`append_var` | 是 | 是 | 是 | 已对齐；更适合作为编排语义糖而不是底层原语 |
| Flow 控制流 | `retry`、`if`、`foreach`、`on_error`、`wait_until` | 是 | 否 | 是 | 不要求同步到 Lua |
| Lua 回调型能力 | `intercept_request` | 否 | 是 | 否 | 保持 Lua 专属更自然 |

推荐的判断原则：

- 像 `HTTP / Redis / 数据库 / 文件读写` 这种“原子数据动作”，最好在 `Flow` 和 `Lua` 两边都可用，这样探索、固化、接入三条路径不会断层。
- 像 `retry / foreach / on_error / wait_until` 这种“编排能力”，更适合保留在 `Flow DSL`，不需要硬翻成 Lua 扩展函数。
- 像 `extract_text / assert_text / assert_visible` 这种“语义增强动作”，可以先作为 `Flow` 的高层封装；如果 Lua 侧频繁手搓同样逻辑，再补一层语法糖。

## 快速开始

### 环境要求

- Go `1.23.6+`
- 能运行 Playwright Chromium 的系统环境
- 首次执行浏览器相关能力时，程序会自动执行 `playwright.Install()` 下载浏览器

### 安装依赖

```bash
go mod download
```

### 选一种方式跑起来

| 想做什么 | 命令 |
| --- | --- |
| 启动交互式 CLI | `go run . -action cli` |
| 运行 Lua 脚本 | `go run . -script script/open_url.lua` |
| 运行 Flow | `go run . -flow script/demo_baidu.flow.yaml` |
| 启动内置静态文件服务 | `go run . -action file-srv -addr :8000` |
| 直接调用一个 TSPlay MCP 工具 | `go run . -action mcp-tool -tool tsplay.list_actions` |
| 列出 macOS 录屏设备 | `go run . -action list-record-devices` |
| 录整个桌面屏幕 | `go run . -action record-screen -record-cmd "go run . -flow script/tutorials/10_assert_page_state.flow.yaml"` |
| 只录浏览器页面内容 | `go run . -flow script/tutorials/10_assert_page_state.flow.yaml -browser-video-output artifacts/recordings/lesson-10-assert-page-state.webm` |
| 列出二进制内置资源 | `go run . -action list-assets` |
| 释放内置 docs/script/demo | `go run . -action extract-assets -extract-root ./tsplay-assets` |
| 启动 MCP Server | `go run . -action srv` |

如果想隐藏浏览器窗口，可以追加 `-headless`。

`record-screen` 录的是整个 macOS 屏幕，适合录桌面演示。  
`-browser-video-output` 录的是 Playwright 页面内容，更适合浏览器教程。  
默认还会额外保留一小段页面停留时间，避免视频短到来不及看。  
更完整的讲师用法见
[docs/training/tutorial-video-recording.md](docs/training/tutorial-video-recording.md)。

当你构建 `./tsplay` 二进制后，`ReadMe.md`、`docs/`、`script/`、`demo/` 会一并打包进二进制：

- 可以直接运行内置示例：`./tsplay -script script/tutorials/01_hello_world.lua`
- 可以直接运行内置 Flow：`./tsplay -flow script/tutorials/01_hello_world.flow.yaml`
- 可以直接服务内置 demo：`./tsplay -action file-srv -addr :8000`
- 也可以把参考资料释放到本地目录：`./tsplay -action extract-assets -extract-root ./tsplay-assets`

### 先跑一个 Flow

仓库里已经带了一个最小示例：

```bash
go run . -flow script/demo_baidu.flow.yaml
```

对应 Flow 大致长这样：

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

执行后会输出结构化 JSON，包含变量、step trace、耗时和失败现场路径。

### 用 CLI 探路

先启动：

```bash
go run . -action cli
```

进入后先输入：

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

### 启动 MCP Server

```bash
go run . -action srv
go run . -action srv -addr :8081
go run . -action srv -flow-root script -artifact-root artifacts
go run . -action mcp-stdio -flow-root script -artifact-root artifacts
go run . -action mcp-tool -tool tsplay.list_actions
go run . -action mcp-tool -tool tsplay.observe_page -args-file script/tutorials/113_mcp_observe_page_template_release.args.json
```

默认约束：

- `flow_path` 默认只允许读取 `script/` 下的文件，可用 `-flow-root` 调整
- 文件类输入输出默认限制在 artifact root 下
- `run_flow` 默认 `headless=true`

如果你想从“用户描述意图，模型帮忙草拟和执行 Flow”这条路径入门，建议配合阅读：
[docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md)

## 为什么以 Flow 为主线

相比裸 Lua，Flow 更适合做长期维护的业务资产：

- 更容易被 AI 严格生成
- 更容易做人工 review 和版本 diff
- 更容易做 schema 校验和结构化 issue 提示
- 更容易记录失败上下文和 repair 线索
- 更适合通过 MCP 暴露给 Agent

常见 Flow 能力包括：

- 变量：`vars`、`save_as`、`set_var`、`append_var`
- 控制流：`retry`、`if`、`foreach`、`on_error`、`wait_until`
- 页面动作：点击、输入、等待、断言、截图、上传、下载
- 数据动作：`http_request`、`json_extract`、`read_csv`、`read_excel`、`write_json`、`write_csv`
- 浏览器状态：`use_session`、`storage_state`、`save_storage_state`

## 核心能力

- 基于 Playwright 驱动 Chromium
- 支持 Lua 直接控制浏览器
- 支持结构化 Flow YAML / JSON
- 支持页面观察、Flow 草拟、执行校验和失败修复
- 支持命名浏览器会话保存、复用和导出
- 支持 `redis_get/set/del/incr` 和 `db_insert/db_query/db_transaction` 等数据动作
- 支持失败时自动落盘截图、HTML、DOM snapshot
- 支持通过 MCP 提供显式安全边界和能力授权

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

如果只是想给模型一条更短的默认路径，优先用 `tsplay.finalize_flow`。  
如果需要更细粒度控制，再走完整链路：`observe_page -> draft_flow -> validate_flow -> run_flow -> repair_flow_context -> repair_flow`。

`tsplay.finalize_flow` 的常见状态：

- `ready`：可以直接执行
- `needs_input`：还缺变量或用户输入
- `needs_permission`：命中了安全边界，需要补授权
- `needs_repair`：Flow 已成形，但还需要先修正

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

1. `tsplay.flow_schema`
2. `tsplay.flow_examples`
3. `tsplay.observe_page`
4. `tsplay.draft_flow`
5. `tsplay.validate_flow`
6. `tsplay.run_flow`
7. `tsplay.repair_flow_context` / `tsplay.repair_flow`
8. `tsplay.save_session`

黄金路径工具会尽量返回统一 envelope，顶层通常包含：

- `ok`
- `tool`
- `summary`
- `artifacts`
- `next_action`
- `warnings`
- `run`

## Flow 编写小贴士

- 用 `type_text`，不要写成 `fill`
- 用 `save_as`，不要写成 `result_var`
- 需要文件读写、上传、下载、截图时，在 MCP 模式下加 `allow_file_access=true`，或使用 `security_preset=browser_write`
- 页面级超时优先放在 `browser.timeout`，不要给 `navigate` 单独加不支持的 `timeout`
- 不确定 action 名时，先查 `tsplay.list_actions` 和 `tsplay.flow_schema`

## 浏览器会话与 Flow 顶层配置

如果业务流程依赖登录态，推荐把浏览器配置放在 Flow 顶层，而不是散落在步骤里：

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

MCP 模式默认不是全放开。高风险能力需要按请求显式授权。

### `security_preset`

- `readonly`：默认最小权限
- `browser_write`：开启文件读写和浏览器状态能力，适合上传、下载、截图、Storage State 复用
- `full_automation`：开启全部 MCP 安全能力

显式传入的 `allow_*` 会覆盖 `security_preset` 中对应字段。

### 常见授权参数

| 授权参数 | 放行动作 |
| --- | --- |
| `allow_lua=true` | `lua` |
| `allow_javascript=true` | `execute_script`、`evaluate` |
| `allow_file_access=true` | `screenshot`、`save_html`、`read_csv`、`read_excel`、上传下载、`write_json`、`write_csv` |
| `allow_browser_state=true` | Cookie / Storage State / `browser.use_session` / persistent profile |
| `allow_http=true` | `http_request` |
| `allow_redis=true` | `redis_get`、`redis_set`、`redis_del`、`redis_incr`、`foreach.with.progress_key` |
| `allow_database=true` | `db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction` |

补充说明：

- 文件类动作即使被授权，也只能在 artifact root 范围内读写
- Flow 顶层 `browser` 里的相对路径也会解析到 artifact root 下
- Lua 里的 `http_request`、`redis_*`、`db_*` 在 Flow / MCP 安全上下文里也会继承对应的 `allow_*` 约束
- 本地命令行运行 `go run . -flow ...` 仍保持更灵活的本地使用方式

## 外部系统集成

除了页面动作，TSPlay 也支持把浏览器流程和数据动作放进同一条 Flow 里。

### HTTP

可直接用 `http_request` 调外部 API，再用 `json_extract` 继续编排，适合：

- OCR 验证码识别
- 内部查单 / 补数接口
- webhook / 通知接口

补充说明：

- `Flow` 和 `Lua` 两边都支持 `http_request`
- 当 `Lua http_request` 运行在 `Flow` / MCP 安全上下文中时，也会遵守 `allow_http`
- 如果 `http_request` 使用了 `save_path` 或 `multipart_files`，Lua 侧也会和 Flow 一样遵守 `allow_file_access`
- 在受限运行模式下，`save_path` 和 `multipart_files` 的相对路径会解析到配置的文件根目录下

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

- `Flow` 和 `Lua` 两边都支持 `db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction`
- 当 `Lua db_*` 或 `Lua db_transaction` 运行在 `Flow` / MCP 安全上下文中时，也会遵守 `allow_database=true`
- `db_transaction` 会在同一个事务作用域里执行内部数据库操作，成功时自动 commit，失败时自动 rollback
- `db_*` 动作在 MCP 模式下需要 `allow_database=true`
- SQL Server / Oracle 需要带对应 build tags 构建带驱动的二进制

## 文档入口

README 负责项目定位和快速上手，更偏培训、落地和 Enablement 的材料放在 `docs/`。

| 内容 | 说明 | 入口 |
| --- | --- | --- |
| 文档索引 | 仓库文档地图和推荐阅读顺序 | [docs/README.md](docs/README.md) |
| 培训体系总览 | 面向实施、测试、开发和讲师的统一入口 | [docs/training/README.md](docs/training/README.md) |
| AI 无感入门 | 面向 Agent 的“用户意图 -> MCP -> Flow -> 执行修复”实战教程 | [docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md) |
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
4. 想接 Agent / MCP 时重点看 [docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md)
