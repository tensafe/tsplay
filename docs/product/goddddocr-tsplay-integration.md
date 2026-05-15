# goddddocr 与 tsplay 集成清单

## 已推进

- `goddddocr` 以独立 Go HTTP 服务运行，tsplay 不嵌入 Python，也不直接加载 ONNX Runtime。
- tsplay 新增 `ocr_ready` Flow action，负责在识别前检查 goddddocr `/ready`。
- tsplay 新增 `ocr_request` Flow action，负责上传本地图片到 goddddocr 兼容服务。
- tsplay 新增验证码登录 demo 与端到端 Flow，覆盖截图、OCR、填表和断言链路。
- `ocr_request` 返回 `text/result/confidence/request_id/processing_time_ms`，业务 Flow 可以直接读取 `{{ocr_result.text}}`；需要排查准确率时可开启 `probability` 返回完整概率矩阵。
- 安全边界沿用 tsplay 现有策略：需要 `allow_http=true`，读图片和保存响应时需要 `allow_file_access=true`。

## 待办

- 扩充 Python ddddocr 与 Go goddddocr 的 golden fixtures 对比集，记录更多真实验证码样本的准确率和耗时差异。
- 给 tsplay 文档补一套部署说明：Windows、macOS、Linux 下启动 goddddocr 服务并配置 `GODDDDOCR_URL`。
- 增加真实站点适配模板：验证码元素定位、刷新验证码、识别失败重试、低置信度人工接管。
- 评估是否需要把 `ocr_request` 扩展为多引擎接口，预留将来接云 OCR 或项目内私有 OCR 服务。
