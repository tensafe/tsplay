---
hide:
  - toc
search:
  boost: 2
---

<div class="tsplay-hero">
  <h1>面向交付与 AI 协作的<br>浏览器自动化文档站</h1>
  <p>
    本文档汇总 <code>Flow DSL</code>、<code>MCP</code>、<code>Skills</code>、会话复用、失败修复与教程交付，
    供研发、实施、测试、培训和 Agent 集成使用。
  </p>
  <div class="tsplay-hero__actions">
    <a class="md-button md-button--primary" href="getting-started/">快速开始</a>
  </div>
  <div class="tsplay-hero__minor">
    Agent 集成：<a href="docs/training/ai-intent-to-flow/">查看 AI 协作入门</a>
    <span> · </span>
    Codex 中使用 Skills 生成或修改 Flow：<a href="docs/skills/">查看 Skills 介绍</a>
  </div>
</div>

## 在 Codex 中使用 Skills

如果你的目标是让 Codex 直接帮你生成、修改或修复 TSPlay Flow，优先从 `Skills` 入口开始。
仓库附带的 `tsplay-flow-authoring` 支持这类工作：

- 根据自然语言需求生成新的 `.flow.yaml`
- 修改已有 Flow 的 selector、等待、断言和变量链路
- 按 `finalize -> run -> repair` 的方式继续收敛

在支持 skills 的 Codex 环境里，可以直接这样提：

```text
请使用 tsplay-flow-authoring，帮我生成一条 TSPlay Flow。
- 页面: <URL 或本地页面>
- 目标: <要完成的动作>
- 输入: <关键词 / 文件 / 条件>
- 输出: <JSON / CSV / save_as / artifact>
- 授权: <readonly / browser_write / full_automation / allow_*>
```

```text
请使用 tsplay-flow-authoring，修改这条 Flow。
- 文件: <flow 文件路径>
- 问题: <超时 / selector 失效 / assert 失败 / 输出为空>
- 预期: <修完后得到什么结果>
```

[查看 Skills 介绍与 Codex 用法](docs/skills/README.md)

## 首次使用

第一次接触 TSPlay，先不要纠结 `Lua / Flow / MCP`。
先执行下面这条路径，完成一次运行后再继续：

```bash
curl -L -o tsplay https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-arm64
chmod +x tsplay
./tsplay -action quickstart-demo
```

这条会直接下载二进制、生成最小 demo Flow 并立刻执行，而且不需要先下载 Playwright。
[查看完整快速开始](getting-started.md)

## 按目标进入

<div class="grid cards" markdown>

-   :material-file-code-outline:{ .lg .middle } __Flow 编写__
    
    ---

    将页面动作整理成可审阅、可复用的 `.flow.yaml`。

    [进入 Flow 学习路径](docs/tutorials/track-newbie.zh-CN.md)

-   :material-robot-outline:{ .lg .middle } __Agent 集成__

    ---

    将 `observe -> draft -> validate -> run -> repair` 接入 MCP 或 Agent 产品。

    [进入 AI 协作入门](docs/training/ai-intent-to-flow.md)

-   :material-shape-plus-outline:{ .lg .middle } __Skills 说明__

    ---

    在 Codex 中生成、修改和修复 Flow，并补充统一的说明、规则和参考材料。

    [进入 Skills 介绍](docs/skills/README.md)

-   :material-school-outline:{ .lg .middle } __培训与推广__

    ---

    用于团队培训、onboarding、Bootcamp 或交付推广。

    [进入培训体系](docs/training/README.md)

</div>

## Skills 是什么

<div class="grid cards" markdown>

-   :material-lightbulb-on-outline:{ .lg .middle } __不是 action 清单__

    `skill` 不是单个能力参数表，而是一套让 Codex 更稳定协作的提示、规则和参考资料。

-   :material-file-tree-outline:{ .lg .middle } __不是替代 Flow__

    `Flow` 仍然是最终交付物；`skill` 解决的是“如何更稳定地得到这条 Flow”。

-   :material-rocket-launch-outline:{ .lg .middle } __仓库附带的 Skill__

    仓库附带 `tsplay-flow-authoring`，可直接用于 Codex 中的 Flow 生成、修改、审阅和邮件通知场景。

    [查看 Skills 介绍](docs/skills/README.md)

</div>

## 主要能力

<div class="grid cards" markdown>

-   :material-layers-triple-outline:{ .lg .middle } __三层入口统一__

    `Lua CLI / Script`、`Flow DSL`、`MCP Server` 同仓库同运行时，适合探索、固化和集成三条路径并行。

-   :material-shield-check-outline:{ .lg .middle } __安全边界明确__

    用 `security_preset` 和 `allow_*` 授权高风险能力，适合把浏览器能力暴露给 Agent。

-   :material-image-filter-hdr-outline:{ .lg .middle } __失败现场可复盘__

    失败时自动留下截图、HTML、DOM snapshot 和 trace，让 repair 与 review 有抓手。

-   :material-database-sync-outline:{ .lg .middle } __浏览器和外部系统一体化__

    同一条 Flow 里处理 HTTP、CSV、Excel、Redis、数据库、邮件通知，不必拆成多段脚本。

</div>

## 推荐路径

| 目标 | 第一步 | 第二步 | 第三步 |
| --- | --- | --- | --- |
| 完成首次自动化运行 | [快速开始](getting-started.md) | [Lesson 01](docs/tutorials/01-hello-world.md) | [基础学习路线](docs/tutorials/track-newbie.zh-CN.md) |
| 使用 Codex 生成或修改 Flow | [Skills 介绍](docs/skills/README.md) | [AI 协作入门](docs/training/ai-intent-to-flow.md) | [MCP 教程入口](docs/tutorials/119-mcp-chain-overview.md) |
| 给团队做培训 | [培训体系总览](docs/training/README.md) | [训练营课程表](docs/training/bootcamp-plan.md) | [讲师手册](docs/training/trainer-playbook.md) |

## 站内重点入口

<div class="grid cards" markdown>

-   __快速开始__
    [getting-started.md](getting-started.md)

-   __中文项目总览__
    [README.zh-CN.md](README.zh-CN.md)

-   __教程总站__
    [docs/tutorials/README.zh-CN.md](docs/tutorials/README.zh-CN.md)

-   __Skills 介绍__
    [docs/skills/README.md](docs/skills/README.md)

-   __英文课程地图__
    [docs/tutorials/README.md](docs/tutorials/README.md)

-   __文档总地图__
    [docs/README.md](docs/README.md)

</div>
