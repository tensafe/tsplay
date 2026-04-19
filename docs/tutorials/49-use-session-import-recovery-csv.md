# Lesson 49: 用命名会话做带恢复的 CSV 批量导入

上一节我们处理的是全量成功的数据。  
这一节开始把坏数据也放进来，看受保护流程里的局部恢复怎么写。

目标：

- `use_session`
- `read_csv`
- `on_error`
- 错误行结果回写

## 开始前

建议先跑完：

- [Lesson 48](48-use-session-batch-import-csv.md)

默认输入文件：

- [../../demo/data/import_users_with_error.csv](../../demo/data/import_users_with_error.csv)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/49_use_session_import_recovery_csv.lua](../../script/tutorials/49_use_session_import_recovery_csv.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/49_use_session_import_recovery_csv.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/49_use_session_import_recovery_csv.lua
```

预期结果：

- 会生成 `artifacts/tutorials/49-use-session-import-recovery-csv-lua.json`
- 会生成 `artifacts/tutorials/49-use-session-import-recovery-csv-lua.csv`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/49_use_session_import_recovery_csv.flow.yaml](../../script/tutorials/49_use_session_import_recovery_csv.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/49_use_session_import_recovery_csv.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/49_use_session_import_recovery_csv.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/49-use-session-import-recovery-csv-flow.json`
- 会生成 `artifacts/tutorials/49-use-session-import-recovery-csv-flow.csv`

## Step 3: 这节要理解什么

这一步不是为了“让错误消失”。  
而是为了让错误变成结构化结果：

- 哪一行成功了
- 哪一行失败了
- 为什么失败

一旦你能把失败写回结果文件，后续复盘和补数就都容易了。

## 下一节

下一节把同样的思路迁移到 Excel 输入：
[Lesson 50](50-use-session-batch-import-excel.md)
