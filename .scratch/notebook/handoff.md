# Handoff

## Goal

Close notebook issue 08 and workbench issue 11 once the vault shows real use on three separate days; until then, pick the next work from `/workbench:plan-day`.

## State

- Branch `main`, level with `origin/main` once this file is pushed. No uncommitted changes.
- Done: notebook issues 01 to 07. The `notebook` tool searches, recalls, follows links, files, validates, serves read-only MCP and pushes notes at session start.
- Open: notebook issue 08 (GitHub #1), the last open issue in `.scratch/`; workbench issue 11 (GitHub #2) closes with it. Close both on GitHub too, with the evidence.
- Installed here: plugin `c2bb62d` in Claude Code, plus Codex and Pi; binary from release `notebook-1a42e3560db85a95`.
- Not in this repo: `delegatr` distribution decision, paused mid-grilling at round 2 (`plans/2026-10-31.md` in the notebook).

## Claims

- verified: `sh scripts/check.sh` exits 0 at `c2bb62d`.
- verified: `~/.local/share/workbench/bin/notebook validate` exits 0 with no output on the real vault.
- verified: with the binary deleted, a new Claude session connects `plugin:workbench:notebook`, calls `search`, and the binary returns (`issues/05-release-and-install.md`).
- verified: `rm` of the binary alone reinstalls it from the release; `scripts/test-notebook-shim.sh` covers it.
- verified: in this repo `claude mcp list` shows only the plugin's server, Connected, after the config moved into `plugin.json`.
- verified: real `record`, `plan-day` and `study-topic` runs commit through `notebook file`, and `record` writes directly with no binary (`issues/07-skills-file-through-the-tool.md`).
- verified: Codex 0.166.0 gets no notebook MCP tools because it leaves `${CLAUDE_PLUGIN_ROOT}` unexpanded; its CLI use works, as does Pi's (`issues/05-release-and-install.md`).
- verified: the push is 368 characters in this repo today, and the contract is under the 8,000 that leaves room for it.
- verified: the owner approved pushing and releases on 2026-10-31. Auto mode blocked the first release push until they said so, so a later push touching `notebook/` may be stopped again; it is covered by the standing approval in `spec.md`.
- assumed: the darwin binaries run. They were built and checksummed, never executed, and macOS is out of daily use.
- assumed: the plain `record` and `study-topic` wording works in interactive sessions. The real runs used `claude -p` with tools allowed by flag.

- verified: tracking runs on GitHub. The `tracker (<repo>)` lines in `~/.local/share/workbench/sources.md` map workbench and delegatr to project 3, dotfiles to 4, homelab to 2 and salon to 6.

## Known gaps

- "Prefers lessons that work as rules" in the push has no check; lessons go newest first (`issues/06-session-push.md`).
- `--update` replaces a whole file, so `study-topic` must read, add and pipe back. Only skill wording guards `never overwrite`.
- The vault's two logs from 2026-10-31 describe 07 as open and the MCP failure as unsolved; both were later fixed. Left as written.
- The quiz has one open miss in lesson 01, recorded on purpose as `wrong`, and lesson 02 is unquizzed.
- `state/pending-corrections.jsonl` holds 2 unfiled lines; the next plain `record` run handles them.
- Stale plugin versions pile up in `~/.claude/plugins/cache/workbench/workbench/` (26 folders). Claude Code owns that folder; nothing here cleans it.

## Next

1. Run `/workbench:plan-day` at the start of a work day. Done when today's plan is committed in the notebook.
2. When `git -C ~/.local/share/workbench/notebook log --format='%ad %s' --date=short` shows the three skills on three separate days, finish `issues/08-daily-use-check.md`. Done when its done-when is quoted there and issues 08 and 11 are `done` in both specs.
3. Settle `delegatr` round 2 in its own repo. Done when that decision is written down.

## Read first

- `.scratch/notebook/spec.md`, `## Decisions`
- `.scratch/notebook/issues/08-daily-use-check.md`
- `.scratch/workbench/issues/11-notebook.md`
- `contract.md`, `## Notebook`
- `AGENTS.md`, `GLOSSARY.md`
- `notebook/`, `notebook/run.sh`, `.github/workflows/notebook-release.yml`
- commits `254ed67`, `7aa7618`, `4650071`

## Skills

- `/workbench:plan-day`
- `/workbench:record`
- `/workbench:study-topic posix-sh-parameter-expansion`
