template_catalog_file = os.getenv("TSPLAY_TEMPLATE_CATALOG_FILE") or "artifacts/tutorials/92-build-template-artifact-catalog-flow.csv"
template_index_file = os.getenv("TSPLAY_TEMPLATE_INDEX_FILE") or "artifacts/tutorials/96-build-template-index-flow.csv"
template_preflight_file = os.getenv("TSPLAY_TEMPLATE_PREFLIGHT_FILE") or "artifacts/tutorials/99-build-template-preflight-checklist-flow.csv"
handoff_file = os.getenv("TSPLAY_HANDOFF_FILE") or "artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv"

catalog_rows = read_csv(template_catalog_file)
index_rows = read_csv(template_index_file)
preflight_rows = read_csv(template_preflight_file)
handoff_rows = read_csv(handoff_file)

if #catalog_rows == 0 then
    error("template catalog is empty, run Lesson 92 first")
end
if #index_rows == 0 then
    error("template index is empty, run Lesson 96 first")
end
if #preflight_rows == 0 then
    error("template preflight checklist is empty, run Lesson 99 first")
end
if #handoff_rows == 0 then
    error("handoff file is empty, run Lesson 90 first")
end

handoff_row = handoff_rows[1]
ready_count = 0
blocked_count = 0
for _, row in ipairs(preflight_rows) do
    if tostring(row.status or "") == "ready" then
        ready_count = ready_count + 1
    else
        blocked_count = blocked_count + 1
    end
end

final_status = blocked_count == 0 and "ready" or "blocked"

summary_rows = {
    {
        original_batch_id = tostring(handoff_row.original_batch_id or ""),
        replay_batch_id = tostring(handoff_row.replay_batch_id or ""),
        template_count = tostring(#index_rows),
        catalog_slot_count = tostring(#catalog_rows),
        ready_check_count = tostring(ready_count),
        blocked_check_count = tostring(blocked_count),
        final_status = final_status,
    }
}

write_csv("artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-lua.csv", summary_rows, {
    "original_batch_id",
    "replay_batch_id",
    "template_count",
    "catalog_slot_count",
    "ready_check_count",
    "blocked_check_count",
    "final_status",
})

write_json("artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-lua.json", {
    lesson = "100",
    mode = "lua",
    handoff_file = handoff_file,
    template_catalog_file = template_catalog_file,
    template_index_file = template_index_file,
    template_preflight_file = template_preflight_file,
    summary_rows = summary_rows,
})

print("completed template round trip from handoff artifacts")
print("wrote artifacts/tutorials/100-template-round-trip-from-handoff-artifacts-lua.json")
