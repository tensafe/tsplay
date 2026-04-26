package tsplay_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

var workbenchUUIDPathPattern = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`)
var workbenchNumericPathPattern = regexp.MustCompile(`/\d+(/|$)`)

func normalizeWorkbenchRoute(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return strings.TrimSpace(rawURL)
	}
	path := parsed.EscapedPath()
	if path == "" {
		path = "/"
	}
	path = workbenchUUIDPathPattern.ReplaceAllString(path, "{id}")
	path = workbenchNumericPathPattern.ReplaceAllString(path, "/{id}$1")
	return path
}

func workbenchPathTemplate(rawURL string) string {
	return normalizeWorkbenchRoute(rawURL)
}

func classifyWorkbenchActionRisk(label string) string {
	value := strings.ToLower(strings.TrimSpace(label))
	switch {
	case containsAny(value, "删除", "移除", "禁用", "停用", "发布", "审批", "通过", "驳回", "支付", "转账", "修改密码"):
		return "write_high"
	case containsAny(value, "保存", "提交", "确认", "创建", "新增", "编辑", "修改", "更新", "备注"):
		return "write_low"
	case containsAny(value, "导出", "下载"):
		return "read_download"
	default:
		return "read"
	}
}

func classifyWorkbenchAPIRisk(method string, path string) string {
	method = strings.ToUpper(strings.TrimSpace(method))
	pathValue := strings.ToLower(strings.TrimSpace(path))
	switch {
	case method == "GET" || method == "HEAD":
		if containsAny(pathValue, "export", "download", "导出", "下载") {
			return "read_download"
		}
		return "read"
	case containsAny(pathValue, "export", "download", "导出", "下载"):
		return "read_download"
	case containsAny(pathValue, "delete", "remove", "disable", "approve", "publish", "payment", "pay", "删除", "禁用", "审批", "发布", "支付"):
		return "write_high"
	default:
		return "write_low"
	}
}

func workbenchOperationTypeFromRisk(risk string) string {
	switch strings.TrimSpace(risk) {
	case "read", "read_download":
		return "read"
	default:
		return "write"
	}
}

func inferWorkbenchSchemaFromText(body string, contentType string) any {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if strings.Contains(contentType, "json") || strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[") {
		var value any
		if json.Unmarshal([]byte(body), &value) == nil {
			return inferWorkbenchSchemaFromValue(value)
		}
	}
	values, err := url.ParseQuery(body)
	if err == nil && len(values) > 0 {
		result := map[string]any{}
		for key, items := range values {
			if len(items) == 1 {
				result[key] = inferWorkbenchScalarType(items[0])
			} else {
				itemTypes := []any{}
				for _, item := range items {
					itemTypes = append(itemTypes, inferWorkbenchScalarType(item))
				}
				result[key] = itemTypes
			}
		}
		return result
	}
	return map[string]any{"type": "string"}
}

func inferWorkbenchSchemaFromBytes(body []byte, contentType string) any {
	if len(body) == 0 {
		return nil
	}
	return inferWorkbenchSchemaFromText(string(body), contentType)
}

func inferWorkbenchSchemaFromValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return inferWorkbenchScalarType(typed)
	case float64:
		if float64(int64(typed)) == typed {
			return "integer"
		}
		return "number"
	case int, int8, int16, int32, int64:
		return "integer"
	case uint, uint8, uint16, uint32, uint64:
		return "integer"
	case []any:
		if len(typed) == 0 {
			return []any{"any"}
		}
		return []any{inferWorkbenchSchemaFromValue(typed[0])}
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		result := map[string]any{}
		for _, key := range keys {
			result[key] = inferWorkbenchSchemaFromValue(typed[key])
		}
		return result
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%T", value)
		}
		var decoded any
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			return fmt.Sprintf("%T", value)
		}
		return inferWorkbenchSchemaFromValue(decoded)
	}
}

func inferWorkbenchScalarType(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "string"
	}
	if matched, _ := regexp.MatchString(`^-?\d+$`, value); matched {
		return "integer"
	}
	if matched, _ := regexp.MatchString(`^-?\d+\.\d+$`, value); matched {
		return "number"
	}
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, value); matched {
		return "date"
	}
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}[tT ]\d{2}:\d{2}`, value); matched {
		return "datetime"
	}
	if strings.EqualFold(value, "true") || strings.EqualFold(value, "false") {
		return "boolean"
	}
	return "string"
}

func flattenWorkbenchPageCard(card WorkbenchPageCard) string {
	parts := []string{card.Title, card.Summary, card.NormalizedRoute, strings.Join(card.MenuPath, " "), strings.Join(card.Breadcrumbs, " ")}
	for _, form := range card.Forms {
		parts = append(parts, form.Name)
		for _, field := range form.Fields {
			parts = append(parts, field.Name, field.Label)
		}
	}
	for _, table := range card.Tables {
		parts = append(parts, table.Name, strings.Join(table.Columns, " "))
	}
	for _, action := range card.Actions {
		parts = append(parts, action.Label, action.Kind)
	}
	return strings.Join(parts, " ")
}

func flattenWorkbenchAPICard(card WorkbenchAPICard) string {
	parts := []string{card.SemanticName, card.Method, card.PathTemplate, card.TriggerRoute, card.TriggerAction}
	return strings.Join(parts, " ")
}

func scoreWorkbenchText(intent string, corpus string) int {
	intent = strings.ToLower(strings.TrimSpace(intent))
	corpus = strings.ToLower(strings.TrimSpace(corpus))
	if intent == "" || corpus == "" {
		return 0
	}
	score := 0
	if strings.Contains(intent, corpus) && len([]rune(corpus)) > 1 {
		score += 20
	}
	for _, token := range workbenchKeywordTokens(corpus) {
		if len([]rune(token)) < 2 {
			continue
		}
		if strings.Contains(intent, token) {
			score += 3
		}
	}
	return score
}

func workbenchKeywordTokens(value string) []string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return nil
	}
	seen := map[string]struct{}{}
	tokens := []string{}
	current := strings.Builder{}
	flush := func() {
		token := strings.TrimSpace(current.String())
		current.Reset()
		if token == "" {
			return
		}
		if _, ok := seen[token]; ok {
			return
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	for _, r := range value {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			current.WriteRune(r)
		case unicode.In(r, unicode.Han):
			current.WriteRune(r)
		default:
			flush()
		}
	}
	flush()

	runes := []rune(value)
	for size := 2; size <= 4; size++ {
		if len(runes) < size {
			continue
		}
		for i := 0; i+size <= len(runes); i++ {
			fragment := strings.TrimSpace(string(runes[i : i+size]))
			if fragment == "" {
				continue
			}
			if _, ok := seen[fragment]; ok {
				continue
			}
			if !containsHan(fragment) {
				continue
			}
			seen[fragment] = struct{}{}
			tokens = append(tokens, fragment)
		}
	}
	sort.Strings(tokens)
	return tokens
}

func containsHan(value string) bool {
	for _, r := range value {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}

func containsAny(value string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(value, strings.ToLower(strings.TrimSpace(keyword))) {
			return true
		}
	}
	return false
}

func deriveWorkbenchSemanticName(method string, pathTemplate string, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	base := filepath.Base(strings.TrimSpace(pathTemplate))
	base = strings.Trim(base, "/")
	base = strings.Trim(base, "{}")
	base = strings.ReplaceAll(base, "-", "_")
	if base == "" || base == "." {
		base = strings.ToLower(strings.TrimSpace(method))
	}
	return base
}

func encodeWorkbenchFlowYAML(flow *Flow) (string, error) {
	encoded, err := yaml.Marshal(flow)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
