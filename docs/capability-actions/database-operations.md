# 能力动作类别：数据库操作

这组动作把页面结果、接口结果和数据库状态连成一条闭环。  
在 Flow / MCP 中，重点授权是 `allow_database`。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `db_insert` | 是 | 是 | 是 | `action: db_insert` / `db_insert({table=..., row=...})` | 插入一行。适合写结果摘要、审计行。 |
| `db_insert_many` | 是 | 是 | 是 | `action: db_insert_many` / `db_insert_many({table=..., rows=...})` | 批量插入多行。 |
| `db_upsert` | 是 | 是 | 是 | `action: db_upsert` / `db_upsert({table=..., row=..., key_columns=...})` | 按主键或自然键写入或更新。 |
| `db_query` | 是 | 是 | 是 | `action: db_query` / `db_query({sql=..., args=...})` | 查询多行，返回 `list<object>`。 |
| `db_query_one` | 是 | 是 | 是 | `action: db_query_one` / `db_query_one({sql=..., args=...})` | 查询单行，返回对象或 `null`。 |
| `db_execute` | 是 | 是 | 是 | `action: db_execute` / `db_execute({sql=..., args=...})` | 执行非查询 SQL，返回执行元信息。 |
| `db_transaction` | 是 | 是 | 是 | `action: db_transaction` + `steps` / `db_transaction(function() ... end, timeout)` | 在事务里跑一组数据库动作。Flow 侧是嵌套 `steps`，Lua 侧是回调函数。成功自动提交，失败自动回滚。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: db_demo
steps:
  - action: db_insert
    connection: reporting
    save_as: insert_result
    with:
      table: public.tutorial_import_batches
      row:
        batch_id: "demo-001"
        row_count: 3

  - action: db_query_one
    connection: reporting
    save_as: batch_row
    with:
      sql: SELECT batch_id, row_count FROM public.tutorial_import_batches WHERE batch_id = $1
      args:
        - "demo-001"
```

### Lua

```lua
local result = db_transaction(function()
  return db_insert({
    table = "public.tutorial_import_batches",
    row = {batch_id = "demo-001", row_count = 3},
    connection = "reporting",
    driver = "pgsql",
  })
end, 5000)
print(result)
```

## 使用建议

- `db_insert_many` 适合批量落明细，`db_upsert` 适合幂等更新
- `db_query_one` 比 `db_query` 更适合“明确只取一行”的意图
- 需要确保一组写入同成同败时，用 `db_transaction`
- SQL Server / Oracle 构建时要注意对应 build tags

## 相关教程

- [Lesson 61](../tutorials/61-db-insert-import-batch-summary.md)
- [Lesson 64](../tutorials/64-db-transaction-import-batch-and-rows.md)
- [Lesson 67](../tutorials/67-transaction-store-shared-batch-rows.md)
- [Lesson 126](../tutorials/126-allow-database-boundary.md)
