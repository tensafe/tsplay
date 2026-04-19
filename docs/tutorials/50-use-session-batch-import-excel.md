# Lesson 50: 用命名会话驱动 Excel 批量导入

这一节把认证过的批量流程再往前推进一步。  
输入不再是 CSV，而是 Excel。

目标：

- `use_session`
- `read_excel`
- `foreach`
- 受保护流程里的 Excel 批量导入

## 开始前

建议先跑完：

- [Lesson 46](46-save-import-session.md)
- [Lesson 48](48-use-session-batch-import-csv.md)

默认输入文件：

- [../../demo/data/import_users.xlsx](../../demo/data/import_users.xlsx)

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/50_use_session_batch_import_excel.lua](../../script/tutorials/50_use_session_batch_import_excel.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/50_use_session_batch_import_excel.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/50_use_session_batch_import_excel.lua
```

预期结果：

- 会生成 `artifacts/tutorials/50-use-session-batch-import-excel-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/50_use_session_batch_import_excel.flow.yaml](../../script/tutorials/50_use_session_batch_import_excel.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/50_use_session_batch_import_excel.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/50_use_session_batch_import_excel.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/50-use-session-batch-import-excel-flow.json`

## Step 3: 这节要理解什么

做到这里，整个认证导入链路已经比较完整了：

- 会登录
- 会保存状态
- 会注册命名会话
- 会在受保护页面里做单条和批量导入
- 会把 CSV 和 Excel 都接起来

这就是一个很像真实业务自动化的小闭环。

## 下一节

下一轮最自然的继续方向，就是把 Excel 坏数据恢复也接进这条认证主线。
