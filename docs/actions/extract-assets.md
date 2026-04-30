# Action: `extract-assets`

`extract-assets` 会把二进制里的内置 docs、script、demo 释放到本地目录，适合离线学习、培训包和单二进制交付。

## 最小命令

```bash
go run . -action extract-assets -extract-root ./tsplay-assets
```

## 常用参数

- `-extract-root`：释放目标目录，默认是 `tsplay-assets`

## 适合什么时候用

- 用户手上只有一个二进制，但还需要配套教程和示例
- 培训前想先把参考资料统一解到本地
- 想对 release 产物做离线 smoke check

## 输出结果

- 命令行会返回释放到哪个目录
- 释放后通常会看到 `docs/`、`script/`、`demo/` 等内容

## 注意事项

- 想直接服务内置 demo 时，不一定要先解压；也可以直接用 [file-srv](file-srv.md)
- 如果目录已存在，建议先确认是否要覆盖旧内容

## 相关文档

- [Lesson 143](../tutorials/143-extract-assets-for-beginners.md)
- [Lesson 144](../tutorials/144-single-binary-delivery-flow.md)
- [Lesson 145](../tutorials/145-offline-learning-delivery-flow.md)
