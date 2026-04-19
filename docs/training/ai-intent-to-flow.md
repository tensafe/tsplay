# TSPlay AI 无感入门：从用户意图到 Flow（以 Codex 为例）

这篇教程服务于这样一类新手：

- 不想先学完整 Flow 语法
- 不想自己找 selector
- 更习惯直接告诉 AI“我想做什么”
- 希望把“所思即所想，所想即实现”尽量落成稳定的自动化流程

核心思路不是让用户手写 Flow，而是让 AI 通过 TSPlay MCP 工具链把用户意图逐步收敛成可运行、可校验、可修复的 Flow。

## 这条路线解决什么问题

传统入门路径通常是：

1. 学 CLI 动作
2. 学 selector
3. 学 Flow 结构
4. 学 MCP 工具

对新手来说，这条路径没有错，但门槛偏高。

如果目标是“先让用户把事做成”，更适合先走 AI 路线：

1. 用户只描述业务意图
2. Codex 调 `draft_flow` 或 `observe_page`
3. AI 自动做一轮 `validate_flow`
4. AI 帮用户执行 `run_flow`
5. 如果失败，再进入 `repair_flow_context` / `repair_flow`

这样用户先获得结果，再逐步理解 Flow。

## 新手最少只需要提供什么

理想情况下，用户只需要告诉 Codex 4 件事：

- 页面在哪：URL
- 想做什么：例如“搜索订单并导出”“上传文件并提交”
- 输入是什么：关键词、文件、筛选条件、账号角色
- 是否可信可授权：例如是否允许文件访问、登录态复用

如果这 4 件事说清楚，用户通常不需要自己提供：

- selector
- Flow 骨架
- 校验命令
- 失败修复策略

## 准备条件

开始前建议确认这几项：

- 已能启动 TSPlay MCP Server：`go run . -action srv -flow-root script -artifact-root artifacts`
- Codex 已连接到 TSPlay MCP 工具
- 本地或内网页面可访问
- 如果要练习仓库内 demo，确保能访问 `demo/` 下页面

建议新手优先使用这些稳定素材：

- [../../demo/demo.html](../../demo/demo.html)
- [../../demo/tables.html](../../demo/tables.html)
- [../../demo/upload.html](../../demo/upload.html)
- [../../demo/multi_upfile.html](../../demo/multi_upfile.html)

如果需要先补项目背景，可先看 [../../ReadMe.md](../../ReadMe.md) 的 MCP 章节。

## Codex 的标准工作流

当用户只说意图时，Codex 最好按这条顺序工作：

1. 先用 `tsplay.flow_schema` 和 `tsplay.flow_examples` 建立约束，不靠猜
2. 如果用户给了明确 URL，优先尝试 `tsplay.draft_flow`
3. 如果页面复杂、元素不明显或草稿不稳，再补 `tsplay.observe_page`
4. 查看 `draft_flow` 返回的 `validation`、`unresolved`、`warnings`、`repair_hints`
5. 需要单独确认结构时，再显式调 `tsplay.validate_flow`
6. 校验通过后再调 `tsplay.run_flow`
7. 执行失败时，用 `tsplay.repair_flow_context` / `tsplay.repair_flow` 收敛修复
8. 如果流程依赖登录态，再用 `tsplay.save_session` 沉淀会话

这套顺序的好处是：

- 用户不必自己找 selector
- `draft_flow` 已经会自动做一轮校验和 selector 修正
- 修复时不会把整页 HTML 原样塞给模型
- 最终产物仍然是可审阅的 Flow，而不是一次性对话产物

## 给 Codex 的推荐提示词

下面这段话可以作为 Codex 的工作指令，帮助它更稳定地把用户意图收敛成 Flow：

```text
你现在是 TSPlay Flow 助手，目标是让用户只通过自然语言描述任务，而不是手写 selector 或 Flow。

工作原则：
1. 先使用 tsplay.flow_schema 和 tsplay.flow_examples 获取约束。
2. 用户给了 URL 且目标清晰时，优先使用 tsplay.draft_flow。
3. 当页面复杂、selector 不确定、或 draft 结果有 unresolved/warnings 时，再使用 tsplay.observe_page。
4. 始终检查 validation、repair_hints、unresolved，不要把草稿直接当成最终答案。
5. 成功前按 validate -> run -> repair 的顺序推进。
6. 只有在用户场景明确需要时，才申请 allow_file_access、allow_browser_state、allow_http、allow_database 等高风险授权。
7. 如果场景依赖登录态，优先建议使用 tsplay.save_session，并在 Flow 顶层使用 browser.use_session。
8. 对用户输出时优先讲“你现在可以做什么”和“还缺什么输入”，不要要求用户理解底层 selector 细节。
```

## 新手操作模板

新手和 Codex 对话时，建议使用这种表达方式：

```text
帮我在 <URL> 上完成下面的任务：
- 目标：<我想做什么>
- 输入：<关键词 / 文件 / 条件>
- 结果要求：<我希望拿到什么结果>
- 授权说明：<是否允许文件访问、登录态、HTTP、数据库>
```

例如：

```text
帮我在 http://127.0.0.1:8000/demo/tables.html 上提取表格数据。
- 目标：抓取表头和所有行
- 输入：无
- 结果要求：给我一条可运行的 Flow，并执行一次
- 授权说明：只允许普通页面读取，不允许文件写入
```

## 三个最适合新手的演示场景

### 场景 1：选择下拉项并验证

用户只需要说：

```text
帮我在 demo 页面里选择“选项 5”，并确认它已经被选中。
```

Codex 的理想动作：

1. 用 `tsplay.draft_flow` 根据页面和意图生成草稿
2. 检查是否出现 `select_option` 与 `is_selected`
3. 如有需要，再用 `tsplay.observe_page` 修正 selector
4. 通过后执行 `tsplay.run_flow`

这个场景适合让用户建立一个直觉：
用户描述的是“业务目标”，不是“动作细节”。

### 场景 2：提取表格

用户只需要说：

```text
帮我把页面里的表格提取成结构化结果。
```

Codex 的理想动作：

1. 观察页面或直接草拟 Flow
2. 优先选择 `capture_table` 这类结构化动作
3. 生成带 `save_as` 的 Flow
4. 执行一次并把结果摘要解释给用户

这个场景适合让新手理解：
AI 不是模拟“复制网页源码”，而是在挑更适合自动化的动作。

### 场景 3：上传文件并提交

用户只需要说：

```text
帮我上传这个文件并提交。
```

Codex 这时要额外做一件事：

- 主动说明该 Flow 需要 `allow_file_access=true`

推荐动作顺序：

1. `tsplay.draft_flow`
2. 看 `repair_hints` 是否已经提示文件授权
3. 明确授权后再 `validate_flow` / `run_flow`

这个场景适合教新手理解：
“无感”不等于“无边界”，高风险能力仍然要显式授权。

## 授权原则

为了让体验足够顺滑，又不失控，建议按最小授权原则处理：

| 授权 | 什么时候开 |
| --- | --- |
| `allow_file_access` | 上传、下载、截图、保存 HTML、读写 CSV/Excel 时 |
| `allow_browser_state` | 读写 Cookie、Storage、保存登录态时 |
| `allow_http` | Flow 里要主动请求外部 API 时 |
| `allow_database` | 有 `db_*` 写入动作时 |

对新手最重要的一句话是：

不要一开始就把所有授权全开，而是让 Codex 根据 Flow 里的动作按需申请。

## 失败时怎么保持“无感”

无感体验不是“永不失败”，而是失败时用户也不用自己排底层细节。

推荐让 Codex 按这个顺序处理：

1. 先向用户解释失败发生在哪一步
2. 再调用 `tsplay.repair_flow_context`
3. 根据 `repair_hints` 调 `tsplay.repair_flow`
4. 修复后重新 `validate_flow`
5. 必要时再重新 `run_flow`

对用户的输出重点应是：

- 哪一步失败了
- AI 准备如何修
- 是否需要用户补充输入或授权

而不是直接把一大段 trace 或 HTML 扔给用户。

## 让体验更像“所想即实现”的落地建议

如果你希望这套教程继续往产品化方向走，建议优先做这几件事：

- 固定一条默认提示词，让 Codex 永远优先走 `draft -> validate -> run -> repair`
- 把常见页面先沉淀成稳定的 demo 或业务模板
- 把登录场景沉淀成 `save_session + use_session`
- 对用户界面只暴露“意图、输入、结果、授权确认”四类信息
- 保留最终 Flow 作为审阅资产，而不是只保留对话记录

这样做的价值是：

- 新手能先把事情做成
- 交付团队仍然有可审查的 Flow 资产
- AI 输出不是黑盒，而是能沉淀、复用、修复的自动化流程

## 推荐阅读顺序

如果你要按这条 AI 路线学习，建议顺序如下：

1. [../../ReadMe.md](../../ReadMe.md) 的 MCP 章节
2. 本文
3. [labs.md](labs.md) 的 Lab 6
4. [capstone-briefs.md](capstone-briefs.md) 中的 MCP 场景

如果之后需要补底层理解，再回看：

- [learning-path.md](learning-path.md)
- [trainer-playbook.md](trainer-playbook.md)
