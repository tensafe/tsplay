# TSPlay Flow 中文速查

## 这份文件适合什么时候看

当用户主要用中文描述需求，或者希望直接用中文把业务需求整理成 TSPlay Flow 时，先读这份文件。

如果你已经知道具体 action 怎么写，再去看：

- `actions.md`
- `examples.md`
- `example-index.md`

如果你已经进入中文排错或中文业务套模板阶段，再去看：

- `zh-cn-troubleshooting.md`
- `zh-cn-business-templates.md`
- `zh-cn-selectors.md`
- `zh-cn-review-checklist.md`

## 中文触发词

这些说法都应该触发这份 skill：

- 帮我写一条 TSPlay Flow
- 把这个需求转成 Flow
- 帮我改一下 flow
- 帮我修这条 flow
- 帮我补一条教程 flow
- 帮我做一个导入 flow
- 帮我做一个会话复用 flow
- 帮我排查这个 flow 为什么跑不通
- 帮我把 selector 修好
- 帮我按 MCP 思路把需求收敛成 Flow

## 中文提问最小模板

优先把需求整理成下面 5 个字段：

```text
- 页面: <URL 或本地页面>
- 目标: <要完成的业务动作>
- 输入: <关键词 / 文件 / 行号 / 条件，没有就写无>
- 输出: <save_as / JSON / CSV / Excel / artifact 路径>
- 授权: <readonly / browser_write / full_automation / allow_*>
```

这是最推荐的中文输入形状，因为它最容易直接落成 Flow。

## 常见中文需求到 Flow 动作的映射

- 打开页面: `navigate`
- 等页面加载到某个元素出现: `wait_for_selector`
- 输入搜索词或表单字段: `type_text`
- 点击搜索、提交、导出: `click`
- 判断页面是不是成功了: `assert_visible`, `assert_text`
- 提取标题、计数、状态文本: `extract_text`
- 抓表格: `capture_table`
- 接管真实 Chrome/Chromium/Edge: 顶层 `browser.cdp_launch`, `browser.cdp_port`, `browser.cdp_endpoint`
- 保存一个对象变量: `set_var`
- 累积结果列表: `append_var`
- 遍历 CSV 或 Excel 多行: `foreach`
- 某一行失败但整体继续: `on_error`
- 页面状态要轮询: `wait_until`
- 页面易抖动、要重试: `retry`
- 读本地 JSON / CSV / Excel: `read_json`, `read_csv`, `read_excel`
- 写结果到 JSON 或 CSV: `write_json`, `write_csv`
- 压缩或解压 ZIP: `zip_compress`, `zip_extract`
- 发邮件通知: `send_email`
- 复用登录态: 顶层 `browser.use_session`

## 中文写 Flow 的推荐顺序

1. 先明确页面和目标。
2. 再明确输入和输出。
3. 选最小可行 action 组合，不要一开始就写太复杂。
4. 如果页面有登录态，优先考虑顶层 `browser.use_session`。
5. 如果用户要复用真实 Chrome、扩展、缓存或人工登录状态，优先考虑 `browser.cdp_launch: true`；已经有远程调试端口时再用 `browser.cdp_port` / `browser.cdp_endpoint`。
6. 如果是批量处理，优先考虑 `read_csv` 或 `read_excel` 加 `foreach`。
7. 如果某个步骤容易失败但不该拖垮整条 Flow，优先考虑 `on_error`。
8. 如果用户不知道 selector，优先考虑 MCP 的 `observe_page` 路线。

## 中文场景起手建议

### 表单填写或搜索

先看：

- `example-index.md` 里的 Form Flows
- `actions.md` 里的 `navigate`, `wait_for_selector`, `type_text`, `click`, `assert_text`

### 表格抓取

先看：

- `example-index.md` 里的 Table Flows
- `actions.md` 里的 `capture_table`

### CSV 或 Excel 导入

先看：

- `example-index.md` 里的 Import Flows
- `actions.md` 里的 `read_csv`, `read_excel`, `foreach`, `append_var`

### 本地 JSON 输入

先看：

- `example-index.md` 里的 JSON And Artifact Input Flows
- `actions.md` 里的 `read_json`
- `zh-cn-business-templates.md` 里的本地 JSON 模板

### 登录态复用

先看：

- `example-index.md` 里的 Session Flows
- `actions.md` 里的顶层 `browser.use_session`

### 真实浏览器 / CDP 接管

先看：

- `actions.md` 里的 `browser.cdp_launch`, `browser.cdp_port`, `browser.cdp_endpoint`
- 如果当前 Chrome 没有监听 `9222`，不要强行重启用户日常窗口；用 `browser.cdp_launch: true` 启动独立 profile
- MCP 下记得 `allow_browser_state=true`
- 如果要写 JSON / CSV / 截图，再加 `allow_file_access=true`
- 如果要用 `execute_script` 或 `evaluate`，再加 `allow_javascript=true`

### 邮件通知

先看：

- `example-index.md` 里的 Email And Notification Flows
- `actions.md` 里的 `send_email`
- `zh-cn-business-templates.md` 里的邮件通知模板

### 容错与恢复

先看：

- `example-index.md` 里的 Recovery Flows
- `actions.md` 里的 `on_error`, `retry`, `wait_until`

### MCP 驱动生成或修 Flow

先看：

- `example-index.md` 里的 MCP Flows
- `examples.md` 里的 MCP 模板

## 中文排错与业务模板入口

- 常见报错和修 Flow 对策: `zh-cn-troubleshooting.md`
- 登录、搜索、导入、导出、表格抓取、MCP finalize 标准问法: `zh-cn-business-templates.md`
- selector 策略速查: `zh-cn-selectors.md`
- Flow review 清单: `zh-cn-review-checklist.md`

## 中文写作风格建议

- `name` 用任务意图，不要只写 `tmp`、`test`、`demo2`
- `description` 说清楚最后交付什么结果
- `save_as` 用业务语义，比如 `import_results`、`auth_status`、`page_title`
- artifact 路径尽量稳定，不要每次都换一套随意命名
- 能用 Flow 原生 action 就别先绕到 Lua
- 发邮件场景下，不要把邮箱密码直接硬编码进长期复用模板，优先走 `connection`

## 一句总结

中文需求先整理成 “页面 + 目标 + 输入 + 输出 + 授权”，再从最小 action 组合起手，最后再补可维护性。
