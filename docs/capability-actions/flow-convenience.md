# 能力动作类别：Flow 便捷动作

这组动作更像“高层语义糖”，它们让 Flow 更可读、更易被 AI 生成，也更适合交付 review。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `extract_text` | 是 | 是 | 是 | `action: extract_text` + `selector` / `extract_text(selector, timeout, pattern)` | 读文本，可先等、可再做一次正则提取。很适合直接配 `save_as`。 |
| `assert_visible` | 是 | 是 | 是 | `action: assert_visible` + `selector` / `assert_visible(selector, timeout)` | 把“元素必须可见”变成明确断言。 |
| `assert_text` | 是 | 是 | 是 | `action: assert_text` + `selector,text` / `assert_text(selector, text, timeout)` | 把“文本必须包含某值”变成明确断言。 |
| `assert_number` | 是 | 否 | 是 | `action: assert_number` + `value,op,expected` | 把分数、置信度、数量这类数值阈值变成明确断言。 |
| `set_var` | 是 | 是 | 是 | `action: set_var` + `save_as` / `set_var(name, value)` | 保存变量。Flow 侧更强调 `save_as` + `value`。 |
| `append_var` | 是 | 是 | 是 | `action: append_var` + `save_as` / `append_var(name, value)` | 追加到列表变量。Flow 侧会自动初始化列表。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: flow_convenience_demo
vars:
  keywords: []
steps:
  - action: extract_text
    selector: "title"
    save_as: page_title

  - action: assert_visible
    selector: "#kw"
    timeout: 5000

  - action: assert_number
    save_as: confidence_gate
    with:
      value: "{{ocr_result.confidence}}"
      op: ">="
      expected: 0.8
      label: OCR confidence

  - action: append_var
    save_as: keywords
    with:
      value: "{{page_title}}"
```

### Lua

```lua
local page_title = extract_text("title")
assert_visible("#kw", 5000)
append_var("keywords", page_title)
```

## 使用建议

- `extract_text + save_as` 很适合把页面值转成后续步骤输入
- `assert_*` 最适合补在关键状态节点，而不是到处堆
- `assert_number` 适合给 OCR 置信度、目标检测 score、导入行数这类数值加闸门
- `set_var / append_var` 可以把 Flow 从“串命令”变成“可读的编排逻辑”

## 相关教程

- [Lesson 04](../tutorials/04-extract-text-and-html.md)
- [Lesson 10](../tutorials/10-assert-page-state.md)
- [Lesson 94](../tutorials/94-build-collect-verify-save-template.md)
