# Lesson 142: 新手怎么读懂 `list-assets`

`Lesson 141` 先证明了二进制里确实带着一整套参考资料。  
这一节继续往前走：新手第一次看到 `list-assets` 输出时，该怎么读。

目标：

- 知道 `list-assets` 输出的意义
- 学会从输出里快速找到教程入口
- 建立“先找资源，再跑命令”的习惯

## 准备工作

如果上一节还没跑，可以直接执行：

```bash
# 方式 A：直接运行源码
go run . -action list-assets > artifacts/tutorials/142-list-assets.txt

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action list-assets > artifacts/tutorials/142-list-assets.txt
```

## Step 1: 先看教程入口

```bash
rg -n '^docs/tutorials/' artifacts/tutorials/142-list-assets.txt | head -n 20
```

预期结果：

- 你会先看到教程 markdown 的主入口
- 新用户最应该先认这个前缀

## Step 2: 再看示例脚本入口

```bash
rg -n '^script/tutorials/' artifacts/tutorials/142-list-assets.txt | head -n 20
```

预期结果：

- 你会看到 Flow、Lua、参数文件、配套资源

## Step 3: 最后看 demo 入口

```bash
rg -n '^demo/' artifacts/tutorials/142-list-assets.txt | head -n 20
```

预期结果：

- 你会看到本地 demo 页面和数据文件

## Step 4: 这一节的最小结论

对新手来说，`list-assets` 不是“看一大堆文件名”。  
它真正回答的是：

- 我手里的二进制里有没有教程
- 有没有示例脚本
- 有没有 demo 页面

## 下一步

继续看：
[Lesson 143](143-extract-assets-for-beginners.md)
