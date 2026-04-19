# Lesson 01: Hello World

这一节只做一件事：先把 TSPlay 跑起来，而且先不碰网页。

你会同时看到：

- 一个最小 `Lua` 示例
- 一个最小 `Flow` 示例
- 两者都把结果写到 `artifacts/tutorials/`

## 这节要学会什么

- 知道怎么运行 `Lua` 脚本
- 知道怎么运行 `Flow`
- 理解 `Lua` 和 `Flow` 都不一定要从网页动作开始

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/01_hello_world.lua](../../script/tutorials/01_hello_world.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/01_hello_world.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/01_hello_world.lua
```

预期结果：

- 终端会打印一条 hello world 消息
- 会生成 `artifacts/tutorials/01-hello-world-lua.json`

这一版的重点不是浏览器，而是先认识 `set_var` 和 `write_json` 这类“无页面也能跑”的能力。

## Step 2: 看一下产物

打开：

```text
artifacts/tutorials/01-hello-world-lua.json
```

你会看到一个最小 JSON，对应这次练习输出。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/01_hello_world.flow.yaml](../../script/tutorials/01_hello_world.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/01_hello_world.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/01_hello_world.flow.yaml
```

预期结果：

- 终端会输出结构化 JSON 结果
- 会生成 `artifacts/tutorials/01-hello-world-flow.json`

## Step 4: 对比两种写法

先看 `Lua` 时，重点感受：

- 写法更直接
- 适合临时试一下、快速做一个脚本

再看 `Flow` 时，重点感受：

- 每一步都是结构化的
- 输出变量、步骤 trace 更容易被 review
- 后面接 MCP 或 AI 时，这种形式更稳定

## 下一节

下一节开始接触页面，但先不用外部网站，只操作仓库自带的本地 demo 页：
[Lesson 02](02-local-page-select-option.md)
