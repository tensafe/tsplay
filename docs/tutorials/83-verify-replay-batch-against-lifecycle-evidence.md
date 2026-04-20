# Lesson 83: 用生命周期证据验证回放批次

`Lesson 82` 已经把回放批次重新建出来了。

这一节不继续加新动作，而是做一轮比对：

- 生命周期里记录的原始行数
- 回放 payload 里的行数
- Postgres 摘要里的行数
- Postgres 明细里的行数

要确认这几份事实仍然一致。

目标：

- `read_csv`
- `redis_get`
- `json_extract`
- `db_query_one`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 80](80-external-sync-lifecycle-round-trip.md)
- [Lesson 82](82-replay-batch-from-lifecycle-evidence.md)

如果前一节跑的是 Lua 版本，记得切一下回放文件：

```bash
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
```

如果生命周期文件来自 Lua，也一并切一下：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua](../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua
```

预期结果：

- 会生成 `artifacts/tutorials/83-verify-replay-batch-against-lifecycle-evidence-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml](../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/83-verify-replay-batch-against-lifecycle-evidence-flow.json`

## Step 3: 这节意味着什么

这一节开始明确一件事：

回放不是“重新写进去就算完”，而是要能说明：

- 为什么这次回放是可信的
- 它和原始生命周期证据到底有没有偏差

## 下一节

[Lesson 84](84-write-replay-audit-row.md)
