# Lesson 122: 用 `allow_http` 放行一条最小 HTTP Flow

`Lesson 121` 先建立了最基本的 blocked / allowed 手势。  
这一节继续用同样节奏看第二类边界：

- `allow_http`

目标：

- `tsplay.validate_flow`
- `allow_http`
- 一条最小 `http_request` Flow

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/122_security_allow_http.flow.yaml](../../script/tutorials/122_security_allow_http.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/122_mcp_validate_allow_http_blocked.args.json](../../script/tutorials/122_mcp_validate_allow_http_blocked.args.json)
- allowed:
  [../../script/tutorials/122_mcp_validate_allow_http_allowed.args.json](../../script/tutorials/122_mcp_validate_allow_http_allowed.args.json)

## Step 1: 先看默认边界下为什么被拦

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/122_mcp_validate_allow_http_blocked.args.json > artifacts/tutorials/122-mcp-validate-allow-http-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/122_mcp_validate_allow_http_blocked.args.json > artifacts/tutorials/122-mcp-validate-allow-http-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/122-mcp-validate-allow-http-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_http`

## Step 2: 只打开 `allow_http`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/122_mcp_validate_allow_http_allowed.args.json > artifacts/tutorials/122-mcp-validate-allow-http-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/122_mcp_validate_allow_http_allowed.args.json > artifacts/tutorials/122-mcp-validate-allow-http-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/122-mcp-validate-allow-http-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

到这里你应该开始感受到：

- “页面动作能跑”不等于“外部请求就默认能跑”
- MCP 的边界是按能力类别拆开的

## 下一步

继续看：
[Lesson 123](123-allow-file-access-boundary.md)
