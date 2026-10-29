# record trials

Both runs used a fresh headless session with this working copy loaded as the
plugin, `preview` set, Write, Edit and the vault's filing tools
(`record_work`, `remember`) disallowed. Nothing filed: checksums of the
correction ledger and `pending.jsonl` match before and after, and
`~/.local/state/harvest/last-run` was not created. Entries from work repos
are left out.

`scripts/sessions.sh` now also marks sessions that ran `record` as recorded.
A two-line fixture showed it bites: `workbench:record` true, `rephrase`
false, `record` true.

## session, preview

Run on a fork of the session that built `onboard` (`--resume <id>
--fork-session --no-session-persistence`), so the original transcript was
not touched. It says the vault's write tools "aren't connected"; they were
disallowed for the trial.

The run named the destination of each kind, drafted log entries and lessons
for the vault, and drafted three correction lines in the ledger format:

```
- 2026-10-29 | workbench | claude | ended a finished task with loose observations and no direction | one recommended next step with a paste-ready prompt | open
- 2026-10-29 | workbench | claude | asked approval to push a private personal repo the rules already allowed | push when it is the natural next step of agreed work | open
- 2026-10-29 | workbench | claude | asked a decision question in jargon and left unclear which machine runs what | plain-terms questions that name where each step runs | open
```

It listed what it dropped and why: a mistake the agent caught itself, an
unverified claim, and work still in progress.

## since-last, preview, window from 2026-10-28

A first run could not read `~/.claude/projects`: the trial's allow rules
did not match the script paths and the folder was not added. The rerun
added it.

The rerun covered 13 sessions, skipped 3 already recorded, and folded a
one-prompt trial into another entry. It drafted 5 log entries, 1 decision,
2 lessons and 4 corrections, grouped by repo, each log entry with changes,
verification and open threads. It linked decisions already in a spec
instead of copying them, and noted that one correction repeated a rule the
global instructions already had. Every source was queried; background
subagents could not run Bash, so one was rerun in the foreground.

Lessons it drafted:

1. A skill that says "show the draft" gets a summary; it must say "print
   verbatim, even if unchanged". Two of three onboard trials summarised the
   draft.
2. Adding a column to an incremental dbt model fills only the rows later
   runs recompute; a full refresh rewrites every column.
