# Lesson 74: 遇到坏数据时，保留有效行并写出异常台账

前面几节处理的都是“正常批次”。  
这一节故意换一份带问题的 CSV，练习结果同步前很常见的一件事：

- 好数据继续往下走
- 坏数据不要吞掉
- 单独写成一份异常台账

目标：

- `read_csv`
- `on_error`
- `db_insert`
- `write_csv`

## 开始前

建议先跑完：

- [Lesson 73](73-verify-rerun-does-not-duplicate-rows.md)

默认输入文件是：
[../../demo/data/import_report_with_issue.csv](../../demo/data/import_report_with_issue.csv)

这份文件里故意放了一条缺少手机号的记录。

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua](../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua
```

预期结果：

- 会生成 `artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-lua.csv`
- 会生成 `artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml](../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-flow.csv`
- 会生成 `artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-flow.json`

补充说明：

- 这一节的 Flow 版本第一次用了一个很小的 `lua` 校验块
- 它只负责把坏数据显式变成失败
- 主流程仍然留在 `Flow + on_error + db_insert + write_csv` 这条主线上

## Step 3: 这节意味着什么

到这里，你已经不是只会处理“理想输入”，而是开始具备最小的异常恢复能力。

## 下一节

下一节继续往“可运维”方向走，把正常批次写成审计记录：
[Lesson 75](75-write-external-sync-audit-row.md)
