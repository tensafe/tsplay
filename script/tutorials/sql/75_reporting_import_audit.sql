CREATE TABLE IF NOT EXISTS public.tutorial_import_audits (
  audit_id TEXT PRIMARY KEY,
  batch_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  status TEXT NOT NULL,
  detail TEXT NOT NULL,
  success_count INTEGER NOT NULL,
  failure_count INTEGER NOT NULL,
  source_lesson TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tutorial_import_audits_batch_id
  ON public.tutorial_import_audits(batch_id);
