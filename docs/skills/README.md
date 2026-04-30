# TSPlay Skills 介绍

这页讲的 `skills`，不是运行时 action，也不是 `.flow.yaml` 本身。  
它更像一层“协作工作法”:

- 告诉 Codex / Agent 什么时候该怎么做
- 帮团队把高频任务收成稳定套路
- 让“写 Flow / 修 Flow / review Flow”不再每次都从零开始提示

## 一句话理解

如果说：

- `action` 是单个能力
- `Flow` 是可执行流程
- `MCP` 是给 Agent 调用的工具入口

那 `skill` 就是“让 Agent 更稳定地使用这些能力和流程的方法包”。

## Skills 和其他概念怎么分

| 概念 | 解决什么问题 | 典型例子 | 产物是什么 |
| --- | --- | --- | --- |
| Action | 单个动作能不能做 | `click`、`read_csv`、`db_query` | 一步能力 |
| Flow | 多个动作怎么串起来 | `import.flow.yaml` | 可执行流程 |
| MCP | Agent 怎么安全调用 TSPlay | `finalize_flow`、`run_flow`、`repair_flow` | 工具接口 |
| Skill | Agent 该按什么套路协作 | `tsplay-flow-authoring` | 一套提示、规则、参考和工作流 |

## 为什么要单独引入 Skills

- 团队里同一类任务会反复出现，比如“写一条 Flow”“修 selector”“补邮件通知”
- 如果每次都临时写 prompt，质量会漂
- 如果只给 action 列表，Agent 知道“能做什么”，但不一定知道“先做什么、后做什么”
- Skill 刚好补上这一层，让协作方式可以复用

## 当前仓库里已经提供什么

当前仓库已经附带一个可分享 skill：

| Skill 名称 | 适合什么场景 | 典型结果 |
| --- | --- | --- |
| `tsplay-flow-authoring` | 写 Flow、修 Flow、review Flow、把中文需求收敛成 Flow、排查 artifact / session 问题 | 一条更可 review、可维护的 TSPlay Flow，或一组清晰的修复建议 |

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

## 最小触发示例

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

仓库里当前的 skill 目录是：

```text
skills/tsplay-flow-authoring/
```

常见做法有两种：

1. 运行目录里的安装脚本，把它安装到本地 Codex `skills/` 目录
2. 直接把整个目录复制到 `~/.codex/skills/` 或 `$CODEX_HOME/skills/`

## 推荐阅读顺序

1. 先看这页，弄清 Skill 解决的是哪一层问题
2. 再看 [AI 无感入门](../training/ai-intent-to-flow.md)，理解它和 MCP 的配合方式
3. 再看 [支持行为清单](../capability-actions/README.md)，理解 Skill 底下实际在调用什么能力
4. 真要落到教程和交付，再回到 [教程总览](../tutorials/README.zh-CN.md)

## 相关入口

- [AI 无感入门](../training/ai-intent-to-flow.md)
- [支持行为清单](../capability-actions/README.md)
- [CLI `-action` 参考](../actions/README.md)
- [教程总览](../tutorials/README.zh-CN.md)
