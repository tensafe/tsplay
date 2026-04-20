# Lesson 102: 继续确认模板发布状态文字对不对

`Lesson 101` 先确认了卡片可见。  
这一节继续往前，只多加一个概念：文本断言。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `assert_text`
- `extract_text`
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
[../../script/tutorials/102_assert_text_template_release_status.lua](../../script/tutorials/102_assert_text_template_release_status.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/102_assert_text_template_release_status.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/102_assert_text_template_release_status.lua
```

预期结果：

- 会生成 `artifacts/tutorials/102-assert-text-template-release-status-lua.json`

## Step 2: 这一节和上一节的差别

上一节回答的是“元素在不在”。  
这一节回答的是“它说的是不是你以为的那句话”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/102_assert_text_template_release_status.flow.yaml](../../script/tutorials/102_assert_text_template_release_status.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/102_assert_text_template_release_status.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/102_assert_text_template_release_status.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/102-assert-text-template-release-status-flow.json`

## 下一节

下一节开始进入第一类不稳定场景：  
按钮第一次点不成功，第二次才过。
[Lesson 103](103-retry-template-release-gate.md)
