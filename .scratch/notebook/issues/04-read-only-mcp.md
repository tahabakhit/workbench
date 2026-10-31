# 04 Read-only notebook mcp

Status: done
Blocked by: 02

## What

Add `notebook mcp` with `mcp-go` v1.5.0: a stdio MCP server with tools
`search`, `recall`, `links` and `backlinks`, calling the same code as the CLI.
Add no write tool. Do not add `.mcp.json` to the plugin yet; issue 05 does.

## Done when

- `sh scripts/check.sh` passes, with a test that starts the server over stdio,
  lists exactly those four tools, and calls `search`.
- `claude --mcp-config <temporary file pointing at a local build> -p` answers a
  question using the `search` tool, and the output is quoted in this issue.

## Result

`claude --mcp-config <file> --strict-mcp-config --allowedTools mcp__notebook__search -p`
on a local `CGO_ENABLED=0` build (17 MB), asked what was done about force push.

Server `notebook`: `connected`. Tool call:
`{"name":"mcp__notebook__search","input":{"query":"force push"}}`.

> A hook now asks for confirmation before a force push (`--force`) or remote
> branch delete (`:b`) and gives the reason, while plain pushes run unprompted;
> this was commit `aac1770`, pushed to main, recorded in
> `log/2026-10-29-workbench-hooks-01-force-push.md`.
