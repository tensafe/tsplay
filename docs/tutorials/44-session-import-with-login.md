# Lesson 44: 登录受会话保护的导入页并完成一条导入

前面 `36-43` 讲的是“怎么保存和复用登录态”。  
这一节开始把这些状态真正用到业务流程里。

我们会进入一个新的本地 demo 页：
[../../demo/session_import_workflow.html](../../demo/session_import_workflow.html)

目标：

- 先认识“受会话保护的页面”长什么样
- 没有登录态时，先完成一次登录
- 登录后提交一条导入数据

## 开始前

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
[../../script/tutorials/44_session_import_with_login.lua](../../script/tutorials/44_session_import_with_login.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/44_session_import_with_login.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/44_session_import_with_login.lua
```

预期结果：

- 会先清掉当前页面上的旧会话
- 会自动登录为 `demo-user`
- 会提交一条导入记录
- 会生成 `artifacts/tutorials/44-session-import-with-login-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/44_session_import_with_login.flow.yaml](../../script/tutorials/44_session_import_with_login.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/44_session_import_with_login.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/44_session_import_with_login.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/44-session-import-with-login-flow.json`

## Step 3: 这节要理解什么

这一节的重点不是“再写一次登录”。  
重点是先看清楚受保护页面的边界：

- 没有会话时，要先登录
- 有会话时，应该直接进入导入表单

先把这个边界看清，后面复用状态文件和命名会话时才不会迷糊。

## 下一节

下一节开始用前面保存下来的状态文件，直接跳过登录：
[Lesson 45](45-storage-state-auth-import.md)
