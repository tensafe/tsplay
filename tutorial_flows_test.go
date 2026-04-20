package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAllTutorialFlows(t *testing.T) {
	var paths []string
	err := filepath.Walk("script/tutorials", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".flow.yaml") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk tutorial flows: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("expected tutorial flow files")
	}

	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()
			if _, err := loadFlowDefinition(path); err != nil {
				t.Fatalf("loadFlowDefinition(%q): %v", path, err)
			}
		})
	}
}
