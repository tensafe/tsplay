package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWorkbenchServerServesUIAndAPI(t *testing.T) {
	handler, sourceLabel, err := newWorkbenchServer("", t.TempDir())
	if err != nil {
		t.Fatalf("newWorkbenchServer: %v", err)
	}
	if !strings.Contains(sourceLabel, "bundled assets") {
		t.Fatalf("sourceLabel = %q", sourceLabel)
	}

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 302 {
		t.Fatalf("root status = %d, want 302", rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/demo/workbench.html" {
		t.Fatalf("root location = %q", location)
	}

	req = httptest.NewRequest("GET", "/demo/workbench.html", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("ui status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "TSPlay Workbench") {
		t.Fatalf("unexpected ui body: %q", body)
	}

	req = httptest.NewRequest("GET", "/api/workbench/health", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("health status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "\"ok\": true") {
		t.Fatalf("unexpected health body: %q", body)
	}

	req = httptest.NewRequest("GET", "/api/workbench/app-meta", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("app-meta status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "\"artifact_base_path\"") {
		t.Fatalf("unexpected app-meta body: %q", body)
	}
}

func TestNewWorkbenchServerServesArtifacts(t *testing.T) {
	artifactRoot := t.TempDir()
	samplePath := filepath.Join(artifactRoot, "runs", "sample.txt")
	if err := os.MkdirAll(filepath.Dir(samplePath), 0755); err != nil {
		t.Fatalf("mkdir artifact dir: %v", err)
	}
	if err := os.WriteFile(samplePath, []byte("artifact body"), 0644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}

	handler, _, err := newWorkbenchServer("", artifactRoot)
	if err != nil {
		t.Fatalf("newWorkbenchServer: %v", err)
	}
	req := httptest.NewRequest("GET", "/workbench-artifacts/runs/sample.txt", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("artifact status = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); body != "artifact body" {
		t.Fatalf("unexpected artifact body: %q", body)
	}
}
