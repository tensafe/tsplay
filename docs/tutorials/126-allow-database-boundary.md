# Lesson 126: 用 `allow_database` 放行数据库动作

`Lesson 125` 看的是 Redis。  
这一节继续往前，进入最后一个常见高权限类别：

- `allow_database`

目标：

- `tsplay.validate_flow`
- `allow_database`
- `db_insert`

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/126_security_allow_database.flow.yaml](../../script/tutorials/126_security_allow_database.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/126_mcp_validate_allow_database_blocked.args.json](../../script/tutorials/126_mcp_validate_allow_database_blocked.args.json)
- allowed:
  [../../script/tutorials/126_mcp_validate_allow_database_allowed.args.json](../../script/tutorials/126_mcp_validate_allow_database_allowed.args.json)

## Step 1: 先看默认边界下为什么被拦

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/126_mcp_validate_allow_database_blocked.args.json > artifacts/tutorials/126-mcp-validate-allow-database-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/126_mcp_validate_allow_database_blocked.args.json > artifacts/tutorials/126-mcp-validate-allow-database-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/126-mcp-validate-allow-database-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_database`

## Step 2: 只打开 `allow_database`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/126_mcp_validate_allow_database_allowed.args.json > artifacts/tutorials/126-mcp-validate-allow-database-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/126_mcp_validate_allow_database_allowed.args.json > artifacts/tutorials/126-mcp-validate-allow-database-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/126-mcp-validate-allow-database-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

数据库动作和 Redis 一样，都会改变或读取外部系统状态。  
所以高级阶段要开始很自然地把它看成“默认关闭、按需放开”的能力。

## 下一步

继续看：
[Lesson 127](127-compare-local-flow-and-mcp-boundaries.md)
