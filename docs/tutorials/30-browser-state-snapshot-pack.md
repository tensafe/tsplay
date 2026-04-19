# Lesson 30: 生成一份浏览器状态快照

前两节分别读了：

- `storage state`
- `cookie header`

这一节把它们合并成一个更接近真实排障和交付的“状态快照”。

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 读取页面状态文本
- 同时抓到 `storage state` 和 `cookie header`
- 写成一份统一 JSON

## 准备工作

先确认本地静态文件服务还在运行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/30_browser_state_snapshot_pack.lua](../../script/tutorials/30_browser_state_snapshot_pack.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/30_browser_state_snapshot_pack.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/30_browser_state_snapshot_pack.lua
```

预期结果：

- 会生成 `artifacts/tutorials/30-browser-state-snapshot-pack-lua.json`

## Step 2: 这一节真正练的是什么

不是多学一个 action，  
而是开始练“把多个观察结果整理成统一交付物”。

这个习惯后面会直接用在：

- 调试产物整理
- 登录态排查
- 失败现场回传

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/30_browser_state_snapshot_pack.flow.yaml](../../script/tutorials/30_browser_state_snapshot_pack.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/30_browser_state_snapshot_pack.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/30_browser_state_snapshot_pack.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/30-browser-state-snapshot-pack-flow.json`

## 下一节

下一节从“状态观察”切到“调试产物”，先学最常用的一种：整页截图。
[Lesson 31](31-full-page-screenshot.md)
