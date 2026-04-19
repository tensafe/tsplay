# Lesson 29: 读取当前浏览器的 Cookie 字符串

这一节延续上一节的本地登录态页面，  
但视角从完整 `storage state` 切到更贴近 HTTP 请求的 `Cookie` 字符串。

使用页面：
[../../demo/session_lab.html](../../demo/session_lab.html)

目标：

- 在本地页面里写入 cookie
- 读取浏览器当前的 cookie header 字符串
- 把结果写到 `artifacts/tutorials/`

## 准备工作

如果静态文件服务还没启动，先执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/29_read_cookies_string.lua](../../script/tutorials/29_read_cookies_string.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/29_read_cookies_string.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/29_read_cookies_string.lua
```

预期结果：

- 会生成 `artifacts/tutorials/29-read-cookies-string-lua.json`

## Step 2: 为什么这一节要单独讲

很多时候你并不是要整份浏览器状态，  
而只是想知道：

- 当前 cookie 到底有没有被写进去
- 如果后面发请求，浏览器会带什么 header

所以这一节专门把 `Cookie` 单独拎出来看。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/29_read_cookies_string.flow.yaml](../../script/tutorials/29_read_cookies_string.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/29_read_cookies_string.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/29_read_cookies_string.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/29-read-cookies-string-flow.json`

## 下一节

下一节把状态文本、`storage state` 和 `cookie header` 合并成一份完整快照。
[Lesson 30](30-browser-state-snapshot-pack.md)
