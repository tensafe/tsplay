# Lesson 94: 把交接链整理成 Collect -> Verify -> Save 模板

`Lesson 93` 是按输入、处理、输出来拆。

这一节换一个更贴近 review 的视角：

- 先收集
- 再验证
- 最后保存

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 91](91-read-handoff-manifest-roles.md)
- [Lesson 89](89-build-pre-release-checklist.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_ROLE_FILE=artifacts/tutorials/91-read-handoff-manifest-roles-lua.csv
export TSPLAY_RUNTIME_CHECKLIST_FILE=artifacts/tutorials/89-build-pre-release-checklist-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/94_build_collect_verify_save_template.lua](../../script/tutorials/94_build_collect_verify_save_template.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/94_build_collect_verify_save_template.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/94_build_collect_verify_save_template.lua
```

预期结果：

- 会生成 `artifacts/tutorials/94-build-collect-verify-save-template-lua.csv`
- 会生成 `artifacts/tutorials/94-build-collect-verify-save-template-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/94_build_collect_verify_save_template.flow.yaml](../../script/tutorials/94_build_collect_verify_save_template.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/94_build_collect_verify_save_template.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/94_build_collect_verify_save_template.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/94-build-collect-verify-save-template-flow.csv`
- 会生成 `artifacts/tutorials/94-build-collect-verify-save-template-flow.json`

## Step 3: 这节意味着什么

这一节开始让模板不只服务“写流程”，也服务“review 流程”。

以后团队看这类模板时，就能更快对齐：

- 先看收集段
- 再看验证段
- 最后看保存段

## 下一节

[Lesson 95](95-build-replay-audit-handoff-template.md)
