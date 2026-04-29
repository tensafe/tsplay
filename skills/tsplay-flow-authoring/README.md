# TSPlay Flow Authoring

面向 Codex 的可分享 skill，用来帮助开发者编写、修复、review TSPlay Flow。

这份 skill 重点支持：

- 用自然语言把需求收敛成 TSPlay `.flow.yaml`
- 修已有 Flow 的 selector、等待、断言、变量链路、会话复用
- 用中文或英文提问
- 按场景快速找到起手示例
- 按 MCP 的 `finalize -> run -> repair` 思路组织 Flow 工作流

## 适合什么场景

- 写新的 TSPlay Flow
- 修一条已经存在的 Flow
- 给团队做教程风格 Flow
- review Flow 的可维护性
- 把中文业务需求整理成稳定的 Flow

## 安装方式

### 方式 1：复制到 Codex skills 目录

把整个目录复制到：

```text
~/.codex/skills/tsplay-flow-authoring
```

如果你使用了自定义 `CODEX_HOME`，则复制到：

```text
$CODEX_HOME/skills/tsplay-flow-authoring
```

### 方式 2：直接分发压缩包后解压

把压缩包解压到：

```text
~/.codex/skills/
```

解压后目录结构应当像这样：

```text
~/.codex/skills/tsplay-flow-authoring/
  SKILL.md
  agents/openai.yaml
  references/
```

## 运行前假设

- 优先使用当前 PATH 里的 `tsplay`
- 如果在 TSPlay 仓库内，也可以用 `./tsplay` 或 `go run .`
- 不要求 skill 绑定某个固定绝对路径的可执行文件

## 如何触发这份 skill

典型中文触发词：

- 帮我写一条 TSPlay Flow
- 帮我修这条 Flow
- 把这个需求转成 Flow
- 帮我排查这条 Flow 为什么跑不通
- 帮我按 MCP 思路收敛成 Flow

典型英文触发词：

- write flow
- fix flow
- review flow
- generate flow
- repair selector
- convert a requirement into Flow

## 建议阅读顺序

### 中文用户

1. `references/zh-cn.md`
2. `references/zh-cn-business-templates.md`
3. `references/zh-cn-selectors.md`
4. `references/zh-cn-troubleshooting.md`
5. `references/zh-cn-review-checklist.md`

### 通用入口

1. `references/flow-authoring.md`
2. `references/actions.md`
3. `references/examples.md`
4. `references/example-index.md`

## 文件说明

- `SKILL.md`: skill 触发条件和主工作流
- `agents/openai.yaml`: Codex UI 元数据
- `references/flow-authoring.md`: Flow 编写主指南
- `references/actions.md`: 高频 action 速查
- `references/examples.md`: 提示词模板和最小示例
- `references/example-index.md`: 按场景分类的 repo 起手示例
- `references/zh-cn.md`: 中文总入口
- `references/zh-cn-business-templates.md`: 中文业务场景模板
- `references/zh-cn-troubleshooting.md`: 中文报错与修复对策
- `references/zh-cn-selectors.md`: 中文 selector 策略速查
- `references/zh-cn-review-checklist.md`: 中文 Flow review 清单
- `references/repo-map.md`: 在 TSPlay 仓库里工作时的路径和命令索引

## 最小使用模板

```text
帮我写一条 TSPlay Flow。
- 页面: <URL 或本地页面>
- 目标: <要完成的业务动作>
- 输入: <关键词 / 文件 / 条件，没有就写无>
- 输出: <save_as / JSON / CSV / Excel / artifact 路径>
- 授权: <readonly / browser_write / full_automation / allow_*>
```

## 一个中文修 Flow 模板

```text
帮我修这条 TSPlay Flow。
- 文件: <flow 文件路径>
- 问题: <超时 / 找不到 selector / assert_text 失败 / 输出为空>
- 预期: <修完后应该得到什么结果>
- 限制: <不要改业务意图 / 不要转 Lua / 保持 artifact 路径>
```

## 分享建议

- 对团队内部分发：直接发整个 `tsplay-flow-authoring` 目录或压缩包
- 对已经在用 Codex 的同事：让对方解压到 `~/.codex/skills/`
- 对中文团队：优先让大家从 `references/zh-cn.md` 开始

## 版本说明

这是一份偏“Flow authoring / repair / review”的 TSPlay skill，不是通用浏览器自动化百科。它的目标是帮助开发者更快写出可 review、可维护的 TSPlay Flow。
