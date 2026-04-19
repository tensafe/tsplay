# 中级教程

中级教程的关键，不是再堆更多 action，而是开始建立“可维护性”。

从这一层开始，问题会变成：

- 这条 Flow 能不能复用
- 下个月再看还认不认得
- 换一个同事接手会不会很难
- 页面变了以后修起来痛不痛

## 适合谁

- 已经能独立写基础 Flow
- 开始遇到页面波动、业务变量膨胀、流程越来越长的问题
- 希望把 TSPlay 变成团队可维护资产的人

## 中级阶段的主题

### 1. 模板化

把一次性脚本提炼成：

- 可复用脚本
- 可复用 Flow 模板
- 可复用环境示例

### 2. 数据驱动

让流程不只是处理一条数据，而是开始处理一批数据：

- CSV
- Excel
- `foreach`
- 批量结果写回

其中 CSV 和 Excel 两条最小实践线，在当前仓库里已经可以先从 `Lesson 13-15` 和 `Lesson 24-27` 跑起来。

### 3. 健壮性设计

要系统掌握：

- `retry`
- `wait_until`
- `on_error`
- `assert_visible`
- `assert_text`

其中 `retry`、`wait_until`、上传下载前后的断言链路，在当前仓库里已经可以先从
`Lesson 16-20` 打基础；`if`、`foreach`、`on_error` 的最小导入链路，则可以从 `Lesson 21-27` 接上。

### 4. MCP 基础链路

开始理解：

- `observe_page`
- `draft_flow`
- `validate_flow`
- `run_flow`
- `repair_flow_context`
- `repair_flow`

## 中级阶段的交付物

- 一个数据驱动 Flow
- 一个带失败恢复的 Flow
- 一份“为什么这样拆步骤”的说明
- 一份修复前 / 修复后的对比记录

## 中级阶段的评估重点

不是“会不会更多动作”，而是：

- 会不会把流程拆清楚
- 会不会让变量稳定
- 会不会让失败变得可观察
- 会不会给后续 repair 留空间

## 中级阶段的退出标准

- 能维护一组模板化 Flow
- 能处理一批数据而不是一条数据
- 能对 Flaky 流程给出合理的 `retry / on_error / wait_until` 设计
- 能说清楚什么时候值得引入 MCP

## 学完之后去哪里

下一站是：
[track-advanced.md](track-advanced.md)
