# 05 Merge retro and harvest into record

Status: done
Blocked by: 02, 04

## What

`skills/record/`, user-invoked, from
`~/dev/agents/home/dot_agents/skills/retro/` and `.../harvest/` (with its
`scripts/sessions.sh` and `scripts/digest.sh`). One skill, two scopes:

- `session` (default): the current session, as retro does now.
- `since-last`: every source since the last run, as harvest does now.

Both scopes share the four kinds (log, decisions, lessons, corrections), the
filter and the destinations from `## Knowledge`. Ad hoc work has no feature
folder, so both read the session transcript and git log. A missing required
destination points to `onboard`. Keep harvest's state folder
(`~/.local/state/harvest`) so the next run finds the last window.

The Mac's `## Knowledge` section has a `harvest:` line, which is not one of
the kinds `contract.md` defines. Turn it into kind lines limited to the
`since-last` scope (for example `log (for record since-last)`) when
`record` replaces harvest. That section now lives in the personal dotfiles
repo, so this step is tracked there in `.scratch/mac-dotfiles/issues/03-instructions.md`.

## Done when

- `sh scripts/check.sh` exits 0.
- `.scratch/workbench/trials/record.md` holds the preview output of one
  `session` run and one `since-last` run on this Mac, with nothing filed.
