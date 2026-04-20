# Lesson 125: 用 `allow_redis` 放行 Redis 动作

`Lesson 124` 看的是浏览器状态。  
这一节继续按同样节奏，把边界切到 Redis：

- `allow_redis`

目标：

- `tsplay.validate_flow`
- `allow_redis`
- `redis_get`

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/125_security_allow_redis.flow.yaml](../../script/tutorials/125_security_allow_redis.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/125_mcp_validate_allow_redis_blocked.args.json](../../script/tutorials/125_mcp_validate_allow_redis_blocked.args.json)
- allowed:
  [../../script/tutorials/125_mcp_validate_allow_redis_allowed.args.json](../../script/tutorials/125_mcp_validate_allow_redis_allowed.args.json)

## Step 1: 先看默认边界下为什么被拦

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/125_mcp_validate_allow_redis_blocked.args.json > artifacts/tutorials/125-mcp-validate-allow-redis-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/125_mcp_validate_allow_redis_blocked.args.json > artifacts/tutorials/125-mcp-validate-allow-redis-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/125-mcp-validate-allow-redis-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_redis`

## Step 2: 只打开 `allow_redis`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/125_mcp_validate_allow_redis_allowed.args.json > artifacts/tutorials/125-mcp-validate-allow-redis-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/125_mcp_validate_allow_redis_allowed.args.json > artifacts/tutorials/125-mcp-validate-allow-redis-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/125-mcp-validate-allow-redis-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

到这里，高级阶段的边界已经开始进入“外部系统状态”层。  
这也是为什么 Redis 不应该和本地页面动作混成一个默认开关。

## 下一步

继续看：
[Lesson 126](126-allow-database-boundary.md)
