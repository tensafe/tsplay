# Action: `record-screen`

`record-screen` 录的是整个 macOS 桌面，不只是浏览器页面。它适合录教程、录桌面演示、录窗口切换过程。

## 最小命令

```bash
go run . -action record-screen -record-cmd "go run . -flow script/tutorials/10_assert_page_state.flow.yaml"
```

## 常用参数

- `-record-cmd`：录屏期间顺手执行的命令
- `-record-input`：ffmpeg avfoundation 输入，默认 `Capture screen 0:none`
- `-record-output`：输出视频路径
- `-record-fps`：帧率
- `-record-size`：可选的视频尺寸
- `-record-cursor`：是否录鼠标
- `-record-warmup-ms`：正式执行前预热时长
- `-record-cooldown-ms`：命令结束后的收尾停留时长
- `-record-duration-ms`：录制时长上限
- `-record-crf`：编码质量参数
- `-record-preset`：编码速度预设
- `-record-shell`：执行 `-record-cmd` 时使用的 shell

## 适合什么时候用

- 要演示桌面切换、文件选择器、系统窗口
- 要给教程录完整桌面视频
- 要把某条命令的执行过程录成交付素材

## 注意事项

- 当前只支持 macOS
- 依赖 `ffmpeg`
- 如果你只想录浏览器页面内容，通常优先考虑 `-browser-video-output`

## 相关文档

- [list-record-devices](list-record-devices.md)
- [教程自动录屏](../training/tutorial-video-recording.md)
