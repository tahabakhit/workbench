# workbench

A Claude Code plugin, and its own marketplace, that holds one person's skills,
shared contract and the notebook (the markdown vault the skills file into)
for working with AI coding agents on any machine. Codex installs the same
plugin; Pi installs this repo as a package, through `package.json` and
`pi/workbench.ts`.

## Validate

Run after every change to the plugin, and quote the output:

```sh
sh scripts/check.sh
```

It needs `claude`, `jq` and `git`; the Pi extension tests skip without `node`.
Sessions load the plugin from
the GitHub marketplace, so a change reaches them only after it is pushed. To
pick it up right away, push, run `claude plugin marketplace update
workbench`, then start a new session.

## Plugin rules

- `hooks/session-context.sh` adds `contract.md` to every session; keep it
  under the hook's 10,000-character cap.
- Skill names are plain English verbs that name the outcome.
- `choose-skill` names every skill; change its map when a skill changes.
- A skill never loads another skill. Shared reference material is a side
  file copied into each skill folder that reads it.
- Records go to the notebook layout in `contract.md`. No skill hardcodes a
  tracker, chat or calendar; skills read `~/.local/share/workbench/sources.md`.
- `plugin.json` has no `version`, so every commit counts as an update.
- `pi/workbench.ts` is the one TypeScript file: Pi extensions cannot be shell.
  It mirrors `hooks/hooks.json`, so change them together.
- Scripts are POSIX `sh` with `jq`, executable: a skill's under its `scripts/`,
  plugin hooks under `hooks/`, repo checks under `scripts/`.

## Work tracking

Follows `contract.md`. Open work is in `.scratch/workbench/`.
Background input: `.scratch/bootstrap/decisions.md` and the per-skill verdicts
in `.scratch/bootstrap/mattpocock-analysis.md`. On conflict, the newest
`## Decisions` entry in a spec wins. Terms are in `GLOSSARY.md`.

## Git and naming

- Commits: plain message, no `Co-Authored-By` trailer, from the personal
  noreply email.
- Git pushes through the `gh` credential helper, as the active `gh` account.
  This repo belongs to `tahabakhit`, so that account must be the active one.
- Never name a file `claude.md` or `agents.md` in lower or mixed case. macOS
  matches them to `CLAUDE.md` and `AGENTS.md`, and harnesses load those as
  instructions.
