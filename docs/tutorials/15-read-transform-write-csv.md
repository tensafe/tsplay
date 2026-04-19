# Lesson 15: 读取、整理、再写出 CSV

这一节开始出现真正的小流程：

- 输入是一份 CSV
- 中间有一层整理
- 输出还是一份 CSV

目标：

- 学会 `read_csv -> transform -> write_csv`
- 认识 `start_row`、`limit`、`row_number_field`
- 开始把 CSV 当成批量流程的中间工件

## 开始前

这一节和 Lesson 13 一样，需要本地真实存在 CSV 文件。  
如果你不是在仓库根目录学习，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

如果你已经进入了 `./tsplay-assets` 目录，运行示例时请把命令里的 `./tsplay` 换成你的实际二进制路径，例如：

```bash
../tsplay -script script/tutorials/15_read_transform_write_csv.lua
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/15_read_transform_write_csv.lua](../../script/tutorials/15_read_transform_write_csv.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/15_read_transform_write_csv.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/15_read_transform_write_csv.lua
```

预期结果：

- 会生成 `artifacts/tutorials/15-transformed-contacts-lua.csv`
- 会生成 `artifacts/tutorials/15-read-transform-write-csv-lua.json`

## Step 2: 看这节为什么要引入 `start_row` 和 `limit`

这节不是把所有数据都读进来，而是只读一段：

- `start_row` 控制从哪一行开始
- `limit` 控制最多处理多少行
- `row_number_field` 让你知道每条结果来自原始文件的哪一行

这三个字段组合起来，后面特别适合做：

- 分批处理
- 断点续跑
- 失败记录回写

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/15_read_transform_write_csv.flow.yaml](../../script/tutorials/15_read_transform_write_csv.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/15_read_transform_write_csv.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/15_read_transform_write_csv.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/15-transformed-contacts-flow.csv`
- 会生成 `artifacts/tutorials/15-read-transform-write-csv-flow.json`

## Step 4: 这节真正重要的点

- CSV 不只是输入格式，也可以是中间结果格式
- 一旦引入 `row_number_field`，后续排错会轻松很多
- `foreach + append_var + write_csv` 是批量流程的一个很自然的主线

## 下一节

下一节回到浏览器，但开始进入健壮性主题：`retry`。
[Lesson 16](16-retry-flaky-action.md)
