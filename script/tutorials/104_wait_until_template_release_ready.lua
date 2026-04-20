page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#stage-check-button", 5000)
click("#stage-check-button")

poll_attempts = 0
ready_visible = false

for attempt = 1, 25 do
    poll_attempts = attempt
    if is_visible("#stage-ready-badge") then
        ready_visible = true
        break
    end
    sleep(0.2)
end

if not ready_visible then
    error("expected stage check to become ready within timeout")
end

stage_status = get_text("#stage-check-status")

write_json("artifacts/tutorials/104-wait-until-template-release-ready-lua.json", {
    lesson = "104",
    mode = "lua",
    page_url = page_url,
    poll_attempts = poll_attempts,
    ready_visible = ready_visible,
    stage_status = stage_status
})

print("stage status:", tostring(stage_status))
print("wrote artifacts/tutorials/104-wait-until-template-release-ready-lua.json")
