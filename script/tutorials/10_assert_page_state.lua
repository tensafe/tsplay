page_url = os.getenv("TSPLAY_ASSERT_URL") or "http://127.0.0.1:8000/demo/extract.html"

navigate(page_url)
wait_for_selector("#notice", 5000)

notice_visible = assert_visible("#notice", 5000)
notice_ready = assert_text("#notice", "Ready", 5000)
page_title = extract_text("#page-title", 5000)
order_count = extract_text("#summary-count", 5000, "([0-9]+)")

write_json("artifacts/tutorials/10-assert-page-state-lua.json", {
    lesson = "10",
    mode = "lua",
    page_url = page_url,
    page_title = page_title,
    order_count = order_count,
    notice_visible = notice_visible,
    notice_ready = notice_ready
})

print("page title:", tostring(page_title))
print("order count:", tostring(order_count))
print("wrote artifacts/tutorials/10-assert-page-state-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
