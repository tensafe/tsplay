package tsplay_core

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const defaultCDPLaunchTimeout = 15 * time.Second

var localCDPProfileCounter atomic.Uint64

type localCDPBrowser struct {
	cmd      *exec.Cmd
	waitDone chan error
	close    sync.Once
}

func normalizeCDPEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", nil
	}
	if !strings.Contains(endpoint, "://") {
		if _, _, err := net.SplitHostPort(endpoint); err == nil || looksLikeBareCDPEndpoint(endpoint) {
			endpoint = "http://" + endpoint
		}
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse browser.cdp_endpoint %q: %w", endpoint, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("browser.cdp_endpoint must be an http(s) or ws(s) URL, or host:port[/path]")
	}
	switch parsed.Scheme {
	case "http", "https", "ws", "wss":
	default:
		return "", fmt.Errorf("browser.cdp_endpoint must use http, https, ws, or wss")
	}
	if portText := parsed.Port(); portText != "" {
		port, err := strconv.Atoi(portText)
		if err != nil || port < 1 || port > 65535 {
			return "", fmt.Errorf("browser.cdp_endpoint has invalid port %q", portText)
		}
	}
	if parsed.Scheme == "http" || parsed.Scheme == "https" {
		if parsed.Path == "/json" || strings.HasPrefix(parsed.Path, "/json/") || strings.HasPrefix(parsed.Path, "/devtools/") {
			parsed.Path = ""
		}
		parsed.RawQuery = ""
		parsed.Fragment = ""
	}
	return parsed.String(), nil
}

func looksLikeBareCDPEndpoint(endpoint string) bool {
	parsed, err := url.Parse("http://" + endpoint)
	if err != nil || parsed.Host == "" || parsed.Port() == "" {
		return false
	}
	return true
}

func ensureLocalCDPBrowser(config FlowBrowserConfig, options FlowRunOptions) (*localCDPBrowser, string, error) {
	endpoint, err := config.cdpEndpointURL()
	if err != nil {
		return nil, "", err
	}
	if endpoint != "" && !cdpEndpointIsLocal(endpoint) {
		return nil, "", fmt.Errorf("browser.cdp_launch can only start or reuse a local browser; endpoint %q is not local", endpoint)
	}
	if endpoint != "" && cdpEndpointReachable(endpoint, 750*time.Millisecond) {
		return nil, endpoint, nil
	}

	port := config.CDPPort
	if port == 0 && endpoint != "" {
		port, err = cdpPortFromEndpoint(endpoint)
		if err != nil {
			return nil, "", fmt.Errorf("browser.cdp_launch with browser.cdp_endpoint requires a local endpoint with an explicit port: %w", err)
		}
	}
	if port == 0 {
		port, err = freeLocalTCPPort()
		if err != nil {
			return nil, "", err
		}
	}

	endpoint = fmt.Sprintf("http://127.0.0.1:%d", port)
	if cdpEndpointReachable(endpoint, 750*time.Millisecond) {
		return nil, endpoint, nil
	}

	executable, err := resolveLocalCDPBrowserExecutable(config.CDPExecutable)
	if err != nil {
		return nil, "", err
	}
	userDataDir, err := config.cdpLaunchUserDataDir(options)
	if err != nil {
		return nil, "", err
	}
	launched, err := startLocalCDPBrowser(executable, userDataDir, port, config)
	if err != nil {
		return nil, "", err
	}
	timeout := defaultCDPLaunchTimeout
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Millisecond
	}
	if err := waitForCDPEndpoint(endpoint, timeout, launched.waitDone); err != nil {
		_ = launched.Close()
		return nil, "", err
	}
	return launched, endpoint, nil
}

func startLocalCDPBrowser(executable string, userDataDir string, port int, config FlowBrowserConfig) (*localCDPBrowser, error) {
	args := []string{
		fmt.Sprintf("--remote-debugging-port=%d", port),
		"--user-data-dir=" + userDataDir,
		"--no-first-run",
		"--no-default-browser-check",
	}
	if config.headlessValue() {
		args = append(args, "--headless=new", "--disable-gpu")
	}
	args = append(args, "about:blank")

	cmd := exec.Command(executable, args...)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	configureLocalCDPBrowserCommand(cmd)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start browser %q for CDP launch: %w", executable, err)
	}
	launched := &localCDPBrowser{
		cmd:      cmd,
		waitDone: make(chan error, 1),
	}
	go func() {
		launched.waitDone <- cmd.Wait()
	}()
	return launched, nil
}

func (browser *localCDPBrowser) Close() error {
	if browser == nil || browser.cmd == nil {
		return nil
	}
	var closeErr error
	browser.close.Do(func() {
		select {
		case <-browser.waitDone:
			return
		default:
		}
		if err := terminateLocalCDPBrowserCommand(browser.cmd); err != nil {
			closeErr = err
		}
		select {
		case <-browser.waitDone:
			return
		case <-time.After(3 * time.Second):
		}
		if err := killLocalCDPBrowserCommand(browser.cmd); err != nil && closeErr == nil {
			closeErr = err
		}
		select {
		case <-browser.waitDone:
		case <-time.After(2 * time.Second):
		}
	})
	return closeErr
}

func waitForCDPEndpoint(endpoint string, timeout time.Duration, waitDone <-chan error) error {
	if timeout <= 0 {
		timeout = defaultCDPLaunchTimeout
	}
	deadline := time.Now().Add(timeout)
	for {
		select {
		case err := <-waitDone:
			if err != nil {
				return fmt.Errorf("browser exited before CDP endpoint became ready: %w", err)
			}
			return fmt.Errorf("browser exited before CDP endpoint became ready")
		default:
		}
		if cdpEndpointReachable(endpoint, 500*time.Millisecond) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for CDP endpoint %q", endpoint)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func cdpEndpointReachable(endpoint string, timeout time.Duration) bool {
	base, err := cdpHTTPBase(endpoint)
	if err != nil {
		return false
	}
	if timeout <= 0 {
		timeout = 500 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(base, "/")+"/json/version", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func cdpHTTPBase(endpoint string) (string, error) {
	endpoint, err := normalizeCDPEndpoint(endpoint)
	if err != nil {
		return "", err
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	switch parsed.Scheme {
	case "ws":
		parsed.Scheme = "http"
	case "wss":
		parsed.Scheme = "https"
	case "http", "https":
	default:
		return "", fmt.Errorf("CDP endpoint %q is not http(s) or ws(s)", endpoint)
	}
	parsed.Path = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func cdpPortFromEndpoint(endpoint string) (int, error) {
	base, err := cdpHTTPBase(endpoint)
	if err != nil {
		return 0, err
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return 0, err
	}
	portText := parsed.Port()
	if portText == "" {
		return 0, fmt.Errorf("endpoint %q has no port", endpoint)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("endpoint %q has invalid port %q", endpoint, portText)
	}
	return port, nil
}

func cdpEndpointIsLocal(endpoint string) bool {
	base, err := cdpHTTPBase(endpoint)
	if err != nil {
		return false
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return false
	}
	host := strings.Trim(strings.ToLower(parsed.Hostname()), "[]")
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func freeLocalTCPPort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("find free local CDP port: %w", err)
	}
	defer listener.Close()
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("find free local CDP port: unexpected addr %T", listener.Addr())
	}
	return addr.Port, nil
}

func (browser FlowBrowserConfig) cdpLaunchUserDataDir(options FlowRunOptions) (string, error) {
	dir := strings.TrimSpace(browser.CDPUserDataDir)
	if dir == "" {
		runSegment := strings.TrimSpace(options.RunID)
		if runSegment == "" {
			runSegment = fmt.Sprintf("run-%s-%d", time.Now().Format("20060102-150405.000000000"), localCDPProfileCounter.Add(1))
		}
		dir = filepath.Join("browser-state", "cdp-launch", sanitizeArtifactSegment(runSegment))
	}
	var resolved string
	var err error
	if options.Security != nil {
		resolved, err = resolveFlowBrowserStatePath(dir, flowFileOutputPath, options.Security)
		if err != nil {
			return "", err
		}
	} else if filepath.IsAbs(expandHomePath(dir)) {
		resolved, err = filepath.Abs(expandHomePath(dir))
		if err != nil {
			return "", err
		}
	} else {
		rootReal, err := prepareRuntimeFileRoot(flowBrowserStateRoot(options))
		if err != nil {
			return "", fmt.Errorf("resolve CDP launch profile root: %w", err)
		}
		resolved, err = filepath.Abs(filepath.Join(rootReal, dir))
		if err != nil {
			return "", err
		}
		if err := ensurePathInsideRoot(resolved, rootReal); err != nil {
			return "", fmt.Errorf("resolve browser.cdp_user_data_dir %q: %w", dir, err)
		}
	}
	if err := os.MkdirAll(resolved, 0755); err != nil {
		return "", fmt.Errorf("create browser.cdp_user_data_dir %q: %w", resolved, err)
	}
	return resolved, nil
}

func resolveLocalCDPBrowserExecutable(explicit string) (string, error) {
	if strings.TrimSpace(explicit) != "" {
		resolved, err := resolveExecutableCandidate(explicit)
		if err != nil {
			return "", fmt.Errorf("resolve browser.cdp_executable %q: %w", explicit, err)
		}
		return resolved, nil
	}

	candidates := localCDPBrowserExecutableCandidates()
	searched := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		resolved, err := resolveExecutableCandidate(candidate)
		if err == nil {
			return resolved, nil
		}
		searched = append(searched, candidate)
	}
	return "", fmt.Errorf("could not find Chrome/Chromium/Edge executable automatically; install a supported browser or set browser.cdp_executable / -browser-cdp-executable. searched: %s", strings.Join(searched, ", "))
}

func resolveExecutableCandidate(candidate string) (string, error) {
	candidate = expandHomePath(strings.TrimSpace(candidate))
	if candidate == "" {
		return "", fmt.Errorf("blank executable path")
	}
	if runtime.GOOS == "darwin" && strings.HasSuffix(strings.ToLower(candidate), ".app") {
		if resolved := macAppExecutable(candidate); resolved != "" {
			return resolved, nil
		}
	}
	if !strings.Contains(candidate, string(os.PathSeparator)) {
		if resolved, err := exec.LookPath(candidate); err == nil {
			return resolved, nil
		}
	}
	if isExecutableFile(candidate) {
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			return "", err
		}
		return resolved, nil
	}
	return "", fmt.Errorf("not an executable file")
}

func localCDPBrowserExecutableCandidates() []string {
	candidates := []string{
		os.Getenv("TSPLAY_BROWSER_EXECUTABLE"),
		os.Getenv("CHROME_EXECUTABLE"),
		os.Getenv("CHROME_PATH"),
	}
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		for _, root := range []string{"/Applications", filepath.Join(home, "Applications")} {
			candidates = append(candidates,
				filepath.Join(root, "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
				filepath.Join(root, "Chromium.app", "Contents", "MacOS", "Chromium"),
				filepath.Join(root, "Microsoft Edge.app", "Contents", "MacOS", "Microsoft Edge"),
				filepath.Join(root, "Brave Browser.app", "Contents", "MacOS", "Brave Browser"),
			)
		}
		candidates = append(candidates, "google-chrome", "chrome", "chromium", "chromium-browser", "microsoft-edge", "msedge", "brave-browser")
	case "windows":
		for _, root := range []string{os.Getenv("LOCALAPPDATA"), os.Getenv("PROGRAMFILES"), os.Getenv("PROGRAMFILES(X86)")} {
			if strings.TrimSpace(root) == "" {
				continue
			}
			candidates = append(candidates,
				filepath.Join(root, "Google", "Chrome", "Application", "chrome.exe"),
				filepath.Join(root, "Microsoft", "Edge", "Application", "msedge.exe"),
				filepath.Join(root, "Chromium", "Application", "chrome.exe"),
				filepath.Join(root, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
			)
		}
		candidates = append(candidates, "chrome.exe", "msedge.exe", "chromium.exe", "brave.exe")
	default:
		candidates = append(candidates,
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
			"microsoft-edge",
			"microsoft-edge-stable",
			"msedge",
			"brave-browser",
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		)
	}
	return candidates
}

func macAppExecutable(appPath string) string {
	info, err := os.Stat(appPath)
	if err != nil || !info.IsDir() {
		return ""
	}
	name := strings.TrimSuffix(filepath.Base(appPath), ".app")
	candidate := filepath.Join(appPath, "Contents", "MacOS", name)
	if isExecutableFile(candidate) {
		return candidate
	}
	return ""
}

func isExecutableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0111 != 0
}

func expandHomePath(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			if path == "~" {
				return home
			}
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
