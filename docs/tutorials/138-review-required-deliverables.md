# Lesson 138: review 时，一节教程至少该带哪些交付物

`Lesson 137` 先把“什么时候新增 demo”说清楚了。  
这一节继续往交付层收口：一节教程不能只交一段代码。

目标：

- 建立“教程交付物”的最小清单
- 避免出现“只有脚本，没有文档”的半成品
- 让后面新增 lesson 更容易沿同一标准扩展

## 准备工作

如果你想直接从二进制里释放整套参考资料，可以先执行：

```bash
# 方式 A：直接运行源码
go run . -action extract-assets -extract-root ./tsplay-assets-review

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action extract-assets -extract-root ./tsplay-assets-review
```

交付物清单：

- 源码目录：
  [../../script/tutorials/review_pack/checklists/138_required_deliverables.md](../../script/tutorials/review_pack/checklists/138_required_deliverables.md)
- 释放后的二进制资源：
  `./tsplay-assets-review/script/tutorials/review_pack/checklists/138_required_deliverables.md`

## Step 1: 先拿最近几节做对照

建议回头对照：

- `Lesson 131-136`

看它们分别交付了：

- 文档
- Flow / Lua 示例
- 命令
- 预期输出

## Step 2: 把“交付物”当成 review 项

以后新补一节 lesson 时，不只问“脚本写了没”，还要问：

- 命令有没有
- 预期结果有没有
- 承接关系有没有
- 需要的 demo/data 有没有

## 下一步

继续看：
[Lesson 139](139-large-flow-package-layout.md)
