# Lesson 129: 理解 `security_preset` 和显式 `allow_*` 覆盖

`Lesson 121-128` 先把单个 `allow_*` 的边界都跑清楚了。  
这一节开始理解更适合日常使用的组合方式：

- `security_preset`
- 显式 `allow_*` 覆盖

目标：

- `browser_write`
- `full_automation`
- 显式 `allow_http=false` 覆盖 preset

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

参数文件：

- browser_write:
  [../../script/tutorials/129_mcp_validate_file_access_browser_write.args.json](../../script/tutorials/129_mcp_validate_file_access_browser_write.args.json)
- full_automation:
  [../../script/tutorials/129_mcp_validate_http_full_automation.args.json](../../script/tutorials/129_mcp_validate_http_full_automation.args.json)
- full_automation + override:
  [../../script/tutorials/129_mcp_validate_http_full_automation_override.args.json](../../script/tutorials/129_mcp_validate_http_full_automation_override.args.json)

## Step 1: 先看 `browser_write` 如何放行文件输出

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_file_access_browser_write.args.json > artifacts/tutorials/129-mcp-validate-file-access-browser-write.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_file_access_browser_write.args.json > artifacts/tutorials/129-mcp-validate-file-access-browser-write.json
```

预期结果：

- 会生成 `artifacts/tutorials/129-mcp-validate-file-access-browser-write.json`
- 里面会看到 `valid=true`
- `security.preset` 会是 `browser_write`

## Step 2: 再看 `full_automation` 如何放行 HTTP

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_http_full_automation.args.json > artifacts/tutorials/129-mcp-validate-http-full-automation.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_http_full_automation.args.json > artifacts/tutorials/129-mcp-validate-http-full-automation.json
```

预期结果：

- 会生成 `artifacts/tutorials/129-mcp-validate-http-full-automation.json`
- 里面会看到 `valid=true`

## Step 3: 最后看显式 `allow_*` 如何覆盖 preset

运行命令：

```bash
# 方式 A：直接运行源码
go run . -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_http_full_automation_override.args.json > artifacts/tutorials/129-mcp-validate-http-full-automation-override.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action mcp-tool -tool tsplay.validate_flow -args-file script/tutorials/129_mcp_validate_http_full_automation_override.args.json > artifacts/tutorials/129-mcp-validate-http-full-automation-override.json
```

预期结果：

- 会生成 `artifacts/tutorials/129-mcp-validate-http-full-automation-override.json`
- 里面会看到 `valid=false`
- 错误信息里仍然会出现 `allow_http`

## Step 4: 这一节意味着什么

到这里，安全边界这一段就从“单个开关”走到了“组合策略”：

- preset 是快捷方式
- 显式 `allow_*` 才是最后的精确覆盖

## 下一步

继续看：
[Lesson 130](130-security-boundary-learning-checkpoint.md)
