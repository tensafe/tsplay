package main

import (
	"github.com/playwright-community/playwright-go"
	"github.com/yuin/gopher-lua"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tsplay/tsplay"
)

func main() {
	// 初始化 Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	// 启动浏览器并打开新页面
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// 创建 Lua 状态机
	L := lua.NewState()
	defer L.Close()

	// 将 Playwright Page 对象传递给 Lua
	ud_b := L.NewUserData()
	ud_b.Value = browser
	L.SetGlobal("browser", ud_b)

	// 将 Playwright Page 对象传递给 Lua
	ud_p := L.NewUserData()
	ud_p.Value = page
	L.SetGlobal("page", ud_p)

	// 注册 Go 函数到 Lua
	for _, fn := range tsplay.GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}

	// 执行 Lua 脚本
	script := `
        -- 使用 Lua 调用注册的函数
        navigate("https://www.baidu.com")
		type_text("#kw","山东")
        click("#su")
		wait_for_network_idle()
		html = get_html()
		print(html)
    `
	script = `-- 打开百度首页
navigate("https://www.baidu.com")

-- 等待搜索框加载完成
wait_for_selector("#kw")

-- 在搜索框中输入“山东”
type_text("#kw", "山东")

-- 点击“百度一下”按钮
click("#su")

-- 等待搜索结果页面加载完成
wait_for_navigation()
`
	if err := L.DoString(script); err != nil {
		log.Fatalf("error running Lua script: %v", err)
	}

	// 捕捉系统信号，以便优雅地关闭程序
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号以便优雅地退出
	<-sigChan
}
