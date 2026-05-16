# goddddocr 点选与滑块验证码模板

这条 Flow 演示把页面验证码元素截图后，交给本地 goddddocr 服务完成目标检测和滑块匹配：点选框用 `click_box` 点击 `det_result.boxes[0]` 的中心，并通过 `image_path` 自动换算截图像素和页面元素坐标；滑块距离用 `slide_result.target_x` 接到 `drag.delta_x`。

启动 goddddocr：

```bash
cd ../goddddocr
go run ./cmd/goddddocr-server -addr :8088 -det
```

仓库里带了一个本地滑块 demo，先启动静态服务：

```bash
go run . -action file-srv -addr :8000 -serve-root .
```

然后运行模板：

```bash
go run . -flow script/tutorials/goddddocr_det_slide.flow.yaml -headless
```

如果要处理识别失败、点选偏移或滑块未验证后的自动重试，看恢复版：
[goddddocr 点选与滑块失败恢复模板](goddddocr-det-slide-recovery.md)。
如果要在低置信度时停止自动点击/拖动并交给人工处理，看人工接管版：
[goddddocr 低置信度人工接管模板](goddddocr-det-slide-manual-review.md)。

接入真实站点时，把 Flow 里的 `page_url` 和选择器替换成目标页面。

关键产物：

- `artifacts/goddddocr/det-response.json`
- `artifacts/goddddocr/slide-match-response.json`
- `artifacts/goddddocr/det-slide-flow-result.json`

点选类验证码通常使用 `det_result.boxes` 交给 `click_box`，由 TSPlay 计算框中心并点击；检测图来自页面截图时建议同时传 `image_path`，这样 Retina、高 DPR 或浏览器缩放导致的坐标差异会自动校准。滑块类验证码通常使用 `slide_result.target_x` 作为拖动距离，必要时按页面比例或滑块起点做一次偏移校准。
