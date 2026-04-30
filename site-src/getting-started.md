---
hide:
  - toc
---

# 快速开始

这页先陪你把 TSPlay 跑通，再帮你顺着目标找到更合适的下一站。

## 一步到位默认路径

如果你现在只是想确认 TSPlay 能不能跑，不要先在 `Lua / Flow / MCP / 二进制` 之间做选择。
也不需要先装 Go，或者先等 Playwright 下载完。

先走这条最短路径：

- `curl` 直接下载最新 release 二进制
- 自动生成一条最小 demo Flow
- 立刻执行它
- 在 `artifacts/quickstart/` 下留下结果

=== "macOS Apple 芯片"

    ```bash
    curl -L -o tsplay https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-arm64
    chmod +x tsplay
    ./tsplay -action quickstart-demo
    ```

=== "macOS Intel"

    ```bash
    curl -L -o tsplay https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-amd64
    chmod +x tsplay
    ./tsplay -action quickstart-demo
    ```

=== "Linux x86_64"

    ```bash
    curl -L -o tsplay https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-amd64
    chmod +x tsplay
    ./tsplay -action quickstart-demo
    ```

=== "Linux ARM64"

    ```bash
    curl -L -o tsplay https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-arm64
    chmod +x tsplay
    ./tsplay -action quickstart-demo
    ```

=== "Windows x86_64"

    ```powershell
    curl.exe -L -o tsplay.exe https://github.com/tensafe/tsplay/releases/latest/download/tsplay-windows-amd64.exe
    .\tsplay.exe -action quickstart-demo
    ```

### 这一步会帮你完成什么

- 生成 `artifacts/quickstart/quickstart-demo.flow.yaml`
- 立刻执行这条 Flow
- 写出 `artifacts/quickstart/quickstart-demo-output.json`
- 在终端打印结构化执行结果

### 这一步暂时不会做什么

- 不会先下载 Playwright
- 不会先打开浏览器
- 不要求你先准备 Go 开发环境

## 如果你想立刻切到页面自动化

`quickstart-demo` 跑通后，再多走两步就能切到真正的页面动作。

### 1. 先起本地 demo 页面

```bash
./tsplay -action file-srv -addr :8000
```

然后访问：

- `http://127.0.0.1:8000/demo/demo.html`
- `http://127.0.0.1:8000/demo/extract.html`

### 2. 再跑第一条浏览器 Flow

```bash
./tsplay -flow script/tutorials/10_assert_page_state.flow.yaml
```

这一步第一次如果需要浏览器，TSPlay 才会自动下载 Playwright。

## 如果你更想从源码开始

如果你已经在仓库里，或者本来就准备参与开发，再切到源码模式会更顺手。

### 1. 先准备这两个基础条件

- Go `1.23.6+`
- 一台能启动 Chromium 的机器

### 2. 在仓库根目录安装依赖

```bash
go mod download
```

### 3. 先跑第一条最小 Flow

```bash
go run . -flow script/tutorials/01_hello_world.flow.yaml
```

### 4. 如果你想练本地 demo 页面，再开一个终端

```bash
go run . -action file-srv -addr :8000
```

## 跑通后的判断标准

满足下面 2 到 3 条，就算环境已经准备好了：

- `quickstart-demo` 或 `go run . -flow ...` 能正常结束
- 命令输出里能看到结构化执行结果
- `artifacts/quickstart/` 或 `artifacts/` 下出现运行产物
- 如果你启了 `file-srv`，本地 demo 页面能打开

## 跑通后的常见分支

默认路径跑通后，再按你现在最想继续的方向往下走。

=== "继续用二进制模式"

    ```bash
    ./tsplay -action extract-assets -extract-root ./tsplay-assets
    ./tsplay -action file-srv -addr :8000
    ```

    适合：

    - 你想继续保持“下载即用”的方式
    - 想把内置 docs、demo、script 释放到本地
    - 后面准备做单二进制交付或培训包

=== "继续用源码模式"

    ```bash
    go run . -flow script/demo_baidu.flow.yaml
    ```

    适合：

    - 你已经跑通默认路径
    - 想继续在仓库里改代码或补教程

=== "先构建二进制"

    ```bash
    go build -o tsplay .
    ./tsplay -flow script/tutorials/01_hello_world.flow.yaml
    ./tsplay -action list-assets
    ```

    适合：

    - 你想把 `./tsplay` 当作主入口
    - 想验证内置文档、示例和 demo 资源是否一起打包
    - 后面准备做单二进制分发

=== "先走 MCP / Agent"

    ```bash
    go run . -action srv
    go run . -action mcp-tool -tool tsplay.list_actions
    ```

    适合：

    - 你更关心 Agent 集成
    - 想先看工具链和能力面
    - 后面会走 `finalize_flow`、repair 和命名会话路线

## 跑通后下一步做什么

<div class="grid cards" markdown>

-   :material-play-circle-outline:{ .lg .middle } __我想继续熟悉基础动作__

    看 [Lesson 01](docs/tutorials/01-hello-world.md) 和 [教程总站](docs/tutorials/README.zh-CN.md)。

-   :material-file-document-edit-outline:{ .lg .middle } __我想尽快学写 Flow__

    看 [新手学习路线](docs/tutorials/track-newbie.zh-CN.md) 和 [学习路径](docs/training/learning-path.md)。

-   :material-robot-love-outline:{ .lg .middle } __我想先走 Agent 路线__

    看 [AI 无感入门](docs/training/ai-intent-to-flow.md) 和 [MCP 教程入口](docs/tutorials/119-mcp-chain-overview.md)。

-   :material-package-variant-closed:{ .lg .middle } __我想先看单二进制交付__

    看 [Lesson 142](docs/tutorials/142-list-assets-for-beginners.md)、[Lesson 144](docs/tutorials/144-single-binary-delivery-flow.md) 和 [核心功能路线图](docs/product/core-feature-roadmap.md)。

</div>

## 几个常见卡点

第一次跑通后，下面这 4 个地方最容易让人停下来犹豫一下：

- 已经能跑 `Lesson 01`，但不知道下一步看哪条线。
  默认继续去 [教程总站](docs/tutorials/README.zh-CN.md) 或 [新手学习路线](docs/tutorials/track-newbie.zh-CN.md)。
- 想直接让 Agent 出 Flow，却还没搞懂默认 MCP 主路径。
  先看 [AI 无感入门](docs/training/ai-intent-to-flow.md)，再看 [Lesson 120](docs/tutorials/120-mcp-finalize-flow.md)。
- 构建了二进制，但不知道 `list-assets`、`extract-assets`、`file-srv` 该先用哪个。
  先看 [Lesson 142](docs/tutorials/142-list-assets-for-beginners.md)、[Lesson 143](docs/tutorials/143-extract-assets-for-beginners.md)、[Lesson 147](docs/tutorials/147-file-srv-dev-vs-release.md)。
- 看到 `artifacts/` 有文件，但不知道哪些才是交付证据。
  先看 [Lesson 87](docs/tutorials/87-build-handoff-artifact-manifest.md)、[Lesson 88](docs/tutorials/88-build-handoff-summary.md)。

## 常用命令速查

```bash
./tsplay -action quickstart-demo
./tsplay -action file-srv -addr :8000
./tsplay -flow script/tutorials/10_assert_page_state.flow.yaml
go run . -action cli
go run . -flow script/demo_baidu.flow.yaml
go run . -action srv
go run . -action mcp-tool -tool tsplay.list_actions
go run . -action extract-assets -extract-root ./tsplay-assets
```

## 如果你想继续看全貌

- [项目总览（中文）](README.zh-CN.md)
- [Project Overview](ReadMe.md)
- [文档入口](docs/README.md)
- [文档健康检查](docs/doc-health-audit.md)
- [核心功能执行面板](docs/product/core-feature-execution-board.md)
