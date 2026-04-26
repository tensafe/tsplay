package tsplay_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

type workbenchPageShape struct {
	Title       string                `json:"title"`
	Breadcrumbs []string              `json:"breadcrumbs"`
	Forms       []WorkbenchFormCard   `json:"forms"`
	Tables      []WorkbenchTableCard  `json:"tables"`
	Actions     []WorkbenchActionCard `json:"actions"`
	Links       []WorkbenchLinkCard   `json:"links"`
}

type workbenchNetworkRecorder struct {
	mu         sync.Mutex
	nextID     int
	indexByReq map[playwright.Request]int
	records    []workbenchNetworkRecord
}

func ExploreWorkbenchSite(options WorkbenchExploreOptions) (*WorkbenchExploreResult, error) {
	site, err := resolveWorkbenchExploreSite(options)
	if err != nil {
		return nil, err
	}
	if len(site.AllowedDomains) == 0 {
		parsed, parseErr := url.Parse(site.StartURL)
		if parseErr == nil && parsed.Hostname() != "" {
			site.AllowedDomains = []string{strings.ToLower(parsed.Hostname())}
		}
	}
	maxPages := options.MaxPages
	if maxPages <= 0 {
		maxPages = 8
	}
	timeoutMS := options.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 30000
	}

	runID := "workbench-explore-" + time.Now().Format("20060102-150405.000000000")
	runRoot, err := workbenchSiteRunRoot(site.SiteID, options.ArtifactRoot, runID)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(runRoot, 0755); err != nil {
		return nil, fmt.Errorf("create workbench run root: %w", err)
	}

	startedAt := time.Now().Format(time.RFC3339Nano)
	pw, _, context, page, closeFn, err := launchWorkbenchBrowser(site, options.ArtifactRoot, options.Headless)
	if err != nil {
		return nil, err
	}
	_ = pw
	defer func() {
		_ = closeFn()
	}()

	recorder := newWorkbenchNetworkRecorder(page)
	queue := []string{site.StartURL}
	seen := map[string]struct{}{}
	explored := []string{}
	pageCards := []WorkbenchPageCard{}
	apiCardsByID := map[string]WorkbenchAPICard{}
	entityCardsByID := map[string]WorkbenchEntityCard{}

	for len(queue) > 0 && len(pageCards) < maxPages {
		currentURL := queue[0]
		queue = queue[1:]
		normalizedURL := normalizeWorkbenchExploreURL(currentURL)
		if normalizedURL == "" {
			continue
		}
		if _, ok := seen[normalizedURL]; ok {
			continue
		}
		seen[normalizedURL] = struct{}{}

		pageIndex := len(pageCards) + 1
		segmentStart := recorder.Len()
		if _, err := page.Goto(currentURL, playwright.PageGotoOptions{
			Timeout:   playwright.Float(float64(timeoutMS)),
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); err != nil {
			continue
		}
		_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateNetworkidle,
			Timeout: playwright.Float(float64(timeoutMS)),
		})

		pageRunRoot := filepath.Join(runRoot, fmt.Sprintf("%02d-%s", pageIndex, sanitizeArtifactSegment(normalizeWorkbenchRoute(page.URL()))))
		if err := os.MkdirAll(pageRunRoot, 0755); err != nil {
			return nil, fmt.Errorf("create page run root: %w", err)
		}

		shape, err := extractWorkbenchPageShape(page)
		if err != nil {
			return nil, err
		}
		observation, err := ObserveLoadedPage(page, PageObservationOptions{
			URL:          page.URL(),
			Headless:     options.Headless,
			ArtifactRoot: options.ArtifactRoot,
			TimeoutMS:    timeoutMS,
			RunRoot:      pageRunRoot,
		})
		if err != nil {
			return nil, err
		}
		observationPath := filepath.Join(pageRunRoot, "observation.json")
		if err := writeWorkbenchJSON(observationPath, observation); err != nil {
			return nil, err
		}

		pageCard := buildWorkbenchPageCard(site, runID, page.URL(), shape, observation, observationPath)
		pageCards = append(pageCards, pageCard)
		explored = append(explored, page.URL())

		for _, apiCard := range buildWorkbenchAPICards(site, pageCard, recorder.Since(segmentStart)) {
			apiCardsByID[apiCard.ID] = apiCard
			for _, entity := range buildWorkbenchEntityCards(site.SiteID, apiCard) {
				entityCardsByID[entity.ID] = entity
			}
		}

		for _, link := range shape.Links {
			target, ok := workbenchAllowedExploreLink(link.Href, site.AllowedDomains)
			if !ok {
				continue
			}
			if _, ok := seen[target]; ok {
				continue
			}
			queue = append(queue, target)
		}
		for _, target := range probeWorkbenchNavigationTargets(context, page.URL(), shape.Links, site.AllowedDomains, timeoutMS) {
			if _, ok := seen[target]; ok {
				continue
			}
			queue = append(queue, target)
		}
	}

	apiCards := make([]WorkbenchAPICard, 0, len(apiCardsByID))
	for _, item := range apiCardsByID {
		apiCards = append(apiCards, item)
	}
	sort.Slice(apiCards, func(i, j int) bool {
		if apiCards[i].Method != apiCards[j].Method {
			return apiCards[i].Method < apiCards[j].Method
		}
		return apiCards[i].PathTemplate < apiCards[j].PathTemplate
	})

	entityCards := make([]WorkbenchEntityCard, 0, len(entityCardsByID))
	for _, item := range entityCardsByID {
		entityCards = append(entityCards, item)
	}
	sort.Slice(entityCards, func(i, j int) bool {
		return entityCards[i].Name < entityCards[j].Name
	})

	result := WorkbenchExploreResult{
		Site:         site,
		RunID:        runID,
		RunRoot:      runRoot,
		StartedAt:    startedAt,
		FinishedAt:   time.Now().Format(time.RFC3339Nano),
		ExploredURLs: explored,
		Pages:        pageCards,
		APIs:         apiCards,
		Entities:     entityCards,
	}
	return SaveWorkbenchExploreResult(result, options.ArtifactRoot)
}

func resolveWorkbenchExploreSite(options WorkbenchExploreOptions) (WorkbenchSiteConfig, error) {
	siteID := normalizeWorkbenchSiteID(options.SiteID)
	if siteID != "" {
		site, err := LoadWorkbenchSiteConfig(siteID, options.ArtifactRoot)
		if err != nil {
			return WorkbenchSiteConfig{}, err
		}
		if strings.TrimSpace(options.StartURL) != "" {
			site.StartURL = strings.TrimSpace(options.StartURL)
		}
		if len(options.AllowedDomains) > 0 {
			site.AllowedDomains = normalizeAllowedDomains(options.AllowedDomains, "")
		}
		if strings.TrimSpace(options.SessionName) != "" {
			site.SessionName = strings.TrimSpace(options.SessionName)
		}
		return *site, nil
	}

	startURL := strings.TrimSpace(options.StartURL)
	if startURL == "" {
		return WorkbenchSiteConfig{}, fmt.Errorf("start_url is required when site_id is not provided")
	}
	derivedSiteID := normalizeWorkbenchSiteID(options.Name)
	if derivedSiteID == "" {
		if parsed, err := url.Parse(startURL); err == nil {
			hostname := strings.TrimSpace(parsed.Hostname())
			if hostname != "" {
				derivedSiteID = normalizeWorkbenchSiteID(hostname)
			}
		}
	}
	if derivedSiteID == "" {
		return WorkbenchSiteConfig{}, fmt.Errorf("failed to derive site_id from name or start_url")
	}
	config, err := SaveWorkbenchSiteConfig(WorkbenchSiteConfig{
		SiteID:         derivedSiteID,
		Name:           options.Name,
		StartURL:       startURL,
		AllowedDomains: options.AllowedDomains,
		SessionName:    options.SessionName,
	}, options.ArtifactRoot)
	if err != nil {
		return WorkbenchSiteConfig{}, err
	}
	return *config, nil
}

func launchWorkbenchBrowser(site WorkbenchSiteConfig, artifactRoot string, headless bool) (*playwright.Playwright, playwright.Browser, playwright.BrowserContext, playwright.Page, func() error, error) {
	pw, err := StartPlaywright()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	config := FlowBrowserConfig{
		Headless: &headless,
	}
	if strings.TrimSpace(site.SessionName) != "" {
		savedConfig, err := ResolveFlowSavedSessionBrowserConfig(site.SessionName, artifactRoot)
		if err != nil {
			_ = pw.Stop()
			return nil, nil, nil, nil, nil, err
		}
		if savedConfig != nil {
			if savedConfig.StorageState != "" {
				config.StorageState = savedConfig.StorageState
			}
			if savedConfig.StorageStatePath != "" {
				config.StorageStatePath = savedConfig.StorageStatePath
			}
			config.Persistent = savedConfig.Persistent
			config.Profile = savedConfig.Profile
			config.Session = savedConfig.Session
		}
	}

	options := FlowRunOptions{
		Headless:     headless,
		ArtifactRoot: artifactRoot,
	}
	var browser playwright.Browser
	var context playwright.BrowserContext
	var page playwright.Page
	if config.wantsPersistentContext() {
		context, page, err = launchPersistentFlowBrowser(pw, config, flowBrowserStateRoot(options), nil)
		if err != nil {
			_ = pw.Stop()
			return nil, nil, nil, nil, nil, err
		}
	} else {
		browser, context, page, err = launchFlowBrowser(pw, config, options, nil)
		if err != nil {
			_ = pw.Stop()
			return nil, nil, nil, nil, nil, err
		}
	}
	closeFn := func() error {
		var closeErr error
		if page != nil {
			if err := page.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
		if context != nil {
			if err := context.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
		if browser != nil {
			if err := browser.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
		if err := pw.Stop(); err != nil && closeErr == nil {
			closeErr = err
		}
		return closeErr
	}
	return pw, browser, context, page, closeFn, nil
}

func newWorkbenchNetworkRecorder(page playwright.Page) *workbenchNetworkRecorder {
	recorder := &workbenchNetworkRecorder{
		indexByReq: map[playwright.Request]int{},
		records:    []workbenchNetworkRecord{},
	}
	if page == nil {
		return recorder
	}

	page.OnRequest(func(request playwright.Request) {
		if request == nil {
			return
		}
		headers := redactWorkbenchHeaders(request.Headers())
		contentType := strings.TrimSpace(fmt.Sprint(headers["content-type"]))
		postData, _ := request.PostData()
		record := workbenchNetworkRecord{
			URL:            request.URL(),
			Method:         request.Method(),
			ResourceType:   request.ResourceType(),
			IsNavigation:   request.IsNavigationRequest(),
			RequestHeaders: headers,
			RequestSchema:  inferWorkbenchSchemaFromText(postData, contentType),
		}
		recorder.mu.Lock()
		record.ID = fmt.Sprintf("request-%04d", recorder.nextID+1)
		recorder.nextID++
		recorder.records = append(recorder.records, record)
		recorder.indexByReq[request] = len(recorder.records) - 1
		recorder.mu.Unlock()
	})
	page.OnResponse(func(response playwright.Response) {
		if response == nil || response.Request() == nil {
			return
		}
		headers := redactWorkbenchHeaders(response.Headers())
		recorder.mu.Lock()
		index, ok := recorder.indexByReq[response.Request()]
		if ok && index >= 0 && index < len(recorder.records) {
			recorder.records[index].Status = response.Status()
			recorder.records[index].ContentType = strings.TrimSpace(fmt.Sprint(headers["content-type"]))
			recorder.records[index].ResponseHeaders = headers
		}
		recorder.mu.Unlock()
	})
	page.OnRequestFinished(func(request playwright.Request) {
		if request == nil {
			return
		}
		response, err := request.Response()
		if err != nil || response == nil {
			return
		}
		body, err := response.Body()
		if err != nil {
			return
		}
		recorder.mu.Lock()
		index, ok := recorder.indexByReq[request]
		if ok && index >= 0 && index < len(recorder.records) {
			recorder.records[index].ResponseSchema = inferWorkbenchSchemaFromBytes(limitWorkbenchBody(body), recorder.records[index].ContentType)
		}
		recorder.mu.Unlock()
	})
	page.OnRequestFailed(func(request playwright.Request) {
		if request == nil {
			return
		}
		recorder.mu.Lock()
		index, ok := recorder.indexByReq[request]
		if ok && index >= 0 && index < len(recorder.records) {
			if failure := request.Failure(); failure != nil {
				recorder.records[index].Error = failure.Error()
			}
		}
		recorder.mu.Unlock()
	})

	return recorder
}

func (recorder *workbenchNetworkRecorder) Len() int {
	if recorder == nil {
		return 0
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	return len(recorder.records)
}

func (recorder *workbenchNetworkRecorder) Since(index int) []workbenchNetworkRecord {
	if recorder == nil {
		return nil
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	if index < 0 {
		index = 0
	}
	if index >= len(recorder.records) {
		return nil
	}
	items := make([]workbenchNetworkRecord, 0, len(recorder.records)-index)
	for _, item := range recorder.records[index:] {
		items = append(items, item)
	}
	return items
}

func limitWorkbenchBody(body []byte) []byte {
	const maxBodyBytes = 128 * 1024
	if len(body) <= maxBodyBytes {
		return body
	}
	return body[:maxBodyBytes]
}

func redactWorkbenchHeaders(headers map[string]string) map[string]any {
	if len(headers) == 0 {
		return nil
	}
	result := map[string]any{}
	for key, value := range headers {
		normalized := strings.ToLower(strings.TrimSpace(key))
		switch normalized {
		case "cookie", "authorization", "set-cookie", "proxy-authorization", "x-api-key":
			result[normalized] = "[redacted]"
		default:
			result[normalized] = value
		}
	}
	return result
}

func extractWorkbenchPageShape(page playwright.Page) (*workbenchPageShape, error) {
	value, err := page.Evaluate(`() => {
		const clean = (value, max = 160) => {
			if (!value) return '';
			const text = String(value).replace(/\s+/g, ' ').trim();
			return text.length > max ? text.slice(0, max) + '...' : text;
		};
		const cssEscape = (value) => {
			if (window.CSS && window.CSS.escape) return window.CSS.escape(value);
			return String(value).replace(/["\\#.:>+~*^[\]$()=|/@]/g, '\\$&');
		};
		const quote = (value) => JSON.stringify(String(value));
		const uniqueStrings = (items) => Array.from(new Set(items.filter(Boolean)));
		const targetFor = (element) => {
			if (!element) return '';
			const attrs = ['href', 'data-href', 'data-path', 'data-route', 'data-url', 'to', 'router-link', 'index'];
			for (const attr of attrs) {
				const raw = clean(element.getAttribute(attr) || '', 512);
				if (!raw) continue;
				try {
					return new URL(raw, window.location.href).href;
				} catch (_) {}
			}
			const onClick = clean(element.getAttribute('onclick') || '', 512);
			const match = onClick.match(/['"]([^'"]+)['"]/);
			if (match && match[1]) {
				try {
					return new URL(match[1], window.location.href).href;
				} catch (_) {}
			}
			return '';
		};
		const selectorFor = (element) => {
			if (!element) return '';
			for (const attr of ['data-testid', 'data-test', 'data-cy']) {
				const value = element.getAttribute(attr);
				if (value) return '[' + attr + '=' + quote(value) + ']';
			}
			if (element.id) return '#' + cssEscape(element.id);
			const directTarget = targetFor(element);
			if (directTarget && element.tagName && element.tagName.toLowerCase() === 'a') return 'a[href=' + quote(directTarget) + ']';
			if (element.name) return element.tagName.toLowerCase() + '[name=' + quote(element.name) + ']';
			if (element.getAttribute('aria-label')) return element.tagName.toLowerCase() + '[aria-label=' + quote(element.getAttribute('aria-label')) + ']';
			const role = element.getAttribute('role');
			const text = clean(element.innerText || element.textContent || element.value || '', 48);
			if (role && text) return '[' + 'role=' + quote(role) + ']:has-text(' + quote(text) + ')';
			if (text && element.tagName) return element.tagName.toLowerCase() + ':has-text(' + quote(text) + ')';
			return element.tagName.toLowerCase();
		};
		const textFor = (element) => clean(element.innerText || element.textContent || element.value || '');
		const breadcrumbs = uniqueStrings(Array.from(document.querySelectorAll('[aria-label*="breadcrumb" i] a, [aria-label*="breadcrumb" i] li, nav.breadcrumb a, nav.breadcrumb li, .breadcrumb a, .breadcrumb li, [class*="breadcrumb"] a, [class*="breadcrumb"] li')).map(textFor));
		const forms = Array.from(document.querySelectorAll('form')).map((form, index) => ({
			name: clean(form.getAttribute('aria-label') || form.getAttribute('name') || 'form_' + (index + 1)),
			selector: selectorFor(form),
			fields: Array.from(form.querySelectorAll('input, select, textarea')).map((field, fieldIndex) => ({
				name: clean(field.getAttribute('name') || field.getAttribute('id') || field.getAttribute('placeholder') || 'field_' + (fieldIndex + 1)),
				label: clean((field.labels && field.labels[0] && (field.labels[0].innerText || field.labels[0].textContent)) || field.getAttribute('aria-label') || field.getAttribute('placeholder') || ''),
				selector: selectorFor(field)
			}))
		})).filter((form) => form.fields.length > 0);
		const tables = Array.from(document.querySelectorAll('table')).map((table, index) => ({
			name: clean(table.getAttribute('aria-label') || table.getAttribute('summary') || 'table_' + (index + 1)),
			selector: selectorFor(table),
			columns: uniqueStrings(Array.from(table.querySelectorAll('th')).map(textFor))
		}));
		const actions = Array.from(document.querySelectorAll('button, input[type="button"], input[type="submit"], [role="button"]')).map((element) => ({
			label: textFor(element) || clean(element.getAttribute('value') || element.getAttribute('aria-label')),
			kind: clean(element.getAttribute('type') || element.tagName.toLowerCase()),
			selector: selectorFor(element)
		})).filter((action) => action.label);
		const linkCandidates = Array.from(document.querySelectorAll([
			'a[href]',
			'[data-href]',
			'[data-path]',
			'[data-route]',
			'[data-url]',
			'[to]',
			'[router-link]',
			'[index]',
			'[role="menuitem"]',
			'[role="tab"]',
			'nav a',
			'nav button',
			'nav [role="button"]',
			'aside a',
			'aside button',
			'aside [role="button"]',
			'[class*="menu"] a',
			'[class*="menu"] button',
			'[class*="menu"] [role="menuitem"]',
			'[class*="sidebar"] a',
			'[class*="sidebar"] button',
			'[class*="sidebar"] [role="button"]',
			'[class*="nav"] a',
			'[class*="nav"] button',
			'[class*="nav"] [role="button"]'
		].join(',')));
		const links = linkCandidates.map((element) => ({
			text: textFor(element),
			href: targetFor(element),
			selector: selectorFor(element)
		})).filter((link) => link.href || link.text);
		return {
			title: document.title || '',
			breadcrumbs,
			forms,
			tables,
			actions,
			links
		};
	}`)
	if err != nil {
		return nil, fmt.Errorf("extract workbench page shape: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal workbench page shape: %w", err)
	}
	var shape workbenchPageShape
	if err := json.Unmarshal(encoded, &shape); err != nil {
		return nil, fmt.Errorf("decode workbench page shape: %w", err)
	}
	for i := range shape.Actions {
		shape.Actions[i].Risk = classifyWorkbenchActionRisk(shape.Actions[i].Label)
	}
	return &shape, nil
}

func buildWorkbenchPageCard(site WorkbenchSiteConfig, runID string, rawURL string, shape *workbenchPageShape, observation *PageObservation, observationPath string) WorkbenchPageCard {
	route := normalizeWorkbenchRoute(rawURL)
	card := WorkbenchPageCard{
		ID:              fmt.Sprintf("route:%s:%s", site.SiteID, route),
		SiteID:          site.SiteID,
		URL:             rawURL,
		NormalizedRoute: route,
		Title:           firstNonEmpty(strings.TrimSpace(shape.Title), strings.TrimSpace(observation.Title)),
		MenuPath:        append([]string{}, shape.Breadcrumbs...),
		Breadcrumbs:     append([]string{}, shape.Breadcrumbs...),
		Summary:         firstNonEmpty(strings.TrimSpace(observation.PageSummary), strings.TrimSpace(shape.Title)),
		Forms:           append([]WorkbenchFormCard{}, shape.Forms...),
		Tables:          append([]WorkbenchTableCard{}, shape.Tables...),
		Actions:         append([]WorkbenchActionCard{}, shape.Actions...),
		Links:           append([]WorkbenchLinkCard{}, shape.Links...),
		Risk:            "read",
		ObservationPath: observationPath,
		ExploreRunID:    runID,
		UpdatedAt:       time.Now().Format(time.RFC3339Nano),
	}
	if observation != nil {
		card.ScreenshotPath = observation.ScreenshotPath
		card.DOMSnapshotPath = observation.DOMSnapshotPath
	}
	for _, action := range card.Actions {
		if action.Risk == "write_high" || action.Risk == "critical" {
			card.Risk = "write_high"
			break
		}
		if action.Risk == "write_low" && card.Risk == "read" {
			card.Risk = "write_low"
		}
		if action.Risk == "read_download" && card.Risk == "read" {
			card.Risk = "read_download"
		}
	}
	return card
}

func buildWorkbenchAPICards(site WorkbenchSiteConfig, pageCard WorkbenchPageCard, records []workbenchNetworkRecord) []WorkbenchAPICard {
	items := []WorkbenchAPICard{}
	seen := map[string]struct{}{}
	now := time.Now().Format(time.RFC3339Nano)
	for _, record := range records {
		resourceType := strings.ToLower(strings.TrimSpace(record.ResourceType))
		if resourceType != "xhr" && resourceType != "fetch" {
			continue
		}
		pathTemplate := workbenchPathTemplate(record.URL)
		card := WorkbenchAPICard{
			ID:             fmt.Sprintf("api:%s:%s", strings.ToUpper(strings.TrimSpace(record.Method)), pathTemplate),
			SiteID:         site.SiteID,
			Method:         strings.ToUpper(strings.TrimSpace(record.Method)),
			PathTemplate:   pathTemplate,
			SemanticName:   deriveWorkbenchSemanticName(record.Method, pathTemplate, ""),
			TriggerRoute:   pageCard.NormalizedRoute,
			TriggerAction:  "page_navigation",
			OperationType:  workbenchOperationTypeFromRisk(classifyWorkbenchAPIRisk(record.Method, pathTemplate)),
			RequestSchema:  record.RequestSchema,
			ResponseSchema: record.ResponseSchema,
			Risk:           classifyWorkbenchAPIRisk(record.Method, pathTemplate),
			ResourceType:   record.ResourceType,
			Status:         record.Status,
			ContentType:    record.ContentType,
			URL:            record.URL,
			UpdatedAt:      now,
		}
		if _, ok := seen[card.ID]; ok {
			continue
		}
		seen[card.ID] = struct{}{}
		items = append(items, card)
	}
	return items
}

func buildWorkbenchEntityCards(siteID string, apiCard WorkbenchAPICard) []WorkbenchEntityCard {
	schemaMap, ok := apiCard.ResponseSchema.(map[string]any)
	if !ok || len(schemaMap) == 0 {
		return nil
	}
	entityName, fields := deriveWorkbenchEntityFields(schemaMap)
	if entityName == "" || len(fields) == 0 {
		return nil
	}
	card := WorkbenchEntityCard{
		ID:        fmt.Sprintf("entity:%s:%s", siteID, entityName),
		SiteID:    siteID,
		Name:      entityName,
		Label:     entityName,
		Fields:    fields,
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
	}
	return []WorkbenchEntityCard{card}
}

func deriveWorkbenchEntityFields(schema map[string]any) (string, []WorkbenchEntityField) {
	if items, ok := schema["items"].([]any); ok && len(items) > 0 {
		if itemSchema, ok := items[0].(map[string]any); ok {
			return "items", workbenchFieldsFromSchemaMap(itemSchema)
		}
	}
	for key, value := range schema {
		if nested, ok := value.(map[string]any); ok {
			return key, workbenchFieldsFromSchemaMap(nested)
		}
		if list, ok := value.([]any); ok && len(list) > 0 {
			if nested, ok := list[0].(map[string]any); ok {
				return key, workbenchFieldsFromSchemaMap(nested)
			}
		}
	}
	return "response", workbenchFieldsFromSchemaMap(schema)
}

func workbenchFieldsFromSchemaMap(schema map[string]any) []WorkbenchEntityField {
	keys := make([]string, 0, len(schema))
	for key := range schema {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fields := make([]WorkbenchEntityField, 0, len(keys))
	for _, key := range keys {
		typeName := fmt.Sprint(schema[key])
		switch schema[key].(type) {
		case map[string]any:
			typeName = "object"
		case []any:
			typeName = "array"
		}
		fields = append(fields, WorkbenchEntityField{
			Name:  key,
			Label: key,
			Type:  typeName,
		})
	}
	return fields
}

func normalizeWorkbenchExploreURL(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	parsed.Fragment = ""
	return parsed.String()
}

func workbenchAllowedExploreLink(rawURL string, allowedDomains []string) (string, bool) {
	normalized := normalizeWorkbenchExploreURL(rawURL)
	if normalized == "" {
		return "", false
	}
	parsed, err := url.Parse(normalized)
	if err != nil {
		return "", false
	}
	host := strings.ToLower(parsed.Hostname())
	for _, allowed := range allowedDomains {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if allowed == "" {
			continue
		}
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return normalized, true
		}
	}
	return "", false
}

func probeWorkbenchNavigationTargets(context playwright.BrowserContext, currentURL string, candidates []WorkbenchLinkCard, allowedDomains []string, timeoutMS int) []string {
	if context == nil || len(candidates) == 0 {
		return nil
	}
	currentURL = normalizeWorkbenchExploreURL(currentURL)
	if currentURL == "" {
		return nil
	}
	clickTimeout := timeoutMS / 4
	if clickTimeout < 1500 {
		clickTimeout = 1500
	}
	if clickTimeout > 5000 {
		clickTimeout = 5000
	}
	waitAfterClickMS := clickTimeout
	if waitAfterClickMS < 1200 {
		waitAfterClickMS = 1200
	}

	targets := []string{}
	seenTargets := map[string]struct{}{}
	seenSelectors := map[string]struct{}{}
	probeBudget := 0
	for _, candidate := range candidates {
		if probeBudget >= 6 {
			break
		}
		selector := strings.TrimSpace(candidate.Selector)
		label := strings.TrimSpace(candidate.Text)
		if selector == "" || label == "" {
			continue
		}
		if candidate.Href != "" {
			continue
		}
		if classifyWorkbenchActionRisk(label) != "read" {
			continue
		}
		if _, ok := seenSelectors[selector]; ok {
			continue
		}
		seenSelectors[selector] = struct{}{}
		probeBudget++

		page, err := context.NewPage()
		if err != nil {
			continue
		}
		func() {
			defer func() { _ = page.Close() }()
			if _, err := page.Goto(currentURL, playwright.PageGotoOptions{
				Timeout:   playwright.Float(float64(clickTimeout)),
				WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			}); err != nil {
				return
			}
			page.WaitForTimeout(float64(300))
			before := normalizeWorkbenchExploreURL(page.URL())
			if err := page.Click(selector, playwright.PageClickOptions{
				Timeout: playwright.Float(float64(clickTimeout)),
			}); err != nil {
				return
			}
			_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
				State:   playwright.LoadStateDomcontentloaded,
				Timeout: playwright.Float(float64(clickTimeout)),
			})
			page.WaitForTimeout(float64(waitAfterClickMS))
			after := normalizeWorkbenchExploreURL(page.URL())
			if after == "" || after == before {
				return
			}
			target, ok := workbenchAllowedExploreLink(after, allowedDomains)
			if !ok {
				return
			}
			if _, ok := seenTargets[target]; ok {
				return
			}
			seenTargets[target] = struct{}{}
			targets = append(targets, target)
		}()
	}
	return targets
}
