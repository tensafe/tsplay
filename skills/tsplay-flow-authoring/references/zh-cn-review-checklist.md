# TSPlay Flow 中文 Review 清单

## 这份文件适合什么时候看

当用户说“帮我 review 这条 Flow”“这条 Flow 能跑，但好不好维护”“想让团队更容易接手”时，先看这份文件。

这份清单的重点不是挑语法细枝末节，而是判断一条 Flow 是否可读、可 review、可接手、可修复。

## Review 的总原则

- 不只看“能不能跑”
- 更要看“别人能不能快速看懂、接手、修”
- 优先指出会增加维护成本的地方

## 第一层：只看入口，不看实现细节

先看这 3 件事：

1. `name`
2. `description`
3. 顶层 `browser` 配置

### 检查点

- `name` 能不能直接说明任务意图
- `description` 是不是在说交付结果，而不是泛泛地说“做点什么”
- 会话、超时、浏览器配置是不是放在顶层，而不是散在各个 step 里

### 常见差写法

- `name: tmp_review`
- `description: Do the thing`
- 登录态相关逻辑散在很多步骤里

### 常见好写法

- `name` 直接体现业务任务
- `description` 说清最后产出
- `browser.use_session` 放在顶层

## 第二层：看变量名和输出名

重点看：

- `save_as`
- `set_var`
- `append_var`
- 输出文件名

### 检查点

- `save_as` 有没有表达业务角色
- 变量名是不是一看就能猜到用途
- 输出文件名是不是人能猜到内容

### 差写法

- `x`
- `tmp2`
- `result1`

### 好写法

- `review_payload`
- `write_result`
- `import_results`
- `auth_status`

## 第三层：看 selector 和等待策略

重点看：

- selector 是否稳定
- 是否有必要的等待
- 是否用业务断言而不是只点按钮

### 检查点

- selector 是不是优先用了 `data-testid`、`id`、placeholder、文本等稳定定位
- 是否少了 `wait_for_selector`
- 是否用了 `assert_visible`、`assert_text` 去断业务结果
- 是否一上来就堆了很长的 XPath

### review 时常见问题

- 页面还没准备好就 click
- selector 很脆
- 没有断言业务结果

## 第四层：看 Flow 结构是不是 Flow 原生

重点看：

- 是不是用 Flow 原生 action 就能表达
- 有没有不必要的 Lua 绕路
- 控制流是不是清楚

### 检查点

- 简单编排是不是已经用 `foreach`、`on_error`、`retry`、`wait_until`
- 如果只是在组织步骤、写变量、写文件，为什么还留着 Lua
- 是否把可以抽回 Flow 的逻辑还停留在 Lua 里

### 最小规则

如果一段逻辑主要是在：

- 组织步骤
- 保存变量
- 写本地文件

那它大概率更适合 Flow，而不是 Lua。

## 第五层：看 artifact 路径和输出布局

重点看：

- 输出路径
- 目录组织
- 结果文件是不是稳定

### 检查点

- artifact 路径是不是稳定可预测
- 教程型输出是不是落在 lesson 自己的目录里
- JSON、CSV、截图、manifest 是否放得清楚

### 差写法

- `artifacts/output.json`
- 文件全平铺在一个目录

### 好写法

- `artifacts/tutorials/133/review-layout/output.json`
- 同一任务的输出靠近放置

## 第六层：看错误恢复和修复友好度

重点看：

- 局部失败是否会拖垮整批任务
- 失败信息是否能留下
- Flow 是否利于后续 repair

### 检查点

- 批量流程里是否该用 `on_error`
- 失败路径是否保留 `{{last_error}}`
- step 是否够小、够清楚，方便 trace 和 artifact 定位

## 第七层：看是否贴业务目标

重点看：

- 每一步是不是服务于用户目标
- 有没有多余动作
- 输出是不是用户真正关心的结果

### 检查点

- 步骤有没有明显“为了凑动作而凑动作”
- 最终输出是不是和业务目标对应
- Flow 是否过于依赖页面实现细节，而不是业务结果

## 中文 review 标准问法

可以直接这样提：

```text
帮我 review 这条 TSPlay Flow，重点看可维护性。
- 文件: <flow 文件路径>
- 关注点:
  - name 和 description 是否清楚
  - save_as 是否表达业务角色
  - selector 和等待是否稳
  - 是否用了不必要的 Lua
  - artifact 路径是否可 review
```

## 中文 review 输出建议结构

如果你在做 review，优先按这个顺序反馈：

1. 最影响维护性的点
2. 最影响稳定性的点
3. 最值得保留的优点
4. 最小修复建议

## 最小 checklist

拿到一条 Flow，至少先过这几项：

1. `name` 能不能说明任务意图
2. `description` 能不能说明交付结果
3. `save_as` 是否表达变量角色
4. selector 是否稳定
5. 是否有必要等待和业务断言
6. artifact 路径是否稳定
7. 能否继续用 Flow，而不是不必要地绕到 Lua
8. 交给别人接手时，最困惑的点会是什么

## 一句总结

一条值得长期保留的 TSPlay Flow，不只是能跑，还要让别人看得懂、改得动、接得住。
