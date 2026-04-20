lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"

lifecycle_rows = read_csv(lifecycle_file)
if #lifecycle_rows == 0 then
    error("lifecycle file is empty, run Lesson 80 first")
end

lifecycle_row = lifecycle_rows[1]
batch_id = tostring(lifecycle_row.batch_id or "")
audit_id = tostring(lifecycle_row.audit_id or "")
input_file = tostring(lifecycle_row.input_file or "")

if batch_id == "" or audit_id == "" or input_file == "" then
    error("lifecycle file is missing batch_id, audit_id, or input_file")
end

audit_row = db_query_one({
    connection = "reporting",
    sql = "SELECT audit_id, batch_id, event_type, status, detail, success_count, failure_count, source_lesson FROM public.tutorial_import_audits WHERE audit_id = $1",
    args = {audit_id}
})

write_json("artifacts/tutorials/81-read-lifecycle-evidence-lua.json", {
    lesson = "81",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    lifecycle_row = lifecycle_row,
    batch_id = batch_id,
    audit_id = audit_id,
    input_file = input_file,
    audit_row = audit_row
})

print("read lifecycle evidence for batch:", batch_id)
print("wrote artifacts/tutorials/81-read-lifecycle-evidence-lua.json")
