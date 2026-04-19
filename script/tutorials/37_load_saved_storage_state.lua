page_url = os.getenv("TSPLAY_SESSION_LAB_URL") or "http://127.0.0.1:8000/demo/session_lab.html"
state_path = os.getenv("TSPLAY_STATE_FILE") or "artifacts/tutorials/36-session-lab-lua-state.json"

load_storage_state(state_path)
navigate(page_url)
wait_for_selector("#session-status", 5000)
session_status = get_text("#session-status")
cookie_header = get_cookies_string()

write_json("artifacts/tutorials/37-load-saved-storage-state-lua.json", {
    lesson = "37",
    mode = "lua",
    page_url = page_url,
    loaded_state_path = state_path,
    session_status = session_status,
    cookie_header = cookie_header
})

print("loaded storage state:", state_path)
print("wrote artifacts/tutorials/37-load-saved-storage-state-lua.json")
