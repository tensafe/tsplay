page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#retry-gate-button", 5000)

attempts = 0
gate_passed = false
gate_status = ""

for attempt = 1, 3 do
    attempts = attempt
    click("#retry-gate-button")
    gate_status = get_text("#retry-gate-status")
    if string.find(gate_status or "", "Gate passed", 1, true) ~= nil then
        gate_passed = true
        break
    end
    sleep(0.2)
end

if not gate_passed then
    error("expected retry gate to pass within 3 attempts")
end

badge_visible = is_visible("#retry-gate-badge")
if not badge_visible then
    error("expected #retry-gate-badge to be visible")
end

write_json("artifacts/tutorials/103-retry-template-release-gate-lua.json", {
    lesson = "103",
    mode = "lua",
    page_url = page_url,
    attempts = attempts,
    gate_status = gate_status,
    gate_passed = gate_passed,
    badge_visible = badge_visible
})

print("retry gate attempts:", tostring(attempts))
print("wrote artifacts/tutorials/103-retry-template-release-gate-lua.json")
