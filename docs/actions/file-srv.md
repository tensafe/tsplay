# Action: `file-srv`

`file-srv` 用来启动内置静态文件服务。它既能服务仓库里的本地 demo，也能在单二进制发布态下直接服务内置资源。

## 最小命令

```bash
go run . -action file-srv -addr :8000
```

## 常见用法

```bash
go run . -action file-srv -addr :8000 -serve-root .
```

## 常用参数

- `-addr`：监听地址
- `-serve-root`：可选，本地目录优先；不传时回退到二进制内置资源

## 适合什么时候用

- 教程练习前先起本地 demo
- 想在 release 包里直接服务 `demo/`
- 想比较开发态和发布态的资源来源

## 运行后会看到什么

- 命令行会打印 demo 地址
- 常见页面包括 `/demo/demo.html`、`/demo/tables.html`、`/demo/extract.html`

## 注意事项

- 本地正在改 demo 时，更建议传 `-serve-root`
- 做单二进制交付演示时，通常不传 `-serve-root`

## 相关文档

- [Lesson 09](../tutorials/09-local-demo-anatomy.md)
- [Lesson 147](../tutorials/147-file-srv-dev-vs-release.md)
- [快速开始](../../getting-started.md)
