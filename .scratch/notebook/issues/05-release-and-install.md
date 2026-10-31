# 05 Release workflow, shim and install

Status: open
Blocked by: 03, 04

## What

Get the binary onto a machine with no Go. Follow the first two rows of
`## Decisions` and `## Layout and paths` in `.scratch/notebook/spec.md`.

The owner has approved this workflow creating releases on every change to
`notebook/`, so it needs no further sign-off (see `## Decisions` in the spec).

- `notebook/source-hash.sh`: prints one hash of the files in `notebook/`, the
  same on Linux and macOS (`sha256sum` or `shasum -a 256`).
- `.github/workflows/notebook-release.yml`: on a push to `main` touching
  `notebook/`, runs the tests, builds with `CGO_ENABLED=0` for linux and darwin
  on amd64 and arm64, and publishes release `notebook-<hash>` with a checksums
  file.
- `notebook/run.sh`: if the installed binary matches the source hash, starts
  it. Otherwise downloads the release for this platform, verifies the checksum,
  installs atomically to `~/.local/share/workbench/bin/notebook` and starts it.
  If the release isn't up yet, it keeps using the installed binary. A second
  concurrent run must not corrupt the install.
- `.mcp.json` at the plugin root: runs `${CLAUDE_PLUGIN_ROOT}/notebook/run.sh`
  with `mcp`.
- `hooks/hooks.json` and `pi/workbench.ts`, changed together: `SessionStart`
  starts the install in the background, since the hook has a 10-second timeout.
- `scripts/check.sh`: tests for `source-hash.sh` and `run.sh`, with the
  download mocked as a local file, including a bad checksum.

Check whether Codex loads the plugin's `.mcp.json` and record the result here.
If it doesn't, Codex uses the CLI, the same as Pi. Do not add a top-level
`bin/`, which claude.ai and Cowork refuse.

## Done when

- `sh scripts/check.sh` passes, with the bad-checksum test refusing the install.
- `gh release view notebook-<hash>` lists four binaries and the checksums file.
- With the installed binary deleted, a new Claude Code session on the pushed
  plugin calls the MCP `search` tool successfully, and the binary is back.
- A plain `codex exec` and `pi -p` session each start with the binary present,
  or the gap is named here.
