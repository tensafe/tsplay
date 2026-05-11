package tsplay_core

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestRunFlowZipCompressAndExtractWithPassword(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	nestedDir := filepath.Join(inputDir, "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("create input dirs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "alpha.txt"), []byte("alpha"), 0644); err != nil {
		t.Fatalf("write alpha: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "beta.txt"), []byte("beta"), 0644); err != nil {
		t.Fatalf("write beta: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "zip_round_trip",
		Steps: []FlowStep{
			{
				Action:   "zip_compress",
				FilePath: "bundle.zip",
				Password: "secret",
				SaveAs:   "archive",
				Files:    []string{"input/alpha.txt"},
				Folders:  []string{"input/nested"},
			},
			{
				Action:   "zip_extract",
				FilePath: "{{archive.file_path}}",
				SavePath: "out",
				Password: "secret",
				SaveAs:   "extracted",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	archive, ok := result.Vars["archive"].(map[string]any)
	if !ok {
		t.Fatalf("archive = %#v", result.Vars["archive"])
	}
	if archive["encrypted"] != true {
		t.Fatalf("expected encrypted archive metadata, got %#v", archive)
	}
	if got, err := os.ReadFile(filepath.Join(root, "out", "alpha.txt")); err != nil || string(got) != "alpha" {
		t.Fatalf("alpha output = %q, %v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(root, "out", "nested", "beta.txt")); err != nil || string(got) != "beta" {
		t.Fatalf("beta output = %q, %v", got, err)
	}
}

func TestValidateFlowSecurityRejectsZipWithoutAllow(t *testing.T) {
	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "zip_policy",
		Steps: []FlowStep{
			{Action: "zip_compress", FilePath: "bundle.zip", Files: []string{"input/a.txt"}},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, FlowSecurityPolicy{
		FileInputRoot:  root,
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected file access security policy error")
	}
	if !strings.Contains(err.Error(), "allow_file_access") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestZipExtractRejectsUnsafeArchivePath(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "unsafe.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	writer := zip.NewWriter(file)
	entry, err := writer.Create("../escape.txt")
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if _, err := entry.Write([]byte("escape")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}

	_, err = executeZipExtract(zipExtractConfig{
		FilePath: archivePath,
		SavePath: filepath.Join(root, "out"),
	})
	if err == nil {
		t.Fatalf("expected unsafe archive path error")
	}
	if !strings.Contains(err.Error(), "not a relative path") && !strings.Contains(err.Error(), "unsafe") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestZipExtractRejectsSymlinkAncestor(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "symlink.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	writer := zip.NewWriter(file)
	entry, err := writer.Create("link/pwn.txt")
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if _, err := entry.Write([]byte("pwn")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}

	outputRoot := filepath.Join(root, "out")
	if err := os.MkdirAll(outputRoot, 0755); err != nil {
		t.Fatalf("create output root: %v", err)
	}
	external := filepath.Join(root, "external")
	if err := os.MkdirAll(external, 0755); err != nil {
		t.Fatalf("create external: %v", err)
	}
	if err := os.Symlink(external, filepath.Join(outputRoot, "link")); err != nil {
		t.Skipf("symlink not available: %v", err)
	}

	_, err = executeZipExtract(zipExtractConfig{
		FilePath: archivePath,
		SavePath: outputRoot,
	})
	if err == nil {
		t.Fatalf("expected symlink ancestor error")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}
