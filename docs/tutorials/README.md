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
```

如果你暂时不想构建二进制，也可以继续用：

```bash
go run . -script script/tutorials/01_hello_world.lua
go run . -flow script/tutorials/01_hello_world.flow.yaml
```

补充说明：

- 构建出来的 `./tsplay` 已经内置了 `ReadMe.md`、`docs/`、`script/`、`demo/`
- 所以即使你把单个二进制拿到别的目录，也还能直接跑这些示例路径
- 如果你想把这套参考资料释放出来，可以执行：

```bash
./tsplay -action list-assets
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

Lesson 01、Lesson 08、Lesson 14、Lesson 24 和 Lesson 25 可以直接开始。  
从 Lesson 02 到 Lesson 05、Lesson 09 到 Lesson 12、Lesson 16 到 Lesson 23、Lesson 26 到 Lesson 38、Lesson 42、Lesson 44 到 Lesson 57，建议另开一个终端，在仓库根目录启动 TSPlay 内置静态文件服务：

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

从 Lesson 58 到 Lesson 71，默认会继续复用 `Lesson 57` 产出的导出 CSV。  
如果你在前一节跑的是 Lua 版本，记得先把：

```bash
export TSPLAY_IMPORTED_REPORT=artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv
```

切到对应路径再继续。

补充说明：

- `Lesson 58-60` 主要需要 Redis
- `Lesson 61-64` 主要需要 Postgres
- `Lesson 65-71` 需要 Redis 和 Postgres 同时准备好

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
