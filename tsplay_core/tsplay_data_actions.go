package tsplay_core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func write_json(L *lua.LState) int {
	filePath := L.CheckString(1)
	value := luaValueToGo(L.Get(2))

	result, err := writeJSONValue(filePath, value)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func write_csv(L *lua.LState) int {
	filePath := L.CheckString(1)
	value := luaValueToGo(L.Get(2))

	var headers []string
	if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
		var err error
		headers, err = stringListValue(luaValueToGo(L.Get(3)))
		if err != nil {
			L.RaiseError("write_csv headers %v", err)
			return 0
		}
	}

	result, err := writeCSVValue(filePath, value, headers)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func writeJSONValue(filePath string, value any) (map[string]any, error) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("write_json marshal %q: %w", filePath, err)
	}
	encoded = append(encoded, '\n')
	if err := writeOutputFile(filePath, encoded); err != nil {
		return nil, fmt.Errorf("write_json write %q: %w", filePath, err)
	}
	return map[string]any{
		"file_path": filePath,
		"bytes":     len(encoded),
	}, nil
}

func writeCSVValue(filePath string, value any, headers []string) (map[string]any, error) {
	rows, resolvedHeaders, err := csvRowsFromValue(value, headers)
	if err != nil {
		return nil, fmt.Errorf("write_csv normalize %q: %w", filePath, err)
	}

	if err := writeOutputFile(filePath, nil); err != nil {
		return nil, fmt.Errorf("write_csv prepare %q: %w", filePath, err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("write_csv create %q: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if len(resolvedHeaders) > 0 {
		if err := writer.Write(resolvedHeaders); err != nil {
			return nil, fmt.Errorf("write_csv header %q: %w", filePath, err)
		}
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("write_csv row %q: %w", filePath, err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("write_csv flush %q: %w", filePath, err)
	}

	columns := 0
	if len(resolvedHeaders) > 0 {
		columns = len(resolvedHeaders)
	} else {
		for _, row := range rows {
			if len(row) > columns {
				columns = len(row)
			}
		}
	}
	return map[string]any{
		"file_path": filePath,
		"rows":      len(rows),
		"columns":   columns,
	}, nil
}

func writeOutputFile(filePath string, content []byte) error {
	parent := filepath.Dir(filePath)
	if parent != "." && parent != "" {
		if err := os.MkdirAll(parent, 0755); err != nil {
			return err
		}
	}
	if content == nil {
		return nil
	}
	return os.WriteFile(filePath, content, 0644)
}

func csvRowsFromValue(value any, headers []string) ([][]string, []string, error) {
	rows := normalizeCSVTopLevel(value)
	resolvedHeaders := append([]string(nil), headers...)
	if len(resolvedHeaders) == 0 && csvRowsContainObjects(rows) {
		resolvedHeaders = deriveCSVHeaders(rows)
	}

	records := make([][]string, 0, len(rows))
	for _, row := range rows {
		record, err := csvRecordFromValue(row, resolvedHeaders)
		if err != nil {
			return nil, nil, err
		}
		records = append(records, record)
	}
	return records, resolvedHeaders, nil
}

func normalizeCSVTopLevel(value any) []any {
	switch typed := value.(type) {
	case nil:
		return []any{}
	case []any:
		return typed
	case []string:
		rows := make([]any, 0, len(typed))
		for _, item := range typed {
			rows = append(rows, item)
		}
		return rows
	default:
		return []any{typed}
	}
}

func csvRowsContainObjects(rows []any) bool {
	for _, row := range rows {
		if _, ok := row.(map[string]any); ok {
			return true
		}
	}
	return false
}

func deriveCSVHeaders(rows []any) []string {
	keys := map[string]struct{}{}
	for _, row := range rows {
		record, ok := row.(map[string]any)
		if !ok {
			continue
		}
		for key := range record {
			keys[key] = struct{}{}
		}
	}
	headers := make([]string, 0, len(keys))
	for key := range keys {
		headers = append(headers, key)
	}
	sort.Strings(headers)
	return headers
}

func csvRecordFromValue(value any, headers []string) ([]string, error) {
	switch typed := value.(type) {
	case map[string]any:
		record := make([]string, 0, len(headers))
		for _, header := range headers {
			record = append(record, stringifyCSVCell(typed[header]))
		}
		return record, nil
	case []any:
		record := make([]string, 0, len(typed))
		for _, item := range typed {
			record = append(record, stringifyCSVCell(item))
		}
		return record, nil
	case []string:
		record := make([]string, 0, len(typed))
		for _, item := range typed {
			record = append(record, item)
		}
		return record, nil
	default:
		return []string{stringifyCSVCell(typed)}, nil
	}
}

func stringifyCSVCell(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case []any, []string, map[string]any:
		encoded, err := json.Marshal(typed)
		if err == nil {
			return string(encoded)
		}
	}
	return strings.TrimSpace(fmt.Sprint(value))
}
