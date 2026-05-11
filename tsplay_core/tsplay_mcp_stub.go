//go:build win7

package tsplay_core

import (
	"context"
	"fmt"
	"log"
)

const DefaultMCPFlowPathRoot = "script"
const DefaultMCPArtifactRoot = DefaultFlowArtifactRoot

type TSPlayMCPServerOptions struct {
	FlowPathRoot                       string
	ArtifactRoot                       string
	MaxConcurrentBrowserRuns           int
	MaxConcurrentBrowserRunsPerSession int
	DefaultRunTimeoutMS                int
	MaxRunTimeoutMS                    int
	QueueTimeoutMS                     int
}

func DefaultTSPlayMCPServerOptions() TSPlayMCPServerOptions {
	return TSPlayMCPServerOptions{
		FlowPathRoot: DefaultMCPFlowPathRoot,
		ArtifactRoot: DefaultMCPArtifactRoot,
	}
}

func InvokeTSPlayTool(
	ctx context.Context,
	tool string,
	arguments map[string]any,
	options ...TSPlayMCPServerOptions,
) (map[string]any, error) {
	return nil, mcpDisabledError()
}

func McpServerMCP(addr string, options ...TSPlayMCPServerOptions) {
	log.Fatal(mcpDisabledError())
}

func McpServerStdio(options ...TSPlayMCPServerOptions) {
	log.Fatal(mcpDisabledError())
}

func McpServerSSE() {
	log.Fatal(mcpDisabledError())
}

func mcpDisabledError() error {
	return fmt.Errorf("MCP support is disabled in this build; rebuild without -tags win7 to enable MCP")
}

func watchContextCancel(ctx context.Context, cancel func()) func() {
	if ctx == nil || cancel == nil {
		return func() {}
	}
	stopped := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-stopped:
		}
	}()
	return func() {
		close(stopped)
	}
}

func flowResultForTool(result *FlowResult) *FlowResult {
	if result == nil {
		return nil
	}
	sanitized := *result
	if vars, ok := compactTraceValue(result.Vars, 0).(map[string]any); ok {
		sanitized.Vars = vars
	}
	return &sanitized
}

func buildFlowActionManifest() []map[string]any {
	descriptions := map[string]string{}
	for _, fn := range GlobalPlayWrightFunc {
		descriptions[fn.Name] = fn.Description_en
	}
	descriptions["lua"] = "Run an inline Lua code block. Prefer structured actions for normal browser steps and use lua only as an escape hatch."

	actions := make([]map[string]any, 0, len(flowActionSpecs))
	for _, name := range FlowActionNames() {
		spec := flowActionSpecs[name]
		args := make([]map[string]any, 0, len(spec.Args))
		for _, arg := range spec.Args {
			args = append(args, map[string]any{
				"name":     arg.Name,
				"type":     flowParamType(arg.Name),
				"required": arg.Required,
			})
		}
		item := map[string]any{
			"name":        name,
			"description": descriptions[name],
			"args":        args,
		}
		if capabilities, ok := flowActionCapabilitiesFor(name); ok {
			item["capabilities"] = capabilities.manifestValue()
		}
		if group := flowActionSecurityGroup(name); group != "" {
			item["security_group"] = group
			item["requires_allow"] = flowActionSecurityOption(group)
		}
		if spec.VarArgName != "" {
			item["var_arg"] = spec.VarArgName
			item["var_arg_type"] = flowParamType(spec.VarArgName)
		}
		actions = append(actions, item)
	}
	return actions
}
