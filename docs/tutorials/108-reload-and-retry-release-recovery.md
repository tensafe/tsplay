# Lesson 108: `reload` 之后再验证一次恢复结果

到这里，页面上的常见不稳定性已经差不多凑齐了。  
这一节继续补最后一类：某些恢复结果只有刷新页面后才看得到。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `reload`
- `retry`
- `assert_text`
- `write_json`

## 准备工作

先确认本地静态文件服务已经启动：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/108_reload_and_retry_release_recovery.lua](../../script/tutorials/108_reload_and_retry_release_recovery.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/108_reload_and_retry_release_recovery.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/108_reload_and_retry_release_recovery.lua
```

预期结果：

- 会生成 `artifacts/tutorials/108-reload-and-retry-release-recovery-lua.json`

## Step 2: 这一节的重点

这里的重点不是单独学 `reload`。  
而是学会把“刷新页面”也纳入一个可验证的恢复链。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/108_reload_and_retry_release_recovery.flow.yaml](../../script/tutorials/108_reload_and_retry_release_recovery.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/108_reload_and_retry_release_recovery.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/108_reload_and_retry_release_recovery.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/108-reload-and-retry-release-recovery-flow.json`

## 下一节

前面的稳定性动作都学过以后，  
下一节开始把证据包重新带回来。
[Lesson 109](109-template-release-artifact-pack.md)
