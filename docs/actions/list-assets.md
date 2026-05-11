# Action: `list-assets`

`list-assets` 会把当前二进制里内置的资源路径列出来，适合先确认 release 里到底带了什么。

## 最小命令

```bash
go run . -action list-assets
```

## 适合什么时候用

- 想确认 `docs/`、`script/`、`demo/` 有没有一起打进二进制
- 想给 release smoke check 留一份资源清单
- 想先判断用户是不是缺资源，而不是先怀疑命令本身

## 输出结果

- 逐行打印资源路径
- 适合重定向到文件留证据

## 注意事项

- 它只列路径，不会把资源释放到本地
- 如果你下一步要实际打开这些文件，通常再用 [extract-assets](extract-assets.md)

## 相关文档

- [Lesson 142](../tutorials/142-list-assets-for-beginners.md)
- [Lesson 148](../tutorials/148-first-run-entry-strategy.md)
