# Lesson 111: 先用 `tsplay.list_actions` 看清 MCP 到底能做什么

`Lesson 101-110` 解决的是“模板发布页怎么稳定跑”。  
从这一节开始，我们沿着同一张页面继续往上走一层：先不手写 Flow，而是先看 `TSPlay MCP` 已经暴露了哪些能力。

目标：

- `tsplay.list_actions`
- `mcp-tool`
- `artifacts/tutorials/` 里的 MCP 输出

## 准备工作

先确认已经构建过二进制，或者准备直接用 `go run .`。  
为了把 MCP 输出留存下来，先创建输出目录：

```bash
mkdir -p artifacts/tutorials
```

## Step 1: 列出 TSPlay MCP 工具

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.list_actions > artifacts/tutorials/111-mcp-list-actions.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.list_actions > artifacts/tutorials/111-mcp-list-actions.json
```

预期结果：

- 会生成 `artifacts/tutorials/111-mcp-list-actions.json`
- 里面会看到 `tool`、`actions`
- `actions` 会把观察、起草、校验、执行、修复、会话这些 MCP 能力列出来

## Step 2: 这一节意味着什么

到这里，你先建立的是“能力地图”，不是马上写 Flow。  
这一步很像先看 API 目录，后面写 `observe_page`、`draft_flow`、`run_flow` 时就不容易跳步。

## 下一步

继续看：
[Lesson 112](112-mcp-flow-schema-and-examples.md)
