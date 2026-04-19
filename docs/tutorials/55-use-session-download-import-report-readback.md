# Lesson 55: 把认证导出 CSV 下载后再读回来

这一节把下载串成完整闭环：

- 先从认证页导出
- 再把文件落到本地
- 然后读回内容

目标：

- `download_file`
- `read_csv`
- 认证态导出文件回读

## 开始前

建议先跑完：

- [Lesson 54](54-use-session-download-import-report.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/55_use_session_download_import_report_readback.lua](../../script/tutorials/55_use_session_download_import_report_readback.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/55_use_session_download_import_report_readback.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/55_use_session_download_import_report_readback.lua
```

预期结果：

- 会生成 `artifacts/tutorials/55-use-session-import-report-lua.csv`
- 会生成 `artifacts/tutorials/55-use-session-download-import-report-readback-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/55_use_session_download_import_report_readback.flow.yaml](../../script/tutorials/55_use_session_download_import_report_readback.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/55_use_session_download_import_report_readback.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/55_use_session_download_import_report_readback.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/55-use-session-import-report-flow.csv`
- 会生成 `artifacts/tutorials/55-use-session-download-import-report-readback-flow.json`

## Step 3: 这节要理解什么

点击下载按钮不等于结果可用。  
真正更稳的判断方式还是：

1. 文件落到本地
2. 能被重新读取
3. 内容结构确实正常

## 下一节

下一节把“页面表格”和“下载文件”放到一份结果里一起看：
[Lesson 56](56-use-session-compare-table-and-download.md)
