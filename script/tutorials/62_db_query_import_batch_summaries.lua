input_file = os.getenv("TSPLAY_IMPORTED_REPORT") or "artifacts/tutorials/57-use-session-import-export-round-trip-flow.csv"
batch_prefix = os.getenv("TSPLAY_IMPORT_BATCH_PREFIX") or "lesson-62-import-batch-"

rows = read_csv(input_file)
row_count = #rows
operator_name = "unknown"

if row_count > 0 then
    operator_name = rows[row_count].operator or "unknown"
end

cleanup = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_import_batches WHERE batch_id LIKE $1",
    args = {batch_prefix .. "%"}
})

insert_one = db_insert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_prefix .. "a",
        report_file = input_file .. "#a",
        source_lesson = "62",
        row_count = row_count,
        operator_name = operator_name
    }
})

insert_two = db_insert({
    connection = "reporting",
    table = "public.tutorial_import_batches",
    columns = {"batch_id", "report_file", "source_lesson", "row_count", "operator_name"},
    row = {
        batch_id = batch_prefix .. "b",
        report_file = input_file .. "#b",
        source_lesson = "62",
        row_count = row_count,
        operator_name = operator_name
    }
})

batch_rows = db_query({
    connection = "reporting",
    sql = "SELECT batch_id, report_file, row_count, operator_name FROM public.tutorial_import_batches WHERE batch_id LIKE $1 ORDER BY batch_id",
    args = {batch_prefix .. "%"}
})

write_json("artifacts/tutorials/62-db-query-import-batch-summaries-lua.json", {
    lesson = "62",
    mode = "lua",
    input_file = input_file,
    batch_prefix = batch_prefix,
    cleanup = cleanup,
    insert_one = insert_one,
    insert_two = insert_two,
    batch_rows = batch_rows
})

print("queried import batch summaries:", tostring(#batch_rows))
print("wrote artifacts/tutorials/62-db-query-import-batch-summaries-lua.json")
