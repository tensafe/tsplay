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
从 Lesson 02 到 Lesson 05、Lesson 09 到 Lesson 12、Lesson 16 到 Lesson 23、Lesson 26 到 Lesson 27，建议另开一个终端，在仓库根目录启动 TSPlay 内置静态文件服务：

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

- 运行带浏览器的 `Lua` 脚本时，脚本执行完成后浏览器会继续保持打开，方便你观察页面；结束时按 `Ctrl+C` 即可
- 这些教程默认把输出写到 `artifacts/tutorials/`，这样不会把练习产物混进仓库版本里
