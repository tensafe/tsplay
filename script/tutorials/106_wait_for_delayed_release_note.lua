page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"

navigate(page_url)
wait_for_selector("#show-delay-note-button", 5000)
click("#show-delay-note-button")
wait_for_selector("#delayed-release-note", 5000)

note_visible = is_visible("#delayed-release-note")
if not note_visible then
    error("expected #delayed-release-note to be visible")
end
note_text = get_text("#delayed-release-note")
delay_status = get_text("#delay-note-status")

write_json("artifacts/tutorials/106-wait-for-delayed-release-note-lua.json", {
    lesson = "106",
    mode = "lua",
    page_url = page_url,
    note_visible = note_visible,
    note_text = note_text,
    delay_status = delay_status
})

print("delayed note:", tostring(note_text))
print("wrote artifacts/tutorials/106-wait-for-delayed-release-note-lua.json")
