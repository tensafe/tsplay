latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 71 first")
end

audit_id = tostring(latest_batch_id) .. "-sync-audit"
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

upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_audits",
    columns = {"audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    key_columns = {"audit_id"},
    update_columns = {"batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    row = {
        audit_id = audit_id,
        batch_id = tostring(latest_batch_id),
        event_type = "external_sync",
        status = "ok",
        detail = "latest shared batch synchronized and checked",
        success_count = tonumber(detail_count_row.row_count) or 0,
        failure_count = 0,
        source_lesson = "75"
    }
})

audit_row = db_query_one({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE audit_id = $1",
    args = {audit_id}
})

write_json("artifacts/tutorials/75-write-external-sync-audit-row-lua.json", {
    lesson = "75",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    audit_id = audit_id,
    summary_row = summary_row,
    detail_count_row = detail_count_row,
    upsert_result = upsert_result,
    audit_row = audit_row
})

print("wrote audit row for batch:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/75-write-external-sync-audit-row-lua.json")
