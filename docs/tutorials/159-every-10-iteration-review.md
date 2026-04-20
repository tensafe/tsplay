# Lesson 159: 设计一套“每 10 次迭代回看一次”的检查机制

`Lesson 158` 已经把月度缺口复盘机制收出来了。  
这一节继续把频率再细化：每推进 10 次迭代，就做一次结构检查。

目标：

- 生成一份 every-10-iterations review JSON
- 把“术语、命令、产物、顺序、前置条件”纳入同一个检查表
- 让教程增长时不失去主线

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/159_every_10_iteration_review.flow.yaml](../../script/tutorials/159_every_10_iteration_review.flow.yaml)
- Checklist:
  [../../script/tutorials/capstone_pack/checklists/159_every_10_iteration_checklist.md](../../script/tutorials/capstone_pack/checklists/159_every_10_iteration_checklist.md)

建议一起看：

- [evolution-playbook.md](evolution-playbook.md)
- [iteration-roadmap-160.md](iteration-roadmap-160.md)

## Step 1: 先生成 every-10-iterations review JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/159_every_10_iteration_review.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/159_every_10_iteration_review.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/159/every-10-iteration-review.json`

## Step 2: 再把它和 `Lesson 158` 区分开

`Lesson 158` 更偏月度复盘。  
这一节更偏节奏控制：

- 每推进 10 次，就强制回头看一次结构

## 下一步

继续看：
[Lesson 160](160-curriculum-continuation-plan.md)
