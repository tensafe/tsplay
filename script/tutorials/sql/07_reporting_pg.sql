CREATE TABLE IF NOT EXISTS public.tutorial_orders (
  order_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  amount NUMERIC(10, 2) NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
