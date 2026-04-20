package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tsplay/tsplay_core"
)

func runMCPToolAction(
	toolName string,
	argsJSON string,
	argsFile string,
	flowRoot string,
	artifactRoot string,
) error {
	arguments, err := loadMCPToolArguments(argsJSON, argsFile)
	if err != nil {
		return err
	}

	suppressedPayload, err := invokeMCPToolWithoutStdout(func() (map[string]any, error) {
		return tsplay_core.InvokeTSPlayTool(
			context.Background(),
			toolName,
			arguments,
			tsplay_core.TSPlayMCPServerOptions{
				FlowPathRoot: flowRoot,
				ArtifactRoot: artifactRoot,
			},
		)
	})
	if err != nil {
		return err
	}

	printJSON(suppressedPayload)
	return nil
}

func invokeMCPToolWithoutStdout(run func() (map[string]any, error)) (map[string]any, error) {
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	outputCh := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = io.Copy(&buffer, reader)
		outputCh <- buffer.String()
	}()

	os.Stdout = writer
	payload, runErr := run()
	_ = writer.Close()
	os.Stdout = originalStdout
	_ = reader.Close()
	<-outputCh

	if runErr != nil {
		return nil, runErr
	}
	return payload, nil
}

func loadMCPToolArguments(argsJSON string, argsFile string) (map[string]any, error) {
	if strings.TrimSpace(argsJSON) != "" && strings.TrimSpace(argsFile) != "" {
		return nil, fmt.Errorf("-args-json and -args-file are mutually exclusive")
	}

	raw := strings.TrimSpace(argsJSON)
	baseDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(argsFile) != "" {
		path, err := filepath.Abs(strings.TrimSpace(argsFile))
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read args file %q: %w", path, err)
		}
		raw = string(content)
		baseDir = filepath.Dir(path)
	}

	if strings.TrimSpace(raw) == "" {
		return map[string]any{}, nil
	}

	var arguments map[string]any
	if err := json.Unmarshal([]byte(raw), &arguments); err != nil {
		return nil, fmt.Errorf("args must be a JSON object: %w", err)
	}

	resolved, err := resolveMCPToolArgumentValue(arguments, baseDir)
	if err != nil {
		return nil, err
	}
	result, ok := resolved.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("resolved MCP arguments must stay a JSON object")
	}
	return result, nil
}

func resolveMCPToolArgumentValue(value any, baseDir string) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		resolved := make(map[string]any, len(typed))
		for key, nested := range typed {
			value, err := resolveMCPToolArgumentValue(nested, baseDir)
			if err != nil {
				return nil, err
			}
			resolved[key] = value
		}
		return resolved, nil
	case []any:
		resolved := make([]any, 0, len(typed))
		for _, nested := range typed {
			value, err := resolveMCPToolArgumentValue(nested, baseDir)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, value)
		}
		return resolved, nil
	case string:
		return resolveMCPToolArgumentString(typed, baseDir)
	default:
		return value, nil
	}
}

func resolveMCPToolArgumentString(value string, baseDir string) (any, error) {
	switch {
	case strings.HasPrefix(value, "@file:"):
		path, err := resolveMCPToolArgumentPath(strings.TrimPrefix(value, "@file:"), baseDir)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read @file %q: %w", path, err)
		}
		return string(content), nil
	case strings.HasPrefix(value, "@jsonfile:"):
		path, err := resolveMCPToolArgumentPath(strings.TrimPrefix(value, "@jsonfile:"), baseDir)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read @jsonfile %q: %w", path, err)
		}
		var decoded any
		if err := json.Unmarshal(content, &decoded); err != nil {
			return nil, fmt.Errorf("decode @jsonfile %q: %w", path, err)
		}
		return resolveMCPToolArgumentValue(decoded, filepath.Dir(path))
	case strings.HasPrefix(value, "@jsonpathfile:"):
		spec := strings.TrimPrefix(value, "@jsonpathfile:")
		fileSpec, jsonPath, ok := strings.Cut(spec, "#")
		if !ok || strings.TrimSpace(fileSpec) == "" || strings.TrimSpace(jsonPath) == "" {
			return nil, fmt.Errorf("@jsonpathfile requires the form @jsonpathfile:path#field.path")
		}
		path, err := resolveMCPToolArgumentPath(fileSpec, baseDir)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read @jsonpathfile %q: %w", path, err)
		}
		var decoded any
		if err := json.Unmarshal(content, &decoded); err != nil {
			return nil, fmt.Errorf("decode @jsonpathfile %q: %w", path, err)
		}
		selected, err := selectMCPToolArgumentJSONPath(decoded, jsonPath)
		if err != nil {
			return nil, fmt.Errorf("resolve @jsonpathfile %q: %w", path, err)
		}
		return resolveMCPToolArgumentValue(selected, filepath.Dir(path))
	default:
		return value, nil
	}
}

func resolveMCPToolArgumentPath(path string, baseDir string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if filepath.IsAbs(path) {
		return path, nil
	}
	if strings.TrimSpace(baseDir) == "" {
		return filepath.Abs(path)
	}
	return filepath.Abs(filepath.Join(baseDir, path))
}

func selectMCPToolArgumentJSONPath(value any, path string) (any, error) {
	current := value
	for _, part := range strings.Split(strings.TrimSpace(path), ".") {
		if strings.TrimSpace(part) == "" {
			return nil, fmt.Errorf("json path %q contains an empty segment", path)
		}
		switch typed := current.(type) {
		case map[string]any:
			next, ok := typed[part]
			if !ok {
				return nil, fmt.Errorf("field %q is missing", part)
			}
			current = next
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("segment %q must be an array index", part)
			}
			if index < 0 || index >= len(typed) {
				return nil, fmt.Errorf("array index %d is out of range", index)
			}
			current = typed[index]
		default:
			return nil, fmt.Errorf("segment %q cannot be applied to %T", part, current)
		}
	}
	return current, nil
}
