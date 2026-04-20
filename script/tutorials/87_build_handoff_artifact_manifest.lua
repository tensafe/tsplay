lifecycle_file = os.getenv("TSPLAY_LIFECYCLE_FILE") or "artifacts/tutorials/80-external-sync-lifecycle-round-trip-flow.csv"
replay_file = os.getenv("TSPLAY_REPLAY_FILE") or "artifacts/tutorials/82-replay-batch-from-lifecycle-evidence-flow.csv"
audit_compare_file = os.getenv("TSPLAY_AUDIT_COMPARE_FILE") or "artifacts/tutorials/85-export-original-and-replay-audits-flow.csv"
reconciliation_file = os.getenv("TSPLAY_RECONCILIATION_FILE") or "artifacts/tutorials/86-build-post-replay-reconciliation-pack-flow.csv"

lifecycle_rows = read_csv(lifecycle_file)
replay_rows = read_csv(replay_file)
audit_rows = read_csv(audit_compare_file)
reconciliation_rows = read_csv(reconciliation_file)
if #lifecycle_rows == 0 or #replay_rows == 0 or #reconciliation_rows == 0 then
    error("required CSV artifacts are missing, run Lessons 80, 82, 85, and 86 first")
end

lifecycle_row = lifecycle_rows[1]
replay_row = replay_rows[1]
reconciliation_row = reconciliation_rows[1]

manifest_rows = {
    {
        artifact_key = "lifecycle_evidence",
        file_path = lifecycle_file,
        batch_id = tostring(lifecycle_row.batch_id or ""),
        related_batch_id = tostring(lifecycle_row.batch_id or ""),
        row_count = tostring(lifecycle_row.pre_cleanup_detail or "0"),
        note = "original lifecycle evidence exported from lesson 80"
    },
    {
        artifact_key = "replay_result",
        file_path = replay_file,
        batch_id = tostring(replay_row.replay_batch_id or ""),
        related_batch_id = tostring(replay_row.source_batch_id or ""),
        row_count = tostring(replay_row.row_count or "0"),
        note = "replayed batch created from lifecycle evidence"
    },
    {
        artifact_key = "audit_comparison",
        file_path = audit_compare_file,
        batch_id = tostring(lifecycle_row.batch_id or ""),
        related_batch_id = tostring(replay_row.replay_batch_id or ""),
        row_count = tostring(#audit_rows),
        note = "combined original and replay audit history"
    },
    {
        artifact_key = "reconciliation_pack",
        file_path = reconciliation_file,
        batch_id = tostring(lifecycle_row.batch_id or ""),
        related_batch_id = tostring(replay_row.replay_batch_id or ""),
        row_count = tostring(reconciliation_row.db_row_count or "0"),
        note = "post-replay reconciliation summary"
    }
}

write_csv("artifacts/tutorials/87-build-handoff-artifact-manifest-lua.csv", manifest_rows, {"artifact_key", "file_path", "batch_id", "related_batch_id", "row_count", "note"})

write_json("artifacts/tutorials/87-build-handoff-artifact-manifest-lua.json", {
    lesson = "87",
    mode = "lua",
    lifecycle_file = lifecycle_file,
    replay_file = replay_file,
    audit_compare_file = audit_compare_file,
    reconciliation_file = reconciliation_file,
    manifest_rows = manifest_rows
})

print("built handoff artifact manifest for batch:", tostring(replay_row.replay_batch_id or ""))
print("wrote artifacts/tutorials/87-build-handoff-artifact-manifest-lua.json")
