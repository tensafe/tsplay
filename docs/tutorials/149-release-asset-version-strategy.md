# Lesson 149: 课程包版本号和资源版本号应该怎么维护

`Lesson 148` 已经把 first-run 的用户路径收清楚了。  
这一节继续到发布维护层：版本策略。

目标：

- 生成一份版本策略 JSON
- 区分“二进制行为变化”和“资源内容变化”
- 让教程更新也能被追踪，而不是静默漂移

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/149_release_asset_version_strategy.flow.yaml](../../script/tutorials/149_release_asset_version_strategy.flow.yaml)
- Checklist:
  [../../script/tutorials/release_pack/checklists/149_version_strategy.md](../../script/tutorials/release_pack/checklists/149_version_strategy.md)

## Step 1: 先生成版本策略 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/149_release_asset_version_strategy.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/149_release_asset_version_strategy.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/149/release-asset-version-strategy.json`

## Step 2: 再把“两层版本”记清楚

最小分层是：

- 二进制版本
- 资源版本

## Step 3: 这一节的最小结论

如果命令没变，但内置教程、demo、脚本已经变了一轮，那也应该是可追踪变化。  
这就是为什么高级阶段要单独讲资源版本。

## 下一步

继续看：
[Lesson 150](150-single-binary-delivery-summary.md)
