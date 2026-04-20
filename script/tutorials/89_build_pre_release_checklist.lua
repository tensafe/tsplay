manifest_file = os.getenv("TSPLAY_MANIFEST_FILE") or "artifacts/tutorials/87-build-handoff-artifact-manifest-flow.csv"

manifest_rows = read_csv(manifest_file)
if #manifest_rows == 0 then
    error("manifest file is empty, run Lesson 87 first")
end

required_items = {
    {"lifecycle_evidence", "lifecycle evidence is available"},
    {"replay_result", "replay batch export is available"},
    {"audit_comparison", "audit comparison export is available"},
    {"reconciliation_pack", "reconciliation pack is available"},
}

seen = {}
for _, row in ipairs(manifest_rows) do
    seen[tostring(row.artifact_key or "")] = true
end

checklist_rows = {}
overall_status = "ready"
for _, item in ipairs(required_items) do
    local artifact_key = item[1]
    local note = item[2]
    local status = seen[artifact_key] and "ready" or "missing"
    if status ~= "ready" then
        overall_status = "blocked"
    end
    table.insert(checklist_rows, {
        artifact_key = artifact_key,
        status = status,
        note = note
    })
end

write_csv("artifacts/tutorials/89-build-pre-release-checklist-lua.csv", checklist_rows, {"artifact_key", "status", "note"})

write_json("artifacts/tutorials/89-build-pre-release-checklist-lua.json", {
    lesson = "89",
    mode = "lua",
    manifest_file = manifest_file,
    overall_status = overall_status,
    checklist_rows = checklist_rows
})

print("built pre-release checklist from manifest:", manifest_file)
print("wrote artifacts/tutorials/89-build-pre-release-checklist-lua.json")
