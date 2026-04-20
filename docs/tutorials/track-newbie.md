# Newbie Track

English | [简体中文](track-newbie.zh-CN.md)

> The job of this track is simple:
> help someone who has never used TSPlay accumulate a first, second, and third success in the right order.

This track is intentionally opinionated:

- no jumping ahead
- no racing toward advanced topics
- no abstract terminology before hands-on wins
- every new action should produce an immediate visible result

Most lesson pages linked below are still Chinese, but this page gives you the English learning path and exit criteria.

## Who This Is For

- people touching TSPlay on day one
- people who only know “it is an automation tool” and do not yet know where to start
- people who do not yet have stable intuitions for terms like `Lua`, `Flow`, `MCP`, or `artifact`

## Learning Goals

By the end of this track, you should be able to answer:

- how to build and run `tsplay`
- why a single binary can still carry `docs/`, `script/`, and `demo/`
- what the boundary is between `Lua` and `Flow`
- why local pages come before external systems
- why lesson outputs should land in `artifacts/`

## Recommended Order

1. [Lesson 01](01-hello-world.md)
2. [Lesson 08](08-bundled-assets-and-artifacts.md)
3. [Lesson 02](02-local-page-select-option.md)
4. [Lesson 09](09-local-demo-anatomy.md)
5. [Lesson 03](03-capture-table.md)
6. [Lesson 04](04-extract-text-and-html.md)
7. [Lesson 05](05-http-request-json.md)
8. [Lesson 10](10-assert-page-state.md)
9. [Lesson 11](11-select-another-option.md)
10. [Lesson 12](12-custom-json-output.md)

## Core Actions In This Track

- `set_var`
- `write_json`
- `navigate`
- `wait_for_selector`
- `select_option`
- `capture_table`
- `extract_text`
- `get_html`
- `http_request`
- `json_extract`
- `assert_visible`
- `assert_text`
- `is_selected`

## Deliverables

This layer does not require complex business workflows.
It only requires that every topic produces at least one runnable result:

- one `Lua` script
- one `Flow`
- one JSON output under `artifacts/tutorials/`
- one short reflection in your own words

## Common Mistakes

### Mistake 1: Starting With The Biggest Business Scenario

That usually means the page, system, permissions, and environment all break at once.
The newbie track should begin with the smallest local success path instead.

### Mistake 2: Learning Only Flow And Skipping Lua

Flow is the long-term main path, but for a beginner Lua often feels more direct.
First learn “how to do it”, then learn “how to express it structurally”.

### Mistake 3: Reading Without Running Commands

TSPlay is not a tool you truly learn by reading concepts alone.
At this stage, every small step should be executed.

## Exit Criteria

You are ready to move on when you can:

- build `./tsplay` independently
- start `./tsplay -action file-srv -addr :8000` independently
- run the newbie local practice chain on your own: `Lessons 01-05` and `08-12`
- explain the difference between `Lua` and `Flow`
- write page text, HTML fragments, tables, and JSON request results into files

## Where To Go Next

Next stop:
[track-junior.md](track-junior.md)
