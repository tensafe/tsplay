# 能力动作类别：文件与表格 I/O

这组动作负责把页面现场、JSON/CSV/Excel 数据和本地产物连起来。  
在 MCP 场景下，它们最常受 `allow_file_access` 影响。

| 动作 | Flow | Lua | MCP | 典型写法 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `screenshot` | 是 | 是 | 是 | `action: screenshot` + `path` / `screenshot(path)` | 保存整页截图。常用于失败证据和交付附件。 |
| `screenshot_element` | 是 | 是 | 是 | `action: screenshot_element` + `selector,path` / `screenshot_element(selector, path)` | 保存局部元素截图。 |
| `save_html` | 是 | 是 | 是 | `action: save_html` + `path` / `save_html(path)` | 保存当前页面 HTML。适合排障、交接、离线审查。 |
| `read_json` | 是 | 是 | 是 | `action: read_json` + `file_path` / `read_json(path)` | 读取本地 JSON，返回对象、列表或原始值。 |
| `read_csv` | 是 | 是 | 是 | `action: read_csv` + `file_path` / `read_csv(path, start_row, limit, row_field)` | 读取 CSV 为行对象列表。支持续跑和保留源行号。 |
| `read_excel` | 是 | 是 | 是 | `action: read_excel` + `file_path` / `read_excel(path, sheet, range, headers, start_row, limit, row_field)` | 读取 Excel 为行对象列表。 |
| `write_json` | 是 | 是 | 是 | `action: write_json` + `file_path,value` / `write_json(path, value)` | 把任意值写成 JSON。 |
| `write_csv` | 是 | 是 | 是 | `action: write_csv` + `file_path,value` / `write_csv(path, rows, headers)` | 把行对象写成 CSV。 |
| `write_excel` | 是 | 是 | 是 | `action: write_excel` + `file_path,value` / `write_excel(path, rows, headers, sheet)` | 把单表或多表数据写成 `.xlsx`。 |

## 最小示例小代码

### Flow

```yaml
schema_version: "1"
name: file_io_demo
steps:
  - action: read_csv
    file_path: artifacts/input/orders.csv
    save_as: rows

  - action: write_json
    file_path: artifacts/output/orders.json
    with:
      value: "{{rows}}"
```

### Lua

```lua
local rows = read_csv("artifacts/input/orders.csv")
write_json("artifacts/output/orders.json", rows)
screenshot("artifacts/output/orders-page.png")
```

## 使用建议

- 教程、排障、交付三类场景里，优先把关键产物写到稳定的 `artifacts/` 路径
- `read_csv / read_excel` 适合配合 `foreach` 做分批处理
- `write_json / write_csv / write_excel` 更适合留结构化交付物，不只打印终端输出

## 相关教程

- [Lesson 12](../tutorials/12-custom-json-output.md)
- [Lesson 13](../tutorials/13-read-csv-basics.md)
- [Lesson 24](../tutorials/24-read-excel-basics.md)
- [Lesson 31](../tutorials/31-full-page-screenshot.md)
- [Lesson 32](../tutorials/32-element-screenshot.md)
- [Lesson 33](../tutorials/33-save-html-basics.md)
