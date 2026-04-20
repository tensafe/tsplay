# Lesson 156: 设计一套“实施同学 2 周计划”

`Lesson 155` 已经把新人 7 天计划收出来了。  
这一节继续往更贴近交付的角色走：实施同学的 2 周计划。

目标：

- 生成一份 2 周计划
- 把教程熟练度和真实业务试点连起来
- 让“学完教程”自然接到“开始交付”

## 准备工作

样例文件：

- [../../script/tutorials/156_implementer_2_week_plan.flow.yaml](../../script/tutorials/156_implementer_2_week_plan.flow.yaml)

建议一起看：

- [../training/README.md](../training/README.md)
- [../training/bootcamp-plan.md](../training/bootcamp-plan.md)

## Step 1: 先生成 2 周计划 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/156_implementer_2_week_plan.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/156_implementer_2_week_plan.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/156/implementer-2-week-plan.json`

## Step 2: 再看这份计划的重点

这一层已经不只是“学会教程”，而是：

- 先练结构和 review
- 再试真实业务
- 再收交付和 handoff

## 下一步

继续看：
[Lesson 157](157-trainer-prep-sequence.md)
