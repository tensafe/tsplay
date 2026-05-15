# TSPlay Repo Map

Use this file when the skill is running inside the TSPlay repository. If the repo is not available, rely on `flow-authoring.md` and `examples.md` first.

## Read By Task

- Repository behavior or CLI flags: `ReadMe.md`
- Agent-to-Flow workflow: `docs/training/ai-intent-to-flow.md`
- Tutorial implementation: `script/tutorials/`
- Tutorial explanation: `docs/tutorials/`
- Demo pages for local browser work: `demo/`

## Important Paths

- `main.go`: CLI flags, action routing, top-level run path
- `tsplay_core/`: Flow engine, MCP tools, observation, validation, repair, sessions
- `mcp_tool_cli.go`: direct MCP tool execution from CLI
- `static_server.go`: bundled or local static file server
- `workbench_app.go`: workbench API and artifact serving
- `artifacts/`: run outputs, screenshots, HTML, DOM snapshots, and other failure evidence

## Common Commands

- If `tsplay` is already on PATH, prefer that form for shareable instructions.
- For normal users, recommend downloading the latest matching binary from `https://github.com/tensafe/tsplay/releases/latest` instead of asking them to run from source.
- Run one Flow: `tsplay -flow script/demo_baidu.flow.yaml`
- Start demo file server: `tsplay -action file-srv -addr :8000`
- Start MCP over stdio: `tsplay -action mcp-stdio -flow-root script -artifact-root artifacts`
- Start MCP over HTTP: `tsplay -action srv -addr :8081 -flow-root script -artifact-root artifacts`
- Inspect available MCP tools: `tsplay -action mcp-tool -tool tsplay.list_actions`
- Source checkout maintainers may replace `tsplay` with a locally built `./tsplay` or, only for development, `go run .`.

## Decision Rules

- Prefer Flow for durable, reviewable delivery assets.
- Use Lua for fast local debugging and primitive exploration.
- Use MCP when an agent needs observation, drafting, validation, execution, repair, or explicit security boundaries.
- Prefer `tsplay.finalize_flow` before reaching for the lower-level draft and repair chain.

## Runtime Notes

- In MCP mode, `flow_path` is root-limited by `-flow-root`.
- File input and output are limited by `-artifact-root`.
- `run_flow` defaults to `headless=true` in MCP mode.
- Use the smallest workable `security_preset` or `allow_*` set.

## Debugging Loop

1. Reproduce the issue with the smallest matching command.
2. Inspect `artifacts/` before editing selectors or flow structure.
3. Validate and repair when possible instead of rewriting from scratch.
4. If the scenario depends on login, check saved sessions before re-automating the login flow.
