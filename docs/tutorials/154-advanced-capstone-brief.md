# Lesson 154: 设计一个高级结业题

`Lesson 153` 已经进入了 MCP 和修复闭环。  
这一节继续往高级阶段收口：高级结业题要开始检验安全边界、review 和交付策略。

目标：

- 生成一份高级结业题 brief
- 把安全边界、review 规则和发布包思路收成一个团队级场景
- 让高级阶段的结业题更像“交付方案”而不是“脚本练习”

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/154_advanced_capstone_brief.flow.yaml](../../script/tutorials/154_advanced_capstone_brief.flow.yaml)
- Checklist:
  [../../script/tutorials/capstone_pack/checklists/154_advanced_capstone_checklist.md](../../script/tutorials/capstone_pack/checklists/154_advanced_capstone_checklist.md)

建议一起看：

- [../training/capstone-briefs.md](../training/capstone-briefs.md)
- [../training/assessment.md](../training/assessment.md)

## Step 1: 先生成高级结业题 brief

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/154_advanced_capstone_brief.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/154_advanced_capstone_brief.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/154/advanced-capstone-brief.json`

## Step 2: 再看它和中级题的差别

高级结业题开始要求的不只是“做出来”，而是：

- 为什么这样放权
- 为什么这样 review
- 为什么这样交付
- 新同学怎么接手

## Step 3: 这一节的最小结论

高级结业题的核心，是把 TSPlay 从“会写 Flow”推进到“会组织团队交付”。

## 下一步

继续看：
[Lesson 155](155-new-hire-7-day-plan.md)
