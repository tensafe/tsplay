# Action: `list-sessions`

`list-sessions` 会列出当前 `artifact-root` 下所有已保存的命名会话。

## 最小命令

```bash
go run . -action list-sessions
```

## 适合什么时候用

- 想知道当前可复用的会话有哪些
- 想确认新保存的会话是不是已经写进注册表
- 想给 Workbench、培训环境或排障留一份会话列表

## 输出结果

- 返回 JSON
- 通常会包含 `artifact_root` 和 `sessions`

## 注意事项

- 它列的是注册表里的会话，不会展开成可直接贴进 Flow 的片段
- 如果要看某个会话细节，继续用 [get-session](get-session.md)

## 相关文档

- [Lesson 41](../tutorials/41-inspect-named-session.md)
- [get-session](get-session.md)
