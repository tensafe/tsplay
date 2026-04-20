# Lesson 121: 用 `allow_lua` 放行一条最小 Lua Flow

`Lesson 120` 已经把 MCP 主线收口了。  
从这一节开始，我们进入高级阶段的第一块内容：安全边界。

第一步先看最直接的一条边界：

- `allow_lua`

目标：

- `tsplay.validate_flow`
- `allow_lua`
- 同一条 Flow 的 blocked / allowed 对照

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/121_security_allow_lua.flow.yaml](../../script/tutorials/121_security_allow_lua.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/121_mcp_validate_allow_lua_blocked.args.json](../../script/tutorials/121_mcp_validate_allow_lua_blocked.args.json)
- allowed:
  [../../script/tutorials/121_mcp_validate_allow_lua_allowed.args.json](../../script/tutorials/121_mcp_validate_allow_lua_allowed.args.json)

## Step 1: 先在默认边界下校验

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/121_mcp_validate_allow_lua_blocked.args.json > artifacts/tutorials/121-mcp-validate-allow-lua-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/121_mcp_validate_allow_lua_blocked.args.json > artifacts/tutorials/121-mcp-validate-allow-lua-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/121-mcp-validate-allow-lua-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_lua`

## Step 2: 只打开匹配的 `allow_lua`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/121_mcp_validate_allow_lua_allowed.args.json > artifacts/tutorials/121-mcp-validate-allow-lua-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/121_mcp_validate_allow_lua_allowed.args.json > artifacts/tutorials/121-mcp-validate-allow-lua-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/121-mcp-validate-allow-lua-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

这里最重要的不是 Lua 本身。  
而是先建立一个高级阶段的固定动作：

- 先看默认边界为什么拦
- 再只打开和当前动作匹配的一项授权

## 下一步

继续看：
[Lesson 122](122-allow-http-boundary.md)
