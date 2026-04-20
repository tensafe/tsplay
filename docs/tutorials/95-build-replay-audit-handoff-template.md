# Lesson 95: 把交接链整理成 Replay -> Audit -> Handoff 模板

到这一步，我们已经有两种模板视角了：

- `Input -> Process -> Output`
- `Collect -> Verify -> Save`

这一节继续保留业务语义，再做一份更贴近交付现场的模板：

- Replay
- Audit
- Handoff

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 92](92-build-template-artifact-catalog.md)
- [Lesson 90](90-handoff-round-trip-from-lifecycle-evidence.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_TEMPLATE_CATALOG_FILE=artifacts/tutorials/92-build-template-artifact-catalog-lua.csv
export TSPLAY_HANDOFF_FILE=artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/95_build_replay_audit_handoff_template.lua](../../script/tutorials/95_build_replay_audit_handoff_template.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/95_build_replay_audit_handoff_template.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/95_build_replay_audit_handoff_template.lua
```

预期结果：

- 会生成 `artifacts/tutorials/95-build-replay-audit-handoff-template-lua.csv`
- 会生成 `artifacts/tutorials/95-build-replay-audit-handoff-template-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/95_build_replay_audit_handoff_template.flow.yaml](../../script/tutorials/95_build_replay_audit_handoff_template.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/95_build_replay_audit_handoff_template.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/95_build_replay_audit_handoff_template.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/95-build-replay-audit-handoff-template-flow.csv`
- 会生成 `artifacts/tutorials/95-build-replay-audit-handoff-template-flow.json`

## Step 3: 这节意味着什么

这一节把“模板”重新拉回了业务语言。

也就是说，以后你不一定总是从技术动作切入，也可以直接从：

- 回放
- 审计
- 交接

这三个业务块来组织整条链。

## 下一节

[Lesson 96](96-build-template-index.md)
