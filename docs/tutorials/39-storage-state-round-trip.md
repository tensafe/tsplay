# Lesson 39: 把“保存状态”和“复用状态”连成一次完整 round trip

这一节不新增脚本。  
它的重点是把前面三节真正连起来，形成一条可复盘的两阶段链路。

目标：

- 第一次运行：保存状态文件
- 第二次运行：加载状态文件
- 对比两次结果，确认“状态复用”已经成立

## Step 1: 先执行保存阶段

优先建议直接用 Flow 版本：

```bash
./tsplay -flow script/tutorials/36_save_storage_state.flow.yaml -headless
```

也可以用 Lua 版本：

```bash
./tsplay -script script/tutorials/36_save_storage_state.lua -headless
```

预期产物：

- `artifacts/tutorials/36-session-lab-*.json`
- `artifacts/tutorials/36-save-storage-state-*.json`

## Step 2: 再执行复用阶段

Flow：

```bash
./tsplay -flow script/tutorials/37_load_saved_storage_state.flow.yaml -headless
```

Lua：

```bash
./tsplay -script script/tutorials/37_load_saved_storage_state.lua -headless
```

## Step 3: 看什么才算 round trip 成功

至少要同时看到两件事：

1. 第一阶段真的生成了状态文件
2. 第二阶段在不重新输入用户名的前提下，页面直接显示 `Logged in as ...`

## Step 4: 为什么这里要单独停一下

因为从课程结构上说，  
这一节是“文件级浏览器状态复用”阶段的小结。

再往后，我们就要把文件路径进一步抽象成：

- 命名会话
- `use_session`

## 下一节

下一节开始把保存好的状态文件注册成一个可复用的命名会话。
[Lesson 40](40-save-named-session.md)
