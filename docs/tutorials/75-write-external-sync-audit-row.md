# Lesson 75: 给外部同步写入一条审计记录

到这里我们已经能同步、重跑、处理异常。  
下一步自然就是留痕。

这一节先做最小审计：

- 为最新批次写一条 audit row
- 把业务数据和审计数据分开

目标：

- `db_upsert`
- `db_query_one`

## 开始前

建议先跑完：

- [Lesson 71](71-external-system-round-trip.md)

先执行这节新增的审计表 SQL：

```bash
psql "postgres://collector:secret@127.0.0.1:5432/analytics?sslmode=disable" \
  -f script/tutorials/sql/75_reporting_import_audit.sql
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/75_write_external_sync_audit_row.lua](../../script/tutorials/75_write_external_sync_audit_row.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/75_write_external_sync_audit_row.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/75_write_external_sync_audit_row.lua
```

预期结果：

- 会生成 `artifacts/tutorials/75-write-external-sync-audit-row-lua.json`
- `public.tutorial_import_audits` 里会有一条当前批次的审计记录

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/75_write_external_sync_audit_row.flow.yaml](../../script/tutorials/75_write_external_sync_audit_row.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/75_write_external_sync_audit_row.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/75_write_external_sync_audit_row.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/75-write-external-sync-audit-row-flow.json`

## Step 3: 这节意味着什么

到这里，批次结果已经不只是“在系统里有数据”，而是开始具备“事后能追”的能力。

## 下一节

下一节把同一个批次的审计历史完整读回来：
[Lesson 76](76-query-external-sync-audit-history.md)
