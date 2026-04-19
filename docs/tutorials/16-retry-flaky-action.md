# Lesson 16: 用 `retry` 处理偶发失败动作

这一节开始进入控制流。  
我们使用仓库里的 [../../demo/retry_wait_until.html](../../demo/retry_wait_until.html)。

这个 demo 的设计很简单：

- 第一次点击只会进入准备态
- 第二次点击才会真正成功

这正适合拿来讲 `retry`。

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

页面地址：

```text
http://127.0.0.1:8000/demo/retry_wait_until.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/16_retry_flaky_action.lua](../../script/tutorials/16_retry_flaky_action.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/16_retry_flaky_action.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/16_retry_flaky_action.lua
```

预期结果：

- 会生成 `artifacts/tutorials/16-retry-flaky-action-lua.json`

这份 Lua 版本故意用“显式循环”的方式写出来，目的是让你看清：

- 什么时候再次尝试
- 什么时候停止

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/16_retry_flaky_action.flow.yaml](../../script/tutorials/16_retry_flaky_action.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/16_retry_flaky_action.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/16_retry_flaky_action.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/16-retry-flaky-action-flow.json`

## Step 3: 这节真正想让你理解什么

- Lua 里做重试，通常要自己写循环
- Flow 里有原生 `retry`
- `retry` 更适合和 `assert_text`、`assert_visible` 组合

也就是说，`retry` 不是“多点几次”，而是：

- 再做一次动作
- 再做一次验证
- 直到成功或彻底失败

## 下一节

下一节继续控制流，但主题换成“等一个异步状态真正完成”，也就是 `wait_until`。
[Lesson 17](17-wait-until-ready.md)
