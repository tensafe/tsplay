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

func TestGoddddocrDemoAssetsAreBundled(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{path: "demo/captcha_login.html", want: "Captcha Login Demo"},
		{path: "script/tutorials/goddddocr_login.flow.yaml", want: "goddddocr_login_flow"},
		{path: "docs/tutorials/goddddocr-captcha-login.md", want: "goddddocr 验证码登录示例"},
	}

	for _, tc := range cases {
		content, err := readBundledAsset(tc.path)
		if err != nil {
			t.Fatalf("readBundledAsset(%q): %v", tc.path, err)
		}
		if !strings.Contains(string(content), tc.want) {
			t.Fatalf("bundled asset %q does not contain %q", tc.path, tc.want)
		}
	}

	image, err := readBundledAsset("demo/data/captcha_3n3d.png")
	if err != nil {
		t.Fatalf("readBundledAsset captcha image: %v", err)
	}
	if len(image) < 8 || string(image[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("captcha image is not a PNG")
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

	for _, path := range []string{
		filepath.Join(root, "demo", "captcha_login.html"),
		filepath.Join(root, "demo", "data", "captcha_3n3d.png"),
		filepath.Join(root, "script", "tutorials", "goddddocr_login.flow.yaml"),
		filepath.Join(root, "docs", "tutorials", "goddddocr-captcha-login.md"),
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("expected extracted goddddocr asset %s: %v", path, err)
		}
		if info.Size() == 0 {
			t.Fatalf("expected extracted goddddocr asset %s to be non-empty", path)
		}
	}
}
