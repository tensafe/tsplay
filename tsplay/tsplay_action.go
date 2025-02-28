package tsplay

import (
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	lua "github.com/yuin/gopher-lua"
	"os"
	"strings"
	"time"
)

type LuaFunction struct {
	Name           string                  // Function name
	Func           func(L *lua.LState) int // Function implementation
	Description_cn string                  // Function description (Chinese)
	Description_en string                  // Function description (English)
}

var GlobalPlayWrightFunc = []LuaFunction{
	// 导航类 / Navigation
	{"navigate", navigate, "导航到指定的 URL", "Navigate to a specified URL. Example: navigate('https://example.com'). Parameters: url (string) - The URL to navigate to."},
	{"click", click, "点击页面上的元素", "Click on an element on the page. Example: click('#button-id'). Parameters: selector (string) - The selector of the element to click."},
	{"reload", reload, "重新加载当前页面", "Reload the current page. Example: reload(). No parameters."},
	{"go_back", go_back, "返回到上一个页面", "Go back to the previous page. Example: go_back(). No parameters."},
	{"go_forward", go_forward, "前进到下一个页面", "Go forward to the next page. Example: go_forward(). No parameters."},

	// 行为类 / Actions
	{"type_text", type_text, "在指定元素中输入文本", "Type text into a specified element. Example: type_text('#input-id', 'Hello World'). Parameters: selector (string) - The selector of the input element; text (string) - The text to type."},
	{"get_text", get_text, "获取指定元素的文本内容", "Get the text content of a specified element. Example: get_text('#element-id'). Parameters: selector (string) - The selector of the element to retrieve text from."},
	{"set_value", set_value, "设置指定元素的值", "Set the value of a specified element. Example: set_value('#input-id', 'new value'). Parameters: selector (string) - The selector of the input element; value (string) - The value to set."},
	{"select_option", select_option, "选择下拉框中的选项", "Select an option in a dropdown. Example: select_option('#dropdown-id', 'option-value'). Parameters: selector (string) - The selector of the dropdown; value (string) - The value of the option to select."},
	{"hover", hover, "将鼠标悬停在指定元素上", "Hover the mouse over a specified element. Example: hover('#element-id'). Parameters: selector (string) - The selector of the element to hover over."},
	{"scroll_to", scroll_to, "滚动页面到指定位置", "Scroll the page to a specified position. Example: scroll_to('#element-id'). Parameters: selector (string) - The selector of the element to scroll to."},

	// 等待操作 / Waiting
	{"wait_for_network_idle", wait_for_network_idle, "等待网络空闲", "Wait for the network to be idle. Example: wait_for_network_idle(). No parameters."},
	{"wait_for_selector", wait_for_selector, "等待指定选择器匹配的元素出现", "Wait for an element matching the specified selector to appear. Example: wait_for_selector('#element-id', 5000). Parameters: selector (string) - The selector to wait for; timeout (int, optional) - Timeout in milliseconds (default is 30000)."},
	{"wait_for_text", wait_for_text, "等待指定文本出现在页面中", "Wait for specified text to appear on the page. Example: wait_for_text('#element-id', 'Hello World', 5000). Parameters: selector (string) - The selector of the element; text (string) - The expected text; timeout (int, optional) - Timeout in milliseconds (default is 30000)."},
	{"sleep", sleep, "暂停执行指定的时间", "Pause execution for a specified duration. Example: sleep(2). Parameters: seconds (number) - The duration to sleep in seconds."},

	// 页面截图 / Screenshots
	{"screenshot", screenshot, "截取整个页面的截图", "Take a screenshot of the entire page. Example: screenshot('screenshot.png'). Parameters: path (string) - The file path to save the screenshot."},
	{"screenshot_element", screenshot_element, "截取指定元素的截图", "Take a screenshot of a specific element. Example: screenshot_element('#element-id', 'element.png'). Parameters: selector (string) - The selector of the element; path (string) - The file path to save the screenshot."},
	{"save_html", save_html, "保存当前页面的 HTML 内容", "Save the HTML content of the page. Example: save_html('page.html'). Parameters: path (string) - The file path to save the HTML content."},

	// 处理弹窗和对话框 / Handling Dialogs
	{"accept_alert", accept_alert, "接受弹窗（点击确定）", "Accept an alert dialog. Example: accept_alert(). No parameters."},
	{"dismiss_alert", dismiss_alert, "关闭弹窗（点击取消）", "Dismiss an alert dialog. Example: dismiss_alert(). No parameters."},
	{"set_alert_text", set_alert_text, "在弹窗中输入文本", "Set the text in an alert dialog. Example: set_alert_text('Hello'). Parameters: text (string) - The text to set in the alert dialog."},

	// 执行 JavaScript / JavaScript Execution
	{"execute_script", execute_script, "在页面中执行 JavaScript 代码", "Execute JavaScript code on the page. Example: execute_script('alert(\"Hello World\")'). Parameters: script (string) - The JavaScript code to execute."},
	{"evaluate", evaluate, "在页面中执行 JavaScript 表达式并返回结果", "Evaluate JavaScript expression and return the result. Example: evaluate('#element-id', 'element => element.textContent'). Parameters: selector (string) - The selector of the element; script (string) - The JavaScript expression to evaluate."},

	// 上传文件 / File Upload/Download
	{"upload_file", upload_file, "上传单个文件到指定元素", "Upload a single file. Example: upload_file('#file-input', 'file.txt'). Parameters: selector (string) - The selector of the file input; file_path (string) - The path to the file to upload."},
	{"upload_multiple_files", upload_multiple_files, "上传多个文件到指定元素", "Upload multiple files. Example: upload_multiple_files('#file-input', 'file1.txt', 'file2.txt'). Parameters: selector (string) - The selector of the file input; files (string[]) - A list of file paths to upload."},
	{"download_file", download_file, "下载文件到本地", "Download a file from the page. Example: download_file('https://example.com/file.txt', 'file.txt'). Parameters: url (string) - The file URL; save_path (string) - The path to save the downloaded file."},

	// 提取数据 / Data Extraction
	{"get_attribute", get_attribute, "获取指定元素的属性值", "Get the value of a specified attribute of an element. Example: get_attribute('#element-id', 'href'). Parameters: selector (string) - The selector of the element; attribute (string) - The attribute name."},
	{"get_html", get_html, "获取指定元素的 HTML 内容", "Get the HTML content of an element. Example: get_html('#element-id'). Parameters: selector (string, optional) - The selector of the element (if omitted, returns the entire page's HTML)."},
	{"get_all_links", get_all_links, "获取页面中所有链接", "Extract all links from the page. Example: get_all_links(). No parameters."},
	{"capture_table", capture_table, "提取表格数据", "Capture and extract data from a table element. Example: capture_table('#table-id'). Parameters: selector (string) - The selector of the table element."},

	// 页面状态检查 / Page State Checks
	{"is_visible", is_visible, "检查元素是否可见", "Check if an element is visible. Example: is_visible('#element-id'). Parameters: selector (string) - The selector of the element."},
	{"is_enabled", is_enabled, "检查元素是否可用", "Check if an element is enabled. Example: is_enabled('#element-id'). Parameters: selector (string) - The selector of the element."},
	{"is_checked", is_checked, "检查复选框或单选按钮是否被选中", "Check if a checkbox or radio button is checked. Example: is_checked('#checkbox-id'). Parameters: selector (string) - The selector of the element."},
	{"is_selected", is_selected, "检查下拉框选项是否被选中", "Check if an option in a dropdown is selected. Example: is_selected('#dropdown-id'). Parameters: selector (string) - The selector of the dropdown."},
	{"is_aria_selected", is_aria_selected, "检查 ARIA 属性是否被选中", "Check if an element has the ARIA 'selected' attribute. Example: is_aria_selected('#element-id'). Parameters: selector (string) - The selector of the element."},

	// 多标签页和窗口管理 / Tab and Window Management
	{"new_tab", new_tab, "打开一个新标签页", "Open a new browser tab. Example: new_tab('https://example.com'). Parameters: url (string) - The URL to open in the new tab."},
	{"close_tab", close_tab, "关闭当前标签页", "Close the current browser tab. Example: close_tab(). No parameters."},
	{"switch_to_tab", switch_to_tab, "切换到指定的标签页", "Switch to a specific browser tab. Example: switch_to_tab(2). Parameters: index (int) - The index of the tab to switch to."},

	// 网络请求与拦截 / Network Request Handling
	{"intercept_request", intercept_request, "拦截网络请求", "Intercept and modify network requests. Example: intercept_request(function(request) return 'https://example.com' end). Parameters: callback (function) - A Lua function to handle intercepted requests."},
	{"block_request", block_request, "阻止指定的网络请求", "Block specific network requests. Example: block_request('*.png'). Parameters: pattern (string) - The pattern of requests to block."},
	{"get_response", get_response, "获取网络请求的响应", "Get the response of a network request. Example: get_response('https://example.com/api'). Parameters: url (string) - The URL of the request to get the response for."},

	// StateStorage 管理 / State Storage Management
	{"get_storage_state", get_storage_state, "获取当前页面的存储状态", "Get the current browser storage state. Example: get_storage_state(). No parameters."},
	{"get_cookies_string", get_cookies_string, "获取当前页面的 Cookie 字符串", "Get cookies as a string. Example: get_cookies_string(). No parameters."},
}

// 安全获得page
func safe_page(L *lua.LState) playwright.Page {
	pageUserData := L.GetGlobal("page")
	if pageUserData == lua.LNil {
		L.RaiseError("No 'page' object found in Lua context")
		return nil
	}

	page, ok := pageUserData.(*lua.LUserData)
	if !ok {
		L.RaiseError("'page' is not of the expected type")
		return nil
	}

	playwrightPage, ok := page.Value.(playwright.Page)
	if !ok {
		L.RaiseError("'page' does not contain a valid Playwright Page object")
		return nil
	}
	return playwrightPage
}

func safe_browser(L *lua.LState) playwright.Browser {
	browserUserData := L.GetGlobal("browser")
	if browserUserData == lua.LNil {
		L.RaiseError("No 'page' object found in Lua context")
		return nil
	}

	browser, ok := browserUserData.(*lua.LUserData)
	if !ok {
		L.RaiseError("'page' is not of the expected type")
		return nil
	}

	playwrightBrowser, ok := browser.Value.(playwright.Browser)
	if !ok {
		L.RaiseError("'page' does not contain a valid Playwright Page object")
		return nil
	}
	return playwrightBrowser
}

//(1)页面导航
//操作页面加载和导航的行为。
//navigate(url)，打开指定的 URL 页面。
//reload()，刷新当前页面。
//go_back()，返回到上一个页面。
//go_forward()，前进到下一个页面。

// Navigate to a URL
func navigate(L *lua.LState) int {
	url := L.ToString(1)

	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 调用 Playwright 页面导航功能
	_, err := page.Goto(url)
	if err != nil {
		fmt.Println("Failed to navigate to %s: %v", url, err)
		return 0
	}

	fmt.Println("Successfully navigated to:", url)
	return 0
}

func reload(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}
	response, err := page.Reload()
	if err != nil {
		L.RaiseError("Failed to reload")
		return 0
	}

	fmt.Println("Successfully reload", response)
	return 0
}

func go_back(L *lua.LState) int {
	// 获取安全页面对象
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 调用 Page.GoBack() 方法
	response, err := page.GoBack()
	if err != nil {
		fmt.Println("Failed to go back: %v", err)
		return 0
	}

	// 检查是否成功返回上一页
	if response != nil {
		fmt.Println("Successfully went back, Response Status:", response.Status())
	} else {
		fmt.Println("No history to go back to.")
	}

	return 0
}

func go_forward(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	response, err := page.GoForward()
	if err != nil {
		fmt.Println("Failed to go forward: %v", err)
		return 0
	}

	if response != nil {
		fmt.Println("Successfully went forward, Response Status:", response.Status())
	} else {
		fmt.Println("No history to go forward to.")
	}

	return 0
}

// (2)元素查找与交互
// 操作页面元素的行为。
// click(selector)，点击指定的元素。
// type_text(selector, text)，在指定输入框中输入文本。
// get_text(selector)，获取指定元素的文本内容。
// set_value(selector, value)，设置输入框的值。
// select_option(selector, value)，从下拉菜单中选择指定选项。
// hover(selector)，将鼠标悬停在指定元素上。
// scroll_to(selector)，滚动到页面上的指定元素。
func click(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 参数
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 调用 Playwright 的 Click 方法
	if err := page.Click(selector); err != nil {
		L.RaiseError("Failed to click on selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully clicked on selector: %s\n", selector)
	return 0
}
func type_text(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 和 text 参数
	selector := L.CheckString(1)
	text := L.CheckString(2)

	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 调用 Playwright 的 Fill 方法，在输入框中输入文本
	if err := page.Fill(selector, text); err != nil {
		L.RaiseError("Failed to type text into selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully typed text '%s' into selector: %s\n", text, selector)
	return 0
}

func get_text(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 参数
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 获取所有匹配的元素
	elements, err := page.QuerySelectorAll(selector)
	if err != nil {
		L.RaiseError("Failed to query elements with selector '%s': %v", selector, err)
		return 0
	}

	// 如果没有匹配的元素，返回 nil
	if len(elements) == 0 {
		L.Push(lua.LNil)
		return 1
	}

	// 如果只有一个元素，返回其文本内容
	if len(elements) == 1 {
		text, err := elements[0].TextContent()
		if err != nil {
			L.RaiseError("Failed to get text from element: %v", err)
			return 0
		}
		L.Push(lua.LString(text))
		return 1
	}

	// 如果有多个元素，返回文本内容列表
	textsTable := L.NewTable()
	for i, element := range elements {
		text, err := element.TextContent()
		if err != nil {
			L.RaiseError("Failed to get text from element %d: %v", i+1, err)
			return 0
		}
		L.RawSetInt(textsTable, i+1, lua.LString(text))
	}
	L.Push(textsTable)
	return 1
}

func set_value(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 和 value 参数
	selector := L.CheckString(1)
	value := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 设置输入框的值
	if err := page.Fill(selector, value); err != nil {
		L.RaiseError("Failed to set value for selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully set value '%s' for selector: %s\n", value, selector)
	return 0
}
func select_option(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 和 value 参数
	selector := L.CheckString(1)
	value := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 从下拉菜单中选择指定选项
	_, err := page.SelectOption(selector, playwright.SelectOptionValues{Values: &[]string{value}})
	if err != nil {
		L.RaiseError("Failed to select option '%s' for selector '%s': %v", value, selector, err)
		return 0
	}

	fmt.Printf("Successfully selected option '%s' for selector: %s\n", value, selector)
	return 0
}
func hover(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 参数
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 将鼠标悬停在指定元素上
	if err := page.Hover(selector); err != nil {
		L.RaiseError("Failed to hover on selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully hovered on selector: %s\n", selector)
	return 0
}
func scroll_to(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 参数
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 滚动到指定元素
	if _, err := page.EvalOnSelector(selector, "element => element.scrollIntoView()", nil); err != nil {
		L.RaiseError("Failed to scroll to selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully scrolled to selector: %s\n", selector)
	return 0
}

// (3)等待操作
// 等待页面加载完成或某些条件满足。
// wait_for_network_idle 等待网页加载完成
// wait_for_selector(selector, timeout)，等待指定的元素出现在页面上。
// wait_for_text(selector, text, timeout)，等待指定元素的文本变为某个值。
// wait_for_navigation(timeout)，等待页面完成导航操作。
// sleep(seconds)，暂停脚本执行指定的秒数。

func wait_for_network_idle(L *lua.LState) int {
	// 使用 safe_page 方法获取页面对象
	page := safe_page(L)
	if page == nil {
		L.RaiseError("Failed to get page object")
		return 0
	}

	// 等待页面达到网络空闲状态
	err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	if err != nil {
		L.RaiseError("Failed to wait for network idle: %v", err)
		return 0
	}

	// 返回成功状态
	L.Push(lua.LBool(true))
	return 1
}

func wait_for_selector(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 和 timeout 参数
	selector := L.CheckString(1)
	timeout := L.OptInt(2, 30000) // 默认超时时间为 30 秒（30000 毫秒）

	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 等待元素出现在页面上
	_, err := page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(float64(timeout)),
	})
	if err != nil {
		L.RaiseError("Failed to wait for selector '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Successfully waited for selector: %s\n", selector)
	return 0
}

func wait_for_text(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector、text 和 timeout 参数
	selector := L.CheckString(1)
	expectedText := L.CheckString(2)
	timeout := L.OptInt(3, 30000) // 默认超时时间为 30 秒（30000 毫秒）

	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 等待元素的文本变为指定值
	err, _ := page.WaitForFunction(
		fmt.Sprintf(`() => document.querySelector('%s')?.textContent.trim() === '%s'`, selector, expectedText),
		playwright.PageWaitForFunctionOptions{Timeout: playwright.Float(float64(timeout))},
	)
	if err != nil {
		L.RaiseError("Failed to wait for text '%s' in selector '%s': %v", expectedText, selector, err)
		return 0
	}

	fmt.Printf("Successfully waited for text '%s' in selector: %s\n", expectedText, selector)
	return 0
}

func sleep(L *lua.LState) int {
	// 从 Lua 获取 seconds 参数
	seconds := L.CheckNumber(1) // 返回 lua.LNumber 类型
	if seconds <= 0 {
		L.RaiseError("Sleep duration must be greater than zero")
		return 0
	}

	// 将 lua.LNumber 转换为 float64
	duration := time.Duration(float64(seconds) * float64(time.Second))
	time.Sleep(duration)

	fmt.Printf("Slept for %.2f seconds\n", float64(seconds))
	return 0
}

// (4)页面截图
// 捕获页面截图或保存 HTML。
// screenshot(path)，截取当前页面并保存为图片文件。
// screenshot_element(selector, path)，截取指定元素并保存为图片文件。
// save_html(path)，保存当前页面的 HTML 源代码到文件。
func screenshot(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取保存路径参数
	path := L.CheckString(1)
	if path == "" {
		L.RaiseError("Path cannot be empty")
		return 0
	}

	// 截取页面截图
	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(path),
	})
	if err != nil {
		L.RaiseError("Failed to take screenshot: %v", err)
		return 0
	}

	fmt.Printf("Screenshot saved to: %s\n", path)
	return 0
}

func screenshot_element(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 selector 和保存路径参数
	selector := L.CheckString(1)
	path := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}
	if path == "" {
		L.RaiseError("Path cannot be empty")
		return 0
	}

	// 找到目标元素
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find element with selector '%s': %v", selector, err)
		return 0
	}

	// 截取元素截图
	_, err = element.Screenshot(playwright.ElementHandleScreenshotOptions{
		Path: playwright.String(path),
	})
	if err != nil {
		L.RaiseError("Failed to take screenshot of element '%s': %v", selector, err)
		return 0
	}

	fmt.Printf("Screenshot of element '%s' saved to: %s\n", selector, path)
	return 0
}

func save_html(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取保存路径参数
	path := L.CheckString(1)
	if path == "" {
		L.RaiseError("Path cannot be empty")
		return 0
	}

	// 获取页面的 HTML 源代码
	content, err := page.Content()
	if err != nil {
		L.RaiseError("Failed to get page content: %v", err)
		return 0
	}

	// 将 HTML 写入文件
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		L.RaiseError("Failed to save HTML to '%s': %v", path, err)
		return 0
	}

	fmt.Printf("HTML saved to: %s\n", path)
	return 0
}

// (5)处理弹窗和对话框
// 处理页面上的弹窗或提示框。
// accept_alert()，接受弹窗（如 alert 或 confirm 框）。
// dismiss_alert()，关闭弹窗。
// set_alert_text(text)，向弹窗输入框中设置文本。
func accept_alert(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 监听弹窗事件并接受弹窗
	page.OnDialog(func(dialog playwright.Dialog) {
		fmt.Printf("Alert detected: %s\n", dialog.Message())
		dialog.Accept()
	})

	fmt.Println("Listening for dialogs to accept...")
	return 0
}

func dismiss_alert(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 监听弹窗事件并关闭弹窗
	page.OnDialog(func(dialog playwright.Dialog) {
		fmt.Printf("Alert detected: %s\n", dialog.Message())
		dialog.Dismiss()
	})

	fmt.Println("Listening for dialogs to dismiss...")
	return 0
}

func set_alert_text(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取输入的文本
	text := L.CheckString(1)

	// 监听弹窗事件并设置文本
	page.OnDialog(func(dialog playwright.Dialog) {
		fmt.Printf("Prompt detected: %s\n", dialog.Message())
		dialog.Accept(text) // 使用 Accept(text) 方法设置文本
	})

	fmt.Printf("Listening for dialogs to set text: %s\n", text)
	return 0
}

// (6)执行 JavaScript
// 支持直接执行自定义 JavaScript。
// execute_script(script)，执行自定义的 JavaScript 脚本。
// evaluate(selector, script)，在指定元素上执行 JavaScript 脚本。
func execute_script(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 JavaScript 脚本
	script := L.CheckString(1)
	if script == "" {
		L.RaiseError("Script cannot be empty")
		return 0
	}

	// 执行 JavaScript 脚本
	result, err := page.Evaluate(script)
	if err != nil {
		L.RaiseError("Failed to execute script: %v", err)
		return 0
	}

	// 将执行结果返回给 Lua
	L.Push(lua.LString(fmt.Sprintf("%v", result)))
	return 1
}

func evaluate(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器和 JavaScript 脚本
	selector := L.CheckString(1)
	script := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}
	if script == "" {
		L.RaiseError("Script cannot be empty")
		return 0
	}

	// 获取指定元素
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find element with selector '%s': %v", selector, err)
		return 0
	}

	// 在指定元素上执行 JavaScript 脚本
	result, err := element.Evaluate(script)
	if err != nil {
		L.RaiseError("Failed to evaluate script on element '%s': %v", selector, err)
		return 0
	}

	// 将执行结果返回给 Lua
	L.Push(lua.LString(fmt.Sprintf("%v", result)))
	return 1
}

// (7)处理文件上传与下载
// 处理文件上传或下载。
// upload_file(selector, file_path) 向指定的文件上传控件上传文件。
// download_file(url, save_path) 下载指定的文件到本地。
func upload_file(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器和文件路径
	selector := L.CheckString(1)
	filePath := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}
	if filePath == "" {
		L.RaiseError("File path cannot be empty")
		return 0
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		L.RaiseError("File does not exist: %s", filePath)
		return 0
	}

	// 找到文件输入控件
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find file input element with selector '%s': %v", selector, err)
		return 0
	}

	// 创建 InputFile 对象
	inputFile := playwright.InputFile{
		Name: filePath, // 文件路径
	}

	// 上传文件
	err = element.SetInputFiles([]playwright.InputFile{inputFile})
	if err != nil {
		L.RaiseError("Failed to upload file '%s': %v", filePath, err)
		return 0
	}

	fmt.Printf("File '%s' uploaded to element '%s'\n", filePath, selector)
	return 0
}
func upload_multiple_files(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器和多个文件路径
	selector := L.CheckString(1)
	files := []string{}
	for i := 2; i <= L.GetTop(); i++ {
		files = append(files, L.CheckString(i))
	}

	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}
	if len(files) == 0 {
		L.RaiseError("No files provided for upload")
		return 0
	}

	// 检查文件是否存在
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			L.RaiseError("File does not exist: %s", file)
			return 0
		}
	}

	// 找到文件输入控件
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find file input element with selector '%s': %v", selector, err)
		return 0
	}

	// 创建 InputFile 对象数组
	inputFiles := []playwright.InputFile{}
	for _, file := range files {
		inputFiles = append(inputFiles, playwright.InputFile{Name: file})
	}

	// 上传文件
	err = element.SetInputFiles(inputFiles)
	if err != nil {
		L.RaiseError("Failed to upload files: %v", err)
		return 0
	}

	fmt.Printf("Files '%v' uploaded to element '%s'\n", files, selector)
	return 0
}
func download_file(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取 URL 和保存路径
	url := L.CheckString(1)
	savePath := L.CheckString(2)
	if url == "" {
		L.RaiseError("URL cannot be empty")
		return 0
	}
	if savePath == "" {
		L.RaiseError("Save path cannot be empty")
		return 0
	}

	// 启用下载功能并捕获下载事件
	download, err := page.ExpectDownload(func() error {
		_, err := page.Goto(url)
		return err
	})
	if err != nil {
		L.RaiseError("Failed to download file from URL '%s': %v", url, err)
		return 0
	}

	// 保存文件到指定路径
	err = download.SaveAs(savePath)
	if err != nil {
		L.RaiseError("Failed to save downloaded file: %v", err)
		return 0
	}

	fmt.Printf("File downloaded from '%s' to '%s'\n", url, savePath)
	return 0
}

// (8)提取数据
// 获取页面上的数据。
// get_attribute(selector, attribute) 获取指定元素的某个属性值。
// get_html(selector) 获取指定元素的 HTML 内容。
// get_all_links() 获取页面上所有链接。
// capture_table(selector) 提取表格数据，返回表格的二维数组。
func get_attribute(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器和属性名
	selector := L.CheckString(1)
	attribute := L.CheckString(2)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}
	if attribute == "" {
		L.RaiseError("Attribute cannot be empty")
		return 0
	}

	// 获取指定元素
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find element with selector '%s': %v", selector, err)
		return 0
	}

	// 获取属性值
	attrValue, err := element.GetAttribute(attribute)
	if err != nil {
		L.RaiseError("Failed to get attribute '%s' of element '%s': %v", attribute, selector, err)
		return 0
	}

	// 将属性值返回给 Lua
	L.Push(lua.LString(attrValue))
	return 1
}
func get_html(L *lua.LState) int {
	// 获取页面对象
	page := safe_page(L)
	if page == nil {
		L.RaiseError("Failed to get page object")
		return 0
	}

	// 检查是否传递了选择器参数
	selector := L.OptString(1, "") // 若未传递参数，默认为空字符串

	var html string
	var err error

	if selector == "" {
		// 如果选择器为空，返回整个页面的 HTML 内容
		html, err = page.Content()
		if err != nil {
			L.RaiseError("Failed to get page content: %v", err)
			return 0
		}
	} else {
		// 如果选择器不为空，获取对应元素的 HTML 内容
		element, err := page.QuerySelector(selector)
		if err != nil {
			L.RaiseError("Failed to find element for selector '%s': %v", selector, err)
			return 0
		}
		if element == nil {
			L.RaiseError("Element not found for selector: %s", selector)
			return 0
		}

		// 获取元素的内部 HTML
		html, err = element.InnerHTML()
		if err != nil {
			L.RaiseError("Failed to get element inner HTML: %v", err)
			return 0
		}
	}

	// 将 HTML 内容返回给 Lua
	L.Push(lua.LString(html))
	return 1
}

func get_all_links(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 获取所有链接元素
	elements, err := page.QuerySelectorAll("a")
	if err != nil {
		L.RaiseError("Failed to get all links: %v", err)
		return 0
	}

	// 提取链接地址
	links := []string{}
	for _, element := range elements {
		href, err := element.GetAttribute("href")
		if err == nil && href != "" {
			links = append(links, href)
		}
	}

	// 将链接列表返回给 Lua
	linkTable := L.NewTable()
	for _, link := range links {
		linkTable.Append(lua.LString(link))
	}
	L.Push(linkTable)
	return 1
}
func capture_table(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 获取表格元素
	tableElement, err := page.QuerySelector(selector)
	if err != nil || tableElement == nil {
		L.RaiseError("Failed to find table element with selector '%s': %v", selector, err)
		return 0
	}

	// 获取所有行元素
	rowElements, err := tableElement.QuerySelectorAll("tr")
	if err != nil {
		L.RaiseError("Failed to get table rows: %v", err)
		return 0
	}

	// 提取表格数据
	tableData := L.NewTable()
	for _, row := range rowElements {
		// 获取每一行的单元格
		cellElements, err := row.QuerySelectorAll("td, th")
		if err != nil {
			L.RaiseError("Failed to get table cells: %v", err)
			return 0
		}

		rowData := L.NewTable()
		for _, cell := range cellElements {
			text, err := cell.TextContent()
			if err == nil {
				rowData.Append(lua.LString(text))
			}
		}
		tableData.Append(rowData)
	}

	// 将表格数据返回给 Lua
	L.Push(tableData)
	return 1
}

// (9)页面状态检查
// 检查页面的状态是否符合预期。
// is_visible(selector) 检查指定元素是否可见。
// is_enabled(selector) 检查指定元素是否可交互。
// is_checked(selector) 检查复选框是否被选中。
// is_selected(selector) 检查下拉菜单中的选项是否被选中。
func is_visible(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 检查元素是否可见
	isVisible, err := page.IsVisible(selector)
	if err != nil {
		L.RaiseError("Failed to check visibility for selector '%s': %v", selector, err)
		return 0
	}

	// 将结果返回给 Lua
	L.Push(lua.LBool(isVisible))
	return 1
}

func is_enabled(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 检查元素是否可交互
	isEnabled, err := page.IsEnabled(selector)
	if err != nil {
		L.RaiseError("Failed to check enabled state for selector '%s': %v", selector, err)
		return 0
	}

	// 将结果返回给 Lua
	L.Push(lua.LBool(isEnabled))
	return 1
}

func is_checked(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 检查复选框是否被选中
	isChecked, err := page.IsChecked(selector)
	if err != nil {
		L.RaiseError("Failed to check if checkbox is checked for selector '%s': %v", selector, err)
		return 0
	}

	// 将结果返回给 Lua
	L.Push(lua.LBool(isChecked))
	return 1
}
func is_selected(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 获取目标元素
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find element with selector '%s': %v", selector, err)
		return 0
	}

	// 检查是否被选中，通过 Evaluate 获取 `selected` 属性
	isSelected, err := element.Evaluate(`element => element.selected`)
	if err != nil {
		L.RaiseError("Failed to check if element is selected for selector '%s': %v", selector, err)
		return 0
	}

	// 将结果返回给 Lua
	L.Push(lua.LBool(isSelected.(bool)))
	return 1
}
func is_aria_selected(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取选择器
	selector := L.CheckString(1)
	if selector == "" {
		L.RaiseError("Selector cannot be empty")
		return 0
	}

	// 获取目标元素
	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		L.RaiseError("Failed to find element with selector '%s': %v", selector, err)
		return 0
	}

	// 检查 `aria-selected` 属性
	ariaSelected, err := element.Evaluate(`element => element.getAttribute('aria-selected') === 'true'`)
	if err != nil {
		L.RaiseError("Failed to check aria-selected for selector '%s': %v", selector, err)
		return 0
	}

	// 将结果返回给 Lua
	L.Push(lua.LBool(ariaSelected.(bool)))
	return 1
}

//(10)多标签页和窗口管理
//操作浏览器的标签页和窗口。
//new_tab(url) 打开一个新的标签页并加载指定 URL。
//close_tab() 关闭当前标签页。
//switch_to_tab(index) 切换到指定的标签页。

func new_tab(L *lua.LState) int {
	browser := safe_browser(L)
	if browser == nil {
		return 0
	}

	// 从 Lua 获取 URL
	url := L.CheckString(1)
	if url == "" {
		L.RaiseError("URL cannot be empty")
		return 0
	}

	// 打开新标签页
	contexts := browser.Contexts()
	if len(contexts) == 0 {
		L.RaiseError("No browser contexts found")
		return 0
	}

	// 使用第一个浏览器上下文
	context := contexts[0]
	page, err := context.NewPage()
	if err != nil {
		L.RaiseError("Failed to open new tab: %v", err)
		return 0
	}

	// 加载指定 URL
	_, err = page.Goto(url)
	if err != nil {
		L.RaiseError("Failed to load URL '%s': %v", url, err)
		return 0
	}

	fmt.Printf("Opened new tab with URL: %s\n", url)
	return 0
}
func close_tab(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 关闭当前标签页
	err := page.Close()
	if err != nil {
		L.RaiseError("Failed to close current tab: %v", err)
		return 0
	}

	fmt.Println("Current tab closed")
	return 0
}

func switch_to_tab(L *lua.LState) int {
	browser := safe_browser(L)
	if browser == nil {
		return 0
	}

	// 从 Lua 获取标签页索引
	index := L.CheckInt(1)
	if index < 0 {
		L.RaiseError("Invalid tab index: %d", index)
		return 0
	}

	// 获取所有的页面
	contexts := browser.Contexts()
	if len(contexts) == 0 {
		L.RaiseError("No browser contexts found")
		return 0
	}

	context := contexts[0]
	pages := context.Pages()
	if index >= len(pages) {
		L.RaiseError("Tab index out of range: %d (total tabs: %d)", index, len(pages))
		return 0
	}

	// 切换到指定的标签页
	page := pages[index]
	err := page.BringToFront()
	if err != nil {
		L.RaiseError("Failed to switch to tab at index %d: %v", index, err)
		return 0
	}

	fmt.Printf("Switched to tab at index: %d\n", index)
	return 0
}

// (11)网络请求与拦截
// 处理页面的网络请求。
// intercept_request(callback) 拦截网络请求，允许修改或阻止。
// block_request(pattern) 阻止与指定模式匹配的网络请求。
// get_response(url) 获取指定请求的响应数据。
func intercept_request(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 检查并获取 Lua 中的回调函数
	callback := L.CheckFunction(1)
	if callback == nil {
		L.RaiseError("Callback function must be provided")
		return 0
	}

	// 设置请求拦截器
	err := page.Route("**/*", func(route playwright.Route) {
		// 获取请求对象
		request := route.Request()

		// 将请求信息传递给 Lua 回调
		L.Push(callback)
		L.Push(lua.LString(request.URL()))
		L.Push(lua.LString(request.Method()))
		L.Push(lua.LString(request.ResourceType()))
		if err := L.PCall(3, 1, nil); err != nil {
			fmt.Printf("Error calling Lua callback: %v\n", err)
			route.Continue()
			return
		}

		// 获取 Lua 回调的返回值
		ret := L.Get(-1)
		L.Pop(1)

		// 检查返回值，决定如何处理请求
		if ret.Type() == lua.LTString {
			newURL := ret.String()
			route.Continue(playwright.RouteContinueOptions{URL: &newURL})
		} else {
			route.Continue()
		}
	})
	if err != nil {
		L.RaiseError("Failed to set request interceptor: %v", err)
		return 0
	}

	fmt.Println("Request interceptor set")
	return 0
}

func block_request(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取模式
	pattern := L.CheckString(1)
	if pattern == "" {
		L.RaiseError("Pattern cannot be empty")
		return 0
	}

	// 设置请求拦截器
	err := page.Route(pattern, func(route playwright.Route) {
		fmt.Printf("Blocked request: %s\n", route.Request().URL())
		route.Abort("blockedbyclient")
	})
	if err != nil {
		L.RaiseError("Failed to block requests with pattern '%s': %v", pattern, err)
		return 0
	}

	fmt.Printf("Blocking requests matching pattern: %s\n", pattern)
	return 0
}

func get_response(L *lua.LState) int {
	page := safe_page(L)
	if page == nil {
		return 0
	}

	// 从 Lua 获取目标 URL
	targetURL := L.CheckString(1)
	if targetURL == "" {
		L.RaiseError("URL cannot be empty")
		return 0
	}

	// 创建一个通道用于等待目标响应
	responseChannel := make(chan playwright.Response, 1)

	// 监听 response 事件
	page.On("response", func(response playwright.Response) {
		if response.URL() == targetURL {
			responseChannel <- response
		}
	})

	// 等待目标响应
	var response playwright.Response
	select {
	case response = <-responseChannel:
		// 收到目标响应，继续处理
	case <-time.After(30 * time.Second): // 设置超时时间为 30 秒
		L.RaiseError("Timed out waiting for response for URL: %s", targetURL)
		return 0
	}

	// 获取响应的文本内容
	body, err := response.Text()
	if err != nil {
		L.RaiseError("Failed to read response body for URL '%s': %v", targetURL, err)
		return 0
	}

	// 将响应数据返回给 Lua
	L.Push(lua.LString(body))
	return 1
}

func get_storage_state(L *lua.LState) int {
	// 获取浏览器对象
	browser := safe_browser(L)
	if browser == nil {
		L.RaiseError("Failed to get browser object")
		return 0
	}

	// 获取上下文索引（从 Lua 参数中获取，默认为第一个上下文）
	contextIndex := L.OptInt(1, 1) - 1 // Lua 索引从 1 开始，Go 数组索引从 0 开始
	contexts := browser.Contexts()

	// 检查索引是否合法
	if contextIndex < 0 || contextIndex >= len(contexts) {
		L.RaiseError("Invalid context index: %d", contextIndex+1)
		return 0
	}

	// 选择指定的上下文
	context := contexts[contextIndex]

	// 获取存储状态
	storageState, err := context.StorageState()
	if err != nil {
		L.RaiseError("Failed to get storage state: %v", err)
		return 0
	}

	// 将存储状态转换为 JSON 字符串并返回
	jsonState, err := json.Marshal(storageState)
	if err != nil {
		L.RaiseError("Failed to marshal storage state to JSON: %v", err)
		return 0
	}

	// 返回 JSON 数据给 Lua
	L.Push(lua.LString(string(jsonState)))
	return 1
}

//	func set_storage_state(L *lua.LState) int {
//		// 获取浏览器对象
//		browser := safe_browser(L)
//		if browser == nil {
//			L.RaiseError("Failed to get browser object")
//			return 0
//		}
//
//		// 获取上下文索引（从 Lua 参数中获取，默认为第一个上下文）
//		contextIndex := L.OptInt(1, 1) - 1 // Lua 索引从 1 开始，Go 数组索引从 0 开始
//		contexts := browser.Contexts()
//
//		// 检查上下文索引是否合法
//		if contextIndex < 0 || contextIndex >= len(contexts) {
//			L.RaiseError("Invalid context index: %d", contextIndex+1)
//			return 0
//		}
//
//		// 获取指定的上下文
//		context := contexts[contextIndex]
//
//		// 获取传入的存储状态 JSON 字符串
//		storageStateJSON := L.CheckString(2)
//
//		// 设置存储状态到指定的上下文
//		err := context.SetStorageState(playwright.BrowserContextSetStorageStateOptions{
//			StorageState: playwright.String(storageStateJSON),
//		})
//		if err != nil {
//			L.RaiseError("Failed to set storage state: %v", err)
//			return 0
//		}
//
//		// 返回成功状态
//		L.Push(lua.LBool(true))
//		return 1
//	}
func get_cookies_string(L *lua.LState) int {
	// 获取浏览器对象
	browser := safe_browser(L)
	if browser == nil {
		L.RaiseError("Failed to get browser object")
		return 0
	}

	// 获取上下文索引（默认为第一个上下文）
	contextIndex := L.OptInt(1, 1) - 1
	contexts := browser.Contexts()

	// 检查上下文索引是否合法
	if contextIndex < 0 || contextIndex >= len(contexts) {
		L.RaiseError("Invalid context index: %d", contextIndex+1)
		return 0
	}

	// 获取目标上下文
	context := contexts[contextIndex]

	// 获取存储状态
	storageState, err := context.StorageState()
	if err != nil {
		L.RaiseError("Failed to get storage state: %v", err)
		return 0
	}

	// 构建 Cookie 字符串
	var cookieParts []string
	for _, cookie := range storageState.Cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	cookieString := strings.Join(cookieParts, "; ")

	// 返回 Cookie 字符串给 Lua
	L.Push(lua.LString(cookieString))
	return 1
}
