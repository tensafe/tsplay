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
latest_operator = json_extract(payload_text, "$.latest_operator")
rows = read_csv(input_file, nil, nil, "line_no")
row_count = #rows

upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    key_columns = {"batch_id"},
    update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
    returning = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = tostring(latest_batch_id),
        report_file = tostring(input_file),
        source_lesson = "65",
        row_count = row_count,
        operator_name = tostring(latest_operator)
    }
})

summary_row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, source_lesson, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})

write_json("artifacts/tutorials/65-sync-latest-redis-batch-to-postgres-summary-lua.json", {
    lesson = "65",
    mode = "lua",
    latest_key = latest_key,
    payload_key = payload_key,
    latest_batch_id = latest_batch_id,
    input_file = input_file,
    row_count = row_count,
    latest_operator = latest_operator,
    upsert_result = upsert_result,
    summary_row = summary_row
})

print("synced Redis batch to Postgres summary:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/65-sync-latest-redis-batch-to-postgres-summary-lua.json")
