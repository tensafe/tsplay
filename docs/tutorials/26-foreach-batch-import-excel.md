# Lesson 26: 用 Excel 驱动批量导入

这一节把前面的两条线真正合到一起：

- `read_excel`
- `foreach`
- 页面表单导入

目标：

- 从 Excel 读取用户数据
- 一行一行导入到本地 demo 表单
- 记录每一行的结果

## 开始前

这一节同时需要：

1. 本地静态文件服务
2. 本地 Excel 文件

先在仓库根目录启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

如果你只有单个二进制，记得先：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/26_foreach_batch_import_excel.lua](../../script/tutorials/26_foreach_batch_import_excel.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/26_foreach_batch_import_excel.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/26_foreach_batch_import_excel.lua
```

预期结果：

- 会生成 `artifacts/tutorials/26-foreach-batch-import-excel-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/26_foreach_batch_import_excel.flow.yaml](../../script/tutorials/26_foreach_batch_import_excel.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/26_foreach_batch_import_excel.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/26_foreach_batch_import_excel.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/26-foreach-batch-import-excel-flow.json`

## Step 3: 这节真正代表什么

到这里，你已经从“单条动作”走进了“数据驱动流程”：

- 输入不是一条数据，而是一批
- 页面动作开始在循环里重复执行
- 结果也开始天然适合做台账

## 下一节

下一节继续留在 Excel，但开始加入失败恢复和结果回写。
[Lesson 27](27-on-error-import-excel-writeback.md)
