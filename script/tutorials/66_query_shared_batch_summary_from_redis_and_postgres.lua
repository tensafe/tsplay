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
redis_row_count = tonumber(json_extract(payload_text, "$.row_count")) or 0
redis_operator = tostring(json_extract(payload_text, "$.latest_operator"))
source_rows = read_csv(input_file, nil, nil, "line_no")

summary_row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, source_lesson, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})

db_row_count = tonumber(summary_row.row_count) or 0
if db_row_count ~= redis_row_count then
    error("row_count mismatch: redis=" .. tostring(redis_row_count) .. " db=" .. tostring(db_row_count))
end
if #source_rows ~= redis_row_count then
    error("source rows mismatch: csv=" .. tostring(#source_rows) .. " redis=" .. tostring(redis_row_count))
end
if tostring(summary_row.operator_name) ~= redis_operator then
    error("operator mismatch: redis=" .. redis_operator .. " db=" .. tostring(summary_row.operator_name))
end

write_json("artifacts/tutorials/66-query-shared-batch-summary-from-redis-and-postgres-lua.json", {
    lesson = "66",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    payload_key = payload_key,
    input_file = input_file,
    redis_row_count = redis_row_count,
    csv_row_count = #source_rows,
    redis_operator = redis_operator,
    summary_row = summary_row,
    status = "matched"
})

print("verified shared summary across Redis and Postgres:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/66-query-shared-batch-summary-from-redis-and-postgres-lua.json")
