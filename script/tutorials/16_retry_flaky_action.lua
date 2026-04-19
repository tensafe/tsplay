page_url = os.getenv("TSPLAY_RETRY_URL") or "http://127.0.0.1:8000/demo/retry_wait_until.html"

navigate(page_url)
wait_for_selector("#retry-button", 5000)

retry_success = false
attempts = 0
status_text = ""

for attempt = 1, 3 do
    attempts = attempt
    click("#retry-button")
    status_text = extract_text("#retry-status", 5000)
    if string.find(status_text, "Export complete", 1, true) then
        retry_success = true
        break
    end
    sleep(0.2)
end

if not retry_success then
    error("expected retry demo to complete within 3 attempts")
end

badge_visible = assert_visible("#retry-result", 5000)

write_json("artifacts/tutorials/16-retry-flaky-action-lua.json", {
    lesson = "16",
    mode = "lua",
    page_url = page_url,
    attempts = attempts,
    status_text = status_text,
    badge_visible = badge_visible
})

print("attempts:", tostring(attempts))
print("wrote artifacts/tutorials/16-retry-flaky-action-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
