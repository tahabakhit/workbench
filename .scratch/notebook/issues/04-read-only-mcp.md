# 04 Read-only notebook mcp

Status: open
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
