package tsplay_core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/playwright-community/playwright-go"
)

const (
	playwrightBundledDirName  = "playwright"
	playwrightDriverDirName   = "driver"
	playwrightBrowsersDirName = "browsers"
	playwrightBrowserName     = "chromium"

	envPlaywrightBundlePath  = "TSPLAY_PLAYWRIGHT_BUNDLE_PATH"
	envPlaywrightDriverPath  = "PLAYWRIGHT_DRIVER_PATH"
	envPlaywrightBrowsersDir = "PLAYWRIGHT_BROWSERS_PATH"
)

type PlaywrightRuntimeInfo struct {
	Browser        string `json:"browser"`
	BundlePath     string `json:"bundle_path,omitempty"`
	BundleSource   string `json:"bundle_source,omitempty"`
	DriverPath     string `json:"driver_path,omitempty"`
	DriverSource   string `json:"driver_source"`
	BrowsersPath   string `json:"browsers_path,omitempty"`
	BrowsersSource string `json:"browsers_source"`
}

var (
	playwrightInstallFunc = func() error {
		return playwright.Install(newPlaywrightRunOptions())
	}
	playwrightRunFunc = func() (*playwright.Playwright, error) {
		return playwright.Run(newPlaywrightRunOptions())
	}

	playwrightInstallMu            sync.Mutex
	playwrightInstallDone          bool
	playwrightBundleCandidatesFunc = defaultPlaywrightBundleCandidates
)

func DescribePlaywrightRuntime() PlaywrightRuntimeInfo {
	return resolvePlaywrightRuntimeInfo()
}

func InstallPlaywrightRuntime() (PlaywrightRuntimeInfo, error) {
	info := configurePlaywrightRuntime()
	if err := EnsurePlaywrightInstalled(); err != nil {
		return info, err
	}
	return info, nil
}

func EnsurePlaywrightInstalled() error {
	playwrightInstallMu.Lock()
	defer playwrightInstallMu.Unlock()

	if playwrightInstallDone {
		return nil
	}
	if err := playwrightInstallFunc(); err != nil {
		return fmt.Errorf("could not install Playwright browsers: %w", err)
	}
	playwrightInstallDone = true
	return nil
}

func StartPlaywright() (*playwright.Playwright, error) {
	configurePlaywrightRuntime()
	if err := EnsurePlaywrightInstalled(); err != nil {
		return nil, err
	}
	pw, err := playwrightRunFunc()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	return pw, nil
}

func newPlaywrightRunOptions() *playwright.RunOptions {
	info := configurePlaywrightRuntime()
	options := &playwright.RunOptions{
		Browsers: []string{playwrightBrowserName},
	}
	if info.DriverPath != "" {
		options.DriverDirectory = info.DriverPath
	}
	return options
}

func configurePlaywrightRuntime() PlaywrightRuntimeInfo {
	info := resolvePlaywrightRuntimeInfo()
	if info.DriverPath != "" && strings.TrimSpace(os.Getenv(envPlaywrightDriverPath)) == "" {
		_ = os.Setenv(envPlaywrightDriverPath, info.DriverPath)
	}
	if info.BrowsersPath != "" && strings.TrimSpace(os.Getenv(envPlaywrightBrowsersDir)) == "" {
		_ = os.Setenv(envPlaywrightBrowsersDir, info.BrowsersPath)
	}
	return info
}

func resolvePlaywrightRuntimeInfo() PlaywrightRuntimeInfo {
	info := PlaywrightRuntimeInfo{
		Browser:        playwrightBrowserName,
		DriverSource:   "playwright-default-cache",
		BrowsersSource: "playwright-default-cache",
	}

	driverPath := strings.TrimSpace(os.Getenv(envPlaywrightDriverPath))
	browsersPath := strings.TrimSpace(os.Getenv(envPlaywrightBrowsersDir))
	if driverPath != "" {
		info.DriverPath = driverPath
		info.DriverSource = "env:" + envPlaywrightDriverPath
	}
	if browsersPath != "" {
		info.BrowsersPath = browsersPath
		info.BrowsersSource = "env:" + envPlaywrightBrowsersDir
	}
	if driverPath != "" && browsersPath != "" {
		return info
	}

	explicitBundle := strings.TrimSpace(os.Getenv(envPlaywrightBundlePath))
	if explicitBundle != "" {
		bundlePath := filepath.Clean(explicitBundle)
		info.BundlePath = bundlePath
		info.BundleSource = "env:" + envPlaywrightBundlePath
		if driverPath == "" {
			info.DriverPath = filepath.Join(bundlePath, playwrightDriverDirName)
			info.DriverSource = "bundle"
		}
		if browsersPath == "" {
			info.BrowsersPath = filepath.Join(bundlePath, playwrightBrowsersDirName)
			info.BrowsersSource = "bundle"
		}
		return info
	}

	for _, bundlePath := range playwrightBundleCandidatesFunc() {
		bundlePath = strings.TrimSpace(bundlePath)
		if bundlePath == "" {
			continue
		}
		driverCandidate := filepath.Join(bundlePath, playwrightDriverDirName)
		browsersCandidate := filepath.Join(bundlePath, playwrightBrowsersDirName)
		driverExists := dirExists(driverCandidate)
		browsersExists := dirExists(browsersCandidate)
		if !driverExists && !browsersExists {
			continue
		}
		info.BundlePath = filepath.Clean(bundlePath)
		info.BundleSource = "executable-neighbor"
		if driverPath == "" && driverExists {
			info.DriverPath = driverCandidate
			info.DriverSource = "bundle"
		}
		if browsersPath == "" && browsersExists {
			info.BrowsersPath = browsersCandidate
			info.BrowsersSource = "bundle"
		}
		return info
	}

	return info
}

func defaultPlaywrightBundleCandidates() []string {
	executablePath, err := os.Executable()
	if err != nil {
		return nil
	}
	if resolved, err := filepath.EvalSymlinks(executablePath); err == nil {
		executablePath = resolved
	}
	return []string{filepath.Join(filepath.Dir(executablePath), playwrightBundledDirName)}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
