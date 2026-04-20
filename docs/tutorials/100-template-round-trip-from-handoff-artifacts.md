# Lesson 100: 跑通一条“交接产物 -> 模板包”的完整 round trip

这一节是 `91-99` 这一整段的收口。

它会重新把这些动作串起来：

- 从 handoff 产物开始
- 整理模板目录
- 建模板索引
- 生成模板检查清单
- 最后汇总成一份模板包结果

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 92](92-build-template-artifact-catalog.md)
- [Lesson 96](96-build-template-index.md)
- [Lesson 99](99-build-template-preflight-checklist.md)
- [Lesson 90](90-handoff-round-trip-from-lifecycle-evidence.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_TEMPLATE_CATALOG_FILE=artifacts/tutorials/92-build-template-artifact-catalog-lua.csv
export TSPLAY_TEMPLATE_INDEX_FILE=artifacts/tutorials/96-build-template-index-lua.csv
export TSPLAY_TEMPLATE_PREFLIGHT_FILE=artifacts/tutorials/99-build-template-preflight-checklist-lua.csv
export TSPLAY_HANDOFF_FILE=artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/100_template_round_trip_from_handoff_artifacts.lua](../../script/tutorials/100_template_round_trip_from_handoff_artifacts.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/100_template_round_trip_from_handoff_artifacts.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/100_template_round_trip_from_handoff_artifacts.lua
```

预期结果：

- 会生成 `artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-lua.csv`
- 会生成 `artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml](../../script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-flow.csv`
- 会生成 `artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-flow.json`

## Step 3: 这节意味着什么

到这里，教程又形成了一条新的完整主线：

- 先把业务链跑通
- 再把交接包做出来
- 再把交接包整理成模板资产
- 最后把模板资产也做成可交付结果

下一段就很自然会进入更细的“健壮性、等待、恢复和可维护性”专题。

## 下一节

[Lesson 101](101-assert-visible-template-release-card.md)
