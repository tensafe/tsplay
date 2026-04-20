# TSPlay 培训体系

这套培训体系的目标，不是只让学员“会跑一个例子”，而是让团队形成稳定的浏览器自动化交付能力：

- 新人能在半天内跑通 TSPlay 的基础动作
- 交付成员能在 2 天内写出可审阅、可复用的 Flow
- 平台或 AI 团队能把 `observe -> draft -> validate -> run -> repair` 串起来
- 内部讲师能按统一口径讲授、评审和复盘

## 适用对象

| 角色 | 典型目标 | 推荐入口 |
| --- | --- | --- |
| 实施 / 运营 / 测试 | 能运行现成脚本、定位页面元素、完成基础自动化 | [learning-path.md](learning-path.md) 的 L0-L2 |
| 自动化开发 / RPA 工程师 | 能设计健壮 Flow，处理变量、控制流、失败恢复和会话 | [learning-path.md](learning-path.md) 的 L1-L3 |
| AI / 平台工程师 | 能接入 MCP 工具链和 Agent 工作流 | [learning-path.md](learning-path.md) 的 L2-L4 |
| 讲师 / Enablement | 能组织课程、带实验、做评审、维护版本 | [trainer-playbook.md](trainer-playbook.md) |

## 培训体系包含什么

| 模块 | 用途 | 入口 |
| --- | --- | --- |
| 学习路径 | 定义等级、能力边界、出入门标准 | [learning-path.md](learning-path.md) |
| AI 无感入门 | 面向 Codex、OpenClaw 等 Agent，把“用户意图 -> MCP -> Flow -> 执行修复”串成新手入口 | [ai-intent-to-flow.md](ai-intent-to-flow.md) |
| Bootcamp 课程表 | 帮你把培训排成可执行课程 | [bootcamp-plan.md](bootcamp-plan.md) |
| Labs | 提供基于仓库现有素材的实验任务 | [labs.md](labs.md) |
| Capstone | 提供结业项目场景和交付要求 | [capstone-briefs.md](capstone-briefs.md) |
| 考核体系 | 给出评分标准、证据和晋级门槛 | [assessment.md](assessment.md) |
| 自动录屏 | 帮讲师和维护者把教程过程稳定录成视频素材 | [tutorial-video-recording.md](tutorial-video-recording.md) |
| 讲师手册 | 帮讲师备课、控节奏、收证据和复盘 | [trainer-playbook.md](trainer-playbook.md) |
| 模板 | 学员提交和讲师评审的统一模板 | [templates/](templates/) |

## 建议的交付模式

### 模式 0：AI 无感入门

适合希望“先做成，再理解底层细节”的业务用户、测试同学或产品同学。

- 目标：让用户只描述意图，由 Codex 等 Agent 帮他草拟、校验、执行和修复 Flow
- 时间：1-2 小时
- 输出：1 条 AI 草拟并执行过的 Flow、1 次失败修复记录
- 推荐先阅读：[ai-intent-to-flow.md](ai-intent-to-flow.md)

### 模式 A：半天入门

适合首次接触 TSPlay 的业务侧、测试侧或管理者。

- 目标：理解 TSPlay 三层能力，亲手跑通 1 个 Lua 例子和 1 个 Flow
- 时间：3-4 小时
- 输出：环境检查通过、1 个基础脚本、1 个基础 Flow

### 模式 B：2 天 Bootcamp

适合要开始实际交付 Flow 的成员。

- 目标：完成从 CLI 到 Flow 再到 MCP 的完整链路
- 时间：2 天
- 输出：3-4 个实验、1 个 Capstone、1 次评审

### 模式 C：4 周应用落地

适合团队从试点走向稳定交付。

- 第 1 周：环境铺设、角色分层、完成 L1-L2
- 第 2 周：用真实业务页面改写 1-2 条 Flow
- 第 3 周：引入 MCP、repair 和命名会话
- 第 4 周：结业评审、标准沉淀、指定讲师

## 成功指标

培训是否有效，不看 PPT 是否讲完，主要看这些指标：

- 学员能否独立运行 `go run . -flow ...`
- 学员提交的 Flow 是否通过 `validate_flow`
- 学员是否能用 `retry`、`if`、`on_error`、`wait_until` 处理页面不稳定
- 学员是否理解高风险能力授权边界
- 团队是否能把至少 1 条真实业务流程沉淀到仓库中
- 内部是否出现至少 1 位能带训练营的人

## 培训所依赖的仓库素材

- 页面素材：`demo/`
- 示例脚本：`script/`
- Flow 示例：`script/demo_baidu.flow.yaml`
- MCP 能力：见根目录 [../../ReadMe.md](../../ReadMe.md) 的 MCP 章节

如果希望优先走 AI 路线，建议先阅读 [ai-intent-to-flow.md](ai-intent-to-flow.md)。

这篇文档现在已经覆盖：

- 用户如何描述需求
- 模型如何调用 TSPlay MCP
- 推荐的 system prompt
- `security_preset` 和授权边界
- 失败后的 repair 闭环

如果希望系统学习，再按 [learning-path.md](learning-path.md)、[bootcamp-plan.md](bootcamp-plan.md) 和 [labs.md](labs.md) 执行。
