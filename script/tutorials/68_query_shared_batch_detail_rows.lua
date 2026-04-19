latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"

latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 59 first")
end

stored_rows = db_query({
    connection = "reporting",
    sql = "SELECT line_no, name, phone, status, operator_name FROM public.tutorial_import_rows WHERE batch_id = $1 ORDER BY line_no",
    args = {tostring(latest_batch_id)}
})

detail_count_row = db_query_one({
    connection = "reporting",
    sql = "SELECT COUNT(*) AS row_count FROM public.tutorial_import_rows WHERE batch_id = $1",
    args = {tostring(latest_batch_id)}
})

if #stored_rows == 0 then
    error("no detail rows found for batch: " .. tostring(latest_batch_id))
end

write_json("artifacts/tutorials/68-query-shared-batch-detail-rows-lua.json", {
    lesson = "68",
    mode = "lua",
    latest_batch_id = latest_batch_id,
    detail_count_row = detail_count_row,
    stored_rows = stored_rows
})

print("queried shared batch detail rows:", tostring(#stored_rows))
print("wrote artifacts/tutorials/68-query-shared-batch-detail-rows-lua.json")
