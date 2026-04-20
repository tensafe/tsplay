template_index_file = os.getenv("TSPLAY_TEMPLATE_INDEX_FILE") or "artifacts/tutorials/96-build-template-index-flow.csv"
handoff_file = os.getenv("TSPLAY_HANDOFF_FILE") or "artifacts/tutorials/90-handoff-round-trip-from-lifecycle-evidence-flow.csv"

index_rows = read_csv(template_index_file)
handoff_rows = read_csv(handoff_file)

if #index_rows == 0 then
    error("template index is empty, run Lesson 96 first")
end
if #handoff_rows == 0 then
    error("handoff file is empty, run Lesson 90 first")
end

handoff_row = handoff_rows[1]
seen = {}
for _, row in ipairs(index_rows) do
    seen[tostring(row.template_name or "")] = true
end

verification_rows = {
    {
        check_name = "handoff_status_is_ok",
        status = tostring(handoff_row.status or "") == "ok" and "ready" or "blocked",
        detail = "handoff summary from lesson 90 must be ok",
    },
    {
        check_name = "has_input_process_output_template",
        status = seen["handoff_input_process_output"] and "ready" or "blocked",
        detail = "template index must contain the stage-based template",
    },
    {
        check_name = "has_collect_verify_save_template",
        status = seen["handoff_collect_verify_save"] and "ready" or "blocked",
        detail = "template index must contain the review-oriented template",
    },
    {
        check_name = "has_replay_audit_handoff_template",
        status = seen["handoff_replay_audit_handoff"] and "ready" or "blocked",
        detail = "template index must contain the delivery-oriented template",
    },
    {
        check_name = "template_count_is_three_or_more",
        status = #index_rows >= 3 and "ready" or "blocked",
        detail = "template index should keep the three core template families",
    },
}

write_csv("artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.csv", verification_rows, {
    "check_name",
    "status",
    "detail",
})

write_json("artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.json", {
    lesson = "97",
    mode = "lua",
    template_index_file = template_index_file,
    handoff_file = handoff_file,
    verification_rows = verification_rows,
})

print("verified template coverage for handoff chain")
print("wrote artifacts/tutorials/97-verify-template-covers-handoff-chain-lua.json")
