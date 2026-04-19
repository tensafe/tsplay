# Lesson 31: 截一张完整页面截图

这一节开始进入调试产物。  
最先学的不是复杂修复，而是先留下一张清楚的整页截图。

使用页面：
[../../demo/debug_artifacts.html](../../demo/debug_artifacts.html)

目标：

- 打开一个稳定本地页面
- 保存完整页面截图
- 把截图路径写进 JSON

## 准备工作

如果静态文件服务没有运行，先执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/31_full_page_screenshot.lua](../../script/tutorials/31_full_page_screenshot.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/31_full_page_screenshot.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/31_full_page_screenshot.lua
```

预期结果：

- 会生成 `artifacts/tutorials/31-debug-artifacts-full-page-lua.png`
- 会生成 `artifacts/tutorials/31-full-page-screenshot-lua.json`

## Step 2: 为什么整页截图通常是第一张图

因为它最适合回答两个问题：

- 当时页面整体长什么样
- 目标元素到底在不在、位置对不对

这是后面做元素截图和 HTML 保存的基线。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/31_full_page_screenshot.flow.yaml](../../script/tutorials/31_full_page_screenshot.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/31_full_page_screenshot.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/31_full_page_screenshot.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/31-debug-artifacts-full-page-flow.png`
- 会生成 `artifacts/tutorials/31-full-page-screenshot-flow.json`

## 下一节

下一节把整页截图收窄成元素级截图。
[Lesson 32](32-element-screenshot.md)
