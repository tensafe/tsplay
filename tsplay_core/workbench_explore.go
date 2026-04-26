package tsplay_core

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

type workbenchPageShape struct {
	Title        string                `json:"title"`
	Breadcrumbs  []string              `json:"breadcrumbs"`
	Forms        []WorkbenchFormCard   `json:"forms"`
	Tables       []WorkbenchTableCard  `json:"tables"`
	Actions      []WorkbenchActionCard `json:"actions"`
	Links        []WorkbenchLinkCard   `json:"links"`
	TextSnippets []string              `json:"text_snippets,omitempty"`
}

type workbenchNetworkRecorder struct {
	mu         sync.Mutex
	nextID     int
	indexByReq map[playwright.Request]int
	records    []workbenchNetworkRecord
}

type workbenchEventRecorder struct {
	mu        sync.Mutex
	maxEvents int
	events    []WorkbenchPageEvent
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
	if context != nil {
		timeout := float64(timeoutMS)
		context.SetDefaultTimeout(timeout)
		context.SetDefaultNavigationTimeout(timeout)
	}
	if page != nil {
		timeout := float64(timeoutMS)
		page.SetDefaultTimeout(timeout)
		page.SetDefaultNavigationTimeout(timeout)
	}
	defer func() {
		_ = closeFn()
	}()

	recorder := newWorkbenchNetworkRecorder(page)
	eventRecorder := newWorkbenchEventRecorder(page)
	queue := []string{site.StartURL}
	seen := map[string]struct{}{}
	explored := []string{}
	pageCards := []WorkbenchPageCard{}
	apiCardsByID := map[string]WorkbenchAPICard{}
	entityCardsByID := map[string]WorkbenchEntityCard{}
	exploreMode := workbenchExploreModeForSite(site)
	log.Printf("workbench explore start site=%s run_id=%s mode=%s start_url=%s headless=%v max_pages=%d timeout_ms=%d", site.SiteID, runID, exploreMode, site.StartURL, options.Headless, maxPages, timeoutMS)
	writeWorkbenchExploreStatus(runRoot, map[string]any{
		"site_id":      site.SiteID,
		"run_id":       runID,
		"explore_mode": exploreMode,
		"stage":        "start",
		"started_at":   startedAt,
		"updated_at":   time.Now().Format(time.RFC3339Nano),
	})

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
		eventSegmentStart := eventRecorder.Len()
		log.Printf("workbench explore visit site=%s run_id=%s page_index=%d url=%s", site.SiteID, runID, pageIndex, currentURL)
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":      site.SiteID,
			"run_id":       runID,
			"explore_mode": exploreMode,
			"stage":        "goto_start",
			"page_index":   pageIndex,
			"url":          currentURL,
			"updated_at":   time.Now().Format(time.RFC3339Nano),
		})
		if _, err := page.Goto(currentURL, playwright.PageGotoOptions{
			Timeout:   playwright.Float(float64(timeoutMS)),
			WaitUntil: playwright.WaitUntilStateCommit,
		}); err != nil {
			log.Printf("workbench explore goto failed site=%s run_id=%s url=%s err=%v", site.SiteID, runID, currentURL, err)
			writeWorkbenchExploreStatus(runRoot, map[string]any{
				"site_id":      site.SiteID,
				"run_id":       runID,
				"explore_mode": exploreMode,
				"stage":        "goto_failed",
				"page_index":   pageIndex,
				"url":          currentURL,
				"error":        err.Error(),
				"updated_at":   time.Now().Format(time.RFC3339Nano),
			})
			continue
		}
		log.Printf("workbench explore goto done site=%s run_id=%s page_index=%d landed_url=%s", site.SiteID, runID, pageIndex, page.URL())
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":      site.SiteID,
			"run_id":       runID,
			"explore_mode": exploreMode,
			"stage":        "goto_done",
			"page_index":   pageIndex,
			"url":          page.URL(),
			"updated_at":   time.Now().Format(time.RFC3339Nano),
		})
		settleTimeout := timeoutMS / 3
		if settleTimeout < 1500 {
			settleTimeout = 1500
		}
		if settleTimeout > 5000 {
			settleTimeout = 5000
		}
		log.Printf("workbench explore wait_domcontentloaded site=%s run_id=%s page_index=%d timeout_ms=%d", site.SiteID, runID, pageIndex, settleTimeout)
		_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateDomcontentloaded,
			Timeout: playwright.Float(float64(settleTimeout)),
		})
		log.Printf("workbench explore wait_load site=%s run_id=%s page_index=%d timeout_ms=%d", site.SiteID, runID, pageIndex, minWorkbenchInt(settleTimeout, 2500))
		_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateLoad,
			Timeout: playwright.Float(float64(minWorkbenchInt(settleTimeout, 2500))),
		})
		page.WaitForTimeout(float64(minWorkbenchInt(settleTimeout/2, 1200)))
		log.Printf("workbench explore settled site=%s run_id=%s page_index=%d url=%s", site.SiteID, runID, pageIndex, page.URL())
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":      site.SiteID,
			"run_id":       runID,
			"explore_mode": exploreMode,
			"stage":        "settled",
			"page_index":   pageIndex,
			"url":          page.URL(),
			"updated_at":   time.Now().Format(time.RFC3339Nano),
		})

		pageRunRoot := filepath.Join(runRoot, fmt.Sprintf("%02d-%s", pageIndex, sanitizeArtifactSegment(normalizeWorkbenchRoute(page.URL()))))
		if err := os.MkdirAll(pageRunRoot, 0755); err != nil {
			return nil, fmt.Errorf("create page run root: %w", err)
		}

		shapeURL := page.URL()
		log.Printf("workbench explore shape start site=%s run_id=%s page_index=%d url=%s", site.SiteID, runID, pageIndex, shapeURL)
		shape := buildWorkbenchFallbackShape(shapeURL, "dom probe disabled in safe mode")
		if strings.TrimSpace(site.SessionName) == "" {
			if httpShape, httpErr := extractWorkbenchPageShapeFromHTTP(shapeURL); httpErr == nil && httpShape != nil {
				httpShape.Title = firstNonEmpty(strings.TrimSpace(httpShape.Title), workbenchSafePageTitle(page), normalizeWorkbenchRoute(shapeURL), shapeURL)
				shape = httpShape
			} else if httpErr != nil {
				log.Printf("workbench explore http shape fallback failed site=%s run_id=%s url=%s err=%v", site.SiteID, runID, shapeURL, httpErr)
			}
			if workbenchShapeNeedsElementProbe(shape) {
				log.Printf("workbench explore public key probe start site=%s run_id=%s url=%s", site.SiteID, runID, shapeURL)
				if probedShape, probeErr := probeWorkbenchPublicKeyElementsSafe(context, shapeURL, timeoutMS); probeErr == nil && probedShape != nil {
					shape = mergeWorkbenchShapes(shape, probedShape, shapeURL)
					log.Printf("workbench explore public key probe done site=%s run_id=%s url=%s forms=%d actions=%d links=%d", site.SiteID, runID, shapeURL, len(shape.Forms), len(shape.Actions), len(shape.Links))
				} else if probeErr != nil {
					log.Printf("workbench explore public key probe failed site=%s run_id=%s url=%s err=%v", site.SiteID, runID, shapeURL, probeErr)
				}
			}
		}
		shape.Title = firstNonEmpty(strings.TrimSpace(shape.Title), workbenchSafePageTitle(page), normalizeWorkbenchRoute(shapeURL), shapeURL)
		log.Printf("workbench explore shape done site=%s run_id=%s page_index=%d forms=%d actions=%d links=%d", site.SiteID, runID, pageIndex, len(shape.Forms), len(shape.Actions), len(shape.Links))
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":      site.SiteID,
			"run_id":       runID,
			"explore_mode": exploreMode,
			"stage":        "shape_done",
			"page_index":   pageIndex,
			"url":          page.URL(),
			"forms":        len(shape.Forms),
			"actions":      len(shape.Actions),
			"links":        len(shape.Links),
			"updated_at":   time.Now().Format(time.RFC3339Nano),
		})
		log.Printf("workbench explore observe start site=%s run_id=%s page_index=%d url=%s", site.SiteID, runID, pageIndex, shapeURL)
		observation := observeWorkbenchPageNoEvaluate(nil, shape, PageObservationOptions{
			URL:          shapeURL,
			Headless:     options.Headless,
			ArtifactRoot: options.ArtifactRoot,
			TimeoutMS:    timeoutMS,
			RunRoot:      pageRunRoot,
		})
		log.Printf("workbench explore observe done site=%s run_id=%s page_index=%d interactive=%d content=%d errors=%d", site.SiteID, runID, pageIndex, len(observation.Elements), len(observation.ContentElements), len(observation.Errors))
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":            site.SiteID,
			"run_id":             runID,
			"explore_mode":       exploreMode,
			"stage":              "observe_done",
			"page_index":         pageIndex,
			"url":                shapeURL,
			"interactive_count":  len(observation.Elements),
			"content_count":      len(observation.ContentElements),
			"observation_errors": len(observation.Errors),
			"updated_at":         time.Now().Format(time.RFC3339Nano),
		})
		observationPath := filepath.Join(pageRunRoot, "observation.json")
		if err := writeWorkbenchJSON(observationPath, observation); err != nil {
			return nil, err
		}

		pageRecords := workbenchPostProcessNetworkRecords(recorder.Since(segmentStart), exploreMode)
		pageCard := buildWorkbenchPageCard(site, runID, exploreMode, shapeURL, shape, observation, observationPath, pageRecords, eventRecorder.Since(eventSegmentStart))
		pageCards = append(pageCards, pageCard)
		explored = append(explored, shapeURL)
		log.Printf("workbench explore observed site=%s run_id=%s route=%s inputs=%d actions=%d links=%d events=%d", site.SiteID, runID, pageCard.NormalizedRoute, len(pageCard.InputFields), len(pageCard.Actions), len(pageCard.Links), len(pageCard.Events))

		for _, apiCard := range buildWorkbenchAPICards(site, pageCard, pageRecords) {
			apiCardsByID[apiCard.ID] = apiCard
			for _, entity := range buildWorkbenchEntityCards(site.SiteID, apiCard) {
				entityCardsByID[entity.ID] = entity
			}
		}

		if workbenchShouldFollowDiscoveredLinks(exploreMode) {
			for _, link := range pageCard.Links {
				target, ok := workbenchAllowedExploreLink(link.Href, site.AllowedDomains)
				if !ok {
					continue
				}
				if _, ok := seen[target]; ok {
					continue
				}
				queue = append(queue, target)
			}
			for _, target := range probeWorkbenchNavigationTargets(context, page.URL(), pageCard.Links, site.AllowedDomains, timeoutMS) {
				if _, ok := seen[target]; ok {
					continue
				}
				queue = append(queue, target)
			}
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
		ExploreMode:  exploreMode,
		StartedAt:    startedAt,
		FinishedAt:   time.Now().Format(time.RFC3339Nano),
		ExploredURLs: explored,
		Pages:        pageCards,
		APIs:         apiCards,
		Entities:     entityCards,
	}
	log.Printf("workbench explore save site=%s run_id=%s pages=%d apis=%d entities=%d", site.SiteID, runID, len(pageCards), len(apiCards), len(entityCards))
	writeWorkbenchExploreStatus(runRoot, map[string]any{
		"site_id":      site.SiteID,
		"run_id":       runID,
		"explore_mode": exploreMode,
		"stage":        "persist_start",
		"pages":        len(pageCards),
		"apis":         len(apiCards),
		"entities":     len(entityCards),
		"updated_at":   time.Now().Format(time.RFC3339Nano),
	})
	savedResult, err := SaveWorkbenchExploreResult(result, options.ArtifactRoot)
	if err != nil {
		writeWorkbenchExploreStatus(runRoot, map[string]any{
			"site_id":      site.SiteID,
			"run_id":       runID,
			"explore_mode": exploreMode,
			"stage":        "persist_failed",
			"error":        err.Error(),
			"updated_at":   time.Now().Format(time.RFC3339Nano),
		})
		return nil, err
	}
	writeWorkbenchExploreStatus(runRoot, map[string]any{
		"site_id":      site.SiteID,
		"run_id":       runID,
		"explore_mode": exploreMode,
		"stage":        "persist_done",
		"pages":        len(pageCards),
		"apis":         len(apiCards),
		"entities":     len(entityCards),
		"updated_at":   time.Now().Format(time.RFC3339Nano),
	})
	log.Printf("workbench explore persist done site=%s run_id=%s result=%s", site.SiteID, runID, filepath.Join(runRoot, "result.json"))
	return savedResult, nil
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
		if !workbenchShouldCaptureNetworkRequest(request.URL(), request.ResourceType()) {
			return
		}
		headers := redactWorkbenchHeaders(request.Headers())
		record := workbenchNetworkRecord{
			URL:            request.URL(),
			Method:         request.Method(),
			ResourceType:   request.ResourceType(),
			IsNavigation:   request.IsNavigationRequest(),
			RequestHeaders: headers,
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

	// 重要：不要在 Playwright 事件回调里调用 request.Response()/response.Body()。
	// 这些都是 Playwright RPC，容易在事件分发链路中形成自锁，表现为
	// page.Evaluate/context.NewPage/page.Screenshot/page.Close 等调用长期不返回。
	// Workbench 稳定版只采集 request/response 元信息，响应 body/schema 后续应改为异步后处理。

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

func newWorkbenchEventRecorder(page playwright.Page) *workbenchEventRecorder {
	recorder := &workbenchEventRecorder{
		maxEvents: 120,
		events:    []WorkbenchPageEvent{},
	}
	if page == nil {
		return recorder
	}

	page.OnFrameNavigated(func(frame playwright.Frame) {
		if frame == nil {
			return
		}
		frameURL := strings.TrimSpace(frame.URL())
		if frameURL == "" {
			return
		}
		recorder.append(WorkbenchPageEvent{
			Type:      "frame_navigated",
			Level:     "info",
			Message:   "frame navigated",
			URL:       frameURL,
			Detail:    workbenchCleanText(frame.Name(), 80),
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
	})

	page.OnPopup(func(popup playwright.Page) {
		if popup == nil {
			return
		}
		recorder.append(WorkbenchPageEvent{
			Type:      "popup",
			Level:     "info",
			Message:   "popup opened",
			URL:       strings.TrimSpace(popup.URL()),
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
	})

	page.OnDownload(func(download playwright.Download) {
		if download == nil {
			return
		}
		recorder.append(WorkbenchPageEvent{
			Type:      "download",
			Level:     "info",
			Message:   "download started",
			URL:       strings.TrimSpace(download.URL()),
			Detail:    workbenchCleanText(download.SuggestedFilename(), 120),
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
	})

	page.OnConsole(func(message playwright.ConsoleMessage) {
		if message == nil {
			return
		}
		recorder.append(WorkbenchPageEvent{
			Type:      "console",
			Level:     workbenchConsoleLevel(message.Type()),
			Message:   workbenchCleanText(message.Text(), 240),
			Detail:    workbenchCleanText(message.Type(), 40),
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
	})

	page.OnPageError(func(pageErr error) {
		if pageErr == nil {
			return
		}
		recorder.append(WorkbenchPageEvent{
			Type:      "page_error",
			Level:     "error",
			Message:   workbenchCleanText(pageErr.Error(), 240),
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
	})

	page.OnWebSocket(func(ws playwright.WebSocket) {
		if ws == nil {
			return
		}
		wsURL := strings.TrimSpace(ws.URL())
		recorder.append(WorkbenchPageEvent{
			Type:      "websocket",
			Level:     "info",
			Message:   "websocket opened",
			URL:       wsURL,
			Timestamp: time.Now().Format(time.RFC3339Nano),
		})
		ws.OnSocketError(func(errText string) {
			recorder.append(WorkbenchPageEvent{
				Type:      "websocket_error",
				Level:     "error",
				Message:   workbenchCleanText(errText, 240),
				URL:       wsURL,
				Timestamp: time.Now().Format(time.RFC3339Nano),
			})
		})
		ws.OnClose(func(playwright.WebSocket) {
			recorder.append(WorkbenchPageEvent{
				Type:      "websocket_closed",
				Level:     "info",
				Message:   "websocket closed",
				URL:       wsURL,
				Timestamp: time.Now().Format(time.RFC3339Nano),
			})
		})
	})

	return recorder
}

func observeWorkbenchPageLightUnsafe(page playwright.Page, options PageObservationOptions) (*PageObservation, error) {
	if page == nil {
		return nil, fmt.Errorf("page is nil")
	}
	log.Printf("workbench observe start url=%s run_root=%s", page.URL(), options.RunRoot)

	fullOptions := options
	if fullOptions.MaxElements <= 0 {
		fullOptions.MaxElements = 120
	}
	observation, err := ObserveLoadedPage(page, fullOptions)
	if err == nil {
		log.Printf("workbench observe content url=%s skipped=false count=%d errors=%d", observation.URL, len(observation.ContentElements), len(observation.Errors))
		log.Printf("workbench observe elements url=%s skipped=false count=%d errors=%d", observation.URL, len(observation.Elements), len(observation.Errors))
		return observation, nil
	}
	log.Printf("workbench observe full failed url=%s err=%v", page.URL(), err)

	artifactRoot := strings.TrimSpace(options.ArtifactRoot)
	if artifactRoot == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	root, rootErr := prepareRuntimeFileRoot(artifactRoot)
	if rootErr != nil {
		return nil, fmt.Errorf("prepare artifact root %q: %w", artifactRoot, rootErr)
	}

	dir := filepath.Join(root, "observe-"+time.Now().Format("20060102-150405.000000000"))
	if strings.TrimSpace(options.RunRoot) != "" {
		dir = filepath.Join(options.RunRoot, "observe")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create observation directory: %w", err)
	}

	observation = &PageObservation{
		URL:          page.URL(),
		Title:        firstNonEmpty(workbenchSafePageTitle(page), normalizeWorkbenchRoute(page.URL()), page.URL()),
		ArtifactRoot: firstNonEmpty(strings.TrimSpace(options.RunRoot), root),
		Elements:     []PageObservationElement{},
		Errors:       []string{fmt.Sprintf("full_observe: %v", err)},
	}

	screenshotPath := filepath.Join(dir, "observe.png")
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(screenshotPath),
		FullPage: playwright.Bool(true),
	}); err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("screenshot: %v", err))
	} else {
		observation.ScreenshotPath = screenshotPath
	}

	contentElements, err := extractWorkbenchContentElements(page)
	if err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("content: %v", err))
	} else {
		observation.ContentElements = contentElements
	}
	log.Printf("workbench observe content url=%s skipped=false count=%d errors=%d", observation.URL, len(observation.ContentElements), len(observation.Errors))

	elements, err := observeWorkbenchInteractiveElementsLight(page)
	if err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("elements: %v", err))
	} else {
		observation.Elements = elements
	}
	log.Printf("workbench observe elements url=%s skipped=false count=%d errors=%d", observation.URL, len(observation.Elements), len(observation.Errors))

	observation.PageSummary = buildObservationPageSummary(observation)
	return observation, nil
}

func observeWorkbenchPageNoEvaluate(page playwright.Page, shape *workbenchPageShape, options PageObservationOptions) *PageObservation {
	// 稳定版观察：不调用任何 Playwright RPC。
	// 注意：page.URL()/page.Title()/page.Screenshot() 也可能在 Playwright 通道异常时阻塞。
	_ = page

	rawURL := strings.TrimSpace(options.URL)
	title := ""
	if shape != nil {
		title = strings.TrimSpace(shape.Title)
	}

	observation := &PageObservation{
		URL:             rawURL,
		Title:           firstNonEmpty(title, normalizeWorkbenchRoute(rawURL), rawURL),
		ArtifactRoot:    strings.TrimSpace(options.RunRoot),
		PageSummary:     firstNonEmpty(title, normalizeWorkbenchRoute(rawURL), rawURL),
		Elements:        []PageObservationElement{},
		ContentElements: []PageObservationContentElement{},
		Errors:          []string{"observe skipped: no Playwright RPC safe mode"},
	}

	artifactRoot := strings.TrimSpace(options.ArtifactRoot)
	if artifactRoot == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	root, rootErr := prepareRuntimeFileRoot(artifactRoot)
	if rootErr != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("prepare artifact root %q: %v", artifactRoot, rootErr))
		return observation
	}
	if observation.ArtifactRoot == "" {
		observation.ArtifactRoot = root
	}

	dir := filepath.Join(root, "observe-"+time.Now().Format("20060102-150405.000000000"))
	if strings.TrimSpace(options.RunRoot) != "" {
		dir = filepath.Join(options.RunRoot, "observe")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("create observation directory: %v", err))
		return observation
	}

	return observation
}

func extractWorkbenchPageShapeSafe(context playwright.BrowserContext, currentURL string, site WorkbenchSiteConfig, timeoutMS int) (*workbenchPageShape, error) {
	// 稳定版 shape：不再创建 helperPage，不再执行 page.Evaluate。
	// 之前卡在 context.NewPage()，说明 Playwright RPC 通道已经不可靠；此处必须彻底避开 Playwright RPC。
	_ = context
	_ = timeoutMS

	currentURL = strings.TrimSpace(currentURL)
	if currentURL == "" {
		return buildWorkbenchFallbackShape("", "empty url"), nil
	}

	// 未登录/公开页面场景，优先用 HTTP 静态 HTML 兜底；这不会占用 Playwright 通道。
	if strings.TrimSpace(site.SessionName) == "" {
		httpShape, httpErr := extractWorkbenchPageShapeFromHTTP(currentURL)
		if httpErr == nil && httpShape != nil {
			httpShape.Title = firstNonEmpty(strings.TrimSpace(httpShape.Title), normalizeWorkbenchRoute(currentURL), currentURL)
			return httpShape, nil
		}
		return buildWorkbenchFallbackShape(currentURL, "http shape fallback failed"), httpErr
	}

	// 登录态/授权站点先走保守模式：只返回 URL/route 画像，保证探索主流程不被 DOM 识别拖死。
	return buildWorkbenchFallbackShape(currentURL, "dom probe disabled for authorized site safe mode"), nil
}

func probeWorkbenchPublicKeyElementsSafe(context playwright.BrowserContext, currentURL string, timeoutMS int) (*workbenchPageShape, error) {
	currentURL = strings.TrimSpace(currentURL)
	if currentURL == "" {
		return nil, fmt.Errorf("empty url")
	}
	if context == nil {
		return nil, fmt.Errorf("browser context is nil")
	}
	budgetMS := timeoutMS / 6
	if budgetMS < 1200 {
		budgetMS = 1200
	}
	if budgetMS > 2500 {
		budgetMS = 2500
	}

	type result struct {
		shape *workbenchPageShape
		err   error
	}
	done := make(chan result, 1)

	helperPage, err := context.NewPage()
	if err != nil {
		return nil, err
	}
	helperPage.SetDefaultTimeout(float64(budgetMS))
	helperPage.SetDefaultNavigationTimeout(float64(budgetMS))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- result{err: fmt.Errorf("panic: %v", r)}
			}
		}()
		if _, gotoErr := helperPage.Goto(currentURL, playwright.PageGotoOptions{
			Timeout:   playwright.Float(float64(budgetMS)),
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); gotoErr != nil {
			done <- result{err: gotoErr}
			return
		}
		shape, shapeErr := extractWorkbenchEssentialShapeFromSelectors(helperPage)
		if shapeErr != nil {
			done <- result{err: shapeErr}
			return
		}
		done <- result{shape: shape}
	}()

	select {
	case outcome := <-done:
		closeWorkbenchPageSafely(helperPage, 900*time.Millisecond)
		return outcome.shape, outcome.err
	case <-time.After(time.Duration(budgetMS) * time.Millisecond):
		closeWorkbenchPageSafely(helperPage, 900*time.Millisecond)
		return nil, fmt.Errorf("public key probe timeout after %dms", budgetMS)
	}
}

func extractWorkbenchEssentialShapeFromSelectors(page playwright.Page) (*workbenchPageShape, error) {
	if page == nil {
		return nil, fmt.Errorf("page is nil")
	}
	fields, err := collectWorkbenchInputFieldsLight(page, 10)
	if err != nil {
		return nil, err
	}
	actions, err := collectWorkbenchActionsLight(page, 10)
	if err != nil {
		return nil, err
	}
	links, err := collectWorkbenchLinksLight(page, 16)
	if err != nil {
		return nil, err
	}
	forms := []WorkbenchFormCard{}
	if len(fields) > 0 {
		forms = append(forms, WorkbenchFormCard{
			Name:   "页面输入控件",
			Fields: fields,
		})
	}
	return &workbenchPageShape{
		Title:   firstNonEmpty(workbenchSafePageTitle(page), normalizeWorkbenchRoute(page.URL()), page.URL()),
		Forms:   forms,
		Actions: actions,
		Links:   links,
	}, nil
}

func workbenchShapeNeedsElementProbe(shape *workbenchPageShape) bool {
	if shape == nil {
		return true
	}
	if len(shape.Forms) > 0 || len(shape.Actions) > 0 || len(shape.Links) > 0 {
		return false
	}
	return true
}

func mergeWorkbenchShapes(base *workbenchPageShape, overlay *workbenchPageShape, rawURL string) *workbenchPageShape {
	if base == nil && overlay == nil {
		return buildWorkbenchFallbackShape(rawURL, "merged shape is nil")
	}
	if base == nil {
		return overlay
	}
	if overlay == nil {
		return base
	}
	merged := *base
	merged.Title = firstNonEmpty(strings.TrimSpace(base.Title), strings.TrimSpace(overlay.Title), normalizeWorkbenchRoute(rawURL), rawURL)
	if len(merged.Breadcrumbs) == 0 {
		merged.Breadcrumbs = append([]string{}, overlay.Breadcrumbs...)
	}
	if len(merged.Forms) == 0 {
		merged.Forms = append([]WorkbenchFormCard{}, overlay.Forms...)
	}
	if len(merged.Tables) == 0 {
		merged.Tables = append([]WorkbenchTableCard{}, overlay.Tables...)
	}
	if len(merged.Actions) == 0 {
		merged.Actions = append([]WorkbenchActionCard{}, overlay.Actions...)
	}
	if len(merged.Links) == 0 {
		merged.Links = append([]WorkbenchLinkCard{}, overlay.Links...)
	}
	if len(merged.TextSnippets) == 0 {
		merged.TextSnippets = append([]string{}, overlay.TextSnippets...)
	}
	return &merged
}

func extractWorkbenchPageShapeLightUnsafe(page playwright.Page, site WorkbenchSiteConfig) (*workbenchPageShape, error) {
	if page == nil {
		return nil, fmt.Errorf("page is nil")
	}

	value, err := page.Evaluate(`() => {
		const LIMITS = {
			maxNodes: 1500,
			maxMS: 1200,
			breadcrumbs: 8,
			forms: 4,
			fields: 24,
			tables: 4,
			columns: 12,
			actions: 24,
			links: 24
		};

		const started = performance.now();
		const expired = () => performance.now() - started > LIMITS.maxMS;

		const clean = (value, max = 160) => {
			if (value === undefined || value === null) return '';
			const text = String(value).replace(/\s+/g, ' ').trim();
			return text.length > max ? text.slice(0, max) + '...' : text;
		};

		const quote = (value) => JSON.stringify(String(value));

		const cssEscape = (value) => {
			if (window.CSS && window.CSS.escape) return window.CSS.escape(value);
			return String(value).replace(/["\\#.:>+~*^[\]$()=|/@]/g, '\\$&');
		};

		const attr = (el, name, max = 160) => clean(el && el.getAttribute && el.getAttribute(name), max);

		const textLite = (el, max = 120) => {
			if (!el) return '';
			return clean(
				el.value ||
				el.getAttribute?.('aria-label') ||
				el.getAttribute?.('title') ||
				el.getAttribute?.('placeholder') ||
				el.textContent ||
				'',
				max
			);
		};

		const visibleLite = (el) => {
			if (!el) return false;
			if (el.hidden) return false;
			if (el.getAttribute && el.getAttribute('aria-hidden') === 'true') return false;
			const s = el.style;
			if (s && (s.display === 'none' || s.visibility === 'hidden' || s.opacity === '0')) return false;

			// 只对候选元素做轻量几何判断，避免全量 getComputedStyle
			const rects = el.getClientRects && el.getClientRects();
			if (!rects || rects.length === 0) return false;

			const rect = rects[0];
			return rect.width > 0 && rect.height > 0;
		};

		const targetFor = (el) => {
			if (!el || !el.getAttribute) return '';
			const attrs = ['href', 'data-href', 'data-path', 'data-route', 'data-url', 'to', 'router-link', 'index'];
			for (const name of attrs) {
				const raw = attr(el, name, 512);
				if (!raw) continue;
				if (/^javascript:/i.test(raw)) continue;
				try {
					return new URL(raw, window.location.href).href;
				} catch (_) {}
			}
			return '';
		};

		const selectorFor = (el) => {
			if (!el || !el.tagName) return '';
			const tag = String(el.tagName).toLowerCase();

			for (const name of ['data-testid', 'data-test', 'data-cy']) {
				const value = el.getAttribute(name);
				if (value) return '[' + name + '=' + quote(value) + ']';
			}

			if (el.id) return '#' + cssEscape(el.id);

			const name = el.getAttribute('name');
			if (name) return tag + '[name=' + quote(name) + ']';

			const placeholder = el.getAttribute('placeholder');
			if (placeholder) return tag + '[placeholder=' + quote(placeholder) + ']';

			const aria = el.getAttribute('aria-label');
			if (aria) return tag + '[aria-label=' + quote(aria) + ']';

			const href = targetFor(el);
			if (href && tag === 'a') return 'a[href=' + quote(href) + ']';

			return tag;
		};

		const pushUnique = (arr, item, key, limit) => {
			if (!item || !key) return;
			if (!arr.__seen) Object.defineProperty(arr, '__seen', { value: new Set(), enumerable: false });
			if (arr.__seen.has(key)) return;
			arr.__seen.add(key);
			arr.push(item);
			if (limit > 0 && arr.length > limit) arr.length = limit;
		};

		const isInputLike = (el, tag, role) => {
			if (tag === 'textarea' || tag === 'select') return true;
			if (tag === 'input') {
				const type = attr(el, 'type', 40).toLowerCase();
				return type !== 'hidden';
			}
			return role === 'textbox' || role === 'combobox' || role === 'searchbox';
		};

		const isActionLike = (el, tag, role) => {
			if (tag === 'button') return true;
			if (tag === 'input') {
				const type = attr(el, 'type', 40).toLowerCase();
				return type === 'button' || type === 'submit' || type === 'reset';
			}
			return role === 'button' || role === 'tab';
		};

		const isLinkLike = (el, tag, role) => {
			if (tag === 'a' && attr(el, 'href', 512)) return true;
			if (role === 'link' || role === 'menuitem' || role === 'tab') return true;
			return !!targetFor(el);
		};

		const isBreadcrumbLike = (el) => {
			let cur = el;
			let depth = 0;
			while (cur && depth < 4) {
				const cls = clean(cur.className || '', 120).toLowerCase();
				const aria = attr(cur, 'aria-label', 120).toLowerCase();
				const role = attr(cur, 'role', 40).toLowerCase();
				if (cls.includes('breadcrumb') || aria.includes('breadcrumb') || role === 'navigation') return true;
				cur = cur.parentElement;
				depth++;
			}
			return false;
		};

		const formsByKey = new Map();
		const looseFields = [];
		const tables = [];
		const actions = [];
		const links = [];
		const breadcrumbs = [];

		const appendField = (el, index) => {
			if (!visibleLite(el)) return;

			const tag = String(el.tagName || '').toLowerCase();
			const inputType = attr(el, 'type', 40).toLowerCase();

			if (tag === 'input' && inputType === 'hidden') return;

			const name = clean(
				attr(el, 'name', 80) ||
				el.id ||
				attr(el, 'placeholder', 80) ||
				attr(el, 'aria-label', 80) ||
				tag ||
				('field_' + index),
				80
			);

			const label = clean(
				attr(el, 'aria-label', 120) ||
				attr(el, 'placeholder', 120) ||
				attr(el, 'title', 120) ||
				name ||
				'input',
				120
			);

			const field = {
				name: name || 'input',
				label: label || name || 'input',
				selector: selectorFor(el)
			};

			const form = el.closest && el.closest('form');
			if (form) {
				const formSelector = selectorFor(form);
				const formName = clean(
					attr(form, 'aria-label', 80) ||
					attr(form, 'name', 80) ||
					formSelector ||
					('form_' + (formsByKey.size + 1)),
					80
				);
				const formKey = formSelector || formName;

				if (!formsByKey.has(formKey)) {
					formsByKey.set(formKey, {
						name: formName,
						selector: formSelector,
						fields: [],
						__seen: new Set()
					});
				}

				const bucket = formsByKey.get(formKey);
				const key = field.selector || field.label || field.name;
				if (!bucket.__seen.has(key) && bucket.fields.length < LIMITS.fields) {
					bucket.__seen.add(key);
					bucket.fields.push(field);
				}
			} else {
				pushUnique(looseFields, field, field.selector || field.label || field.name, LIMITS.fields);
			}
		};

		const appendAction = (el) => {
			if (!visibleLite(el)) return;

			const label = clean(
				textLite(el, 80) ||
				attr(el, 'value', 80) ||
				attr(el, 'aria-label', 80) ||
				attr(el, 'title', 80),
				80
			);
			if (!label) return;

			const kind = clean(
				attr(el, 'type', 40) ||
				attr(el, 'role', 40) ||
				String(el.tagName || 'button').toLowerCase(),
				40
			).toLowerCase();

			const item = {
				label,
				kind: kind || 'button',
				selector: selectorFor(el)
			};

			pushUnique(actions, item, item.selector || item.label, LIMITS.actions);
		};

		const appendLink = (el) => {
			if (!visibleLite(el)) return;

			const href = targetFor(el);
			const text = clean(
				textLite(el, 80) ||
				attr(el, 'aria-label', 80) ||
				attr(el, 'title', 80) ||
				href,
				80
			);

			if (!href && !text) return;

			const item = {
				text: text || href,
				href,
				selector: selectorFor(el)
			};

			pushUnique(links, item, item.href || item.selector || item.text, LIMITS.links);

			if (isBreadcrumbLike(el)) {
				pushUnique(breadcrumbs, text, text, LIMITS.breadcrumbs);
			}
		};

		const appendTable = (el) => {
			if (!visibleLite(el)) return;

			const columns = [];
			const headers = el.querySelectorAll('th, [role="columnheader"]');
			for (let i = 0; i < headers.length && columns.length < LIMITS.columns; i++) {
				const value = textLite(headers[i], 120);
				if (value && !columns.includes(value)) columns.push(value);
			}

			const item = {
				name: clean(
					attr(el, 'aria-label', 80) ||
					attr(el, 'summary', 80) ||
					('table_' + (tables.length + 1)),
					80
				),
				selector: selectorFor(el),
				columns
			};

			pushUnique(tables, item, item.selector || item.name, LIMITS.tables);
		};

		const root = document.body || document.documentElement;
		if (!root) {
			return {
				title: clean(document.title || '', 160),
				breadcrumbs: [],
				forms: [],
				tables: [],
				actions: [],
				links: []
			};
		}

		const walker = document.createTreeWalker(root, NodeFilter.SHOW_ELEMENT);
		let node;
		let scanned = 0;
		let fieldIndex = 1;

		while ((node = walker.nextNode())) {
			scanned++;
			if (scanned > LIMITS.maxNodes || expired()) break;

			const tag = String(node.tagName || '').toLowerCase();
			const role = attr(node, 'role', 40).toLowerCase();

			if (isInputLike(node, tag, role)) {
				appendField(node, fieldIndex++);
			}

			if (actions.length < LIMITS.actions && isActionLike(node, tag, role)) {
				appendAction(node);
			}

			if (links.length < LIMITS.links && isLinkLike(node, tag, role)) {
				appendLink(node);
			}

			if (tables.length < LIMITS.tables && (tag === 'table' || role === 'table' || role === 'grid')) {
				appendTable(node);
			}

			if (
				formsByKey.size >= LIMITS.forms &&
				looseFields.length >= LIMITS.fields &&
				actions.length >= LIMITS.actions &&
				links.length >= LIMITS.links &&
				tables.length >= LIMITS.tables
			) {
				break;
			}
		}

		const forms = [];

		for (const form of formsByKey.values()) {
			if (forms.length >= LIMITS.forms) break;
			if (!form.fields || !form.fields.length) continue;
			delete form.__seen;
			forms.push(form);
		}

		if (looseFields.length && forms.length < LIMITS.forms) {
			forms.push({
				name: '页面输入控件',
				selector: '',
				fields: looseFields
			});
		}

		return {
			title: clean(document.title || '', 160),
			breadcrumbs,
			forms,
			tables,
			actions,
			links
		};
	}`)
	if err != nil {
		return nil, fmt.Errorf("extract lightweight workbench page shape: %w", err)
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal lightweight workbench page shape: %w", err)
	}

	var shape workbenchPageShape
	if err := json.Unmarshal(encoded, &shape); err != nil {
		return nil, fmt.Errorf("decode lightweight workbench page shape: %w", err)
	}

	for i := range shape.Actions {
		shape.Actions[i].Risk = classifyWorkbenchActionRisk(shape.Actions[i].Label)
	}

	shape.Title = firstNonEmpty(
		workbenchCleanText(shape.Title, 160),
		workbenchSafePageTitle(page),
		normalizeWorkbenchRoute(page.URL()),
		page.URL(),
	)

	if len(shape.Forms) > 0 ||
		len(shape.Actions) > 0 ||
		len(shape.Links) > 0 ||
		len(shape.Tables) > 0 ||
		len(shape.Breadcrumbs) > 0 {
		return &shape, nil
	}

	if strings.TrimSpace(site.SessionName) == "" {
		httpShape, httpErr := extractWorkbenchPageShapeFromHTTP(page.URL())
		if httpErr == nil && httpShape != nil {
			return httpShape, nil
		}
	}

	return &shape, nil
}

var (
	workbenchTitlePattern           = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	workbenchMetaDescriptionPattern = regexp.MustCompile(`(?is)<meta\b[^>]*(?:name|property)\s*=\s*["'](?:description|og:description)["'][^>]*content\s*=\s*["']([^"']+)["'][^>]*>`)
	workbenchAnchorPattern          = regexp.MustCompile(`(?is)<a\b([^>]*)href\s*=\s*["']([^"'#]+)["']([^>]*)>(.*?)</a>`)
	workbenchButtonPattern          = regexp.MustCompile(`(?is)<button\b([^>]*)>(.*?)</button>`)
	workbenchInputPattern           = regexp.MustCompile(`(?is)<input\b([^>]*)>`)
	workbenchTextareaPattern        = regexp.MustCompile(`(?is)<textarea\b([^>]*)>(.*?)</textarea>`)
	workbenchSelectPattern          = regexp.MustCompile(`(?is)<select\b([^>]*)>(.*?)</select>`)
	workbenchHeadlinePattern        = regexp.MustCompile(`(?is)<h[1-3]\b[^>]*>(.*?)</h[1-3]>`)
	workbenchParagraphPattern       = regexp.MustCompile(`(?is)<p\b[^>]*>(.*?)</p>`)
	workbenchListItemPattern        = regexp.MustCompile(`(?is)<li\b[^>]*>(.*?)</li>`)
	workbenchScriptPattern          = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	workbenchStylePattern           = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	workbenchTagPattern             = regexp.MustCompile(`(?is)<[^>]+>`)
)

func extractWorkbenchPageShapeFromHTTP(rawURL string) (*workbenchPageShape, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	shape, err := extractWorkbenchPageShapeFromHTTPWithClient(rawURL, client)
	if err == nil {
		return shape, nil
	}
	if workbenchShouldRetryDirectHTTPFetch(err) {
		directShape, directErr := extractWorkbenchPageShapeFromHTTPWithClient(rawURL, &http.Client{
			Timeout: 12 * time.Second,
			Transport: &http.Transport{
				Proxy: nil,
			},
		})
		if directErr == nil {
			return directShape, nil
		}
		return nil, directErr
	}
	return nil, err
}

func extractWorkbenchPageShapeFromHTTPWithClient(rawURL string, client *http.Client) (*workbenchPageShape, error) {
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build http request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 TSPlay Workbench")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch html: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 900_000))
	if err != nil {
		return nil, fmt.Errorf("read html: %w", err)
	}
	htmlBody := string(body)
	cleanHTMLBody := workbenchScriptPattern.ReplaceAllString(htmlBody, " ")
	cleanHTMLBody = workbenchStylePattern.ReplaceAllString(cleanHTMLBody, " ")
	baseURL, _ := url.Parse(rawURL)

	fields := make([]WorkbenchFieldCard, 0, 16)
	actions := make([]WorkbenchActionCard, 0, 16)
	links := make([]WorkbenchLinkCard, 0, 24)

	fieldSeen := map[string]struct{}{}
	actionSeen := map[string]struct{}{}
	linkSeen := map[string]struct{}{}

	appendField := func(name string, label string, selector string) {
		key := firstNonEmpty(selector, label, name)
		if key == "" {
			return
		}
		if _, ok := fieldSeen[key]; ok {
			return
		}
		fieldSeen[key] = struct{}{}
		fields = append(fields, WorkbenchFieldCard{
			Name:     firstNonEmpty(name, "input"),
			Label:    firstNonEmpty(label, name, "input"),
			Selector: selector,
		})
	}
	appendAction := func(label string, kind string, selector string) {
		key := firstNonEmpty(selector, label)
		if key == "" {
			return
		}
		if _, ok := actionSeen[key]; ok {
			return
		}
		actionSeen[key] = struct{}{}
		actions = append(actions, WorkbenchActionCard{
			Label:    label,
			Kind:     firstNonEmpty(kind, "button"),
			Selector: selector,
			Risk:     classifyWorkbenchActionRisk(label),
		})
	}
	appendLink := func(text string, href string, selector string) {
		key := firstNonEmpty(href, selector, text)
		if key == "" {
			return
		}
		if _, ok := linkSeen[key]; ok {
			return
		}
		linkSeen[key] = struct{}{}
		links = append(links, WorkbenchLinkCard{
			Text:     firstNonEmpty(text, href),
			Href:     href,
			Selector: selector,
		})
	}

	for _, match := range workbenchInputPattern.FindAllStringSubmatch(cleanHTMLBody, 24) {
		attrs := match[1]
		inputType := workbenchAttrValue(attrs, "type")
		if strings.EqualFold(inputType, "hidden") {
			continue
		}
		name := workbenchAttrValue(attrs, "name")
		id := workbenchAttrValue(attrs, "id")
		placeholder := workbenchAttrValue(attrs, "placeholder")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		appendField(
			firstNonEmpty(name, id, inputType, "input"),
			firstNonEmpty(ariaLabel, placeholder, name, id, inputType, "input"),
			workbenchSimpleSelector(id, name, placeholder, ariaLabel, ""),
		)
	}
	for _, match := range workbenchTextareaPattern.FindAllStringSubmatch(cleanHTMLBody, 12) {
		attrs := match[1]
		name := workbenchAttrValue(attrs, "name")
		id := workbenchAttrValue(attrs, "id")
		placeholder := workbenchAttrValue(attrs, "placeholder")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		appendField(
			firstNonEmpty(name, id, "textarea"),
			firstNonEmpty(ariaLabel, placeholder, name, id, "textarea"),
			workbenchSimpleSelector(id, name, placeholder, ariaLabel, ""),
		)
	}
	for _, match := range workbenchSelectPattern.FindAllStringSubmatch(cleanHTMLBody, 12) {
		attrs := match[1]
		name := workbenchAttrValue(attrs, "name")
		id := workbenchAttrValue(attrs, "id")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		appendField(
			firstNonEmpty(name, id, "select"),
			firstNonEmpty(ariaLabel, name, id, "select"),
			workbenchSimpleSelector(id, name, "", ariaLabel, ""),
		)
	}
	for _, match := range workbenchButtonPattern.FindAllStringSubmatch(cleanHTMLBody, 20) {
		attrs := match[1]
		name := workbenchAttrValue(attrs, "name")
		id := workbenchAttrValue(attrs, "id")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		label := firstNonEmpty(workbenchStripHTML(match[2]), ariaLabel, name, id)
		if label == "" {
			continue
		}
		appendAction(label, "button", workbenchSimpleSelector(id, name, "", ariaLabel, ""))
	}
	for _, match := range workbenchInputPattern.FindAllStringSubmatch(cleanHTMLBody, 16) {
		attrs := match[1]
		inputType := strings.ToLower(strings.TrimSpace(workbenchAttrValue(attrs, "type")))
		if inputType != "button" && inputType != "submit" {
			continue
		}
		name := workbenchAttrValue(attrs, "name")
		id := workbenchAttrValue(attrs, "id")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		value := workbenchAttrValue(attrs, "value")
		label := firstNonEmpty(value, ariaLabel, name, id, inputType)
		appendAction(label, inputType, workbenchSimpleSelector(id, name, "", ariaLabel, ""))
	}
	for _, match := range workbenchAnchorPattern.FindAllStringSubmatch(cleanHTMLBody, 24) {
		href := strings.TrimSpace(match[2])
		if href == "" {
			continue
		}
		resolved := href
		if baseURL != nil {
			if ref, err := url.Parse(href); err == nil {
				resolved = baseURL.ResolveReference(ref).String()
			}
		}
		attrs := match[1] + " " + match[3]
		id := workbenchAttrValue(attrs, "id")
		ariaLabel := workbenchAttrValue(attrs, "aria-label")
		text := firstNonEmpty(workbenchStripHTML(match[4]), ariaLabel, resolved)
		appendLink(text, resolved, workbenchSimpleSelector(id, "", "", ariaLabel, resolved))
	}

	title := firstNonEmpty(workbenchExtractTitle(cleanHTMLBody), normalizeWorkbenchRoute(rawURL), rawURL)
	forms := []WorkbenchFormCard{}
	if len(fields) > 0 {
		forms = append(forms, WorkbenchFormCard{Name: "页面输入控件", Fields: fields})
	}
	return &workbenchPageShape{
		Title:        title,
		Forms:        forms,
		Actions:      actions,
		Links:        links,
		TextSnippets: extractWorkbenchTextSnippetsFromHTML(cleanHTMLBody, 8),
	}, nil
}

func workbenchShouldRetryDirectHTTPFetch(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
	}
	if strings.Contains(message, "proxyconnect") {
		return true
	}
	if strings.Contains(message, "127.0.0.1:7890") {
		return true
	}
	if strings.Contains(message, "connection refused") && strings.Contains(message, "proxy") {
		return true
	}
	return false
}

func collectWorkbenchInputFieldsLight(page playwright.Page, limit int) ([]WorkbenchFieldCard, error) {
	elements, err := page.QuerySelectorAll("input, textarea, select")
	if err != nil {
		return nil, fmt.Errorf("query inputs: %w", err)
	}
	items := make([]WorkbenchFieldCard, 0, limit)
	seen := map[string]struct{}{}
	for _, element := range elements {
		if limit > 0 && len(items) >= limit {
			break
		}
		inputType, _ := element.GetAttribute("type")
		if strings.EqualFold(strings.TrimSpace(inputType), "hidden") {
			continue
		}
		name, _ := element.GetAttribute("name")
		id, _ := element.GetAttribute("id")
		placeholder, _ := element.GetAttribute("placeholder")
		ariaLabel, _ := element.GetAttribute("aria-label")
		label := firstNonEmpty(
			workbenchCleanText(ariaLabel, 80),
			workbenchCleanText(placeholder, 80),
			workbenchCleanText(name, 80),
			workbenchCleanText(id, 80),
			workbenchCleanText(inputType, 80),
			"input",
		)
		selector := workbenchSimpleSelector(id, name, placeholder, ariaLabel, "")
		key := firstNonEmpty(selector, label)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, WorkbenchFieldCard{
			Name:     firstNonEmpty(workbenchCleanText(name, 80), workbenchCleanText(id, 80), workbenchCleanText(inputType, 80), "input"),
			Label:    label,
			Selector: selector,
		})
	}
	return items, nil
}

func collectWorkbenchActionsLight(page playwright.Page, limit int) ([]WorkbenchActionCard, error) {
	elements, err := page.QuerySelectorAll("button, input[type='button'], input[type='submit'], [role='button']")
	if err != nil {
		return nil, fmt.Errorf("query actions: %w", err)
	}
	items := make([]WorkbenchActionCard, 0, limit)
	seen := map[string]struct{}{}
	for _, element := range elements {
		if limit > 0 && len(items) >= limit {
			break
		}
		text, _ := element.TextContent()
		value, _ := element.GetAttribute("value")
		name, _ := element.GetAttribute("name")
		id, _ := element.GetAttribute("id")
		ariaLabel, _ := element.GetAttribute("aria-label")
		actionType, _ := element.GetAttribute("type")
		label := firstNonEmpty(
			workbenchCleanText(text, 80),
			workbenchCleanText(value, 80),
			workbenchCleanText(ariaLabel, 80),
			workbenchCleanText(name, 80),
			workbenchCleanText(id, 80),
		)
		if label == "" {
			continue
		}
		selector := workbenchSimpleSelector(id, name, "", ariaLabel, "")
		key := firstNonEmpty(selector, label)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, WorkbenchActionCard{
			Label:    label,
			Kind:     firstNonEmpty(workbenchCleanText(actionType, 40), "button"),
			Selector: selector,
			Risk:     classifyWorkbenchActionRisk(label),
		})
	}
	return items, nil
}

func collectWorkbenchLinksLight(page playwright.Page, limit int) ([]WorkbenchLinkCard, error) {
	elements, err := page.QuerySelectorAll("a[href]")
	if err != nil {
		return nil, fmt.Errorf("query links: %w", err)
	}
	items := make([]WorkbenchLinkCard, 0, limit)
	seen := map[string]struct{}{}
	baseURL, _ := url.Parse(page.URL())
	for _, element := range elements {
		if limit > 0 && len(items) >= limit {
			break
		}
		href, _ := element.GetAttribute("href")
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(strings.ToLower(href), "javascript:") {
			continue
		}
		resolved := href
		if baseURL != nil {
			if ref, err := url.Parse(href); err == nil {
				resolved = baseURL.ResolveReference(ref).String()
			}
		}
		text, _ := element.TextContent()
		id, _ := element.GetAttribute("id")
		ariaLabel, _ := element.GetAttribute("aria-label")
		label := firstNonEmpty(
			workbenchCleanText(text, 80),
			workbenchCleanText(ariaLabel, 80),
			workbenchCleanText(resolved, 120),
		)
		selector := workbenchSimpleSelector(id, "", "", ariaLabel, resolved)
		key := firstNonEmpty(resolved, selector, label)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, WorkbenchLinkCard{
			Text:     label,
			Href:     resolved,
			Selector: selector,
		})
	}
	return items, nil
}

func workbenchSimpleSelector(id string, name string, placeholder string, ariaLabel string, href string) string {
	id = strings.TrimSpace(id)
	if id != "" {
		return "#" + id
	}
	name = strings.TrimSpace(name)
	if name != "" {
		return "[name=\"" + name + "\"]"
	}
	placeholder = strings.TrimSpace(placeholder)
	if placeholder != "" {
		return "[placeholder=\"" + placeholder + "\"]"
	}
	ariaLabel = strings.TrimSpace(ariaLabel)
	if ariaLabel != "" {
		return "[aria-label=\"" + ariaLabel + "\"]"
	}
	href = strings.TrimSpace(href)
	if href != "" {
		return "a[href=\"" + href + "\"]"
	}
	return ""
}

func workbenchCleanText(value string, max int) string {
	value = strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if max > 0 && len([]rune(value)) > max {
		runes := []rune(value)
		return string(runes[:max]) + "..."
	}
	return value
}

func workbenchPageCardSummary(title string, shape *workbenchPageShape, observation *PageObservation, textSnippets []string, keyElements []string, linkGroups []string, events []WorkbenchPageEvent, apiHits []WorkbenchPageAPIHit) string {
	inputCount := 0
	actionCount := 0
	linkCount := 0
	if shape != nil {
		for _, form := range shape.Forms {
			inputCount += len(form.Fields)
		}
		actionCount = len(shape.Actions)
		linkCount = len(shape.Links)
	}

	parts := make([]string, 0, 4)
	if snippet := buildWorkbenchSummarySnippet(title, textSnippets); snippet != "" {
		parts = append(parts, "页面内容："+snippet)
	}
	if len(keyElements) > 0 {
		parts = append(parts, "关键元素："+strings.Join(keyElements[:minWorkbenchInt(len(keyElements), 3)], "；"))
	} else if inputCount > 0 || actionCount > 0 || linkCount > 0 {
		parts = append(parts, fmt.Sprintf("页面包含 %d 个输入控件、%d 个动作、%d 个链接", inputCount, actionCount, linkCount))
	}
	if len(linkGroups) > 0 {
		parts = append(parts, "导航线索："+strings.Join(linkGroups[:minWorkbenchInt(len(linkGroups), 2)], "；"))
	}
	if len(apiHits) > 0 || len(events) > 0 {
		captureParts := make([]string, 0, 2)
		if len(apiHits) > 0 {
			captureParts = append(captureParts, fmt.Sprintf("%d 个接口", len(apiHits)))
		}
		if len(events) > 0 {
			captureParts = append(captureParts, fmt.Sprintf("%d 条页面事件", len(events)))
		}
		if len(captureParts) > 0 {
			parts = append(parts, "运行线索："+strings.Join(captureParts, "、"))
		}
	}
	if len(parts) > 0 {
		return workbenchCleanText(strings.Join(parts, "。")+"。", 220)
	}
	if observation != nil {
		return firstNonEmpty(strings.TrimSpace(observation.PageSummary), strings.TrimSpace(observation.Title), strings.TrimSpace(observation.URL))
	}
	if shape != nil {
		return strings.TrimSpace(shape.Title)
	}
	return strings.TrimSpace(title)
}

func buildWorkbenchSummarySnippet(title string, snippets []string) string {
	title = workbenchCleanText(strings.TrimSpace(title), 120)
	items := make([]string, 0, 2)
	seen := map[string]struct{}{}
	appendItem := func(value string) {
		value = workbenchCleanText(strings.TrimSpace(value), 100)
		if value == "" || len([]rune(value)) < 4 {
			return
		}
		normalized := strings.ToLower(strings.Join(strings.Fields(value), " "))
		if _, ok := seen[normalized]; ok {
			return
		}
		seen[normalized] = struct{}{}
		items = append(items, value)
	}
	for _, item := range snippets {
		if len(items) >= 2 {
			break
		}
		if title != "" && workbenchEquivalentSummaryText(title, item) {
			continue
		}
		appendItem(item)
	}
	return strings.Join(items, "；")
}

func workbenchEquivalentSummaryText(left string, right string) bool {
	normalize := func(value string) string {
		value = strings.ToLower(strings.TrimSpace(value))
		value = strings.ReplaceAll(value, "，", "")
		value = strings.ReplaceAll(value, "。", "")
		value = strings.ReplaceAll(value, "：", "")
		value = strings.ReplaceAll(value, " ", "")
		return value
	}
	lv := normalize(left)
	rv := normalize(right)
	if lv == "" || rv == "" {
		return false
	}
	return lv == rv
}

func workbenchExtractTitle(htmlBody string) string {
	match := workbenchTitlePattern.FindStringSubmatch(htmlBody)
	if len(match) < 2 {
		return ""
	}
	return workbenchStripHTML(match[1])
}

func workbenchAttrValue(attrs string, name string) string {
	pattern := regexp.MustCompile(`(?is)\b` + regexp.QuoteMeta(name) + `\s*=\s*["']([^"']*)["']`)
	match := pattern.FindStringSubmatch(attrs)
	if len(match) < 2 {
		return ""
	}
	return workbenchCleanText(html.UnescapeString(match[1]), 160)
}

func workbenchStripHTML(value string) string {
	value = workbenchTagPattern.ReplaceAllString(value, " ")
	value = html.UnescapeString(value)
	return workbenchCleanText(value, 160)
}

func extractWorkbenchTextSnippetsFromHTML(htmlBody string, limit int) []string {
	if strings.TrimSpace(htmlBody) == "" || limit <= 0 {
		return nil
	}
	items := make([]string, 0, limit)
	seen := map[string]struct{}{}
	appendSnippet := func(value string) {
		value = workbenchCleanText(workbenchStripHTML(value), 180)
		if value == "" || len([]rune(value)) < 4 {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}

	if match := workbenchMetaDescriptionPattern.FindStringSubmatch(htmlBody); len(match) >= 2 {
		appendSnippet(match[1])
	}
	for _, match := range workbenchHeadlinePattern.FindAllStringSubmatch(htmlBody, limit) {
		if len(match) >= 2 {
			appendSnippet(match[1])
		}
		if len(items) >= limit {
			return items
		}
	}
	for _, match := range workbenchParagraphPattern.FindAllStringSubmatch(htmlBody, limit*2) {
		if len(match) >= 2 {
			appendSnippet(match[1])
		}
		if len(items) >= limit {
			return items
		}
	}
	for _, match := range workbenchListItemPattern.FindAllStringSubmatch(htmlBody, limit*3) {
		if len(match) >= 2 {
			appendSnippet(match[1])
		}
		if len(items) >= limit {
			return items
		}
	}
	return items
}

func observeWorkbenchInteractiveElementsLight(page playwright.Page) ([]PageObservationElement, error) {
	value, err := page.Evaluate(`() => {
		const clean = (value, max = 160) => {
			if (!value) return '';
			const text = String(value).replace(/\s+/g, ' ').trim();
			return text.length > max ? text.slice(0, max) + '...' : text;
		};
		const visible = (element) => {
			const style = window.getComputedStyle(element);
			return style.display !== 'none' &&
				style.visibility !== 'hidden' &&
				style.opacity !== '0' &&
				element.offsetWidth > 0 &&
				element.offsetHeight > 0;
		};
		const quote = (value) => JSON.stringify(String(value));
		const cssEscape = (value) => {
			if (window.CSS && window.CSS.escape) return window.CSS.escape(value);
			return String(value).replace(/["\\#.:>+~*^[\]$()=|/@]/g, '\\$&');
		};
		const selectorCandidates = (element, tag, type, text, placeholder, ariaLabel, href) => {
			const items = [];
			const add = (value) => {
				if (value && !items.includes(value)) items.push(value);
			};
			if (element.id) add('#' + cssEscape(element.id));
			if (element.name) add(tag + '[name=' + quote(element.name) + ']');
			if (placeholder) add(tag + '[placeholder=' + quote(placeholder) + ']');
			if (ariaLabel) add(tag + '[aria-label=' + quote(ariaLabel) + ']');
			if (href && tag === 'a') add('a[href=' + quote(href) + ']');
			if (text && (tag === 'button' || tag === 'a')) add('text=' + quote(text));
			return items;
		};
		const elements = [];
		for (const element of document.querySelectorAll('a[href], button, input, textarea, select, [role="button"], [role="link"], [role="textbox"]')) {
			if (elements.length >= 120) break;
			if (!visible(element)) continue;
			const tag = String(element.tagName || '').toLowerCase();
			const type = tag === 'input' ? String(element.getAttribute('type') || 'text').toLowerCase() : tag;
			const text = clean(element.innerText || element.textContent || element.value || '');
			const placeholder = clean(element.getAttribute('placeholder') || '');
			const ariaLabel = clean(element.getAttribute('aria-label') || '');
			const href = clean(element.getAttribute('href') || '', 512);
			elements.push({
				tag,
				type,
				role: clean(element.getAttribute('role') || ''),
				id: clean(element.id || ''),
				name: clean(element.getAttribute('name') || ''),
				text,
				label: clean(element.getAttribute('title') || ''),
				placeholder,
				aria_label: ariaLabel,
				href,
				value: clean(element.value || ''),
				visible: true,
				enabled: !element.disabled && element.getAttribute('aria-disabled') !== 'true',
				selector_candidates: selectorCandidates(element, tag, type, text, placeholder, ariaLabel, href)
			});
		}
		return elements;
	}`)
	if err != nil {
		return nil, fmt.Errorf("observe lightweight interactive elements: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal lightweight interactive elements: %w", err)
	}
	var elements []PageObservationElement
	if err := json.Unmarshal(encoded, &elements); err != nil {
		return nil, fmt.Errorf("decode lightweight interactive elements: %w", err)
	}
	for i := range elements {
		elements[i].Index = i + 1
		normalizeObservedSelectorDiagnostics(&elements[i])
	}
	return elements, nil
}

func extractWorkbenchContentElements(page playwright.Page) ([]PageObservationContentElement, error) {
	value, err := page.Evaluate(`() => {
		const clean = (value, max = 180) => {
			if (!value) return '';
			const text = String(value).replace(/\s+/g, ' ').trim();
			return text.length > max ? text.slice(0, max) + '...' : text;
		};
		const visible = (element) => {
			const style = window.getComputedStyle(element);
			const rect = element.getBoundingClientRect();
			return style.display !== 'none' &&
				style.visibility !== 'hidden' &&
				style.opacity !== '0' &&
				rect.width > 0 &&
				rect.height > 0;
		};
		const items = [];
		for (const element of document.querySelectorAll('h1, h2, h3, p, li, label, a[href]')) {
			if (items.length >= 20) break;
			if (!visible(element)) continue;
			const text = clean(element.innerText || element.textContent || '');
			if (!text) continue;
			const tag = String(element.tagName || '').toLowerCase();
			let kind = 'text';
			if (tag === 'a') kind = 'link';
			if (tag === 'h1' || tag === 'h2' || tag === 'h3') kind = 'headline';
			items.push({
				kind,
				tag,
				text,
				href: tag === 'a' ? (element.getAttribute('href') || '') : ''
			});
		}
		return items;
	}`)
	if err != nil {
		return nil, fmt.Errorf("extract lightweight content: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal lightweight content: %w", err)
	}
	var items []PageObservationContentElement
	if err := json.Unmarshal(encoded, &items); err != nil {
		return nil, fmt.Errorf("decode lightweight content: %w", err)
	}
	for i := range items {
		items[i].Index = i + 1
	}
	return items, nil
}

func workbenchShouldCaptureNetworkRequest(rawURL string, resourceType string) bool {
	resourceType = strings.ToLower(strings.TrimSpace(resourceType))
	if resourceType != "xhr" && resourceType != "fetch" {
		return false
	}
	if workbenchIsInternalControlAPI(rawURL) {
		return false
	}
	return !workbenchLooksLikeStaticAsset(rawURL)
}

func workbenchIsInternalControlAPI(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(parsed.Path), "/api/workbench/")
}

func workbenchLooksLikeStaticAsset(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}
	pathValue := strings.ToLower(strings.TrimSpace(parsed.Path))
	if pathValue == "" {
		return false
	}
	switch {
	case strings.HasSuffix(pathValue, ".html"),
		strings.HasSuffix(pathValue, ".htm"),
		strings.HasSuffix(pathValue, ".css"),
		strings.HasSuffix(pathValue, ".js"),
		strings.HasSuffix(pathValue, ".mjs"),
		strings.HasSuffix(pathValue, ".png"),
		strings.HasSuffix(pathValue, ".jpg"),
		strings.HasSuffix(pathValue, ".jpeg"),
		strings.HasSuffix(pathValue, ".gif"),
		strings.HasSuffix(pathValue, ".svg"),
		strings.HasSuffix(pathValue, ".ico"),
		strings.HasSuffix(pathValue, ".webp"),
		strings.HasSuffix(pathValue, ".avif"),
		strings.HasSuffix(pathValue, ".woff"),
		strings.HasSuffix(pathValue, ".woff2"),
		strings.HasSuffix(pathValue, ".ttf"),
		strings.HasSuffix(pathValue, ".otf"),
		strings.HasSuffix(pathValue, ".map"),
		strings.HasSuffix(pathValue, ".mp4"),
		strings.HasSuffix(pathValue, ".mp3"),
		strings.HasSuffix(pathValue, ".wav"):
		return true
	default:
		return false
	}
}

func minWorkbenchInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func workbenchSafePageTitle(page playwright.Page) string {
	if page == nil {
		return ""
	}
	title, err := page.Title()
	if err != nil {
		return ""
	}
	return title
}

func closeWorkbenchPageSafely(page playwright.Page, wait time.Duration) {
	if page == nil {
		return
	}
	done := make(chan struct{}, 1)
	go func() {
		_ = page.Close()
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-time.After(wait):
	}
}

func workbenchExploreModeForSite(site WorkbenchSiteConfig) string {
	if strings.TrimSpace(site.SessionName) != "" {
		return "authorized_dom_api"
	}
	return "public_html_fallback"
}

func workbenchShouldFollowDiscoveredLinks(exploreMode string) bool {
	return strings.TrimSpace(exploreMode) == "authorized_dom_api"
}

func workbenchConsoleLevel(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "error", "assert":
		return "error"
	case "warning":
		return "warning"
	case "trace":
		return "debug"
	default:
		return "info"
	}
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

func (recorder *workbenchEventRecorder) append(event WorkbenchPageEvent) {
	if recorder == nil {
		return
	}
	event.Type = strings.TrimSpace(event.Type)
	event.Level = firstNonEmpty(strings.TrimSpace(event.Level), "info")
	event.Message = strings.TrimSpace(event.Message)
	event.URL = strings.TrimSpace(event.URL)
	event.Detail = strings.TrimSpace(event.Detail)
	event.Timestamp = firstNonEmpty(strings.TrimSpace(event.Timestamp), time.Now().Format(time.RFC3339Nano))
	if event.Type == "" || (event.Message == "" && event.URL == "" && event.Detail == "") {
		return
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	if recorder.maxEvents > 0 && len(recorder.events) >= recorder.maxEvents {
		return
	}
	recorder.events = append(recorder.events, event)
}

func (recorder *workbenchEventRecorder) Len() int {
	if recorder == nil {
		return 0
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	return len(recorder.events)
}

func (recorder *workbenchEventRecorder) Since(index int) []WorkbenchPageEvent {
	if recorder == nil {
		return nil
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	if index < 0 {
		index = 0
	}
	if index >= len(recorder.events) {
		return nil
	}
	items := make([]WorkbenchPageEvent, 0, len(recorder.events)-index)
	for _, item := range recorder.events[index:] {
		items = append(items, item)
	}
	return items
}

func workbenchPostProcessNetworkRecords(records []workbenchNetworkRecord, exploreMode string) []workbenchNetworkRecord {
	if len(records) == 0 {
		return nil
	}
	items := make([]workbenchNetworkRecord, 0, len(records))
	for _, record := range records {
		item := record
		if item.RequestSchema == nil {
			item.RequestSchema = inferWorkbenchSchemaFromURLQuery(item.URL)
		}
		if strings.TrimSpace(exploreMode) == "public_html_fallback" && workbenchShouldSuppressNoisyAPI(item) {
			continue
		}
		items = append(items, item)
	}
	return items
}

func inferWorkbenchSchemaFromURLQuery(rawURL string) any {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	values := parsed.Query()
	if len(values) == 0 {
		return nil
	}
	result := map[string]any{}
	for key, items := range values {
		key = strings.TrimSpace(key)
		if key == "" || len(items) == 0 {
			continue
		}
		sample := strings.TrimSpace(items[0])
		switch {
		case sample == "":
			result[key] = "string"
		case strings.EqualFold(sample, "true") || strings.EqualFold(sample, "false"):
			result[key] = "boolean"
		case regexp.MustCompile(`^-?\d+$`).MatchString(sample):
			result[key] = "number"
		case regexp.MustCompile(`^-?\d+\.\d+$`).MatchString(sample):
			result[key] = "number"
		default:
			result[key] = "string"
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func workbenchShouldSuppressNoisyAPI(record workbenchNetworkRecord) bool {
	rawURL := strings.TrimSpace(record.URL)
	if rawURL == "" {
		return false
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	pathTemplate := strings.ToLower(strings.TrimSpace(workbenchPathTemplate(rawURL)))
	method := strings.ToUpper(strings.TrimSpace(record.Method))
	contentType := strings.ToLower(strings.TrimSpace(record.ContentType))

	if host == "" && pathTemplate == "" {
		return false
	}
	if workbenchContainsAny(host, "analytics.", ".analytics.", "doubleclick", "googlesyndication", "google-analytics", "umeng", "growingio", "sensorsdata", "mixpanel", "segment", "hotjar", "clarity", "revive.", "adservice", "adserver", "tracking", "tracker", "beacon") {
		return true
	}
	if workbenchContainsAny(pathTemplate, "/collect", "/track", "/tracking", "/beacon", "/metric", "/metrics", "/report", "/log", "/logs", "/gtr/", "/ads/", "/ad/", "/news/g") {
		if method == http.MethodPost || method == http.MethodGet {
			return true
		}
	}
	if method == http.MethodPost && workbenchContainsAny(contentType, "application/octet-stream", "text/plain") && workbenchContainsAny(host, "analytics", "tracker", "tracking", "beacon") {
		return true
	}
	return false
}

func workbenchContainsAny(value string, fragments ...string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return false
	}
	for _, fragment := range fragments {
		fragment = strings.ToLower(strings.TrimSpace(fragment))
		if fragment != "" && strings.Contains(value, fragment) {
			return true
		}
	}
	return false
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

func buildWorkbenchFallbackShape(rawURL string, reason string) *workbenchPageShape {
	_ = strings.TrimSpace(reason)
	title := firstNonEmpty(
		normalizeWorkbenchRoute(rawURL),
		rawURL,
		"unknown page",
	)
	return &workbenchPageShape{
		Title:       title,
		Breadcrumbs: []string{},
		Forms:       []WorkbenchFormCard{},
		Tables:      []WorkbenchTableCard{},
		Actions:     []WorkbenchActionCard{},
		Links:       []WorkbenchLinkCard{},
	}
}

func buildWorkbenchPageCard(site WorkbenchSiteConfig, runID string, discoveryMode string, rawURL string, shape *workbenchPageShape, observation *PageObservation, observationPath string, rawRecords []workbenchNetworkRecord, rawEvents []WorkbenchPageEvent) WorkbenchPageCard {
	if shape == nil {
		shape = buildWorkbenchFallbackShape(rawURL, "shape is nil")
	}
	inputFields := collectWorkbenchInputFields(observation)
	if len(inputFields) == 0 {
		for _, form := range shape.Forms {
			inputFields = append(inputFields, form.Fields...)
		}
	}
	actions := collectWorkbenchObservedActions(shape.Actions, observation)
	links := collectWorkbenchObservedLinks(shape.Links, observation)
	events := collectWorkbenchPageEvents(rawEvents, rawURL)
	apiHits := buildWorkbenchPageAPIHits(rawRecords)
	linkGroups := buildWorkbenchLinkGroups(links)
	keyElements := buildWorkbenchKeyElements(inputFields, actions, shape.Tables, links)
	textSnippets := collectWorkbenchTextSnippets(shape, observation)
	forms := append([]WorkbenchFormCard{}, shape.Forms...)
	if len(forms) == 0 && len(inputFields) > 0 {
		forms = append(forms, WorkbenchFormCard{
			Name:   "页面输入控件",
			Fields: append([]WorkbenchFieldCard{}, inputFields...),
		})
	}
	route := normalizeWorkbenchRoute(rawURL)
	observationTitle := ""
	if observation != nil {
		observationTitle = strings.TrimSpace(observation.Title)
	}
	title := firstNonEmpty(strings.TrimSpace(shape.Title), observationTitle)
	if title == "" || title == "/" || title == rawURL || title == route {
		if len(textSnippets) > 0 {
			title = textSnippets[0]
		}
	}
	title = firstNonEmpty(title, route, rawURL)
	card := WorkbenchPageCard{
		ID:              fmt.Sprintf("route:%s:%s", site.SiteID, route),
		SiteID:          site.SiteID,
		DiscoveryMode:   discoveryMode,
		URL:             rawURL,
		NormalizedRoute: route,
		Title:           title,
		MenuPath:        append([]string{}, shape.Breadcrumbs...),
		Breadcrumbs:     append([]string{}, shape.Breadcrumbs...),
		Summary:         workbenchPageCardSummary(title, shape, observation, textSnippets, keyElements, linkGroups, events, apiHits),
		Forms:           forms,
		InputFields:     inputFields,
		Tables:          append([]WorkbenchTableCard{}, shape.Tables...),
		Actions:         actions,
		Links:           links,
		Events:          events,
		APIHits:         apiHits,
		CaptureSummary:  buildWorkbenchPageCaptureSummary(observation, rawRecords, events),
		TextSnippets:    textSnippets,
		LinkGroups:      linkGroups,
		KeyElements:     keyElements,
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
	for _, event := range card.Events {
		if event.Type == "download" && card.Risk == "read" {
			card.Risk = "read_download"
		}
	}
	for _, hit := range card.APIHits {
		if hit.Risk == "write_high" || hit.Risk == "critical" {
			card.Risk = "write_high"
			break
		}
		if hit.Risk == "write_low" && card.Risk == "read" {
			card.Risk = "write_low"
		}
		if hit.Risk == "read_download" && card.Risk == "read" {
			card.Risk = "read_download"
		}
	}
	return card
}

func buildWorkbenchPageCaptureSummary(observation *PageObservation, records []workbenchNetworkRecord, events []WorkbenchPageEvent) *WorkbenchPageCaptureSummary {
	if observation == nil && len(records) == 0 && len(events) == 0 {
		return nil
	}
	summary := &WorkbenchPageCaptureSummary{
		FilterRule:      "仅保留 xhr/fetch 元信息，请求阶段自动排除 html/css/js/image/font/media 等静态资源；当前不会在 Playwright 回调里读取 response body。",
		ObservationMode: "默认安全模式：不执行 DOM Evaluate，仅保留标题、截图、页面事件与接口元信息。",
		EventCount:      len(events),
	}
	for _, record := range records {
		summary.NetworkRequestCount++
		if strings.TrimSpace(record.Error) != "" {
			summary.NetworkFailureCount++
		}
		if record.ResponseSchema != nil {
			summary.ReadableResponseCount++
		}
	}
	if observation != nil {
		summary.InteractiveElementCount = len(observation.Elements)
		summary.ContentElementCount = len(observation.ContentElements)
		summary.ObservationErrorCount = len(observation.Errors)
		summary.ObservationSummary = strings.TrimSpace(observation.PageSummary)
		if len(observation.Elements) > 0 || len(observation.ContentElements) > 0 {
			summary.ObservationMode = "页面 load 完成后再做渲染态观察，补抓交互元素、内容块、截图和 DOM 快照。"
		}
	}
	return summary
}

func buildWorkbenchPageAPIHits(records []workbenchNetworkRecord) []WorkbenchPageAPIHit {
	if len(records) == 0 {
		return nil
	}
	items := make([]WorkbenchPageAPIHit, 0, minWorkbenchInt(len(records), 10))
	seen := map[string]struct{}{}
	for _, record := range records {
		resourceType := strings.ToLower(strings.TrimSpace(record.ResourceType))
		if resourceType != "xhr" && resourceType != "fetch" {
			continue
		}
		pathTemplate := workbenchPathTemplate(record.URL)
		method := strings.ToUpper(strings.TrimSpace(record.Method))
		risk := classifyWorkbenchAPIRisk(method, pathTemplate)
		hit := WorkbenchPageAPIHit{
			Method:        method,
			PathTemplate:  pathTemplate,
			URL:           strings.TrimSpace(record.URL),
			Status:        record.Status,
			ContentType:   strings.TrimSpace(record.ContentType),
			ResourceType:  record.ResourceType,
			OperationType: workbenchOperationTypeFromRisk(risk),
			Risk:          risk,
			Error:         workbenchCleanText(record.Error, 160),
		}
		key := strings.Join([]string{hit.Method, hit.PathTemplate, fmt.Sprint(hit.Status), hit.Error}, "|")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, hit)
		if len(items) >= 10 {
			break
		}
	}
	return items
}

func collectWorkbenchPageEvents(rawEvents []WorkbenchPageEvent, rawURL string) []WorkbenchPageEvent {
	if len(rawEvents) == 0 {
		return nil
	}
	pageURL := normalizeWorkbenchExploreURL(rawURL)
	items := make([]WorkbenchPageEvent, 0, minWorkbenchInt(len(rawEvents), 12))
	seen := map[string]struct{}{}
	appendEvent := func(event WorkbenchPageEvent) {
		event.Type = strings.TrimSpace(event.Type)
		event.Level = firstNonEmpty(strings.TrimSpace(event.Level), "info")
		event.Message = workbenchCleanText(event.Message, 240)
		event.URL = strings.TrimSpace(event.URL)
		event.Detail = workbenchCleanText(event.Detail, 140)
		event.Timestamp = strings.TrimSpace(event.Timestamp)
		key := strings.Join([]string{event.Type, event.Level, event.Message, event.URL, event.Detail}, "|")
		if event.Type == "" || key == "" {
			return
		}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		items = append(items, event)
	}

	for _, event := range rawEvents {
		eventURL := normalizeWorkbenchExploreURL(event.URL)
		if eventURL != "" && pageURL != "" {
			if event.Type == "frame_navigated" && eventURL != pageURL {
				continue
			}
			if event.Type == "popup" && !workbenchEventURLRelated(pageURL, eventURL) {
				continue
			}
			if strings.HasPrefix(event.Type, "websocket") && !workbenchEventURLRelated(pageURL, eventURL) {
				continue
			}
		}
		if event.Type == "console" && strings.EqualFold(event.Level, "info") && strings.TrimSpace(event.Message) == "" {
			continue
		}
		appendEvent(event)
		if len(items) >= 12 {
			break
		}
	}
	return items
}

func workbenchEventURLRelated(pageURL string, eventURL string) bool {
	pageURL = normalizeWorkbenchExploreURL(pageURL)
	eventURL = normalizeWorkbenchExploreURL(eventURL)
	if pageURL == "" || eventURL == "" {
		return false
	}
	pageParsed, pageErr := url.Parse(pageURL)
	eventParsed, eventErr := url.Parse(eventURL)
	if pageErr != nil || eventErr != nil {
		return pageURL == eventURL
	}
	if !strings.EqualFold(pageParsed.Hostname(), eventParsed.Hostname()) {
		return false
	}
	return true
}

func collectWorkbenchInputFields(observation *PageObservation) []WorkbenchFieldCard {
	if observation == nil || len(observation.Elements) == 0 {
		return nil
	}
	fields := make([]WorkbenchFieldCard, 0, 8)
	seen := map[string]struct{}{}
	for _, element := range observation.Elements {
		tag := strings.ToLower(strings.TrimSpace(element.Tag))
		role := strings.ToLower(strings.TrimSpace(element.Role))
		typ := strings.ToLower(strings.TrimSpace(element.Type))
		if !(tag == "input" || tag == "textarea" || tag == "select" || role == "textbox" || role == "combobox") {
			continue
		}
		selector := firstNonEmpty(strings.TrimSpace(element.PrimarySelector), firstWorkbenchSelectorCandidate(element.SelectorCandidates))
		key := firstNonEmpty(selector, strings.TrimSpace(element.Name), strings.TrimSpace(element.Label), strings.TrimSpace(element.Placeholder))
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		label := firstNonEmpty(strings.TrimSpace(element.Label), strings.TrimSpace(element.Placeholder), strings.TrimSpace(element.AriaLabel), strings.TrimSpace(element.Text), strings.TrimSpace(element.Name))
		fields = append(fields, WorkbenchFieldCard{
			Name:     firstNonEmpty(strings.TrimSpace(element.Name), typ, "input"),
			Label:    label,
			Selector: selector,
		})
		if len(fields) >= 10 {
			break
		}
	}
	return fields
}

func collectWorkbenchObservedActions(existing []WorkbenchActionCard, observation *PageObservation) []WorkbenchActionCard {
	items := append([]WorkbenchActionCard{}, existing...)
	if observation == nil || len(observation.Elements) == 0 {
		return items
	}
	seen := map[string]struct{}{}
	for _, item := range items {
		key := firstNonEmpty(strings.TrimSpace(item.Selector), strings.TrimSpace(item.Label))
		if key != "" {
			seen[key] = struct{}{}
		}
	}
	for _, element := range observation.Elements {
		tag := strings.ToLower(strings.TrimSpace(element.Tag))
		role := strings.ToLower(strings.TrimSpace(element.Role))
		typ := strings.ToLower(strings.TrimSpace(element.Type))
		if !(tag == "button" || role == "button" || typ == "button" || typ == "submit") {
			continue
		}
		label := firstNonEmpty(strings.TrimSpace(element.Text), strings.TrimSpace(element.Label), strings.TrimSpace(element.AriaLabel), strings.TrimSpace(element.Name))
		if label == "" {
			continue
		}
		selector := firstNonEmpty(strings.TrimSpace(element.PrimarySelector), firstWorkbenchSelectorCandidate(element.SelectorCandidates))
		key := firstNonEmpty(selector, label)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, WorkbenchActionCard{
			Label:    label,
			Kind:     firstNonEmpty(typ, tag, role, "button"),
			Selector: selector,
			Risk:     classifyWorkbenchActionRisk(label),
		})
	}
	return items
}

func collectWorkbenchObservedLinks(existing []WorkbenchLinkCard, observation *PageObservation) []WorkbenchLinkCard {
	items := append([]WorkbenchLinkCard{}, existing...)
	if observation == nil || len(observation.Elements) == 0 {
		return items
	}
	seen := map[string]struct{}{}
	for _, item := range items {
		key := firstNonEmpty(strings.TrimSpace(item.Href), strings.TrimSpace(item.Selector), strings.TrimSpace(item.Text))
		if key != "" {
			seen[key] = struct{}{}
		}
	}
	for _, element := range observation.Elements {
		href := strings.TrimSpace(element.Href)
		if href == "" {
			continue
		}
		text := firstNonEmpty(strings.TrimSpace(element.Text), strings.TrimSpace(element.Label), strings.TrimSpace(element.AriaLabel), href)
		selector := firstNonEmpty(strings.TrimSpace(element.PrimarySelector), firstWorkbenchSelectorCandidate(element.SelectorCandidates))
		key := firstNonEmpty(href, selector, text)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, WorkbenchLinkCard{
			Text:     text,
			Href:     href,
			Selector: selector,
		})
	}
	return items
}

func buildWorkbenchLinkGroups(links []WorkbenchLinkCard) []string {
	if len(links) == 0 {
		return nil
	}
	shortLabels := make([]string, 0, 8)
	serviceLabels := make([]string, 0, 6)
	domainLabels := make([]string, 0, 6)
	shortSeen := map[string]struct{}{}
	serviceSeen := map[string]struct{}{}
	domainSeen := map[string]struct{}{}

	appendUnique := func(target *[]string, seen map[string]struct{}, value string, limit int) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		*target = append(*target, value)
		if limit > 0 && len(*target) > limit {
			*target = (*target)[:limit]
		}
	}

	for _, link := range links {
		label := workbenchCleanText(firstNonEmpty(link.Text, link.Href), 40)
		href := strings.TrimSpace(link.Href)
		if label != "" && len([]rune(label)) >= 2 && len([]rune(label)) <= 10 && !strings.Contains(label, "http") {
			appendUnique(&shortLabels, shortSeen, label, 8)
		}
		if workbenchContainsAny(label, "登录", "注册", "下载", "邮箱", "视频", "直播", "文档", "指南", "api", "sdk", "控制台", "帮助", "客服", "社区") {
			appendUnique(&serviceLabels, serviceSeen, label, 6)
		}
		if href != "" {
			if parsed, err := url.Parse(href); err == nil {
				host := strings.TrimSpace(parsed.Hostname())
				host = strings.TrimPrefix(strings.ToLower(host), "www.")
				if host != "" {
					appendUnique(&domainLabels, domainSeen, host, 5)
				}
			}
		}
	}

	groups := make([]string, 0, 3)
	if len(shortLabels) >= 4 {
		groups = append(groups, "高频入口："+strings.Join(shortLabels[:minWorkbenchInt(len(shortLabels), 6)], "、"))
	}
	if len(serviceLabels) >= 2 {
		groups = append(groups, "功能入口："+strings.Join(serviceLabels[:minWorkbenchInt(len(serviceLabels), 5)], "、"))
	}
	if len(domainLabels) >= 2 {
		groups = append(groups, "链接域名："+strings.Join(domainLabels[:minWorkbenchInt(len(domainLabels), 4)], "、"))
	}
	return groups
}

func buildWorkbenchKeyElements(inputs []WorkbenchFieldCard, actions []WorkbenchActionCard, tables []WorkbenchTableCard, links []WorkbenchLinkCard) []string {
	items := make([]string, 0, 12)
	seen := map[string]struct{}{}
	appendItem := func(prefix string, label string, detail string) {
		label = workbenchCleanText(label, 80)
		detail = workbenchCleanText(detail, 80)
		if label == "" {
			return
		}
		value := prefix + " · " + label
		if detail != "" {
			value += " · " + detail
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}
	for _, field := range inputs {
		appendItem("输入", firstNonEmpty(field.Label, field.Name), field.Selector)
		if len(items) >= 12 {
			return items
		}
	}
	for _, action := range actions {
		appendItem("动作", action.Label, action.Selector)
		if len(items) >= 12 {
			return items
		}
	}
	for _, table := range tables {
		detail := ""
		if len(table.Columns) > 0 {
			detail = strings.Join(table.Columns[:minWorkbenchInt(len(table.Columns), 4)], " / ")
		}
		appendItem("表格", firstNonEmpty(table.Name, table.Selector, "table"), detail)
		if len(items) >= 12 {
			return items
		}
	}
	for _, link := range links {
		appendItem("链接", firstNonEmpty(link.Text, link.Href), link.Href)
		if len(items) >= 12 {
			return items
		}
	}
	return items
}

func collectWorkbenchTextSnippets(shape *workbenchPageShape, observation *PageObservation) []string {
	snippets := make([]string, 0, 8)
	seen := map[string]struct{}{}
	appendSnippet := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" || len([]rune(value)) < 2 {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		snippets = append(snippets, value)
	}
	if observation != nil {
		for _, item := range observation.ContentElements {
			appendSnippet(item.Text)
			if len(snippets) >= 8 {
				return snippets
			}
		}
		for _, element := range observation.Elements {
			appendSnippet(firstNonEmpty(strings.TrimSpace(element.Label), strings.TrimSpace(element.Text), strings.TrimSpace(element.NearText)))
			if len(snippets) >= 8 {
				return snippets
			}
		}
	}
	if shape != nil {
		for _, item := range shape.TextSnippets {
			appendSnippet(item)
			if len(snippets) >= 8 {
				return snippets
			}
		}
	}
	return snippets
}

func firstWorkbenchSelectorCandidate(candidates []string) string {
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" {
			return candidate
		}
	}
	return ""
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
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	if parsed.Path == "/" && parsed.RawQuery == "" {
		parsed.Path = ""
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

func writeWorkbenchExploreStatus(runRoot string, payload map[string]any) {
	runRoot = strings.TrimSpace(runRoot)
	if runRoot == "" || payload == nil {
		return
	}
	path := filepath.Join(runRoot, "status.json")
	if err := writeWorkbenchJSON(path, payload); err != nil {
		log.Printf("workbench explore status write_failed run_root=%s err=%v", runRoot, err)
	}
}

func probeWorkbenchNavigationTargets(context playwright.BrowserContext, currentURL string, candidates []WorkbenchLinkCard, allowedDomains []string, timeoutMS int) []string {
	// 稳定版禁用点击探测。
	// 该函数原先会 context.NewPage()+Goto()+Click()，在 Playwright RPC 通道异常时同样可能卡死。
	// 后续需要恢复时，建议放到独立 browser/context 的 worker 里做，而不是复用主探索 context。
	_ = context
	_ = currentURL
	_ = candidates
	_ = allowedDomains
	_ = timeoutMS
	return nil
}
