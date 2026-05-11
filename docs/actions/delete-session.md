# Action: `delete-session`

`delete-session` 用来删除命名会话的注册信息，并在适用时清理对应的 storage state 文件。

## 最小命令

```bash
go run . -action delete-session -session-name demo_login
```

## 常用参数

- `-session-name`：必填，要删除的会话名
- `-artifact-root`：会话注册表根目录

## 适合什么时候用

- 清理测试用、临时用的命名会话
- 培训或演示结束后做环境收口
- 想避免后续 Flow 误用旧登录态

## 输出结果

- 返回删除结果 JSON
- 如果是 storage state 类型，会说明是否删除了对应文件

## 注意事项

- 对 persistent profile 类型，会保留 profile 目录本体，只删除注册信息
- 真要删除 profile 数据时，建议人工确认后再处理

## 相关文档

- [Lesson 43](../tutorials/43-delete-named-session.md)
- [list-sessions](list-sessions.md)
