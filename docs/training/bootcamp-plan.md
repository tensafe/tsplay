# TSPlay Bootcamp 课程表

这份课程表默认对应“2 天训练营 + 4 周应用落地”。你可以按团队成熟度裁剪，但建议保留“讲解 -> 演示 -> 实操 -> 评审”的节奏。

## Day 0：预习与环境检查

目标：避免正式开课时大家把时间全花在装环境上。

### 学员预习包

- 阅读 [../../ReadMe.md](../../ReadMe.md) 的“运行模式”“快速开始”
- 阅读 [learning-path.md](learning-path.md) 的 L0-L2
- 准备好 Go 与 Playwright 运行环境
- 熟悉仓库中的：
  - [../../script/open_url.lua](../../script/open_url.lua)
  - [../../script/demo_baidu.flow.yaml](../../script/demo_baidu.flow.yaml)
  - [../../demo/demo.html](../../demo/demo.html)

### 讲师检查项

- 确认每位学员都能访问仓库
- 确认至少 1 位助教能现场处理环境问题
- 预先安排好本地静态页面的暴露方式，确保能访问 `/demo/*.html`

## Day 1：从动作到基础 Flow

| 时段 | 模块 | 目标 | 方式 | 输出 |
| --- | --- | --- | --- | --- |
| 09:30-10:15 | 模块 1：架构导览 | 讲清楚 Lua、Flow、MCP 三层能力 | 讲解 + 演示 | 学员知道用什么模式解决什么问题 |
| 10:15-11:15 | 模块 2：动作基础 | 练会 `navigate/click/type_text/wait_for_selector` | 演示 + 跟练 | 跑通 CLI |
| 11:15-12:00 | Lab 1 | 在本地 demo 页面上完成基础交互 | 个人实操 | 1 个 Lua 脚本 |
| 13:30-14:30 | 模块 3：Flow 入门 | 讲清楚 `vars/steps/save_as/args` | 讲解 + 改写示例 | 读懂 Flow |
| 14:30-15:30 | Lab 2 | 把 Lua 脚本改成 Flow | 个人实操 | 1 条基础 Flow |
| 15:45-16:45 | 模块 4：校验与 trace | 学会 `validate_flow`、`run_flow`、artifact 定位 | 演示 + 讲解 | 知道如何看失败现场 |
| 16:45-17:30 | 评审环节 | 讲师挑 2-3 份作品做现场 review | 小组复盘 | 修正命名、selector、结构问题 |

## Day 2：从健壮 Flow 到 MCP

| 时段 | 模块 | 目标 | 方式 | 输出 |
| --- | --- | --- | --- | --- |
| 09:30-10:30 | 模块 5：高级控制流 | 学会 `extract_text/set_var/retry/if/foreach` | 讲解 + 小例子 | 1 条增强版 Flow |
| 10:30-11:30 | Lab 3 | 给已有 Flow 加变量与控制流 | 个人实操 | 1 条可复用 Flow |
| 11:30-12:00 | 模块 6：失败恢复 | 学会 `on_error/wait_until` 与 artifact 复盘 | 演示 | 失败恢复思路 |
| 13:30-14:30 | 模块 7：MCP 工具链 | 讲清楚 `observe -> draft -> validate -> run -> repair` | 演示 + Q&A | 看懂 Agent 闭环 |
| 14:30-15:30 | Lab 4 | 用 MCP 工具完成一次草拟和一次修复 | 小组实操 | 1 份 MCP 操作记录 |
| 15:45-16:45 | Capstone | 选 1 个业务场景做结业项目 | 小组实操 | Capstone 方案 |
| 16:45-17:30 | 结业评审 | 按评分表评 Capstone | 讲师评审 | 结业结果 |

## 4 周落地节奏

### 第 1 周：共识建立

- 完成 Day 0 与 Day 1 内容
- 每位学员至少提交 1 条基础 Flow

### 第 2 周：真实业务试点

- 选择 1-2 条真实业务流程
- 用 Flow 替换原始零散脚本
- 开始记录失败现场和修复方式

### 第 3 周：MCP 与会话能力

- 引入 `finalize_flow`、`observe_page`、`draft_flow`、`repair_flow`
- 对需要登录态的场景试用 `save_session` / `use_session`

### 第 4 周：标准化与传帮带

- 完成 Capstone
- 指定内部讲师
- 将有效流程、讲义和评分表沉淀回仓库

## 课程交付原则

- 每 45-60 分钟必须有一次实操
- 每个模块至少要产出一个可检查的文件、Flow 或截图
- 讲师不要只讲“怎么写”，要讲“为什么这样更稳”
- 评审优先讲 selector 质量、控制流设计和失败恢复，不只看是否跑通
