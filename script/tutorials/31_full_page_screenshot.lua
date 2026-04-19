page_url = os.getenv("TSPLAY_DEBUG_PAGE_URL") or "http://127.0.0.1:8000/demo/debug_artifacts.html"
result_path = "artifacts/tutorials/31-full-page-screenshot-lua.json"
screenshot_path = "artifacts/tutorials/31-debug-artifacts-full-page-lua.png"

write_json(result_path, {
    lesson = "31",
    mode = "lua",
    status = "starting",
    screenshot_path = screenshot_path
})

navigate(page_url)
wait_for_selector("#artifact-card", 5000)
screenshot(screenshot_path)
artifact_title = get_text("#artifact-title")

write_json(result_path, {
    lesson = "31",
    mode = "lua",
    page_url = page_url,
    artifact_title = artifact_title,
    screenshot_path = screenshot_path
})

print("saved screenshot:", screenshot_path)
print("wrote " .. result_path)
