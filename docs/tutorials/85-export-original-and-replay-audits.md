# Lesson 85: 把原批次和回放批次的审计导出成对照 CSV

这节继续沿着审计线往下走。

既然原批次和 replay 批次都已经有自己的审计记录了，下一步最自然就是：

- 把两边的审计一起读出来
- 放进一份对照 CSV
- 方便人直接复盘

目标：

- `db_query`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 80](80-external-sync-lifecycle-round-trip.md)
- [Lesson 84](84-write-replay-audit-row.md)

如果你前面跑的是 Lua 版本，先切一下文件：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/85_export_original_and_replay_audits.lua](../../script/tutorials/85_export_original_and_replay_audits.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/85_export_original_and_replay_audits.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/85_export_original_and_replay_audits.lua
```

预期结果：

- 会生成 `artifacts/tutorials/85-export-original-and-replay-audits-lua.csv`
- 会生成 `artifacts/tutorials/85-export-original-and-replay-audits-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/85_export_original_and_replay_audits.flow.yaml](../../script/tutorials/85_export_original_and_replay_audits.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/85_export_original_and_replay_audits.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/85_export_original_and_replay_audits.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/85-export-original-and-replay-audits-flow.csv`
- 会生成 `artifacts/tutorials/85-export-original-and-replay-audits-flow.json`

## Step 3: 这节意味着什么

这一节把“机器里的审计”正式转成了“人能看的对照表”。

后面做对账包和交接 manifest，就都会继续复用这份审计对照结果。

## 下一节

[Lesson 86](86-build-post-replay-reconciliation-pack.md)
