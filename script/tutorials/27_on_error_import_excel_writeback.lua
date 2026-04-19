page_url = os.getenv("TSPLAY_IMPORT_URL") or "http://127.0.0.1:8000/demo/import_workflow.html"
input_file = os.getenv("TSPLAY_EXCEL_INPUT") or "demo/data/import_users.xlsx"

rows = read_excel(input_file, "ErrorBatch", nil, nil, nil, nil, "source_row")
results = {}

navigate(page_url)
wait_for_selector("#import-form", 5000)

for _, row in ipairs(rows) do
    ok, err = pcall(function()
        type_text("#name", row.name)
        type_text("#phone", row.phone)
        click("#submit")
        assert_text("#submit-status", "Imported", 1000)

        table.insert(results, {
            source_row = row.source_row,
            name = row.name,
            phone = row.phone,
            status = "success",
            error = ""
        })

        click("#clear-form")
    end)

    if not ok then
        table.insert(results, {
            source_row = row.source_row,
            name = row.name,
            phone = row.phone,
            status = "failed",
            error = tostring(err)
        })
        click("#clear-form")
    end
end

write_json("artifacts/tutorials/27-on-error-import-excel-writeback-lua.json", {
    lesson = "27",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    results = results
})

write_csv("artifacts/tutorials/27-on-error-import-excel-writeback-lua.csv", results, {"source_row", "name", "phone", "status", "error"})

print("processed excel rows:", tostring(#results))
print("wrote artifacts/tutorials/27-on-error-import-excel-writeback-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
