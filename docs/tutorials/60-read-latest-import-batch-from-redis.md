# Lesson 60: 把最新 Redis 批次重新读回本地

这一节是 Redis 这一小段的收口。

我们不再写新数据，而是只做读取和验证：

- 先拿到最新批次 id
- 再拼出 payload key
- 最后把 payload 读回本地 JSON

目标：

- `redis_get`
- `json_extract`
- 本地 checkpoint 结果整理

## 开始前

建议先跑完：

- [Lesson 59](59-save-import-batch-key-to-redis.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/60_read_latest_import_batch_from_redis.lua](../../script/tutorials/60_read_latest_import_batch_from_redis.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/60_read_latest_import_batch_from_redis.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/60_read_latest_import_batch_from_redis.lua
```

预期结果：

- 会生成 `artifacts/tutorials/60-read-latest-import-batch-from-redis-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml](../../script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/60-read-latest-import-batch-from-redis-flow.json`

## Step 3: 这节意味着什么

到这里，Redis 这一段已经形成一个完整的最小链路：

- 本地 CSV 能写成摘要
- 摘要能升级成批次化 key
- 最新批次能被重新读回本地

这也是接下来进入数据库前，最自然的一次“先缓存、再持久化”的过渡。

## 下一节

下一节开始把同一份导出结果写入 Postgres：
[Lesson 61](61-db-insert-import-batch-summary.md)
