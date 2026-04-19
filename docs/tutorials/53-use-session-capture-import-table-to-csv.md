# Lesson 53: 把认证导入页上的结果表写成本地 CSV

上一节我们已经把结果表抓出来了。  
这一节继续顺着这条线，把页面表格真正落到本地文件。

目标：

- `capture_table`
- `write_csv`
- 从页面结果到本地交付物

## 开始前

建议先跑完：

- [Lesson 52](52-use-session-capture-import-table.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/53_use_session_capture_import_table_to_csv.lua](../../script/tutorials/53_use_session_capture_import_table_to_csv.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/53_use_session_capture_import_table_to_csv.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/53_use_session_capture_import_table_to_csv.lua
```

预期结果：

- 会生成 `artifacts/tutorials/53-use-session-captured-import-table-lua.csv`
- 会生成 `artifacts/tutorials/53-use-session-capture-import-table-to-csv-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml](../../script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/53-use-session-captured-import-table-flow.csv`
- 会生成 `artifacts/tutorials/53-use-session-capture-import-table-to-csv-flow.json`

## Step 3: 这节要理解什么

从这一节开始，你就可以把“页面上看见的东西”真正变成可归档、可传递的本地文件了。  
这对交付、复盘、补数都很有价值。

## 下一节

下一节把页面自带的导出按钮真正用起来：
[Lesson 54](54-use-session-download-import-report.md)
