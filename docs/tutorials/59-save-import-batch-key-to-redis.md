# Lesson 59: 给认证导出结果分配 Redis 批次 key

上一节我们只写了一份固定摘要。  
这一节继续往前走，把导出结果升级成“带批次号”的 Redis 记录：

- 用计数器分配批次号
- 用批次号拼出 namespaced key
- 再保存一份“最新批次指针”

目标：

- `read_csv`
- `redis_incr`
- `redis_set`
- `redis_get`

## 开始前

建议先跑完：

- [Lesson 58](58-sync-import-report-summary-to-redis.md)

如果你前一节导出文件不是默认路径，也可以先覆盖：

```bash
export TSPLAY_IMPORTED_REPORT=artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/59_save_import_batch_key_to_redis.lua](../../script/tutorials/59_save_import_batch_key_to_redis.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/59_save_import_batch_key_to_redis.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/59_save_import_batch_key_to_redis.lua
```

预期结果：

- Redis 里会新增一个 `tutorial:session_import:batch:*` key
- Redis 里还会写入 `tutorial:session_import:latest_batch`
- 会生成 `artifacts/tutorials/59-save-import-batch-key-to-redis-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/59_save_import_batch_key_to_redis.flow.yaml](../../script/tutorials/59_save_import_batch_key_to_redis.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/59_save_import_batch_key_to_redis.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/59_save_import_batch_key_to_redis.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/59-save-import-batch-key-to-redis-flow.json`
- 输出里能看到 `batch_id`、`payload_key` 和最新批次指针

## Step 3: 这节意味着什么

现在 Redis 里的内容不再只是“一份最新摘要”，而是开始有：

- 一个计数器
- 一个当前最新批次指针
- 多个历史批次 payload

这就是后面做 checkpoint、断点续跑、外部回写时最常见的最小形态。

## 下一节

下一节只做一件事：把“最新批次”完整读回来：
[Lesson 60](60-read-latest-import-batch-from-redis.md)
