template_catalog_file = os.getenv("TSPLAY_TEMPLATE_CATALOG_FILE") or "artifacts/tutorials/92-build-template-artifact-catalog-flow.csv"
handoff_file = os.getenv("TSPLAY_HANDOFF_FILE") or "artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv"

catalog_rows = read_csv(template_catalog_file)
handoff_rows = read_csv(handoff_file)

if #catalog_rows == 0 then
    error("template catalog is empty, run Lesson 92 first")
end
if #handoff_rows == 0 then
    error("handoff file is empty, run Lesson 90 first")
end

handoff_row = handoff_rows[1]
replay_batch_id = tostring(handoff_row.replay_batch_id or "")

template_rows = {}
for _, row in ipairs(catalog_rows) do
    local block_name = "handoff_release"
    local step_order = "4"
    if row.artifact_key == "lifecycle_evidence" then
        block_name = "replay_batch"
        step_order = "1"
    elseif row.artifact_key == "replay_result" then
        block_name = "replay_batch"
        step_order = "2"
    elseif row.artifact_key == "audit_comparison" then
        block_name = "replay_audit"
        step_order = "3"
    elseif row.artifact_key == "reconciliation_pack" then
        block_name = "handoff_release"
        step_order = "4"
    end

    table.insert(template_rows, {
        block_name = block_name,
        step_order = step_order,
        template_slot = tostring(row.template_slot or ""),
        action_focus = tostring(row.action_focus or ""),
        expected_artifact = tostring(row.artifact_key or ""),
        owner_batch_id = replay_batch_id,
        note = tostring(row.note or ""),
    })
end

table.insert(template_rows, {
    block_name = "handoff_release",
    step_order = "5",
    template_slot = "output.handoff_summary_csv",
    action_focus = "read_csv,write_json",
    expected_artifact = "handoff_summary",
    owner_batch_id = replay_batch_id,
    note = "final handoff summary for template release",
})

write_csv("artifacts/tutorials/95-build-replay-audit-handoff-template-lua.csv", template_rows, {
    "block_name",
    "step_order",
    "template_slot",
    "action_focus",
    "expected_artifact",
    "owner_batch_id",
    "note",
})

write_json("artifacts/tutorials/95-build-replay-audit-handoff-template-lua.json", {
    lesson = "95",
    mode = "lua",
    template_catalog_file = template_catalog_file,
    handoff_file = handoff_file,
    template_rows = template_rows,
})

print("built replay-audit-handoff template for batch:", replay_batch_id)
print("wrote artifacts/tutorials/95-build-replay-audit-handoff-template-lua.json")
