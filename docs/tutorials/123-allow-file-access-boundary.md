# Lesson 123: 用 `allow_file_access` 放行一条最小文件输出 Flow

`Lesson 122` 处理的是外部请求边界。  
这一节继续往前走，但换成更贴近日常交付的一类：

- `allow_file_access`

目标：

- `tsplay.validate_flow`
- `allow_file_access`
- `write_json`

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

示例 Flow：
[../../script/tutorials/123_security_allow_file_access.flow.yaml](../../script/tutorials/123_security_allow_file_access.flow.yaml)

参数文件：

- blocked:
  [../../script/tutorials/123_mcp_validate_allow_file_access_blocked.args.json](../../script/tutorials/123_mcp_validate_allow_file_access_blocked.args.json)
- allowed:
  [../../script/tutorials/123_mcp_validate_allow_file_access_allowed.args.json](../../script/tutorials/123_mcp_validate_allow_file_access_allowed.args.json)

## Step 1: 先看 MCP 默认边界下为什么被拦

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/123_mcp_validate_allow_file_access_blocked.args.json > artifacts/tutorials/123-mcp-validate-allow-file-access-blocked.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/123_mcp_validate_allow_file_access_blocked.args.json > artifacts/tutorials/123-mcp-validate-allow-file-access-blocked.json
```

预期结果：

- 会生成 `artifacts/tutorials/123-mcp-validate-allow-file-access-blocked.json`
- 里面会看到 `valid=false`
- 错误信息里会出现 `allow_file_access`

## Step 2: 只打开 `allow_file_access`

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/123_mcp_validate_allow_file_access_allowed.args.json > artifacts/tutorials/123-mcp-validate-allow-file-access-allowed.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/123_mcp_validate_allow_file_access_allowed.args.json > artifacts/tutorials/123-mcp-validate-allow-file-access-allowed.json
```

预期结果：

- 会生成 `artifacts/tutorials/123-mcp-validate-allow-file-access-allowed.json`
- 里面会看到 `valid=true`

## Step 3: 这一节意味着什么

从这一节开始，安全边界已经不再只是“抽象概念”。  
因为文件输出、截图、上传、下载，都是业务里非常常见的一类高频动作。

## 下一步

继续看：
[Lesson 124](124-allow-browser-state-boundary.md)
