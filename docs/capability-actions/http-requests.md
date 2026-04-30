# 能力动作类别：HTTP 请求

这组动作把“浏览器自动化”和“直接调接口”接到一起。  
在 Flow / MCP 安全上下文里，重点关注 `allow_http`，而 `save_path` 这类落文件行为还会继续受 `allow_file_access` 约束。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `http_request` | 是 | 是 | 是 | `action: http_request` / `http_request({url=..., method=..., json=...})` | 发起 HTTP 请求。支持 headers、query、json、form、multipart、保存响应文件，以及复用浏览器 cookies / referer / UA。 |
| `json_extract` | 是 | 是 | 是 | `action: json_extract` + `from,path` / `json_extract(value, '$.items[0]')` | 从 JSON 或 JSON 字符串里取值。适合把接口结果再拆成可复用变量。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: http_demo
steps:
  - action: http_request
    url: http://127.0.0.1:8000/demo/data/order_summary.json
    save_as: api_result
    with:
      response_as: json

  - action: json_extract
    from: "{{api_result}}"
    path: "$.body.summary.open"
    save_as: open_count
```

### Lua

```lua
local response = http_request({
  url = "http://127.0.0.1:8000/demo/data/order_summary.json",
  response_as = "json",
})
local open_count = json_extract(response, "$.body.summary.open")
print(open_count)
```

## 使用建议

- 页面能直接抓 API 时，`http_request` 往往比“继续点页面”更稳定
- `json_extract` 很适合和 `save_as`、`set_var` 串起来，把响应拆成后续步骤要用的字段
- `use_browser_cookies=true` 时，意味着这条请求会依赖浏览器上下文

## 相关教程

- [Lesson 05](../tutorials/05-http-request-json.md)
- [Lesson 71](../tutorials/71-external-system-round-trip.md)
- [Lesson 122](../tutorials/122-allow-http-boundary.md)
