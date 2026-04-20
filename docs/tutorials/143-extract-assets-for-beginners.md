# Lesson 143: 新手怎么跑通 `extract-assets`

`Lesson 142` 先学会了怎么看二进制里有什么。  
这一节继续下一步：把这些内置资源释放到本地目录，方便学习和交付。

目标：

- 跑通一次 `extract-assets`
- 确认释放后的目录结构
- 建立“二进制内置资源”和“本地释放目录”的对应关系

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

## Step 1: 先把二进制资源释放到本地

```bash
# 方式 A：直接运行源码
go run . -action extract-assets -extract-root ./tsplay-assets-143

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action extract-assets -extract-root ./tsplay-assets-143
```

预期结果：

- 会生成 `./tsplay-assets-143`
- 目录里会出现：
  - `ReadMe.md`
  - `docs/`
  - `script/`
  - `demo/`

## Step 2: 再确认释放后的关键入口

```bash
rg -n '^(ReadMe.md|docs/tutorials/README.md|script/tutorials/01_hello_world.flow.yaml|demo/demo.html)$' <(cd tsplay-assets-143 && find . -type f | sed 's#^./##')
```

如果你不想用进程替换，也可以拆成两步：

```bash
cd tsplay-assets-143
find . -type f | sed 's#^./##' > ../artifacts/tutorials/143-extracted-assets.txt
cd ..
rg -n '^(ReadMe.md|docs/tutorials/README.md|script/tutorials/01_hello_world.flow.yaml|demo/demo.html)$' artifacts/tutorials/143-extracted-assets.txt
```

## Step 3: 这一节意味着什么

`extract-assets` 让“单二进制交付”和“本地可阅读目录”同时成立：

- 没有源码目录，也能拿到完整参考资料
- 但释放出来后，又能像普通文件一样查看、检索、培训

## 下一步

继续看：
[Lesson 144](144-single-binary-delivery-flow.md)
