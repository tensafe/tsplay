# Lesson 146: 设计“哪些资源必须内置，哪些资源可以不带”的规则

`Lesson 145` 已经把离线学习路径说清楚了。  
这一节继续往交付边界收口：到底哪些资源必须跟着二进制一起走。

目标：

- 生成一份内置资源策略清单
- 明确必须内置的最小集合
- 避免发布包越来越大，但关键入口反而缺失

## 准备工作

样例文件：

- [../../script/tutorials/146_embedded_asset_policy.flow.yaml](../../script/tutorials/146_embedded_asset_policy.flow.yaml)

## Step 1: 先生成内置资源策略 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/146_embedded_asset_policy.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/146_embedded_asset_policy.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/146/embedded-asset-policy.json`

## Step 2: 再把这份策略和 `Lesson 141` 对照

重点确认：

- `ReadMe.md`
- `docs/`
- `script/`
- `demo/`

是不是都还在“必须内置”的集合里。

## Step 3: 这一节的最小结论

发布包不应该只追求“小”。  
更重要的是：

- 关键入口不能丢
- 教程和示例不能脱节

## 下一步

继续看：
[Lesson 147](147-file-srv-dev-vs-release.md)
