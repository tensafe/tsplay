template_index_file = os.getenv("TSPLAY_TEMPLATE_INDEX_FILE") or "artifacts/tutorials/96-build-template-index-flow.csv"
template_catalog_file = os.getenv("TSPLAY_TEMPLATE_CATALOG_FILE") or "artifacts/tutorials/92-build-template-artifact-catalog-flow.csv"

index_rows = read_csv(template_index_file)
catalog_rows = read_csv(template_catalog_file)

if #index_rows == 0 then
    error("template index is empty, run Lesson 96 first")
end
if #catalog_rows == 0 then
    error("template catalog is empty, run Lesson 92 first")
end

artifact_by_phase = {}
for _, row in ipairs(catalog_rows) do
    artifact_by_phase[tostring(row.phase or "")] = tostring(row.artifact_key or "")
end

matrix_rows = {
    {
        scenario = "restore_from_lifecycle",
        template_name = "handoff_input_process_output",
        starting_artifact = artifact_by_phase["input"] or "lifecycle_evidence",
        expected_output = "reconciliation_pack",
        lesson_range = "91-93",
    },
    {
        scenario = "review_handoff_quality",
        template_name = "handoff_collect_verify_save",
        starting_artifact = artifact_by_phase["verify"] or "audit_comparison",
        expected_output = "template_review_summary",
        lesson_range = "94,97",
    },
    {
        scenario = "rebuild_delivery_chain",
        template_name = "handoff_replay_audit_handoff",
        starting_artifact = artifact_by_phase["process"] or "replay_result",
        expected_output = "handoff_summary",
        lesson_range = "95-96",
    },
}

write_csv("artifacts/tutorials/98-build-template-lesson-matrix-lua.csv", matrix_rows, {
    "scenario",
    "template_name",
    "starting_artifact",
    "expected_output",
    "lesson_range",
})

write_json("artifacts/tutorials/98-build-template-lesson-matrix-lua.json", {
    lesson = "98",
    mode = "lua",
    matrix_rows = matrix_rows,
})

print("built template lesson matrix from index and catalog")
print("wrote artifacts/tutorials/98-build-template-lesson-matrix-lua.json")
