page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"

rows = read_csv(input_file, nil, nil, "source_row")

use_session(session_name)
navigate(page_url)
wait_for_selector("#import-form", 5000)
auth_status = get_text("#auth-status")
if not string.find(auth_status, "Logged in as", 1, true) then
    error("expected authenticated status, got: " .. tostring(auth_status))
end

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
import_count = get_text("#import-count")

write_json("artifacts/tutorials/52-use-session-capture-import-table-lua.json", {
    lesson = "52",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    auth_status = auth_status,
    import_count = import_count,
    table_rows = table_rows
})

print("captured protected table rows:", tostring(#table_rows))
print("wrote artifacts/tutorials/52-use-session-capture-import-table-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
