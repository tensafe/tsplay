package tsplay_core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
