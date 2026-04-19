# Lesson 57: 跑通认证导入到导出的完整 round trip

这一节把前面所有关键动作串成一个完整闭环：

- 用命名会话进入受保护页面
- 用 Excel 批量导入
- 抓取页面表格
- 下载当前导出文件
- 再把导出文件读回来

目标：

- `use_session`
- `read_excel`
- `capture_table`
- `download_file`
- `read_csv`

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)
- [Lesson 50](50-use-session-batch-import-excel.md)
- [Lesson 55](55-use-session-download-import-report-readback.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/57_use_session_import_export_round_trip.lua](../../script/tutorials/57_use_session_import_export_round_trip.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/57_use_session_import_export_round_trip.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/57_use_session_import_export_round_trip.lua
```

预期结果：

- 会生成 `artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv`
- 会生成 `artifacts/tutorials/57-use-session-import-export-round-trip-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/57_use_session_import_export_round_trip.flow.yaml](../../script/tutorials/57_use_session_import_export_round_trip.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/57_use_session_import_export_round_trip.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/57_use_session_import_export_round_trip.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv`
- 会生成 `artifacts/tutorials/57-use-session-import-export-round-trip-flow.json`

## Step 3: 这节意味着什么

走到这里，认证导入主线已经形成了一条完整的业务型教程链：

- 登录态可以保存和复用
- 受保护页面可以做单条和批量导入
- 页面结果可以被抓取和落盘
- 页面导出可以被下载和回读

这已经不再是单点 action 练习，而是一条很像真实业务交付的最小闭环。

## 下一节

下一节先从最轻的一步开始，把这份导出 CSV 的摘要写入 Redis：
[Lesson 58](58-sync-import-report-summary-to-redis.md)
