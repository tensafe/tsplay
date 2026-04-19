# TSPlay Capstone 场景

Capstone 用来检验学员是否已经具备“独立交付”的能力。建议每位学员或每个小组至少完成 1 个场景，并接受一次评审。

## 统一交付要求

每个 Capstone 都至少要提交：

- 1 条可运行的 Flow
- 1 份设计说明
- 1 份实验或评审记录
- 关键 artifact 或 trace 截图

## 场景 A：本地演示站点自动化套件

- 难度：L2-L3
- 背景：团队需要把仓库里的本地 demo 页面变成一套稳定的回归验证样例
- 要求：
  - 至少覆盖 3 个页面
  - 至少包含 1 个文件上传动作
  - 至少包含 1 个数据提取动作
  - 至少包含 1 个控制流动作
- 推荐素材：
  - [../../demo/demo.html](../../demo/demo.html)
  - [../../demo/tables.html](../../demo/tables.html)
  - [../../demo/upload.html](../../demo/upload.html)
- 评分重点：
  - 结构清晰
  - selector 稳定
  - 健壮性设计合理

## 场景 B：业务页面观察到 Flow 交付

- 难度：L3-L4
- 背景：给定一个真实业务 URL，让学员从页面观察一路走到可运行 Flow
- 要求：
  - 先做 `finalize_flow`；需要细粒度控制时再做 `observe_page` 或 `draft_flow`
  - 再做 `validate_flow` 和 `run_flow`
  - 至少演示一次失败后修复
- 评分重点：
  - 是否理解 MCP 工具链顺序
  - 是否能利用 artifact 和 repair_hints
  - 是否能收敛到可复用 Flow

## 场景 C：团队级长期复用方案

- 难度：L4-L5
- 背景：团队需要把 TSPlay 引入到长期运营中，而不是只做一次性脚本
- 要求：
  - 定义一条业务 Flow 的长期维护策略
  - 说明会话保存策略
  - 说明培训和评审怎样接入项目流程
  - 给出最少 1 条团队规范
- 评分重点：
  - 方案的可维护性
  - 安全边界意识
  - 培训体系是否能支撑业务复制

## 结业展示建议

每组展示控制在 10-15 分钟：

1. 业务背景
2. Flow 或 MCP 方案
3. 关键设计点
4. 一次失败与修复
5. 复用和维护建议

讲师评审时，建议使用 [templates/capstone-review.md](templates/capstone-review.md)。
