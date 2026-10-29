# Workbench contract

Rules shared by every machine. A repository's `AGENTS.md` adds to them and
wins on conflict.

## Work tracking

- Work that spans sessions lives in `.scratch/<feature>/`: `spec.md` with a
  task table and a `## Decisions` section, and one `issues/NN-slug.md` per
  task. Each issue has a plain `Status: open|done|wontfix` line (not bold),
  an optional `Blocked by: NN` line, and a "Done when" check.
- Ad hoc work records nothing while it happens; filing reads the session
  transcript and git log afterwards.
- Finish at the done-when: run each check and quote its output. Before
  closing, diff the result against the done-when and Out of scope in a fresh
  subagent. Then the session that met it sets `Status: done` in the issue and
  in the spec table.
- Write the failing test first. Mock only at system boundaries.
- Tracker updates go to the trackers the `## Knowledge` sections declare,
  with evidence, applied only on approval.
- Decisions from a decision session are written only on approval: to
  `## Decisions` in `spec.md` when a feature folder exists, otherwise to the
  `decisions:` destination.
- Write an ADR in `docs/adr/` only when all three hold: hard to reverse,
  surprising without context, a real trade-off.

## Knowledge destinations

Skills that file or gather things read where to do it from `## Knowledge`
sections: one in each machine's global instructions, and optionally one in a
repository's `AGENTS.md` that overrides it kind by kind. Each entry is one
line, `- <kind>: <where>, <how>`. A kind can be limited to some skills, as in
`sources (for plan-day)`.

| Kind | Holds |
|---|---|
| `log` | records of finished work |
| `decisions` | decisions made outside a feature folder |
| `lessons` | lessons that would help in another repository |
| `corrections` | the ledger of corrections the user made |
| `plan` | the daily plan |
| `goals` | long-term goals the plan ranks against |
| `sources` | where work is gathered from: tracker, chat, calendar, notes, code host |

- The vault a session writes to is the one the machine's section names; a
  repository's section can name another.
- A skill that needs a destination no section declares stops and points the
  user to `onboard`. Any other skill carries on and says which fallback it
  used.

## Writing that lands in a repository

Committed text must make sense to someone who wasn't there. Before committing,
reread comments, docs, and PR bodies and remove:

- dates, run or job numbers, ticket IDs, and people's names (a ticket ID in a
  commit subject or PR title is fine);
- history: "used to", "no longer", "the old X". Describe what is true now; the
  past belongs in the commit message;
- narration of the investigation or of what a reviewer should notice.

A comment states a non-obvious constraint in one or two sentences, then stops.
When unsure whether a comment is needed, it isn't.

A PR body has three short parts. Summary: the smallest view that makes the
point, such as pseudocode, a call tree, a file tree or a diff sketch.
Evidence: before and after, as command output, a test run or a screenshot.
Merge risk: a one-way or two-way door, and the blast radius in a few words.
