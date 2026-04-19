page_url = os.getenv("TSPLAY_SESSION_LAB_URL") or "http://127.0.0.1:8000/demo/session_lab.html"
session_name = os.getenv("TSPLAY_SAVED_SESSION") or "session_lab_demo"

use_session(session_name)
navigate(page_url)
wait_for_selector("#session-status", 5000)
session_status = get_text("#session-status")
cookie_header = get_cookies_string()

write_json("artifacts/tutorials/42-use-named-session-lua.json", {
    lesson = "42",
    mode = "lua",
    page_url = page_url,
    session_name = session_name,
    session_status = session_status,
    cookie_header = cookie_header
})

print("reused named session:", session_name)
print("wrote artifacts/tutorials/42-use-named-session-lua.json")
