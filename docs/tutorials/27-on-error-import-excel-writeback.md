# Lesson 27: Excel 批量导入、局部恢复与结果回写

这一节是前面几节的一个小汇总：

- `read_excel`
- `foreach`
- `on_error`
- `write_json`
- `write_csv`

目标：

- 读取一批 Excel 行
- 导入时允许局部失败
- 把成功和失败都回写成结果台账

## 开始前

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/27_on_error_import_excel_writeback.lua](../../script/tutorials/27_on_error_import_excel_writeback.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/27_on_error_import_excel_writeback.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/27_on_error_import_excel_writeback.lua
```

预期结果：

- 会生成 `artifacts/tutorials/27-on-error-import-excel-writeback-lua.json`
- 会生成 `artifacts/tutorials/27-on-error-import-excel-writeback-lua.csv`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/27_on_error_import_excel_writeback.flow.yaml](../../script/tutorials/27_on_error_import_excel_writeback.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/27_on_error_import_excel_writeback.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/27_on_error_import_excel_writeback.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/27-on-error-import-excel-writeback-flow.json`
- 会生成 `artifacts/tutorials/27-on-error-import-excel-writeback-flow.csv`

## Step 3: 这节是一个阶段性分界点

这节意味着你已经不只是“会点页面”了。  
你开始具备下面这些交付能力：

- 用表格驱动一批数据
- 允许局部失败而不是整批报废
- 给结果留出 JSON 和 CSV 两种形态

这已经很接近真实业务里的批量导入主线了。

## 下一步

继续往下走，最自然的就是：

- 断点续跑
- Redis checkpoint
- 批量数据库写入
- `observe -> draft -> validate -> run -> repair`
