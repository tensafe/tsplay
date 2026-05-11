package tsplay_core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// InvokeTSPlayTool runs a TSPlay MCP tool handler directly and returns the
// JSON payload that the MCP server would emit for the same request.
func InvokeTSPlayTool(
	ctx context.Context,
	tool string,
	arguments map[string]any,
	options ...TSPlayMCPServerOptions,
) (map[string]any, error) {
	tool = strings.TrimSpace(tool)
	if tool == "" {
		return nil, fmt.Errorf("tool is required")
	}
	if arguments == nil {
		arguments = map[string]any{}
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      tool,
			Arguments: arguments,
		},
	}
	normalized := normalizeTSPlayMCPServerOptions(options)

	var (
		result *mcp.CallToolResult
		err    error
	)

	switch tool {
	case "tsplay.list_actions":
		result, err = handleFlowListActionsTool(ctx, request)
	case "tsplay.list_sessions":
		result, err = handleListSessionsToolWithOptions(ctx, request, normalized)
	case "tsplay.get_session":
		result, err = handleGetSessionToolWithOptions(ctx, request, normalized)
	case "tsplay.export_session_flow_snippet":
		result, err = handleExportSessionFlowSnippetToolWithOptions(ctx, request, normalized)
	case "tsplay.delete_session":
		result, err = handleDeleteSessionToolWithOptions(ctx, request, normalized)
	case "tsplay.save_session":
		result, err = handleSaveSessionToolWithOptions(ctx, request, normalized)
	case "tsplay.flow_schema":
		result, err = handleFlowSchemaTool(ctx, request)
	case "tsplay.flow_examples":
		result, err = handleFlowExamplesTool(ctx, request)
	case "tsplay.observe_page":
		result, err = handleObservePageToolWithOptions(ctx, request, normalized)
	case "tsplay.draft_flow":
		result, err = handleDraftFlowToolWithOptions(ctx, request, normalized)
	case "tsplay.finalize_flow":
		result, err = handleFinalizeFlowToolWithOptions(ctx, request, normalized)
	case "tsplay.validate_flow":
		result, err = handleValidateFlowToolWithOptions(ctx, request, normalized)
	case "tsplay.run_flow":
		result, err = handleRunFlowToolWithOptions(ctx, request, normalized)
	case "tsplay.repair_flow_context":
		result, err = handleRepairFlowContextToolWithOptions(ctx, request, normalized)
	case "tsplay.repair_flow":
		result, err = handleRepairFlowToolWithOptions(ctx, request, normalized)
	default:
		return nil, fmt.Errorf("unsupported TSPlay MCP tool %q", tool)
	}
	if err != nil {
		return nil, err
	}
	return decodeTSPlayToolResult(result)
}

func decodeTSPlayToolResult(result *mcp.CallToolResult) (map[string]any, error) {
	if result == nil {
		return nil, fmt.Errorf("tool returned no result")
	}
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("tool returned no content")
	}

	var text string
	switch content := result.Content[0].(type) {
	case mcp.TextContent:
		text = content.Text
	case *mcp.TextContent:
		if content != nil {
			text = content.Text
		}
	default:
		return nil, fmt.Errorf("unsupported tool result content type %T", result.Content[0])
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("tool returned empty text content")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return nil, fmt.Errorf("decode tool JSON payload: %w", err)
	}
	return payload, nil
}
