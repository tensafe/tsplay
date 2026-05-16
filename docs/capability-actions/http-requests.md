# 能力动作类别：HTTP 请求

这组动作把“浏览器自动化”和“直接调接口”接到一起。  
在 Flow / MCP 安全上下文里，重点关注 `allow_http`，而 `save_path` 这类落文件行为还会继续受 `allow_file_access` 约束。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `http_request` | 是 | 是 | 是 | `action: http_request` / `http_request({url=..., method=..., json=...})` | 发起 HTTP 请求。支持 headers、query、json、form、multipart、保存响应文件，以及复用浏览器 cookies / referer / UA。 |
| `ocr_ready` | 是 | 否 | 是 | `action: ocr_ready` + `url` | 调用 goddddocr `/ready`，在识别前确认服务和模型可用。 |
| `ocr_request` | 是 | 否 | 是 | `action: ocr_request` + `file_path,url` | 调用本地或内网的 goddddocr 兼容服务，把图片识别结果整理成 `text/result/confidence`。 |
| `ocr_detect` | 是 | 否 | 是 | `action: ocr_detect` + `file_path,url` | 调用 goddddocr `/det/file`，返回目标框 `result/boxes`。 |
| `ocr_slide_comparison` | 是 | 否 | 是 | `action: ocr_slide_comparison` + `target_file_path,background_file_path` | 调用 goddddocr `/slide_comparison/file`，返回滑块缺口中心点。 |
| `ocr_slide_match` | 是 | 否 | 是 | `action: ocr_slide_match` + `target_file_path,background_file_path` | 调用 goddddocr `/slide_match/file`，返回滑块匹配中心点和置信度。 |
| `json_extract` | 是 | 是 | 是 | `action: json_extract` + `from,path` / `json_extract(value, '$.items[0]')` | 从 JSON 或 JSON 字符串里取值。适合把接口结果再拆成可复用变量。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: http_demo
steps:
  - action: http_request
    url: http://127.0.0.1:8000/demo/data/order_summary.json
    save_as: api_result
    with:
      response_as: json

  - action: json_extract
    from: "{{api_result}}"
    path: "$.body.summary.open"
    save_as: open_count
```

### goddddocr OCR

`ocr_request` 是给 tsplay 调验证码识别服务的轻量封装。它默认读取 `GODDDDOCR_URL`，没有配置时使用 `http://127.0.0.1:8088/ocr/file`；如果传入的是服务根地址，例如 `http://127.0.0.1:8088`，会自动补成 `/ocr/file`。
默认只返回识别文本和置信度；调试识别差异时可以加 `probability: true`，结果会放到 `ocr_result.probability`。
遇到彩色验证码时，可以通过 `color_filter_colors` 或 `color_filter_custom_ranges` 先按 HSV 颜色范围保留目标颜色。

```yaml
schema_version: "1"
name: goddddocr_ocr_demo
steps:
  - action: ocr_ready
    url: http://127.0.0.1:8088
    save_as: ocr_service
    with:
      timeout: 3000

  - action: ocr_request
    url: http://127.0.0.1:8088
    file_path: ../goddddocr/samples/yzm1.png
    save_as: ocr_result
    with:
      charset_range: 0123456789abcdefghijklmnopqrstuvwxyz
      color_filter_colors: [red, blue]
      color_filter_custom_ranges:
        - [[90, 30, 30], [110, 255, 255]]
      confidence: true
      probability: false

  - action: set_var
    save_as: captcha_text
    with:
      value: "{{ocr_result.text}}"
```

### goddddocr Det / Slide

`ocr_detect`、`ocr_slide_comparison`、`ocr_slide_match` 复用同一个 `GODDDDOCR_URL`。传服务根地址时会自动补成对应的 `/det/file`、`/slide_comparison/file` 或 `/slide_match/file`。

```yaml
schema_version: "1"
name: goddddocr_det_slide_demo
steps:
  - action: ocr_detect
    url: http://127.0.0.1:8088
    file_path: ../goddddocr/samples/yzm2.jpeg
    save_as: det_result
    with:
      detailed: true
      score_threshold: 0.2
      nms_threshold: 0.45

  - action: ocr_slide_match
    url: http://127.0.0.1:8088
    target_file_path: artifacts/captcha/slider.png
    background_file_path: artifacts/captcha/background.png
    save_as: slide_result
    with:
      simple_target: true

  - action: set_var
    save_as: drag_x
    with:
      value: "{{slide_result.target_x}}"
```

### Lua

```lua
local response = http_request({
  url = "http://127.0.0.1:8000/demo/data/order_summary.json",
  response_as = "json",
})
local open_count = json_extract(response, "$.body.summary.open")
print(open_count)
```

## 使用建议

- 页面能直接抓 API 时，`http_request` 往往比“继续点页面”更稳定
- goddddocr 这类本地服务建议先跑 `ocr_ready`，失败时更容易定位是服务问题还是图片识别问题
- 验证码识别优先用 `ocr_request`，业务 Flow 只处理 `ocr_result.text` 和必要的置信度判断；彩色验证码可先用颜色过滤，完整概率矩阵只在排查准确率时打开
- 点选类验证码优先用 `ocr_detect` 拿框，再按业务页面坐标体系转换点击点；滑块验证码优先用 `ocr_slide_match`，同尺寸差分图再用 `ocr_slide_comparison`
- `json_extract` 很适合和 `save_as`、`set_var` 串起来，把响应拆成后续步骤要用的字段
- `use_browser_cookies=true` 时，意味着这条请求会依赖浏览器上下文

## 相关教程

- [Lesson 05](../tutorials/05-http-request-json.md)
- [Lesson 71](../tutorials/71-external-system-round-trip.md)
- [Lesson 122](../tutorials/122-allow-http-boundary.md)
- [goddddocr 验证码登录示例](../tutorials/goddddocr-captcha-login.md)
