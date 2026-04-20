# TSPlay Step-by-Step 教程

这套教程专门解决一个上手问题：

- 不从大而全的概念开始
- 同一个功能，同时给出 `Lua` 和 `Flow` 两种写法
- 先跑通，再逐步理解为什么日常交付更推荐 `Flow`

仓库里当前使用的是 `docs/` 目录，所以这套教程统一放在这里，而不是单独建 `doc/`。

## 两条学习路线

### 路线 A：今天就能开始跑的实战线

适合今天就想把 `tsplay` 跑起来、先建立手感的同学。

1. [Lesson 01: Hello World，不打开网页也能先跑通](01-hello-world.md)
2. [Lesson 02: 打开本地页面并选择下拉框选项](02-local-page-select-option.md)
3. [Lesson 03: 抓取本地表格并写出 JSON](03-capture-table.md)
4. [Lesson 04: 提取文本和 HTML 片段](04-extract-text-and-html.md)
5. [Lesson 05: 请求本地 JSON 并提取字段](05-http-request-json.md)
6. [Lesson 06: Redis 基础读写和计数](06-redis-round-trip.md)
7. [Lesson 07: Postgres 基础查询与写入](07-db-postgres-basics.md)
8. [Lesson 08: 理解内置资源和 `artifacts/` 输出目录](08-bundled-assets-and-artifacts.md)
9. [Lesson 09: 启动本地 demo 服务并拆解页面动作](09-local-demo-anatomy.md)
10. [Lesson 10: 对本地页面做可见性和文本断言](10-assert-page-state.md)
11. [Lesson 11: 改选另一个选项并验证结果](11-select-another-option.md)
12. [Lesson 12: 把页面交互结果整理成自定义 JSON](12-custom-json-output.md)
13. [Lesson 13: 读取本地 CSV 并写出 JSON](13-read-csv-basics.md)
14. [Lesson 14: 写出第一份 CSV](14-write-csv-basics.md)
15. [Lesson 15: 读取、整理、再写出 CSV](15-read-transform-write-csv.md)
16. [Lesson 16: 用 `retry` 处理偶发失败动作](16-retry-flaky-action.md)
17. [Lesson 17: 用 `wait_until` 等异步状态完成](17-wait-until-ready.md)
18. [Lesson 18: 上传单个本地文件](18-upload-single-file.md)
19. [Lesson 19: 上传多个本地文件](19-upload-multiple-files.md)
20. [Lesson 20: 下载本地报表并回读验证](20-download-report.md)
21. [Lesson 21: 用 `if` 处理可选登录分支](21-if-optional-login.md)
22. [Lesson 22: 用 `foreach` 批量导入 CSV](22-foreach-batch-import-csv.md)
23. [Lesson 23: 用 `on_error` 做局部恢复并回写结果](23-on-error-import-recovery.md)
24. [Lesson 24: 读取第一份 Excel](24-read-excel-basics.md)
25. [Lesson 25: 读取 Excel 指定区域并显式声明表头](25-read-excel-range-headers.md)
26. [Lesson 26: 用 Excel 驱动批量导入](26-foreach-batch-import-excel.md)
27. [Lesson 27: Excel 批量导入、局部恢复与结果回写](27-on-error-import-excel-writeback.md)
28. [Lesson 28: 读取当前浏览器的 Storage State](28-inspect-storage-state.md)
29. [Lesson 29: 读取当前浏览器的 Cookie 字符串](29-read-cookies-string.md)
30. [Lesson 30: 生成一份浏览器状态快照](30-browser-state-snapshot-pack.md)
31. [Lesson 31: 截一张完整页面截图](31-full-page-screenshot.md)
32. [Lesson 32: 截一张元素级截图](32-element-screenshot.md)
33. [Lesson 33: 保存当前页面的 HTML](33-save-html-basics.md)
34. [Lesson 34: 生成一份调试产物包](34-debug-artifact-pack.md)
35. [Lesson 35: 在失败分支里保存错误现场](35-error-evidence-pack.md)
36. [Lesson 36: 把当前浏览器状态保存到文件](36-save-storage-state.md)
37. [Lesson 37: 从保存好的状态文件直接复用登录态](37-load-saved-storage-state.md)
38. [Lesson 38: 验证加载后的状态到底是不是你想要的](38-verify-loaded-storage-state.md)
39. [Lesson 39: 把“保存状态”和“复用状态”连成一次完整 round trip](39-storage-state-round-trip.md)
40. [Lesson 40: 把状态文件注册成命名会话](40-save-named-session.md)
41. [Lesson 41: 查看和导出命名会话信息](41-inspect-named-session.md)
42. [Lesson 42: 用命名会话直接复用登录态](42-use-named-session.md)
43. [Lesson 43: 删除一个已经不用的命名会话](43-delete-named-session.md)
44. [Lesson 44: 登录受会话保护的导入页并完成一条导入](44-session-import-with-login.md)
45. [Lesson 45: 用状态文件直接跳过登录进入受保护导入页](45-storage-state-auth-import.md)
46. [Lesson 46: 把状态文件注册成导入专用命名会话](46-save-import-session.md)
47. [Lesson 47: 用命名会话直接进入受保护导入页](47-use-session-import-single.md)
48. [Lesson 48: 用命名会话驱动 CSV 批量导入](48-use-session-batch-import-csv.md)
49. [Lesson 49: 用命名会话做带恢复的 CSV 批量导入](49-use-session-import-recovery-csv.md)
50. [Lesson 50: 用命名会话驱动 Excel 批量导入](50-use-session-batch-import-excel.md)
51. [Lesson 51: 用命名会话做带恢复的 Excel 批量导入](51-use-session-import-recovery-excel.md)
52. [Lesson 52: 抓取认证导入页上的结果表](52-use-session-capture-import-table.md)
53. [Lesson 53: 把认证导入页上的结果表写成本地 CSV](53-use-session-capture-import-table-to-csv.md)
54. [Lesson 54: 下载认证导入页当前导出的 CSV](54-use-session-download-import-report.md)
55. [Lesson 55: 把认证导出 CSV 下载后再读回来](55-use-session-download-import-report-readback.md)
56. [Lesson 56: 把认证页面结果表和下载文件放在一起比对](56-use-session-compare-table-and-download.md)
57. [Lesson 57: 跑通认证导入到导出的完整 round trip](57-use-session-import-export-round-trip.md)
58. [Lesson 58: 把认证导出 CSV 的摘要写入 Redis](58-sync-import-report-summary-to-redis.md)
59. [Lesson 59: 给认证导出结果分配 Redis 批次 key](59-save-import-batch-key-to-redis.md)
60. [Lesson 60: 把最新 Redis 批次重新读回本地](60-read-latest-import-batch-from-redis.md)
61. [Lesson 61: 把认证导出结果写成一条 Postgres 批次摘要](61-db-insert-import-batch-summary.md)
62. [Lesson 62: 查询多条 Postgres 批次摘要](62-db-query-import-batch-summaries.md)
63. [Lesson 63: 用 `db_upsert` 更新 Postgres 批次摘要](63-db-upsert-import-batch-summary.md)
64. [Lesson 64: 在一个事务里写入批次摘要和明细行](64-db-transaction-import-batch-and-rows.md)
65. [Lesson 65: 把最新 Redis 批次摘要同步到 Postgres](65-sync-latest-redis-batch-to-postgres-summary.md)
66. [Lesson 66: 一次读回 Redis 和 Postgres 的共享批次摘要](66-query-shared-batch-summary-from-redis-and-postgres.md)
67. [Lesson 67: 用共享批次号把明细行写入 Postgres](67-transaction-store-shared-batch-rows.md)
68. [Lesson 68: 读回共享批次的 Postgres 明细行](68-query-shared-batch-detail-rows.md)
69. [Lesson 69: 把源 CSV 和 DB 明细行放到一起比](69-compare-source-csv-and-db-rows.md)
70. [Lesson 70: 生成一份 CSV、Redis、Postgres 三边对账包](70-build-reconciliation-pack-from-csv-redis-db.md)
71. [Lesson 71: 跑通一次完整的外部系统 round trip](71-external-system-round-trip.md)
72. [Lesson 72: 用同一个批次号重跑同步，但不产生重复数据](72-rerun-shared-batch-idempotently.md)
73. [Lesson 73: 验证重跑后没有重复行](73-verify-rerun-does-not-duplicate-rows.md)
74. [Lesson 74: 遇到坏数据时，保留有效行并写出异常台账](74-recover-external-sync-with-anomaly-ledger.md)
75. [Lesson 75: 给外部同步写入一条审计记录](75-write-external-sync-audit-row.md)
76. [Lesson 76: 读回某个批次的审计历史](76-query-external-sync-audit-history.md)
77. [Lesson 77: 把审计历史导出成 CSV](77-export-external-sync-audit-history.md)
78. [Lesson 78: 清理最新批次的运行数据，但保留审计](78-cleanup-latest-external-batch.md)
79. [Lesson 79: 验证批次清理后，审计仍然保留](79-verify-external-batch-cleanup.md)
80. [Lesson 80: 跑通一条完整的外部同步生命周期](80-external-sync-lifecycle-round-trip.md)
81. [Lesson 81: 从生命周期 CSV 里读回批次证据](81-read-lifecycle-evidence.md)
82. [Lesson 82: 按生命周期证据回放一个新批次](82-replay-batch-from-lifecycle-evidence.md)
83. [Lesson 83: 用生命周期证据验证回放批次](83-verify-replay-batch-against-lifecycle-evidence.md)
84. [Lesson 84: 给回放批次补写一条审计记录](84-write-replay-audit-row.md)
85. [Lesson 85: 把原批次和回放批次的审计导出成对照 CSV](85-export-original-and-replay-audits.md)
86. [Lesson 86: 生成一份回放后的对账包](86-build-post-replay-reconciliation-pack.md)
87. [Lesson 87: 生成一份交接 artifact manifest](87-build-handoff-artifact-manifest.md)
88. [Lesson 88: 把交接 manifest 整理成交付摘要](88-build-handoff-summary.md)
89. [Lesson 89: 生成发布前检查清单](89-build-pre-release-checklist.md)
90. [Lesson 90: 跑通一条“生命周期证据 -> 回放 -> 交接包”的完整 round trip](90-handoff-round-trip-from-lifecycle-evidence.md)
91. [Lesson 91: 读交接 manifest，识别每份产物的角色](91-read-handoff-manifest-roles.md)
92. [Lesson 92: 把交接产物整理成模板目录](92-build-template-artifact-catalog.md)
93. [Lesson 93: 把交接链整理成 Input -> Process -> Output 模板](93-build-input-process-output-template.md)
94. [Lesson 94: 把交接链整理成 Collect -> Verify -> Save 模板](94-build-collect-verify-save-template.md)
95. [Lesson 95: 把交接链整理成 Replay -> Audit -> Handoff 模板](95-build-replay-audit-handoff-template.md)
96. [Lesson 96: 把几份模板整理成统一索引](96-build-template-index.md)
97. [Lesson 97: 验证模板索引仍然覆盖完整交接链](97-verify-template-covers-handoff-chain.md)
98. [Lesson 98: 生成一份“场景 -> 模板”的学习矩阵](98-build-template-lesson-matrix.md)
99. [Lesson 99: 给模板包生成发布前检查清单](99-build-template-preflight-checklist.md)
100. [Lesson 100: 跑通一条“交接产物 -> 模板包”的完整 round trip](100-template-round-trip-from-handoff-artifacts.md)
101. [Lesson 101: 先确认模板发布卡片真的在页面上](101-assert-visible-template-release-card.md)
102. [Lesson 102: 继续确认模板发布状态文字对不对](102-assert-text-template-release-status.md)
103. [Lesson 103: 用 `retry` 跑通模板发布 gate](103-retry-template-release-gate.md)
104. [Lesson 104: 用 `wait_until` 等模板发布检查完成](104-wait-until-template-release-ready.md)
105. [Lesson 105: 用 `on_error` 接住模板发布校验失败](105-on-error-template-release-validation.md)
106. [Lesson 106: 等一条延迟出现的发布说明项](106-wait-for-delayed-release-note.md)
107. [Lesson 107: 用 `retry` 接住一次偶发失败点击](107-retry-flaky-publish-click.md)
108. [Lesson 108: `reload` 之后再验证一次恢复结果](108-reload-and-retry-release-recovery.md)
109. [Lesson 109: 给模板发布页留一份调试证据包](109-template-release-artifact-pack.md)
110. [Lesson 110: 跑通一条完整的模板发布稳定性 round trip](110-template-release-robustness-round-trip.md)
111. [Lesson 111: 先用 `tsplay.list_actions` 看清 MCP 到底能做什么](111-mcp-list-actions.md)
112. [Lesson 112: 用 `flow_schema` 和 `flow_examples` 看清 Flow 长什么样](112-mcp-flow-schema-and-examples.md)
113. [Lesson 113: 先用 `observe_page` 观察模板发布页](113-mcp-observe-page.md)
114. [Lesson 114: 用 `draft_flow` 把 observation 变成第一份 Flow 草稿](114-mcp-draft-flow.md)
115. [Lesson 115: 先校验草稿，再决定能不能运行](115-mcp-validate-drafted-flow.md)
116. [Lesson 116: 运行刚刚起草并校验过的 Flow](116-mcp-run-drafted-flow.md)
117. [Lesson 117: 先故意跑坏一次，再生成 repair context](117-mcp-repair-flow-context.md)
118. [Lesson 118: 用 repair context 生成真正可用的修复请求](118-mcp-repair-flow.md)
119. [Lesson 119: 把 `observe -> draft -> validate -> run -> repair` 串成一条线](119-mcp-chain-overview.md)
120. [Lesson 120: 用 `finalize_flow` 收成一份更短的默认入口](120-mcp-finalize-flow.md)
121. [Lesson 121: 用 `allow_lua` 放行一条最小 Lua Flow](121-allow-lua-boundary.md)
122. [Lesson 122: 用 `allow_http` 放行一条最小 HTTP Flow](122-allow-http-boundary.md)
123. [Lesson 123: 用 `allow_file_access` 放行一条最小文件输出 Flow](123-allow-file-access-boundary.md)
124. [Lesson 124: 用 `allow_browser_state` 放行浏览器状态动作](124-allow-browser-state-boundary.md)
125. [Lesson 125: 用 `allow_redis` 放行 Redis 动作](125-allow-redis-boundary.md)
126. [Lesson 126: 用 `allow_database` 放行数据库动作](126-allow-database-boundary.md)
127. [Lesson 127: 对比本地 Flow 和 MCP 的权限边界](127-compare-local-flow-and-mcp-boundaries.md)
128. [Lesson 128: 为什么教程不能跳过权限边界](128-why-security-boundaries-come-first.md)
129. [Lesson 129: 理解 `security_preset` 和显式 `allow_*` 覆盖](129-security-preset-and-override.md)
130. [Lesson 130: 完成安全边界模块的第一轮 checkpoint](130-security-boundary-learning-checkpoint.md)

### 路线 B：完整进阶教程体系

适合要把 TSPlay 做成一套系统学习路径，而不是只跑几个 demo 的同学。

1. [完整课程总览](curriculum-overview.md)
2. [新手教程](track-newbie.md)
3. [初级教程](track-junior.md)
4. [中级教程](track-intermediate.md)
5. [高级教程](track-advanced.md)
6. [160 次递进迭代路线图](iteration-roadmap-160.md)
7. [教程持续进化手册](evolution-playbook.md)

## 教程地图

| Lesson | 功能 | Lua 示例 | Flow 示例 | 说明 |
| --- | --- | --- | --- | --- |
| 01 | Hello World + 写 JSON | [../../script/tutorials/01_hello_world.lua](../../script/tutorials/01_hello_world.lua) | [../../script/tutorials/01_hello_world.flow.yaml](../../script/tutorials/01_hello_world.flow.yaml) | 不需要打开网页 |
| 02 | 打开本地 demo 页并选择选项 | [../../script/tutorials/02_select_option.lua](../../script/tutorials/02_select_option.lua) | [../../script/tutorials/02_select_option.flow.yaml](../../script/tutorials/02_select_option.flow.yaml) | 需要一个本地静态文件服务 |
| 03 | 抓取本地表格并写出 JSON | [../../script/tutorials/03_capture_table.lua](../../script/tutorials/03_capture_table.lua) | [../../script/tutorials/03_capture_table.flow.yaml](../../script/tutorials/03_capture_table.flow.yaml) | 继续复用本地静态文件服务 |
| 04 | 提取文本和 HTML 片段 | [../../script/tutorials/04_extract_text_and_html.lua](../../script/tutorials/04_extract_text_and_html.lua) | [../../script/tutorials/04_extract_text_and_html.flow.yaml](../../script/tutorials/04_extract_text_and_html.flow.yaml) | 继续复用本地静态文件服务 |
| 05 | 请求本地 JSON 并提取字段 | [../../script/tutorials/05_http_request_json.lua](../../script/tutorials/05_http_request_json.lua) | [../../script/tutorials/05_http_request_json.flow.yaml](../../script/tutorials/05_http_request_json.flow.yaml) | 继续复用本地静态文件服务 |
| 06 | Redis 基础读写和计数 | [../../script/tutorials/06_redis_round_trip.lua](../../script/tutorials/06_redis_round_trip.lua) | [../../script/tutorials/06_redis_round_trip.flow.yaml](../../script/tutorials/06_redis_round_trip.flow.yaml) | 需要本地 Redis |
| 07 | Postgres 基础查询与写入 | [../../script/tutorials/07_db_postgres_basics.lua](../../script/tutorials/07_db_postgres_basics.lua) | [../../script/tutorials/07_db_postgres_basics.flow.yaml](../../script/tutorials/07_db_postgres_basics.flow.yaml) | 需要本地 Postgres |
| 08 | 理解内置资源和 `artifacts/` 输出目录 | - | - | 命令和理解型章节，复用 Lesson 01 |
| 09 | 启动本地 demo 服务并拆解页面动作 | - | - | 命令和理解型章节，复用 Lesson 02 |
| 10 | 对本地页面做可见性和文本断言 | [../../script/tutorials/10_assert_page_state.lua](../../script/tutorials/10_assert_page_state.lua) | [../../script/tutorials/10_assert_page_state.flow.yaml](../../script/tutorials/10_assert_page_state.flow.yaml) | 继续复用本地静态文件服务 |
| 11 | 改选另一个选项并验证结果 | [../../script/tutorials/11_select_another_option.lua](../../script/tutorials/11_select_another_option.lua) | [../../script/tutorials/11_select_another_option.flow.yaml](../../script/tutorials/11_select_another_option.flow.yaml) | 继续复用本地静态文件服务 |
| 12 | 把页面交互结果整理成自定义 JSON | [../../script/tutorials/12_custom_json_output.lua](../../script/tutorials/12_custom_json_output.lua) | [../../script/tutorials/12_custom_json_output.flow.yaml](../../script/tutorials/12_custom_json_output.flow.yaml) | 继续复用本地静态文件服务 |
| 13 | 读取本地 CSV 并写出 JSON | [../../script/tutorials/13_read_csv_basics.lua](../../script/tutorials/13_read_csv_basics.lua) | [../../script/tutorials/13_read_csv_basics.flow.yaml](../../script/tutorials/13_read_csv_basics.flow.yaml) | 需要本地 CSV 文件 |
| 14 | 写出第一份 CSV | [../../script/tutorials/14_write_csv_basics.lua](../../script/tutorials/14_write_csv_basics.lua) | [../../script/tutorials/14_write_csv_basics.flow.yaml](../../script/tutorials/14_write_csv_basics.flow.yaml) | 不需要打开网页 |
| 15 | 读取、整理、再写出 CSV | [../../script/tutorials/15_read_transform_write_csv.lua](../../script/tutorials/15_read_transform_write_csv.lua) | [../../script/tutorials/15_read_transform_write_csv.flow.yaml](../../script/tutorials/15_read_transform_write_csv.flow.yaml) | 需要本地 CSV 文件 |
| 16 | 用 `retry` 处理偶发失败动作 | [../../script/tutorials/16_retry_flaky_action.lua](../../script/tutorials/16_retry_flaky_action.lua) | [../../script/tutorials/16_retry_flaky_action.flow.yaml](../../script/tutorials/16_retry_flaky_action.flow.yaml) | 需要本地静态文件服务 |
| 17 | 用 `wait_until` 等异步状态完成 | [../../script/tutorials/17_wait_until_ready.lua](../../script/tutorials/17_wait_until_ready.lua) | [../../script/tutorials/17_wait_until_ready.flow.yaml](../../script/tutorials/17_wait_until_ready.flow.yaml) | 需要本地静态文件服务 |
| 18 | 上传单个本地文件 | [../../script/tutorials/18_upload_single_file.lua](../../script/tutorials/18_upload_single_file.lua) | [../../script/tutorials/18_upload_single_file.flow.yaml](../../script/tutorials/18_upload_single_file.flow.yaml) | 需要本地静态文件服务和本地文件 |
| 19 | 上传多个本地文件 | [../../script/tutorials/19_upload_multiple_files.lua](../../script/tutorials/19_upload_multiple_files.lua) | [../../script/tutorials/19_upload_multiple_files.flow.yaml](../../script/tutorials/19_upload_multiple_files.flow.yaml) | 需要本地静态文件服务和本地文件 |
| 20 | 下载本地报表并回读验证 | [../../script/tutorials/20_download_report.lua](../../script/tutorials/20_download_report.lua) | [../../script/tutorials/20_download_report.flow.yaml](../../script/tutorials/20_download_report.flow.yaml) | 需要本地静态文件服务 |
| 21 | 用 `if` 处理可选登录分支 | [../../script/tutorials/21_if_optional_login.lua](../../script/tutorials/21_if_optional_login.lua) | [../../script/tutorials/21_if_optional_login.flow.yaml](../../script/tutorials/21_if_optional_login.flow.yaml) | 需要本地静态文件服务 |
| 22 | 用 `foreach` 批量导入 CSV | [../../script/tutorials/22_foreach_batch_import_csv.lua](../../script/tutorials/22_foreach_batch_import_csv.lua) | [../../script/tutorials/22_foreach_batch_import_csv.flow.yaml](../../script/tutorials/22_foreach_batch_import_csv.flow.yaml) | 需要本地静态文件服务和本地 CSV |
| 23 | 用 `on_error` 做局部恢复并回写结果 | [../../script/tutorials/23_on_error_import_recovery.lua](../../script/tutorials/23_on_error_import_recovery.lua) | [../../script/tutorials/23_on_error_import_recovery.flow.yaml](../../script/tutorials/23_on_error_import_recovery.flow.yaml) | 需要本地静态文件服务和本地 CSV |
| 24 | 读取第一份 Excel | [../../script/tutorials/24_read_excel_basics.lua](../../script/tutorials/24_read_excel_basics.lua) | [../../script/tutorials/24_read_excel_basics.flow.yaml](../../script/tutorials/24_read_excel_basics.flow.yaml) | 需要本地 Excel 文件 |
| 25 | 读取 Excel 指定区域并显式声明表头 | [../../script/tutorials/25_read_excel_range_headers.lua](../../script/tutorials/25_read_excel_range_headers.lua) | [../../script/tutorials/25_read_excel_range_headers.flow.yaml](../../script/tutorials/25_read_excel_range_headers.flow.yaml) | 需要本地 Excel 文件 |
| 26 | 用 Excel 驱动批量导入 | [../../script/tutorials/26_foreach_batch_import_excel.lua](../../script/tutorials/26_foreach_batch_import_excel.lua) | [../../script/tutorials/26_foreach_batch_import_excel.flow.yaml](../../script/tutorials/26_foreach_batch_import_excel.flow.yaml) | 需要本地静态文件服务和本地 Excel |
| 27 | Excel 批量导入、局部恢复与结果回写 | [../../script/tutorials/27_on_error_import_excel_writeback.lua](../../script/tutorials/27_on_error_import_excel_writeback.lua) | [../../script/tutorials/27_on_error_import_excel_writeback.flow.yaml](../../script/tutorials/27_on_error_import_excel_writeback.flow.yaml) | 需要本地静态文件服务和本地 Excel |
| 28 | 读取当前浏览器的 Storage State | [../../script/tutorials/28_inspect_storage_state.lua](../../script/tutorials/28_inspect_storage_state.lua) | [../../script/tutorials/28_inspect_storage_state.flow.yaml](../../script/tutorials/28_inspect_storage_state.flow.yaml) | 需要本地静态文件服务 |
| 29 | 读取当前浏览器的 Cookie 字符串 | [../../script/tutorials/29_read_cookies_string.lua](../../script/tutorials/29_read_cookies_string.lua) | [../../script/tutorials/29_read_cookies_string.flow.yaml](../../script/tutorials/29_read_cookies_string.flow.yaml) | 需要本地静态文件服务 |
| 30 | 生成一份浏览器状态快照 | [../../script/tutorials/30_browser_state_snapshot_pack.lua](../../script/tutorials/30_browser_state_snapshot_pack.lua) | [../../script/tutorials/30_browser_state_snapshot_pack.flow.yaml](../../script/tutorials/30_browser_state_snapshot_pack.flow.yaml) | 需要本地静态文件服务 |
| 31 | 截一张完整页面截图 | [../../script/tutorials/31_full_page_screenshot.lua](../../script/tutorials/31_full_page_screenshot.lua) | [../../script/tutorials/31_full_page_screenshot.flow.yaml](../../script/tutorials/31_full_page_screenshot.flow.yaml) | 需要本地静态文件服务 |
| 32 | 截一张元素级截图 | [../../script/tutorials/32_element_screenshot.lua](../../script/tutorials/32_element_screenshot.lua) | [../../script/tutorials/32_element_screenshot.flow.yaml](../../script/tutorials/32_element_screenshot.flow.yaml) | 需要本地静态文件服务 |
| 33 | 保存当前页面的 HTML | [../../script/tutorials/33_save_html_basics.lua](../../script/tutorials/33_save_html_basics.lua) | [../../script/tutorials/33_save_html_basics.flow.yaml](../../script/tutorials/33_save_html_basics.flow.yaml) | 需要本地静态文件服务 |
| 34 | 生成一份调试产物包 | [../../script/tutorials/34_debug_artifact_pack.lua](../../script/tutorials/34_debug_artifact_pack.lua) | [../../script/tutorials/34_debug_artifact_pack.flow.yaml](../../script/tutorials/34_debug_artifact_pack.flow.yaml) | 需要本地静态文件服务 |
| 35 | 在失败分支里保存错误现场 | [../../script/tutorials/35_error_evidence_pack.lua](../../script/tutorials/35_error_evidence_pack.lua) | [../../script/tutorials/35_error_evidence_pack.flow.yaml](../../script/tutorials/35_error_evidence_pack.flow.yaml) | 需要本地静态文件服务 |
| 36 | 把当前浏览器状态保存到文件 | [../../script/tutorials/36_save_storage_state.lua](../../script/tutorials/36_save_storage_state.lua) | [../../script/tutorials/36_save_storage_state.flow.yaml](../../script/tutorials/36_save_storage_state.flow.yaml) | 需要本地静态文件服务 |
| 37 | 从保存好的状态文件直接复用登录态 | [../../script/tutorials/37_load_saved_storage_state.lua](../../script/tutorials/37_load_saved_storage_state.lua) | [../../script/tutorials/37_load_saved_storage_state.flow.yaml](../../script/tutorials/37_load_saved_storage_state.flow.yaml) | 默认依赖 Lesson 36 生成的状态文件 |
| 38 | 验证加载后的状态到底是不是你想要的 | [../../script/tutorials/38_verify_loaded_storage_state.lua](../../script/tutorials/38_verify_loaded_storage_state.lua) | [../../script/tutorials/38_verify_loaded_storage_state.flow.yaml](../../script/tutorials/38_verify_loaded_storage_state.flow.yaml) | 默认依赖 Lesson 36 生成的状态文件 |
| 39 | 把“保存状态”和“复用状态”连成一次完整 round trip | - | - | 命令和复盘型章节，复用 Lesson 36-37 |
| 40 | 把状态文件注册成命名会话 | - | - | 命令型章节，使用 `tsplay -action save-session` |
| 41 | 查看和导出命名会话信息 | - | - | 命令型章节，使用 `list/get/export-session` |
| 42 | 用命名会话直接复用登录态 | [../../script/tutorials/42_use_named_session.lua](../../script/tutorials/42_use_named_session.lua) | [../../script/tutorials/42_use_named_session.flow.yaml](../../script/tutorials/42_use_named_session.flow.yaml) | 默认依赖 Lesson 40 注册好的 `session_lab_demo` |
| 43 | 删除一个已经不用的命名会话 | - | - | 命令型章节，使用 `tsplay -action delete-session` |
| 44 | 登录受会话保护的导入页并完成一条导入 | [../../script/tutorials/44_session_import_with_login.lua](../../script/tutorials/44_session_import_with_login.lua) | [../../script/tutorials/44_session_import_with_login.flow.yaml](../../script/tutorials/44_session_import_with_login.flow.yaml) | 新建认证导入主线，从页面登录开始 |
| 45 | 用状态文件直接跳过登录进入受保护导入页 | [../../script/tutorials/45_storage_state_auth_import.lua](../../script/tutorials/45_storage_state_auth_import.lua) | [../../script/tutorials/45_storage_state_auth_import.flow.yaml](../../script/tutorials/45_storage_state_auth_import.flow.yaml) | 默认依赖 Lesson 36 生成的状态文件 |
| 46 | 把状态文件注册成导入专用命名会话 | - | - | 命令型章节，统一注册 `session_import_demo` |
| 47 | 用命名会话直接进入受保护导入页 | [../../script/tutorials/47_use_session_import_single.lua](../../script/tutorials/47_use_session_import_single.lua) | [../../script/tutorials/47_use_session_import_single.flow.yaml](../../script/tutorials/47_use_session_import_single.flow.yaml) | 默认依赖 Lesson 46 注册好的 `session_import_demo` |
| 48 | 用命名会话驱动 CSV 批量导入 | [../../script/tutorials/48_use_session_batch_import_csv.lua](../../script/tutorials/48_use_session_batch_import_csv.lua) | [../../script/tutorials/48_use_session_batch_import_csv.flow.yaml](../../script/tutorials/48_use_session_batch_import_csv.flow.yaml) | 把 `use_session` 接到 CSV 批量导入 |
| 49 | 用命名会话做带恢复的 CSV 批量导入 | [../../script/tutorials/49_use_session_import_recovery_csv.lua](../../script/tutorials/49_use_session_import_recovery_csv.lua) | [../../script/tutorials/49_use_session_import_recovery_csv.flow.yaml](../../script/tutorials/49_use_session_import_recovery_csv.flow.yaml) | 在认证流程里体验 `on_error` |
| 50 | 用命名会话驱动 Excel 批量导入 | [../../script/tutorials/50_use_session_batch_import_excel.lua](../../script/tutorials/50_use_session_batch_import_excel.lua) | [../../script/tutorials/50_use_session_batch_import_excel.flow.yaml](../../script/tutorials/50_use_session_batch_import_excel.flow.yaml) | 把认证导入主线继续扩到 Excel |
| 51 | 用命名会话做带恢复的 Excel 批量导入 | [../../script/tutorials/51_use_session_import_recovery_excel.lua](../../script/tutorials/51_use_session_import_recovery_excel.lua) | [../../script/tutorials/51_use_session_import_recovery_excel.flow.yaml](../../script/tutorials/51_use_session_import_recovery_excel.flow.yaml) | 把 Excel 坏数据恢复接进认证流程 |
| 52 | 抓取认证导入页上的结果表 | [../../script/tutorials/52_use_session_capture_import_table.lua](../../script/tutorials/52_use_session_capture_import_table.lua) | [../../script/tutorials/52_use_session_capture_import_table.flow.yaml](../../script/tutorials/52_use_session_capture_import_table.flow.yaml) | 回到页面事实，直接抓结果表 |
| 53 | 把认证导入页上的结果表写成本地 CSV | [../../script/tutorials/53_use_session_capture_import_table_to_csv.lua](../../script/tutorials/53_use_session_capture_import_table_to_csv.lua) | [../../script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml](../../script/tutorials/53_use_session_capture_import_table_to_csv.flow.yaml) | 从页面表格落到本地 CSV |
| 54 | 下载认证导入页当前导出的 CSV | [../../script/tutorials/54_use_session_download_import_report.lua](../../script/tutorials/54_use_session_download_import_report.lua) | [../../script/tutorials/54_use_session_download_import_report.flow.yaml](../../script/tutorials/54_use_session_download_import_report.flow.yaml) | 正式进入认证页面自带导出 |
| 55 | 把认证导出 CSV 下载后再读回来 | [../../script/tutorials/55_use_session_download_import_report_readback.lua](../../script/tutorials/55_use_session_download_import_report_readback.lua) | [../../script/tutorials/55_use_session_download_import_report_readback.flow.yaml](../../script/tutorials/55_use_session_download_import_report_readback.flow.yaml) | 下载闭环重新回到本地读文件 |
| 56 | 把认证页面结果表和下载文件放在一起比对 | [../../script/tutorials/56_use_session_compare_table_and_download.lua](../../script/tutorials/56_use_session_compare_table_and_download.lua) | [../../script/tutorials/56_use_session_compare_table_and_download.flow.yaml](../../script/tutorials/56_use_session_compare_table_and_download.flow.yaml) | 同时保存页面事实和文件事实 |
| 57 | 跑通认证导入到导出的完整 round trip | [../../script/tutorials/57_use_session_import_export_round_trip.lua](../../script/tutorials/57_use_session_import_export_round_trip.lua) | [../../script/tutorials/57_use_session_import_export_round_trip.flow.yaml](../../script/tutorials/57_use_session_import_export_round_trip.flow.yaml) | 把认证导入、抓表、导出、回读串成完整闭环 |
| 58 | 把认证导出 CSV 的摘要写入 Redis | [../../script/tutorials/58_sync_import_report_summary_to_redis.lua](../../script/tutorials/58_sync_import_report_summary_to_redis.lua) | [../../script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml](../../script/tutorials/58_sync_import_report_summary_to_redis.flow.yaml) | 从 Lesson 57 的导出文件过渡到 Redis 摘要缓存 |
| 59 | 给认证导出结果分配 Redis 批次 key | [../../script/tutorials/59_save_import_batch_key_to_redis.lua](../../script/tutorials/59_save_import_batch_key_to_redis.lua) | [../../script/tutorials/59_save_import_batch_key_to_redis.flow.yaml](../../script/tutorials/59_save_import_batch_key_to_redis.flow.yaml) | 把单份摘要升级成批次化 Redis payload |
| 60 | 把最新 Redis 批次重新读回本地 | [../../script/tutorials/60_read_latest_import_batch_from_redis.lua](../../script/tutorials/60_read_latest_import_batch_from_redis.lua) | [../../script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml](../../script/tutorials/60_read_latest_import_batch_from_redis.flow.yaml) | 用最新批次指针完成一次 Redis 读回闭环 |
| 61 | 把认证导出结果写成一条 Postgres 批次摘要 | [../../script/tutorials/61_db_insert_import_batch_summary.lua](../../script/tutorials/61_db_insert_import_batch_summary.lua) | [../../script/tutorials/61_db_insert_import_batch_summary.flow.yaml](../../script/tutorials/61_db_insert_import_batch_summary.flow.yaml) | 从导出文件进入结构化持久化 |
| 62 | 查询多条 Postgres 批次摘要 | [../../script/tutorials/62_db_query_import_batch_summaries.lua](../../script/tutorials/62_db_query_import_batch_summaries.lua) | [../../script/tutorials/62_db_query_import_batch_summaries.flow.yaml](../../script/tutorials/62_db_query_import_batch_summaries.flow.yaml) | 把单条读回扩展成多条列表查询 |
| 63 | 用 `db_upsert` 更新 Postgres 批次摘要 | [../../script/tutorials/63_db_upsert_import_batch_summary.lua](../../script/tutorials/63_db_upsert_import_batch_summary.lua) | [../../script/tutorials/63_db_upsert_import_batch_summary.flow.yaml](../../script/tutorials/63_db_upsert_import_batch_summary.flow.yaml) | 练习“先占位再补全”的数据库更新模式 |
| 64 | 在一个事务里写入批次摘要和明细行 | [../../script/tutorials/64_db_transaction_import_batch_and_rows.lua](../../script/tutorials/64_db_transaction_import_batch_and_rows.lua) | [../../script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml](../../script/tutorials/64_db_transaction_import_batch_and_rows.flow.yaml) | 用事务把批次摘要和明细行一次写入完成 |
| 65 | 把最新 Redis 批次摘要同步到 Postgres | [../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua](../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.lua) | [../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml](../../script/tutorials/65_sync_latest_redis_batch_to_postgres_summary.flow.yaml) | 从 Redis 最新批次进入共享的 Postgres 摘要 |
| 66 | 一次读回 Redis 和 Postgres 的共享批次摘要 | [../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua](../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.lua) | [../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml](../../script/tutorials/66_query_shared_batch_summary_from_redis_and_postgres.flow.yaml) | 同时核对 Redis、CSV、Postgres 的摘要事实 |
| 67 | 用共享批次号把明细行写入 Postgres | [../../script/tutorials/67_transaction_store_shared_batch_rows.lua](../../script/tutorials/67_transaction_store_shared_batch_rows.lua) | [../../script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml](../../script/tutorials/67_transaction_store_shared_batch_rows.flow.yaml) | 让 Redis 里的共享 batch id 进入明细层 |
| 68 | 读回共享批次的 Postgres 明细行 | [../../script/tutorials/68_query_shared_batch_detail_rows.lua](../../script/tutorials/68_query_shared_batch_detail_rows.lua) | [../../script/tutorials/68_query_shared_batch_detail_rows.flow.yaml](../../script/tutorials/68_query_shared_batch_detail_rows.flow.yaml) | 单独观察共享批次的明细事实 |
| 69 | 把源 CSV 和 DB 明细行放到一起比 | [../../script/tutorials/69_compare_source_csv_and_db_rows.lua](../../script/tutorials/69_compare_source_csv_and_db_rows.lua) | [../../script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml](../../script/tutorials/69_compare_source_csv_and_db_rows.flow.yaml) | 从“写成功”进入“结果逐行核对” |
| 70 | 生成一份 CSV、Redis、Postgres 三边对账包 | [../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua](../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.lua) | [../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml](../../script/tutorials/70_build_reconciliation_pack_from_csv_redis_db.flow.yaml) | 把三边事实压成一份可复盘的对账结果 |
| 71 | 跑通一次完整的外部系统 round trip | [../../script/tutorials/71_external_system_round_trip.lua](../../script/tutorials/71_external_system_round_trip.lua) | [../../script/tutorials/71_external_system_round_trip.flow.yaml](../../script/tutorials/71_external_system_round_trip.flow.yaml) | 从 CSV 出发，重新串起 Redis、Postgres 和对账输出 |
| 72 | 用同一个批次号重跑同步，但不产生重复数据 | [../../script/tutorials/72_rerun_shared_batch_idempotently.lua](../../script/tutorials/72_rerun_shared_batch_idempotently.lua) | [../../script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml](../../script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml) | 从“能跑一次”继续走到“能安全重跑” |
| 73 | 验证重跑后没有重复行 | [../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua](../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua) | [../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml](../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml) | 把幂等重跑变成可验证结果 |
| 74 | 遇到坏数据时，保留有效行并写出异常台账 | [../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua](../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.lua) | [../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml](../../script/tutorials/74_recover_external_sync_with_anomaly_ledger.flow.yaml) | 从正常批次进入异常恢复 |
| 75 | 给外部同步写入一条审计记录 | [../../script/tutorials/75_write_external_sync_audit_row.lua](../../script/tutorials/75_write_external_sync_audit_row.lua) | [../../script/tutorials/75_write_external_sync_audit_row.flow.yaml](../../script/tutorials/75_write_external_sync_audit_row.flow.yaml) | 把运行态数据和审计留痕拆开 |
| 76 | 读回某个批次的审计历史 | [../../script/tutorials/76_query_external_sync_audit_history.lua](../../script/tutorials/76_query_external_sync_audit_history.lua) | [../../script/tutorials/76_query_external_sync_audit_history.flow.yaml](../../script/tutorials/76_query_external_sync_audit_history.flow.yaml) | 从单条审计进入批次历史视角 |
| 77 | 把审计历史导出成 CSV | [../../script/tutorials/77_export_external_sync_audit_history.lua](../../script/tutorials/77_export_external_sync_audit_history.lua) | [../../script/tutorials/77_export_external_sync_audit_history.flow.yaml](../../script/tutorials/77_export_external_sync_audit_history.flow.yaml) | 把审计结果落盘给人复盘 |
| 78 | 清理最新批次的运行数据，但保留审计 | [../../script/tutorials/78_cleanup_latest_external_batch.lua](../../script/tutorials/78_cleanup_latest_external_batch.lua) | [../../script/tutorials/78_cleanup_latest_external_batch.flow.yaml](../../script/tutorials/78_cleanup_latest_external_batch.flow.yaml) | 引入运行数据清理 |
| 79 | 验证批次清理后，审计仍然保留 | [../../script/tutorials/79_verify_external_batch_cleanup.lua](../../script/tutorials/79_verify_external_batch_cleanup.lua) | [../../script/tutorials/79_verify_external_batch_cleanup.flow.yaml](../../script/tutorials/79_verify_external_batch_cleanup.flow.yaml) | 验证“删运行态，不删审计” |
| 80 | 跑通一条完整的外部同步生命周期 | [../../script/tutorials/80_external_sync_lifecycle_round_trip.lua](../../script/tutorials/80_external_sync_lifecycle_round_trip.lua) | [../../script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml](../../script/tutorials/80_external_sync_lifecycle_round_trip.flow.yaml) | 把创建、审计、清理、验证重新串成闭环 |
| 81 | 从生命周期 CSV 里读回批次证据 | [../../script/tutorials/81_read_lifecycle_evidence.lua](../../script/tutorials/81_read_lifecycle_evidence.lua) | [../../script/tutorials/81_read_lifecycle_evidence.flow.yaml](../../script/tutorials/81_read_lifecycle_evidence.flow.yaml) | 从 `Lesson 80` 的证据重新读回原始批次事实 |
| 82 | 按生命周期证据回放一个新批次 | [../../script/tutorials/82_replay_batch_from_lifecycle_evidence.lua](../../script/tutorials/82_replay_batch_from_lifecycle_evidence.lua) | [../../script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml](../../script/tutorials/82_replay_batch_from_lifecycle_evidence.flow.yaml) | 从清理后的证据重建一条 replay 批次 |
| 83 | 用生命周期证据验证回放批次 | [../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua](../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.lua) | [../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml](../../script/tutorials/83_verify_replay_batch_against_lifecycle_evidence.flow.yaml) | 对齐生命周期、Redis 和 Postgres 三边事实 |
| 84 | 给回放批次补写一条审计记录 | [../../script/tutorials/84_write_replay_audit_row.lua](../../script/tutorials/84_write_replay_audit_row.lua) | [../../script/tutorials/84_write_replay_audit_row.flow.yaml](../../script/tutorials/84_write_replay_audit_row.flow.yaml) | 让 replay 批次拥有独立审计留痕 |
| 85 | 把原批次和回放批次的审计导出成对照 CSV | [../../script/tutorials/85_export_original_and_replay_audits.lua](../../script/tutorials/85_export_original_and_replay_audits.lua) | [../../script/tutorials/85_export_original_and_replay_audits.flow.yaml](../../script/tutorials/85_export_original_and_replay_audits.flow.yaml) | 把原链和 replay 链的审计放在一起复盘 |
| 86 | 生成一份回放后的对账包 | [../../script/tutorials/86_build_post_replay_reconciliation_pack.lua](../../script/tutorials/86_build_post_replay_reconciliation_pack.lua) | [../../script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml](../../script/tutorials/86_build_post_replay_reconciliation_pack.flow.yaml) | 压缩 replay 后的多边事实结果 |
| 87 | 生成一份交接 artifact manifest | [../../script/tutorials/87_build_handoff_artifact_manifest.lua](../../script/tutorials/87_build_handoff_artifact_manifest.lua) | [../../script/tutorials/87_build_handoff_artifact_manifest.flow.yaml](../../script/tutorials/87_build_handoff_artifact_manifest.flow.yaml) | 把关键产物组织成交接目录 |
| 88 | 把交接 manifest 整理成交付摘要 | [../../script/tutorials/88_build_handoff_summary.lua](../../script/tutorials/88_build_handoff_summary.lua) | [../../script/tutorials/88_build_handoff_summary.flow.yaml](../../script/tutorials/88_build_handoff_summary.flow.yaml) | 生成一份更适合快速扫读的 summary |
| 89 | 生成发布前检查清单 | [../../script/tutorials/89_build_pre_release_checklist.lua](../../script/tutorials/89_build_pre_release_checklist.lua) | [../../script/tutorials/89_build_pre_release_checklist.flow.yaml](../../script/tutorials/89_build_pre_release_checklist.flow.yaml) | 把交接包变成可核对的 checklist |
| 90 | 跑通一条“生命周期证据 -> 回放 -> 交接包”的完整 round trip | [../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua](../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.lua) | [../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml](../../script/tutorials/90_handoff_round_trip_from_lifecycle_evidence.flow.yaml) | 把回放、审计、交接收成一条完整交付链 |
| 91 | 读交接 manifest，识别每份产物的角色 | [../../script/tutorials/91_read_handoff_manifest_roles.lua](../../script/tutorials/91_read_handoff_manifest_roles.lua) | [../../script/tutorials/91_read_handoff_manifest_roles.flow.yaml](../../script/tutorials/91_read_handoff_manifest_roles.flow.yaml) | 先把交接包里的产物分清角色 |
| 92 | 把交接产物整理成模板目录 | [../../script/tutorials/92_build_template_artifact_catalog.lua](../../script/tutorials/92_build_template_artifact_catalog.lua) | [../../script/tutorials/92_build_template_artifact_catalog.flow.yaml](../../script/tutorials/92_build_template_artifact_catalog.flow.yaml) | 给产物分配稳定槽位和环境变量入口 |
| 93 | 把交接链整理成 Input -> Process -> Output 模板 | [../../script/tutorials/93_build_input_process_output_template.lua](../../script/tutorials/93_build_input_process_output_template.lua) | [../../script/tutorials/93_build_input_process_output_template.flow.yaml](../../script/tutorials/93_build_input_process_output_template.flow.yaml) | 从三段式视角整理模板 |
| 94 | 把交接链整理成 Collect -> Verify -> Save 模板 | [../../script/tutorials/94_build_collect_verify_save_template.lua](../../script/tutorials/94_build_collect_verify_save_template.lua) | [../../script/tutorials/94_build_collect_verify_save_template.flow.yaml](../../script/tutorials/94_build_collect_verify_save_template.flow.yaml) | 从 review 视角整理模板 |
| 95 | 把交接链整理成 Replay -> Audit -> Handoff 模板 | [../../script/tutorials/95_build_replay_audit_handoff_template.lua](../../script/tutorials/95_build_replay_audit_handoff_template.lua) | [../../script/tutorials/95_build_replay_audit_handoff_template.flow.yaml](../../script/tutorials/95_build_replay_audit_handoff_template.flow.yaml) | 保留业务语义重新抽模板 |
| 96 | 把几份模板整理成统一索引 | [../../script/tutorials/96_build_template_index.lua](../../script/tutorials/96_build_template_index.lua) | [../../script/tutorials/96_build_template_index.flow.yaml](../../script/tutorials/96_build_template_index.flow.yaml) | 形成一个可浏览的模板目录 |
| 97 | 验证模板索引仍然覆盖完整交接链 | [../../script/tutorials/97_verify_template_covers_handoff_chain.lua](../../script/tutorials/97_verify_template_covers_handoff_chain.lua) | [../../script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml](../../script/tutorials/97_verify_template_covers_handoff_chain.flow.yaml) | 确认模板没有把原链路丢掉 |
| 98 | 生成一份“场景 -> 模板”的学习矩阵 | [../../script/tutorials/98_build_template_lesson_matrix.lua](../../script/tutorials/98_build_template_lesson_matrix.lua) | [../../script/tutorials/98_build_template_lesson_matrix.flow.yaml](../../script/tutorials/98_build_template_lesson_matrix.flow.yaml) | 帮新人按场景选模板 |
| 99 | 给模板包生成发布前检查清单 | [../../script/tutorials/99_build_template_preflight_checklist.lua](../../script/tutorials/99_build_template_preflight_checklist.lua) | [../../script/tutorials/99_build_template_preflight_checklist.flow.yaml](../../script/tutorials/99_build_template_preflight_checklist.flow.yaml) | 给模板包建立自己的发布门槛 |
| 100 | 跑通一条“交接产物 -> 模板包”的完整 round trip | [../../script/tutorials/100_template_round_trip_from_handoff_artifacts.lua](../../script/tutorials/100_template_round_trip_from_handoff_artifacts.lua) | [../../script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml](../../script/tutorials/100_template_round_trip_from_handoff_artifacts.flow.yaml) | 把交接产物收成可复用模板包 |
| 101 | 先确认模板发布卡片真的在页面上 | [../../script/tutorials/101_assert_visible_template_release_card.lua](../../script/tutorials/101_assert_visible_template_release_card.lua) | [../../script/tutorials/101_assert_visible_template_release_card.flow.yaml](../../script/tutorials/101_assert_visible_template_release_card.flow.yaml) | 把 Lesson 100 的模板包正式带进发布检查页 |
| 102 | 继续确认模板发布状态文字对不对 | [../../script/tutorials/102_assert_text_template_release_status.lua](../../script/tutorials/102_assert_text_template_release_status.lua) | [../../script/tutorials/102_assert_text_template_release_status.flow.yaml](../../script/tutorials/102_assert_text_template_release_status.flow.yaml) | 从“在不在”继续走到“文案对不对” |
| 103 | 用 `retry` 跑通模板发布 gate | [../../script/tutorials/103_retry_template_release_gate.lua](../../script/tutorials/103_retry_template_release_gate.lua) | [../../script/tutorials/103_retry_template_release_gate.flow.yaml](../../script/tutorials/103_retry_template_release_gate.flow.yaml) | 练习第一次不通过、第二次才通过的 gate |
| 104 | 用 `wait_until` 等模板发布检查完成 | [../../script/tutorials/104_wait_until_template_release_ready.lua](../../script/tutorials/104_wait_until_template_release_ready.lua) | [../../script/tutorials/104_wait_until_template_release_ready.flow.yaml](../../script/tutorials/104_wait_until_template_release_ready.flow.yaml) | 练习异步 ready 而不是固定 sleep |
| 105 | 用 `on_error` 接住模板发布校验失败 | [../../script/tutorials/105_on_error_template_release_validation.lua](../../script/tutorials/105_on_error_template_release_validation.lua) | [../../script/tutorials/105_on_error_template_release_validation.flow.yaml](../../script/tutorials/105_on_error_template_release_validation.flow.yaml) | 把失败分支正式纳入模板发布流程 |
| 106 | 等一条延迟出现的发布说明项 | [../../script/tutorials/106_wait_for_delayed_release_note.lua](../../script/tutorials/106_wait_for_delayed_release_note.lua) | [../../script/tutorials/106_wait_for_delayed_release_note.flow.yaml](../../script/tutorials/106_wait_for_delayed_release_note.flow.yaml) | 练习等待“元素本身”晚一点出现 |
| 107 | 用 `retry` 接住一次偶发失败点击 | [../../script/tutorials/107_retry_flaky_publish_click.lua](../../script/tutorials/107_retry_flaky_publish_click.lua) | [../../script/tutorials/107_retry_flaky_publish_click.flow.yaml](../../script/tutorials/107_retry_flaky_publish_click.flow.yaml) | 把 retry 从 gate 推进到 click 本身 |
| 108 | `reload` 之后再验证一次恢复结果 | [../../script/tutorials/108_reload_and_retry_release_recovery.lua](../../script/tutorials/108_reload_and_retry_release_recovery.lua) | [../../script/tutorials/108_reload_and_retry_release_recovery.flow.yaml](../../script/tutorials/108_reload_and_retry_release_recovery.flow.yaml) | 把刷新页面也纳入恢复链 |
| 109 | 给模板发布页留一份调试证据包 | [../../script/tutorials/109_template_release_artifact_pack.lua](../../script/tutorials/109_template_release_artifact_pack.lua) | [../../script/tutorials/109_template_release_artifact_pack.flow.yaml](../../script/tutorials/109_template_release_artifact_pack.flow.yaml) | 保存整页图、卡片图和 HTML |
| 110 | 跑通一条完整的模板发布稳定性 round trip | [../../script/tutorials/110_template_release_robustness_round_trip.lua](../../script/tutorials/110_template_release_robustness_round_trip.lua) | [../../script/tutorials/110_template_release_robustness_round_trip.flow.yaml](../../script/tutorials/110_template_release_robustness_round_trip.flow.yaml) | 把断言、等待、重试、恢复和证据留存收成一条线 |
| 111 | 先用 `tsplay.list_actions` 看清 MCP 到底能做什么 | - | - | 命令型章节，先建立 MCP 能力地图 |
| 112 | 用 `flow_schema` 和 `flow_examples` 看清 Flow 长什么样 | - | - | 命令型章节，先建立 schema 和 example 心智模型 |
| 113 | 先用 `observe_page` 观察模板发布页 | - | [../../script/tutorials/113_mcp_observe_page_template_release.args.json](../../script/tutorials/113_mcp_observe_page_template_release.args.json) | 这一段开始切到 `mcp-tool + args.json` 方式 |
| 114 | 用 `draft_flow` 把 observation 变成第一份 Flow 草稿 | - | [../../script/tutorials/114_mcp_draft_flow_template_release.args.json](../../script/tutorials/114_mcp_draft_flow_template_release.args.json) | 直接复用 `Lesson 113` 的 observation 输出 |
| 115 | 先校验草稿，再决定能不能运行 | - | [../../script/tutorials/115_mcp_validate_drafted_template_release.args.json](../../script/tutorials/115_mcp_validate_drafted_template_release.args.json) | 通过 `@jsonpathfile` 直接抽取 `draft.flow_yaml` |
| 116 | 运行刚刚起草并校验过的 Flow | - | [../../script/tutorials/116_mcp_run_drafted_template_release.args.json](../../script/tutorials/116_mcp_run_drafted_template_release.args.json) | 把 `observe -> draft -> validate` 接进真实执行 |
| 117 | 先故意跑坏一次，再生成 repair context | - | [../../script/tutorials/117_mcp_repair_flow_context_template_release.args.json](../../script/tutorials/117_mcp_repair_flow_context_template_release.args.json) | 会同时用到破坏版 Flow 和失败运行结果 |
| 118 | 用 repair context 生成真正可用的修复请求 | - | [../../script/tutorials/118_mcp_repair_flow_template_release.args.json](../../script/tutorials/118_mcp_repair_flow_template_release.args.json) | 从失败上下文进入统一 repair request |
| 119 | 把 `observe -> draft -> validate -> run -> repair` 串成一条线 | - | - | 复盘型章节，把前面 6 步重新串成完整 MCP 主线 |
| 120 | 用 `finalize_flow` 收成一份更短的默认入口 | - | [../../script/tutorials/120_mcp_finalize_flow_template_release.args.json](../../script/tutorials/120_mcp_finalize_flow_template_release.args.json) | 从 observation 直接收成更接近交付态的结果 |
| 121 | 用 `allow_lua` 放行一条最小 Lua Flow | - | [../../script/tutorials/121_security_allow_lua.flow.yaml](../../script/tutorials/121_security_allow_lua.flow.yaml) | 先看默认边界拦截，再只打开 `allow_lua` |
| 122 | 用 `allow_http` 放行一条最小 HTTP Flow | - | [../../script/tutorials/122_security_allow_http.flow.yaml](../../script/tutorials/122_security_allow_http.flow.yaml) | 继续练习 blocked / allowed 对照 |
| 123 | 用 `allow_file_access` 放行一条最小文件输出 Flow | - | [../../script/tutorials/123_security_allow_file_access.flow.yaml](../../script/tutorials/123_security_allow_file_access.flow.yaml) | 为后面的本地 Flow / MCP 对照打基础 |
| 124 | 用 `allow_browser_state` 放行浏览器状态动作 | - | [../../script/tutorials/124_security_allow_browser_state.flow.yaml](../../script/tutorials/124_security_allow_browser_state.flow.yaml) | 把 Cookie / Storage State 单独当一类边界理解 |
| 125 | 用 `allow_redis` 放行 Redis 动作 | - | [../../script/tutorials/125_security_allow_redis.flow.yaml](../../script/tutorials/125_security_allow_redis.flow.yaml) | 从浏览器状态推进到外部系统状态 |
| 126 | 用 `allow_database` 放行数据库动作 | - | [../../script/tutorials/126_security_allow_database.flow.yaml](../../script/tutorials/126_security_allow_database.flow.yaml) | 把 DB 动作纳入同样的最小授权思路 |
| 127 | 对比本地 Flow 和 MCP 的权限边界 | - | - | 命令型章节，复用 `Lesson 123` 的 Flow 和 blocked / allowed 结果 |
| 128 | 为什么教程不能跳过权限边界 | - | - | 复盘型章节，重新解释为什么要先学边界再学放权 |
| 129 | 理解 `security_preset` 和显式 `allow_*` 覆盖 | - | [../../script/tutorials/129_mcp_validate_file_access_browser_write.args.json](../../script/tutorials/129_mcp_validate_file_access_browser_write.args.json) | 从单个 allow 进入 preset 和 override 组合策略 |
| 130 | 完成安全边界模块的第一轮 checkpoint | - | - | 收口型章节，把 `121-129` 重新整理成高级阶段 checkpoint |

## 完整课程体系地图

| 层级 | 面向谁 | 目标 | 入口 |
| --- | --- | --- | --- |
| 新手教程 | 第一天接触 TSPlay 的同学 | 跑通二进制、能理解 Lua / Flow 的基本关系 | [track-newbie.md](track-newbie.md) |
| 初级教程 | 已经能跑 demo、想进入真实业务动作的人 | 把文件、变量、控制流、HTTP、Redis、DB 基础连起来 | [track-junior.md](track-junior.md) |
| 中级教程 | 已能独立写基础 Flow 的实施 / QA / 开发 | 形成可复用模板、数据驱动流程、健壮性设计、MCP 基础 | [track-intermediate.md](track-intermediate.md) |
| 高级教程 | 要交付、评审、培训、集成 TSPlay 的同学 | 建立规范、架构、安全边界、发布包、repair 和课程演进能力 | [track-advanced.md](track-advanced.md) |
| 160 次迭代路线图 | 希望长期持续推进课程建设的人 | 把教程演进拆成 160 个不跳跃的迭代点 | [iteration-roadmap-160.md](iteration-roadmap-160.md) |

## 开始前

先在仓库根目录执行：

```bash
go mod download
```

如果你想直接用 `tsplay` 命令测试，先构建一次：

```bash
go build -o tsplay .
```

构建完成后，下面这些命令都可以直接跑：

```bash
./tsplay -script script/tutorials/01_hello_world.lua
./tsplay -flow script/tutorials/01_hello_world.flow.yaml
./tsplay -action cli
./tsplay -action file-srv -addr :8000
./tsplay -action mcp-tool -tool tsplay.list_actions
```

如果你暂时不想构建二进制，也可以继续用：

```bash
go run . -script script/tutorials/01_hello_world.lua
go run . -flow script/tutorials/01_hello_world.flow.yaml
go run . -action mcp-tool -tool tsplay.list_actions
```

补充说明：

- 构建出来的 `./tsplay` 已经内置了 `ReadMe.md`、`docs/`、`script/`、`demo/`
- 所以即使你把单个二进制拿到别的目录，也还能直接跑这些示例路径
- 如果你想把这套参考资料释放出来，可以执行：

```bash
./tsplay -action list-assets
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

Lesson 01、Lesson 08、Lesson 14、Lesson 24、Lesson 25、Lesson 111、Lesson 112、Lesson 121 到 Lesson 130 都可以直接开始。  
从 Lesson 02 到 Lesson 05、Lesson 09 到 Lesson 12、Lesson 16 到 Lesson 23、Lesson 26 到 Lesson 38、Lesson 42、Lesson 44 到 Lesson 57、Lesson 101 到 Lesson 120，建议另开一个终端，在仓库根目录启动 TSPlay 内置静态文件服务：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

如果你已经只有单个 `./tsplay` 二进制、没有完整仓库目录，也可以直接执行上面这条命令。  
此时服务的是二进制里内置的 `demo/`、`docs/`、`script/` 资源。

然后通过下面这些本地地址访问仓库里的 demo 页面：

- `http://127.0.0.1:8000/demo/demo.html`
- `http://127.0.0.1:8000/demo/demo_alt.html`
- `http://127.0.0.1:8000/demo/tables.html`
- `http://127.0.0.1:8000/demo/extract.html`
- `http://127.0.0.1:8000/demo/retry_wait_until.html`
- `http://127.0.0.1:8000/demo/upload.html`
- `http://127.0.0.1:8000/demo/multi_upfile.html`
- `http://127.0.0.1:8000/demo/download.html`
- `http://127.0.0.1:8000/demo/import_workflow.html`
- `http://127.0.0.1:8000/demo/import_workflow.html?login=1`
- `http://127.0.0.1:8000/demo/session_lab.html`
- `http://127.0.0.1:8000/demo/session_import_workflow.html`
- `http://127.0.0.1:8000/demo/debug_artifacts.html`
- `http://127.0.0.1:8000/demo/template_release_lab.html`
- `http://127.0.0.1:8000/demo/data/order_summary.json`

如果你不是在仓库根目录启动服务，可以显式指定：

```bash
./tsplay -action file-srv -addr :8000 -serve-root /path/to/tsplay
```

从 Lesson 06 起，开始接 Redis 和数据库。  
为了让命令更好抄，可以直接复用这些辅助文件：

- Redis 环境变量示例：[../../script/tutorials/env/06_redis_example.sh](../../script/tutorials/env/06_redis_example.sh)
- Postgres 环境变量示例：[../../script/tutorials/env/07_reporting_pg_example.sh](../../script/tutorials/env/07_reporting_pg_example.sh)
- Postgres 初始化 SQL：[../../script/tutorials/sql/07_reporting_pg.sql](../../script/tutorials/sql/07_reporting_pg.sql)
- 导出结果同步到 Postgres 的初始化 SQL：[../../script/tutorials/sql/61_reporting_import_sync.sql](../../script/tutorials/sql/61_reporting_import_sync.sql)
- 外部同步审计表 SQL：[../../script/tutorials/sql/75_reporting_import_audit.sql](../../script/tutorials/sql/75_reporting_import_audit.sql)

从 Lesson 58 到 Lesson 80，默认会继续复用 `Lesson 57` 产出的导出 CSV。  
如果你在前一节跑的是 Lua 版本，记得先把：

```bash
export TSPLAY_IMPORTED_REPORT=artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv
```

切到对应路径再继续。

从 Lesson 81 到 Lesson 90，默认会继续复用 `Lesson 80` 产出的生命周期 CSV。  
如果你在前一节跑的是 Lua 版本，记得先把：

```bash
export TSPLAY_LIFECYCLE_FILE=artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv
```

切到对应路径再继续。  
其中 `Lesson 83-89` 还会继续复用前一节输出的 replay / audit / manifest CSV，如果你跑的是 Lua 版本，也记得把：

```bash
export TSPLAY_REPLAY_FILE=artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv
export TSPLAY_AUDIT_COMPARE_FILE=artifacts/tutorials/85-export-original-and-replay-audits-lua.csv
export TSPLAY_RECONCILIATION_FILE=artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.csv
export TSPLAY_MANIFEST_FILE=artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv
```

切到对应路径再继续。

从 Lesson 91 到 Lesson 100，默认会继续复用 `Lesson 87`、`Lesson 89`、`Lesson 90` 产出的 CSV。  
如果你在前一节跑的是 Lua 版本，记得先把：

```bash
export TSPLAY_MANIFEST_FILE=artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv
export TSPLAY_HANDOFF_FILE=artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv
export TSPLAY_ROLE_FILE=artifacts/tutorials/91-read-handoff-manifest-roles-lua.csv
export TSPLAY_TEMPLATE_CATALOG_FILE=artifacts/tutorials/92-build-template-artifact-catalog-lua.csv
export TSPLAY_RUNTIME_CHECKLIST_FILE=artifacts/tutorials/89-build-pre-release-checklist-lua.csv
export TSPLAY_INPUT_PROCESS_OUTPUT_FILE=artifacts/tutorials/93-build-input-process-output-template-lua.csv
export TSPLAY_COLLECT_VERIFY_SAVE_FILE=artifacts/tutorials/94-build-collect-verify-save-template-lua.csv
export TSPLAY_REPLAY_AUDIT_HANDOFF_FILE=artifacts/tutorials/95-build-replay-audit-handoff-template-lua.csv
export TSPLAY_TEMPLATE_INDEX_FILE=artifacts/tutorials/96-build-template-index-lua.csv
export TSPLAY_TEMPLATE_VERIFICATION_FILE=artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.csv
export TSPLAY_TEMPLATE_LESSON_MATRIX_FILE=artifacts/tutorials/98-build-template-lesson-matrix-lua.csv
export TSPLAY_TEMPLATE_PREFLIGHT_FILE=artifacts/tutorials/99-build-template-preflight-checklist-lua.csv
```

切到对应路径再继续。

从 Lesson 101 到 Lesson 110，默认会继续复用同一张本地模板发布练习页。  
如果你想显式覆盖地址，可以先切一下：

```bash
export TSPLAY_TEMPLATE_RELEASE_URL=http://127.0.0.1:8000/demo/template_release_lab.html
```

从 Lesson 113 到 Lesson 120，默认还会继续复用同一张模板发布练习页，并把 MCP 输出写进 `artifacts/tutorials/`。  
如果你准备按 lesson 顺序推进，建议先把这个目录建好：

```bash
mkdir -p artifacts/tutorials
```

补充说明：

- `Lesson 58-60` 主要需要 Redis
- `Lesson 61-64` 主要需要 Postgres
- `Lesson 65-73` 需要 Redis 和 Postgres 同时准备好
- `Lesson 74` 主要需要 Postgres
- `Lesson 75-90` 需要 Redis 和 Postgres 同时准备好，并执行过审计表 SQL
- `Lesson 91-100` 主要继续消费本地 CSV 产物；如果前一段已经跑完，通常不需要再新开 Redis / Postgres 连接
- `Lesson 101-110` 主要继续复用本地静态文件服务和 `demo/template_release_lab.html`，通常不再需要新开 Redis / Postgres 连接
- `Lesson 111-120` 主要继续复用同一张模板发布练习页，但切到 `mcp-tool` 路径，重点练习 observation、draft、validate、run、repair 和 finalize
- `Lesson 121-130` 主要继续练习 `validate_flow`、`security_preset` 和边界对照，通常不需要本地静态文件服务，也不要求真的起 Redis / Postgres / DB 服务

从 Lesson 13 起，开始接本地文件输入输出。  
如果你只有单个 `./tsplay` 二进制，记得先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

这样 `demo/data/*.csv`、`demo/data/*.xlsx`、`demo/data/*.pdf`、`demo/data/*.png` 这些示例文件才会真正出现在本地磁盘上，方便 `read_csv`、`read_excel`、`upload_file`、`upload_multiple_files`、`download_file` 这些动作直接复用。

## 学的时候重点看什么

- `Lua` 更像“我现在一步一步怎么做”
- `Flow` 更像“这个流程本身是什么，可以怎么被 review、复用和交给 AI 生成”
- 同一个功能都先看 `Lua` 再看 `Flow`，更容易体会两者的边界
- 先走“今天能跑通”的实战线，再进入“完整进阶教程体系”，学习曲线会更稳
- 如果你要长期建设教程，直接结合 [160 次递进迭代路线图](iteration-roadmap-160.md) 和 [教程持续进化手册](evolution-playbook.md)

## 两个小提醒

- 运行带浏览器的 `Lua` 脚本时，如果不带 `-headless`，脚本执行完成后浏览器会继续保持打开，方便你观察页面；结束时按 `Ctrl+C` 即可
- 这些教程默认把输出写到 `artifacts/tutorials/`，这样不会把练习产物混进仓库版本里
