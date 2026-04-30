# 能力动作类别：邮件通知

邮件动作只有一个核心入口，但在交付里很常见：跑完通知、失败告警、发送附件、发送汇总。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `send_email` | 是 | 是 | 是 | `action: send_email` / `send_email({to=..., subject=..., body=...})` | 通过 SMTP 发邮件。支持 `to/cc/bcc`、文本或 HTML、附件、命名连接、超时和自定义头。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: email_demo
steps:
  - action: send_email
    save_as: email_result
    with:
      to:
        - ops@example.com
      subject: "TSPlay run finished"
      body: "Import completed."
      connection: alerts
```

### Lua

```lua
local result = send_email({
  to = {"ops@example.com"},
  subject = "TSPlay run finished",
  body = "Import completed.",
  connection = "alerts",
})
print(result)
```

## 使用建议

- 一般优先把 SMTP 连接配成命名环境变量，不要把敏感信息硬写进 Flow
- 带附件时，同时要满足 `allow_email` 和 `allow_file_access`
- 更适合在“结果已产出”的节点发送，不要把关键业务逻辑藏进邮件里

## 相关教程

- [AI 无感入门](../training/ai-intent-to-flow.md)
- [学习路径](../training/learning-path.md)
