page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#release-status", 5000)

release_status = get_text("#release-status")
release_badge = get_text("#release-badge")
status_ready = release_status == "Draft checks complete"
badge_ready = release_badge == "Ready for staged checks"
if not status_ready then
    error("expected release status to equal 'Draft checks complete', actual=" .. tostring(release_status))
end
if not badge_ready then
    error("expected release badge to equal 'Ready for staged checks', actual=" .. tostring(release_badge))
end
release_summary = get_text("#release-summary")
release_title = get_text("#release-title")

write_json("artifacts/tutorials/102-assert-text-template-release-status-lua.json", {
    lesson = "102",
    mode = "lua",
    page_url = page_url,
    status_ready = status_ready,
    badge_ready = badge_ready,
    release_title = release_title,
    release_summary = release_summary
})

print("release summary:", tostring(release_summary))
print("wrote artifacts/tutorials/102-assert-text-template-release-status-lua.json")
