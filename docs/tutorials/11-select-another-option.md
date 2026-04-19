# Lesson 11: 改选另一个选项并验证结果

这一节仍然使用本地下拉框 demo，不过目标从“跑通”变成“改值再验证”。

目标：

- 不再固定选择 `选项 5`
- 改成选择另一个值
- 继续把结果写到 `artifacts/tutorials/`

这一节会默认选择 `value="7"`，也就是分组选项里的第二项。

## 准备工作

先确认 TSPlay 内置静态文件服务还在运行。  
如果没有运行，就在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

页面地址：

```text
http://127.0.0.1:8000/demo/demo.html
```

## Step 1: 运行 Lua 版本

示例文件：
[../../script/tutorials/11_select_another_option.lua](../../script/tutorials/11_select_another_option.lua)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/11_select_another_option.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/11_select_another_option.lua
```

预期结果：

- 会把下拉框切到 `value="7"`
- 会生成 `artifacts/tutorials/11-select-another-option-lua.json`

如果你想再换一个值，可以直接覆盖环境变量：

```bash
TSPLAY_OPTION_VALUE=5 ./tsplay -script script/tutorials/11_select_another_option.lua
```

## Step 2: 看这次为什么一定要保留验证

这节和 Lesson 02 的差异不在于动作本身有多复杂，  
而在于你开始把“动作”和“验证”明确拆开了：

- `select_option` 负责执行选择
- `is_selected` 负责确认目标选项真的被选中了

这一步很小，但它会直接影响你后面写业务流程时，是否能快速定位“到底是没点到，还是点到了但页面状态没更新”。

## Step 3: 运行 Flow 版本

示例文件：
[../../script/tutorials/11_select_another_option.flow.yaml](../../script/tutorials/11_select_another_option.flow.yaml)

运行命令：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/11_select_another_option.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/11_select_another_option.flow.yaml -headless
```

如果你想改成别的值，最直接的做法是先把这份 YAML 里的 `target_value` 改掉，再重新运行。

## Step 4: 再做一次 URL 覆盖

如果你想证明“同一份逻辑可以打到另一个本地页面”，可以执行：

```bash
TSPLAY_DEMO_URL=http://127.0.0.1:8000/demo/demo_alt.html \
TSPLAY_OPTION_VALUE=7 \
./tsplay -script script/tutorials/11_select_another_option.lua
```

这样你会把同样的选择逻辑，跑到备用 demo 页上。

## 下一节

下一节继续复用这次交互，不过重点改成“把结果整理成你自己想要的 JSON 结构”：
[Lesson 12](12-custom-json-output.md)
