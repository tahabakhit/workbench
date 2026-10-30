# workbench

A Claude Code plugin, and its own marketplace, that holds one person's skills,
shared contract and the notebook (the markdown vault the skills file into)
for working with AI coding agents on any machine. Codex and Pi do not load
this plugin.

## Validate

Run after every change to the plugin, and quote the output:

```sh
sh scripts/check.sh
```

It needs `claude`, `jq` and `git`. Sessions load the plugin from
the GitHub marketplace, so a change reaches them only after it is pushed. To
pick it up right away, push, run `claude plugin marketplace update
workbench`, then start a new session.

## Plugin rules

- `hooks/session-context.sh` adds `contract.md` to every session; keep it
  under the hook's 10,000-character cap.
- Skill names are plain English verbs that name the outcome.
- A skill never loads another skill. Shared reference material is a side
  file copied into each skill folder that reads it.
- Records go to the notebook layout in `contract.md`. No skill hardcodes a
  tracker, chat or calendar; skills read `~/.config/workbench/sources.md`.
- `plugin.json` has no `version`, so every commit counts as an update.
- Scripts are POSIX `sh` with `jq`, executable, under the skill's `scripts/`.

## Work tracking

Follows `contract.md`. Open work is in `.scratch/workbench/`.
Settled input: `.scratch/bootstrap/decisions.md`,
which wins on conflict, and the per-skill verdicts in
`.scratch/bootstrap/mattpocock-analysis.md`. Terms are in `GLOSSARY.md`.

## Git and naming

- Commits: plain message, no `Co-Authored-By` trailer. This clone's local
  git config sets the personal noreply email and makes git push as the
  `tahabakhit` account. For `gh`, prefix
  `GH_TOKEN=$(gh auth token -u tahabakhit)`; the active gh account is the
  work one.
- Never name a file `claude.md` or `agents.md` in any case. macOS matches
  them to `CLAUDE.md` and `AGENTS.md`, and harnesses load those as
  instructions.
