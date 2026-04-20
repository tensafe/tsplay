# Lesson 71: 跑通一次完整的外部系统 round trip

这一节是 `58-70` 这一整段的收口。

它不再依赖你提前跑好某一个中间 lesson，而是把整条外部系统主线重新串起来：

- 从本地导出 CSV 出发
- 生成新的 Redis 批次
- 保存最新批次指针
- 用同一个 `batch_id` 持久化 Postgres 摘要和明细
- 最后把结果重新读回来

目标：

- `read_csv`
- `redis_incr`
- `redis_set`
- `db_transaction`
- `db_upsert`
- `db_insert_many`

## 开始前

建议先跑完：

- [Lesson 57](57-use-session-import-export-round-trip.md)
- [Lesson 70](70-build-reconciliation-pack-from-csv-redis-db.md)

这节需要 Redis 和 Postgres 都准备好：

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
[../../script/tutorials/71_external_system_round_trip.lua](../../script/tutorials/71_external_system_round_trip.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/71_external_system_round_trip.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/71_external_system_round_trip.lua
```

预期结果：

- 会生成 `artifacts/tutorials/71-external-system-round-trip-lua.csv`
- 会生成 `artifacts/tutorials/71-external-system-round-trip-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/71_external_system_round_trip.flow.yaml](../../script/tutorials/71_external_system_round_trip.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/71_external_system_round_trip.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/71_external_system_round_trip.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/71-external-system-round-trip-flow.csv`
- 会生成 `artifacts/tutorials/71-external-system-round-trip-flow.json`

## Step 3: 这节意味着什么

到这里，教程已经把一条非常完整的业务链路跑通了：

- 浏览器页面产生结果
- 本地文件保存结果
- Redis 管批次和最新指针
- Postgres 存摘要和明细
- 本地对账包负责复盘

这一条线已经很接近真实交付里的“最小结果同步流程”。

## 下一节

下一节先从运维视角继续往前走，练习“同一个批次号重跑”：
[Lesson 72](72-rerun-shared-batch-idempotently.md)
