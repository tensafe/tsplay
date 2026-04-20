lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
counter_key = os.getenv("TSPLAY_REDIS_BATCH_COUNTER_KEY") or "tutorial:session_import:batch_counter"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

lifecycle_rows = read_csv(lifecycle_file)
if #lifecycle_rows == 0 then
    error("lifecycle file is empty, run Lesson 80 first")
end

lifecycle_row = lifecycle_rows[1]
original_batch_id = tostring(lifecycle_row.batch_id or "")
original_audit_id = tostring(lifecycle_row.audit_id or "")
input_file = tostring(lifecycle_row.input_file or "")
expected_row_count = tonumber(lifecycle_row.pre_cleanup_detail) or 0

rows = read_csv(input_file, nil, nil, "line_no")
if #rows == 0 then
    error("input file is empty: " .. input_file)
end

latest_operator = "unknown"
detail_rows = {}
for _, row in ipairs(rows) do
    latest_operator = row.operator or latest_operator
end

handoff_number = redis_incr(counter_key)
replay_batch_id = original_batch_id .. "-handoff-" .. tostring(handoff_number)
payload_key = payload_prefix .. replay_batch_id

for _, row in ipairs(rows) do
    table.insert(detail_rows, {
        batch_id = replay_batch_id,
        line_no = row.line_no,
        name = row.name,
        phone = row.phone,
        status = row.status,
        operator_name = row.operator
    })
end

redis_set(payload_key, {
    lesson = "90",
    batch_id = replay_batch_id,
    source_batch_id = original_batch_id,
    source_audit_id = original_audit_id,
    input_file = input_file,
    row_count = #rows,
    latest_operator = latest_operator
}, 3600)
redis_set(latest_key, replay_batch_id, 3600)

transaction_result = db_transaction(function()
    db_upsert({
        connection = "reporting",
        table = "public.tutorial_import_batches",
        columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
        key_columns = {"batch_id"},
        update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
        row = {
            batch_id = replay_batch_id,
            report_file = input_file,
            source_lesson = "90",
            row_count = #rows,
            operator_name = latest_operator
        }
    })

    db_execute({
        connection = "reporting",
        sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
        args = {replay_batch_id}
    })

    return db_insert_many({
        connection = "reporting",
        table = "public.tutorial_import_rows",
        columns = {"batch_id", "line_no", "name", "phone", "status", "operator_name"},
        rows = detail_rows
    })
end, 5000)

replay_audit_id = replay_batch_id .. "-handoff-audit"
audit_upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_audits",
    columns = {"audit_id", "batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    key_columns = {"audit_id"},
    update_columns = {"batch_id", "event_type", "status", "detail", "success_count", "failure_count", "source_lesson"},
    row = {
        audit_id = replay_audit_id,
        batch_id = replay_batch_id,
        event_type = "handoff_round_trip",
        status = "ok",
        detail = "replayed from lifecycle evidence and prepared handoff package",
        success_count = #rows,
        failure_count = 0,
        source_lesson = "90"
    }
})

detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {replay_batch_id}
})
audit_row = db_query_one({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE audit_id = $1",
    args = {replay_audit_id}
})

db_row_count = tonumber(detail_count_row.row_count) or 0
if db_row_count ~= expected_row_count then
    error("handoff row count mismatch: db=" .. tostring(db_row_count) .. " expected=" .. tostring(expected_row_count))
end

manifest_rows = {
    {
        artifact_key = "lifecycle_evidence",
        batch_id = original_batch_id,
        related_batch_id = replay_batch_id,
        row_count = expected_row_count,
        note = "source lifecycle evidence from lesson 80"
    },
    {
        artifact_key = "handoff_replay",
        batch_id = replay_batch_id,
        related_batch_id = original_batch_id,
        row_count = db_row_count,
        note = "new replay batch prepared for handoff"
    },
    {
        artifact_key = "handoff_audit",
        batch_id = replay_batch_id,
        related_batch_id = replay_audit_id,
        row_count = #rows,
        note = "audit row written for handoff batch"
    }
}

checklist_rows = {
    {artifact_key = "lifecycle_evidence", status = "ready", note = "lifecycle evidence loaded"},
    {artifact_key = "handoff_replay", status = "ready", note = "replay batch persisted"},
    {artifact_key = "handoff_audit", status = "ready", note = "audit row written"},
}

write_csv("artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.csv", {
    {
        original_batch_id = original_batch_id,
        replay_batch_id = replay_batch_id,
        replay_audit_id = replay_audit_id,
        expected_row_count = expected_row_count,
        replay_row_count = db_row_count,
        latest_batch_id = redis_get(latest_key),
        status = "ok"
    }
}, {"original_batch_id", "replay_batch_id", "replay_audit_id", "expected_row_count", "replay_row_count", "latest_batch_id", "status"})

write_json("artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.json", {
    lesson = "90",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    input_file = input_file,
    original_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    replay_audit_id = replay_audit_id,
    transaction_result = transaction_result,
    audit_upsert_result = audit_upsert_result,
    detail_count_row = detail_count_row,
    audit_row = audit_row,
    manifest_rows = manifest_rows,
    checklist_rows = checklist_rows
})

print("completed handoff round trip from lifecycle evidence:", replay_batch_id)
print("wrote artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-lua.json")
