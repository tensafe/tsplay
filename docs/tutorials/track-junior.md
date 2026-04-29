# Junior Track

English | [简体中文](track-junior.zh-CN.md)

> The junior track is where isolated actions become a small workflow.
> If the newbie track answers “Can TSPlay run?”, this track answers “Can TSPlay start handling real work?”

Most lesson pages below are still Chinese, but the sequence and expectations here define the English path through them.

## Who This Is For

- people who can already run the local newbie practice chain reliably
- people who want to add files, variables, control flow, Redis, and databases
- people moving from “I can write an example” to “I can complete a simple task”

## Main Learning Spine

Move through this layer in this order:

1. file input and output
2. variable organization
3. control flow
4. HTTP / Redis / DB basics
5. session reuse and artifact management

Why this order:

- first you need inputs and outputs
- then you need state
- then you need flow control
- only then is it worth connecting external systems

## Current Entry Ranges

- `Lessons 13-15`: CSV input, transformation, and output
- `Lessons 16-17`: `retry` and `wait_until`
- `Lessons 18-20`: upload and download actions
- `Lessons 21-23`: `if`, `foreach`, and `on_error`
- `Lessons 24-27`: Excel-driven batch work
- `Lessons 28-39`: browser state inspection, snapshots, and storage-state round trips
- `Lessons 40-57`: named sessions and authenticated import flows
- `Lessons 58-64`: Redis / Postgres summaries and transactions
- `Lessons 65-71`: shared batch IDs and multi-system reconciliation
- `Lessons 72-80`: idempotent reruns, anomaly handling, audit, cleanup, and lifecycle closure

Good first stops:

- [Lesson 13](13-read-csv-basics.md)
- [Lesson 16](16-retry-flaky-action.md)
- [Lesson 22](22-foreach-batch-import-csv.md)
- [Lesson 28](28-inspect-storage-state.md)
- [Lesson 42](42-use-named-session.md)
- [Lesson 58](58-sync-import-report-summary-to-redis.md)
- [Lesson 61](61-db-insert-import-batch-summary.md)
- [Lesson 71](71-external-system-round-trip.md)

## What You Must Learn Here

### 1. Variable Design

You need more than `save_as`.
You need stable variable names that make a Flow reviewable, repairable, and reusable.

### 2. Workflow Boundaries

You should be able to state clearly:

- what the input is
- what the intermediate variables are
- what the output is
- where to inspect artifacts when a step fails

### 3. At Least One External System

You do not need all of them at once, but you should connect at least one of:

- HTTP
- Redis
- Postgres

## Deliverables

Each theme in this layer should produce at least one of these:

- a variable-driven Flow
- a Flow with control flow
- a minimal Redis / DB / HTTP example
- an input/output explanation
- a failure-artifact explanation

## Exit Criteria

You are ready for the next layer when you can:

- write a small Flow with roughly 5 to 10 steps
- use `save_as`, `set_var`, `assert_*`, `read_csv`, `write_csv`, `read_excel`, `get_storage_state`, `get_cookies_string`, `save_storage_state`, and `load_storage_state` independently
- use `use_session` to connect saved browser state to a real authenticated batch-import flow
- preserve both “page table output” and “exported file output” for review
- continue an exported CSV into Redis and Postgres and explain the progression from browser result to external-system summary / persistence
- explain why a shared batch ID matters across Redis, Postgres summaries, Postgres detail rows, and local reconciliation artifacts
- explain why runtime data and audit records should be separate, then complete one minimal cycle of rerun, audit, cleanup, and verification
- explain why one step belongs in `Lua` and another in `Flow`
- explain where to inspect failures and when screenshots / HTML / JSON evidence should be preserved
- explain the difference between writing a storage-state path directly and using a named session

## Where To Go Next

Next stop:
[track-intermediate.md](track-intermediate.md)
