# Lesson 110: 跑通一条完整的模板发布稳定性 round trip

这一节是 `101-109` 的收口。

它会重新把这些动作串起来：

- 先做可见性断言
- 再做文本断言
- 跑异步 stage check
- 跑 retry gate
- 等延迟说明项出现
- 处理一次失败恢复
- 通过 reload 确认恢复状态
- 最后保存截图和 HTML

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `assert_visible`
- `assert_text`
- `retry`
- `wait_until`
- `on_error`
- `reload`
- `screenshot`
- `save_html`
- `write_json`

## 准备工作

先确认本地静态文件服务已经启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/110_template_release_robustness_round_trip.lua](../../script/tutorials/110_template_release_robustness_round_trip.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/110_template_release_robustness_round_trip.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/110_template_release_robustness_round_trip.lua
```

预期结果：

- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-lua.png`
- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-lua.html`
- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-lua.json`

## Step 2: 这一节意味着什么

到这里，你不只是会单独用某个 action。  
而是已经能把“断言、等待、重试、恢复、证据留存”组织成一条完整的小型发布链。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/110_template_release_robustness_round_trip.flow.yaml](../../script/tutorials/110_template_release_robustness_round_trip.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/110_template_release_robustness_round_trip.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/110_template_release_robustness_round_trip.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-flow.png`
- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-flow.html`
- 会生成 `artifacts/tutorials/110-template-release-robustness-round-trip-flow.json`

## 下一步

如果你是按课程体系推进，  
可以回到：
[track-intermediate.md](track-intermediate.md)

如果你想继续往 MCP 这一层推进，  
可以继续看：
[Lesson 111](111-mcp-list-actions.md)

如果你是按长期演化推进，  
可以继续看：
[iteration-roadmap-160.md](iteration-roadmap-160.md)
