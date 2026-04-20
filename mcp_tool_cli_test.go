package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMCPToolArgumentsResolvesFileReferences(t *testing.T) {
	dir := t.TempDir()

	runPath := filepath.Join(dir, "run.json")
	if err := os.WriteFile(runPath, []byte(`{"ok":false,"result":{"trace":[{"index":1}]}}`), 0600); err != nil {
		t.Fatalf("write run.json: %v", err)
	}

	draftPath := filepath.Join(dir, "draft.json")
	if err := os.WriteFile(draftPath, []byte(`{"draft":{"flow_yaml":"schema_version: \"1\"\nsteps:\n  - action: set_var\n    save_as: payload\n    with:\n      value:\n        lesson: \"draft\"\n"}}`), 0600); err != nil {
		t.Fatalf("write draft.json: %v", err)
	}

	observationPath := filepath.Join(dir, "observation.json")
	if err := os.WriteFile(observationPath, []byte(`{"observation":{"url":"http://127.0.0.1:8000/demo/template_release_lab.html","title":"Template Release Lab"}}`), 0600); err != nil {
		t.Fatalf("write observation.json: %v", err)
	}

	argsPath := filepath.Join(dir, "args.json")
	argsContent := `{
  "observation": "@jsonfile:observation.json",
  "flow": "@jsonpathfile:draft.json#draft.flow_yaml",
  "run_result": "@file:run.json"
}`
	if err := os.WriteFile(argsPath, []byte(argsContent), 0600); err != nil {
		t.Fatalf("write args.json: %v", err)
	}

	arguments, err := loadMCPToolArguments("", argsPath)
	if err != nil {
		t.Fatalf("load args: %v", err)
	}

	observation, ok := arguments["observation"].(map[string]any)
	if !ok {
		t.Fatalf("expected observation object, got %#v", arguments["observation"])
	}
	wrapper, ok := observation["observation"].(map[string]any)
	if !ok || wrapper["title"] != "Template Release Lab" {
		t.Fatalf("unexpected observation payload: %#v", arguments["observation"])
	}

	flow, ok := arguments["flow"].(string)
	if !ok || flow == "" {
		t.Fatalf("expected flow YAML string, got %#v", arguments["flow"])
	}

	runResult, ok := arguments["run_result"].(string)
	if !ok || runResult == "" {
		t.Fatalf("expected raw run_result text, got %#v", arguments["run_result"])
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(runResult), &parsed); err != nil {
		t.Fatalf("run_result should still be JSON text: %v", err)
	}
}
