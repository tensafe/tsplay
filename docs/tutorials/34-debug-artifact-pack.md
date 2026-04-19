# Lesson 34: 生成一份调试产物包

前面三节分别留了：

- 整页截图
- 元素截图
- HTML

这一节把它们合成一份更像真实交付物的调试产物包。

使用页面：
[../../demo/debug_artifacts.html](../../demo/debug_artifacts.html)

目标：

- 在一条脚本里生成多种证据
- 用 JSON 记录所有产物路径
- 建立“证据包”思维，而不是单文件思维

## 准备工作

先确认本地静态文件服务还在运行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/34_debug_artifact_pack.lua](../../script/tutorials/34_debug_artifact_pack.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/34_debug_artifact_pack.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/34_debug_artifact_pack.lua
```

预期结果：

- 会生成整页截图
- 会生成元素截图
- 会生成 HTML
- 会生成 `artifacts/tutorials/34-debug-artifact-pack-lua.json`

## Step 2: 这一节最重要的变化

不是多学了几个动作，  
而是开始把“证据本身”当成正式输出。

这在真实项目里很重要，因为：

- 你要把问题发给别人看
- 你要让下次复盘还能找到现场
- 你要让 Flow 的失败不只剩一条报错字符串

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/34_debug_artifact_pack.flow.yaml](../../script/tutorials/34_debug_artifact_pack.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/34_debug_artifact_pack.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/34_debug_artifact_pack.flow.yaml -headless
```

预期结果：

- 会生成整页截图
- 会生成元素截图
- 会生成 HTML
- 会生成 `artifacts/tutorials/34-debug-artifact-pack-flow.json`

## 下一节

下一节把证据包放进一个“失败场景”里，做成真正的错误现场保留。
[Lesson 35](35-error-evidence-pack.md)
