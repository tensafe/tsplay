# Lesson 86: 生成一份回放后的对账包

这一节开始把几份分散证据重新压成一份更容易交接的结果。

这里会把三类事实放到一起：

- 生命周期里记录的原始行数
- replay 批次现在的 Redis / Postgres 行数
- 原始 / replay 两边的审计数量

目标：

- `read_csv`
- `redis_get`
- `json_extract`
- `db_query_one`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 82](82-replay-batch-from-lifecycle-evidence.md)
- [Lesson 85](85-export-original-and-replay-audits.md)

如果你前面跑的是 Lua 版本，记得切路径：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
export TSPLAY_AUDIT_COMPARE_FILE=artifacts/tutorials/85-export-original-and-replay-audits-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/86_build_post_replay_reconciliation_pack.lua](../../script/tutorials/86_build_post_replay_reconciliation_pack.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/86_build_post_replay_reconciliation_pack.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/86_build_post_replay_reconciliation_pack.lua
```

预期结果：

- 会生成 `artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.csv`
- 会生成 `artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml](../../script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/86-build-post-replay-reconciliation-pack-flow.csv`
- 会生成 `artifacts/tutorials/86-build-post-replay-reconciliation-pack-flow.json`

## Step 3: 这节意味着什么

这一步开始进入“面向交接”的视角。

因为真正交付时，团队通常不会想看很多零散 JSON；更想先看到一份压缩后的对账结论。

## 下一节

[Lesson 87](87-build-handoff-artifact-manifest.md)
