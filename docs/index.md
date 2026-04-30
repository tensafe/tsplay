# TSPlay 文档站

面向 AI Agent 和交付团队的浏览器自动化文档入口。

这个站点把 `项目总览 / 教程 / 培训 / MCP / Flow 交付` 整理在一起，方便新同学、实施同学、讲师和集成人员沿同一条路径进入。

## 从哪里开始

### 第一次接触 TSPlay

- 看 [项目总览（中文）](../README.zh-CN.md)
- 如果你更习惯英文，看 [Project Overview](../ReadMe.md)
- 如果你只想先跑起来，优先看 [Step-by-Step 教程](tutorials/README.zh-CN.md)

### 想先进入 Flow 主线

- 从 [教程索引（中文）](tutorials/README.zh-CN.md) 进入
- 想走 Agent 路线，直接看 [AI 协作入门](training/ai-intent-to-flow.md)
- 想系统提升，继续看 [学习路径](training/learning-path.md)

### 想先查命令入口

- 看 [CLI `-action` 参考](actions/README.md)
- 适合先弄清楚命令行 `-action` 现在支持什么
- 如果你更关心单二进制入口，也可以顺手看 `list-assets / extract-assets / file-srv`

### 想先查支持行为

- 看 [支持行为清单](capability-actions/README.md)
- 适合先弄清楚 `navigate / click / read_csv / db_query / retry` 这层能力在 `Flow / Lua / MCP` 里怎么对应
- 如果你正在写教程、做 Agent 接入或补 action 文档，这一层通常更值得先看

### 想先看 Skills

- 看 [Skills 介绍](skills/README.md)
- 适合先在 Codex 中生成、修改或修复 Flow
- 如果你希望 Codex 按统一方式协作，这一层值得先看

### 想准备培训或团队推广

- 看 [培训体系总览](training/README.md)
- 配合 [训练营课程表](training/bootcamp-plan.md)
- 讲师再看 [讲师手册](training/trainer-playbook.md)

## 5 分钟跑起来

最短路径已经单独收进了 [快速开始](../getting-started.md)：

- 直接 `curl` 下载最新 release 二进制
- 运行 `./tsplay -action quickstart-demo`
- 自动生成并执行最小 demo Flow
- 不需要先装 Go，也不需要先下载 Playwright

如果你更喜欢从源码仓库开始，再看 [项目总览（中文）](../README.zh-CN.md) 里的完整命令表。

## 三条推荐路径

| 目标 | 先看什么 | 下一步 |
| --- | --- | --- |
| 完成首次自动化运行 | [项目总览（中文）](../README.zh-CN.md) | [教程 Lesson 01](tutorials/01-hello-world.md) |
| 让 Codex / Agent 协助生成 Flow | [AI 协作入门](training/ai-intent-to-flow.md) | MCP 相关 `101-120` 课 |
| 给团队准备培训路径 | [培训体系总览](training/README.md) | [Bootcamp 课程表](training/bootcamp-plan.md) |

## 这个站点包含什么

- [项目总览（中文）](../README.zh-CN.md)：产品定位、三层能力、快速开始、MCP、安全边界
- [Project Overview](../ReadMe.md)：英文版本总览
- [文档入口](README.md)：站内文档地图和推荐阅读顺序
- [Skills 介绍](skills/README.md)：`skill` 在 Codex / Agent 协作里的角色说明
- [支持行为清单](capability-actions/README.md)：`navigate / click / read_csv / db_query / retry` 这层动作词典
- [CLI `-action` 参考](actions/README.md)：命令行 `-action` 支持列表与逐项说明
- [教程索引（中文）](tutorials/README.zh-CN.md)：160 节递进式教程
- [金字塔课程总图](tutorials/curriculum-overview.zh-CN.md)：先跑通、再结构化、再交付、再标准化的总图
- [English Tutorial Map](tutorials/README.md)：英文课程地图
- [培训体系](training/README.md)：学习路径、课程表、实验、考核、讲师材料
- [产品方案](product/README.md)：产品定位和工作台思路

## 仓库与站点

- GitHub 仓库：[tensafe/tsplay](https://github.com/tensafe/tsplay)
