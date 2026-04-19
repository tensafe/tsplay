# Lesson 40: 把状态文件注册成命名会话

前面几节我们一直在直接使用状态文件路径。  
这一节开始把它抽象成一个更容易复用的名字。

目标：

- 用 `tsplay` 命令把状态文件注册成命名会话
- 后面不再反复写具体文件路径
- 为 `use_session` 做准备

## 开始前

建议先跑完：

- [Lesson 36](36-save-storage-state.md)

默认会话名我们统一用：

- `session_lab_demo`

## Step 1: 用 Flow 产物注册一个命名会话

如果你刚才跑的是 Flow 版 `Lesson 36`，执行：

```bash
./tsplay -action save-session \
  -session-name session_lab_demo \
  -storage-state-path tutorials/36-session-lab-flow-state.json
```

这里要注意：

- `-storage-state-path` 写的是相对于 `artifact_root` 的路径
- 默认 `artifact_root` 就是 `artifacts`

所以最终它指向的是：

- `artifacts/tutorials/36-session-lab-flow-state.json`

## Step 2: 如果你用的是 Lua 产物

也可以直接注册 Lua 版状态文件：

```bash
./tsplay -action save-session \
  -session-name session_lab_demo \
  -storage-state-path tutorials/36-session-lab-lua-state.json
```

## Step 3: 这一步做完会发生什么

TSPlay 会在 `artifacts/sessions/` 下保存：

- 会话元数据
- 一份可复用的 storage state 副本

从这一步开始，你就不需要在每条 Flow 里反复硬编码状态文件路径了。

## 下一节

下一节继续看这个命名会话，但重点变成：

- 怎么列出来
- 怎么查看详情
- 怎么导出可直接复用的 snippet

[Lesson 41](41-inspect-named-session.md)
