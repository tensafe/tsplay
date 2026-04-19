# TSPlay Docs

`ReadMe.md` 负责项目入口和快速上手，`docs/` 负责承载更适合团队长期使用的材料，尤其是培训、Enablement 和交付规范。

## 推荐阅读顺序

1. [项目总览](../ReadMe.md)
2. [Step-by-Step 教程](tutorials/README.md)
3. [完整课程总览](tutorials/curriculum-overview.md)
4. [160 次递进迭代路线图](tutorials/iteration-roadmap-160.md)
5. [培训体系总览](training/README.md)
6. [AI 无感入门](training/ai-intent-to-flow.md)
7. [学习路径](training/learning-path.md)
8. [训练营课程表](training/bootcamp-plan.md)
9. [实训实验](training/labs.md)
10. [考核与认证](training/assessment.md)

## 文档地图

| 类别 | 说明 | 入口 |
| --- | --- | --- |
| 项目入口 | TSPlay 的核心概念、运行方式、Flow 和 MCP 能力 | [../ReadMe.md](../ReadMe.md) |
| Step-by-Step 教程 | 面向使用者的分步上手教程；同一个功能同时给出 Lua 和 Flow 写法 | [tutorials/README.md](tutorials/README.md) |
| 完整进阶教程 | 按新手 / 初级 / 中级 / 高级组织的一整套课程体系 | [tutorials/curriculum-overview.md](tutorials/curriculum-overview.md) |
| 160 次迭代路线图 | 把教程建设拆成 160 个渐进迭代点，适合持续演进 | [tutorials/iteration-roadmap-160.md](tutorials/iteration-roadmap-160.md) |
| 培训总览 | 培训对象、交付模式、成功指标和文档清单 | [training/README.md](training/README.md) |
| AI 新手教程 | 面向 Codex、OpenClaw 等 Agent 的“用户意图 -> MCP -> Flow -> 执行修复”实战教程 | [training/ai-intent-to-flow.md](training/ai-intent-to-flow.md) |
| 学习路径 | 从新人到讲师的分层路线图 | [training/learning-path.md](training/learning-path.md) |
| 课程安排 | 2 天 Bootcamp 和 4 周应用节奏 | [training/bootcamp-plan.md](training/bootcamp-plan.md) |
| 实操实验 | 结合本仓库 `demo/` 和 `script/` 的实验清单 | [training/labs.md](training/labs.md) |
| Capstone 场景 | 结业项目说明和交付要求 | [training/capstone-briefs.md](training/capstone-briefs.md) |
| 考核与认证 | 评分、门槛、证据和复盘机制 | [training/assessment.md](training/assessment.md) |
| 讲师手册 | 讲师备课、授课、辅导和版本维护 | [training/trainer-playbook.md](training/trainer-playbook.md) |

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

- 对第一次上手、想按功能对照学习的同学：先看“Step-by-Step 教程”
- 对要系统学习、要把教程一直迭代下去的同学：直接从“完整课程总览”进入
- 对想直接用 AI 做事的新手：先看“AI 无感入门”，再做 MCP 相关实验
- 对想把 TSPlay 接到大模型产品里的同学：先看“AI 无感入门”，重点关注接入方式、system prompt、授权策略和失败闭环
- 对个人学习者：按“总览 -> 学习路径 -> Labs”走最快
- 对项目经理或 Enablement 负责人：先看“培训体系总览”和“课程安排”
- 对讲师：先看“讲师手册”，再按 cohort 目标挑 Labs 和 Capstone
