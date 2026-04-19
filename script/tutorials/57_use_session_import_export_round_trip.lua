page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_EXCEL_INPUT") or "demo/data/import_users.xlsx"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"
save_path = os.getenv("TSPLAY_DOWNLOAD_SAVE_PATH") or "artifacts/tutorials/57-use-session-import-export-round-trip-lua.csv"

rows = read_excel(input_file, "Users", nil, nil, nil, nil, "source_row")

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

page_rows = capture_table("#import-results")
download_file("#download-import-report", save_path)
downloaded_rows = read_csv(save_path)
import_count = get_text("#import-count")
export_status = get_text("#export-status")

if #page_rows ~= #downloaded_rows then
    error("round trip mismatch: " .. tostring(#page_rows) .. " page rows vs " .. tostring(#downloaded_rows) .. " downloaded rows")
end

write_json("artifacts/tutorials/57-use-session-import-export-round-trip-lua.json", {
    lesson = "57",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    save_path = save_path,
    auth_status = auth_status,
    import_count = import_count,
    export_status = export_status,
    page_row_count = #page_rows,
    downloaded_row_count = #downloaded_rows,
    page_rows = page_rows,
    downloaded_rows = downloaded_rows
})

print("authenticated round trip complete:", tostring(#downloaded_rows))
print("wrote artifacts/tutorials/57-use-session-import-export-round-trip-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
