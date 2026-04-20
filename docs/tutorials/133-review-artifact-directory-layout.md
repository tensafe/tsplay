# Lesson 133: 把 artifact 目录也整理成可 review 的结构

`Lesson 131-132` 先把 Flow 名字和变量名字理顺了。  
这一节继续往外层走：artifact 目录也要能被 review。

目标：

- 对比“能写文件”和“可维护的输出布局”之间的区别
- 先建立 lesson 级目录习惯
- 让输出路径本身就能说明用途

## 准备工作

样例文件：

- before:
  [../../script/tutorials/133_review_artifact_layout_before.flow.yaml](../../script/tutorials/133_review_artifact_layout_before.flow.yaml)
- after:
  [../../script/tutorials/133_review_artifact_layout_after.flow.yaml](../../script/tutorials/133_review_artifact_layout_after.flow.yaml)

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

## Step 1: 先跑一个“平铺输出”的 before 版本

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/133_review_artifact_layout_before.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/133_review_artifact_layout_before.flow.yaml
```

预期结果：

- 会成功运行
- 会写出 `artifacts/output.json`

这个输出不是错，但它的问题是：

- lesson 编号不明显
- 后面再加别的教程时，容易混在一起

## Step 2: 再跑一个“lesson 级目录”的 after 版本

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/133_review_artifact_layout_after.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/133_review_artifact_layout_after.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/133/review-layout/output.json`
- 还会写出 `artifacts/tutorials/133/review-layout/manifest.json`

## Step 3: 这一节的最小结论

对于教程型输出，优先用这种结构：

- `artifacts/tutorials/<lesson>/...`

这样做的好处是：

- lesson 粒度更清楚
- review 和复跑更容易定位
- 后面补 manifest、checklist、截图时也更好放

## 下一步

继续看：
[Lesson 134](134-review-example-with-checklist.md)
