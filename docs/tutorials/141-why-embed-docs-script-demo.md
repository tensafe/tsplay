# Lesson 141: 为什么要把 `ReadMe.md`、`docs/`、`script/`、`demo/` 一起打进二进制

`Lesson 121-140` 先把安全边界、review、组织方式讲清楚了。  
这一节开始进入下一条高级主线：

- 发布包
- 内置资产
- 单二进制交付

目标：

- 先用真实命令证明二进制里到底带了什么
- 理解为什么交付时不能只发运行时，不发参考资料
- 为后面的 `list-assets`、`extract-assets`、`file-srv` 铺路

## 准备工作

先确认输出目录存在：

```bash
mkdir -p artifacts/tutorials
```

配套说明：

- [../../script/tutorials/release_pack/checklists/141_embedded_asset_inventory.md](../../script/tutorials/release_pack/checklists/141_embedded_asset_inventory.md)
- [../../embedded_assets.go](../../embedded_assets.go)

## Step 1: 先把二进制内置资源列出来

```bash
# 方式 A：直接运行源码
go run . -action list-assets > artifacts/tutorials/141-bundled-assets.txt

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action list-assets > artifacts/tutorials/141-bundled-assets.txt
```

预期结果：

- 会生成 `artifacts/tutorials/141-bundled-assets.txt`
- 里面至少能找到：
  - `ReadMe.md`
  - `docs/tutorials/README.md`
  - `script/tutorials/01_hello_world.flow.yaml`
  - `demo/demo.html`

## Step 2: 再用一个最小过滤确认“四类入口”都在

```bash
rg -n '^(ReadMe.md|docs/tutorials/README.md|script/tutorials/01_hello_world.flow.yaml|demo/demo.html)$' artifacts/tutorials/141-bundled-assets.txt
```

预期结果：

- 这四条至少都会被匹配到

## Step 3: 这一节意味着什么

把 `ReadMe.md`、`docs/`、`script/`、`demo/` 一起打进二进制，不只是“省事一点”。  
它解决的是交付一致性：

- 运行入口在二进制里
- 教程入口在二进制里
- 示例脚本在二进制里
- demo 页面也在二进制里

## 下一步

继续看：
[Lesson 142](142-list-assets-for-beginners.md)
