# 08 Daily-use check on this PC and close workbench issue 11

Status: open
Blocked by: 06, 07

## What

Confirm the tool held up in ordinary work on this PC, then close the parent
issue. This asks nothing special of the owner: they work as usual, and a later
session reads what happened from the vault's git log. It needs commits from
`record`, `plan-day` and `study-topic` on at least three separate days, so it
can't finish in one sitting.

The push on/off comparison is dropped; see `## Decisions` in
`.scratch/notebook/spec.md`.

## Done when

- This issue quotes one dated run of each of the three skills, from
  `git -C ~/.local/share/workbench/notebook log --format='%ad %h %s' --date=short`,
  and the final `~/.local/share/workbench/bin/notebook validate` output.
- `.scratch/workbench/issues/11-notebook.md` and its row in
  `.scratch/workbench/spec.md` are `Status: done` with the evidence linked.

## Runs so far

Day 1, 2026-10-31: `record` `68ae543`, `plan-day` `e9fd6c8`, `study-topic`
`e7af6e2`, as checks for issue 07.
