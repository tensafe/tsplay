input_file = os.getenv("TSPLAY_EXCEL_INPUT") or "demo/data/import_users.xlsx"

rows = read_excel(input_file, "RangeDemo", "A2:B4", {"name", "phone"}, nil, nil, "source_row")

write_json("artifacts/tutorials/25-read-excel-range-headers-lua.json", {
    lesson = "25",
    mode = "lua",
    input_file = input_file,
    row_count = #rows,
    rows = rows
})

print("read excel range rows:", tostring(#rows))
print("wrote artifacts/tutorials/25-read-excel-range-headers-lua.json")
