# Lesson 130: 完成安全边界模块的第一轮 checkpoint

`Lesson 121-129` 已经把高级阶段的第一块内容跑通了。  
这一节不再引入新动作，而是把这 10 节重新收成一个可复盘的 checkpoint。

## Step 1: 先把这一段的主线重新说一遍

这段高级教程的顺序不是随便排的，而是：

1. `121-126` 先把 6 类常见 `allow_*` 一项项跑清楚
2. `127` 再把本地 Flow 和 MCP 放在一起比较
3. `128` 解释为什么教程不能跳过边界
4. `129` 最后再理解 `security_preset` 和显式覆盖

## Step 2: 这一段最关键的结论

如果你现在要看一条新 Flow，推荐顺序应该是：

1. 先识别它属于哪类能力
2. 默认从 `readonly` 开始
3. 只打开最小匹配的 `allow_*`
4. 只有在确实值得时，才考虑 `security_preset`
5. 即使用了 preset，也要知道显式 `allow_*` 仍然可以覆盖它

## Step 3: 用这一段的产物做自检

你现在至少应该能解释这些文件为什么会一部分 `blocked`、一部分 `allowed`：

- `121-mcp-validate-allow-lua-*.json`
- `122-mcp-validate-allow-http-*.json`
- `123-mcp-validate-allow-file-access-*.json`
- `124-mcp-validate-allow-browser-state-*.json`
- `125-mcp-validate-allow-redis-*.json`
- `126-mcp-validate-allow-database-*.json`
- `129-mcp-validate-file-access-browser-write.json`
- `129-mcp-validate-http-full-automation*.json`

## Step 4: 这一节意味着什么

到这里，高级阶段已经不是“再学几个 action”。  
而是开始学会：

- 怎么给 Flow 划边界
- 怎么给教程安排边界
- 怎么给 AI 协作入口安排边界

## 下一步

如果你是按课程体系推进，  
可以继续看：
[track-advanced.zh-CN.md](track-advanced.zh-CN.md)

如果你是按长期演化推进，  
可以继续看：
[iteration-roadmap-160.md](iteration-roadmap-160.md)
