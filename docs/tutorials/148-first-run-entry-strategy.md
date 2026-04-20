# Lesson 148: 新用户第一次打开二进制时，应该先看到什么

`Lesson 147` 已经把服务入口分成开发态和发布态。  
这一节继续往 onboarding 走：一个新用户第一次拿到 `tsplay` 时，最应该先走哪条路径。

目标：

- 生成一份 first-run entry manifest
- 把第一次接触的顺序收成固定入口
- 降低“我拿到二进制但不知道先干嘛”的摩擦

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/148_first_run_entry_manifest.flow.yaml](../../script/tutorials/148_first_run_entry_manifest.flow.yaml)
- Checklist:
  [../../script/tutorials/release_pack/checklists/148_first_run_entry_checklist.md](../../script/tutorials/release_pack/checklists/148_first_run_entry_checklist.md)

## Step 1: 先生成 first-run entry manifest

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/148_first_run_entry_manifest.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/148_first_run_entry_manifest.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/148/first-run-entry-manifest.json`

## Step 2: 再把顺序和你前面学过的内容对齐

推荐顺序不是随机的，而是：

1. `list-assets`
2. `extract-assets`
3. `Lesson 01`
4. `docs/tutorials/README.md`

## Step 3: 这一节的最小结论

第一步入口的任务不是“展示所有能力”，而是：

- 先看见资源
- 再建立第一次成功体验
- 最后再进入完整课程

## 下一步

继续看：
[Lesson 149](149-release-asset-version-strategy.md)
