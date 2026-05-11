# TSPlay CLI `-action` 参考

这组文档专门回答一个很实际的问题：命令行里的 `-action` 现在到底支持什么，以及每个入口什么时候该用。

如果你想查的是 `navigate`、`click`、`read_csv`、`db_query`、`retry` 这一层能力动作，不要停在这里，直接看 [支持行为清单](../capability-actions/README.md)。

如果你只是想先跑起来，优先看：

- [快速开始](../../getting-started.md)
- [项目总览（中文）](../../README.zh-CN.md)

如果你已经知道自己要查的是某个命令入口，就从下面这张表直接进去。

## 先说明边界

这里整理的是 `go run . -action ...` 这条入口下支持的动作。  
`-flow ...`、`-script ...`、`-browser-video-output ...` 这些能力很重要，但它们不是 `action` 值，所以不放在这张清单里。

## Action 列表

| Action | 用途 | 什么时候优先用 | 说明文档 |
| --- | --- | --- | --- |
| `cli` | 启动交互式 Lua CLI | 临时调试、探索页面、快速试动作 | [cli](cli.md) |
| `gpt` | 预留占位入口 | 只在你明确要看当前占位状态时 | [gpt](gpt.md) |
| `srv` | 启动 MCP 网络服务 | 给 Agent 或外部客户端接 TSPlay 工具 | [srv](srv.md) |
| `workbench-api` | 启动 Workbench UI + API | 想直接打开内置工作台页面 | [workbench-api](workbench-api.md) |
| `mcp-stdio` | 启动 MCP stdio 服务 | 桌面 Agent / 本地 MCP 集成 | [mcp-stdio](mcp-stdio.md) |
| `mcp-tool` | 直接调用一个 MCP 工具 | 验证单个工具输入输出 | [mcp-tool](mcp-tool.md) |
| `file-srv` | 启动内置静态文件服务 | 本地 demo、内置资源、教程练习 | [file-srv](file-srv.md) |
| `demo-srv` | `file-srv` 的兼容别名 | 只在兼容旧命令时 | [demo-srv](demo-srv.md) |
| `list-assets` | 列出二进制内置资源 | 想先确认 release 里带了什么 | [list-assets](list-assets.md) |
| `extract-assets` | 释放内置 docs/script/demo | 单二进制交付、离线学习、培训包 | [extract-assets](extract-assets.md) |
| `quickstart-demo` | 生成并执行最小 demo Flow | 想下载二进制后立刻体验一次，不先碰浏览器 | [quickstart-demo](quickstart-demo.md) |
| `list-record-devices` | 列出 macOS 录屏设备 | 录屏前先查设备和权限 | [list-record-devices](list-record-devices.md) |
| `record-screen` | 录制整个 macOS 桌面 | 录教程、录桌面演示、录窗口切换 | [record-screen](record-screen.md) |
| `save-session` | 保存可复用会话 | 把登录态或 profile 注册下来 | [save-session](save-session.md) |
| `list-sessions` | 列出已保存会话 | 看当前有哪些命名会话 | [list-sessions](list-sessions.md) |
| `get-session` | 查看单个会话详情 | 想确认某个会话保存了什么 | [get-session](get-session.md) |
| `export-session` | 导出可复用片段 | 想把命名会话接回 Flow | [export-session](export-session.md) |
| `delete-session` | 删除命名会话 | 清理不用的登录态或注册记录 | [delete-session](delete-session.md) |

## 按场景选

### 今天先跑通

- 只想下载二进制就立刻体验：先用 [quickstart-demo](quickstart-demo.md)
- 先用 [file-srv](file-srv.md) 起 demo 页面
- 再配合 [mcp-tool](mcp-tool.md) 或 `-flow` 跑第一条例子

### 我要接 Agent

- 本地桌面集成：优先 [mcp-stdio](mcp-stdio.md)
- 通过网络暴露工具：优先 [srv](srv.md)
- 先看工具清单：优先 [mcp-tool](mcp-tool.md)

### 我要做单二进制交付

- 先看 [list-assets](list-assets.md)
- 再看 [extract-assets](extract-assets.md)
- 本地直接服务 demo 时看 [file-srv](file-srv.md)

### 我要复用登录态

- 先存： [save-session](save-session.md)
- 再查： [list-sessions](list-sessions.md)、[get-session](get-session.md)
- 要接回 Flow： [export-session](export-session.md)

## 相关教程

- MCP 主线从 [Lesson 111](../tutorials/111-mcp-list-actions.md) 开始
- 单二进制主线从 [Lesson 142](../tutorials/142-list-assets-for-beginners.md) 开始
- 会话主线从 [Lesson 40](../tutorials/40-save-named-session.md) 开始
- 录屏与讲师用法看 [教程自动录屏](../training/tutorial-video-recording.md)
