# 06 Session push

Status: done
Blocked by: 02, 05

## What

Add `notebook push --repo <name>` and call it from the session hook, following
the push row in `## Decisions` in `.scratch/notebook/spec.md`.

- Output: paths and titles only, at most 2,000 characters. The 3 latest logs
  for this repo, then lessons with this repo's `scope`, then `global` lessons.
  Superseded notes are left out.
- Over the limit, drop global lessons first, then repo lessons, then logs, and
  end with one line saying how many were dropped.
- `hooks/session-context.sh` appends the output after `contract.md`, with the
  repo name from the git root. A missing binary, an empty result or
  `WORKBENCH_NOTEBOOK_PUSH=off` prints the contract alone, and the hook never
  fails because of the push.
- Keep `pi/workbench.ts` in step with the hook.

`contract.md` plus the push must stay under the hook's 10,000-character cap.

## Done when

- `sh scripts/check.sh` passes, with tests for the drop order, the dropped
  line, the 2,000-character limit and the off switch, and a check that the
  contract plus 2,000 characters is under 10,000.
- `sh hooks/session-context.sh` run in this repo prints the contract followed
  by the `log/2026-10-29-workbench-*` notes, quoted in this issue.
- The same hook run with the binary removed prints the contract and exits 0.

## Result

`sh hooks/session-context.sh` in this repo, after the contract:

```
## Notebook

Notes for workbench. Read one with "notebook recall <path>".

- log/2026-10-29-workbench-topic-tutor-01-05.md: topic-tutor spec and issues 01-05
- log/2026-10-29-workbench-hooks-02-04-gitleaks-dirty-rm.md: hooks issues 02-04, gitleaks deny, dirty-tree and rm asks
- log/2026-10-29-workbench-hooks-01-force-push.md: hooks issue 01, ask on force push and remote branch delete
```

With the binary moved away, the output equals `contract.md`, exit 0. The
contract is 5,306 characters, so with the push it stays under 10,000.

Inside the lessons tiers, notes go newest first. "Prefers lessons that work as
rules" in `## Decisions` has no check here: nothing marks a lesson as a rule.
