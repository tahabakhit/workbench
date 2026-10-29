# 02 contract.md in this Mac's global instructions

Status: done
Blocked by: 01

## What

Write `contract.md` at the repo root: the always-on rules every machine
shares, in three sections. Keep it short; each line must change behaviour.

1. **Work tracking.** `.scratch/<feature>/spec.md` with a task table, one
   `issues/NN-slug.md` per task with a plain `Status: open|done|wontfix` line
   (not bold) and an optional `Blocked by: NN`. The session that meets a
   done-when sets `Status: done` in the issue and the spec table. Closing
   check: diff the result against done-when and Out of scope in a fresh
   subagent. Tracker updates go to the declared trackers, with evidence,
   applied on approval. Ad hoc work records nothing during the work.
   Decisions from a decision session are written only on approval: to
   `## Decisions` in `spec.md` when a feature folder exists, otherwise to the
   `decisions:` destination. ADR only when all three hold: hard to reverse,
   surprising without context, a real trade-off. From the `implement` and
   `tdd` verdicts: finish at done-when and quote the output, failing test
   first, mock only at boundaries.
2. **Knowledge destinations.** The `## Knowledge` section format and its
   kinds (`log`, `decisions`, `lessons`, `corrections`, `plan`, `goals`,
   `sources`), the vault per machine with repo override, and the
   hard/soft dependency rule (missing required destination: point to
   `onboard`).
3. **Writing rules** for text that lands in a repo, taken from the section
   "Writing that lands in a repository" in
   `~/dev/agents/home/.chezmoitemplates/shared-instructions.md`, plus the
   `pr` verdict (PR body as a format paragraph).

Then include it on this Mac. `~/dev/agents/home/private_dot_claude/private_CLAUDE.md.tmpl`
and `~/dev/agents/home/private_dot_codex/private_AGENTS.md.tmpl` build the
global files from `.chezmoitemplates/*.md`. Add `contract.md` from the
workbench checkout to both (check in the chezmoi docs how a template includes
a file outside the source directory), and move rules that the contract now
holds out of `shared-instructions.md`, `claude-instructions.md` and
`codex-instructions.md` so nothing is said twice. Commit in `~/dev/agents`.

## Done when

- `contract.md` exists, under 80 lines, with the three sections.
- `chezmoi diff` is empty after `chezmoi apply`, and both `~/.claude/CLAUDE.md`
  and `~/.codex/AGENTS.md` contain the contract's first heading (`rg` output
  quoted).
- `rg -n "Writing that lands" ~/dev/agents/home/.chezmoitemplates/` prints
  nothing.
