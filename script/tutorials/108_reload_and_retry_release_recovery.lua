page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#reload-recovery-button", 5000)
click("#reload-recovery-button")

status_before_reload = get_text("#reload-status")
reload_attempts = 0
recovered = false
status_after_reload = ""

for attempt = 1, 3 do
    reload_attempts = attempt
    reload()
    wait_for_selector("#reload-status", 5000)
    status_after_reload = get_text("#reload-status")
    if string.find(status_after_reload or "", "Recovered after reload", 1, true) ~= nil and is_visible("#reload-badge") then
        recovered = true
        break
    end
    sleep(0.2)
end

if not recovered then
    error("expected reload recovery to finish within 3 attempts")
end

write_json("artifacts/tutorials/108-reload-and-retry-release-recovery-lua.json", {
    lesson = "108",
    mode = "lua",
    page_url = page_url,
    status_before_reload = status_before_reload,
    status_after_reload = status_after_reload,
    reload_attempts = reload_attempts,
    recovered = recovered
})

print("reload recovery:", tostring(status_after_reload))
print("wrote artifacts/tutorials/108-reload-and-retry-release-recovery-lua.json")
