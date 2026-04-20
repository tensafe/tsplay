# Lesson 112: 用 `flow_schema` 和 `flow_examples` 看清 Flow 长什么样

`Lesson 111` 先看了 MCP 工具列表。  
这一节继续往前走，但还不直接起草 Flow，而是先看“Flow 结构”和“官方示例”。

目标：

- `tsplay.flow_schema`
- `tsplay.flow_examples`
- `action_manifest`
- `authoring_checklist`

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

## Step 1: 导出 Flow schema

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.flow_schema > artifacts/tutorials/112-mcp-flow-schema.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.flow_schema > artifacts/tutorials/112-mcp-flow-schema.json
```

预期结果：

- 会生成 `artifacts/tutorials/112-mcp-flow-schema.json`
- 里面会看到 `schema`
- 还会看到 `action_manifest`、`generation_rules`、`authoring_checklist`

## Step 2: 导出 Flow 示例

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.flow_examples > artifacts/tutorials/112-mcp-flow-examples.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.flow_examples > artifacts/tutorials/112-mcp-flow-examples.json
```

预期结果：

- 会生成 `artifacts/tutorials/112-mcp-flow-examples.json`
- 里面会看到 `examples`
- 还会看到 `example_selection_hints`

## Step 3: 为什么这一节放在 draft 之前

先看 schema 和 example，后面 `draft_flow` 才不会变成“黑盒猜测”。  
你会更容易理解：为什么生成出来的是结构化 Flow，而不是一段随意拼出来的 YAML。

## 下一步

继续看：
[Lesson 113](113-mcp-observe-page.md)
