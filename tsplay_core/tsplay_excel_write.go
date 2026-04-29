package tsplay_core

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultExcelSheetName = "Sheet1"

var xmlAttributeEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

type excelWriteOptions struct {
	Headers []string
	Sheet   string
}

type excelWorkbookWriteData struct {
	Sheets []excelSheetWriteData
}

type excelSheetWriteData struct {
	Name     string
	Rows     [][]excelCellValue
	DataRows int
	Columns  int
}

type excelSheetInput struct {
	Name    string
	Value   any
	Headers []string
}

type excelCellValue struct {
	Kind  excelCellKind
	Value string
}

type excelCellKind string

const (
	excelCellBlank  excelCellKind = "blank"
	excelCellString excelCellKind = "string"
	excelCellNumber excelCellKind = "number"
	excelCellBool   excelCellKind = "bool"
)

type xlsxArchiveEntry struct {
	Name    string
	Content string
}

func applyExcelWriteOption(options *excelWriteOptions, value any) error {
	if options == nil || value == nil {
		return nil
	}
	switch typed := value.(type) {
	case string:
		if options.Sheet != "" {
			return fmt.Errorf("sheet was provided more than once")
		}
		options.Sheet = typed
		return nil
	case map[string]any:
		if rawSheet, ok := typed["sheet"]; ok {
			sheet, ok := rawSheet.(string)
			if !ok {
				return fmt.Errorf("sheet must be a string")
			}
			if options.Sheet != "" {
				return fmt.Errorf("sheet was provided more than once")
			}
			options.Sheet = sheet
		}
		if rawHeaders, ok := typed["headers"]; ok {
			headers, err := stringListValue(rawHeaders)
			if err != nil {
				return fmt.Errorf("headers %w", err)
			}
			if len(options.Headers) > 0 {
				return fmt.Errorf("headers were provided more than once")
			}
			options.Headers = headers
		}
		if _, ok := typed["sheets"]; ok {
			return fmt.Errorf("multi-sheet workbooks belong in value.sheets, not an extra option argument")
		}
		return nil
	default:
		headers, err := stringListValue(typed)
		if err != nil {
			return fmt.Errorf("must be a sheet name string, headers list, or options object: %w", err)
		}
		if len(options.Headers) > 0 {
			return fmt.Errorf("headers were provided more than once")
		}
		options.Headers = headers
		return nil
	}
}

func writeExcelValue(filePath string, value any, options excelWriteOptions) (map[string]any, error) {
	if strings.ToLower(filepath.Ext(filePath)) != ".xlsx" {
		return nil, fmt.Errorf("write_excel currently supports only .xlsx files")
	}

	workbook, err := excelWorkbookFromValue(value, options)
	if err != nil {
		return nil, fmt.Errorf("write_excel normalize %q: %w", filePath, err)
	}

	entries := buildXLSXArchiveEntries(workbook.Sheets)
	if err := writeXLSXArchive(filePath, entries); err != nil {
		return nil, fmt.Errorf("write_excel write %q: %w", filePath, err)
	}

	result := map[string]any{
		"file_path":   filePath,
		"sheet_count": len(workbook.Sheets),
	}
	totalRows := 0
	sheetSummaries := make([]any, 0, len(workbook.Sheets))
	for _, sheet := range workbook.Sheets {
		totalRows += sheet.DataRows
		sheetSummaries = append(sheetSummaries, map[string]any{
			"name":    sheet.Name,
			"rows":    sheet.DataRows,
			"columns": sheet.Columns,
		})
	}
	result["sheets"] = sheetSummaries
	if len(workbook.Sheets) == 1 {
		result["sheet"] = workbook.Sheets[0].Name
		result["rows"] = workbook.Sheets[0].DataRows
		result["columns"] = workbook.Sheets[0].Columns
	} else {
		result["rows"] = totalRows
	}
	if info, statErr := os.Stat(filePath); statErr == nil {
		result["bytes"] = info.Size()
	}
	return result, nil
}

func excelWorkbookFromValue(value any, options excelWriteOptions) (excelWorkbookWriteData, error) {
	workbookSheets, isWorkbook, err := workbookSheetInputsFromValue(value)
	if err != nil {
		return excelWorkbookWriteData{}, err
	}
	if isWorkbook {
		if strings.TrimSpace(options.Sheet) != "" || len(options.Headers) > 0 {
			return excelWorkbookWriteData{}, fmt.Errorf("cannot combine workbook value.sheets with top-level sheet or headers options")
		}
		return buildExcelWorkbookWriteData(workbookSheets)
	}

	sheetName, err := normalizeExcelSheetName(options.Sheet)
	if err != nil {
		return excelWorkbookWriteData{}, fmt.Errorf("options: %w", err)
	}
	sheet, err := buildExcelSheetWriteData(excelSheetInput{
		Name:    sheetName,
		Value:   value,
		Headers: options.Headers,
	})
	if err != nil {
		return excelWorkbookWriteData{}, err
	}
	return excelWorkbookWriteData{
		Sheets: []excelSheetWriteData{sheet},
	}, nil
}

func workbookSheetInputsFromValue(value any) ([]excelSheetInput, bool, error) {
	record, ok := value.(map[string]any)
	if !ok || len(record) != 1 {
		return nil, false, nil
	}
	rawSheets, ok := record["sheets"]
	if !ok {
		return nil, false, nil
	}
	items, err := anyListValue(rawSheets)
	if err != nil {
		return nil, true, fmt.Errorf("workbook sheets must be a list")
	}
	if len(items) == 0 {
		return nil, true, fmt.Errorf("workbook sheets must contain at least one sheet")
	}

	sheets := make([]excelSheetInput, 0, len(items))
	for index, item := range items {
		sheetMap, ok := item.(map[string]any)
		if !ok {
			return nil, true, fmt.Errorf("workbook sheet #%d must be an object", index+1)
		}
		name := ""
		if rawName, ok := sheetMap["name"]; ok {
			text, ok := rawName.(string)
			if !ok {
				return nil, true, fmt.Errorf("workbook sheet #%d name must be a string", index+1)
			}
			name = text
		} else if rawName, ok := sheetMap["sheet"]; ok {
			text, ok := rawName.(string)
			if !ok {
				return nil, true, fmt.Errorf("workbook sheet #%d sheet must be a string", index+1)
			}
			name = text
		}

		headers := []string(nil)
		if rawHeaders, ok := sheetMap["headers"]; ok {
			resolved, err := stringListValue(rawHeaders)
			if err != nil {
				return nil, true, fmt.Errorf("workbook sheet #%d headers %w", index+1, err)
			}
			headers = resolved
		}

		rawValue, hasValue := sheetMap["value"]
		rawRows, hasRows := sheetMap["rows"]
		if hasValue && hasRows {
			return nil, true, fmt.Errorf("workbook sheet #%d cannot specify both value and rows", index+1)
		}
		if !hasValue && !hasRows {
			return nil, true, fmt.Errorf("workbook sheet #%d requires value or rows", index+1)
		}
		if hasRows {
			rawValue = rawRows
		}

		sheets = append(sheets, excelSheetInput{
			Name:    name,
			Value:   rawValue,
			Headers: headers,
		})
	}
	return sheets, true, nil
}

func anyListValue(value any) ([]any, error) {
	switch typed := value.(type) {
	case nil:
		return nil, nil
	case []any:
		return typed, nil
	case []string:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("must be a list")
	}
}

func buildExcelWorkbookWriteData(inputs []excelSheetInput) (excelWorkbookWriteData, error) {
	sheets := make([]excelSheetWriteData, 0, len(inputs))
	usedNames := map[string]struct{}{}
	nextDefaultIndex := 1
	for _, input := range inputs {
		name, err := resolveWorkbookSheetName(input.Name, usedNames, &nextDefaultIndex)
		if err != nil {
			return excelWorkbookWriteData{}, err
		}
		input.Name = name
		sheet, err := buildExcelSheetWriteData(input)
		if err != nil {
			return excelWorkbookWriteData{}, err
		}
		sheets = append(sheets, sheet)
	}
	return excelWorkbookWriteData{Sheets: sheets}, nil
}

func resolveWorkbookSheetName(name string, usedNames map[string]struct{}, nextDefaultIndex *int) (string, error) {
	if strings.TrimSpace(name) != "" {
		normalized, err := normalizeExcelSheetName(name)
		if err != nil {
			return "", err
		}
		key := strings.ToLower(normalized)
		if _, exists := usedNames[key]; exists {
			return "", fmt.Errorf("sheet %q was provided more than once", normalized)
		}
		usedNames[key] = struct{}{}
		return normalized, nil
	}

	for {
		candidate := fmt.Sprintf("Sheet%d", *nextDefaultIndex)
		*nextDefaultIndex++
		key := strings.ToLower(candidate)
		if _, exists := usedNames[key]; exists {
			continue
		}
		usedNames[key] = struct{}{}
		return candidate, nil
	}
}

func buildExcelSheetWriteData(input excelSheetInput) (excelSheetWriteData, error) {
	rows := normalizeCSVTopLevel(input.Value)
	resolvedHeaders := append([]string(nil), input.Headers...)
	if len(resolvedHeaders) == 0 && csvRowsContainObjects(rows) {
		resolvedHeaders = deriveCSVHeaders(rows)
	}

	records := make([][]excelCellValue, 0, len(rows)+1)
	columns := len(resolvedHeaders)
	if len(resolvedHeaders) > 0 {
		headerRow := make([]excelCellValue, 0, len(resolvedHeaders))
		for _, header := range resolvedHeaders {
			headerRow = append(headerRow, excelStringCell(header))
		}
		records = append(records, headerRow)
	}

	for _, row := range rows {
		record, err := excelRecordFromValue(row, resolvedHeaders)
		if err != nil {
			return excelSheetWriteData{}, err
		}
		records = append(records, record)
		if len(record) > columns {
			columns = len(record)
		}
	}

	return excelSheetWriteData{
		Name:     input.Name,
		Rows:     records,
		DataRows: len(rows),
		Columns:  columns,
	}, nil
}

func excelRecordFromValue(value any, headers []string) ([]excelCellValue, error) {
	switch typed := value.(type) {
	case map[string]any:
		record := make([]excelCellValue, 0, len(headers))
		for _, header := range headers {
			record = append(record, excelCellFromValue(typed[header]))
		}
		return record, nil
	case []any:
		record := make([]excelCellValue, 0, len(typed))
		for _, item := range typed {
			record = append(record, excelCellFromValue(item))
		}
		return record, nil
	case []string:
		record := make([]excelCellValue, 0, len(typed))
		for _, item := range typed {
			record = append(record, excelStringCell(item))
		}
		return record, nil
	default:
		return []excelCellValue{excelCellFromValue(typed)}, nil
	}
}

func excelCellFromValue(value any) excelCellValue {
	switch typed := value.(type) {
	case nil:
		return excelBlankCell()
	case string:
		return excelStringCell(typed)
	case bool:
		if typed {
			return excelCellValue{Kind: excelCellBool, Value: "1"}
		}
		return excelCellValue{Kind: excelCellBool, Value: "0"}
	case int:
		return excelNumberCell(strconv.Itoa(typed))
	case int8:
		return excelNumberCell(strconv.FormatInt(int64(typed), 10))
	case int16:
		return excelNumberCell(strconv.FormatInt(int64(typed), 10))
	case int32:
		return excelNumberCell(strconv.FormatInt(int64(typed), 10))
	case int64:
		return excelNumberCell(strconv.FormatInt(typed, 10))
	case uint:
		return excelNumberCell(strconv.FormatUint(uint64(typed), 10))
	case uint8:
		return excelNumberCell(strconv.FormatUint(uint64(typed), 10))
	case uint16:
		return excelNumberCell(strconv.FormatUint(uint64(typed), 10))
	case uint32:
		return excelNumberCell(strconv.FormatUint(uint64(typed), 10))
	case uint64:
		return excelNumberCell(strconv.FormatUint(typed, 10))
	case float32:
		return excelFloatCell(float64(typed), 32)
	case float64:
		return excelFloatCell(typed, 64)
	case []any, []string, map[string]any:
		encoded, err := json.Marshal(typed)
		if err == nil {
			return excelStringCell(string(encoded))
		}
	}
	return excelStringCell(strings.TrimSpace(fmt.Sprint(value)))
}

func excelFloatCell(value float64, bitSize int) excelCellValue {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return excelStringCell(strings.TrimSpace(fmt.Sprint(value)))
	}
	return excelNumberCell(strconv.FormatFloat(value, 'f', -1, bitSize))
}

func excelBlankCell() excelCellValue {
	return excelCellValue{Kind: excelCellBlank}
}

func excelStringCell(value string) excelCellValue {
	if value == "" {
		return excelBlankCell()
	}
	return excelCellValue{Kind: excelCellString, Value: value}
}

func excelNumberCell(value string) excelCellValue {
	return excelCellValue{Kind: excelCellNumber, Value: value}
}

func normalizeExcelSheetName(sheetName string) (string, error) {
	if sheetName == "" {
		return defaultExcelSheetName, nil
	}
	trimmed := strings.TrimSpace(sheetName)
	if trimmed == "" {
		return "", fmt.Errorf("sheet cannot be blank")
	}
	if err := validateExcelSheetName(trimmed); err != nil {
		return "", err
	}
	return trimmed, nil
}

func validateExcelSheetName(sheetName string) error {
	if sheetName == "" {
		return fmt.Errorf("sheet cannot be blank")
	}
	if len([]rune(sheetName)) > 31 {
		return fmt.Errorf("sheet cannot exceed 31 characters")
	}
	if strings.ContainsAny(sheetName, `[]:*?/\`) {
		return fmt.Errorf("sheet contains invalid characters; avoid []:*?/\\")
	}
	return nil
}

func buildXLSXArchiveEntries(sheets []excelSheetWriteData) []xlsxArchiveEntry {
	entries := []xlsxArchiveEntry{
		{
			Name:    "[Content_Types].xml",
			Content: buildXLSXContentTypesXML(sheets),
		},
		{
			Name: "_rels/.rels",
			Content: `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,
		},
		{
			Name:    "xl/workbook.xml",
			Content: buildXLSXWorkbookXML(sheets),
		},
		{
			Name:    "xl/_rels/workbook.xml.rels",
			Content: buildXLSXWorkbookRelationshipsXML(sheets),
		},
	}
	for index, sheet := range sheets {
		entries = append(entries, xlsxArchiveEntry{
			Name:    fmt.Sprintf("xl/worksheets/sheet%d.xml", index+1),
			Content: buildXLSXSheetXML(sheet.Rows),
		})
	}
	return entries
}

func buildXLSXContentTypesXML(sheets []excelSheetWriteData) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` + "\n")
	builder.WriteString(`  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` + "\n")
	builder.WriteString(`  <Default Extension="xml" ContentType="application/xml"/>` + "\n")
	builder.WriteString(`  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>` + "\n")
	for index := range sheets {
		builder.WriteString(fmt.Sprintf(`  <Override PartName="/xl/worksheets/sheet%d.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`+"\n", index+1))
	}
	builder.WriteString(`</Types>`)
	return builder.String()
}

func buildXLSXWorkbookXML(sheets []excelSheetWriteData) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` + "\n")
	builder.WriteString(`  <sheets>` + "\n")
	for index, sheet := range sheets {
		builder.WriteString(fmt.Sprintf(`    <sheet name="%s" sheetId="%d" r:id="rId%d"/>`+"\n", escapeXMLAttribute(sheet.Name), index+1, index+1))
	}
	builder.WriteString(`  </sheets>` + "\n")
	builder.WriteString(`</workbook>`)
	return builder.String()
}

func buildXLSXWorkbookRelationshipsXML(sheets []excelSheetWriteData) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + "\n")
	for index := range sheets {
		builder.WriteString(fmt.Sprintf(`  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`+"\n", index+1, index+1))
	}
	builder.WriteString(`</Relationships>`)
	return builder.String()
}

func buildXLSXSheetXML(records [][]excelCellValue) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">` + "\n")
	builder.WriteString(`  <sheetData>` + "\n")
	for rowIndex, row := range records {
		rowNumber := rowIndex + 1
		builder.WriteString(fmt.Sprintf(`    <row r="%d">`, rowNumber))
		builder.WriteByte('\n')
		for columnIndex, cell := range row {
			if cell.isBlank() {
				continue
			}
			builder.WriteString(buildXLSXCellXML(xlsxCellRef(rowNumber, columnIndex+1), cell))
			builder.WriteByte('\n')
		}
		builder.WriteString(`    </row>` + "\n")
	}
	builder.WriteString(`  </sheetData>` + "\n")
	builder.WriteString(`</worksheet>`)
	return builder.String()
}

func buildXLSXCellXML(cellRef string, cell excelCellValue) string {
	switch cell.Kind {
	case excelCellNumber:
		return fmt.Sprintf(`      <c r="%s"><v>%s</v></c>`, cellRef, escapeXLSXText(cell.Value))
	case excelCellBool:
		return fmt.Sprintf(`      <c r="%s" t="b"><v>%s</v></c>`, cellRef, cell.Value)
	default:
		return fmt.Sprintf(`      <c r="%s" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, cellRef, escapeXLSXText(cell.Value))
	}
}

func (cell excelCellValue) isBlank() bool {
	return cell.Kind == excelCellBlank || (cell.Kind == excelCellString && cell.Value == "")
}

func writeXLSXArchive(filePath string, entries []xlsxArchiveEntry) error {
	if err := writeOutputFile(filePath, nil); err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	for _, entry := range entries {
		fileWriter, err := writer.Create(entry.Name)
		if err != nil {
			_ = writer.Close()
			return err
		}
		if _, err := fileWriter.Write([]byte(entry.Content)); err != nil {
			_ = writer.Close()
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func xlsxCellRef(rowNumber int, columnNumber int) string {
	return xlsxColumnName(columnNumber) + strconv.Itoa(rowNumber)
}

func xlsxColumnName(columnNumber int) string {
	if columnNumber < 1 {
		return ""
	}
	name := ""
	for columnNumber > 0 {
		columnNumber--
		name = string(rune('A'+(columnNumber%26))) + name
		columnNumber /= 26
	}
	return name
}

func escapeXLSXText(value string) string {
	var buffer bytes.Buffer
	_ = xml.EscapeText(&buffer, []byte(sanitizeXLSXText(value)))
	return buffer.String()
}

func escapeXMLAttribute(value string) string {
	return xmlAttributeEscaper.Replace(sanitizeXLSXText(value))
}

func sanitizeXLSXText(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))
	for _, r := range value {
		switch {
		case r == '\t' || r == '\n' || r == '\r':
			builder.WriteRune(r)
		case r >= 0x20 && r <= 0xD7FF:
			builder.WriteRune(r)
		case r >= 0xE000 && r <= 0xFFFD:
			builder.WriteRune(r)
		case r >= 0x10000 && r <= 0x10FFFF:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
