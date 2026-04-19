package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"tsplay/tsplay_core"
)

//go:embed ReadMe.md demo docs script
var bundledAssets embed.FS

func loadScriptSource(path string) (string, error) {
	content, err := readAssetOrFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func loadFlowDefinition(path string) (*tsplay_core.Flow, error) {
	content, err := readAssetOrFile(path)
	if err != nil {
		return nil, err
	}

	flow, err := tsplay_core.ParseFlow(content, strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), "."))
	if err != nil {
		return nil, fmt.Errorf("parse flow %s: %w", path, err)
	}
	return flow, nil
}

func readAssetOrFile(name string) ([]byte, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if content, err := os.ReadFile(name); err == nil {
		return content, nil
	}
	return readBundledAsset(name)
}

func readBundledAsset(name string) ([]byte, error) {
	normalized, err := normalizeBundledAssetPath(name)
	if err != nil {
		return nil, err
	}
	content, err := bundledAssets.ReadFile(normalized)
	if err != nil {
		return nil, fmt.Errorf("open bundled asset %q: %w", normalized, err)
	}
	return content, nil
}

func bundledAssetExists(name string) bool {
	normalized, err := normalizeBundledAssetPath(name)
	if err != nil {
		return false
	}
	info, err := fs.Stat(bundledAssets, normalized)
	return err == nil && !info.IsDir()
}

func normalizeBundledAssetPath(name string) (string, error) {
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	normalized = filepath.ToSlash(normalized)
	normalized = strings.TrimPrefix(normalized, "./")
	normalized = strings.TrimPrefix(normalized, "/")
	normalized = path.Clean(normalized)
	if normalized == "." || normalized == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if !fs.ValidPath(normalized) {
		return "", fmt.Errorf("path %q is invalid", name)
	}
	return normalized, nil
}

func bundledAssetNames() ([]string, error) {
	names := []string{}
	err := fs.WalkDir(bundledAssets, ".", func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		names = append(names, current)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func extractBundledAssets(root string) (int, error) {
	if strings.TrimSpace(root) == "" {
		root = "tsplay-assets"
	}
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return 0, fmt.Errorf("resolve extract root %q: %w", root, err)
	}
	if err := os.MkdirAll(resolvedRoot, 0755); err != nil {
		return 0, fmt.Errorf("create extract root %q: %w", resolvedRoot, err)
	}

	count := 0
	err = fs.WalkDir(bundledAssets, ".", func(current string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		target := filepath.Join(resolvedRoot, filepath.FromSlash(current))
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		content, err := bundledAssets.ReadFile(current)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0644); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}
