role_file = os.getenv("TSPLAY_ROLE_FILE") or "artifacts/tutorials/91-read-handoff-manifest-roles-flow.csv"
checklist_file = os.getenv("TSPLAY_RUNTIME_CHECKLIST_FILE") or "artifacts/tutorials/89-build-pre-release-checklist-flow.csv"

role_rows = read_csv(role_file)
checklist_rows = read_csv(checklist_file)

if #role_rows == 0 then
    error("role file is empty, run Lesson 91 first")
end
if #checklist_rows == 0 then
    error("runtime checklist is empty, run Lesson 89 first")
end

template_rows = {}
for _, row in ipairs(role_rows) do
    local step_group = "collect"
    local focus_action = "read_csv"
    if row.phase == "verify" then
        step_group = "verify"
        focus_action = "read_csv"
    elseif row.phase == "deliver" then
        step_group = "save"
        focus_action = "write_csv,write_json"
    end

    table.insert(template_rows, {
        step_group = step_group,
        source_key = tostring(row.artifact_key or ""),
        source_file = tostring(row.file_path or ""),
        focus_action = focus_action,
        status_hint = "ready",
        note = tostring(row.note or ""),
    })
end

for _, row in ipairs(checklist_rows) do
    table.insert(template_rows, {
        step_group = "verify",
        source_key = tostring(row.artifact_key or ""),
        source_file = checklist_file,
        focus_action = "read_csv",
        status_hint = tostring(row.status or ""),
        note = tostring(row.note or ""),
    })
end

write_csv("artifacts/tutorials/94-build-collect-verify-save-template-lua.csv", template_rows, {
    "step_group",
    "source_key",
    "source_file",
    "focus_action",
    "status_hint",
    "note",
})

write_json("artifacts/tutorials/94-build-collect-verify-save-template-lua.json", {
    lesson = "94",
    mode = "lua",
    role_file = role_file,
    checklist_file = checklist_file,
    template_rows = template_rows,
})

print("built collect-verify-save template from:", role_file)
print("wrote artifacts/tutorials/94-build-collect-verify-save-template-lua.json")
