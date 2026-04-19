# Lesson 32: 截一张元素级截图

整页截图适合看上下文，  
元素级截图更适合聚焦到真正有问题的区域。

使用页面：
[../../demo/debug_artifacts.html](../../demo/debug_artifacts.html)

目标：

- 对指定元素截图
- 理解“整页截图”和“元素截图”的分工
- 把截图路径写回 JSON

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
[../../script/tutorials/32_element_screenshot.lua](../../script/tutorials/32_element_screenshot.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/32_element_screenshot.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/32_element_screenshot.lua
```

预期结果：

- 会生成 `artifacts/tutorials/32-debug-artifacts-card-lua.png`
- 会生成 `artifacts/tutorials/32-element-screenshot-lua.json`

## Step 2: 这节和上一节的边界

- 整页截图：看“现场全貌”
- 元素截图：看“问题核心区域”

真实排障里，这两种图通常是一起留的，不是二选一。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/32_element_screenshot.flow.yaml](../../script/tutorials/32_element_screenshot.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/32_element_screenshot.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/32_element_screenshot.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/32-debug-artifacts-card-flow.png`
- 会生成 `artifacts/tutorials/32-element-screenshot-flow.json`

## 下一节

下一节继续把现场保存下来，不过这次保存的是 HTML。
[Lesson 33](33-save-html-basics.md)
