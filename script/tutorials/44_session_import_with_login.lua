page_url = os.getenv("TSPLAY_SESSION_IMPORT_URL") or "http://127.0.0.1:8000/demo/session_import_workflow.html"
username = os.getenv("TSPLAY_SESSION_USER") or "demo-user"

navigate(page_url)
wait_for_selector("#page-title", 5000)
click("#clear-session")
type_text("#username", username)
click("#login")
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
import_count = get_text("#import-count")

write_json("artifacts/tutorials/44-session-import-with-login-lua.json", {
    lesson = "44",
    mode = "lua",
    page_url = page_url,
    username = username,
    auth_status = auth_status,
    submit_status = submit_status,
    import_count = import_count
})

print("authenticated import complete")
print("wrote artifacts/tutorials/44-session-import-with-login-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
