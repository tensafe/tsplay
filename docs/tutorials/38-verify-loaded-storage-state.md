# Lesson 38: 验证加载后的状态到底是不是你想要的

光“能加载”还不够。  
这一节要进一步确认：

- 页面文本对不对
- `storage state` 对不对
- `cookie header` 对不对

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 继续复用上一节保存的状态文件
- 再次抓取页面状态、cookies、storage state
- 把验证结果写回 JSON

## 开始前

建议先跑完：

- [Lesson 36](36-save-storage-state.md)
- [Lesson 37](37-load-saved-storage-state.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/38_verify_loaded_storage_state.lua](../../script/tutorials/38_verify_loaded_storage_state.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/38_verify_loaded_storage_state.lua -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/38_verify_loaded_storage_state.lua -headless
```

预期结果：

- 会生成 `artifacts/tutorials/38-verify-loaded-storage-state-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/38_verify_loaded_storage_state.flow.yaml](../../script/tutorials/38_verify_loaded_storage_state.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/38_verify_loaded_storage_state.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/38_verify_loaded_storage_state.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/38-verify-loaded-storage-state-flow.json`

## Step 3: 这节新引入了什么

这一节把 Flow 顶层的：

- `storage_state`

换成了：

- `load_storage_state`

它们在这里表达的是同一个意思：  
运行前先把保存好的浏览器状态加载进来。

## 下一节

下一节不新增脚本，而是把 `36` 和 `37` 连起来，做成一次完整的“两阶段 round trip” 复盘。
[Lesson 39](39-storage-state-round-trip.md)
