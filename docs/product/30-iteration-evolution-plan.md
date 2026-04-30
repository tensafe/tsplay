# TSPlay 30 轮持续进化计划

> 目标：把“分析、梳理、优化、完善、继续迭代”收成一份 30 轮可执行计划，让这件事以后不是靠感觉推进。

## 使用方式

- 每一轮只做一个明确焦点，不贪多。
- 每一轮至少交付 2 类成果：`代码 / 文档 / 检查脚本 / 测试 / 站点入口` 中的任意两类。
- 每 5 轮做一次阶段回看，每 10 轮做一次结构复盘。

## 30 轮总览

### Phase 1：先把入口、主路径、检查机制收紧（01-05）

1. `01` 统一 quick-start、README、教程总站的默认下一步。
2. `02` 补文档健康页，区分坏链和内容断层。
3. `03` 加文档连续性检查脚本，并收成统一 docs suite 接进 CI。
4. `04` 补核心功能路线图。
5. `05` 把路线图继续压成执行面板。

### Phase 2：先把主路径做成产品感（06-10）

6. `06` 收紧 `finalize_flow` 状态机。
7. `07` 补 `needs_input / needs_permission / needs_repair` 的统一输出。
8. `08` 让 Workbench 默认先走 `plan -> finalize -> run`。
9. `09` 给主路径补最小回归测试。
10. `10` 回看 README、AI 入门、MCP 教程是否仍然同口径。

### Phase 3：把 observation 和 selector 真正做成原材料（11-15）

11. `11` 补 observation 的表单、结果区、分页、提示条摘要。
12. `12` 补 selector 候选稳定性排序。
13. `13` 让 draft 优先使用高稳定性 selector。
14. `14` 让 repair 保留次优 selector 作为回退链。
15. `15` 回看 observation、draft、repair 三处输出口径。

### Phase 4：把 repair 闭环从“能修”推到“能证明修了什么”（16-20）

16. `16` repair 输出修改范围。
17. `17` repair 后自动 validate。
18. `18` 输出修复前后差异摘要。
19. `19` 按 selector / wait / var / permission 分类修复建议。
20. `20` 回看 repair 文档、教程、Workbench 页面是否一致。

### Phase 5：把 session、artifact、batch 这些交付底座补稳（21-25）

21. `21` 给 session 增加健康状态和验证时间。
22. `22` Workbench 展示 session 健康状态。
23. `23` 给运行产物补统一 manifest。
24. `24` 给 batch 流程补标准 ledger 和 resume 统计。
25. `25` 回看 handoff、artifact、resume、single-binary 四条交付线是否串起来。

### Phase 6：把外部连接、发布体验、长期演进闭环补全（26-30）

26. `26` 统一 HTTP / Redis / DB / SMTP 的超时、重试、脱敏风格。
27. `27` 给外部连接补 preflight / connection test 入口。
28. `28` 收紧单二进制 first-run 顺序和 release 说明。
29. `29` 做一次 10 条核心功能的阶段验收回看。
30. `30` 生成下一轮 continuation plan，决定下一阶段主题。

## 每轮固定动作

每一轮都建议按同样顺序推进：

1. 先分析当前断层或缺口。
2. 再梳理受影响的入口、模块和文档。
3. 先补最小结构，再补细节优化。
4. 最后跑脚本或测试，留下可复用证据。

## 每轮最小交付标准

一轮结束时，至少满足下面 4 条中的 3 条：

- 有一处真实文件改动
- 有一条可重复执行的检查命令
- 有一条文档入口或代码主路径变得更顺
- 有一条新的“以后不用再靠人工记忆”的约束

## 建议的阶段产物

### 第 5 轮交付

- 文档健康检查
- 核心功能路线图
- 核心功能执行面板

### 第 10 轮交付

- 主路径状态机
- Workbench 默认任务流
- 主路径回归测试

### 第 15 轮交付

- observation 增强
- selector 稳定性机制
- selector 回退策略

### 第 20 轮交付

- repair diff
- repair validate 串联
- repair 分类提示

### 第 25 轮交付

- session 生命周期
- artifact manifest
- batch ledger / resume 规范

### 第 30 轮交付

- 外部连接可靠性层
- release 首跑体验回看
- 下一轮 continuation plan

## 这份计划和现有文档怎么配合

- 路线图：[core-feature-roadmap.md](core-feature-roadmap.md)
- 执行面板：[core-feature-execution-board.md](core-feature-execution-board.md)
- 文档健康页：[../doc-health-audit.md](../doc-health-audit.md)
- 教程演进手册：[../tutorials/evolution-playbook.md](../tutorials/evolution-playbook.md)
- 160 次迭代路线图：[../tutorials/iteration-roadmap-160.md](../tutorials/iteration-roadmap-160.md)
