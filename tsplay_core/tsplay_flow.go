package tsplay_core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// Flow is the structured workflow format used by TSPlay.
// It keeps most business logic declarative, while still allowing lua steps as
// an escape hatch for advanced cases.
type Flow struct {
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Vars        map[string]any `json:"vars,omitempty" yaml:"vars,omitempty"`
	Steps       []FlowStep     `json:"steps" yaml:"steps"`
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
	Name  string          `json:"name"`
	Vars  map[string]any  `json:"vars,omitempty"`
	Trace []FlowStepTrace `json:"trace"`
}

type FlowStepTrace struct {
	Index      int    `json:"index"`
	Name       string `json:"name,omitempty"`
	Action     string `json:"action"`
	Status     string `json:"status"`
	SaveAs     string `json:"save_as,omitempty"`
	Error      string `json:"error,omitempty"`
	Output     any    `json:"output,omitempty"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	DurationMS int64  `json:"duration_ms"`
}

type FlowContext struct {
	Vars map[string]any
}

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

func LoadFlowFile(path string) (*Flow, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var flow Flow
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		err = json.Unmarshal(content, &flow)
	default:
		err = yaml.Unmarshal(content, &flow)
	}
	if err != nil {
		return nil, fmt.Errorf("parse flow %s: %w", path, err)
	}
	return &flow, nil
}

func ValidateFlow(flow *Flow) error {
	if flow == nil {
		return fmt.Errorf("flow is nil")
	}
	if len(flow.Steps) == 0 {
		return fmt.Errorf("flow must contain at least one step")
	}
	for i, step := range flow.Steps {
		if strings.TrimSpace(step.Action) == "" {
			return fmt.Errorf("step %d action is required", i+1)
		}
		spec, ok := flowActionSpecs[step.Action]
		if !ok {
			return fmt.Errorf("step %d uses unsupported action %q", i+1, step.Action)
		}
		if len(step.Args) > 0 {
			continue
		}
		for _, arg := range spec.Args {
			if !arg.Required {
				continue
			}
			if _, ok := step.param(arg.Name); !ok {
				return fmt.Errorf("step %d action %q requires %q", i+1, step.Action, arg.Name)
			}
		}
		if spec.VarArgName != "" {
			if values, ok := step.param(spec.VarArgName); !ok || listLen(values) == 0 {
				return fmt.Errorf("step %d action %q requires %q", i+1, step.Action, spec.VarArgName)
			}
		}
	}
	return nil
}

func FlowActionNames() []string {
	names := make([]string, 0, len(flowActionSpecs))
	for name := range flowActionSpecs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func RunFlowInState(L *lua.LState, flow *Flow) (*FlowResult, error) {
	if err := ValidateFlow(flow); err != nil {
		return nil, err
	}

	ctx := &FlowContext{Vars: map[string]any{}}
	for key, value := range flow.Vars {
		ctx.Vars[key] = value
		L.SetGlobal(key, goValueToLua(L, value))
	}

	result := &FlowResult{Name: flow.Name, Vars: ctx.Vars}
	for i, step := range flow.Steps {
		trace := FlowStepTrace{
			Index:     i + 1,
			Name:      step.Name,
			Action:    step.Action,
			SaveAs:    step.SaveAs,
			Status:    "running",
			StartedAt: time.Now().Format(time.RFC3339Nano),
		}

		output, err := runFlowStep(L, ctx, step)
		trace.FinishedAt = time.Now().Format(time.RFC3339Nano)
		started, _ := time.Parse(time.RFC3339Nano, trace.StartedAt)
		finished, _ := time.Parse(time.RFC3339Nano, trace.FinishedAt)
		trace.DurationMS = finished.Sub(started).Milliseconds()

		if err != nil {
			trace.Status = "error"
			trace.Error = err.Error()
			result.Trace = append(result.Trace, trace)
			if !step.ContinueOnError {
				return result, fmt.Errorf("step %d %q failed: %w", i+1, step.Action, err)
			}
			continue
		}

		trace.Status = "ok"
		trace.Output = output
		if step.SaveAs != "" {
			ctx.Vars[step.SaveAs] = output
			L.SetGlobal(step.SaveAs, goValueToLua(L, output))
		}
		result.Trace = append(result.Trace, trace)
	}

	return result, nil
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
