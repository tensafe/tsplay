# Lesson 03: 抓取本地表格并写出 JSON

这一节继续复用本地静态文件服务和仓库里的 [../../demo/tables.html](../../demo/tables.html)。

目标：

- 打开本地表格页
- 抓取 `#myTable`
- 把结构化表格内容写到 `artifacts/tutorials/`

## 准备工作

如果上一节启动的 TSPlay 内置静态文件服务还在运行，可以直接复用。  
如果没有运行，就在仓库根目录重新执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

表格页地址：

```text
http://127.0.0.1:8000/demo/tables.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/03_capture_table.lua](../../script/tutorials/03_capture_table.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/03_capture_table.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/03_capture_table.lua
```

预期结果：

- 浏览器会打开本地表格页
- 脚本会抓取 `#myTable`
- 会生成 `artifacts/tutorials/03-capture-table-lua.json`

和上一节一样，浏览器会保持打开，结束时按 `Ctrl+C` 即可。

## Step 2: 看抓出来的数据

打开：

```text
artifacts/tutorials/03-capture-table-lua.json
```

你会看到一个结构化表格结果。  
这就是 `capture_table` 比直接拿整页 HTML 更适合很多交付场景的原因之一。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/03_capture_table.flow.yaml](../../script/tutorials/03_capture_table.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/03_capture_table.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/03_capture_table.flow.yaml -headless
```

预期结果：

- Flow 会抓取同一张表
- 结果会写到 `artifacts/tutorials/03-capture-table-flow.json`
- 终端会输出每一步的执行结果

## Step 4: 为什么这节很关键

到这里，你已经有了一个很典型的迁移路径：

1. 先用 `Lua` 试出可行路径
2. 再把同样的能力固化成 `Flow`
3. 后续再继续往 `write_csv`、`foreach`、断言、失败修复扩展

这也是 TSPlay 在真实交付里很自然的一条主线。

## 下一节

下一节改成“读页面内容”，分别看 `extract_text` 和 `get_html`：
[Lesson 04](04-extract-text-and-html.md)
