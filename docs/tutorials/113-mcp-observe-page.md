# Lesson 113: 先用 `observe_page` 观察模板发布页

`Lesson 112` 看的是“Flow 长什么样”。  
这一节开始第一次真正把 MCP 用到页面上，但仍然不直接写 Flow，而是先拿一份结构化观察结果。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `tsplay.observe_page`
- `mcp-tool`
- 结构化 observation

## 准备工作

先确认本地静态文件服务已经启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

再确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

## Step 1: 查看参数文件

示例文件：
[../../script/tutorials/113_mcp_observe_page_template_release.args.json](../../script/tutorials/113_mcp_observe_page_template_release.args.json)

这个参数文件会访问：

- `http://127.0.0.1:8000/demo/template_release_lab.html`

## Step 2: 运行 `observe_page`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.observe_page -args-file script/tutorials/113_mcp_observe_page_template_release.args.json > artifacts/tutorials/113-mcp-observe-page-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.observe_page -args-file script/tutorials/113_mcp_observe_page_template_release.args.json > artifacts/tutorials/113-mcp-observe-page-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/113-mcp-observe-page-template-release.json`
- 里面会看到 `ok=true`
- 还会看到 `observation`、`run`
- `observation` 里会包含页面标题和结构化元素信息

## Step 3: 这一节意味着什么

到这里，你拿到的还不是 Flow。  
它更像一份“页面结构快照”，后面的 `draft_flow` 会直接把这份 observation 当输入继续用。

## 下一步

继续看：
[Lesson 114](114-mcp-draft-flow.md)
