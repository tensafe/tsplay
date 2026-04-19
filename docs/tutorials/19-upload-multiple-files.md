# Lesson 19: 上传多个本地文件

这一节继续讲上传，不过把单文件换成多文件。  
我们使用的是 [../../demo/multi_upfile.html](../../demo/multi_upfile.html)。

目标：

- 一次上传两个本地文件
- 验证页面都已经显示出来
- 把结果写到 `artifacts/tutorials/`

## 开始前

如果你在仓库根目录学习，示例文件已经在：

- [../../demo/data/upload_receipt.pdf](../../demo/data/upload_receipt.pdf)
- [../../demo/data/upload_avatar.png](../../demo/data/upload_avatar.png)

如果你只有单个 `./tsplay` 二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

如果你已经进入了 `./tsplay-assets` 目录，运行示例时请把命令里的 `./tsplay` 换成你的实际二进制路径，例如：

```bash
../tsplay -script script/tutorials/19_upload_multiple_files.lua
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/19_upload_multiple_files.lua](../../script/tutorials/19_upload_multiple_files.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/19_upload_multiple_files.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/19_upload_multiple_files.lua
```

预期结果：

- 会生成 `artifacts/tutorials/19-upload-multiple-files-lua.json`

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/19_upload_multiple_files.flow.yaml](../../script/tutorials/19_upload_multiple_files.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/19_upload_multiple_files.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/19_upload_multiple_files.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/19-upload-multiple-files-flow.json`

## Step 3: 这节和上一节的关键差异

- 单文件上传更像“一个输入框塞一个文件”
- 多文件上传更像“一个输入框接收一个文件列表”

这也是为什么 Flow 这里会显式出现：

- `files`

而不是继续只有一个 `file_path`。

## 下一节

下一节把文件动作切到另一边：从页面下载报表到本地，然后再读回来验证。
[Lesson 20](20-download-report.md)
