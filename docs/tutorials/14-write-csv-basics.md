# Lesson 14: 写出第一份 CSV

这一节把方向反过来：不是从文件里读，而是先在脚本或 Flow 里准备几行数据，再写成 CSV。

目标：

- 理解 `write_csv`
- 学会显式指定表头顺序
- 让产物同时落成 `.csv` 和 `.json`

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/14_write_csv_basics.lua](../../script/tutorials/14_write_csv_basics.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/14_write_csv_basics.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/14_write_csv_basics.lua
```

预期结果：

- 会生成 `artifacts/tutorials/14-write-csv-basics-lua.csv`
- 会生成 `artifacts/tutorials/14-write-csv-basics-lua.json`

## Step 2: 为什么这节要强调表头顺序

CSV 和 JSON 不一样。  
JSON 更关注对象结构，CSV 更关注列顺序。

所以这节特意显式传入了：

- `name`
- `city`
- `status`

这样即使你后面给对象加了别的字段，最终导出的列顺序仍然是稳定的。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/14_write_csv_basics.flow.yaml](../../script/tutorials/14_write_csv_basics.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/14_write_csv_basics.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/14_write_csv_basics.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/14-write-csv-basics-flow.csv`
- 会生成 `artifacts/tutorials/14-write-csv-basics-flow.json`

## Step 4: 这节要记住什么

- `write_csv` 更像“结构化导出”
- `headers` 很重要，它决定最终列顺序
- 很多交付场景会同时保留 `.csv` 和 `.json`

## 下一节

下一节开始把读和写串起来：读一份 CSV，做一层整理，再写出新的 CSV。
[Lesson 15](15-read-transform-write-csv.md)
