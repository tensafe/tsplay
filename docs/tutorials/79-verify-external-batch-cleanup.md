# Lesson 79: 验证批次清理后，审计仍然保留

这一节是清理动作的复盘。

重点不是“删了没删”，而是：

- Redis payload 是不是没了
- Postgres 摘要和明细是不是都没了
- 审计表里的记录是不是还在

目标：

- `read_csv`
- `redis_get`
- `db_query`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 78](78-cleanup-latest-external-batch.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/79_verify_external_batch_cleanup.lua](../../script/tutorials/79_verify_external_batch_cleanup.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/79_verify_external_batch_cleanup.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/79_verify_external_batch_cleanup.lua
```

预期结果：

- 会生成 `artifacts/tutorials/79-verify-external-batch-cleanup-lua.json`

如果你前一节跑的是 Lua 版本，记得先切一下清理记录文件：

```bash
export TSPLAY_CLEANUP_FILE=artifacts/tutorials/78-cleanup-latest-external-batch-lua.csv
```

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/79_verify_external_batch_cleanup.flow.yaml](../../script/tutorials/79_verify_external_batch_cleanup.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/79_verify_external_batch_cleanup.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/79_verify_external_batch_cleanup.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/79-verify-external-batch-cleanup-flow.json`

## 下一节

最后一节把“生成批次、写审计、清理运行数据、验证审计保留”重新串成一条完整生命周期：
[Lesson 80](80-external-sync-lifecycle-round-trip.md)
