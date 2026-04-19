# 160 次递进迭代路线图

这份路线图的目的，不是追求“看起来很多”，而是把教程建设拆成 160 个不跳跃、可执行、可验证、可持续增长的迭代点。

使用方式：

- 学员：按阶段顺序推进，不要跨阶段乱跳
- 课程作者：每次至少落实 1 到 3 个迭代点
- 讲师 / 实施 / 负责人：拿它做课程排期和缺口检查

补充说明：

- 当前仓库里已经直接落地、可运行的是 `Lesson 01-71`
- 这份路线图负责把完整课程继续往前推进，避免后续文档生长失去顺序

总结构：

- 新手教程：`001-040`
- 初级教程：`041-080`
- 中级教程：`081-120`
- 高级教程：`121-160`

## 新手教程

## 模块 01：环境、二进制、第一条成功路径

- `[001]` 构建 `./tsplay` 并执行 `./tsplay -action list-assets`。交付物：确认二进制内置 `docs/`、`script/`、`demo/`。
- `[002]` 执行 `./tsplay -action extract-assets -extract-root ./tsplay-assets`。交付物：能在本地看到释放出的参考资料目录。
- `[003]` 阅读 `ReadMe.md` 和 `docs/README.md`，理解仓库入口与文档入口。交付物：一段 5 句话以内的口头或文字总结。
- `[004]` 跑通 `Lesson 01` 的 Lua 版本。交付物：`artifacts/tutorials/01-hello-world-lua.json`。
- `[005]` 跑通 `Lesson 01` 的 Flow 版本。交付物：`artifacts/tutorials/01-hello-world-flow.json`。
- `[006]` 对比 `Lua` 和 `Flow` 在 Hello World 场景的写法差异。交付物：一段对照说明。
- `[007]` 理解 `set_var` 和 `write_json` 为什么不依赖浏览器。交付物：自己改一个字段并重新运行。
- `[008]` 理解 `artifacts/tutorials/` 为什么是默认输出位置。交付物：说明“输入、输出、产物”三者关系。
- `[009]` 能独立解释 `./tsplay -script ...` 和 `./tsplay -flow ...` 的区别。交付物：一段简短说明。
- `[010]` 能独立从空目录调用二进制内置示例。交付物：在非仓库目录运行内置 `script/tutorials/01_hello_world.lua`。

## 模块 02：本地页面与最小交互

- `[011]` 启动 `./tsplay -action file-srv -addr :8000`。交付物：能访问 `http://127.0.0.1:8000/demo/demo.html`。
- `[012]` 阅读 `demo/demo.html`，先用肉眼理解页面结构。交付物：指出下拉框的关键元素。
- `[013]` 跑通 `Lesson 02` 的 Lua 版本。交付物：`artifacts/tutorials/02-select-option-lua.json`。
- `[014]` 跑通 `Lesson 02` 的 Flow 版本。交付物：`artifacts/tutorials/02-select-option-flow.json`。
- `[015]` 理解 `navigate` 与 `wait_for_selector` 的先后关系。交付物：复述为什么不建议上来就 `sleep`。
- `[016]` 理解 `select_option` 和 `is_selected` 的区别。交付物：解释“动作”和“验证”的边界。
- `[017]` 把 `选项 5` 改成另一个值再跑一遍。交付物：新的 JSON 输出。
- `[018]` 学会通过环境变量覆盖 demo URL。交付物：使用 `TSPLAY_DEMO_URL` 运行脚本。
- `[019]` 能解释为什么教程先用本地 demo，不先访问公网网站。交付物：一段复盘。
- `[020]` 能把本地页面交互结果重新写回 JSON。交付物：自定义输出字段。

## 模块 03：读取页面内容

- `[021]` 跑通 `Lesson 03` 的 Lua 版本。交付物：`03-capture-table-lua.json`。
- `[022]` 跑通 `Lesson 03` 的 Flow 版本。交付物：`03-capture-table-flow.json`。
- `[023]` 理解 `capture_table` 为什么优先于 `get_html`。交付物：写一句“什么时候该用哪个动作”的判断标准。
- `[024]` 跑通 `Lesson 04` 的 Lua 版本。交付物：`04-extract-text-and-html-lua.json`。
- `[025]` 跑通 `Lesson 04` 的 Flow 版本。交付物：`04-extract-text-and-html-flow.json`。
- `[026]` 理解 `extract_text` 与 `get_html` 的用途差异。交付物：说明“拿文本”和“拿 DOM 片段”的边界。
- `[027]` 修改 `extract.html` 中的标题文本并重新运行。交付物：新的提取结果。
- `[028]` 修改 `extract.html` 中的计数值并重新运行。交付物：新的 `order_count`。
- `[029]` 学会从同一页面同时提取文本、数字和 HTML 片段。交付物：一份包含三类值的 JSON。
- `[030]` 能解释为什么“先提取页面事实，再进入分支逻辑”更稳。交付物：一段复盘。

## 模块 04：本地 JSON、结果整理、学习闭环

- `[031]` 跑通 `Lesson 05` 的 Lua 版本。交付物：`05-http-request-json-lua.json`。
- `[032]` 跑通 `Lesson 05` 的 Flow 版本。交付物：`05-http-request-json-flow.json`。
- `[033]` 理解 `http_request` 返回对象中 `status`、`headers`、`body` 的结构。交付物：一段说明。
- `[034]` 理解为什么 `json_extract` 路径会写成 `$.body.summary.open`。交付物：解释路径来源。
- `[035]` 修改 `demo/data/order_summary.json` 中一条字段并重新运行。交付物：变化后的 JSON 输出。
- `[036]` 把一个新字段加入本地 JSON，再补一条 `json_extract`。交付物：脚本或 Flow 的增量改动。
- `[037]` 把 `Lesson 02-05` 的结果统一放进一个自己的学习目录。交付物：自定义命名的结果文件集合。
- `[038]` 复盘新手阶段学过的动作。交付物：列出自己已经真正跑过的 action 清单。
- `[039]` 写出“Lua 版更像什么，Flow 版更像什么”的个人理解。交付物：一段 100 字以内总结。
- `[040]` 通过新手阶段检查。交付物：能从头独立跑通 `Lesson 01-05`。

## 初级教程

## 模块 05：文件读写、CSV、输入输出

- `[041]` 跑通 `Lesson 13` 的 Lua 版本。交付物：`13-read-csv-basics-lua.json`。
- `[042]` 跑通 `Lesson 13` 的 Flow 版本。交付物：`13-read-csv-basics-flow.json`。
- `[043]` 理解 CSV 表头、数据行和 `row_number_field` 的关系。交付物：一段说明。
- `[044]` 跑通 `Lesson 14` 的 Lua 版本。交付物：`14-write-csv-basics-lua.csv` 和 `14-write-csv-basics-lua.json`。
- `[045]` 跑通 `Lesson 14` 的 Flow 版本。交付物：`14-write-csv-basics-flow.csv` 和 `14-write-csv-basics-flow.json`。
- `[046]` 理解为什么 `write_csv` 里要显式写 `headers`。交付物：解释列顺序的意义。
- `[047]` 跑通 `Lesson 15` 的 Lua 版本。交付物：`15-transformed-contacts-lua.csv`。
- `[048]` 跑通 `Lesson 15` 的 Flow 版本。交付物：`15-transformed-contacts-flow.csv`。
- `[049]` 理解 `start_row`、`limit`、`row_number_field` 这三个字段的配合方式。交付物：一段批处理说明。
- `[050]` 复盘文件动作的最小心智模型。交付物：输入、处理、输出三段式总结。

## 模块 06：断言、控制流、上传下载

- `[051]` 跑通 `Lesson 16` 的 Lua 版本。交付物：`16-retry-flaky-action-lua.json`。
- `[052]` 跑通 `Lesson 16` 的 Flow 版本。交付物：`16-retry-flaky-action-flow.json`。
- `[053]` 对比“Lua 显式循环”和 `Flow retry`。交付物：一段控制流分工说明。
- `[054]` 跑通 `Lesson 17` 的 Lua 版本。交付物：`17-wait-until-ready-lua.json`。
- `[055]` 跑通 `Lesson 17` 的 Flow 版本。交付物：`17-wait-until-ready-flow.json`。
- `[056]` 理解 `wait_until` 和 `sleep` 的边界。交付物：一段轮询说明。
- `[057]` 跑通 `Lesson 18` 的 Lua 版本。交付物：`18-upload-single-file-lua.json`。
- `[058]` 跑通 `Lesson 18` 的 Flow 版本。交付物：`18-upload-single-file-flow.json`。
- `[059]` 跑通 `Lesson 19` 的 Lua 版本。交付物：`19-upload-multiple-files-lua.json`。
- `[060]` 跑通 `Lesson 19` 的 Flow 版本。交付物：`19-upload-multiple-files-flow.json`。

## 模块 07：下载闭环与外部系统基础接入

- `[061]` 跑通 `Lesson 20` 的 Lua 版本。交付物：`20-download-report-lua.json` 和下载下来的 CSV。
- `[062]` 跑通 `Lesson 20` 的 Flow 版本。交付物：`20-download-report-flow.json` 和下载下来的 CSV。
- `[063]` 理解 `download_file` 和 `download_url` 的区别。交付物：一段说明。
- `[064]` 理解为什么文件类教程经常要先 `extract-assets`。交付物：一段说明。
- `[065]` 跑通 `Lesson 06` 的 Lua 版本。交付物：`06-redis-round-trip-lua.json`。
- `[066]` 跑通 `Lesson 06` 的 Flow 版本。交付物：`06-redis-round-trip-flow.json`。
- `[067]` 理解 `redis_set`、`redis_get`、`redis_incr`、`redis_del` 的最小用法。交付物：一段动作说明。
- `[068]` 跑通 `Lesson 07` 的 Lua 版本。交付物：`07-db-postgres-basics-lua.json`。
- `[069]` 跑通 `Lesson 07` 的 Flow 版本。交付物：`07-db-postgres-basics-flow.json`。
- `[070]` 复盘“文件动作 + 外部系统动作”的最小边界。交付物：CSV、下载、Redis、DB 四者的使用场景对照。

## 模块 08：会话、产物、复盘机制

- `[071]` 结合 `Lesson 28-30`、`Lesson 36-71` 学会根据输出 JSON 回看浏览器状态结果。交付物：指出一条状态快照里的关键字段。
- `[072]` 学会从 `artifact_root` 和 `run_root` 理解一次运行的上下文。交付物：一段说明。
- `[073]` 结合 `Lesson 31-35` 为教程补“失败时先看哪里”的统一说明。交付物：一段通用排障话术。
- `[074]` 学会记录“我改了什么，结果怎么变了”。交付物：一条结构化复盘记录。
- `[075]` 用一个自己的小页面或小 JSON 替换仓库 demo。交付物：一个自定义练习资源。
- `[076]` 把一条 Lua 脚本改写成对应 Flow。交付物：一对对照示例。
- `[077]` 把一条 Flow 改得更可读。交付物：命名更清晰的版本。
- `[078]` 学会写“本阶段退出标准”式复盘。交付物：一页自评清单。
- `[079]` 整理初级阶段主题缺口。交付物：下一步想补的 5 个 lesson 主题。
- `[080]` 通过初级阶段检查。交付物：能独立组织一个包含变量和外部系统基础的 5 到 10 步小流程。

## 中级教程

## 模块 09：模板化与可复用结构

- `[081]` 识别哪些脚本只是一次性试验，哪些值得沉淀成模板。交付物：一个分类清单。
- `[082]` 为一个已有 lesson 写出“模板版目标”。交付物：主题重构说明。
- `[083]` 统一 Flow 顶层结构。交付物：一套推荐字段顺序。
- `[084]` 统一步骤命名规则。交付物：一份步骤命名规范。
- `[085]` 统一变量命名规则。交付物：一份变量命名规范。
- `[086]` 统一输出文件命名规则。交付物：一份文件命名规范。
- `[087]` 为“输入、处理、输出”三段式流程设计模板。交付物：模板草案。
- `[088]` 为“采集、断言、保存”三段式流程设计模板。交付物：模板草案。
- `[089]` 为“请求、解析、落盘”三段式流程设计模板。交付物：模板草案。
- `[090]` 复盘为什么模板化比“示例越多越好”更重要。交付物：一段说明。

## 模块 10：数据驱动与批量处理

- `[091]` 规划 `read_csv` 入门 lesson。交付物：主题说明。
- `[092]` 规划 `write_csv` 入门 lesson。交付物：主题说明。
- `[093]` 规划 `read_excel` 入门 lesson。交付物：主题说明。
- `[094]` 规划 `foreach` 基础 lesson。交付物：主题说明。
- `[095]` 把“单条数据处理”改造成“列表处理”的流程草案。交付物：一条 Flow 草稿。
- `[096]` 设计“输入 CSV，输出 JSON 报告”的最小课题。交付物：课题说明。
- `[097]` 设计“输入 Excel，遍历上传”的最小课题。交付物：课题说明。
- `[098]` 设计“批量 HTTP 请求后汇总结果”的最小课题。交付物：课题说明。
- `[099]` 设计“批量数据库写入”的最小课题。交付物：课题说明。
- `[100]` 复盘批量处理为什么需要比单步示例更强的结构感。交付物：一段说明。

## 模块 11：健壮性、等待、恢复

- `[101]` 系统补 `assert_visible` lesson。交付物：lesson 草稿。
- `[102]` 系统补 `assert_text` lesson。交付物：lesson 草稿。
- `[103]` 系统补 `retry` lesson。交付物：lesson 草稿。
- `[104]` 系统补 `wait_until` lesson。交付物：lesson 草稿。
- `[105]` 系统补 `on_error` lesson。交付物：lesson 草稿。
- `[106]` 设计“延迟出现的元素”场景。交付物：一页 demo / 需求说明。
- `[107]` 设计“偶发失败点击”场景。交付物：一页 demo / 需求说明。
- `[108]` 设计“失败后 reload 再重试”场景。交付物：一页 Flow 草稿。
- `[109]` 整理 artifact、截图、HTML、DOM snapshot 的教学顺序。交付物：一页讲解提纲。
- `[110]` 复盘“健壮性设计不是多写几个 sleep”。交付物：一段反例说明。

## 模块 12：MCP 基础链路与 repair 入门

- `[111]` 解释 `tsplay.list_actions` 的作用。交付物：一段工具定位说明。
- `[112]` 解释 `tsplay.flow_schema` 和 `tsplay.flow_examples` 的作用。交付物：一段说明。
- `[113]` 设计 `observe_page` 入门 lesson。交付物：lesson 草稿。
- `[114]` 设计 `draft_flow` 入门 lesson。交付物：lesson 草稿。
- `[115]` 设计 `validate_flow` 入门 lesson。交付物：lesson 草稿。
- `[116]` 设计 `run_flow` 入门 lesson。交付物：lesson 草稿。
- `[117]` 设计 `repair_flow_context` 入门 lesson。交付物：lesson 草稿。
- `[118]` 设计 `repair_flow` 入门 lesson。交付物：lesson 草稿。
- `[119]` 写出“observe -> draft -> validate -> run -> repair”的一页总图。交付物：流程总览。
- `[120]` 通过中级阶段检查。交付物：能解释 TSPlay 从脚本、Flow 到 MCP 的主线关系。

## 高级教程

## 模块 13：安全边界与运行边界

- `[121]` 为 `allow_lua` 设计一条专门 lesson。交付物：lesson 草稿。
- `[122]` 为 `allow_http` 设计一条专门 lesson。交付物：lesson 草稿。
- `[123]` 为 `allow_file_access` 设计一条专门 lesson。交付物：lesson 草稿。
- `[124]` 为 `allow_browser_state` 设计一条专门 lesson。交付物：lesson 草稿。
- `[125]` 为 `allow_redis` 设计一条专门 lesson。交付物：lesson 草稿。
- `[126]` 为 `allow_database` 设计一条专门 lesson。交付物：lesson 草稿。
- `[127]` 写出“本地 CLI / 本地 Flow / MCP 模式”三者的边界对照。交付物：一页说明。
- `[128]` 写出“为什么教程不能一开始就忽略授权边界”。交付物：一段说明。
- `[129]` 为每类高风险动作整理默认教学前置。交付物：一份前置条件清单。
- `[130]` 复盘安全边界在交付中的价值。交付物：一段交付视角说明。

## 模块 14：大型 Flow、规范、评审

- `[131]` 制定步骤命名评审规则。交付物：规则草案。
- `[132]` 制定变量命名评审规则。交付物：规则草案。
- `[133]` 制定 artifact 目录评审规则。交付物：规则草案。
- `[134]` 制定教程示例代码评审规则。交付物：规则草案。
- `[135]` 制定“什么时候允许 Lua escape hatch”的评审规则。交付物：规则草案。
- `[136]` 制定“什么时候必须抽成 Flow”的评审规则。交付物：规则草案。
- `[137]` 制定“什么时候应该新增 demo 页面”的评审规则。交付物：规则草案。
- `[138]` 制定“示例必须带什么交付物”的评审规则。交付物：规则草案。
- `[139]` 设计一个大型 Flow 示例目录结构。交付物：目录草图。
- `[140]` 复盘教程为什么也需要 code review 思维。交付物：一段说明。

## 模块 15：发布包、内置资产、交付体验

- `[141]` 解释为什么把 `docs/`、`script/`、`demo/` 一起打包到二进制里。交付物：一段说明。
- `[142]` 为 `list-assets` 写一条面向新人的使用说明。交付物：一页帮助文案。
- `[143]` 为 `extract-assets` 写一条面向新人的使用说明。交付物：一页帮助文案。
- `[144]` 设计“只发一个二进制给用户”的交付流程。交付物：交付说明。
- `[145]` 设计“离线环境也能跑基础教程”的交付流程。交付物：交付说明。
- `[146]` 设计“发布包里哪些资源必须带，哪些可以不带”的规则。交付物：清单。
- `[147]` 为 `file-srv` 写一条“开发态”和“发布态”的对照说明。交付物：一页说明。
- `[148]` 设计“新用户第一次打开二进制时怎么找到教程”的入口策略。交付物：入口草案。
- `[149]` 设计“课程包版本号与资源版本号”的维护策略。交付物：版本策略草案。
- `[150]` 复盘“单二进制 + 内置教程”对交付的意义。交付物：一段总结。

## 模块 16：Capstone、培训、持续演化

- `[151]` 设计一个新手结业题。交付物：题目说明。
- `[152]` 设计一个初级结业题。交付物：题目说明。
- `[153]` 设计一个中级结业题。交付物：题目说明。
- `[154]` 设计一个高级结业题。交付物：题目说明。
- `[155]` 设计一套“新人 7 天计划”。交付物：7 天排期。
- `[156]` 设计一套“实施同学 2 周计划”。交付物：2 周排期。
- `[157]` 设计一套“讲师备课顺序”。交付物：备课提纲。
- `[158]` 设计一套“教程缺口复盘机制”。交付物：月度复盘模板。
- `[159]` 设计一套“每 10 次迭代回看一次”的检查机制。交付物：检查清单。
- `[160]` 通过高级阶段检查。交付物：能继续沿同一逻辑再扩 40 次以上而不失去结构。

## 怎么继续往后生长

160 次不是终点，而是第一圈。

如果后面还要继续扩，可以继续按同样结构往外长：

- 第 2 圈：把每个模块的“草稿主题”落成可运行 lesson
- 第 3 圈：把每条 lesson 补成 Lua / Flow / MCP 三视角
- 第 4 圈：把课程和交付规范、评审规范、讲师素材打通

配套文档：

- [curriculum-overview.md](curriculum-overview.md)
- [track-newbie.md](track-newbie.md)
- [track-junior.md](track-junior.md)
- [track-intermediate.md](track-intermediate.md)
- [track-advanced.md](track-advanced.md)
- [evolution-playbook.md](evolution-playbook.md)
