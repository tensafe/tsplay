page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"
save_path = os.getenv("TSPLAY_DOWNLOAD_SAVE_PATH") or "artifacts/tutorials/56-use-session-import-report-lua.csv"

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

page_rows = capture_table("#import-results")
download_file("#download-import-report", save_path)
downloaded_rows = read_csv(save_path)

if #page_rows ~= #downloaded_rows then
    error("page rows and downloaded rows count mismatch: " .. tostring(#page_rows) .. " vs " .. tostring(#downloaded_rows))
end

write_json("artifacts/tutorials/56-use-session-compare-table-and-download-lua.json", {
    lesson = "56",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    save_path = save_path,
    page_row_count = #page_rows,
    downloaded_row_count = #downloaded_rows,
    page_rows = page_rows,
    downloaded_rows = downloaded_rows
})

print("page rows match downloaded rows:", tostring(#page_rows))
print("wrote artifacts/tutorials/56-use-session-compare-table-and-download-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
