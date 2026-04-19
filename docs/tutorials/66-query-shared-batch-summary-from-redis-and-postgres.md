# Lesson 66: 一次读回 Redis 和 Postgres 的共享批次摘要

上一节我们把最新 Redis 批次同步进了 Postgres。  
这一节先不写新数据，只做一件事：把同一个 `batch_id` 从两边都读回来。

目标：

- `redis_get`
- `json_extract`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 65](65-sync-latest-redis-batch-to-postgres-summary.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua](../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua
```

预期结果：

- 会生成 `artifacts/tutorials/66-query-shared-batch-summary-from-redis-and-postgres-lua.json`
- Lua 版本会额外校验 Redis、CSV、Postgres 三边的摘要数量是不是一致

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml](../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/66-query-shared-batch-summary-from-redis-and-postgres-flow.json`

## Step 3: 这节意味着什么

现在你已经可以确认：

- Redis 里的“最新批次”
- 本地 CSV 的来源文件
- Postgres 里的持久化摘要

三者到底是不是指向同一个业务批次。

## 下一节

下一节继续往前，把这个共享 `batch_id` 再扩展到明细行：
[Lesson 67](67-transaction-store-shared-batch-rows.md)
