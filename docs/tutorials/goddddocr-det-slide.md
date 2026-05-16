# goddddocr 点选与滑块验证码模板

这条 Flow 演示把页面验证码元素截图后，交给本地 goddddocr 服务完成目标检测和滑块匹配，再把 `slide_result.target_x` 接到 `drag.delta_x`。

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

接入真实站点时，把 Flow 里的 `page_url` 和选择器替换成目标页面。

关键产物：

- `artifacts/goddddocr/det-response.json`
- `artifacts/goddddocr/slide-match-response.json`
- `artifacts/goddddocr/det-slide-flow-result.json`

点选类验证码通常使用 `det_result.boxes` 生成点击点；滑块类验证码通常使用 `slide_result.target_x` 作为拖动距离，必要时按页面比例或滑块起点做一次偏移校准。
