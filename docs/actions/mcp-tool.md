# Action: `mcp-tool`

`mcp-tool` 是最适合“单个工具直接验证”的 action。它不会先起完整交互环境，而是直接调用一次指定的 MCP 工具并输出 JSON。

## 最小命令

```bash
go run . -action mcp-tool -tool tsplay.list_actions
```

## 常见用法

```bash
go run . -action mcp-tool \
  -tool tsplay.observe_page \
  -args-file script/tutorials/113_mcp_observe_page_template_release.args.json
```

## 常用参数

- `-tool`：必填，要调用的 MCP 工具名
- `-args-json`：直接传入 JSON 参数
- `-args-file`：从 JSON 文件读取参数
- `-flow-root`：限制工具可访问的 Flow 根目录
- `-artifact-root`：工具运行产物根目录

## 适合什么时候用

- 想看 `tsplay.list_actions`、`tsplay.observe_page`、`tsplay.finalize_flow` 的原始输出
- 想给教程、测试或排障留一份结构化 JSON
- 想先确认工具是否工作，再接进 Agent

## 注意事项

- `-tool` 不传时，这个 action 无法工作
- `-args-json` 和 `-args-file` 二选一更清晰

## 相关文档

- [Lesson 111](../tutorials/111-mcp-list-actions.md)
- [Lesson 113](../tutorials/113-mcp-observe-page.md)
- [Lesson 120](../tutorials/120-mcp-finalize-flow.md)
