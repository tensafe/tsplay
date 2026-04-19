input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
summary_key = os.getenv("TSPLAY_REDIS_SUMMARY_KEY") or "tutorial:session_import:latest_summary"

rows = read_csv(input_file)
row_count = #rows
latest_name = ""
latest_operator = ""

if row_count > 0 then
    latest_name = rows[row_count].name or ""
    latest_operator = rows[row_count].operator or ""
end

redis_del(summary_key)
redis_set(summary_key, {
    lesson = "58",
    source_lesson = "57",
    input_file = input_file,
    row_count = row_count,
    latest_name = latest_name,
    latest_operator = latest_operator
}, 3600)

summary_text = redis_get(summary_key)
stored_row_count = json_extract(summary_text, "$.row_count")
stored_operator = json_extract(summary_text, "$.latest_operator")

write_json("artifacts/tutorials/58-sync-import-report-summary-to-redis-lua.json", {
    lesson = "58",
    mode = "lua",
    input_file = input_file,
    summary_key = summary_key,
    row_count = row_count,
    latest_name = latest_name,
    latest_operator = latest_operator,
    stored_row_count = stored_row_count,
    stored_operator = stored_operator,
    summary_text = summary_text
})

print("cached Redis summary for rows:", tostring(row_count))
print("wrote artifacts/tutorials/58-sync-import-report-summary-to-redis-lua.json")
