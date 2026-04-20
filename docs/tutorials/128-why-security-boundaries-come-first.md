# Lesson 128: 为什么教程不能跳过权限边界

`Lesson 127` 已经把本地 Flow 和 MCP 放到了一起。  
这一节专门回答一个很容易被新手忽略的问题：

- 为什么教程不能一开始就把权限全开

## Step 1: 回看 `121-126` 的 blocked 结果

重点回看这些文件：

- `121-mcp-validate-allow-lua-blocked.json`
- `122-mcp-validate-allow-http-blocked.json`
- `123-mcp-validate-allow-file-access-blocked.json`
- `124-mcp-validate-allow-browser-state-blocked.json`
- `125-mcp-validate-allow-redis-blocked.json`
- `126-mcp-validate-allow-database-blocked.json`

## Step 2: 这些 blocked 结果在提醒什么

它们其实都在提醒同一件事：

- 这条 Flow 不是不能做
- 而是不能在“默认最小权限”下直接做

也就是说，MCP 的教学顺序应该是：

1. 先识别这条 Flow 触碰了哪类风险能力
2. 再决定要不要打开对应边界
3. 最后才进入执行

## Step 3: 为什么这比“一上来全开”更重要

因为一旦教程跳过边界，新手很容易学到一种危险习惯：

- 不是先理解动作属于哪类能力
- 而是先把 `full_automation` 当默认值

这会让后面的评审、交付和接入都变得模糊。

## 下一步

继续看：
[Lesson 129](129-security-preset-and-override.md)
