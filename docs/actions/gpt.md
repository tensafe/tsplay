# Action: `gpt`

`gpt` 目前还是一个预留占位入口，不是完整工作流。

## 最小命令

```bash
go run . -action gpt
```

## 当前行为

现在执行它时，只会输出一条占位提示，不会像 `srv`、`mcp-tool`、`workbench-api` 那样真正启动可用能力。

## 什么时候才需要看它

- 你在梳理仓库现状，想确认这个入口是不是已经实现
- 你准备后续补全这个入口的产品设计或代码实现

## 更推荐用什么

- 要给 Agent 暴露能力：看 [srv](srv.md)
- 要做本地 MCP 集成：看 [mcp-stdio](mcp-stdio.md)
- 要直接验证单个 MCP 工具：看 [mcp-tool](mcp-tool.md)

## 注意事项

- 不建议把 `gpt` 当成当前可交付入口
- 如果页面或 README 里要给新用户推荐入口，通常不应优先提它

## 相关文档

- [核心功能路线图](../product/core-feature-roadmap.md)
- [核心功能执行面板](../product/core-feature-execution-board.md)
