---
name: tsplay-flow-authoring
description: Shareable TSPlay Flow authoring skill for coders. Use when Codex needs to write, modify, fix, review, or explain TSPlay `.flow.yaml` files, choose between Flow DSL, Lua script mode, and MCP workflows, draft or review tutorial flows, or debug artifacts and saved sessions. Typical English triggers include write flow, fix flow, review flow, generate flow, repair selector, and convert a requirement into Flow. Typical Chinese triggers include 编写 Flow, 写 Flow, 改 Flow, 修 Flow, 生成 flow, 补一条 flow, 调试 flow, 排查 flow, 设计教程 flow, and 把需求转成 TSPlay Flow.
---

# TSPlay

## Overview

Use this skill when the task is about the TSPlay codebase or, most importantly, when a coder needs help producing a solid TSPlay Flow. Prefer Flow as the default delivery asset, use Lua for quick local exploration, and use MCP when the task is agent-facing or needs observation, drafting, validation, execution, or repair loops.

## Quick Start

- 如果用户主要用中文提需求，先读 [references/zh-cn.md](references/zh-cn.md)。
- 如果用户已经在中文排错，继续读 [references/zh-cn-troubleshooting.md](references/zh-cn-troubleshooting.md)。
- 如果用户想直接套中文业务模板，继续读 [references/zh-cn-business-templates.md](references/zh-cn-business-templates.md)。
- 如果用户主要卡在 selector，继续读 [references/zh-cn-selectors.md](references/zh-cn-selectors.md)。
- 如果用户是在做中文 review，继续读 [references/zh-cn-review-checklist.md](references/zh-cn-review-checklist.md)。
- Read [references/flow-authoring.md](references/flow-authoring.md) first for the self-contained Flow writing and repair guide.
- Read [references/actions.md](references/actions.md) when you need to choose the right TSPlay action or remember the exact YAML shape.
- Read [references/examples.md](references/examples.md) when the user request is still vague and you need a strong prompt shape for authoring or repairing a Flow.
- Read [references/example-index.md](references/example-index.md) when you want the fastest repo-backed starting point by category such as form, table, import, session, recovery, or MCP.
- Read [references/repo-map.md](references/repo-map.md) when you are inside the TSPlay repository and need exact paths or commands.
- If the local repo is available and the task is about product behavior, CLI flags, or repository entry points, read `ReadMe.md`.
- If the local repo is available and the task is about agent-driven Flow generation or repair, read `docs/training/ai-intent-to-flow.md`.
- If the local repo is available and the task is about a specific tutorial, inspect both `script/tutorials/<lesson>.flow.yaml` or `script/tutorials/<lesson>.lua` and the matching `docs/tutorials/<lesson>.md`.

## Decision Rules

- Default to Flow when the result should be reviewed, versioned, reused, or repaired later.
- Use Lua CLI or standalone Lua scripts only for one-off debugging, quick page exploration, or primitive-level experiments.
- Use MCP when another agent or product needs tool calling, page observation, security boundaries, or a `finalize -> run -> repair` loop.
- Prefer `tsplay.finalize_flow` as the default agent path. Split the workflow into `observe_page`, `draft_flow`, `validate_flow`, `run_flow`, `repair_flow_context`, and `repair_flow` only when the status, warnings, or unresolved items require tighter control.

## Runtime Assumptions

- Do not hardcode one absolute binary path into the skill.
- Prefer `tsplay` from the current PATH when it is already installed.
- If the user is inside the TSPlay repository, `./tsplay` or `go run .` is an acceptable fallback.
- If neither PATH nor repo-local entry points exist, ask the user where the TSPlay executable lives before assuming a custom path.

## Flow Authoring Priorities

- Start from the user goal, page URL, required inputs, desired output file or variables, and the minimum security boundary.
- Prefer stable and readable Flow structure over clever Lua detours.
- Use clear `name`, optional `description`, and meaningful `save_as` names so the Flow is reviewable.
- Reach for common Flow-native actions first: `navigate`, `wait_for_selector`, `type_text`, `click`, `assert_visible`, `assert_text`, `extract_text`, `set_var`, `foreach`, `on_error`, `write_json`, and `write_csv`.
- For Chinese requests, normalize the task into five fields first: 页面, 目标, 输入, 输出, 授权.
- Keep artifact paths stable and task-scoped, especially for tutorial or handoff work.
- Reuse existing tutorial patterns before inventing a brand-new structure from scratch.

## Operating Workflow

1. Clarify the page or system target, the user intent, required inputs, and the minimum authorization boundary.
2. Pick the smallest viable entry point: `tsplay -flow ...` when the binary is already on PATH, `./tsplay -flow ...` or `go run . -flow ...` inside the repo, `file-srv` for local demo pages, `mcp-stdio` for local agent integration, or `srv` for HTTP MCP.
3. Check validation results, warnings, unresolved selectors, and repair hints before treating a draft as complete.
4. If execution fails, inspect artifacts under `artifacts/` and prefer validation or repair tools over blindly rewriting selectors.
5. If login state matters, prefer saved sessions and put `browser.use_session` in top-level browser config instead of scattering login steps.
