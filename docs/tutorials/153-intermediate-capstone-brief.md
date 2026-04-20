# Lesson 153: 设计一个中级结业题

`Lesson 152` 已经把初级阶段收成了“本地 demo + 文件 + 控制流”的组合题。  
这一节继续进入中级：开始检验 MCP 链路和失败修复能力。

目标：

- 生成一份中级结业题 brief
- 把 `observe -> draft -> validate -> run -> repair` 收成一条考核链
- 让中级阶段有一个真正体现“可复用”和“可修复”的项目

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/153_intermediate_capstone_brief.flow.yaml](../../script/tutorials/153_intermediate_capstone_brief.flow.yaml)
- Checklist:
  [../../script/tutorials/capstone_pack/checklists/153_intermediate_capstone_checklist.md](../../script/tutorials/capstone_pack/checklists/153_intermediate_capstone_checklist.md)

建议一起看：

- [../training/capstone-briefs.md](../training/capstone-briefs.md)
- [../training/assessment.md](../training/assessment.md)

## Step 1: 先生成中级结业题 brief

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/153_intermediate_capstone_brief.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/153_intermediate_capstone_brief.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/153/intermediate-capstone-brief.json`

## Step 2: 再看这道题到底在检验什么

它不是只检验“会不会用 MCP”，而是检验：

- 会不会从 observation 收敛到 Flow
- 会不会从失败收敛到 repair
- 会不会留证据给别人复盘

## Step 3: 这一节的最小结论

中级结业题的重点，是让学员把“草拟、执行、失败、修复”当成一条完整闭环，而不是只做一次成功演示。

## 下一步

继续看：
[Lesson 154](154-advanced-capstone-brief.md)
