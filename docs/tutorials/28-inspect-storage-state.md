# Lesson 28: 读取当前浏览器的 Storage State

这一节开始进入“浏览器状态”主题。  
目标不是先讲复杂会话复用，而是先学会把当前页面里的本地状态真正读出来。

这一节使用的页面是：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 在本地页面里种下一份登录态
- 读取当前浏览器的 `storage state`
- 把原始结果写到 `artifacts/tutorials/`

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
[../../script/tutorials/28_inspect_storage_state.lua](../../script/tutorials/28_inspect_storage_state.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/28_inspect_storage_state.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/28_inspect_storage_state.lua
```

预期结果：

- 会生成 `artifacts/tutorials/28-inspect-storage-state-lua.json`

## Step 2: 这节重点看什么

`storage state` 不是一句抽象概念。  
跑完之后，你会在输出里直接看到：

- 当前 origin
- cookies
- local storage

也就是说，这一节是在帮你建立“浏览器状态是可观察对象”的直觉。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/28_inspect_storage_state.flow.yaml](../../script/tutorials/28_inspect_storage_state.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/28_inspect_storage_state.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/28_inspect_storage_state.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/28-inspect-storage-state-flow.json`

## 下一节

下一节继续看浏览器状态，但会切到更贴近请求头的视角：`Cookie` 字符串。
[Lesson 29](29-read-cookies-string.md)
