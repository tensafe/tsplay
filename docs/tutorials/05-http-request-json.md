# Lesson 05: 请求本地 JSON 并提取字段

这一节开始用 `http_request` 和 `json_extract`。  
为了让新手也能稳定复现，我们仍然不访问外网，而是直接请求仓库里的静态 JSON：
[../../demo/data/order_summary.json](../../demo/data/order_summary.json)

目标：

- 通过 HTTP 读取本地 JSON
- 解析返回结果里的状态码和业务字段
- 把结果写到 `artifacts/tutorials/`

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

JSON 地址：

```text
http://127.0.0.1:8000/demo/data/order_summary.json
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/05_http_request_json.lua](../../script/tutorials/05_http_request_json.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/05_http_request_json.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/05_http_request_json.lua
```

预期结果：

- 会生成 `artifacts/tutorials/05-http-request-json-lua.json`
- 输出里会包含 `http_status`、`open_count`、`first_order_id`

## Step 2: 理解这次返回值结构

这次 `http_request` 返回的不是裸 JSON，而是一个结果对象。  
常用字段至少包括：

- `status`
- `headers`
- `body`

所以这里的 JSONPath 才会写成：

```text
$.body.summary.open
$.body.orders[0].id
```

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/05_http_request_json.flow.yaml](../../script/tutorials/05_http_request_json.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/05_http_request_json.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/05_http_request_json.flow.yaml
```

预期结果：

- 会生成 `artifacts/tutorials/05-http-request-json-flow.json`
- 终端会输出每一步的结构化结果

## Step 4: 什么时候适合用这组动作

- 页面流程里顺便调一个业务 API
- 从内部接口拉 JSON，再继续写 Flow
- OCR、Webhook、补数接口这类“非页面动作”

补充说明：

- 本地直接跑 `go run . -flow ...` 不需要显式加 `allow_http`
- 如果你以后改成 MCP 模式运行，再补 `allow_http=true`

## 下一节

下一节开始准备 Redis 环境，体验最基础的读写和计数：
[Lesson 06](06-redis-round-trip.md)
