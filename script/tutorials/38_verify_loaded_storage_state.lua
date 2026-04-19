page_url = os.getenv("TSPLAY_SESSION_LAB_URL") or "http://127.0.0.1:8000/demo/session_lab.html"
state_path = os.getenv("TSPLAY_STATE_FILE") or "artifacts/tutorials/36-session-lab-lua-state.json"

load_storage_state(state_path)
navigate(page_url)
wait_for_selector("#session-status", 5000)
session_status = get_text("#session-status")
storage_state_json = get_storage_state()
storage_origin = json_extract(storage_state_json, "$.origins[0].origin")
cookie_header = get_cookies_string()

write_json("artifacts/tutorials/38-verify-loaded-storage-state-lua.json", {
    lesson = "38",
    mode = "lua",
    page_url = page_url,
    loaded_state_path = state_path,
    session_status = session_status,
    storage_origin = storage_origin,
    cookie_header = cookie_header,
    storage_state_json = storage_state_json
})

print("verified loaded storage state:", state_path)
print("wrote artifacts/tutorials/38-verify-loaded-storage-state-lua.json")
