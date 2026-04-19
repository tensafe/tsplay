package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildStaticAssetFS(t *testing.T) {
	root := t.TempDir()

	assetFS, label, err := buildStaticAssetFS(root)
	if err != nil {
		t.Fatalf("buildStaticAssetFS: %v", err)
	}
	if assetFS == nil {
		t.Fatalf("expected assetFS")
	}
	if !strings.Contains(label, root) {
		t.Fatalf("label = %q", label)
	}
}

func TestNewStaticFileServerServesFiles(t *testing.T) {
	root := t.TempDir()
	demoRoot := filepath.Join(root, "demo")
	if err := os.MkdirAll(demoRoot, 0755); err != nil {
		t.Fatalf("mkdir demo root: %v", err)
	}
	pagePath := filepath.Join(demoRoot, "demo.html")
	if err := os.WriteFile(pagePath, []byte("hello tsplay"), 0644); err != nil {
		t.Fatalf("write page: %v", err)
	}

	handler, resolvedRoot, err := newStaticFileServer(root)
	if err != nil {
		t.Fatalf("newStaticFileServer: %v", err)
	}
	if !strings.Contains(resolvedRoot, root) {
		t.Fatalf("resolvedRoot = %q, want contains %q", resolvedRoot, root)
	}

	req := httptest.NewRequest("GET", "/demo/demo.html", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); body != "hello tsplay" {
		t.Fatalf("body = %q", body)
	}
}

func TestNewStaticFileServerServesBundledAssets(t *testing.T) {
	handler, sourceLabel, err := newStaticFileServer("")
	if err != nil {
		t.Fatalf("newStaticFileServer: %v", err)
	}
	if !strings.Contains(sourceLabel, "bundled assets") {
		t.Fatalf("sourceLabel = %q", sourceLabel)
	}

	req := httptest.NewRequest("GET", "/demo/demo.html", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "选项选择示例") {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestStaticServerBaseURL(t *testing.T) {
	cases := []struct {
		addr string
		want string
	}{
		{addr: ":8000", want: "http://127.0.0.1:8000"},
		{addr: "0.0.0.0:9000", want: "http://127.0.0.1:9000"},
		{addr: "127.0.0.1:7000", want: "http://127.0.0.1:7000"},
	}

	for _, tc := range cases {
		if got := staticServerBaseURL(tc.addr); got != tc.want {
			t.Fatalf("staticServerBaseURL(%q) = %q, want %q", tc.addr, got, tc.want)
		}
	}
}
