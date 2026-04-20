# Lesson 104: 用 `wait_until` 等模板发布检查完成

`Lesson 103` 解决的是“偶发失败，需要再试一次”。  
这一节解决的是另一个很常见的问题：动作已经触发了，但结果不会立刻出现。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `wait_until`
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
[../../script/tutorials/104_wait_until_template_release_ready.lua](../../script/tutorials/104_wait_until_template_release_ready.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/104_wait_until_template_release_ready.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/104_wait_until_template_release_ready.lua
```

预期结果：

- 会生成 `artifacts/tutorials/104-wait-until-template-release-ready-lua.json`

## Step 2: 为什么这里不用 `sleep`

因为真正稳定的写法，不是先猜一个时间，  
而是等“页面事实”出现。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/104_wait_until_template_release_ready.flow.yaml](../../script/tutorials/104_wait_until_template_release_ready.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/104_wait_until_template_release_ready.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/104_wait_until_template_release_ready.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/104-wait-until-template-release-ready-flow.json`

## 下一节

下一节把失败正式纳入流程：  
故意触发一次错误，再在错误分支里恢复。
[Lesson 105](105-on-error-template-release-validation.md)
