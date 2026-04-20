# Lesson 155: 设计一套“新人 7 天计划”

`Lesson 151-154` 已经把四层结业题都收出来了。  
这一节继续往培训排期走：如果一个新人要在 7 天内建立稳定起步，节奏应该怎么排。

目标：

- 生成一份 7 天计划
- 把 lesson、capstone 和 review 串成日程
- 让 onboarding 也保持同样的循序渐进

## 准备工作

样例文件：

- [../../script/tutorials/155_new_hire_7_day_plan.flow.yaml](../../script/tutorials/155_new_hire_7_day_plan.flow.yaml)

建议一起看：

- [../training/README.md](../training/README.md)
- [../training/bootcamp-plan.md](../training/bootcamp-plan.md)

## Step 1: 先生成 7 天计划 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/155_new_hire_7_day_plan.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/155_new_hire_7_day_plan.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/155/new-hire-7-day-plan.json`

## Step 2: 再看这 7 天为什么这样排

它不是把所有主题平均切成 7 份，而是按难度渐进：

- 先资源和第一条 Flow
- 再文件和本地 demo
- 再控制流和恢复
- 最后才进入边界和结业题

## 下一步

继续看：
[Lesson 156](156-implementer-2-week-plan.md)
