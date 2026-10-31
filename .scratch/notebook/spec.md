# notebook

## Goal

One Go binary, `notebook`, that indexes, searches and files into the markdown
vault at `~/.local/share/workbench/notebook`. It serves a CLI for every harness
and an MCP server for ad hoc lookups, ships in the workbench plugin, and takes
over filing from `record`, `plan-day` and `study-topic`. Files stay plain
markdown in git; Obsidian stays an optional viewer.

Settled input: the `notebook (vault tool)` section of
`.scratch/bootstrap/decisions.md` and the Notebook section of `contract.md`.
On conflict, the newest entry under `## Decisions` here wins, then
`contract.md`, then the decisions file. Terms are in `GLOSSARY.md`.

## Principles

- Influence, not force: refuse only writes that would damage the vault; warn on
  the rest.
- Markdown files are the truth. The index can be deleted and rebuilt at any
  time, and nothing is stored only in it.
- The tool works without Obsidian, a background process or a network, except
  to fetch its own binary.
- `notebook file` is the only write path in the tool. MCP is read-only.
- Skills reach the tool through the CLI at its fixed path. If the binary is
  missing they write the file directly, as they do today.
- Write the failing test first. Mock only at system boundaries. Test fixtures
  are a small vault in `notebook/testdata/`, with expected results written by
  hand.
- No Python. Go for the tool; POSIX `sh` with `jq` for scripts.

## Layout and paths

```
notebook/                 Go module: source, tests, testdata
notebook/run.sh           shim: fetch, verify, start the binary
notebook/source-hash.sh   hash of notebook/, shared by the shim and CI
.mcp.json                 plugin MCP config, points at the shim
.github/workflows/        release workflow
```

| Path | Holds |
|---|---|
| `~/.local/share/workbench/bin/notebook` | the installed binary, outside the plugin folder so updates keep it |
| `~/.local/share/workbench/state/notebook.db` | the search index, rebuildable, out of git |

## Commands

| Command | Does | MCP |
|---|---|---|
| `search <query>` | ranked snippets with a stable ID; `--repo`, `--kind`; `--all` includes superseded notes | yes |
| `recall <path or title>` | one whole note; fails on no match or more than one | yes |
| `links <note>` | outgoing wikilinks, flagging broken and ambiguous | yes |
| `backlinks <note>` | notes that link to it | yes |
| `file` | validates, writes, commits | no |
| `validate [path]` | checks one note or the vault | no |
| `push --repo <name>` | the session-start text | no |
| `mcp` | serves the MCP tools over stdio | no |

## Index and search

- Pure-Go SQLite `modernc.org/sqlite` v1.61.1 with FTS5 and `mcp-go` v1.5.0.
  Both built together with `CGO_ENABLED=0` in a trial, 15 MB, 9s from cold.
- FTS tokenizer `porter unicode61 remove_diacritics 1`, title ranked far above
  body, 20-token snippets, starting from zk's values and tuned by tests.
- Each command refreshes the index first: walk the files and reindex any whose
  modification time differs from the stored one. No file watcher.
- Raw link targets are stored so a rebuild can resolve them again. A wikilink
  resolves to exactly one note or stays unresolved; it never picks the first.

## Validation

| Case | Result |
|---|---|
| Path outside the `contract.md` layout | refuse |
| Missing or malformed frontmatter (`date`; `scope` on a lesson) | refuse |
| `supersedes` target does not exist | refuse |
| A note of the same name already exists | refuse |
| Broken or ambiguous wikilink | warn |

## Tasks

| # | Task | Blocked by | Status |
|---|---|---|---|
| 01 | Go module, index and `notebook search` | | done |
| 02 | `recall`, `links`, `backlinks` | 01 | done |
| 03 | `file` and `validate` | 01 | done |
| 04 | Read-only `notebook mcp` | 02 | done |
| 05 | Release workflow, shim and install | 03, 04 | done |
| 06 | Session push | 02, 05 | done |
| 07 | Skills file through the tool | 03, 05 | done |
| 08 | Daily-use check on this PC and close workbench issue 11 | 06, 07 | open |

## Not yet specified

- How the notebook syncs between machines, if a second machine ever uses it.
- What values lesson `scope` takes beyond a repo name and `global`.

## Out of scope

- Semantic search and rank fusion: later, after keyword search has been used.
- MCP write tools: added only if the CLI proves not to be enough.
- Porting obsidian-mind commands: dropped by the decisions file.
- Periodic vault lint and wiki-style page updates (the chief-of-staff verdict):
  later, not built here.
- A hook that checks direct edits to the vault: `notebook validate` covers the
  rare fallback write.
- MCP in Pi: Pi uses the CLI.
- Switching other machines: the owner dropped them; see Decisions.

## Decisions

| Decision | Reason |
|---|---|
| CI builds one binary per platform (linux and darwin, amd64 and arm64), published as a release tagged with a hash of `notebook/` | No Go needed on a machine; a commit outside `notebook/` fetches nothing. Building locally needs Go, and committing binaries adds 15 MB per platform per change to public history |
| The owner gives standing approval for `notebook-release.yml` to run and create `notebook-<hash>` releases on every change to `notebook/` | One approval for an automatic step; it covers this workflow only |
| `notebook/run.sh` fetches and verifies on first call and backs the plugin's MCP config; `SessionStart` installs the binary in the background | An MCP server may start before the hook finishes, and Codex runs hooks only once trusted. A top-level `bin/` is avoided because claude.ai and Cowork refuse plugins that have one |
| Skills use only the CLI; MCP is read-only at first | One instruction works in every harness, and writes have one path |
| `notebook file` validates, writes and commits | Skills stop repeating the commit step. No tool examined auto-commits, so failure handling gets its own test |
| Refusals and warnings as in Validation | Follows influence, not force |
| Flat `log/` with a `repo` field, no project folders | `contract.md` is newer than the decisions file's "project folder", and the vault already uses it |
| `recall` needs exactly one match; `search --all` shows superseded notes | First-match lookup can return the wrong note silently, and hidden notes need a way out |
| The push sends at most 2,000 characters of paths and titles: the 3 latest logs for this repo, then lessons scoped to it, then global lessons. It drops in that reverse order, notes what it dropped, prefers lessons that work as rules, and has an off switch | The contract uses 5,306 of the hook's 10,000 characters. One study found repo context files raised cost over 20% without generally improving success, so the push must earn its place |
| Done-when covers this PC only, with no switch-over slice | This PC has no Obsidian and no obsidian-mind fork, and the owner dropped other machines |
| The MCP config is inline in `plugin.json`, not a root `.mcp.json` | Claude Code reads a root `.mcp.json` as this repo's project config too, where `CLAUDE_PLUGIN_ROOT` is unset, so every session here showed a failed server. Codex read that file but could not expand the variable, so it loses nothing |
| The push stays on, and the on/off session comparison is dropped | The owner will not run manual trials, and agent runs could not show whether the push helps real work. It costs 368 characters today, capped at 2,000, against a 10,000 cap, and `WORKBENCH_NOTEBOOK_PUSH=off` remains |
