---
hide:
  - toc
  - title
search:
  boost: 2
---

<div class="tsplay-landing" markdown="1">

<section class="tsplay-cosmos tsplay-cosmos--hero" markdown="1">
<div class="tsplay-hero-grid" markdown="1">
<div markdown="1">

<p class="tsplay-cosmos__kicker">TSPlay Docs</p>

# 给人和 AI 的浏览器自动化平台

<p class="tsplay-cosmos__lede">用一套运行时把 <code>Flow DSL</code>、<code>MCP</code>、<code>Skills</code>、会话复用、失败修复与交付证据收成同一条主线。你可以先一键跑起来，也可以直接让 Codex 生成或修改 Flow，再继续扩展到培训、交付和 Agent 集成。</p>

<div class="tsplay-pill-row">
<span class="tsplay-pill tsplay-pill--accent">Flow DSL</span>
<span class="tsplay-pill">MCP</span>
<span class="tsplay-pill">Skills</span>
<span class="tsplay-pill tsplay-pill--mint">single binary</span>
</div>

<div class="tsplay-cta-row">
<a class="md-button md-button--primary" href="getting-started/">快速开始</a>
<a class="md-button" href="docs/skills/">在 Codex 中使用 Skills</a>
</div>

<p class="tsplay-terminal__note">更偏教程学习，从 <a href="docs/tutorials/README.zh-CN/">教程总览</a> 进入；更偏 Agent 路线，从 <a href="docs/training/ai-intent-to-flow/">AI 协作入门</a> 进入。</p>

</div>
<div class="tsplay-install-card" markdown="1">

<p class="tsplay-install-card__eyebrow">One-line Quickstart</p>

```bash
curl -fsSL https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.sh | sh
```

<div class="tsplay-metric-row">
  <div class="tsplay-metric">
    <strong>自动识别系统</strong>
    <span>macOS / Linux 用同一条命令进入</span>
  </div>
  <div class="tsplay-metric">
    <strong>不先等 Playwright</strong>
    <span>先生成 demo Flow 并拿到结构化结果</span>
  </div>
  <div class="tsplay-metric">
    <strong>直接写 artifacts</strong>
    <span>结果输出到 <code>artifacts/quickstart/</code></span>
  </div>
  <div class="tsplay-metric">
    <strong>适合团队分发</strong>
    <span>继续扩展到培训包和单二进制交付</span>
  </div>
</div>

<p class="tsplay-terminal__note">Windows 可在 <a href="getting-started/">快速开始</a> 页直接切到 PowerShell 命令。</p>

</div>
</div>
</section>

<section class="tsplay-landing-section" markdown="1">

<p class="tsplay-section-kicker">What It Does</p>

## 你可以直接用它做什么

<div class="grid cards" markdown="1">

-   __在 Codex 中生成或修改 Flow__

    仓库附带 `tsplay-flow-authoring`，可以直接把页面目标、输入输出和权限边界整理成 TSPlay Flow。

    [查看 Skills 介绍与 Codex 用法](docs/skills/README.md)

-   __跑浏览器动作并留下可复盘证据__

    同时处理页面动作、断言、提取、下载、截图、HTML、trace 和输出产物，不必拆成多段脚本。

    [查看支持行为清单](docs/capability-actions/README.md)

-   __接到外部系统和批处理流程__

    在同一条 Flow 里处理 HTTP、CSV、Excel、Redis、数据库、邮件通知和 artifact 汇总。

    [查看能力动作参考](docs/capability-actions/README.md)

-   __进入 MCP / Agent 主链__

    把 `observe -> draft -> validate -> run -> repair` 接进 Agent 产品，而不是只停留在本地脚本。

    [进入 AI 协作入门](docs/training/ai-intent-to-flow.md)

</div>

</section>

<section class="tsplay-landing-section" markdown="1">

<p class="tsplay-section-kicker">How To Start</p>

## 三条最常用的进入方式

<div class="grid cards" markdown="1">

-   __先跑通一次__

    第一次接触 TSPlay，先不要纠结 `Lua / Flow / MCP` 的差别，先完成一次运行，再决定下一条线。

    [进入快速开始](getting-started.md)

-   __先走 Codex + Skills__

    如果你的目标是“让 AI 帮你产出 Flow”，先把 `tsplay-flow-authoring` 的使用方式吃透。

    [进入 Skills 介绍](docs/skills/README.md)

-   __先走教程金字塔__

    如果你想系统学，教程不会先把 `160` 节课摊给你自己筛，而是先按“跑通动作 -> 组织成 Flow -> 接进真实交付 -> 进入 Agent / MCP / 标准化”四层往上走。

    [进入教程总览](docs/tutorials/README.zh-CN.md)

-   __先走 Agent / MCP__

    如果你更关心能力暴露和主链工具，把服务起起来，再看 `finalize_flow`、repair 和 session。

    [进入 AI 协作入门](docs/training/ai-intent-to-flow.md)

</div>

</section>

<section class="tsplay-landing-section" markdown="1">

<p class="tsplay-section-kicker">Core Surface</p>

## 核心能力面

<div class="grid cards" markdown="1">

-   __页面动作、提取与断言__

    `navigate`、`click`、`type_text`、`extract_text`、`assert_visible`、`assert_text` 等动作可直接组成 Flow。

-   __Flow 控制流与批处理__

    `retry`、`if`、`foreach`、`on_error`、`wait_until` 等能力适合把一次成功变成稳定重复运行。

-   __浏览器状态与会话复用__

    支持 `get_storage_state`、cookie、保存命名 session、导入导出登录态，适合进入真实系统。

-   __外部系统和交付物整合__

    支持文件、表格、HTTP、Redis、数据库、邮件，以及面向 handoff 的 artifacts 输出。

</div>

</section>

<section class="tsplay-landing-section" markdown="1">

<p class="tsplay-section-kicker">Key Entry</p>

## 站内重点入口

<div class="grid cards" markdown="1">

-   __快速开始__

    一条命令下载对应平台二进制，自动生成 demo Flow 并立刻执行。

    [打开快速开始](getting-started.md)

-   __Skills 介绍__

    重点看 “在 Codex 中如何提要求、如何让 skill 生成或修改 Flow”。

    [打开 Skills 介绍](docs/skills/README.md)

-   __支持行为清单__

    按能力面查询支持的 action、约束边界和 Flow / Lua / MCP 的对应关系。

    [打开支持行为清单](docs/capability-actions/README.md)

-   __教程总览__

    先用四层金字塔确定“今天该走哪一层”，再决定是否展开逐课 lesson 地图，而不是先面对一整张 160 节课的大表。

    [打开教程总览](docs/tutorials/README.zh-CN.md)

</div>

</section>

</div>
