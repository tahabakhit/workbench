---
name: improve
description: Use when the user asks to review repeated corrections, whether earlier fixes held, or which skills go unused. Turns corrections made more than once into approved fixes (a hook, an instruction line, or a skill change) and proposes removals.
---

# Improve

Every change here is proposed, shown, and applied only after the user
approves it. Small edits only, never rewrites.

## 1. Read the ledger

Open `~/.local/share/workbench/notebook/corrections.md`. Group lines that
describe the same agent behaviour, even when worded differently.

- **Candidates**: groups with entries from two or more sessions and fix
  `open`.
- **Fixes that did not hold**: groups with a fix recorded and a newer entry.

Sort each candidate into one type, which hints at the fix:

- **Navigation**: the agent was slow to find a file or fact; a pointer helps.
- **Missing check**: a lint, test or hook could have caught the mistake. Look
  for an existing check that is unwired or broken first.
- **Rule**: a convention the agent broke. If a script can detect it, it gets
  a check, not prose. Only judgement calls become instruction lines.
- **No-op**: an instruction that changes nothing; propose cutting it.
- **Tool cost**: an expensive call that a cheaper one would replace.
- **Missing access**: information the agent needed but could not reach.

Report a repo with no checks at all (no hook, no CI job) as a finding beside
the candidates, even without ledger entries.

Done when every open line belongs to a group or is noted as seen once, and
each candidate has a type.

## 2. Propose fixes

For each candidate, pick the cheapest form that enforces it, using the table
in `cheapest-form.md` in this skill's folder: hook or permission rule for
"never", a check for anything a script can detect, an instruction line for
other "always" rules, the repo's `AGENTS.md` for one repo, a skill only for a
procedure or domain knowledge. A standard that only a review needs goes in
the file the reviewer reads, not in instructions every session loads. A fix that did not hold moves one step
stronger (instruction line to hook).

Before drafting, read the target. If it already says this, propose a clearer
wording or a better place, not a second copy. When a file is generated from a
source (a dotfiles repo, a template; the file usually says so at the top),
edit the source, never the generated file. A proposal that adds lines to an
instruction file also names a line to cut or merge, or says why none can go.

Each proposal shows: the ledger lines as evidence, the scope, the target
file, the diff, and why this form.

Done when every candidate has a proposal or a reason it needs none.

## 3. Propose removals

Run `scripts/skill-usage.sh 60`. A skill in the collection with zero uses in
60 days is proposed for removal, or for manual-only if it is model-invoked
and still wanted. In Claude Code, `/skill-doctor` also shows the context cost
of each skill. Also propose merging instruction lines that say the same
thing.

Done when each unused skill and each duplicate has a proposal.

## 4. Apply what is approved

Show all proposals in one list and let the user approve, edit or reject each.
For each approved one:

1. Apply it, one commit per fix; the commit message names the correction it
   addresses.
2. A skill change passes `claude plugin validate` on its skills folder. A
   change to a generated file's source is deployed the way that source documents.
3. Set the fix column of its ledger lines to the commit hash.

Rejected proposals get `rejected` in the fix column so they are not proposed
again unless new entries arrive.

Done when every approved fix is committed and applied and the ledger is
updated. Report what changed, what was rejected, and fixes that did not hold.
