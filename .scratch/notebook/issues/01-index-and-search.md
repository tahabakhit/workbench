# 01 Go module, index and notebook search

Status: done

## What

Create the Go module in `notebook/` (`go 1.26`, `modernc.org/sqlite` v1.61.1)
with the `notebook` command and its `search` subcommand. Follow `## Index and
search`, `## Commands` and `## Layout and paths` in `.scratch/notebook/spec.md`.

- Read the vault from `~/.local/share/workbench/notebook`; a `--vault` flag
  overrides it for tests. Index frontmatter (`date`, `repo`, `scope`,
  `supersedes`), title, body, and kind (the top folder).
- Store the index at `~/.local/share/workbench/state/notebook.db`; `--index`
  overrides it. Refresh it by file modification time at each command.
- Hide notes named by another note's `supersedes` unless `--all` is passed.
- Results print path, stable ID, score and snippet, plain by default, with
  `--json` for scripts.
- Add a fixture vault in `notebook/testdata/`, with expected results by hand.
- Add `go test ./...` for `notebook/` to `scripts/check.sh`, skipped with a
  message when `go` is missing, the way the Pi tests skip without `node`.

Do not change skills, hooks or `contract.md`. Do not write to the vault.

## Done when

- `sh scripts/check.sh` passes and its output includes the Go tests.
- `cd notebook && go run . search "force push"` lists
  `log/2026-10-29-workbench-hooks-01-force-push.md` first.
- A test edits a fixture note's modification time and content, and the next
  search finds the new text.
- Deleting the index file and searching again gives the same results.
