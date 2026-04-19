page_url = os.getenv("TSPLAY_DEBUG_PAGE_URL") or "http://127.0.0.1:8000/demo/debug_artifacts.html"
result_path = "artifacts/tutorials/32-element-screenshot-lua.json"
screenshot_path = "artifacts/tutorials/32-debug-artifacts-card-lua.png"

write_json(result_path, {
    lesson = "32",
    mode = "lua",
    status = "starting",
    screenshot_path = screenshot_path
})

navigate(page_url)
wait_for_selector("#artifact-card", 5000)
screenshot_element("#artifact-card", screenshot_path)
artifact_summary = get_text("#artifact-summary")

write_json(result_path, {
    lesson = "32",
    mode = "lua",
    page_url = page_url,
    selector = "#artifact-card",
    artifact_summary = artifact_summary,
    screenshot_path = screenshot_path
})

print("saved element screenshot:", screenshot_path)
print("wrote " .. result_path)
