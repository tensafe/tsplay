# Lesson 68: 读回共享批次的 Postgres 明细行

上一节已经把共享 `batch_id` 的明细行写入 Postgres。  
这一节只做读取，先把“明细事实”单独看清楚。

目标：

- `redis_get`
- `db_query`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 67](67-transaction-store-shared-batch-rows.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/68_query_shared_batch_detail_rows.lua](../../script/tutorials/68_query_shared_batch_detail_rows.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/68_query_shared_batch_detail_rows.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/68_query_shared_batch_detail_rows.lua
```

预期结果：

- 会生成 `artifacts/tutorials/68-query-shared-batch-detail-rows-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/68_query_shared_batch_detail_rows.flow.yaml](../../script/tutorials/68_query_shared_batch_detail_rows.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/68_query_shared_batch_detail_rows.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/68_query_shared_batch_detail_rows.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/68-query-shared-batch-detail-rows-flow.json`

## Step 3: 这节意味着什么

到这里，你已经能分别看清：

- 摘要层
- 明细层

接下来就可以开始做真正的对比，而不是只看有没有“写成功”。

## 下一节

下一节把源 CSV 和 DB 明细放到一起比：
[Lesson 69](69-compare-source-csv-and-db-rows.md)
