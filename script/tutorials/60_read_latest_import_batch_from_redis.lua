latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

latest_batch_id = redis_get(latest_key)
if latest_batch_id == nil or tostring(latest_batch_id) == "" then
    error("latest batch id is empty, run Lesson 59 first")
end

payload_key = payload_prefix .. tostring(latest_batch_id)
batch_payload_text = redis_get(payload_key)
if batch_payload_text == nil or tostring(batch_payload_text) == "" then
    error("batch payload missing for key: " .. payload_key)
end

stored_batch_id = json_extract(batch_payload_text, "$.batch_id")
stored_row_count = json_extract(batch_payload_text, "$.row_count")
stored_latest_name = json_extract(batch_payload_text, "$.latest_name")
stored_latest_operator = json_extract(batch_payload_text, "$.latest_operator")

write_json("artifacts/tutorials/60-read-latest-import-batch-from-redis-lua.json", {
    lesson = "60",
    mode = "lua",
    latest_key = latest_key,
    payload_key = payload_key,
    latest_batch_id = latest_batch_id,
    stored_batch_id = stored_batch_id,
    stored_row_count = stored_row_count,
    stored_latest_name = stored_latest_name,
    stored_latest_operator = stored_latest_operator,
    batch_payload_text = batch_payload_text
})

print("loaded latest Redis batch:", tostring(latest_batch_id))
print("wrote artifacts/tutorials/60-read-latest-import-batch-from-redis-lua.json")
