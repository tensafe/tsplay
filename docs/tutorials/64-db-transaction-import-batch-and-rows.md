# Lesson 64: 在一个事务里写入批次摘要和明细行

这一节是这一小段数据库主线的收口。

前面我们已经有了：

- 一条批次摘要
- 多条批次查询
- 已存在记录的更新

现在把它们进一步组合成一个更像真实业务交付的动作：

- 先整理导出 CSV 的明细行
- 再用一个事务写入批次摘要
- 同时批量写入明细行
- 最后把两张表再查回来

目标：

- `read_csv`
- `db_transaction`
- `db_insert_many`
- `db_query`

## 开始前

建议先跑完：

- [Lesson 61](61-db-insert-import-batch-summary.md)
- [Lesson 63](63-db-upsert-import-batch-summary.md)

如果你换过终端，记得重新加载：

```bash
source script/tutorials/env/07_reporting_pg_example.sh
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/64_db_transaction_import_batch_and_rows.lua](../../script/tutorials/64_db_transaction_import_batch_and_rows.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/64_db_transaction_import_batch_and_rows.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/64_db_transaction_import_batch_and_rows.lua
```

预期结果：

- 会生成 `artifacts/tutorials/64-db-transaction-import-batch-and-rows-lua.json`
- `public.tutorial_import_batches` 和 `public.tutorial_import_rows` 都会有同一个 `batch_id` 的数据

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml](../../script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/64-db-transaction-import-batch-and-rows-flow.json`

## Step 3: 这节意味着什么

到这里，`Lesson 57` 那份认证导出结果已经可以沿着一条很顺的主线流动：

- 浏览器页面导出
- 本地 CSV 回读
- Redis 摘要和批次 key
- Postgres 摘要持久化
- Postgres 明细批量入库

这一段结束后，教程已经从“浏览器自动化”自然进入了“业务结果同步与持久化”。

## 下一节

下一节先从最轻的一步开始，把“最新 Redis 批次摘要”同步成一条共享的 Postgres 摘要：
[Lesson 65](65-sync-latest-redis-batch-to-postgres-summary.md)
