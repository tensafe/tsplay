# Lesson 33: 保存当前页面的 HTML

截图更像视觉证据，  
HTML 更像结构证据。

这一节先不做复杂分析，只先把页面 HTML 留下来。

使用页面：
[../../demo/debug_artifacts.html](../../demo/debug_artifacts.html)

目标：

- 保存当前页面 HTML
- 把 HTML 文件路径写到 JSON
- 建立“图像证据 + 结构证据”的双视角

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
[../../script/tutorials/33_save_html_basics.lua](../../script/tutorials/33_save_html_basics.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/33_save_html_basics.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/33_save_html_basics.lua
```

预期结果：

- 会生成 `artifacts/tutorials/33-debug-artifacts-page-lua.html`
- 会生成 `artifacts/tutorials/33-save-html-basics-lua.json`

## Step 2: 为什么 HTML 不能省

因为很多问题只看截图是看不出来的，比如：

- 真实 DOM 结构
- class 和 attribute
- 文本是不是在隐藏节点里

所以截图和 HTML 常常要成对保存。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/33_save_html_basics.flow.yaml](../../script/tutorials/33_save_html_basics.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/33_save_html_basics.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/33_save_html_basics.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/33-debug-artifacts-page-flow.html`
- 会生成 `artifacts/tutorials/33-save-html-basics-flow.json`

## 下一节

下一节把整页截图、元素截图和 HTML 一次性打成一个调试产物包。
[Lesson 34](34-debug-artifact-pack.md)
