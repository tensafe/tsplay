# 160 次递进迭代路线图

这份路线图的目的，不是追求“看起来很多”，而是把教程建设拆成 160 个不跳跃、可执行、可验证、可持续增长的迭代点。

使用方式：

- 学员：按阶段顺序推进，不要跨阶段乱跳
- 课程作者：每次至少落实 1 到 3 个迭代点
- 讲师 / 实施 / 负责人：拿它做课程排期和缺口检查

补充说明：

- 当前仓库里已经直接落地、可运行的是 `Lesson 01-110`
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

- `[071]` 结合 `Lesson 28-30`、`Lesson 36-80` 学会根据输出 JSON 回看浏览器状态结果。交付物：指出一条状态快照里的关键字段。
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

## 模块 09：证据回放、交接与发布前检查

- `[081]` 跑通 `Lesson 81`，从 `Lesson 80` 的生命周期 CSV 里重新读回批次证据。交付物：`81-read-lifecycle-evidence-*.json`。
- `[082]` 跑通 `Lesson 82`，根据生命周期证据回放一条新批次。交付物：`82-replay-batch-from-lifecycle-evidence-*.csv/json`。
- `[083]` 跑通 `Lesson 83`，验证回放批次和生命周期证据仍然一致。交付物：`83-verify-replay-batch-against-lifecycle-evidence-*.json`。
- `[084]` 跑通 `Lesson 84`，为 replay 批次补一条独立审计记录。交付物：`84-write-replay-audit-row-*.json`。
- `[085]` 跑通 `Lesson 85`，把原批次和 replay 批次的审计导出成一份对照 CSV。交付物：`85-export-original-and-replay-audits-*.csv/json`。
- `[086]` 跑通 `Lesson 86`，把生命周期、Redis、Postgres、审计对照压成一份回放对账包。交付物：`86-build-post-replay-reconciliation-pack-*.csv/json`。
- `[087]` 跑通 `Lesson 87`，把关键产物整理成交接 manifest。交付物：`87-build-handoff-artifact-manifest-*.csv/json`。
- `[088]` 跑通 `Lesson 88`，把 manifest 再压成一份交付摘要。交付物：`88-build-handoff-summary-*.json`。
- `[089]` 跑通 `Lesson 89`，把交接包整理成发布前检查清单。交付物：`89-build-pre-release-checklist-*.csv/json`。
- `[090]` 跑通 `Lesson 90`，把“生命周期证据 -> 回放 -> 交接包”重新串成一条完整 round trip。交付物：`90-handoff-round-trip-from-lifecycle-evidence-*.csv/json`。

## 模块 10：模板目录、模板索引与模板发布前检查

- `[091]` 跑通 `Lesson 91`，从交接 manifest 里识别每份产物的角色。交付物：`91-read-handoff-manifest-roles-*.csv/json`。
- `[092]` 跑通 `Lesson 92`，把交接产物整理成模板目录。交付物：`92-build-template-artifact-catalog-*.csv/json`。
- `[093]` 跑通 `Lesson 93`，把交接链整理成 `Input -> Process -> Output` 模板。交付物：`93-build-input-process-output-template-*.csv/json`。
- `[094]` 跑通 `Lesson 94`，把交接链整理成 `Collect -> Verify -> Save` 模板。交付物：`94-build-collect-verify-save-template-*.csv/json`。
- `[095]` 跑通 `Lesson 95`，把交接链整理成 `Replay -> Audit -> Handoff` 模板。交付物：`95-build-replay-audit-handoff-template-*.csv/json`。
- `[096]` 跑通 `Lesson 96`，把几份模板整理成统一索引。交付物：`96-build-template-index-*.csv/json`。
- `[097]` 跑通 `Lesson 97`，验证模板索引仍然覆盖完整交接链。交付物：`97-verify-template-covers-handoff-chain-*.csv/json`。
- `[098]` 跑通 `Lesson 98`，生成一份“场景 -> 模板”的学习矩阵。交付物：`98-build-template-lesson-matrix-*.csv/json`。
- `[099]` 跑通 `Lesson 99`，给模板包生成发布前检查清单。交付物：`99-build-template-preflight-checklist-*.csv/json`。
- `[100]` 跑通 `Lesson 100`，把交接产物重新收成一份模板包 round trip。交付物：`100-template-round-trip-from-handoff-artifacts-*.csv/json`。

## 模块 11：健壮性、等待、恢复

- `[101]` 跑通 `Lesson 101`，先确认模板发布卡片和关键 badge 真的在页面上。交付物：`101-assert-visible-template-release-card-*.json`。
- `[102]` 跑通 `Lesson 102`，继续确认模板发布状态和摘要文字是对的。交付物：`102-assert-text-template-release-status-*.json`。
- `[103]` 跑通 `Lesson 103`，用 `retry` 处理模板发布 gate 的临时不一致。交付物：`103-retry-template-release-gate-*.json`。
- `[104]` 跑通 `Lesson 104`，用 `wait_until` 等待异步 stage check 变成 ready。交付物：`104-wait-until-template-release-ready-*.json`。
- `[105]` 跑通 `Lesson 105`，故意触发一次 ticket 校验失败并在错误分支里恢复。交付物：`105-on-error-template-release-validation-*.json`。
- `[106]` 跑通 `Lesson 106`，等待一条延迟出现的发布说明项。交付物：`106-wait-for-delayed-release-note-*.json`。
- `[107]` 跑通 `Lesson 107`，用 `retry` 接住一次偶发失败点击。交付物：`107-retry-flaky-publish-click-*.json`。
- `[108]` 跑通 `Lesson 108`，通过 `reload + retry` 验证恢复状态已经真正生效。交付物：`108-reload-and-retry-release-recovery-*.json`。
- `[109]` 跑通 `Lesson 109`，给模板发布页保存整页图、元素图和 HTML。交付物：`109-template-release-artifact-pack-*.png/html/json`。
- `[110]` 跑通 `Lesson 110`，把断言、等待、重试、恢复和证据留存重新串成一条完整 round trip。交付物：`110-template-release-robustness-round-trip-*.png/html/json`。

## 模块 12：MCP 基础链路与 repair 入门

- `[111]` 跑通 `Lesson 111`，先用 `tsplay.list_actions` 建立 MCP 能力地图。交付物：`111-mcp-list-actions.json`。
- `[112]` 跑通 `Lesson 112`，把 `flow_schema` 和 `flow_examples` 都留成本地参考。交付物：`112-mcp-flow-schema.json`、`112-mcp-flow-examples.json`。
- `[113]` 跑通 `Lesson 113`，对模板发布页生成第一份 observation。交付物：`113-mcp-observe-page-template-release.json`。
- `[114]` 跑通 `Lesson 114`，把 observation 变成第一份 `draft.flow_yaml`。交付物：`114-mcp-draft-flow-template-release.json`。
- `[115]` 跑通 `Lesson 115`，先校验 draft 再进入运行。交付物：`115-mcp-validate-drafted-template-release.json`。
- `[116]` 跑通 `Lesson 116`，执行 draft flow 并拿到结果 trace。交付物：`116-mcp-run-drafted-template-release.json`。
- `[117]` 跑通 `Lesson 117`，先记录一次失败运行，再生成 repair context。交付物：`117-mcp-run-broken-template-release.json`、`117-mcp-repair-flow-context-template-release.json`。
- `[118]` 跑通 `Lesson 118`，把 repair context 收成统一 repair request。交付物：`118-mcp-repair-flow-template-release.json`。
- `[119]` 跑通 `Lesson 119`，把 `observe -> draft -> validate -> run -> repair` 串成完整心智模型。交付物：一页 MCP 基础链路复盘。
- `[120]` 跑通 `Lesson 120`，用 `finalize_flow` 收成更短的默认入口。交付物：`120-mcp-finalize-flow-template-release.json`。

## 高级教程

## 模块 13：安全边界与运行边界

- `[121]` 跑通 `Lesson 121`，先把 `allow_lua` 的 blocked / allowed 对照跑清楚。交付物：`121-mcp-validate-allow-lua-blocked.json`、`121-mcp-validate-allow-lua-allowed.json`。
- `[122]` 跑通 `Lesson 122`，把 `allow_http` 的 blocked / allowed 对照跑清楚。交付物：`122-mcp-validate-allow-http-blocked.json`、`122-mcp-validate-allow-http-allowed.json`。
- `[123]` 跑通 `Lesson 123`，把 `allow_file_access` 的 blocked / allowed 对照跑清楚。交付物：`123-mcp-validate-allow-file-access-blocked.json`、`123-mcp-validate-allow-file-access-allowed.json`。
- `[124]` 跑通 `Lesson 124`，把 `allow_browser_state` 的 blocked / allowed 对照跑清楚。交付物：`124-mcp-validate-allow-browser-state-blocked.json`、`124-mcp-validate-allow-browser-state-allowed.json`。
- `[125]` 跑通 `Lesson 125`，把 `allow_redis` 的 blocked / allowed 对照跑清楚。交付物：`125-mcp-validate-allow-redis-blocked.json`、`125-mcp-validate-allow-redis-allowed.json`。
- `[126]` 跑通 `Lesson 126`，把 `allow_database` 的 blocked / allowed 对照跑清楚。交付物：`126-mcp-validate-allow-database-blocked.json`、`126-mcp-validate-allow-database-allowed.json`。
- `[127]` 跑通 `Lesson 127`，写出“本地 Flow 和 MCP”两种入口的边界对照。交付物：一页对照说明和本地 `123` Flow 运行结果。
- `[128]` 跑通 `Lesson 128`，写出“为什么教程不能一开始就忽略授权边界”。交付物：一页说明。
- `[129]` 跑通 `Lesson 129`，把 `security_preset` 和显式 `allow_*` 覆盖关系跑清楚。交付物：`129-mcp-validate-file-access-browser-write.json`、`129-mcp-validate-http-full-automation.json`、`129-mcp-validate-http-full-automation-override.json`。
- `[130]` 跑通 `Lesson 130`，把安全边界模块收成第一轮 checkpoint。交付物：一页边界复盘。

## 模块 14：大型 Flow、规范、评审

- `[131]` 跑通 `Lesson 131`，对比一条 before / after Flow，先把 `name` 和 `description` 的 review 入口立起来。交付物：`artifacts/review-output.json`、`artifacts/tutorials/131/review-summary.json`。
- `[132]` 跑通 `Lesson 132`，继续复用 `Lesson 131` 的样例，把 `save_as` 命名规则收成最小标准。交付物：一页变量命名复盘。
- `[133]` 跑通 `Lesson 133`，把 artifact 平铺输出改成 lesson 级目录输出。交付物：`artifacts/output.json`、`artifacts/tutorials/133/review-layout/output.json`、`artifacts/tutorials/133/review-layout/manifest.json`。
- `[134]` 跑通 `Lesson 134`，用统一 checklist review `Lesson 131` 和 `Lesson 133`。交付物：一页 review checklist 复盘。
- `[135]` 跑通 `Lesson 135`，先运行一个最小 Lua escape hatch，再判断它是不是结构上应该保留。交付物：`artifacts/tutorials/135/review-summary-from-lua.json`。
- `[136]` 跑通 `Lesson 136`，把上一节的简单 Lua 编排抽回 Flow。交付物：`artifacts/tutorials/136/review-summary-from-flow.json`。
- `[137]` 跑通 `Lesson 137`，先盘点现有 demo，再判断什么时候才该新增页面。交付物：`artifacts/tutorials/137-list-assets.txt` 和一页 demo 决策复盘。
- `[138]` 跑通 `Lesson 138`，把教程最小交付物收成一份稳定 checklist。交付物：一页交付物清单复盘。
- `[139]` 跑通 `Lesson 139`，把大型 Flow 包拆成 `collect -> verify -> publish` 三段目录骨架。交付物：`artifacts/tutorials/139/collect/raw-items.json`、`artifacts/tutorials/139/verify/verification-summary.json`、`artifacts/tutorials/139/publish/publish-manifest.json`。
- `[140]` 跑通 `Lesson 140`，把安全边界和 review 思维重新收成一张高级阶段心智图。交付物：一页复盘说明。

## 模块 15：发布包、内置资产、交付体验

- `[141]` 跑通 `Lesson 141`，先用 `list-assets` 证明二进制里真的带着 `ReadMe.md`、`docs/`、`script/`、`demo/`。交付物：`artifacts/tutorials/141-bundled-assets.txt`。
- `[142]` 跑通 `Lesson 142`，把 `list-assets` 输出重新整理成“教程 / 脚本 / demo”三类入口心智模型。交付物：一页 `list-assets` 复盘。
- `[143]` 跑通 `Lesson 143`，把内置资源真正释放到本地目录并核对关键入口。交付物：`./tsplay-assets-143` 和 `artifacts/tutorials/143-extracted-assets.txt`。
- `[144]` 跑通 `Lesson 144`，把单二进制交付流程收成稳定 manifest。交付物：`artifacts/tutorials/144/single-binary-delivery-manifest.json`。
- `[145]` 跑通 `Lesson 145`，把离线学习路径收成 manifest。交付物：`artifacts/tutorials/145/offline-learning-manifest.json`。
- `[146]` 跑通 `Lesson 146`，把必须内置和可选资源收成策略 JSON。交付物：`artifacts/tutorials/146/embedded-asset-policy.json`。
- `[147]` 跑通 `Lesson 147`，把 `file-srv` 的开发态和发布态对照说明写清楚。交付物：一页 dev / release 模式说明。
- `[148]` 跑通 `Lesson 148`，把新用户 first-run 路径收成固定入口 manifest。交付物：`artifacts/tutorials/148/first-run-entry-manifest.json`。
- `[149]` 跑通 `Lesson 149`，把二进制版本和资源版本拆成两层维护策略。交付物：`artifacts/tutorials/149/release-asset-version-strategy.json`。
- `[150]` 跑通 `Lesson 150`，把“单二进制 + 内置教程”收成完整交付复盘。交付物：`artifacts/tutorials/150/single-binary-delivery-summary.json`。

## 模块 16：Capstone、培训、持续演化

- `[151]` 跑通 `Lesson 151`，把新手阶段收成一份 capstone brief。交付物：`artifacts/tutorials/151/newbie-capstone-brief.json`。
- `[152]` 跑通 `Lesson 152`，把初级阶段收成一份本地 demo 与文件能力的 capstone brief。交付物：`artifacts/tutorials/152/junior-capstone-brief.json`。
- `[153]` 跑通 `Lesson 153`，把 MCP 草拟、执行、修复链收成一份中级 capstone brief。交付物：`artifacts/tutorials/153/intermediate-capstone-brief.json`。
- `[154]` 跑通 `Lesson 154`，把安全边界、review 和交付策略收成一份高级 capstone brief。交付物：`artifacts/tutorials/154/advanced-capstone-brief.json`。
- `[155]` 跑通 `Lesson 155`，把新人的第一周收成固定 onboarding 节奏。交付物：`artifacts/tutorials/155/new-hire-7-day-plan.json`。
- `[156]` 跑通 `Lesson 156`，把实施同学的前两周收成“教程 -> 试点 -> handoff”节奏。交付物：`artifacts/tutorials/156/implementer-2-week-plan.json`。
- `[157]` 跑通 `Lesson 157`，把教程、Bootcamp、Capstone、评审和讲师手册接成备课顺序。交付物：`artifacts/tutorials/157/trainer-prep-sequence.json`。
- `[158]` 跑通 `Lesson 158`，把月度教程缺口复盘收成稳定模板。交付物：`artifacts/tutorials/158/tutorial-gap-review-cycle.json`。
- `[159]` 跑通 `Lesson 159`，把“每 10 次迭代回看一次”的检查机制写成稳定 checklist。交付物：`artifacts/tutorials/159/every-10-iteration-review.json`。
- `[160]` 跑通 `Lesson 160`，把下一圈教程扩展计划收成 continuation plan。交付物：`artifacts/tutorials/160/curriculum-continuation-plan.json`。

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
