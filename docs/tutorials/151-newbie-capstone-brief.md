# Lesson 151: 设计一个新手结业题

`Lesson 141-150` 已经把“单二进制 + 内置教程”的交付路径收清楚了。  
这一节继续往下接：既然新手已经有了稳定入口，那结业题应该怎么设计。

目标：

- 生成一份新手结业题 brief
- 把“资源发现 -> 资源释放 -> 最小成功运行”收成一条考核链
- 让新手阶段有一个自然的收口点

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/151_newbie_capstone_brief.flow.yaml](../../script/tutorials/151_newbie_capstone_brief.flow.yaml)
- Checklist:
  [../../script/tutorials/capstone_pack/checklists/151_newbie_capstone_checklist.md](../../script/tutorials/capstone_pack/checklists/151_newbie_capstone_checklist.md)

建议一起看：

- [../training/capstone-briefs.md](../training/capstone-briefs.md)
- [../training/assessment.md](../training/assessment.md)

## Step 1: 先生成新手结业题 brief

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/151_newbie_capstone_brief.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/151_newbie_capstone_brief.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/151/newbie-capstone-brief.json`

## Step 2: 再看它为什么只收前面这一小段能力

新手结业题重点不是“做复杂页面自动化”，而是证明：

- 会找资源
- 会释放资源
- 会跑第一条 Flow
- 会留下第一份产物

## Step 3: 这一节的最小结论

新手结业题应该检验“能不能独立迈出第一步”，而不是提前把初级、中级内容全塞进来。

## 下一步

继续看：
[Lesson 152](152-junior-capstone-brief.md)
