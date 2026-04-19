# Lesson 06: Redis 基础读写和计数

这一节开始接触外部系统。  
目标不是一下子讲复杂场景，而是先把最常见的几件事跑通：

- `redis_set`
- `redis_get`
- `redis_incr`
- `redis_del`

## 先准备 Redis

如果你本机已经有 Redis，可以直接跳到下一步。  
如果只是为了跟教程走一遍，最简单的方式之一是：

```bash
docker run --name tsplay-redis -p 6379:6379 -d redis:7
```

然后在仓库根目录加载示例环境变量：

```bash
source script/tutorials/env/06_redis_example.sh
```

这个示例文件当前只做一件事：

```bash
export TSPLAY_REDIS_ADDR=127.0.0.1:6379
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/06_redis_round_trip.lua](../../script/tutorials/06_redis_round_trip.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/06_redis_round_trip.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/06_redis_round_trip.lua
```

预期结果：

- 会写入几个 `tutorial:*` 前缀的 key
- 会生成 `artifacts/tutorials/06-redis-round-trip-lua.json`

这一版除了字符串，还顺手写了一个 JSON payload，然后用 `json_extract` 继续把里面的 `status` 拿出来。

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/06_redis_round_trip.flow.yaml](../../script/tutorials/06_redis_round_trip.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/06_redis_round_trip.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/06_redis_round_trip.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/06-redis-round-trip-flow.json`
- 结果里能看到字符串值、JSON 字段和计数器结果

## Step 3: 这节适合迁移到哪些真实场景

- 暂存 cookie / token
- 去重标记
- 任务计数器
- 断点续跑 checkpoint

补充说明：

- 本地直接跑 `go run . -flow ...` 不需要显式加 `allow_redis`
- 如果以后改成 MCP 模式运行，再补 `allow_redis=true`

## 下一节

下一节继续往外部系统走，改成 Postgres：
[Lesson 07](07-db-postgres-basics.md)
