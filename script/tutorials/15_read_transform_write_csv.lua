input_file = os.getenv("TSPLAY_CSV_INPUT") or "demo/data/tutorial_contacts.csv"
output_file = "artifacts/tutorials/15-transformed-contacts-lua.csv"

rows = read_csv(input_file, 3, 2, "source_row")
transformed_rows = {}

for _, row in ipairs(rows) do
    table.insert(transformed_rows, {
        source_row = row.source_row,
        name = row.name,
        city = row.city,
        status = row.status,
        review_note = "needs_follow_up"
    })
end

write_result = write_csv(output_file, transformed_rows, {"source_row", "name", "city", "status", "review_note"})

write_json("artifacts/tutorials/15-read-transform-write-csv-lua.json", {
    lesson = "15",
    mode = "lua",
    input_file = input_file,
    output_file = output_file,
    write_result = write_result,
    rows = transformed_rows
})

print("transformed rows:", tostring(#transformed_rows))
print("wrote", output_file)
print("wrote artifacts/tutorials/15-read-transform-write-csv-lua.json")
