# Lesson 17: 用 `wait_until` 等异步状态完成

这一节继续使用 [../../demo/retry_wait_until.html](../../demo/retry_wait_until.html)，不过换到另一块区域。

这个 demo 的异步任务不会立刻完成，而是延迟一小段时间才把页面状态改成完成态。  
这正好用来理解 `wait_until`。

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
[../../script/tutorials/17_wait_until_ready.lua](../../script/tutorials/17_wait_until_ready.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/17_wait_until_ready.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/17_wait_until_ready.lua
```

预期结果：

- 会生成 `artifacts/tutorials/17-wait-until-ready-lua.json`

这份 Lua 版本仍然是显式轮询，好让你先把机制看明白。

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/17_wait_until_ready.flow.yaml](../../script/tutorials/17_wait_until_ready.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/17_wait_until_ready.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/17_wait_until_ready.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/17-wait-until-ready-flow.json`

## Step 3: `wait_until` 和 `sleep` 的区别

这节最关键的一点是：

- `sleep` 只是盲等
- `wait_until` 是带条件的等

所以更稳定的写法通常是：

- 先声明我在等什么
- 再设定总超时和轮询间隔

而不是“先睡两秒，看看运气”。

## 下一节

下一节把文件动作接回浏览器：上传单个本地文件。
[Lesson 18](18-upload-single-file.md)
