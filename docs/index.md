# TSPlay 文档站

面向 AI Agent 和交付团队的浏览器自动化文档入口。

这个站点把 `项目总览 / 教程 / 培训 / MCP / Flow 交付` 收到一套可直接浏览的网页里，适合给新同学、实施同学、讲师和集成人员统一分发。

## 从哪里开始

### 第一次接触 TSPlay

- 看 [项目总览（中文）](../README.zh-CN.md)
- 如果你更习惯英文，看 [Project Overview](../ReadMe.md)
- 如果你只想先跑起来，优先看 [Step-by-Step 教程](tutorials/README.zh-CN.md)

### 想直接学写 Flow

- 从 [教程索引（中文）](tutorials/README.zh-CN.md) 进入
- 想走 Agent 路线，直接看 [AI 无感入门](training/ai-intent-to-flow.md)
- 想系统提升，继续看 [学习路径](training/learning-path.md)

### 想做培训或团队推广

- 看 [培训体系总览](training/README.md)
- 配合 [训练营课程表](training/bootcamp-plan.md)
- 讲师再看 [讲师手册](training/trainer-playbook.md)

## 5 分钟跑起来

```bash
go mod download
go run . -flow script/demo_baidu.flow.yaml
```

如果你更喜欢先构建二进制：

```bash
go build -o tsplay .
./tsplay -flow script/tutorials/01_hello_world.flow.yaml
```

## 三条推荐路径

| 目标 | 先看什么 | 下一步 |
| --- | --- | --- |
| 今天先跑通一条自动化 | [项目总览（中文）](../README.zh-CN.md) | [教程 Lesson 01](tutorials/01-hello-world.md) |
| 让 Codex / Agent 帮你出 Flow | [AI 无感入门](training/ai-intent-to-flow.md) | MCP 相关 `101-120` 课 |
| 给团队做培训 | [培训体系总览](training/README.md) | [Bootcamp 课程表](training/bootcamp-plan.md) |

## 这个站点包含什么

- [项目总览（中文）](../README.zh-CN.md)：产品定位、三层能力、快速开始、MCP、安全边界
- [Project Overview](../ReadMe.md)：英文版本总览
- [文档入口](README.md)：站内文档地图和推荐阅读顺序
- [教程索引（中文）](tutorials/README.zh-CN.md)：160 节递进式教程
- [English Tutorial Map](tutorials/README.md)：英文课程地图
- [培训体系](training/README.md)：学习路径、课程表、实验、考核、讲师材料
- [产品方案](product/README.md)：产品定位和工作台思路

## 仓库与站点

- GitHub 仓库：[tensafe/tsplay](https://github.com/tensafe/tsplay)
- 站点发布说明：[GitHub Pages 发布说明](github-pages.md)
