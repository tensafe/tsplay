package tsplay_core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
	if observation.ScreenshotPath == "" || observation.DOMSnapshotPath == "" {
		t.Fatalf("expected artifact paths: %#v", observation)
	}
	if _, err := os.Stat(observation.ScreenshotPath); err != nil {
		t.Fatalf("expected screenshot artifact: %v", err)
	}
	if _, err := os.Stat(observation.DOMSnapshotPath); err != nil {
		t.Fatalf("expected dom snapshot artifact: %v", err)
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
