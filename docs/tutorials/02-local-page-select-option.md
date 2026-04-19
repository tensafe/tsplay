# Lesson 02: 打开本地页面并选择选项

这一节开始接触页面，但仍然不依赖外部网站。  
我们直接使用仓库里的 [../../demo/demo.html](../../demo/demo.html)。

目标很简单：

- 打开本地 demo 页
- 选择下拉框里的 `选项 5`
- 把结果写到 `artifacts/tutorials/`

## 先准备 TSPlay 内置静态文件服务

在仓库根目录另开一个终端，执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

准备好之后，页面地址就是：

```text
http://127.0.0.1:8000/demo/demo.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/02_select_option.lua](../../script/tutorials/02_select_option.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/02_select_option.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/02_select_option.lua
```

预期结果：

- 浏览器会打开本地 demo 页
- 脚本会把下拉框切换到 `选项 5`
- 会生成 `artifacts/tutorials/02-select-option-lua.json`

提示：

- 这个脚本执行完后，浏览器会保持打开，方便你看结果
- 观察完成后按 `Ctrl+C` 结束

## Step 2: 看 Lua 做了什么

这一版主要串起了这些动作：

- `navigate`
- `wait_for_selector`
- `select_option`
- `is_selected`
- `write_json`

这是典型的“过程式”写法：怎么做，就怎么一行一行写出来。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/02_select_option.flow.yaml](../../script/tutorials/02_select_option.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/02_select_option.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/02_select_option.flow.yaml -headless
```

如果你想直接看浏览器过程，可以去掉 `-headless`。

预期结果：

- Flow 会访问同一个本地 demo 页
- 结果会写到 `artifacts/tutorials/02-select-option-flow.json`
- 终端里还能看到这次执行的结构化 trace

## Step 4: 对比两种写法

这节里同一个功能的核心差异已经很明显了：

- `Lua` 强在顺手，尤其适合边试边改
- `Flow` 强在可读、可审查、可复用

如果把这个动作继续往业务流程扩展，比如“选完选项后再断言、截图、写台账”，`Flow` 的优势会越来越明显。

## 下一节

下一节继续用本地 demo 页，不过目标换成“抓表格”：
[Lesson 03](03-capture-table.md)
