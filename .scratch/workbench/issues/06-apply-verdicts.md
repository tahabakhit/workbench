# 06 Apply the "merge into own" verdicts

Status: done
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

## Closing note

| Row | Change |
|---|---|
| grilling | `settle` step 1, items 3 and 4: lookups in a subagent that block only dependent questions; questions worded so "yes" accepts the recommendation; wait for answers; act on nothing before the confirmed list |
| grill-with-docs | `settle` step 2: optional GLOSSARY.md and ADR proposals, written only on approval |
| wayfinder | `settle` step 1, items 1 and 2: destination, Not yet specified vs Out of scope, a sharp but unsettled question becomes an open item |
| loop-me | `settle` intro: optional output file kept current; Workflows technique: trigger, checkpoint, brief, push checkpoints late |
| codebase-design | `settle` techniques: "Design it twice" |
| decisions file | `settle` step 2: `## Decisions` in `spec.md`, else the `decisions:` destination, else the final summary |
| domain-modeling | new `skills/settle/glossary-adr-format.md`, read in step 2; copied into `onboard` by issue 03 |
| diagnosing-bugs | `diagnose`: redaction line at the top; step 1 stop-and-ask and minimise; step 2 show hypotheses; step 4 cause in the commit message |
| to-questionnaire | `draft-questions`: "why this matters" only where needed, "I don't know" welcome, catch-all last question, coverage done-check |
| wait-what | `rephrase`: ASD-STE100 as the leading word, GLOSSARY.md terms, a missing premise may be added |
| retro | `improve` step 1: candidate types; a repo with no checks is reported as a finding; "if a script can detect it, it gets a check", also in step 2 |
| writing-for-agents | new `skills/write-skill/vocabulary.md`, read at the start of step 3 |
| setup-ts-deep-modules | `write-skill` step 5, check 4: prove a new check bites (pass, fail, pass) |

Not done: the loop-me row also asks for its terms in a plugin `GLOSSARY.md`.
No such file exists and no issue creates one, so the terms are defined inline
in `settle` instead.
