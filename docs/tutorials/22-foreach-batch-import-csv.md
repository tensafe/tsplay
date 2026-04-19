# Lesson 22: 用 `foreach` 批量导入 CSV

这一节把单条导入动作扩成批量处理。  
我们会读取 [../../demo/data/import_users.csv](../../demo/data/import_users.csv)，然后一行一行提交到 [../../demo/import_workflow.html](../../demo/import_workflow.html)。

目标：

- `read_csv`
- `foreach`
- `append_var`

## 开始前

这一节同时需要：

1. 本地静态文件服务
2. 本地 CSV 文件

如果你只有单个二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/22_foreach_batch_import_csv.lua](../../script/tutorials/22_foreach_batch_import_csv.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/22_foreach_batch_import_csv.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/22_foreach_batch_import_csv.lua
```

预期结果：

- 会生成 `artifacts/tutorials/22-foreach-batch-import-csv-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/22_foreach_batch_import_csv.flow.yaml](../../script/tutorials/22_foreach_batch_import_csv.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/22_foreach_batch_import_csv.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/22_foreach_batch_import_csv.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/22-foreach-batch-import-csv-flow.json`

## Step 3: 这节的关键点

`foreach` 的价值不只是“重复做几次”，而是：

- 把一批输入变成一批动作
- 把一批动作再整理成一份结构化结果

## 下一节

下一节会故意引入坏数据，体验 `on_error` 怎么做局部恢复：
[Lesson 23](23-on-error-import-recovery.md)
