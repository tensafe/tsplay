# Lesson 21: 用 `if` 处理可选登录分支

这一节开始进入真正的控制流。  
我们使用 [../../demo/import_workflow.html](../../demo/import_workflow.html) 的登录分支模式：

```text
http://127.0.0.1:8000/demo/import_workflow.html?login=1
```

目标：

- 识别页面上是否有登录弹层
- 有就先登录
- 没有就直接进入导入表单

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/21_if_optional_login.lua](../../script/tutorials/21_if_optional_login.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/21_if_optional_login.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/21_if_optional_login.lua
```

预期结果：

- 会生成 `artifacts/tutorials/21-if-optional-login-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/21_if_optional_login.flow.yaml](../../script/tutorials/21_if_optional_login.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/21_if_optional_login.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/21_if_optional_login.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/21-if-optional-login-flow.json`

## Step 3: 这节要理解什么

`if` 最适合解决“页面可能出现，也可能不出现”的状态分支。  
它不是为了炫技，而是为了把这种业务差异写清楚。

## 下一节

下一节继续控制流，但把单次动作扩成“列表里的每一行都做一遍”：
[Lesson 22](22-foreach-batch-import-csv.md)
