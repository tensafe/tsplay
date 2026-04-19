# Lesson 62: 查询多条 Postgres 批次摘要

上一节我们只查回一条摘要。  
这一节继续往前走，把它扩成“多条结果列表”：

- 先准备两条批次摘要
- 再用 `db_query` 一次查回多行

目标：

- `db_insert`
- `db_query`

## 开始前

建议先跑完：

- [Lesson 61](61-db-insert-import-batch-summary.md)

这节继续复用同一套数据库连接和同一张表。  
如果你换过终端，记得重新执行：

```bash
source script/tutorials/env/07_reporting_pg_example.sh
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/62_db_query_import_batch_summaries.lua](../../script/tutorials/62_db_query_import_batch_summaries.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/62_db_query_import_batch_summaries.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/62_db_query_import_batch_summaries.lua
```

预期结果：

- 会生成 `artifacts/tutorials/62-db-query-import-batch-summaries-lua.json`
- 结果里会有 `batch_rows`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/62_db_query_import_batch_summaries.flow.yaml](../../script/tutorials/62_db_query_import_batch_summaries.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/62_db_query_import_batch_summaries.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/62_db_query_import_batch_summaries.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/62-db-query-import-batch-summaries-flow.json`

## Step 3: 这节意味着什么

到这里，数据库这一段不再只是“写一条再读一条”，而是开始进入真正的列表查询视角。  
这一步很重要，因为后面做批次管理、审计、回查时，基本都离不开 `db_query`。

## 下一节

下一节继续补上“已存在时更新”的能力：
[Lesson 63](63-db-upsert-import-batch-summary.md)
