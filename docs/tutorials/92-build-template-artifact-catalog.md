# Lesson 92: 把交接产物整理成模板目录

`Lesson 91` 已经把 manifest 里的每份产物都分了角色。

这一节再往前走一步：

- 给这些产物分配稳定的模板槽位
- 给每种槽位配一个环境变量入口
- 让后面写模板时，不用再反复猜路径和命名

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 91](91-read-handoff-manifest-roles.md)

如果前一节跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_ROLE_FILE=artifacts/tutorials/91-read-handoff-manifest-roles-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/92_build_template_artifact_catalog.lua](../../script/tutorials/92_build_template_artifact_catalog.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/92_build_template_artifact_catalog.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/92_build_template_artifact_catalog.lua
```

预期结果：

- 会生成 `artifacts/tutorials/92-build-template-artifact-catalog-lua.csv`
- 会生成 `artifacts/tutorials/92-build-template-artifact-catalog-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/92_build_template_artifact_catalog.flow.yaml](../../script/tutorials/92_build_template_artifact_catalog.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/92_build_template_artifact_catalog.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/92_build_template_artifact_catalog.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/92-build-template-artifact-catalog-flow.csv`
- 会生成 `artifacts/tutorials/92-build-template-artifact-catalog-flow.json`

## Step 3: 这节意味着什么

到这里，交接产物已经不再只是“跑完后留下的几个文件”。

它们已经开始变成：

- 稳定命名的模板槽位
- 可覆写路径的输入点
- 可以复用的模板资产

## 下一节

[Lesson 93](93-build-input-process-output-template.md)
