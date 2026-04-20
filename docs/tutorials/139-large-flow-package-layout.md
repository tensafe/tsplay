# Lesson 139: 设计一个大型 Flow 示例包的目录骨架

`Lesson 131-138` 已经把 review 规则收得差不多了。  
这一节开始把这些规则落到一个更大的目录骨架里。

目标：

- 看懂一个大型教程包为什么要拆 stage
- 跑通 `collect -> verify -> publish` 三个独立 stage
- 建立一个可 review、可复跑、可交接的目录样板

## 准备工作

目录样板：

- [../../script/tutorials/review_pack/139_large_flow_layout/README.md](../../script/tutorials/review_pack/139_large_flow_layout/README.md)
- [../../script/tutorials/review_pack/139_large_flow_layout/flows/collect.flow.yaml](../../script/tutorials/review_pack/139_large_flow_layout/flows/collect.flow.yaml)
- [../../script/tutorials/review_pack/139_large_flow_layout/flows/verify.flow.yaml](../../script/tutorials/review_pack/139_large_flow_layout/flows/verify.flow.yaml)
- [../../script/tutorials/review_pack/139_large_flow_layout/flows/publish.flow.yaml](../../script/tutorials/review_pack/139_large_flow_layout/flows/publish.flow.yaml)

## Step 1: 先跑 collect stage

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/review_pack/139_large_flow_layout/flows/collect.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/review_pack/139_large_flow_layout/flows/collect.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/139/collect/raw-items.json`

## Step 2: 再跑 verify stage

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/review_pack/139_large_flow_layout/flows/verify.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/review_pack/139_large_flow_layout/flows/verify.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/139/verify/verification-summary.json`

## Step 3: 最后跑 publish stage

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/review_pack/139_large_flow_layout/flows/publish.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/review_pack/139_large_flow_layout/flows/publish.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/139/publish/publish-manifest.json`

## Step 4: 这一节意味着什么

一个大的教程包，不应该只是一条越来越长的 Flow。  
更稳的做法通常是：

- collect
- verify
- publish

拆开之后，每个 stage 都更容易 review，也更容易复跑。

## 下一步

继续看：
[Lesson 140](140-why-tutorials-need-code-review-thinking.md)
