# Lesson 106: 等一条延迟出现的发布说明项

`Lesson 104` 等的是“已有元素状态变化”。  
这一节等的是“元素本身晚一点才出现”。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `wait_for_selector`
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
[../../script/tutorials/106_wait_for_delayed_release_note.lua](../../script/tutorials/106_wait_for_delayed_release_note.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/106_wait_for_delayed_release_note.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/106_wait_for_delayed_release_note.lua
```

预期结果：

- 会生成 `artifacts/tutorials/106-wait-for-delayed-release-note-lua.json`

## Step 2: 这一节的关键词

如果元素还没出现，  
最稳的写法通常不是先 `sleep`，而是先把“等它出现”这件事显式写出来。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/106_wait_for_delayed_release_note.flow.yaml](../../script/tutorials/106_wait_for_delayed_release_note.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/106_wait_for_delayed_release_note.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/106_wait_for_delayed_release_note.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/106-wait-for-delayed-release-note-flow.json`

## 下一节

下一节继续补“偶发失败”，  
但这次故障点从 gate 变成 click 本身。
[Lesson 107](107-retry-flaky-publish-click.md)
