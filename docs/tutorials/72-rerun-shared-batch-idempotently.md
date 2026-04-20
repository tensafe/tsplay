# Lesson 72: 用同一个批次号重跑同步，但不产生重复数据

上一节我们已经能把一条完整的外部系统链路跑通。  
这一节先不加新系统，而是练一个交付里非常常见的动作：重跑。

重点不是“再跑一遍”，而是：

- 继续使用同一个 `batch_id`
- 重写明细行
- 最终结果仍然只保留一份

目标：

- `redis_get`
- `read_csv`
- `db_transaction`
- `db_upsert`
- `db_insert_many`

## 开始前

建议先跑完：

- [Lesson 71](71-external-system-round-trip.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/72_rerun_shared_batch_idempotently.lua](../../script/tutorials/72_rerun_shared_batch_idempotently.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/72_rerun_shared_batch_idempotently.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/72_rerun_shared_batch_idempotently.lua
```

预期结果：

- 会生成 `artifacts/tutorials/72-rerun-shared-batch-idempotently-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml](../../script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/72_rerun_shared_batch_idempotently.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/72-rerun-shared-batch-idempotently-flow.json`

## Step 3: 这节意味着什么

到这里，你已经开始接触“幂等重跑”的最小形态。  
这比单次成功更接近真实交付，因为线上流程几乎一定会遇到重跑场景。

## 下一节

下一节专门验证这次重跑到底有没有产生重复行：
[Lesson 73](73-verify-rerun-does-not-duplicate-rows.md)
