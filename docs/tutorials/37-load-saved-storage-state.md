# Lesson 37: 从保存好的状态文件直接复用登录态

这一节是上一节的直接后续。  
我们不再重新输入用户名，而是直接加载已经保存好的浏览器状态文件。

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 加载上一节保存的状态文件
- 直接进入已登录状态
- 把复用结果写到 JSON

## 开始前

这节默认复用上一节产物，所以建议先跑完：

- [Lesson 36](36-save-storage-state.md)

如果你使用默认路径，那么会直接复用：

- Lua: `artifacts/tutorials/36-session-lab-lua-state.json`
- Flow: `artifacts/tutorials/36-session-lab-flow-state.json`

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/37_load_saved_storage_state.lua](../../script/tutorials/37_load_saved_storage_state.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/37_load_saved_storage_state.lua -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/37_load_saved_storage_state.lua -headless
```

预期结果：

- 会生成 `artifacts/tutorials/37-load-saved-storage-state-lua.json`

如果你想改用别的状态文件，可以覆盖：

```bash
TSPLAY_STATE_FILE=/absolute/path/to/state.json ./tsplay -script script/tutorials/37_load_saved_storage_state.lua -headless
```

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/37_load_saved_storage_state.flow.yaml](../../script/tutorials/37_load_saved_storage_state.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/37_load_saved_storage_state.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/37_load_saved_storage_state.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/37-load-saved-storage-state-flow.json`

## Step 3: 这一节在课程里的位置

`Lesson 36` 是“保存状态”。  
`Lesson 37` 是“复用状态”。

这两节连起来，才算真正进入了浏览器会话复用的起点。

## 下一节

下一节继续复用这份状态，但会进一步验证：

- `storage state` 本身
- `cookie header`

[Lesson 38](38-verify-loaded-storage-state.md)
