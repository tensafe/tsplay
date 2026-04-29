package tsplay_core

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
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
	default:
		headers, err := stringListValue(typed)
		if err != nil {
			return fmt.Errorf("must be a sheet name string or headers list: %w", err)
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

	sheetName, err := normalizeExcelSheetName(options.Sheet)
	if err != nil {
		return nil, fmt.Errorf("write_excel options: %w", err)
	}

	rows, resolvedHeaders, err := csvRowsFromValue(value, options.Headers)
	if err != nil {
		return nil, fmt.Errorf("write_excel normalize %q: %w", filePath, err)
	}

	entries := buildXLSXArchiveEntries(sheetName, rows, resolvedHeaders)
	if err := writeXLSXArchive(filePath, entries); err != nil {
		return nil, fmt.Errorf("write_excel write %q: %w", filePath, err)
	}

	columns := 0
	if len(resolvedHeaders) > 0 {
		columns = len(resolvedHeaders)
	} else {
		columns = maxCSVRecordWidth(rows)
	}

	result := map[string]any{
		"file_path": filePath,
		"rows":      len(rows),
		"columns":   columns,
		"sheet":     sheetName,
	}
	if info, statErr := os.Stat(filePath); statErr == nil {
		result["bytes"] = info.Size()
	}
	return result, nil
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

func buildXLSXArchiveEntries(sheetName string, rows [][]string, headers []string) []xlsxArchiveEntry {
	return []xlsxArchiveEntry{
		{
			Name: "[Content_Types].xml",
			Content: `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`,
		},
		{
			Name: "_rels/.rels",
			Content: `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,
		},
		{
			Name: "xl/workbook.xml",
			Content: fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets>
    <sheet name="%s" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`, escapeXMLAttribute(sheetName)),
		},
		{
			Name: "xl/_rels/workbook.xml.rels",
			Content: `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`,
		},
		{
			Name:    "xl/worksheets/sheet1.xml",
			Content: buildXLSXSheetXML(buildXLSXSheetRows(rows, headers)),
		},
	}
}

func buildXLSXSheetRows(rows [][]string, headers []string) [][]string {
	records := make([][]string, 0, len(rows)+1)
	if len(headers) > 0 {
		records = append(records, append([]string(nil), headers...))
	}
	for _, row := range rows {
		records = append(records, append([]string(nil), row...))
	}
	return records
}

func buildXLSXSheetXML(records [][]string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">` + "\n")
	builder.WriteString(`  <sheetData>` + "\n")
	for rowIndex, row := range records {
		rowNumber := rowIndex + 1
		builder.WriteString(fmt.Sprintf(`    <row r="%d">`, rowNumber))
		builder.WriteByte('\n')
		for columnIndex, cell := range row {
			if cell == "" {
				continue
			}
			builder.WriteString(fmt.Sprintf(
				`      <c r="%s" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`+"\n",
				xlsxCellRef(rowNumber, columnIndex+1),
				escapeXLSXText(cell),
			))
		}
		builder.WriteString(`    </row>` + "\n")
	}
	builder.WriteString(`  </sheetData>` + "\n")
	builder.WriteString(`</worksheet>`)
	return builder.String()
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

func maxCSVRecordWidth(rows [][]string) int {
	width := 0
	for _, row := range rows {
		if len(row) > width {
			width = len(row)
		}
	}
	return width
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
