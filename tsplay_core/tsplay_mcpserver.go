package tsplay_core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const MCP_TINY_IMAGE = "iVBORw0KGgoAAAANSUhEUgAAARgAAAEYCAIAAAAI7H7bAAAZyUlEQVR4nOzce1RVZd4H8MM5BwERQUDxQpCoI0RajDWjomSEkOaltDBvaaIVy5aJltNkadkSdXJoWs6IKZko6bh0aABXxDTCKFgwgwalOKCICiJyEY7cz+Fw3rV63nnWb/a5eNSfWNP389fZt2dvNvu797Of5zlHazKZVABwZ9T3+gAA/hcgSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABho7/UBwM9L9w9M/43OkZ/FyhaXqlQqOp+uJrYy/qCrq0t87urqMhqN3d3dKpWq6wdiUi7t6uoSJZhvJRaZTCYxKTY0Go0eHh7Lly/v06eP+LsQpJ8vcZUYDAb9D8SFJSfF5SU+GwwGcQnq/0NuaDAYxIaKRWKp0Wg0mUzyYqUXrtFoFBe9nJRXv7hY5YaKRWJDOikS0pO8vLwyMzNlin56QZJ3I4vzzT/f6srimuj6D/n/MxgM8o5lMBjkZSEW0f863Zbe6hRligLpYciixFJ6uSgORnH7VCxSXLt0qVikOI2KU2r/pO01/1e5uLjMmzfv9ddfDwwMpPNvEiSDwXD06FHxH6VPUvn0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2rVrtPNq4MCB1tqIc3Nzb+Xf/l/i4uIslqnoJzhz5gz9lgSvkJAQnU73u9/97uDBg6J6cuzYsfj4+NWrV+fk5BQVFcmjio+PF52tPj4+u3btMr9a7om7HqScnBxaSTPvWBDkqIIDBw6YL1X8uMK6deto/dDZ2dnaaOIvv/ySbrh79276ew/+/v4Wt1L0HcfGxup0OtpTPm/ePIsbKrpEXVxcTp8+be3MKL4PsmTJkhMnTtBOgtDQ0JaWFtoxOnbsWGulKb6Lbj+tViu/fqeg+GJlSkqK4jnMZfTo0eIYXnrppdLS0uLiYrVaffXq1V//+tfiZwLE92EbGhrEC3avXr0WLFhAv3B5z931xgbR2iM+9+nTZ9y4cRZXmzJlimg8ffrpp82XKt6RCgoK6IDr8PBwa+PKFJX4oqIi2ixm7UdU5s6dS4cL7dmzp6SkhDbjWvyJQJVK9eSTT7722mtysr29feHChdaGOSve1oYOHXr+/HnaRPHwww8rvvwTFBRksSgxtMzaItvmzZtn8Yux5iNNs7KybPzy220QnXjbt28/ceKEOIby8vIPP/xQ9G5t3LgxODj4hRdeWLx4cUVFRV5e3iOPPPLVV1+99dZbpaWl+/btszik61656z/HJX6G9/935uBgrQ1H/OipWq221tZko3NGo9FYexURdUs5qVar6ZVqY3ei7VhOarVag8EgH61ardbar+SIn8+mc6z15yjWFH8CbZURe+no6JD7tfGX0vN8S+w/ew4ODFeLqDzLv06j0dCHsDjtckdardZoNDo4OOTm5s6YMaO1tfXw4cOzZ8++w2O4G/C7dvCj1tbW9u67737wwQdBQUFbt26dOnXqvT4iy35iP6IPPysVFRXR0dFarTYtLS0iIsLen5i7F36SP+oNPwepqamRkZELFiz45ptvZsyY8WNOEZ5I8GOUn58fFxcXGBiYl5d3S1/av4fwjgQ/LhcvXoyIiNi/f/89Gelz2/BEgh+XwsLCo0ePKoa//PjhiQTAAI0NAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAwQJAAGCBIAAwQJgAGCBMAAQQJggCABMECQABggSAAMECQABggSAAMECYABggTAAEECYIAgATBAkAAYIEgADBAkAAYIEgADBAmAAYIEwABBAmCAIAEwQJAAGCBIAAz+LwAA///FzJto8JNVBwAAAABJRU5ErkJggg=="

func McpServerSSE() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "address to listen on")
	flag.Parse()

	mcpServer := server.NewMCPServer("dynamic-path-example", "1.0.0")

	// Add a trivial tool for demonstration
	mcpServer.AddTool(mcp.NewTool("echo"), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Println(req.GetArguments())
		return mcp.NewToolResultText(fmt.Sprintf("Echo: %v", req.GetArguments()["message"])), nil
	})

	mcpServer.AddTool(mcp.NewTool(string("open_excel"),
		mcp.WithDescription("打开Excel文件"),
		mcp.WithString("file_path",
			mcp.Description("打开Excel文件"),
			mcp.Required(),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(fmt.Sprintf("Echo: %v", req.GetArguments()["file_path"])), nil
	})

	// Use a dynamic base path based on a path parameter (Go 1.22+)
	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithDynamicBasePath(func(r *http.Request, sessionID string) string {
			tenant := r.PathValue("tenant")
			return "/api/" + tenant
		}),
		server.WithBaseURL(fmt.Sprintf("http://localhost%s", addr)),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/{tenant}/sse", sseServer.SSEHandler())
	mux.Handle("/api/{tenant}/message", sseServer.MessageHandler())

	log.Printf("Dynamic SSE server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type ToolName string

const (
	ECHO                   ToolName = "echo"
	ADD                    ToolName = "add"
	LONG_RUNNING_OPERATION ToolName = "longRunningOperation"
	SAMPLE_LLM             ToolName = "sampleLLM"
	GET_TINY_IMAGE         ToolName = "getTinyImage"
)

type PromptName string

const (
	SIMPLE  PromptName = "simple_prompt"
	COMPLEX PromptName = "complex_prompt"
)

const DefaultMCPFlowPathRoot = "script"
const DefaultMCPArtifactRoot = DefaultFlowArtifactRoot

type TSPlayMCPServerOptions struct {
	FlowPathRoot string
	ArtifactRoot string
}

func DefaultTSPlayMCPServerOptions() TSPlayMCPServerOptions {
	return TSPlayMCPServerOptions{
		FlowPathRoot: DefaultMCPFlowPathRoot,
		ArtifactRoot: DefaultMCPArtifactRoot,
	}
}

func normalizeTSPlayMCPServerOptions(options []TSPlayMCPServerOptions) TSPlayMCPServerOptions {
	normalized := DefaultTSPlayMCPServerOptions()
	if len(options) == 0 {
		return normalized
	}
	if options[0].FlowPathRoot != "" {
		normalized.FlowPathRoot = options[0].FlowPathRoot
	}
	if options[0].ArtifactRoot != "" {
		normalized.ArtifactRoot = options[0].ArtifactRoot
	}
	return normalized
}

func NewTSPlayMCPServer(options ...TSPlayMCPServerOptions) *server.MCPServer {
	normalizedOptions := normalizeTSPlayMCPServerOptions(options)
	mcpServer := server.NewMCPServer(
		"tsplay",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	registerTSPlayFlowTools(mcpServer, normalizedOptions)
	return mcpServer
}

// NewMCPServer keeps the upstream example server available for local MCP demos.
// Production TSPlay/OpenClaw integrations should use NewTSPlayMCPServer.
func NewMCPServer() *server.MCPServer {
	hooks := &server.Hooks{}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("beforeAny: %s, %v, %v\n", method, id, message)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		fmt.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		fmt.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		fmt.Printf("beforeInitialize: %v, %v\n", id, message)
	})
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		fmt.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
		// authorization verification and other preprocessing tasks are performed.
		return nil
	})
	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		fmt.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		fmt.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		fmt.Printf("beforeCallTool: %v, %v\n", id, message)
	})

	mcpServer := server.NewMCPServer(
		"tsplay-demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	mcpServer.AddResource(mcp.NewResource("test://static/resource",
		"Static Resource",
		mcp.WithMIMEType("text/plain"),
	), handleReadResource)
	mcpServer.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"test://dynamic/resource/{id}",
			"Dynamic Resource",
		),
		handleResourceTemplate,
	)

	resources := generateResources()
	for _, resource := range resources {
		mcpServer.AddResource(resource, handleGeneratedResource)
	}

	mcpServer.AddPrompt(mcp.NewPrompt(string(SIMPLE),
		mcp.WithPromptDescription("A simple prompt"),
	), handleSimplePrompt)
	mcpServer.AddPrompt(mcp.NewPrompt(string(COMPLEX),
		mcp.WithPromptDescription("A complex prompt"),
		mcp.WithArgument("temperature",
			mcp.ArgumentDescription("The temperature parameter for generation"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("style",
			mcp.ArgumentDescription("The style to use for the response"),
			mcp.RequiredArgument(),
		),
	), handleComplexPrompt)
	mcpServer.AddTool(mcp.NewTool(string(ECHO),
		mcp.WithDescription("Echoes back the input"),
		mcp.WithString("message",
			mcp.Description("Message to echo"),
			mcp.Required(),
		),
	), handleEchoTool)

	mcpServer.AddTool(
		mcp.NewTool("notify"),
		handleSendNotification,
	)

	mcpServer.AddTool(mcp.NewTool(string(ADD),
		mcp.WithDescription("Adds two numbers"),
		mcp.WithNumber("a",
			mcp.Description("First number"),
			mcp.Required(),
		),
		mcp.WithNumber("b",
			mcp.Description("Second number"),
			mcp.Required(),
		),
	), handleAddTool)
	mcpServer.AddTool(mcp.NewTool(
		string(LONG_RUNNING_OPERATION),
		mcp.WithDescription(
			"Demonstrates a long running operation with progress updates",
		),
		mcp.WithNumber("duration",
			mcp.Description("Duration of the operation in seconds"),
			mcp.DefaultNumber(10),
		),
		mcp.WithNumber("steps",
			mcp.Description("Number of steps in the operation"),
			mcp.DefaultNumber(5),
		),
	), handleLongRunningOperationTool)

	// s.server.AddTool(mcp.Tool{
	// 	Name:        string(SAMPLE_LLM),
	// 	Description: "Samples from an LLM using MCP's sampling feature",
	// 	InputSchema: mcp.ToolInputSchema{
	// 		Type: "object",
	// 		Properties: map[string]any{
	// 			"prompt": map[string]any{
	// 				"type":        "string",
	// 				"description": "The prompt to send to the LLM",
	// 			},
	// 			"maxTokens": map[string]any{
	// 				"type":        "number",
	// 				"description": "Maximum number of tokens to generate",
	// 				"default":     100,
	// 			},
	// 		},
	// 	},
	// }, s.handleSampleLLMTool)
	mcpServer.AddTool(mcp.NewTool(string(GET_TINY_IMAGE),
		mcp.WithDescription("Returns the MCP_TINY_IMAGE"),
	), handleGetTinyImageTool)

	mcpServer.AddNotificationHandler("notification", handleNotification)

	return mcpServer
}

func registerTSPlayFlowTools(mcpServer *server.MCPServer, options TSPlayMCPServerOptions) {
	mcpServer.AddTool(mcp.NewTool("tsplay.list_actions",
		mcp.WithDescription("List structured TSPlay Flow actions and their argument schema."),
		mcp.WithReadOnlyHintAnnotation(true),
	), handleFlowListActionsTool)

	mcpServer.AddTool(mcp.NewTool("tsplay.flow_schema",
		mcp.WithDescription("Return the JSON Schema and generation rules for TSPlay Flow. Use this before generating or repairing flows."),
		mcp.WithReadOnlyHintAnnotation(true),
	), handleFlowSchemaTool)

	mcpServer.AddTool(mcp.NewTool("tsplay.flow_examples",
		mcp.WithDescription("Return typical TSPlay Flow examples for AI generation and repair."),
		mcp.WithReadOnlyHintAnnotation(true),
	), handleFlowExamplesTool)

	mcpServer.AddTool(mcp.NewTool("tsplay.draft_flow",
		mcp.WithDescription("Draft a TSPlay Flow from user intent plus page observation. If observation is omitted, the tool will open the page first."),
		mcp.WithString("intent",
			mcp.Description("User intent in natural language, for example 搜索订单并导出 or upload a file and submit."),
			mcp.Required(),
		),
		mcp.WithString("url",
			mcp.Description("Optional page URL. Required when observation is not provided."),
		),
		mcp.WithString("observation",
			mcp.Description("Optional PageObservation JSON, or a wrapper that contains an observation field."),
		),
		mcp.WithString("flow_name",
			mcp.Description("Optional explicit flow name. Defaults to an auto-generated draft name."),
		),
		mcp.WithBoolean("headless",
			mcp.Description("Run browser in headless mode when url is provided. Defaults to true."),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Navigation timeout in milliseconds when url is provided. Defaults to 30000."),
		),
		mcp.WithNumber("max_elements",
			mcp.Description("Maximum interactive elements to observe when url is provided. Defaults to 100."),
		),
		mcp.WithOpenWorldHintAnnotation(true),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleDraftFlowToolWithOptions(ctx, request, options)
	})

	mcpServer.AddTool(mcp.NewTool("tsplay.repair_flow_context",
		mcp.WithDescription("Build an AI-friendly repair context from a Flow and failed run trace. Returns summaries and artifact paths without embedding full HTML."),
		mcp.WithString("flow",
			mcp.Description("Flow content as YAML or JSON. Use this or flow_path."),
		),
		mcp.WithString("flow_path",
			mcp.Description("Flow YAML or JSON file path relative to the configured flow root. Use this or flow."),
		),
		mcp.WithString("format",
			mcp.Description("Optional format hint: yaml or json."),
			mcp.Enum("yaml", "json"),
		),
		mcp.WithString("run_result",
			mcp.Description("JSON returned by tsplay.run_flow, or its result field. Preferred input."),
		),
		mcp.WithString("trace",
			mcp.Description("JSON trace array when run_result is not available."),
		),
		mcp.WithString("error",
			mcp.Description("Optional top-level error message from the failed run."),
		),
		mcp.WithNumber("max_artifact_excerpt",
			mcp.Description("Maximum characters of simplified DOM snapshot to include. Defaults to 4000."),
		),
		mcp.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleRepairFlowContextToolWithOptions(ctx, request, options)
	})

	mcpServer.AddTool(mcp.NewTool("tsplay.observe_page",
		mcp.WithDescription("Open a page and return an AI-friendly observation: screenshot path, DOM snapshot path, and interactive elements with selector candidates."),
		mcp.WithString("url",
			mcp.Description("URL to open and observe."),
			mcp.Required(),
		),
		mcp.WithBoolean("headless",
			mcp.Description("Run browser in headless mode. Defaults to true."),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Navigation timeout in milliseconds. Defaults to 30000."),
		),
		mcp.WithNumber("max_elements",
			mcp.Description("Maximum interactive elements to return. Defaults to 100."),
		),
		mcp.WithOpenWorldHintAnnotation(true),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleObservePageToolWithOptions(ctx, request, options)
	})

	mcpServer.AddTool(mcp.NewTool("tsplay.validate_flow",
		mcp.WithDescription("Validate a TSPlay Flow YAML or JSON document without launching a browser."),
		mcp.WithString("flow",
			mcp.Description("Flow content as YAML or JSON. Use this or flow_path."),
		),
		mcp.WithString("flow_path",
			mcp.Description("Flow YAML or JSON file path relative to the configured flow root. Use this or flow."),
		),
		mcp.WithString("format",
			mcp.Description("Optional format hint: yaml or json."),
			mcp.Enum("yaml", "json"),
		),
		mcp.WithBoolean("allow_lua",
			mcp.Description("Allow lua steps for this request. Defaults to false."),
		),
		mcp.WithBoolean("allow_javascript",
			mcp.Description("Allow execute_script/evaluate steps for this request. Defaults to false."),
		),
		mcp.WithBoolean("allow_file_access",
			mcp.Description("Allow local file read/write actions for this request. Defaults to false. File paths are constrained to the configured artifact root."),
		),
		mcp.WithBoolean("allow_browser_state",
			mcp.Description("Allow browser storage/cookie export actions for this request. Defaults to false."),
		),
		mcp.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleValidateFlowToolWithOptions(ctx, request, options)
	})

	mcpServer.AddTool(mcp.NewTool("tsplay.run_flow",
		mcp.WithDescription("Run a TSPlay Flow YAML or JSON document in Playwright and return the execution trace."),
		mcp.WithString("flow",
			mcp.Description("Flow content as YAML or JSON. Use this or flow_path."),
		),
		mcp.WithString("flow_path",
			mcp.Description("Flow YAML or JSON file path relative to the configured flow root. Use this or flow."),
		),
		mcp.WithString("format",
			mcp.Description("Optional format hint: yaml or json."),
			mcp.Enum("yaml", "json"),
		),
		mcp.WithBoolean("headless",
			mcp.Description("Run browser in headless mode. Defaults to true."),
		),
		mcp.WithBoolean("allow_lua",
			mcp.Description("Allow lua steps for this request. Defaults to false."),
		),
		mcp.WithBoolean("allow_javascript",
			mcp.Description("Allow execute_script/evaluate steps for this request. Defaults to false."),
		),
		mcp.WithBoolean("allow_file_access",
			mcp.Description("Allow local file read/write actions for this request. Defaults to false. File paths are constrained to the configured artifact root."),
		),
		mcp.WithBoolean("allow_browser_state",
			mcp.Description("Allow browser storage/cookie export actions for this request. Defaults to false."),
		),
		mcp.WithOpenWorldHintAnnotation(true),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleRunFlowToolWithOptions(ctx, request, options)
	})
}

func handleFlowListActionsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return newJSONToolResult(map[string]any{
		"actions": buildFlowActionManifest(),
	})
}

func handleFlowSchemaTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return newJSONToolResult(map[string]any{
		"schema":              BuildFlowJSONSchema(),
		"action_manifest":     buildFlowActionManifest(),
		"generation_rules":    flowSchemaGenerationRules(),
		"selector_strategy":   flowSelectorStrategy(),
		"authoring_checklist": flowAuthoringChecklist(),
		"repair_checklist":    flowRepairValidationChecklist(),
	})
}

func handleFlowExamplesTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return newJSONToolResult(map[string]any{
		"examples":                BuildFlowExamples(),
		"example_selection_hints": flowExampleSelectionHints(),
	})
}

func handleDraftFlowTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return handleDraftFlowToolWithOptions(ctx, request, DefaultTSPlayMCPServerOptions())
}

func handleDraftFlowToolWithOptions(
	ctx context.Context,
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (*mcp.CallToolResult, error) {
	intent := strings.TrimSpace(request.GetString("intent", ""))
	if intent == "" {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": "intent is required",
		})
	}

	observation, err := ParseObservationForDraft(request.GetString("observation", ""))
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}

	url := strings.TrimSpace(request.GetString("url", ""))
	if observation == nil {
		if url == "" {
			return newJSONToolResult(map[string]any{
				"ok":    false,
				"error": "url or observation is required",
			})
		}
		observation, err = ObservePage(PageObservationOptions{
			URL:          url,
			Headless:     request.GetBool("headless", true),
			ArtifactRoot: options.ArtifactRoot,
			TimeoutMS:    request.GetInt("timeout", 30000),
			MaxElements:  request.GetInt("max_elements", 100),
		})
		if err != nil {
			return newJSONToolResult(map[string]any{
				"ok":    false,
				"error": err.Error(),
			})
		}
	}

	draft, err := BuildDraftFlow(FlowDraftOptions{
		Intent:       intent,
		URL:          url,
		FlowName:     request.GetString("flow_name", ""),
		ArtifactRoot: options.ArtifactRoot,
		Observation:  observation,
	})
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}

	return newJSONToolResult(map[string]any{
		"ok":          true,
		"observation": observation,
		"draft":       draft,
	})
}

func handleRepairFlowContextTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return handleRepairFlowContextToolWithOptions(ctx, request, DefaultTSPlayMCPServerOptions())
}

func handleRepairFlowContextToolWithOptions(
	ctx context.Context,
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (*mcp.CallToolResult, error) {
	flow, err := flowFromToolRequestWithOptions(request, options)
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}

	result, runError, err := ParseFlowRunResultForRepair(
		request.GetString("run_result", ""),
		request.GetString("trace", ""),
	)
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}

	context, err := BuildFlowRepairContext(FlowRepairContextOptions{
		Flow:               flow,
		Result:             result,
		Error:              firstNonEmpty(request.GetString("error", ""), runError),
		ArtifactRoot:       options.ArtifactRoot,
		MaxArtifactExcerpt: request.GetInt("max_artifact_excerpt", defaultFlowRepairArtifactExcerpt),
	})
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}
	return newJSONToolResult(map[string]any{
		"ok":      true,
		"context": context,
	})
}

func handleObservePageTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return handleObservePageToolWithOptions(ctx, request, DefaultTSPlayMCPServerOptions())
}

func handleObservePageToolWithOptions(
	ctx context.Context,
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (*mcp.CallToolResult, error) {
	observation, err := ObservePage(PageObservationOptions{
		URL:          request.GetString("url", ""),
		Headless:     request.GetBool("headless", true),
		ArtifactRoot: options.ArtifactRoot,
		TimeoutMS:    request.GetInt("timeout", 30000),
		MaxElements:  request.GetInt("max_elements", 100),
	})
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}
	return newJSONToolResult(map[string]any{
		"ok":          true,
		"observation": observation,
	})
}

func handleValidateFlowTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return handleValidateFlowToolWithOptions(ctx, request, DefaultTSPlayMCPServerOptions())
}

func handleValidateFlowToolWithOptions(
	ctx context.Context,
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (*mcp.CallToolResult, error) {
	flow, err := flowFromToolRequestWithOptions(request, options)
	if err != nil {
		return newJSONToolResult(map[string]any{
			"valid": false,
			"error": err.Error(),
		})
	}
	security := flowSecurityPolicyFromToolRequest(request, options)
	if err := ValidateFlow(flow); err != nil {
		return newJSONToolResult(map[string]any{
			"valid": false,
			"name":  flow.Name,
			"error": err.Error(),
		})
	}
	if err := ValidateFlowSecurity(flow, security); err != nil {
		return newJSONToolResult(map[string]any{
			"valid": false,
			"name":  flow.Name,
			"error": err.Error(),
		})
	}
	return newJSONToolResult(map[string]any{
		"valid": true,
		"name":  flow.Name,
		"steps": len(flow.Steps),
	})
}

func handleRunFlowTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return handleRunFlowToolWithOptions(ctx, request, DefaultTSPlayMCPServerOptions())
}

func handleRunFlowToolWithOptions(
	ctx context.Context,
	request mcp.CallToolRequest,
	options TSPlayMCPServerOptions,
) (*mcp.CallToolResult, error) {
	flow, err := flowFromToolRequestWithOptions(request, options)
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}

	security := flowSecurityPolicyFromToolRequest(request, options)
	result, err := RunFlow(flow, FlowRunOptions{
		Headless:     request.GetBool("headless", true),
		Security:     &security,
		ArtifactRoot: options.ArtifactRoot,
	})
	if err != nil {
		return newJSONToolResult(map[string]any{
			"ok":     false,
			"error":  err.Error(),
			"result": flowResultForTool(result),
		})
	}
	return newJSONToolResult(map[string]any{
		"ok":     true,
		"result": flowResultForTool(result),
	})
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

func flowFromToolRequest(request mcp.CallToolRequest) (*Flow, error) {
	return flowFromToolRequestWithOptions(request, DefaultTSPlayMCPServerOptions())
}

func flowFromToolRequestWithOptions(request mcp.CallToolRequest, options TSPlayMCPServerOptions) (*Flow, error) {
	flowPath := request.GetString("flow_path", "")
	flowContent := request.GetString("flow", "")
	if flowPath != "" && flowContent != "" {
		return nil, fmt.Errorf("use either flow or flow_path, not both")
	}
	if flowPath != "" {
		resolvedPath, err := resolveMCPFlowPath(flowPath, options.FlowPathRoot)
		if err != nil {
			return nil, err
		}
		return LoadFlowFile(resolvedPath)
	}
	if flowContent == "" {
		return nil, fmt.Errorf("either flow or flow_path is required")
	}
	return ParseFlow([]byte(flowContent), request.GetString("format", "yaml"))
}

func flowSecurityPolicyFromToolRequest(request mcp.CallToolRequest, options TSPlayMCPServerOptions) FlowSecurityPolicy {
	return FlowSecurityPolicy{
		AllowLua:          request.GetBool("allow_lua", false),
		AllowJavaScript:   request.GetBool("allow_javascript", false),
		AllowFileAccess:   request.GetBool("allow_file_access", false),
		AllowBrowserState: request.GetBool("allow_browser_state", false),
		FileInputRoot:     options.ArtifactRoot,
		FileOutputRoot:    options.ArtifactRoot,
	}
}

func resolveMCPFlowPath(flowPath string, flowRoot string) (string, error) {
	if strings.TrimSpace(flowPath) == "" {
		return "", fmt.Errorf("flow_path is required")
	}
	if strings.TrimSpace(flowRoot) == "" {
		flowRoot = DefaultMCPFlowPathRoot
	}

	rootAbs, err := filepath.Abs(flowRoot)
	if err != nil {
		return "", fmt.Errorf("resolve flow root %q: %w", flowRoot, err)
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", fmt.Errorf("flow root %q is not accessible: %w", rootAbs, err)
	}

	candidate := flowPath
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(rootReal, candidate)
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve flow_path %q: %w", flowPath, err)
	}
	candidateReal, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", fmt.Errorf("flow_path %q is not accessible: %w", flowPath, err)
	}

	rel, err := filepath.Rel(rootReal, candidateReal)
	if err != nil {
		return "", fmt.Errorf("compare flow_path %q with flow root %q: %w", candidateReal, rootReal, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("flow_path %q is outside allowed flow root %q", flowPath, rootReal)
	}
	return candidateReal, nil
}

func buildFlowActionManifest() []map[string]any {
	descriptions := map[string]string{}
	for _, fn := range GlobalPlayWrightFunc {
		descriptions[fn.Name] = fn.Description_en
	}
	descriptions["lua"] = "Run an inline Lua code block. Prefer structured actions for normal browser steps and use lua only as an escape hatch."
	descriptions["extract_text"] = "Read text from a selector, optionally wait first, and optionally extract the first regex match."
	descriptions["set_var"] = "Set a flow variable from a resolved value. Requires save_as; for non-string literals use with.value."
	descriptions["assert_visible"] = "Fail the flow unless the selector is visible. Optional timeout waits before asserting."
	descriptions["assert_text"] = "Fail the flow unless the selected element text contains the expected text. Optional timeout polls before asserting."
	descriptions["retry"] = "Retry nested Flow steps until they succeed or the retry count is exhausted."
	descriptions["if"] = "Run then or else nested Flow steps based on a condition step output."
	descriptions["foreach"] = "Run nested Flow steps once for each item in a list."
	descriptions["on_error"] = "Run nested Flow steps and execute an error handler block if they fail."
	descriptions["wait_until"] = "Poll a condition step until it returns a truthy result or times out."

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
		if name == "extract_text" {
			item["returns"] = "string|string[]"
			item["encouraged_save_as"] = true
		}
		if name == "set_var" {
			item["requires_save_as"] = true
			item["returns"] = "any"
			item["notes"] = []string{
				"Use value for strings or placeholders such as {{order_count}}.",
				"Use with.value when the literal is a boolean, number, list, or object.",
			}
		}
		if name == "retry" {
			item["args"] = []map[string]any{
				{"name": "times", "type": "int", "required": false, "default": 3},
				{"name": "interval_ms", "type": "int", "required": false, "default": 0},
				{"name": "steps", "type": "steps", "required": true},
			}
		}
		if name == "if" {
			item["args"] = []map[string]any{
				{"name": "condition", "type": "condition", "required": true},
				{"name": "then", "type": "steps", "required": false},
				{"name": "else", "type": "steps", "required": false},
			}
		}
		if name == "foreach" {
			item["args"] = []map[string]any{
				{"name": "items", "type": "items", "required": true},
				{"name": "item_var", "type": "string", "required": true},
				{"name": "index_var", "type": "string", "required": false},
				{"name": "steps", "type": "steps", "required": true},
			}
		}
		if name == "on_error" {
			item["args"] = []map[string]any{
				{"name": "steps", "type": "steps", "required": true},
				{"name": "on_error", "type": "steps", "required": true},
			}
		}
		if name == "wait_until" {
			item["args"] = []map[string]any{
				{"name": "condition", "type": "condition", "required": true},
				{"name": "timeout", "type": "int", "required": false, "default": 30000},
				{"name": "interval_ms", "type": "int", "required": false, "default": 500},
			}
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

func newJSONToolResult(value any) (*mcp.CallToolResult, error) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(encoded)), nil
}

func generateResources() []mcp.Resource {
	resources := make([]mcp.Resource, 100)
	for i := 0; i < 100; i++ {
		uri := fmt.Sprintf("test://static/resource/%d", i+1)
		if i%2 == 0 {
			resources[i] = mcp.Resource{
				URI:      uri,
				Name:     fmt.Sprintf("Resource %d", i+1),
				MIMEType: "text/plain",
			}
		} else {
			resources[i] = mcp.Resource{
				URI:      uri,
				Name:     fmt.Sprintf("Resource %d", i+1),
				MIMEType: "application/octet-stream",
			}
		}
	}
	return resources
}

func handleReadResource(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "test://static/resource",
			MIMEType: "text/plain",
			Text:     "This is a sample resource",
		},
	}, nil
}

func handleResourceTemplate(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     "This is a sample resource",
		},
	}, nil
}

func handleGeneratedResource(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI

	var resourceNumber string
	if _, err := fmt.Sscanf(uri, "test://static/resource/%s", &resourceNumber); err != nil {
		return nil, fmt.Errorf("invalid resource URI format: %w", err)
	}

	num, err := strconv.Atoi(resourceNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid resource number: %w", err)
	}

	index := num - 1

	if index%2 == 0 {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "text/plain",
				Text:     fmt.Sprintf("Text content for resource %d", num),
			},
		}, nil
	} else {
		return []mcp.ResourceContents{
			mcp.BlobResourceContents{
				URI:      uri,
				MIMEType: "application/octet-stream",
				Blob:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Binary content for resource %d", num))),
			},
		}, nil
	}
}

func handleSimplePrompt(
	ctx context.Context,
	request mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "A simple prompt without arguments",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: "This is a simple prompt without arguments.",
				},
			},
		},
	}, nil
}

func handleComplexPrompt(
	ctx context.Context,
	request mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	arguments := request.Params.Arguments
	return &mcp.GetPromptResult{
		Description: "A complex prompt with arguments",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(
						"This is a complex prompt with arguments: temperature=%s, style=%s",
						arguments["temperature"],
						arguments["style"],
					),
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Type: "text",
					Text: "I understand. You've provided a complex prompt with temperature and style arguments. How would you like me to proceed?",
				},
			},
			{
				Role: mcp.RoleUser,
				Content: mcp.ImageContent{
					Type:     "image",
					Data:     MCP_TINY_IMAGE,
					MIMEType: "image/png",
				},
			},
		},
	}, nil
}

func handleEchoTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	message, ok := arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid message argument")
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Echo: %s", message),
			},
		},
	}, nil
}

func handleAddTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	a, ok1 := arguments["a"].(float64)
	b, ok2 := arguments["b"].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid number arguments")
	}
	sum := a + b
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("The sum of %f and %f is %f.", a, b, sum),
			},
		},
	}, nil
}

func handleSendNotification(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	server := server.ServerFromContext(ctx)

	err := server.SendNotificationToClient(
		ctx,
		"notifications/progress",
		map[string]any{
			"progress":      10,
			"total":         10,
			"progressToken": 0,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "notification sent successfully",
			},
		},
	}, nil
}

func handleLongRunningOperationTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	progressToken := "" //request.Params.Meta.ProgressToken
	duration, _ := arguments["duration"].(float64)
	steps, _ := arguments["steps"].(float64)
	stepDuration := duration / steps
	server := server.ServerFromContext(ctx)

	for i := 1; i < int(steps)+1; i++ {
		time.Sleep(time.Duration(stepDuration * float64(time.Second)))
		if progressToken != "" {
			err := server.SendNotificationToClient(
				ctx,
				"notifications/progress",
				map[string]any{
					"progress":      i,
					"total":         int(steps),
					"progressToken": progressToken,
					"message":       fmt.Sprintf("Server progress %v%%", int(float64(i)*100/steps)),
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to send notification: %w", err)
			}
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Long running operation completed. Duration: %f seconds, Steps: %d.",
					duration,
					int(steps),
				),
			},
		},
	}, nil
}

// func (s *MCPServer) handleSampleLLMTool(arguments map[string]any) (*mcp.CallToolResult, error) {
// 	prompt, _ := arguments["prompt"].(string)
// 	maxTokens, _ := arguments["maxTokens"].(float64)

// 	// This is a mock implementation. In a real scenario, you would use the server's RequestSampling method.
// 	result := fmt.Sprintf(
// 		"Sample LLM result for prompt: '%s' (max tokens: %d)",
// 		prompt,
// 		int(maxTokens),
// 	)

// 	return &mcp.CallToolResult{
// 		Content: []any{
// 			mcp.TextContent{
// 				Type: "text",
// 				Text: fmt.Sprintf("LLM sampling result: %s", result),
// 			},
// 		},
// 	}, nil
// }

func handleGetTinyImageTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "This is a tiny image:",
			},
			mcp.ImageContent{
				Type:     "image",
				Data:     MCP_TINY_IMAGE,
				MIMEType: "image/png",
			},
			mcp.TextContent{
				Type: "text",
				Text: "The image above is the MCP tiny image.",
			},
		},
	}, nil
}

func handleNotification(
	ctx context.Context,
	notification mcp.JSONRPCNotification,
) {
	log.Printf("Received notification: %s", notification.Method)
}

func McpServerMCP(addr string, options ...TSPlayMCPServerOptions) {
	if addr == "" {
		addr = ":8080"
	}

	normalizedOptions := normalizeTSPlayMCPServerOptions(options)
	mcpServer := NewTSPlayMCPServer(normalizedOptions)

	httpServer := server.NewStreamableHTTPServer(mcpServer)
	log.Printf("HTTP server listening on %s/mcp", addr)
	log.Printf("MCP flow_path root: %s", normalizedOptions.FlowPathRoot)
	log.Printf("MCP artifact root: %s", normalizedOptions.ArtifactRoot)
	if err := httpServer.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
