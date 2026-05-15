# 能力动作类别：浏览器状态

这组能力负责保存、读取和复用浏览器状态。  
在 Flow / MCP 中，重点授权是 `allow_browser_state`。

| 动作 / 配置 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `get_storage_state` | 是 | 是 | 是 | `action: get_storage_state` / `get_storage_state()` | 读取当前浏览器上下文的 storage state。 |
| `get_cookies_string` | 是 | 是 | 是 | `action: get_cookies_string` / `get_cookies_string()` | 把 cookies 导成字符串，适合喂给接口或日志。 |
| `browser.use_session` | 是 | 否 | 是 | `browser.use_session: demo_admin` | Flow 顶层浏览器配置，不是普通 step。推荐作为复用命名会话的默认写法。 |
| `browser.cdp_launch` | 是 | 通过 CLI 参数 | 是 | `browser.cdp_launch: true` | TSPlay 自动查找本机 Chrome/Chromium/Edge，启动独立 profile 和远程调试端口，再通过 CDP 接管。MCP 下需要 `allow_browser_state=true`。 |
| `browser.cdp_endpoint` / `browser.cdp_port` | 是 | 通过 CLI 参数 | 是 | `browser.cdp_port: 9222` | 通过 CDP 接管真实 Chrome/Chromium，复用用户数据、登录态和扩展。MCP 下需要 `allow_browser_state=true`。 |
| `save_storage_state` | 否 | 是 | 否 | `save_storage_state('states/admin.json')` | Lua 辅助能力，把当前 state 保存到本地文件。 |
| `load_storage_state` | 否 | 是 | 否 | `load_storage_state('states/admin.json')` | Lua 辅助能力，从文件加载 state 到新上下文。 |
| `use_session` | 否 | 是 | 否 | `use_session('demo_admin')` | Lua 辅助能力，复用已保存命名会话。Flow 等价物是 `browser.use_session`。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: browser_state_demo
browser:
  use_session: session_lab_demo
steps:
  - action: navigate
    url: http://127.0.0.1:8000/demo/session_lab.html

  - action: get_cookies_string
    save_as: cookie_header
```

### 接管外部 Chrome

先用远程调试端口启动浏览器：

```bash
google-chrome --remote-debugging-port=9222
```

然后在 Flow 顶层接入：

```yaml
schema_version: "1"
name: cdp_attach_demo
browser:
  cdp_port: 9222
steps:
  - action: navigate
    url: https://example.com
```

`cdp_endpoint` 也可以直接传 `ws://127.0.0.1:9222/devtools/browser/...`、`http://127.0.0.1:9222`，或直接粘贴 `127.0.0.1:9222/json/version` / `127.0.0.1:9222/devtools/browser/...` 这类本机调试地址。CDP 接管会复用外部浏览器的默认 context 和第一个页面；运行结束后 TSPlay 只断开 Playwright 连接，不会关闭真实浏览器。

### 自动启动本机浏览器

如果不想手动找 Chrome 路径或自己加 `--remote-debugging-port`，可以让 TSPlay 启动一个独立 profile：

```yaml
schema_version: "1"
name: cdp_launch_demo
browser:
  cdp_launch: true
  cdp_port: 9222
  cdp_user_data_dir: profiles/cdp-demo
steps:
  - action: navigate
    url: https://example.com
```

`cdp_port` 和 `cdp_user_data_dir` 都可以省略；省略端口时 TSPlay 会挑一个空闲本地端口，省略 profile 时会放到 artifact root 下。浏览器可执行文件不写时会自动搜索 macOS、Windows、Linux 常见位置；找不到时再用 `cdp_executable` 或 CLI 的 `-browser-cdp-executable` 手动指定。

### Lua

```lua
use_session("session_lab_demo")
navigate("http://127.0.0.1:8000/demo/session_lab.html")
local cookie_header = get_cookies_string()
print(cookie_header)
```

## 使用建议

- 长期复用时，优先 `browser.use_session`，而不是把登录步骤散落在每条 Flow 里
- 只想临时拿 cookie 或 state 时，`get_*` 动作就够了
- 要导出成可复用片段时，再配合 `save-session / export-session` 那条 CLI / MCP 入口

## 相关教程

- [Lesson 36](../tutorials/36-save-storage-state.md)
- [Lesson 40](../tutorials/40-save-named-session.md)
- [Lesson 42](../tutorials/42-use-named-session.md)
- [Lesson 57](../tutorials/57-use-session-import-export-round-trip.md)
