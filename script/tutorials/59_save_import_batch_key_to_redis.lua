input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
counter_key = os.getenv("TSPLAY_REDIS_BATCH_COUNTER_KEY") or "tutorial:session_import:batch_counter"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"
payload_prefix = os.getenv("TSPLAY_REDIS_BATCH_PREFIX") or "tutorial:session_import:batch:"

rows = read_csv(input_file)
row_count = #rows
latest_name = ""
latest_operator = ""

if row_count > 0 then
    latest_name = rows[row_count].name or ""
    latest_operator = rows[row_count].operator or ""
end

batch_number = redis_incr(counter_key)
batch_id = "session-import-batch-" .. tostring(batch_number)
payload_key = payload_prefix .. batch_id

redis_set(payload_key, {
    lesson = "59",
    batch_id = batch_id,
    input_file = input_file,
    row_count = row_count,
    latest_name = latest_name,
    latest_operator = latest_operator
}, 3600)
redis_set(latest_key, batch_id, 3600)

latest_batch_id = redis_get(latest_key)
batch_payload_text = redis_get(payload_key)
stored_batch_id = json_extract(batch_payload_text, "$.batch_id")
stored_row_count = json_extract(batch_payload_text, "$.row_count")

write_json("artifacts/tutorials/59-save-import-batch-key-to-redis-lua.json", {
    lesson = "59",
    mode = "lua",
    input_file = input_file,
    counter_key = counter_key,
    latest_key = latest_key,
    payload_key = payload_key,
    batch_id = batch_id,
    latest_batch_id = latest_batch_id,
    stored_batch_id = stored_batch_id,
    stored_row_count = stored_row_count,
    batch_payload_text = batch_payload_text
})

print("saved Redis batch key:", tostring(batch_id))
print("wrote artifacts/tutorials/59-save-import-batch-key-to-redis-lua.json")
