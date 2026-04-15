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
- Flow 支持变量替换、`save_as`、`extract_text`、`set_var`、`retry`、`if`、`foreach`、`on_error`、`wait_until`、断言和失败 trace
- 失败时自动落盘现场资料：截图、HTML、DOM snapshot
- 支持基于“用户意图 + 页面观察”自动草拟 Flow
- 可作为 MCP Server 暴露给 Agent 调用
- MCP 模式带安全边界，可按能力显式授权

## 运行模式

| 模式 | 入口 | 适合场景 |
| --- | --- | --- |
| 交互式 CLI | `go run . -action cli` | 手动调试、边试边写 |
| Lua 脚本 | `go run . -script script/open_url.lua` | 自定义逻辑、一次性任务 |
| Flow DSL | `go run . -flow script/demo_baidu.flow.yaml` | 结构化流程、版本管理、AI 生成 |
| MCP Server | `go run . -action srv` | 接入 OpenClaw / Codex / 其他 Agent |

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
| `tsplay.observe_page` | 打开页面并返回截图路径、DOM snapshot、可交互元素和候选 selector |
| `tsplay.repair_flow_context` | 根据失败 Flow 和 run result 组织修复上下文，附带失败分类、修复焦点和校验清单 |
| `tsplay.validate_flow` | 只校验 Flow，不启动浏览器 |
| `tsplay.run_flow` | 启动 Playwright 执行 Flow，并返回 trace |

### 推荐给 Agent 的调用顺序

1. 调 `tsplay.flow_schema`，拿到严格约束
2. 调 `tsplay.flow_examples`，拿到参考模板
3. 如果已经有明确 URL 和用户意图，优先调 `tsplay.draft_flow`
4. 如果需要更细粒度控制，也可以先调 `tsplay.observe_page` 再把 observation 传给 `tsplay.draft_flow`
5. 查看 `tsplay.draft_flow` 返回里的 `validation`、`selector_repairs` 和 `repair_hints`
6. 如需单独二次校验，再调 `tsplay.validate_flow`
7. 成功后再调用 `tsplay.run_flow`
8. 失败时把原 Flow 和 run result 交给 `tsplay.repair_flow_context`

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

如果已经提前调过 `tsplay.observe_page`，也可以把 observation JSON 直接传给 `tsplay.draft_flow`，避免重复打开页面。

### 安全边界

MCP 模式默认不是全放开。高风险能力需要在单次请求里显式授权：

| 授权参数 | 放行动作 |
| --- | --- |
| `allow_lua=true` | `lua` |
| `allow_javascript=true` | `execute_script`、`evaluate` |
| `allow_file_access=true` | `screenshot`、`screenshot_element`、`save_html`、`upload_file`、`upload_multiple_files`、`download_file`、`download_url` |
| `allow_browser_state=true` | `get_storage_state`、`get_cookies_string` |

补充说明：

- `flow_path` 默认只允许读取 `script/` 目录内的文件，可用 `-flow-root` 调整
- 文件类动作即使被授权，也只允许在 artifact root 内读写
- `run_flow` 默认 `headless=true`
- 直接命令行执行 `go run . -flow ...` 仍保持本地使用的灵活性

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
| `retry` | 重试一组嵌套步骤 | `times`, `interval_ms`, `steps` | Flow |
| `if` | 条件分支 | `condition`, `then`, `else` | Flow |
| `foreach` | 遍历列表并执行嵌套步骤 | `items`, `item_var`, `index_var`, `steps` | Flow |
| `on_error` | 局部错误处理 | `steps`, `on_error` | Flow |
| `wait_until` | 轮询条件直到满足 | `condition`, `timeout`, `interval_ms` | Flow |
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

## 项目结构

```text
.
├── main.go
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
