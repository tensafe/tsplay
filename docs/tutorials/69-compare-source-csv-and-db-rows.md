# Lesson 69: 把源 CSV 和 DB 明细行放到一起比

这一节开始进入真正的对账动作。

我们不再只问“有没有写进去”，而是开始问：

- 行数是不是一样
- 每一行是不是同一个人
- 电话、状态是不是也一致

目标：

- `read_csv`
- `db_query`
- 行级对比

## 开始前

建议先跑完：

- [Lesson 67](67-transaction-store-shared-batch-rows.md)
- [Lesson 68](68-query-shared-batch-detail-rows.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/69_compare_source_csv_and_db_rows.lua](../../script/tutorials/69_compare_source_csv_and_db_rows.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/69_compare_source_csv_and_db_rows.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/69_compare_source_csv_and_db_rows.lua
```

预期结果：

- 会生成 `artifacts/tutorials/69-compare-source-csv-and-db-rows-lua.json`
- Lua 版本会逐行检查 CSV 和 DB 明细是不是一致

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml](../../script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/69-compare-source-csv-and-db-rows-flow.json`

## Step 3: 这节意味着什么

现在你已经从“流程跑通”进入了“结果核对”的阶段。  
这一步是业务自动化里非常重要的分水岭。

## 下一节

下一节把 CSV、Redis、Postgres 三边放进一份统一的对账包：
[Lesson 70](70-build-reconciliation-pack-from-csv-redis-db.md)
