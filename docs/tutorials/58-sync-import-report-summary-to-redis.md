# Lesson 58: 把认证导出 CSV 的摘要写入 Redis

这一节不再重新讲“怎么导入页面”，而是接着 [Lesson 57](57-use-session-import-export-round-trip.md) 往下走：

- 先复用 `Lesson 57` 产出的导出 CSV
- 再提炼出一份最小摘要
- 最后把这份摘要缓存到 Redis

目标：

- `read_csv`
- `redis_set`
- `redis_get`
- `json_extract`

## 开始前

建议先跑完：

- [Lesson 57](57-use-session-import-export-round-trip.md)
- [Lesson 06](06-redis-round-trip.md)

默认输入文件是 `artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv`。  
如果你前一节跑的是 Lua 版本，可以先切换：

```bash
export TSPLAY_IMPORTED_REPORT=artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/58_sync_import_report_summary_to_redis.lua](../../script/tutorials/58_sync_import_report_summary_to_redis.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/58_sync_import_report_summary_to_redis.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/58_sync_import_report_summary_to_redis.lua
```

预期结果：

- Redis 里会写入 `tutorial:session_import:latest_summary`
- 会生成 `artifacts/tutorials/58-sync-import-report-summary-to-redis-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml](../../script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/58-sync-import-report-summary-to-redis-flow.json`
- 输出里会同时看到本地摘要和 Redis 读回的字段

## Step 3: 这节意味着什么

到这里，`Lesson 57` 的本地导出结果第一次真正进入了外部系统。  
这一步还很轻，只是先放一份摘要，不急着做批次号和多批次管理。

## 下一节

下一节把这份摘要提升成“批次化 key”：
[Lesson 59](59-save-import-batch-key-to-redis.md)
