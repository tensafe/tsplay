package tsplay_core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseObservationForDraftAllowsContentOnlyObservation(t *testing.T) {
	observation, err := ParseObservationForDraft(`{
  "title": "网易财经",
  "content_elements": [
    {
      "index": 1,
      "kind": "headline",
      "tag": "h2",
      "text": "财经要闻",
      "selector": "xpath=/html/body/main/section[1]/h2[1]"
    }
  ]
}`)
	if err != nil {
		t.Fatalf("parse observation: %v", err)
	}
	if observation == nil || observation.Title != "网易财经" {
		t.Fatalf("expected parsed observation title, got %#v", observation)
	}
	if len(observation.ContentElements) != 1 {
		t.Fatalf("expected content elements, got %#v", observation.ContentElements)
	}
}

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

func TestBuildDraftFlowUsesContentElementsForNewsList(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://money.163.com/",
		Title:        "网易财经",
		ArtifactRoot: t.TempDir(),
		ContentElements: []PageObservationContentElement{
			{
				Index:    1,
				Kind:     "headline",
				Tag:      "h2",
				Text:     "财经要闻",
				XPath:    "/html/body/main/section[1]/h2[1]",
				Selector: "xpath=/html/body/main/section[1]/h2[1]",
			},
			{
				Index:    2,
				Kind:     "article_link",
				Tag:      "a",
				Text:     "头条新闻一",
				Href:     "https://money.163.com/story/1",
				XPath:    "/html/body/main/section[1]/ul[1]/li[1]/a[1]",
				Selector: "xpath=/html/body/main/section[1]/ul[1]/li[1]/a[1]",
			},
			{
				Index:    3,
				Kind:     "article_link",
				Tag:      "a",
				Text:     "头条新闻二",
				Href:     "https://money.163.com/story/2",
				XPath:    "/html/body/main/section[1]/ul[1]/li[2]/a[1]",
				Selector: "xpath=/html/body/main/section[1]/ul[1]/li[2]/a[1]",
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "查看财经要闻内容",
		Observation: observation,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if !strings.Contains(draft.FlowYAML, `action: get_all_links`) {
		t.Fatalf("expected get_all_links step from content elements: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `save_as: content_links`) {
		t.Fatalf("expected content_links variable: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `save_as: content_section_title`) {
		t.Fatalf("expected content section title extraction: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `xpath=/html/body/main/section[1]`) {
		t.Fatalf("expected derived content container selector: %s", draft.FlowYAML)
	}
	if !draftHasAnyPlannedAction(draft.PlannedActions, "view_content") {
		t.Fatalf("expected view_content planned action, got %#v", draft.PlannedActions)
	}
	flow, err := ParseFlow([]byte(draft.FlowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse drafted flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate drafted flow: %v", err)
	}
}

func TestBuildDraftFlowUsesContentElementsForArticlePage(t *testing.T) {
	observation := &PageObservation{
		URL:          "https://money.163.com/story/1",
		Title:        "文章详情",
		ArtifactRoot: t.TempDir(),
		ContentElements: []PageObservationContentElement{
			{
				Index:    1,
				Kind:     "headline",
				Tag:      "h1",
				Text:     "A股午盘观察",
				XPath:    "/html/body/main/article[1]/h1[1]",
				Selector: "xpath=/html/body/main/article[1]/h1[1]",
			},
			{
				Index:    2,
				Kind:     "summary_text",
				Tag:      "p",
				Text:     "今日市场围绕科技和消费板块展开，成交额较上一交易日明显放大。",
				XPath:    "/html/body/main/article[1]/p[1]",
				Selector: "xpath=/html/body/main/article[1]/p[1]",
			},
		},
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:      "阅读文章内容",
		Observation: observation,
	})
	if err != nil {
		t.Fatalf("build draft flow: %v", err)
	}
	if !strings.Contains(draft.FlowYAML, `save_as: article_title`) {
		t.Fatalf("expected article_title extraction: %s", draft.FlowYAML)
	}
	if !strings.Contains(draft.FlowYAML, `save_as: article_summary`) {
		t.Fatalf("expected article_summary extraction: %s", draft.FlowYAML)
	}
	flow, err := ParseFlow([]byte(draft.FlowYAML), "yaml")
	if err != nil {
		t.Fatalf("parse drafted flow: %v", err)
	}
	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate drafted flow: %v", err)
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
