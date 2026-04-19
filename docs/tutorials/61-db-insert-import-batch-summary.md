# Lesson 61: 把认证导出结果写成一条 Postgres 批次摘要

这一节从 Redis 过渡到数据库。

我们先不急着写多行明细，而是先保存一条批次摘要：

- 批次 id
- 导出文件路径
- 行数
- 操作人

目标：

- `read_csv`
- `db_execute`
- `db_insert`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 57](57-use-session-import-export-round-trip.md)
- [Lesson 07](07-db-postgres-basics.md)

先加载数据库连接：

```bash
source script/tutorials/env/07_reporting_pg_example.sh
```

然后执行这节新增的初始化 SQL：

```bash
psql "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable" \
  -f script/tutorials/sql/61_reporting_import_sync.sql
```

默认输入文件是 `artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv`。  
如果你前一节跑的是 Lua 版本，可以先切换：

```bash
export TSPLAY_IMPORTED_REPORT=artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/61_db_insert_import_batch_summary.lua](../../script/tutorials/61_db_insert_import_batch_summary.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/61_db_insert_import_batch_summary.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/61_db_insert_import_batch_summary.lua
```

预期结果：

- 会生成 `artifacts/tutorials/61-db-insert-import-batch-summary-lua.json`
- `public.tutorial_import_batches` 里会有一条 `lesson-61-import-batch`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/61_db_insert_import_batch_summary.flow.yaml](../../script/tutorials/61_db_insert_import_batch_summary.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/61_db_insert_import_batch_summary.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/61_db_insert_import_batch_summary.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/61-db-insert-import-batch-summary-flow.json`

## Step 3: 这节意味着什么

到这里，认证导出结果第一次进入了“可查询的持久层”。  
跟 Redis 不同的是，这里开始强调结构化表、固定列和后续查询能力。

## 下一节

下一节继续把“单条写入”扩成“多条查询”：
[Lesson 62](62-db-query-import-batch-summaries.md)
