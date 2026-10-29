# workbench

A Claude Code plugin, and its own marketplace, that holds one person's skills
and shared contract for working with AI coding agents across repos and
machines. Codex and Pi do not load this plugin.

## Validate

Run after every change to the plugin, and quote the output:

```sh
sh scripts/check.sh
```

It needs `claude`, `jq`, `git` and `gitleaks`. Sessions load the plugin from
the GitHub marketplace, so a change reaches them only after it is pushed. To
pick it up right away, push, run `claude plugin marketplace update
workbench`, then start a new session.

## Plugin rules

- `hooks/session-context.sh` adds `contract.md` and the machine's
  `~/.config/workbench/knowledge.md` to every session; keep the two under its
  10,000-character cap.
- Skill names are plain English verbs that name the outcome.
- A skill never loads another skill. Shared reference material is a side
  file copied into each skill folder that reads it.
- No skill hardcodes a tracker, chat, calendar or vault. Skills read the
  `## Knowledge` destinations.
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
