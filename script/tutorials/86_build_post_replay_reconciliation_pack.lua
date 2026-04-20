lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
replay_file = os.getenv("TSPLAY_REPLAY_FILE") or "artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.csv"
audit_compare_file = os.getenv("TSPLAY_AUDIT_COMPARE_FILE") or "artifacts/tutorials/85-export-original-and-replay-audits-flow.csv"

lifecycle_rows = read_csv(lifecycle_file)
replay_rows = read_csv(replay_file)
audit_rows = read_csv(audit_compare_file)
if #lifecycle_rows == 0 or #replay_rows == 0 then
    error("lifecycle or replay file is empty, run Lesson 80 and Lesson 82 first")
end

lifecycle_row = lifecycle_rows[1]
replay_row = replay_rows[1]
original_batch_id = tostring(lifecycle_row.batch_id or "")
replay_batch_id = tostring(replay_row.replay_batch_id or "")
expected_row_count = tonumber(lifecycle_row.pre_cleanup_detail) or 0
payload_key = tostring(replay_row.payload_key or "")
replay_row_count = tonumber(replay_row.row_count) or 0

payload_text = redis_get(payload_key)
if payload_text == nil or tostring(payload_text) == "" then
    error("replay payload is empty, run Lesson 82 first")
end

payload_row_count = tonumber(json_extract(payload_text, "$.row_count")) or 0
summary_row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {replay_batch_id}
})
detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {replay_batch_id}
})

db_row_count = tonumber(detail_count_row.row_count) or 0
summary_row_count = tonumber(summary_row.row_count) or 0
if replay_row_count ~= expected_row_count then
    error("replay csv row count mismatch: replay=" .. tostring(replay_row_count) .. " expected=" .. tostring(expected_row_count))
end
if payload_row_count ~= expected_row_count then
    error("payload row count mismatch: payload=" .. tostring(payload_row_count) .. " expected=" .. tostring(expected_row_count))
end
if db_row_count ~= expected_row_count or summary_row_count ~= expected_row_count then
    error("database row count mismatch: summary=" .. tostring(summary_row_count) .. " detail=" .. tostring(db_row_count) .. " expected=" .. tostring(expected_row_count))
end

original_audit_row_count = 0
replay_audit_row_count = 0
for _, row in ipairs(audit_rows) do
    if tostring(row.group or "") == "original" then
        original_audit_row_count = original_audit_row_count + 1
    elseif tostring(row.group or "") == "replay" then
        replay_audit_row_count = replay_audit_row_count + 1
    end
end

report_rows = {
    {
        original_batch_id = original_batch_id,
        replay_batch_id = replay_batch_id,
        expected_row_count = expected_row_count,
        replay_row_count = replay_row_count,
        payload_row_count = payload_row_count,
        db_row_count = db_row_count,
        original_audit_row_count = original_audit_row_count,
        replay_audit_row_count = replay_audit_row_count,
        status = "ok"
    }
}

write_csv("artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.csv", report_rows, {"original_batch_id", "replay_batch_id", "expected_row_count", "replay_row_count", "payload_row_count", "db_row_count", "original_audit_row_count", "replay_audit_row_count", "status"})

write_json("artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.json", {
    lesson = "86",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    replay_file = replay_file,
    audit_compare_file = audit_compare_file,
    summary_row = summary_row,
    detail_count_row = detail_count_row,
    report_rows = report_rows
})

print("built post-replay reconciliation pack for batch:", replay_batch_id)
print("wrote artifacts/tutorials/86-build-post-replay-reconciliation-pack-lua.json")
