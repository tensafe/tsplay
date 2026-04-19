latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 72 first")
end

payload_text = redis_get(payload_prefix .. tostring(latest_batch_id))
if payload_text == nil or tostring(payload_text) == "" then
    error("batch payload missing for latest batch")
end

input_file = tostring(json_extract(payload_text, "$.input_file"))
source_rows = read_csv(input_file, nil, nil, "line_no")
db_rows = db_query({
    connection = "reporting",
    sql = "SELECT line_no, name, phone, status, operator_name FROM public.tutorial_import_rows WHERE batch_id = $1 ORDER BY line_no",
    args = {tostring(latest_batch_id)}
})

line_seen = {}
for index, row in ipairs(db_rows) do
    local line_no = tonumber(row.line_no)
    if line_seen[line_no] then
        error("duplicate line_no found in db rows: " .. tostring(line_no))
    end
    line_seen[line_no] = true
    if source_rows[index] == nil then
        error("db has more rows than source file")
    end
end

if #source_rows ~= #db_rows then
    error("row count mismatch after rerun: csv=" .. tostring(#source_rows) .. " db=" .. tostring(#db_rows))
end

write_json("artifacts/tutorials/73-verify-rerun-does-not-duplicate-rows-lua.json", {
    lesson = "73",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    input_file = input_file,
    source_row_count = #source_rows,
    db_row_count = #db_rows,
    status = "deduplicated",
    db_rows = db_rows
})

print("verified rerun has no duplicate rows:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/73-verify-rerun-does-not-duplicate-rows-lua.json")
