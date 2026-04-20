# Lesson 157: 设计一套“讲师备课顺序”

`Lesson 156` 已经把实施同学的 2 周路径排出来了。  
这一节继续往 Enablement 角色走：讲师应该先读什么、先排什么、先准备什么。

目标：

- 生成一份讲师备课顺序
- 把教程主线和 `docs/training/` 正式接起来
- 让培训材料不再和教程主线分家

## 准备工作

样例文件：

- [../../script/tutorials/157_trainer_prep_sequence.flow.yaml](../../script/tutorials/157_trainer_prep_sequence.flow.yaml)

建议一起看：

- [../training/README.md](../training/README.md)
- [../training/trainer-playbook.md](../training/trainer-playbook.md)

## Step 1: 先生成讲师备课顺序 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/157_trainer_prep_sequence.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/157_trainer_prep_sequence.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/157/trainer-prep-sequence.json`

## Step 2: 再看为什么这条顺序是从教程走到培训

顺序不是随机的，而是：

- 先全局入口
- 再教程主线
- 再训练营材料
- 再评分与评审
- 最后才到讲师手册

## 下一步

继续看：
[Lesson 158](158-tutorial-gap-review-cycle.md)
