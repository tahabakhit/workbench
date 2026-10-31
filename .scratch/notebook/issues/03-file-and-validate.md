# 03 file and validate

Status: open
Blocked by: 01

## What

Add `file` and `validate` to `notebook/`, following `## Validation` in
`.scratch/notebook/spec.md` and the layout table and frontmatter rule in the
Notebook section of `contract.md`.

- `file` takes the target path and reads the note from stdin. It validates,
  writes, runs `git add` on that path and commits with a plain message naming
  it. It creates the vault repository and the folder when missing.
- `corrections.md`, `goals.md` and `plans/` entries may be updated, since
  skills append to them; `--update` is required for an existing path, and
  without it the same-name refusal applies.
- Refusals exit non-zero and change nothing. Warnings print and do not stop.
- If the commit fails, the file is left unwritten or restored, and the command
  exits non-zero with git's error. A test covers this.
- `validate` with no path checks the whole vault and reports each problem.

Tests run against a temporary git repository. Never write to the real vault in
a test.

## Done when

- `sh scripts/check.sh` passes, with one test per row of the Validation table
  plus the failed-commit case.
- A test files a good note and sees exactly one new commit holding only it.
- `cd notebook && go run . validate` on the real vault reports each problem by
  path, and the report is quoted in this issue.
