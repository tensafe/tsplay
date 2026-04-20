# Lesson 117: 先故意跑坏一次，再生成 repair context

前面 `113-116` 走的是成功链。  
这一节开始故意制造一个坏 selector，让你第一次看到 MCP repair 不是“凭空修”，而是先基于失败现场整理上下文。

目标：

- `tsplay.run_flow`
- `tsplay.repair_flow_context`
- 失败 trace
- repair context

## 准备工作

先确认本地静态文件服务仍然在 `:8000`。

示例文件：

- 破坏版 Flow：
  [../../script/tutorials/117_mcp_template_release_summary_broken.flow.yaml](../../script/tutorials/117_mcp_template_release_summary_broken.flow.yaml)
- 失败运行参数：
  [../../script/tutorials/117_mcp_run_broken_template_release.args.json](../../script/tutorials/117_mcp_run_broken_template_release.args.json)
- repair context 参数：
  [../../script/tutorials/117_mcp_repair_flow_context_template_release.args.json](../../script/tutorials/117_mcp_repair_flow_context_template_release.args.json)

## Step 1: 先运行一份故意写坏的 Flow

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.run_flow -args-file script/tutorials/117_mcp_run_broken_template_release.args.json > artifacts/tutorials/117-mcp-run-broken-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.run_flow -args-file script/tutorials/117_mcp_run_broken_template_release.args.json > artifacts/tutorials/117-mcp-run-broken-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/117-mcp-run-broken-template-release.json`
- 里面通常会看到 `ok=false`
- 会记录失败 step 和运行 trace

## Step 2: 基于失败结果生成 repair context

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.repair_flow_context -args-file script/tutorials/117_mcp_repair_flow_context_template_release.args.json > artifacts/tutorials/117-mcp-repair-flow-context-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.repair_flow_context -args-file script/tutorials/117_mcp_repair_flow_context_template_release.args.json > artifacts/tutorials/117-mcp-repair-flow-context-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/117-mcp-repair-flow-context-template-release.json`
- 里面会看到 `ok=true`
- 还会看到 `context`
- `context` 里会整理失败类别、trace 摘要、相关 selector 和 repair hints

## Step 3: 这一节意味着什么

repair 的第一步不是“直接改 YAML”。  
而是先把失败变成一份对 AI 和人都能读懂的结构化上下文。

## 下一步

继续看：
[Lesson 118](118-mcp-repair-flow.md)
