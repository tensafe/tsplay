# Action: `get-session`

`get-session` 用来查看一个命名会话的详细信息。

## 最小命令

```bash
go run . -action get-session -session-name demo_login
```

## 常用参数

- `-session-name`：必填，要查看的会话名
- `-artifact-root`：会话注册表根目录

## 适合什么时候用

- 想确认某个会话保存的是 storage state 还是 persistent profile
- 想看创建时间、更新时间、最近使用时间
- 想排查会话是不是保存到了预期目录

## 输出结果

- 返回这个会话的详细 JSON

## 注意事项

- 如果会话不存在，会直接报错
- 如果只是想批量扫一眼，先用 [list-sessions](list-sessions.md)

## 相关文档

- [Lesson 41](../tutorials/41-inspect-named-session.md)
- [save-session](save-session.md)
