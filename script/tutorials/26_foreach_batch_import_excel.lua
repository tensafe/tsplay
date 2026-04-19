page_url = os.getenv("TSPLAY_IMPORT_URL") or "http://127.0.0.1:8000/demo/import_workflow.html"
input_file = os.getenv("TSPLAY_EXCEL_INPUT") or "demo/data/import_users.xlsx"

rows = read_excel(input_file, "Users", nil, nil, nil, nil, "source_row")
results = {}

navigate(page_url)
wait_for_selector("#import-form", 5000)

for _, row in ipairs(rows) do
    type_text("#name", row.name)
    type_text("#phone", row.phone)
    click("#submit")
    assert_text("#submit-status", "Imported", 5000)

    table.insert(results, {
        source_row = row.source_row,
        name = row.name,
        phone = row.phone,
        status = "success"
    })

    click("#clear-form")
end

write_json("artifacts/tutorials/26-foreach-batch-import-excel-lua.json", {
    lesson = "26",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    results = results
})

print("imported excel rows:", tostring(#results))
print("wrote artifacts/tutorials/26-foreach-batch-import-excel-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
