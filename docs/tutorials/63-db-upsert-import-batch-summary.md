# Lesson 63: 用 `db_upsert` 更新 Postgres 批次摘要

这一节把数据库动作再往前推一步。

现实里经常不是“永远插新行”，而是：

- 先有一条占位记录
- 后面再用真实结果补全或覆盖

所以这节专门练 `db_upsert`。

目标：

- `db_insert`
- `db_upsert`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 61](61-db-insert-import-batch-summary.md)
- [Lesson 62](62-db-query-import-batch-summaries.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/63_db_upsert_import_batch_summary.lua](../../script/tutorials/63_db_upsert_import_batch_summary.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/63_db_upsert_import_batch_summary.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/63_db_upsert_import_batch_summary.lua
```

预期结果：

- 会生成 `artifacts/tutorials/63-db-upsert-import-batch-summary-lua.json`
- 同一个 `batch_id` 会先被插入，再被更新

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/63_db_upsert_import_batch_summary.flow.yaml](../../script/tutorials/63_db_upsert_import_batch_summary.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/63_db_upsert_import_batch_summary.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/63_db_upsert_import_batch_summary.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/63-db-upsert-import-batch-summary-flow.json`

## Step 3: 这节意味着什么

现在你已经有了数据库这条线里的三种基本动作：

- `insert`
- `query`
- `upsert`

这三步一旦走顺，后面再进入“批量明细入库”和“事务打包”就会自然很多。

## 下一节

下一节把批次摘要和明细行一起放进事务里：
[Lesson 64](64-db-transaction-import-batch-and-rows.md)
