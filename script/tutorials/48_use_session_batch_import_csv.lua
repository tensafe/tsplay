page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
input_file = os.getenv("TSPLAY_IMPORT_CSV") or "demo/data/import_users.csv"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_import_demo"

rows = read_csv(input_file, nil, nil, "source_row")
results = {}

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

    table.insert(results, {
        source_row = row.source_row,
        name = row.name,
        phone = row.phone,
        status = "success"
    })

    click("#clear-form")
end

import_count = get_text("#import-count")

write_json("artifacts/tutorials/48-use-session-batch-import-csv-lua.json", {
    lesson = "48",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    session_name = session_name,
    auth_status = auth_status,
    import_count = import_count,
    results = results
})

print("authenticated CSV batch import complete:", tostring(#results))
print("wrote artifacts/tutorials/48-use-session-batch-import-csv-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
