package tsplay_core

type WorkbenchSiteConfig struct {
	SiteID         string   `json:"site_id"`
	Name           string   `json:"name"`
	StartURL       string   `json:"start_url"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	SessionName    string   `json:"session_name,omitempty"`
	ProviderID     string   `json:"provider_id,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

type WorkbenchProviderConfig struct {
	ProviderID   string `json:"provider_id"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type"`
	BaseURL      string `json:"base_url,omitempty"`
	Model        string `json:"model,omitempty"`
	APIKey       string `json:"api_key,omitempty"`
	APIKeyEnv    string `json:"api_key_env,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type WorkbenchProviderView struct {
	ProviderID           string `json:"provider_id"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type"`
	BaseURL              string `json:"base_url,omitempty"`
	Model                string `json:"model,omitempty"`
	APIKeyEnv            string `json:"api_key_env,omitempty"`
	SystemPrompt         string `json:"system_prompt,omitempty"`
	Enabled              bool   `json:"enabled"`
	HasAPIKey            bool   `json:"has_api_key,omitempty"`
	APIKeyMasked         string `json:"api_key_masked,omitempty"`
	ResolvedBaseURL      string `json:"resolved_base_url,omitempty"`
	ResolvedModel        string `json:"resolved_model,omitempty"`
	ResolvedAPIKeySource string `json:"resolved_api_key_source,omitempty"`
	Ready                bool   `json:"ready,omitempty"`
	Status               string `json:"status,omitempty"`
	Error                string `json:"error,omitempty"`
	Detected             bool   `json:"detected,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
}

type WorkbenchExploreOptions struct {
	SiteID         string
	Name           string
	StartURL       string
	AllowedDomains []string
	SessionName    string
	ArtifactRoot   string
	Headless       bool
	TimeoutMS      int
	MaxPages       int
}

type WorkbenchFieldCard struct {
	Name     string `json:"name,omitempty"`
	Label    string `json:"label,omitempty"`
	Selector string `json:"selector,omitempty"`
}

type WorkbenchFormCard struct {
	Name     string               `json:"name,omitempty"`
	Selector string               `json:"selector,omitempty"`
	Fields   []WorkbenchFieldCard `json:"fields,omitempty"`
}

type WorkbenchTableCard struct {
	Name     string   `json:"name,omitempty"`
	Selector string   `json:"selector,omitempty"`
	Columns  []string `json:"columns,omitempty"`
}

type WorkbenchActionCard struct {
	Label    string `json:"label,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Selector string `json:"selector,omitempty"`
	Risk     string `json:"risk,omitempty"`
}

type WorkbenchLinkCard struct {
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	Selector string `json:"selector,omitempty"`
}

type WorkbenchPageCard struct {
	ID              string                `json:"id"`
	SiteID          string                `json:"site_id"`
	URL             string                `json:"url"`
	NormalizedRoute string                `json:"normalized_route,omitempty"`
	Title           string                `json:"title,omitempty"`
	MenuPath        []string              `json:"menu_path,omitempty"`
	Breadcrumbs     []string              `json:"breadcrumbs,omitempty"`
	Summary         string                `json:"summary,omitempty"`
	Forms           []WorkbenchFormCard   `json:"forms,omitempty"`
	Tables          []WorkbenchTableCard  `json:"tables,omitempty"`
	Actions         []WorkbenchActionCard `json:"actions,omitempty"`
	Links           []WorkbenchLinkCard   `json:"links,omitempty"`
	Risk            string                `json:"risk,omitempty"`
	ObservationPath string                `json:"observation_path,omitempty"`
	ScreenshotPath  string                `json:"screenshot_path,omitempty"`
	DOMSnapshotPath string                `json:"dom_snapshot_path,omitempty"`
	ExploreRunID    string                `json:"explore_run_id,omitempty"`
	UpdatedAt       string                `json:"updated_at,omitempty"`
}

type WorkbenchAPICard struct {
	ID             string `json:"id"`
	SiteID         string `json:"site_id"`
	Method         string `json:"method"`
	PathTemplate   string `json:"path_template"`
	SemanticName   string `json:"semantic_name,omitempty"`
	TriggerRoute   string `json:"trigger_route,omitempty"`
	TriggerAction  string `json:"trigger_action,omitempty"`
	OperationType  string `json:"operation_type,omitempty"`
	RequestSchema  any    `json:"request_schema,omitempty"`
	ResponseSchema any    `json:"response_schema,omitempty"`
	Risk           string `json:"risk,omitempty"`
	ResourceType   string `json:"resource_type,omitempty"`
	Status         int    `json:"status,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
	URL            string `json:"url,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type WorkbenchEntityField struct {
	Name  string `json:"name"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`
}

type WorkbenchEntityCard struct {
	ID        string                 `json:"id"`
	SiteID    string                 `json:"site_id"`
	Name      string                 `json:"name"`
	Label     string                 `json:"label,omitempty"`
	Fields    []WorkbenchEntityField `json:"fields,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

type WorkbenchExploreResult struct {
	Site         WorkbenchSiteConfig   `json:"site"`
	RunID        string                `json:"run_id"`
	RunRoot      string                `json:"run_root"`
	StartedAt    string                `json:"started_at,omitempty"`
	FinishedAt   string                `json:"finished_at,omitempty"`
	ExploredURLs []string              `json:"explored_urls,omitempty"`
	Pages        []WorkbenchPageCard   `json:"pages,omitempty"`
	APIs         []WorkbenchAPICard    `json:"apis,omitempty"`
	Entities     []WorkbenchEntityCard `json:"entities,omitempty"`
}

type WorkbenchTaskPlanOptions struct {
	SiteID       string `json:"site_id"`
	ArtifactRoot string `json:"artifact_root,omitempty"`
	Intent       string `json:"intent"`
}

type WorkbenchTaskCandidate struct {
	Kind   string `json:"kind"`
	ID     string `json:"id"`
	Label  string `json:"label,omitempty"`
	URL    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Score  int    `json:"score"`
}

type WorkbenchTaskPlan struct {
	SiteID              string                   `json:"site_id"`
	Intent              string                   `json:"intent"`
	MatchedPages        []WorkbenchTaskCandidate `json:"matched_pages,omitempty"`
	MatchedAPIs         []WorkbenchTaskCandidate `json:"matched_apis,omitempty"`
	Strategy            string                   `json:"strategy,omitempty"`
	Reason              string                   `json:"reason,omitempty"`
	FlowName            string                   `json:"flow_name,omitempty"`
	FlowYAML            string                   `json:"flow_yaml,omitempty"`
	Flow                *Flow                    `json:"flow,omitempty"`
	RequiresUserConfirm bool                     `json:"requires_user_confirm,omitempty"`
}

type workbenchNetworkRecord struct {
	ID              string         `json:"id"`
	URL             string         `json:"url"`
	Method          string         `json:"method"`
	ResourceType    string         `json:"resource_type,omitempty"`
	IsNavigation    bool           `json:"is_navigation,omitempty"`
	Status          int            `json:"status,omitempty"`
	ContentType     string         `json:"content_type,omitempty"`
	RequestSchema   any            `json:"request_schema,omitempty"`
	ResponseSchema  any            `json:"response_schema,omitempty"`
	RequestHeaders  map[string]any `json:"request_headers,omitempty"`
	ResponseHeaders map[string]any `json:"response_headers,omitempty"`
	Error           string         `json:"error,omitempty"`
}
