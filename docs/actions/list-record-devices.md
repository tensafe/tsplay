# Action: `list-record-devices`

`list-record-devices` 用来列出 macOS 上可用的录屏设备信息，适合在正式录制前先检查 `ffmpeg`、设备名和系统权限。

## 最小命令

```bash
go run . -action list-record-devices
```

## 适合什么时候用

- 第一次配置 `record-screen`
- 不确定 `ffmpeg` 能不能找到 avfoundation 设备
- 想确认终端是否已经拿到系统的屏幕录制权限

## 输出结果

- 返回 JSON
- 通常会包含 `video_devices`、`audio_devices`、`permission_hint`

## 注意事项

- 当前只支持 macOS
- 需要系统里能找到 `ffmpeg`
- 如果没有列出设备，先检查终端的屏幕录制权限

## 相关文档

- [record-screen](record-screen.md)
- [教程自动录屏](../training/tutorial-video-recording.md)
