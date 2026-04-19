latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 71 first")
end

payload_key = payload_prefix .. tostring(latest_batch_id)
payload_deleted = redis_del(payload_key)
latest_deleted = redis_del(latest_key)
rows_deleted = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})
batch_deleted = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})

write_csv("artifacts/tutorials/78-cleanup-latest-external-batch-lua.csv", {
    {
        batch_id = tostring(latest_batch_id),
        payload_key = payload_key,
        payload_deleted = payload_deleted,
        latest_deleted = latest_deleted,
        rows_deleted = rows_deleted.rows_affected or 0,
        batch_deleted = batch_deleted.rows_affected or 0
    }
}, {"batch_id", "payload_key", "payload_deleted", "latest_deleted", "rows_deleted", "batch_deleted"})

write_json("artifacts/tutorials/78-cleanup-latest-external-batch-lua.json", {
    lesson = "78",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    payload_key = payload_key,
    payload_deleted = payload_deleted,
    latest_deleted = latest_deleted,
    rows_deleted = rows_deleted,
    batch_deleted = batch_deleted
})

print("cleaned latest external batch:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/78-cleanup-latest-external-batch-lua.json")
