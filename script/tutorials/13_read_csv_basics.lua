input_file = os.getenv("TSPLAY_CSV_INPUT") or "demo/data/tutorial_contacts.csv"

rows = read_csv(input_file, nil, nil, "source_row")

write_json("artifacts/tutorials/13-read-csv-basics-lua.json", {
    lesson = "13",
    mode = "lua",
    input_file = input_file,
    row_count = #rows,
    rows = rows
})

print("read rows:", tostring(#rows))
print("wrote artifacts/tutorials/13-read-csv-basics-lua.json")
