# Action: `quickstart-demo`

`quickstart-demo` 会自动生成一条最小 demo Flow，并立刻执行它。  
这条 demo 只用 `set_var + write_json`，所以不需要先下载 Playwright。

## 最小命令

```bash
./tsplay -action quickstart-demo
```

如果你还没有二进制，最短路径直接看 [快速开始](../../getting-started.md) 里的 `sh / PowerShell` 安装脚本。

## 它会做什么

- 在 `artifacts/quickstart/` 下生成 `quickstart-demo.flow.yaml`
- 立刻执行这条 Flow
- 再生成 `artifacts/quickstart/quickstart-demo-output.json`
- 在终端输出这次 quickstart 的结构化结果

## 适合什么时候用

- 想先确认二进制能不能直接跑
- 想先体验 Flow 执行结果，但暂时不想碰浏览器
- 想避免第一次就等 Playwright 下载完成
- 想给新用户一个“下载即跑”的最短入口

## 注意事项

- 这条 demo 不打开浏览器，所以不会验证页面动作
- 如果你下一步想练页面自动化，再去跑 `file-srv` 或浏览器 Flow
- 首次真正执行浏览器相关 Flow 时，TSPlay 仍会自动安装 Playwright 浏览器

## 相关文档

- [快速开始](../../getting-started.md)
- [file-srv](file-srv.md)
- [extract-assets](extract-assets.md)
