# Lesson 82: 按生命周期证据回放一个新批次

这一节正式把 `Lesson 80` 留下来的证据重新转回一个可运行批次。

动作很直接：

- 从生命周期 CSV 里拿到 `input_file`
- 再读一次原始导出 CSV
- 新建一个 replay 批次
- 把 Redis 和 Postgres 里的运行态重新建回来

目标：

- `read_csv`
- `redis_incr`
- `redis_set`
- `db_transaction`
- `db_upsert`
- `db_insert_many`

## 开始前

建议先跑完：

- [Lesson 80](80-external-sync-lifecycle-round-trip.md)
- [Lesson 81](81-read-lifecycle-evidence.md)

如果你前一节跑的是 Lua 版本，先切一下生命周期文件：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/82_replay_batch_from_lifecycle_evidence.lua](../../script/tutorials/82_replay_batch_from_lifecycle_evidence.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/82_replay_batch_from_lifecycle_evidence.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/82_replay_batch_from_lifecycle_evidence.lua
```

预期结果：

- 会生成 `artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv`
- 会生成 `artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml](../../script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.csv`
- 会生成 `artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.json`

## Step 3: 这节意味着什么

到这里，教程已经不只是“把一条业务流程跑一次”。

而是开始进入更真实的交付场景：

- 原批次已经结束
- 运行态已经被清理
- 但你仍然能根据证据把一条批次重新回放出来

## 下一节

[Lesson 83](83-verify-replay-batch-against-lifecycle-evidence.md)
