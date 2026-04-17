package tsplay_core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// Flow is the structured workflow format used by TSPlay.
// It keeps most business logic declarative, while still allowing lua steps as
// an escape hatch for advanced cases.
type Flow struct {
	SchemaVersion string             `json:"schema_version" yaml:"schema_version"`
	Name          string             `json:"name" yaml:"name"`
	Description   string             `json:"description,omitempty" yaml:"description,omitempty"`
	Browser       *FlowBrowserConfig `json:"browser,omitempty" yaml:"browser,omitempty"`
	Vars          map[string]any     `json:"vars,omitempty" yaml:"vars,omitempty"`
	Steps         []FlowStep         `json:"steps" yaml:"steps"`
}

type FlowBrowserConfig struct {
	Headless         *bool         `json:"headless,omitempty" yaml:"headless,omitempty"`
	UseSession       string        `json:"use_session,omitempty" yaml:"use_session,omitempty"`
	StorageState     string        `json:"storage_state,omitempty" yaml:"storage_state,omitempty"`
	StorageStatePath string        `json:"storage_state_path,omitempty" yaml:"storage_state_path,omitempty"`
	LoadStorageState string        `json:"load_storage_state,omitempty" yaml:"load_storage_state,omitempty"`
	SaveStorageState string        `json:"save_storage_state,omitempty" yaml:"save_storage_state,omitempty"`
	Persistent       bool          `json:"persistent,omitempty" yaml:"persistent,omitempty"`
	Profile          string        `json:"profile,omitempty" yaml:"profile,omitempty"`
	Session          string        `json:"session,omitempty" yaml:"session,omitempty"`
	Timeout          int           `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	UserAgent        string        `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	Viewport         *FlowViewport `json:"viewport,omitempty" yaml:"viewport,omitempty"`
}

type FlowViewport struct {
	Width  int `json:"width" yaml:"width"`
	Height int `json:"height" yaml:"height"`
}

type FlowStep struct {
	Name            string         `json:"name,omitempty" yaml:"name,omitempty"`
	Action          string         `json:"action" yaml:"action"`
	Args            []any          `json:"args,omitempty" yaml:"args,omitempty"`
	With            map[string]any `json:"with,omitempty" yaml:"with,omitempty"`
	SaveAs          string         `json:"save_as,omitempty" yaml:"save_as,omitempty"`
	ContinueOnError bool           `json:"continue_on_error,omitempty" yaml:"continue_on_error,omitempty"`
	Steps           []FlowStep     `json:"steps,omitempty" yaml:"steps,omitempty"`
	Condition       *FlowStep      `json:"condition,omitempty" yaml:"condition,omitempty"`
	Then            []FlowStep     `json:"then,omitempty" yaml:"then,omitempty"`
	Else            []FlowStep     `json:"else,omitempty" yaml:"else,omitempty"`
	OnError         []FlowStep     `json:"on_error,omitempty" yaml:"on_error,omitempty"`

	// Common named parameters. Flow accepts both args: [...] and these named
	// fields because named fields are easier for humans and AI to review.
	URL          string   `json:"url,omitempty" yaml:"url,omitempty"`
	Selector     string   `json:"selector,omitempty" yaml:"selector,omitempty"`
	Text         string   `json:"text,omitempty" yaml:"text,omitempty"`
	Value        string   `json:"value,omitempty" yaml:"value,omitempty"`
	Timeout      int      `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Seconds      float64  `json:"seconds,omitempty" yaml:"seconds,omitempty"`
	Path         string   `json:"path,omitempty" yaml:"path,omitempty"`
	Range        string   `json:"range,omitempty" yaml:"range,omitempty"`
	Script       string   `json:"script,omitempty" yaml:"script,omitempty"`
	Code         string   `json:"code,omitempty" yaml:"code,omitempty"`
	Attribute    string   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Sheet        string   `json:"sheet,omitempty" yaml:"sheet,omitempty"`
	Key          string   `json:"key,omitempty" yaml:"key,omitempty"`
	FilePath     string   `json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Files        []string `json:"files,omitempty" yaml:"files,omitempty"`
	SavePath     string   `json:"save_path,omitempty" yaml:"save_path,omitempty"`
	Pattern      string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	From         any      `json:"from,omitempty" yaml:"from,omitempty"`
	Connection   string   `json:"connection,omitempty" yaml:"connection,omitempty"`
	Index        int      `json:"index,omitempty" yaml:"index,omitempty"`
	ContextIndex int      `json:"context_index,omitempty" yaml:"context_index,omitempty"`
	Delta        int      `json:"delta,omitempty" yaml:"delta,omitempty"`
	TTLSeconds   int      `json:"ttl_seconds,omitempty" yaml:"ttl_seconds,omitempty"`
	Times        int      `json:"times,omitempty" yaml:"times,omitempty"`
	IntervalMS   int      `json:"interval_ms,omitempty" yaml:"interval_ms,omitempty"`
	Items        any      `json:"items,omitempty" yaml:"items,omitempty"`
	ItemVar      string   `json:"item_var,omitempty" yaml:"item_var,omitempty"`
	IndexVar     string   `json:"index_var,omitempty" yaml:"index_var,omitempty"`
}

type FlowResult struct {
	Name         string           `json:"name"`
	Vars         map[string]any   `json:"vars,omitempty"`
	Trace        []FlowStepTrace  `json:"trace"`
	ArtifactRoot string           `json:"artifact_root,omitempty"`
	RunID        string           `json:"run_id,omitempty"`
	RunRoot      string           `json:"run_root,omitempty"`
	SessionID    string           `json:"session_id,omitempty"`
	Playwright   *PlaywrightUsage `json:"playwright,omitempty"`
}

type FlowStepTrace struct {
	Index         int                `json:"index"`
	Path          string             `json:"path,omitempty"`
	Attempt       int                `json:"attempt,omitempty"`
	Iteration     int                `json:"iteration,omitempty"`
	Branch        string             `json:"branch,omitempty"`
	Name          string             `json:"name,omitempty"`
	Action        string             `json:"action"`
	Args          any                `json:"args,omitempty"`
	ArgsSummary   string             `json:"args_summary,omitempty"`
	Status        string             `json:"status"`
	SaveAs        string             `json:"save_as,omitempty"`
	Error         string             `json:"error,omitempty"`
	ErrorStack    string             `json:"error_stack,omitempty"`
	Output        any                `json:"output,omitempty"`
	OutputSummary string             `json:"output_summary,omitempty"`
	PageURL       string             `json:"page_url,omitempty"`
	Artifacts     *FlowStepArtifacts `json:"artifacts,omitempty"`
	Condition     *FlowStepTrace     `json:"condition,omitempty"`
	Children      []FlowStepTrace    `json:"children,omitempty"`
	Attempts      []FlowStepTrace    `json:"attempts,omitempty"`
	StartedAt     string             `json:"started_at"`
	FinishedAt    string             `json:"finished_at"`
	DurationMS    int64              `json:"duration_ms"`
}

type FlowStepArtifacts struct {
	Directory       string `json:"directory,omitempty"`
	ScreenshotPath  string `json:"screenshot_path,omitempty"`
	HTMLPath        string `json:"html_path,omitempty"`
	DOMSnapshotPath string `json:"dom_snapshot_path,omitempty"`
	CaptureError    string `json:"capture_error,omitempty"`
}

func (artifacts FlowStepArtifacts) empty() bool {
	return artifacts.Directory == "" &&
		artifacts.ScreenshotPath == "" &&
		artifacts.HTMLPath == "" &&
		artifacts.DOMSnapshotPath == "" &&
		artifacts.CaptureError == ""
}

type FlowContext struct {
	Vars          map[string]any
	Security      *FlowSecurityPolicy
	ArtifactRoot  string
	RunID         string
	RunRoot       string
	Context       context.Context
	DBTransaction *flowDBTransactionScope
	SessionID     string
	ClientName    string
	ClientVersion string
}

type FlowRunOptions struct {
	Headless      bool
	Security      *FlowSecurityPolicy
	ArtifactRoot  string
	Context       context.Context
	RunID         string
	RunRoot       string
	SessionID     string
	ClientName    string
	ClientVersion string
}

type FlowSecurityPolicy struct {
	AllowLua          bool   `json:"allow_lua"`
	AllowJavaScript   bool   `json:"allow_javascript"`
	AllowFileAccess   bool   `json:"allow_file_access"`
	AllowBrowserState bool   `json:"allow_browser_state"`
	AllowHTTP         bool   `json:"allow_http"`
	AllowRedis        bool   `json:"allow_redis"`
	AllowDatabase     bool   `json:"allow_database"`
	FileInputRoot     string `json:"file_input_root,omitempty"`
	FileOutputRoot    string `json:"file_output_root,omitempty"`
}

const CurrentFlowSchemaVersion = "1"
const DefaultFlowArtifactRoot = "artifacts"

type flowActionSpec struct {
	Args       []flowArgSpec
	VarArgName string
}

type flowArgSpec struct {
	Name     string
	Required bool
}

var placeholderPattern = regexp.MustCompile(`^\{\{\s*([^{}]+?)\s*\}\}$`)
var replacePattern = regexp.MustCompile(`\{\{\s*([^{}]+?)\s*\}\}`)
var flowIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

var flowActionSpecs = map[string]flowActionSpec{
	"navigate":              {Args: []flowArgSpec{{Name: "url", Required: true}}},
	"click":                 {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"reload":                {},
	"go_back":               {},
	"go_forward":            {},
	"type_text":             {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "text", Required: true}}},
	"get_text":              {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"extract_text":          {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "timeout"}, {Name: "pattern"}}},
	"set_var":               {Args: []flowArgSpec{{Name: "value", Required: true}}},
	"append_var":            {Args: []flowArgSpec{{Name: "value", Required: true}}},
	"set_value":             {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "value", Required: true}}},
	"select_option":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "value", Required: true}}},
	"hover":                 {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"scroll_to":             {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"wait_for_network_idle": {},
	"wait_for_selector":     {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "timeout"}}},
	"wait_for_text":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "text", Required: true}, {Name: "timeout"}}},
	"sleep":                 {Args: []flowArgSpec{{Name: "seconds", Required: true}}},
	"assert_visible":        {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "timeout"}}},
	"assert_text":           {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "text", Required: true}, {Name: "timeout"}}},
	"retry":                 {},
	"if":                    {},
	"foreach":               {},
	"on_error":              {},
	"wait_until":            {},
	"db_transaction":        {},
	"screenshot":            {Args: []flowArgSpec{{Name: "path", Required: true}}},
	"screenshot_element":    {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "path", Required: true}}},
	"save_html":             {Args: []flowArgSpec{{Name: "path", Required: true}}},
	"read_csv":              {Args: []flowArgSpec{{Name: "file_path", Required: true}}},
	"read_excel":            {Args: []flowArgSpec{{Name: "file_path", Required: true}, {Name: "sheet"}, {Name: "range"}}},
	"write_json":            {Args: []flowArgSpec{{Name: "file_path", Required: true}, {Name: "value", Required: true}}},
	"write_csv":             {Args: []flowArgSpec{{Name: "file_path", Required: true}, {Name: "value", Required: true}, {Name: "headers"}}},
	"accept_alert":          {},
	"dismiss_alert":         {},
	"set_alert_text":        {Args: []flowArgSpec{{Name: "text", Required: true}}},
	"execute_script":        {Args: []flowArgSpec{{Name: "script", Required: true}}},
	"evaluate":              {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "script", Required: true}}},
	"upload_file":           {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "file_path", Required: true}}},
	"upload_multiple_files": {Args: []flowArgSpec{{Name: "selector", Required: true}}, VarArgName: "files"},
	"download_file":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "save_path", Required: true}}},
	"download_url":          {Args: []flowArgSpec{{Name: "url", Required: true}, {Name: "save_path", Required: true}}},
	"get_attribute":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "attribute", Required: true}}},
	"get_html":              {Args: []flowArgSpec{{Name: "selector"}}},
	"get_all_links":         {Args: []flowArgSpec{{Name: "selector"}}},
	"capture_table":         {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"http_request":          {Args: []flowArgSpec{{Name: "url", Required: true}, {Name: "method"}, {Name: "headers"}, {Name: "query"}, {Name: "body"}, {Name: "json"}, {Name: "form"}, {Name: "multipart_files"}, {Name: "multipart_fields"}, {Name: "timeout"}, {Name: "response_as"}, {Name: "use_browser_cookies"}, {Name: "use_browser_referer"}, {Name: "use_browser_user_agent"}, {Name: "save_path"}}},
	"json_extract":          {Args: []flowArgSpec{{Name: "from", Required: true}, {Name: "path", Required: true}}},
	"redis_get":             {Args: []flowArgSpec{{Name: "key", Required: true}, {Name: "connection"}}},
	"redis_set":             {Args: []flowArgSpec{{Name: "key", Required: true}, {Name: "value", Required: true}, {Name: "ttl_seconds"}, {Name: "connection"}}},
	"redis_del":             {Args: []flowArgSpec{{Name: "key", Required: true}, {Name: "connection"}}},
	"redis_incr":            {Args: []flowArgSpec{{Name: "key", Required: true}, {Name: "delta"}, {Name: "connection"}}},
	"db_insert":             {Args: []flowArgSpec{{Name: "table", Required: true}, {Name: "row", Required: true}, {Name: "columns"}, {Name: "connection"}, {Name: "driver"}, {Name: "returning"}, {Name: "timeout"}}},
	"db_insert_many":        {Args: []flowArgSpec{{Name: "table", Required: true}, {Name: "rows", Required: true}, {Name: "columns"}, {Name: "connection"}, {Name: "driver"}, {Name: "returning"}, {Name: "timeout"}}},
	"db_upsert":             {Args: []flowArgSpec{{Name: "table", Required: true}, {Name: "row", Required: true}, {Name: "key_columns", Required: true}, {Name: "columns"}, {Name: "update_columns"}, {Name: "do_nothing"}, {Name: "connection"}, {Name: "driver"}, {Name: "returning"}, {Name: "timeout"}}},
	"db_query":              {Args: []flowArgSpec{{Name: "sql", Required: true}, {Name: "args"}, {Name: "connection"}, {Name: "driver"}, {Name: "timeout"}}},
	"db_query_one":          {Args: []flowArgSpec{{Name: "sql", Required: true}, {Name: "args"}, {Name: "connection"}, {Name: "driver"}, {Name: "timeout"}}},
	"db_execute":            {Args: []flowArgSpec{{Name: "sql", Required: true}, {Name: "args"}, {Name: "connection"}, {Name: "driver"}, {Name: "timeout"}}},
	"find_element":          {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"find_elements":         {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"is_visible":            {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"is_enabled":            {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"is_checked":            {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"is_selected":           {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"is_aria_selected":      {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"new_tab":               {Args: []flowArgSpec{{Name: "url", Required: true}}},
	"close_tab":             {},
	"switch_to_tab":         {Args: []flowArgSpec{{Name: "index", Required: true}}},
	"block_request":         {Args: []flowArgSpec{{Name: "pattern", Required: true}}},
	"get_response":          {Args: []flowArgSpec{{Name: "url", Required: true}}},
	"get_storage_state":     {Args: []flowArgSpec{{Name: "context_index"}}},
	"get_cookies_string":    {Args: []flowArgSpec{{Name: "context_index"}}},
	"lua":                   {Args: []flowArgSpec{{Name: "code", Required: true}}},
}

func DefaultFlowSecurityPolicy() FlowSecurityPolicy {
	return FlowSecurityPolicy{}
}

func TrustedFlowSecurityPolicy() FlowSecurityPolicy {
	return FlowSecurityPolicy{
		AllowLua:          true,
		AllowJavaScript:   true,
		AllowFileAccess:   true,
		AllowBrowserState: true,
		AllowHTTP:         true,
		AllowRedis:        true,
		AllowDatabase:     true,
	}
}

func LoadFlowFile(path string) (*Flow, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	flow, err := ParseFlow(content, strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), "."))
	if err != nil {
		return nil, fmt.Errorf("parse flow %s: %w", path, err)
	}
	return flow, nil
}

func ParseFlow(content []byte, format string) (*Flow, error) {
	var flow Flow
	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(content, &flow); err != nil {
			return nil, err
		}
	default:
		if err := yaml.Unmarshal(content, &flow); err != nil {
			return nil, err
		}
	}
	return &flow, nil
}

func ValidateFlow(flow *Flow) error {
	return ValidateFlowStrict(flow)
}

func ValidateFlowStrict(flow *Flow) error {
	if flow == nil {
		return fmt.Errorf("flow is nil")
	}
	if strings.TrimSpace(flow.SchemaVersion) == "" {
		return fmt.Errorf("flow schema_version is required")
	}
	if flow.SchemaVersion != CurrentFlowSchemaVersion {
		return fmt.Errorf("unsupported flow schema_version %q, expected %q", flow.SchemaVersion, CurrentFlowSchemaVersion)
	}
	if len(flow.Steps) == 0 {
		return fmt.Errorf("flow must contain at least one step")
	}

	knownVars := map[string]any{}
	for name, value := range flow.Vars {
		if !flowIdentifierPattern.MatchString(name) {
			return fmt.Errorf("vars key %q is not a valid variable name", name)
		}
		knownVars[name] = value
	}

	if err := validateFlowBrowserConfig(flow.Browser); err != nil {
		return err
	}

	return validateFlowStepSequence(flow.Steps, knownVars, "")
}

func validateFlowStepSequence(steps []FlowStep, knownVars map[string]any, parentPath string) error {
	for i, step := range steps {
		stepPath := flowStepPath(parentPath, i+1)
		if strings.TrimSpace(step.Action) == "" {
			return fmt.Errorf("step %s action is required", stepPath)
		}
		spec, ok := flowActionSpecs[step.Action]
		if !ok {
			return fmt.Errorf("step %s uses unsupported action %q", stepPath, step.Action)
		}

		if step.SaveAs != "" && !flowIdentifierPattern.MatchString(step.SaveAs) {
			return fmt.Errorf("step %s save_as %q is not a valid variable name", stepPath, step.SaveAs)
		}

		if step.Action == "set_var" {
			if err := validateSetVarFlowStep(stepPath, step, knownVars); err != nil {
				return err
			}
			knownVars[step.SaveAs] = nil
			continue
		}
		if step.Action == "append_var" {
			if err := validateAppendVarFlowStep(stepPath, step, knownVars); err != nil {
				return err
			}
			knownVars[step.SaveAs] = nil
			continue
		}
		if step.Action == "http_request" {
			if err := validateHTTPRequestFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}
		if step.Action == "redis_set" {
			if err := validateRedisSetFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}
		if step.Action == "read_csv" {
			if err := validateReadCSVFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}
		if step.Action == "read_excel" {
			if err := validateReadExcelFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}
		if step.Action == "write_json" {
			if err := validateWriteJSONFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}
		if step.Action == "write_csv" {
			if err := validateWriteCSVFlowStep(stepPath, step, spec, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}

		if isFlowControlAction(step.Action) {
			if err := validateFlowControlStep(stepPath, step, knownVars); err != nil {
				return err
			}
			if step.SaveAs != "" {
				knownVars[step.SaveAs] = nil
			}
			continue
		}

		if len(step.Args) > 0 {
			if err := validateFlowStepArgs(stepPath, step, spec, knownVars); err != nil {
				return err
			}
		} else {
			if err := validateFlowStepNamedParams(stepPath, step, spec, knownVars); err != nil {
				return err
			}
		}

		if step.SaveAs != "" {
			knownVars[step.SaveAs] = nil
		}
	}
	return nil
}

func isFlowControlAction(action string) bool {
	switch action {
	case "retry", "if", "foreach", "on_error", "wait_until", "db_transaction":
		return true
	default:
		return false
	}
}

func flowStepPath(parentPath string, index int) string {
	if parentPath == "" {
		return fmt.Sprint(index)
	}
	return fmt.Sprintf("%s.%d", parentPath, index)
}

func copyKnownVars(knownVars map[string]any) map[string]any {
	copied := map[string]any{}
	for name, value := range knownVars {
		copied[name] = value
	}
	return copied
}

func validateFlowBrowserConfig(browser *FlowBrowserConfig) error {
	if browser == nil {
		return nil
	}

	if browser.UseSession != "" && strings.TrimSpace(browser.UseSession) == "" {
		return fmt.Errorf("browser.use_session cannot be blank")
	}
	loadPath, err := browser.loadStorageStatePath()
	if err != nil {
		return err
	}
	if browser.UseSession != "" {
		if loadPath != "" {
			return fmt.Errorf("browser.use_session cannot be combined with browser.storage_state/load_storage_state")
		}
		if browser.wantsPersistentContext() {
			return fmt.Errorf("browser.use_session cannot be combined with browser.persistent/profile/session")
		}
	}
	if browser.wantsPersistentContext() && loadPath != "" {
		return fmt.Errorf("browser.storage_state/load_storage_state is not supported together with persistent profile/session")
	}
	if browser.Timeout < 0 {
		return fmt.Errorf("browser.timeout must be at least 0")
	}
	if browser.SaveStorageState != "" && strings.TrimSpace(browser.SaveStorageState) == "" {
		return fmt.Errorf("browser.save_storage_state cannot be blank")
	}
	if browser.Viewport != nil {
		if browser.Viewport.Width < 1 {
			return fmt.Errorf("browser.viewport.width must be at least 1")
		}
		if browser.Viewport.Height < 1 {
			return fmt.Errorf("browser.viewport.height must be at least 1")
		}
	}
	if strings.TrimSpace(browser.Profile) == "" && browser.Profile != "" {
		return fmt.Errorf("browser.profile cannot be blank")
	}
	if strings.TrimSpace(browser.Session) == "" && browser.Session != "" {
		return fmt.Errorf("browser.session cannot be blank")
	}
	return nil
}

func validateFlowControlStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	switch step.Action {
	case "retry":
		return validateRetryFlowStep(stepPath, step, knownVars)
	case "if":
		return validateIfFlowStep(stepPath, step, knownVars)
	case "foreach":
		return validateForeachFlowStep(stepPath, step, knownVars)
	case "on_error":
		return validateOnErrorFlowStep(stepPath, step, knownVars)
	case "wait_until":
		return validateWaitUntilFlowStep(stepPath, step, knownVars)
	case "db_transaction":
		return validateDBTransactionFlowStep(stepPath, step, knownVars)
	default:
		return fmt.Errorf("step %s action %q is not a control action", stepPath, step.Action)
	}
}

func validateRetryFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use named fields times, interval_ms, and steps", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{"times": true, "interval_ms": true, "steps": true}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if name == "steps" {
			continue
		}
		if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if len(step.Steps) == 0 {
		return fmt.Errorf("step %s action %q requires nested steps", stepPath, step.Action)
	}

	if value, ok := step.param("times"); ok && len(flowReferences(value)) == 0 {
		times, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "times", err)
		}
		if times < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "times")
		}
	}
	if value, ok := step.param("interval_ms"); ok && len(flowReferences(value)) == 0 {
		intervalMS, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "interval_ms", err)
		}
		if intervalMS < 0 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 0", stepPath, step.Action, "interval_ms")
		}
	}
	return validateFlowStepSequence(step.Steps, knownVars, stepPath)
}

func validateIfFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use condition, then, and else", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{"condition": true, "then": true, "else": true}
	for name := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
	}
	if step.Condition == nil {
		return fmt.Errorf("step %s action %q requires condition", stepPath, step.Action)
	}
	if len(step.Then) == 0 && len(step.Else) == 0 {
		return fmt.Errorf("step %s action %q requires then or else steps", stepPath, step.Action)
	}
	if err := validateFlowStepSequence([]FlowStep{*step.Condition}, copyKnownVars(knownVars), stepPath+".condition"); err != nil {
		return err
	}
	if len(step.Then) > 0 {
		if err := validateFlowStepSequence(step.Then, copyKnownVars(knownVars), stepPath+".then"); err != nil {
			return err
		}
	}
	if len(step.Else) > 0 {
		if err := validateFlowStepSequence(step.Else, copyKnownVars(knownVars), stepPath+".else"); err != nil {
			return err
		}
	}
	return nil
}

func validateForeachFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use items, item_var, index_var, and steps", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{
		"items":               true,
		"item_var":            true,
		"index_var":           true,
		"steps":               true,
		"progress_key":        true,
		"progress_connection": true,
		"progress_value":      true,
	}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if name == "steps" {
			continue
		}
		if name == "progress_value" {
			if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
			continue
		}
		if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if _, ok := step.param("items"); !ok {
		return fmt.Errorf("step %s action %q requires items", stepPath, step.Action)
	}
	itemVar, ok := step.param("item_var")
	if !ok {
		return fmt.Errorf("step %s action %q requires item_var", stepPath, step.Action)
	}
	if !flowIdentifierPattern.MatchString(fmt.Sprint(itemVar)) {
		return fmt.Errorf("step %s action %q item_var %q is not a valid variable name", stepPath, step.Action, itemVar)
	}
	if indexVar, ok := step.param("index_var"); ok && !flowIdentifierPattern.MatchString(fmt.Sprint(indexVar)) {
		return fmt.Errorf("step %s action %q index_var %q is not a valid variable name", stepPath, step.Action, indexVar)
	}
	if _, hasProgressConnection := step.param("progress_connection"); hasProgressConnection {
		if _, hasProgressKey := step.param("progress_key"); !hasProgressKey {
			return fmt.Errorf("step %s action %q progress_connection requires progress_key", stepPath, step.Action)
		}
	}
	if _, hasProgressValue := step.param("progress_value"); hasProgressValue {
		if _, hasProgressKey := step.param("progress_key"); !hasProgressKey {
			return fmt.Errorf("step %s action %q progress_value requires progress_key", stepPath, step.Action)
		}
	}
	if progressKey, ok := step.param("progress_key"); ok {
		if len(flowReferences(progressKey)) == 0 && strings.TrimSpace(fmt.Sprint(progressKey)) == "" {
			return fmt.Errorf("step %s action %q progress_key cannot be blank", stepPath, step.Action)
		}
	}
	if len(step.Steps) == 0 {
		return fmt.Errorf("step %s action %q requires nested steps", stepPath, step.Action)
	}
	localVars := copyKnownVars(knownVars)
	localVars[fmt.Sprint(itemVar)] = nil
	if indexVar, ok := step.param("index_var"); ok {
		localVars[fmt.Sprint(indexVar)] = nil
	}
	return validateFlowStepSequence(step.Steps, localVars, stepPath)
}

func validateOnErrorFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use steps and on_error", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{"steps": true, "on_error": true}
	for name := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
	}
	if len(step.Steps) == 0 {
		return fmt.Errorf("step %s action %q requires steps", stepPath, step.Action)
	}
	if len(step.OnError) == 0 {
		return fmt.Errorf("step %s action %q requires on_error steps", stepPath, step.Action)
	}
	if err := validateFlowStepSequence(step.Steps, copyKnownVars(knownVars), stepPath+".try"); err != nil {
		return err
	}
	handlerVars := copyKnownVars(knownVars)
	handlerVars["last_error"] = nil
	return validateFlowStepSequence(step.OnError, handlerVars, stepPath+".on_error")
}

func validateWaitUntilFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use condition, timeout, and interval_ms", stepPath, step.Action)
	}
	present := step.presentNamedParams()
	allowed := map[string]bool{"condition": true, "timeout": true, "interval_ms": true}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if name == "condition" {
			continue
		}
		if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if step.Condition == nil {
		return fmt.Errorf("step %s action %q requires condition", stepPath, step.Action)
	}
	if value, ok := step.param("timeout"); ok && len(flowReferences(value)) == 0 {
		timeout, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "timeout", err)
		}
		if timeout < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "timeout")
		}
	}
	if value, ok := step.param("interval_ms"); ok && len(flowReferences(value)) == 0 {
		intervalMS, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "interval_ms", err)
		}
		if intervalMS < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "interval_ms")
		}
	}
	return validateFlowStepSequence([]FlowStep{*step.Condition}, copyKnownVars(knownVars), stepPath+".condition")
}

func validateSetVarFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use save_as and value", stepPath, step.Action)
	}
	if strings.TrimSpace(step.SaveAs) == "" {
		return fmt.Errorf("step %s action %q requires save_as", stepPath, step.Action)
	}

	present := step.presentNamedParams()
	allowed := map[string]bool{"value": true}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if _, ok := present["value"]; !ok {
		return fmt.Errorf("step %s action %q requires value", stepPath, step.Action)
	}
	return nil
}

func validateAppendVarFlowStep(stepPath string, step FlowStep, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		return fmt.Errorf("step %s action %q does not support args; use save_as and value", stepPath, step.Action)
	}
	if strings.TrimSpace(step.SaveAs) == "" {
		return fmt.Errorf("step %s action %q requires save_as", stepPath, step.Action)
	}

	present := step.presentNamedParams()
	allowed := map[string]bool{"value": true}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	if _, ok := present["value"]; !ok {
		return fmt.Errorf("step %s action %q requires value", stepPath, step.Action)
	}
	return nil
}

func validateHTTPRequestFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		if err := validateFlowStepArgs(stepPath, step, spec, knownVars); err != nil {
			return err
		}
	} else {
		if err := validateFlowStepNamedParams(stepPath, step, spec, knownVars); err != nil {
			return err
		}
	}

	bodyModes := 0
	for _, name := range []string{"body", "json", "form"} {
		value, ok := step.param(name)
		if !ok || value == nil {
			continue
		}
		bodyModes++
	}
	hasMultipart := false
	for _, name := range []string{"multipart_files", "multipart_fields"} {
		value, ok := step.param(name)
		if !ok || value == nil {
			continue
		}
		if listLen(value) == 0 {
			if typed, ok := value.(map[string]any); ok && len(typed) == 0 {
				continue
			}
		}
		hasMultipart = true
		break
	}
	if hasMultipart {
		bodyModes++
	}
	if bodyModes > 1 {
		return fmt.Errorf("step %s action %q accepts only one of body, json, form, or multipart data", stepPath, step.Action)
	}

	if value, ok := step.param("response_as"); ok && len(flowReferences(value)) == 0 {
		responseAs := strings.ToLower(strings.TrimSpace(fmt.Sprint(value)))
		switch responseAs {
		case "", "auto", "text", "json":
		default:
			return fmt.Errorf("step %s action %q parameter %q must be one of auto, text, or json", stepPath, step.Action, "response_as")
		}
	}

	if value, ok := step.param("timeout"); ok && len(flowReferences(value)) == 0 {
		timeout, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "timeout", err)
		}
		if timeout < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "timeout")
		}
	}
	return nil
}

func validateRedisSetFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	if len(step.Args) > 0 {
		if len(step.presentNamedParams()) > 0 {
			return fmt.Errorf("step %s action %q cannot mix args with named parameters", stepPath, step.Action)
		}
		if len(step.Args) < 2 || len(step.Args) > 4 {
			return fmt.Errorf("step %s action %q expects between 2 and 4 args, got %d", stepPath, step.Action, len(step.Args))
		}
		for i, value := range step.Args {
			name := spec.Args[i].Name
			if name == "value" {
				if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
					return err
				}
				continue
			}
			if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		}
	} else {
		present := step.presentNamedParams()
		allowed := allowedFlowParamNames(spec)
		for name, value := range present {
			if !allowed[name] {
				return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
			}
			if name == "value" {
				if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
					return err
				}
				continue
			}
			if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		}
		for _, required := range []string{"key", "value"} {
			if _, ok := present[required]; !ok {
				return fmt.Errorf("step %s action %q requires %q", stepPath, step.Action, required)
			}
		}
	}

	if value, ok := step.param("ttl_seconds"); ok && len(flowReferences(value)) == 0 {
		ttlSeconds, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "ttl_seconds", err)
		}
		if ttlSeconds < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, "ttl_seconds")
		}
	}
	return nil
}

func validateReadExcelFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	return validateReadTableFlowStep(stepPath, step, spec, knownVars, true)
}

func validateReadCSVFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	return validateReadTableFlowStep(stepPath, step, spec, knownVars, false)
}

func validateReadTableFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any, allowExcelParams bool) error {
	validateRangeValue := func(value any) error {
		if len(flowReferences(value)) > 0 {
			return nil
		}
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("step %s action %q parameter %q must be a string", stepPath, step.Action, "range")
		}
		if _, err := parseExcelRangeSpec(text); err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "range", err)
		}
		return nil
	}
	validateHeadersValue := func(value any) error {
		if err := validateFlowReferences(stepPath, step.Action, "headers", value, knownVars); err != nil {
			return err
		}
		if _, ok := fullPlaceholderExpression(value); ok {
			return nil
		}
		headers, err := stringListValue(resolveStaticFlowValue(value, knownVars))
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "headers", err)
		}
		if len(headers) == 0 {
			return fmt.Errorf("step %s action %q parameter %q must contain at least one header", stepPath, step.Action, "headers")
		}
		return nil
	}
	validateRowNumberField := func(value any) error {
		if err := validateFlowParamValue(stepPath, step.Action, "row_number_field", value, knownVars); err != nil {
			return err
		}
		if len(flowReferences(value)) > 0 {
			return nil
		}
		if strings.TrimSpace(fmt.Sprint(value)) == "" {
			return fmt.Errorf("step %s action %q parameter %q cannot be blank", stepPath, step.Action, "row_number_field")
		}
		return nil
	}
	validatePositiveInt := func(name string, value any) error {
		if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
			return err
		}
		if len(flowReferences(value)) > 0 {
			return nil
		}
		number, err := intParam(value)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, name, err)
		}
		if number < 1 {
			return fmt.Errorf("step %s action %q parameter %q must be at least 1", stepPath, step.Action, name)
		}
		return nil
	}

	if len(step.Args) > 0 {
		if len(step.presentNamedParams()) > 0 {
			return fmt.Errorf("step %s action %q cannot mix args with named parameters", stepPath, step.Action)
		}
		if err := validateFlowStepArgs(stepPath, step, spec, knownVars); err != nil {
			return err
		}
		if allowExcelParams {
			if value, ok := step.param("range"); ok {
				if err := validateRangeValue(value); err != nil {
					return err
				}
			}
		}
		return nil
	}

	present := step.presentNamedParams()
	allowed := map[string]bool{"file_path": true, "start_row": true, "limit": true, "row_number_field": true}
	if allowExcelParams {
		allowed["sheet"] = true
		allowed["range"] = true
		allowed["headers"] = true
	}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		switch name {
		case "headers":
			if err := validateHeadersValue(value); err != nil {
				return err
			}
		case "range":
			if err := validateRangeValue(value); err != nil {
				return err
			}
		case "start_row", "limit":
			if err := validatePositiveInt(name, value); err != nil {
				return err
			}
		case "row_number_field":
			if err := validateRowNumberField(value); err != nil {
				return err
			}
		default:
			if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		}
	}
	if _, ok := present["file_path"]; !ok {
		return fmt.Errorf("step %s action %q requires %q", stepPath, step.Action, "file_path")
	}
	return nil
}

func validateWriteJSONFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	return validateWriteValueFlowStep(stepPath, step, spec, knownVars, false)
}

func validateWriteCSVFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	return validateWriteValueFlowStep(stepPath, step, spec, knownVars, true)
}

func validateWriteValueFlowStep(stepPath string, step FlowStep, spec flowActionSpec, knownVars map[string]any, allowHeaders bool) error {
	validateHeadersValue := func(value any) error {
		if err := validateFlowReferences(stepPath, step.Action, "headers", value, knownVars); err != nil {
			return err
		}
		if _, ok := fullPlaceholderExpression(value); ok {
			return nil
		}
		headers, err := stringListValue(resolveStaticFlowValue(value, knownVars))
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q %w", stepPath, step.Action, "headers", err)
		}
		if len(headers) == 0 {
			return fmt.Errorf("step %s action %q parameter %q must contain at least one header", stepPath, step.Action, "headers")
		}
		return nil
	}

	if len(step.Args) > 0 {
		if len(step.presentNamedParams()) > 0 {
			return fmt.Errorf("step %s action %q cannot mix args with named parameters", stepPath, step.Action)
		}
		if len(step.Args) < 2 || len(step.Args) > len(spec.Args) {
			return fmt.Errorf("step %s action %q expects between 2 and %d args, got %d", stepPath, step.Action, len(spec.Args), len(step.Args))
		}
		for i, value := range step.Args {
			name := spec.Args[i].Name
			if name == "value" {
				if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
					return err
				}
				continue
			}
			if name == "headers" {
				if err := validateHeadersValue(value); err != nil {
					return err
				}
				continue
			}
			if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		}
		return nil
	}

	present := step.presentNamedParams()
	allowed := map[string]bool{"file_path": true, "value": true}
	if allowHeaders {
		allowed["headers"] = true
	}
	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepPath, step.Action, name)
		}
		switch name {
		case "value":
			if err := validateFlowReferences(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		case "headers":
			if err := validateHeadersValue(value); err != nil {
				return err
			}
		default:
			if err := validateFlowParamValue(stepPath, step.Action, name, value, knownVars); err != nil {
				return err
			}
		}
	}
	for _, required := range []string{"file_path", "value"} {
		if _, ok := present[required]; !ok {
			return fmt.Errorf("step %s action %q requires %q", stepPath, step.Action, required)
		}
	}
	return nil
}

func resolveStaticFlowValue(value any, knownVars map[string]any) any {
	if known, ok := resolveKnownPlaceholderValue(value, knownVars); ok {
		return known
	}
	return value
}

func validateFlowStepArgs(stepIndex string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	if len(step.presentNamedParams()) > 0 {
		return fmt.Errorf("step %s action %q cannot mix args with named parameters", stepIndex, step.Action)
	}

	requiredCount := requiredArgCount(spec)
	minCount := requiredCount
	maxCount := len(spec.Args)
	if spec.VarArgName != "" {
		minCount = len(spec.Args) + 1
		maxCount = -1
	}

	if len(step.Args) < minCount {
		return fmt.Errorf("step %s action %q expects at least %d args, got %d", stepIndex, step.Action, minCount, len(step.Args))
	}
	if maxCount >= 0 && len(step.Args) > maxCount {
		return fmt.Errorf("step %s action %q expects at most %d args, got %d", stepIndex, step.Action, maxCount, len(step.Args))
	}

	for i, value := range step.Args {
		name := ""
		if i < len(spec.Args) {
			name = spec.Args[i].Name
		} else {
			name = spec.VarArgName
			if name == "files" {
				name = "file_path"
			}
		}
		if err := validateFlowParamValue(stepIndex, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}
	return nil
}

func validateFlowStepNamedParams(stepIndex string, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	present := step.presentNamedParams()
	allowed := allowedFlowParamNames(spec)

	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %s action %q does not accept parameter %q", stepIndex, step.Action, name)
		}
		if err := validateFlowParamValue(stepIndex, step.Action, name, value, knownVars); err != nil {
			return err
		}
	}

	for _, arg := range spec.Args {
		if !arg.Required {
			continue
		}
		if _, ok := present[arg.Name]; !ok {
			return fmt.Errorf("step %s action %q requires %q", stepIndex, step.Action, arg.Name)
		}
	}
	if spec.VarArgName != "" {
		value, ok := present[spec.VarArgName]
		if !ok || listLen(value) == 0 {
			return fmt.Errorf("step %s action %q requires %q", stepIndex, step.Action, spec.VarArgName)
		}
	}
	return nil
}

func validateFlowParamValue(stepIndex string, action string, name string, value any, knownVars map[string]any) error {
	if name == "" {
		return fmt.Errorf("step %s action %q has too many arguments", stepIndex, action)
	}
	if err := validateFlowReferences(stepIndex, action, name, value, knownVars); err != nil {
		return err
	}
	if err := validateFlowParamType(name, value, knownVars); err != nil {
		return fmt.Errorf("step %s action %q parameter %q %w", stepIndex, action, name, err)
	}
	return nil
}

func validateFlowReferences(stepIndex string, action string, name string, value any, knownVars map[string]any) error {
	for _, ref := range flowReferenceExpressions(value) {
		base, _, err := parseFlowVariableReference(ref)
		if err != nil {
			return fmt.Errorf("step %s action %q parameter %q has invalid placeholder %q: %w", stepIndex, action, name, ref, err)
		}
		if _, ok := knownVars[base]; !ok {
			return fmt.Errorf("step %s action %q parameter %q references unknown variable %q", stepIndex, action, name, base)
		}
	}
	return nil
}

func validateFlowParamType(name string, value any, knownVars map[string]any) error {
	if resolved, ok := resolveKnownPlaceholderValue(value, knownVars); ok {
		value = resolved
	} else if _, ok := fullPlaceholderExpression(value); ok {
		return nil
	}

	switch flowParamType(name) {
	case "string":
		if _, ok := value.(string); ok {
			return nil
		}
		return fmt.Errorf("must be a string")
	case "bool":
		if _, ok := value.(bool); ok {
			return nil
		}
		return fmt.Errorf("must be a boolean")
	case "int":
		if isIntegerValue(value) {
			return nil
		}
		return fmt.Errorf("must be an integer")
	case "number":
		if isNumberValue(value) {
			return nil
		}
		return fmt.Errorf("must be a number")
	case "string_list":
		if isStringListValue(value, knownVars) {
			return nil
		}
		return fmt.Errorf("must be a list of strings")
	case "object":
		if _, ok := value.(map[string]any); ok {
			return nil
		}
		return fmt.Errorf("must be an object")
	case "steps":
		if _, ok := value.([]FlowStep); ok {
			return nil
		}
		return fmt.Errorf("must be a list of flow steps")
	case "condition":
		if _, ok := value.(*FlowStep); ok {
			return nil
		}
		if _, ok := value.(FlowStep); ok {
			return nil
		}
		return fmt.Errorf("must be a flow step")
	case "items":
		return nil
	default:
		return nil
	}
}

func flowParamType(name string) string {
	switch name {
	case "url", "selector", "text", "value", "path", "range", "script", "code", "attribute", "sheet", "key", "connection", "file_path", "save_path", "pattern", "item_var", "index_var", "method", "response_as", "body", "row_number_field", "progress_key", "progress_connection", "table", "driver", "sql":
		return "string"
	case "use_browser_cookies", "use_browser_referer", "use_browser_user_agent", "do_nothing":
		return "bool"
	case "timeout", "index", "context_index", "delta", "ttl_seconds", "times", "interval_ms", "start_row", "limit", "timeout_ms", "timeout_seconds":
		return "int"
	case "seconds":
		return "number"
	case "headers", "query", "form", "multipart_files", "multipart_fields", "row":
		return "object"
	case "files", "columns", "returning", "key_columns", "update_columns":
		return "string_list"
	case "steps":
		return "steps"
	case "items", "rows":
		return "items"
	case "condition":
		return "condition"
	case "from", "json", "progress_value":
		return "any"
	default:
		return ""
	}
}

func ValidateFlowSecurity(flow *Flow, policy FlowSecurityPolicy) error {
	if flow == nil {
		return fmt.Errorf("flow is nil")
	}

	if err := validateFlowStepSequenceSecurity(flow.Steps, policy, ""); err != nil {
		return err
	}
	if err := validateFlowFileAccessRoots(flow, policy); err != nil {
		return err
	}
	if err := validateFlowBrowserConfigSecurity(flow.Browser, policy); err != nil {
		return err
	}
	return nil
}

func validateFlowStepSequenceSecurity(steps []FlowStep, policy FlowSecurityPolicy, parentPath string) error {
	for i, step := range steps {
		stepPath := flowStepPath(parentPath, i+1)
		group := flowActionSecurityGroup(step.Action)
		if group != "" && !flowSecurityPolicyAllows(group, policy) {
			option := flowActionSecurityOption(group)
			return fmt.Errorf("step %s action %q is disabled by security policy; set %s=true only for trusted flows", stepPath, step.Action, option)
		}
		if stepUsesRedisCheckpoint(step) && !policy.AllowRedis {
			return fmt.Errorf("step %s action %q progress checkpoint is disabled by security policy; set allow_redis=true only for trusted flows", stepPath, step.Action)
		}
		if stepRequiresFileAccess(step) && !policy.AllowFileAccess {
			return fmt.Errorf("step %s action %q is disabled by security policy; set allow_file_access=true only for trusted flows", stepPath, step.Action)
		}
		if err := forEachNestedFlowStepSequence(step, stepPath, func(nestedSteps []FlowStep, nestedPath string) error {
			return validateFlowStepSequenceSecurity(nestedSteps, policy, nestedPath)
		}); err != nil {
			return err
		}
	}
	return nil
}

func flowActionSecurityGroup(action string) string {
	switch action {
	case "lua":
		return "lua"
	case "execute_script", "evaluate":
		return "javascript"
	case "http_request":
		return "http"
	case "redis_get", "redis_set", "redis_del", "redis_incr":
		return "redis"
	case "db_insert", "db_insert_many", "db_upsert", "db_query", "db_query_one", "db_execute", "db_transaction":
		return "database"
	case "screenshot", "screenshot_element", "save_html", "read_csv", "read_excel", "write_json", "write_csv", "upload_file", "upload_multiple_files", "download_file", "download_url":
		return "file_access"
	case "get_storage_state", "get_cookies_string":
		return "browser_state"
	default:
		return ""
	}
}

func flowActionSecurityOption(group string) string {
	switch group {
	case "lua":
		return "allow_lua"
	case "javascript":
		return "allow_javascript"
	case "file_access":
		return "allow_file_access"
	case "browser_state":
		return "allow_browser_state"
	case "http":
		return "allow_http"
	case "redis":
		return "allow_redis"
	case "database":
		return "allow_database"
	default:
		return "allow_unsafe"
	}
}

func flowSecurityPolicyAllows(group string, policy FlowSecurityPolicy) bool {
	switch group {
	case "lua":
		return policy.AllowLua
	case "javascript":
		return policy.AllowJavaScript
	case "file_access":
		return policy.AllowFileAccess
	case "browser_state":
		return policy.AllowBrowserState
	case "http":
		return policy.AllowHTTP
	case "redis":
		return policy.AllowRedis
	case "database":
		return policy.AllowDatabase
	default:
		return true
	}
}

func stepUsesRedisCheckpoint(step FlowStep) bool {
	if step.Action != "foreach" {
		return false
	}
	value, ok := step.param("progress_key")
	if !ok {
		return false
	}
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text) != "" || len(flowReferences(value)) > 0
	}
	return value != nil
}

func stepRequiresFileAccess(step FlowStep) bool {
	required := false
	_ = forEachFlowFilePathValue(step, func(_ string, _ flowFilePathRole, value any) error {
		if value == nil {
			return nil
		}
		if typed, ok := value.(string); ok && strings.TrimSpace(typed) == "" {
			return nil
		}
		if typed, ok := value.([]string); ok && len(typed) == 0 {
			return nil
		}
		if typed, ok := value.([]any); ok && len(typed) == 0 {
			return nil
		}
		if typed, ok := value.(map[string]any); ok && len(typed) == 0 {
			return nil
		}
		required = true
		return nil
	})
	return required
}

type flowFilePathRole string

const (
	flowFileInputPath  flowFilePathRole = "input"
	flowFileOutputPath flowFilePathRole = "output"
)

func validateFlowFileAccessRoots(flow *Flow, policy FlowSecurityPolicy) error {
	if !policy.AllowFileAccess {
		return nil
	}
	return validateFlowFileAccessRootsForSteps(flow.Steps, policy, "")
}

func validateFlowBrowserConfigSecurity(browser *FlowBrowserConfig, policy FlowSecurityPolicy) error {
	if browser == nil {
		return nil
	}
	if strings.TrimSpace(policy.FileInputRoot) == "" {
		policy.FileInputRoot = DefaultFlowArtifactRoot
	}
	if strings.TrimSpace(policy.FileOutputRoot) == "" {
		policy.FileOutputRoot = DefaultFlowArtifactRoot
	}
	if browser.UseSession != "" {
		if !policy.AllowBrowserState {
			return fmt.Errorf("flow browser config requires allow_browser_state=true only for trusted flows")
		}
		savedSession, err := LoadFlowSavedSession(browser.UseSession, policy.FileOutputRoot)
		if err != nil {
			return fmt.Errorf("resolve browser.use_session %q: %w", browser.UseSession, err)
		}
		if savedSession.Kind == flowSavedSessionKindStorageState && savedSession.StorageStatePath != "" {
			if err := validateFlowFilePathValue("browser", "browser", "use_session", flowFileInputPath, savedSession.StorageStatePath, policy); err != nil {
				return err
			}
		}
	}
	loadPath, err := browser.loadStorageStatePath()
	if err != nil {
		return err
	}
	savePath := strings.TrimSpace(browser.SaveStorageState)
	if !browser.usesBrowserState() {
		return nil
	}
	if !policy.AllowBrowserState {
		return fmt.Errorf("flow browser config requires allow_browser_state=true only for trusted flows")
	}
	if loadPath != "" {
		if err := validateFlowFilePathValue("browser", "browser", "storage_state", flowFileInputPath, loadPath, policy); err != nil {
			return err
		}
	}
	if savePath != "" {
		if err := validateFlowFilePathValue("browser", "browser", "save_storage_state", flowFileOutputPath, savePath, policy); err != nil {
			return err
		}
	}
	return nil
}

func validateFlowFileAccessRootsForSteps(steps []FlowStep, policy FlowSecurityPolicy, parentPath string) error {
	for i, step := range steps {
		stepPath := flowStepPath(parentPath, i+1)
		if err := validateFlowStepFileAccessRoots(stepPath, step, policy); err != nil {
			return err
		}
		if err := forEachNestedFlowStepSequence(step, stepPath, func(nestedSteps []FlowStep, nestedPath string) error {
			return validateFlowFileAccessRootsForSteps(nestedSteps, policy, nestedPath)
		}); err != nil {
			return err
		}
	}
	return nil
}

func forEachNestedFlowStepSequence(step FlowStep, stepPath string, visit func([]FlowStep, string) error) error {
	if step.Condition != nil {
		if err := visit([]FlowStep{*step.Condition}, stepPath+".condition"); err != nil {
			return err
		}
	}
	if len(step.Steps) > 0 {
		nestedPath := stepPath
		if step.Action == "on_error" {
			nestedPath = stepPath + ".try"
		}
		if err := visit(step.Steps, nestedPath); err != nil {
			return err
		}
	}
	if len(step.Then) > 0 {
		if err := visit(step.Then, stepPath+".then"); err != nil {
			return err
		}
	}
	if len(step.Else) > 0 {
		if err := visit(step.Else, stepPath+".else"); err != nil {
			return err
		}
	}
	if len(step.OnError) > 0 {
		if err := visit(step.OnError, stepPath+".on_error"); err != nil {
			return err
		}
	}
	return nil
}

func validateFlowStepFileAccessRoots(stepIndex string, step FlowStep, policy FlowSecurityPolicy) error {
	return forEachFlowFilePathValue(step, func(name string, role flowFilePathRole, value any) error {
		return validateFlowFilePathValue(stepIndex, step.Action, name, role, value, policy)
	})
}

func validateFlowFilePathValue(stepIndex string, action string, name string, role flowFilePathRole, value any, policy FlowSecurityPolicy) error {
	if flowFilePathValueIsDynamic(value) {
		return nil
	}

	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if err := validateFlowFilePathValue(stepIndex, action, name, role, item, policy); err != nil {
				return err
			}
		}
		return nil
	case []string:
		for _, item := range typed {
			if err := validateFlowFilePathValue(stepIndex, action, name, role, item, policy); err != nil {
				return err
			}
		}
		return nil
	case map[string]any:
		for _, item := range typed {
			if err := validateFlowFilePathValue(stepIndex, action, name, role, item, policy); err != nil {
				return err
			}
		}
		return nil
	case string:
		root := flowFileRootForRole(role, policy)
		if root == "" {
			return nil
		}
		if err := validatePathWithinRoot(typed, root); err != nil {
			return fmt.Errorf("step %s action %q parameter %q is outside allowed file %s root %q: %w", stepIndex, action, name, role, root, err)
		}
	}
	return nil
}

func forEachFlowFilePathValue(step FlowStep, visit func(name string, role flowFilePathRole, value any) error) error {
	if len(step.Args) > 0 {
		return forEachFlowFilePathArgValue(step, visit)
	}

	params := step.presentNamedParams()
	for name, role := range flowFilePathParams(step.Action) {
		value, ok := params[name]
		if !ok {
			continue
		}
		if err := visit(name, role, value); err != nil {
			return err
		}
	}
	return nil
}

func forEachFlowFilePathArgValue(step FlowStep, visit func(name string, role flowFilePathRole, value any) error) error {
	params := flowFilePathParams(step.Action)
	spec := flowActionSpecs[step.Action]
	for i, value := range step.Args {
		name := ""
		if i < len(spec.Args) {
			name = spec.Args[i].Name
		} else {
			name = spec.VarArgName
		}
		role, ok := params[name]
		if !ok {
			continue
		}
		if err := visit(name, role, value); err != nil {
			return err
		}
	}
	return nil
}

func flowFilePathParams(action string) map[string]flowFilePathRole {
	switch action {
	case "screenshot", "save_html":
		return map[string]flowFilePathRole{"path": flowFileOutputPath}
	case "screenshot_element":
		return map[string]flowFilePathRole{"path": flowFileOutputPath}
	case "read_csv", "read_excel":
		return map[string]flowFilePathRole{"file_path": flowFileInputPath}
	case "write_json", "write_csv":
		return map[string]flowFilePathRole{"file_path": flowFileOutputPath}
	case "download_file", "download_url":
		return map[string]flowFilePathRole{"save_path": flowFileOutputPath}
	case "http_request":
		return map[string]flowFilePathRole{"multipart_files": flowFileInputPath, "save_path": flowFileOutputPath}
	case "upload_file":
		return map[string]flowFilePathRole{"file_path": flowFileInputPath}
	case "upload_multiple_files":
		return map[string]flowFilePathRole{"files": flowFileInputPath}
	default:
		return nil
	}
}

func flowFileRootForRole(role flowFilePathRole, policy FlowSecurityPolicy) string {
	switch role {
	case flowFileInputPath:
		return policy.FileInputRoot
	case flowFileOutputPath:
		return policy.FileOutputRoot
	default:
		return ""
	}
}

func flowFilePathValueIsDynamic(value any) bool {
	return len(flowReferences(value)) > 0
}

func validatePathWithinRoot(path string, root string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root: %w", err)
	}
	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootAbs, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	return ensurePathInsideRoot(candidateAbs, rootAbs)
}

func ensurePathInsideRoot(path string, root string) error {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return fmt.Errorf("%q is outside %q", path, root)
	}
	return nil
}

func requiredArgCount(spec flowActionSpec) int {
	count := 0
	for _, arg := range spec.Args {
		if arg.Required {
			count++
		}
	}
	return count
}

func allowedFlowParamNames(spec flowActionSpec) map[string]bool {
	allowed := map[string]bool{}
	for _, arg := range spec.Args {
		allowed[arg.Name] = true
	}
	if spec.VarArgName != "" {
		allowed[spec.VarArgName] = true
	}
	return allowed
}

func FlowActionNames() []string {
	names := make([]string, 0, len(flowActionSpecs))
	for name := range flowActionSpecs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func RunFlow(flow *Flow, options FlowRunOptions) (*FlowResult, error) {
	if err := ValidateFlow(flow); err != nil {
		return nil, err
	}
	if options.Security != nil {
		if err := ValidateFlowSecurity(flow, *options.Security); err != nil {
			return nil, err
		}
	}
	playwrightUsage := AnalyzeFlowPlaywrightUsage(flow)
	needsPlaywright := playwrightUsage.NeedsPlaywright

	browserConfig, err := resolveFlowBrowserConfig(flow, options)
	if err != nil {
		return nil, err
	}

	L := lua.NewState()
	defer L.Close()

	for _, fn := range GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}

	var pw *playwright.Playwright
	var browser playwright.Browser
	var context playwright.BrowserContext
	var page playwright.Page
	closePlaywright := sync.OnceFunc(func() {
		if page != nil {
			_ = page.Close()
		}
		if context != nil {
			_ = context.Close()
		}
		if browser != nil {
			_ = browser.Close()
		}
		if pw != nil {
			_ = pw.Stop()
		}
	})
	defer closePlaywright()
	stopWatcher := func() {}
	if needsPlaywright {
		pw, err = StartPlaywright()
		if err != nil {
			summary := playwrightUsage.Summary(3)
			if summary != "" {
				return nil, fmt.Errorf("flow requires Playwright because %s: %w", summary, err)
			}
			return nil, err
		}
		if browserConfig.wantsPersistentContext() {
			context, page, err = launchPersistentFlowBrowser(pw, browserConfig, flowBrowserStateRoot(options))
			if err != nil {
				return nil, err
			}
		} else {
			browser, context, page, err = launchFlowBrowser(pw, browserConfig, options)
			if err != nil {
				return nil, err
			}
		}
		stopWatcher = watchContextCancel(options.Context, closePlaywright)
		setFlowBrowserGlobals(L, browser, context, page)
	}
	defer stopWatcher()

	return RunFlowInStateWithOptions(L, flow, options)
}

func resolveFlowBrowserConfig(flow *Flow, options FlowRunOptions) (FlowBrowserConfig, error) {
	config := FlowBrowserConfig{}
	if flow != nil && flow.Browser != nil {
		config = *flow.Browser
	}
	if strings.TrimSpace(config.UseSession) != "" {
		actor := FlowSavedSessionAccessInfo{
			SessionID:     options.SessionID,
			ClientName:    options.ClientName,
			ClientVersion: options.ClientVersion,
			RunID:         options.RunID,
		}
		savedConfig, err := ResolveFlowSavedSessionBrowserConfig(config.UseSession, flowBrowserStateRoot(options), actor)
		if err != nil {
			return FlowBrowserConfig{}, fmt.Errorf("resolve browser.use_session %q: %w", config.UseSession, err)
		}
		if _, err := MarkFlowSavedSessionUsed(config.UseSession, flowBrowserStateRoot(options), actor); err != nil {
			return FlowBrowserConfig{}, fmt.Errorf("mark browser.use_session %q as used: %w", config.UseSession, err)
		}
		if savedConfig != nil {
			if savedConfig.StorageState != "" {
				config.StorageState = savedConfig.StorageState
			}
			if savedConfig.StorageStatePath != "" {
				config.StorageStatePath = savedConfig.StorageStatePath
			}
			if savedConfig.LoadStorageState != "" {
				config.LoadStorageState = savedConfig.LoadStorageState
			}
			if savedConfig.Persistent {
				config.Persistent = true
			}
			if savedConfig.Profile != "" {
				config.Profile = savedConfig.Profile
			}
			if savedConfig.Session != "" {
				config.Session = savedConfig.Session
			}
		}
	}
	if config.Headless == nil {
		headless := options.Headless
		config.Headless = &headless
	}
	if _, err := config.loadStorageStatePath(); err != nil {
		return FlowBrowserConfig{}, err
	}
	return config, nil
}

func launchFlowBrowser(pw *playwright.Playwright, config FlowBrowserConfig, options FlowRunOptions) (playwright.Browser, playwright.BrowserContext, playwright.Page, error) {
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(config.headlessValue()),
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not launch browser: %w", err)
	}

	contextOptions := playwright.BrowserNewContextOptions{}
	if err := applyFlowBrowserContextOptions(&contextOptions, config, options.Security); err != nil {
		browser.Close()
		return nil, nil, nil, err
	}
	context, err := browser.NewContext(contextOptions)
	if err != nil {
		browser.Close()
		return nil, nil, nil, fmt.Errorf("could not create browser context: %w", err)
	}
	page, err := context.NewPage()
	if err != nil {
		context.Close()
		browser.Close()
		return nil, nil, nil, fmt.Errorf("could not create page: %w", err)
	}
	applyFlowBrowserTimeouts(context, page, config)
	return browser, context, page, nil
}

func launchPersistentFlowBrowser(pw *playwright.Playwright, config FlowBrowserConfig, stateRoot string) (playwright.BrowserContext, playwright.Page, error) {
	loadPath, err := config.loadStorageStatePath()
	if err != nil {
		return nil, nil, err
	}
	if loadPath != "" {
		return nil, nil, fmt.Errorf("browser storage_state/load_storage_state is not supported together with persistent profile/session")
	}
	userDataDir, err := config.persistentContextDir(stateRoot)
	if err != nil {
		return nil, nil, err
	}
	contextOptions := playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(config.headlessValue()),
	}
	if config.UserAgent != "" {
		contextOptions.UserAgent = playwright.String(config.UserAgent)
	}
	if config.Viewport != nil {
		contextOptions.Viewport = &playwright.Size{Width: config.Viewport.Width, Height: config.Viewport.Height}
	}
	if config.Timeout > 0 {
		contextOptions.Timeout = playwright.Float(float64(config.Timeout))
	}
	context, err := pw.Chromium.LaunchPersistentContext(userDataDir, contextOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch persistent browser context: %w", err)
	}
	var page playwright.Page
	pages := context.Pages()
	if len(pages) > 0 {
		page = pages[0]
	} else {
		page, err = context.NewPage()
		if err != nil {
			context.Close()
			return nil, nil, fmt.Errorf("could not create page: %w", err)
		}
	}
	applyFlowBrowserTimeouts(context, page, config)
	return context, page, nil
}

func applyFlowBrowserContextOptions(options *playwright.BrowserNewContextOptions, config FlowBrowserConfig, security *FlowSecurityPolicy) error {
	if options == nil {
		return nil
	}
	if config.UserAgent != "" {
		options.UserAgent = playwright.String(config.UserAgent)
	}
	if config.Viewport != nil {
		options.Viewport = &playwright.Size{Width: config.Viewport.Width, Height: config.Viewport.Height}
	}
	loadPath, err := config.runtimeLoadStorageStatePath(security)
	if err != nil {
		return err
	}
	if loadPath != "" {
		options.StorageStatePath = playwright.String(loadPath)
	}
	return nil
}

func applyFlowBrowserTimeouts(context playwright.BrowserContext, page playwright.Page, config FlowBrowserConfig) {
	if config.Timeout <= 0 {
		return
	}
	timeout := float64(config.Timeout)
	if context != nil {
		context.SetDefaultTimeout(timeout)
		context.SetDefaultNavigationTimeout(timeout)
	}
	if page != nil {
		page.SetDefaultTimeout(timeout)
		page.SetDefaultNavigationTimeout(timeout)
	}
}

func RunFlowInState(L *lua.LState, flow *Flow) (*FlowResult, error) {
	return runFlowInState(L, flow, FlowRunOptions{})
}

func RunFlowInStateWithOptions(L *lua.LState, flow *Flow, options FlowRunOptions) (*FlowResult, error) {
	return runFlowInState(L, flow, options)
}

func runFlowInState(L *lua.LState, flow *Flow, options FlowRunOptions) (*FlowResult, error) {
	if err := ValidateFlow(flow); err != nil {
		return nil, err
	}
	if options.Security != nil {
		if err := ValidateFlowSecurity(flow, *options.Security); err != nil {
			return nil, err
		}
	}
	ensureFlowActionGlobals(L)

	artifactRoot := flowArtifactRoot(options)
	runID := strings.TrimSpace(options.RunID)
	if runID == "" {
		runID = newFlowRunID(flow)
	}
	runRoot := strings.TrimSpace(options.RunRoot)
	if runRoot == "" && strings.TrimSpace(artifactRoot) != "" {
		if root, err := prepareRuntimeFileRoot(artifactRoot); err == nil {
			runRoot = filepath.Join(root, runID)
		}
	}
	runContext := options.Context
	if runContext == nil {
		runContext = context.Background()
	}
	ctx := &FlowContext{
		Vars:          map[string]any{},
		Security:      options.Security,
		ArtifactRoot:  artifactRoot,
		RunID:         runID,
		RunRoot:       runRoot,
		Context:       runContext,
		SessionID:     options.SessionID,
		ClientName:    options.ClientName,
		ClientVersion: options.ClientVersion,
	}
	for key, value := range flow.Vars {
		ctx.Vars[key] = value
		L.SetGlobal(key, goValueToLua(L, value))
	}

	result := &FlowResult{
		Name:         flow.Name,
		Vars:         ctx.Vars,
		ArtifactRoot: artifactRoot,
		RunID:        runID,
		RunRoot:      runRoot,
		SessionID:    options.SessionID,
	}
	playwrightUsage := AnalyzeFlowPlaywrightUsage(flow)
	if playwrightUsage.NeedsPlaywright {
		usageCopy := playwrightUsage
		result.Playwright = &usageCopy
	}
	traces, err := runFlowStepSequence(L, ctx, flow.Steps, "", 0, 0)
	result.Trace = append(result.Trace, traces...)
	saveErr := saveFlowBrowserStateFromConfig(L, flow, options)
	if err != nil {
		if saveErr != nil {
			return result, fmt.Errorf("%w (also failed to save storage state: %v)", err, saveErr)
		}
		return result, err
	}
	if saveErr != nil {
		return result, saveErr
	}

	return result, nil
}

func ensureFlowActionGlobals(L *lua.LState) {
	if L == nil {
		return
	}
	for _, fn := range GlobalPlayWrightFunc {
		if L.GetGlobal(fn.Name) != lua.LNil {
			continue
		}
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}
}

func runFlowStepSequence(L *lua.LState, ctx *FlowContext, steps []FlowStep, parentPath string, attempt int, iteration int) ([]FlowStepTrace, error) {
	traces := make([]FlowStepTrace, 0, len(steps))
	for i, step := range steps {
		if err := flowRunContextError(ctx); err != nil {
			return traces, err
		}
		stepPath := flowStepPath(parentPath, i+1)
		trace, err := runFlowStepWithTrace(L, ctx, step, i+1, stepPath, attempt, iteration)
		traces = append(traces, trace)
		if err != nil && !step.ContinueOnError {
			return traces, fmt.Errorf("step %s %q failed: %w", stepPath, step.Action, err)
		}
	}
	return traces, nil
}

func runFlowStepWithTrace(L *lua.LState, ctx *FlowContext, step FlowStep, index int, stepPath string, attempt int, iteration int) (FlowStepTrace, error) {
	if err := flowRunContextError(ctx); err != nil {
		return FlowStepTrace{
			Index:      index,
			Path:       stepPath,
			Attempt:    attempt,
			Iteration:  iteration,
			Name:       step.Name,
			Action:     step.Action,
			Status:     "error",
			Error:      err.Error(),
			StartedAt:  time.Now().Format(time.RFC3339Nano),
			FinishedAt: time.Now().Format(time.RFC3339Nano),
		}, err
	}
	traceArgs := traceStepParams(step, ctx)
	trace := FlowStepTrace{
		Index:       index,
		Path:        stepPath,
		Attempt:     attempt,
		Iteration:   iteration,
		Name:        step.Name,
		Action:      step.Action,
		Args:        compactTraceValue(traceArgs, 0),
		ArgsSummary: summarizeTraceValue(traceArgs),
		SaveAs:      step.SaveAs,
		Status:      "running",
		StartedAt:   time.Now().Format(time.RFC3339Nano),
	}

	var output any
	var err error
	switch step.Action {
	case "retry":
		output, trace.Attempts, err = runFlowRetryStep(L, ctx, step, stepPath)
	case "if":
		output, trace.Condition, trace.Children, trace.Branch, err = runFlowIfStep(L, ctx, step, stepPath)
	case "foreach":
		output, trace.Children, err = runFlowForeachStep(L, ctx, step, stepPath)
	case "on_error":
		output, trace.Children, trace.Branch, err = runFlowOnErrorStep(L, ctx, step, stepPath)
	case "wait_until":
		output, trace.Attempts, err = runFlowWaitUntilStep(L, ctx, step, stepPath)
	case "db_transaction":
		output, trace.Children, err = runFlowDBTransactionStep(L, ctx, step, stepPath)
	default:
		output, err = runFlowStep(L, ctx, step)
	}
	trace.FinishedAt = time.Now().Format(time.RFC3339Nano)
	started, _ := time.Parse(time.RFC3339Nano, trace.StartedAt)
	finished, _ := time.Parse(time.RFC3339Nano, trace.FinishedAt)
	trace.DurationMS = finished.Sub(started).Milliseconds()
	trace.PageURL = currentFlowPageURL(L)

	if err != nil {
		trace.Status = "error"
		trace.Error = err.Error()
		trace.ErrorStack = string(debug.Stack())
		artifacts := captureFlowFailureArtifacts(L, ctx, trace)
		if !artifacts.empty() {
			trace.Artifacts = &artifacts
		}
		return trace, err
	}

	trace.Status = "ok"
	trace.Output = compactTraceValue(output, 0)
	trace.OutputSummary = summarizeTraceValue(output)
	if step.SaveAs != "" {
		ctx.Vars[step.SaveAs] = output
		L.SetGlobal(step.SaveAs, goValueToLua(L, output))
	}
	return trace, nil
}

func runFlowRetryStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, []FlowStepTrace, error) {
	times, err := retryTimes(ctx, step)
	if err != nil {
		return nil, nil, err
	}
	if times < 1 {
		return nil, nil, fmt.Errorf("retry times must be at least 1")
	}
	intervalMS, err := retryIntervalMS(ctx, step)
	if err != nil {
		return nil, nil, err
	}
	if intervalMS < 0 {
		return nil, nil, fmt.Errorf("retry interval_ms must be at least 0")
	}

	allAttempts := []FlowStepTrace{}
	var lastErr error
	for attempt := 1; attempt <= times; attempt++ {
		snapshot := snapshotFlowVars(ctx)
		traces, err := runFlowStepSequence(L, ctx, step.Steps, stepPath, attempt, 0)
		allAttempts = append(allAttempts, traces...)
		if err == nil {
			return map[string]any{
				"attempts": attempt,
				"status":   "succeeded",
			}, allAttempts, nil
		}
		restoreFlowVars(L, ctx, snapshot)
		lastErr = err
		if intervalMS > 0 && attempt < times {
			if err := sleepWithFlowContext(ctx, time.Duration(intervalMS)*time.Millisecond); err != nil {
				return nil, allAttempts, err
			}
		}
	}
	return map[string]any{
		"attempts": times,
		"status":   "failed",
	}, allAttempts, fmt.Errorf("retry failed after %d attempts: %w", times, lastErr)
}

func runFlowIfStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, *FlowStepTrace, []FlowStepTrace, string, error) {
	if step.Condition == nil {
		return nil, nil, nil, "", fmt.Errorf("if requires condition")
	}
	conditionTrace, err := runFlowStepWithTrace(L, ctx, *step.Condition, 0, stepPath+".condition", 0, 0)
	conditionResult := err == nil && flowValueTruthy(conditionTrace.Output)
	branch := "else"
	branchSteps := step.Else
	if conditionResult {
		branch = "then"
		branchSteps = step.Then
	}
	if len(branchSteps) == 0 {
		return map[string]any{
			"condition": conditionResult,
			"branch":    branch,
			"status":    "skipped",
		}, &conditionTrace, nil, branch, nil
	}
	children, childErr := runFlowStepSequence(L, ctx, branchSteps, stepPath+"."+branch, 0, 0)
	if childErr != nil {
		return nil, &conditionTrace, children, branch, childErr
	}
	return map[string]any{
		"condition": conditionResult,
		"branch":    branch,
		"status":    "completed",
	}, &conditionTrace, children, branch, nil
}

func runFlowForeachStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, []FlowStepTrace, error) {
	itemsValue, ok := step.param("items")
	if !ok {
		return nil, nil, fmt.Errorf("foreach requires items")
	}
	resolvedItems, err := resolveValue(itemsValue, ctx)
	if err != nil {
		return nil, nil, err
	}
	items, err := toList(resolvedItems)
	if err != nil {
		return nil, nil, fmt.Errorf("foreach items must be a list: %w", err)
	}
	itemVarValue, ok := step.param("item_var")
	if !ok {
		return nil, nil, fmt.Errorf("foreach requires item_var")
	}
	itemVar := fmt.Sprint(itemVarValue)
	indexVar := ""
	if indexVarValue, ok := step.param("index_var"); ok {
		indexVar = fmt.Sprint(indexVarValue)
	}
	checkpoint, err := newFlowForeachCheckpoint(ctx, step)
	if err != nil {
		return nil, nil, err
	}

	itemSnapshot, hadItem := snapshotSingleFlowVar(ctx, itemVar)
	indexSnapshot, hadIndex := snapshotSingleFlowVar(ctx, indexVar)
	defer restoreSingleFlowVar(L, ctx, itemVar, itemSnapshot, hadItem)
	defer restoreSingleFlowVar(L, ctx, indexVar, indexSnapshot, hadIndex)

	children := []FlowStepTrace{}
	for index, item := range items {
		if err := flowRunContextError(ctx); err != nil {
			return nil, children, err
		}
		setFlowVar(L, ctx, itemVar, item)
		if indexVar != "" {
			setFlowVar(L, ctx, indexVar, index+1)
		}
		iteration := index + 1
		traces, err := runFlowStepSequence(L, ctx, step.Steps, fmt.Sprintf("%s[%d]", stepPath, iteration), 0, iteration)
		children = append(children, traces...)
		if err != nil {
			return nil, children, err
		}
		checkpoint.recordSuccess(L, ctx, item, iteration)
	}
	result := map[string]any{
		"iterations": len(items),
		"status":     "completed",
	}
	if summary := checkpoint.summary(); summary != nil {
		result["checkpoint"] = summary
	}
	return result, children, nil
}

type flowForeachCheckpoint struct {
	enabled    bool
	key        string
	connection string
	rawValue   any
	value      any
	writes     int
	skipped    bool
	skipReason string
	lastError  string
}

func newFlowForeachCheckpoint(ctx *FlowContext, step FlowStep) (*flowForeachCheckpoint, error) {
	progressKey, err := flowStepOptionalStringParam(ctx, step, "progress_key")
	if err != nil {
		return nil, err
	}
	progressKey = strings.TrimSpace(progressKey)
	if progressKey == "" {
		return &flowForeachCheckpoint{}, nil
	}
	connection, err := flowStepOptionalStringParam(ctx, step, "progress_connection")
	if err != nil {
		return nil, err
	}
	progressValue, _ := step.param("progress_value")
	return &flowForeachCheckpoint{
		enabled:    true,
		key:        progressKey,
		connection: connection,
		rawValue:   progressValue,
	}, nil
}

func (checkpoint *flowForeachCheckpoint) recordSuccess(_ *lua.LState, ctx *FlowContext, item any, iteration int) {
	if checkpoint == nil || !checkpoint.enabled {
		return
	}
	if !redisConnectionHasConfig(checkpoint.connection) {
		checkpoint.skipped = true
		checkpoint.skipReason = "redis connection not configured"
		return
	}
	value, err := checkpoint.resolveValue(ctx, item, iteration)
	if err != nil {
		checkpoint.lastError = err.Error()
		return
	}
	if _, err := redisSet(checkpoint.key, value, 0, checkpoint.connection); err != nil {
		checkpoint.lastError = err.Error()
		return
	}
	checkpoint.writes++
	checkpoint.value = value
	checkpoint.lastError = ""
}

func (checkpoint *flowForeachCheckpoint) resolveValue(ctx *FlowContext, item any, iteration int) (any, error) {
	if checkpoint == nil || !checkpoint.enabled {
		return nil, nil
	}
	if checkpoint.rawValue != nil {
		return resolveValue(checkpoint.rawValue, ctx)
	}
	if nextRow, ok := flowForeachCheckpointNextRow(item); ok {
		return nextRow, nil
	}
	return iteration + 1, nil
}

func (checkpoint *flowForeachCheckpoint) summary() map[string]any {
	if checkpoint == nil || !checkpoint.enabled {
		return nil
	}
	summary := map[string]any{
		"key":        checkpoint.key,
		"connection": firstNonEmpty(strings.TrimSpace(checkpoint.connection), redisDefaultConnection),
		"writes":     checkpoint.writes,
	}
	if checkpoint.value != nil {
		summary["last_value"] = checkpoint.value
	}
	if checkpoint.skipped {
		summary["status"] = "skipped"
		summary["reason"] = checkpoint.skipReason
		return summary
	}
	if checkpoint.lastError != "" {
		summary["status"] = "error"
		summary["error"] = checkpoint.lastError
		return summary
	}
	summary["status"] = "ok"
	return summary
}

func flowForeachCheckpointNextRow(item any) (int, bool) {
	record, ok := item.(map[string]any)
	if !ok {
		return 0, false
	}
	for _, field := range []string{"source_row", "row_number", "row"} {
		value, ok := record[field]
		if !ok {
			continue
		}
		number, ok := flowValueAsInt(value)
		if !ok || number < 1 {
			continue
		}
		return number + 1, true
	}
	return 0, false
}

func flowValueAsInt(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int8:
		return int(typed), true
	case int16:
		return int(typed), true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case uint:
		return int(typed), true
	case uint8:
		return int(typed), true
	case uint16:
		return int(typed), true
	case uint32:
		return int(typed), true
	case uint64:
		return int(typed), true
	case float32:
		number := int(typed)
		return number, float32(number) == typed
	case float64:
		number := int(typed)
		return number, float64(number) == typed
	case string:
		number, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0, false
		}
		return number, true
	default:
		return 0, false
	}
}

func runFlowOnErrorStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, []FlowStepTrace, string, error) {
	children, err := runFlowStepSequence(L, ctx, step.Steps, stepPath+".try", 0, 0)
	if err == nil {
		return map[string]any{
			"status": "succeeded",
		}, children, "try", nil
	}

	setFlowVar(L, ctx, "last_error", err.Error())
	handlerTraces, handlerErr := runFlowStepSequence(L, ctx, step.OnError, stepPath+".on_error", 0, 0)
	children = append(children, handlerTraces...)
	if handlerErr != nil {
		return nil, children, "on_error", fmt.Errorf("on_error handler failed after original error %q: %w", err.Error(), handlerErr)
	}
	return map[string]any{
		"status": "handled",
		"error":  err.Error(),
	}, children, "on_error", nil
}

func runFlowWaitUntilStep(L *lua.LState, ctx *FlowContext, step FlowStep, stepPath string) (any, []FlowStepTrace, error) {
	if step.Condition == nil {
		return nil, nil, fmt.Errorf("wait_until requires condition")
	}
	timeout, err := waitUntilTimeoutMS(ctx, step)
	if err != nil {
		return nil, nil, err
	}
	intervalMS, err := waitUntilIntervalMS(ctx, step)
	if err != nil {
		return nil, nil, err
	}

	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	attempts := []FlowStepTrace{}
	var lastErr error
	for attempt := 1; ; attempt++ {
		if err := flowRunContextError(ctx); err != nil {
			return nil, attempts, err
		}
		conditionTrace, err := runFlowStepWithTrace(L, ctx, *step.Condition, 0, stepPath+".condition", attempt, 0)
		attempts = append(attempts, conditionTrace)
		if err == nil && flowValueTruthy(conditionTrace.Output) {
			return map[string]any{
				"attempts": attempt,
				"status":   "satisfied",
			}, attempts, nil
		}
		if err != nil {
			lastErr = err
		}
		if time.Now().After(deadline) {
			if lastErr != nil {
				return nil, attempts, fmt.Errorf("wait_until timed out after %dms: %w", timeout, lastErr)
			}
			return nil, attempts, fmt.Errorf("wait_until timed out after %dms", timeout)
		}
		if err := sleepWithFlowContext(ctx, time.Duration(intervalMS)*time.Millisecond); err != nil {
			return nil, attempts, err
		}
	}
}

func flowRunContextError(ctx *FlowContext) error {
	if ctx == nil || ctx.Context == nil {
		return nil
	}
	err := ctx.Context.Err()
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("flow run timed out")
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("flow run canceled")
	default:
		return err
	}
}

func sleepWithFlowContext(ctx *FlowContext, duration time.Duration) error {
	if duration <= 0 {
		return flowRunContextError(ctx)
	}
	if ctx == nil || ctx.Context == nil {
		time.Sleep(duration)
		return nil
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Context.Done():
		return flowRunContextError(ctx)
	}
}

func snapshotFlowVars(ctx *FlowContext) map[string]any {
	snapshot := map[string]any{}
	if ctx == nil {
		return snapshot
	}
	for key, value := range ctx.Vars {
		snapshot[key] = value
	}
	return snapshot
}

func restoreFlowVars(L *lua.LState, ctx *FlowContext, snapshot map[string]any) {
	if ctx == nil {
		return
	}
	for key := range ctx.Vars {
		if _, ok := snapshot[key]; !ok {
			delete(ctx.Vars, key)
			L.SetGlobal(key, lua.LNil)
		}
	}
	for key, value := range snapshot {
		ctx.Vars[key] = value
		L.SetGlobal(key, goValueToLua(L, value))
	}
}

func setFlowVar(L *lua.LState, ctx *FlowContext, name string, value any) {
	if ctx == nil || strings.TrimSpace(name) == "" {
		return
	}
	ctx.Vars[name] = value
	L.SetGlobal(name, goValueToLua(L, value))
}

func snapshotSingleFlowVar(ctx *FlowContext, name string) (any, bool) {
	if ctx == nil || strings.TrimSpace(name) == "" {
		return nil, false
	}
	value, ok := ctx.Vars[name]
	return value, ok
}

func restoreSingleFlowVar(L *lua.LState, ctx *FlowContext, name string, value any, existed bool) {
	if ctx == nil || strings.TrimSpace(name) == "" {
		return
	}
	if !existed {
		delete(ctx.Vars, name)
		L.SetGlobal(name, lua.LNil)
		return
	}
	setFlowVar(L, ctx, name, value)
}

func flowValueTruthy(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed != ""
	case int:
		return typed != 0
	case int8:
		return typed != 0
	case int16:
		return typed != 0
	case int32:
		return typed != 0
	case int64:
		return typed != 0
	case uint:
		return typed != 0
	case uint8:
		return typed != 0
	case uint16:
		return typed != 0
	case uint32:
		return typed != 0
	case uint64:
		return typed != 0
	case float32:
		return typed != 0
	case float64:
		return typed != 0
	case []any:
		return len(typed) > 0
	case []string:
		return len(typed) > 0
	case map[string]any:
		return len(typed) > 0
	default:
		return true
	}
}

func flowArtifactRoot(options FlowRunOptions) string {
	if strings.TrimSpace(options.ArtifactRoot) != "" {
		return options.ArtifactRoot
	}
	if options.Security != nil && strings.TrimSpace(options.Security.FileOutputRoot) != "" {
		return options.Security.FileOutputRoot
	}
	return ""
}

func flowBrowserStateRoot(options FlowRunOptions) string {
	if root := flowArtifactRoot(options); strings.TrimSpace(root) != "" {
		return root
	}
	return DefaultFlowArtifactRoot
}

func (browser FlowBrowserConfig) headlessValue() bool {
	if browser.Headless == nil {
		return false
	}
	return *browser.Headless
}

func (browser FlowBrowserConfig) loadStorageStatePath() (string, error) {
	candidates := []string{}
	for _, value := range []string{browser.StorageState, browser.StorageStatePath, browser.LoadStorageState} {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		candidates = append(candidates, trimmed)
	}
	if len(candidates) == 0 {
		return "", nil
	}
	selected := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate != selected {
			return "", fmt.Errorf("browser.storage_state, browser.storage_state_path, and browser.load_storage_state must point to the same file when combined")
		}
	}
	return selected, nil
}

func (browser FlowBrowserConfig) runtimeLoadStorageStatePath(security *FlowSecurityPolicy) (string, error) {
	path, err := browser.loadStorageStatePath()
	if err != nil || path == "" {
		return path, err
	}
	return resolveFlowBrowserStatePath(path, flowFileInputPath, security)
}

func (browser FlowBrowserConfig) runtimeSaveStorageStatePath(security *FlowSecurityPolicy) (string, error) {
	path := strings.TrimSpace(browser.SaveStorageState)
	if path == "" {
		return "", nil
	}
	return resolveFlowBrowserStatePath(path, flowFileOutputPath, security)
}

func (browser FlowBrowserConfig) usesBrowserState() bool {
	loadPath, _ := browser.loadStorageStatePath()
	return loadPath != "" || strings.TrimSpace(browser.SaveStorageState) != "" || browser.wantsPersistentContext()
}

func (browser FlowBrowserConfig) wantsPersistentContext() bool {
	return browser.Persistent || strings.TrimSpace(browser.Profile) != "" || strings.TrimSpace(browser.Session) != ""
}

func (browser FlowBrowserConfig) persistentContextDir(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		root = DefaultFlowArtifactRoot
	}
	rootReal, err := prepareRuntimeFileRoot(root)
	if err != nil {
		return "", fmt.Errorf("resolve browser state root %q: %w", root, err)
	}
	profile := strings.TrimSpace(browser.Profile)
	if profile == "" {
		profile = "default"
	}
	dir, err := flowSavedSessionProfileDir(rootReal, profile, browser.Session)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create persistent browser profile %q: %w", dir, err)
	}
	return dir, nil
}

func resolveFlowBrowserStatePath(path string, role flowFilePathRole, security *FlowSecurityPolicy) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}
	policy := FlowSecurityPolicy{}
	if security != nil {
		policy = *security
	}
	root := flowFileRootForRole(role, policy)
	if strings.TrimSpace(root) == "" {
		root = DefaultFlowArtifactRoot
		switch role {
		case flowFileInputPath:
			policy.FileInputRoot = root
		case flowFileOutputPath:
			policy.FileOutputRoot = root
		}
	}
	return resolveRuntimeFilePath(path, role, policy)
}

func saveFlowBrowserStateFromConfig(L *lua.LState, flow *Flow, options FlowRunOptions) error {
	if flow == nil || flow.Browser == nil {
		return nil
	}
	path, err := flow.Browser.runtimeSaveStorageStatePath(options.Security)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	context, ok := flowBrowserContextFromState(L)
	if !ok || context == nil {
		return fmt.Errorf("browser.save_storage_state requires an active browser context")
	}
	if _, err := context.StorageState(path); err != nil {
		return fmt.Errorf("save browser storage state to %q: %w", path, err)
	}
	return nil
}

func newFlowRunID(flow *Flow) string {
	name := "flow"
	if flow != nil && strings.TrimSpace(flow.Name) != "" {
		name = flow.Name
	}
	return sanitizeArtifactSegment(name) + "-" + time.Now().Format("20060102-150405.000000000")
}

func traceStepParams(step FlowStep, ctx *FlowContext) any {
	var value any
	if len(step.Args) > 0 {
		value = append([]any(nil), step.Args...)
	} else {
		value = step.presentNamedParams()
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return value
	}
	return resolved
}

func summarizeTraceValue(value any) string {
	encoded, err := json.Marshal(compactTraceValue(value, 0))
	if err != nil {
		text := fmt.Sprint(value)
		if len(text) > 300 {
			return text[:300] + "...(truncated)"
		}
		return text
	}
	text := string(encoded)
	if len(text) > 500 {
		return text[:500] + "...(truncated)"
	}
	return text
}

func compactTraceValue(value any, depth int) any {
	if depth > 4 {
		return fmt.Sprintf("<%T>", value)
	}
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		if len(typed) > 1000 {
			return typed[:1000] + "...(truncated)"
		}
		return typed
	case []any:
		limit := len(typed)
		if limit > 20 {
			limit = 20
		}
		items := make([]any, 0, limit+1)
		for i := 0; i < limit; i++ {
			items = append(items, compactTraceValue(typed[i], depth+1))
		}
		if len(typed) > limit {
			items = append(items, fmt.Sprintf("...(%d more)", len(typed)-limit))
		}
		return items
	case []string:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return compactTraceValue(items, depth)
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if len(keys) > 20 {
			keys = keys[:20]
		}
		result := map[string]any{}
		for _, key := range keys {
			result[key] = compactTraceValue(typed[key], depth+1)
		}
		if len(typed) > len(keys) {
			result["_truncated"] = fmt.Sprintf("%d more keys", len(typed)-len(keys))
		}
		return result
	default:
		return value
	}
}

func captureFlowFailureArtifacts(L *lua.LState, ctx *FlowContext, trace FlowStepTrace) FlowStepArtifacts {
	artifacts := FlowStepArtifacts{}
	if ctx == nil {
		return artifacts
	}

	root := strings.TrimSpace(ctx.RunRoot)
	if root == "" {
		if strings.TrimSpace(ctx.ArtifactRoot) == "" {
			return artifacts
		}
		preparedRoot, err := prepareRuntimeFileRoot(ctx.ArtifactRoot)
		if err != nil {
			artifacts.CaptureError = fmt.Sprintf("prepare artifact root: %v", err)
			return artifacts
		}
		root = filepath.Join(preparedRoot, ctx.RunID)
	}
	stepID := fmt.Sprintf("%02d", trace.Index)
	if strings.Contains(trace.Path, ".") {
		stepID = sanitizeArtifactSegment(trace.Path)
	}
	if trace.Attempt > 0 {
		stepID = fmt.Sprintf("%s-attempt-%02d", stepID, trace.Attempt)
	}
	dir := filepath.Join(root, fmt.Sprintf("%s-%s", stepID, sanitizeArtifactSegment(trace.Action)))
	if err := os.MkdirAll(dir, 0755); err != nil {
		artifacts.CaptureError = fmt.Sprintf("create artifact directory: %v", err)
		return artifacts
	}
	artifacts.Directory = dir

	page, ok := flowPageFromState(L)
	if !ok {
		artifacts.CaptureError = "page is not available"
		return artifacts
	}

	captureErrors := []string{}
	screenshotPath := filepath.Join(dir, "failure.png")
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{Path: playwright.String(screenshotPath)}); err != nil {
		captureErrors = append(captureErrors, fmt.Sprintf("screenshot: %v", err))
	} else {
		artifacts.ScreenshotPath = screenshotPath
	}

	htmlPath := filepath.Join(dir, "page.html")
	if content, err := page.Content(); err != nil {
		captureErrors = append(captureErrors, fmt.Sprintf("html: %v", err))
	} else if err := os.WriteFile(htmlPath, []byte(content), 0644); err != nil {
		captureErrors = append(captureErrors, fmt.Sprintf("html write: %v", err))
	} else {
		artifacts.HTMLPath = htmlPath
	}

	domPath := filepath.Join(dir, "dom_snapshot.json")
	if snapshot, err := ExtractSimplifiedElementWithXPathResult(page); err != nil {
		captureErrors = append(captureErrors, fmt.Sprintf("dom snapshot: %v", err))
	} else if err := os.WriteFile(domPath, []byte(snapshot), 0644); err != nil {
		captureErrors = append(captureErrors, fmt.Sprintf("dom snapshot write: %v", err))
	} else {
		artifacts.DOMSnapshotPath = domPath
	}

	if len(captureErrors) > 0 {
		artifacts.CaptureError = strings.Join(captureErrors, "; ")
	}
	return artifacts
}

func currentFlowPageURL(L *lua.LState) string {
	page, ok := flowPageFromState(L)
	if !ok {
		return ""
	}
	return page.URL()
}

func flowPageFromState(L *lua.LState) (playwright.Page, bool) {
	if L == nil {
		return nil, false
	}
	value := L.GetGlobal("page")
	userData, ok := value.(*lua.LUserData)
	if !ok || userData == nil {
		return nil, false
	}
	page, ok := userData.Value.(playwright.Page)
	return page, ok && page != nil
}

var artifactSegmentPattern = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func sanitizeArtifactSegment(value string) string {
	cleaned := artifactSegmentPattern.ReplaceAllString(strings.TrimSpace(value), "_")
	cleaned = strings.Trim(cleaned, "._-")
	if cleaned == "" {
		return "flow"
	}
	if len(cleaned) > 80 {
		return cleaned[:80]
	}
	return cleaned
}

func runFlowStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	switch step.Action {
	case "extract_text":
		return runFlowExtractTextStep(L, ctx, step)
	case "set_var":
		return runFlowSetVarStep(ctx, step)
	case "append_var":
		return runFlowAppendVarStep(ctx, step)
	case "http_request":
		return runFlowHTTPRequestStep(L, ctx, step)
	case "json_extract":
		return runFlowJSONExtractStep(ctx, step)
	case "db_insert":
		return runFlowDBInsertStep(ctx, step, "")
	case "db_insert_many":
		return runFlowDBInsertManyStep(ctx, step, "")
	case "db_upsert":
		return runFlowDBUpsertStep(ctx, step, "")
	case "db_query":
		return runFlowDBQueryStep(ctx, step, "")
	case "db_query_one":
		return runFlowDBQueryOneStep(ctx, step, "")
	case "db_execute":
		return runFlowDBExecuteStep(ctx, step, "")
	case "read_csv":
		return runFlowReadCSVStep(ctx, step)
	case "read_excel":
		return runFlowReadExcelStep(ctx, step)
	case "assert_visible":
		return runFlowAssertVisibleStep(L, ctx, step)
	case "assert_text":
		return runFlowAssertTextStep(L, ctx, step)
	case "retry", "if", "foreach", "on_error", "wait_until", "db_transaction":
		return nil, fmt.Errorf("control action %q can only be executed by the flow step runner", step.Action)
	}

	if step.Action == "lua" {
		code, err := step.luaCode(ctx)
		if err != nil {
			return nil, err
		}
		return callLuaChunk(L, code)
	}

	args, err := buildActionArgs(L, ctx, step)
	if err != nil {
		return nil, err
	}
	args, err = rewriteFlowFileAccessArgs(step, args, ctx.Security)
	if err != nil {
		return nil, err
	}

	fn := L.GetGlobal(step.Action)
	if fn == lua.LNil {
		return nil, fmt.Errorf("lua action %q is not registered", step.Action)
	}

	top := L.GetTop()
	if err := L.CallByParam(lua.P{Fn: fn, NRet: lua.MultRet, Protect: true}, args...); err != nil {
		L.SetTop(top)
		return nil, err
	}
	return collectReturns(L, top), nil
}

func runFlowAssertVisibleStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	selector, err := flowStepStringParam(ctx, step, "selector")
	if err != nil {
		return nil, err
	}
	timeout, err := flowStepOptionalIntParam(ctx, step, "timeout")
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		if _, err := runFlowStep(L, ctx, FlowStep{Action: "wait_for_selector", Selector: selector, Timeout: timeout}); err != nil {
			return nil, err
		}
	}

	visible, err := runFlowStep(L, ctx, FlowStep{Action: "is_visible", Selector: selector})
	if err != nil {
		return nil, err
	}
	if ok, _ := visible.(bool); !ok {
		return nil, fmt.Errorf("assert_visible failed: selector %q is not visible", selector)
	}
	return true, nil
}

func runFlowAssertTextStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	selector, err := flowStepStringParam(ctx, step, "selector")
	if err != nil {
		return nil, err
	}
	expectedText, err := flowStepStringParam(ctx, step, "text")
	if err != nil {
		return nil, err
	}
	timeout, err := flowStepOptionalIntParam(ctx, step, "timeout")
	if err != nil {
		return nil, err
	}

	deadline := time.Now()
	if timeout > 0 {
		deadline = deadline.Add(time.Duration(timeout) * time.Millisecond)
	}
	for {
		actual, err := runFlowStep(L, ctx, FlowStep{Action: "get_text", Selector: selector})
		if err == nil && flowTextContains(actual, expectedText) {
			return map[string]any{
				"selector": selector,
				"text":     expectedText,
				"actual":   compactTraceValue(actual, 0),
			}, nil
		}
		if timeout <= 0 || time.Now().After(deadline) {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("assert_text failed: selector %q text does not contain %q; actual=%s", selector, expectedText, summarizeTraceValue(actual))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func runFlowExtractTextStep(L *lua.LState, ctx *FlowContext, step FlowStep) (any, error) {
	selector, err := flowStepStringParam(ctx, step, "selector")
	if err != nil {
		return nil, err
	}
	timeout, err := flowStepOptionalIntParam(ctx, step, "timeout")
	if err != nil {
		return nil, err
	}
	pattern, err := flowStepOptionalStringParam(ctx, step, "pattern")
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		if _, err := runFlowStep(L, ctx, FlowStep{Action: "wait_for_selector", Selector: selector, Timeout: timeout}); err != nil {
			return nil, err
		}
	}

	actual, err := runFlowStep(L, ctx, FlowStep{Action: "get_text", Selector: selector})
	if err != nil {
		return nil, err
	}
	texts := flowTextValues(actual)
	if pattern == "" {
		return compactExtractedTexts(texts), nil
	}
	matched, err := extractTextByPattern(texts, pattern)
	if err != nil {
		return nil, fmt.Errorf("extract_text selector %q pattern %q %w", selector, pattern, err)
	}
	return matched, nil
}

func runFlowSetVarStep(ctx *FlowContext, step FlowStep) (any, error) {
	if strings.TrimSpace(step.SaveAs) == "" {
		return nil, fmt.Errorf("set_var requires save_as")
	}
	value, ok := step.param("value")
	if !ok {
		return nil, fmt.Errorf("set_var requires value")
	}
	return resolveValue(value, ctx)
}

func runFlowAppendVarStep(ctx *FlowContext, step FlowStep) (any, error) {
	if strings.TrimSpace(step.SaveAs) == "" {
		return nil, fmt.Errorf("append_var requires save_as")
	}
	value, ok := step.param("value")
	if !ok {
		return nil, fmt.Errorf("append_var requires value")
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return nil, err
	}
	current, ok := ctx.Vars[step.SaveAs]
	if !ok || current == nil {
		return []any{resolved}, nil
	}
	items, err := toList(current)
	if err != nil {
		return nil, fmt.Errorf("append_var save_as %q must already be a list, got %T", step.SaveAs, current)
	}
	items = append(append([]any(nil), items...), resolved)
	return items, nil
}

func runFlowReadCSVStep(ctx *FlowContext, step FlowStep) (any, error) {
	filePath, err := flowStepStringParam(ctx, step, "file_path")
	if err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Security != nil {
		filePath, err = resolveRuntimeFilePath(filePath, flowFileInputPath, *ctx.Security)
		if err != nil {
			return nil, fmt.Errorf("action %q parameter %q %w", step.Action, "file_path", err)
		}
	}
	readOptions, err := flowStepTableReadOptions(ctx, step)
	if err != nil {
		return nil, err
	}
	return loadCSVRows(filePath, readOptions)
}

func runFlowReadExcelStep(ctx *FlowContext, step FlowStep) (any, error) {
	filePath, err := flowStepStringParam(ctx, step, "file_path")
	if err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Security != nil {
		filePath, err = resolveRuntimeFilePath(filePath, flowFileInputPath, *ctx.Security)
		if err != nil {
			return nil, fmt.Errorf("action %q parameter %q %w", step.Action, "file_path", err)
		}
	}
	sheet, err := flowStepOptionalStringParam(ctx, step, "sheet")
	if err != nil {
		return nil, err
	}
	rangeSpec, err := flowStepOptionalStringParam(ctx, step, "range")
	if err != nil {
		return nil, err
	}
	headers, err := flowStepOptionalStringListParam(ctx, step, "headers")
	if err != nil {
		return nil, err
	}
	readOptions, err := flowStepTableReadOptions(ctx, step)
	if err != nil {
		return nil, err
	}
	return loadExcelRows(filePath, excelReadOptions{
		Sheet:          sheet,
		Range:          rangeSpec,
		Headers:        headers,
		StartRow:       readOptions.StartRow,
		Limit:          readOptions.Limit,
		RowNumberField: readOptions.RowNumberField,
	})
}

func flowStepTableReadOptions(ctx *FlowContext, step FlowStep) (tableReadOptions, error) {
	startRow, err := flowStepOptionalIntParam(ctx, step, "start_row")
	if err != nil {
		return tableReadOptions{}, err
	}
	limit, err := flowStepOptionalIntParam(ctx, step, "limit")
	if err != nil {
		return tableReadOptions{}, err
	}
	rowNumberField, err := flowStepOptionalStringParam(ctx, step, "row_number_field")
	if err != nil {
		return tableReadOptions{}, err
	}
	options := tableReadOptions{
		StartRow:       startRow,
		Limit:          limit,
		RowNumberField: rowNumberField,
	}
	if err := validateTableReadOptions(options); err != nil {
		return tableReadOptions{}, err
	}
	return options, nil
}

func flowStepStringParam(ctx *FlowContext, step FlowStep, name string) (string, error) {
	value, ok := step.param(name)
	if !ok {
		return "", fmt.Errorf("action %q requires %q", step.Action, name)
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return "", err
	}
	text, ok := resolved.(string)
	if !ok {
		return "", fmt.Errorf("action %q %q must be a string", step.Action, name)
	}
	return text, nil
}

func flowStepOptionalStringParam(ctx *FlowContext, step FlowStep, name string) (string, error) {
	value, ok := step.param(name)
	if !ok {
		return "", nil
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return "", err
	}
	text, ok := resolved.(string)
	if !ok {
		return "", fmt.Errorf("action %q %q must be a string", step.Action, name)
	}
	return text, nil
}

func flowStepOptionalStringListParam(ctx *FlowContext, step FlowStep, name string) ([]string, error) {
	value, ok := step.param(name)
	if !ok {
		return nil, nil
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return nil, err
	}
	switch typed := resolved.(type) {
	case []string:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return items, nil
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("action %q %q must be a list of strings", step.Action, name)
			}
			items = append(items, text)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("action %q %q must be a list of strings", step.Action, name)
	}
}

func flowStepOptionalIntParam(ctx *FlowContext, step FlowStep, name string) (int, error) {
	value, ok := step.param(name)
	if !ok {
		return 0, nil
	}
	resolved, err := resolveValue(value, ctx)
	if err != nil {
		return 0, err
	}
	return intParam(resolved)
}

func flowTextContains(actual any, expected string) bool {
	for _, text := range flowTextValues(actual) {
		if strings.Contains(text, expected) {
			return true
		}
	}
	return false
}

func flowTextValues(actual any) []string {
	switch typed := actual.(type) {
	case nil:
		return nil
	case string:
		return []string{strings.TrimSpace(typed)}
	case []string:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			values = append(values, strings.TrimSpace(item))
		}
		return values
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			values = append(values, flowTextValues(item)...)
		}
		return values
	default:
		return []string{strings.TrimSpace(fmt.Sprint(typed))}
	}
}

func compactExtractedTexts(texts []string) any {
	switch len(texts) {
	case 0:
		return ""
	case 1:
		return texts[0]
	default:
		return texts
	}
}

func extractTextByPattern(texts []string, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}
	for _, text := range texts {
		match := re.FindStringSubmatch(text)
		if len(match) == 0 {
			continue
		}
		if len(match) > 1 {
			return strings.TrimSpace(match[1]), nil
		}
		return strings.TrimSpace(match[0]), nil
	}
	return "", fmt.Errorf("did not match any extracted text")
}

func buildActionArgs(L *lua.LState, ctx *FlowContext, step FlowStep) ([]lua.LValue, error) {
	values := make([]any, 0)
	if len(step.Args) > 0 {
		values = append(values, step.Args...)
	} else {
		spec := flowActionSpecs[step.Action]
		for _, arg := range spec.Args {
			value, ok := step.param(arg.Name)
			if !ok {
				if arg.Required {
					return nil, fmt.Errorf("action %q requires %q", step.Action, arg.Name)
				}
				continue
			}
			values = append(values, value)
		}
		if spec.VarArgName != "" {
			value, ok := step.param(spec.VarArgName)
			if !ok {
				return nil, fmt.Errorf("action %q requires %q", step.Action, spec.VarArgName)
			}
			items, err := toList(value)
			if err != nil {
				return nil, fmt.Errorf("action %q %q must be a list: %w", step.Action, spec.VarArgName, err)
			}
			values = append(values, items...)
		}
	}

	args := make([]lua.LValue, 0, len(values))
	for _, value := range values {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return nil, err
		}
		args = append(args, goValueToLua(L, resolved))
	}
	return args, nil
}

func rewriteFlowFileAccessArgs(step FlowStep, args []lua.LValue, policy *FlowSecurityPolicy) ([]lua.LValue, error) {
	if policy == nil || !policy.AllowFileAccess {
		return args, nil
	}
	params := flowFilePathParams(step.Action)
	if len(params) == 0 {
		return args, nil
	}

	rewritten := append([]lua.LValue(nil), args...)
	spec := flowActionSpecs[step.Action]
	for i, value := range rewritten {
		name := ""
		if i < len(spec.Args) {
			name = spec.Args[i].Name
		} else {
			name = spec.VarArgName
		}
		role, ok := params[name]
		if !ok {
			continue
		}
		path := value.String()
		resolved, err := resolveRuntimeFilePath(path, role, *policy)
		if err != nil {
			return nil, fmt.Errorf("action %q parameter %q %w", step.Action, name, err)
		}
		rewritten[i] = lua.LString(resolved)
	}
	return rewritten, nil
}

func resolveRuntimeFilePath(path string, role flowFilePathRole, policy FlowSecurityPolicy) (string, error) {
	root := flowFileRootForRole(role, policy)
	if root == "" {
		return path, nil
	}
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	rootReal, err := prepareRuntimeFileRoot(root)
	if err != nil {
		return "", fmt.Errorf("resolve file %s root %q: %w", role, root, err)
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootReal, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", path, err)
	}
	if err := ensurePathInsideRoot(candidateAbs, rootReal); err != nil {
		return "", fmt.Errorf("path %q is outside allowed file %s root %q", path, role, rootReal)
	}

	switch role {
	case flowFileInputPath:
		candidateReal, err := filepath.EvalSymlinks(candidateAbs)
		if err != nil {
			return "", fmt.Errorf("input path %q is not accessible: %w", path, err)
		}
		if err := ensurePathInsideRoot(candidateReal, rootReal); err != nil {
			return "", fmt.Errorf("input path %q is outside allowed file input root %q", path, rootReal)
		}
		return candidateReal, nil
	case flowFileOutputPath:
		parent := filepath.Dir(candidateAbs)
		if err := os.MkdirAll(parent, 0755); err != nil {
			return "", fmt.Errorf("create output directory %q: %w", parent, err)
		}
		parentReal, err := filepath.EvalSymlinks(parent)
		if err != nil {
			return "", fmt.Errorf("resolve output directory %q: %w", parent, err)
		}
		if err := ensurePathInsideRoot(parentReal, rootReal); err != nil {
			return "", fmt.Errorf("output path %q is outside allowed file output root %q", path, rootReal)
		}
		if info, err := os.Lstat(candidateAbs); err == nil && info.Mode()&os.ModeSymlink != 0 {
			targetReal, err := filepath.EvalSymlinks(candidateAbs)
			if err != nil {
				return "", fmt.Errorf("resolve output symlink %q: %w", path, err)
			}
			if err := ensurePathInsideRoot(targetReal, rootReal); err != nil {
				return "", fmt.Errorf("output path %q is outside allowed file output root %q", path, rootReal)
			}
		}
		return candidateAbs, nil
	default:
		return path, nil
	}
}

func prepareRuntimeFileRoot(root string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(rootAbs, 0755); err != nil {
		return "", err
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", err
	}
	return rootReal, nil
}

func callLuaChunk(L *lua.LState, code string) (any, error) {
	fn, err := L.LoadString(code)
	if err != nil {
		return nil, err
	}
	top := L.GetTop()
	if err := L.CallByParam(lua.P{Fn: fn, NRet: lua.MultRet, Protect: true}); err != nil {
		L.SetTop(top)
		return nil, err
	}
	return collectReturns(L, top), nil
}

func collectReturns(L *lua.LState, top int) any {
	count := L.GetTop() - top
	if count <= 0 {
		return nil
	}

	values := make([]any, 0, count)
	for i := top + 1; i <= L.GetTop(); i++ {
		values = append(values, luaValueToGo(L.Get(i)))
	}
	L.SetTop(top)

	if len(values) == 1 {
		return values[0]
	}
	return values
}

func (step FlowStep) presentNamedParams() map[string]any {
	params := map[string]any{}
	add := func(name string, value any, ok bool) {
		if ok {
			params[name] = value
		}
	}
	addString := func(name string, value string) {
		param, ok := stringParam(value)
		add(name, param, ok)
	}

	addString("url", step.URL)
	addString("selector", step.Selector)
	addString("text", step.Text)
	addString("value", step.Value)
	if step.Timeout != 0 {
		params["timeout"] = step.Timeout
	}
	if step.Seconds != 0 {
		params["seconds"] = step.Seconds
	}
	addString("path", step.Path)
	addString("range", step.Range)
	addString("script", step.Script)
	addString("code", step.Code)
	addString("attribute", step.Attribute)
	addString("sheet", step.Sheet)
	addString("key", step.Key)
	addString("file_path", step.FilePath)
	if len(step.Files) > 0 {
		params["files"] = step.Files
	}
	addString("save_path", step.SavePath)
	addString("pattern", step.Pattern)
	if step.From != nil {
		params["from"] = step.From
	}
	addString("connection", step.Connection)
	if step.Index != 0 {
		params["index"] = step.Index
	}
	if step.ContextIndex != 0 {
		params["context_index"] = step.ContextIndex
	}
	if step.Delta != 0 {
		params["delta"] = step.Delta
	}
	if step.TTLSeconds != 0 {
		params["ttl_seconds"] = step.TTLSeconds
	}
	if step.Times != 0 {
		params["times"] = step.Times
	}
	if step.IntervalMS != 0 {
		params["interval_ms"] = step.IntervalMS
	}
	if len(step.Steps) > 0 {
		params["steps"] = step.Steps
	}
	if step.Condition != nil {
		params["condition"] = step.Condition
	}
	if len(step.Then) > 0 {
		params["then"] = step.Then
	}
	if len(step.Else) > 0 {
		params["else"] = step.Else
	}
	if len(step.OnError) > 0 {
		params["on_error"] = step.OnError
	}
	if step.Items != nil {
		params["items"] = step.Items
	}
	addString("item_var", step.ItemVar)
	addString("index_var", step.IndexVar)
	for name, value := range step.With {
		params[name] = value
	}
	return params
}

func (step FlowStep) luaCode(ctx *FlowContext) (string, error) {
	code := step.Code
	if code == "" {
		if value, ok := step.param("code"); ok {
			code = fmt.Sprint(value)
		}
	}
	resolved, err := resolveValue(code, ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(resolved), nil
}

func (step FlowStep) param(name string) (any, bool) {
	if step.With != nil {
		if value, ok := step.With[name]; ok {
			return value, true
		}
	}
	if len(step.Args) > 0 {
		spec, ok := flowActionSpecs[step.Action]
		if ok {
			for i, value := range step.Args {
				if i < len(spec.Args) && spec.Args[i].Name == name {
					return value, true
				}
			}
			if spec.VarArgName == name {
				if len(step.Args) <= len(spec.Args) {
					return nil, false
				}
				items := make([]any, 0, len(step.Args)-len(spec.Args))
				items = append(items, step.Args[len(spec.Args):]...)
				return items, true
			}
		}
	}

	switch name {
	case "url":
		return stringParam(step.URL)
	case "selector":
		return stringParam(step.Selector)
	case "text":
		return stringParam(step.Text)
	case "value":
		return stringParam(step.Value)
	case "timeout":
		if step.Timeout == 0 {
			return nil, false
		}
		return step.Timeout, true
	case "seconds":
		if step.Seconds == 0 {
			return nil, false
		}
		return step.Seconds, true
	case "path":
		return stringParam(step.Path)
	case "range":
		return stringParam(step.Range)
	case "script":
		return stringParam(step.Script)
	case "code":
		return stringParam(step.Code)
	case "attribute":
		return stringParam(step.Attribute)
	case "sheet":
		return stringParam(step.Sheet)
	case "key":
		return stringParam(step.Key)
	case "file_path":
		return stringParam(step.FilePath)
	case "files":
		if len(step.Files) == 0 {
			return nil, false
		}
		return step.Files, true
	case "save_path":
		return stringParam(step.SavePath)
	case "pattern":
		return stringParam(step.Pattern)
	case "from":
		if step.From == nil {
			return nil, false
		}
		return step.From, true
	case "connection":
		return stringParam(step.Connection)
	case "index":
		if step.Index == 0 {
			return nil, false
		}
		return step.Index, true
	case "context_index":
		if step.ContextIndex == 0 {
			return nil, false
		}
		return step.ContextIndex, true
	case "delta":
		if step.Delta == 0 {
			return nil, false
		}
		return step.Delta, true
	case "ttl_seconds":
		if step.TTLSeconds == 0 {
			return nil, false
		}
		return step.TTLSeconds, true
	case "times":
		if step.Times == 0 {
			return nil, false
		}
		return step.Times, true
	case "interval_ms":
		if step.IntervalMS == 0 {
			return nil, false
		}
		return step.IntervalMS, true
	case "steps":
		if len(step.Steps) == 0 {
			return nil, false
		}
		return step.Steps, true
	case "condition":
		if step.Condition == nil {
			return nil, false
		}
		return step.Condition, true
	case "then":
		if len(step.Then) == 0 {
			return nil, false
		}
		return step.Then, true
	case "else":
		if len(step.Else) == 0 {
			return nil, false
		}
		return step.Else, true
	case "on_error":
		if len(step.OnError) == 0 {
			return nil, false
		}
		return step.OnError, true
	case "items":
		if step.Items == nil {
			return nil, false
		}
		return step.Items, true
	case "item_var":
		return stringParam(step.ItemVar)
	case "index_var":
		return stringParam(step.IndexVar)
	default:
		return nil, false
	}
}

func stringParam(value string) (any, bool) {
	if value == "" {
		return nil, false
	}
	return value, true
}

func retryTimes(ctx *FlowContext, step FlowStep) (int, error) {
	value, ok := step.param("times")
	if !ok {
		return 3, nil
	}
	if ctx != nil {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return 0, err
		}
		value = resolved
	}
	return intParam(value)
}

func retryIntervalMS(ctx *FlowContext, step FlowStep) (int, error) {
	value, ok := step.param("interval_ms")
	if !ok {
		return 0, nil
	}
	if ctx != nil {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return 0, err
		}
		value = resolved
	}
	return intParam(value)
}

func waitUntilTimeoutMS(ctx *FlowContext, step FlowStep) (int, error) {
	value, ok := step.param("timeout")
	if !ok {
		return 30000, nil
	}
	if ctx != nil {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return 0, err
		}
		value = resolved
	}
	timeout, err := intParam(value)
	if err != nil {
		return 0, err
	}
	if timeout < 1 {
		return 0, fmt.Errorf("wait_until timeout must be at least 1")
	}
	return timeout, nil
}

func waitUntilIntervalMS(ctx *FlowContext, step FlowStep) (int, error) {
	value, ok := step.param("interval_ms")
	if !ok {
		return 500, nil
	}
	if ctx != nil {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return 0, err
		}
		value = resolved
	}
	intervalMS, err := intParam(value)
	if err != nil {
		return 0, err
	}
	if intervalMS < 1 {
		return 0, fmt.Errorf("wait_until interval_ms must be at least 1")
	}
	return intervalMS, nil
}

func flowReferenceExpressions(value any) []string {
	refs := []string{}
	switch typed := value.(type) {
	case string:
		matches := replacePattern.FindAllStringSubmatch(typed, -1)
		for _, match := range matches {
			if len(match) == 2 {
				refs = append(refs, strings.TrimSpace(match[1]))
			}
		}
	case []any:
		for _, item := range typed {
			refs = append(refs, flowReferenceExpressions(item)...)
		}
	case []string:
		for _, item := range typed {
			refs = append(refs, flowReferenceExpressions(item)...)
		}
	case []FlowStep:
		for _, item := range typed {
			refs = append(refs, flowReferenceExpressions(item.presentNamedParams())...)
		}
	case *FlowStep:
		if typed != nil {
			refs = append(refs, flowReferenceExpressions(typed.presentNamedParams())...)
		}
	case FlowStep:
		refs = append(refs, flowReferenceExpressions(typed.presentNamedParams())...)
	case map[string]any:
		for _, item := range typed {
			refs = append(refs, flowReferenceExpressions(item)...)
		}
	}
	return refs
}

func parseFlowVariableReference(ref string) (string, string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", "", fmt.Errorf("placeholder reference cannot be empty")
	}

	end := 0
	for end < len(ref) {
		r := rune(ref[end])
		if end == 0 && !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_') {
			return "", "", fmt.Errorf("placeholder reference %q must start with a variable name", ref)
		}
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			break
		}
		end++
	}
	base := ref[:end]
	if !flowIdentifierPattern.MatchString(base) {
		return "", "", fmt.Errorf("placeholder reference %q starts with invalid variable name %q", ref, base)
	}

	rest := strings.TrimSpace(ref[end:])
	if rest == "" {
		return base, "$", nil
	}
	path := "$" + rest
	if _, err := parseJSONPath(path); err != nil {
		return "", "", fmt.Errorf("placeholder reference %q %w", ref, err)
	}
	return base, path, nil
}

func resolveFlowVariableReference(ref string, vars map[string]any) (any, error) {
	base, path, err := parseFlowVariableReference(ref)
	if err != nil {
		return nil, err
	}
	value, ok := vars[base]
	if !ok {
		return nil, fmt.Errorf("unknown flow variable %q", base)
	}
	if path == "$" {
		return value, nil
	}
	return extractJSONPathValue(value, path)
}

func resolveValue(value any, ctx *FlowContext) (any, error) {
	switch typed := value.(type) {
	case string:
		if matches := placeholderPattern.FindStringSubmatch(typed); len(matches) == 2 {
			return resolveFlowVariableReference(matches[1], ctx.Vars)
		}
		var err error
		resolved := replacePattern.ReplaceAllStringFunc(typed, func(token string) string {
			if err != nil {
				return token
			}
			matches := replacePattern.FindStringSubmatch(token)
			if len(matches) != 2 {
				return token
			}
			value, resolveErr := resolveFlowVariableReference(matches[1], ctx.Vars)
			if resolveErr != nil {
				err = resolveErr
				return token
			}
			return fmt.Sprint(value)
		})
		return resolved, err
	case []any:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			resolved, err := resolveValue(item, ctx)
			if err != nil {
				return nil, err
			}
			items = append(items, resolved)
		}
		return items, nil
	case []string:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			resolved, err := resolveValue(item, ctx)
			if err != nil {
				return nil, err
			}
			items = append(items, resolved)
		}
		return items, nil
	case map[string]any:
		resolved := map[string]any{}
		for key, item := range typed {
			value, err := resolveValue(item, ctx)
			if err != nil {
				return nil, err
			}
			resolved[key] = value
		}
		return resolved, nil
	default:
		return value, nil
	}
}

func goValueToLua(L *lua.LState, value any) lua.LValue {
	switch typed := value.(type) {
	case nil:
		return lua.LNil
	case lua.LValue:
		return typed
	case string:
		return lua.LString(typed)
	case bool:
		return lua.LBool(typed)
	case int:
		return lua.LNumber(typed)
	case int64:
		return lua.LNumber(typed)
	case float64:
		return lua.LNumber(typed)
	case float32:
		return lua.LNumber(typed)
	case []any:
		table := L.NewTable()
		for _, item := range typed {
			table.Append(goValueToLua(L, item))
		}
		return table
	case []string:
		table := L.NewTable()
		for _, item := range typed {
			table.Append(lua.LString(item))
		}
		return table
	case map[string]any:
		table := L.NewTable()
		for key, item := range typed {
			table.RawSetString(key, goValueToLua(L, item))
		}
		return table
	default:
		return lua.LString(fmt.Sprint(typed))
	}
}

func luaValueToGo(value lua.LValue) any {
	switch typed := value.(type) {
	case lua.LBool:
		return bool(typed)
	case lua.LNumber:
		return float64(typed)
	case lua.LString:
		return string(typed)
	case *lua.LNilType:
		return nil
	case *lua.LTable:
		if isArrayTable(typed) {
			values := make([]any, 0, typed.Len())
			for i := 1; i <= typed.Len(); i++ {
				values = append(values, luaValueToGo(typed.RawGetInt(i)))
			}
			return values
		}
		values := map[string]any{}
		typed.ForEach(func(key lua.LValue, value lua.LValue) {
			values[fmt.Sprint(luaValueToGo(key))] = luaValueToGo(value)
		})
		return values
	default:
		return value.String()
	}
}

func isArrayTable(table *lua.LTable) bool {
	length := table.Len()
	count := 0
	array := true
	table.ForEach(func(key lua.LValue, _ lua.LValue) {
		count++
		number, ok := key.(lua.LNumber)
		if !ok {
			array = false
			return
		}
		index := int(number)
		if float64(index) != float64(number) || index < 1 || index > length {
			array = false
		}
	})
	return array && count == length
}

func toList(value any) ([]any, error) {
	switch typed := value.(type) {
	case []any:
		return typed, nil
	case []string:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("got %T", value)
	}
}

func listLen(value any) int {
	items, err := toList(value)
	if err != nil {
		return 0
	}
	return len(items)
}

func flowReferences(value any) []string {
	refs := []string{}
	for _, ref := range flowReferenceExpressions(value) {
		base, _, err := parseFlowVariableReference(ref)
		if err != nil {
			refs = append(refs, ref)
			continue
		}
		refs = append(refs, base)
	}
	return refs
}

func fullPlaceholderRef(value any) (string, bool) {
	expr, ok := fullPlaceholderExpression(value)
	if !ok {
		return "", false
	}
	base, _, err := parseFlowVariableReference(expr)
	if err != nil {
		return "", false
	}
	return base, true
}

func fullPlaceholderExpression(value any) (string, bool) {
	text, ok := value.(string)
	if !ok {
		return "", false
	}
	matches := placeholderPattern.FindStringSubmatch(text)
	if len(matches) != 2 {
		return "", false
	}
	expr := strings.TrimSpace(matches[1])
	if _, _, err := parseFlowVariableReference(expr); err != nil {
		return "", false
	}
	return expr, true
}

func resolveKnownPlaceholderValue(value any, knownVars map[string]any) (any, bool) {
	expr, ok := fullPlaceholderExpression(value)
	if !ok {
		return nil, false
	}
	resolved, err := resolveFlowVariableReference(expr, knownVars)
	if err != nil || resolved == nil {
		return nil, false
	}
	return resolved, true
}

func isIntegerValue(value any) bool {
	switch typed := value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32:
		return math.Trunc(float64(typed)) == float64(typed)
	case float64:
		return math.Trunc(typed) == typed
	default:
		return false
	}
}

func intParam(value any) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int8:
		return int(typed), nil
	case int16:
		return int(typed), nil
	case int32:
		return int(typed), nil
	case int64:
		return int(typed), nil
	case uint:
		return int(typed), nil
	case uint8:
		return int(typed), nil
	case uint16:
		return int(typed), nil
	case uint32:
		return int(typed), nil
	case uint64:
		return int(typed), nil
	case float32:
		if math.Trunc(float64(typed)) == float64(typed) {
			return int(typed), nil
		}
	case float64:
		if math.Trunc(typed) == typed {
			return int(typed), nil
		}
	}
	return 0, fmt.Errorf("must be an integer")
}

func isNumberValue(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	default:
		return false
	}
}

func isStringListValue(value any, knownVars map[string]any) bool {
	if resolved, ok := resolveKnownPlaceholderValue(value, knownVars); ok {
		value = resolved
	} else if _, ok := fullPlaceholderExpression(value); ok {
		return true
	}

	switch typed := value.(type) {
	case []string:
		for _, item := range typed {
			if resolved, ok := resolveKnownPlaceholderValue(item, knownVars); ok {
				if _, ok := resolved.(string); !ok {
					return false
				}
				continue
			}
			if _, ok := fullPlaceholderExpression(item); ok {
				continue
			}
		}
		return true
	case []any:
		for _, item := range typed {
			if resolved, ok := resolveKnownPlaceholderValue(item, knownVars); ok {
				item = resolved
			} else if _, ok := fullPlaceholderExpression(item); ok {
				continue
			}
			if _, ok := item.(string); !ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}
