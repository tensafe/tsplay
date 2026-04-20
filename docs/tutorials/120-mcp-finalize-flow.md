# Lesson 120: 用 `finalize_flow` 收成一份更短的默认入口

`Lesson 113-119` 已经把 MCP 基础链路拆开讲完了。  
这一节回到“更适合日常使用”的视角：直接用 `finalize_flow` 把 observation 收成一份更接近可交付状态的结果。

目标：

- `tsplay.finalize_flow`
- `status`
- `flow_yaml`
- 一个更短的内容提取默认入口

## 准备工作

先确认：

- `artifacts/tutorials/113-mcp-observe-page-template-release.json` 已存在

示例参数文件：
[../../script/tutorials/120_mcp_finalize_flow_template_release.args.json](../../script/tutorials/120_mcp_finalize_flow_template_release.args.json)

这个参数文件会直接复用 `Lesson 113` 的 observation，  
同时把页面 URL 显式带回去，并使用和 `Lesson 114` 同一条更稳的 intent：

- `查看模板发布内容`

## Step 1: 运行 `finalize_flow`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.finalize_flow -args-file script/tutorials/120_mcp_finalize_flow_template_release.args.json > artifacts/tutorials/120-mcp-finalize-flow-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.finalize_flow -args-file script/tutorials/120_mcp_finalize_flow_template_release.args.json > artifacts/tutorials/120-mcp-finalize-flow-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/120-mcp-finalize-flow-template-release.json`
- 里面会看到 `status`
- 还会看到 `flow_yaml`
- 如果还有阻塞点，也会看到 `blocking_reason`、`unresolved` 或推荐示例

## Step 2: 这一节意味着什么

到这里，中级教程里的 MCP 主线就完整了：

- 你知道怎么先观察页面
- 知道怎么拿草稿
- 知道为什么要先校验再运行
- 也知道失败时怎么整理 repair context 和 repair request
- 最后还知道什么时候可以直接走 `finalize_flow`

## 下一步

如果你是按课程体系推进，  
可以回到：
[track-intermediate.md](track-intermediate.md)

如果你是按长期演化推进，  
可以继续看：
[iteration-roadmap-160.md](iteration-roadmap-160.md)
