# TSPlay Flow Example Index

## How To Use This File

Use this file when you want the fastest recommended starting Flow from the TSPlay repository. Pick the category closest to the user task, open the suggested Flow, and then adapt names, variables, selectors, and artifact paths.

If the repository is not available, fall back to `examples.md` and `flow-authoring.md`.

## Form Flows / 表单类 Flow

### Start Here: Simple search or submit

- `script/tutorials/10_assert_page_state.flow.yaml`
- Use when you need a compact browser flow with `navigate`, `wait_for_selector`, `assert_visible`, `assert_text`, `extract_text`, and `write_json`.
- Good first template for single-page form validation or smoke checks.

### Start Here: Optional login before form work

- `script/tutorials/21_if_optional_login.flow.yaml`
- Use when the page may or may not show a login step before the main form.
- Good first template for branching with `if`, then continuing the main flow.

### Start Here: Select or change one field value

- `script/tutorials/11_select_another_option.flow.yaml`
- Use when the job is mainly choosing one option and validating the selected value.
- Good first template for focused field interaction without extra control flow.

## Table Flows / 表格类 Flow

### Start Here: Capture a stable table

- `script/tutorials/03_capture_table.flow.yaml`
- Use when the page already exposes a usable table and later steps need structured rows.
- Good first template for `capture_table` plus simple structured extraction.

### Start Here: Compare page data with another artifact

- `script/tutorials/56_use_session_compare_table_and_download.flow.yaml`
- Use when you need table comparison plus a related download or post-check step.
- Good first template when the table is only part of a broader business flow.

## Import Flows / 导入类 Flow

### Start Here: CSV-driven batch import

- `script/tutorials/22_foreach_batch_import_csv.flow.yaml`
- Use when rows come from CSV and the page repeats the same submit pattern per row.
- Good first template for `read_csv`, `foreach`, `append_var`, and JSON summary output.

### Start Here: Excel-driven batch import

- `script/tutorials/50_use_session_batch_import_excel.flow.yaml`
- Use when rows come from Excel and the page is behind an authenticated session.
- Good first template for `read_excel`, `foreach`, `browser.use_session`, and JSON output.

### Start Here: Read one Excel slice cleanly

- `script/tutorials/25_read_excel_range_headers.flow.yaml`
- Use when the workbook contains a larger sheet but only one bounded range matters.
- Good first template for `read_excel` with `sheet`, `range`, and explicit header handling.

## Session Flows / 会话类 Flow

### Start Here: Reuse a saved login session

- `script/tutorials/38_verify_loaded_storage_state.flow.yaml`
- Use when you need a minimal proof that the session or storage state is loaded correctly.
- Good first template for top-level browser session configuration and verification.

### Start Here: Full session-backed import

- `script/tutorials/51_use_session_import_recovery_excel.flow.yaml`
- Use when the flow both reuses a session and needs per-row recovery.
- Good first template for session reuse plus `foreach` and `on_error`.

### Start Here: Session plus export or round-trip work

- `script/tutorials/57_use_session_import_export_round_trip.flow.yaml`
- Use when one authenticated flow includes import, verification, and exported result checks.
- Good first template for longer session-based business flows.

## Recovery Flows / 恢复与容错类 Flow

### Start Here: Per-row local recovery

- `script/tutorials/27_on_error_import_excel_writeback.flow.yaml`
- Use when one bad row should be recorded and skipped instead of stopping the whole batch.
- Good first template for `on_error`, `append_var`, JSON ledger, and CSV ledger output.

### Start Here: Failure artifacts for later repair

- `script/tutorials/35_error_evidence_pack.flow.yaml`
- Use when the real goal is to preserve screenshots, HTML, or error context for later debugging.
- Good first template when repairability matters as much as immediate success.

### Start Here: Retry flaky interactions

- `script/tutorials/103_retry_template_release_gate.flow.yaml`
- Use when the page is dynamic and short retries are safer than rewriting the flow.
- Good first template for `retry` plus business-level assertions.

## MCP Flows / MCP 驱动类 Flow

### Start Here: Learn the authoring contract

- `docs/tutorials/112-mcp-flow-schema-and-examples.md`
- Use when you want the official schema, action manifest, and example selection hints before drafting.
- Good first reference before building repo-specific MCP-facing flows.

### Start Here: Observe the page before drafting

- `docs/tutorials/113-mcp-observe-page.md`
- Use when the user knows the intent but not the selectors or DOM details.
- Good first reference for observation-driven flow drafting.

### Start Here: Draft and validate through MCP

- `docs/tutorials/114-mcp-draft-flow.md`
- `docs/tutorials/115-mcp-validate-drafted-flow.md`
- Use when the goal is to turn intent into Flow and then check structural issues before execution.
- Good first reference for the MCP draft and validation loop.

### Start Here: Finalize or repair through MCP

- `docs/tutorials/117-mcp-repair-flow-context.md`
- `docs/tutorials/118-mcp-repair-flow.md`
- `docs/tutorials/120-mcp-finalize-flow.md`
- Use when the draft is close but still needs repair context, repair steps, or the shorter finalize path.
- Good first reference for the `finalize -> run -> repair` workflow.

## Review And Readability / 可维护性与可评审性

### Start Here: Make the Flow easier to review

- `script/tutorials/131_review_readability_after.flow.yaml`
- Use when the Flow already works but names, descriptions, or artifact layout are too vague.
- Good first template for making a Flow shareable across a team.
