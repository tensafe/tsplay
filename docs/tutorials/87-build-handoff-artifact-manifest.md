# Lesson 87: 生成一份交接 artifact manifest

上一节已经有了压缩后的对账包。

这一节继续往“交接”推进：

- 把生命周期证据
- 回放结果
- 审计对照
- 对账包

四份产物放进一份 manifest 里，明确告诉别人“这次交接包里到底有什么”。

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 85](85-export-original-and-replay-audits.md)
- [Lesson 86](86-build-post-replay-reconciliation-pack.md)

如果你前面跑的是 Lua 版本，记得切：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
export TSPLAY_AUDIT_COMPARE_FILE=artifacts/tutorials/85-export-original-and-replay-audits-lua.csv
export TSPLAY_RECONCILIATION_FILE=artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/87_build_handoff_artifact_manifest.lua](../../script/tutorials/87_build_handoff_artifact_manifest.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/87_build_handoff_artifact_manifest.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/87_build_handoff_artifact_manifest.lua
```

预期结果：

- 会生成 `artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv`
- 会生成 `artifacts/tutorials/87-build-handoff-artifact-manifest-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/87_build_handoff_artifact_manifest.flow.yaml](../../script/tutorials/87_build_handoff_artifact_manifest.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/87_build_handoff_artifact_manifest.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/87_build_handoff_artifact_manifest.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/87-build-handoff-artifact-manifest-flow.csv`
- 会生成 `artifacts/tutorials/87-build-handoff-artifact-manifest-flow.json`

## Step 3: 这节意味着什么

到这里，你已经不只是“有几份结果文件”，而是开始把它们组织成交付物。

这一步对新人非常重要，因为以后交接别人时，对方首先要看的就是 manifest。

## 下一节

[Lesson 88](88-build-handoff-summary.md)
