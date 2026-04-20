# Lesson 93: 把交接链整理成 Input -> Process -> Output 模板

这一节开始第一次真正“抽模板”。

不是再关心某个具体批次，而是把整条链拆成三段：

- 输入
- 处理
- 输出

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
[../../script/tutorials/93_build_input_process_output_template.lua](../../script/tutorials/93_build_input_process_output_template.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/93_build_input_process_output_template.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/93_build_input_process_output_template.lua
```

预期结果：

- 会生成 `artifacts/tutorials/93-build-input-process-output-template-lua.csv`
- 会生成 `artifacts/tutorials/93-build-input-process-output-template-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/93_build_input_process_output_template.flow.yaml](../../script/tutorials/93_build_input_process_output_template.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/93_build_input_process_output_template.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/93_build_input_process_output_template.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/93-build-input-process-output-template-flow.csv`
- 会生成 `artifacts/tutorials/93-build-input-process-output-template-flow.json`

## Step 3: 这节意味着什么

这一节完成以后，你已经不只是“知道有哪些文件”。

而是开始知道：

- 哪个文件先出现
- 哪个文件是处理中间态
- 哪个文件是应该留给交付方的最终产物

## 下一节

[Lesson 94](94-build-collect-verify-save-template.md)
