# Lesson 134: 用一份 checklist 正式 review 教程示例

`Lesson 131-133` 先把名字、变量名、artifact 布局理顺了。  
这一节开始把这些经验收成一份真正能复用的 review checklist。

目标：

- 学会用 checklist review 教程示例
- 不再只凭“感觉”
- 让新加的 lesson 也能沿同一标准迭代

## 准备工作

如果你手里只有二进制，没有源码目录，可以先把内置资源释放出来：

```bash
# 方式 A：直接运行源码
go run . -action extract-assets -extract-root ./tsplay-assets

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

review checklist：

- 源码目录：
  [../../script/tutorials/review_pack/checklists/134_example_review_checklist.md](../../script/tutorials/review_pack/checklists/134_example_review_checklist.md)
- 释放后的二进制资源：
  `./tsplay-assets/script/tutorials/review_pack/checklists/134_example_review_checklist.md`

## Step 1: 选两个已经跑过的样例

建议直接 review 这两个：

- [../../script/tutorials/131_review_readability_after.flow.yaml](../../script/tutorials/131_review_readability_after.flow.yaml)
- [../../script/tutorials/133_review_artifact_layout_after.flow.yaml](../../script/tutorials/133_review_artifact_layout_after.flow.yaml)

## Step 2: 按 checklist 逐项过一遍

不要一上来挑实现细节，先从这些入口过：

1. `name`
2. `description`
3. `save_as`
4. artifact 路径
5. 输出文件名

## Step 3: 这一节意味着什么

从这里开始，高级教程不再只是“我觉得这样更好”，而是：

- 可以复用
- 可以解释
- 可以让团队一起照着 review

## 下一步

继续看：
[Lesson 135](135-review-when-lua-is-allowed.md)
