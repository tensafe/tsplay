package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tsplay/tsplay_core"
)

func serveWorkbenchApp(addr string, staticRoot string, artifactRoot string) error {
	handler, sourceLabel, err := newWorkbenchServer(staticRoot, artifactRoot)
	if err != nil {
		return err
	}
	baseURL := staticServerBaseURL(addr)
	fmt.Printf("Serving Workbench UI from %s\n", sourceLabel)
	fmt.Printf("Workbench listening on %s\n", baseURL)
	fmt.Printf("Workbench page: %s/demo/workbench.html\n", baseURL)
	fmt.Printf("Workbench API health: %s/api/workbench/health\n", baseURL)
	return http.ListenAndServe(addr, handler)
}

func newWorkbenchServer(staticRoot string, artifactRoot string) (http.Handler, string, error) {
	staticHandler, sourceLabel, err := newStaticFileServer(staticRoot)
	if err != nil {
		return nil, "", err
	}
	resolvedArtifactRoot, err := resolveWorkbenchArtifactRoot(artifactRoot)
	if err != nil {
		return nil, "", err
	}
	apiHandler := tsplay_core.NewWorkbenchAPIHandler(artifactRoot)
	artifactHandler := http.StripPrefix("/workbench-artifacts/", http.FileServer(http.Dir(resolvedArtifactRoot)))
	mux := http.NewServeMux()
	mux.HandleFunc("/api/workbench/app-meta", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"artifact_root":      resolvedArtifactRoot,
			"artifact_base_path": "/workbench-artifacts/",
		})
	})
	mux.Handle("/api/workbench/", apiHandler)
	mux.Handle("/workbench-artifacts/", artifactHandler)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.TrimSpace(r.URL.Path) == "" {
			http.Redirect(w, r, "/demo/workbench.html", http.StatusFound)
			return
		}
		staticHandler.ServeHTTP(w, r)
	}))
	return mux, sourceLabel, nil
}

func resolveWorkbenchArtifactRoot(root string) (string, error) {
	resolved := strings.TrimSpace(root)
	if resolved == "" {
		resolved = tsplay_core.DefaultFlowArtifactRoot
	}
	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", fmt.Errorf("resolve artifact root %q: %w", resolved, err)
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", fmt.Errorf("create artifact root %q: %w", abs, err)
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("evaluate artifact root %q: %w", abs, err)
	}
	return real, nil
}
