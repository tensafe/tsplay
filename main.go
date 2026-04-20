package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	"tsplay/tsplay_core"

	"github.com/c-bata/go-prompt"
	"github.com/chzyer/readline"
	"github.com/playwright-community/playwright-go"
	"github.com/yuin/gopher-lua"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, fn := range tsplay_core.GlobalPlayWrightFunc {
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
	for _, fn := range tsplay_core.GlobalPlayWrightFunc {
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

var g_headless = false
var g_artifactRoot = tsplay_core.DefaultFlowArtifactRoot
var g_browserVideoOutput = ""
var g_browserVideoWidth = 0
var g_browserVideoHeight = 0
var g_browserVideoCooldownMS = 1200

func main() {
	action := flag.String("action", "cli", "Start Cli Mod | Web Mod | GPT Mod | MCP Stdio | MCP Tool | File Server")
	tsfile := flag.String("script", "", "tsplay script file")
	flowfile := flag.String("flow", "", "tsplay flow file")
	addr := flag.String("addr", ":8082", "server listen address")
	flowRoot := flag.String("flow-root", tsplay_core.DefaultMCPFlowPathRoot, "allowed root directory for MCP flow_path")
	artifactRoot := flag.String("artifact-root", tsplay_core.DefaultMCPArtifactRoot, "allowed root directory for MCP file input/output paths")
	serveRoot := flag.String("serve-root", "", "optional local root directory for built-in static file server; when omitted tsplay serves bundled assets from the binary")
	extractRoot := flag.String("extract-root", "tsplay-assets", "target directory for extracting bundled docs/demo/script assets")
	toolName := flag.String("tool", "", "TSPlay MCP tool name for -action mcp-tool")
	argsJSON := flag.String("args-json", "", "JSON object arguments for -action mcp-tool")
	argsFile := flag.String("args-file", "", "JSON file containing arguments for -action mcp-tool")
	recordInput := flag.String("record-input", defaultScreenRecordInput, "ffmpeg avfoundation input spec for -action record-screen, for example 'Capture screen 0:none'")
	recordOutput := flag.String("record-output", defaultScreenRecordOutput, "video output path for -action record-screen")
	recordCommand := flag.String("record-cmd", "", "shell command to run while -action record-screen is recording")
	recordShell := flag.String("record-shell", "/bin/zsh", "shell used to launch -record-cmd")
	recordFrameRate := flag.Int("record-fps", 30, "frame rate for -action record-screen")
	recordSize := flag.String("record-size", "", "optional ffmpeg video size for -action record-screen, for example 1728x1117")
	recordCursor := flag.Bool("record-cursor", true, "whether -action record-screen should capture the mouse cursor")
	recordWarmupMS := flag.Int("record-warmup-ms", 1200, "warmup delay in milliseconds before -record-cmd starts")
	recordCooldownMS := flag.Int("record-cooldown-ms", 900, "cooldown delay in milliseconds after -record-cmd ends")
	recordDurationMS := flag.Int("record-duration-ms", 0, "optional hard limit in milliseconds for ffmpeg recording duration")
	recordCRF := flag.Int("record-crf", 23, "ffmpeg libx264 CRF for -action record-screen")
	recordPreset := flag.String("record-preset", "veryfast", "ffmpeg encoding preset for -action record-screen")
	browserVideoOutput := flag.String("browser-video-output", "", "save Playwright page video to this path when running -flow or -script; use .webm for the cleanest result")
	browserVideoWidth := flag.Int("browser-video-width", 0, "optional browser video width in pixels when using -browser-video-output")
	browserVideoHeight := flag.Int("browser-video-height", 0, "optional browser video height in pixels when using -browser-video-output")
	browserVideoCooldownMS := flag.Int("browser-video-cooldown-ms", 1200, "keep the page open for this many milliseconds before saving -browser-video-output")
	sessionName := flag.String("session-name", "", "saved session name for session management actions")
	storageStatePath := flag.String("storage-state-path", "", "storage state path for save-session actions")
	storageStateJSON := flag.String("storage-state-json", "", "inline storage state JSON for save-session actions")
	profileName := flag.String("profile-name", "", "persistent profile name for save-session actions")
	profileSession := flag.String("profile-session", "", "persistent profile session name for save-session actions")
	sessionFormat := flag.String("session-format", "all", "snippet format for export-session action")
	isheadless := flag.Bool("headless", false, "is hide browser")

	// 解析命令行参数
	flag.Parse()

	g_headless = *isheadless
	g_artifactRoot = *artifactRoot
	g_browserVideoOutput = strings.TrimSpace(*browserVideoOutput)
	g_browserVideoWidth = *browserVideoWidth
	g_browserVideoHeight = *browserVideoHeight
	g_browserVideoCooldownMS = *browserVideoCooldownMS

	if len(*flowfile) != 0 {
		flow, err := loadFlowDefinition(*flowfile)
		if err != nil {
			log.Fatal(err)
		}
		run_flow(flow)
	} else if len(*tsfile) != 0 {
		content, err := loadScriptSource(*tsfile)
		if err != nil {
			log.Fatal(err)
		}
		run_script(content)
	} else {
		switch *action {
		case "cli":
			//fmt.Println("Start As Cli.")
			cli_mode()
		case "gpt":
			fmt.Println("Start As GPT.")
		case "srv":
			fmt.Println("Start As Web.")
			tsplay_core.McpServerMCP(*addr, tsplay_core.TSPlayMCPServerOptions{
				FlowPathRoot: *flowRoot,
				ArtifactRoot: *artifactRoot,
			})
		case "mcp-stdio":
			tsplay_core.McpServerStdio(tsplay_core.TSPlayMCPServerOptions{
				FlowPathRoot: *flowRoot,
				ArtifactRoot: *artifactRoot,
			})
		case "mcp-tool":
			if err := runMCPToolAction(*toolName, *argsJSON, *argsFile, *flowRoot, *artifactRoot); err != nil {
				log.Fatal(err)
			}
		case "file-srv", "demo-srv":
			if err := serveStaticFiles(*addr, *serveRoot); err != nil {
				log.Fatal(err)
			}
		case "list-assets":
			names, err := bundledAssetNames()
			if err != nil {
				log.Fatal(err)
			}
			for _, name := range names {
				fmt.Println(name)
			}
		case "extract-assets":
			count, err := extractBundledAssets(*extractRoot)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Extracted %d bundled assets to %s\n", count, *extractRoot)
		case "list-record-devices":
			probe, err := listScreenRecordDevices()
			if probe != nil {
				printJSON(probe)
			}
			if err != nil {
				log.Fatal(err)
			}
		case "record-screen":
			result, err := runScreenRecordAction(screenRecordOptions{
				InputSpec:     *recordInput,
				OutputPath:    *recordOutput,
				Command:       *recordCommand,
				Shell:         *recordShell,
				FrameRate:     *recordFrameRate,
				VideoSize:     *recordSize,
				CaptureCursor: *recordCursor,
				Warmup:        time.Duration(*recordWarmupMS) * time.Millisecond,
				Cooldown:      time.Duration(*recordCooldownMS) * time.Millisecond,
				MaxDuration:   time.Duration(*recordDurationMS) * time.Millisecond,
				CRF:           *recordCRF,
				Preset:        *recordPreset,
			})
			if result != nil {
				printJSON(result)
			}
			if err != nil {
				log.Fatal(err)
			}
		case "save-session":
			if strings.TrimSpace(*sessionName) == "" {
				log.Fatal("-session-name is required for -action save-session")
			}
			session, err := tsplay_core.SaveFlowSavedSession(tsplay_core.FlowSavedSessionSaveOptions{
				Name:             *sessionName,
				ArtifactRoot:     *artifactRoot,
				StorageStateJSON: *storageStateJSON,
				StorageStatePath: *storageStatePath,
				Profile:          *profileName,
				Session:          *profileSession,
			})
			if err != nil {
				log.Fatal(err)
			}
			printJSON(tsplay_core.BuildFlowSavedSessionDetail(*session, *artifactRoot))
		case "list-sessions":
			sessions, err := tsplay_core.ListFlowSavedSessions(*artifactRoot)
			if err != nil {
				log.Fatal(err)
			}
			items := make([]map[string]any, 0, len(sessions))
			for _, session := range sessions {
				items = append(items, tsplay_core.BuildFlowSavedSessionView(session, *artifactRoot))
			}
			printJSON(map[string]any{
				"artifact_root": *artifactRoot,
				"sessions":      items,
			})
		case "get-session":
			if strings.TrimSpace(*sessionName) == "" {
				log.Fatal("-session-name is required for -action get-session")
			}
			session, err := tsplay_core.LoadFlowSavedSession(*sessionName, *artifactRoot)
			if err != nil {
				log.Fatal(err)
			}
			printJSON(tsplay_core.BuildFlowSavedSessionDetail(*session, *artifactRoot))
		case "export-session":
			if strings.TrimSpace(*sessionName) == "" {
				log.Fatal("-session-name is required for -action export-session")
			}
			session, err := tsplay_core.LoadFlowSavedSession(*sessionName, *artifactRoot)
			if err != nil {
				log.Fatal(err)
			}
			exported, err := tsplay_core.ExportFlowSavedSessionFlowSnippet(*session, *artifactRoot, *sessionFormat)
			if err != nil {
				log.Fatal(err)
			}
			printJSON(exported)
		case "delete-session":
			if strings.TrimSpace(*sessionName) == "" {
				log.Fatal("-session-name is required for -action delete-session")
			}
			deleted, err := tsplay_core.DeleteFlowSavedSession(*sessionName, *artifactRoot)
			if err != nil {
				log.Fatal(err)
			}
			printJSON(deleted)
		}
	}
}

func printJSON(value any) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(encoded))
}

func run_flow(flow *tsplay_core.Flow) {
	result, err := tsplay_core.RunFlow(flow, tsplay_core.FlowRunOptions{
		Headless:               g_headless,
		ArtifactRoot:           g_artifactRoot,
		BrowserVideoOutputPath: g_browserVideoOutput,
		BrowserVideoWidth:      g_browserVideoWidth,
		BrowserVideoHeight:     g_browserVideoHeight,
		BrowserVideoCooldownMS: g_browserVideoCooldownMS,
	})
	if result != nil {
		encoded, marshalErr := json.MarshalIndent(result, "", "  ")
		if marshalErr != nil {
			log.Printf("could not encode flow result: %v", marshalErr)
		} else {
			fmt.Println(string(encoded))
		}
	}
	if err != nil {
		log.Fatalf("error running flow: %v", err)
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
	for _, fn := range tsplay_core.GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}
	L.SetGlobal("artifact_root", lua.LString(g_artifactRoot))

	var pw *playwright.Playwright
	var browser playwright.Browser
	var page playwright.Page

	clearPlaywrightGlobals := func() {
		L.SetGlobal("browser", lua.LNil)
		L.SetGlobal("context", lua.LNil)
		L.SetGlobal("page", lua.LNil)
	}
	setPlaywrightGlobals := func() {
		ud_b := L.NewUserData()
		ud_b.Value = browser
		L.SetGlobal("browser", ud_b)

		ud_c := L.NewUserData()
		ud_c.Value = page.Context()
		L.SetGlobal("context", ud_c)

		ud_p := L.NewUserData()
		ud_p.Value = page
		L.SetGlobal("page", ud_p)
	}
	clearPlaywrightGlobals()

	stopPlaywright := func() {
		clearPlaywrightGlobals()
		if page != nil {
			if err := page.Close(); err != nil {
				log.Printf("failed to close page: %v", err)
			}
			page = nil
		}
		if browser != nil {
			if err := browser.Close(); err != nil {
				log.Printf("failed to close browser: %v", err)
			}
			browser = nil
		}
		if pw != nil {
			if err := pw.Stop(); err != nil {
				log.Printf("failed to stop Playwright runtime: %v", err)
			}
			pw = nil
		}
	}
	defer stopPlaywright()

	ensurePlaywrightRuntime := func(reason string) error {
		if pw != nil {
			return nil
		}
		if strings.TrimSpace(reason) == "" {
			fmt.Println("Starting Playwright runtime...")
		} else {
			fmt.Printf("Starting Playwright runtime because %s...\n", reason)
		}
		runtimeHandle, err := tsplay_core.StartPlaywright()
		if err != nil {
			return err
		}
		pw = runtimeHandle
		return nil
	}

	ensurePlaywrightPage := func(reason string) error {
		if page != nil && browser != nil {
			setPlaywrightGlobals()
			return nil
		}
		if err := ensurePlaywrightRuntime(reason); err != nil {
			return err
		}
		if strings.TrimSpace(reason) == "" {
			fmt.Println("Launching Playwright browser and page...")
		} else {
			fmt.Printf("Launching Playwright browser and page because %s...\n", reason)
		}
		launchedBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(g_headless),
		})
		if err != nil {
			return fmt.Errorf("could not launch browser: %v", err)
		}
		launchedPage, err := launchedBrowser.NewPage()
		if err != nil {
			_ = launchedBrowser.Close()
			return fmt.Errorf("could not create page: %v", err)
		}
		browser = launchedBrowser
		page = launchedPage
		setPlaywrightGlobals()
		fmt.Println("Playwright initialized. Browser, context, and page are ready.")
		return nil
	}
	fmt.Println("Please input the 'start' command to run and launch tsplay")

	var rl *readline.Instance
	if os_type == "windows" {
		var err error
		rl, err = readline.NewEx(&readline.Config{
			Prompt:       "> ",
			AutoComplete: createReadlineCompleter(),
		})
		if err != nil {
			log.Printf("failed to initialize readline: %v", err)
		}
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
			if pw == nil && browser == nil && page == nil {
				fmt.Println("Playwright is already idle.")
				continue
			}
			stopPlaywright()
			fmt.Println("Playwright has been reset. It will start again on the next 'start' command or browser action.")
			continue
		}

		// 处理 start 命令
		if input == "start" {
			if err := ensurePlaywrightPage("the CLI start command was requested"); err != nil {
				fmt.Printf("Failed to initialize Playwright: %v\n", err)
			} else {
				fmt.Println("Playwright started. Browser and page objects are now available in Lua.")
			}
			continue
		}

		runLuaScript := func(script string) {
			usage := tsplay_core.AnalyzeLuaScriptPlaywrightUsage(script)
			var err error
			switch {
			case usage.NeedsBrowser():
				err = ensurePlaywrightPage(usage.Summary(3))
			case usage.NeedsRuntime:
				err = ensurePlaywrightRuntime(usage.Summary(3))
			}
			if err != nil {
				fmt.Printf("Failed to initialize Playwright: %v\n", err)
				return
			}
			if err := L.DoString(script); err != nil {
				fmt.Printf("Lua error: %v\n", err)
			}
		}

		// 处理 Lua 脚本
		if strings.HasPrefix(input, "lua ") {
			script := strings.TrimPrefix(input, "lua ")
			runLuaScript(script)
			continue
		}
		// 默认行为：将输入内容作为 Lua 脚本执行
		if input != "" {
			runLuaScript(input)
		}
	}
}

func run_script(script string) {
	// 创建 Lua 状态机
	L := lua.NewState()
	defer L.Close()

	// 注册 Go 函数到 Lua
	for _, fn := range tsplay_core.GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}
	L.SetGlobal("artifact_root", lua.LString(g_artifactRoot))

	usage := tsplay_core.AnalyzeLuaScriptPlaywrightUsage(script)
	var pw *playwright.Playwright
	var browser playwright.Browser
	var page playwright.Page
	var browserVideo *tsplay_core.BrowserVideoRecording
	setPlaywrightGlobals := func() {
		ud_b := L.NewUserData()
		ud_b.Value = browser
		L.SetGlobal("browser", ud_b)

		ud_c := L.NewUserData()
		ud_c.Value = page.Context()
		L.SetGlobal("context", ud_c)

		ud_p := L.NewUserData()
		ud_p.Value = page
		L.SetGlobal("page", ud_p)
	}
	stopPlaywright := func() {
		L.SetGlobal("browser", lua.LNil)
		L.SetGlobal("context", lua.LNil)
		L.SetGlobal("page", lua.LNil)
		if page != nil {
			_ = page.Close()
			page = nil
		}
		if browser != nil {
			_ = browser.Close()
			browser = nil
		}
		if pw != nil {
			_ = pw.Stop()
			pw = nil
		}
	}
	defer stopPlaywright()

	if usage.NeedsRuntime {
		var err error
		pw, err = tsplay_core.StartPlaywright()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	if usage.NeedsBrowser() {
		if pw == nil {
			var err error
			pw, err = tsplay_core.StartPlaywright()
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
		var err error
		browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(g_headless),
		})
		if err != nil {
			log.Fatalf("could not launch browser: %v", err)
		}
		browserVideo, err = tsplay_core.PrepareBrowserVideoRecording(g_browserVideoOutput, g_browserVideoWidth, g_browserVideoHeight)
		if err != nil {
			log.Fatalf("could not prepare browser video: %v", err)
		}
		if browserVideo != nil {
			page, err = browser.NewPage(playwright.BrowserNewPageOptions{
				RecordVideo: browserVideo.RecordVideo,
			})
		} else {
			page, err = browser.NewPage()
		}
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}
		setPlaywrightGlobals()
	}

	if err := L.DoString(script); err != nil {
		log.Fatalf("error running Lua script: %v", err)
	}

	if !usage.NeedsBrowser() {
		return
	}

	if browserVideo != nil && page != nil {
		if g_browserVideoCooldownMS > 0 {
			time.Sleep(time.Duration(g_browserVideoCooldownMS) * time.Millisecond)
		}
		savedPath, err := tsplay_core.SaveBrowserVideo(page, browserVideo.OutputPath)
		if err != nil {
			log.Fatalf("could not save browser video: %v", err)
		}
		fmt.Printf("saved browser video: %s\n", savedPath)
		page = nil
		return
	}

	if g_headless {
		return
	}

	// 捕捉系统信号，以便优雅地关闭程序
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号以便优雅地退出
	<-sigChan
}
