# Lesson 88: 把交接 manifest 整理成交付摘要

manifest 更像“完整目录”。

这一节再补一层更短的摘要：

- 这次交接包一共有几项
- 涉及哪几个 artifact key
- 对应哪些文件路径

目标：

- `read_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 87](87-build-handoff-artifact-manifest.md)

如果前一节跑的是 Lua 版本，先切一下 manifest：

```bash
export TSPLAY_MANIFEST_FILE=artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/88_build_handoff_summary.lua](../../script/tutorials/88_build_handoff_summary.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/88_build_handoff_summary.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/88_build_handoff_summary.lua
```

预期结果：

- 会生成 `artifacts/tutorials/88-build-handoff-summary-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/88_build_handoff_summary.flow.yaml](../../script/tutorials/88_build_handoff_summary.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/88_build_handoff_summary.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/88_build_handoff_summary.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/88-build-handoff-summary-flow.json`

## Step 3: 这节意味着什么

这一步是在练一种很常见的交付习惯：

- manifest 留给追细节的人
- summary 留给先扫全局的人

## 下一节

[Lesson 89](89-build-pre-release-checklist.md)
