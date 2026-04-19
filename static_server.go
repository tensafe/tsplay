package main

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func serveStaticFiles(addr string, root string) error {
	handler, sourceLabel, err := newStaticFileServer(root)
	if err != nil {
		return err
	}

	baseURL := staticServerBaseURL(addr)
	fmt.Printf("Serving static files from %s\n", sourceLabel)
	fmt.Printf("Static server listening on %s\n", baseURL)
	if bundledAssetExists("demo/demo.html") {
		fmt.Printf("Demo page: %s/demo/demo.html\n", baseURL)
		fmt.Printf("Table page: %s/demo/tables.html\n", baseURL)
		fmt.Printf("Extract page: %s/demo/extract.html\n", baseURL)
		fmt.Printf("JSON data: %s/demo/data/order_summary.json\n", baseURL)
	}
	return http.ListenAndServe(addr, handler)
}

func newStaticFileServer(root string) (http.Handler, string, error) {
	assetFS, sourceLabel, err := buildStaticAssetFS(root)
	if err != nil {
		return nil, "", err
	}
	return http.FileServer(http.FS(assetFS)), sourceLabel, nil
}

func buildStaticAssetFS(root string) (fs.FS, string, error) {
	if strings.TrimSpace(root) == "" {
		return bundledAssets, "bundled assets in tsplay binary", nil
	}
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, "", fmt.Errorf("resolve serve root %q: %w", root, err)
	}
	info, err := os.Stat(resolvedRoot)
	if err != nil {
		return nil, "", fmt.Errorf("open serve root %q: %w", resolvedRoot, err)
	}
	if !info.IsDir() {
		return nil, "", fmt.Errorf("serve root %q must be a directory", resolvedRoot)
	}
	return overlayFS{root: resolvedRoot}, fmt.Sprintf("local root %s with bundled asset fallback", resolvedRoot), nil
}

func staticServerBaseURL(addr string) string {
	trimmed := strings.TrimSpace(addr)
	if trimmed == "" {
		return "http://127.0.0.1:8082"
	}
	if strings.HasPrefix(trimmed, ":") {
		return "http://127.0.0.1" + trimmed
	}
	host, port, err := net.SplitHostPort(trimmed)
	if err != nil {
		return "http://" + trimmed
	}
	switch strings.TrimSpace(host) {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port)
}

type overlayFS struct {
	root string
}

func (ofs overlayFS) Open(name string) (fs.File, error) {
	normalized := normalizeStaticOpenPath(name)
	if ofs.root != "" {
		candidate := filepath.Join(ofs.root, filepath.FromSlash(normalized))
		if file, err := os.Open(candidate); err == nil {
			return file, nil
		}
	}
	return bundledAssets.Open(normalized)
}

func normalizeStaticOpenPath(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "."
	}
	trimmed = strings.TrimPrefix(trimmed, "/")
	trimmed = filepath.ToSlash(trimmed)
	trimmed = strings.TrimPrefix(trimmed, "./")
	trimmed = path.Clean(trimmed)
	if trimmed == "" {
		return "."
	}
	return trimmed
}
