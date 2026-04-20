# TSPlay 教程自动录屏

这份说明给讲师、Enablement 同学和教程维护者使用。  
目标不是做影视级后期，而是稳定地把“教程动作 + 讲解画面”录成可复用素材。

当前这套能力先支持：

- macOS
- 本机已安装 `ffmpeg`
- 用 `tsplay` 自动跑教程命令
- 在命令开始前自动开录，命令结束后自动收尾

## 什么时候适合用它

最适合下面这两类场景：

- 你已经有一条稳定可跑的 `Lua` / `Flow`，现在想把演示过程录下来
- 你要批量录教程，不想每次手动开、手动停录屏

如果是需要大量旁白、剪辑、字幕和镜头切换的正式课程，  
更推荐把这里的自动录屏当成“原始素材采集层”，再交给 OBS、Screen Studio、Descript 等工具做后期。

## Step 1: 先看本机能不能识别录屏设备

```bash
# 方式 A：直接运行源码
go run . -action list-record-devices

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action list-record-devices
```

预期结果：

- 会返回 `video_devices` 和 `audio_devices`
- 常见视频输入通常会类似 `Capture screen 0`

如果 `video_devices` 为空，通常先检查两件事：

- Terminal / Codex 是否已经拿到系统的屏幕录制权限
- `ffmpeg` 是否真的在当前 shell 的 `PATH` 里

## Step 2: 录一条“命令跑完就自动结束”的教程视频

如果你要录浏览器类教程，建议先在另一个终端把本地 demo 服务启动好：

```bash
./tsplay -action file-srv -addr :8000
```

然后在当前终端执行自动录屏：

```bash
./tsplay -action record-screen \
  -record-input 'Capture screen 0:none' \
  -record-output artifacts/recordings/lesson-10-assert-page-state.mp4 \
  -record-cmd './tsplay -flow script/tutorials/10_assert_page_state.flow.yaml'
```

这条命令的含义是：

- 用 `Capture screen 0:none` 作为 `ffmpeg` 录屏输入
- 先开始录屏
- 稍等一小段 warmup
- 自动运行 `Lesson 10` 的 Flow
- Flow 跑完后再留一点 cooldown
- 自动结束录屏

预期结果：

- 会生成 `artifacts/recordings/lesson-10-assert-page-state.mp4`
- 终端还会输出一份 JSON，总结这次录屏的输入、输出、fps、命令退出状态

## Step 3: 如果你只想手动演示，但不想手动开关录屏

也可以不传 `-record-cmd`：

```bash
./tsplay -action record-screen \
  -record-input 'Capture screen 0:none' \
  -record-output artifacts/recordings/manual-demo.mp4
```

这时：

- `tsplay` 会立即开始录屏
- 你可以自己手动操作页面、终端、讲解步骤
- 结束时按 `Ctrl-C`，录屏会自动收尾

## Step 4: 常用参数怎么调

- `-record-fps 30`
  默认 30 帧，教程演示通常已经够用
- `-record-size 1728x1117`
  需要固定分辨率时再加；不加时由设备自己决定
- `-record-cursor=false`
  如果你不想把鼠标指针录进去，可以关掉
- `-record-warmup-ms 1500`
  让录屏比命令早一点开始，避免开头被切掉
- `-record-cooldown-ms 1200`
  让录屏比命令晚一点结束，避免结尾太急
- `-record-duration-ms 300000`
  需要硬性限制最大时长时再加

## 一条更适合教学出片的建议

为了让学员更容易看懂，建议录制时尽量保持这几个习惯：

- 不要给教程命令加 `-headless`
- 录屏前先把浏览器和终端摆到固定位置
- 一条视频只讲一类动作，不要一口气跨太多 lesson
- 浏览器类视频尽量提前起好 `file-srv`
- 数据类、Redis、DB 类视频更适合后期补字幕或旁白

## 推荐的第一批录制对象

如果你要先录第一批最有价值的视频，建议先做：

1. `Lesson 01-05`
2. `Lesson 10-12`
3. `Lesson 13-20`
4. `Lesson 28-39`
5. `Lesson 40-57`
6. `Lesson 111-120`

这些内容最容易形成“看一遍就能跟着复现”的教学体验。
