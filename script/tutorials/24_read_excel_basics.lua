input_file = os.getenv("TSPLAY_EXCEL_INPUT") or "demo/data/import_users.xlsx"

rows = read_excel(input_file, "Users", nil, nil, nil, nil, "source_row")

write_json("artifacts/tutorials/24-read-excel-basics-lua.json", {
    lesson = "24",
    mode = "lua",
    input_file = input_file,
    row_count = #rows,
    rows = rows
})

print("read excel rows:", tostring(#rows))
print("wrote artifacts/tutorials/24-read-excel-basics-lua.json")
