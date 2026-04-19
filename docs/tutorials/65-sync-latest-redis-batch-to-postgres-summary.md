# Lesson 65: 把最新 Redis 批次摘要同步到 Postgres

上一段我们已经有了两套外部系统能力：

- `Lesson 58-60`：Redis 里有了批次摘要和最新批次指针
- `Lesson 61-64`：Postgres 里能写摘要、写明细、做事务

这一节开始把两边真正接起来。  
先做最轻的一步：把“最新 Redis 批次摘要”同步成一条 Postgres 摘要记录。

目标：

- `redis_get`
- `json_extract`
- `read_csv`
- `db_upsert`

## 开始前

建议先跑完：

- [Lesson 59](59-save-import-batch-key-to-redis.md)
- [Lesson 61](61-db-insert-import-batch-summary.md)

这节同时需要 Redis 和 Postgres 环境：

```bash
source script/tutorials/env/06_redis_example.sh
source script/tutorials/env/07_reporting_pg_example.sh
```

如果你还没执行过导入同步用的建表 SQL，先跑一次：

```bash
psql "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable" \
  -f script/tutorials/sql/61_reporting_import_sync.sql
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua](../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua
```

预期结果：

- 会生成 `artifacts/tutorials/65-sync-latest-redis-batch-to-postgres-summary-lua.json`
- 最新 Redis 批次会在 `public.tutorial_import_batches` 里出现同一个 `batch_id`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml](../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/65-sync-latest-redis-batch-to-postgres-summary-flow.json`

## Step 3: 这节意味着什么

到这里，Redis 不再只是临时缓存。  
它已经开始承担“给数据库提供同步来源”的角色。

## 下一节

下一节继续用同一个 `batch_id`，把 Redis 和 Postgres 两边的摘要放到一起核对：
[Lesson 66](66-query-shared-batch-summary-from-redis-and-postgres.md)
