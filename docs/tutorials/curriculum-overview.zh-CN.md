# TSPlay 金字塔课程总图

[English](curriculum-overview.md) | 简体中文

这页不再把教程当成一长串 lesson 目录来展示。  
它只回答一个问题：

TSPlay 的教程体系，为什么要按金字塔来组织，而不是按功能把内容平铺开来。

## 先给结论

TSPlay 的教程更适合分成四层：

1. 先跑通基本动作
2. 再学会组织小流程
3. 再接真实交付闭环
4. 最后进入 Agent、边界和标准化

如果顺序反过来，新手会一下子同时撞上：

- 环境问题
- 页面问题
- 外部系统问题
- 安全边界问题
- MCP / repair / review 这些高阶抽象

这就是“教程混杂感”的主要来源。

## 金字塔四层总表

| 层级 | Lesson 范围 | 这一层只回答什么问题 | 核心产物 | 入口 |
| --- | --- | --- | --- | --- |
| 第 1 层：跑通基本动作 | `01-12` | 我能不能稳定跑起来，并看懂 `Lua / Flow / artifacts` 的关系 | 本地页面交互结果、JSON 输出、基础断言 | [track-newbie.zh-CN.md](track-newbie.zh-CN.md) |
| 第 2 层：结构化小流程 | `13-57` | 我能不能把文件、变量、控制流、会话串成一个小业务流程 | CSV / Excel / 控制流 / 会话复用的小闭环 | [track-junior.zh-CN.md](track-junior.zh-CN.md) |
| 第 3 层：真实交付闭环 | `58-100` | 我能不能把浏览器结果接进 Redis、数据库、对账、交接和模板 | 外部系统同步、审计、交接包、模板索引 | [track-intermediate.zh-CN.md](track-intermediate.zh-CN.md) |
| 第 4 层：Agent 与标准化 | `101-160` | 我能不能把 TSPlay 提升成团队能力，而不是个人脚本 | MCP 闭环、授权边界、review 规范、单二进制交付、课程演进 | [track-advanced.zh-CN.md](track-advanced.zh-CN.md) |

## 为什么必须按金字塔组织

### 1. 先确定性，后复杂性

先用本地页面、本地文件和固定 demo 建手感。  
只有确定性足够高，后面的 Redis、数据库、MCP、repair 才不会把学习成本放大。

### 2. 先动作，后编排

先会一个动作，再会几个动作怎么串。  
这比一开始就把 `retry / foreach / on_error / session / batch` 全部塞进来更稳。

### 3. 先交付物，后抽象词

每一层都要求留下证据：

- 一条脚本或 Flow
- 一份写进 `artifacts/` 的输出
- 一份可以复盘的结果

没有产物，抽象概念很快就会失真。

### 4. 上层不能抢跑

如果还没完成第 1 层，就不应该把主要精力放在第 4 层。  
如果第 2 层还不稳，也不应该直接把 Redis、Postgres、交接模板和 MCP 当主线。

## 四层展开

### 第 1 层：先跑通

适合谁：

- 第一天接触 TSPlay
- 还不熟悉 `Lua / Flow / MCP`
- 需要先建立最小成功体验的人

这一层的重点：

- 跑通第一条 Flow
- 看懂本地 demo 页面
- 学会页面提取、文本断言、JSON 输出
- 知道为什么结果要写进 `artifacts/`

建议先看：

- [Lesson 01](01-hello-world.md)
- [Lesson 03](03-capture-table.md)
- [Lesson 10](10-assert-page-state.md)
- [新手路线](track-newbie.zh-CN.md)

### 第 2 层：会组织小流程

适合谁：

- 已经跑通本地 demo
- 想开始处理文件、变量、控制流和会话
- 想从“会跑示例”进入“会做一条小流程”

这一层的重点：

- `read_csv / write_csv / read_excel`
- `retry / if / foreach / on_error / wait_until`
- `storage_state / named session`
- 认证页面导入、导出、回读闭环

建议先看：

- [Lesson 13](13-read-csv-basics.md)
- [Lesson 16](16-retry-flaky-action.md)
- [Lesson 22](22-foreach-batch-import-csv.md)
- [Lesson 36](36-save-storage-state.md)
- [Lesson 57](57-use-session-import-export-round-trip.md)
- [初级路线](track-junior.zh-CN.md)

### 第 3 层：能做真实交付

适合谁：

- 已能维护一条基础 Flow
- 开始接真实导入导出和外部系统
- 需要沉淀对账、审计、交接和模板的人

这一层的重点：

- Redis 摘要和批次 key
- Postgres 摘要、明细、事务和 upsert
- 三边对账、异常台账、审计留痕
- handoff、manifest、模板目录、模板索引

建议先看：

- [Lesson 58](58-sync-import-report-summary-to-redis.md)
- [Lesson 61](61-db-insert-import-batch-summary.md)
- [Lesson 71](71-external-system-round-trip.md)
- [Lesson 87](87-build-handoff-artifact-manifest.md)
- [Lesson 96](96-build-template-index.md)
- [中级路线](track-intermediate.zh-CN.md)

### 第 4 层：进入 Agent 与标准化

适合谁：

- 要把 TSPlay 带进团队、产品或客户环境
- 要做 review、规范、发布、培训和演进的人
- 要把它变成系统能力，而不是一组个人脚本的人

这一层的重点：

- `observe -> draft -> validate -> run -> repair -> finalize`
- `allow_*`、`security_preset`、边界对照
- review、命名、artifact 布局、大型 Flow 包
- 单二进制交付、离线学习、capstone、教程持续演进

建议先看：

- [Lesson 111](111-mcp-list-actions.md)
- [Lesson 120](120-mcp-finalize-flow.md)
- [Lesson 127](127-compare-local-flow-and-mcp-boundaries.md)
- [Lesson 134](134-review-example-with-checklist.md)
- [Lesson 144](144-single-binary-delivery-flow.md)
- [Lesson 160](160-curriculum-continuation-plan.md)
- [高级路线](track-advanced.zh-CN.md)

## 角色怎么映射到这座金字塔

| 角色 | 推荐主路径 | 为什么 |
| --- | --- | --- |
| 测试 / 运营 / 实施 | 第 1 层 -> 第 2 层 | 先把运行、提取、断言、批量和会话练稳 |
| 自动化开发 / Flow 编写者 | 第 1 层关键课 -> 第 2 层 -> 第 3 层 | 真正的工作量主要在结构化和交付闭环 |
| AI / 平台工程师 | 第 2 层 -> 第 3 层 -> 第 4 层 | 先理解 Flow 和交付，再进入 MCP 和边界 |
| 讲师 / Enablement | 先读整张图，再按 1-4 层设计课程 | 需要保证学习顺序不跳跃、评估标准可复用 |

## 课程作者怎么使用这张图

1. 先决定一节课属于哪一层，不要先想它属于哪个 action。
2. 再决定这节课在这一层里回答哪个问题，不要一节课同时解决三个层级的问题。
3. 每一层都要先给默认主线，再给补充课，不要把补充课插进主线里抢前置位置。
4. 每加一节新课，都要回答“它是在加宽哪一层，还是在抬高哪一层”。

## 如果你现在要继续往下看

- 想实际开始学：看 [教程首页](README.zh-CN.md)
- 想按层走：看 [新手路线](track-newbie.zh-CN.md)、[初级路线](track-junior.zh-CN.md)、[中级路线](track-intermediate.zh-CN.md)、[高级路线](track-advanced.zh-CN.md)
- 想按等级看能力要求：看 [学习路径](../training/learning-path.md)
- 想看更细的 lesson 编号扩展：回到 [教程首页](README.zh-CN.md) 里的完整教程地图
- 想长期维护这套课程：看 [160 次迭代路线图](iteration-roadmap-160.md) 和 [演进手册](evolution-playbook.md)
