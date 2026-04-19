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
        batch_id = "",
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
for _, row in ipairs(detail_rows) do
    row.batch_id = batch_id
end

redis_set(payload_key, {
    lesson = "80",
    batch_id = batch_id,
    input_file = input_file,
    row_count = #rows,
    latest_name = latest_name,
    latest_operator = latest_operator
}, 3600)
redis_set(latest_key, batch_id, 3600)

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
            source_lesson = "80",
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
        rows = detail_rows
    })
end, 5000)

audit_id = batch_id .. "-lifecycle-audit"
audit_upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_audits",
    columns = {"audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    key_columns = {"audit_id"},
    update_columns = {"batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    row = {
        audit_id = audit_id,
        batch_id = batch_id,
        event_type = "lifecycle_round_trip",
        status = "ok",
        detail = "batch created, persisted, cleaned, and verified",
        success_count = #rows,
        failure_count = 0,
        source_lesson = "80"
    }
})

pre_cleanup_summary = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})
pre_cleanup_detail = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

payload_deleted = redis_del(payload_key)
latest_deleted = redis_del(latest_key)
rows_deleted = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})
batch_deleted = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

post_cleanup_summary = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})
post_cleanup_detail = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})
audit_row = db_query_one({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE audit_id = $1",
    args = {audit_id}
})

write_csv("artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.csv", {
    {
        batch_id = batch_id,
        audit_id = audit_id,
        pre_cleanup_summary = pre_cleanup_summary.row_count,
        pre_cleanup_detail = pre_cleanup_detail.row_count,
        post_cleanup_summary = post_cleanup_summary.row_count,
        post_cleanup_detail = post_cleanup_detail.row_count,
        payload_deleted = payload_deleted,
        latest_deleted = latest_deleted
    }
}, {"batch_id", "audit_id", "pre_cleanup_summary", "pre_cleanup_detail", "post_cleanup_summary", "post_cleanup_detail", "payload_deleted", "latest_deleted"})

write_json("artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.json", {
    lesson = "80",
    mode = "lua",
    input_file = input_file,
    batch_id = batch_id,
    audit_id = audit_id,
    payload_key = payload_key,
    transaction_result = transaction_result,
    audit_upsert_result = audit_upsert_result,
    pre_cleanup_summary = pre_cleanup_summary,
    pre_cleanup_detail = pre_cleanup_detail,
    rows_deleted = rows_deleted,
    batch_deleted = batch_deleted,
    post_cleanup_summary = post_cleanup_summary,
    post_cleanup_detail = post_cleanup_detail,
    audit_row = audit_row
})

print("completed external sync lifecycle round trip:", tostring(batch_id))
print("wrote artifacts/tutorials/80-external-sync-lifecycle-round-trip-lua.json")
