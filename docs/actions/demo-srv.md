# Action: `demo-srv`

`demo-srv` 是 `file-srv` 的兼容别名，两者走的是同一套实现。

## 最小命令

```bash
go run . -action demo-srv -addr :8000
```

## 建议怎么用

- 如果你只是兼容旧命令，继续用 `demo-srv` 也可以
- 如果你在写新文档、新教程或新脚本，更推荐统一写成 [file-srv](file-srv.md)

## 为什么更推荐 `file-srv`

- 现在仓库里的教程、快速开始、release 说明基本都以 `file-srv` 为主
- 统一名字后，培训、交付和排障时更容易对齐

## 相关文档

- [file-srv](file-srv.md)
- [Lesson 147](../tutorials/147-file-srv-dev-vs-release.md)
