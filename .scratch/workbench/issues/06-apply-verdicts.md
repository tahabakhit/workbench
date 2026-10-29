# 06 Apply the "merge into own" verdicts

Status: open
Blocked by: 04

## What

Change each moved skill as the verdict rows in
`.scratch/bootstrap/mattpocock-analysis.md` say. Small edits, no rewrites.
The source for each change is the named row.

| Skill | Rows |
|---|---|
| settle | grilling, grill-with-docs, wayfinder, loop-me, codebase-design ("design it twice" line) |
| diagnose | diagnosing-bugs |
| draft-questions | to-questionnaire |
| rephrase | wait-what |
| improve | retro (candidate types, "an unguarded repo is a finding", "mechanical gets a check") |
| write-skill | writing-for-agents (vocabulary side file), setup-ts-deep-modules |

Also, from the decisions file: `settle` writes decisions only on approval,
to `## Decisions` in `spec.md` when a feature folder exists, otherwise to
the `decisions:` destination. From the domain-modeling row: one format side
file for GLOSSARY.md and ADRs, copied into `settle` and `onboard`, proposed
at the end and written only on approval. If `onboard` doesn't exist yet,
issue 03 picks up the copy.

## Done when

- `sh scripts/check.sh` exits 0, including the side-file copy check.
- For each row in the table, `git diff` shows the change, and the issue lists
  row to diff hunk in a closing note.
