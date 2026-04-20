# Lesson 84: 给回放批次补写一条审计记录

`Lesson 82-83` 已经把回放批次重新建起来，并且做了第一轮核对。

这一节再往前走一步：

- 给 replay 批次补一条单独的审计记录
- 把“它是从哪条原始批次回放来的”写清楚

目标：

- `db_query_one`
- `db_upsert`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 82](82-replay-batch-from-lifecycle-evidence.md)
- [Lesson 83](83-verify-replay-batch-against-lifecycle-evidence.md)

如果前一节跑的是 Lua 版本，记得先切：

```bash
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/84_write_replay_audit_row.lua](../../script/tutorials/84_write_replay_audit_row.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/84_write_replay_audit_row.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/84_write_replay_audit_row.lua
```

预期结果：

- 会生成 `artifacts/tutorials/84-write-replay-audit-row-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/84_write_replay_audit_row.flow.yaml](../../script/tutorials/84_write_replay_audit_row.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/84_write_replay_audit_row.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/84_write_replay_audit_row.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/84-write-replay-audit-row-flow.json`

## Step 3: 这节意味着什么

到这里，“原始生命周期”和“回放生命周期”就不再混在一起了。

你已经开始有能力把两条链分别留痕：

- 原链有原链的审计
- 回放链有回放链的审计

## 下一节

[Lesson 85](85-export-original-and-replay-audits.md)
