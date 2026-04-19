# TSPlay 完整课程总览

这份总览不是“再加几篇教程”，而是把 TSPlay 教程升级成一套能持续演进的学习系统。

目标有三个：

- 让第一次接触 TSPlay 的人也能按顺序学，不会一上来就跳到 MCP、repair、数据库事务这些高级主题
- 让已经能跑 demo 的同学，知道下一步应该补哪一块，不会东学一点西学一点
- 让课程作者、实施同学、讲师有一条可以不断迭代的主线，而不是每次临时加几节文档

## 课程总原则

### 1. 先本地，后外部系统

先用本地页面、本地 JSON、本地二进制内置资源建立确定性。  
只有当前一层完全走通后，再进入 Redis、数据库、MCP、会话复用、repair。

### 2. 先可运行，后可抽象

先学会“怎么跑”，再学会“怎么抽象成 Flow”。  
所以绝大多数主题都建议先看 `Lua`，再看 `Flow`。

### 3. 先单点，后组合

先把一个动作学透，再把多个动作串起来。  
例如先学 `extract_text`，再学 `extract_text + set_var + if`，最后再学带 `retry / on_error` 的组合流程。

### 4. 先稳定页，后不稳定页

先用仓库自带 `demo/` 和本地 JSON 建基本功。  
这样教程重点落在 TSPlay 本身，而不是被外部页面变动打断。

### 5. 先交付物，后“懂了”

每一节都要有交付物：

- 一个脚本
- 一条 Flow
- 一份 JSON / CSV / 截图 / trace
- 一段复盘说明

只有这样，学习过程才是可评估、可复现、可迭代的。

## 四层课程结构

## 第一层：新手教程

入口：
[track-newbie.md](track-newbie.md)

适合谁：

- 第一天接触 TSPlay 的人
- 还不熟悉 Lua / Flow 的人
- 需要通过最小成功体验建立信心的人

这一层解决的问题：

- `tsplay` 二进制怎么跑
- 内置资源怎么用
- 本地静态服务怎么开
- `Lua` 和 `Flow` 到底是什么关系
- 页面文本、HTML、表格、JSON 到底怎么拿

完成标志：

- 能独立跑通当前新手本地练习链路：`Lesson 01-05` 和 `Lesson 08-12`
- 能解释为什么“先 Lua 再 Flow”更容易上手
- 能把结果写到 `artifacts/tutorials/`

## 第二层：初级教程

入口：
[track-junior.md](track-junior.md)

适合谁：

- 已经能跑本地 demo
- 想开始接文件、变量、控制流、外部系统基础
- 想从“会跑”走向“能做一个小业务流程”

这一层解决的问题：

- 文件输入输出怎么处理
- 变量、分支、循环、断言怎么组织
- HTTP、Redis、数据库基础动作怎么接入
- 会话、artifact、输出目录怎么管理

完成标志：

- 能独立写出一个包含变量和控制流的基础 Flow
- 能接通一个 Redis 或 Postgres 最小例子
- 能解释一个流程的输入、输出和失败现场

## 第三层：中级教程

入口：
[track-intermediate.md](track-intermediate.md)

适合谁：

- 已能独立写基础 Flow
- 开始关注复用、稳定性、可维护性
- 希望把单条脚本提升成项目资产

这一层解决的问题：

- 如何把临时脚本变成模板
- 如何用 CSV / Excel / foreach 做数据驱动
- 如何用 `retry / on_error / wait_until` 做健壮性设计
- 如何理解 `observe -> draft -> validate -> run -> repair`

完成标志：

- 能维护一组可复用 Flow 模板
- 能对 Flaky 流程做最小修复设计
- 能解释什么时候应该引入 MCP，而不是继续手写脚本

## 第四层：高级教程

入口：
[track-advanced.md](track-advanced.md)

适合谁：

- 要交付 TSPlay 能力到团队 / 项目 / 客户环境
- 要做规范、评审、培训、集成、发布
- 希望把 TSPlay 作为“系统能力”而不是“个人脚本工具”

这一层解决的问题：

- 安全边界和授权怎么设计
- 大型 Flow 怎么分层、命名、评审、版本化
- 内置资产、二进制发布、交付包怎么组织
- 如何把课程长期演进，而不是停在第一版

完成标志：

- 能制定团队教程结构和代码评审规则
- 能解释 TSPlay 在交付链路里的位置
- 能持续推进教程体系，而不是一次性写完就停

## 当前已落地的“立即可跑”部分

今天仓库里已经可直接运行的是：

- [Lesson 01](01-hello-world.md)
- [Lesson 02](02-local-page-select-option.md)
- [Lesson 03](03-capture-table.md)
- [Lesson 04](04-extract-text-and-html.md)
- [Lesson 05](05-http-request-json.md)
- [Lesson 06](06-redis-round-trip.md)
- [Lesson 07](07-db-postgres-basics.md)
- [Lesson 08](08-bundled-assets-and-artifacts.md)
- [Lesson 09](09-local-demo-anatomy.md)
- [Lesson 10](10-assert-page-state.md)
- [Lesson 11](11-select-another-option.md)
- [Lesson 12](12-custom-json-output.md)
- [Lesson 13](13-read-csv-basics.md)
- [Lesson 14](14-write-csv-basics.md)
- [Lesson 15](15-read-transform-write-csv.md)
- [Lesson 16](16-retry-flaky-action.md)
- [Lesson 17](17-wait-until-ready.md)
- [Lesson 18](18-upload-single-file.md)
- [Lesson 19](19-upload-multiple-files.md)
- [Lesson 20](20-download-report.md)
- [Lesson 21](21-if-optional-login.md)
- [Lesson 22](22-foreach-batch-import-csv.md)
- [Lesson 23](23-on-error-import-recovery.md)
- [Lesson 24](24-read-excel-basics.md)
- [Lesson 25](25-read-excel-range-headers.md)
- [Lesson 26](26-foreach-batch-import-excel.md)
- [Lesson 27](27-on-error-import-excel-writeback.md)
- [Lesson 28](28-inspect-storage-state.md)
- [Lesson 29](29-read-cookies-string.md)
- [Lesson 30](30-browser-state-snapshot-pack.md)
- [Lesson 31](31-full-page-screenshot.md)
- [Lesson 32](32-element-screenshot.md)
- [Lesson 33](33-save-html-basics.md)
- [Lesson 34](34-debug-artifact-pack.md)
- [Lesson 35](35-error-evidence-pack.md)
- [Lesson 36](36-save-storage-state.md)
- [Lesson 37](37-load-saved-storage-state.md)
- [Lesson 38](38-verify-loaded-storage-state.md)
- [Lesson 39](39-storage-state-round-trip.md)
- [Lesson 40](40-save-named-session.md)
- [Lesson 41](41-inspect-named-session.md)
- [Lesson 42](42-use-named-session.md)
- [Lesson 43](43-delete-named-session.md)
- [Lesson 44](44-session-import-with-login.md)
- [Lesson 45](45-storage-state-auth-import.md)
- [Lesson 46](46-save-import-session.md)
- [Lesson 47](47-use-session-import-single.md)
- [Lesson 48](48-use-session-batch-import-csv.md)
- [Lesson 49](49-use-session-import-recovery-csv.md)
- [Lesson 50](50-use-session-batch-import-excel.md)
- [Lesson 51](51-use-session-import-recovery-excel.md)
- [Lesson 52](52-use-session-capture-import-table.md)
- [Lesson 53](53-use-session-capture-import-table-to-csv.md)
- [Lesson 54](54-use-session-download-import-report.md)
- [Lesson 55](55-use-session-download-import-report-readback.md)
- [Lesson 56](56-use-session-compare-table-and-download.md)
- [Lesson 57](57-use-session-import-export-round-trip.md)
- [Lesson 58](58-sync-import-report-summary-to-redis.md)
- [Lesson 59](59-save-import-batch-key-to-redis.md)
- [Lesson 60](60-read-latest-import-batch-from-redis.md)
- [Lesson 61](61-db-insert-import-batch-summary.md)
- [Lesson 62](62-db-query-import-batch-summaries.md)
- [Lesson 63](63-db-upsert-import-batch-summary.md)
- [Lesson 64](64-db-transaction-import-batch-and-rows.md)
- [Lesson 65](65-sync-latest-redis-batch-to-postgres-summary.md)
- [Lesson 66](66-query-shared-batch-summary-from-redis-and-postgres.md)
- [Lesson 67](67-transaction-store-shared-batch-rows.md)
- [Lesson 68](68-query-shared-batch-detail-rows.md)
- [Lesson 69](69-compare-source-csv-and-db-rows.md)
- [Lesson 70](70-build-reconciliation-pack-from-csv-redis-db.md)
- [Lesson 71](71-external-system-round-trip.md)

这 71 节是完整课程的起点，不是终点。  
后续扩展的依据统一放在：

- [160 次递进迭代路线图](iteration-roadmap-160.md)
- [教程持续进化手册](evolution-playbook.md)

## 推荐使用方式

如果你是第一次接触 TSPlay：

1. 先走 [README.md](README.md) 里的本地基础教程
2. 再读 [track-newbie.md](track-newbie.md)
3. 跑完之后再看 [iteration-roadmap-160.md](iteration-roadmap-160.md) 的新手阶段

如果你已经有一点自动化经验：

1. 先扫一遍 [track-junior.md](track-junior.md)
2. 挑和自己项目最接近的块开始补
3. 每完成一个块，就回到路线图继续下一迭代

如果你是课程作者或交付负责人：

1. 先读这份总览
2. 再读 [track-advanced.md](track-advanced.md)
3. 最后用 [evolution-playbook.md](evolution-playbook.md) 约束后续文档持续生长的方式
