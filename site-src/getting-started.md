---
hide:
  - toc
---

# 快速开始

这页只做一件事：帮你尽快跑通 TSPlay，然后根据你的目标跳到正确的下一站。

## 先选一种起步方式

=== "直接跑源码"

    ```bash
    go mod download
    go run . -flow script/demo_baidu.flow.yaml
    ```

    适合：

    - 你已经在仓库里
    - 想最快验证环境和 Flow 能不能跑
    - 后面还会继续改代码或补教程

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
    go mod download
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
