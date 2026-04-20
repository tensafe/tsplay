# Lesson 140: 复盘为什么教程也需要 code review 思维

`Lesson 121-130` 讲的是安全边界。  
`Lesson 131-139` 讲的是 review 和组织方式。  
这一节把这两条线收在一起。

目标：

- 重新理解“能跑”和“可交付”之间的差别
- 把高级阶段前半段收成一张简单心智图
- 为后面的更大规模教程演进打底

## 准备工作

如果你想把这一整段配套资源释放出来统一看，可以执行：

```bash
# 方式 A：直接运行源码
go run . -action extract-assets -extract-root ./tsplay-assets-140

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action extract-assets -extract-root ./tsplay-assets-140
```

建议一起回看：

- `Lesson 121-130`
- `Lesson 131-139`
- [../../script/tutorials/review_pack/README.md](../../script/tutorials/review_pack/README.md)

## Step 1: 先复盘安全边界那一段

那一段回答的是：

- 哪些动作默认不能放行
- 为什么要最小授权

## Step 2: 再复盘 review 这一段

这一段回答的是：

- 即使能跑，结构是不是清楚
- 变量名、artifact、目录骨架是不是能被别人接手

## Step 3: 这一节的最小结论

教程也需要 code review 思维，因为教程本身就是一种长期资产。  
如果它只在“作者本人刚写完时”能看懂，那它就还没有真正交付完成。

## 下一步

继续看高级路线说明：
[track-advanced.md](track-advanced.md)
