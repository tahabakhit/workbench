# 01 Add the `learning` Knowledge kind

Status: done

## What

Read `../spec.md` (Principles and Decisions) first.

- Add a `learning` row to the kinds table in `contract.md`: "missions,
  source lists, lessons and quiz records for topics being studied".
- Add `learning` to the kinds `skills/onboard/SKILL.md` asks about (the
  list at the start of the skill, now `log`, `decisions`, `lessons`,
  `corrections`, `plan`, `goals` and `sources`), so machine and repo scopes
  offer it.
- Add a `learning` line to this machine's
  `~/.config/workbench/knowledge.md`, by running `onboard` at machine scope
  or by hand, only with the owner's approval of the destination (a vault
  folder such as `~/docs/agents/learning/`).
- Don't change any other skill.

## Done when

- `sh scripts/check.sh` passes.
- `grep -c learning contract.md skills/onboard/SKILL.md` prints a nonzero
  count for each file.
- `sh hooks/session-context.sh | wc -c` prints a number under 10000.
- `grep '^- learning:' ~/.config/workbench/knowledge.md` prints the approved
  destination.
