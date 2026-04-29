package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadBundledAsset(t *testing.T) {
	content, err := readBundledAsset("script/tutorials/01_hello_world.lua")
	if err != nil {
		t.Fatalf("readBundledAsset: %v", err)
	}
	if !strings.Contains(string(content), "lesson 01: hello world from lua") {
		t.Fatalf("unexpected bundled content: %s", string(content))
	}
}

func TestLoadBundledFlowDefinition(t *testing.T) {
	flow, err := loadFlowDefinition("script/tutorials/01_hello_world.flow.yaml")
	if err != nil {
		t.Fatalf("loadFlowDefinition: %v", err)
	}
	if flow.Name != "lesson_01_hello_world_flow" {
		t.Fatalf("flow.Name = %q", flow.Name)
	}
}

func TestExtractBundledAssets(t *testing.T) {
	root := filepath.Join(t.TempDir(), "bundle")
	count, err := extractBundledAssets(root)
	if err != nil {
		t.Fatalf("extractBundledAssets: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected extracted files")
	}

	content, err := os.ReadFile(filepath.Join(root, "docs", "tutorials", "README.md"))
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if !strings.Contains(string(content), "TSPlay Step-by-Step Tutorials") {
		t.Fatalf("unexpected extracted content")
	}
}
