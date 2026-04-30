# TSPlay Step-by-Step Tutorials

English | [简体中文](README.zh-CN.md)

> This tutorial set is built to solve one onboarding problem:
> help people get TSPlay running quickly, then grow from single actions into maintainable delivery workflows.

This page is the English curriculum layer for the tutorial system.
Most individual lesson pages are still written in Chinese for now, so the lesson links below point to Chinese lesson files.
Use this page to understand the structure, pick the right lesson range, and enter the Chinese lesson set with a clear map.

## What This Tutorial Set Optimizes For

- start with runnable examples instead of abstract concepts
- show both `Lua` and `Flow` for the same capability whenever that helps learning
- let people see success first, then understand why day-to-day delivery usually favors `Flow`
- keep outputs observable through `artifacts/` so every lesson has something concrete to inspect

The repository already uses `docs/`, so the tutorial system lives here instead of a separate `doc/` folder.

## If You Only Want One Starting Point Today

Do not try to absorb the full curriculum first.
Pick the shortest path that matches your goal:

| Goal | Go here first | What counts as success today |
| --- | --- | --- |
| confirm TSPlay runs locally | [Lesson 01](01-hello-world.md) | run the first Flow and see output |
| learn page extraction first | [Lesson 03](03-capture-table.md) | write a local table into JSON |
| learn assertions first | [Lesson 10](10-assert-page-state.md) | verify page state against expectations |
| learn batch processing first | [Lesson 22](22-foreach-batch-import-csv.md) | finish one CSV-driven `foreach` run |
| learn session reuse first | [Lesson 36](36-save-storage-state.md) | save and reload browser state |
| jump to Agent / MCP | [Lesson 111](111-mcp-list-actions.md) | understand the MCP toolchain and next step |

## Before You Pick Lessons

If you have not brought the local environment up yet, start with [../../getting-started.md](../../getting-started.md) first.

The shortest first-run path is:

```bash
go mod download
go run . -flow script/tutorials/01_hello_world.flow.yaml
# only run this when you want local demo pages
go run . -action file-srv -addr :8000
```

The first two commands confirm TSPlay can run.
The third is only needed when you want to practice against local demo pages.

## Fast Entry Points

<div class="grid cards" markdown>

-   :material-run-fast:{ .lg .middle } __I Want A First Runnable Example__

    Start with the smallest possible lesson and confirm your local runtime first.

    [Lesson 01](01-hello-world.md)

-   :material-table-eye:{ .lg .middle } __I Want Data Extraction__

    Start with a local table and see the full page-to-JSON loop.

    [Lesson 03](03-capture-table.md)

-   :material-check-decagram-outline:{ .lg .middle } __I Want Assertions__

    Learn visibility and text checks early so later lessons feel easier.

    [Lesson 10](10-assert-page-state.md)

-   :material-file-sync-outline:{ .lg .middle } __I Want Batch Processing__

    Jump into CSV-driven `foreach` and local recovery patterns.

    [Lesson 22](22-foreach-batch-import-csv.md)

-   :material-account-lock-open-outline:{ .lg .middle } __I Want Sessions__

    Learn storage state, named sessions, and protected-page flows.

    [Lesson 36](36-save-storage-state.md)

-   :material-robot-outline:{ .lg .middle } __I Want MCP__

    Jump to the MCP toolchain and `finalize_flow` path.

    [Lesson 111](111-mcp-list-actions.md)

</div>

## Quick Topic Map

This section is a set of high-frequency entry points, not a strict learning sequence.
If you want the stable first-time path, start with the [Newbie Track](track-newbie.md).

### Foundations

- first runnable lesson: [01](01-hello-world.md), [02](02-local-page-select-option.md)
- extraction and local outputs: [03](03-capture-table.md), [04](04-extract-text-and-html.md), [12](12-custom-json-output.md)
- orchestration and resilience: [16](16-retry-flaky-action.md), [17](17-wait-until-ready.md), [21](21-if-optional-login.md), [23](23-on-error-import-recovery.md)
- files and batch import: [13](13-read-csv-basics.md), [22](22-foreach-batch-import-csv.md), [24](24-read-excel-basics.md), [27](27-on-error-import-excel-writeback.md)

### Sessions And Auth

- browser-state save and reuse: [36](36-save-storage-state.md), [42](42-use-named-session.md)
- login-driven import flow: [44](44-session-import-with-login.md)
- authenticated import/export round trip: [57](57-use-session-import-export-round-trip.md)

### MCP And Agent Workflow

- capability overview: [111](111-mcp-list-actions.md)
- page observation entry: [113](113-mcp-observe-page.md)
- `finalize_flow` convergence path: [120](120-mcp-finalize-flow.md)

## Recommended Learning Order

### 1. The 30-Minute Starter Path

This is the best path for day one.
It is intentionally short so you can accumulate a few wins before the curriculum widens.

1. [Lesson 01](01-hello-world.md)
2. [Lesson 02](02-local-page-select-option.md)
3. [Lesson 03](03-capture-table.md)
4. [Lesson 10](10-assert-page-state.md)
5. [Lesson 12](12-custom-json-output.md)

After these, you should already know:

- whether `tsplay` runs on your machine
- how local interaction, extraction, and assertions feel
- why outputs are written into `artifacts/`
- why the later Flow path is worth learning

### 2. The Half-Day Foundations Path

Use this when you want the first practical delivery layer, not just the first demo:

1. file input/output: [13](13-read-csv-basics.md), [14](14-write-csv-basics.md), [15](15-read-transform-write-csv.md)
2. reliability and control flow: [16](16-retry-flaky-action.md), [17](17-wait-until-ready.md), [21](21-if-optional-login.md), [23](23-on-error-import-recovery.md)
3. batch and Excel flows: [22](22-foreach-batch-import-csv.md), [24](24-read-excel-basics.md), [26](26-foreach-batch-import-excel.md), [27](27-on-error-import-excel-writeback.md)

### 3. Sessions And Protected Pages

Use this when you are moving toward real authenticated pages:

1. browser-state basics: [28](28-inspect-storage-state.md), [36](36-save-storage-state.md), [37](37-load-saved-storage-state.md)
2. named sessions: [40](40-save-named-session.md), [42](42-use-named-session.md)
3. authenticated round trip: [44](44-session-import-with-login.md), [47](47-use-session-import-single.md), [57](57-use-session-import-export-round-trip.md)

### 4. Agent / MCP Path

This path works best after you already have some local Flow intuition:

1. tool surface first: [111](111-mcp-list-actions.md), [112](112-mcp-flow-schema-and-examples.md)
2. shortest closed loop: [113](113-mcp-observe-page.md), [114](114-mcp-draft-flow.md), [115](115-mcp-validate-drafted-flow.md)
3. converge to delivery default: [119](119-mcp-chain-overview.md), [120](120-mcp-finalize-flow.md)

## Full Curriculum Map

If you are not a first-time learner, use this table instead of browsing lesson by lesson from the top:

| Phase | Lesson Range | Focus | Good First Stops |
| --- | --- | --- | --- |
| Foundations | `01-12` | local pages, extraction, assertions, JSON output | [01](01-hello-world.md), [03](03-capture-table.md), [10](10-assert-page-state.md) |
| Files and orchestration | `13-27` | CSV, Excel, upload/download, `retry`, `foreach`, `on_error` | [13](13-read-csv-basics.md), [16](16-retry-flaky-action.md), [22](22-foreach-batch-import-csv.md) |
| Sessions and auth | `28-57` | storage state, named sessions, protected pages | [28](28-inspect-storage-state.md), [36](36-save-storage-state.md), [42](42-use-named-session.md) |
| External systems | `58-80` | Redis, Postgres, reconciliation, audit, cleanup | [58](58-sync-import-report-summary-to-redis.md), [61](61-db-insert-import-batch-summary.md), [71](71-external-system-round-trip.md) |
| Handoff and templates | `81-100` | replay evidence, handoff packs, template indexes, preflight checks | [81](81-read-lifecycle-evidence.md), [87](87-build-handoff-artifact-manifest.md), [96](96-build-template-index.md) |
| MCP chain | `101-120` | release-page robustness, observe/draft/run/repair/finalize | [101](101-assert-visible-template-release-card.md), [111](111-mcp-list-actions.md), [120](120-mcp-finalize-flow.md) |
| Security and review | `121-140` | `allow_*`, `security_preset`, Flow review, larger package structure | [121](121-allow-lua-boundary.md), [127](127-compare-local-flow-and-mcp-boundaries.md), [134](134-review-example-with-checklist.md) |
| Delivery and evolution | `141-160` | single-binary delivery, offline learning, capstones, trainer operations | [141](141-why-embed-docs-script-demo.md), [144](144-single-binary-delivery-flow.md), [160](160-curriculum-continuation-plan.md) |

### 5. The Curriculum-Designer Path

Use this when you want a structured learning system instead of a single tutorial chain.

1. [Curriculum Overview](curriculum-overview.md)
2. [Newbie Track](track-newbie.md)
3. [Junior Track](track-junior.md)
4. [Intermediate Track](track-intermediate.md)
5. [Advanced Track](track-advanced.md)
6. [160-Iteration Roadmap](iteration-roadmap-160.md)
7. [Evolution Playbook](evolution-playbook.md)

## Full Curriculum System

| Layer | Audience | Outcome | Entry |
| --- | --- | --- | --- |
| Newbie | first-time TSPlay users | run the binary and understand the relationship between Lua and Flow | [track-newbie.md](track-newbie.md) |
| Junior | users who can already run demos | connect files, variables, control flow, HTTP, Redis, and DB basics | [track-junior.md](track-junior.md) |
| Intermediate | implementers / QA / developers who can already write basic Flows | build reusable templates, data-driven processes, robustness, and MCP basics | [track-intermediate.md](track-intermediate.md) |
| Advanced | people responsible for delivery, review, enablement, or integration | design standards, security boundaries, packaged delivery, repair strategy, and curriculum evolution | [track-advanced.md](track-advanced.md) |
| Long-term evolution | people growing the curriculum over time | turn the curriculum into a stable, iterative system | [iteration-roadmap-160.md](iteration-roadmap-160.md) |

## Before You Start

Run this once from the repository root:

```bash
go mod download
```

If you want to use the `tsplay` binary directly, build it once:

```bash
go build -o tsplay .
```

After that, these commands are ready to use:

```bash
./tsplay -script script/tutorials/01_hello_world.lua
./tsplay -flow script/tutorials/01_hello_world.flow.yaml
./tsplay -action cli
./tsplay -action file-srv -addr :8000
./tsplay -action mcp-tool -tool tsplay.list_actions
```

If you do not want to build first, the source-based equivalents also work:

```bash
go run . -script script/tutorials/01_hello_world.lua
go run . -flow script/tutorials/01_hello_world.flow.yaml
go run . -action mcp-tool -tool tsplay.list_actions
```

Additional notes:

- the built `./tsplay` binary already bundles `ReadMe.md`, `docs/`, `script/`, and `demo/`
- you can still run the bundled example paths even when the binary is copied elsewhere
- if you want to inspect the bundled reference assets on disk, use:

```bash
./tsplay -action list-assets
./tsplay -action extract-assets -extract-root ./tsplay-assets
```

For many browser-based lessons, you should keep the built-in static file server running in another terminal:

```bash
# Option A: run from source
go run . -action file-srv -addr :8000

# Option B: use the built binary
./tsplay -action file-srv -addr :8000
```

This is especially helpful for lessons `02-05`, `09-12`, `16-23`, `26-38`, `42`, `44-57`, and `101-120`.

## How To Use This Index

- If you are brand new to TSPlay, start with [track-newbie.md](track-newbie.md), then come back here when you need the wider map.
- If you already have automation experience, start with [track-junior.md](track-junior.md) and pick the lesson range closest to your current project.
- If you are designing training or enablement, read [curriculum-overview.md](curriculum-overview.md), then [track-advanced.md](track-advanced.md), and finally [evolution-playbook.md](evolution-playbook.md).
- If your main goal is AI-agent collaboration, the shortest useful jump is usually [111](111-mcp-list-actions.md) through [120](120-mcp-finalize-flow.md), after you are already comfortable with local Flow basics.

## Two Small Reminders

- Tutorial outputs usually go under `artifacts/tutorials/` so practice artifacts stay out of the repository source tree.
- This English layer currently maps the curriculum; most linked lesson pages are still Chinese and will be translated incrementally.
