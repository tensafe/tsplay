---
hide:
  - toc
---

# 快速开始

这页只做一件事：帮你尽快跑通 TSPlay，然后根据你的目标跳到正确的下一站。

## 零基础默认路径

如果你第一次接触 TSPlay，不要先在 `Lua / Flow / MCP / 二进制` 之间做选择。
先按下面这条默认路径走，跑通一次再分叉。

### 1. 先准备两个东西

- Go `1.23.6+`
- 一台能启动 Chromium 的机器

首次执行浏览器相关命令时，TSPlay 会自动调用 `playwright.Install()` 下载浏览器。

### 2. 在仓库根目录安装依赖

```bash
go mod download
```

### 3. 先跑第一条最小 Flow

```bash
go run . -flow script/tutorials/01_hello_world.flow.yaml
```

### 4. 如果你要练本地 demo 页面，再开一个终端

```bash
go run . -action file-srv -addr :8000
```

然后访问：

- `http://127.0.0.1:8000/demo/demo.html`
- `http://127.0.0.1:8000/demo/tables.html`

## 跑通算成功

满足下面 2 到 3 条，就算环境已经准备好了：

- `go run . -flow ...` 能正常结束
- 命令输出里能看到结构化执行结果
- `artifacts/` 下出现运行产物
- 如果你启了 `file-srv`，本地 demo 页面能打开

## 进阶起步方式

默认路径跑通后，再按你的目标切换下面这些入口。

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

-   :material-play-circle-outline:{ .lg .middle } __我想继续学基础动作__

    看 [Lesson 01](docs/tutorials/01-hello-world.md) 和 [教程总站](docs/tutorials/README.zh-CN.md)。

-   :material-file-document-edit-outline:{ .lg .middle } __我想尽快学写 Flow__

    看 [新手学习路线](docs/tutorials/track-newbie.zh-CN.md) 和 [学习路径](docs/training/learning-path.md)。

-   :material-robot-love-outline:{ .lg .middle } __我想让 Agent 直接帮我做__

    看 [AI 无感入门](docs/training/ai-intent-to-flow.md) 和 [MCP 教程入口](docs/tutorials/119-mcp-chain-overview.md)。

</div>

## 常用命令速查

```bash
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
