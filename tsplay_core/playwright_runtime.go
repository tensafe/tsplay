package tsplay_core

import (
	"fmt"
	"sync"

	"github.com/playwright-community/playwright-go"
)

var (
	playwrightInstallFunc = func() error {
		return playwright.Install()
	}
	playwrightRunFunc = func() (*playwright.Playwright, error) {
		return playwright.Run()
	}

	playwrightInstallMu   sync.Mutex
	playwrightInstallDone bool
)

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
	if err := EnsurePlaywrightInstalled(); err != nil {
		return nil, err
	}
	pw, err := playwrightRunFunc()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	return pw, nil
}
