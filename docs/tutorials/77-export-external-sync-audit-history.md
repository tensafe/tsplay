# Lesson 77: 把审计历史导出成 CSV

这一节不再只是查询，而是把查询结果导出成更容易复盘和传递的格式。

目标：

- `db_query`
- `write_csv`

## 开始前

建议先跑完：

- [Lesson 76](76-query-external-sync-audit-history.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/77_export_external_sync_audit_history.lua](../../script/tutorials/77_export_external_sync_audit_history.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/77_export_external_sync_audit_history.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/77_export_external_sync_audit_history.lua
```

预期结果：

- 会生成 `artifacts/tutorials/77-export-external-sync-audit-history-lua.csv`
- 会生成 `artifacts/tutorials/77-export-external-sync-audit-history-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/77_export_external_sync_audit_history.flow.yaml](../../script/tutorials/77_export_external_sync_audit_history.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/77_export_external_sync_audit_history.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/77_export_external_sync_audit_history.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/77-export-external-sync-audit-history-flow.csv`
- 会生成 `artifacts/tutorials/77-export-external-sync-audit-history-flow.json`

## 下一节

下一节进入清理环节，把最新批次的运行数据删掉，但保留审计记录：
[Lesson 78](78-cleanup-latest-external-batch.md)
