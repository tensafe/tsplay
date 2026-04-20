# Lesson 144: 设计“只发一个二进制给用户”的交付流程

`Lesson 141-143` 先把“内置了什么”和“怎么释放出来”讲清楚了。  
这一节开始把这些动作收成一条真正的交付流程。

目标：

- 生成一份单二进制交付 manifest
- 明确“用户拿到二进制后的第一组命令”
- 让交付动作不再只靠口头说明

## 准备工作

样例文件：

- Flow:
  [../../script/tutorials/144_single_binary_delivery_manifest.flow.yaml](../../script/tutorials/144_single_binary_delivery_manifest.flow.yaml)
- Checklist:
  [../../script/tutorials/release_pack/checklists/144_single_binary_delivery_checklist.md](../../script/tutorials/release_pack/checklists/144_single_binary_delivery_checklist.md)

## Step 1: 先生成单二进制交付 manifest

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/144_single_binary_delivery_manifest.flow.yaml

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/144_single_binary_delivery_manifest.flow.yaml
```

预期结果：

- 会写出 `artifacts/tutorials/144/single-binary-delivery-manifest.json`

## Step 2: 再对照 checklist

这一步重点不是再跑别的能力，而是确认 manifest 里有没有把这些动作交代清楚：

- `list-assets`
- `extract-assets`
- 教程入口
- 最小 runnable 示例

## Step 3: 这一节的最小结论

“只发一个二进制”不是一句口号。  
它至少要回答：

- 用户拿到什么
- 用户第一步敲什么
- 用户怎么找到教程

## 下一步

继续看：
[Lesson 145](145-offline-learning-delivery-flow.md)
