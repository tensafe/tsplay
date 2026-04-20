lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
replay_file = os.getenv("TSPLAY_REPLAY_FILE") or "artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.csv"
latest_key = os.getenv("TSPLAY_REDIS_LATEST_BATCH_KEY") or "tutorial:session_import:latest_batch"

lifecycle_rows = read_csv(lifecycle_file)
replay_rows = read_csv(replay_file)
if #lifecycle_rows == 0 or #replay_rows == 0 then
    error("lifecycle or replay file is empty, run Lesson 80 and Lesson 82 first")
end

lifecycle_row = lifecycle_rows[1]
replay_row = replay_rows[1]

original_batch_id = tostring(lifecycle_row.batch_id or "")
expected_row_count = tonumber(lifecycle_row.pre_cleanup_detail) or 0
input_file = tostring(lifecycle_row.input_file or "")
replay_batch_id = tostring(replay_row.replay_batch_id or "")
payload_key = tostring(replay_row.payload_key or "")

if original_batch_id == "" or replay_batch_id == "" or payload_key == "" then
    error("lifecycle or replay file is missing required fields")
end

payload_text = redis_get(payload_key)
if payload_text == nil or tostring(payload_text) == "" then
    error("replay payload is empty, run Lesson 82 first")
end

latest_batch_after = redis_get(latest_key)
payload_row_count = tonumber(json_extract(payload_text, "$.row_count")) or 0
payload_source_batch_id = tostring(json_extract(payload_text, "$.source_batch_id") or "")
payload_input_file = tostring(json_extract(payload_text, "$.input_file") or "")

source_rows = read_csv(input_file)
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

summary_row_count = tonumber(summary_row.row_count) or 0
detail_row_count = tonumber(detail_count_row.row_count) or 0

if #source_rows ~= expected_row_count then
    error("source row count mismatch: source=" .. tostring(#source_rows) .. " expected=" .. tostring(expected_row_count))
end
if payload_row_count ~= expected_row_count then
    error("payload row count mismatch: payload=" .. tostring(payload_row_count) .. " expected=" .. tostring(expected_row_count))
end
if summary_row_count ~= expected_row_count then
    error("summary row count mismatch: summary=" .. tostring(summary_row_count) .. " expected=" .. tostring(expected_row_count))
end
if detail_row_count ~= expected_row_count then
    error("detail row count mismatch: detail=" .. tostring(detail_row_count) .. " expected=" .. tostring(expected_row_count))
end
if payload_source_batch_id ~= original_batch_id then
    error("source batch mismatch: payload=" .. payload_source_batch_id .. " lifecycle=" .. original_batch_id)
end
if payload_input_file ~= input_file then
    error("input file mismatch: payload=" .. payload_input_file .. " lifecycle=" .. input_file)
end
if tostring(latest_batch_after or "") ~= replay_batch_id then
    error("latest batch mismatch: latest=" .. tostring(latest_batch_after) .. " replay=" .. replay_batch_id)
end

write_json("artifacts/tutorials/83-verify-replay-batch-against-lifecycle-evidence-lua.json", {
    lesson = "83",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    replay_file = replay_file,
    original_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    input_file = input_file,
    expected_row_count = expected_row_count,
    payload_row_count = payload_row_count,
    summary_row = summary_row,
    detail_count_row = detail_count_row,
    latest_batch_after = latest_batch_after
})

print("verified replay batch against lifecycle evidence:", replay_batch_id)
print("wrote artifacts/tutorials/83-verify-replay-batch-against-lifecycle-evidence-lua.json")
