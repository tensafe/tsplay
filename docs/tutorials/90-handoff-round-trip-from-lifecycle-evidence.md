# Lesson 90: 跑通一条“生命周期证据 -> 回放 -> 交接包”的完整 round trip

这一节是 `81-89` 这一整段的收口。

它会把下面这些动作重新串起来：

- 从 `Lesson 80` 的生命周期证据开始
- 新建一条 handoff replay 批次
- 写一条 handoff audit
- 产出最终交接摘要

目标：

- `read_csv`
- `redis_incr`
- `redis_set`
- `db_transaction`
- `db_upsert`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 80](80-external-sync-lifecycle-round-trip.md)
- [Lesson 89](89-build-pre-release-checklist.md)

如果你前一节跑的是 Lua 版本，先切一下生命周期文件：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua](../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua
```

预期结果：

- 会生成 `artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv`
- 会生成 `artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml](../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv`
- 会生成 `artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.json`

## Step 3: 这节意味着什么

到这里，中间这一整段已经形成了另一条完整主线：

- 生命周期结束
- 证据保留
- 证据回放
- 回放审计
- 对账压缩
- 交接整理
- 发布前检查

这条线比前一段更接近真实交付场景，因为它开始关心“怎么把一次运行变成可以交给别人继续用的结果”。

## 下一节

下一段会继续从这里往下接，开始把交接产物真正提炼成可复用模板：
[Lesson 91](91-read-handoff-manifest-roles.md)
