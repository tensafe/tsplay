# Lesson 116: 运行刚刚起草并校验过的 Flow

`Lesson 115` 已经确认草稿结构是合法的。  
这一节开始真正执行它，看看 observation 生成出来的 Flow 能不能跑通模板发布页。

目标：

- `tsplay.run_flow`
- `run.trace`
- 运行结果里的提取文本

## 准备工作

先确认：

- 本地静态文件服务仍然在 `:8000`
- `artifacts/tutorials/114-mcp-draft-flow-template-release.json` 已存在

示例参数文件：
[../../script/tutorials/116_mcp_run_drafted_template_release.args.json](../../script/tutorials/116_mcp_run_drafted_template_release.args.json)

## Step 1: 运行 `run_flow`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.run_flow -args-file script/tutorials/116_mcp_run_drafted_template_release.args.json > artifacts/tutorials/116-mcp-run-drafted-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.run_flow -args-file script/tutorials/116_mcp_run_drafted_template_release.args.json > artifacts/tutorials/116-mcp-run-drafted-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/116-mcp-run-drafted-template-release.json`
- 里面会看到 `ok=true`
- 还会看到运行结果、trace 和提取出来的 headline / summary 数据

## Step 2: 这一节意味着什么

到这里，你已经完整走过了一次：

- `observe_page`
- `draft_flow`
- `validate_flow`
- `run_flow`

这就是 MCP 起草 Flow 的最小闭环。

## 下一步

继续看：
[Lesson 117](117-mcp-repair-flow-context.md)
