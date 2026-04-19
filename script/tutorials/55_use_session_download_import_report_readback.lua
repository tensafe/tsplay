page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"
save_path = os.getenv("TSPLAY_DOWNLOAD_SAVE_PATH") or "artifacts/tutorials/55-use-session-import-report-lua.csv"

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

download_file("#download-import-report", save_path)
downloaded_rows = read_csv(save_path)
import_count = get_text("#import-count")
export_status = get_text("#export-status")

write_json("artifacts/tutorials/55-use-session-download-import-report-readback-lua.json", {
    lesson = "55",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    save_path = save_path,
    import_count = import_count,
    export_status = export_status,
    downloaded_row_count = #downloaded_rows,
    downloaded_rows = downloaded_rows
})

print("downloaded rows:", tostring(#downloaded_rows))
print("wrote artifacts/tutorials/55-use-session-download-import-report-readback-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
