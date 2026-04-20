# Lesson 124: 用 `allow_browser_state` 放行浏览器状态动作

`Lesson 123` 看的是文件输出。  
这一节继续往前走，但换成浏览器状态这一类：

- `allow_browser_state`

目标：

- `tsplay.validate_flow`
- `allow_browser_state`
- `get_storage_state`

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/124_security_allow_browser_state.flow.yaml](../../script/tutorials/124_security_allow_browser_state.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/124_mcp_validate_allow_browser_state_blocked.args.json](../../script/tutorials/124_mcp_validate_allow_browser_state_blocked.args.json)
- allowed:
  [../../script/tutorials/124_mcp_validate_allow_browser_state_allowed.args.json](../../script/tutorials/124_mcp_validate_allow_browser_state_allowed.args.json)

## Step 1: 先看默认边界下为什么被拦

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/124_mcp_validate_allow_browser_state_blocked.args.json > artifacts/tutorials/124-mcp-validate-allow-browser-state-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/124_mcp_validate_allow_browser_state_blocked.args.json > artifacts/tutorials/124-mcp-validate-allow-browser-state-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/124-mcp-validate-allow-browser-state-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_browser_state`

## Step 2: 只打开 `allow_browser_state`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/124_mcp_validate_allow_browser_state_allowed.args.json > artifacts/tutorials/124-mcp-validate-allow-browser-state-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/124_mcp_validate_allow_browser_state_allowed.args.json > artifacts/tutorials/124-mcp-validate-allow-browser-state-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/124-mcp-validate-allow-browser-state-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

浏览器状态不是普通文本。  
它经常牵涉登录态、Cookie、Storage State 和命名会话，所以需要单独一层边界。

## 下一步

继续看：
[Lesson 125](125-allow-redis-boundary.md)
