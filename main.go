package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/chzyer/readline"
	"github.com/playwright-community/playwright-go"
	"github.com/yuin/gopher-lua"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"tsplay/tsplay"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, fn := range tsplay.GlobalPlayWrightFunc {
		sug := prompt.Suggest{
			Text:        fn.Name,
			Description: fn.Description_en,
		}
		s = append(s, sug)
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func createReadlineCompleter() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface
	for _, fn := range tsplay.GlobalPlayWrightFunc {
		items = append(items, &readline.PrefixCompleter{
			Name:    []rune(fn.Name),
			Dynamic: false,
		})
	}
	items = append(items, &readline.PrefixCompleter{
		Name:    []rune("start"),
		Dynamic: false,
	})
	items = append(items, &readline.PrefixCompleter{
		Name:    []rune("reset"),
		Dynamic: false,
	})
	items = append(items, &readline.PrefixCompleter{
		Name:    []rune("exit"),
		Dynamic: false,
	})
	return readline.NewPrefixCompleter(items...)
}

func main() {
	action := flag.String("action", "cli", "Start Cli Mod | Web Mod | GPT Mod")
	err := playwright.Install()
	if err != nil {
		log.Println("could not install playwright browsers: %v", err)
	}
	// 解析命令行参数
	flag.Parse()
	switch *action {
	case "cli":
		//fmt.Println("Start As Cli.")
		cli_mode()
	case "gpt":
		fmt.Println("Start As GPT.")
	case "srv":
		fmt.Println("Start As Web.")
	}

}

func cli_mode() {
	os_type := "windows"
	switch runtime.GOOS {
	case "windows":
		os_type = "windows"
	case "darwin":
		os_type = "darwin"
	case "linux":
		os_type = "linux"
	default:
		os_type = "windows"
	}

	L := lua.NewState()
	defer L.Close()

	// 注册 Go 函数到 Lua
	for _, fn := range tsplay.GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	var browser playwright.Browser
	var page playwright.Page

	// 初始化浏览器和页面
	initPlaywright := func() error {
		var err error
		browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("could not launch browser: %v", err)
		}
		page, err = browser.NewPage()
		if err != nil {
			return fmt.Errorf("could not create page: %v", err)
		}
		fmt.Println("Playwright initialized. Browser and page are ready.")
		return nil
	}

	// 将浏览器和页面对象传递给 Lua
	startPlaywright := func() {
		if browser == nil || page == nil {
			fmt.Println("Playwright is not initialized. Initializing now...")
			if err := initPlaywright(); err != nil {
				fmt.Printf("Failed to initialize Playwright: %v\n", err)
				return
			}
		}
		// 将 Playwright 对象传递给 Lua
		ud_b := L.NewUserData()
		ud_b.Value = browser
		L.SetGlobal("browser", ud_b)

		ud_p := L.NewUserData()
		ud_p.Value = page
		L.SetGlobal("page", ud_p)
		fmt.Println("Playwright started. Browser and page objects are now available in Lua.")
	}
	fmt.Println("Please input the 'start' command to run and launch tsplay")

	var rl *readline.Instance
	if os_type == "windows" {
		rl, err = readline.NewEx(&readline.Config{
			Prompt:       "> ",
			AutoComplete: createReadlineCompleter(),
		})
	}

	if rl != nil {
		defer rl.Close()
	}

	for {
		// 动态 CLI 提示符
		prefix := "> "
		if page != nil {
			prefix = "(playwright) > "
		}

		input := ""
		if os_type == "windows" {
			line, err := rl.Readline()
			if err != nil { // 处理 Ctrl+D 或 Ctrl+C
				break
			}
			input = line
		} else {
			// 启动 prompt
			input = prompt.Input(prefix, completer)
		}
		// 检查输入是否为 exit
		if input == "exit" {
			fmt.Println("Exiting the shell. Goodbye!")
			break
		}

		// 处理 reset 命令
		if input == "reset" {
			fmt.Println("Resetting Playwright...")
			if browser != nil {
				if err := browser.Close(); err != nil {
					log.Printf("failed to close browser: %v", err)
				}
				browser = nil
				page = nil
			}
			if err := initPlaywright(); err != nil {
				log.Printf("Failed to reset Playwright: %v\n", err)
				continue
			}
			startPlaywright()
			continue
		}

		// 处理 start 命令
		if input == "start" {
			startPlaywright()
			continue
		}

		// 处理 Lua 脚本
		if strings.HasPrefix(input, "lua ") {
			script := strings.TrimPrefix(input, "lua ")
			if err := L.DoString(script); err != nil {
				fmt.Printf("Lua error: %v\n", err)
			}
			continue
		}
		// 默认行为：将输入内容作为 Lua 脚本执行
		if input != "" {
			if err := L.DoString(input); err != nil {
				fmt.Printf("Lua error: %v\n", err)
			}
		}
	}
}

func main_old() {
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
wait_for_network_idle()

screenshot("./data.png")

print(get_storage_state())

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
