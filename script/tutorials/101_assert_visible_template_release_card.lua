page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#template-release-card", 5000)

card_visible = is_visible("#template-release-card")
badge_visible = is_visible("#release-badge")
if not card_visible then
    error("expected #template-release-card to be visible")
end
if not badge_visible then
    error("expected #release-badge to be visible")
end
release_title = get_text("#release-title")
release_status = get_text("#release-status")

write_json("artifacts/tutorials/101-assert-visible-template-release-card-lua.json", {
    lesson = "101",
    mode = "lua",
    page_url = page_url,
    card_visible = card_visible,
    badge_visible = badge_visible,
    release_title = release_title,
    release_status = release_status
})

print("release title:", tostring(release_title))
print("wrote artifacts/tutorials/101-assert-visible-template-release-card-lua.json")
