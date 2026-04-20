# TSPlay Curriculum Overview

English | [简体中文](curriculum-overview.zh-CN.md)

> This overview is not about adding a few more tutorial pages.
> It turns TSPlay tutorials into a learning system that can keep evolving over time.

This English page gives you the curriculum model and learning order.
Most lesson-level pages are still Chinese, but the track entry documents and lesson numbers below tell you where to go next.

The curriculum has three goals:

- help first-time TSPlay users learn in order instead of jumping straight into MCP, repair, or database transactions
- help people who can already run demos understand what capability gap to close next
- give course authors, implementers, and trainers a stable spine they can extend instead of adding lessons ad hoc

## Core Principles

### 1. Local First, External Systems Second

Start with local pages, local JSON, and bundled binary assets to build certainty.
Only after a layer is stable do you move into Redis, databases, MCP, session reuse, and repair.

### 2. Runnable First, Abstraction Second

Learn how to make something run before you learn how to abstract it into `Flow`.
That is why many topics are easiest to understand as `Lua` first, then `Flow`.

### 3. Single Action First, Composition Second

Understand one action clearly before chaining several together.
For example: learn `extract_text`, then `extract_text + set_var + if`, and only later combine it with `retry` or `on_error`.

### 4. Stable Pages First, Unstable Pages Later

Use the repository’s bundled `demo/` pages and local JSON to build fundamentals.
That keeps the tutorial focused on TSPlay instead of getting derailed by changing external pages.

### 5. Deliverables First, “I Get It” Second

Every lesson should leave behind something inspectable:

- a script
- a Flow
- a JSON / CSV / screenshot / trace artifact
- a short reflection or review note

That is what makes the learning path assessable, reproducible, and improvable.

## Four Curriculum Layers

### Layer 1: Newbie Track

Entry:
[track-newbie.md](track-newbie.md)

Best for:

- people touching TSPlay for the first time
- users who do not yet have stable intuitions for Lua or Flow
- learners who need a string of small wins before tackling broader systems

Problems this layer solves:

- how to run the `tsplay` binary
- how bundled assets work
- how to start the local static file server
- how Lua and Flow relate to each other
- how to capture page text, HTML, tables, and JSON output

Completion signals:

- you can run the current newbie chain independently: `Lessons 01-05` and `08-12`
- you can explain why “Lua first, Flow second” is often easier for onboarding
- you can save outputs into `artifacts/tutorials/`

### Layer 2: Junior Track

Entry:
[track-junior.md](track-junior.md)

Best for:

- people who can already run the local demo flows
- people who want to add files, variables, control flow, and basic external-system integration
- people moving from “I can run it” to “I can complete a small business workflow”

Problems this layer solves:

- file input and output
- variables, branching, looping, and assertions
- basic HTTP, Redis, and database actions
- session handling, artifact management, and output layout

Completion signals:

- you can write a basic Flow with variables and control flow
- you can connect a minimal Redis or Postgres example
- you can explain a workflow’s inputs, outputs, and failure artifacts

### Layer 3: Intermediate Track

Entry:
[track-intermediate.md](track-intermediate.md)

Best for:

- people who can already write basic Flows independently
- people thinking about reuse, robustness, and maintainability
- people turning one-off automation into team assets

Problems this layer solves:

- how to turn temporary scripts into templates
- how to use CSV / Excel / `foreach` for data-driven processes
- how to design robustness with `retry / on_error / wait_until`
- how to understand `observe -> draft -> validate -> run -> repair`

Completion signals:

- you can maintain a small set of reusable Flow templates
- you can design a minimal fix strategy for flaky workflows
- you can explain when MCP is worth introducing instead of writing more scripts by hand

### Layer 4: Advanced Track

Entry:
[track-advanced.md](track-advanced.md)

Best for:

- people delivering TSPlay capability into teams, projects, or client environments
- people designing standards, reviews, training, integrations, and release workflows
- people who want TSPlay to operate as a system capability rather than a private scripting tool

Problems this layer solves:

- how to design security boundaries and permission grants
- how to layer, name, review, and version larger Flows
- how to organize bundled assets, binary delivery, and delivery packs
- how to keep the curriculum evolving instead of freezing at version one

Completion signals:

- you can define team tutorial structure and review rules
- you can explain where TSPlay fits into the delivery chain
- you can keep extending the curriculum with the same logic instead of improvising each time

## What Is Already Runnable Today

The current repository already contains a full first-pass lesson system from `Lesson 01` through `Lesson 160`.
The easiest way to reason about it in English is by capability bands:

- `01-12`: local fundamentals, bundled assets, demo pages, extraction, and assertions
- `13-27`: files, control flow, uploads/downloads, and batch inputs
- `28-57`: browser state, storage files, named sessions, and authenticated import flows
- `58-80`: Redis / Postgres integration, reconciliation, audit, cleanup, and external-sync lifecycles
- `81-100`: replay from lifecycle evidence, handoff artifacts, template systems, and preflight checks
- `101-120`: MCP chain and template-release robustness
- `121-140`: permission boundaries, `security_preset`, review quality, artifact layout, and larger package structure
- `141-160`: bundled delivery, offline onboarding, capstones, trainer prep, and curriculum iteration loops

The full expansion logic for this system continues in:

- [160-Iteration Roadmap](iteration-roadmap-160.md)
- [Evolution Playbook](evolution-playbook.md)

## Recommended Usage

If this is your first time with TSPlay:

1. Start with [README.md](README.md) for the tutorial index.
2. Continue with [track-newbie.md](track-newbie.md).
3. After that, use the newbie stage inside [iteration-roadmap-160.md](iteration-roadmap-160.md).

If you already have some automation experience:

1. Scan [track-junior.md](track-junior.md).
2. Pick the lesson block closest to your real project.
3. After each block, return to the roadmap and continue the next iteration.

If you are a course author or delivery lead:

1. Read this overview first.
2. Continue with [track-advanced.md](track-advanced.md).
3. Use [evolution-playbook.md](evolution-playbook.md) to control how the documentation keeps growing.
