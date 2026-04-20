# Lesson 96: 把几份模板整理成统一索引

`Lesson 93-95` 已经做出了三种不同视角的模板。

这一节不再补新模板，而是先给模板做一份统一索引：

- 它叫什么
- 适合解决什么问题
- 对应哪份源文件

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 93](93-build-input-process-output-template.md)
- [Lesson 94](94-build-collect-verify-save-template.md)
- [Lesson 95](95-build-replay-audit-handoff-template.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_INPUT_PROCESS_OUTPUT_FILE=artifacts/tutorials/93-build-input-process-output-template-lua.csv
export TSPLAY_COLLECT_VERIFY_SAVE_FILE=artifacts/tutorials/94-build-collect-verify-save-template-lua.csv
export TSPLAY_REPLAY_AUDIT_HANDOFF_FILE=artifacts/tutorials/95-build-replay-audit-handoff-template-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/96_build_template_index.lua](../../script/tutorials/96_build_template_index.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/96_build_template_index.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/96_build_template_index.lua
```

预期结果：

- 会生成 `artifacts/tutorials/96-build-template-index-lua.csv`
- 会生成 `artifacts/tutorials/96-build-template-index-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/96_build_template_index.flow.yaml](../../script/tutorials/96_build_template_index.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/96_build_template_index.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/96_build_template_index.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/96-build-template-index-flow.csv`
- 会生成 `artifacts/tutorials/96-build-template-index-flow.json`

## Step 3: 这节意味着什么

这一节完成以后，模板已经从“几份分散文件”变成了“一个可以浏览的模板目录”。

这一步对后面的验证和发布前检查很关键。

## 下一节

[Lesson 97](97-verify-template-covers-handoff-chain.md)
