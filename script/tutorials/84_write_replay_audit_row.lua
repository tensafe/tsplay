lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
replay_file = os.getenv("TSPLAY_REPLAY_FILE") or "artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.csv"

lifecycle_rows = read_csv(lifecycle_file)
replay_rows = read_csv(replay_file)
if #lifecycle_rows == 0 or #replay_rows == 0 then
    error("lifecycle or replay file is empty, run Lesson 80 and Lesson 82 first")
end

lifecycle_row = lifecycle_rows[1]
replay_row = replay_rows[1]
original_batch_id = tostring(lifecycle_row.batch_id or "")
original_audit_id = tostring(lifecycle_row.audit_id or "")
input_file = tostring(lifecycle_row.input_file or "")
replay_batch_id = tostring(replay_row.replay_batch_id or "")

if original_batch_id == "" or replay_batch_id == "" then
    error("lifecycle or replay file is missing batch identifiers")
end

summary_row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {replay_batch_id}
})
detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {replay_batch_id}
})

audit_id = replay_batch_id .. "-replay-audit"
audit_upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_audits",
    columns = {"audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    key_columns = {"audit_id"},
    update_columns = {"batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    row = {
        audit_id = audit_id,
        batch_id = replay_batch_id,
        event_type = "replay_from_lifecycle",
        status = "ok",
        detail = "replayed from " .. original_batch_id .. " using " .. input_file .. " and source audit " .. original_audit_id,
        success_count = tonumber(detail_count_row.row_count) or 0,
        failure_count = 0,
        source_lesson = "84"
    }
})

audit_row = db_query_one({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE audit_id = $1",
    args = {audit_id}
})

write_json("artifacts/tutorials/84-write-replay-audit-row-lua.json", {
    lesson = "84",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    replay_file = replay_file,
    source_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    audit_id = audit_id,
    summary_row = summary_row,
    detail_count_row = detail_count_row,
    audit_upsert_result = audit_upsert_result,
    audit_row = audit_row
})

print("wrote replay audit row for batch:", replay_batch_id)
print("wrote artifacts/tutorials/84-write-replay-audit-row-lua.json")
