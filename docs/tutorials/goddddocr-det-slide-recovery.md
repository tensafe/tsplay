# goddddocr 点选与滑块失败恢复模板

这条 Flow 是 `goddddocr_det_slide.flow.yaml` 的恢复版：把点选检测、滑块识别、拖动验证放进 `retry`，每次失败都由 `on_error` 写出诊断 JSON、截图和 HTML，再进入下一次尝试。

启动 goddddocr：

```bash
cd ../goddddocr
go run ./cmd/goddddocr-server -addr :8088 -det
```

启动本地 demo：

```bash
go run . -action file-srv -addr :8000 -serve-root .
```

运行恢复模板：

```bash
go run . -flow script/tutorials/goddddocr_det_slide_recovery.flow.yaml -headless
```

关键产物：

- `artifacts/goddddocr/det-slide-recovery-result.json`
- `artifacts/goddddocr/det-slide-recovery-diagnostic.json`
- `artifacts/goddddocr/det-slide-recovery-failure.png`
- `artifacts/goddddocr/det-slide-recovery-failure.html`

模板里的 `current_phase` 会标记失败发生在 `ready / capture / detect / slide` 哪一段。真实接入时，把选择器和截图路径替换成目标页面；如果失败集中在 `detect`，优先检查检测框和 `click_box.image_path`；如果集中在 `slide`，优先检查 `slide_result.target_x` 是否需要页面比例或滑块起点校准。
