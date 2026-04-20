# Lesson 99: 给模板包生成发布前检查清单

前面 `89` 是给交接包做发布前检查。

这一节继续同样的思路，但对象换成模板包本身：

- 模板索引在不在
- 模板验证过没过
- 场景矩阵在不在

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 96](96-build-template-index.md)
- [Lesson 97](97-verify-template-covers-handoff-chain.md)
- [Lesson 98](98-build-template-lesson-matrix.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_TEMPLATE_INDEX_FILE=artifacts/tutorials/96-build-template-index-lua.csv
export TSPLAY_TEMPLATE_VERIFICATION_FILE=artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.csv
export TSPLAY_TEMPLATE_LESSON_MATRIX_FILE=artifacts/tutorials/98-build-template-lesson-matrix-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/99_build_template_preflight_checklist.lua](../../script/tutorials/99_build_template_preflight_checklist.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/99_build_template_preflight_checklist.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/99_build_template_preflight_checklist.lua
```

预期结果：

- 会生成 `artifacts/tutorials/99-build-template-preflight-checklist-lua.csv`
- 会生成 `artifacts/tutorials/99-build-template-preflight-checklist-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/99_build_template_preflight_checklist.flow.yaml](../../script/tutorials/99_build_template_preflight_checklist.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/99_build_template_preflight_checklist.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/99_build_template_preflight_checklist.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/99-build-template-preflight-checklist-flow.csv`
- 会生成 `artifacts/tutorials/99-build-template-preflight-checklist-flow.json`

## Step 3: 这节意味着什么

到这里，模板已经开始拥有自己的发布前门槛。

也就是说，我们不只是在发布业务结果，也是在发布“可复用的教程骨架”。

## 下一节

[Lesson 100](100-template-round-trip-from-handoff-artifacts.md)
