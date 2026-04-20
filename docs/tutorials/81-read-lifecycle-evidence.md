# Lesson 81: 从生命周期 CSV 里读回批次证据

这一节不新建批次，也不重新写数据。

先做一件更稳的事：

- 直接读取 `Lesson 80` 留下来的生命周期 CSV
- 重新拿到原始 `batch_id`
- 再去 Postgres 里把保留下来的审计记录读回来

目标：

- `read_csv`
- `db_query_one`
- `write_json`

## 开始前

建议先跑完：

- [Lesson 80](80-external-sync-lifecycle-round-trip.md)

如果你前一节跑的是 Lua 版本，先切一下生命周期文件：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/81_read_lifecycle_evidence.lua](../../script/tutorials/81_read_lifecycle_evidence.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/81_read_lifecycle_evidence.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/81_read_lifecycle_evidence.lua
```

预期结果：

- 会生成 `artifacts/tutorials/81-read-lifecycle-evidence-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/81_read_lifecycle_evidence.flow.yaml](../../script/tutorials/81_read_lifecycle_evidence.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/81_read_lifecycle_evidence.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/81_read_lifecycle_evidence.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/81-read-lifecycle-evidence-flow.json`

## Step 3: 这节意味着什么

从这一节开始，主线进入“生命周期结束以后还剩下什么”。

也就是：

- 运行态数据可以清理
- 但证据和审计还在
- 后面的回放和交接，都会基于这些留下来的证据继续做

## 下一节

[Lesson 82](82-replay-batch-from-lifecycle-evidence.md)
