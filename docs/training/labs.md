# TSPlay 实训实验

这些实验默认复用仓库现成素材，尽量减少对外部网站的依赖。建议先把仓库根目录通过任意静态文件服务器暴露出来，再访问 `demo/` 目录下的页面。

例如，当仓库根目录被映射到 `<host>` 时，页面路径会是：

- `<host>/demo/demo.html`
- `<host>/demo/tables.html`
- `<host>/demo/upload.html`
- `<host>/demo/multi_upfile.html`

## 实验设计原则

- 先本地页面，后线上页面
- 先 CLI / Lua，后 Flow
- 先稳定页面，后不稳定页面
- 每个实验都要有提交物和通过标准

## Lab 1：CLI 热身

- 难度：L1
- 目标：在 CLI 中完成基础导航、等待和元素状态判断
- 素材：
  - [../../demo/demo.html](../../demo/demo.html)
  - [../../script/is_sel.lua](../../script/is_sel.lua)
- 任务：
  - 启动 CLI
  - 打开 `demo/demo.html`
  - 判断默认选中的选项
  - 切换到另一个选项并再次验证
- 提交物：
  - 1 个 Lua 脚本
  - 1 份运行截图或终端输出
- 通过标准：
  - 脚本中至少使用 `navigate`、`wait_for_network_idle`、`select_option`、`is_selected`

## Lab 2：表格提取

- 难度：L1-L2
- 目标：从本地表格页面提取结构化数据
- 素材：
  - [../../demo/tables.html](../../demo/tables.html)
- 任务：
  - 打开页面并等待表格可见
  - 使用 `capture_table`
  - 打印提取出的结构化结果
- 提交物：
  - 1 个 Lua 脚本或 1 条 Flow
- 通过标准：
  - 结果中必须包含表头和至少 2 行数据
  - 学员能解释为什么这里 `capture_table` 比 `get_html` 更合适

## Lab 3：单文件与多文件上传

- 难度：L2
- 目标：练会文件动作和文件路径管理
- 素材：
  - [../../demo/upload.html](../../demo/upload.html)
  - [../../demo/multi_upfile.html](../../demo/multi_upfile.html)
- 任务：
  - 用 `upload_file` 完成单文件上传场景
  - 用 `upload_multiple_files` 完成多文件上传场景
  - 观察文件信息区域是否变化
- 提交物：
  - 2 条 Flow，或 1 条包含两个场景的 Flow
- 通过标准：
  - 至少 1 条 Flow 启用了文件访问授权说明
  - 学员能解释为什么文件类动作在 MCP 下需要 `allow_file_access`

## Lab 4：从 Lua 改写成 Flow

- 难度：L2
- 目标：把命令式脚本转成结构化 Flow
- 素材：
  - [../../script/open_url.lua](../../script/open_url.lua)
- 任务：
  - 把 Lua 脚本改写为 Flow
  - 引入 `vars`
  - 将提取出的链接保存到 `save_as`
- 提交物：
  - 1 条 Flow
- 通过标准：
  - Flow 能通过 `validate_flow`
  - 使用 `vars` 和 `save_as`
  - 命名清楚、步骤顺序可审阅

## Lab 5：增强健壮性

- 难度：L3
- 目标：学会给 Flow 加控制流和失败恢复
- 素材：
  - Lab 4 的 Flow
  - 任一需要等待状态变化的 demo 页面
- 任务：
  - 至少加入一种控制流：`retry` / `if` / `foreach`
  - 至少加入一种恢复或等待机制：`on_error` / `wait_until`
  - 用 `extract_text` 和 `set_var` 生成一个业务变量
- 提交物：
  - 1 条增强版 Flow
  - 1 份设计说明，解释你为什么这样组织控制流
- 通过标准：
  - Flow 中同时出现“业务变量”和“健壮性设计”
  - 学员能说明为什么没有直接用 `sleep` 替代所有等待逻辑

## Lab 6：MCP 草拟与修复

- 难度：L4
- 目标：练会从“用户意图”进入 MCP 工具链
- 素材：
  - 任一 `demo/` 页面
  - MCP server
- 任务：
  - 用 `tsplay.observe_page` 或 `tsplay.draft_flow` 草拟一条 Flow
  - 用 `tsplay.validate_flow` 校验
  - 人为制造一个 selector 错误
  - 用 `tsplay.repair_flow_context` 或 `tsplay.repair_flow` 组织修复
- 提交物：
  - 1 份 MCP 操作记录
  - 1 条修复前 Flow
  - 1 条修复后 Flow
- 通过标准：
  - 能解释 MCP 工具之间的顺序关系
  - 能指出 repair 依据来自哪些 artifact 或 trace

## Stretch Lab：命名会话与长期复用

- 难度：L4
- 目标：理解命名会话的价值
- 素材：
  - 任意需要保存登录态的业务环境
- 任务：
  - 保存一个命名会话
  - 在 Flow 顶层使用 `browser.use_session`
  - 导出推荐的 Flow 片段
- 提交物：
  - 1 个命名会话说明
  - 1 条使用会话的 Flow
- 通过标准：
  - 学员能说明 `use_session`、`storage_state`、`persistent profile` 的适用场景差异

## 实验提交要求

每个实验都建议附一份实验报告，模板见 [templates/lab-report.md](templates/lab-report.md)。最低需要包含：

- 实验目标
- 使用的模式：Lua / Flow / MCP
- 最终脚本或 Flow
- 运行结果
- 遇到的问题与修复方式

## 讲师使用建议

- L1-L2 学员优先做 Lab 1-4
- L3 学员至少要完成 Lab 5
- L4 学员必须完成 Lab 6 或 Stretch Lab
- 训练营中，讲师最好先统一演示 1 次，再让学员分组完成
