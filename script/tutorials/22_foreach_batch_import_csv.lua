page_url = os.getenv("TSPLAY_IMPORT_URL") or "http://127.0.0.1:8000/demo/import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"

rows = read_csv(input_file, nil, nil, "source_row")
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

write_json("artifacts/tutorials/22-foreach-batch-import-csv-lua.json", {
    lesson = "22",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    results = results
})

print("imported rows:", tostring(#results))
print("wrote artifacts/tutorials/22-foreach-batch-import-csv-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
