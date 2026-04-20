page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"
invalid_ticket = os.getenv("TSPLAY_INVALID_RELEASE_TICKET") or "BAD-TICKET"
result_path = "artifacts/tutorials/110-template-release-robustness-round-trip-lua.json"
screenshot_path = "artifacts/tutorials/110-template-release-robustness-round-trip-lua.png"
html_path = "artifacts/tutorials/110-template-release-robustness-round-trip-lua.html"

navigate(page_url)
wait_for_selector("#template-release-card", 5000)

card_visible = is_visible("#template-release-card")
if not card_visible then
    error("expected #template-release-card to be visible")
end
base_status = get_text("#release-status")
base_status_ok = base_status == "Draft checks complete"
if not base_status_ok then
    error("expected #release-status to equal 'Draft checks complete', actual=" .. tostring(base_status))
end

click("#stage-check-button")
stage_polls = 0
for attempt = 1, 25 do
    stage_polls = attempt
    if is_visible("#stage-ready-badge") then
        break
    end
    sleep(0.2)
end
stage_status = get_text("#stage-check-status")

gate_attempts = 0
gate_status = ""
for attempt = 1, 3 do
    gate_attempts = attempt
    click("#retry-gate-button")
    gate_status = get_text("#retry-gate-status")
    if string.find(gate_status or "", "Gate passed", 1, true) ~= nil then
        break
    end
    sleep(0.2)
end
gate_badge_visible = is_visible("#retry-gate-badge")
if not gate_badge_visible then
    error("expected #retry-gate-badge to be visible")
end

click("#show-delay-note-button")
wait_for_selector("#delayed-release-note", 5000)
delayed_note = get_text("#delayed-release-note")

publish_attempts = 0
publish_status = ""
for attempt = 1, 3 do
    publish_attempts = attempt
    click("#flaky-publish-button")
    publish_status = get_text("#flaky-publish-status")
    if string.find(publish_status or "", "Publish preview opened", 1, true) ~= nil then
        break
    end
    sleep(0.2)
end
publish_preview_visible = is_visible("#flaky-publish-badge")
if not publish_preview_visible then
    error("expected #flaky-publish-badge to be visible")
end

type_text("#release-ticket", invalid_ticket)
click("#finalize-release-button")

error_branch = "unexpected_success"
failure_status = ""
error_message = ""

ok, err = pcall(function()
    status_text = get_text("#finalize-status")
    if string.find(status_text or "", "Release finalized", 1, true) == nil then
        error("expected finalize status to contain 'Release finalized', actual=" .. tostring(status_text))
    end
end)

if not ok then
    error_branch = "handled_error"
    failure_status = get_text("#finalize-status")
    error_message = tostring(err)
    click("#reset-finalize")
end

reset_status = get_text("#finalize-status")

click("#reload-recovery-button")
reload_attempts = 0
reload_status = ""
for attempt = 1, 3 do
    reload_attempts = attempt
    reload()
    wait_for_selector("#reload-status", 5000)
    reload_status = get_text("#reload-status")
    if string.find(reload_status or "", "Recovered after reload", 1, true) ~= nil and is_visible("#reload-badge") then
        break
    end
    sleep(0.2)
end
reload_recovered = is_visible("#reload-badge")
if not reload_recovered then
    error("expected #reload-badge to be visible after reload recovery")
end

screenshot(screenshot_path)
save_html(html_path)

write_json(result_path, {
    lesson = "110",
    mode = "lua",
    page_url = page_url,
    invalid_ticket = invalid_ticket,
    card_visible = card_visible,
    base_status_ok = base_status_ok,
    stage_polls = stage_polls,
    stage_status = stage_status,
    gate_attempts = gate_attempts,
    gate_status = gate_status,
    gate_badge_visible = gate_badge_visible,
    delayed_note = delayed_note,
    publish_attempts = publish_attempts,
    publish_status = publish_status,
    publish_preview_visible = publish_preview_visible,
    error_branch = error_branch,
    failure_status = failure_status,
    reset_status = reset_status,
    error_message = error_message,
    reload_attempts = reload_attempts,
    reload_status = reload_status,
    reload_recovered = reload_recovered,
    screenshot_path = screenshot_path,
    html_path = html_path
})

print("template release robustness round trip complete")
print("wrote " .. result_path)
