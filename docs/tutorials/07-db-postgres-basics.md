# Lesson 07: Postgres 基础查询与写入

这一节用一个最小 Postgres 例子，把数据库动作串起来：

- `db_execute`
- `db_insert`
- `db_query_one`

为了让命令更直接，这一节默认使用命名连接 `reporting` 和 `pgsql` 驱动。

## 先准备 Postgres

如果你本机已经有 Postgres，可以直接复用。  
如果只是为了跟教程跑通一遍，下面这个命令比较省事：

```bash
docker run --name tsplay-pg \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_USER=collector \
  -e POSTGRES_DB=analytics \
  -p 5432:5432 \
  -d postgres:16
```

然后在仓库根目录加载示例环境变量：

```bash
source script/tutorials/env/07_reporting_pg_example.sh
```

这个示例文件会设置：

```bash
export TSPLAY_DB_REPORTING_DRIVER=pgsql
export TSPLAY_DB_REPORTING_DSN=postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable
```

接着执行初始化 SQL：

```bash
psql "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable" \
  -f script/tutorials/sql/07_reporting_pg.sql
```

如果你没有 `psql`，也可以用任何熟悉的 Postgres 客户端执行同一个 SQL 文件。

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/07_db_postgres_basics.lua](../../script/tutorials/07_db_postgres_basics.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/07_db_postgres_basics.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/07_db_postgres_basics.lua
```

这份脚本会先删掉旧数据，再插入一行固定记录，最后查回这一行。

预期结果：

- 会生成 `artifacts/tutorials/07-db-postgres-basics-lua.json`
- 输出里会包含 `cleanup`、`insert_result` 和 `row`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/07_db_postgres_basics.flow.yaml](../../script/tutorials/07_db_postgres_basics.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/07_db_postgres_basics.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/07_db_postgres_basics.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/07-db-postgres-basics-flow.json`
- 终端会显示每一步数据库动作的结构化结果

## Step 3: 这节要特别注意什么

- 这份教程默认是 Postgres，所以参数占位符写成 `$1`
- 如果你换成 MySQL / SQL Server / Oracle，要同时调整 driver、DSN 和 SQL 占位符风格
- SQL Server / Oracle 还可能需要额外的 build tags

补充说明：

- 本地直接跑 `go run . -flow ...` 不需要显式加 `allow_database`
- 如果以后改成 MCP 模式运行，再补 `allow_database=true`

## 接下来可以继续补的章节

这一节之后，最自然的下一段不是凭空再讲数据库动作，而是把认证导出结果真正接进数据库：

- [Lesson 61: 把认证导出结果写成一条 Postgres 批次摘要](61-db-insert-import-batch-summary.md)
- [Lesson 62: 查询多条 Postgres 批次摘要](62-db-query-import-batch-summaries.md)
- [Lesson 63: 用 `db_upsert` 更新 Postgres 批次摘要](63-db-upsert-import-batch-summary.md)
- [Lesson 64: 在一个事务里写入批次摘要和明细行](64-db-transaction-import-batch-and-rows.md)
- [Lesson 65: 把最新 Redis 批次摘要同步到 Postgres](65-sync-latest-redis-batch-to-postgres-summary.md)
- [Lesson 66: 一次读回 Redis 和 Postgres 的共享批次摘要](66-query-shared-batch-summary-from-redis-and-postgres.md)
- [Lesson 67: 用共享批次号把明细行写入 Postgres](67-transaction-store-shared-batch-rows.md)
- [Lesson 68: 读回共享批次的 Postgres 明细行](68-query-shared-batch-detail-rows.md)
- [Lesson 69: 把源 CSV 和 DB 明细行放到一起比](69-compare-source-csv-and-db-rows.md)
- [Lesson 70: 生成一份 CSV、Redis、Postgres 三边对账包](70-build-reconciliation-pack-from-csv-redis-db.md)
- [Lesson 71: 跑通一次完整的外部系统 round trip](71-external-system-round-trip.md)
