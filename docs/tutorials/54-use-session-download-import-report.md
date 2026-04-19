# Lesson 54: 下载认证导入页当前导出的 CSV

这一节开始进入“页面自带导出能力”。

目标：

- `use_session`
- `download_file`
- 从认证页面直接拿到当前导出文件

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)
- [Lesson 48](48-use-session-batch-import-csv.md)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/54_use_session_download_import_report.lua](../../script/tutorials/54_use_session_download_import_report.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/54_use_session_download_import_report.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/54_use_session_download_import_report.lua
```

预期结果：

- 会生成 `artifacts/tutorials/54-use-session-import-report-lua.csv`
- 会生成 `artifacts/tutorials/54-use-session-download-import-report-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/54_use_session_download_import_report.flow.yaml](../../script/tutorials/54_use_session_download_import_report.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/54_use_session_download_import_report.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/54_use_session_download_import_report.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/54-use-session-import-report-flow.csv`
- 会生成 `artifacts/tutorials/54-use-session-download-import-report-flow.json`

## Step 3: 这节要理解什么

这一节先只学一件事：

- 页面导出的文件，能不能真正被你接住

先把下载落地，再谈内容校验，会更稳。

## 下一节

下一节把下载下来的 CSV 再读回来：
[Lesson 55](55-use-session-download-import-report-readback.md)
