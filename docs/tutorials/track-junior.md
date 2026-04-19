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
- [Lesson 06: Redis 基础读写和计数](06-redis-round-trip.md)
- [Lesson 07: Postgres 基础查询与写入](07-db-postgres-basics.md)

这一层现在已经有了一条比较完整的最小链路。  
建议顺序是：

1. `Lesson 13-15` 先把文件输入输出跑顺
2. `Lesson 16-17` 再把 `retry` / `wait_until` 吃透
3. `Lesson 18-20` 再把上传 / 下载动作接上
4. `Lesson 21-23` 再把 `if` / `foreach` / `on_error` 串成小流程
5. `Lesson 24-27` 再把 Excel 导入链路打通
6. `Lesson 06-07` 最后接 Redis / Postgres

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
- 能独立使用 `save_as`、`set_var`、`assert_*`、`read_csv`、`write_csv`、`read_excel`
- 能解释为什么某一步放在 `Lua`，某一步放在 `Flow`
- 能至少接通一个外部系统
- 能说明一个流程失败后应该看哪里

## 学完之后去哪里

下一站是：
[track-intermediate.md](track-intermediate.md)
