# TSPlay Skills 介绍

如果你当前最关心的是“在 Codex 中如何自动生成或修改 Flow”，先看这一节。
仓库附带的 `tsplay-flow-authoring` 就是这类工作的主要入口。

## Codex 中可以直接做什么

- 根据自然语言需求自动生成新的 TSPlay `.flow.yaml`
- 修改已有 Flow 的 selector、等待、断言、变量链路和会话复用
- 结合 `finalize -> run -> repair` 继续修复失败的 Flow
- 把需求、修复建议和交付结构整理成更容易审阅的结果

## 在 Codex 中怎么提

在支持 skills 的 Codex 环境里，可以直接在请求中写出 `tsplay-flow-authoring`，也可以写成 `$tsplay-flow-authoring`。

### 生成新 Flow

```text
请使用 tsplay-flow-authoring，帮我生成一条 TSPlay Flow。
- 页面: <URL 或本地页面>
- 目标: <要完成的动作>
- 输入: <关键词 / 文件 / 条件，没有就写无>
- 输出: <JSON / CSV / save_as / artifact>
- 授权: <readonly / browser_write / full_automation / allow_*>
```

### 修改已有 Flow

```text
请使用 tsplay-flow-authoring，修改这条 TSPlay Flow。
- 文件: <flow 文件路径>
- 问题: <超时 / selector 失效 / assert 失败 / 输出为空>
- 预期: <修完后得到什么结果>
- 限制: <不要改业务意图 / 保持 artifact 路径 / 不转 Lua>
```

### 让 Codex 先修再解释

```text
请使用 tsplay-flow-authoring 先修这条 Flow，再告诉我改了什么。
如果缺输入或授权，只告诉我最关键的一项。
```

这页讲的 `skills`，不是运行时 action，也不是 `.flow.yaml` 本身。  
它更像一层“协作说明”：

- 告诉 Codex / Agent 什么时候该怎么做
- 帮团队把高频任务收成统一做法
- 让“写 Flow / 修 Flow / review Flow”有可复用的说明和步骤

## 一句话理解

如果说：

- `action` 是单个能力
- `Flow` 是可执行流程
- `MCP` 是给 Agent 调用的工具入口

那 `skill` 就是“让 Agent 更稳定地使用这些能力和流程的说明包”。

## Skills 和其他概念怎么分

| 概念 | 解决什么问题 | 典型例子 | 产物是什么 |
| --- | --- | --- | --- |
| Action | 单个动作能不能做 | `click`、`read_csv`、`db_query` | 一步能力 |
| Flow | 多个动作怎么串起来 | `import.flow.yaml` | 可执行流程 |
| MCP | Agent 怎么安全调用 TSPlay | `finalize_flow`、`run_flow`、`repair_flow` | 工具接口 |
| Skill | Agent 该按什么方式协作 | `tsplay-flow-authoring` | 一套提示、规则、参考和工作流 |

## 为什么要单独引入 Skills

- 团队里同一类任务会反复出现，比如“写一条 Flow”“修 selector”“补邮件通知”
- 如果每次都临时写 prompt，结果容易不一致
- 如果只给 action 列表，Agent 知道“能做什么”，但不一定知道“先做什么、后做什么”
- Skill 刚好补上这一层，让协作方式可以复用

## 当前仓库里已经提供什么

当前仓库附带一个 skill：

| Skill 名称 | 适合什么场景 | 典型结果 |
| --- | --- | --- |
| `tsplay-flow-authoring` | 写 Flow、修 Flow、review Flow、把中文需求收敛成 Flow、排查 artifact / session 问题 | 一条更可 review、可维护的 TSPlay Flow，或一组清晰的修复建议 |

## 为什么先从 `tsplay-flow-authoring` 开始

如果你在 Codex 中的主要目标就是“让模型帮我做 Flow”，这页应该先于 action 列表阅读。

- action 列表回答的是“单个动作支不支持”
- `tsplay-flow-authoring` 回答的是“如何把需求收敛成一条可执行、可审阅、可修复的 Flow”
- 对生成、修改和修复 Flow 来说，后者通常更接近真实工作方式

## `tsplay-flow-authoring` 会帮什么

- 用自然语言把需求收敛成 TSPlay `.flow.yaml`
- 修已有 Flow 的 selector、等待、断言、变量链路、会话复用
- 编写 `send_email` 相关 Flow
- 按 `finalize -> run -> repair` 的思路组织 Agent 协作
- 中英文都可以触发

## 什么时候优先用 Skill

- 你要“写一条新 Flow”，而不是只查某个 action 参数
- 你要“修一条现有 Flow”，而不是只跑一次命令
- 你要给团队沉淀一类固定协作方式
- 你希望 Codex 输出的是“可 review 的结果”，而不是零散动作

## 什么时候不必先用 Skill

- 你只是想确认某个 action 是否支持
- 你只是想查某个 CLI `-action` 命令
- 你已经有现成 Flow，只差手工改一两个字面量参数

这几种情况，通常先看：

- [支持行为清单](../capability-actions/README.md)
- [CLI `-action` 参考](../actions/README.md)

## 其他触发示例

你可以直接这样提：

```text
帮我写一条 TSPlay Flow。
- 页面: <URL 或本地页面>
- 目标: <要完成的动作>
- 输入: <关键词 / 文件 / 条件，没有就写无>
- 输出: <JSON / CSV / save_as / artifact>
- 授权: <readonly / browser_write / full_automation / allow_*>
```

或者这样提：

```text
帮我修这条 TSPlay Flow。
- 文件: <flow 文件路径>
- 问题: <超时 / selector 失效 / assert 失败 / 输出为空>
- 预期: <修完后应该得到什么结果>
- 限制: <不要改业务意图 / 保持 artifact 路径 / 不转 Lua>
```

## 如果你要在本地安装

普通用户推荐直接从 [TSPlay Releases](https://github.com/tensafe/tsplay/releases/latest) 下载：

- 适合自己系统和 CPU 架构的 `tsplay` 二进制包
- `tsplay-flow-authoring-codex_<version>.zip` 或通用 `tsplay-flow-authoring_<version>.zip`

安装 skill 的常见做法有两种：

1. 解压 release 里的 skill 包，运行目录里的安装脚本，把它安装到本地 Codex `skills/` 目录
2. 直接把整个目录复制到 `~/.codex/skills/` 或 `$CODEX_HOME/skills/`

只有在开发或调试 TSPlay 源码仓库时，才需要从仓库里的 `skills/tsplay-flow-authoring/` 目录安装。

## Release 包里会生成什么

推送 `v*` 版本 tag 或发布 GitHub Release 时，release workflow 会把这份 skill 打成可分发压缩包：

| 资源 | 用途 |
| --- | --- |
| `tsplay-flow-authoring_<version>.zip` | 通用 skill 压缩包 |
| `tsplay-flow-authoring-codex_<version>.zip` | 给 Codex 分发时更直观的命名 |
| `tsplay-flow-authoring-openclaw_<version>.zip` | 给 OpenClaw 或同类 Agent 分发时更直观的命名 |
| `tsplay-skills_<version>.json` | skill 发布 manifest，列出入口文件、安装目录和适配的 Agent 类型 |

安装时把压缩包解到 Codex 的 `skills/` 目录即可，解压后应保留 `tsplay-flow-authoring/SKILL.md` 这一层目录结构。

同一个 Release 里也会提供 macOS、Linux、Windows 的 `tsplay` 二进制包；skill 文档默认假设用户使用这些 release binary 或 PATH 里的 `tsplay`，而不是 `go run .`。

## 推荐阅读顺序

1. 先看这页，弄清 Skill 解决的是哪一层问题
2. 再看 [AI 协作入门](../training/ai-intent-to-flow.md)，理解它和 MCP 的配合方式
3. 再看 [支持行为清单](../capability-actions/README.md)，理解 Skill 底下实际在调用什么能力
4. 真要落到教程和交付，再回到 [教程总览](../tutorials/README.zh-CN.md)

## 相关入口

- [AI 协作入门](../training/ai-intent-to-flow.md)
- [支持行为清单](../capability-actions/README.md)
- [CLI `-action` 参考](../actions/README.md)
- [教程总览](../tutorials/README.zh-CN.md)
