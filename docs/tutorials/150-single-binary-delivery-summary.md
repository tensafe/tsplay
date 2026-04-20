# Lesson 150: 复盘“单二进制 + 内置教程”对交付的意义

`Lesson 141-149` 这一整段，已经把：

- 内置资源
- 资源释放
- 单二进制交付
- 离线学习
- `file-srv`
- first-run 入口
- 版本策略

连成了一条主线。  
这一节把它收口成一个最终 summary。

目标：

- 生成一份交付总结 JSON
- 把这一段高级主题压成一个可复盘结论
- 顺势接到后面的 capstone 和培训模块

## 准备工作

样例文件：

- [../../script/tutorials/150_single_binary_delivery_summary.flow.yaml](../../script/tutorials/150_single_binary_delivery_summary.flow.yaml)
- [../../script/tutorials/release_pack/README.md](../../script/tutorials/release_pack/README.md)

## Step 1: 先生成交付总结 JSON

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/150_single_binary_delivery_summary.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/150_single_binary_delivery_summary.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/150/single-binary-delivery-summary.json`

## Step 2: 再回看这一段真正解决了什么

它解决的不是“多了几篇文档”，而是：

- 交付物和教程绑定
- 二进制和参考资料绑定
- 新手入口和高级维护入口绑定

## Step 3: 这一节的最小结论

单二进制 + 内置教程的价值，在于把：

- 运行
- 学习
- 示例
- demo
- 交付

这些原本容易分散的东西，收成一个更稳定的交付物。

## 下一步

继续看：
[Lesson 151](151-newbie-capstone-brief.md)
