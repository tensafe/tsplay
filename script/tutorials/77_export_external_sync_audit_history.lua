latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 75 first")
end

audit_rows = db_query({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE batch_id = $1 ORDER BY audit_id",
    args = {tostring(latest_batch_id)}
})

write_csv("artifacts/tutorials/77-export-external-sync-audit-history-lua.csv", audit_rows, {"audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"})

write_json("artifacts/tutorials/77-export-external-sync-audit-history-lua.json", {
    lesson = "77",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    audit_row_count = #audit_rows,
    audit_rows = audit_rows
})

print("exported audit history for batch:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/77-export-external-sync-audit-history-lua.json")
