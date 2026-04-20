# 初级教程

初级教程的重点，是把“单点动作”提升成“一个小流程”。

如果新手教程回答的是“TSPlay 能不能跑起来”，  
那初级教程回答的是“TSPlay 能不能开始接真实业务动作”。

## 适合谁

- 已经能稳定跑通新手本地练习链路
- 想开始接文件、变量、控制流、Redis、数据库
- 想从“会写一段示例”走向“能完成一个简单任务”

## 初级阶段的主线

这一层的主题建议按下面顺序推进：

1. 文件读写
2. 变量组织
3. 控制流
4. HTTP / Redis / DB 基础
5. 会话和 artifact 管理

为什么是这个顺序：

- 先有输入输出
- 再有状态
- 再有流程
- 最后再接外部系统

## 当前已落地的直接入口

- [Lesson 13: 读取本地 CSV 并写出 JSON](13-read-csv-basics.md)
- [Lesson 14: 写出第一份 CSV](14-write-csv-basics.md)
- [Lesson 15: 读取、整理、再写出 CSV](15-read-transform-write-csv.md)
- [Lesson 16: 用 `retry` 处理偶发失败动作](16-retry-flaky-action.md)
- [Lesson 17: 用 `wait_until` 等异步状态完成](17-wait-until-ready.md)
- [Lesson 18: 上传单个本地文件](18-upload-single-file.md)
- [Lesson 19: 上传多个本地文件](19-upload-multiple-files.md)
- [Lesson 20: 下载本地报表并回读验证](20-download-report.md)
- [Lesson 21: 用 `if` 处理可选登录分支](21-if-optional-login.md)
- [Lesson 22: 用 `foreach` 批量导入 CSV](22-foreach-batch-import-csv.md)
- [Lesson 23: 用 `on_error` 做局部恢复并回写结果](23-on-error-import-recovery.md)
- [Lesson 24: 读取第一份 Excel](24-read-excel-basics.md)
- [Lesson 25: 读取 Excel 指定区域并显式声明表头](25-read-excel-range-headers.md)
- [Lesson 26: 用 Excel 驱动批量导入](26-foreach-batch-import-excel.md)
- [Lesson 27: Excel 批量导入、局部恢复与结果回写](27-on-error-import-excel-writeback.md)
- [Lesson 28: 读取当前浏览器的 Storage State](28-inspect-storage-state.md)
- [Lesson 29: 读取当前浏览器的 Cookie 字符串](29-read-cookies-string.md)
- [Lesson 30: 生成一份浏览器状态快照](30-browser-state-snapshot-pack.md)
- [Lesson 31: 截一张完整页面截图](31-full-page-screenshot.md)
- [Lesson 32: 截一张元素级截图](32-element-screenshot.md)
- [Lesson 33: 保存当前页面的 HTML](33-save-html-basics.md)
- [Lesson 34: 生成一份调试产物包](34-debug-artifact-pack.md)
- [Lesson 35: 在失败分支里保存错误现场](35-error-evidence-pack.md)
- [Lesson 36: 把当前浏览器状态保存到文件](36-save-storage-state.md)
- [Lesson 37: 从保存好的状态文件直接复用登录态](37-load-saved-storage-state.md)
- [Lesson 38: 验证加载后的状态到底是不是你想要的](38-verify-loaded-storage-state.md)
- [Lesson 39: 把“保存状态”和“复用状态”连成一次完整 round trip](39-storage-state-round-trip.md)
- [Lesson 40: 把状态文件注册成命名会话](40-save-named-session.md)
- [Lesson 41: 查看和导出命名会话信息](41-inspect-named-session.md)
- [Lesson 42: 用命名会话直接复用登录态](42-use-named-session.md)
- [Lesson 43: 删除一个已经不用的命名会话](43-delete-named-session.md)
- [Lesson 44: 登录受会话保护的导入页并完成一条导入](44-session-import-with-login.md)
- [Lesson 45: 用状态文件直接跳过登录进入受保护导入页](45-storage-state-auth-import.md)
- [Lesson 46: 把状态文件注册成导入专用命名会话](46-save-import-session.md)
- [Lesson 47: 用命名会话直接进入受保护导入页](47-use-session-import-single.md)
- [Lesson 48: 用命名会话驱动 CSV 批量导入](48-use-session-batch-import-csv.md)
- [Lesson 49: 用命名会话做带恢复的 CSV 批量导入](49-use-session-import-recovery-csv.md)
- [Lesson 50: 用命名会话驱动 Excel 批量导入](50-use-session-batch-import-excel.md)
- [Lesson 51: 用命名会话做带恢复的 Excel 批量导入](51-use-session-import-recovery-excel.md)
- [Lesson 52: 抓取认证导入页上的结果表](52-use-session-capture-import-table.md)
- [Lesson 53: 把认证导入页上的结果表写成本地 CSV](53-use-session-capture-import-table-to-csv.md)
- [Lesson 54: 下载认证导入页当前导出的 CSV](54-use-session-download-import-report.md)
- [Lesson 55: 把认证导出 CSV 下载后再读回来](55-use-session-download-import-report-readback.md)
- [Lesson 56: 把认证页面结果表和下载文件放在一起比对](56-use-session-compare-table-and-download.md)
- [Lesson 57: 跑通认证导入到导出的完整 round trip](57-use-session-import-export-round-trip.md)
- [Lesson 58: 把认证导出 CSV 的摘要写入 Redis](58-sync-import-report-summary-to-redis.md)
- [Lesson 59: 给认证导出结果分配 Redis 批次 key](59-save-import-batch-key-to-redis.md)
- [Lesson 60: 把最新 Redis 批次重新读回本地](60-read-latest-import-batch-from-redis.md)
- [Lesson 61: 把认证导出结果写成一条 Postgres 批次摘要](61-db-insert-import-batch-summary.md)
- [Lesson 62: 查询多条 Postgres 批次摘要](62-db-query-import-batch-summaries.md)
- [Lesson 63: 用 `db_upsert` 更新 Postgres 批次摘要](63-db-upsert-import-batch-summary.md)
- [Lesson 64: 在一个事务里写入批次摘要和明细行](64-db-transaction-import-batch-and-rows.md)
- [Lesson 65: 把最新 Redis 批次摘要同步到 Postgres](65-sync-latest-redis-batch-to-postgres-summary.md)
- [Lesson 66: 一次读回 Redis 和 Postgres 的共享批次摘要](66-query-shared-batch-summary-from-redis-and-postgres.md)
- [Lesson 67: 用共享批次号把明细行写入 Postgres](67-transaction-store-shared-batch-rows.md)
- [Lesson 68: 读回共享批次的 Postgres 明细行](68-query-shared-batch-detail-rows.md)
- [Lesson 69: 把源 CSV 和 DB 明细行放到一起比](69-compare-source-csv-and-db-rows.md)
- [Lesson 70: 生成一份 CSV、Redis、Postgres 三边对账包](70-build-reconciliation-pack-from-csv-redis-db.md)
- [Lesson 71: 跑通一次完整的外部系统 round trip](71-external-system-round-trip.md)
- [Lesson 72: 用同一个批次号重跑同步，但不产生重复数据](72-rerun-shared-batch-idempotently.md)
- [Lesson 73: 验证重跑后没有重复行](73-verify-rerun-does-not-duplicate-rows.md)
- [Lesson 74: 遇到坏数据时，保留有效行并写出异常台账](74-recover-external-sync-with-anomaly-ledger.md)
- [Lesson 75: 给外部同步写入一条审计记录](75-write-external-sync-audit-row.md)
- [Lesson 76: 读回某个批次的审计历史](76-query-external-sync-audit-history.md)
- [Lesson 77: 把审计历史导出成 CSV](77-export-external-sync-audit-history.md)
- [Lesson 78: 清理最新批次的运行数据，但保留审计](78-cleanup-latest-external-batch.md)
- [Lesson 79: 验证批次清理后，审计仍然保留](79-verify-external-batch-cleanup.md)
- [Lesson 80: 跑通一条完整的外部同步生命周期](80-external-sync-lifecycle-round-trip.md)
- [Lesson 06: Redis 基础读写和计数](06-redis-round-trip.md)
- [Lesson 07: Postgres 基础查询与写入](07-db-postgres-basics.md)

这一层现在已经有了一条比较完整的最小链路。  
建议顺序是：

1. `Lesson 13-15` 先把文件输入输出跑顺
2. `Lesson 16-17` 再把 `retry` / `wait_until` 吃透
3. `Lesson 18-20` 再把上传 / 下载动作接上
4. `Lesson 21-23` 再把 `if` / `foreach` / `on_error` 串成小流程
5. `Lesson 24-27` 再把 Excel 导入链路打通
6. `Lesson 28-30` 再把浏览器状态观察吃透
7. `Lesson 31-35` 再把截图 / HTML / 错误现场保留接上
8. `Lesson 36-39` 再把状态文件保存和复用串起来
9. `Lesson 40-43` 再把命名会话和 `use_session` 打通
10. `Lesson 44-50` 再把命名会话真正接到受保护业务流程里
11. `Lesson 51-57` 再把认证导入结果做成页面表格、导出文件和回读闭环
12. 在进入 `Lesson 58-64` 之前，先回头补跑 `Lesson 06-07`，把 Redis / Postgres 的最小动作热身一遍
13. `Lesson 58-60` 再把认证导出结果接进 Redis 摘要和批次 key
14. `Lesson 61-64` 最后把同一份导出结果接进 Postgres 摘要、查询、upsert 和事务写入
15. `Lesson 65-71` 再把 Redis 批次、Postgres 摘要/明细和三边对账真正串成一次完整外部同步链
16. `Lesson 72-73` 再把这条链推进到“幂等重跑”和“无重复验证”
17. `Lesson 74` 再补一条异常输入恢复链，学会保留异常台账
18. `Lesson 75-77` 再把外部同步接进审计留痕和审计导出
19. `Lesson 78-80` 最后把清理、验证和完整生命周期闭环补齐

## 初级阶段必须形成的能力

### 1. 能设计变量

不仅要会 `save_as`，还要会给变量起稳定名字。  
变量名稳定，Flow 才好 review、好修复、好复用。

### 2. 能设计流程边界

要开始明确：

- 输入是什么
- 中间变量是什么
- 输出是什么
- 失败时去哪里看 artifact

### 3. 能接一个外部系统

不用一下子全学，但至少要真正接通一种：

- HTTP
- Redis
- Postgres

## 初级阶段的交付物

建议这一层每个主题都产出下面这些内容之一：

- 一条带变量的 Flow
- 一条带控制流的 Flow
- 一条 Redis / DB / HTTP 最小可用示例
- 一份输入输出说明
- 一份失败现场说明

## 初级阶段的退出标准

- 能写出一个 5 到 10 步的小 Flow
- 能独立使用 `save_as`、`set_var`、`assert_*`、`read_csv`、`write_csv`、`read_excel`、`get_storage_state`、`get_cookies_string`、`save_storage_state`、`load_storage_state`
- 能独立使用 `use_session` 把已保存的浏览器状态接进一个真实的批量导入流程
- 能把认证页面里的“页面表格结果”和“导出文件结果”同时保存下来做复盘
- 能把认证导出的 CSV 继续接进 Redis 和 Postgres，理解“浏览器结果 -> 外部系统摘要/持久化”的递进关系
- 能说明“共享 batch id”为什么重要，并能把同一个批次从 Redis 一路接到 Postgres 摘要、Postgres 明细和本地对账包
- 能解释“运行态数据”和“审计留痕”为什么要分开，并能完成一次重跑、审计、清理、验证的最小生命周期
- 能解释为什么某一步放在 `Lua`，某一步放在 `Flow`
- 能至少接通一个外部系统
- 能说明一个流程失败后应该看哪里，并知道什么时候要留截图 / HTML / JSON 证据
- 能解释“直接写状态文件路径”和“使用命名会话”之间的差别，并知道什么时候该把登录态提升成业务专用会话名

## 学完之后去哪里

下一站是：
[track-intermediate.md](track-intermediate.md)
