# Lesson 51: 用命名会话做带恢复的 Excel 批量导入

`Lesson 50` 处理的是全量成功的 Excel 数据。  
这一节继续沿着同一条认证导入主线往前走，把坏数据也接进来。

目标：

- `use_session`
- `read_excel`
- `on_error`
- 认证态页面里的 Excel 局部恢复

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)
- [Lesson 50](50-use-session-batch-import-excel.md)

默认输入文件：

- [../../demo/data/import_users.xlsx](../../demo/data/import_users.xlsx)

这一节会读取其中的 `ErrorBatch` 工作表。

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/51_use_session_import_recovery_excel.lua](../../script/tutorials/51_use_session_import_recovery_excel.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/51_use_session_import_recovery_excel.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/51_use_session_import_recovery_excel.lua
```

预期结果：

- 会生成 `artifacts/tutorials/51-use-session-import-recovery-excel-lua.json`
- 会生成 `artifacts/tutorials/51-use-session-import-recovery-excel-lua.csv`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/51_use_session_import_recovery_excel.flow.yaml](../../script/tutorials/51_use_session_import_recovery_excel.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/51_use_session_import_recovery_excel.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/51_use_session_import_recovery_excel.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/51-use-session-import-recovery-excel-flow.json`
- 会生成 `artifacts/tutorials/51-use-session-import-recovery-excel-flow.csv`

## Step 3: 这节要理解什么

到这里，认证导入链不只是“能成功跑”。  
它已经开始具备更像真实流程的两个特征：

- 会区分成功行和失败行
- 会把失败原因写回结果文件

## 下一节

下一节先不急着导出文件，而是回到页面本身，抓一份导入结果表：
[Lesson 52](52-use-session-capture-import-table.md)
