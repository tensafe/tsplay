package tsplay_core

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestValidateFlowSecurityRejectsReadCSVWithoutAllow(t *testing.T) {
	root := t.TempDir()
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "csv_policy",
		Steps: []FlowStep{
			{Action: "read_csv", FilePath: "imports/users.csv"},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, FlowSecurityPolicy{
		FileInputRoot:  root,
		FileOutputRoot: root,
	})
	if err == nil {
		t.Fatalf("expected file access security policy error")
	}
	if !strings.Contains(err.Error(), "allow_file_access") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFlowReadCSVRows(t *testing.T) {
	root := t.TempDir()
	csvPath := filepath.Join(root, "users.csv")
	if err := os.WriteFile(csvPath, []byte("name,phone\nAlice,13800000000\nBob,13900000000\n"), 0644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "read_csv_rows",
		Steps: []FlowStep{
			{Action: "read_csv", FilePath: "users.csv", SaveAs: "rows"},
			{Action: "set_var", SaveAs: "first_name", Value: "{{rows[0].name}}"},
			{Action: "set_var", SaveAs: "summary", Value: "hello {{rows[1].name}}"},
			{
				Action:  "foreach",
				Items:   "{{rows}}",
				ItemVar: "row",
				Steps: []FlowStep{
					{Action: "set_var", SaveAs: "last_phone", Value: "{{row.phone}}"},
				},
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["first_name"]; got != "Alice" {
		t.Fatalf("first_name = %#v", got)
	}
	if got := result.Vars["summary"]; got != "hello Bob" {
		t.Fatalf("summary = %#v", got)
	}
	if got := result.Vars["last_phone"]; got != "13900000000" {
		t.Fatalf("last_phone = %#v", got)
	}
}

func TestRunFlowReadExcelRows(t *testing.T) {
	root := t.TempDir()
	xlsxPath := filepath.Join(root, "users.xlsx")
	if err := writeTestXLSXFile(xlsxPath); err != nil {
		t.Fatalf("write xlsx: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "read_excel_rows",
		Steps: []FlowStep{
			{Action: "read_excel", FilePath: "users.xlsx", Sheet: "Users", SaveAs: "rows"},
			{Action: "set_var", SaveAs: "first_user", Value: `{{rows[0]["User Name"]}}`},
			{Action: "set_var", SaveAs: "second_email", Value: "{{rows[1].Email}}"},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if got := result.Vars["first_user"]; got != "Alice Chen" {
		t.Fatalf("first_user = %#v", got)
	}
	if got := result.Vars["second_email"]; got != "bob@example.com" {
		t.Fatalf("second_email = %#v", got)
	}
}

func writeTestXLSXFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`,
		"xl/workbook.xml": `<?xml version="1.0" encoding="UTF-8"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets>
    <sheet name="Users" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`,
		"xl/worksheets/sheet1.xml": `<?xml version="1.0" encoding="UTF-8"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <sheetData>
    <row r="1">
      <c r="A1" t="inlineStr"><is><t>User Name</t></is></c>
      <c r="B1" t="inlineStr"><is><t>Email</t></is></c>
    </row>
    <row r="2">
      <c r="A2" t="inlineStr"><is><t>Alice Chen</t></is></c>
      <c r="B2" t="inlineStr"><is><t>alice@example.com</t></is></c>
    </row>
    <row r="3">
      <c r="A3" t="inlineStr"><is><t>Bob Li</t></is></c>
      <c r="B3" t="inlineStr"><is><t>bob@example.com</t></is></c>
    </row>
  </sheetData>
</worksheet>`,
	}

	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			_ = writer.Close()
			return err
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			_ = writer.Close()
			return err
		}
	}

	return writer.Close()
}
