CREATE TABLE IF NOT EXISTS public.tutorial_import_batches (
  batch_id TEXT PRIMARY KEY,
  report_file TEXT NOT NULL,
  source_lesson TEXT NOT NULL,
  row_count INTEGER NOT NULL,
  operator_name TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.tutorial_import_rows (
  batch_id TEXT NOT NULL REFERENCES public.tutorial_import_batches(batch_id) ON DELETE CASCADE,
  line_no INTEGER NOT NULL,
  name TEXT NOT NULL,
  phone TEXT NOT NULL,
  status TEXT NOT NULL,
  operator_name TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (batch_id, line_no)
);
