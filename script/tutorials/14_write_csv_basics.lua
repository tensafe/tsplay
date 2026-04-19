output_file = "artifacts/tutorials/14-write-csv-basics-lua.csv"

rows = {
    {name = "Alice", city = "Shanghai", status = "ready"},
    {name = "Bob", city = "Hangzhou", status = "reviewing"},
    {name = "Carol", city = "Suzhou", status = "queued"}
}

write_result = write_csv(output_file, rows, {"name", "city", "status"})

write_json("artifacts/tutorials/14-write-csv-basics-lua.json", {
    lesson = "14",
    mode = "lua",
    output_file = output_file,
    write_result = write_result,
    rows = rows
})

print("wrote csv:", output_file)
print("wrote artifacts/tutorials/14-write-csv-basics-lua.json")
