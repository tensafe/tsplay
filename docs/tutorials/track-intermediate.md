# Intermediate Track

English | [简体中文](track-intermediate.zh-CN.md)

> The key question at this layer is no longer “Can I add more actions?”
> It becomes “Can I make this automation maintainable?”

From here on, the real questions are:

- can the Flow be reused
- will we still recognize it next month
- can another teammate take it over
- how painful will repair be after the page changes

## Who This Is For

- people who can already write basic Flows independently
- people hitting flaky pages, expanding variable sets, and longer workflows
- people who want TSPlay to become a maintainable team asset

## Core Themes

### 1. Template Thinking

Turn one-off scripts into:

- reusable scripts
- reusable Flow templates
- reusable environment examples

### 2. Data-Driven Workflows

Move from handling one record to handling batches:

- CSV
- Excel
- `foreach`
- writing batch results back out

Today’s repository already lets you start this through:

- `Lessons 13-15` for CSV
- `Lessons 24-27` for Excel

### 3. Robustness Design

You should become systematic with:

- `retry`
- `wait_until`
- `on_error`
- `assert_visible`
- `assert_text`

The shortest progression looks like this:

- `Lessons 16-20`: retries, waits, uploads/downloads, and assertion chains
- `Lessons 21-27`: `if`, `foreach`, `on_error`, and import flows
- `Lessons 28-30` and `36-43`: browser state, state files, and named sessions
- `Lessons 44-57`: authenticated import flows with reusable sessions
- `Lessons 58-64`: Redis / Postgres synchronization basics
- `Lessons 65-71`: shared batch IDs, detail persistence, and three-way reconciliation
- `Lessons 72-80`: reruns, anomaly recovery, audit, and cleanup
- `Lessons 81-100`: lifecycle evidence replay, handoff packs, templates, and preflight checks
- `Lessons 101-110`: template-release robustness, assertions, waits, retries, reloads, and evidence capture

### 4. MCP Foundation Chain

This is where you learn:

- `observe_page`
- `draft_flow`
- `validate_flow`
- `run_flow`
- `repair_flow_context`
- `repair_flow`
- `finalize_flow`

The core MCP run-up is `Lessons 111-120`.

## Deliverables

- one data-driven Flow
- one Flow with explicit failure recovery
- one explanation of why the workflow was split into these steps
- one repair-before / repair-after comparison
- one MCP artifact set: `observation / draft / validate / run / repair / finalize`

## What Evaluation Should Focus On

This stage is not about “Do you know more actions?”
It is about:

- whether you can decompose a process clearly
- whether variable names stay stable
- whether failures remain observable
- whether the Flow leaves room for repair
- whether you can explain MCP input/output relationships instead of calling tools blindly

## Exit Criteria

You are ready for the advanced layer when you can:

- maintain a small reusable template set
- handle batches instead of single records
- design reasonable `retry / on_error / wait_until` strategies for flaky flows
- explain when MCP is worth introducing
- complete one end-to-end `observe -> draft -> validate -> run -> repair -> finalize` cycle independently

## Where To Go Next

Next stop:
[track-advanced.md](track-advanced.md)
