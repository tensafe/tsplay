# Lesson 107: 用 `retry` 接住一次偶发失败点击

`Lesson 103` 的 retry 更像是 gate 通过重试。  
这一节更贴近日常页面问题：按钮第一次点了没反应，第二次才成功。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `retry`
- `assert_text`
- `assert_visible`
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
[../../script/tutorials/107_retry_flaky_publish_click.lua](../../script/tutorials/107_retry_flaky_publish_click.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/107_retry_flaky_publish_click.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/107_retry_flaky_publish_click.lua
```

预期结果：

- 会生成 `artifacts/tutorials/107-retry-flaky-publish-click-lua.json`

## Step 2: 和 `Lesson 103` 的区别

两节都在讲 `retry`，  
但这节更强调“把重试包在具体交互动作外面”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/107_retry_flaky_publish_click.flow.yaml](../../script/tutorials/107_retry_flaky_publish_click.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/107_retry_flaky_publish_click.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/107_retry_flaky_publish_click.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/107-retry-flaky-publish-click-flow.json`

## 下一节

下一节把恢复流程再往前推一步：  
有些状态必须刷新页面后才能真正生效。
[Lesson 108](108-reload-and-retry-release-recovery.md)
