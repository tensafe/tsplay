# TSPlay Flow 中文 Selector 策略速查

## 这份文件适合什么时候看

当用户说“不会写 selector”“这个 selector 老是失效”“页面结构太复杂”“点击不到正确按钮”时，先看这份文件。

这份文件的目标不是教你死记 selector 语法，而是帮你在写 Flow 时优先选更稳定、更贴业务意图的定位方式。

## 一句话原则

先选和业务意图最接近、最稳定的 selector，不要一上来就写很深的 DOM 路径或 XPath。

## TSPlay 推荐的 selector 优先级

按优先顺序，通常优先考虑：

1. `data-testid`
2. `data-cy`
3. `id`
4. `placeholder`
5. `aria-label`
6. `role` 或可见文本
7. 稳定 class 组合
8. XPath 只作为最后兜底

这条优先级特别适合写可维护的 Flow，因为越靠前的 selector 往往越稳定、越接近业务语义。

## 中文理解这套优先级

- `data-testid` / `data-cy`: 最适合自动化，通常最稳
- `id`: 如果页面开发写得规范，也很稳
- `placeholder`: 很适合输入框
- `aria-label`: 适合无明显文本的输入框或按钮
- 文本或 role: 适合按钮、链接、标签页
- class: 只有在 class 足够稳定时才考虑
- XPath: 能不用就别先用

## 什么时候用哪一种

### 1. 输入框

优先考虑：

- `#id`
- `[data-testid="..."]`
- `[placeholder="..."]`
- `[aria-label="..."]`

示例：

```yaml
- action: type_text
  selector: "#keyword"
  text: "{{keyword}}"
```

```yaml
- action: type_text
  selector: '[placeholder="请输入订单号"]'
  text: "{{order_id}}"
```

### 2. 按钮

优先考虑：

- `data-testid`
- `id`
- 文本按钮，比如 `text="搜索"`

示例：

```yaml
- action: click
  selector: 'text="搜索"'
```

### 3. 表格

优先考虑：

- 表格 id
- `data-testid`
- 明确的容器 selector

示例：

```yaml
- action: capture_table
  selector: "#orders-table"
  save_as: orders
```

### 4. 状态区、结果区、提示区

优先考虑：

- 稳定 id
- `data-testid`
- 业务区块 selector

因为这些位置后面经常要配合 `assert_visible` 或 `assert_text`。

## 先看业务意图，不要先看 DOM 深度

坏思路：

- 先看页面层级，硬写一长串 `div > div > span`

好思路：

- 先问自己“我要点的是搜索按钮，还是我要断言的是导入结果区域”
- 再去找最贴这个业务角色的定位方式

## 常见不稳 selector 的特征

下面这些写法更容易脆：

- 过长的层级路径
- 很多 `nth-child`
- 只依赖视觉布局顺序
- 用了临时 class 名
- XPath 写到很深、很长

这类 selector 一旦页面布局微调，就很容易失效。

## 常见更稳 selector 的特征

- 对应明确业务角色
- 不依赖页面层级顺序
- 输入框、按钮、表格、结果区各自有清楚标识
- 别人一看 selector 就能猜到是定位什么

## selector 和等待要一起想

很多问题看起来像 selector 错，其实是页面还没准备好。

优先修法：

- selector 可能晚出现时，配 `wait_for_selector`
- 页面状态有延迟时，配 `assert_visible`、`assert_text`
- 页面偶发抖动时，配 `retry` 或 `wait_until`

示例：

```yaml
- action: wait_for_selector
  selector: "#import-form"
  timeout: 5000

- action: click
  selector: "#submit"
```

## 如果用户不知道 selector，怎么办

不要逼用户自己贴一整页 HTML。

优先路线：

1. 用户给页面 URL 和业务目标
2. 先走 MCP 的 `observe_page`
3. 从 selector candidates 里选最合适的
4. 再落回 Flow

这比让用户手工猜 selector 更稳定。

## 常见中文场景的 selector 策略

### 搜索框

- 优先 `id`
- 其次 `placeholder`
- 再其次 `aria-label`

### 搜索按钮

- 优先 `id` 或 `data-testid`
- 其次 `text="搜索"` 这类业务文本

### 登录框

- 用户名、密码输入框优先 `id`、`placeholder`、`aria-label`
- 登录按钮优先文本或稳定 id

### 表格

- 先找整张表的稳定容器
- 不要先写每个单元格的复杂 selector

### 导入结果或提示条

- 优先结果区 id 或 data-testid
- 再用 `assert_text` 校验业务结果

## selector 失效时的修复顺序

1. 先确认页面是不是到对地方了
2. 再确认是否少了等待
3. 再检查 selector 是否选错角色
4. 最后才考虑改成更复杂的 XPath

## review 时怎么判断 selector 写得好不好

问这几个问题：

1. 这个 selector 一看能不能猜到业务角色
2. 它是不是依赖脆弱层级
3. 如果页面轻微改版，它会不会立刻坏
4. 它是不是和 `wait_for_selector`、`assert_text` 组合得当

## 最后建议

- 先稳，再短
- 先业务语义，再 DOM 技巧
- 先 `observe_page` 和 selector candidates，再考虑自己硬猜
- XPath 留作最后兜底，不要做第一选择
