latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 59 first")
end

payload_key = payload_prefix .. tostring(latest_batch_id)
payload_text = redis_get(payload_key)
if payload_text == nil or tostring(payload_text) == "" then
    error("batch payload missing for key: " .. payload_key)
end

input_file = json_extract(payload_text, "$.input_file")
latest_operator = tostring(json_extract(payload_text, "$.latest_operator"))
source_rows = read_csv(input_file, nil, nil, "line_no")
db_rows = db_query({
    connection = "reporting",
    sql = "SELECT line_no, name, phone, status, operator_name FROM public.tutorial_import_rows WHERE batch_id = $1 ORDER BY line_no",
    args = {tostring(latest_batch_id)}
})

if #source_rows ~= #db_rows then
    error("row count mismatch: csv=" .. tostring(#source_rows) .. " db=" .. tostring(#db_rows))
end

for index, source_row in ipairs(source_rows) do
    local db_row = db_rows[index]
    if db_row == nil then
        error("missing db row at index " .. tostring(index))
    end
    if tonumber(source_row.line_no) ~= tonumber(db_row.line_no) then
        error("line_no mismatch at index " .. tostring(index))
    end
    if tostring(source_row.name) ~= tostring(db_row.name) then
        error("name mismatch at index " .. tostring(index))
    end
    if tostring(source_row.phone) ~= tostring(db_row.phone) then
        error("phone mismatch at index " .. tostring(index))
    end
    if tostring(source_row.status) ~= tostring(db_row.status) then
        error("status mismatch at index " .. tostring(index))
    end
    if tostring(db_row.operator_name) ~= latest_operator then
        error("operator mismatch at index " .. tostring(index))
    end
end

write_json("artifacts/tutorials/69-compare-source-csv-and-db-rows-lua.json", {
    lesson = "69",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    input_file = input_file,
    latest_operator = latest_operator,
    row_count = #source_rows,
    status = "matched",
    source_rows = source_rows,
    db_rows = db_rows
})

print("compared source CSV and DB rows:", tostring(#source_rows))
print("wrote artifacts/tutorials/69-compare-source-csv-and-db-rows-lua.json")
