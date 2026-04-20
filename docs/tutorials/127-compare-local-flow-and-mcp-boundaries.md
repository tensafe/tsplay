# Lesson 127: 对比本地 Flow 和 MCP 的权限边界

`Lesson 121-126` 已经把每个常见 `allow_*` 都拆开看过一遍。  
这一节不再增加新动作，而是把“本地 Flow”和“MCP”放在一起看。

目标：

- 理解本地 `-flow` 和 `mcp-tool validate_flow` 的差异
- 理解为什么同一条 Flow，在不同入口会有不同边界

## 准备工作

这一节继续复用：

- [../../script/tutorials/123_security_allow_file_access.flow.yaml](../../script/tutorials/123_security_allow_file_access.flow.yaml)
- `Lesson 123` 产出的 blocked / allowed 校验结果

## Step 1: 先在本地直接运行同一条 Flow

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/123_security_allow_file_access.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/123_security_allow_file_access.flow.yaml
```

预期结果：

- 会直接生成 `artifacts/tutorials/123-security-allow-file-access-flow.json`

## Step 2: 再回看 MCP 里的 blocked / allowed 对照

重点回看：

- `artifacts/tutorials/123-mcp-validate-allow-file-access-blocked.json`
- `artifacts/tutorials/123-mcp-validate-allow-file-access-allowed.json`

## Step 3: 这一节意味着什么

这里最重要的不是“哪种入口更好”。  
而是开始建立一条清晰边界：

- 本地 `-flow` 更像你在自己机器上做开发和固化
- `mcp-tool` / MCP Server 更像把能力暴露给 AI 或外部协作者
- 一旦进入 MCP，这些高风险能力就应该被显式声明

## 下一步

继续看：
[Lesson 128](128-why-security-boundaries-come-first.md)
