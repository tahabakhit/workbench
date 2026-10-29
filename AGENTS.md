# workbench

A Claude Code plugin, and its own marketplace, that holds one person's skills
and shared contract for working with AI coding agents across repos and
machines (a work Mac and a Fedora PC). A `SessionStart` hook in `hooks/`
adds `contract.md` and the machine's `~/.config/workbench/knowledge.md` to
every Claude Code session; keep the two under the hook's 10,000-character
cap. Codex and Pi do not load this plugin.

## Work tracking

- One folder per feature under `.scratch/<feature>/`: `spec.md` with a task
  table and an empty `## Decisions` section, plus one `issues/NN-slug.md` per
  task.
- Each issue has a plain `Status: open|done|wontfix` line (not bold), an
  optional `Blocked by: NN` line, and a "Done when" check.
- The session that meets a done-when sets `Status: done` in the issue and in
  the spec table. Before that, diff the result against done-when and Out of
  scope in a fresh subagent.
- Current work: `.scratch/workbench/spec.md`. Settled input:
  `.scratch/bootstrap/decisions.md`, which wins on conflict, and the
  per-skill verdicts in `.scratch/bootstrap/mattpocock-analysis.md`.
- Decisions are written only on approval, to `## Decisions` in the feature's
  `spec.md`. ADR in `docs/adr/` only when a decision is hard to reverse,
  surprising without context, and a real trade-off.

## Plugin rules

- Skill names are plain English verbs that name the outcome.
- A skill never loads another skill. Shared reference material is a side
  file copied into each skill folder that reads it.
- No skill hardcodes a tracker, chat, calendar or vault. Skills read the
  `## Knowledge` destinations.
- `plugin.json` has no `version`, so every commit counts as an update.
- Scripts are POSIX `sh` with `jq`, executable, under the skill's `scripts/`.

## Validate

Run after every change to the plugin, and quote the output:

```sh
sh scripts/check.sh
```

It runs `claude plugin validate --strict`, which turns warnings (unknown
fields, missing metadata) into failures, except the missing-version warning.
It also fails if copies of a shared side file differ between skill folders.
Sessions load the plugin from the GitHub marketplace, so a change reaches
them only after it is pushed. To pick it up right away, push, run `claude
plugin marketplace update workbench`, then start a new session.

## Writing in this repo

- Text that lands here must make sense to someone who wasn't there. Leave out
  dates, run numbers, ticket IDs and people's names. Describe what is true
  now; history goes in the commit message.
- A comment states a non-obvious constraint in one or two sentences.
- Commits: plain message, no `Co-Authored-By` trailer. This clone's local
  git config sets the personal noreply email and makes git push as the
  `tahabakhit` account. For `gh`, prefix
  `GH_TOKEN=$(gh auth token -u tahabakhit)`; the active gh account is the
  work one.
- Never name a file `claude.md` or `agents.md` in any case. macOS matches
  them to `CLAUDE.md` and `AGENTS.md`, and harnesses load those as
  instructions.
