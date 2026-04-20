# Lesson 160: 通过高级阶段检查，并为下一圈教程扩展做准备

`Lesson 151-159` 已经把：

- 结业题
- 学习计划
- 讲师备课
- 缺口复盘
- 每 10 次回看

连成了一条高级阶段的收口主线。  
这一节把它收成一个 continuation plan。

目标：

- 生成一份继续扩展课程的计划 JSON
- 把高级阶段的退出标准写得可检查
- 让 `160` 不是终点，而是下一圈的起点

## 准备工作

样例文件：

- [../../script/tutorials/160_curriculum_continuation_plan.flow.yaml](../../script/tutorials/160_curriculum_continuation_plan.flow.yaml)
- [../../script/tutorials/capstone_pack/README.md](../../script/tutorials/capstone_pack/README.md)

建议一起看：

- [iteration-roadmap-160.md](iteration-roadmap-160.md)
- [evolution-playbook.md](evolution-playbook.md)

## Step 1: 先生成 continuation plan JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/160_curriculum_continuation_plan.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/160_curriculum_continuation_plan.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/160/curriculum-continuation-plan.json`

## Step 2: 再看 `160` 为什么不是终点

这一节真正要证明的是：

- 你能继续补前置认知
- 你能继续补连接点
- 你能继续补 training / review / release 素材
- 你不会把课程补成一条越来越乱的长链

## Step 3: 这一节的最小结论

高级阶段通过，不是因为“写到了第 160 节”。  
而是因为你已经能继续沿同一逻辑，再扩下一圈而不失去结构。

## 下一步

回到：
[iteration-roadmap-160.md](iteration-roadmap-160.md)
