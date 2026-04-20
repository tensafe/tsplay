# Lesson 103: 用 `retry` 跑通模板发布 gate

到这里，页面本身已经确认没问题。  
下一步就很自然会遇到“第一次检查没过，第二次才过”的场景。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `retry`
- `assert_text`
- `assert_visible`
- `write_json`

## 准备工作

先确认本地静态文件服务已经启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/103_retry_template_release_gate.lua](../../script/tutorials/103_retry_template_release_gate.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/103_retry_template_release_gate.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/103_retry_template_release_gate.lua
```

预期结果：

- 会生成 `artifacts/tutorials/103-retry-template-release-gate-lua.json`

## Step 2: 这一节想建立什么感觉

`retry` 不是“瞎重试”。  
它的核心是：每次重试都带着一个明确的通过标准。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/103_retry_template_release_gate.flow.yaml](../../script/tutorials/103_retry_template_release_gate.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/103_retry_template_release_gate.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/103_retry_template_release_gate.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/103-retry-template-release-gate-flow.json`

## 下一节

下一节换一种不稳定性：  
按钮点下去后不是立刻完成，而是异步变成 ready。
[Lesson 104](104-wait-until-template-release-ready.md)
