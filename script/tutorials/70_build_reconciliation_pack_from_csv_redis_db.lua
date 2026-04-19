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

detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})

summary_row_count = tonumber(summary_row.row_count) or 0
detail_row_count = tonumber(detail_count_row.row_count) or 0
status = "matched"
if redis_row_count ~= #source_rows or summary_row_count ~= #source_rows or detail_row_count ~= #source_rows then
    status = "mismatch"
    error("reconciliation mismatch for batch " .. tostring(latest_batch_id))
end
if tostring(summary_row.operator_name) ~= redis_operator then
    status = "mismatch"
    error("operator mismatch for batch " .. tostring(latest_batch_id))
end

write_csv("artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-lua.csv", {
    {
        batch_id = tostring(latest_batch_id),
        input_file = tostring(input_file),
        source_row_count = #source_rows,
        redis_row_count = redis_row_count,
        db_summary_row_count = summary_row_count,
        db_detail_row_count = detail_row_count,
        operator_name = redis_operator,
        status = status
    }
}, {"batch_id", "input_file", "source_row_count", "redis_row_count", "db_summary_row_count", "db_detail_row_count", "operator_name", "status"})

write_json("artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-lua.json", {
    lesson = "70",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    payload_key = payload_key,
    input_file = input_file,
    redis_row_count = redis_row_count,
    source_row_count = #source_rows,
    summary_row = summary_row,
    detail_count_row = detail_count_row,
    status = status
})

print("built reconciliation pack for batch:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/70-build-reconciliation-pack-from-csv-redis-db-lua.json")
