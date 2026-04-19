# Lesson 04: 提取文本和 HTML 片段

这一节继续使用本地静态文件服务，不过目标从“点页面”变成“读页面”。

我们使用仓库里的 [../../demo/extract.html](../../demo/extract.html)，重点体验两个动作：

- `extract_text`
- `get_html`

目标：

- 提取标题文本
- 从计数区域里提取数字
- 读取一个局部 HTML 片段
- 把结果写到 `artifacts/tutorials/`

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
http://127.0.0.1:8000/demo/extract.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/04_extract_text_and_html.lua](../../script/tutorials/04_extract_text_and_html.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/04_extract_text_and_html.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/04_extract_text_and_html.lua
```

预期结果：

- 浏览器会打开本地页面
- 会生成 `artifacts/tutorials/04-extract-text-and-html-lua.json`

## Step 2: 看输出里有什么

输出里会包含这些字段：

- `page_title`
- `order_count`
- `notice_html`

这三个字段正好能说明区别：

- `extract_text` 适合拿“人看到的文本”
- `get_html` 适合拿“某一块 DOM 片段”

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/04_extract_text_and_html.flow.yaml](../../script/tutorials/04_extract_text_and_html.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/04_extract_text_and_html.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/04_extract_text_and_html.flow.yaml -headless
```

如果你想看浏览器过程，可以去掉 `-headless`。

预期结果：

- 会生成 `artifacts/tutorials/04-extract-text-and-html-flow.json`
- 终端会输出结构化 trace

## Step 4: 这节要记住什么

- 要拿文本，优先想 `extract_text`
- 要拿某块 HTML，优先想 `get_html`
- 如果只是为了得到一个数字或状态词，通常没必要先拿整页 HTML

## 下一节

下一节把页面动作换成 HTTP 请求，不过依然不依赖外网：
[Lesson 05](05-http-request-json.md)
