page_url = os.getenv("TSPLAY_SESSION_LAB_URL") or "http://127.0.0.1:8000/demo/session_lab.html"
username = os.getenv("TSPLAY_SESSION_USER") or "demo-user"

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

storage_state_json = get_storage_state()
storage_origin = json_extract(storage_state_json, "$.origins[0].origin")

write_json("artifacts/tutorials/28-inspect-storage-state-lua.json", {
    lesson = "28",
    mode = "lua",
    page_url = page_url,
    username = username,
    session_status = session_status,
    storage_origin = storage_origin,
    storage_state_json = storage_state_json
})

print("storage origin:", tostring(storage_origin))
print("wrote artifacts/tutorials/28-inspect-storage-state-lua.json")
