# Lesson 115: 先校验草稿，再决定能不能运行

`Lesson 114` 生成了第一份 Flow 草稿。  
这一节不急着运行，而是先检查这份草稿在结构和安全预设上是不是站得住。

目标：

- `tsplay.validate_flow`
- `@jsonpathfile`
- `draft.flow_yaml`

## 准备工作

先确认上一节输出已经存在：

- `artifacts/tutorials/114-mcp-draft-flow-template-release.json`

示例参数文件：
[../../script/tutorials/115_mcp_validate_drafted_template_release.args.json](../../script/tutorials/115_mcp_validate_drafted_template_release.args.json)

这个参数文件会通过 `@jsonpathfile`，直接把 `114` 输出里的 `draft.flow_yaml` 抽出来做校验。

## Step 1: 运行 `validate_flow`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/115_mcp_validate_drafted_template_release.args.json > artifacts/tutorials/115-mcp-validate-drafted-template-release.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/115_mcp_validate_drafted_template_release.args.json > artifacts/tutorials/115-mcp-validate-drafted-template-release.json
```

预期结果：

- 会生成 `artifacts/tutorials/115-mcp-validate-drafted-template-release.json`
- 里面会看到 `valid=true`
- 还会看到 `name`、`steps`、`security`

## Step 2: 这一节意味着什么

到这里，你先建立了一个很重要的顺序：  
`observe -> draft -> validate`，而不是“起草完就直接跑”。

这会让后面的失败定位和 repair 更清楚。

## 下一步

继续看：
[Lesson 116](116-mcp-run-drafted-flow.md)
