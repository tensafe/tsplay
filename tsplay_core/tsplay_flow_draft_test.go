package tsplay_core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDraftFlowSearchAndExport(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://example.com/orders",
		Title:        "Orders",
		ArtifactRoot: t.TempDir(),
		Elements: []PageObservationElement{
			{
				Index:              1,
				Tag:                "input",
				Type:               "text",
				ID:                 "query",
				Label:              "Order keyword",
				Placeholder:        "Search orders",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`[data-testid="order-query"]`, `#query`},
				Attributes:         map[string]string{"data-testid": "order-query"},
			},
			{
				Index:              2,
				Tag:                "button",
				Type:               "button",
				ID:                 "search-button",
				Text:               "Search",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`#search-button`, `text="Search"`},
			},
			{
				Index:              3,
				Tag:                "a",
				Type:               "link",
				Text:               "Export orders",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`[data-cy="export-link"]`, `text="Export orders"`},
				Attributes:         map[string]string{"data-cy": "export-link"},
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "搜索订单并导出",
		Observation: observation,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if !strings.Contains(draft.FlowYAML, `action: type_text`) {
		t.Fatalf("expected type_text step: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `[data-testid="order-query"]`) {
		t.Fatalf("expected search selector: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `[data-cy="export-link"]`) {
		t.Fatalf("expected export selector: %s", draft.FlowYAML)
	}
	flow, err := ParseFlow([]byte(draft.FlowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse drafted flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate drafted flow: %v", err)
	}
	if draft.SuggestedVars["order_query"] != "TODO" {
		t.Fatalf("expected TODO order_query, got %#v", draft.SuggestedVars["order_query"])
	}
	if draft.Validation == nil || !draft.Validation.Valid {
		t.Fatalf("expected final validation to pass, got %#v", draft.Validation)
	}
}

func TestBuildDraftFlowAutoRepairsSelectors(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://example.com/orders",
		Title:        "Orders",
		ArtifactRoot: t.TempDir(),
		Elements: []PageObservationElement{
			{
				Index:              1,
				Tag:                "input",
				Type:               "text",
				ID:                 "query",
				Label:              "Order keyword",
				Placeholder:        "Search orders",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`#query`, `[data-testid="order-query"]`},
				Attributes:         map[string]string{"data-testid": "order-query"},
			},
			{
				Index:              2,
				Tag:                "button",
				Type:               "button",
				ID:                 "search-button",
				Text:               "Search",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`#search-button`, `text="Search"`},
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "搜索订单",
		Observation: observation,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if draft.InitialValidation == nil || !draft.InitialValidation.Valid {
		t.Fatalf("expected initial validation to pass, got %#v", draft.InitialValidation)
	}
	if draft.AutoRepaired {
		t.Fatalf("expected stable selector ranking to avoid auto repair, got %#v", draft)
	}
	if len(draft.SelectorRepairs) != 0 {
		t.Fatalf("expected selector repairs to be empty, got %#v", draft.SelectorRepairs)
	}
	if !strings.Contains(draft.FlowYAML, `[data-testid="order-query"]`) {
		t.Fatalf("expected stable selector in yaml: %s", draft.FlowYAML)
	}
	if strings.Contains(draft.FlowYAML, "\n    selector: \"#query\"") || strings.Contains(draft.FlowYAML, "\n    selector: '#query'") {
		t.Fatalf("expected weak selector to be avoided: %s", draft.FlowYAML)
	}
	if draft.Validation == nil || !draft.Validation.Valid {
		t.Fatalf("expected final validation to pass, got %#v", draft.Validation)
	}
}

func TestBuildDraftFlowTitleAndTableExtraction(t *testing.T) {
	artifactRoot := t.TempDir()
	domPath := filepath.Join(artifactRoot, "dom_snapshot.json")
	content := `{
  "tag": "HTML",
  "xpath": "/html",
  "text": "Order Center Table with rows",
  "children": [
    {
      "tag": "BODY",
      "xpath": "/html/body",
      "text": "Order Center Table with rows",
      "children": [
        {
          "tag": "H1",
          "xpath": "//*[@id=\"pageTitle\"]",
          "text": "Order Center",
          "children": []
        },
        {
          "tag": "TABLE",
          "xpath": "//*[@id=\"myTable\"]",
          "text": "Row 1 Row 2",
          "children": []
        }
      ]
    }
  ]
}`
	if err := os.WriteFile(domPath, []byte(content), 0600); err != nil {
		t.Fatalf("write dom snapshot: %v", err)
	}

	observation := &PageObservation{
		URL:             "https://example.com/tables",
		Title:           "Tables",
		ArtifactRoot:    artifactRoot,
		DOMSnapshotPath: domPath,
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:       "提取标题并抓取表格",
		Observation:  observation,
		ArtifactRoot: artifactRoot,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if !strings.Contains(draft.FlowYAML, `action: extract_text`) {
		t.Fatalf("expected extract_text step: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `action: capture_table`) {
		t.Fatalf("expected capture_table step: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `#myTable`) {
		t.Fatalf("expected table selector from dom snapshot: %s", draft.FlowYAML)
	}
	flow, err := ParseFlow([]byte(draft.FlowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse drafted flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate drafted flow: %v", err)
	}
}

func TestBuildDraftFlowUploadIntent(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://example.com/upload",
		Title:        "Upload",
		ArtifactRoot: t.TempDir(),
		Elements: []PageObservationElement{
			{
				Index:              1,
				Tag:                "input",
				Type:               "file",
				ID:                 "fileInput",
				Label:              "选择文件",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`#fileInput`},
			},
			{
				Index:              2,
				Tag:                "button",
				Type:               "submit",
				Text:               "上传文件",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`text="上传文件"`},
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "上传文件并提交",
		Observation: observation,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if !strings.Contains(draft.FlowYAML, `action: upload_file`) {
		t.Fatalf("expected upload_file step: %s", draft.FlowYAML)
	}
	if draft.SuggestedVars["upload_file_path"] != "TODO" {
		t.Fatalf("expected TODO upload path, got %#v", draft.SuggestedVars["upload_file_path"])
	}
	if draft.Validation == nil || !draft.Validation.Valid {
		t.Fatalf("expected structural validation to pass, got %#v", draft.Validation)
	}
}

func TestBuildDraftFlowValidationFailureProducesRepairHints(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://example.com/upload",
		Title:        "Upload",
		ArtifactRoot: t.TempDir(),
		Elements: []PageObservationElement{
			{
				Index:              1,
				Tag:                "input",
				Type:               "file",
				ID:                 "fileInput",
				Label:              "选择文件",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`#fileInput`},
			},
			{
				Index:              2,
				Tag:                "button",
				Type:               "submit",
				Text:               "上传文件",
				Visible:            true,
				Enabled:            true,
				SelectorCandidates: []string{`text="上传文件"`},
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "上传文件并提交",
		Observation: observation,
		Security:    &FlowSecurityPolicy{},
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if draft.Validation == nil || draft.Validation.Valid {
		t.Fatalf("expected validation failure, got %#v", draft.Validation)
	}
	if len(draft.RepairHints) == 0 {
		t.Fatalf("expected repair hints, got %#v", draft.RepairHints)
	}
	firstHint := draft.RepairHints[0]
	if firstHint.StepPath != "3" {
		t.Fatalf("expected first hint to point at upload_file step, got %#v", firstHint)
	}
	if !strings.Contains(firstHint.Suggestion, "allow_file_access=true") {
		t.Fatalf("expected file access hint, got %#v", firstHint)
	}
	if firstHint.Action != "upload_file" {
		t.Fatalf("expected upload_file action, got %#v", firstHint.Action)
	}
}
