# 能力动作类别：Flow 控制流

这组能力本质上属于 Flow DSL 的编排层，而不是要强行复制到 Lua 的底层原语。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `retry` | 是 | 否 | 是 | `action: retry` + `times,interval_ms,steps` | 重试一组嵌套步骤。适合临时抖动和弱一致页面。 |
| `if` | 是 | 否 | 是 | `action: if` + `condition,then,else` | 按条件走分支。 |
| `foreach` | 是 | 否 | 是 | `action: foreach` + `items,item_var,steps` | 对列表逐项执行。支持进度 checkpoint。 |
| `on_error` | 是 | 否 | 是 | `action: on_error` + `steps,on_error` | 主步骤失败时执行错误处理块。 |
| `wait_until` | 是 | 否 | 是 | `action: wait_until` + `condition,timeout,interval_ms` | 轮询条件直到成功或超时。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: flow_control_demo
steps:
  - action: retry
    times: 3
    interval_ms: 200
    steps:
      - action: click
        selector: "#retry-button"
      - action: assert_text
        selector: "#retry-status"
        text: "Export complete"
        timeout: 500

  - action: wait_until
    timeout: 5000
    interval_ms: 200
    condition:
      action: is_visible
      selector: "#job-ready"
```
这类能力属于 Flow DSL 的编排层，不需要额外硬补一份 Lua 等价写法。

## 使用建议

- 页面偶发抖动时，优先 `retry`，不要一上来改 selector
- 批量导入、批量回放时，优先 `foreach`
- 恢复动作要清晰表达时，用 `on_error`
- 需要“直到变成真”为止时，用 `wait_until`，不要到处散 `sleep`

## 相关教程

- [Lesson 16](../tutorials/16-retry-flaky-action.md)
- [Lesson 21](../tutorials/21-if-optional-login.md)
- [Lesson 22](../tutorials/22-foreach-batch-import-csv.md)
- [Lesson 23](../tutorials/23-on-error-import-recovery.md)
- [Lesson 17](../tutorials/17-wait-until-ready.md)
