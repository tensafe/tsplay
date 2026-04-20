role_file = os.getenv("TSPLAY_ROLE_FILE") or "artifacts/tutorials/91-read-handoff-manifest-roles-flow.csv"

role_rows = read_csv(role_file)
if #role_rows == 0 then
    error("role file is empty, run Lesson 91 first")
end

catalog_rows = {}
for _, row in ipairs(role_rows) do
    local template_slot = "other.review_file"
    local required_env = "TSPLAY_ARTIFACT_FILE"
    local action_focus = "read_csv"
    local batch_scope = "replay"
    if row.artifact_key == "lifecycle_evidence" then
        template_slot = "input.lifecycle_csv"
        required_env = "TSPLAY_LIFECYCLE_FILE"
        action_focus = "read_csv"
        batch_scope = "original"
    elseif row.artifact_key == "replay_result" then
        template_slot = "process.replay_result_csv"
        required_env = "TSPLAY_REPLAY_FILE"
        action_focus = "read_csv,write_csv"
        batch_scope = "replay"
    elseif row.artifact_key == "audit_comparison" then
        template_slot = "verify.audit_compare_csv"
        required_env = "TSPLAY_AUDIT_COMPARE_FILE"
        action_focus = "read_csv,write_csv"
        batch_scope = "replay"
    elseif row.artifact_key == "reconciliation_pack" then
        template_slot = "deliver.reconciliation_csv"
        required_env = "TSPLAY_RECONCILIATION_FILE"
        action_focus = "read_csv,write_csv,write_json"
        batch_scope = "replay"
    end

    table.insert(catalog_rows, {
        template_slot = template_slot,
        artifact_key = tostring(row.artifact_key or ""),
        phase = tostring(row.phase or ""),
        required_env = required_env,
        action_focus = action_focus,
        default_path = tostring(row.file_path or ""),
        batch_scope = batch_scope,
        note = tostring(row.note or ""),
    })
end

write_csv("artifacts/tutorials/92-build-template-artifact-catalog-lua.csv", catalog_rows, {
    "template_slot",
    "artifact_key",
    "phase",
    "required_env",
    "action_focus",
    "default_path",
    "batch_scope",
    "note",
})

write_json("artifacts/tutorials/92-build-template-artifact-catalog-lua.json", {
    lesson = "92",
    mode = "lua",
    role_file = role_file,
    catalog_rows = catalog_rows,
})

print("built template artifact catalog from:", role_file)
print("wrote artifacts/tutorials/92-build-template-artifact-catalog-lua.json")
