page_url = os.getenv("TSPLAY_EXTRACT_URL") or "http://127.0.0.1:8000/demo/extract.html"

navigate(page_url)
wait_for_selector("#summary-count", 5000)

page_title = extract_text("#page-title", 5000)
order_count = extract_text("#summary-count", 5000, "([0-9]+)")
notice_html = get_html("#notice")

write_json("artifacts/tutorials/04-extract-text-and-html-lua.json", {
    lesson = "04",
    mode = "lua",
    page_url = page_url,
    page_title = page_title,
    order_count = order_count,
    notice_html = notice_html
})

print("page title:", tostring(page_title))
print("order count:", tostring(order_count))
print("wrote artifacts/tutorials/04-extract-text-and-html-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
