# TSPlay Flow Authoring Guide

## Purpose

Use this guide when Codex needs to write, repair, or review a TSPlay Flow and should still succeed even without the full repository nearby.

## Default Mindset

- Treat Flow as the primary delivery asset.
- Prefer readable YAML over clever one-off shortcuts.
- Start from user intent, then map it into stable steps, variables, and artifacts.
- Keep the result reviewable by another coder.
- Prefer standard Flow actions over ad hoc Lua when the task fits normal browser, file, control-flow, or extraction patterns.

## Runtime Assumptions

- Prefer a `tsplay` executable already available on PATH.
- If the task is happening inside the TSPlay repository, `./tsplay` or `go run .` is a valid fallback.
- Do not assume a fixed absolute executable path unless the user or environment explicitly provides one.

## Default Workflow

1. Identify the page or system target.
2. Clarify the business goal.
3. List runtime inputs and decide whether they belong in `vars`.
4. Decide what the Flow should output: `save_as`, JSON, CSV, Excel, screenshots, or other artifacts.
5. Pick the minimum required browser or security configuration.
6. Draft the smallest Flow that proves the task works.
7. Improve names, `save_as` fields, and artifact paths so the Flow is easy to review.

## Base Skeleton

```yaml
schema_version: "1"
name: example_flow
description: Explain the business result of this Flow in one sentence.
vars:
  page_url: http://127.0.0.1:8000/demo/example.html
steps:
  - action: navigate
    url: "{{page_url}}"

  - action: wait_for_selector
    selector: "#ready"
    timeout: 5000

  - action: extract_text
    selector: "#title"
    timeout: 5000
    save_as: page_title

  - action: set_var
    save_as: payload
    with:
      value:
        page_url: "{{page_url}}"
        page_title: "{{page_title}}"

  - action: write_json
    file_path: artifacts/example-flow.json
    value: "{{payload}}"
```

## Common Patterns

For action-level syntax, examples, and notes, read `actions.md`.

### Pattern 5: Report And Notify

Use when the Flow should send a completion notice, failure summary, or exported report by email.

- build the result payload first
- write artifacts such as JSON, CSV, or Excel when needed
- call `send_email`
- use `connection` for environment-driven SMTP or `with.smtp` for inline SMTP
- remember that restricted Flow or MCP contexts require `allow_email=true`

### Pattern 6: Read Local JSON And Continue

Use when a previous task already produced a JSON artifact, or when local configuration should drive later browser or data steps.

- call `read_json`
- save the decoded object into a business variable such as `payload` or `config`
- extract only the fields you need with variable expressions
- continue with assertions, loops, writes, or email notification

### Pattern 1: Page Assert And Extract

Use for readiness checks, smoke checks, and lightweight data capture.

- `navigate`
- `wait_for_selector`
- `assert_visible` or `assert_text`
- `extract_text` or table capture
- `set_var`
- `write_json`

### Pattern 2: Batch Input Processing

Use for CSV or Excel driven browser work.

- `read_csv` or `read_excel`
- `foreach`
- one row variable such as `row`
- `append_var` into a result ledger
- `write_json` and `write_csv`

### Pattern 3: Session Reuse

Use when login state should not be replayed every run.

```yaml
browser:
  use_session: admin
```

Keep session configuration at the top-level `browser` block, not scattered throughout steps.

### Pattern 4: Local Recovery

Use `on_error` inside `foreach` when one failing row should not kill the whole batch.

- append a success row in the happy path
- append a failed row in the `on_error` block
- preserve the error text with `{{last_error}}`

## Naming Rules

- `name` should say what the Flow does, not just that it is temporary.
- `description` should explain the delivery result, not restate the action list.
- `save_as` names should reflect business meaning such as `import_results`, `page_title`, or `auth_status`.
- Avoid vague names like `tmp`, `data`, `result1`, or `foo`.

## Artifact Rules

- Use stable, task-scoped output paths.
- For tutorial-style work, prefer `artifacts/tutorials/...`.
- Keep JSON, CSV, screenshots, and related outputs near each other when they belong to the same task.
- Do not hide key outputs in random file names that are hard to review later.

## Browser Rules

- Put reusable browser configuration in the top-level `browser` block.
- Use `browser.use_session` for saved login state.
- Use `browser.cdp_launch: true` when a trusted local run should use the real installed Chrome/Chromium/Edge with an isolated profile.
- Use `browser.cdp_port` or `browser.cdp_endpoint` only when the external browser was already started with `--remote-debugging-port`.
- Use `browser.timeout` for page-level timeout defaults.
- Prefer explicit waits or assertions before extracting or clicking unstable elements.
- CDP mode requires `allow_browser_state=true` in MCP. File outputs still require `allow_file_access=true`; page scripts still require `allow_javascript=true`.
- Do not combine CDP mode with `browser.use_session`, `storage_state/load_storage_state`, persistent `profile/session`, `user_agent`, or browser video recording.

## Repair Rules

When fixing a Flow:

1. Preserve the original business goal.
2. Keep good names and artifact paths unless they are part of the problem.
3. Fix selectors, waits, or variable wiring before rewriting the entire Flow.
4. Prefer the smallest repair that restores readability and execution.

## Review Checklist

- Can another coder understand the purpose from `name` and `description` alone
- Do `save_as` fields communicate business meaning
- Does the Flow use Flow-native structure instead of unnecessary Lua detours
- Are artifact paths stable and predictable
- Is the session or browser config placed at the top where it belongs
- Is the Flow short and direct enough for the task
