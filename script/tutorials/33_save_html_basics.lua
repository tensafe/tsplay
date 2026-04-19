page_url = os.getenv("TSPLAY_DEBUG_PAGE_URL") or "http://127.0.0.1:8000/demo/debug_artifacts.html"
result_path = "artifacts/tutorials/33-save-html-basics-lua.json"
html_path = "artifacts/tutorials/33-debug-artifacts-page-lua.html"

write_json(result_path, {
    lesson = "33",
    mode = "lua",
    status = "starting",
    html_path = html_path
})

navigate(page_url)
wait_for_selector("#artifact-note", 5000)
save_html(html_path)
artifact_note = get_text("#artifact-note")

write_json(result_path, {
    lesson = "33",
    mode = "lua",
    page_url = page_url,
    artifact_note = artifact_note,
    html_path = html_path
})

print("saved html:", html_path)
print("wrote " .. result_path)
