page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#flaky-publish-button", 5000)

attempts = 0
preview_opened = false
publish_status = ""

for attempt = 1, 3 do
    attempts = attempt
    click("#flaky-publish-button")
    publish_status = get_text("#flaky-publish-status")
    if string.find(publish_status or "", "Publish preview opened", 1, true) ~= nil then
        preview_opened = true
        break
    end
    sleep(0.2)
end

if not preview_opened then
    error("expected publish preview to open within 3 attempts")
end

preview_visible = is_visible("#flaky-publish-badge")
if not preview_visible then
    error("expected #flaky-publish-badge to be visible")
end

write_json("artifacts/tutorials/107-retry-flaky-publish-click-lua.json", {
    lesson = "107",
    mode = "lua",
    page_url = page_url,
    attempts = attempts,
    publish_status = publish_status,
    preview_opened = preview_opened,
    preview_visible = preview_visible
})

print("publish attempts:", tostring(attempts))
print("wrote artifacts/tutorials/107-retry-flaky-publish-click-lua.json")
