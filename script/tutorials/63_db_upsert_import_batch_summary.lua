input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
batch_id = os.getenv("TSPLAY_IMPORT_BATCH_ID") or "lesson-63-import-batch"

rows = read_csv(input_file)
row_count = #rows
operator_name = "unknown"

if row_count > 0 then
    operator_name = rows[row_count].operator or "unknown"
end

cleanup = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

seed_insert = db_insert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_id,
        report_file = "seed.csv",
        source_lesson = "63-seed",
        row_count = 1,
        operator_name = "seed-user"
    }
})

upsert_result = db_upsert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    key_columns = {"batch_id"},
    update_columns = {"report_file", "source_lesson", "row_count", "operator_name"},
    returning = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_id,
        report_file = input_file,
        source_lesson = "57",
        row_count = row_count,
        operator_name = operator_name
    }
})

row = db_query_one({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, source_lesson, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id = $1",
    args = {batch_id}
})

write_json("artifacts/tutorials/63-db-upsert-import-batch-summary-lua.json", {
    lesson = "63",
    mode = "lua",
    input_file = input_file,
    batch_id = batch_id,
    cleanup = cleanup,
    seed_insert = seed_insert,
    upsert_result = upsert_result,
    row = row
})

print("upserted import batch summary:", tostring(batch_id))
print("wrote artifacts/tutorials/63-db-upsert-import-batch-summary-lua.json")
