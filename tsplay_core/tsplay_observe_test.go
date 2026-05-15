package tsplay_core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/playwright-community/playwright-go"
)

func TestObservePageCapturesInteractiveElements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html>
<head><title>Orders</title></head>
<body>
  <main>
    <h1>Order Center</h1>
    <label for="query">Order keyword</label>
    <input id="query" name="query" data-testid="order-query" placeholder="Search orders">
    <button id="search-button">Search</button>
    <a href="/export" data-cy="export-link">Export orders</a>
  </main>
</body>
</html>`)
	}))
	defer server.Close()

	observation, err := ObservePage(PageObservationOptions{
		URL:          server.URL,
		Headless:     true,
		ArtifactRoot: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("observe page: %v", err)
	}
	if observation.Title != "Orders" {
		t.Fatalf("title = %q", observation.Title)
	}
	if strings.TrimSpace(observation.PageSummary) == "" {
		t.Fatalf("expected page summary, got %#v", observation.PageSummary)
	}
	if observation.ScreenshotPath == "" || observation.DOMSnapshotPath == "" {
		t.Fatalf("expected artifact paths: %#v", observation)
	}
	if strings.TrimSpace(observation.DOMSnapshotExcerpt) == "" || !strings.Contains(observation.DOMSnapshotExcerpt, "Order Center") {
		t.Fatalf("expected dom snapshot excerpt, got %#v", observation.DOMSnapshotExcerpt)
	}
	if _, err := os.Stat(observation.ScreenshotPath); err != nil {
		t.Fatalf("expected screenshot artifact: %v", err)
	}
	if _, err := os.Stat(observation.DOMSnapshotPath); err != nil {
		t.Fatalf("expected dom snapshot artifact: %v", err)
	}
	if len(observation.ContentElements) == 0 {
		t.Fatalf("expected content elements, got %#v", observation.ContentElements)
	}
	headline := findObservedContentElement(observation.ContentElements, func(element PageObservationContentElement) bool {
		return element.Kind == "headline" && strings.Contains(element.Text, "Order Center")
	})
	if headline == nil {
		t.Fatalf("expected headline content element, got %#v", observation.ContentElements)
	}
	linkContent := findObservedContentElement(observation.ContentElements, func(element PageObservationContentElement) bool {
		return element.Kind == "article_link" && strings.Contains(element.Text, "Export orders")
	})
	if linkContent == nil || linkContent.Selector == "" {
		t.Fatalf("expected link content element with selector, got %#v", observation.ContentElements)
	}

	input := findObservedElement(observation.Elements, func(element PageObservationElement) bool {
		return element.ID == "query"
	})
	if input == nil {
		t.Fatalf("query input not found: %#v", observation.Elements)
	}
	if input.Label != "Order keyword" {
		t.Fatalf("input label = %q", input.Label)
	}
	if input.PrimarySelector != `[data-testid="order-query"]` {
		t.Fatalf("input primary selector = %q", input.PrimarySelector)
	}
	if input.SelectorCandidates[0] != input.PrimarySelector {
		t.Fatalf("expected primary selector first, got %#v", input.SelectorCandidates)
	}
	if input.SelectorRationale == "" {
		t.Fatalf("expected selector rationale, got %#v", input)
	}
	if !containsString(input.SelectorCandidates, `[data-testid="order-query"]`) {
		t.Fatalf("missing data-testid selector: %#v", input.SelectorCandidates)
	}
	if !containsString(input.SelectorCandidates, `input[placeholder="Search orders"]`) {
		t.Fatalf("missing placeholder selector: %#v", input.SelectorCandidates)
	}

	button := findObservedElement(observation.Elements, func(element PageObservationElement) bool {
		return strings.Contains(element.Text, "Search")
	})
	if button == nil {
		t.Fatalf("search button not found: %#v", observation.Elements)
	}
	if button.PrimarySelector != `button#search-button` && button.PrimarySelector != `#search-button` {
		t.Fatalf("unexpected button primary selector: %#v", button)
	}
	if !containsString(button.SelectorCandidates, `text="Search"`) {
		t.Fatalf("missing text selector: %#v", button.SelectorCandidates)
	}

	link := findObservedElement(observation.Elements, func(element PageObservationElement) bool {
		return strings.Contains(element.Text, "Export orders")
	})
	if link == nil {
		t.Fatalf("export link not found: %#v", observation.Elements)
	}
	if !containsString(link.SelectorCandidates, `[data-cy="export-link"]`) {
		t.Fatalf("missing data-cy selector: %#v", link.SelectorCandidates)
	}
}

func findObservedElement(elements []PageObservationElement, match func(PageObservationElement) bool) *PageObservationElement {
	for i := range elements {
		if match(elements[i]) {
			return &elements[i]
		}
	}
	return nil
}

func TestObservePageWithCDPLaunch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html><html><head><title>CDP Observe</title></head><body><button id="ready">Observed over CDP</button></body></html>`)
	}))
	defer server.Close()

	pw, err := StartPlaywright()
	if err != nil {
		t.Fatalf("start playwright: %v", err)
	}
	executable := pw.Chromium.ExecutablePath()
	if err := pw.Stop(); err != nil {
		t.Fatalf("stop playwright: %v", err)
	}
	if executable == "" {
		t.Skip("playwright chromium executable path is empty")
	}

	profileRoot, err := prepareRuntimeFileRoot(t.TempDir())
	if err != nil {
		t.Fatalf("prepare profile root: %v", err)
	}
	observation, err := ObservePage(PageObservationOptions{
		URL:            server.URL,
		Headless:       true,
		CDPLaunch:      true,
		CDPExecutable:  executable,
		CDPUserDataDir: filepath.Join(profileRoot, "observe-cdp-profile"),
		ArtifactRoot:   t.TempDir(),
		TimeoutMS:      15000,
	})
	if err != nil {
		t.Fatalf("observe page over CDP launch: %v", err)
	}
	if observation.Title != "CDP Observe" {
		t.Fatalf("title = %q", observation.Title)
	}
	foundButton := false
	for _, element := range observation.Elements {
		if strings.Contains(element.Text, "Observed over CDP") {
			foundButton = true
			break
		}
	}
	if !foundButton {
		t.Fatalf("expected observed CDP button, got %#v", observation.Elements)
	}
}

func TestObservePageRejectsInvalidCDPEndpointBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, nil)
	defer restore()

	_, err := ObservePage(PageObservationOptions{
		URL:         "https://example.com",
		CDPEndpoint: "127.0.0.1:70000/json/version",
	})
	if err == nil {
		t.Fatalf("expected invalid CDP endpoint error")
	}
	if !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled {
		t.Fatalf("invalid CDP endpoint should be rejected before starting Playwright")
	}
}

func TestObservePageRejectsRemoteCDPLaunchBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, func() (*playwright.Playwright, error) {
		runCalled = true
		return nil, fmt.Errorf("unexpected playwright startup")
	})
	defer restore()

	_, err := ObservePage(PageObservationOptions{
		URL:         "https://example.com",
		CDPLaunch:   true,
		CDPEndpoint: "http://192.0.2.1:9222",
		Security:    &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected remote CDP launch endpoint error")
	}
	if !strings.Contains(err.Error(), "only start or reuse a local browser") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("remote CDP launch endpoint should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func TestObservePageRejectsLocalCDPLaunchEndpointWithoutPortBeforePlaywrightStart(t *testing.T) {
	installCalled := false
	runCalled := false
	restore := stubPlaywrightRuntime(t, func() error {
		installCalled = true
		return fmt.Errorf("unexpected playwright install")
	}, func() (*playwright.Playwright, error) {
		runCalled = true
		return nil, fmt.Errorf("unexpected playwright startup")
	})
	defer restore()

	_, err := ObservePage(PageObservationOptions{
		URL:         "https://example.com",
		CDPLaunch:   true,
		CDPEndpoint: "http://127.0.0.1",
		Security:    &FlowSecurityPolicy{AllowBrowserState: true},
	})
	if err == nil {
		t.Fatalf("expected local CDP launch endpoint without port error")
	}
	if !strings.Contains(err.Error(), "explicit port") {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled || runCalled {
		t.Fatalf("local CDP launch endpoint without port should be rejected before Playwright starts, install=%v run=%v", installCalled, runCalled)
	}
}

func findObservedContentElement(elements []PageObservationContentElement, match func(PageObservationContentElement) bool) *PageObservationContentElement {
	for i := range elements {
		if match(elements[i]) {
			return &elements[i]
		}
	}
	return nil
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
