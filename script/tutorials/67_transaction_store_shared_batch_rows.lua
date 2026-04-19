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
rows = read_csv(input_file, nil, nil, "line_no")
detail_rows = {}

for _, row in ipairs(rows) do
    table.insert(detail_rows, {
        batch_id = tostring(latest_batch_id),
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = latest_operator
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
            batch_id = tostring(latest_batch_id),
            report_file = tostring(input_file),
            source_lesson = "67",
            row_count = #rows,
            operator_name = latest_operator
        }
    })

    db_execute({
        connection = "reporting",
        sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
        args = {tostring(latest_batch_id)}
    })

    return db_insert_many({
        connection = "reporting",
        table = "public.tutorial_import_rows",
        columns = {"batch_id", "line_no", "name", "phone", "status", "operator_name"},
        rows = detail_rows
    })
end, 5000)

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

write_json("artifacts/tutorials/67-transaction-store-shared-batch-rows-lua.json", {
    lesson = "67",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    payload_key = payload_key,
    input_file = input_file,
    latest_operator = latest_operator,
    transaction_result = transaction_result,
    summary_row = summary_row,
    detail_count_row = detail_count_row
})

print("stored shared batch rows in Postgres:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/67-transaction-store-shared-batch-rows-lua.json")
