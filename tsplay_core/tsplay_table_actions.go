package tsplay_core

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"path/filepath"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type xlsxSheetInfo struct {
	Name string
	Path string
}

func read_csv(L *lua.LState) int {
	filePath := L.CheckString(1)
	rows, err := loadCSVRows(filePath)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, rows))
	return 1
}

func read_excel(L *lua.LState) int {
	filePath := L.CheckString(1)
	sheet := ""
	if L.GetTop() >= 2 {
		sheet = L.OptString(2, "")
	}
	rows, err := loadExcelRows(filePath, sheet)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, rows))
	return 1
}

func loadCSVRows(filePath string) ([]any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("read_csv open %q: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	if bom, err := reader.Peek(3); err == nil && len(bom) == 3 && bytes.Equal(bom, []byte{0xEF, 0xBB, 0xBF}) {
		if _, discardErr := reader.Discard(3); discardErr != nil {
			return nil, fmt.Errorf("read_csv discard utf-8 bom from %q: %w", filePath, discardErr)
		}
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read_csv parse %q: %w", filePath, err)
	}
	return tableRowsToObjects(records), nil
}

func loadExcelRows(filePath string, sheet string) ([]any, error) {
	if strings.ToLower(filepath.Ext(filePath)) != ".xlsx" {
		return nil, fmt.Errorf("read_excel currently supports only .xlsx files")
	}

	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("read_excel open %q: %w", filePath, err)
	}
	defer archive.Close()

	sharedStrings, err := readXLSXSharedStrings(&archive.Reader)
	if err != nil {
		return nil, fmt.Errorf("read_excel load shared strings from %q: %w", filePath, err)
	}

	sheets, err := readXLSXWorkbookSheets(&archive.Reader)
	if err != nil {
		return nil, fmt.Errorf("read_excel load workbook metadata from %q: %w", filePath, err)
	}
	if len(sheets) == 0 {
		return nil, fmt.Errorf("read_excel found no sheets in %q", filePath)
	}

	selected := sheets[0]
	if strings.TrimSpace(sheet) != "" {
		found := false
		for _, candidate := range sheets {
			if candidate.Name == sheet {
				selected = candidate
				found = true
				break
			}
		}
		if !found {
			names := make([]string, 0, len(sheets))
			for _, candidate := range sheets {
				names = append(names, candidate.Name)
			}
			return nil, fmt.Errorf("read_excel sheet %q not found in %q; available sheets: %s", sheet, filePath, strings.Join(names, ", "))
		}
	}

	records, err := readXLSXSheetRows(&archive.Reader, selected.Path, sharedStrings)
	if err != nil {
		return nil, fmt.Errorf("read_excel read sheet %q from %q: %w", selected.Name, filePath, err)
	}
	return tableRowsToObjects(records), nil
}

func tableRowsToObjects(records [][]string) []any {
	headerIndex := -1
	for i, record := range records {
		if rowIsEmpty(record) {
			continue
		}
		headerIndex = i
		break
	}
	if headerIndex < 0 {
		return []any{}
	}

	headers := normalizeTableHeaders(records[headerIndex])
	rows := make([]any, 0, len(records)-headerIndex-1)
	for _, record := range records[headerIndex+1:] {
		if rowIsEmpty(record) {
			continue
		}
		row := map[string]any{}
		width := len(headers)
		if len(record) > width {
			width = len(record)
		}
		for i := 0; i < width; i++ {
			header := tableHeaderName(headers, i)
			value := ""
			if i < len(record) {
				value = strings.TrimSpace(record[i])
			}
			row[header] = value
		}
		rows = append(rows, row)
	}
	return rows
}

func normalizeTableHeaders(headerRow []string) []string {
	headers := make([]string, 0, len(headerRow))
	seen := map[string]int{}
	for i, raw := range headerRow {
		name := strings.TrimSpace(raw)
		if name == "" {
			name = fmt.Sprintf("column_%d", i+1)
		}
		seen[name]++
		if seen[name] > 1 {
			name = fmt.Sprintf("%s_%d", name, seen[name])
		}
		headers = append(headers, name)
	}
	return headers
}

func tableHeaderName(headers []string, index int) string {
	if index < len(headers) && strings.TrimSpace(headers[index]) != "" {
		return headers[index]
	}
	return fmt.Sprintf("column_%d", index+1)
}

func rowIsEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func readXLSXWorkbookSheets(archive *zip.Reader) ([]xlsxSheetInfo, error) {
	workbookXML, ok, err := readZipFile(archive, "xl/workbook.xml")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("missing xl/workbook.xml")
	}

	relsXML, ok, err := readZipFile(archive, "xl/_rels/workbook.xml.rels")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("missing xl/_rels/workbook.xml.rels")
	}

	relationships, err := readXLSXRelationships(relsXML, "xl/workbook.xml")
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(workbookXML))
	sheets := []xlsxSheetInfo{}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sheet" {
			continue
		}
		name := xmlAttrValue(start, "name")
		relationshipID := xmlAttrValue(start, "id")
		target, ok := relationships[relationshipID]
		if !ok || target == "" {
			continue
		}
		sheets = append(sheets, xlsxSheetInfo{
			Name: name,
			Path: target,
		})
	}
	return sheets, nil
}

func readXLSXRelationships(content []byte, basePath string) (map[string]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	relationships := map[string]string{}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}
		id := xmlAttrValue(start, "Id")
		target := xmlAttrValue(start, "Target")
		if id == "" || target == "" {
			continue
		}
		relationships[id] = resolveXLSXPartPath(basePath, target)
	}
	return relationships, nil
}

func readXLSXSharedStrings(archive *zip.Reader) ([]string, error) {
	content, ok, err := readZipFile(archive, "xl/sharedStrings.xml")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	decoder := xml.NewDecoder(bytes.NewReader(content))
	values := []string{}
	var builder strings.Builder
	inSharedString := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			if typed.Name.Local == "si" {
				builder.Reset()
				inSharedString = true
				continue
			}
			if inSharedString && typed.Name.Local == "t" {
				var text string
				if err := decoder.DecodeElement(&text, &typed); err != nil {
					return nil, err
				}
				builder.WriteString(text)
			}
		case xml.EndElement:
			if typed.Name.Local == "si" && inSharedString {
				values = append(values, builder.String())
				inSharedString = false
			}
		}
	}
	return values, nil
}

func readXLSXSheetRows(archive *zip.Reader, partPath string, sharedStrings []string) ([][]string, error) {
	content, ok, err := readZipFile(archive, partPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("missing %s", partPath)
	}

	decoder := xml.NewDecoder(bytes.NewReader(content))
	rows := [][]string{}
	currentRow := []string(nil)
	inRow := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "row":
				currentRow = []string{}
				inRow = true
			case "c":
				if !inRow {
					continue
				}
				columnIndex, value, err := readXLSXCell(decoder, typed, sharedStrings, len(currentRow))
				if err != nil {
					return nil, err
				}
				for len(currentRow) <= columnIndex {
					currentRow = append(currentRow, "")
				}
				currentRow[columnIndex] = value
			}
		case xml.EndElement:
			if typed.Name.Local == "row" && inRow {
				rows = append(rows, trimTrailingEmptyCells(currentRow))
				currentRow = nil
				inRow = false
			}
		}
	}
	return rows, nil
}

func readXLSXCell(decoder *xml.Decoder, start xml.StartElement, sharedStrings []string, defaultColumn int) (int, string, error) {
	columnIndex := defaultColumn
	if cellRef := xmlAttrValue(start, "r"); cellRef != "" {
		index, err := xlsxColumnIndex(cellRef)
		if err != nil {
			return 0, "", err
		}
		columnIndex = index
	}

	cellType := xmlAttrValue(start, "t")
	rawValue := ""
	inlineValue := ""
	for {
		token, err := decoder.Token()
		if err != nil {
			return 0, "", err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "v":
				if err := decoder.DecodeElement(&rawValue, &typed); err != nil {
					return 0, "", err
				}
			case "t":
				var text string
				if err := decoder.DecodeElement(&text, &typed); err != nil {
					return 0, "", err
				}
				inlineValue += text
			}
		case xml.EndElement:
			if typed.Name.Local == "c" {
				return columnIndex, decodeXLSXCellValue(cellType, rawValue, inlineValue, sharedStrings), nil
			}
		}
	}
}

func decodeXLSXCellValue(cellType string, rawValue string, inlineValue string, sharedStrings []string) string {
	switch cellType {
	case "s":
		index, err := strconv.Atoi(strings.TrimSpace(rawValue))
		if err == nil && index >= 0 && index < len(sharedStrings) {
			return sharedStrings[index]
		}
		return rawValue
	case "inlineStr":
		return inlineValue
	case "b":
		if strings.TrimSpace(rawValue) == "1" {
			return "true"
		}
		return "false"
	default:
		if inlineValue != "" {
			return inlineValue
		}
		return rawValue
	}
}

func trimTrailingEmptyCells(row []string) []string {
	end := len(row)
	for end > 0 && strings.TrimSpace(row[end-1]) == "" {
		end--
	}
	if end == len(row) {
		return row
	}
	return append([]string(nil), row[:end]...)
}

func xlsxColumnIndex(cellRef string) (int, error) {
	cellRef = strings.TrimSpace(cellRef)
	if cellRef == "" {
		return 0, fmt.Errorf("xlsx cell reference cannot be blank")
	}
	letters := strings.Builder{}
	for _, r := range cellRef {
		if r >= '0' && r <= '9' {
			break
		}
		if r >= 'a' && r <= 'z' {
			r = r - 'a' + 'A'
		}
		if r < 'A' || r > 'Z' {
			return 0, fmt.Errorf("invalid xlsx cell reference %q", cellRef)
		}
		letters.WriteRune(r)
	}
	if letters.Len() == 0 {
		return 0, fmt.Errorf("invalid xlsx cell reference %q", cellRef)
	}
	column := 0
	for _, r := range letters.String() {
		column = column*26 + int(r-'A'+1)
	}
	return column - 1, nil
}

func resolveXLSXPartPath(basePath string, target string) string {
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(pathpkg.Clean(target), "/")
	}
	return pathpkg.Clean(pathpkg.Join(pathpkg.Dir(basePath), target))
}

func readZipFile(archive *zip.Reader, filePath string) ([]byte, bool, error) {
	if archive == nil {
		return nil, false, fmt.Errorf("zip archive is nil")
	}
	for _, file := range archive.File {
		if file.Name != filePath {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			return nil, false, err
		}
		defer reader.Close()
		content, err := io.ReadAll(reader)
		if err != nil {
			return nil, false, err
		}
		return content, true, nil
	}
	return nil, false, nil
}

func xmlAttrValue(start xml.StartElement, name string) string {
	for _, attr := range start.Attr {
		if strings.EqualFold(attr.Name.Local, name) {
			return attr.Value
		}
	}
	return ""
}
