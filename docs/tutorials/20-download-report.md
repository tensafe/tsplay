# Lesson 20: 下载本地报表并回读验证

这一节把下载串成一个完整闭环：

- 从页面点下载
- 保存到本地文件
- 再把下载下来的 CSV 读回来

我们使用仓库里的 [../../demo/download.html](../../demo/download.html) 和 [../../demo/data/monthly_report.csv](../../demo/data/monthly_report.csv)。

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/20_download_report.lua](../../script/tutorials/20_download_report.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/20_download_report.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/20_download_report.lua
```

预期结果：

- 会生成 `artifacts/tutorials/20-downloaded-monthly-report-lua.csv`
- 会生成 `artifacts/tutorials/20-download-report-lua.json`

## Step 2: 这节为什么一定要“下载后再回读”

很多人会把“下载成功”理解成“点击没报错”。  
但真正更稳的方式是：

1. 把文件真的落到本地
2. 再把它读回来
3. 确认内容真的可用

所以这节同时用到了：

- `download_file`
- `read_csv`

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/20_download_report.flow.yaml](../../script/tutorials/20_download_report.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/20_download_report.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/20_download_report.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/20-downloaded-monthly-report-flow.csv`
- 会生成 `artifacts/tutorials/20-download-report-flow.json`

## Step 4: 这一段学完意味着什么

到这里，初级阶段最核心的一段已经打通了：

- 本地文件输入输出
- 页面重试与等待
- 上传与下载

这已经非常接近真实业务里的“最小可交付链路”。

## 下一步

接下来继续往外部系统走，可以回到：

- [Lesson 06](06-redis-round-trip.md)
- [Lesson 07](07-db-postgres-basics.md)

如果你要继续沿着课程体系推进，可以进入：
[track-junior.zh-CN.md](track-junior.zh-CN.md)
