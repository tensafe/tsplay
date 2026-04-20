# Lesson 118: 用 repair context 生成真正可用的修复请求

`Lesson 117` 已经把失败现场整理成 repair context。  
这一节继续往前走，把它变成一份真正给 AI 或人工修复用的统一 repair request。

目标：

- `tsplay.repair_flow`
- `repair.prompt`
- `repair.repair_hints`

## 准备工作

先确认上一节输出已经存在：

- `artifacts/tutorials/117-mcp-repair-flow-context-template-release.json`

示例参数文件：
[../../script/tutorials/118_mcp_repair_flow_template_release.args.json](../../script/tutorials/118_mcp_repair_flow_template_release.args.json)

## Step 1: 运行 `repair_flow`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.repair_flow -args-file script/tutorials/118_mcp_repair_flow_template_release.args.json > artifacts/tutorials/118-mcp-repair-flow-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.repair_flow -args-file script/tutorials/118_mcp_repair_flow_template_release.args.json > artifacts/tutorials/118-mcp-repair-flow-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/118-mcp-repair-flow-template-release.json`
- 里面会看到 `ok=true`
- 还会看到 `repair`
- `repair` 里会包含统一的 `prompt` 和 `repair_hints`

## Step 2: 这一节意味着什么

到这里，repair 这一段已经不是一句“帮我修一下”。  
而是变成了一份带失败上下文、建议线索和修复提示的结构化请求。

## 下一步

继续看：
[Lesson 119](119-mcp-chain-overview.md)
