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

cookie_header = get_cookies_string()
cookie_has_user = string.find(cookie_header, "tsplay_user=", 1, true) ~= nil
cookie_has_session = string.find(cookie_header, "tsplay_session=", 1, true) ~= nil

write_json("artifacts/tutorials/29-read-cookies-string-lua.json", {
    lesson = "29",
    mode = "lua",
    page_url = page_url,
    username = username,
    session_status = session_status,
    cookie_header = cookie_header,
    cookie_has_user = cookie_has_user,
    cookie_has_session = cookie_has_session
})

print("cookie header:", tostring(cookie_header))
print("wrote artifacts/tutorials/29-read-cookies-string-lua.json")
