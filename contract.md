# Workbench contract

Rules shared by every machine.

## Work tracking

- Work that spans sessions lives in `.scratch/<feature>/`: `spec.md` with a
  task table and a `## Decisions` section, and one `issues/NN-slug.md` per
  task. Each issue has a plain `Status: open|done|wontfix` line (not bold),
  an optional `Blocked by: NN` line, and a "Done when" check.
- An issue set to `wontfix` takes the open issues it blocks with it in the
  same change: each becomes `wontfix` too or names another blocker.
- Ad hoc work gets no feature folder and no notes while it happens.
- Finish at the done-when: run each check and quote its output. Before
  closing, diff the result against the done-when and Out of scope in a fresh
  subagent. Then the session that met it sets `Status: done` in the issue and
  in the spec table.
- Write the failing test first. Mock only at system boundaries.
- Tracker updates go to the trackers named in the `sources` entries, with
  evidence, applied only on approval.
- Decisions from a decision session are written only on approval: to
  `## Decisions` in `spec.md` when a feature folder exists, otherwise to
  `decisions/` in the notebook.
- Write an ADR in `docs/adr/` only when all three hold: hard to reverse,
  surprising without context, a real trade-off.

## Notebook

Records live in the notebook: plain markdown in a git repository at
`~/.agents/notebook`, one folder per kind, the same on every machine. The plugin's
`notebook` tool will search and file into it; until it does, skills write the
files directly. Create the repository with `git init` and a folder with its
first file, and commit each write there.

| Path | Holds |
|---|---|
| `log/YYYY-MM-DD-<repo>-<slug>.md` | one finished task |
| `decisions/YYYY-MM-DD-<slug>.md` | a decision made outside a feature folder, with the alternatives rejected |
| `lessons/<slug>.md` | a lesson that would help in another repository |
| `corrections.md` | the ledger of corrections the user made |
| `plans/YYYY-MM-DD.md` | the daily plan |
| `goals.md` | long-term goals the plan ranks against; changed only when the user asks |
| `learning/<topic>/` | a topic being studied: `mission.md`, `sources.md`, `lessons/`, `quiz.md` |

Each note opens with frontmatter holding `date`, plus `repo` when it came
from one; a lesson also holds `scope` and, when it replaces one, `supersedes`.

## Sources

Where work comes from (tracker, chat, calendar, meeting notes, code host)
differs by person, so each machine lists it in
`~/.config/workbench/sources.md`, which `onboard` writes: one line per
source, `- <tool>: <exact command or tool call, with filters and limits>`.
A skill that gathers work reads that file; without it, the skill says so and
gathers from git and `.scratch/` alone.

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

Keep the whole body under about 15 lines, one line per point. Leave out
anything the diff already shows (file lists, change-by-change summaries),
rejected alternatives, and the story of how you got there. A reviewer should
read it in under a minute.
