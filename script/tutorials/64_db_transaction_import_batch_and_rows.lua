input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
batch_id = os.getenv("TSPLAY_IMPORT_BATCH_ID") or "lesson-64-import-batch"

rows = read_csv(input_file, nil, nil, "line_no")
row_count = #rows
operator_name = "unknown"
detail_rows = {}

for _, row in ipairs(rows) do
    operator_name = row.operator or operator_name
    table.insert(detail_rows, {
        batch_id = batch_id,
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = row.operator
    })
end

cleanup_rows = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

cleanup_batch = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

transaction_result = db_transaction(function()
    db_insert({
        connection = "reporting",
        table = "public.tutorial_import_batches",
        columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
        row = {
            batch_id = batch_id,
            report_file = input_file,
            source_lesson = "64",
            row_count = row_count,
            operator_name = operator_name
        }
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
    args = {batch_id}
})

stored_rows = db_query({
    connection = "reporting",
    sql = "SELECT line_no, name, phone, status, operator_name FROM public.tutorial_import_rows WHERE batch_id = $1 ORDER BY line_no",
    args = {batch_id}
})

write_json("artifacts/tutorials/64-db-transaction-import-batch-and-rows-lua.json", {
    lesson = "64",
    mode = "lua",
    input_file = input_file,
    batch_id = batch_id,
    cleanup_rows = cleanup_rows,
    cleanup_batch = cleanup_batch,
    transaction_result = transaction_result,
    summary_row = summary_row,
    stored_rows = stored_rows
})

print("transactionally stored import batch:", tostring(batch_id))
print("wrote artifacts/tutorials/64-db-transaction-import-batch-and-rows-lua.json")
