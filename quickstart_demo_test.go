package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunQuickstartDemoAction(t *testing.T) {
	root := t.TempDir()

	result, err := runQuickstartDemoAction(root, true)
	if err != nil {
		t.Fatalf("runQuickstartDemoAction: %v", err)
	}
	if result == nil {
		t.Fatalf("expected result")
	}
	if result.Action != "quickstart-demo" {
		t.Fatalf("result.Action = %q", result.Action)
	}
	if result.PlaywrightRequired {
		t.Fatalf("expected quickstart demo to skip Playwright")
	}
	if result.RunResult == nil {
		t.Fatalf("expected run result")
	}
	if result.RunResult.Playwright != nil {
		t.Fatalf("expected no Playwright usage for data-only quickstart: %#v", result.RunResult.Playwright)
	}

	flowPath := filepath.Join(root, quickstartDemoDirName, quickstartDemoFlowName)
	outputPath := filepath.Join(root, quickstartDemoDirName, quickstartDemoJSONName)
	if result.GeneratedFlowPath != flowPath {
		t.Fatalf("GeneratedFlowPath = %q, want %q", result.GeneratedFlowPath, flowPath)
	}
	if result.OutputJSONPath != outputPath {
		t.Fatalf("OutputJSONPath = %q, want %q", result.OutputJSONPath, outputPath)
	}

	flowContent, err := os.ReadFile(flowPath)
	if err != nil {
		t.Fatalf("read generated flow: %v", err)
	}
	flowText := string(flowContent)
	for _, want := range []string{"name: quickstart_demo", "action: set_var", "action: write_json"} {
		if !strings.Contains(flowText, want) {
			t.Fatalf("generated flow missing %q:\n%s", want, flowText)
		}
	}

	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read quickstart output: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(outputContent, &payload); err != nil {
		t.Fatalf("unmarshal quickstart output: %v", err)
	}
	if payload["greeting"] != "hello from tsplay quickstart" {
		t.Fatalf("payload[greeting] = %#v", payload["greeting"])
	}
	if payload["playwright_required"] != false {
		t.Fatalf("payload[playwright_required] = %#v", payload["playwright_required"])
	}
	if payload["demo"] != "quickstart" {
		t.Fatalf("payload[demo] = %#v", payload["demo"])
	}

	if len(result.NextSteps) != 3 {
		t.Fatalf("len(NextSteps) = %d", len(result.NextSteps))
	}
	if !strings.Contains(result.NextSteps[0], "file-srv") {
		t.Fatalf("unexpected first next step: %q", result.NextSteps[0])
	}
	if !strings.Contains(result.NextSteps[1], "10_assert_page_state.flow.yaml") {
		t.Fatalf("unexpected browser next step: %q", result.NextSteps[1])
	}
}
