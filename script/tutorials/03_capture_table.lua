table_url = os.getenv("TSPLAY_TABLE_URL") or "http://127.0.0.1:8000/demo/tables.html"

navigate(table_url)
wait_for_selector("#myTable", 5000)

rows = capture_table("#myTable")
write_json("artifacts/tutorials/03-capture-table-lua.json", rows)

print("captured rows:", #rows)
print("wrote artifacts/tutorials/03-capture-table-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
