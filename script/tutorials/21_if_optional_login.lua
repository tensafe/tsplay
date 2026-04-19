page_url = os.getenv("TSPLAY_IMPORT_LOGIN_URL") or "http://127.0.0.1:8000/demo/import_workflow.html?login=1"

navigate(page_url)
wait_for_selector("#page-title", 5000)

login_path = "direct_entry"
if is_visible(".login-dialog") then
    type_text("#username", "demo-user")
    click("#login")
    login_path = "login_required"
else
    wait_for_selector("#import-form", 5000)
end

wait_for_selector("#import-form", 5000)
submit_status = extract_text("#submit-status", 5000)
workflow_mode = extract_text("#workflow-mode", 5000)

write_json("artifacts/tutorials/21-if-optional-login-lua.json", {
    lesson = "21",
    mode = "lua",
    page_url = page_url,
    login_path = login_path,
    workflow_mode = workflow_mode,
    submit_status = submit_status
})

print("login path:", tostring(login_path))
print("wrote artifacts/tutorials/21-if-optional-login-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
