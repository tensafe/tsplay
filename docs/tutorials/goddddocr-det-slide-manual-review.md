# goddddocr 低置信度人工接管模板

这条 Flow 在点选和滑块动作前加置信度闸门：`det_result.boxes[0].score` 低于 `min_det_score` 时不会点击，`slide_result.confidence` 低于 `min_slide_confidence` 时不会拖动，而是保存证据并输出 `manual_review` 结构化结果。

启动 goddddocr：

```bash
cd ../goddddocr
go run ./cmd/goddddocr-server -addr :8088 -det
```

启动本地 demo：

```bash
go run . -action file-srv -addr :8000 -serve-root .
```

运行模板：

```bash
go run . -flow script/tutorials/goddddocr_det_slide_manual_review.flow.yaml -headless
```

关键产物：

- `artifacts/goddddocr/manual-review-result.json`
- `artifacts/goddddocr/manual-review-evidence.png`
- `artifacts/goddddocr/manual-review-evidence.html`
- `artifacts/goddddocr/manual-review-det-response.json`
- `artifacts/goddddocr/manual-review-slide-response.json`

`manual-review-result.json` 的 `status` 为 `manual_review` 时，`phase` 会说明阻断点是 `detect_score` 还是 `slide_confidence`，`reason` 会记录阈值断言失败详情，`evidence` 里保留截图、HTML 和原始识别响应路径。接真实站点时先把 `min_det_score`、`min_slide_confidence` 调到保守值，再逐步结合业务容忍度放宽。

在 CLI / MCP / Workbench 里，TSPlay 也会把这个 payload 提取成运行结果的标准字段：

```json
{
  "status": "manual_review_required",
  "manual_review": {
    "required": true,
    "action": "manual_review_required",
    "phase": "detect_score",
    "reason": "assert_number failed...",
    "artifacts": [
      {
        "name": "screenshot",
        "relative_path": "artifacts/goddddocr/manual-review-evidence.png",
        "content_type": "image/png",
        "exists": true
      }
    ]
  }
}
```

Workbench 的 `/api/workbench/tasks/run` 会额外给 `manual_review.artifacts[]` 补 `url_path`，业务操作台可以直接用它预览 evidence。
