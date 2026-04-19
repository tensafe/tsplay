package main

import (
	"path/filepath"
	"testing"
)

func TestLoadAllTutorialFlows(t *testing.T) {
	paths, err := filepath.Glob("script/tutorials/*.flow.yaml")
	if err != nil {
		t.Fatalf("glob tutorial flows: %v", err)
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
