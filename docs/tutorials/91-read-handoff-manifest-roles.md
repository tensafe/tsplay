# Lesson 91: 读交接 manifest，识别每份产物的角色

`Lesson 87-90` 已经把交接包做出来了。

这一节开始进入下一段主线：

- 不再继续补业务动作
- 而是把已经产出的交接文件重新整理
- 先看清楚每份文件在整条链里扮演什么角色

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 87](87-build-handoff-artifact-manifest.md)
- [Lesson 90](90-handoff-round-trip-from-lifecycle-evidence.md)

如果你前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_MANIFEST_FILE=artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv
export TSPLAY_HANDOFF_FILE=artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/91_read_handoff_manifest_roles.lua](../../script/tutorials/91_read_handoff_manifest_roles.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/91_read_handoff_manifest_roles.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/91_read_handoff_manifest_roles.lua
```

预期结果：

- 会生成 `artifacts/tutorials/91-read-handoff-manifest-roles-lua.csv`
- 会生成 `artifacts/tutorials/91-read-handoff-manifest-roles-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/91_read_handoff_manifest_roles.flow.yaml](../../script/tutorials/91_read_handoff_manifest_roles.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/91_read_handoff_manifest_roles.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/91_read_handoff_manifest_roles.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/91-read-handoff-manifest-roles-flow.csv`
- 会生成 `artifacts/tutorials/91-read-handoff-manifest-roles-flow.json`

## Step 3: 这节意味着什么

这一节是在给后面的模板化做准备。

因为只有先分清楚：

- 哪些是输入证据
- 哪些是中间结果
- 哪些是核对材料
- 哪些是最终交接产物

后面才能把它们提炼成可复用模板。

## 下一节

[Lesson 92](92-build-template-artifact-catalog.md)
