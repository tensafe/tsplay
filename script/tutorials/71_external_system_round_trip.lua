input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
counter_key = os.getenv("TSPLAY_REDIS_BATCH_COUNTER_KEY") or "tutorial:session_import:batch_counter"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

rows = read_csv(input_file, nil, nil, "line_no")
detail_rows = {}
latest_name = ""
latest_operator = "unknown"

for _, row in ipairs(rows) do
    latest_name = row.name or latest_name
    latest_operator = row.operator or latest_operator
    table.insert(detail_rows, {
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = row.operator
    })
end

batch_number = redis_incr(counter_key)
batch_id = "session-import-batch-" .. tostring(batch_number)
payload_key = payload_prefix .. batch_id

redis_set(payload_key, {
    lesson = "71",
    batch_id = batch_id,
    input_file = input_file,
    row_count = #rows,
    latest_name = latest_name,
    latest_operator = latest_operator
}, 3600)
redis_set(latest_key, batch_id, 3600)

db_detail_rows = {}
for _, row in ipairs(detail_rows) do
    table.insert(db_detail_rows, {
        batch_id = batch_id,
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = row.operator_name
    })
end

transaction_result = db_transaction(function()
    db_upsert({
        connection = "reporting",
        table = "public.tutorial_import_batches",
        columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
        key_columns = {"batch_id"},
        update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
        row = {
            batch_id = batch_id,
            report_file = input_file,
            source_lesson = "71",
            row_count = #rows,
            operator_name = latest_operator
        }
    })

    db_execute({
        connection = "reporting",
        sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
        args = {batch_id}
    })

    return db_insert_many({
        connection = "reporting",
        table = "public.tutorial_import_rows",
        columns = {"batch_id", "line_no", "name", "phone", "status", "operator_name"},
        rows = db_detail_rows
    })
end, 5000)

latest_batch_id = redis_get(latest_key)
payload_text = redis_get(payload_key)
summary_row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, source_lesson, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})
detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

write_csv("artifacts/tutorials/71-external-system-round-trip-lua.csv", {
    {
        batch_id = batch_id,
        input_file = input_file,
        row_count = #rows,
        latest_name = latest_name,
        latest_operator = latest_operator,
        latest_batch_id = latest_batch_id
    }
}, {"batch_id", "input_file", "row_count", "latest_name", "latest_operator", "latest_batch_id"})

write_json("artifacts/tutorials/71-external-system-round-trip-lua.json", {
    lesson = "71",
    mode = "lua",
    input_file = input_file,
    batch_id = batch_id,
    payload_key = payload_key,
    transaction_result = transaction_result,
    latest_batch_id = latest_batch_id,
    payload_text = payload_text,
    summary_row = summary_row,
    detail_count_row = detail_count_row
})

print("completed external system round trip:", tostring(batch_id))
print("wrote artifacts/tutorials/71-external-system-round-trip-lua.json")
