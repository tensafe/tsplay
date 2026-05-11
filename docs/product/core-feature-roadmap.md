# TSPlay 核心功能路线图（Top 10）

> 目标：把前面提炼出的 10 条核心功能，从“方向清单”补成可排期、可验收、可对代码落点的版本。

## 这份文档怎么用

- 如果你现在要定版本优先级，先看“优先级总览”和“建议里程碑”。
- 如果你现在要开工实现，直接看每条能力下面的 `P0 / P1 / 验收 / 代码落点`。
- 如果你现在要对外讲产品，配合 [PRD 摘要版](PRD-summary.md) 和 [V1 落地任务清单](V1-task-list.md) 一起用。
- 如果你现在要按周排活或拆实施面，再看 [核心功能执行面板](core-feature-execution-board.md)。
- 如果你现在要按轮推进长期演化，再看 [30 轮持续进化计划](30-iteration-evolution-plan.md)。

## 优先级总览

| 编号 | 核心功能 | 建议优先级 | 为什么现在做 |
| --- | --- | --- | --- |
| 1 | 自然语言一键收敛成可运行 Flow | `P0` | 这是用户最短主路径，直接决定“能不能先跑起来” |
| 2 | 页面观察能力升级 | `P0` | 观察质量决定后续 draft、repair、Workbench 卡片质量 |
| 3 | Selector 稳定性与多候选回退 | `P0` | 页面一改就坏是自动化最大痛点 |
| 4 | 自动修复闭环质量 | `P0` | 没有 repair 闭环，Flow 很难长期交付 |
| 5 | Flow 校验器从语法校验走向语义校验 | `P1` | 能减少无效运行和低质量草稿 |
| 6 | 会话与登录态生命周期管理 | `P0` | 登录态是企业场景真正可复用的入口 |
| 7 | 批量任务的断点续跑与幂等控制 | `P1` | 一旦进入真实导入/同步场景，这条会迅速变成刚需 |
| 8 | 运行产物与证据中心 | `P1` | 交付、review、repair、培训都会依赖证据索引 |
| 9 | 外部系统连接可靠性层 | `P2` | 影响 HTTP、Redis、DB、邮件等混合流程的稳定性 |
| 10 | 单二进制首跑体验 | `P2` | 影响传播、试用、离线交付和 release 包体验 |

## 建议里程碑

### M1：先把默认主路径收紧

- 覆盖能力：`1 + 2 + 3 + 4`
- 目标：让用户从一句需求进入 `finalize_flow -> run -> repair`，第一次就更容易跑通

### M2：把登录态和校验变成稳定底座

- 覆盖能力：`5 + 6`
- 目标：减少“明明有 session 却接不进去”和“Flow 看起来对，运行才发现不对”

### M3：把真实交付场景补齐

- 覆盖能力：`7 + 8 + 9`
- 目标：让批量导入、外部同步、失败追查和交接都更像正式产品

### M4：把首跑和发布体验补顺

- 覆盖能力：`10`
- 目标：让 release 包更像“拿来就能试”的产品，而不是“源码用户专属”

## 1. 自然语言一键收敛成可运行 Flow

### 目标

让用户优先走一条默认主线：输入意图后，系统直接返回“能不能跑、缺什么、下一步是什么”，而不是让用户自己拼 `observe / draft / validate / run / repair`。

### 现有基础

- 已有 `tsplay.finalize_flow`、`tsplay.draft_flow`、`tsplay.validate_flow`、`tsplay.run_flow`
- Workbench 已有任务规划和 provider 生成能力
- Flow schema、generation rules、action manifest 已经能喂给模型

### 当前缺口

- `finalize_flow` 还更像“工具链包装”，还不是“默认产品主路径”
- 缺少统一的状态机视图，例如 `ready / needs_input / needs_permission / needs_repair`
- 缺少“为什么这样收敛”的解释和建议下一步

### P0

- 把 `finalize_flow` 的返回标准化成面向用户的状态机
- 当缺输入、缺权限、缺 session 时，返回统一补齐建议
- 在 Workbench 里把 `plan -> finalize -> run` 收成一个默认动作

### P1

- 支持基于历史成功 Flow 的相似任务收敛
- 给出“为什么优先选 UI / API / hybrid”的简短解释

### 验收

- 同一个需求，默认只需要调一次主入口就能得到下一步
- 用户能区分“能直接跑”和“要先补条件”
- Workbench 能展示可读状态，而不是只吐原始 YAML

### 代码落点

`tsplay_core/tsplay_mcpserver.go`、`tsplay_core/tsplay_flow_ai.go`、`tsplay_core/workbench_plan.go`、`tsplay_core/workbench_plan_ai.go`、`tsplay_core/workbench_server.go`

## 2. 页面观察能力升级

### 目标

让 `observe_page` 和 Workbench explore 输出真正成为 Flow 草拟和修复的可靠原材料，而不是只给一份“可看但不够用”的 observation。

### 现有基础

- `PageObservation` 已包含 `page_summary / screenshot / dom_snapshot / content_elements / selector_candidates`
- Workbench explore 已能生成 page card、api card、entity card

### 当前缺口

- 表单、表格、分页、结果区、登录态线索的抽取还不够系统
- 页面观察和站点级探索之间还没形成统一卡片模型
- 对复杂页面的“关键区域”识别还偏弱

### P0

- 强化表单字段、表格列、主要按钮、分页器、空状态、提示条抽取
- 补充登录态、筛选区、结果区、详情入口等业务线索
- 让 Workbench page card 直接复用 observation 里的关键字段

### P1

- 补“页面关键区块”评分
- 补菜单树、面包屑、关联链接的更稳定抽取

### 验收

- 对典型后台页，观察结果能直接回答“去哪输入、点什么、结果在哪”
- draft/repair 使用 observation 后，对用户追问次数下降
- Workbench 页面卡片更像“业务页面摘要”，不是 DOM 片段堆砌

### 代码落点

`tsplay_core/tsplay_observe.go`、`tsplay_core/tsplay_selector_observation.go`、`tsplay_core/workbench_explore.go`、`tsplay_core/workbench_models.go`

## 3. Selector 稳定性与多候选回退

### 目标

让 selector 从“单点命中”升级成“可回退的候选链”，降低页面轻微改版后的失效率。

### 现有基础

- observation 已返回 `primary_selector`、`selector_candidates`、`selector_details`
- generation rules 已明确偏好稳定 selector
- draft/repair 已有 selector repair 基础逻辑

### 当前缺口

- 缺少统一的 selector 稳定性评分
- 还没有把候选 selector 真正贯穿到 draft、run、repair
- 缺少“这个 selector 为什么更稳”的解释信息

### P0

- 给 selector 增加来源类型和稳定性排序，例如 `testid / role / label / placeholder / text`
- draft 时优先选高稳定性候选，repair 时保留次优回退链
- 对常见脆弱 selector 给出显式 warning

### P1

- 为页面卡片和 repair context 输出“推荐 selector 策略”
- 对同一元素保留跨次观察的一致 selector 记忆

### 验收

- 页面做轻微结构改动时，Flow 仍能通过候选链恢复
- repair 输出里能看见“旧 selector -> 新 selector”的依据
- observation、draft、repair 三处 selector 排序一致

### 代码落点

`tsplay_core/tsplay_selector_observation.go`、`tsplay_core/tsplay_flow_draft.go`、`tsplay_core/tsplay_flow_guidance.go`、`tsplay_core/tsplay_flow_repair.go`

## 4. 自动修复闭环质量

### 目标

让 repair 从“给模型一段失败上下文”升级成“按失败步精准修、修完就能再验证”的闭环。

### 现有基础

- 已有 `repair_flow_context`、`repair_flow`
- repair context 已能输出失败分类、失败步、附近步骤、artifact 摘要、repair hints
- Workbench API 已有 repair 相关入口

### 当前缺口

- 还缺修复前后 diff、修复置信度、再验证结果
- 自动修复仍偏一次性，不够 step-scoped
- UI / API / hybrid 失败后的切换策略还不够清楚

### P0

- 输出修复影响范围，例如改了哪些 step/path
- repair 后自动串一次 validate
- 对 selector、等待条件、变量引用、权限问题做分类型修复建议

### P1

- 支持“单步修复”和“整条 Flow 修复”两种模式
- 失败后自动建议切到 `api_first` 或 `hybrid`

### 验收

- 一次 repair 结果里能看到失败点、修复点、验证结果
- 不需要把整份 HTML 塞给模型，也能完成常见修复
- 用户能读懂“为什么这次修成了/还没修成”

### 代码落点

`tsplay_core/tsplay_flow_repair.go`、`tsplay_core/tsplay_flow_repair_request.go`、`tsplay_core/tsplay_repair_hint.go`、`tsplay_core/tsplay_mcp_tool_response.go`、`tsplay_core/workbench_server.go`

## 5. Flow 校验器从语法校验走向语义校验

### 目标

在浏览器实际运行前，尽可能把“会失败但语法没错”的 Flow 提前拦下来。

### 现有基础

- `ValidateFlow` 已覆盖 schema、动作参数、安全策略和部分结构约束
- `flow_issue` 和 `action_capabilities` 已能输出较好的 issue/suggestion

### 当前缺口

- 变量依赖、artifact 目录规范、步骤顺序、命名可读性等语义校验还不够
- 还没有统一的 review 级 warning 模型
- 缺少“修复建议优先级”

### P0

- 补变量引用链和 `save_as` 依赖检查
- 补 `use_session / storage_state / persistent profile` 组合冲突检查
- 补 artifact 输出路径、目录骨架、命名可读性 warning

### P1

- 增加“交付可 review 性”规则，例如步骤命名、证据留痕、危险动作隔离
- 在 Workbench 里按 severity 展示 validation issues

### 验收

- 一批明显会失败的 Flow 能在运行前被拦住
- validation 输出能明确区分 `error / warning / suggestion`
- AI 生成的 Flow 更少出现“语法对、业务不对”的情况

### 代码落点

`tsplay_core/tsplay_flow.go`、`tsplay_core/tsplay_flow_issue.go`、`tsplay_core/tsplay_action_capabilities.go`、`tsplay_core/tsplay_mcpserver.go`

## 6. 会话与登录态生命周期管理

### 目标

让 `save_session -> list_sessions -> use_session` 变成真正可长期复用的登录态体系，而不是一组分散命令。

### 现有基础

- 已有保存、读取、列出、导出、删除命名会话
- 已有 session ownership、last_used_at 等元数据
- Workbench API 已经可以列出和保存 session

### 当前缺口

- 缺少“这个 session 还能不能用”的健康状态
- 缺少与站点配置的强绑定关系
- 缺少 refresh / clone / import / export 的一致体验

### P0

- 给 session 增加 `site_id / last_verified_at / status` 等生命周期信息
- 增加会话健康检查或轻量验证入口
- Workbench 里补“测试会话”“绑定站点”“失效提示”

### P1

- 支持 session 导入导出标准包
- 支持多角色 session 对比和切换

### 验收

- 用户能知道哪个 session 是可用的、最近用于哪个站点
- Flow 失败时能区分“页面问题”还是“登录态失效”
- Workbench 能把保存会话和使用会话连成一条顺手路径

### 代码落点

`tsplay_core/tsplay_session_registry.go`、`tsplay_core/tsplay_session_export.go`、`tsplay_core/workbench_server.go`、`main.go`

## 7. 批量任务的断点续跑与幂等控制

### 目标

让 CSV/Excel 导入、批量处理、长链路同步这类任务从“能跑”变成“能重跑、能续跑、不重复写”。

### 现有基础

- `foreach` 已支持 `progress_key`
- `read_csv / read_excel` 已支持 `start_row`
- 已有 `append_var / on_error / write_json / write_csv`
- 已有 Redis checkpoint 测试和示例模板

### 当前缺口

- 缺少统一 ledger 结构
- 缺少“已成功项跳过”和“失败项重试”的标准模式
- 缺少运行后摘要统计

### P0

- 约定标准 ledger 输出结构，例如 `source_row / status / error / artifact`
- `foreach` 输出增加成功、失败、跳过统计
- Workbench 和教程里沉淀一套标准 resume/import 模板

### P1

- 支持按 ledger 自动跳过已完成项
- 支持 chunk 级恢复和重试窗口

### 验收

- 中途失败后，用户知道从哪一行继续
- 重跑不会重复提交已经成功的项
- 产物里能明确看到处理总数、成功数、失败数、跳过数

### 代码落点

`tsplay_core/tsplay_flow.go`、`tsplay_core/tsplay_data_actions.go`、`tsplay_core/tsplay_action.go`、`tsplay_core/tsplay_flow_ai.go`

## 8. 运行产物与证据中心

### 目标

把 `artifacts/` 从“运行顺手落文件”升级成“可索引、可回放、可交接”的证据中心。

### 现有基础

- 已有 `artifact_root / run_root`
- 观察和运行已经会产出截图、HTML、DOM snapshot、trace
- 教程体系里已经有 manifest、summary、handoff 的思路

### 当前缺口

- 还缺统一 run manifest
- 缺少 step 级 artifact 索引和浏览入口
- 缺少“这次运行到底生成了什么”的总览

### P0

- 每次运行生成统一 `manifest.json`
- 把 step、artifact、输出文件、错误证据索引起来
- Workbench 增加 artifact 列表和预览入口

### P1

- 支持跨 run 搜索 artifact
- 支持按页面、API、任务、站点聚合证据

### 验收

- 用户不需要翻目录就能知道本次运行留下了哪些证据
- repair context 能直接引用 manifest 中的关键 artifact
- 交付和培训可以直接复用一份标准 evidence 包

### 代码落点

`tsplay_core/tsplay_observe.go`、`tsplay_core/tsplay_mcp_tool_response.go`、`tsplay_core/workbench_server.go`、`tsplay_core/workbench_explore.go`

## 9. 外部系统连接可靠性层

### 目标

把 HTTP、Redis、DB、邮件这些能力从“分别可用”补成“统一可靠”。

### 现有基础

- 已有 `http_request`、Redis、DB、邮件动作
- 已有安全策略和 capability 分析
- 已有多项动作级测试

### 当前缺口

- 缺少统一的连接配置视图
- 缺少超时、重试、脱敏、错误分类的一致约定
- 缺少 preflight 自检

### P0

- 统一连接配置命名和错误消息风格
- 为外部连接补 timeout/retry/backoff 基线
- 对敏感头、cookie、token、SMTP 凭据做更一致的脱敏输出

### P1

- 增加 `connection test / preflight` 能力
- 让 Workbench 能展示连接状态和最近失败原因

### 验收

- 外部连接失败时，错误信息能明确落到“配置错 / 网络错 / 权限错 / 业务错”
- 运行日志里不会随手泄露敏感信息
- 混合 Flow 的稳定性不再只依赖单个动作实现质量

### 代码落点

`tsplay_core/tsplay_http_actions.go`、`tsplay_core/tsplay_redis_actions.go`、`tsplay_core/tsplay_db_actions.go`、`tsplay_core/tsplay_db_actions_ext.go`、`tsplay_core/tsplay_email_actions.go`

## 10. 单二进制首跑体验

### 目标

让 release 包用户不依赖源码仓库，也能知道“第一步该做什么、哪些资料已经在二进制里、什么时候需要 file-srv”。

### 现有基础

- 已有 `list-assets`、`extract-assets`、`file-srv`
- 已有内置资源打包和 release workflow
- 教程 `141-150` 已经覆盖单二进制与 first-run 心智

### 当前缺口

- 首跑入口还分散在 README、教程、release 资产里
- 命令可用，但“该先跑哪条”仍需要用户自己判断
- 缺少 release 包内的统一 onboarding 说明

### P0

- 统一首跑入口文案和命令顺序
- release 产物里固定附带 `README`、`LICENSE`、教程入口说明
- 把 `list-assets -> extract-assets -> getting-started` 收成一条默认路径

### P1

- 增加更清晰的首跑 manifest 或内置帮助
- 区分开发态 `file-srv` 和发布态 `file-srv` 的推荐用法

### 验收

- 没看源码的人也能在 release 包里找到第一步
- 新用户能理解什么时候直接跑 Flow，什么时候先 `extract-assets`
- 单二进制交付更适合培训、试用和离线环境

### 代码落点

`main.go`、`embedded_assets.go`、`workbench_app.go`、`.github/workflows/release-binaries.yml`、`getting-started.md`

## 建议下一步

如果接下来要直接进入实现，我建议先按下面顺序开工：

1. `P0-1`：先补 `finalize_flow` 主状态机和 Workbench 默认主路径。
2. `P0-2`：同时增强 observation 和 selector 稳定性，因为这两项会直接影响生成质量。
3. `P0-3`：把 repair 闭环和 session 生命周期补到“可持续用”的程度。
4. `P1-1`：再补语义校验、批量续跑和 artifact manifest。

这四步完成后，TSPlay 会明显从“能力很多”进入“主路径更顺、交付更稳”的阶段。
