catalog_file = os.getenv("TSPLAY_TEMPLATE_CATALOG_FILE") or "artifacts/tutorials/92-build-template-artifact-catalog-flow.csv"
handoff_file = os.getenv("TSPLAY_HANDOFF_FILE") or "artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv"

catalog_rows = read_csv(catalog_file)
handoff_rows = read_csv(handoff_file)

if #catalog_rows == 0 then
    error("template catalog is empty, run Lesson 92 first")
end
if #handoff_rows == 0 then
    error("handoff file is empty, run Lesson 90 first")
end

stage_rows = {}
for _, row in ipairs(catalog_rows) do
    local stage = "process"
    if row.artifact_key == "lifecycle_evidence" then
        stage = "input"
    elseif row.artifact_key == "reconciliation_pack" then
        stage = "output"
    end

    table.insert(stage_rows, {
        stage = stage,
        template_slot = tostring(row.template_slot or ""),
        artifact_key = tostring(row.artifact_key or ""),
        required_env = tostring(row.required_env or ""),
        default_path = tostring(row.default_path or ""),
        recommended_actions = tostring(row.action_focus or ""),
        note = tostring(row.note or ""),
    })
end

table.insert(stage_rows, {
    stage = "output",
    template_slot = "output.handoff_summary_csv",
    artifact_key = "handoff_summary",
    required_env = "TSPLAY_HANDOFF_FILE",
    default_path = handoff_file,
    recommended_actions = "read_csv,write_json",
    note = "final handoff round trip summary",
})

write_csv("artifacts/tutorials/93-build-input-process-output-template-lua.csv", stage_rows, {
    "stage",
    "template_slot",
    "artifact_key",
    "required_env",
    "default_path",
    "recommended_actions",
    "note",
})

write_json("artifacts/tutorials/93-build-input-process-output-template-lua.json", {
    lesson = "93",
    mode = "lua",
    catalog_file = catalog_file,
    handoff_file = handoff_file,
    handoff_row = handoff_rows[1],
    stage_rows = stage_rows,
})

print("built input-process-output template from:", catalog_file)
print("wrote artifacts/tutorials/93-build-input-process-output-template-lua.json")
