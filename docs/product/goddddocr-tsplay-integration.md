# goddddocr 与 tsplay 集成清单

## 已推进

- `goddddocr` 以独立 Go HTTP 服务运行，tsplay 不嵌入 Python，也不直接加载 ONNX Runtime。
- tsplay 新增 `ocr_ready` Flow action，负责在识别前检查 goddddocr `/ready`。
- tsplay 新增 `ocr_request` Flow action，负责上传本地图片到 goddddocr 兼容服务。
- tsplay 新增验证码登录 demo 与端到端 Flow，覆盖截图、OCR、填表和断言链路。
- `ocr_request` 返回 `text/result/confidence/request_id/processing_time_ms`，业务 Flow 可以直接读取 `{{ocr_result.text}}`；彩色验证码可传 `color_filter_colors` / `color_filter_custom_ranges`，需要排查准确率时可开启 `probability` 返回完整概率矩阵。
- `goddddocr` 支持外部 ONNX 模型和 charset JSON，服务可通过 `GODDDDOCR_MODEL_PATH` / `GODDDDOCR_CHARSET_PATH` 挂载项目私有验证码模型，tsplay 侧仍按同一个 HTTP action 调用。
- `goddddocr` 提供 `ocrdoctor` 本地诊断命令，部署到 Windows、macOS、Linux 后可先验证 ONNX Runtime、模型、charset 和样本识别，再接入 tsplay Flow。
- 安全边界沿用 tsplay 现有策略：需要 `allow_http=true`，读图片和保存响应时需要 `allow_file_access=true`。

## 自定义模型启动示例

```bash
GODDDDOCR_MODEL_PATH=/opt/models/custom.onnx \
GODDDDOCR_CHARSET_PATH=/opt/models/charset.json \
GODDDDOCR_INPUT_NAME=input1 \
GODDDDOCR_OUTPUT_NAME=387 \
goddddocr-server -addr :8088
```

`charset.json` 必须是 JSON 字符串数组，并且第一个元素是 CTC blank，通常是空字符串。模型 tensor 名如果不是 ddddocr 默认的 `input1` / `387`，同步调整 `GODDDDOCR_INPUT_NAME` 和 `GODDDDOCR_OUTPUT_NAME`。

接入 tsplay 前，先在目标机器跑一次 smoke test：

```bash
ocrdoctor -image /opt/models/smoke.png -expect abcd -json
```

只有 `ocrdoctor` 返回 `ok: true` 后，再启动服务并让 Flow 使用 `ocr_ready` / `ocr_request`。

## 待办

- 扩充 Python ddddocr 与 Go goddddocr 的 golden fixtures 对比集，记录更多真实验证码样本的准确率和耗时差异。
- 给 tsplay 文档补一套部署说明：Windows、macOS、Linux 下安装 ONNX Runtime、运行 `ocrdoctor`、启动 goddddocr 服务并配置 `GODDDDOCR_URL`。
- 增加真实站点适配模板：验证码元素定位、刷新验证码、识别失败重试、低置信度人工接管。
- 评估是否需要把 `ocr_request` 扩展为多引擎接口，预留将来接云 OCR 或项目内私有 OCR 服务。
