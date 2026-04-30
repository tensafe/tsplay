# Action: `workbench-api`

`workbench-api` 会同时启动 Workbench 的页面和 API。它更像一个可直接打开的内置工作台入口，而不只是单个命令。

## 最小命令

```bash
go run . -action workbench-api
```

## 常见用法

```bash
go run . -action workbench-api -addr :8082 -artifact-root artifacts
```

## 常用参数

- `-addr`：监听地址
- `-serve-root`：可选，本地静态目录；不传时使用二进制内置资源
- `-artifact-root`：Workbench 读取和展示运行产物的目录

## 运行后会看到什么

- 根路径会跳到 `/demo/workbench.html`
- 同时会暴露 `/api/workbench/health`
- `artifact-root` 下的内容会通过 `/workbench-artifacts/` 暴露给页面

## 适合什么时候用

- 想直接打开内置 Workbench 页面
- 想边看 UI，边看 artifact 和会话数据
- 想验证工作台层的页面与 API 是否接通

## 注意事项

- 如果你只需要 MCP 工具接口，不一定要先起 Workbench
- 如果你正在本地改页面资源，可以搭配 `-serve-root` 使用

## 相关文档

- [产品方案](../product/README.md)
- [核心功能执行面板](../product/core-feature-execution-board.md)
