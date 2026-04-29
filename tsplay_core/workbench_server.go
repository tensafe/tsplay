package tsplay_core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type workbenchServer struct {
	artifactRoot string
}

func NewWorkbenchAPIHandler(artifactRoot string) http.Handler {
	server := &workbenchServer{
		artifactRoot: strings.TrimSpace(artifactRoot),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/workbench/health", server.handleHealth)
	mux.HandleFunc("/api/workbench/sessions", server.handleSessions)
	mux.HandleFunc("/api/workbench/sessions/", server.handleSessionByName)
	mux.HandleFunc("/api/workbench/providers", server.handleProviders)
	mux.HandleFunc("/api/workbench/providers/", server.handleProviderByID)
	mux.HandleFunc("/api/workbench/sites", server.handleSites)
	mux.HandleFunc("/api/workbench/sites/", server.handleSiteSubroutes)
	mux.HandleFunc("/api/workbench/tasks/plan", server.handleTaskPlan)
	mux.HandleFunc("/api/workbench/tasks/run", server.handleTaskRun)
	mux.HandleFunc("/api/workbench/tasks/repair/auto", server.handleTaskRepairAuto)
	mux.HandleFunc("/api/workbench/tasks/repair", server.handleTaskRepair)

	return withWorkbenchRequestLogging(withWorkbenchCORS(mux))
}

func StartWorkbenchAPIServer(addr string, artifactRoot string) error {
	return http.ListenAndServe(addr, NewWorkbenchAPIHandler(artifactRoot))
}

func withWorkbenchCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withWorkbenchRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		recorder := &workbenchResponseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		duration := time.Since(startedAt)
		if workbenchShouldSkipAccessLog(r.Method, r.URL.Path, recorder.status, duration) {
			return
		}
		log.Printf(
			"workbench api method=%s path=%s status=%d duration_ms=%d",
			r.Method,
			r.URL.Path,
			recorder.status,
			duration.Milliseconds(),
		)
	})
}

func workbenchShouldSkipAccessLog(method string, path string, status int, duration time.Duration) bool {
	if status >= 400 {
		return false
	}
	if strings.ToUpper(strings.TrimSpace(method)) != http.MethodGet {
		return false
	}
	path = strings.TrimSpace(path)
	switch {
	case path == "/api/workbench/health",
		path == "/api/workbench/sessions",
		path == "/api/workbench/providers",
		path == "/api/workbench/sites":
		return true
	case strings.HasSuffix(path, "/pages"),
		strings.HasSuffix(path, "/apis"),
		strings.HasSuffix(path, "/entities"):
		return duration < 250*time.Millisecond
	default:
		return false
	}
}

type workbenchResponseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *workbenchResponseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *workbenchServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		workbenchMethodNotAllowed(w, http.MethodGet)
		return
	}
	writeWorkbenchResponse(w, http.StatusOK, map[string]any{
		"ok":            true,
		"artifact_root": firstNonEmpty(s.artifactRoot, DefaultFlowArtifactRoot),
	})
}

func (s *workbenchServer) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sessions, err := ListFlowSavedSessions(s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusInternalServerError, err)
			return
		}
		items := make([]map[string]any, 0, len(sessions))
		for _, session := range sessions {
			items = append(items, BuildFlowSavedSessionView(session, s.artifactRoot))
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"sessions": items,
		})
	case http.MethodPost:
		var payload struct {
			Name             string `json:"name"`
			StorageState     string `json:"storage_state"`
			StorageStatePath string `json:"storage_state_path"`
			Profile          string `json:"profile"`
			Session          string `json:"session"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode session request: %w", err))
			return
		}
		saved, err := SaveFlowSavedSession(FlowSavedSessionSaveOptions{
			Name:             payload.Name,
			ArtifactRoot:     s.artifactRoot,
			StorageStateJSON: payload.StorageState,
			StorageStatePath: payload.StorageStatePath,
			Profile:          payload.Profile,
			Session:          payload.Session,
		})
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, err)
			return
		}
		log.Printf("workbench session saved name=%s kind=%s", saved.Name, saved.Kind)
		writeWorkbenchResponse(w, http.StatusOK, BuildFlowSavedSessionDetail(*saved, s.artifactRoot))
	default:
		workbenchMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (s *workbenchServer) handleSessionByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		workbenchMethodNotAllowed(w, http.MethodGet)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/api/workbench/sessions/")
	name = strings.TrimSpace(name)
	if name == "" {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("session name is required"))
		return
	}
	session, err := LoadFlowSavedSession(name, s.artifactRoot)
	if err != nil {
		writeWorkbenchError(w, http.StatusNotFound, err)
		return
	}
	writeWorkbenchResponse(w, http.StatusOK, BuildFlowSavedSessionDetail(*session, s.artifactRoot))
}

func (s *workbenchServer) handleProviders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := ListWorkbenchProviderConfigs(s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusInternalServerError, err)
			return
		}
		views := make([]WorkbenchProviderView, 0, len(items))
		for _, item := range items {
			views = append(views, BuildWorkbenchProviderView(item))
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"providers": views,
		})
	case http.MethodPost:
		var payload WorkbenchProviderConfig
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode provider config: %w", err))
			return
		}
		saved, err := SaveWorkbenchProviderConfig(payload, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, err)
			return
		}
		log.Printf("workbench provider saved provider_id=%s type=%s enabled=%v", saved.ProviderID, saved.Type, saved.Enabled)
		writeWorkbenchResponse(w, http.StatusOK, BuildWorkbenchProviderView(*saved))
	default:
		workbenchMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (s *workbenchServer) handleProviderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		workbenchMethodNotAllowed(w, http.MethodGet)
		return
	}
	providerID := strings.TrimPrefix(r.URL.Path, "/api/workbench/providers/")
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("provider_id is required"))
		return
	}
	config, err := LoadWorkbenchProviderConfig(providerID, s.artifactRoot)
	if err != nil {
		writeWorkbenchError(w, http.StatusNotFound, err)
		return
	}
	writeWorkbenchResponse(w, http.StatusOK, BuildWorkbenchProviderView(*config))
}

func (s *workbenchServer) handleSites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := ListWorkbenchSiteConfigs(s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusInternalServerError, err)
			return
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"sites": items,
		})
	case http.MethodPost:
		var payload WorkbenchSiteConfig
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode site config: %w", err))
			return
		}
		saved, err := SaveWorkbenchSiteConfig(payload, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, err)
			return
		}
		log.Printf(
			"workbench site saved site_id=%s start_url=%s session=%s provider=%s domains=%d",
			saved.SiteID,
			saved.StartURL,
			saved.SessionName,
			saved.ProviderID,
			len(saved.AllowedDomains),
		)
		writeWorkbenchResponse(w, http.StatusOK, saved)
	default:
		workbenchMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (s *workbenchServer) handleSiteSubroutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/workbench/sites/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("site_id is required"))
		return
	}
	siteID := parts[0]
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			workbenchMethodNotAllowed(w, http.MethodGet)
			return
		}
		site, err := LoadWorkbenchSiteConfig(siteID, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusNotFound, err)
			return
		}
		writeWorkbenchResponse(w, http.StatusOK, site)
		return
	}

	switch parts[1] {
	case "explore":
		if r.Method != http.MethodPost {
			workbenchMethodNotAllowed(w, http.MethodPost)
			return
		}
		var payload struct {
			Headless  bool `json:"headless"`
			TimeoutMS int  `json:"timeout_ms"`
			MaxPages  int  `json:"max_pages"`
		}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		log.Printf(
			"workbench explore request site=%s headless=%v max_pages=%d timeout_ms=%d",
			siteID,
			payload.Headless,
			payload.MaxPages,
			payload.TimeoutMS,
		)
		result, err := ExploreWorkbenchSite(WorkbenchExploreOptions{
			SiteID:       siteID,
			ArtifactRoot: s.artifactRoot,
			Headless:     payload.Headless,
			TimeoutMS:    payload.TimeoutMS,
			MaxPages:     payload.MaxPages,
		})
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, err)
			return
		}
		log.Printf(
			"workbench explore response site=%s run_id=%s mode=%s pages=%d apis=%d entities=%d",
			siteID,
			result.RunID,
			result.ExploreMode,
			len(result.Pages),
			len(result.APIs),
			len(result.Entities),
		)
		writeWorkbenchResponse(w, http.StatusOK, result)
	case "pages":
		if r.Method != http.MethodGet {
			workbenchMethodNotAllowed(w, http.MethodGet)
			return
		}
		items, err := ListWorkbenchPageCards(siteID, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusNotFound, err)
			return
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{"pages": items})
	case "apis":
		if r.Method != http.MethodGet {
			workbenchMethodNotAllowed(w, http.MethodGet)
			return
		}
		items, err := ListWorkbenchAPICards(siteID, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusNotFound, err)
			return
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{"apis": items})
	case "entities":
		if r.Method != http.MethodGet {
			workbenchMethodNotAllowed(w, http.MethodGet)
			return
		}
		items, err := ListWorkbenchEntityCards(siteID, s.artifactRoot)
		if err != nil {
			writeWorkbenchError(w, http.StatusNotFound, err)
			return
		}
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{"entities": items})
	default:
		writeWorkbenchError(w, http.StatusNotFound, fmt.Errorf("unknown workbench site path"))
	}
}

func (s *workbenchServer) handleTaskPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		workbenchMethodNotAllowed(w, http.MethodPost)
		return
	}
	var payload WorkbenchTaskPlanOptions
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode task plan request: %w", err))
		return
	}
	payload.ArtifactRoot = firstNonEmpty(payload.ArtifactRoot, s.artifactRoot)
	log.Printf(
		"workbench plan request site=%s provider_id=%s intent=%q",
		payload.SiteID,
		payload.ProviderID,
		workbenchCleanText(payload.Intent, 120),
	)
	plan, err := s.buildWorkbenchTaskPlan(payload)
	if err != nil {
		writeWorkbenchError(w, http.StatusBadRequest, err)
		return
	}
	log.Printf(
		"workbench plan response site=%s strategy=%s generation=%s matched_pages=%d matched_apis=%d flow_name=%s requires_confirm=%v warnings=%d",
		plan.SiteID,
		plan.Strategy,
		plan.GenerationMode,
		len(plan.MatchedPages),
		len(plan.MatchedAPIs),
		plan.FlowName,
		plan.RequiresUserConfirm,
		len(plan.Warnings),
	)
	writeWorkbenchResponse(w, http.StatusOK, plan)
}

func (s *workbenchServer) buildWorkbenchTaskPlan(options WorkbenchTaskPlanOptions) (*WorkbenchTaskPlan, error) {
	options.ArtifactRoot = firstNonEmpty(strings.TrimSpace(options.ArtifactRoot), s.artifactRoot)
	plan, err := BuildWorkbenchTaskPlan(options)
	if err != nil {
		return nil, err
	}

	siteProviderID := ""
	if site, loadErr := LoadWorkbenchSiteConfig(options.SiteID, options.ArtifactRoot); loadErr == nil && site != nil {
		siteProviderID = strings.TrimSpace(site.ProviderID)
	}
	providerRequested := strings.TrimSpace(options.ProviderID) != "" || siteProviderID != ""

	providerConfig, providerView, err := s.resolveWorkbenchProvider(options.ProviderID, options.SiteID)
	if err != nil {
		if providerRequested {
			if providerView.ProviderID != "" {
				plan.Provider = &providerView
			}
			plan.GenerationMode = "local_fallback"
			addWorkbenchPlanWarning(plan, fmt.Sprintf("AI provider is not ready: %s. Using the local planner instead.", err.Error()))
		}
		return plan, nil
	}

	aiPlan, runtimeProviderView, modelOutput, err := BuildWorkbenchProviderTaskPlan(options, plan, providerConfig)
	if err != nil {
		plan.Provider = &runtimeProviderView
		plan.ModelOutput = strings.TrimSpace(modelOutput)
		plan.ValidationError = err.Error()
		if plan.Flow != nil {
			plan.GenerationMode = "local_fallback"
			addWorkbenchPlanWarning(plan, fmt.Sprintf("AI-generated flow did not pass validation: %s. Using the local planner instead.", err.Error()))
		} else {
			plan.GenerationMode = "provider_failed"
			addWorkbenchPlanWarning(plan, fmt.Sprintf("AI-generated flow did not pass validation: %s.", err.Error()))
		}
		return plan, nil
	}

	return aiPlan, nil
}

func addWorkbenchPlanWarning(plan *WorkbenchTaskPlan, warning string) {
	if plan == nil {
		return
	}
	warning = strings.TrimSpace(warning)
	if warning == "" {
		return
	}
	for _, existing := range plan.Warnings {
		if existing == warning {
			return
		}
	}
	plan.Warnings = append(plan.Warnings, warning)
}

func (s *workbenchServer) handleTaskRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		workbenchMethodNotAllowed(w, http.MethodPost)
		return
	}
	var payload struct {
		SiteID       string `json:"site_id"`
		ArtifactRoot string `json:"artifact_root,omitempty"`
		Intent       string `json:"intent,omitempty"`
		FlowYAML     string `json:"flow_yaml,omitempty"`
		ProviderID   string `json:"provider_id,omitempty"`
		Headless     *bool  `json:"headless,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode task run request: %w", err))
		return
	}
	log.Printf(
		"workbench run request site=%s flow_yaml=%t headless_override=%t intent=%q",
		payload.SiteID,
		strings.TrimSpace(payload.FlowYAML) != "",
		payload.Headless != nil,
		workbenchCleanText(payload.Intent, 120),
	)

	artifactRoot := firstNonEmpty(strings.TrimSpace(payload.ArtifactRoot), s.artifactRoot)
	flowYAML := strings.TrimSpace(payload.FlowYAML)
	var (
		plan *WorkbenchTaskPlan
		flow *Flow
		err  error
	)
	if flowYAML != "" {
		flow, err = ParseFlow([]byte(flowYAML), "yaml")
		if err != nil {
			log.Printf("workbench run parse_flow_failed site=%s err=%v", payload.SiteID, err)
			writeWorkbenchResponse(w, http.StatusOK, map[string]any{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}
	} else {
		log.Printf("workbench run planning_inline site=%s provider_id=%s", payload.SiteID, payload.ProviderID)
		plan, err = s.buildWorkbenchTaskPlan(WorkbenchTaskPlanOptions{
			SiteID:       payload.SiteID,
			ArtifactRoot: artifactRoot,
			Intent:       payload.Intent,
			ProviderID:   payload.ProviderID,
		})
		if err != nil {
			log.Printf("workbench run planning_failed site=%s err=%v", payload.SiteID, err)
			writeWorkbenchResponse(w, http.StatusOK, map[string]any{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}
		if plan.Flow == nil {
			log.Printf("workbench run no_runnable_flow site=%s strategy=%s reason=%q", payload.SiteID, plan.Strategy, plan.Reason)
			writeWorkbenchResponse(w, http.StatusOK, map[string]any{
				"ok":    false,
				"error": firstNonEmpty(plan.Reason, "planner did not produce a runnable flow"),
				"plan":  plan,
			})
			return
		}
		flow = plan.Flow
		flowYAML = plan.FlowYAML
	}

	if payload.Headless != nil {
		if flow.Browser == nil {
			flow.Browser = &FlowBrowserConfig{}
		}
		flow.Browser.Headless = payload.Headless
	} else if flow.Browser == nil || flow.Browser.Headless == nil {
		headless := true
		if flow.Browser == nil {
			flow.Browser = &FlowBrowserConfig{}
		}
		flow.Browser.Headless = &headless
	}

	security := TrustedFlowSecurityPolicy()
	security.FileInputRoot = artifactRoot
	security.FileOutputRoot = artifactRoot

	result, runErr := RunFlow(flow, FlowRunOptions{
		Headless:     flow.Browser != nil && flow.Browser.Headless != nil && *flow.Browser.Headless,
		Security:     &security,
		ArtifactRoot: artifactRoot,
	})

	if flowYAML == "" {
		if encoded, encodeErr := encodeWorkbenchFlowYAML(flow); encodeErr == nil {
			flowYAML = encoded
		}
	}

	response := map[string]any{
		"ok":        runErr == nil,
		"site_id":   payload.SiteID,
		"intent":    payload.Intent,
		"flow_name": flow.Name,
		"flow_yaml": flowYAML,
		"result":    flowResultForTool(result),
	}
	if plan != nil {
		response["plan"] = plan
	}
	if runErr != nil {
		response["error"] = runErr.Error()
	}
	runID := ""
	stepCount := 0
	runErrorText := ""
	if result != nil {
		runID = result.RunID
		stepCount = len(result.Trace)
	}
	if runErr != nil {
		runErrorText = runErr.Error()
	}
	log.Printf(
		"workbench run response site=%s flow=%s ok=%v run_id=%s trace_steps=%d err=%q",
		payload.SiteID,
		flow.Name,
		runErr == nil,
		runID,
		stepCount,
		runErrorText,
	)
	writeWorkbenchResponse(w, http.StatusOK, response)
}

func (s *workbenchServer) handleTaskRepair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		workbenchMethodNotAllowed(w, http.MethodPost)
		return
	}
	var payload struct {
		ArtifactRoot       string          `json:"artifact_root,omitempty"`
		FlowYAML           string          `json:"flow_yaml,omitempty"`
		RunResult          json.RawMessage `json:"run_result,omitempty"`
		Error              string          `json:"error,omitempty"`
		MaxArtifactExcerpt int             `json:"max_artifact_excerpt,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode task repair request: %w", err))
		return
	}
	log.Printf("workbench repair request artifact_root=%s error=%q", payload.ArtifactRoot, workbenchCleanText(payload.Error, 120))

	repairData, err := s.buildWorkbenchRepairData(workbenchRepairRequest{
		ArtifactRoot:       payload.ArtifactRoot,
		FlowYAML:           payload.FlowYAML,
		RunResult:          payload.RunResult,
		Error:              payload.Error,
		MaxArtifactExcerpt: payload.MaxArtifactExcerpt,
	})
	if err != nil {
		log.Printf("workbench repair build_failed err=%v", err)
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}
	log.Printf(
		"workbench repair response failed_step=%s hints=%d",
		repairData.Context.FailedStepPath,
		len(repairData.Context.RepairHints),
	)
	writeWorkbenchResponse(w, http.StatusOK, map[string]any{
		"ok":      true,
		"context": repairData.Context,
		"repair":  repairData.Repair,
	})
}

func (s *workbenchServer) handleTaskRepairAuto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		workbenchMethodNotAllowed(w, http.MethodPost)
		return
	}
	var payload struct {
		SiteID             string          `json:"site_id,omitempty"`
		ProviderID         string          `json:"provider_id,omitempty"`
		ArtifactRoot       string          `json:"artifact_root,omitempty"`
		FlowYAML           string          `json:"flow_yaml,omitempty"`
		RunResult          json.RawMessage `json:"run_result,omitempty"`
		Error              string          `json:"error,omitempty"`
		MaxArtifactExcerpt int             `json:"max_artifact_excerpt,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeWorkbenchError(w, http.StatusBadRequest, fmt.Errorf("decode auto repair request: %w", err))
		return
	}
	log.Printf(
		"workbench repair_auto request site=%s provider_id=%s error=%q",
		payload.SiteID,
		payload.ProviderID,
		workbenchCleanText(payload.Error, 120),
	)

	repairData, err := s.buildWorkbenchRepairData(workbenchRepairRequest{
		ArtifactRoot:       payload.ArtifactRoot,
		FlowYAML:           payload.FlowYAML,
		RunResult:          payload.RunResult,
		Error:              payload.Error,
		MaxArtifactExcerpt: payload.MaxArtifactExcerpt,
	})
	if err != nil {
		log.Printf("workbench repair_auto build_failed err=%v", err)
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	providerConfig, providerView, err := s.resolveWorkbenchProvider(payload.ProviderID, payload.SiteID)
	if err != nil {
		log.Printf("workbench repair_auto resolve_provider_failed site=%s provider_id=%s err=%v", payload.SiteID, payload.ProviderID, err)
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"ok":       false,
			"error":    err.Error(),
			"provider": providerView,
		})
		return
	}

	modelOutput, runtimeProviderView, err := RunWorkbenchProviderPrompt(
		providerConfig,
		"",
		firstNonEmpty(repairData.Repair.Prompt, repairData.Context.Prompt),
	)
	if err != nil {
		log.Printf("workbench repair_auto provider_failed provider_id=%s err=%v", runtimeProviderView.ProviderID, err)
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"ok":       false,
			"error":    err.Error(),
			"provider": runtimeProviderView,
		})
		return
	}

	repairedFlowYAML := ExtractWorkbenchFlowYAML(modelOutput)
	response := map[string]any{
		"ok":                 true,
		"context":            repairData.Context,
		"repair":             repairData.Repair,
		"provider":           runtimeProviderView,
		"model_output":       modelOutput,
		"repaired_flow_yaml": repairedFlowYAML,
	}
	if _, err := ParseFlow([]byte(repairedFlowYAML), "yaml"); err != nil {
		log.Printf("workbench repair_auto validation_failed provider_id=%s err=%v", runtimeProviderView.ProviderID, err)
		response["ok"] = false
		response["validation_error"] = err.Error()
		writeWorkbenchResponse(w, http.StatusOK, map[string]any{
			"ok":                 false,
			"context":            repairData.Context,
			"repair":             repairData.Repair,
			"provider":           runtimeProviderView,
			"model_output":       modelOutput,
			"repaired_flow_yaml": repairedFlowYAML,
			"validation_error":   err.Error(),
		})
		return
	}
	log.Printf(
		"workbench repair_auto response provider_id=%s model=%s repaired_flow=%t",
		runtimeProviderView.ProviderID,
		runtimeProviderView.ResolvedModel,
		strings.TrimSpace(repairedFlowYAML) != "",
	)
	writeWorkbenchResponse(w, http.StatusOK, response)
}

type workbenchRepairRequest struct {
	ArtifactRoot       string
	FlowYAML           string
	RunResult          json.RawMessage
	Error              string
	MaxArtifactExcerpt int
}

type workbenchRepairData struct {
	Flow    *Flow
	Context *FlowRepairContext
	Repair  *FlowRepairRequest
}

func (s *workbenchServer) buildWorkbenchRepairData(payload workbenchRepairRequest) (*workbenchRepairData, error) {
	flowYAML := strings.TrimSpace(payload.FlowYAML)
	if flowYAML == "" {
		return nil, fmt.Errorf("flow_yaml is required")
	}
	flow, err := ParseFlow([]byte(flowYAML), "yaml")
	if err != nil {
		return nil, err
	}
	if len(payload.RunResult) == 0 || string(payload.RunResult) == "null" {
		return nil, fmt.Errorf("run_result is required")
	}
	artifactRoot := firstNonEmpty(strings.TrimSpace(payload.ArtifactRoot), s.artifactRoot)
	result, runError, err := ParseFlowRunResultForRepair(string(payload.RunResult), "")
	if err != nil {
		return nil, err
	}
	repairContext, err := BuildFlowRepairContext(FlowRepairContextOptions{
		Flow:               flow,
		Result:             result,
		Error:              firstNonEmpty(strings.TrimSpace(payload.Error), runError),
		ArtifactRoot:       artifactRoot,
		MaxArtifactExcerpt: payload.MaxArtifactExcerpt,
	})
	if err != nil {
		return nil, err
	}
	repair, err := BuildFlowRepairRequest(FlowRepairRequestOptions{
		Flow:    flow,
		Context: repairContext,
	})
	if err != nil {
		return nil, err
	}
	return &workbenchRepairData{
		Flow:    flow,
		Context: repairContext,
		Repair:  repair,
	}, nil
}

func (s *workbenchServer) resolveWorkbenchProvider(providerID string, siteID string) (WorkbenchProviderConfig, WorkbenchProviderView, error) {
	loadByID := func(id string) (WorkbenchProviderConfig, WorkbenchProviderView, error) {
		config, err := LoadWorkbenchProviderConfig(id, s.artifactRoot)
		if err != nil {
			return WorkbenchProviderConfig{}, WorkbenchProviderView{}, err
		}
		view := BuildWorkbenchProviderView(*config)
		return *config, view, nil
	}

	providerID = strings.TrimSpace(providerID)
	if providerID != "" {
		return loadByID(providerID)
	}

	siteID = normalizeWorkbenchSiteID(siteID)
	if siteID != "" {
		site, err := LoadWorkbenchSiteConfig(siteID, s.artifactRoot)
		if err == nil && strings.TrimSpace(site.ProviderID) != "" {
			return loadByID(site.ProviderID)
		}
	}

	providers, err := ListWorkbenchProviderConfigs(s.artifactRoot)
	if err != nil {
		return WorkbenchProviderConfig{}, WorkbenchProviderView{}, err
	}
	for _, provider := range providers {
		view := BuildWorkbenchProviderView(provider)
		if provider.Enabled && view.Ready {
			return provider, view, nil
		}
	}
	return WorkbenchProviderConfig{}, WorkbenchProviderView{}, fmt.Errorf("no ready provider found; save a provider or configure OPENAI_API_KEY for codex_auto")
}

func writeWorkbenchResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(status)
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		_, _ = w.Write([]byte(`{"error":"marshal response failed"}`))
		return
	}
	_, _ = w.Write(encoded)
}

func writeWorkbenchError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}
	writeWorkbenchResponse(w, status, map[string]any{
		"ok":    false,
		"error": err.Error(),
	})
}

func workbenchMethodNotAllowed(w http.ResponseWriter, allowed ...string) {
	if len(allowed) > 0 {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
	}
	writeWorkbenchError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
}
