# Action: `save-session`

`save-session` 用来把可复用的登录态保存成命名会话。你可以保存 storage state，也可以注册持久化 profile。

## 最小命令

```bash
go run . -action save-session \
  -session-name demo_login \
  -storage-state-path artifacts/storage-state.json
```

## 另一种常见用法

```bash
go run . -action save-session \
  -session-name demo_profile \
  -profile-name default
```

## 常用参数

- `-session-name`：必填，会话名
- `-storage-state-path`：从已有文件复制 storage state
- `-storage-state-json`：直接传入 storage state JSON
- `-profile-name`：注册持久化 profile
- `-profile-session`：可选，profile 下的具体 session
- `-artifact-root`：会话注册表根目录

## 适合什么时候用

- 登录一次后，后面多条 Flow 反复复用
- 想把浏览器登录态从一次人工操作转成可交付会话
- 想在 training / demo 环境里预放一个稳定会话

## 注意事项

- `storage_state` 和 `profile/session` 两种保存方式不能混用
- 会话名要保持稳定、可读、可复用

## 相关文档

- [Lesson 40](../tutorials/40-save-named-session.md)
- [Lesson 46](../tutorials/46-save-import-session.md)
- [export-session](export-session.md)
