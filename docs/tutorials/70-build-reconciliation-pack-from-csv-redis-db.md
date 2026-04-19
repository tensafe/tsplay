# Lesson 70: 生成一份 CSV、Redis、Postgres 三边对账包

前一节只比较了源 CSV 和 DB 明细。  
这一节再把 Redis 摘要也加进来，形成一份真正的三边对账结果。

目标：

- `read_csv`
- `redis_get`
- `db_query_one`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 66](66-query-shared-batch-summary-from-redis-and-postgres.md)
- [Lesson 69](69-compare-source-csv-and-db-rows.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua](../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua
```

预期结果：

- 会生成 `artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-lua.csv`
- 会生成 `artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml](../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-flow.csv`
- 会生成 `artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-flow.json`

## Step 3: 这节意味着什么

到这里，你已经不是单纯地“把结果写到几个地方”，而是在做一份可以交给别人复盘的外部系统对账包。

## 下一节

最后一节把整段外部系统同步线真正串成一次完整 round trip：
[Lesson 71](71-external-system-round-trip.md)
