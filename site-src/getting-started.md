---
hide:
  - toc
  - title
---

<div class="tsplay-landing" markdown="1">

<section class="tsplay-cosmos" markdown="1">

<p class="tsplay-cosmos__kicker">Quick Start</p>

# 一条命令就能跑通第一次体验

<p class="tsplay-cosmos__lede">这里默认走 <strong>bin-first</strong> 路线，不要求你先装 Go，也不要求你先等 Playwright 下载完。目标很简单：自动识别系统、下载对应二进制、生成一条最小 demo Flow，并在 <code>artifacts/quickstart/</code> 下留下结构化结果。</p>

<div class="tsplay-pill-row">
<span class="tsplay-pill tsplay-pill--accent">一步到位体验</span>
<span class="tsplay-pill">自动识别当前系统</span>
<span class="tsplay-pill tsplay-pill--mint">不先下载 Playwright</span>
<span class="tsplay-pill">先拿到可见结果</span>
</div>

```bash
curl -fsSL https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.sh | sh
```

<p class="tsplay-terminal__note">如果你需要 Windows，往下切到对应平台命令即可；如果你想把二进制装到固定目录，也保留了手动安装写法。</p>

=== "macOS / Linux"

    ```bash
    curl -fsSL https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.sh | sh
    ```

    如果你想把二进制装到指定目录：

    ```bash
    curl -fsSLO https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.sh
    sh ./tsplay-quickstart.sh --install-dir ./bin
    ```

=== "Windows"

    ```powershell
    Invoke-WebRequest https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.ps1 -OutFile tsplay-quickstart.ps1
    powershell -ExecutionPolicy Bypass -File .\tsplay-quickstart.ps1
    ```

    如果你想把二进制装到指定目录：

    ```powershell
    powershell -ExecutionPolicy Bypass -File .\tsplay-quickstart.ps1 -InstallDir .\bin
    ```

<div class="grid cards" markdown="1">

-   __这一步会自动帮你完成什么__

    生成 `artifacts/quickstart/quickstart-demo.flow.yaml`，立刻执行它，并写出 `quickstart-demo-output.json`。

-   __这一步暂时不会做什么__

    不会先下载 Playwright，不会先打开浏览器，也不要求你先准备 Go 开发环境。

-   __如果你要手动挑平台__

    仍然可以直接下载 release 二进制，适合内网分发、培训包或手工部署场景。

</div>

</section>

### 如果你需要手动选择平台二进制

| 平台 | 二进制 |
| --- | --- |
| macOS Apple Silicon | [tsplay-darwin-arm64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-arm64) |
| macOS Intel | [tsplay-darwin-amd64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-amd64) |
| Linux x86_64 | [tsplay-linux-amd64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-amd64) |
| Linux ARM64 | [tsplay-linux-arm64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-arm64) |
| Windows x86_64 | [tsplay-windows-amd64.exe](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-windows-amd64.exe) |
| Windows ARM64 | [tsplay-windows-arm64.exe](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-windows-arm64.exe) |

如果机器不适合首次运行时再下载 Playwright，可以在 GitHub Release 资源里选择匹配平台的 `playwright-offline` 压缩包。这个包会同时包含 `tsplay`、`playwright/driver` 和 `playwright/browsers`；解压后保持 `playwright/` 目录和二进制在同一层即可。

<p class="tsplay-section-kicker">Next Move</p>

## 跑通后下一步做什么

<div class="grid cards" markdown="1">

-   __立刻切到页面自动化__

    先起本地 demo 页，再跑第一条浏览器 Flow。第一次真的需要浏览器时，TSPlay 才会自动下载 Playwright。

    ```bash
    ./tsplay -action file-srv -addr :8000
    ./tsplay -flow script/tutorials/10_assert_page_state.flow.yaml
    ```

    页面入口：`http://127.0.0.1:8000/demo/demo.html`、`http://127.0.0.1:8000/demo/extract.html`

-   __直接切到源码模式__

    如果你已经在仓库里，或者本来就准备参与开发，再切到源码模式会更顺手。

    ```bash
    go mod download
    go run . -flow script/tutorials/01_hello_world.flow.yaml
    go run . -action file-srv -addr :8000
    ```

    基础条件：Go `1.23.6+`，以及一台能启动 Chromium 的机器。

-   __继续保持二进制模式__

    适合下载即用、培训包、单二进制交付和内置资源验证。

    ```bash
    ./tsplay -action extract-assets -extract-root ./tsplay-assets
    ./tsplay -action file-srv -addr :8000
    ```

-   __先走 MCP / Agent 路线__

    如果你更关心工具链和能力面，可以先把 MCP 服务起起来，再看工具清单和 `finalize_flow` 主链。

    ```bash
    go run . -action srv
    go run . -action mcp-tool -tool tsplay.list_actions
    ```

</div>

## 首次运行完成后的检查点

满足下面 2 到 3 条，就算环境已经准备好了：

<div class="grid cards" markdown="1">

-   __可以视为已准备好__

    - `quickstart-demo` 或 `go run . -flow ...` 能正常结束
    - 命令输出里能看到结构化执行结果
    - `artifacts/quickstart/` 或 `artifacts/` 下出现运行产物
    - 如果你启了 `file-srv`，本地 demo 页面能打开

</div>

## 后续路径

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

## 下一步

<div class="grid cards" markdown>

-   :material-play-circle-outline:{ .lg .middle } __基础动作__

    看 [Lesson 01](docs/tutorials/01-hello-world.md) 和 [教程总站](docs/tutorials/README.zh-CN.md)。

-   :material-file-document-edit-outline:{ .lg .middle } __Flow 编写__

    看 [基础学习路线](docs/tutorials/track-newbie.zh-CN.md) 和 [学习路径](docs/training/learning-path.md)。

-   :material-robot-love-outline:{ .lg .middle } __Agent 集成__

    看 [AI 协作入门](docs/training/ai-intent-to-flow.md) 和 [MCP 教程入口](docs/tutorials/119-mcp-chain-overview.md)。

-   :material-package-variant-closed:{ .lg .middle } __单二进制交付__

    看 [Lesson 142](docs/tutorials/142-list-assets-for-beginners.md)、[Lesson 144](docs/tutorials/144-single-binary-delivery-flow.md) 和 [核心功能路线图](docs/product/core-feature-roadmap.md)。

</div>

## 常见问题

首次运行后，下面这 4 个问题最常见：

<div class="grid cards" markdown="1">

-   __已经能跑 `Lesson 01`，但不知道下一步看哪条线__

    默认继续去 [教程总站](docs/tutorials/README.zh-CN.md) 或 [基础学习路线](docs/tutorials/track-newbie.zh-CN.md)。

-   __想直接让 Agent 出 Flow，却还没搞懂默认 MCP 主路径__

    先看 [AI 协作入门](docs/training/ai-intent-to-flow.md)，再看 [Lesson 120](docs/tutorials/120-mcp-finalize-flow.md)。

-   __构建了二进制，但不知道 `list-assets`、`extract-assets`、`file-srv` 该先用哪个__

    先看 [Lesson 142](docs/tutorials/142-list-assets-for-beginners.md)、[Lesson 143](docs/tutorials/143-extract-assets-for-beginners.md)、[Lesson 147](docs/tutorials/147-file-srv-dev-vs-release.md)。

-   __看到 `artifacts/` 有文件，但不知道哪些才是交付证据__

    先看 [Lesson 87](docs/tutorials/87-build-handoff-artifact-manifest.md) 和 [Lesson 88](docs/tutorials/88-build-handoff-summary.md)。

</div>

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

<div class="grid cards" markdown="1">

-   __项目总览（中文）__

    [打开项目总览（中文）](README.zh-CN.md)

-   __Project Overview__

    [Open Project Overview](ReadMe.md)

-   __文档入口__

    [打开文档入口](docs/README.md)

-   __文档健康检查__

    [打开文档健康检查](docs/doc-health-audit.md)

-   __核心功能执行面板__

    [打开核心功能执行面板](docs/product/core-feature-execution-board.md)

</div>

</div>
