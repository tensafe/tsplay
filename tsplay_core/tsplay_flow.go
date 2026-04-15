package tsplay_core

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// Flow is the structured workflow format used by TSPlay.
// It keeps most business logic declarative, while still allowing lua steps as
// an escape hatch for advanced cases.
type Flow struct {
	SchemaVersion string         `json:"schema_version" yaml:"schema_version"`
	Name          string         `json:"name" yaml:"name"`
	Description   string         `json:"description,omitempty" yaml:"description,omitempty"`
	Vars          map[string]any `json:"vars,omitempty" yaml:"vars,omitempty"`
	Steps         []FlowStep     `json:"steps" yaml:"steps"`
}

type FlowStep struct {
	Name            string         `json:"name,omitempty" yaml:"name,omitempty"`
	Action          string         `json:"action" yaml:"action"`
	Args            []any          `json:"args,omitempty" yaml:"args,omitempty"`
	With            map[string]any `json:"with,omitempty" yaml:"with,omitempty"`
	SaveAs          string         `json:"save_as,omitempty" yaml:"save_as,omitempty"`
	ContinueOnError bool           `json:"continue_on_error,omitempty" yaml:"continue_on_error,omitempty"`

	// Common named parameters. Flow accepts both args: [...] and these named
	// fields because named fields are easier for humans and AI to review.
	URL          string   `json:"url,omitempty" yaml:"url,omitempty"`
	Selector     string   `json:"selector,omitempty" yaml:"selector,omitempty"`
	Text         string   `json:"text,omitempty" yaml:"text,omitempty"`
	Value        string   `json:"value,omitempty" yaml:"value,omitempty"`
	Timeout      int      `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Seconds      float64  `json:"seconds,omitempty" yaml:"seconds,omitempty"`
	Path         string   `json:"path,omitempty" yaml:"path,omitempty"`
	Script       string   `json:"script,omitempty" yaml:"script,omitempty"`
	Code         string   `json:"code,omitempty" yaml:"code,omitempty"`
	Attribute    string   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	FilePath     string   `json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Files        []string `json:"files,omitempty" yaml:"files,omitempty"`
	SavePath     string   `json:"save_path,omitempty" yaml:"save_path,omitempty"`
	Pattern      string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Index        int      `json:"index,omitempty" yaml:"index,omitempty"`
	ContextIndex int      `json:"context_index,omitempty" yaml:"context_index,omitempty"`
}

type FlowResult struct {
	Name         string          `json:"name"`
	Vars         map[string]any  `json:"vars,omitempty"`
	Trace        []FlowStepTrace `json:"trace"`
	ArtifactRoot string          `json:"artifact_root,omitempty"`
}

type FlowStepTrace struct {
	Index         int                `json:"index"`
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
	Vars         map[string]any
	Security     *FlowSecurityPolicy
	ArtifactRoot string
	RunID        string
}

type FlowRunOptions struct {
	Headless     bool
	Security     *FlowSecurityPolicy
	ArtifactRoot string
}

type FlowSecurityPolicy struct {
	AllowLua          bool   `json:"allow_lua"`
	AllowJavaScript   bool   `json:"allow_javascript"`
	AllowFileAccess   bool   `json:"allow_file_access"`
	AllowBrowserState bool   `json:"allow_browser_state"`
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

var placeholderPattern = regexp.MustCompile(`^\{\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}\}$`)
var replacePattern = regexp.MustCompile(`\{\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}\}`)
var flowIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

var flowActionSpecs = map[string]flowActionSpec{
	"navigate":              {Args: []flowArgSpec{{Name: "url", Required: true}}},
	"click":                 {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"reload":                {},
	"go_back":               {},
	"go_forward":            {},
	"type_text":             {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "text", Required: true}}},
	"get_text":              {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"set_value":             {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "value", Required: true}}},
	"select_option":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "value", Required: true}}},
	"hover":                 {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"scroll_to":             {Args: []flowArgSpec{{Name: "selector", Required: true}}},
	"wait_for_network_idle": {},
	"wait_for_selector":     {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "timeout"}}},
	"wait_for_text":         {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "text", Required: true}, {Name: "timeout"}}},
	"sleep":                 {Args: []flowArgSpec{{Name: "seconds", Required: true}}},
	"screenshot":            {Args: []flowArgSpec{{Name: "path", Required: true}}},
	"screenshot_element":    {Args: []flowArgSpec{{Name: "selector", Required: true}, {Name: "path", Required: true}}},
	"save_html":             {Args: []flowArgSpec{{Name: "path", Required: true}}},
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

	for i, step := range flow.Steps {
		if strings.TrimSpace(step.Action) == "" {
			return fmt.Errorf("step %d action is required", i+1)
		}
		spec, ok := flowActionSpecs[step.Action]
		if !ok {
			return fmt.Errorf("step %d uses unsupported action %q", i+1, step.Action)
		}

		if step.SaveAs != "" && !flowIdentifierPattern.MatchString(step.SaveAs) {
			return fmt.Errorf("step %d save_as %q is not a valid variable name", i+1, step.SaveAs)
		}

		if len(step.Args) > 0 {
			if err := validateFlowStepArgs(i+1, step, spec, knownVars); err != nil {
				return err
			}
		} else {
			if err := validateFlowStepNamedParams(i+1, step, spec, knownVars); err != nil {
				return err
			}
		}

		if step.SaveAs != "" {
			knownVars[step.SaveAs] = nil
		}
	}
	return nil
}

func validateFlowStepArgs(stepIndex int, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	if len(step.presentNamedParams()) > 0 {
		return fmt.Errorf("step %d action %q cannot mix args with named parameters", stepIndex, step.Action)
	}

	requiredCount := requiredArgCount(spec)
	minCount := requiredCount
	maxCount := len(spec.Args)
	if spec.VarArgName != "" {
		minCount = len(spec.Args) + 1
		maxCount = -1
	}

	if len(step.Args) < minCount {
		return fmt.Errorf("step %d action %q expects at least %d args, got %d", stepIndex, step.Action, minCount, len(step.Args))
	}
	if maxCount >= 0 && len(step.Args) > maxCount {
		return fmt.Errorf("step %d action %q expects at most %d args, got %d", stepIndex, step.Action, maxCount, len(step.Args))
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

func validateFlowStepNamedParams(stepIndex int, step FlowStep, spec flowActionSpec, knownVars map[string]any) error {
	present := step.presentNamedParams()
	allowed := allowedFlowParamNames(spec)

	for name, value := range present {
		if !allowed[name] {
			return fmt.Errorf("step %d action %q does not accept parameter %q", stepIndex, step.Action, name)
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
			return fmt.Errorf("step %d action %q requires %q", stepIndex, step.Action, arg.Name)
		}
	}
	if spec.VarArgName != "" {
		value, ok := present[spec.VarArgName]
		if !ok || listLen(value) == 0 {
			return fmt.Errorf("step %d action %q requires %q", stepIndex, step.Action, spec.VarArgName)
		}
	}
	return nil
}

func validateFlowParamValue(stepIndex int, action string, name string, value any, knownVars map[string]any) error {
	if name == "" {
		return fmt.Errorf("step %d action %q has too many arguments", stepIndex, action)
	}
	if err := validateFlowReferences(stepIndex, action, name, value, knownVars); err != nil {
		return err
	}
	if err := validateFlowParamType(name, value, knownVars); err != nil {
		return fmt.Errorf("step %d action %q parameter %q %w", stepIndex, action, name, err)
	}
	return nil
}

func validateFlowReferences(stepIndex int, action string, name string, value any, knownVars map[string]any) error {
	for _, ref := range flowReferences(value) {
		if _, ok := knownVars[ref]; !ok {
			return fmt.Errorf("step %d action %q parameter %q references unknown variable %q", stepIndex, action, name, ref)
		}
	}
	return nil
}

func validateFlowParamType(name string, value any, knownVars map[string]any) error {
	if ref, ok := fullPlaceholderRef(value); ok {
		resolved, known := knownVars[ref]
		if !known || resolved == nil {
			return nil
		}
		value = resolved
	}

	switch flowParamType(name) {
	case "string":
		if _, ok := value.(string); ok {
			return nil
		}
		return fmt.Errorf("must be a string")
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
	default:
		return nil
	}
}

func flowParamType(name string) string {
	switch name {
	case "url", "selector", "text", "value", "path", "script", "code", "attribute", "file_path", "save_path", "pattern":
		return "string"
	case "timeout", "index", "context_index":
		return "int"
	case "seconds":
		return "number"
	case "files":
		return "string_list"
	default:
		return ""
	}
}

func ValidateFlowSecurity(flow *Flow, policy FlowSecurityPolicy) error {
	if flow == nil {
		return fmt.Errorf("flow is nil")
	}

	for i, step := range flow.Steps {
		group := flowActionSecurityGroup(step.Action)
		if group == "" || flowSecurityPolicyAllows(group, policy) {
			continue
		}
		option := flowActionSecurityOption(group)
		return fmt.Errorf("step %d action %q is disabled by security policy; set %s=true only for trusted flows", i+1, step.Action, option)
	}
	if err := validateFlowFileAccessRoots(flow, policy); err != nil {
		return err
	}
	return nil
}

func flowActionSecurityGroup(action string) string {
	switch action {
	case "lua":
		return "lua"
	case "execute_script", "evaluate":
		return "javascript"
	case "screenshot", "screenshot_element", "save_html", "upload_file", "upload_multiple_files", "download_file", "download_url":
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
	default:
		return true
	}
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
	for i, step := range flow.Steps {
		if err := validateFlowStepFileAccessRoots(i+1, step, policy); err != nil {
			return err
		}
	}
	return nil
}

func validateFlowStepFileAccessRoots(stepIndex int, step FlowStep, policy FlowSecurityPolicy) error {
	return forEachFlowFilePathValue(step, func(name string, role flowFilePathRole, value any) error {
		return validateFlowFilePathValue(stepIndex, step.Action, name, role, value, policy)
	})
}

func validateFlowFilePathValue(stepIndex int, action string, name string, role flowFilePathRole, value any, policy FlowSecurityPolicy) error {
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
	case string:
		root := flowFileRootForRole(role, policy)
		if root == "" {
			return nil
		}
		if err := validatePathWithinRoot(typed, root); err != nil {
			return fmt.Errorf("step %d action %q parameter %q is outside allowed file %s root %q: %w", stepIndex, action, name, role, root, err)
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
	case "download_file", "download_url":
		return map[string]flowFilePathRole{"save_path": flowFileOutputPath}
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

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(options.Headless),
	})
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %w", err)
	}
	defer page.Close()

	L := lua.NewState()
	defer L.Close()

	udBrowser := L.NewUserData()
	udBrowser.Value = browser
	L.SetGlobal("browser", udBrowser)

	udPage := L.NewUserData()
	udPage.Value = page
	L.SetGlobal("page", udPage)

	for _, fn := range GlobalPlayWrightFunc {
		L.SetGlobal(fn.Name, L.NewFunction(fn.Func))
	}

	return RunFlowInStateWithOptions(L, flow, options)
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

	artifactRoot := flowArtifactRoot(options)
	ctx := &FlowContext{
		Vars:         map[string]any{},
		Security:     options.Security,
		ArtifactRoot: artifactRoot,
		RunID:        newFlowRunID(flow),
	}
	for key, value := range flow.Vars {
		ctx.Vars[key] = value
		L.SetGlobal(key, goValueToLua(L, value))
	}

	result := &FlowResult{Name: flow.Name, Vars: ctx.Vars, ArtifactRoot: artifactRoot}
	for i, step := range flow.Steps {
		traceArgs := traceStepParams(step, ctx)
		trace := FlowStepTrace{
			Index:       i + 1,
			Name:        step.Name,
			Action:      step.Action,
			Args:        compactTraceValue(traceArgs, 0),
			ArgsSummary: summarizeTraceValue(traceArgs),
			SaveAs:      step.SaveAs,
			Status:      "running",
			StartedAt:   time.Now().Format(time.RFC3339Nano),
		}

		output, err := runFlowStep(L, ctx, step)
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
			result.Trace = append(result.Trace, trace)
			if !step.ContinueOnError {
				return result, fmt.Errorf("step %d %q failed: %w", i+1, step.Action, err)
			}
			continue
		}

		trace.Status = "ok"
		trace.Output = compactTraceValue(output, 0)
		trace.OutputSummary = summarizeTraceValue(output)
		if step.SaveAs != "" {
			ctx.Vars[step.SaveAs] = output
			L.SetGlobal(step.SaveAs, goValueToLua(L, output))
		}
		result.Trace = append(result.Trace, trace)
	}

	return result, nil
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
	if ctx == nil || strings.TrimSpace(ctx.ArtifactRoot) == "" {
		return artifacts
	}

	root, err := prepareRuntimeFileRoot(ctx.ArtifactRoot)
	if err != nil {
		artifacts.CaptureError = fmt.Sprintf("prepare artifact root: %v", err)
		return artifacts
	}
	dir := filepath.Join(root, ctx.RunID, fmt.Sprintf("%02d-%s", trace.Index, sanitizeArtifactSegment(trace.Action)))
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
	addString("script", step.Script)
	addString("code", step.Code)
	addString("attribute", step.Attribute)
	addString("file_path", step.FilePath)
	if len(step.Files) > 0 {
		params["files"] = step.Files
	}
	addString("save_path", step.SavePath)
	addString("pattern", step.Pattern)
	if step.Index != 0 {
		params["index"] = step.Index
	}
	if step.ContextIndex != 0 {
		params["context_index"] = step.ContextIndex
	}
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
	case "script":
		return stringParam(step.Script)
	case "code":
		return stringParam(step.Code)
	case "attribute":
		return stringParam(step.Attribute)
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

func resolveValue(value any, ctx *FlowContext) (any, error) {
	switch typed := value.(type) {
	case string:
		if matches := placeholderPattern.FindStringSubmatch(typed); len(matches) == 2 {
			value, ok := ctx.Vars[matches[1]]
			if !ok {
				return nil, fmt.Errorf("unknown flow variable %q", matches[1])
			}
			return value, nil
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
			value, ok := ctx.Vars[matches[1]]
			if !ok {
				err = fmt.Errorf("unknown flow variable %q", matches[1])
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
	switch typed := value.(type) {
	case string:
		matches := replacePattern.FindAllStringSubmatch(typed, -1)
		for _, match := range matches {
			if len(match) == 2 {
				refs = append(refs, match[1])
			}
		}
	case []any:
		for _, item := range typed {
			refs = append(refs, flowReferences(item)...)
		}
	case []string:
		for _, item := range typed {
			refs = append(refs, flowReferences(item)...)
		}
	case map[string]any:
		for _, item := range typed {
			refs = append(refs, flowReferences(item)...)
		}
	}
	return refs
}

func fullPlaceholderRef(value any) (string, bool) {
	text, ok := value.(string)
	if !ok {
		return "", false
	}
	matches := placeholderPattern.FindStringSubmatch(text)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
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
	if ref, ok := fullPlaceholderRef(value); ok {
		resolved, known := knownVars[ref]
		if !known || resolved == nil {
			return true
		}
		value = resolved
	}

	switch typed := value.(type) {
	case []string:
		for _, item := range typed {
			if ref, ok := fullPlaceholderRef(item); ok {
				resolved, known := knownVars[ref]
				if known && resolved != nil {
					if _, ok := resolved.(string); !ok {
						return false
					}
				}
			}
		}
		return true
	case []any:
		for _, item := range typed {
			if ref, ok := fullPlaceholderRef(item); ok {
				resolved, known := knownVars[ref]
				if !known || resolved == nil {
					continue
				}
				item = resolved
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
