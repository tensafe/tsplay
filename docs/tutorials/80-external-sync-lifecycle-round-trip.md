# Lesson 80: 跑通一条完整的外部同步生命周期

这一节是 `72-79` 这一整段的收口。

它把几个分散动作重新串起来：

- 新建一条外部同步批次
- 写入 Redis 和 Postgres
- 追加一条审计记录
- 清理运行态数据
- 再确认审计仍然存在

目标：

- `redis_incr`
- `redis_set`
- `db_transaction`
- `db_upsert`
- `redis_del`
- `db_execute`

## 开始前

建议先跑完：

- [Lesson 71](71-external-system-round-trip.md)
- [Lesson 79](79-verify-external-batch-cleanup.md)

这节需要 Redis 和 Postgres 都准备好，并且审计表已经创建好：

```bash
source script/tutorials/env/06_redis_example.sh
source script/tutorials/env/07_reporting_pg_example.sh
psql "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable" \
  -f script/tutorials/sql/75_reporting_import_audit.sql
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/80_external_sync_lifecycle_round_trip.lua](../../script/tutorials/80_external_sync_lifecycle_round_trip.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/80_external_sync_lifecycle_round_trip.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/80_external_sync_lifecycle_round_trip.lua
```

预期结果：

- 会生成 `artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv`
- 会生成 `artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.json`

现在这份 CSV 里还会带上 `input_file` 和 `payload_key`，方便后面继续做“证据回放”和“交接包”这一段。

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml](../../script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv`
- 会生成 `artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.json`

如果你后面准备接 `Lesson 81+`，建议优先跑 Flow 版本；因为后续默认会直接复用这份
`80-external-sync-lifecycle-round-trip-flow.csv`。

## Step 3: 这节意味着什么

到这里，外部系统这一整段已经完整形成了“创建、重跑、异常、审计、清理、验证”的生命周期。

这也意味着初级阶段已经不只是“会接一个 Redis / DB 示例”，而是真的能组织一条小型业务结果流。

下一段会继续从这里往下接：不是重新造一条新链，而是开始学“生命周期结束后，怎么根据留下来的证据再回放、再比对、再交接”。

## 下一节

[Lesson 81](81-read-lifecycle-evidence.md)
