# TSPlay 支持行为清单

这组文档说的不是命令行 `-action file-srv` 这一层，而是 `navigate`、`click`、`read_csv`、`db_query`、`retry` 这一层 Flow / Lua / MCP 能力动作。
你也可以把它理解成“TSPlay 支持的 action 总表”。

如果你要查的是命令行入口，请看 [CLI `-action` 参考](../actions/README.md)。

## 先说明三个维度

- `Flow`：是否可以作为 `steps[].action` 或等价的 Flow 配置能力使用
- `Lua`：是否可以在 Lua CLI / `-script` 里直接调用
- `MCP`：是否可以出现在通过 MCP 生成、校验、运行的 Flow 里

这里的 `MCP=是`，表示它能作为 Flow 能力被 `tsplay.flow_schema`、`tsplay.finalize_flow`、`tsplay.run_flow` 等链路使用；不代表每个动作都有单独同名的 MCP tool。

## 这套文档怎么读

- 每一类页都尽量先给动作表，再给一段“最小示例小代码”
- 有 Lua 对应能力的页面，通常会同时给 `Flow` 和 `Lua` 两种短例子
- 没有 Lua 对应能力的页面，比如 `retry / foreach / on_error / wait_until`，会只给 `Flow` 例子
- 如果你只是想确认字段名，优先看示例；如果你要系统理解边界，再看后面的建议和教程链接

## 核心能力矩阵

| 能力类别 | 典型动作 | Flow | Lua | MCP | 建议 | 说明页 |
| --- | --- | --- | --- | --- | --- | --- |
| 页面原子动作 | `navigate`、`click`、`type_text`、`select_option` | 是 | 是 | 是 | 保持强同步 | [页面原子动作](page-primitives.md) |
| 文件与表格 I/O | `screenshot`、`save_html`、`read_json`、`read_csv`、`read_excel`、`write_json`、`write_csv`、`write_excel` | 是 | 是 | 是 | 保持强同步；MCP 下受 `allow_file_access` 约束 | [文件与表格 I/O](file-and-spreadsheet-io.md) |
| HTTP 请求 | `http_request`、`json_extract` | 是 | 是 | 是 | 保持强同步；Lua 在 Flow / MCP 上下文里也遵守授权边界 | [HTTP 请求](http-requests.md) |
| 邮件通知 | `send_email` | 是 | 是 | 是 | 保持强同步；Lua 在 Flow / MCP 上下文里也遵守 `allow_email` | [邮件通知](email-delivery.md) |
| Redis 操作 | `redis_get`、`redis_set`、`redis_del`、`redis_incr` | 是 | 是 | 是 | 保持强同步；Lua 在 Flow / MCP 上下文里也遵守 `allow_redis` | [Redis 操作](redis-operations.md) |
| 数据库操作 | `db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction` | 是 | 是 | 是 | 保持强同步；`db_transaction` 自动提交或回滚 | [数据库操作](database-operations.md) |
| 浏览器状态 | `get_storage_state`、`get_cookies_string`、`browser.use_session` | 是 | 是 | 是 | 保持强同步；MCP 下受 `allow_browser_state` 约束 | [浏览器状态](browser-state.md) |
| Flow 便捷动作 | `extract_text`、`assert_visible`、`assert_text`、`set_var`、`append_var` | 是 | 是 | 是 | 已对齐；更适合作为编排语义糖 | [Flow 便捷动作](flow-convenience.md) |
| Flow 控制流 | `retry`、`if`、`foreach`、`on_error`、`wait_until` | 是 | 否 | 是 | 不需要硬同步到 Lua | [Flow 控制流](flow-control.md) |
| Lua 回调型能力 | `intercept_request` | 否 | 是 | 否 | 保持 Lua 专属更自然 | [Lua 回调型能力](lua-callbacks.md) |

## 常查 Action 快速索引

- 页面原子动作：`navigate`、`click`、`type_text`、`select_option`
- 文件与表格 I/O：`screenshot`、`save_html`、`read_json`、`read_csv`、`read_excel`、`write_json`、`write_csv`、`write_excel`
- HTTP 请求：`http_request`、`json_extract`
- 邮件通知：`send_email`
- Redis 操作：`redis_get`、`redis_set`、`redis_del`、`redis_incr`
- 数据库操作：`db_insert`、`db_insert_many`、`db_upsert`、`db_query`、`db_query_one`、`db_execute`、`db_transaction`
- 浏览器状态：`get_storage_state`、`get_cookies_string`、`browser.use_session`
- Flow 便捷动作：`extract_text`、`assert_visible`、`assert_text`、`set_var`、`append_var`
- Flow 控制流：`retry`、`if`、`foreach`、`on_error`、`wait_until`
- Lua 专属能力：`intercept_request`
- 其他常用浏览器动作：`get_text`、`get_attribute`、`get_html`、`get_all_links`、`capture_table`、`upload_file`、`upload_multiple_files`、`download_file`、`download_url`、`accept_alert`、`dismiss_alert`、`set_alert_text`、`execute_script`、`evaluate`、`new_tab`、`close_tab`、`switch_to_tab`、`find_element`、`find_elements`、`is_visible`、`is_enabled`、`block_request`、`get_response`

## 补充浏览器动作

除了上面这张核心矩阵，TSPlay 还有一批常用但不一定出现在“第一眼矩阵”里的浏览器动作，比如：

- `get_text`、`get_attribute`、`get_html`、`get_all_links`、`capture_table`
- `upload_file`、`upload_multiple_files`、`download_file`、`download_url`
- `accept_alert`、`dismiss_alert`、`set_alert_text`
- `execute_script`、`evaluate`
- `new_tab`、`close_tab`、`switch_to_tab`
- `find_element`、`find_elements`、`is_visible`、`is_enabled`
- `block_request`、`get_response`

这部分统一见 [补充浏览器动作](supplemental-browser-actions.md)。

## 推荐判断原则

- 像 `HTTP / Redis / 数据库 / 文件读写` 这种原子数据动作，最好在 `Flow` 和 `Lua` 两边都能成立，这样探索、固化、接入三条路径不会断层。
- 像 `retry / foreach / on_error / wait_until` 这种编排能力，更适合保留在 `Flow DSL`。
- 像 `extract_text / assert_text / assert_visible` 这种语义增强动作，可以继续优先作为 Flow 友好的高层封装。
- 像 `intercept_request` 这种需要回调和运行时介入的能力，保留 Lua 专属通常更自然。

## 安全边界提醒

- 文件读写、上传、下载、截图、保存 HTML：重点看 `allow_file_access`
- HTTP：重点看 `allow_http`
- 邮件：重点看 `allow_email`
- Redis：重点看 `allow_redis`
- 数据库：重点看 `allow_database`
- 浏览器状态：重点看 `allow_browser_state`
- `lua` 逃生口：重点看 `allow_lua`

## 相关入口

- [项目总览（中文）](../../README.zh-CN.md)
- [AI 协作入门](../training/ai-intent-to-flow.md)
- [Lesson 111](../tutorials/111-mcp-list-actions.md)
- [Lesson 112](../tutorials/112-mcp-flow-schema-and-examples.md)
