manifest_file = os.getenv("TSPLAY_MANIFEST_FILE") or "artifacts/tutorials/87-build-handoff-artifact-manifest-flow.csv"
handoff_file = os.getenv("TSPLAY_HANDOFF_FILE") or "artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv"

manifest_rows = read_csv(manifest_file)
handoff_rows = read_csv(handoff_file)

if #manifest_rows == 0 then
    error("manifest file is empty, run Lesson 87 first")
end
if #handoff_rows == 0 then
    error("handoff file is empty, run Lesson 90 first")
end

handoff_row = handoff_rows[1]
original_batch_id = tostring(handoff_row.original_batch_id or "")
replay_batch_id = tostring(handoff_row.replay_batch_id or "")

role_rows = {}
for _, row in ipairs(manifest_rows) do
    local phase = "other"
    local owner_batch_id = replay_batch_id
    local template_goal = "review"
    if row.artifact_key == "lifecycle_evidence" then
        phase = "input"
        owner_batch_id = original_batch_id
        template_goal = "rebuild"
    elseif row.artifact_key == "replay_result" then
        phase = "process"
        owner_batch_id = replay_batch_id
        template_goal = "replay"
    elseif row.artifact_key == "audit_comparison" then
        phase = "verify"
        owner_batch_id = replay_batch_id
        template_goal = "compare"
    elseif row.artifact_key == "reconciliation_pack" then
        phase = "deliver"
        owner_batch_id = replay_batch_id
        template_goal = "handoff"
    end

    table.insert(role_rows, {
        phase = phase,
        artifact_key = tostring(row.artifact_key or ""),
        file_path = tostring(row.file_path or ""),
        owner_batch_id = owner_batch_id,
        related_batch_id = tostring(row.related_batch_id or ""),
        row_count = tostring(row.row_count or "0"),
        template_goal = template_goal,
        note = tostring(row.note or ""),
    })
end

write_csv("artifacts/tutorials/91-read-handoff-manifest-roles-lua.csv", role_rows, {
    "phase",
    "artifact_key",
    "file_path",
    "owner_batch_id",
    "related_batch_id",
    "row_count",
    "template_goal",
    "note",
})

write_json("artifacts/tutorials/91-read-handoff-manifest-roles-lua.json", {
    lesson = "91",
    mode = "lua",
    manifest_file = manifest_file,
    handoff_file = handoff_file,
    original_batch_id = original_batch_id,
    replay_batch_id = replay_batch_id,
    role_rows = role_rows,
})

print("read handoff manifest roles for replay batch:", replay_batch_id)
print("wrote artifacts/tutorials/91-read-handoff-manifest-roles-lua.json")
