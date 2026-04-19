# Lesson 24: 读取第一份 Excel

这一节把本地表格输入从 CSV 切到 Excel。  
我们使用的是 [../../demo/data/import_users.xlsx](../../demo/data/import_users.xlsx)。

目标：

- 读取一个本地 `.xlsx`
- 指定工作表
- 把结构化结果写到 `artifacts/tutorials/`

## 开始前

这一节和 CSV 教程一样，不需要浏览器。  
但它需要本地真实存在 Excel 文件。

如果你只有单个 `./tsplay` 二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/24_read_excel_basics.lua](../../script/tutorials/24_read_excel_basics.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/24_read_excel_basics.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/24_read_excel_basics.lua
```

预期结果：

- 会生成 `artifacts/tutorials/24-read-excel-basics-lua.json`

## Step 2: 这节和 `read_csv` 的差异

这节最先要建立的认知是：

- `read_csv` 更像“单表文本文件”
- `read_excel` 还要多一个“工作表”的维度

所以这里会显式指定：

- `sheet: Users`

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/24_read_excel_basics.flow.yaml](../../script/tutorials/24_read_excel_basics.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/24_read_excel_basics.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/24_read_excel_basics.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/24-read-excel-basics-flow.json`

## 下一节

下一节继续读 Excel，不过会开始加入 `range` 和显式 `headers`。
[Lesson 25](25-read-excel-range-headers.md)
