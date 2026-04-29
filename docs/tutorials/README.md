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

## Two Learning Tracks

### Track A: The Hands-On Path You Can Start Today

This is the fastest way to build intuition.
The path is intentionally staged so you do not hit MCP, repair, or delivery design too early.

1. Local foundations: [Lessons 01-12](01-hello-world.md)
   Learn the binary, bundled assets, local pages, extraction, assertions, and basic JSON output.
2. Files and control flow: [Lessons 13-27](13-read-csv-basics.md)
   Add CSV, Excel, upload/download, `retry`, `wait_until`, `if`, `foreach`, and `on_error`.
3. Browser state and sessions: [Lessons 28-57](28-inspect-storage-state.md)
   Learn storage state, screenshots, HTML, named sessions, protected pages, and authenticated import flows.
4. External systems and lifecycle: [Lessons 58-80](58-sync-import-report-summary-to-redis.md)
   Connect browser results to Redis and Postgres, then build reconciliation, audit, cleanup, and replayable lifecycles.
5. Handoff and templates: [Lessons 81-100](81-read-lifecycle-evidence.md)
   Turn lifecycle evidence into replay batches, handoff artifacts, reusable templates, and release checklists.
6. MCP workflow: [Lessons 101-120](101-assert-visible-template-release-card.md)
   Move from template-release robustness into `observe -> draft -> validate -> run -> repair -> finalize`.
7. Security and review thinking: [Lessons 121-140](121-allow-lua-boundary.md)
   Understand `allow_*`, `security_preset`, review rules, naming, artifact layout, and larger Flow packaging.
8. Delivery and curriculum operations: [Lessons 141-160](141-why-embed-docs-script-demo.md)
   Cover bundled assets, single-binary delivery, capstones, onboarding plans, trainer prep, and tutorial evolution loops.

### Track B: The Full Curriculum System

Use this when you want a structured learning system instead of a single tutorial chain.

1. [Curriculum Overview](curriculum-overview.md)
2. [Newbie Track](track-newbie.md)
3. [Junior Track](track-junior.md)
4. [Intermediate Track](track-intermediate.md)
5. [Advanced Track](track-advanced.md)
6. [160-Iteration Roadmap](iteration-roadmap-160.md)
7. [Evolution Playbook](evolution-playbook.md)

## Tutorial Map

| Phase | Lesson Range | Focus | Good First Stops |
| --- | --- | --- | --- |
| Foundations | `01-12` | bundled assets, local demo pages, extraction, assertions | [01](01-hello-world.md), [08](08-bundled-assets-and-artifacts.md), [10](10-assert-page-state.md) |
| Files and orchestration | `13-27` | CSV/Excel, upload/download, variables, retry, loops, recovery | [13](13-read-csv-basics.md), [16](16-retry-flaky-action.md), [22](22-foreach-batch-import-csv.md) |
| Browser state and reuse | `28-57` | storage state, screenshots, HTML, named sessions, protected imports | [28](28-inspect-storage-state.md), [36](36-save-storage-state.md), [42](42-use-named-session.md) |
| External system sync | `58-80` | Redis, Postgres, shared batch IDs, reconciliation, audit, cleanup | [58](58-sync-import-report-summary-to-redis.md), [61](61-db-insert-import-batch-summary.md), [71](71-external-system-round-trip.md) |
| Replay and handoff | `81-100` | evidence replay, handoff manifests, template catalogs, preflight checks | [81](81-read-lifecycle-evidence.md), [87](87-build-handoff-artifact-manifest.md), [96](96-build-template-index.md) |
| MCP chain | `101-120` | release-page robustness, MCP observe/draft/run/repair/finalize | [101](101-assert-visible-template-release-card.md), [111](111-mcp-list-actions.md), [120](120-mcp-finalize-flow.md) |
| Security and review | `121-140` | MCP permission boundaries, presets, Flow review, large-package structure | [121](121-allow-lua-boundary.md), [127](127-compare-local-flow-and-mcp-boundaries.md), [134](134-review-example-with-checklist.md) |
| Delivery and evolution | `141-160` | bundled delivery, offline learning, capstones, trainer and iteration systems | [141](141-why-embed-docs-script-demo.md), [144](144-single-binary-delivery-flow.md), [160](160-curriculum-continuation-plan.md) |

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
