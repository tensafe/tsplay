# Advanced Track

English | [简体中文](track-advanced.zh-CN.md)

> The center of gravity at this layer is systematic delivery.
> The question is no longer just “Can this Flow run?” but “Can this capability survive review, packaging, handoff, and long-term growth?”

## Who This Is For

- people rolling TSPlay out inside a team
- people responsible for integration, review, delivery, training, or standards
- people treating `tsplay` as a product capability or delivery capability instead of a personal scripting tool

## Core Themes

### 1. Security Boundaries

You should understand:

- `allow_lua`
- `allow_http`
- `allow_file_access`
- `allow_redis`
- `allow_database`
- `allow_browser_state`

The goal is not to memorize flags.
It is to understand why these boundaries exist and what risk they are containing.

The current repository already supports a strong progression:

- `Lessons 121-126`: blocked vs allowed comparisons for each `allow_*`
- `Lessons 127-130`: local Flow vs MCP boundaries, `security_preset`, explicit overrides, and why boundary education comes early
- `Lessons 131-140`: naming, review, artifact layout, Lua escape hatches, and larger package structure

### 2. Large-Flow Organization

At this stage you should start defining standards for:

- step naming
- variable naming
- artifact directory layout
- environment variable conventions
- review rules

### 3. Binary Delivery and Delivery Packs

This is where you understand why the project bundles:

- `docs/`
- `script/`
- `demo/`
- `ReadMe.md`

into the binary itself.
The purpose is not minor convenience.
It is to ship the tutorials, examples, and runnable entry points as one deliverable.

The current progression is:

- `Lessons 141-143`: asset discovery and extraction
- `Lessons 144-146`: single-binary delivery and asset policy
- `Lessons 147-150`: `file-srv`, first-run entry, version strategy, and delivery retrospectives

### 4. Continuous Curriculum Evolution

The advanced layer should not stop at “I know how to use it”.
You should also be able to:

- add new themes
- refactor older documents
- arrange learning order by level
- maintain a growth path that does not jump unpredictably

The current progression is:

- `Lessons 151-154`: capstone briefs for four levels
- `Lessons 155-157`: onboarding plans for new hires, implementers, and trainers
- `Lessons 158-160`: gap reviews, every-10-iteration checkpoints, and the next expansion loop

## Deliverables

- one draft of team standards
- one tutorial evolution plan
- one deliverable example pack
- one recommended “where should a new person start” path
- one boundary map from action categories to minimum permission grants
- one review checklist for naming, artifacts, and Lua / Flow tradeoffs
- one release-pack explanation covering single-binary delivery, offline learning, and version strategy
- one continuous-evolution plan connecting capstones, training, retrospectives, and the next expansion loop

## Exit Criteria

You are ready to operate at this layer when you can:

- organize a deliverable tutorial system
- review structural quality instead of checking only whether a Flow runs
- explain TSPlay’s security boundaries and delivery boundaries clearly
- keep extending the curriculum using the same logic instead of starting from scratch each time

## Companion Documents

The advanced layer works best alongside:

- [iteration-roadmap-160.md](iteration-roadmap-160.md)
- [evolution-playbook.md](evolution-playbook.md)
