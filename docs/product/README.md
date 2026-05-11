# TSPlay 本地智能工作台落地规划

> 目标：沉淀一份可直接复用到 README、PRD、方案汇报和项目讨论中的产品底稿。

## 配套文件

- [README 开头可用版本](README-snippet.md)
- [PRD 摘要版](PRD-summary.md)
- [V1 落地任务清单](V1-task-list.md)
- [Top 10 核心功能路线图](core-feature-roadmap.md)
- [核心功能执行面板](core-feature-execution-board.md)
- [30 轮持续进化计划](30-iteration-evolution-plan.md)

## 一、产品定位一页纸

### 1. 产品名称建议

暂定名称：

`TSPlay Workbench: 已授权 Web 系统认知与数据编排本地工作台`

可选简化版本：

- `登录态 Web 系统认知与数据编排工作台`
- `已授权 Web 系统认知与数据编排的本地 AI 工作台`

更偏宣传口径的一句话：

`让 AI 读懂你的 Web 后台。`

### 2. 一句话定位

面向用户已授权访问的 Web 系统，在本地登录态下自动认知页面功能、发现前端接口、沉淀数据结构，并让 AI 基于这些知识完成数据获取、分析和流程编排。

### 3. 核心问题

大量企业后台、SaaS 系统、数据平台和管理系统，普遍存在三类问题：

- 功能藏在页面里
  菜单、表单、按钮、表格、导出、详情页等能力分散在 Web 页面中，AI 默认并不知道这些能力在哪里。
- 接口没有文档
  很多系统没有 Swagger 或 OpenAPI 文档，但前端页面实际上调用了大量可复用接口。
- 数据难以被 AI 编排
  AI 不知道哪些数据存在、如何筛选、如何导出、哪些操作安全、哪些接口值得直接调用。

因此，单纯让 AI “操作浏览器”是不够的。真正需要的是：

`先认知站点，再发现接口，再沉淀数据能力，最后让 AI 基于这些知识可靠地编排流程。`

### 4. 产品价值

TSPlay Workbench 的核心价值不是“替用户点网页”，而是：

`把用户有权访问的 Web 系统，转化为 AI 可理解、可检索、可调用、可分析、可编排的本地知识资产。`

| 能力 | 价值 |
| --- | --- |
| 站点认知 | 理解页面、菜单、表单、按钮、表格和功能路由 |
| 接口发现 | 发现页面动作背后的 XHR / Fetch / API 调用 |
| 数据建模 | 推断请求参数、响应结构、业务实体和字段含义 |
| AI 编排 | 根据自然语言意图生成 UI / API / 混合执行流程 |
| 本地执行 | 登录态、Cookie、Token 和业务数据尽量留在本地 |
| 数据分析 | 获取数据后进行结构化、清洗、查询、分析和报告生成 |
| Flow 沉淀 | 将一次性操作沉淀为可复用、可审阅、可修复的流程资产 |

### 5. 产品边界

它不是传统 RPA，也不是普通爬虫，也不是 API 管理工具。

更准确的定义是：

`AI-native Web App Intelligence Workbench`

也就是：

`已授权站点认知 + 页面接口发现 + 数据能力建模 + AI 流程编排 + 本地安全执行`

### 6. 差异化

它的差异化不在“把网页点得更像人”，而在：

- 先理解站点，再执行任务
- 把页面、接口、字段、流程持续沉淀成知识资产
- 优先走更稳定的接口路径，再决定是否补 UI 自动化
- 用户用得越久，站点知识越完整，AI 编排越可靠

## 二、MVP 范围

MVP 不建议一开始做“大而全”的自动化平台，而是先做一个完整闭环：

`用户本地启动 TSPlay Workbench，登录一个 Web 后台，系统识别页面功能和接口能力，沉淀知识卡片，然后 AI 根据用户需求生成并执行数据获取 / 分析流程。`

### 1. 整体目的

验证 `TSPlay Workbench` 的核心闭环是否成立：

在用户已授权登录态下，系统能够认知 Web 系统的主要页面与接口能力，并让 AI 基于这些知识完成数据获取、基础分析和可执行 Flow 编排。

第一阶段的目标，不是做完整平台，而是跑通这条最小可用闭环：

`登录站点 -> 页面认知 -> 接口发现 -> AI 编排 -> 数据获取 -> 基础分析`

### 2. MVP 目标

第一版目标：

让用户可以在本地登录一个网站，自动识别该网站的页面、菜单、按钮、表单、表格和接口，生成页面/API 知识卡片，并让 AI 基于这些卡片生成可执行的数据获取流程。

### 3. MVP 成功标准

1. 能够成功登录目标网站，并稳定复用 `session`
2. 能够识别 `80%` 以上主要菜单页面，并形成基础站点地图
3. 能够捕获主要 `XHR / Fetch` 接口，并建立页面动作与接口的关联关系
4. 能够生成可读、可检索的页面卡片和 API 卡片
5. 能够根据自然语言需求生成一条可执行的 `TSPlay Flow`
6. 能够完成至少一次数据获取、结构化落地和基础分析输出

### 4. MVP 典型场景

#### 场景一：订单后台数据分析

用户登录电商后台后，系统识别：

- 订单管理
- 商品管理
- 售后管理
- 数据报表
- 客户管理

用户在订单页面执行一次搜索和导出，系统发现：

- `POST /api/orders/search`
- `POST /api/orders/export`
- `GET /api/orders/{id}`

然后用户说：

`帮我分析上周未发货订单，找出超过 3 天未处理的订单。`

系统执行：

1. 识别意图
2. 匹配订单管理页面和订单搜索接口
3. 生成数据获取 Flow
4. 调用 UI 或接口获取订单数据
5. 结构化为本地数据表
6. 分析异常订单
7. 输出表格和结论

### 5. MVP 功能清单

#### 必须有

| 模块 | MVP 功能 |
| --- | --- |
| 本地启动 | 本地启动 TSPlay Runner / Workbench |
| 登录态管理 | 用户手动登录网站，保存 session |
| 页面探索 | 识别菜单、链接、页面标题、表单、按钮、表格 |
| 路由图谱 | 生成站点页面地图和功能路由 |
| 网络监听 | 捕获 XHR / Fetch 请求 |
| 接口发现 | 记录接口 URL、method、参数、响应结构 |
| 接口关联 | 将接口和页面动作关联 |
| 知识卡片 | 生成页面卡片、接口卡片、数据实体卡片 |
| AI 编排 | 根据自然语言需求匹配页面/API 能力 |
| Flow 生成 | 生成 TSPlay Flow |
| 本地执行 | 调用 TSPlay 执行 Flow |
| 结果落地 | 将接口结果或导出文件落到本地 |
| 简单分析 | 支持表格预览、基础统计、自然语言总结 |

#### 暂不做

第一版先不要做：

- 多租户 SaaS
- 复杂权限体系
- 团队协作
- 云端同步
- 大规模任务调度
- 复杂图数据库
- 浏览器插件
- 全自动无监督探索
- 高风险写操作自动执行
- 支付、审批、删除类自动化

原因不是这些能力不重要，而是 MVP 的关键在于先验证：

`站点认知 + 接口发现 + AI 数据编排` 这条链路是否成立。

### 6. MVP 交付物

第一版建议交付以下内容：

1. 本地 Workbench UI
2. TSPlay Runner 集成
3. 登录态保存能力
4. 页面探索器
5. 网络请求记录器
6. 页面/API 知识卡片生成器
7. 简单知识库
8. AI 意图解析与 Flow 生成
9. Flow 执行与日志回放
10. 数据结果表格与分析摘要

## 三、模块分层

整体架构建议分为五层：

```text
┌────────────────────────────────────┐
│  VSCode / Local Workbench UI       │
│  用户入口、可视化、配置、审阅、调试  │
└────────────────────────────────────┘
                    ↓
┌────────────────────────────────────┐
│  AI Orchestration Layer            │
│  意图识别、知识检索、Flow 生成、修复 │
└────────────────────────────────────┘
                    ↓
┌────────────────────────────────────┐
│  Knowledge Layer                   │
│  页面图谱、接口图谱、实体图谱、卡片库 │
└────────────────────────────────────┘
                    ↓
┌────────────────────────────────────┐
│  Discovery Layer                   │
│  站点探索、DOM 摘要、接口监听、Schema │
└────────────────────────────────────┘
                    ↓
┌────────────────────────────────────┐
│  TSPlay Runtime Layer              │
│  浏览器控制、Session、Flow 执行、修复 │
└────────────────────────────────────┘
```

### 1. TSPlay Runtime Layer

#### 定位

TSPlay 是底层执行引擎，不直接承担完整产品职责。

它负责：

- 浏览器启动
- 页面访问
- 登录态复用
- Flow 执行
- 页面动作
- 数据动作
- 失败截图
- DOM / HTML 快照
- 文件下载
- 基础 repair
- MCP / CLI / HTTP 调用入口

TSPlay 不应该过早承担：

- 完整 UI
- 复杂知识图谱
- PRD 级产品体验
- 多模型调度
- 复杂数据分析
- 团队权限系统

#### 职责边界

TSPlay 的边界是：

`给定一个明确的 Flow 或操作指令，稳定、安全、可观测地执行它。`

也就是说，TSPlay 更像：

`Browser Runtime + Flow Engine + Automation Executor + Failure Artifact Collector`

### 2. Discovery Layer

#### 定位

Discovery Layer 是在 TSPlay 之上新增的“认知采集层”。

它负责从已登录网站中采集结构化信息。

#### 核心模块

- `Site Explorer`
- `DOM Summarizer`
- `Action Discoverer`
- `Network Recorder`
- `API Normalizer`
- `Schema Inferencer`
- `Risk Classifier`

#### 采集内容

每个页面采集：

```json
{
  "url": "/admin/orders",
  "title": "订单管理",
  "menu_path": ["业务", "订单管理"],
  "breadcrumbs": ["后台", "订单管理"],
  "forms": ["订单号", "状态", "时间范围"],
  "buttons": ["搜索", "重置", "导出", "查看详情"],
  "tables": ["订单号", "客户", "金额", "状态", "创建时间"],
  "links": [],
  "risk_actions": []
}
```

每个接口采集：

```json
{
  "method": "POST",
  "path": "/api/orders/search",
  "trigger": "订单管理页面点击搜索",
  "request_schema": {
    "status": "string",
    "startDate": "date",
    "endDate": "date",
    "page": "number"
  },
  "response_schema": {
    "items": "Order[]",
    "total": "number"
  },
  "risk": "read"
}
```

#### 职责边界

Discovery Layer 的边界是：

`把网站中实际存在的页面、动作、接口和数据结构采集出来，并转成结构化知识。`

它不负责最终 AI 推理，也不负责完整执行策略。

### 3. Knowledge Layer

#### 定位

Knowledge Layer 是 AI 可用的站点知识库。

它不保存一堆原始 HTML，而是保存压缩后的、结构化的、可检索的知识资产。

#### 知识类型

- `RouteNode` 页面 / 路由节点
- `FeatureNode` 功能节点
- `ActionNode` 页面动作节点
- `ApiNode` 接口节点
- `EntityNode` 业务实体节点
- `FieldNode` 字段节点
- `FlowTemplate` 流程模板

#### 典型关系

- 页面 `HAS_ACTION` 动作
- 动作 `TRIGGERS_API` 接口
- 接口 `RETURNS` 实体
- 页面 `HAS_TABLE` 表格
- 功能 `IMPLEMENTED_BY` 页面
- 功能 `IMPLEMENTED_BY` 接口
- 用户意图 `MATCHES` 功能

#### 第一版存储建议

MVP 阶段不要急着上复杂图数据库。

建议先用：

- `SQLite`
- `JSON` 字段
- `FTS` 全文检索
- 本地向量索引
- 文件目录保存截图和 artifacts

后续再升级到：

- `PostgreSQL + pgvector`
- `Neo4j`
- `ArangoDB`
- `DuckDB`

#### 职责边界

Knowledge Layer 的边界是：

`保存、索引和检索站点能力知识，为 AI 编排提供 grounded context。`

### 4. AI Orchestration Layer

#### 定位

AI 层负责理解用户意图，并把意图转成可执行计划。

它不应该直接盲目操作页面，而应该先检索 Knowledge Layer。

#### 工作流程

1. 用户输入自然语言
2. 意图识别
3. 检索页面/API/实体/Flow 卡片
4. 选择执行策略
5. 生成执行计划
6. 生成 TSPlay Flow
7. 调用 TSPlay 执行
8. 失败时读取 artifacts 进行修复

#### 执行策略

AI 层需要在三种路径中选择：

| 路径 | 适用场景 |
| --- | --- |
| UI Flow | 接口不稳定、动作依赖页面状态、需要人工确认 |
| API Flow | 接口清晰、只读查询、数据获取效率要求高 |
| 混合 Flow | 先用页面确认状态，再用接口拉取数据 |

#### AI 层职责

- 意图解析
- 知识检索
- 任务拆解
- Flow 生成
- Flow 校验
- 执行策略选择
- 失败修复建议
- 数据分析摘要生成

#### AI 层不应该做

- 直接保存用户 token
- 直接执行高风险写操作
- 绕过用户权限
- 猜测未知接口
- 脱离知识库盲目编排

#### 职责边界

AI 层的边界是：

`基于已采集的站点知识和用户授权，生成可审阅、可执行、可修复的数据与流程编排方案。`

### 5. VSCode / Local Workbench UI

#### 为什么可以考虑 VSCode 工作台

VSCode 适合作为早期工作台，因为它天然适合：

- 本地文件
- YAML / JSON 编辑
- 日志查看
- 插件扩展
- 终端命令
- 开发者用户
- Flow 调试
- 知识卡片查看

如果你的早期用户是：

- 开发者
- 实施工程师
- 数据工程师
- 运维工程师
- 测试工程师
- AI Agent 开发者

那么 VSCode 插件 / 本地工作台是很合适的。

## 四、TSPlay / AI 层 / VSCode 工作台职责边界

这是后续设计最关键的部分。

### 1. TSPlay 的职责

TSPlay 负责“执行”。

#### TSPlay 应该做

- 启动浏览器
- 复用 session
- 执行 Flow
- 执行页面动作
- 执行数据动作
- 保存截图
- 保存 HTML / DOM 快照
- 记录执行日志
- 暴露 CLI / MCP / HTTP 接口
- 支持失败 repair 所需 artifacts

#### TSPlay 不应该做

- 不负责完整产品 UI
- 不负责复杂业务知识图谱
- 不负责多模型策略
- 不负责用户意图理解
- 不负责复杂数据分析
- 不负责团队权限系统
- 不负责长期任务管理

#### TSPlay 的定位

`可嵌入的本地浏览器执行引擎和 Flow Runtime。`

### 2. AI 层的职责

AI 层负责“理解和编排”。

#### AI 层应该做

- 理解用户自然语言需求
- 检索站点知识卡片
- 识别用户真实意图
- 判断需要页面操作还是接口调用
- 生成执行计划
- 生成 TSPlay Flow
- 解释执行结果
- 失败后辅助修复
- 生成数据分析摘要

#### AI 层不应该做

- 不直接接管浏览器底层细节
- 不保存敏感登录态
- 不绕过权限
- 不直接执行高风险动作
- 不臆造不存在的页面和接口
- 不把原始敏感数据无控制地发送到云端模型

#### AI 层的定位

`基于站点知识的意图解析与流程编排层。`

### 3. VSCode / Workbench 的职责

Workbench 负责“用户入口和工程化管理”。

#### Workbench 应该做

- 站点配置
- 登录入口
- Session 管理
- 探索任务配置
- 页面地图查看
- 接口列表查看
- 知识卡片查看
- Flow 编辑
- Flow 调试
- 执行日志查看
- 失败回放
- 数据结果预览
- AI 对话入口
- 权限确认

#### Workbench 不应该做

- 不直接实现浏览器自动化核心
- 不直接处理复杂 AI 推理
- 不承载所有后端逻辑
- 不成为臃肿客户端

#### Workbench 的定位

`本地可视化控制台，用于认知站点、审阅知识、调试 Flow、触发 AI 编排。`

## 五、推荐工程目录

可以按这个结构设计：

```text
tsplay-workbench/
  apps/
    vscode-extension/
    desktop-workbench/
    web-console/

  packages/
    tsplay-adapter/
      client.ts
      flow-runner.ts
      session-manager.ts

    discovery/
      site-explorer.ts
      dom-summarizer.ts
      action-discoverer.ts
      network-recorder.ts
      api-normalizer.ts
      schema-inferencer.ts
      risk-classifier.ts

    knowledge/
      route-store.ts
      api-store.ts
      entity-store.ts
      flow-store.ts
      vector-index.ts
      knowledge-card.ts

    ai/
      model-gateway.ts
      intent-resolver.ts
      retrieval.ts
      flow-generator.ts
      repair-agent.ts
      data-analyst.ts

    data/
      duckdb-adapter.ts
      sqlite-store.ts
      file-loader.ts
      table-preview.ts

    security/
      redactor.ts
      permission-policy.ts
      action-risk.ts
      audit-log.ts

  runtime/
    tsplay/
    artifacts/
    sessions/
    sites/
    flows/
    data/
```

如果直接基于 TSPlay 仓库扩展，也可以先简单一些：

```text
tsplay/
  cmd/
  flow/
  mcp/
  browser/
  collector/
    site_explorer.go
    dom_summarizer.go
    network_recorder.go
  knowledge/
    graph_store.go
    route_card.go
    api_card.go
    entity_card.go
  planner/
    intent_resolver.go
    flow_generator.go
  workbench/
    api_server.go
    ui/
```

## 六、关键数据模型

### 1. Site

```json
{
  "site_id": "demo_admin",
  "name": "Demo Admin",
  "start_url": "https://example.com/admin",
  "allowed_domains": ["example.com"],
  "session_name": "demo_admin_user",
  "created_at": "2026-04-25T10:00:00+09:00"
}
```

### 2. Route Card

```json
{
  "id": "route:demo_admin:/orders",
  "site_id": "demo_admin",
  "url": "https://example.com/admin/orders",
  "normalized_route": "/orders",
  "title": "订单管理",
  "menu_path": ["业务管理", "订单管理"],
  "summary": "用于订单查询、筛选、导出和详情查看",
  "forms": [
    {
      "name": "订单筛选表单",
      "fields": ["订单号", "订单状态", "开始时间", "结束时间"]
    }
  ],
  "tables": [
    {
      "name": "订单列表",
      "columns": ["订单号", "客户", "金额", "状态", "创建时间"]
    }
  ],
  "actions": ["搜索", "重置", "导出", "查看详情"],
  "risk": "low"
}
```

### 3. API Card

```json
{
  "id": "api:POST:/api/orders/search",
  "site_id": "demo_admin",
  "method": "POST",
  "path_template": "/api/orders/search",
  "semantic_name": "订单搜索",
  "trigger_route": "/orders",
  "trigger_action": "点击搜索按钮",
  "operation_type": "read",
  "request_schema": {
    "status": "string",
    "startDate": "date",
    "endDate": "date",
    "page": "number",
    "pageSize": "number"
  },
  "response_schema": {
    "items": "Order[]",
    "total": "number"
  },
  "risk": "read"
}
```

### 4. Entity Card

```json
{
  "id": "entity:Order",
  "site_id": "demo_admin",
  "name": "Order",
  "label": "订单",
  "fields": [
    {
      "name": "orderId",
      "label": "订单号",
      "type": "string"
    },
    {
      "name": "status",
      "label": "订单状态",
      "type": "string"
    },
    {
      "name": "amount",
      "label": "金额",
      "type": "number"
    },
    {
      "name": "createdAt",
      "label": "创建时间",
      "type": "datetime"
    }
  ]
}
```

### 5. Flow Plan

```json
{
  "intent": "分析上周未发货订单",
  "matched_features": ["订单管理", "订单搜索", "订单导出"],
  "matched_apis": ["POST /api/orders/search"],
  "strategy": "api_first",
  "reason": "订单搜索接口为只读接口，参数明确，适合直接获取数据",
  "flow_name": "analyze_unshipped_orders_last_week",
  "requires_user_confirm": false
}
```

## 七、权限与安全边界

这个产品必须从第一版就设计安全边界。

### 1. 授权边界

系统只处理：

- 用户自己登录
- 用户有权访问
- 用户主动配置
- 用户明确授权

不做：

- 绕过登录
- 破解接口
- 扫描未授权系统
- 绕过验证码
- 规避风控
- 批量攻击
- 越权访问

### 2. 数据边界

默认本地保存：

- Session
- Cookie
- Token
- 页面截图
- HTML 快照
- 接口日志
- 业务数据
- 分析结果

默认发送给 AI 的内容应该是：

- 脱敏后的页面摘要
- 字段名
- 接口 schema
- 功能卡片
- 风险标签
- 少量脱敏样例

不应直接发送：

- Cookie
- `Authorization`
- 完整 Token
- 身份证号
- 手机号
- 客户姓名
- 订单明细原文
- 大量业务数据

### 3. 操作风险分级

| 风险级别 | 操作类型 | 默认策略 |
| --- | --- | --- |
| `read` | 搜索、查看、分页、筛选 | 允许自动执行 |
| `read_download` | 导出、下载 | 需要用户授权 |
| `write_low` | 保存草稿、修改备注 | 需要确认 |
| `write_high` | 删除、审批、发布、禁用 | 默认禁止自动执行 |
| `critical` | 支付、转账、修改密码 | 禁止自动执行 |

## 八、推荐开发路线

### 阶段 1：本地 Runner + Session

目标：

- TSPlay 能在本地启动
- 用户能手动登录站点
- 系统能保存和复用 session

交付：

- 站点配置
- 登录按钮
- session 保存
- session 测试

### 阶段 2：页面认知

目标：

- 能够探索站点菜单和页面
- 生成页面功能卡片

交付：

- 站点地图
- 页面列表
- 菜单路径
- 表单、按钮、表格识别
- 截图归档

### 阶段 3：接口发现

目标：

- 能监听页面操作产生的 XHR / Fetch 请求
- 推断接口语义和 schema

交付：

- 接口列表
- 接口详情
- 触发动作
- 请求 / 响应 schema
- 风险等级

### 阶段 4：知识检索 + AI 编排

目标：

- 用户自然语言输入后，AI 能匹配相关页面/API，并生成 Flow

交付：

- `Intent Resolver`
- `Knowledge Retrieval`
- `Flow Generator`
- `Flow Preview`
- `Flow Execute`

### 阶段 5：数据分析

目标：

- 获取到的数据可以落表、查询、分析和导出

交付：

- 本地数据表
- `CSV / Excel / JSON` 处理
- `DuckDB / SQLite` 查询
- AI 分析摘要
- 报告导出

## 九、README 可用版本

下面这段可以直接放到 README 开头。

```md
# TSPlay Workbench

TSPlay Workbench is a local AI workbench for authorized web systems.

It helps users understand logged-in web applications, discover frontend APIs, build page/API/data knowledge cards, and let AI agents safely orchestrate data acquisition, analysis, and browser workflows based on real site knowledge.

Unlike traditional RPA or browser agents, TSPlay Workbench does not simply click pages. It first maps the web application, discovers the APIs behind user actions, infers data schemas, and builds a local knowledge layer that AI can use for grounded planning.

## Core Capabilities

- Authorized site cognition
- Logged-in session reuse
- Page, menu, form, table, and action discovery
- XHR / Fetch API discovery
- Request / response schema inference
- Page-to-API relationship mapping
- Local knowledge cards for AI retrieval
- UI / API / hybrid workflow generation
- TSPlay Flow execution
- Local data acquisition and analysis
- Human approval for risky operations

## Positioning

TSPlay Workbench turns web systems you are authorized to access into AI-readable, callable, analyzable, and orchestratable local knowledge assets.
```

## 十、PRD 摘要版

### 产品目标

构建一个本地 AI 工作台，使用户能够在已授权登录态下认知 Web 系统、发现接口、沉淀数据能力，并通过 AI 完成数据获取、分析和流程编排。

### 目标用户

- 实施工程师
- 数据分析师
- 测试工程师
- 运维工程师
- AI Agent 开发者
- 企业内部工具团队
- 经常操作 Web 后台的业务人员

### 核心用户故事

作为一个实施工程师，我希望登录客户的 Web 后台后，系统能自动识别页面功能和接口能力，这样我就可以快速为客户生成自动化数据获取和报表流程。

作为一个数据分析师，我希望不用理解后台复杂页面，只需要说“帮我分析上周订单异常”，系统就能知道从哪个页面或接口取数据，并生成分析结果。

作为一个 AI Agent 开发者，我希望把一个登录态网站转化成 AI 可检索的知识库，这样 Agent 在执行任务时不再盲目点击页面。

### MVP 成功指标

- 能够成功登录并复用 session
- 能够识别 80% 以上主要菜单页面
- 能够捕获主要 XHR / Fetch 接口
- 能够生成可读的页面/API 卡片
- 能够根据自然语言生成一个可执行 Flow
- 能够完成一次数据获取和基础分析
