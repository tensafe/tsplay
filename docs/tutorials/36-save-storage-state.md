# Lesson 36: 把当前浏览器状态保存到文件

前面几节我们已经能看见：

- `storage state`
- `cookie header`

这一节开始把这些状态真正落盘，变成后续可复用的文件。

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 在本地页面里生成登录态
- 把当前浏览器状态保存成文件
- 把状态文件路径写到 JSON

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/36_save_storage_state.lua](../../script/tutorials/36_save_storage_state.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/36_save_storage_state.lua -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/36_save_storage_state.lua -headless
```

预期结果：

- 会生成 `artifacts/tutorials/36-session-lab-lua-state.json`
- 会生成 `artifacts/tutorials/36-save-storage-state-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/36_save_storage_state.flow.yaml](../../script/tutorials/36_save_storage_state.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/36_save_storage_state.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/36_save_storage_state.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/36-session-lab-flow-state.json`
- 会生成 `artifacts/tutorials/36-save-storage-state-flow.json`

## Step 3: 这一节真正带走什么

这里最重要的不是“又学了一个动作”，  
而是开始建立一个新认知：

- 浏览器状态不是只存在内存里
- 它可以被保存成文件
- 后续流程可以直接复用这份文件

## 下一节

下一节就用这份状态文件，直接进入已登录状态。
[Lesson 37](37-load-saved-storage-state.md)
