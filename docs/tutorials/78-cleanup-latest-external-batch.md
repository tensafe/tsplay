# Lesson 78: 清理最新批次的运行数据，但保留审计

前面我们已经有了同步数据和审计数据。  
这一节开始区分两类信息：

- 运行态数据可以清掉
- 审计留痕应该保留

目标：

- `redis_del`
- `db_execute`
- `write_csv`

## 开始前

建议先跑完：

- [Lesson 75](75-write-external-sync-audit-row.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/78_cleanup_latest_external_batch.lua](../../script/tutorials/78_cleanup_latest_external_batch.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/78_cleanup_latest_external_batch.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/78_cleanup_latest_external_batch.lua
```

预期结果：

- 会生成 `artifacts/tutorials/78-cleanup-latest-external-batch-lua.csv`
- 会生成 `artifacts/tutorials/78-cleanup-latest-external-batch-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/78_cleanup_latest_external_batch.flow.yaml](../../script/tutorials/78_cleanup_latest_external_batch.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/78_cleanup_latest_external_batch.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/78_cleanup_latest_external_batch.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/78-cleanup-latest-external-batch-flow.csv`
- 会生成 `artifacts/tutorials/78-cleanup-latest-external-batch-flow.json`

## 下一节

下一节验证这次清理到底是不是只删了运行数据、保住了审计：
[Lesson 79](79-verify-external-batch-cleanup.md)
