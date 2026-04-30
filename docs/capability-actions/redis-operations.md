# 能力动作类别：Redis 操作

这组动作适合做轻量状态同步、进度记录、断点续跑和跨步骤共享。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `redis_get` | 是 | 是 | 是 | `action: redis_get` + `key` / `redis_get(key)` | 读取一个键。支持命名连接。 |
| `redis_set` | 是 | 是 | 是 | `action: redis_set` + `key,value` / `redis_set(key, value, ttl)` | 写入一个键，可选 TTL。 |
| `redis_del` | 是 | 是 | 是 | `action: redis_del` + `key` / `redis_del(key)` | 删除一个键。 |
| `redis_incr` | 是 | 是 | 是 | `action: redis_incr` + `key` / `redis_incr(key, delta)` | 递增计数器。适合进度或批次号。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: redis_demo
steps:
  - action: redis_set
    key: "tutorial:latest_batch"
    ttl_seconds: 3600
    with:
      value:
        batch_id: "demo-001"
        status: "done"

  - action: redis_get
    key: "tutorial:latest_batch"
    save_as: latest_batch
```

### Lua

```lua
redis_set("tutorial:latest_batch", {batch_id = "demo-001", status = "done"}, 3600)
local latest_batch = redis_get("tutorial:latest_batch")
print(latest_batch)
```

## 使用建议

- `foreach` 做断点续跑时，Redis 很适合存 checkpoint
- 用命名连接时，团队里最好统一连接名语义，避免教程和生产环境脱节
- Flow / MCP 中想用这组动作，记得确认 `allow_redis`

## 相关教程

- [Lesson 58](../tutorials/58-sync-import-report-summary-to-redis.md)
- [Lesson 59](../tutorials/59-save-import-batch-key-to-redis.md)
- [Lesson 125](../tutorials/125-allow-redis-boundary.md)
