# Lesson 46: 把状态文件注册成导入专用命名会话

前一节我们已经证明，状态文件可以直接带你跳过登录。  
这一节继续把它抽象成一个更稳定、可读的名字。

目标：

- 不再反复写状态文件路径
- 为后面的受保护导入流程统一一个会话名
- 让 `use_session` 直接变成业务级写法

## 开始前

建议先跑完：

- [Lesson 45](45-storage-state-auth-import.md)

这条线统一使用的命名会话名是：

- `session_import_demo`

## Step 1: 用 Flow 状态文件注册命名会话

```bash
# 方式 A：直接运行源码
go run . -action save-session \
  -session-name session_import_demo \
  -storage-state-path tutorials/36-session-lab-flow-state.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action save-session \
  -session-name session_import_demo \
  -storage-state-path tutorials/36-session-lab-flow-state.json
```

## Step 2: 如果你更想复用 Lua 产物

也可以直接注册 Lua 版状态文件：

```bash
# 方式 A：直接运行源码
go run . -action save-session \
  -session-name session_import_demo \
  -storage-state-path tutorials/36-session-lab-lua-state.json

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action save-session \
  -session-name session_import_demo \
  -storage-state-path tutorials/36-session-lab-lua-state.json
```

## Step 3: 检查会话是否真的保存好了

```bash
# 方式 A：直接运行源码
go run . -action get-session -session-name session_import_demo

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action get-session -session-name session_import_demo
```

预期结果：

- 返回的 JSON 里能看到 `session_import_demo`
- 会显示底层关联的 storage state 信息

## Step 4: 这节要理解什么

状态文件更像“底层实现”。  
命名会话更像“业务层入口”。

从后面的教程开始，我们默认优先写：

- `use_session("session_import_demo")`
- `browser.use_session: session_import_demo`

而不是每次都回到底层路径。

## 下一节

下一节开始正式用这个命名会话进入受保护导入页：
[Lesson 47](47-use-session-import-single.md)
