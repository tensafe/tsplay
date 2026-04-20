# Lesson 136: review 时，什么时候应该把 Lua 抽回 Flow

`Lesson 135` 先看了一个“能跑但不一定该保留 Lua”的例子。  
这一节把它正式抽回 Flow。

目标：

- 跑通与上一节同类的 Flow 版本
- 理解“简单编排优先 Flow”的边界
- 让 review、artifact、变量命名都更清晰

## 准备工作

样例文件：

- Lua:
  [../../script/tutorials/135_review_summary_escape_hatch.lua](../../script/tutorials/135_review_summary_escape_hatch.lua)
- Flow:
  [../../script/tutorials/136_review_summary_extracted.flow.yaml](../../script/tutorials/136_review_summary_extracted.flow.yaml)

## Step 1: 先运行 Flow 版本

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/136_review_summary_extracted.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/136_review_summary_extracted.flow.yaml
```

预期结果：

- 会成功运行
- 会写出 `artifacts/tutorials/136/review-summary-from-flow.json`

## Step 2: 再把它和上一节对照起来

你会看到：

- 行为差不多
- 但 Flow 版本的 `name`、`description`、`save_as`、artifact 路径都更适合 review

## Step 3: 这一节的最小规则

如果一段 Lua 主要在做这些事：

- 组织步骤
- 写本地文件
- 保存中间变量

那它大概率应该抽回 Flow。

## 下一步

继续看：
[Lesson 137](137-review-when-to-add-demo-page.md)
