# 02 recall, links and backlinks

Status: done
Blocked by: 01

## What

Add `recall`, `links` and `backlinks` to `notebook/`, as in `## Commands` and
`## Index and search` in `.scratch/notebook/spec.md`.

- `recall` takes a path or title and prints the whole note. No match, or more
  than one, exits non-zero and lists the candidates.
- A wikilink resolves by path, then exact title or filename, to exactly one
  note. Otherwise it stays unresolved and is reported as broken or ambiguous.
- Store raw link targets in the index so a rebuild resolves them again.

## Done when

- `sh scripts/check.sh` passes, with tests for each case below on the fixture
  vault: a unique match, no match, two matches, a broken link, an ambiguous
  link, and a link that resolves once its target note is added.
- `cd notebook && go run . recall log/2026-10-29-workbench-hooks-01-force-push.md`
  prints that note in full.
