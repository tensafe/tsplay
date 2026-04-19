# Lesson 48: 用命名会话驱动 CSV 批量导入

这一节开始把“会话复用”真正接到批量处理上。

目标：

- `use_session`
- `read_csv`
- `foreach`
- 受保护页面里的批量导入

## 开始前

这一节建议先跑完：

- [Lesson 46](46-save-import-session.md)

默认输入文件：

- [../../demo/data/import_users.csv](../../demo/data/import_users.csv)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/48_use_session_batch_import_csv.lua](../../script/tutorials/48_use_session_batch_import_csv.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/48_use_session_batch_import_csv.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/48_use_session_batch_import_csv.lua
```

预期结果：

- 会生成 `artifacts/tutorials/48-use-session-batch-import-csv-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/48_use_session_batch_import_csv.flow.yaml](../../script/tutorials/48_use_session_batch_import_csv.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/48_use_session_batch_import_csv.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/48_use_session_batch_import_csv.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/48-use-session-batch-import-csv-flow.json`

## Step 3: 这节要理解什么

这一节是前面两条线真正会合的地方：

- `Lesson 22` 的批量导入能力
- `Lesson 36-47` 的会话复用能力

一旦两条线合并，你就已经很接近真实业务自动化了。

## 下一节

下一节在同一条受保护流程里体验坏数据恢复：
[Lesson 49](49-use-session-import-recovery-csv.md)
