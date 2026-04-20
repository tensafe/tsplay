# Lesson 73: 验证重跑后没有重复行

上一节完成了“同一个批次号重跑”。  
这一节把关注点收窄到一个问题：重跑之后，数据库里会不会多出重复行。

目标：

- `redis_get`
- `read_csv`
- `db_query`

## 开始前

建议先跑完：

- [Lesson 72](72-rerun-shared-batch-idempotently.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua](../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/73_verify_rerun_does_not_duplicate_rows.lua
```

预期结果：

- 会生成 `artifacts/tutorials/73-verify-rerun-does-not-duplicate-rows-lua.json`
- Lua 版本会直接检查 `line_no` 有没有重复

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml](../../script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/73_verify_rerun_does_not_duplicate_rows.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/73-verify-rerun-does-not-duplicate-rows-flow.json`

## Step 3: 这节意味着什么

到这里，重跑这件事不只是“又执行了一次”，而是有了明确的结果校验。

## 下一节

下一节故意换成一份带异常的输入，练习“保留有效数据、单独记录异常行”：
[Lesson 74](74-recover-external-sync-with-anomaly-ledger.md)
