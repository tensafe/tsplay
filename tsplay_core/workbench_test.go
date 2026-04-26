package tsplay_core

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWorkbenchStoreRoundTrip(t *testing.T) {
	artifactRoot := t.TempDir()

	savedSite, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:      "demo_admin",
		Name:        "Demo Admin",
		StartURL:    "https://example.com/admin",
		SessionName: "demo_admin_user",
	}, artifactRoot)
	if err != nil {
		t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
	}
	if savedSite.SiteID != "demo_admin" {
		t.Fatalf("unexpected site id: %q", savedSite.SiteID)
	}
	if len(savedSite.AllowedDomains) != 1 || savedSite.AllowedDomains[0] != "example.com" {
		t.Fatalf("unexpected allowed domains: %#v", savedSite.AllowedDomains)
	}

	loadedSite, err := LoadWorkbenchSiteConfig("demo_admin", artifactRoot)
	if err != nil {
		t.Fatalf("LoadWorkbenchSiteConfig() error = %v", err)
	}
	if loadedSite.SessionName != "demo_admin_user" {
		t.Fatalf("unexpected session name: %q", loadedSite.SessionName)
	}

	pages, err := UpsertWorkbenchPageCards("demo_admin", artifactRoot, []WorkbenchPageCard{
		{
			ID:              "route:demo_admin:/orders",
			SiteID:          "demo_admin",
			URL:             "https://example.com/admin/orders",
			NormalizedRoute: "/orders",
			Title:           "订单管理",
			Actions: []WorkbenchActionCard{
				{Label: "搜索", Kind: "button", Risk: "read"},
			},
		},
	})
	if err != nil {
		t.Fatalf("UpsertWorkbenchPageCards() error = %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("unexpected page count: %d", len(pages))
	}

	apis, err := UpsertWorkbenchAPICards("demo_admin", artifactRoot, []WorkbenchAPICard{
		{
			ID:            "api:POST:/api/orders/search",
			SiteID:        "demo_admin",
			Method:        "POST",
			PathTemplate:  "/api/orders/search",
			SemanticName:  "订单搜索",
			TriggerRoute:  "/orders",
			OperationType: "read",
			Risk:          "read",
			ResponseSchema: map[string]any{
				"items": "Order[]",
				"total": "number",
			},
		},
	})
	if err != nil {
		t.Fatalf("UpsertWorkbenchAPICards(first) error = %v", err)
	}
	if len(apis) != 1 {
		t.Fatalf("unexpected api count after first upsert: %d", len(apis))
	}

	apis, err = UpsertWorkbenchAPICards("demo_admin", artifactRoot, []WorkbenchAPICard{
		{
			ID:           "api:POST:/api/orders/search",
			SiteID:       "demo_admin",
			Method:       "POST",
			PathTemplate: "/api/orders/search",
			RequestSchema: map[string]any{
				"status": "string",
				"page":   "number",
			},
			ContentType: "application/json",
			Status:      200,
		},
	})
	if err != nil {
		t.Fatalf("UpsertWorkbenchAPICards(second) error = %v", err)
	}
	if len(apis) != 1 {
		t.Fatalf("unexpected api count after merge: %d", len(apis))
	}
	if apis[0].SemanticName != "订单搜索" {
		t.Fatalf("semantic name was not preserved: %#v", apis[0])
	}
	if apis[0].RequestSchema == nil || apis[0].ResponseSchema == nil {
		t.Fatalf("api schema merge failed: %#v", apis[0])
	}

	entities, err := UpsertWorkbenchEntityCards("demo_admin", artifactRoot, []WorkbenchEntityCard{
		{
			ID:     "entity:Order",
			SiteID: "demo_admin",
			Name:   "Order",
			Label:  "订单",
			Fields: []WorkbenchEntityField{
				{Name: "orderId", Label: "订单号", Type: "string"},
			},
		},
	})
	if err != nil {
		t.Fatalf("UpsertWorkbenchEntityCards() error = %v", err)
	}
	if len(entities) != 1 || entities[0].Name != "Order" {
		t.Fatalf("unexpected entities: %#v", entities)
	}

	sites, err := ListWorkbenchSiteConfigs(artifactRoot)
	if err != nil {
		t.Fatalf("ListWorkbenchSiteConfigs() error = %v", err)
	}
	if len(sites) != 1 {
		t.Fatalf("unexpected site count: %d", len(sites))
	}

	savedProvider, err := SaveWorkbenchProviderConfig(WorkbenchProviderConfig{
		ProviderID: "codex_main",
		Name:       "Codex Main",
		Type:       WorkbenchProviderTypeOpenAICompatible,
		BaseURL:    "https://api.openai.com/v1",
		Model:      "gpt-4.1-mini",
		APIKey:     "sk-test-12345678",
		Enabled:    true,
	}, artifactRoot)
	if err != nil {
		t.Fatalf("SaveWorkbenchProviderConfig() error = %v", err)
	}
	if savedProvider.ProviderID != "codex_main" {
		t.Fatalf("unexpected provider id: %q", savedProvider.ProviderID)
	}

	loadedProvider, err := LoadWorkbenchProviderConfig("codex_main", artifactRoot)
	if err != nil {
		t.Fatalf("LoadWorkbenchProviderConfig() error = %v", err)
	}
	if loadedProvider.Type != WorkbenchProviderTypeOpenAICompatible {
		t.Fatalf("unexpected provider type: %q", loadedProvider.Type)
	}
	if loadedProvider.APIKey == "" {
		t.Fatalf("expected provider api key to be persisted")
	}

	providers, err := ListWorkbenchProviderConfigs(artifactRoot)
	if err != nil {
		t.Fatalf("ListWorkbenchProviderConfigs() error = %v", err)
	}
	if len(providers) != 1 {
		t.Fatalf("unexpected provider count: %d", len(providers))
	}
}

func TestResolveWorkbenchExploreSiteDerivesIDFromURL(t *testing.T) {
	artifactRoot := t.TempDir()

	site, err := resolveWorkbenchExploreSite(WorkbenchExploreOptions{
		Name:         "",
		StartURL:     "https://demo.example.com/admin",
		SessionName:  "demo_admin_user",
		ArtifactRoot: artifactRoot,
	})
	if err != nil {
		t.Fatalf("resolveWorkbenchExploreSite() error = %v", err)
	}
	if site.SiteID == "" {
		t.Fatalf("expected derived site id, got empty")
	}
	if site.StartURL != "https://demo.example.com/admin" {
		t.Fatalf("unexpected start url: %q", site.StartURL)
	}
	if len(site.AllowedDomains) == 0 || site.AllowedDomains[0] != "demo.example.com" {
		t.Fatalf("unexpected allowed domains: %#v", site.AllowedDomains)
	}
}

func TestBuildWorkbenchTaskPlanFallbacks(t *testing.T) {
	t.Run("page fallback", func(t *testing.T) {
		artifactRoot := t.TempDir()
		_, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
			SiteID:      "demo_admin",
			Name:        "Demo Admin",
			StartURL:    "https://example.com/admin",
			SessionName: "demo_admin_user",
		}, artifactRoot)
		if err != nil {
			t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
		}
		_, err = UpsertWorkbenchPageCards("demo_admin", artifactRoot, []WorkbenchPageCard{
			{
				ID:              "route:demo_admin:/orders",
				SiteID:          "demo_admin",
				URL:             "https://example.com/admin/orders",
				NormalizedRoute: "/orders",
				Title:           "订单管理",
				Summary:         "用于订单查询和导出",
				Tables: []WorkbenchTableCard{
					{Name: "订单列表", Selector: ".orders-table", Columns: []string{"订单号", "状态"}},
				},
			},
		})
		if err != nil {
			t.Fatalf("UpsertWorkbenchPageCards() error = %v", err)
		}

		plan, err := BuildWorkbenchTaskPlan(WorkbenchTaskPlanOptions{
			SiteID:       "demo_admin",
			ArtifactRoot: artifactRoot,
			Intent:       "帮我分析订单管理页面的数据",
		})
		if err != nil {
			t.Fatalf("BuildWorkbenchTaskPlan() error = %v", err)
		}
		if plan.Strategy != "ui_first" {
			t.Fatalf("unexpected strategy: %q", plan.Strategy)
		}
		if plan.Flow == nil {
			t.Fatalf("expected flow")
		}
		if plan.Flow.Browser == nil || plan.Flow.Browser.UseSession != "demo_admin_user" {
			t.Fatalf("expected flow to reuse session: %#v", plan.Flow.Browser)
		}
		if len(plan.Flow.Steps) < 1 || plan.Flow.Steps[0].Action != "navigate" {
			t.Fatalf("unexpected flow steps: %#v", plan.Flow.Steps)
		}
		if plan.FlowYAML == "" {
			t.Fatalf("expected flow yaml")
		}
	})

	t.Run("api fallback", func(t *testing.T) {
		artifactRoot := t.TempDir()
		_, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
			SiteID:      "demo_admin",
			Name:        "Demo Admin",
			StartURL:    "https://example.com/admin",
			SessionName: "demo_admin_user",
		}, artifactRoot)
		if err != nil {
			t.Fatalf("SaveWorkbenchSiteConfig() error = %v", err)
		}
		_, err = UpsertWorkbenchAPICards("demo_admin", artifactRoot, []WorkbenchAPICard{
			{
				ID:            "api:POST:/api/orders/search",
				SiteID:        "demo_admin",
				Method:        "POST",
				PathTemplate:  "/api/orders/search",
				SemanticName:  "订单搜索",
				OperationType: "read",
				Risk:          "read",
				RequestSchema: map[string]any{
					"status": "string",
				},
			},
		})
		if err != nil {
			t.Fatalf("UpsertWorkbenchAPICards() error = %v", err)
		}

		plan, err := BuildWorkbenchTaskPlan(WorkbenchTaskPlanOptions{
			SiteID:       "demo_admin",
			ArtifactRoot: artifactRoot,
			Intent:       "请执行订单搜索并分析结果",
		})
		if err != nil {
			t.Fatalf("BuildWorkbenchTaskPlan() error = %v", err)
		}
		if plan.Strategy != "api_first" {
			t.Fatalf("unexpected strategy: %q", plan.Strategy)
		}
		if plan.Flow == nil {
			t.Fatalf("expected flow")
		}
		if plan.Flow.Browser == nil || plan.Flow.Browser.UseSession != "demo_admin_user" {
			t.Fatalf("expected flow to reuse session: %#v", plan.Flow.Browser)
		}
		if len(plan.Flow.Steps) < 2 {
			t.Fatalf("expected navigate + http_request steps, got %#v", plan.Flow.Steps)
		}
		last := plan.Flow.Steps[len(plan.Flow.Steps)-1]
		if last.Action != "http_request" {
			t.Fatalf("unexpected last action: %#v", last)
		}
		if got, ok := plan.Flow.Vars["target_url"].(string); !ok || got != "https://example.com/api/orders/search" {
			t.Fatalf("unexpected target_url: %#v", plan.Flow.Vars["target_url"])
		}
	})
}

func TestBuildWorkbenchFallbackShape(t *testing.T) {
	shape := buildWorkbenchFallbackShape("https://example.com/admin/orders", "timeout")
	if shape == nil {
		t.Fatalf("expected fallback shape")
	}
	if shape.Title != "/admin/orders" {
		t.Fatalf("unexpected fallback title: %q", shape.Title)
	}
	if len(shape.Forms) != 0 || len(shape.Actions) != 0 || len(shape.Links) != 0 {
		t.Fatalf("expected empty fallback shape: %#v", shape)
	}
}

func TestObserveWorkbenchPageNoEvaluate(t *testing.T) {
	observation := observeWorkbenchPageNoEvaluate(nil, &workbenchPageShape{
		Title: "订单管理",
	}, PageObservationOptions{
		URL: "https://example.com/admin/orders",
	})
	if observation == nil {
		t.Fatalf("expected observation")
	}
	if observation.Title != "订单管理" {
		t.Fatalf("unexpected observation title: %q", observation.Title)
	}
	if observation.PageSummary != "订单管理" {
		t.Fatalf("unexpected observation summary: %q", observation.PageSummary)
	}
	if len(observation.Elements) != 0 || len(observation.ContentElements) != 0 {
		t.Fatalf("expected no-evaluate observation to skip element extraction: %#v", observation)
	}
	if len(observation.Errors) == 0 || !strings.Contains(observation.Errors[0], "safe mode") {
		t.Fatalf("expected safe-mode marker error, got: %#v", observation.Errors)
	}
}

func TestWorkbenchShouldRetryDirectHTTPFetch(t *testing.T) {
	if !workbenchShouldRetryDirectHTTPFetch(errors.New(`fetch html: Get "https://www.163.com/": proxyconnect tcp: dial tcp 127.0.0.1:7890: connect: connection refused`)) {
		t.Fatalf("expected proxy failure to trigger direct retry")
	}
	if workbenchShouldRetryDirectHTTPFetch(errors.New("fetch html: tls handshake timeout")) {
		t.Fatalf("did not expect generic timeout to trigger proxy retry")
	}
}

func TestBuildWorkbenchPageCaptureSummarySafeMode(t *testing.T) {
	summary := buildWorkbenchPageCaptureSummary(&PageObservation{
		URL:         "https://example.com",
		Title:       "Example",
		PageSummary: "Example",
		Errors:      []string{"observe skipped: no Playwright RPC safe mode"},
	}, []workbenchNetworkRecord{
		{Method: "GET", URL: "https://example.com/api/demo", ResourceType: "xhr", ContentType: "application/json"},
	}, []WorkbenchPageEvent{
		{Type: "frame_navigated", Message: "frame navigated"},
	})
	if summary == nil {
		t.Fatalf("expected summary")
	}
	if !strings.Contains(summary.FilterRule, "不会在 Playwright 回调里读取 response body") {
		t.Fatalf("unexpected filter rule: %q", summary.FilterRule)
	}
	if !strings.Contains(summary.ObservationMode, "默认安全模式") {
		t.Fatalf("unexpected observation mode: %q", summary.ObservationMode)
	}
	if summary.ReadableResponseCount != 0 {
		t.Fatalf("expected no structured responses in safe mode, got %d", summary.ReadableResponseCount)
	}
}

func TestExtractWorkbenchTextSnippetsFromHTML(t *testing.T) {
	htmlBody := `
	<html>
	  <head>
	    <meta name="description" content="飞书开放平台，帮助开发者构建企业应用。">
	  </head>
	  <body>
	    <h1>飞书开放平台</h1>
	    <p>面向企业协作场景的开发者入口。</p>
	    <ul><li>文档中心</li><li>API 参考</li></ul>
	  </body>
	</html>`
	items := extractWorkbenchTextSnippetsFromHTML(htmlBody, 6)
	if len(items) < 3 {
		t.Fatalf("expected snippets, got %#v", items)
	}
	if items[0] != "飞书开放平台，帮助开发者构建企业应用。" {
		t.Fatalf("unexpected first snippet: %#v", items)
	}
}

func TestInferWorkbenchSchemaFromURLQuery(t *testing.T) {
	schema, ok := inferWorkbenchSchemaFromURLQuery("https://example.com/api/orders/search?status=paid&page=2&debug=true").(map[string]any)
	if !ok {
		t.Fatalf("expected schema map")
	}
	if schema["status"] != "string" || schema["page"] != "number" || schema["debug"] != "boolean" {
		t.Fatalf("unexpected query schema: %#v", schema)
	}
}

func TestWorkbenchPostProcessNetworkRecordsPublicDenoise(t *testing.T) {
	records := workbenchPostProcessNetworkRecords([]workbenchNetworkRecord{
		{
			Method:       "GET",
			URL:          "https://gw.m.163.com/search/api/v1/pc-wap/rolling-word?tab=hot",
			ResourceType: "xhr",
			ContentType:  "application/json",
		},
		{
			Method:       "POST",
			URL:          "https://h5.analytics.126.net/news/g",
			ResourceType: "fetch",
			ContentType:  "application/octet-stream",
		},
		{
			Method:       "GET",
			URL:          "https://revive.outin.cn/www/gtr/gtrspc.php?zones=1",
			ResourceType: "xhr",
			ContentType:  "application/json",
		},
	}, "public_html_fallback")
	if len(records) != 1 {
		t.Fatalf("expected only one public API record after denoise, got %#v", records)
	}
	schema, ok := records[0].RequestSchema.(map[string]any)
	if !ok || schema["tab"] != "string" {
		t.Fatalf("expected query schema to be inferred after post process, got %#v", records[0].RequestSchema)
	}
}

func TestBuildWorkbenchLinkGroups(t *testing.T) {
	groups := buildWorkbenchLinkGroups([]WorkbenchLinkCard{
		{Text: "新闻", Href: "https://news.example.com"},
		{Text: "体育", Href: "https://sports.example.com"},
		{Text: "财经", Href: "https://finance.example.com"},
		{Text: "娱乐", Href: "https://ent.example.com"},
		{Text: "登录", Href: "https://passport.example.com/login"},
		{Text: "下载APP", Href: "https://www.example.com/app"},
	})
	if len(groups) < 2 {
		t.Fatalf("expected grouped link summary, got %#v", groups)
	}
	if !strings.Contains(groups[0], "高频入口") {
		t.Fatalf("expected first group to describe hot entries, got %#v", groups)
	}
}

func TestBuildWorkbenchKeyElements(t *testing.T) {
	items := buildWorkbenchKeyElements(
		[]WorkbenchFieldCard{
			{Name: "keyword", Label: "搜索框", Selector: "#kw"},
		},
		[]WorkbenchActionCard{
			{Label: "百度一下", Selector: "#su"},
		},
		[]WorkbenchTableCard{
			{Name: "订单列表", Columns: []string{"订单号", "状态", "金额"}},
		},
		[]WorkbenchLinkCard{
			{Text: "新闻", Href: "https://news.example.com"},
		},
	)
	if len(items) < 4 {
		t.Fatalf("expected key elements, got %#v", items)
	}
	if !strings.Contains(items[0], "输入 · 搜索框") {
		t.Fatalf("unexpected first key element: %#v", items)
	}
}

func TestWorkbenchPageCardSummaryUsesSemanticContext(t *testing.T) {
	shape := &workbenchPageShape{
		Title: "百度",
		Forms: []WorkbenchFormCard{
			{
				Name: "页面输入控件",
				Fields: []WorkbenchFieldCard{
					{Name: "kw", Label: "搜索框", Selector: "#kw"},
				},
			},
		},
		Actions: []WorkbenchActionCard{
			{Label: "百度一下", Kind: "submit", Selector: "#su"},
		},
		Links: []WorkbenchLinkCard{
			{Text: "新闻", Href: "https://news.baidu.com"},
		},
		TextSnippets: []string{
			"百度一下，你就知道",
			"搜索新闻资讯视频地图贴吧文库",
		},
	}

	summary := workbenchPageCardSummary(
		shape.Title,
		shape,
		nil,
		shape.TextSnippets,
		[]string{"输入 · 搜索框 · #kw", "动作 · 百度一下 · #su"},
		[]string{"高频入口：新闻、视频、贴吧"},
		[]WorkbenchPageEvent{{Type: "console"}},
		[]WorkbenchPageAPIHit{{Method: "GET", PathTemplate: "/s"}},
	)

	if !strings.Contains(summary, "页面内容：百度一下，你就知道") {
		t.Fatalf("expected semantic content summary, got %q", summary)
	}
	if !strings.Contains(summary, "关键元素：输入 · 搜索框 · #kw") {
		t.Fatalf("expected key element summary, got %q", summary)
	}
	if !strings.Contains(summary, "运行线索：1 个接口、1 条页面事件") {
		t.Fatalf("expected runtime clue summary, got %q", summary)
	}
}

func TestBuildWorkbenchPageCardPromotesSnippetToTitle(t *testing.T) {
	card := buildWorkbenchPageCard(
		WorkbenchSiteConfig{SiteID: "open_feishu_cn"},
		"run-1",
		"public_html_fallback",
		"https://open.feishu.cn/",
		&workbenchPageShape{
			Title: "/",
			TextSnippets: []string{
				"飞书开放平台",
				"连接企业业务系统与飞书能力",
			},
		},
		nil,
		"",
		nil,
		nil,
	)

	if card.Title != "飞书开放平台" {
		t.Fatalf("expected snippet-promoted title, got %q", card.Title)
	}
	if !strings.Contains(card.Summary, "页面内容：连接企业业务系统与飞书能力") {
		t.Fatalf("expected semantic summary to include text snippet, got %q", card.Summary)
	}
}

func TestMergeWorkbenchShapes(t *testing.T) {
	base := &workbenchPageShape{
		Title: "/",
	}
	overlay := &workbenchPageShape{
		Title: "百度一下，你就知道",
		Forms: []WorkbenchFormCard{
			{Name: "页面输入控件", Fields: []WorkbenchFieldCard{{Label: "搜索框", Selector: "#kw"}}},
		},
		Actions: []WorkbenchActionCard{
			{Label: "百度一下", Selector: "#su"},
		},
		Links: []WorkbenchLinkCard{
			{Text: "新闻", Href: "https://news.baidu.com"},
		},
	}
	merged := mergeWorkbenchShapes(base, overlay, "https://www.baidu.com")
	if merged == nil {
		t.Fatalf("expected merged shape")
	}
	if len(merged.Forms) != 1 || len(merged.Actions) != 1 || len(merged.Links) != 1 {
		t.Fatalf("expected merged key elements, got %#v", merged)
	}
}

func TestCollectWorkbenchPageEventsAndSummary(t *testing.T) {
	events := collectWorkbenchPageEvents([]WorkbenchPageEvent{
		{
			Type:    "frame_navigated",
			Level:   "info",
			Message: "frame navigated",
			URL:     "https://example.com/admin/orders",
		},
		{
			Type:    "popup",
			Level:   "info",
			Message: "popup opened",
			URL:     "https://other.example.net/dialog",
		},
		{
			Type:    "console",
			Level:   "error",
			Message: "search failed",
		},
	}, "https://example.com/admin/orders")
	if len(events) != 2 {
		t.Fatalf("expected 2 related events, got %d: %#v", len(events), events)
	}
	shape := &workbenchPageShape{
		Title: "订单管理",
		Forms: []WorkbenchFormCard{
			{Name: "查询表单", Fields: []WorkbenchFieldCard{{Label: "订单号"}}},
		},
		Actions: []WorkbenchActionCard{{Label: "搜索"}},
		Links:   []WorkbenchLinkCard{{Text: "详情"}},
	}
	summary := workbenchPageCardSummary(
		shape.Title,
		shape,
		nil,
		nil,
		[]string{"输入 · 订单号", "动作 · 搜索"},
		nil,
		events,
		[]WorkbenchPageAPIHit{{Method: "POST", PathTemplate: "/api/orders/search"}},
	)
	if !strings.Contains(summary, "2 条页面事件") {
		t.Fatalf("expected summary to mention event count, got %q", summary)
	}
	if !strings.Contains(summary, "1 个接口") {
		t.Fatalf("expected summary to mention api count, got %q", summary)
	}
}

func TestFlattenWorkbenchPageCardIncludesEvents(t *testing.T) {
	card := WorkbenchPageCard{
		Title: "订单管理",
		Events: []WorkbenchPageEvent{
			{
				Type:    "console",
				Level:   "error",
				Message: "search failed",
				URL:     "https://example.com/admin/orders",
			},
		},
	}
	corpus := flattenWorkbenchPageCard(card)
	if !strings.Contains(corpus, "search failed") {
		t.Fatalf("expected flattened page card to include event message, got %q", corpus)
	}
}

func TestNormalizeWorkbenchExploreURLCanonicalizesRoot(t *testing.T) {
	left := normalizeWorkbenchExploreURL("https://www.163.com")
	right := normalizeWorkbenchExploreURL("https://www.163.com/")
	if left == "" || right == "" {
		t.Fatalf("expected normalized urls to be non-empty")
	}
	if left != right {
		t.Fatalf("expected same normalized url, got %q vs %q", left, right)
	}
}

func TestWorkbenchShouldFollowDiscoveredLinks(t *testing.T) {
	if workbenchShouldFollowDiscoveredLinks("public_html_fallback") {
		t.Fatalf("public_html_fallback should not follow discovered links")
	}
	if !workbenchShouldFollowDiscoveredLinks("authorized_dom_api") {
		t.Fatalf("authorized_dom_api should follow discovered links")
	}
}

func TestWorkbenchShouldCaptureNetworkRequestFiltersStaticAssets(t *testing.T) {
	if workbenchShouldCaptureNetworkRequest("https://example.com/assets/app.js", "fetch") {
		t.Fatalf("expected js asset request to be ignored")
	}
	if workbenchShouldCaptureNetworkRequest("https://example.com/assets/app.css", "xhr") {
		t.Fatalf("expected css asset request to be ignored")
	}
	if !workbenchShouldCaptureNetworkRequest("https://example.com/api/orders/search", "fetch") {
		t.Fatalf("expected api request to be captured")
	}
	if workbenchShouldCaptureNetworkRequest("https://example.com/admin", "document") {
		t.Fatalf("expected non-xhr request to be ignored")
	}
}

func TestBuildWorkbenchPageCaptureSummary(t *testing.T) {
	summary := buildWorkbenchPageCaptureSummary(&PageObservation{
		PageSummary:     `Observed "订单管理" with 3 interactive elements and 2 content elements.`,
		ContentElements: []PageObservationContentElement{{Text: "订单管理"}},
		Elements:        []PageObservationElement{{Text: "搜索"}, {Text: "导出"}, {Text: "订单号"}},
		Errors:          []string{"minor warning"},
	}, []workbenchNetworkRecord{
		{
			URL:          "https://example.com/api/orders/search",
			Method:       "POST",
			ResourceType: "fetch",
			ContentType:  "application/json",
			ResponseSchema: map[string]any{
				"items": []any{},
			},
		},
		{
			URL:          "https://example.com/api/orders/export",
			Method:       "POST",
			ResourceType: "xhr",
			Error:        "timeout",
		},
	}, []WorkbenchPageEvent{
		{Type: "frame_navigated", Message: "frame navigated"},
		{Type: "console", Message: "loaded"},
	})
	if summary == nil {
		t.Fatalf("expected capture summary")
	}
	if summary.NetworkRequestCount != 2 {
		t.Fatalf("unexpected network count: %#v", summary)
	}
	if summary.NetworkFailureCount != 1 {
		t.Fatalf("unexpected failure count: %#v", summary)
	}
	if summary.InteractiveElementCount != 3 || summary.ContentElementCount != 1 {
		t.Fatalf("unexpected observation counts: %#v", summary)
	}
	if !strings.Contains(summary.FilterRule, "xhr/fetch") {
		t.Fatalf("expected filter rule to mention xhr/fetch: %#v", summary)
	}
}

func TestExploreWorkbenchSiteDiscoversSPAMenuRoutes(t *testing.T) {
	t.Skip("browser e2e coverage for SPA menu discovery is manual for now")

	mux := http.NewServeMux()
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!doctype html>
<html>
  <body>
    <aside class="side-menu">
      <div role="menuitem" class="menu-item">订单管理</div>
    </aside>
    <script>
      const item = document.querySelector('[role="menuitem"]');
      item.addEventListener('click', () => {
        history.pushState({}, '', '/orders');
        document.title = '订单管理';
      });
    </script>
  </body>
</html>`)
	})
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!doctype html>
<html>
  <head><title>订单管理</title></head>
  <body>
    <form>
      <input name="keyword" placeholder="订单号" />
    </form>
    <table aria-label="订单列表">
      <thead><tr><th>订单号</th><th>状态</th></tr></thead>
      <tbody><tr><td>1001</td><td>待处理</td></tr></tbody>
    </table>
  </body>
</html>`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	result, err := ExploreWorkbenchSite(WorkbenchExploreOptions{
		Name:           "spa_admin",
		StartURL:       server.URL + "/admin",
		AllowedDomains: []string{parsed.Hostname()},
		ArtifactRoot:   t.TempDir(),
		Headless:       true,
		TimeoutMS:      4000,
		MaxPages:       4,
	})
	if err != nil {
		t.Fatalf("ExploreWorkbenchSite() error = %v", err)
	}
	if len(result.Pages) < 2 {
		t.Fatalf("expected at least 2 pages, got %d", len(result.Pages))
	}
	foundOrders := false
	for _, page := range result.Pages {
		if page.NormalizedRoute == "/orders" {
			foundOrders = true
			break
		}
	}
	if !foundOrders {
		t.Fatalf("expected /orders route in pages: %#v", result.Pages)
	}
}
