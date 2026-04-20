# Lesson 89: 生成发布前检查清单

到这一步，交接包已经有了。

但真正要交给别人之前，最好再多一步：

- 把关键产物整理成 checklist
- 明确哪些项已经 ready
- 哪些项如果缺了，就不应该往外发

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 87](87-build-handoff-artifact-manifest.md)
- [Lesson 88](88-build-handoff-summary.md)

如果前一节跑的是 Lua 版本，先切一下 manifest：

```bash
export TSPLAY_MANIFEST_FILE=artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/89_build_pre_release_checklist.lua](../../script/tutorials/89_build_pre_release_checklist.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/89_build_pre_release_checklist.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/89_build_pre_release_checklist.lua
```

预期结果：

- 会生成 `artifacts/tutorials/89-build-pre-release-checklist-lua.csv`
- 会生成 `artifacts/tutorials/89-build-pre-release-checklist-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/89_build_pre_release_checklist.flow.yaml](../../script/tutorials/89_build_pre_release_checklist.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/89_build_pre_release_checklist.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/89_build_pre_release_checklist.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/89-build-pre-release-checklist-flow.csv`
- 会生成 `artifacts/tutorials/89-build-pre-release-checklist-flow.json`

## Step 3: 这节意味着什么

这一节开始把“教程跑通”往“可以交付”再推一层。

因为交付里最怕的不是报错，而是少东西。

## 下一节

[Lesson 90](90-handoff-round-trip-from-lifecycle-evidence.md)
