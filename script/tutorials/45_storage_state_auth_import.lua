page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
state_path = os.getenv("TSPLAY_STATE_FILE") or "artifacts/tutorials/36-session-lab-lua-state.json"

load_storage_state(state_path)
navigate(page_url)
wait_for_selector("#auth-status", 5000)
auth_status = get_text("#auth-status")
if not string.find(auth_status, "Logged in as", 1, true) then
    error("expected authenticated status, got: " .. tostring(auth_status))
end
type_text("#name", "Alice")
type_text("#phone", "13800000001")
click("#submit")
wait_for_selector("#submit-status", 5000)
submit_status = get_text("#submit-status")
if not string.find(submit_status, "Imported", 1, true) then
    error("expected imported status, got: " .. tostring(submit_status))
end

session_preview = get_text("#session-preview")
import_count = get_text("#import-count")

write_json("artifacts/tutorials/45-storage-state-auth-import-lua.json", {
    lesson = "45",
    mode = "lua",
    page_url = page_url,
    loaded_state_path = state_path,
    auth_status = auth_status,
    submit_status = submit_status,
    session_preview = session_preview,
    import_count = import_count
})

print("authenticated import via storage_state complete")
print("wrote artifacts/tutorials/45-storage-state-auth-import-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
