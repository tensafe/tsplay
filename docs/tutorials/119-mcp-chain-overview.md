# Lesson 119: 把 `observe -> draft -> validate -> run -> repair` 串成一条线

`Lesson 113-118` 不是 6 个彼此独立的小功能。  
这一节专门把它们重新串起来，帮助你建立 MCP 基础链路的整体感。

目标：

- 理清每一步的输入和输出
- 理清成功链和失败链
- 为 `finalize_flow` 做准备

## Step 1: 回看这一段已经产出的文件

重点看这些文件：

- `artifacts/tutorials/113-mcp-observe-page-template-release.json`
- `artifacts/tutorials/114-mcp-draft-flow-template-release.json`
- `artifacts/tutorials/115-mcp-validate-drafted-template-release.json`
- `artifacts/tutorials/116-mcp-run-drafted-template-release.json`
- `artifacts/tutorials/117-mcp-run-broken-template-release.json`
- `artifacts/tutorials/117-mcp-repair-flow-context-template-release.json`
- `artifacts/tutorials/118-mcp-repair-flow-template-release.json`

## Step 2: 这一段真正的顺序

成功链：

- `Lesson 113` 用 `observe_page` 拿 observation
- `Lesson 114` 用 observation 生成 `draft.flow_yaml`
- `Lesson 115` 先校验这份草稿
- `Lesson 116` 再执行这份草稿

失败链：

- `Lesson 117` 先故意跑坏一份 Flow
- `Lesson 117` 再把失败现场整理成 `repair context`
- `Lesson 118` 把 `repair context` 变成统一 repair request

## Step 3: 为什么还要有 `Lesson 120`

前面这条链的重点是“拆开理解”。  
但真实使用时，你不一定每次都想自己手动走完 `observe -> draft -> validate`。

所以最后一节会把这条链收成一个更短的默认入口：

- `tsplay.finalize_flow`

## 下一步

继续看：
[Lesson 120](120-mcp-finalize-flow.md)
