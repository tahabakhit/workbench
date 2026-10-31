# 08 Daily-use check on this PC and close workbench issue 11

Status: open
Blocked by: 06, 07

## What

Prove the tool in real work on this PC, then close the parent issue. This needs
real sessions across at least three separate days, so it can't finish in one
sitting.

- Use `record`, `plan-day` and `study-topic` in normal work.
- Compare sessions with the push on and off
  (`WORKBENCH_NOTEBOOK_PUSH=off`): whether the pushed notes were opened or
  changed the work, and what the push cost in context. Remove the push if it
  doesn't help, and record that under `## Decisions` in the spec.
- Run `notebook validate` over the vault at the end.

## Done when

- This issue quotes one dated run of each of the three skills and the final
  `notebook validate` output.
- The push decision, keep or remove, is under `## Decisions` in
  `.scratch/notebook/spec.md`.
- `.scratch/workbench/issues/11-notebook.md` and its row in
  `.scratch/workbench/spec.md` are `Status: done` with the evidence linked.
