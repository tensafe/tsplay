cleanup_file = os.getenv("TSPLAY_CLEANUP_FILE") or "artifacts/tutorials/78-cleanup-latest-external-batch-flow.csv"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"

cleanup_rows = read_csv(cleanup_file)
if #cleanup_rows == 0 then
    error("cleanup file is empty: " .. cleanup_file)
end

batch_id = tostring(cleanup_rows[1].batch_id)
payload_key = tostring(cleanup_rows[1].payload_key)
payload_after = redis_get(payload_key)
latest_after = redis_get(latest_key)

summary_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {batch_id}
})

audit_rows = db_query({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE batch_id = $1 ORDER BY audit_id",
    args = {batch_id}
})

if payload_after ~= nil and tostring(payload_after) ~= "" then
    error("expected payload key deleted, got: " .. tostring(payload_after))
end
if latest_after ~= nil and tostring(latest_after) ~= "" then
    error("expected latest pointer deleted, got: " .. tostring(latest_after))
end
if tonumber(summary_count_row.row_count) ~= 0 then
    error("expected no summary rows for cleaned batch")
end
if tonumber(detail_count_row.row_count) ~= 0 then
    error("expected no detail rows for cleaned batch")
end
if #audit_rows == 0 then
    error("expected audit rows to remain for cleaned batch")
end

write_json("artifacts/tutorials/79-verify-external-batch-cleanup-lua.json", {
    lesson = "79",
    mode = "lua",
    cleanup_file = cleanup_file,
    batch_id = batch_id,
    payload_key = payload_key,
    summary_count_row = summary_count_row,
    detail_count_row = detail_count_row,
    audit_rows = audit_rows,
    status = "cleaned_with_audit_retained"
})

print("verified external batch cleanup:", tostring(batch_id))
print("wrote artifacts/tutorials/79-verify-external-batch-cleanup-lua.json")
