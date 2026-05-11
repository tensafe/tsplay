# Action: `export-session`

`export-session` 会把命名会话导出成可直接复用的浏览器片段或 Flow 片段，适合把“保存下来的登录态”重新接回交付脚本。

## 最小命令

```bash
go run . -action export-session -session-name demo_login
```

## 常见用法

```bash
go run . -action export-session \
  -session-name demo_login \
  -session-format flow_yaml
```

## 常用参数

- `-session-name`：必填，会话名
- `-session-format`：导出格式，默认 `all`

## 可用格式

- `all`
- `browser` / `browser_yaml`
- `expanded_browser` / `expanded_browser_yaml`
- `flow` / `flow_yaml`
- `expanded_flow` / `expanded_flow_yaml`
- `browser_json`
- `expanded_browser_json`
- `flow_json`
- `expanded_flow_json`

## 适合什么时候用

- 想快速得到 `browser.use_session` 片段
- 想看推荐写法和展开写法的区别
- 想把会话导出成 YAML 或 JSON 给 Flow / 工具继续用

## 注意事项

- 一般优先复用推荐写法，而不是一上来就用 expanded 版本
- 如果只想看会话元数据，不必先导出，先看 [get-session](get-session.md)

## 相关文档

- [Lesson 41](../tutorials/41-inspect-named-session.md)
- [Lesson 57](../tutorials/57-use-session-import-export-round-trip.md)
