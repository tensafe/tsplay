# Lesson 18: 上传单个本地文件

这一节开始讲文件上传。  
我们复用仓库里的 [../../demo/upload.html](../../demo/upload.html)。

目标：

- 给文件输入框设置一个本地文件
- 断言页面已经显示出文件名
- 把结果写到 `artifacts/tutorials/`

## 开始前

这节既需要浏览器页面，也需要本地文件。  
所以你要同时满足两件事：

1. TSPlay 内置静态文件服务正在运行
2. 本地真的存在待上传文件

如果你是在仓库根目录学习，示例文件已经在：
[../../demo/data/upload_receipt.pdf](../../demo/data/upload_receipt.pdf)

如果你只有单个 `./tsplay` 二进制，先执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
cd ./tsplay-assets
```

如果你已经进入了 `./tsplay-assets` 目录，运行示例时请把命令里的 `./tsplay` 换成你的实际二进制路径，例如：

```bash
../tsplay -script script/tutorials/18_upload_single_file.lua
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/18_upload_single_file.lua](../../script/tutorials/18_upload_single_file.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/18_upload_single_file.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/18_upload_single_file.lua
```

预期结果：

- 会生成 `artifacts/tutorials/18-upload-single-file-lua.json`
- 页面上的 `#fileInfo` 会出现 `upload_receipt.pdf`

如果你想换自己的文件，可以一起覆盖：

```bash
TSPLAY_UPLOAD_FILE=/absolute/path/to/file.pdf \
TSPLAY_UPLOAD_FILENAME=file.pdf \
./tsplay -script script/tutorials/18_upload_single_file.lua
```

## Step 2: 运行 Flow 版本

示例文件：
[../../script/tutorials/18_upload_single_file.flow.yaml](../../script/tutorials/18_upload_single_file.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/18_upload_single_file.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/18_upload_single_file.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/18-upload-single-file-flow.json`

## Step 3: 这节要记住什么

- `upload_file` 解决的是“把本地文件交给文件输入框”
- 真正的验证还是要看页面反馈
- 所以最自然的组合是 `upload_file + assert_text`

## 下一节

下一节继续上传，但改成一次上传多个文件。
[Lesson 19](19-upload-multiple-files.md)
