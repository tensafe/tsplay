# Lesson 105: 用 `on_error` 接住模板发布校验失败

前面几节都还是“成功后再继续”。  
这一节开始故意制造一次失败，让流程学会自己收住，而不是整条链直接崩掉。

使用页面：
[../../demo/template_release_lab.html](../../demo/template_release_lab.html)

目标：

- `on_error`
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
[../../script/tutorials/105_on_error_template_release_validation.lua](../../script/tutorials/105_on_error_template_release_validation.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/105_on_error_template_release_validation.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/105_on_error_template_release_validation.lua
```

预期结果：

- 会生成 `artifacts/tutorials/105-on-error-template-release-validation-lua.json`

## Step 2: 这一节在练什么

不是练“怎么失败”，  
而是练“失败之后流程还能不能自己回到可继续的状态”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/105_on_error_template_release_validation.flow.yaml](../../script/tutorials/105_on_error_template_release_validation.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/105_on_error_template_release_validation.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/105_on_error_template_release_validation.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/105-on-error-template-release-validation-flow.json`

## 下一节

下一节换成另一种等待问题：  
目标元素一开始根本不存在，而是稍后才出现。
[Lesson 106](106-wait-for-delayed-release-note.md)
