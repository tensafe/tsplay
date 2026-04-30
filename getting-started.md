# Quick Start / 快速开始

This repository keeps the full docs-site quick start page in
[site-src/getting-started.md](site-src/getting-started.md).
This root-level alias exists so links from `docs/` and GitHub-rendered
Markdown keep working outside the generated site.

完整的文档站快速开始页在
[site-src/getting-started.md](site-src/getting-started.md)。
这里保留一个根目录入口，主要是为了让 `docs/` 里的相对链接在仓库源码视图中也能正常打开。

## Smallest First Run

If you want the binary-first path, use the release installer script first:

```bash
curl -fsSL https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.sh | sh
```

On Windows, download and run the PowerShell installer:

```powershell
Invoke-WebRequest https://github.com/tensafe/tsplay/releases/latest/download/tsplay-quickstart.ps1 -OutFile tsplay-quickstart.ps1
powershell -ExecutionPolicy Bypass -File .\tsplay-quickstart.ps1
```

The installer script detects the platform, downloads the matching binary, runs `quickstart-demo`, and writes `artifacts/quickstart/quickstart-demo-output.json` without requiring Playwright first.

If you prefer the manual per-platform downloads instead:

| Platform | Download |
| --- | --- |
| macOS Apple Silicon | [tsplay-darwin-arm64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-arm64) |
| macOS Intel | [tsplay-darwin-amd64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-darwin-amd64) |
| Linux x86_64 | [tsplay-linux-amd64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-amd64) |
| Linux ARM64 | [tsplay-linux-arm64](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-linux-arm64) |
| Windows x86_64 | [tsplay-windows-amd64.exe](https://github.com/tensafe/tsplay/releases/latest/download/tsplay-windows-amd64.exe) |

If you prefer the source-first path instead:

```bash
go mod download
go run . -flow script/tutorials/01_hello_world.flow.yaml
```

Only start the local demo server when you want to practice against bundled demo pages:

```bash
go run . -action file-srv -addr :8000
```

## Next Stops

- Project overview: [README.zh-CN.md](README.zh-CN.md), [ReadMe.md](ReadMe.md)
- Tutorial index: [docs/tutorials/README.zh-CN.md](docs/tutorials/README.zh-CN.md), [docs/tutorials/README.md](docs/tutorials/README.md)
- Docs map: [docs/README.md](docs/README.md)
- Doc health audit: [docs/doc-health-audit.md](docs/doc-health-audit.md)
- Core feature execution board: [docs/product/core-feature-execution-board.md](docs/product/core-feature-execution-board.md)

## Default Next Step

After the smallest first run, it is usually easier to pick one next path than to open the whole document map at once.
Choose the direction that feels closest to what you want to do next:

- Learn basic lessons: [docs/tutorials/README.zh-CN.md](docs/tutorials/README.zh-CN.md)
- Learn Flow as the default delivery path: [docs/tutorials/track-newbie.zh-CN.md](docs/tutorials/track-newbie.zh-CN.md)
- Go straight to Agent / MCP: [docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md)
- Learn single-binary delivery: [docs/tutorials/144-single-binary-delivery-flow.md](docs/tutorials/144-single-binary-delivery-flow.md)

## Common Next-Step Questions

- Downloaded the binary and want the one-step first run:
  Use `./tsplay -action quickstart-demo`, then continue with [site-src/getting-started.md](site-src/getting-started.md) for the browser path.
- Built the binary but not sure whether to use `list-assets`, `extract-assets`, or `file-srv` first:
  Start with [docs/tutorials/142-list-assets-for-beginners.md](docs/tutorials/142-list-assets-for-beginners.md), then [docs/tutorials/143-extract-assets-for-beginners.md](docs/tutorials/143-extract-assets-for-beginners.md).
- Want Agent flow generation before understanding the default MCP path:
  Start with [docs/training/ai-intent-to-flow.md](docs/training/ai-intent-to-flow.md), then [docs/tutorials/120-mcp-finalize-flow.md](docs/tutorials/120-mcp-finalize-flow.md).
- See files under `artifacts/` but do not know which ones matter for handoff:
  Start with [docs/tutorials/87-build-handoff-artifact-manifest.md](docs/tutorials/87-build-handoff-artifact-manifest.md).
