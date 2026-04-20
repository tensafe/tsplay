template_index_file = os.getenv("TSPLAY_TEMPLATE_INDEX_FILE") or "artifacts/tutorials/96-build-template-index-flow.csv"
verification_file = os.getenv("TSPLAY_TEMPLATE_VERIFICATION_FILE") or "artifacts/tutorials/97-verify-template-covers-handoff-chain-flow.csv"
lesson_matrix_file = os.getenv("TSPLAY_TEMPLATE_LESSON_MATRIX_FILE") or "artifacts/tutorials/98-build-template-lesson-matrix-flow.csv"

index_rows = read_csv(template_index_file)
verification_rows = read_csv(verification_file)
matrix_rows = read_csv(lesson_matrix_file)

if #index_rows == 0 then
    error("template index is empty, run Lesson 96 first")
end
if #verification_rows == 0 then
    error("template verification is empty, run Lesson 97 first")
end
if #matrix_rows == 0 then
    error("template lesson matrix is empty, run Lesson 98 first")
end

verification_ready = true
for _, row in ipairs(verification_rows) do
    if tostring(row.status or "") ~= "ready" then
        verification_ready = false
    end
end

checklist_rows = {
    {
        check_name = "template_index_present",
        status = #index_rows >= 3 and "ready" or "blocked",
        note = "template index should expose the three reusable template families",
    },
    {
        check_name = "template_verification_passed",
        status = verification_ready and "ready" or "blocked",
        note = "template verification from lesson 97 should be fully ready",
    },
    {
        check_name = "template_lesson_matrix_present",
        status = #matrix_rows >= 3 and "ready" or "blocked",
        note = "template lesson matrix should map scenarios to templates",
    },
}

write_csv("artifacts/tutorials/99-build-template-preflight-checklist-lua.csv", checklist_rows, {
    "check_name",
    "status",
    "note",
})

write_json("artifacts/tutorials/99-build-template-preflight-checklist-lua.json", {
    lesson = "99",
    mode = "lua",
    checklist_rows = checklist_rows,
})

print("built template preflight checklist")
print("wrote artifacts/tutorials/99-build-template-preflight-checklist-lua.json")
