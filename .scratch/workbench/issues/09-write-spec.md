# 09 write-spec skill

Status: open
Blocked by: 02

## What

`skills/write-spec/`, user-invoked. Turns a settled conversation into
`.scratch/<feature>/spec.md` and `issues/NN-slug.md` in one pass, with no
new interview. Sources: the decisions file (Workflow section) and the rows
to-spec, to-tickets and triage in `.scratch/bootstrap/mattpocock-analysis.md`.

- spec: goal, a real Out of scope list, a task table, an empty
  `## Decisions` section. No long user-story lists.
- issues: vertical slices, prefactor first, `Blocked by`, plain `Status:
  open`, a done-when a fresh session can verify.
- Before writing, check whether the concept already exists in the repo
  (triage row).
- Show the slicing as a numbered quiz for approval before writing files.
- Must work for code and non-code work alike (no typecheck or test-runner
  assumptions).

## Done when

- `sh scripts/check.sh` exits 0.
- Running it on this repo's bootstrap input would produce a spec with the
  same sections as `.scratch/workbench/spec.md` (shown in preview, nothing
  written).
