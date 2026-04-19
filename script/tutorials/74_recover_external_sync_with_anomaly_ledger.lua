input_file = os.getenv("TSPLAY_ANOMALY_INPUT") or "demo/data/import_report_with_issue.csv"
batch_id = os.getenv("TSPLAY_ANOMALY_BATCH_ID") or "lesson-74-anomaly-batch"

rows = read_csv(input_file, nil, nil, "line_no")
results = {}
last_operator = "demo-user"

db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

db_insert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_id,
        report_file = input_file,
        source_lesson = "74",
        row_count = 0,
        operator_name = last_operator
    }
})

for _, row in ipairs(rows) do
    local ok, err = pcall(function()
        if row.phone == nil or tostring(row.phone) == "" then
            error("phone is required")
        end
        if row.status == nil or tostring(row.status) ~= "Imported" then
            error("status must be Imported")
        end
        last_operator = row.operator or last_operator
        db_insert({
            connection = "reporting",
            table = "public.tutorial_import_rows",
            columns = {"batch_id", "line_no", "name", "phone", "status", "operator_name"},
            row = {
                batch_id = batch_id,
                line_no = row.line_no,
                name = row.name,
                phone = row.phone,
                status = row.status,
                operator_name = row.operator or last_operator
            }
        })
    end)

    if ok then
        table.insert(results, {
            line_no = row.line_no,
            name = row.name,
            phone = row.phone,
            status = "success",
            error = ""
        })
    else
        table.insert(results, {
            line_no = row.line_no,
            name = row.name,
            phone = row.phone or "",
            status = "failed",
            error = tostring(err)
        })
    end
end

stored_success_count = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    key_columns = {"batch_id"},
    update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_id,
        report_file = input_file,
        source_lesson = "74",
        row_count = tonumber(stored_success_count.row_count) or 0,
        operator_name = last_operator
    }
})

write_csv("artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-lua.csv", results, {"line_no", "name", "phone", "status", "error"})

write_json("artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-lua.json", {
    lesson = "74",
    mode = "lua",
    input_file = input_file,
    batch_id = batch_id,
    stored_success_count = stored_success_count,
    results = results
})

print("recovered anomaly batch with ledger:", tostring(batch_id))
print("wrote artifacts/tutorials/74-recover-external-sync-with-anomaly-ledger-lua.json")
