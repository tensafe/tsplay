# Lesson 145: 设计“离线环境也能跑基础教程”的交付流程

`Lesson 144` 先把单二进制交付流程写成 manifest。  
这一节继续推进到一个更具体的场景：离线学习。

目标：

- 生成一份离线学习 manifest
- 明确哪些 lesson 能先跑
- 让离线环境的学习路径也保持循序渐进

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/145_offline_learning_manifest.flow.yaml](../../script/tutorials/145_offline_learning_manifest.flow.yaml)
- Checklist:
  [../../script/tutorials/release_pack/checklists/145_offline_learning_checklist.md](../../script/tutorials/release_pack/checklists/145_offline_learning_checklist.md)

## Step 1: 先生成离线学习 manifest

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/145_offline_learning_manifest.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/145_offline_learning_manifest.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/145/offline-learning-manifest.json`

## Step 2: 再把 lesson 顺序和前面阶段对上

这一节最关键的不是“全部都能离线跑”，而是：

- 先跑不依赖服务的内容
- 再进入需要 `file-srv`、Redis、Postgres 的内容

## Step 3: 这一节的最小结论

离线学习也要保持同样的节奏：

- 先资源发现
- 再资源释放
- 再最小成功运行
- 最后才进完整课程

## 下一步

继续看：
[Lesson 146](146-embedded-asset-policy.md)
