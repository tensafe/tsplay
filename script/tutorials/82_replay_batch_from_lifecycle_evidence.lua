lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
counter_key = os.getenv("TSPLAY_REDIS_BATCH_COUNTER_KEY") or "tutorial:session_import:batch_counter"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

lifecycle_rows = read_csv(lifecycle_file)
if #lifecycle_rows == 0 then
    error("lifecycle file is empty, run Lesson 80 first")
end

lifecycle_row = lifecycle_rows[1]
original_batch_id = tostring(lifecycle_row.batch_id or "")
original_audit_id = tostring(lifecycle_row.audit_id or "")
input_file = tostring(lifecycle_row.input_file or "")

if original_batch_id == "" or input_file == "" then
    error("lifecycle file is missing batch_id or input_file")
end

rows = read_csv(input_file, nil, nil, "line_no")
if #rows == 0 then
    error("input file is empty: " .. input_file)
end

latest_name = ""
latest_operator = "unknown"
for _, row in ipairs(rows) do
    latest_name = row.name or latest_name
    latest_operator = row.operator or latest_operator
end

replay_number = redis_incr(counter_key)
replay_batch_id = original_batch_id .. "-replay-" .. tostring(replay_number)
payload_key = payload_prefix .. replay_batch_id

detail_rows = {}
for _, row in ipairs(rows) do
    table.insert(detail_rows, {
        batch_id = replay_batch_id,
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = row.operator
    })
end

redis_set(payload_key, {
    lesson = "82",
    batch_id = replay_batch_id,
    source_batch_id = original_batch_id,
    source_audit_id = original_audit_id,
    input_file = input_file,
    row_count = #rows,
    latest_name = latest_name,
    latest_operator = latest_operator
}, 3600)
redis_set(latest_key, replay_batch_id, 3600)

transaction_result = db_transaction(function()
    db_upsert({
        connection = "reporting",
        table = "public.tutorial_import_batches",
        columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
        key_columns = {"batch_id"},
        update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
        row = {
            batch_id = replay_batch_id,
            report_file = input_file,
            source_lesson = "82",
            row_count = #rows,
            operator_name = latest_operator
        }
    })

    db_execute({
        connection = "reporting",
        sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
        args = {replay_batch_id}
    })

    return db_insert_many({
        connection = "reporting",
        table = "public.tutorial_import_rows",
        columns = {"batch_id", "line_no", "name", "phone", "status", "operator_name"},
        rows = detail_rows
    })
end, 5000)

latest_batch_id = redis_get(latest_key)

write_csv("artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.csv", {
    {
        source_batch_id = original_batch_id,
        source_audit_id = original_audit_id,
        replay_batch_id = replay_batch_id,
        payload_key = payload_key,
        input_file = input_file,
        row_count = #rows,
        latest_operator = latest_operator,
        latest_batch_id = latest_batch_id
    }
}, {"source_batch_id", "source_audit_id", "replay_batch_id", "payload_key", "input_file", "row_count", "latest_operator", "latest_batch_id"})

write_json("artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.json", {
    lesson = "82",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    source_batch_id = original_batch_id,
    source_audit_id = original_audit_id,
    replay_batch_id = replay_batch_id,
    payload_key = payload_key,
    input_file = input_file,
    row_count = #rows,
    transaction_result = transaction_result,
    latest_batch_id = latest_batch_id
})

print("replayed lifecycle evidence into batch:", replay_batch_id)
print("wrote artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-lua.json")
