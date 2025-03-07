package tsplay

import (
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"strings"
)

// ElementInfo 表示页面元素的信息
type ElementInfo struct {
	TagName     string            `json:"tagName"`
	ID          string            `json:"id"`
	Class       string            `json:"class"`
	Text        string            `json:"text"`
	XPath       string            `json:"xpath"`
	Role        string            `json:"role"`
	Attributes  map[string]string `json:"attributes"`
	Description string            `json:"description"`
}

// call
// ExtractPageToJson 提取页面元素信息并返回 JSON 字符串
func ExtractPageToJson(page playwright.Page) string {
	// 执行 JavaScript 来获取 DOM 树
	domTree, err := page.Evaluate(`() => {
		function getElementTree(element) {
			return {
				tag: element.tagName,
				attributes: Array.from(element.attributes).reduce((attrs, attr) => {
					attrs[attr.name] = attr.value;
					return attrs;
				}, {}),
				children: Array.from(element.children).map(getElementTree)
			};
		}
		return getElementTree(document.documentElement);
	}`)
	if err != nil {
		log.Fatalf("could not evaluate DOM tree: %v", err)
	}

	// 将结果格式化为 JSON 并打印
	domTreeJSON, err := json.MarshalIndent(domTree, "", "  ")
	if err != nil {
		log.Fatalf("could not marshal DOM tree: %v", err)
	}
	log.Printf("DOM Tree:\n%s", domTreeJSON)
	return string(domTreeJSON)
}

func ExtractElementWithXPath(page playwright.Page) string {
	// 执行 JavaScript，获取 DOM 树并生成 XPath
	// 执行 JavaScript，获取 DOM 树并生成 XPath
	domTreeWithXPath, err := page.Evaluate(`() => {
		function generateXPath(element) {
			if (element.id) {
				return '//*[@id="' + element.id + '"]';
			}
			if (element === document.documentElement) { // 根节点 <html>
				return '/html';
			}
			if (element === document.body) { // <body> 节点
				return '/html/body';
			}

			let ix = 0;
			const siblings = element.parentNode ? Array.from(element.parentNode.childNodes) : [];
			for (let i = 0; i < siblings.length; i++) {
				const sibling = siblings[i];
				if (sibling === element) {
					return generateXPath(element.parentNode) + '/' + element.tagName.toLowerCase() + '[' + (ix + 1) + ']';
				}
				if (sibling.nodeType === 1 && sibling.tagName === element.tagName) {
					ix++;
				}
			}
		}

		function getElementTreeWithXPath(element) {
			return {
				tag: element.tagName,
				xpath: generateXPath(element),
				attributes: Array.from(element.attributes).reduce((attrs, attr) => {
					attrs[attr.name] = attr.value;
					return attrs;
				}, {}),
				children: Array.from(element.children).map(getElementTreeWithXPath)
			};
		}
		return getElementTreeWithXPath(document.documentElement);
	}`)
	if err != nil {
		log.Fatalf("could not evaluate DOM tree with XPath: %v", err)
	}

	// 格式化为 JSON 输出
	domTreeWithXPathJSON, err := json.MarshalIndent(domTreeWithXPath, "", "  ")
	if err != nil {
		log.Fatalf("could not marshal DOM tree with XPath: %v", err)
	}
	log.Printf("DOM Tree with XPath:\n%s", domTreeWithXPathJSON)
	return string(domTreeWithXPathJSON)
}

func ExtractSimplifiedElementWithXPath(page playwright.Page) string {
	// 执行 JavaScript，获取 DOM 树并生成简化的 XPath 树，过滤掉不可见和无用的部分
	simplifiedDomTreeWithXPath, err := page.Evaluate(`() => {
		function generateXPath(element) {
			if (element.id) {
				return '//*[@id="' + element.id + '"]';
			}
			if (element === document.documentElement) { // 根节点 <html>
				return '/html';
			}
			if (element === document.body) { // <body> 节点
				return '/html/body';
			}

			let ix = 0;
			const siblings = element.parentNode ? Array.from(element.parentNode.childNodes) : [];
			for (let i = 0; i < siblings.length; i++) {
				const sibling = siblings[i];
				if (sibling === element) {
					return generateXPath(element.parentNode) + '/' + element.tagName.toLowerCase() + '[' + (ix + 1) + ']';
				}
				if (sibling.nodeType === 1 && sibling.tagName === element.tagName) {
					ix++;
				}
			}
		}

		function isVisible(element) {
			// 检查元素是否在屏幕上可见
			const style = window.getComputedStyle(element);
			return (
				style.display !== 'none' &&
				style.visibility !== 'hidden' &&
				style.opacity !== '0' &&
				element.offsetHeight > 0 &&
				element.offsetWidth > 0
			);
		}

		function truncateText(text, length) {
			// 如果文本长度超出限制，则截断并添加省略号
			return text.length > length ? text.slice(0, length) + '...' : text;
		}

		function getSimplifiedElementTree(element) {
			// 获取文本内容并清理首尾空格
			let text = element.textContent.trim();
			// 对文本内容进行截断，最长 128 个字符
			text = truncateText(text, 128);

			// 如果是无用的标签，元素不可见，或者文本为空，则跳过该节点
			if (
				element.tagName.toLowerCase() === 'style' ||
				element.tagName.toLowerCase() === 'script' ||
				element.tagName.toLowerCase() === 'head' ||
				element.tagName.toLowerCase() === 'textarea' ||
				!isVisible(element) ||
				text === ''
			) {
				return null;
			}

			// 如果是 <a> 标签，提取 href 属性
			let href = null;
			if (element.tagName.toLowerCase() === 'a') {
				href = element.getAttribute('href') || null;
			}

			// 构造当前节点信息
			const current = {
				tag: element.tagName,
				xpath: generateXPath(element),
				text: text,
				href: href, // 添加 href 属性，非 <a> 标签时为 null
				children: Array.from(element.children)
					.map(getSimplifiedElementTree) // 递归获取子元素
					.filter(child => child !== null) // 过滤掉无效的子元素
			};

			return current;
		}

		return getSimplifiedElementTree(document.documentElement);
	}`)
	if err != nil {
		log.Fatalf("could not evaluate simplified DOM tree with XPath: %v", err)
	}

	// 格式化为 JSON 输出
	simplifiedDomTreeWithXPathJSON, err := json.MarshalIndent(simplifiedDomTreeWithXPath, "", "  ")
	if err != nil {
		log.Fatalf("could not marshal simplified DOM tree with XPath: %v", err)
	}
	log.Printf("Simplified DOM Tree with XPath:\n%s", simplifiedDomTreeWithXPathJSON)
	return string(simplifiedDomTreeWithXPathJSON)
}

func DemoBaidu() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("无法启动 Playwright: %v", err)
	}
	defer pw.Stop()

	// 启动浏览器
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("无法启动浏览器: %v", err)
	}
	defer browser.Close()

	// 打开新页面
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("无法创建新页面: %v", err)
	}

	// 导航到目标页面
	_, err = page.Goto("https://www.baidu.com")
	if err != nil {
		log.Fatalf("无法导航到页面: %v", err)
	}

	page.WaitForLoadState()

	// 提取页面元素信息
	rootElement := ExtractSimplifiedElementWithXPath(page)
	if err != nil {
		log.Fatalf("无法提取页面元素: %v", err)
	}

	// 将元素信息转换为 JSON
	jsonData, err := json.MarshalIndent(rootElement, "", "  ")
	if err != nil {
		log.Fatalf("无法转换为 JSON: %v", err)
	}

	// 打印 JSON 数据
	fmt.Println(string(jsonData))
}

// skipTag 判断是否跳过某些无意义的标签
func skipTag(tagName string) bool {
	skipTags := []string{"SCRIPT", "STYLE", "META", "LINK", "HEAD", "IFRAME", "NOSCRIPT", "BODY", "TITLE", "HTML"}
	for _, tag := range skipTags {
		if tagName == tag {
			return true
		}
	}
	return false
}

// isEmptyText 判断文本是否为空（去除空格）
func isEmptyText(text string) bool {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\r\n", " ")
	for _, char := range text {
		if char != ' ' && char != '\n' && char != '\t' {
			return false
		}
	}
	return true
}
