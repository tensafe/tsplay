manifest_file = os.getenv("TSPLAY_MANIFEST_FILE") or "artifacts/tutorials/87-build-handoff-artifact-manifest-flow.csv"

manifest_rows = read_csv(manifest_file)
if #manifest_rows == 0 then
    error("manifest file is empty, run Lesson 87 first")
end

artifact_keys = {}
file_paths = {}
original_batch_id = ""
replay_batch_id = ""
for _, row in ipairs(manifest_rows) do
    table.insert(artifact_keys, row.artifact_key)
    table.insert(file_paths, row.file_path)
    if original_batch_id == "" then
        original_batch_id = tostring(row.batch_id or "")
    end
    replay_batch_id = tostring(row.related_batch_id or replay_batch_id)
end

write_json("artifacts/tutorials/88-build-handoff-summary-lua.json", {
    lesson = "88",
    mode = "lua",
    manifest_file = manifest_file,
    item_count = #manifest_rows,
    original_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    artifact_keys = artifact_keys,
    file_paths = file_paths,
    manifest_rows = manifest_rows
})

print("built handoff summary from manifest:", manifest_file)
print("wrote artifacts/tutorials/88-build-handoff-summary-lua.json")
