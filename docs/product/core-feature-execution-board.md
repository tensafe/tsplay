# TSPlay 核心功能执行面板

> 目标：把 [核心功能路线图](core-feature-roadmap.md) 再压成“可以直接开工”的执行面板。

## 怎么配合路线图使用

- 路线图回答的是：为什么做、先做什么、做到什么算完成。
- 这份执行面板回答的是：第一刀先切哪里、依赖什么、哪些测试和文档要一起补。
- 如果你要按周或按轮推进，再看 [30 轮持续进化计划](30-iteration-evolution-plan.md)。

## 第一优先级开工面

建议先按这个顺序推进：

1. `P0-A`：收紧 `finalize_flow` 默认主路径
2. `P0-B`：增强 observation 与 selector 稳定性
3. `P0-C`：补 repair 闭环的再验证与 diff
4. `P0-D`：补 session 生命周期和健康状态

## 10 条核心功能执行表

| 编号 | 核心功能 | 第一刀先做什么 | 主要依赖 | 归属面 | 必补测试 | 文档联动 |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | 自然语言一键收敛成可运行 Flow | 统一 `finalize_flow` 状态机输出 | provider plan、validate、run | MCP / Workbench | finalize 状态分支测试 | README、AI 入门、MCP 教程 |
| 2 | 页面观察能力升级 | 补表单、表格、分页、结果区抽取 | observe、workbench explore | Runtime / Discovery | observation 结构测试 | 教程 113、产品路线图 |
| 3 | Selector 稳定性与多候选回退 | 给 selector 候选做稳定性排序 | observation、draft、repair | Runtime / AI | selector 排序与回退测试 | Flow 编写规范、review 教程 |
| 4 | 自动修复闭环质量 | repair 后自动 validate，并输出 diff | run trace、repair context | AI / Workbench | repair 回归测试 | AI 入门、repair 教程 |
| 5 | 语义校验器 | 先补变量依赖和 artifact 命名 warning | ValidateFlow、issue hints | Flow Core | validation issue 测试 | Flow 规范、review 清单 |
| 6 | 会话生命周期 | 增加 session 健康状态和验证入口 | save/list/use session | Session / Workbench | session 状态与权限测试 | 会话教程、快速开始 |
| 7 | 断点续跑与幂等 | 标准 ledger 结构和统计摘要 | foreach、read_csv/excel | Flow Core / Data | resume / skip / ledger 测试 | 批量教程、交付模板 |
| 8 | 证据中心 | 每次运行生成统一 manifest | run_root、artifact output | Runtime / Workbench | manifest 生成测试 | handoff 教程、交付文档 |
| 9 | 外部连接可靠性层 | 统一 timeout / retry / 脱敏输出 | HTTP、Redis、DB、SMTP | Data / Infra | connection failure 测试 | 安全边界、外部系统教程 |
| 10 | 单二进制首跑体验 | 固定首跑顺序和入口说明 | list-assets、extract-assets | Release / Docs | release smoke check | 快速开始、142-150 教程 |

## 立即可开的 P0 清单

### P0-A：默认主路径

- [ ] `finalize_flow` 返回统一 `status`
- [ ] `needs_input` 附带缺失字段列表
- [ ] `needs_permission` 附带授权建议
- [ ] `needs_repair` 附带 repair 入口建议
- [ ] Workbench 默认先走 `plan -> finalize -> run`

### P0-B：页面观察与 selector

- [ ] observation 补表单字段与筛选区摘要
- [ ] observation 补结果区、分页器、空状态、提示条
- [ ] selector 候选增加稳定性排序
- [ ] draft 选高稳定性 selector
- [ ] repair 保留次优候选作为回退链

### P0-C：repair 闭环

- [ ] repair 输出修改影响范围
- [ ] repair 后自动 validate
- [ ] 输出修复前后差异摘要
- [ ] 对 selector / wait / var / permission 失败做分类提示

### P0-D：session 生命周期

- [ ] session 增加 `status`
- [ ] session 增加 `last_verified_at`
- [ ] 增加轻量 session verify 入口
- [ ] Workbench 展示 session 健康状态

## 每条功能的完成定义

### 1. 自然语言一键收敛成可运行 Flow

- 用户不需要自己拼工具链顺序
- 同一个意图能明确落在 `ready / needs_input / needs_permission / needs_repair`
- Workbench 上能显示下一步而不是只显示 YAML

### 2. 页面观察能力升级

- observation 能回答“去哪输入、点什么、结果在哪”
- page card 直接复用 observation 结果
- 复杂后台页的摘要明显比现在更业务化

### 3. Selector 稳定性与多候选回退

- 候选 selector 有稳定性顺序
- 页面轻微改版后，repair 能沿候选链恢复
- draft、repair、review 三处 selector 口径一致

### 4. 自动修复闭环质量

- 一次 repair 结果里能看到失败点、修复点、验证结果
- 不需要完整 HTML 也能完成常见修复
- 用户能看懂修复是否真的更接近可运行

### 5. 语义校验器

- 运行前能拦下更多“语法对但业务错”的 Flow
- 校验输出有 `error / warning / suggestion`
- 变量依赖和 artifact 规范能提前暴露

### 6. 会话生命周期

- 用户知道哪个 session 还能用
- 能看出 session 绑定了哪个站点
- 登录态失效和页面 selector 失败能被区分开

### 7. 断点续跑与幂等

- 中断后知道从哪继续
- 成功项不会被重复写
- 失败项有标准 ledger

### 8. 证据中心

- 一次运行生成统一 manifest
- step、artifact、输出文件能被索引
- review 和 handoff 不需要翻目录找材料

### 9. 外部连接可靠性层

- 错误能区分配置、网络、权限、业务
- 日志不泄露敏感头和凭据
- HTTP / Redis / DB / SMTP 的基本行为更一致

### 10. 单二进制首跑体验

- 没看源码的人也知道第一步是什么
- `list-assets -> extract-assets -> getting-started` 成为固定顺序
- release 包和文档入口长期保持同步

## 建议按模块分工

### 方向 A：Flow 与 MCP 主路径

- `finalize_flow`
- `validate_flow`
- `repair_flow_context`
- `repair_flow`

### 方向 B：页面观察与 Knowledge 原材料

- `observe_page`
- `selector_candidates`
- Workbench page/api/entity card

### 方向 C：Session、Artifacts、Workbench

- session 生命周期
- artifact manifest
- Workbench 默认任务流和预览页

### 方向 D：批量与外部系统

- `foreach` resume/ledger
- HTTP / Redis / DB / SMTP 可靠性
- 交付模板与证据包
