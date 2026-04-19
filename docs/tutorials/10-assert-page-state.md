# Lesson 10: 对本地页面做可见性和文本断言

这一节开始补一组非常关键的基础动作：

- `assert_visible`
- `assert_text`

和前面“先点一下、先抓一下”相比，断言更接近真实交付，因为它开始回答：

- 页面是不是到了预期状态
- 结果是不是符合业务预期

这一节继续使用仓库里的 [../../demo/extract.html](../../demo/extract.html)。

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

页面地址：

```text
http://127.0.0.1:8000/demo/extract.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/10_assert_page_state.lua](../../script/tutorials/10_assert_page_state.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/10_assert_page_state.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/10_assert_page_state.lua
```

预期结果：

- 会断言 `#notice` 可见
- 会断言 `#notice` 文本里包含 `Ready`
- 会生成 `artifacts/tutorials/10-assert-page-state-lua.json`

## Step 2: 看懂“等待”和“断言”的边界

这节特意保留了：

- `wait_for_selector`
- `assert_visible`
- `assert_text`

它们不是重复关系，而是职责不同：

- `wait_for_selector` 更偏“等它出现”
- `assert_visible` 更偏“确认它真的可见”
- `assert_text` 更偏“确认业务文案已经对了”

也正因为这样，新手阶段不建议上来就 `sleep`。  
`sleep` 只能拖时间，不能表达你到底在等什么。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/10_assert_page_state.flow.yaml](../../script/tutorials/10_assert_page_state.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/10_assert_page_state.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/10_assert_page_state.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/10-assert-page-state-flow.json`
- 终端会看到结构化 trace

## Step 4: 这节真正要带走什么

- 页面动作不是只有“点”和“填”
- 断言是把流程从“能跑”推进到“可验证”的第一步
- 后面学 `retry`、`wait_until`、`on_error` 时，断言会是最常见的组合对象

## 下一节

下一节继续留在本地 demo，但把目标从“选项 5”改成“另一个选项”：
[Lesson 11](11-select-another-option.md)
