# 能力动作类别：页面原子动作

这页收的是最基础、最常用、最应该在 `Flow / Lua / MCP` 三边保持一致的一组页面动作。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `navigate` | 是 | 是 | 是 | `action: navigate` + `url` / `navigate(url)` | 打开页面。更适合把超时放在浏览器配置或 MCP 调用层，而不是塞到 step 本身。 |
| `click` | 是 | 是 | 是 | `action: click` + `selector` / `click(selector)` | 点击元素。适合与 `wait_for_selector`、`retry` 搭配。 |
| `click_at` | 是 | 是 | 是 | `action: click_at` + `selector,x,y` / `click_at(selector, x, y)` | 点击元素内部相对坐标。适合点选验证码、Canvas、图片热区。 |
| `click_box` | 是 | 是 | 是 | `action: click_box` + `selector,box` / `click_box(selector, box)` | 点击检测框中心。可直接承接 `ocr_detect` 返回的 `boxes`，Flow 里可用 `image_path` 自动换算截图和页面缩放。 |
| `reload` | 是 | 是 | 是 | `action: reload` / `reload()` | 刷新当前页。 |
| `go_back` | 是 | 是 | 是 | `action: go_back` / `go_back()` | 浏览器后退。 |
| `go_forward` | 是 | 是 | 是 | `action: go_forward` / `go_forward()` | 浏览器前进。 |
| `type_text` | 是 | 是 | 是 | `action: type_text` + `selector,text` / `type_text(selector, text)` | 在输入框里输入文本。`fill`、`type` 可以视作常见别名心智模型。 |
| `set_value` | 是 | 是 | 是 | `action: set_value` + `selector,value` / `set_value(selector, value)` | 直接设置元素值，适合不希望逐字输入的场景。 |
| `select_option` | 是 | 是 | 是 | `action: select_option` + `selector,value` / `select_option(selector, value)` | 选择下拉项。 |
| `drag` | 是 | 是 | 是 | `action: drag` + `selector,delta_x` / `drag(selector, dx, dy, steps)` | 按像素偏移拖动元素，适合滑块验证码、拖拽排序这类动作。 |
| `hover` | 是 | 是 | 是 | `action: hover` + `selector` / `hover(selector)` | 鼠标悬停。适合下拉菜单、悬浮操作。 |
| `scroll_to` | 是 | 是 | 是 | `action: scroll_to` + `selector` / `scroll_to(selector)` | 滚动到目标元素。 |
| `wait_for_network_idle` | 是 | 是 | 是 | `action: wait_for_network_idle` / `wait_for_network_idle()` | 等待页面请求基本稳定。适合提交后、跳转后收口。 |
| `wait_for_selector` | 是 | 是 | 是 | `action: wait_for_selector` + `selector` / `wait_for_selector(selector, timeout)` | 等元素出现。是最常见的页面同步动作。 |
| `wait_for_text` | 是 | 是 | 是 | `action: wait_for_text` + `selector,text` / `wait_for_text(selector, text, timeout)` | 等文本出现。适合状态文字、提示语。 |
| `sleep` | 是 | 是 | 是 | `action: sleep` + `seconds` / `sleep(seconds)` | 硬等待。能不用就尽量不用，优先换成显式页面状态等待。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: page_primitive_demo
steps:
  - action: navigate
    url: https://www.baidu.com

  - action: wait_for_selector
    selector: "#kw"
    timeout: 5000

  - action: type_text
    selector: "#kw"
    text: "TSPlay"

  - action: click
    selector: "#su"

  - action: click_box
    selector: "#captcha"
    with:
      box:
        x1: 28
        y1: 28
        x2: 76
        y2: 76
      image_path: artifacts/captcha/det-source.png

  - action: drag
    selector: "#slider-handle"
    with:
      delta_x: 120
      delta_y: 0
      move_steps: 24
```

### Lua

```lua
navigate("https://www.baidu.com")
wait_for_selector("#kw", 5000)
type_text("#kw", "TSPlay")
click("#su")
click_at("#captcha", 52, 52)
click_box("#captcha", {x1=28, y1=28, x2=76, y2=76})
drag("#slider-handle", 120, 0, 24)
```

## 使用建议

- 先用 `wait_for_selector` 建页面同步，再点 `click / type_text / select_option`
- 点选类验证码可以先用 `ocr_detect` 拿 `det_result.boxes`，再用 `click_box` 点击目标框中心；如果检测图片来自 `screenshot_element`，把截图路径传给 `image_path`，TSPlay 会按图片尺寸和元素尺寸自动换算缩放
- 滑块类动作优先把识别出的距离接到 `drag.delta_x`，再用 `move_steps` 控制拖动平滑度
- 页面容易抖动时，优先 `wait_for_selector + retry`，不要直接堆 `sleep`
- 新手路线里，这组动作通常是最先需要跑熟的一层

## 相关教程

- [Lesson 01](../tutorials/01-hello-world.md)
- [Lesson 02](../tutorials/02-local-page-select-option.md)
- [Lesson 10](../tutorials/10-assert-page-state.md)
- [新手学习路线](../tutorials/track-newbie.zh-CN.md)
