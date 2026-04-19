# Lesson 52: 抓取认证导入页上的结果表

前面几节我们一直在看“脚本里记录了什么”。  
这一节开始回到页面本身，直接抓取认证导入页上显示出来的结果表。

目标：

- `use_session`
- `capture_table`
- 从页面事实回看导入结果

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)
- [Lesson 48](48-use-session-batch-import-csv.md)

默认输入文件：

- [../../demo/data/import_users.csv](../../demo/data/import_users.csv)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/52_use_session_capture_import_table.lua](../../script/tutorials/52_use_session_capture_import_table.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/52_use_session_capture_import_table.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/52_use_session_capture_import_table.lua
```

预期结果：

- 会生成 `artifacts/tutorials/52-use-session-capture-import-table-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/52_use_session_capture_import_table.flow.yaml](../../script/tutorials/52_use_session_capture_import_table.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/52_use_session_capture_import_table.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/52_use_session_capture_import_table.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/52-use-session-capture-import-table-flow.json`

## Step 3: 这节要理解什么

这一步的重点是：

- 脚本里记录的结果是一层事实
- 页面上最终呈现出来的表格，是另一层事实

很多排障和复盘，其实都要从“页面到底显示了什么”开始。

## 下一节

下一节把这份页面结果表再落成本地 CSV：
[Lesson 53](53-use-session-capture-import-table-to-csv.md)
