# Lesson 158: 设计一套“教程缺口复盘机制”

`Lesson 157` 已经把讲师备课顺序收出来了。  
这一节继续往长期演进走：教程体系怎么月度复盘、怎么找缺口。

目标：

- 生成一份月度复盘模板
- 把 learner blocker、bridge topic、demo 更新都纳入同一个循环
- 让教程增长有节奏，而不是想到什么补什么

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/158_tutorial_gap_review_cycle.flow.yaml](../../script/tutorials/158_tutorial_gap_review_cycle.flow.yaml)
- Template:
  [../../script/tutorials/capstone_pack/checklists/158_monthly_gap_review_template.md](../../script/tutorials/capstone_pack/checklists/158_monthly_gap_review_template.md)

建议一起看：

- [evolution-playbook.md](evolution-playbook.md)
- [../training/assessment.md](../training/assessment.md)

## Step 1: 先生成月度复盘模板 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/158_tutorial_gap_review_cycle.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/158_tutorial_gap_review_cycle.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/158/tutorial-gap-review-cycle.json`

## Step 2: 再看这份复盘为什么要同时看“学员卡点”和“素材脱节”

因为教程缺口通常不只一种：

- 可能是前置认知缺了
- 可能是连接点没讲
- 也可能是 demo 或脚本已经和文档脱节

## 下一步

继续看：
[Lesson 159](159-every-10-iteration-review.md)
