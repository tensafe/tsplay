# Lesson 35: 在失败分支里保存错误现场

这一节不再只保存“成功时的页面”，  
而是故意触发一次已知失败，然后把失败现场保存下来。

使用页面：
[../../demo/import_workflow.html](../../demo/import_workflow.html)

目标：

- 故意触发一次表单校验失败
- 在错误分支里保存截图和 HTML
- 把错误状态整理成一份 JSON

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/35_error_evidence_pack.lua](../../script/tutorials/35_error_evidence_pack.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/35_error_evidence_pack.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/35_error_evidence_pack.lua
```

预期结果：

- 会生成 `artifacts/tutorials/35-error-evidence-pack-lua.png`
- 会生成 `artifacts/tutorials/35-error-evidence-pack-lua.html`
- 会生成 `artifacts/tutorials/35-error-evidence-pack-lua.json`

## Step 2: 这节和上一节的差别

上一节保存的是“成功时的证据包”。  
这一节保存的是“失败时的证据包”。

这意味着你开始从“会跑 demo”进入“会留排障现场”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/35_error_evidence_pack.flow.yaml](../../script/tutorials/35_error_evidence_pack.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/35_error_evidence_pack.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/35_error_evidence_pack.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/35-error-evidence-pack-flow.png`
- 会生成 `artifacts/tutorials/35-error-evidence-pack-flow.html`
- 会生成 `artifacts/tutorials/35-error-evidence-pack-flow.json`

## 下一步

到这里，初级阶段已经不只是“跑脚本”，  
而是开始具备会话观察、状态快照和失败现场保留能力了。

如果你要继续按课程体系推进，可以回到：
[track-junior.zh-CN.md](track-junior.zh-CN.md)
