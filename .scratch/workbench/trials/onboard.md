# onboard trials

Both runs used a fresh headless session with this working copy loaded as the
plugin (`claude -p "/workbench:onboard <scope> preview" --plugin-dir ~/dev/workbench`),
Write and Edit disallowed, and only read-only tools allowed. Details of the
work machine's tools and of the work repo are left out.

## Machine scope, preview, on the work Mac

Run from an empty folder. Two earlier runs covered every kind but did not
print the drafted section (one called it "identical", one gave no draft);
step 3 of the skill now says to print it verbatim even when unchanged. The
second earlier run also claimed the generated `~/.claude/CLAUDE.md` has no
generated-file marker, which is false (line 1 has one).

The third run listed every connected tool with the kind it serves and the
result of one test call, asked one numbered round with a recommendation per
kind, and printed the drafted `## Knowledge` section in full. Every kind of
the existing section appeared verbatim in the draft, including the
skill-limited `sources` lines. Nothing was written: checksums of the
instruction files match before and after.

## Repo scope, preview, on a work repo

The repo had mattpocock `docs/agents/` files, an `AGENTS.md` over 60 lines,
a `CLAUDE.md` with more than `@AGENTS.md`, and `.scratch/` tickets with bold
`**Status:**` lines. The run surveyed 13 artifacts and proposed a row for
each: delete the three `docs/agents/` files, convert `AGENTS.md`,
`CLAUDE.md`, `CONTEXT.md` (to `GLOSSARY.md`) and the `.scratch/` tracking,
keep the ADR and the plugin parts. It ran the safe commands (typecheck and
tests passed, a format check failed) and listed the unsafe ones as not run.
It also found two tickets marked open whose commits had already implemented
them. Spot-checked claims held. Nothing applied: `git status --short` was
empty before and after.
