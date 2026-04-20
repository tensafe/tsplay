# Lesson 135: review 时，什么时候才允许 Lua escape hatch

`Lesson 121` 已经讲过 `allow_lua` 的安全边界。  
这一节不是再讲一次安全，而是继续往交付层走：

- 安全上能开，不代表结构上就应该开

目标：

- 跑一个最小 Lua 示例
- 理解“能写”不等于“该写”
- 建立 Lua escape hatch 的最小评审标准

## 准备工作

样例脚本：

- [../../script/tutorials/135_review_summary_escape_hatch.lua](../../script/tutorials/135_review_summary_escape_hatch.lua)

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials/135
```

## Step 1: 先运行 Lua 版本

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/135_review_summary_escape_hatch.lua -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/135_review_summary_escape_hatch.lua -headless
```

预期结果：

- 会成功运行
- 会写出 `artifacts/tutorials/135/review-summary-from-lua.json`

## Step 2: 这一节真正要判断什么

这个 Lua 脚本能跑，但它做的事情其实很简单：

- 组一个小 payload
- 写一份 JSON

这种场景通常不是“必须 Lua”，而更像“暂时偷懒用了 Lua”。

## Step 3: Lua escape hatch 的最小标准

优先只有在这些情况下，才考虑保留 Lua：

1. 真的需要 Flow 还不适合表达的逻辑。
2. 逻辑已经明显超出简单编排。
3. 团队愿意承担更高的 review 和维护成本。

如果只是：

- 组 payload
- 写文件
- 简单分支

那一般都应该优先回到 Flow。

## 下一步

继续看：
[Lesson 136](136-review-when-to-extract-lua-to-flow.md)
