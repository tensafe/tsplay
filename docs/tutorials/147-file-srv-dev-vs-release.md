# Lesson 147: `file-srv` 的开发态和发布态有什么区别

`Lesson 141-146` 讲的是“资源怎么进二进制、怎么被交付”。  
这一节回到一个非常具体的入口：`file-srv`。

目标：

- 理解开发态和发布态两种用法
- 知道什么时候该传 `-serve-root`
- 知道为什么不传 `-serve-root` 时也能服务 demo

## 准备工作

配套说明：

- [../../script/tutorials/release_pack/checklists/147_file_srv_dev_vs_release.md](../../script/tutorials/release_pack/checklists/147_file_srv_dev_vs_release.md)
- [../../static_server.go](../../static_server.go)

## Step 1: 先看开发态命令

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000 -serve-root .

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000 -serve-root .
```

这适合：

- 本地正在改 demo
- 希望本地文件优先，内置资源兜底

## Step 2: 再看发布态命令

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

这适合：

- 只发一个二进制
- 不要求用户拿到源码目录

## Step 3: 这一节的最小结论

`file-srv` 不是只有一种模式：

- 开发态：本地目录优先
- 发布态：二进制内置资源优先

## 下一步

继续看：
[Lesson 148](148-first-run-entry-strategy.md)
