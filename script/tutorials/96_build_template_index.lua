input_process_output_file = os.getenv("TSPLAY_INPUT_PROCESS_OUTPUT_FILE") or "artifacts/tutorials/93-build-input-process-output-template-flow.csv"
collect_verify_save_file = os.getenv("TSPLAY_COLLECT_VERIFY_SAVE_FILE") or "artifacts/tutorials/94-build-collect-verify-save-template-flow.csv"
replay_audit_handoff_file = os.getenv("TSPLAY_REPLAY_AUDIT_HANDOFF_FILE") or "artifacts/tutorials/95-build-replay-audit-handoff-template-flow.csv"

input_process_output_rows = read_csv(input_process_output_file)
collect_verify_save_rows = read_csv(collect_verify_save_file)
replay_audit_handoff_rows = read_csv(replay_audit_handoff_file)

if #input_process_output_rows == 0 then
    error("input-process-output template is empty, run Lesson 93 first")
end
if #collect_verify_save_rows == 0 then
    error("collect-verify-save template is empty, run Lesson 94 first")
end
if #replay_audit_handoff_rows == 0 then
    error("replay-audit-handoff template is empty, run Lesson 95 first")
end

index_rows = {
    {
        template_name = "handoff_input_process_output",
        focus_area = "stage_design",
        step_row_count = tostring(#input_process_output_rows),
        source_file = input_process_output_file,
        recommended_when = "you need a stable input -> process -> output skeleton",
    },
    {
        template_name = "handoff_collect_verify_save",
        focus_area = "review_design",
        step_row_count = tostring(#collect_verify_save_rows),
        source_file = collect_verify_save_file,
        recommended_when = "you need to separate collection, verification, and output",
    },
    {
        template_name = "handoff_replay_audit_handoff",
        focus_area = "delivery_design",
        step_row_count = tostring(#replay_audit_handoff_rows),
        source_file = replay_audit_handoff_file,
        recommended_when = "you need replay, audit, and handoff blocks in one chain",
    },
}

write_csv("artifacts/tutorials/96-build-template-index-lua.csv", index_rows, {
    "template_name",
    "focus_area",
    "step_row_count",
    "source_file",
    "recommended_when",
})

write_json("artifacts/tutorials/96-build-template-index-lua.json", {
    lesson = "96",
    mode = "lua",
    index_rows = index_rows,
})

print("built template index from lessons 93-95")
print("wrote artifacts/tutorials/96-build-template-index-lua.json")
