# Lesson 47: 用命名会话直接进入受保护导入页

这一节把 `Lesson 46` 注册好的命名会话真正跑起来。

目标：

- 不再显式写状态文件路径
- 直接用 `session_import_demo` 进入受保护页面
- 完成第一条真正的业务导入

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)

默认会话名：

- `session_import_demo`

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/47_use_session_import_single.lua](../../script/tutorials/47_use_session_import_single.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/47_use_session_import_single.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/47_use_session_import_single.lua
```

如果你想换会话名，可以覆盖：

```bash
TSPLAY_SAVED_SESSION=session_import_demo \
./tsplay -script script/tutorials/47_use_session_import_single.lua
```

预期结果：

- 不会出现登录步骤
- 会直接进入导入表单
- 会生成 `artifacts/tutorials/47-use-session-import-single-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/47_use_session_import_single.flow.yaml](../../script/tutorials/47_use_session_import_single.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/47_use_session_import_single.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/47_use_session_import_single.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/47-use-session-import-single-flow.json`

## Step 3: 这节要理解什么

到这里为止，会话复用已经形成三层清晰结构：

- 页面里先产生登录态
- 状态文件负责持久化
- 命名会话负责给业务流程稳定入口

这一层分清楚之后，后面做批量导入就会简单很多。

## 下一节

下一节把单条导入扩成命名会话驱动的 CSV 批量导入：
[Lesson 48](48-use-session-batch-import-csv.md)
