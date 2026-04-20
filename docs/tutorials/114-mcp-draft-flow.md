# Lesson 114: 用 `draft_flow` 把 observation 变成第一份 Flow 草稿

`Lesson 113` 已经拿到模板发布页的 observation。  
这一节开始让 MCP 帮我们起草第一份“查看模板发布内容”的 Flow，但还不直接执行。

目标：

- `tsplay.draft_flow`
- `observation`
- `draft.flow_yaml`
- 一个更稳的内容提取 intent

## 准备工作

先确认上一节的输出已经存在：

- `artifacts/tutorials/113-mcp-observe-page-template-release.json`

示例参数文件：
[../../script/tutorials/114_mcp_draft_flow_template_release.args.json](../../script/tutorials/114_mcp_draft_flow_template_release.args.json)

这个参数文件会直接复用上一节的 observation 输出，  
同时把页面 URL 显式带回去，保证 draft 里会补上导航步骤。  
这一节的 intent 会收成更稳的一句：

- `查看模板发布内容`

## Step 1: 运行 `draft_flow`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.draft_flow -args-file script/tutorials/114_mcp_draft_flow_template_release.args.json > artifacts/tutorials/114-mcp-draft-flow-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.draft_flow -args-file script/tutorials/114_mcp_draft_flow_template_release.args.json > artifacts/tutorials/114-mcp-draft-flow-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/114-mcp-draft-flow-template-release.json`
- 里面会看到 `draft`
- `draft` 里会包含 `flow_yaml`
- 通常还会包含 `validation`、推荐示例和一些编写提示
- 这份草稿会优先提取页面 headline 和 summary_text

## Step 2: 这一节意味着什么

到这里，MCP 已经把“观察结果”变成了“Flow 草稿”。  
但它还只是草稿，所以接下来先做 `validate_flow`，再决定要不要运行。

## 下一步

继续看：
[Lesson 115](115-mcp-validate-drafted-flow.md)
