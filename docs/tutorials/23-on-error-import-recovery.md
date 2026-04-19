# Lesson 23: 用 `on_error` 做局部恢复并回写结果

这一节会用 [../../demo/data/import_users_with_error.csv](../../demo/data/import_users_with_error.csv) 故意制造一条坏数据。

目标：

- 批量流程里允许某一条失败
- 失败后不要整批中断
- 最后把成功和失败都写成结果台账

## 开始前

这一节同样需要本地静态文件服务和本地 CSV 文件。  
如果你只有单个二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/23_on_error_import_recovery.lua](../../script/tutorials/23_on_error_import_recovery.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/23_on_error_import_recovery.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/23_on_error_import_recovery.lua
```

预期结果：

- 会生成 `artifacts/tutorials/23-on-error-import-recovery-lua.json`
- 会生成 `artifacts/tutorials/23-on-error-import-recovery-lua.csv`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/23_on_error_import_recovery.flow.yaml](../../script/tutorials/23_on_error_import_recovery.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/23_on_error_import_recovery.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/23_on_error_import_recovery.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/23-on-error-import-recovery-flow.json`
- 会生成 `artifacts/tutorials/23-on-error-import-recovery-flow.csv`

## Step 3: 这节要真正带走什么

`on_error` 的重点不是“把错误吞掉”，而是：

- 明确哪一小段可能失败
- 失败后只在局部恢复
- 同时把失败现场也记进结果里

这会直接决定你的批量流程是不是可交付。

## 下一节

下一节把输入源从 CSV 再切到 Excel：
[Lesson 24](24-read-excel-basics.md)
