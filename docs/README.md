# TSPlay Docs

根目录 [../ReadMe.md](../ReadMe.md) 提供英文项目入口，[../README.zh-CN.md](../README.zh-CN.md) 提供中文对应版本；教程入口层现在也采用相同策略：英文默认文件名，中文使用 `.zh-CN.md` 副本。`docs/` 负责承载更适合团队长期使用的材料，尤其是培训、Enablement 和交付规范。

## 常用入口

<div class="grid cards" markdown>

-   :material-rocket-launch-outline:{ .lg .middle } __第一次看 TSPlay__

    想先知道 TSPlay 是什么、三层能力怎么分、怎么 5 分钟跑起来。

    [看中文项目总览](../README.zh-CN.md)

-   :material-play-box-outline:{ .lg .middle } __今天先跑起来__

    想按最短路径把教程跑通，再慢慢理解 Flow、会话和 MCP。

    [进入教程总站](tutorials/README.zh-CN.md)

-   :material-file-document-edit-outline:{ .lg .middle } __主要想写 Flow__

    想看更贴近交付的路线、变量、控制流、认证导入和恢复。

    [看学习路径](training/learning-path.md)

-   :material-robot-outline:{ .lg .middle } __主要想接 Agent__

    想从“用户意图 -> MCP -> Flow -> 执行修复”开始。

    [看 AI 无感入门](training/ai-intent-to-flow.md)

-   :material-school-outline:{ .lg .middle } __我要做培训__

    想组织课程、训练营、实训实验和讲师材料。

    [看培训体系总览](training/README.md)

-   :material-map-search-outline:{ .lg .middle } __我要看完整地图__

    想从文档全图、课程体系和路线图来理解仓库资料。

    [看完整课程总览](tutorials/curriculum-overview.zh-CN.md)

</div>

## 高频跳转

| 场景 | 直接入口 |
| --- | --- |
| 跑第一条 Flow | [Lesson 01](tutorials/01-hello-world.md) |
| 本地抓表到 JSON | [Lesson 03](tutorials/03-capture-table.md) |
| 页面断言与检查 | [Lesson 10](tutorials/10-assert-page-state.md) |
| 批量导入与恢复 | [Lesson 22](tutorials/22-foreach-batch-import-csv.md) |
| 浏览器状态与命名会话 | [Lesson 36](tutorials/36-save-storage-state.md)、[Lesson 42](tutorials/42-use-named-session.md) |
| MCP 能力总入口 | [Lesson 111](tutorials/111-mcp-list-actions.md)、[Lesson 120](tutorials/120-mcp-finalize-flow.md) |
| 培训与课程安排 | [training/README.md](training/README.md)、[training/bootcamp-plan.md](training/bootcamp-plan.md) |

## 推荐阅读顺序

1. [项目总览（中文）](../README.zh-CN.md)
2. [产品定位与工作台方案](product/README.md)
3. [Step-by-Step 教程（中文）](tutorials/README.zh-CN.md)
4. [完整课程总览（中文）](tutorials/curriculum-overview.zh-CN.md)
5. [160 次递进迭代路线图](tutorials/iteration-roadmap-160.md)
6. [培训体系总览](training/README.md)
7. [AI 无感入门](training/ai-intent-to-flow.md)
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
| Step-by-Step 教程 | 面向使用者的分步上手教程；同一个功能同时给出 Lua 和 Flow 写法 | [tutorials/README.zh-CN.md](tutorials/README.zh-CN.md) |
| 完整进阶教程 | 按新手 / 初级 / 中级 / 高级组织的一整套课程体系 | [tutorials/curriculum-overview.zh-CN.md](tutorials/curriculum-overview.zh-CN.md) |
| 160 次迭代路线图 | 把教程建设拆成 160 个渐进迭代点，适合持续演进 | [tutorials/iteration-roadmap-160.md](tutorials/iteration-roadmap-160.md) |
| 培训总览 | 培训对象、交付模式、成功指标和文档清单 | [training/README.md](training/README.md) |
| AI 新手教程 | 面向 Codex、OpenClaw 等 Agent 的“用户意图 -> MCP -> Flow -> 执行修复”实战教程 | [training/ai-intent-to-flow.md](training/ai-intent-to-flow.md) |
| 学习路径 | 从新人到讲师的分层路线图 | [training/learning-path.md](training/learning-path.md) |
| 课程安排 | 2 天 Bootcamp 和 4 周应用节奏 | [training/bootcamp-plan.md](training/bootcamp-plan.md) |
| 实操实验 | 结合本仓库 `demo/` 和 `script/` 的实验清单 | [training/labs.md](training/labs.md) |
| Capstone 场景 | 结业项目说明和交付要求 | [training/capstone-briefs.md](training/capstone-briefs.md) |
| 考核与认证 | 评分、门槛、证据和复盘机制 | [training/assessment.md](training/assessment.md) |
| 教程自动录屏 | 用 `tsplay` + `ffmpeg` 把教程演示稳定录成视频素材 | [training/tutorial-video-recording.md](training/tutorial-video-recording.md) |
| 讲师手册 | 讲师备课、授课、辅导和版本维护 | [training/trainer-playbook.md](training/trainer-playbook.md) |

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
- 对要系统学习、要把教程一直迭代下去的同学：直接从“完整课程总览”进入
- 对想直接用 AI 做事的新手：先看“AI 无感入门”，再做 MCP 相关实验
- 对想把 TSPlay 接到大模型产品里的同学：先看“AI 无感入门”，重点关注接入方式、system prompt、授权策略和失败闭环
- 对个人学习者：按“总览 -> 学习路径 -> Labs”走最快
- 对项目经理或 Enablement 负责人：先看“培训体系总览”和“课程安排”
- 对讲师：先看“讲师手册”，再按 cohort 目标挑 Labs 和 Capstone
