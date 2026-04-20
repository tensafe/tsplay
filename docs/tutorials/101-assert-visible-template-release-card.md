# Lesson 101: 先确认模板发布卡片真的在页面上

`Lesson 100` 已经把交接产物整理成模板包。  
从这一节开始，我们把模板包当成“准备发布的对象”来检查。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `assert_visible`
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
[../../script/tutorials/101_assert_visible_template_release_card.lua](../../script/tutorials/101_assert_visible_template_release_card.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/101_assert_visible_template_release_card.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/101_assert_visible_template_release_card.lua
```

预期结果：

- 会生成 `artifacts/tutorials/101-assert-visible-template-release-card-lua.json`

## Step 2: 为什么先做可见性断言

因为发布前的第一件事，不是立刻点按钮，  
而是先确认页面骨架、核心卡片和关键 badge 真的已经出现。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/101_assert_visible_template_release_card.flow.yaml](../../script/tutorials/101_assert_visible_template_release_card.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/101_assert_visible_template_release_card.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/101_assert_visible_template_release_card.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/101-assert-visible-template-release-card-flow.json`

## 下一节

下一节不只看“在不在”，  
还要开始检查“文字对不对”。
[Lesson 102](102-assert-text-template-release-status.md)
