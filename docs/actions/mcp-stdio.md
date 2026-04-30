# Action: `mcp-stdio`

`mcp-stdio` 会通过标准输入输出启动 MCP 服务，适合本地桌面 Agent 集成。

## 最小命令

```bash
go run . -action mcp-stdio
```

## 常见用法

```bash
go run . -action mcp-stdio -flow-root script -artifact-root artifacts
```

## 常用参数

- `-flow-root`：允许读写 Flow 的根目录
- `-artifact-root`：产物和会话根目录

## 适合什么时候用

- 想把 TSPlay 接进本地 MCP 客户端
- 不想额外占一个网络端口
- 希望和桌面 Agent 保持“子进程 + stdio”这条默认集成方式

## 注意事项

- 如果你需要通过网络地址给别的进程访问，改用 [srv](srv.md)
- 如果你只是想单次验证某个工具输入输出，改用 [mcp-tool](mcp-tool.md)

## 相关文档

- [AI 协作入门](../training/ai-intent-to-flow.md)
- [Lesson 119](../tutorials/119-mcp-chain-overview.md)
