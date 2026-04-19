page_url = os.getenv("TSPLAY_SESSION_LAB_URL") or "http://127.0.0.1:8000/demo/session_lab.html"
username = os.getenv("TSPLAY_SESSION_USER") or "demo-user"
state_path = os.getenv("TSPLAY_STATE_FILE") or "artifacts/tutorials/36-session-lab-lua-state.json"

navigate(page_url)
wait_for_selector("#page-title", 5000)
click("#clear-state")
type_text("#username", username)
click("#login")
wait_for_selector("#session-status", 5000)
session_status = get_text("#session-status")
if string.find(session_status or "", "Logged in as", 1, true) == nil then
    error("expected session status to contain 'Logged in as', actual=" .. tostring(session_status))
end

save_storage_state(state_path)

write_json("artifacts/tutorials/36-save-storage-state-lua.json", {
    lesson = "36",
    mode = "lua",
    page_url = page_url,
    username = username,
    session_status = session_status,
    saved_state_path = state_path
})

print("saved storage state:", state_path)
print("wrote artifacts/tutorials/36-save-storage-state-lua.json")
