# TSPlay Flow 中文排错与修复对策

## 这份文件适合什么时候看

当用户已经有一条 Flow，但出现报错、超时、断言失败、变量不对、权限不够、登录态失效，或者 MCP finalize 没有直接得到可运行结果时，先看这份文件。

## 排错总思路

1. 先看报错属于哪一类：页面选择器、等待时机、断言结果、变量传递、文件权限、登录态、MCP 状态。
2. 优先做最小修复，不要一上来把整条 Flow 推翻重写。
3. 先修页面就绪和 selector，再修业务断言，再修输出结构。
4. 如果是中文需求但用户不知道 selector，优先考虑 MCP 的 `observe_page` 路线。

## 常见报错 1: 找不到 selector 或 `wait_for_selector` 超时

常见现象：

- `wait_for_selector` 超时
- `click` 时报元素不存在
- `type_text` 时报找不到输入框

优先检查：

- 页面是不是跳转到了预期 URL
- 目标元素是不是在登录后才出现
- 是不是少了一步 `wait_for_selector`
- selector 是不是写得太脆弱

优先修法：

- 在关键交互前补 `wait_for_selector`
- 把业务断言前移，确认页面已经到对的状态
- 如果用户不确定 selector，改成先 `observe_page`
- 把过于依赖层级的 selector 改成更稳定的 id、data-testid、文本或业务标识

修法示例：

```yaml
- action: navigate
  url: "{{page_url}}"

- action: wait_for_selector
  selector: "#import-form"
  timeout: 5000

- action: type_text
  selector: "#name"
  text: "{{row.name}}"
```

## 常见报错 2: `assert_text` 失败

常见现象：

- 页面已经点了按钮，但断言的成功文本不匹配
- 断言太早执行，页面还没刷新完成

优先检查：

- 断言文本是不是完整写错了
- 页面上是不是只会出现包含关系，而不是完全相等
- 点击后是不是需要先等待

优先修法：

- 保持 `assert_text` 断业务结果，不要断太底层的瞬时文本
- 适当增加 `timeout`
- 如果页面有明确的结果区域，先 `assert_visible` 再 `assert_text`

修法示例：

```yaml
- action: click
  selector: "#submit"

- action: assert_visible
  selector: "#submit-status"
  timeout: 5000

- action: assert_text
  selector: "#submit-status"
  text: "Imported"
  timeout: 5000
```

## 常见报错 3: `extract_text` 拿到空值或格式不对

常见现象：

- 变量为空
- 提取出来是一整段文本，但其实只想要数字或状态码

优先检查：

- selector 对不对
- 页面内容是不是晚于当前步骤才渲染出来
- 是否应该加 `pattern`

优先修法：

- 在提取前先加等待或断言
- 如果只想提数字或局部文本，补 `pattern`
- 把 `save_as` 改成业务语义名，方便排查后续变量链路

修法示例：

```yaml
- action: extract_text
  selector: "#summary-count"
  timeout: 5000
  pattern: "([0-9]+)"
  save_as: order_count
```

## 常见报错 4: 变量为空、`save_as` 混乱、输出结构不对

常见现象：

- `write_json` 写出来结构不完整
- `append_var` 后结果列表不对
- 后续步骤引用了不存在的变量

优先检查：

- 前面的 `save_as` 名字是不是写错
- `set_var` 和 `append_var` 是不是该用 `with.value`
- 变量名是不是只反映 DOM，而不是业务含义

优先修法：

- 用清楚的变量名，比如 `import_results`、`auth_status`、`page_title`
- 写对象时优先用 `with.value`
- 别把多个不相关结果全塞进一个模糊变量

修法示例：

```yaml
- action: set_var
  save_as: payload
  with:
    value:
      auth_status: "{{auth_status}}"
      import_count: "{{import_count}}"
```

## 常见报错 5: `read_csv` / `read_excel` 读不到文件

常见现象：

- 文件不存在
- 读出来的行数不对
- Excel 明明有数据但结果为空

优先检查：

- 文件路径是不是对的
- Excel 的 `sheet`、`range`、`headers` 是否写对
- MCP 模式下有没有文件访问权限

优先修法：

- 先确认文件路径
- Excel 有复杂布局时明确写 `sheet` 和 `range`
- 数据区域没有表头时明确写 `with.headers`
- MCP 下需要文件读写时补 `browser_write` 或对应 `allow_file_access`

## 常见报错 6: `read_json` 读不到或解析失败

常见现象：

- `read_json open ...` 报错
- `read_json parse ...` 报错
- 读取成功了，但后面的字段引用拿不到值

优先检查：

- JSON 文件路径是不是对的
- 文件内容是不是合法 JSON
- 后续字段路径是不是写对了，比如 `{{payload.meta.status}}`
- 受限上下文里是否补了文件访问权限

优先修法：

- 先确认文件实际存在
- 如果是前一步刚写出的 artifact，确认写入路径和读取路径一致
- 先把整个 JSON `save_as: payload`，再逐步提字段
- 用业务语义变量名，不要一上来就写很长、很难读的字段链路

修法示例：

```yaml
- action: read_json
  file_path: artifacts/payload.json
  save_as: payload

- action: set_var
  save_as: status
  value: "{{payload.meta.status}}"
```

## 常见报错 7: `write_json` / `write_csv` 写不出去

常见现象：

- 路径不允许
- 文件没生成
- 结果写到了不易维护的位置

优先检查：

- MCP 下是否有文件权限
- 路径是不是落在允许范围内
- 输出路径是不是过于随意，不利于后续 review

优先修法：

- 需要文件输出时补最小必要权限
- 保持 artifact 路径稳定
- JSON 和 CSV 尽量靠近同一任务目录

## 常见报错 8: 登录态失效、`use_session` 不生效

常见现象：

- 明明配了 `browser.use_session`，页面还是跳登录
- 断言时看到未登录状态

优先检查：

- session 名称是否正确
- 当前 session 是否已经过期
- 页面是不是换了域名或登录入口

优先修法：

- 先写一条最小验证 Flow 检查会话是否真的可用
- 登录态相关配置放在顶层 `browser`
- 不要把登录步骤散在每个业务 step 里

修法示例：

```yaml
browser:
  use_session: admin
```

## 常见报错 9: `foreach` 中断，整批任务全挂

常见现象：

- 某一行失败后整个导入停止
- 本来想写结果台账，但实际没有失败记录

优先修法：

- 在 `foreach` 内部局部包一层 `on_error`
- 成功和失败都显式写到结果列表
- 失败路径里把 `{{last_error}}` 存下来

修法示例：

```yaml
- action: foreach
  items: "{{rows}}"
  item_var: row
  steps:
    - action: on_error
      steps:
        - action: click
          selector: "#submit"
      on_error:
        - action: append_var
          save_as: import_results
          with:
            value:
              source_row: "{{row.source_row}}"
              error: "{{last_error}}"
```

## 常见报错 10: 页面状态不稳定，需要重试或轮询

常见现象：

- 偶发点击成功，偶发失败
- 页面状态延迟变化

优先修法：

- 短时易抖动用 `retry`
- 状态轮询用 `wait_until`
- 断言业务结果，不要只断 click 执行了

## 常见报错 11: MCP `finalize_flow` 没有直接 ready

常见现象：

- `status=needs_input`
- `status=needs_permission`
- `status=needs_repair`

修法思路：

- `needs_input`: 补页面、变量、输入文件、目标输出等缺失信息
- `needs_permission`: 补最小必要权限，不要默认全开
- `needs_repair`: 进入 `validate_flow`、`repair_flow_context`、`repair_flow`

中文解释模板：

```text
这条 Flow 现在还不能直接跑，原因是 <needs_input / needs_permission / needs_repair>。
下一步优先补 <缺的输入 / 缺的权限 / 修复上下文>，再继续 finalize 或 repair。
```

## 常见报错 12: `send_email` 被安全策略拦住

常见现象：

- 报错里出现 `allow_email`
- Flow 能跑到前面，但发邮件这一步被拒绝

优先修法：

- 在受限 Flow / MCP 上下文里补 `allow_email=true`
- 不要默认全开其他权限，只补最小必要权限

如果还有附件：

- 除了 `allow_email`，还要补文件访问权限

## 常见报错 13: 邮件发不出去或 SMTP 配置不对

常见现象：

- 连接 SMTP 失败
- 用户名密码不对
- TLS 模式不匹配

优先检查：

- 是走 `connection` 还是 `with.smtp`
- `host`、`port`、`username`、`password`、`from`、`tls_mode` 是否正确
- 如果走环境变量，`TSPLAY_EMAIL_*` 或 `TSPLAY_EMAIL_<NAME>_*` 是否已经配置

优先修法：

- 团队复用优先走 `connection`
- 一次性验证可以临时用 `with.smtp`
- 如果是 465 这类隐式 TLS，确认 `tls_mode: tls`

## 常见报错 14: 附件邮件失败

常见现象：

- 邮件步骤本身报附件相关错误
- 文件明明生成了，但附件没带上

优先检查：

- 附件路径是不是存在
- `with.attachments` 结构是否正确
- 受限上下文里是否补了文件访问权限

优先修法：

- 先确认附件文件确实在前一步已经生成
- 单个附件可用路径字符串或对象，多个附件用列表
- 需要更清楚文件名时，用 `{path, name, content_type}`

## 常见报错 15: 把问题修复杂了

常见现象：

- 本来只是 selector 不稳，结果整条 Flow 被大改
- 引入了没必要的 Lua
- 输出结构被一起改乱

优先修法：

- 先修最小范围
- 能补等待就别先重构
- 能补断言就别先换技术路径
- 能继续用 Flow 原生 action 就别先跳 Lua

## 最后排错顺序

1. 页面和登录态对不对
2. selector 和等待对不对
3. 断言是不是断到了业务结果
4. 变量链路和 `save_as` 对不对
5. 文件权限和输出路径对不对
6. 需要时再进入 MCP observe、validate、repair 链路
