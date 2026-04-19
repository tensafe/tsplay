page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"
save_path = os.getenv("TSPLAY_DOWNLOAD_SAVE_PATH") or "artifacts/tutorials/54-use-session-import-report-lua.csv"

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

wait_for_selector("#download-import-report", 5000)
download_file("#download-import-report", save_path)
export_status = get_text("#export-status")
import_count = get_text("#import-count")

write_json("artifacts/tutorials/54-use-session-download-import-report-lua.json", {
    lesson = "54",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    save_path = save_path,
    export_status = export_status,
    import_count = import_count
})

print("downloaded authenticated report to:", save_path)
print("wrote artifacts/tutorials/54-use-session-download-import-report-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
