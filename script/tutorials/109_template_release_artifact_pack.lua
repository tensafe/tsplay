page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"
result_path = "artifacts/tutorials/109-template-release-artifact-pack-lua.json"
full_page_path = "artifacts/tutorials/109-template-release-full-page-lua.png"
card_path = "artifacts/tutorials/109-template-release-card-lua.png"
html_path = "artifacts/tutorials/109-template-release-artifact-pack-lua.html"

navigate(page_url)
wait_for_selector("#template-release-card", 5000)
click("#stage-check-button")

for attempt = 1, 25 do
    if is_visible("#stage-ready-badge") then
        break
    end
    sleep(0.2)
end

screenshot(full_page_path)
screenshot_element("#template-release-card", card_path)
save_html(html_path)

stage_status = get_text("#stage-check-status")
release_summary = get_text("#release-summary")

write_json(result_path, {
    lesson = "109",
    mode = "lua",
    page_url = page_url,
    stage_status = stage_status,
    release_summary = release_summary,
    full_page_path = full_page_path,
    card_path = card_path,
    html_path = html_path
})

print("saved artifact pack:", result_path)
