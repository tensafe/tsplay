# Lesson 97: 验证模板索引仍然覆盖完整交接链

索引建出来以后，还需要确认一件事：

- 这些模板有没有把原来的交接链真的覆盖住

也就是要确认：

- handoff 结果本身还是 ok
- 三份核心模板都还在
- 模板数量没有掉

目标：

- `read_csv`
- `write_csv`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 96](96-build-template-index.md)
- [Lesson 90](90-handoff-round-trip-from-lifecycle-evidence.md)

如果前面跑的是 Lua 版本，先切一下：

```bash
export TSPLAY_TEMPLATE_INDEX_FILE=artifacts/tutorials/96-build-template-index-lua.csv
export TSPLAY_HANDOFF_FILE=artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/97_verify_template_covers_handoff_chain.lua](../../script/tutorials/97_verify_template_covers_handoff_chain.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/97_verify_template_covers_handoff_chain.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/97_verify_template_covers_handoff_chain.lua
```

预期结果：

- 会生成 `artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.csv`
- 会生成 `artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml](../../script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/97-verify-template-covers-handoff-chain-flow.csv`
- 会生成 `artifacts/tutorials/97-verify-template-covers-handoff-chain-flow.json`

## Step 3: 这节意味着什么

这一节是在练一种很关键的习惯：

- 模板不是建出来就完
- 还要确认它确实没有把原链路丢掉

## 下一节

[Lesson 98](98-build-template-lesson-matrix.md)
