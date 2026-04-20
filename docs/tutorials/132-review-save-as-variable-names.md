# Lesson 132: 继续把 `save_as` 变量名写得可 review

`Lesson 131` 先把 `name` 和 `description` 写清楚了。  
这一节继续往里走一步：变量名也要能被 review。

目标：

- 看懂为什么 `save_as: x` 会提高维护成本
- 建立一套最小变量命名标准
- 保持“看输出就能猜到数据角色”

## 准备工作

这节直接复用上一节的两条 Flow：

- [../../script/tutorials/131_review_readability_before.flow.yaml](../../script/tutorials/131_review_readability_before.flow.yaml)
- [../../script/tutorials/131_review_readability_after.flow.yaml](../../script/tutorials/131_review_readability_after.flow.yaml)

如果你还没跑过，可以直接执行：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/131_review_readability_before.flow.yaml
go run . -flow script/tutorials/131_review_readability_after.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/131_review_readability_before.flow.yaml
./tsplay -flow script/tutorials/131_review_readability_after.flow.yaml
```

## Step 1: 先看 before Flow 的变量名

before 版本里最关键的两个名字是：

- `x`
- `y`

问题不是它会不会跑，而是：

- `x` 到底是 payload、rows、summary，还是别的东西
- `y` 到底是写文件结果，还是别的结果

## Step 2: 再看 after Flow 的变量名

after 版本里，对应的是：

- `review_payload`
- `write_result`

这两个名字的价值在于：

- 不用打开实现细节，也能大概猜到变量角色
- trace、artifact、review 评论都更容易写清楚

## Step 3: 这一节的最小规则

优先用“角色名”，不要用“临时名”：

- 好：`review_payload`、`write_result`
- 一般：`result`
- 差：`x`、`tmp2`

## 下一步

继续看：
[Lesson 133](133-review-artifact-directory-layout.md)
