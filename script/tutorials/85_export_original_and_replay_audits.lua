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
replay_batch_id = tostring(replay_row.replay_batch_id or "")

original_audits = db_query({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE batch_id = $1 ORDER BY audit_id",
    args = {original_batch_id}
})
replay_audits = db_query({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE batch_id = $1 ORDER BY audit_id",
    args = {replay_batch_id}
})

comparison_rows = {}
for _, row in ipairs(original_audits) do
    table.insert(comparison_rows, {
        group = "original",
        linked_batch_id = replay_batch_id,
        audit_id = row.audit_id,
        batch_id = row.batch_id,
        event_type = row.event_type,
        status = row.status,
        detail = row.detail,
        success_count = row.success_count,
        failure_count = row.failure_count,
        source_lesson = row.source_lesson
    })
end
for _, row in ipairs(replay_audits) do
    table.insert(comparison_rows, {
        group = "replay",
        linked_batch_id = original_batch_id,
        audit_id = row.audit_id,
        batch_id = row.batch_id,
        event_type = row.event_type,
        status = row.status,
        detail = row.detail,
        success_count = row.success_count,
        failure_count = row.failure_count,
        source_lesson = row.source_lesson
    })
end

write_csv("artifacts/tutorials/85-export-original-and-replay-audits-lua.csv", comparison_rows, {"group", "linked_batch_id", "audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"})

write_json("artifacts/tutorials/85-export-original-and-replay-audits-lua.json", {
    lesson = "85",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    replay_file = replay_file,
    original_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    original_audit_row_count = #original_audits,
    replay_audit_row_count = #replay_audits,
    comparison_rows = comparison_rows
})

print("exported original and replay audits for batches:", original_batch_id, replay_batch_id)
print("wrote artifacts/tutorials/85-export-original-and-replay-audits-lua.json")
