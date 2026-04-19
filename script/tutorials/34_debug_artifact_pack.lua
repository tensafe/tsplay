page_url = os.getenv("TSPLAY_DEBUG_PAGE_URL") or "http://127.0.0.1:8000/demo/debug_artifacts.html"
result_path = "artifacts/tutorials/34-debug-artifact-pack-lua.json"
full_page_path = "artifacts/tutorials/34-debug-artifacts-pack-page-lua.png"
card_path = "artifacts/tutorials/34-debug-artifacts-pack-card-lua.png"
html_path = "artifacts/tutorials/34-debug-artifacts-pack-lua.html"

write_json(result_path, {
    lesson = "34",
    mode = "lua",
    status = "starting",
    full_page_path = full_page_path,
    card_path = card_path,
    html_path = html_path
})

navigate(page_url)
wait_for_selector("#artifact-card", 5000)
screenshot(full_page_path)
screenshot_element("#artifact-card", card_path)
save_html(html_path)
artifact_title = get_text("#artifact-title")
artifact_note = get_text("#artifact-note")

write_json(result_path, {
    lesson = "34",
    mode = "lua",
    page_url = page_url,
    artifact_title = artifact_title,
    artifact_note = artifact_note,
    full_page_path = full_page_path,
    card_path = card_path,
    html_path = html_path
})

print("saved debug artifact pack")
print("wrote " .. result_path)
