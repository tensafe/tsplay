# 新手教程

[English](track-newbie.md) | 简体中文

这条线只做一件事：让一个从没用过 TSPlay 的人，按顺序拿到第一个、第二个、第三个成功体验。

这里强调的是：

- 不跳跃
- 不抢跑
- 不先讲抽象词
- 每学一个动作，就立刻看到结果

## 适合谁

- 第一天接触 TSPlay
- 只知道“它是自动化工具”，还不知道该从哪里下手
- 对 `Lua`、`Flow`、`MCP`、`artifact` 这些词还没有稳定概念

## 这一层的学习目标

学完之后，应该能独立回答这些问题：

- `tsplay` 怎么构建、怎么直接运行
- 为什么单个二进制也能带着 `docs/`、`script/`、`demo/`
- `Lua` 和 `Flow` 的边界分别是什么
- 为什么要先学本地页面，再学外部系统
- 为什么结果一定要落到 `artifacts/`

## 推荐顺序

1. [Lesson 01](01-hello-world.md)
2. [Lesson 08](08-bundled-assets-and-artifacts.md)
3. [Lesson 02](02-local-page-select-option.md)
4. [Lesson 09](09-local-demo-anatomy.md)
5. [Lesson 03](03-capture-table.md)
6. [Lesson 04](04-extract-text-and-html.md)
7. [Lesson 05](05-http-request-json.md)
8. [Lesson 10](10-assert-page-state.md)
9. [Lesson 11](11-select-another-option.md)
10. [Lesson 12](12-custom-json-output.md)

## 新手阶段的核心动作

- `set_var`
- `write_json`
- `navigate`
- `wait_for_selector`
- `select_option`
- `capture_table`
- `extract_text`
- `get_html`
- `http_request`
- `json_extract`
- `assert_visible`
- `assert_text`
- `is_selected`

## 新手阶段的交付物

这一层不要求做复杂业务。  
只要求每学一个主题，都至少交出一种可运行结果：

- 一个 `Lua` 脚本
- 一条 `Flow`
- 一份写入 `artifacts/tutorials/` 的 JSON
- 一段你自己的复盘说明

## 新手阶段的常见误区

### 误区 1：一开始就想学最完整的业务场景

结果往往是页面、系统、权限、环境一起卡住。  
新手阶段应该先把“本地最小成功体验”跑出来。

### 误区 2：一开始就只学 Flow，不学 Lua

Flow 是主线，但对新手来说，Lua 更像一条直线。  
先知道“怎么做”，再知道“怎么结构化表达”，心智负担会小很多。

### 误区 3：只看文档，不跑命令

TSPlay 不是只看概念就能真正学会的工具。  
新手阶段必须坚持“每一小步都跑一下”。

## 新手阶段的退出标准

满足下面这些条件时，就可以进入初级教程：

- 能独立构建 `./tsplay`
- 能独立启动 `./tsplay -action file-srv -addr :8000`
- 能独立运行新手本地练习链路：`Lesson 01-05` 和 `Lesson 08-12`
- 能说清楚 `Lua` 与 `Flow` 的差异
- 能把页面文本、HTML 片段、表格、JSON 请求结果写到文件里

## 学完之后去哪里

下一站是：
[track-junior.zh-CN.md](track-junior.zh-CN.md)
