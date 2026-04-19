# Lesson 56: 把认证页面结果表和下载文件放在一起比对

前面几节我们分别拿到了两类结果：

- 页面上的结果表
- 下载下来的 CSV 文件

这一节把它们放到一起。

目标：

- `capture_table`
- `download_file`
- `read_csv`
- 页面结果与文件结果并排保存

## 开始前

建议先跑完：

- [Lesson 52](52-use-session-capture-import-table.md)
- [Lesson 55](55-use-session-download-import-report-readback.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/56_use_session_compare_table_and_download.lua](../../script/tutorials/56_use_session_compare_table_and_download.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/56_use_session_compare_table_and_download.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/56_use_session_compare_table_and_download.lua
```

预期结果：

- 会生成 `artifacts/tutorials/56-use-session-import-report-lua.csv`
- 会生成 `artifacts/tutorials/56-use-session-compare-table-and-download-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/56_use_session_compare_table_and_download.flow.yaml](../../script/tutorials/56_use_session_compare_table_and_download.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/56_use_session_compare_table_and_download.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/56_use_session_compare_table_and_download.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/56-use-session-import-report-flow.csv`
- 会生成 `artifacts/tutorials/56-use-session-compare-table-and-download-flow.json`

## Step 3: 这节要理解什么

这一步很接近真实交付里的复盘方式：

- 页面怎么显示
- 文件怎么导出
- 两边是不是同一回事

一旦能把这两层事实并排保存，问题定位就会容易很多。

## 下一节

下一节把整条认证导入链真正做成一次完整 round trip：
[Lesson 57](57-use-session-import-export-round-trip.md)
