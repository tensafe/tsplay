# 能力动作类别：补充浏览器动作

这页收的是不一定会出现在“第一页矩阵”里，但在真实自动化里经常要用到的浏览器动作。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `get_text` | 是 | 是 | 是 | `action: get_text` + `selector` / `get_text(selector)` | 直接读取元素文本。 |
| `get_attribute` | 是 | 是 | 是 | `action: get_attribute` + `selector,attribute` / `get_attribute(selector, attr)` | 读取属性值。 |
| `get_html` | 是 | 是 | 是 | `action: get_html` + `selector?` / `get_html(selector)` | 获取 HTML 片段或整页 HTML。 |
| `get_all_links` | 是 | 是 | 是 | `action: get_all_links` + `selector?` / `get_all_links(selector)` | 收集链接列表。 |
| `capture_table` | 是 | 是 | 是 | `action: capture_table` + `selector` / `capture_table(selector)` | 读取表格为结构化行数据。 |
| `find_element` | 是 | 是 | 是 | `action: find_element` + `selector` / `find_element(selector)` | 找单个元素，偏原语级。 |
| `find_elements` | 是 | 是 | 是 | `action: find_elements` + `selector` / `find_elements(selector)` | 找多个元素。 |
| `is_visible` | 是 | 是 | 是 | `action: is_visible` + `selector` / `is_visible(selector)` | 判断元素是否可见。 |
| `is_enabled` | 是 | 是 | 是 | `action: is_enabled` + `selector` / `is_enabled(selector)` | 判断元素是否可用。 |
| `is_checked` | 是 | 是 | 是 | `action: is_checked` + `selector` / `is_checked(selector)` | 判断 checkbox / radio 是否勾选。 |
| `is_selected` | 是 | 是 | 是 | `action: is_selected` + `selector` / `is_selected(selector)` | 判断选项是否被选中。 |
| `is_aria_selected` | 是 | 是 | 是 | `action: is_aria_selected` + `selector` / `is_aria_selected(selector)` | 判断 ARIA 选中态。 |
| `upload_file` | 是 | 是 | 是 | `action: upload_file` + `selector,file_path` / `upload_file(selector, path)` | 上传单文件。通常同时涉及浏览器和本地文件。 |
| `upload_multiple_files` | 是 | 是 | 是 | `action: upload_multiple_files` + `selector,files` / `upload_multiple_files(selector, ...)` | 上传多文件。 |
| `download_file` | 是 | 是 | 是 | `action: download_file` + `selector,save_path` / `download_file(selector, path)` | 点击触发下载并保存到本地。 |
| `download_url` | 是 | 是 | 是 | `action: download_url` + `url,save_path` / `download_url(url, path)` | 直接下载指定 URL。 |
| `accept_alert` | 是 | 是 | 是 | `action: accept_alert` / `accept_alert()` | 接受弹窗。 |
| `dismiss_alert` | 是 | 是 | 是 | `action: dismiss_alert` / `dismiss_alert()` | 关闭弹窗。 |
| `set_alert_text` | 是 | 是 | 是 | `action: set_alert_text` + `text` / `set_alert_text(text)` | 给 prompt 弹窗写入文本。 |
| `execute_script` | 是 | 是 | 是 | `action: execute_script` + `script` / `execute_script(js)` | 在页面上下文执行脚本。 |
| `evaluate` | 是 | 是 | 是 | `action: evaluate` + `selector,script` / `evaluate(selector, js)` | 在选中元素上执行表达式并拿结果。 |
| `new_tab` | 是 | 是 | 是 | `action: new_tab` + `url` / `new_tab(url)` | 打开新标签页。 |
| `close_tab` | 是 | 是 | 是 | `action: close_tab` / `close_tab()` | 关闭当前标签。 |
| `switch_to_tab` | 是 | 是 | 是 | `action: switch_to_tab` + `index` / `switch_to_tab(index)` | 切换标签页。 |
| `block_request` | 是 | 是 | 是 | `action: block_request` + `pattern` / `block_request(pattern)` | 按模式阻止请求。 |
| `get_response` | 是 | 是 | 是 | `action: get_response` + `url` / `get_response(url)` | 取某个请求的响应。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: supplemental_browser_demo
steps:
  - action: capture_table
    selector: "table"
    save_as: table_rows

  - action: write_csv
    file_path: artifacts/output/table_rows.csv
    with:
      value: "{{table_rows}}"
```

### Lua

```lua
local links = get_all_links("body")
local first_enabled = is_enabled("#submit")
local title_text = get_text("title")
print(links, first_enabled, title_text)
```

## 使用建议

- 上传、下载、截图这类动作，经常同时碰到浏览器和本地文件边界
- `execute_script / evaluate` 能解决问题，但要先判断是不是已有结构化动作更合适
- `find_* / is_*` 更偏探索和诊断，不一定每次都要进最终交付 Flow

## 相关教程

- [Lesson 03](../tutorials/03-capture-table.md)
- [Lesson 18](../tutorials/18-upload-single-file.md)
- [Lesson 19](../tutorials/19-upload-multiple-files.md)
- [Lesson 20](../tutorials/20-download-report.md)
- [Lesson 29](../tutorials/29-read-cookies-string.md)
