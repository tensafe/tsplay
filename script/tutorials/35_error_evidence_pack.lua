page_url = os.getenv("TSPLAY_IMPORT_URL") or "http://127.0.0.1:8000/demo/import_workflow.html"
result_path = "artifacts/tutorials/35-error-evidence-pack-lua.json"
screenshot_path = "artifacts/tutorials/35-error-evidence-pack-lua.png"
html_path = "artifacts/tutorials/35-error-evidence-pack-lua.html"

write_json(result_path, {
    lesson = "35",
    mode = "lua",
    status = "starting",
    screenshot_path = screenshot_path,
    html_path = html_path
})

navigate(page_url)
wait_for_selector("#import-form", 5000)
type_text("#name", "Fail Case")
type_text("#phone", "99900000000")
click("#submit")

evidence_status = "unexpected_success"
error_message = ""

ok, err = pcall(function()
    wait_for_selector("#submit-status", 1000)
    status_text = get_text("#submit-status")
    if string.find(status_text or "", "Imported", 1, true) == nil then
        error("expected submit status to contain 'Imported', actual=" .. tostring(status_text))
    end
end)

if not ok then
    evidence_status = "handled_error"
    error_message = tostring(err)
    screenshot(screenshot_path)
    save_html(html_path)
end

submit_status = get_text("#submit-status")

write_json(result_path, {
    lesson = "35",
    mode = "lua",
    page_url = page_url,
    evidence_status = evidence_status,
    submit_status = submit_status,
    error_message = error_message,
    screenshot_path = screenshot_path,
    html_path = html_path
})

print("captured error evidence pack")
print("wrote " .. result_path)
