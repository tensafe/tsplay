# Lesson 137: review 时，什么时候应该新增 demo 页面

`Lesson 131-136` 一直在处理 Flow、Lua、artifact 的组织问题。  
这一节继续把边界讲清楚：不是每个教程都要新开 demo 页面。

目标：

- 知道什么时候该复用现有 demo
- 知道什么时候真的需要新 demo 页面
- 让教程增长得更稳，而不是越长越散

## 准备工作

先把当前内置资源列出来：

```bash
# 方式 A：直接运行源码
go run . -action list-assets > artifacts/tutorials/137-list-assets.txt

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action list-assets > artifacts/tutorials/137-list-assets.txt
```

配套 checklist：

- [../../script/tutorials/review_pack/checklists/137_add_demo_page_checklist.md](../../script/tutorials/review_pack/checklists/137_add_demo_page_checklist.md)

## Step 1: 先看你手里已经有什么 demo

这一步的重点不是逐个打开，而是先建立一个习惯：

- 新增 demo 之前，先盘点仓库里已经有的页面

## Step 2: 再用 checklist 判断是否真要新开页面

优先只在这些情况下新增：

- 你要教页面行为
- 现有 demo 复现不了
- 这个页面能稳定承载新 lesson

## Step 3: 这一节意味着什么

教程增长得太快时，最容易失控的不是代码，而是 demo 页面数量。  
这一节的目的，就是让新增 demo 成为一个“需要理由”的动作。

## 下一步

继续看：
[Lesson 138](138-review-required-deliverables.md)
