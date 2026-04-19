# Lesson 42: 用命名会话直接复用登录态

到这里，我们终于从：

- 直接写状态文件路径

推进到：

- 直接写会话名字

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 用 `session_lab_demo` 直接进入已登录状态
- 同时给出 Lua 和 Flow 两种写法

## 开始前

建议先完成：

- [Lesson 40](40-save-named-session.md)
- [Lesson 41](41-inspect-named-session.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/42_use_named_session.lua](../../script/tutorials/42_use_named_session.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/42_use_named_session.lua -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/42_use_named_session.lua -headless
```

预期结果：

- 会生成 `artifacts/tutorials/42-use-named-session-lua.json`

如果你想改用别的会话名：

```bash
TSPLAY_SAVED_SESSION=your_session_name ./tsplay -script script/tutorials/42_use_named_session.lua -headless
```

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/42_use_named_session.flow.yaml](../../script/tutorials/42_use_named_session.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/42_use_named_session.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/42_use_named_session.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/42-use-named-session-flow.json`

## Step 3: 这节的关键变化

这里真正变化的，不是页面动作。  
而是浏览器配置的抽象层次提高了：

- 以前：记具体文件路径
- 现在：记一个稳定会话名

这也是后面更复杂业务 Flow 更推荐的方式。

## 下一节

最后一节把这条链补完整：会话不用了，怎么删除。
[Lesson 43](43-delete-named-session.md)
