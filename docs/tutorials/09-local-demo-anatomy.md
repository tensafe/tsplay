# Lesson 09: 启动本地 demo 服务并拆解页面动作

这一节不新增脚本，重点是把 `Lesson 02` 背后的页面动作顺序吃透。

目标：

- 启动 TSPlay 内置静态文件服务
- 用肉眼看懂 `demo/demo.html`
- 复述 `navigate -> wait_for_selector -> select_option -> is_selected`
- 学会用 URL 覆盖去复用同一份 Lua 脚本

## Step 1: 启动本地 demo 服务

在仓库根目录执行：

```bash
# 方式 A：直接运行源码
go run . -action file-srv -addr :8000

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -action file-srv -addr :8000
```

然后访问：

```text
http://127.0.0.1:8000/demo/demo.html
```

## Step 2: 先不用脚本，直接读页面

页面里最重要的元素只有几个：

- 下拉框：`#options`
- 默认选中的项：`value="3"`
- 这节最常用的练习项：`value="5"`
- 还能继续往后练的分组选项：`value="7"`

这一层先用肉眼看清楚，后面写选择器时会轻松很多。

## Step 3: 再回到 `Lesson 02`

Lua 版本：

```bash
# 方式 A：直接运行源码
go run . -script script/tutorials/02_select_option.lua

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -script script/tutorials/02_select_option.lua
```

Flow 版本：

```bash
# 方式 A：直接运行源码
go run . -flow script/tutorials/02_select_option.flow.yaml -headless

# 方式 B：先执行 go build -o tsplay .，再直接用 tsplay 命令
./tsplay -flow script/tutorials/02_select_option.flow.yaml -headless
```

## Step 4: 记住动作顺序，而不是先记 API 名字

这节真正要记住的是顺序：

1. `navigate`
2. `wait_for_selector`
3. `select_option`
4. `is_selected`

这里最容易犯的错，是页面一打开就先 `sleep`。  
新手阶段更推荐先建立这个原则：

- 已知元素会出现，就先 `wait_for_selector`
- 要做动作，就用动作型 API，比如 `select_option`
- 要做验证，就用验证型 API，比如 `is_selected`

所以：

- `select_option` 解决的是“做了没有”
- `is_selected` 解决的是“结果对不对”

## Step 5: 试一次 URL 覆盖

仓库里额外准备了一个备用页面：
[../../demo/demo_alt.html](../../demo/demo_alt.html)

先确认服务还在运行，然后执行：

```bash
TSPLAY_DEMO_URL=http://127.0.0.1:8000/demo/demo_alt.html \
./tsplay -script script/tutorials/02_select_option.lua
```

这样你会发现：

- 没有改脚本逻辑
- 只是换了目标 URL
- 产物还是会正常写到 `artifacts/tutorials/`

如果你想在 Flow 里做同样的练习，最简单的办法就是先把
[../../script/tutorials/02_select_option.flow.yaml](../../script/tutorials/02_select_option.flow.yaml)
里的 `vars.demo_url` 改成备用页地址，再重跑一遍。

## Step 6: 为什么教程先用本地 demo，不先上公网

因为新手阶段最怕的是四种问题同时混在一起：

- 站点变了
- 网络波动
- 权限不够
- 选择器本来就没想明白

先用仓库自带 demo，把变量降到最低，才能确认你学到的是 TSPlay 本身，而不是“刚好这次网络没出问题”。

## 下一节

下一节开始进入真正的新动作：断言页面状态，而不是只做交互。
[Lesson 10](10-assert-page-state.md)
