# Action: `srv`

`srv` 用来启动 TSPlay 的 MCP 网络服务。它适合把观察、生成、执行、修复、会话这些能力暴露给外部客户端或 Agent。

## 最小命令

```bash
go run . -action srv
```

## 常见用法

```bash
go run . -action srv -addr :8081 -flow-root script -artifact-root artifacts
```

## 常用参数

- `-addr`：监听地址，默认是 `:8082`
- `-flow-root`：允许 MCP 读取或写入 Flow 的根目录
- `-artifact-root`：运行产物和会话产物的根目录

## 适合什么时候用

- 要把 TSPlay 接到 Agent 或平台
- 想通过网络地址统一暴露工具能力
- 想让外部系统调用 `observe_page`、`finalize_flow`、`repair_flow` 等 MCP 工具

## 注意事项

- 如果你只是想直接测一个工具，通常先用 [mcp-tool](mcp-tool.md) 更轻
- 如果你是在桌面 Agent 环境里做本地集成，通常 [mcp-stdio](mcp-stdio.md) 更顺手

## 相关文档

- [AI 协作入门](../training/ai-intent-to-flow.md)
- [Lesson 111](../tutorials/111-mcp-list-actions.md)
- [Lesson 119](../tutorials/119-mcp-chain-overview.md)
