page_url = os.getenv("TSPLAY_DOWNLOAD_URL") or "http://127.0.0.1:8000/demo/download.html"
save_path = os.getenv("TSPLAY_DOWNLOAD_SAVE_PATH") or "artifacts/tutorials/20-downloaded-monthly-report-lua.csv"

navigate(page_url)
wait_for_selector("#download-report", 5000)
download_file("#download-report", save_path)

rows = read_csv(save_path)

write_json("artifacts/tutorials/20-download-report-lua.json", {
    lesson = "20",
    mode = "lua",
    page_url = page_url,
    save_path = save_path,
    row_count = #rows,
    rows = rows
})

print("downloaded rows:", tostring(#rows))
print("wrote artifacts/tutorials/20-download-report-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
