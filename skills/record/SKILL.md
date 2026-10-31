---
name: record
description: Use when the user asks to log, file, save or keep a task, decision or lesson, or when one is finished and worth keeping. Files them into the notebook, with corrections. Scope `session` (default) covers this session; `since-last` catches up on every session and source since the last run, after a preview the user approves.
---

# Record

Turn work into entries of four kinds. A **log** entry records one task. A
**decision** is a choice worth explaining later. A **lesson** is a claim still
true in a different repo six weeks from now. A **correction** is a time the
user had to fix how the agent worked. Everything else is dropped.

Scope comes from the argument: `session` (default) or `since-last`. With
`preview`, show every draft and file nothing, in either scope.

## 1. Notebook

Entries go to the notebook paths in the contract: `log/`, `decisions/`,
`lessons/` and `corrections.md` under `~/.local/share/workbench/notebook`. Create the repository
or a folder when it is missing.

Done when `~/.local/share/workbench/notebook` is a git repository.

## 2. Gather

**session**: this conversation, plus `git log` and `git diff` since it
started: work done, decisions and rejected alternatives, surprises, how
results were verified, open threads.

**since-last**: the window starts at a date given as an argument, else at
the timestamp in `~/.local/share/workbench/state/record-last-run`; if neither exists, ask
for a start date. It ends now. Query every source:

| Source | How |
|---|---|
| Sessions | `scripts/sessions.sh <start>` in this skill's folder lists Claude Code and interactive Codex sessions, one JSON line each. Skip this session, sessions with `recorded: true`, and one- or two-prompt sessions with no edits. `scripts/digest.sh <file>` prints a session without tool output. |
| Native memories | `find ~/.claude/projects/*/memory -name '*.md' ! -name MEMORY.md -newermt <start>` |
| Git | per repo seen in the sessions: `git -C <repo> log --since=<start> --author="$(git -C <repo> config user.email)" --format='%h %ad %s' --date=short` |
| Declared | every line in `~/.local/share/workbench/sources.md` and in the `## Sources` of each repo seen, with its stated limits; with neither, note the gap |

With more than five sessions, hand each project's digests to a subagent for
candidates; verify each kept candidate against its digest. Group evidence
into tasks: a ticket key, branch or PR ties a session, commits and a ticket
together.

**Both scopes**: flagged corrections in
`~/.local/share/workbench/state/pending-corrections.jsonl` inside the session or window
(for `session`, also matching this repo's `cwd`), plus moments the hook missed
where the user rejected or redid the agent's work, plus struggles nobody
corrected: a file or fact that took many tool calls to find, or a call a
cheaper one would replace. A struggle counts only with a transcript quote
that shows it; file it as a correction with `self-noticed` as what the user
wanted. In `since-last`, confirm each flagged line against its session
digest. Then search the notebook
for what it already holds on these topics.

Done when every task has its evidence, every source was queried or its
absence noted with the gap it leaves, each session, PR, ticket and meeting in
the window is read or skipped with a reason, and every overlapping entry is
named.

## 3. Filter

Keep a candidate only if a future session without this context would act
differently for knowing it. Drop:

- status chatter, one-off fixes, and anything already in the repo's docs,
  ADRs or commit messages;
- secrets, tokens, local absolute paths, session URLs;
- people's names and personal details: describe people by team and role;
- corrections where the user changed their mind rather than the agent
  getting it wrong.

Done when each kept candidate has a one-line reason it survived.

## 4. Draft

- **Log**: one entry per task with title, summary, changes, decisions,
  learned, verification (real results, failures included), open threads and
  sources read. In `since-last`, each entry names its source project and
  opens with the date the work happened.
- **Decisions**: only choices that are hard to reverse, surprising without
  context, or a real trade-off; the alternatives and why they lost.
- **Lessons**: at most five, each titled as the claim, with confidence, the
  narrowest scope that fits, and what it supersedes.
- **Corrections**: one line each in `corrections.md`, in the format at its
  top (create it with this header if missing):

  ```markdown
  # Agent corrections
  <!-- date | repo | harness | agent did | user wanted | fix -->
  - YYYY-MM-DD | <repo> | claude | guessed a CLI flag | check --help first | open
  ```

Done when every kept candidate has a draft in its notebook file's shape.

## 5. File and report

- **session**: running `record` is the approval; file directly.
- **since-last**: show the drafts grouped by project, then what was dropped
  and why, then sources that failed or were out of scope. File only what the
  user approves, with their edits.

Update an overlapping entry rather than adding a near-duplicate; on a
collision keep the existing entry unless the draft corrects it. Commit the
notebook with one message naming the entries. Remove the
pending correction lines handled. In `since-last`, write the window end to
`~/.local/share/workbench/state/record-last-run` (create the folder if needed). If a
correction has now been seen in two or more sessions, use the `improve` skill
once the report is shown.

Report each entry with its notebook path, what was dropped and
why, and anything that failed. Done when the report is shown; in `preview`,
when the drafts are shown and nothing was filed or written.
