# Lesson 13: 读取本地 CSV 并写出 JSON

这一节把“页面读取”切换成“本地文件读取”。  
我们先不做复杂转换，只做一件事：把 CSV 读成结构化行数据。

这一节使用的示例文件是：
[../../demo/data/tutorial_contacts.csv](../../demo/data/tutorial_contacts.csv)

目标：

- 读取本地 CSV
- 理解首行表头和数据行的关系
- 把读取结果写到 `artifacts/tutorials/`

## 开始前

这节不需要浏览器，也不需要静态文件服务。  
但它需要本地真的有 CSV 文件。

如果你是在仓库根目录里学习，直接继续就行。  
如果你只有单个 `./tsplay` 二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

这样 `demo/data/tutorial_contacts.csv` 就会真正出现在本地磁盘上。

如果你已经进入了 `./tsplay-assets` 目录，运行示例时请把命令里的 `./tsplay` 换成你的实际二进制路径，例如：

```bash
../tsplay -script script/tutorials/13_read_csv_basics.lua
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/13_read_csv_basics.lua](../../script/tutorials/13_read_csv_basics.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/13_read_csv_basics.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/13_read_csv_basics.lua
```

预期结果：

- 会生成 `artifacts/tutorials/13-read-csv-basics-lua.json`
- 输出里会包含 `rows`

如果你想换自己的 CSV，可以覆盖路径：

```bash
TSPLAY_CSV_INPUT=/absolute/path/to/your.csv ./tsplay -script script/tutorials/13_read_csv_basics.lua
```

## Step 2: 看这次读出来的结构

这份 CSV 的第一行是表头，所以读出来的每一行会变成对象：

- `name`
- `city`
- `status`
- `source_row`

其中 `source_row` 不是 CSV 原生列，而是这节额外补进去的“源行号”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/13_read_csv_basics.flow.yaml](../../script/tutorials/13_read_csv_basics.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/13_read_csv_basics.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/13_read_csv_basics.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/13-read-csv-basics-flow.json`

## Step 4: 这节要带走什么

- `read_csv` 读出来的是结构化行，而不是整段原始文本
- 表头决定字段名
- `row_number_field` 很适合后面做批量处理、断点续跑和结果回写

## 下一节

下一节反过来做：先在脚本里组织数据，再把它写成 CSV。
[Lesson 14](14-write-csv-basics.md)
