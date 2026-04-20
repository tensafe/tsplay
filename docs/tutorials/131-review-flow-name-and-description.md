# Lesson 131: 先把 Flow 名字和描述写得可 review

`Lesson 121-130` 先把安全边界讲清楚了。  
从这一节开始，我们进入下一条高级主线：

- 不只看“能不能跑”
- 还要看“别人能不能 review、能不能接手”

目标：

- 对比一条“能跑但不可 review”的 Flow
- 对比一条“同样简单，但更适合 review”的 Flow
- 先把 `name` 和 `description` 这两个入口写清楚

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

样例文件：

- before:
  [../../script/tutorials/131_review_readability_before.flow.yaml](../../script/tutorials/131_review_readability_before.flow.yaml)
- after:
  [../../script/tutorials/131_review_readability_after.flow.yaml](../../script/tutorials/131_review_readability_after.flow.yaml)

## Step 1: 先跑一条“能跑但不可 review”的 before Flow

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/131_review_readability_before.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/131_review_readability_before.flow.yaml
```

预期结果：

- 会成功运行
- 会写出 `artifacts/review-output.json`
- 但你会看到 `name=tmp_review`
- `description` 也只是很泛的 “Do the thing and write a file.”

## Step 2: 再跑一条更适合 review 的 after Flow

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/131_review_readability_after.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/131_review_readability_after.flow.yaml
```

预期结果：

- 会成功运行
- 会写出 `artifacts/tutorials/131/review-summary.json`
- `name` 会直接说明 lesson 和任务意图
- `description` 会直接说明这条 Flow 对交付者的结果是什么

## Step 3: 这一节到底在 review 什么

先不要急着看实现细节，先问两个问题：

1. 只看 `name`，我能不能猜到这条 Flow 是干什么的。
2. 只看 `description`，我能不能知道它最后交付什么结果。

如果这两个入口都答不出来，那么这条 Flow 就算能跑，review 成本也会偏高。

## 下一步

继续看：
[Lesson 132](132-review-save-as-variable-names.md)
