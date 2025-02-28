# **操作手册**
本操作手册详细介绍TSPlay的导航、行为操作、等待操作、截图操作等指令功能。

---

## **1. 导航类 / Navigation**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `navigate` | 导航到指定的 URL | `navigate('https://example.com')` | `url` (string): 目标 URL。 |
| `click` | 点击页面上的元素 | `click('#button-id')` | `selector` (string): 要点击的元素选择器。 |
| `reload` | 重新加载当前页面 | `reload()` | 无参数 |
| `go_back` | 返回到上一个页面 | `go_back()` | 无参数 |
| `go_forward` | 前进到下一个页面 | `go_forward()` | 无参数 |


---

## **2. 行为类 / Actions**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `type_text` | 在指定元素中输入文本 | `type_text('#input-id', 'Hello World')` | `selector` (string): 输入框的选择器；`text` (string): 要输入的文本。 |
| `get_text` | 获取指定元素的文本内容 | `get_text('#element-id')` | `selector` (string): 要获取文本内容的元素选择器。 |
| `set_value` | 设置指定元素的值 | `set_value('#input-id', 'new value')` | `selector` (string): 输入框的选择器；`value` (string): 要设置的值。 |
| `select_option` | 选择下拉框中的选项 | `select_option('#dropdown-id', 'option-value')` | `selector` (string): 下拉框选择器；`value` (string): 要选择的选项值。 |
| `hover` | 将鼠标悬停在指定元素上 | `hover('#element-id')` | `selector` (string): 要悬停的元素选择器。 |
| `scroll_to` | 滚动页面到指定位置 | `scroll_to('#element-id')` | `selector` (string): 要滚动到的元素选择器。 |


---

## **3. 等待操作 / Waiting**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `wait_for_network_idle` | 等待网络空闲 | `wait_for_network_idle()` | 无参数 |
| `wait_for_selector` | 等待指定选择器匹配的元素出现 | `wait_for_selector('#element-id', 5000)` | `selector` (string): 要等待的选择器；`timeout` (int, optional): 超时时间（默认 30000 毫秒）。 |
| `wait_for_text` | 等待指定文本出现在页面中 | `wait_for_text('#element-id', 'Hello World', 5000)` | `selector` (string): 元素选择器；`text` (string): 期待的文本；`timeout` (int, optional): 超时时间（默认 30000 毫秒）。 |
| `sleep` | 暂停执行指定的时间 | `sleep(2)` | `seconds` (number): 暂停时间（秒）。 |


---

## **4. 页面截图 / Screenshots**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `screenshot` | 截取整个页面的截图 | `screenshot('screenshot.png')` | `path` (string): 保存截图的文件路径。 |
| `screenshot_element` | 截取指定元素的截图 | `screenshot_element('#element-id', 'element.png')` | `selector` (string): 元素选择器；`path` (string): 保存截图的文件路径。 |
| `save_html` | 保存当前页面的 HTML 内容 | `save_html('page.html')` | `path` (string): 保存 HTML 的文件路径。 |


---

## **5. 处理弹窗和对话框 / Handling Dialogs**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `accept_alert` | 接受弹窗（点击确定） | `accept_alert()` | 无参数 |
| `dismiss_alert` | 关闭弹窗（点击取消） | `dismiss_alert()` | 无参数 |
| `set_alert_text` | 在弹窗中输入文本 | `set_alert_text('Hello')` | `text` (string): 要输入的文本。 |


---

## **6. 执行 JavaScript / JavaScript Execution**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `execute_script` | 在页面中执行 JavaScript 代码 | `execute_script('alert("Hello World")')` | `script` (string): 要执行的 JavaScript 代码。 |
| `evaluate` | 执行 JavaScript 表达式并返回结果 | `evaluate('#element-id', 'element => element.textContent')` | `selector` (string): 元素选择器；`script` (string): JavaScript 表达式。 |


---

## **7. 上传文件 / File Upload/Download**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `upload_file` | 上传单个文件到指定元素 | `upload_file('#file-input', 'file.txt')` | `selector` (string): 文件输入框选择器；`file_path` (string): 要上传的文件路径。 |
| `upload_multiple_files` | 上传多个文件到指定元素 | `upload_multiple_files('#file-input', 'file1.txt', 'file2.txt')` | `selector` (string): 文件输入框选择器；`files` (string[]): 要上传的文件路径列表。 |
| `download_file` | 下载文件到本地 | `download_file('https://example.com/file.txt', 'file.txt')` | `url` (string): 文件 URL；`save_path` (string): 保存文件的路径。 |


---

## **8. 提取数据 / Data Extraction**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `get_attribute` | 获取指定元素的属性值 | `get_attribute('#element-id', 'href')` | `selector` (string): 元素选择器；`attribute` (string): 属性名称。 |
| `get_html` | 获取指定元素的 HTML 内容 | `get_html('#element-id')` | `selector` (string, optional): 元素选择器（如果省略，返回页面的完整 HTML）。 |
| `get_all_links` | 获取页面中所有链接 | `get_all_links()` | 无参数 |
| `capture_table` | 提取表格数据 | `capture_table('#table-id')` | `selector` (string): 表格元素的选择器。 |


---

## **9. 页面状态检查 / Page State Checks**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `is_visible` | 检查元素是否可见 | `is_visible('#element-id')` | `selector` (string): 元素选择器。 |
| `is_enabled` | 检查元素是否可用 | `is_enabled('#element-id')` | `selector` (string): 元素选择器。 |
| `is_checked` | 检查复选框或单选按钮是否被选中 | `is_checked('#checkbox-id')` | `selector` (string): 元素选择器。 |
| `is_selected` | 检查下拉框选项是否被选中 | `is_selected('#dropdown-id')` | `selector` (string): 下拉框选择器。 |
| `is_aria_selected` | 检查 ARIA 属性是否被选中 | `is_aria_selected('#element-id')` | `selector` (string): 元素选择器。 |


---

## **10. 多标签页和窗口管理 / Tab and Window Management**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `new_tab` | 打开一个新标签页 | `new_tab('https://example.com')` | `url` (string): 要在新标签页中打开的 URL。 |
| `close_tab` | 关闭当前标签页 | `close_tab()` | 无参数 |
| `switch_to_tab` | 切换到指定的标签页 | `switch_to_tab(2)` | `index` (int): 要切换到的标签页索引。 |


---

## **11. 网络请求与拦截 / Network Request Handling**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `intercept_request` | 拦截网络请求 | `intercept_request(function(request) return 'https://example.com' end)` | `callback` (function): 用于处理请求的 Lua 函数。 |
| `block_request` | 阻止指定的网络请求 | `block_request('*.png')` | `pattern` (string): 要阻止的请求模式。 |
| `get_response` | 获取网络请求的响应 | `get_response('https://example.com/api')` | `url` (string): 请求的 URL。 |


---

## **12. StateStorage 管理 / State Storage Management**
| **函数名** | **说明** | **使用示例** | **参数** |
| --- | --- | --- | --- |
| `get_storage_state` | 获取当前页面的存储状态 | `get_storage_state()` | 无参数 |
| `get_cookies_string` | 获取当前页面的 Cookie 字符串 | `get_cookies_string()` | 无参数 |


---

以上就是所有操作的详细说明和示例，便于快速上手并有效操作。

