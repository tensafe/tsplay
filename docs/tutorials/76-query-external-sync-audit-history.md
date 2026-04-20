# Lesson 76: 读回某个批次的审计历史

上一节只写了一条审计记录。  
这一节继续往前，改成“按批次查历史”。

目标：

- `db_query`

## 开始前

建议先跑完：

- [Lesson 75](75-write-external-sync-audit-row.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/76_query_external_sync_audit_history.lua](../../script/tutorials/76_query_external_sync_audit_history.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/76_query_external_sync_audit_history.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/76_query_external_sync_audit_history.lua
```

预期结果：

- 会生成 `artifacts/tutorials/76-query-external-sync-audit-history-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/76_query_external_sync_audit_history.flow.yaml](../../script/tutorials/76_query_external_sync_audit_history.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/76_query_external_sync_audit_history.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/76_query_external_sync_audit_history.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/76-query-external-sync-audit-history-flow.json`

## 下一节

下一节把这份审计历史导出成 CSV：
[Lesson 77](77-export-external-sync-audit-history.md)
