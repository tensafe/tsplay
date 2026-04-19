# TSPlay Step-by-Step 教程

这套教程专门解决一个上手问题：

- 不从大而全的概念开始
- 同一个功能，同时给出 `Lua` 和 `Flow` 两种写法
- 先跑通，再逐步理解为什么日常交付更推荐 `Flow`

仓库里当前使用的是 `docs/` 目录，所以这套教程统一放在这里，而不是单独建 `doc/`。

## 推荐学习顺序

1. [Lesson 01: Hello World，不打开网页也能先跑通](01-hello-world.md)
2. [Lesson 02: 打开本地页面并选择下拉框选项](02-local-page-select-option.md)
3. [Lesson 03: 抓取本地表格并写出 JSON](03-capture-table.md)
4. [Lesson 04: 提取文本和 HTML 片段](04-extract-text-and-html.md)
5. [Lesson 05: 请求本地 JSON 并提取字段](05-http-request-json.md)
6. [Lesson 06: Redis 基础读写和计数](06-redis-round-trip.md)
7. [Lesson 07: Postgres 基础查询与写入](07-db-postgres-basics.md)

## 教程地图

| Lesson | 功能 | Lua 示例 | Flow 示例 | 说明 |
| --- | --- | --- | --- | --- |
| 01 | Hello World + 写 JSON | [../../script/tutorials/01_hello_world.lua](../../script/tutorials/01_hello_world.lua) | [../../script/tutorials/01_hello_world.flow.yaml](../../script/tutorials/01_hello_world.flow.yaml) | 不需要打开网页 |
| 02 | 打开本地 demo 页并选择选项 | [../../script/tutorials/02_select_option.lua](../../script/tutorials/02_select_option.lua) | [../../script/tutorials/02_select_option.flow.yaml](../../script/tutorials/02_select_option.flow.yaml) | 需要一个本地静态文件服务 |
| 03 | 抓取本地表格并写出 JSON | [../../script/tutorials/03_capture_table.lua](../../script/tutorials/03_capture_table.lua) | [../../script/tutorials/03_capture_table.flow.yaml](../../script/tutorials/03_capture_table.flow.yaml) | 继续复用本地静态文件服务 |
| 04 | 提取文本和 HTML 片段 | [../../script/tutorials/04_extract_text_and_html.lua](../../script/tutorials/04_extract_text_and_html.lua) | [../../script/tutorials/04_extract_text_and_html.flow.yaml](../../script/tutorials/04_extract_text_and_html.flow.yaml) | 继续复用本地静态文件服务 |
| 05 | 请求本地 JSON 并提取字段 | [../../script/tutorials/05_http_request_json.lua](../../script/tutorials/05_http_request_json.lua) | [../../script/tutorials/05_http_request_json.flow.yaml](../../script/tutorials/05_http_request_json.flow.yaml) | 继续复用本地静态文件服务 |
| 06 | Redis 基础读写和计数 | [../../script/tutorials/06_redis_round_trip.lua](../../script/tutorials/06_redis_round_trip.lua) | [../../script/tutorials/06_redis_round_trip.flow.yaml](../../script/tutorials/06_redis_round_trip.flow.yaml) | 需要本地 Redis |
| 07 | Postgres 基础查询与写入 | [../../script/tutorials/07_db_postgres_basics.lua](../../script/tutorials/07_db_postgres_basics.lua) | [../../script/tutorials/07_db_postgres_basics.flow.yaml](../../script/tutorials/07_db_postgres_basics.flow.yaml) | 需要本地 Postgres |

## 开始前

先在仓库根目录执行：

```bash
go mod download
```

如果你想直接用 `tsplay` 命令测试，先构建一次：

```bash
go build -o tsplay .
```

构建完成后，下面这些命令都可以直接跑：

```bash
./tsplay -script script/tutorials/01_hello_world.lua
./tsplay -flow script/tutorials/01_hello_world.flow.yaml
./tsplay -action cli
./tsplay -action file-srv -addr :8000
```

如果你暂时不想构建二进制，也可以继续用：

```bash
go run . -script script/tutorials/01_hello_world.lua
go run . -flow script/tutorials/01_hello_world.flow.yaml
```

补充说明：

- 构建出来的 `./tsplay` 已经内置了 `ReadMe.md`、`docs/`、`script/`、`demo/`
- 所以即使你把单个二进制拿到别的目录，也还能直接跑这些示例路径
- 如果你想把这套参考资料释放出来，可以执行：

```bash
./tsplay -action list-assets
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

Lesson 01 可以直接开始。  
从 Lesson 02 到 Lesson 05，建议另开一个终端，在仓库根目录启动 TSPlay 内置静态文件服务：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

如果你已经只有单个 `./tsplay` 二进制、没有完整仓库目录，也可以直接执行上面这条命令。  
此时服务的是二进制里内置的 `demo/`、`docs/`、`script/` 资源。

然后通过下面这些本地地址访问仓库里的 demo 页面：

- `http://127.0.0.1:8000/demo/demo.html`
- `http://127.0.0.1:8000/demo/tables.html`
- `http://127.0.0.1:8000/demo/extract.html`
- `http://127.0.0.1:8000/demo/data/order_summary.json`

如果你不是在仓库根目录启动服务，可以显式指定：

```bash
./tsplay -action file-srv -addr :8000 -serve-root /path/to/tsplay
```

从 Lesson 06 起，开始接 Redis 和数据库。  
为了让命令更好抄，可以直接复用这些辅助文件：

- Redis 环境变量示例：[../../script/tutorials/env/06_redis_example.sh](../../script/tutorials/env/06_redis_example.sh)
- Postgres 环境变量示例：[../../script/tutorials/env/07_reporting_pg_example.sh](../../script/tutorials/env/07_reporting_pg_example.sh)
- Postgres 初始化 SQL：[../../script/tutorials/sql/07_reporting_pg.sql](../../script/tutorials/sql/07_reporting_pg.sql)

## 学的时候重点看什么

- `Lua` 更像“我现在一步一步怎么做”
- `Flow` 更像“这个流程本身是什么，可以怎么被 review、复用和交给 AI 生成”
- 同一个功能都先看 `Lua` 再看 `Flow`，更容易体会两者的边界

## 两个小提醒

- 运行带浏览器的 `Lua` 脚本时，脚本执行完成后浏览器会继续保持打开，方便你观察页面；结束时按 `Ctrl+C` 即可
- 这些教程默认把输出写到 `artifacts/tutorials/`，这样不会把练习产物混进仓库版本里
