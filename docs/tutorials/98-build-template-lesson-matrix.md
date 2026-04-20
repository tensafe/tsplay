# Lesson 98: 生成一份“场景 -> 模板”的学习矩阵

模板有了，索引也有了。

这一节再补一个更适合教学和选型的视角：

- 哪种场景
- 适合用哪份模板
- 从哪个产物开始
- 想得到什么输出

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 92](92-build-template-artifact-catalog.md)
- [Lesson 96](96-build-template-index.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_TEMPLATE_CATALOG_FILE=artifacts/tutorials/92-build-template-artifact-catalog-lua.csv
export TSPLAY_TEMPLATE_INDEX_FILE=artifacts/tutorials/96-build-template-index-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/98_build_template_lesson_matrix.lua](../../script/tutorials/98_build_template_lesson_matrix.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/98_build_template_lesson_matrix.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/98_build_template_lesson_matrix.lua
```

预期结果：

- 会生成 `artifacts/tutorials/98-build-template-lesson-matrix-lua.csv`
- 会生成 `artifacts/tutorials/98-build-template-lesson-matrix-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/98_build_template_lesson_matrix.flow.yaml](../../script/tutorials/98_build_template_lesson_matrix.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/98_build_template_lesson_matrix.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/98_build_template_lesson_matrix.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/98-build-template-lesson-matrix-flow.csv`
- 会生成 `artifacts/tutorials/98-build-template-lesson-matrix-flow.json`

## Step 3: 这节意味着什么

这一节开始让模板真正变成“可教、可挑、可解释”的学习资产。

以后新人不一定要先懂所有文件，只要先知道自己遇到的是哪种场景，就能更快找到入口。

## 下一节

[Lesson 99](99-build-template-preflight-checklist.md)
