# Lesson 109: 给模板发布页留一份调试证据包

前面几节已经把稳定性动作一条条拆开。  
这一节开始把“现场证据”重新接回来，方便排障和回看。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `screenshot`
- `screenshot_element`
- `save_html`
- `write_json`

## 准备工作

先确认本地静态文件服务已经启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/109_template_release_artifact_pack.lua](../../script/tutorials/109_template_release_artifact_pack.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/109_template_release_artifact_pack.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/109_template_release_artifact_pack.lua
```

预期结果：

- 会生成 `artifacts/tutorials/109-template-release-full-page-lua.png`
- 会生成 `artifacts/tutorials/109-template-release-card-lua.png`
- 会生成 `artifacts/tutorials/109-template-release-artifact-pack-lua.html`
- 会生成 `artifacts/tutorials/109-template-release-artifact-pack-lua.json`

## Step 2: 为什么这节放在这里

因为到这时你已经知道：

- 页面哪部分最关键
- 哪些状态是成功现场
- 哪些状态最值得留证据

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/109_template_release_artifact_pack.flow.yaml](../../script/tutorials/109_template_release_artifact_pack.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/109_template_release_artifact_pack.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/109_template_release_artifact_pack.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/109-template-release-full-page-flow.png`
- 会生成 `artifacts/tutorials/109-template-release-card-flow.png`
- 会生成 `artifacts/tutorials/109-template-release-artifact-pack-flow.html`
- 会生成 `artifacts/tutorials/109-template-release-artifact-pack-flow.json`

## 下一节

下一节把 `101-109` 真正收成一条完整的模板发布稳定性 round trip。
[Lesson 110](110-template-release-robustness-round-trip.md)
