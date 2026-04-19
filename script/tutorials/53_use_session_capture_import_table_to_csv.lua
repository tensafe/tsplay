page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"
output_file = os.getenv("TSPLAY_CAPTURED_TABLE_CSV") or "artifacts/tutorials/53-use-session-captured-import-table-lua.csv"

rows = read_csv(input_file, nil, nil, "source_row")

use_session(session_name)
navigate(page_url)
wait_for_selector("#import-form", 5000)

for _, row in ipairs(rows) do
    type_text("#name", row.name)
    type_text("#phone", row.phone)
    click("#submit")
    wait_for_selector("#submit-status", 5000)
    submit_status = get_text("#submit-status")
    if not string.find(submit_status, "Imported", 1, true) then
        error("expected imported status, got: " .. tostring(submit_status))
    end
    click("#clear-form")
end

table_rows = capture_table("#import-results")
write_csv(output_file, table_rows, {"Name", "Phone", "Status", "Operator"})

write_json("artifacts/tutorials/53-use-session-capture-import-table-to-csv-lua.json", {
    lesson = "53",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    output_file = output_file,
    table_rows = table_rows
})

print("captured table written to:", output_file)
print("wrote artifacts/tutorials/53-use-session-capture-import-table-to-csv-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
