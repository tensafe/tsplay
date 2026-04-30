# TSPlay Docs

根目录 [../ReadMe.md](../ReadMe.md) 提供英文项目入口，[../README.zh-CN.md](../README.zh-CN.md) 提供中文对应版本；教程入口层现在也采用相同策略：英文默认文件名，中文使用 `.zh-CN.md` 副本。`docs/` 负责承载更适合团队长期使用的材料，尤其是培训、Enablement 和交付规范。

教程默认也不再把 `160` 节课直接摊给你自己筛，而是先按四层金字塔进入：先跑通动作，再组织成 Flow，再接进真实交付，最后再进入 Agent、MCP 和标准化。

## 常用入口

<div class="grid cards" markdown>

-   :material-rocket-launch-outline:{ .lg .middle } __项目概览__

    了解 TSPlay 的定位、三层能力和快速开始。

    [看中文项目总览](../README.zh-CN.md)

-   :material-play-box-outline:{ .lg .middle } __快速开始__

    先拿到一次成功运行，再按四层金字塔决定今天进入哪一层。

    [进入教程总站](tutorials/README.zh-CN.md)

-   :material-file-document-edit-outline:{ .lg .middle } __Flow 路线__

    想看更贴近交付的路线、变量、控制流、认证导入和恢复。

    [看学习路径](training/learning-path.md)

-   :material-robot-outline:{ .lg .middle } __Agent 集成__

    从“用户意图 -> MCP -> Flow -> 执行修复”开始。

    [看 AI 协作入门](training/ai-intent-to-flow.md)

-   :material-shape-outline:{ .lg .middle } __支持行为__

    想先查 `navigate / click / read_csv / db_query / retry` 这些动作在 `Flow / Lua / MCP` 里怎么对应。

    [看支持行为清单](capability-actions/README.md)

-   :material-creation-outline:{ .lg .middle } __Skills 说明__

    了解如何在 Codex 中通过 Skills 生成、修改和修复 Flow，以及当前仓库已提供的协作说明。

    [看 Skills 介绍](skills/README.md)

-   :material-console-line:{ .lg .middle } __CLI 命令入口__

    想先看命令行 `-action` 现在支持什么、每个命令该什么时候用。

    [看 CLI `-action` 参考](actions/README.md)

-   :material-school-outline:{ .lg .middle } __培训材料__

    想组织课程、训练营、实训实验和讲师材料。

    [看培训体系总览](training/README.md)

-   :material-map-search-outline:{ .lg .middle } __文档总图__

    想从文档全图、课程体系和路线图来理解仓库资料。

    [看金字塔课程总图](tutorials/curriculum-overview.zh-CN.md)

</div>

## 如果你现在只想快速进入

先不用把整套文档一口气看完。按你现在最关心的目标，先选一个入口就够：

| 你现在最想做什么 | 入口 |
| --- | --- |
| 先把 TSPlay 跑通 | [tutorials/README.zh-CN.md](tutorials/README.zh-CN.md) |
| 先知道项目是做什么的 | [../README.zh-CN.md](../README.zh-CN.md) |
| 先学怎么写 Flow | [tutorials/track-newbie.zh-CN.md](tutorials/track-newbie.zh-CN.md) |
| 先接 Agent / MCP | [training/ai-intent-to-flow.md](training/ai-intent-to-flow.md) |
| 先准备做培训 | [training/README.md](training/README.md) |

## 高频跳转

这一栏按“最常被点开的入口”组织，不代表推荐阅读顺序。
如果你想按顺序学习教程，不要先去翻 `160` 节课表；先从 [tutorials/README.zh-CN.md](tutorials/README.zh-CN.md) 进入，按四层金字塔选层，再按 lesson 往下走。

### 基础起步

- 第一条可运行示例：[Lesson 01](tutorials/01-hello-world.md)
- 本地抓表到 JSON：[Lesson 03](tutorials/03-capture-table.md)
- 页面断言与检查：[Lesson 10](tutorials/10-assert-page-state.md)
- 批量导入与恢复：[Lesson 22](tutorials/22-foreach-batch-import-csv.md)

### 会话与认证

- 保存浏览器状态：[Lesson 36](tutorials/36-save-storage-state.md)
- 复用命名会话：[Lesson 42](tutorials/42-use-named-session.md)
- 登录后完成导入：[Lesson 44](tutorials/44-session-import-with-login.md)
- 认证导入导出闭环：[Lesson 57](tutorials/57-use-session-import-export-round-trip.md)

### MCP 与 Agent

- MCP 能力总览：[Lesson 111](tutorials/111-mcp-list-actions.md)
- 页面观察入口：[Lesson 113](tutorials/113-mcp-observe-page.md)
- 收敛生成 Flow：[Lesson 120](tutorials/120-mcp-finalize-flow.md)
- 培训与课程安排：[training/README.md](training/README.md)、[training/bootcamp-plan.md](training/bootcamp-plan.md)

## 推荐阅读顺序

1. [项目总览（中文）](../README.zh-CN.md)
2. [产品定位与工作台方案](product/README.md)
3. [Step-by-Step 教程（中文）](tutorials/README.zh-CN.md)
4. [金字塔课程总图（中文）](tutorials/curriculum-overview.zh-CN.md)
5. [160 次递进迭代路线图](tutorials/iteration-roadmap-160.md)
6. [培训体系总览](training/README.md)
7. [AI 协作入门](training/ai-intent-to-flow.md)
8. [学习路径](training/learning-path.md)
9. [训练营课程表](training/bootcamp-plan.md)
10. [实训实验](training/labs.md)
11. [考核与认证](training/assessment.md)
12. [教程自动录屏](training/tutorial-video-recording.md)

## 文档地图

| 类别 | 说明 | 入口 |
| --- | --- | --- |
| 项目入口 | TSPlay 的核心概念、运行方式、Flow 和 MCP 能力 | [../README.zh-CN.md](../README.zh-CN.md) |
| 产品方案 | 已授权 Web 系统认知与数据编排工作台的定位、MVP、分层和职责边界 | [product/README.md](product/README.md) |
| Skills 介绍 | 解释 `skill` 在 Codex / Agent 协作里解决什么问题，以及当前仓库附带的 skill | [skills/README.md](skills/README.md) |
| 支持行为清单 | 汇总 `navigate / click / read_csv / db_query / retry` 这类动作在 Flow、Lua、MCP 三边的对应关系，并提供按类别查询入口 | [capability-actions/README.md](capability-actions/README.md) |
| CLI `-action` 参考 | 汇总命令行 `-action` 支持列表，并给每个命令入口提供单独说明页 | [actions/README.md](actions/README.md) |
| 核心功能路线图 | 把 Top 10 核心能力拆成优先级、验收标准、代码落点和建议里程碑 | [product/core-feature-roadmap.md](product/core-feature-roadmap.md) |
| 核心功能执行面板 | 把 Top 10 核心能力继续拆成第一刀、依赖、测试和文档联动 | [product/core-feature-execution-board.md](product/core-feature-execution-board.md) |
| 30 轮持续进化计划 | 把持续分析、梳理、优化、完善收成 30 轮可执行计划 | [product/30-iteration-evolution-plan.md](product/30-iteration-evolution-plan.md) |
| Step-by-Step 教程 | 面向使用者的分步上手教程；同一个功能同时给出 Lua 和 Flow 写法 | [tutorials/README.zh-CN.md](tutorials/README.zh-CN.md) |
| 金字塔课程总图 | 用“先跑通 -> 再结构化 -> 再交付 -> 再标准化”的四层方式组织整套教程 | [tutorials/curriculum-overview.zh-CN.md](tutorials/curriculum-overview.zh-CN.md) |
| 160 次迭代路线图 | 把教程建设拆成 160 个渐进迭代点，适合持续演进 | [tutorials/iteration-roadmap-160.md](tutorials/iteration-roadmap-160.md) |
| 培训总览 | 培训对象、交付模式、成功指标和文档清单 | [training/README.md](training/README.md) |
| AI 协作入门 | 面向 Codex、OpenClaw 等 Agent 的“用户意图 -> MCP -> Flow -> 执行修复”实战教程 | [training/ai-intent-to-flow.md](training/ai-intent-to-flow.md) |
| 学习路径 | 从新人到讲师的分层路线图 | [training/learning-path.md](training/learning-path.md) |
| 课程安排 | 2 天 Bootcamp 和 4 周应用节奏 | [training/bootcamp-plan.md](training/bootcamp-plan.md) |
| 实操实验 | 结合本仓库 `demo/` 和 `script/` 的实验清单 | [training/labs.md](training/labs.md) |
| Capstone 场景 | 结业项目说明和交付要求 | [training/capstone-briefs.md](training/capstone-briefs.md) |
| 考核与认证 | 评分、门槛、证据和复盘机制 | [training/assessment.md](training/assessment.md) |
| 教程自动录屏 | 用 `tsplay` + `ffmpeg` 把教程演示稳定录成视频素材 | [training/tutorial-video-recording.md](training/tutorial-video-recording.md) |
| 讲师手册 | 讲师备课、授课、辅导和版本维护 | [training/trainer-playbook.md](training/trainer-playbook.md) |
| 文档健康检查 | 汇总链接健康、内容断层和固定检查动作 | [doc-health-audit.md](doc-health-audit.md) |

## 按角色进入

- 实施 / 测试 / 运营：从 [tutorials/README.zh-CN.md](tutorials/README.zh-CN.md) 和 [training/labs.md](training/labs.md) 开始
- 自动化开发 / Flow 编写者：先看 [training/learning-path.md](training/learning-path.md)，再挑 [tutorials/track-junior.zh-CN.md](tutorials/track-junior.zh-CN.md)
- AI / 平台工程师：先看 [training/ai-intent-to-flow.md](training/ai-intent-to-flow.md)，再进 [tutorials/119-mcp-chain-overview.md](tutorials/119-mcp-chain-overview.md)
- 讲师 / Enablement：先看 [training/README.md](training/README.md)、[training/bootcamp-plan.md](training/bootcamp-plan.md)、[training/trainer-playbook.md](training/trainer-playbook.md)

## 仓库内可直接复用的训练素材

- 示例脚本与 Flow：
  [../script/open_url.lua](../script/open_url.lua),
  [../script/demo_baidu.flow.yaml](../script/demo_baidu.flow.yaml),
  [../script/is_sel.lua](../script/is_sel.lua)
- 本地演示页面：
  [../demo/demo.html](../demo/demo.html),
  [../demo/tables.html](../demo/tables.html),
  [../demo/upload.html](../demo/upload.html),
  [../demo/multi_upfile.html](../demo/multi_upfile.html)
- 失败现场与观察产物目录：
  `artifacts/`

## 使用建议

- 对正在定义产品形态、讨论路线和边界的同学：先看“产品定位与工作台方案”
- 对第一次上手、想按功能对照学习的同学：先看“Step-by-Step 教程”
- 对要系统学习、要把教程一直迭代下去的同学：直接从“金字塔课程总图”进入
- 对希望直接使用 AI 协作的人：先看“AI 协作入门”，再做 MCP 相关实验
- 对想把 TSPlay 接到大模型产品里的同学：先看“AI 协作入门”，重点关注接入方式、system prompt、授权策略和失败闭环
- 对个人学习者：按“总览 -> 学习路径 -> Labs”走最快
- 对项目经理或 Enablement 负责人：先看“培训体系总览”和“课程安排”
- 对讲师：先看“讲师手册”，再按 cohort 目标挑 Labs 和 Capstone
