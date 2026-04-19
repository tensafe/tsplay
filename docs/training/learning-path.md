# TSPlay 学习路径

这份学习路径把学员从“能跑例子”一路带到“能带别人”。每一级都要求有可检查的交付物，而不只是听课。

## 等级定义

| 等级 | 名称 | 目标时间 | 学完后应该能做什么 |
| --- | --- | --- | --- |
| L0 | Awareness | 1-2 小时 | 知道 TSPlay 的三层能力，能看懂脚本和 Flow 的关系 |
| L1 | Operator | 半天 | 能启动 CLI，运行脚本，完成基础元素交互 |
| L2 | Flow Author | 1-2 天 | 能写、改、校验和执行结构化 Flow |
| L3 | Delivery Engineer | 2-5 天 | 能设计健壮 Flow，处理异常、变量、会话和失败修复 |
| L4 | MCP Integrator | 1-2 天 | 能接入 Agent / MCP 工具链，设计 observe-draft-run-repair 流程 |
| L5 | Coach / Trainer | 持续实践 | 能带课、做评审、维护标准与训练素材 |

## 各等级的入门与出门标准

### L0 Awareness

- 先修：无
- 学习重点：
  - 理解 `Lua CLI / Flow / MCP` 的区别
  - 认识 `demo/`、`script/`、`artifacts/`
  - 能读懂一个基础 Flow
- 出门交付：
  - 口头讲清楚 TSPlay 的三层能力
  - 指出仓库里一个 Lua 示例和一个 Flow 示例

### L1 Operator

- 先修：L0
- 学习重点：
  - 启动 `go run . -action cli`
  - 使用 `navigate`、`click`、`type_text`、`wait_for_selector`
  - 识别 selector 和页面状态
- 出门交付：
  - 跑通 1 个 CLI 例子
  - 提交 1 个基础 Lua 脚本
  - 完成 Labs 1-2

### L2 Flow Author

- 先修：L1
- 学习重点：
  - `schema_version`、`vars`、`steps`、`save_as`
  - 命名参数与 `args`
  - `validate_flow`、`run_flow`
  - 失败 trace 和 artifact 的读取
- 出门交付：
  - 提交 2 条基础 Flow
  - 至少 1 条 Flow 使用 `save_as`
  - 完成 Labs 3-4

### L3 Delivery Engineer

- 先修：L2
- 学习重点：
  - `extract_text`、`set_var`
  - `retry`、`if`、`foreach`
  - `on_error`、`wait_until`
  - 文件动作、安全授权、会话管理
- 出门交付：
  - 提交 1 条带控制流和失败恢复的 Flow
  - 能解释为什么某一步需要 `retry` 或 `wait_until`
  - 完成 Labs 5-6

### L4 MCP Integrator

- 先修：L3
- 学习重点：
  - `flow_schema`、`flow_examples`
  - `finalize_flow`、`observe_page`、`draft_flow`
  - `validate_flow`、`run_flow`
  - `repair_flow_context`、`repair_flow`
  - `save_session`、`list_sessions`、`use_session`
- 出门交付：
  - 演示一次从用户意图到自动草拟 Flow 的过程
  - 演示一次失败 Flow 的修复闭环
  - 完成至少 1 个 MCP 驱动的 Capstone

### L5 Coach / Trainer

- 先修：L4
- 学习重点：
  - 课程节奏控制
  - 常见卡点诊断
  - 学员作品评审
  - 版本更新后的训练材料维护
- 出门交付：
  - 主持 1 次训练营或内训
  - 完成 3 份 Capstone 评审
  - 提交 1 次培训复盘

## 按角色推荐路线

| 角色 | 推荐等级路径 | 说明 |
| --- | --- | --- |
| 测试 / 运营 / 支持 | L0 -> L1 -> L2 | 先会运行和维护已有流程，再学习结构化 Flow |
| 自动化开发 / RPA | L0 -> L1 -> L2 -> L3 | 重点放在健壮性、异常处理和可维护性 |
| AI / 平台工程师 | L0 -> L2 -> L3 -> L4 | 重点放在 MCP、会话与 Agent 交互 |
| 讲师 / Enablement | L0 -> L2 -> L4 -> L5 | 先掌握交付，再掌握教学和标准化 |

## 学习节奏建议

### 个人自学

- 第 1 天：L0-L1
- 第 2 天：L2
- 第 3-4 天：L3
- 第 5 天：L4 或结业项目

### 团队训练营

- Day 0：预习与环境检查
- Day 1：L0-L2
- Day 2：L3-L4
- Day 3 以后：Capstone 与业务实战

## 证据优先原则

每一级都要求有证据，而不是只看自报掌握程度。推荐证据包括：

- 可运行的 Lua 脚本
- 可校验通过的 Flow
- `artifacts/` 中的截图或 trace
- MCP 调用记录
- 实验报告和 Capstone 评审表

对应模板见 [templates/lab-report.md](templates/lab-report.md) 和 [templates/capstone-review.md](templates/capstone-review.md)。
