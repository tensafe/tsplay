# Lesson 12: 把页面交互结果整理成自定义 JSON

这一节继续复用本地下拉框页面，不过重点已经不再是“会选”或“会验”，而是：

- 能不能把结果组织成你自己想要的结构
- 能不能让输出更接近真实交付物

这一节会把结果写成一个带嵌套字段的 JSON。

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/12_custom_json_output.lua](../../script/tutorials/12_custom_json_output.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/12_custom_json_output.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/12_custom_json_output.lua
```

预期结果：

- 会生成 `artifacts/tutorials/12-custom-json-output-lua.json`
- 输出里不只是平铺字段，还会有 `source` 和 `selection` 这类嵌套结构

## Step 2: 看这节为什么比“直接 write_json 一把梭”更重要

教程前面几节更多是在证明“TSPlay 能拿到数据”。  
但真实交付里更常见的问题其实是：

- 产物要给谁看
- 字段该怎么命名
- 哪些信息要放一层，哪些要嵌套

所以这一节的重点不在动作数量，而在输出组织方式。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/12_custom_json_output.flow.yaml](../../script/tutorials/12_custom_json_output.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/12_custom_json_output.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/12_custom_json_output.flow.yaml -headless
```

预期结果：

- 会生成 `artifacts/tutorials/12-custom-json-output-flow.json`
- 终端会看到 `set_var -> write_json` 的结构化 trace

## Step 4: 自己再改一版

建议你至少再做一次主动改动：

- 把 `target_value` 从 `5` 改成 `7`
- 或者把 `summary` 文案改成你自己的表述
- 或者给 `source` 再加一个字段

只要你主动改过一次，这一节就从“看懂了”变成“真的会用了”。

## 下一步

到这里，新手阶段关于“内置资源、本地 demo、断言、验证、自定义产物”的链路就完整了。  
如果你要继续往外部系统走，可以回到：

- [Lesson 06](06-redis-round-trip.md)
- [Lesson 07](07-db-postgres-basics.md)

如果你要继续沿着课程体系推进，可以进入：
[track-junior.md](track-junior.md)
