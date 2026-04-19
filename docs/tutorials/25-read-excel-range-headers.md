# Lesson 25: 读取 Excel 指定区域并显式声明表头

这一节继续使用 [../../demo/data/import_users.xlsx](../../demo/data/import_users.xlsx)，不过目标更进一步：

- 不是读整张表
- 而是只读一块范围
- 并手动声明这块数据的字段名

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/25_read_excel_range_headers.lua](../../script/tutorials/25_read_excel_range_headers.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/25_read_excel_range_headers.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/25_read_excel_range_headers.lua
```

预期结果：

- 会生成 `artifacts/tutorials/25-read-excel-range-headers-lua.json`

## Step 2: 这节要理解什么

这里同时出现了两个新点：

- `range`
- `headers`

它们适合的场景是：

- 工作表里不止一张逻辑表
- 或者目标区域本身没有规范表头
- 你想只截取一块稳定范围

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/25_read_excel_range_headers.flow.yaml](../../script/tutorials/25_read_excel_range_headers.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/25_read_excel_range_headers.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/25_read_excel_range_headers.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/25-read-excel-range-headers-flow.json`

## 下一节

下一节把 Excel 真正接进批量导入流程：`read_excel + foreach`。
[Lesson 26](26-foreach-batch-import-excel.md)
