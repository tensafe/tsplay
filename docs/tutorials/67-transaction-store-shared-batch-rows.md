# Lesson 67: 用共享批次号把明细行写入 Postgres

上一节我们确认了 Redis 和 Postgres 的摘要已经对上同一个 `batch_id`。  
这一节继续往前，不再用 lesson 自己造的 batch 名，而是直接复用 Redis 的共享批次号，把明细行也写进去。

目标：

- `redis_get`
- `read_csv`
- `db_transaction`
- `db_upsert`
- `db_insert_many`

## 开始前

建议先跑完：

- [Lesson 65](65-sync-latest-redis-batch-to-postgres-summary.md)
- [Lesson 66](66-query-shared-batch-summary-from-redis-and-postgres.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/67_transaction_store_shared_batch_rows.lua](../../script/tutorials/67_transaction_store_shared_batch_rows.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/67_transaction_store_shared_batch_rows.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/67_transaction_store_shared_batch_rows.lua
```

预期结果：

- 会生成 `artifacts/tutorials/67-transaction-store-shared-batch-rows-lua.json`
- 共享 `batch_id` 会同时出现在 `tutorial_import_batches` 和 `tutorial_import_rows`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml](../../script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/67-transaction-store-shared-batch-rows-flow.json`

## Step 3: 这节意味着什么

到这里，Redis 已经不仅给了一个摘要，还给了一个真正可以跨系统复用的 batch 主键。  
这一步很关键，因为真实业务里“跨系统共用一套 batch id”比“每个系统自己起名”稳定得多。

## 下一节

下一节先把这些共享明细行完整读回来：
[Lesson 68](68-query-shared-batch-detail-rows.md)
