# Lesson 45: 用状态文件直接跳过登录进入受保护导入页

上一节我们是“先登录，再导入”。  
这一节开始真正复用 `Lesson 36` 保存下来的浏览器状态文件。

目标：

- 不再手工登录
- 直接进入受会话保护的导入表单
- 完成一条导入，确认状态文件真的可用

## 开始前

建议先跑完：

- [Lesson 36](36-save-storage-state.md)

默认会使用这些状态文件：

- Flow：`artifacts/tutorials/36-session-lab-flow-state.json`
- Lua：`artifacts/tutorials/36-session-lab-lua-state.json`

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/45_storage_state_auth_import.lua](../../script/tutorials/45_storage_state_auth_import.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/45_storage_state_auth_import.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/45_storage_state_auth_import.lua
```

如果你的状态文件路径不同，也可以覆盖：

```bash
TSPLAY_STATE_FILE=artifacts/tutorials/36-session-lab-lua-state.json \
./tsplay -script script/tutorials/45_storage_state_auth_import.lua
```

预期结果：

- 页面会直接显示 `Logged in as ...`
- 不会再停在登录面板
- 会生成 `artifacts/tutorials/45-storage-state-auth-import-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/45_storage_state_auth_import.flow.yaml](../../script/tutorials/45_storage_state_auth_import.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/45_storage_state_auth_import.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/45_storage_state_auth_import.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/45-storage-state-auth-import-flow.json`

## Step 3: 这节要理解什么

这一步验证的是：

- `save_storage_state` 保存下来的不仅是“能看”的状态
- 它还是“能继续做业务动作”的状态

从这一步开始，会话复用就不再只是观察 cookie，而是正式进入流程复用。

## 下一节

下一节把状态文件再往前抽象一层，注册成导入专用的命名会话：
[Lesson 46](46-save-import-session.md)
