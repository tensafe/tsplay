# Action: `cli`

`cli` 是最轻的一条交互入口。它适合你先把浏览器动作试出来，再决定要不要整理成 `-script` 或 `-flow`。

## 最小命令

```bash
go run . -action cli
```

## 适合什么时候用

- 想快速试一个 `navigate`、`click`、`type_text`
- 想临时探索页面，不想先建脚本文件
- 想在同一个会话里反复敲命令看结果

## 进入后怎么用

- 输入 `start`：启动 Playwright 浏览器与页面对象
- 直接输入 Lua：执行一段即时脚本
- 输入 `reset`：重置 Playwright 运行时
- 输入 `exit`：退出 CLI

## 常用参数

- `-headless`：隐藏浏览器窗口
- `-artifact-root`：把产物写到指定目录
- `-browser-cdp-launch`：自动查找本机 Chrome/Chromium/Edge，启动一个带远程调试端口的独立浏览器，再通过 CDP 接管
- `-browser-cdp-endpoint`：通过 CDP endpoint 接管已启动的 Chrome/Chromium，支持 `ws://127.0.0.1:9222/devtools/browser/...`、`http://127.0.0.1:9222` 或直接粘贴 `127.0.0.1:9222/json/version`
- `-browser-cdp-port`：通过本地远程调试端口接管已启动的 Chrome/Chromium，例如 `9222`
- `-browser-cdp-executable`：为 `-browser-cdp-launch` 指定浏览器可执行文件；不传时会自动搜索 macOS / Windows / Linux 常见位置
- `-browser-cdp-user-data-dir`：为 `-browser-cdp-launch` 指定独立用户数据目录；不传时默认在 artifact root 下创建

## 注意事项

- `cli` 更适合探索，不适合长期保存交付逻辑
- 使用 `-browser-cdp-endpoint` / `-browser-cdp-port` 接管外部浏览器时，TSPlay 退出只会断开连接，不会关闭真实浏览器
- 使用 `-browser-cdp-launch` 由 TSPlay 启动浏览器时，TSPlay 会在退出时回收这个独立浏览器进程
- 如果已经跑通一段稳定操作，通常更适合收成 `-script` 或 `-flow`

## 相关文档

- [项目总览（中文）](../../README.zh-CN.md)
- [新手学习路线](../tutorials/track-newbie.zh-CN.md)
- [学习路径](../training/learning-path.md)
