# 能力动作类别：Lua 回调型能力与逃生口

这组能力更贴近“运行时控制”或“最后的 escape hatch”，不适合强求和 Flow 完全对称。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `intercept_request` | 否 | 是 | 否 | `intercept_request(function(req) ... end)` | 拦截请求并用 Lua 回调处理。最适合保留为 Lua 专属。 |
| `lua` | 是 | 否 | 是 | `action: lua` + `code` | 在 Flow 里内联一段 Lua。适合小范围 escape hatch，但要谨慎使用。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: lua_escape_hatch_demo
steps:
  - action: lua
    save_as: computed_value
    code: |
      local base = 40
      return base + 2
```

### Lua

```lua
intercept_request(function(request)
  if string.find(request.url, "/tracking") then
    return ""
  end
  return request.url
end)
```

## 使用建议

- 正常页面步骤优先用结构化动作，不要一上来就写 `lua`
- 确实需要回调、动态改请求、或结构化动作不够表达时，再考虑 `intercept_request`
- `lua` 在 Flow / MCP 里属于高风险能力，通常需要明确开启 `allow_lua`

## 相关教程

- [Lesson 121](../tutorials/121-allow-lua-boundary.md)
- [Lesson 127](../tutorials/127-compare-local-flow-and-mcp-boundaries.md)
