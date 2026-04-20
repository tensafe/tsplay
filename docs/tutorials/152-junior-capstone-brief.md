# Lesson 152: 设计一个初级结业题

`Lesson 151` 先把新手结业题收成了“会找资源、会跑第一条 Flow”。  
这一节继续往前走：初级阶段要检验的是“能不能把几个基础能力组合起来”。

目标：

- 生成一份初级结业题 brief
- 把本地 demo、文件输入输出、控制流和 artifact 串起来
- 让初级阶段有一个比新手更接近业务的小项目

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/152_junior_capstone_brief.flow.yaml](../../script/tutorials/152_junior_capstone_brief.flow.yaml)
- Checklist:
  [../../script/tutorials/capstone_pack/checklists/152_junior_capstone_checklist.md](../../script/tutorials/capstone_pack/checklists/152_junior_capstone_checklist.md)

建议一起看：

- [../training/capstone-briefs.md](../training/capstone-briefs.md)
- [../training/assessment.md](../training/assessment.md)

## Step 1: 先生成初级结业题 brief

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/152_junior_capstone_brief.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/152_junior_capstone_brief.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/152/junior-capstone-brief.json`

## Step 2: 再看它和新手结业题的差别

这一层不再只要求“跑起来”，而是开始要求：

- 本地 demo
- 文件动作
- 控制流
- artifact 或 ledger

## Step 3: 这一节的最小结论

初级结业题的重点，是把多个基础动作收成一个清晰、可检查的小流程。

## 下一步

继续看：
[Lesson 153](153-intermediate-capstone-brief.md)
