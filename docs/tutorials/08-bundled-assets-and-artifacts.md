# Lesson 08: 理解内置资源和 `artifacts/` 输出目录

这一节不引入新网页动作，重点补两件对新手特别重要的事：

- 为什么一个 `./tsplay` 二进制就能带着 `docs/`、`script/`、`demo/`
- 为什么教程默认把结果写到 `artifacts/tutorials/`

这一节更像“命令练习 + 心智模型补齐”。

## Step 1: 先构建一次 `tsplay`

在仓库根目录执行：

```bash
go build -o tsplay .
```

如果你暂时不想构建，也可以继续使用 `go run .`。  
但这一节建议你先拿到一个真实的 `./tsplay`，因为接下来要演示“单个二进制也能带着教程走”。

## Step 2: 看看二进制里内置了什么

先列出内置资源：

```bash
./tsplay -action list-assets
```

你应该能看到这些前缀：

- `ReadMe.md`
- `docs/`
- `script/`
- `demo/`

这就是为什么把 `tsplay` 单独带走之后，仍然还能继续跑教程示例。

## Step 3: 把内置资源释放到本地看看

执行：

```bash
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

释放完成后，重点看这几个目录：

- `tsplay-assets/docs/tutorials/`
- `tsplay-assets/script/tutorials/`
- `tsplay-assets/demo/`

这一步的意义不是“以后必须先释放再跑”，而是先帮你确认：二进制里确实已经带着这些参考资料。

## Step 4: 从仓库外目录调用内置示例

假设你刚构建出的二进制路径是：

```text
/absolute/path/to/tsplay
```

现在切到一个空目录，再直接调用内置示例：

```bash
mkdir -p /tmp/tsplay-lesson-08
cd /tmp/tsplay-lesson-08

/absolute/path/to/tsplay -script script/tutorials/01_hello_world.lua
/absolute/path/to/tsplay -flow script/tutorials/01_hello_world.flow.yaml
```

预期结果：

- 当前目录下会生成 `artifacts/tutorials/01-hello-world-lua.json`
- 当前目录下会生成 `artifacts/tutorials/01-hello-world-flow.json`

这就说明：

- 输入可以来自二进制内置的 `script/tutorials/...`
- 输出会落到你当前工作的目录
- `artifacts/tutorials/` 是教程默认产物目录，不会和源码目录混在一起

## Step 5: 把“输入、输出、产物”三件事分清楚

新手很容易把这三者混成一团。  
这里建议先形成一个固定心智：

- 输入：`script/tutorials/...` 或 `docs/tutorials/...` 里的示例和说明
- 执行器：`./tsplay` 或 `go run .`
- 输出：你这次运行产生的 `artifacts/tutorials/...`

所以教程里反复出现的 `artifacts/tutorials/`，本质上是在提醒你：

- 脚本和文档是“参考输入”
- 运行结果是“练习产物”
- 两者要分开放，后续排错和复盘才会轻松

## 下一节

下一节开始把本地 demo 服务和页面结构彻底拆开讲清楚：
[Lesson 09](09-local-demo-anatomy.md)
