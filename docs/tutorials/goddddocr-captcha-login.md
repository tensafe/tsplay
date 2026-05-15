# goddddocr 验证码登录示例

这条 Flow 演示从页面截取验证码图片、调用本地 goddddocr 服务识别、再把结果填回登录表单。

## 运行

启动 goddddocr：

```bash
cd ../goddddocr
go run ./cmd/goddddocr-server -addr :8088
```

启动 tsplay 本地静态页：

```bash
cd ../tsplay
go run . -action file-srv -addr :8000 -serve-root .
```

运行 Flow：

```bash
cd ../tsplay
go run . -flow script/tutorials/goddddocr_login.flow.yaml -headless
```

## 产物

- `artifacts/goddddocr/login-captcha.png`
- `artifacts/goddddocr/login-ready-response.json`
- `artifacts/goddddocr/login-ocr-response.json`
- `artifacts/goddddocr/login-flow-result.json`
