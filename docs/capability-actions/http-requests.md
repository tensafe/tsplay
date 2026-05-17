# 能力动作类别：HTTP 请求

这组动作把“浏览器自动化”和“直接调接口”接到一起。  
在 Flow / MCP 安全上下文里，重点关注 `allow_http`；`save_path` 这类落文件行为还会继续受 `allow_file_access` 约束。
goddddocr 的 managed sidecar / direct CLI 模式会启动本地进程，需要额外开启 `allow_process`。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `http_request` | 是 | 是 | 是 | `action: http_request` / `http_request({url=..., method=..., json=...})` | 发起 HTTP 请求。支持 headers、query、json、form、multipart、保存响应文件，以及复用浏览器 cookies / referer / UA。 |
| `ocr_ready` | 是 | 否 | 是 | `action: ocr_ready` + `url` 或 `mode: sidecar` | 调用 goddddocr `/ready`，在识别前确认服务和模型可用；sidecar 模式会自动启动 `goddddocr-server`。 |
| `ocr_request` | 是 | 否 | 是 | `action: ocr_request` + `file_path` | 调用 goddddocr HTTP 服务、managed sidecar 或 direct CLI，把图片识别结果整理成 `text/result/confidence`。 |
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

`ocr_request` 是给 tsplay 调验证码识别服务的轻量封装。它支持三种模式：

- HTTP：显式传 `url`，或配置 `GODDDDOCR_URL`，TSPlay 会调用已有服务；传服务根地址时会自动补成 `/ocr/file`
- managed sidecar：`mode: sidecar`，或在可信本地运行且没有 `url/GODDDDOCR_URL` 时自动启动 `goddddocr-server`，自动选端口并等待 `/ready`
- direct CLI：`mode: cli`，TSPlay 执行本地二进制并解析 stdout JSON，适合低频识别和调试

默认只返回识别文本和置信度；调试识别差异时可以加 `probability: true`，结果会放到 `ocr_result.probability`。
遇到彩色验证码时，可以通过 `color_filter_colors` 或 `color_filter_custom_ranges` 先按 HSV 颜色范围保留目标颜色。

Managed sidecar 示例：

```yaml
schema_version: "1"
name: goddddocr_sidecar_demo
steps:
  - action: ocr_request
    file_path: artifacts/captcha/captcha.png
    save_as: ocr_result
    with:
      mode: sidecar
      startup_timeout: 15000
      confidence: true
```

`mode: sidecar` 默认查找 `GODDDDOCR_SERVER_BIN`、`GODDDDOCR_BIN` 或 PATH 里的 `goddddocr-server`。需要传模型参数时，用 `server_args`：

```yaml
with:
  mode: sidecar
  server_args:
    - -model
    - beta
```

Direct CLI 示例：

```yaml
schema_version: "1"
name: goddddocr_cli_demo
steps:
  - action: ocr_request
    file_path: artifacts/captcha/captcha.png
    save_as: ocr_result
    with:
      mode: cli
      executable: ocrdoctor
      confidence: true
```

如果使用统一二进制，也可以完全指定命令模板：

```yaml
with:
  mode: cli
  executable: goddddocr
  cli_args: [ocr, --file, "{file}", --json, --confidence]
```

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

  - action: click_box
    selector: "#captcha-image"
    with:
      box: "{{det_result.boxes[0]}}"
      image_path: ../goddddocr/samples/yzm2.jpeg

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
- goddddocr 这类本地服务可以先跑 `ocr_ready`；如果 Flow 使用 `mode: sidecar`，TSPlay 会自己启动服务并等待 `/ready`
- 验证码识别优先用 `ocr_request`，业务 Flow 只处理 `ocr_result.text` 和必要的置信度判断；彩色验证码可先用颜色过滤，完整概率矩阵只在排查准确率时打开
- 点选类验证码优先用 `ocr_detect` 拿框，再用 `click_box` 承接 `det_result.boxes[0]`；如果 det 图片是页面元素截图，传 `image_path` 可自动处理设备像素比和截图缩放。滑块验证码优先用 `ocr_slide_match`，同尺寸差分图再用 `ocr_slide_comparison`
- `json_extract` 很适合和 `save_as`、`set_var` 串起来，把响应拆成后续步骤要用的字段
- `use_browser_cookies=true` 时，意味着这条请求会依赖浏览器上下文

## 相关教程

- [Lesson 05](../tutorials/05-http-request-json.md)
- [Lesson 71](../tutorials/71-external-system-round-trip.md)
- [Lesson 122](../tutorials/122-allow-http-boundary.md)
- [goddddocr 验证码登录示例](../tutorials/goddddocr-captcha-login.md)
- [goddddocr 点选与滑块验证码模板](../tutorials/goddddocr-det-slide.md)
- [goddddocr 点选与滑块失败恢复模板](../tutorials/goddddocr-det-slide-recovery.md)
- [goddddocr 低置信度人工接管模板](../tutorials/goddddocr-det-slide-manual-review.md)
