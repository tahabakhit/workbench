# topic-tutor

## Goal

A `study-topic` skill tutors the owner through one topic over many
sessions. It sets a mission, gathers cited primary sources and names their
gaps, writes Markdown lessons, quizzes with retrieval questions whose answer
order comes from code, and keeps learning records so the next session picks
up where the last one stopped, misses first. A shared `researcher` agent
does the source work.

Settled input: `.scratch/bootstrap/decisions.md` (wins on conflict), the
`teach` and `research` rows of `.scratch/bootstrap/mattpocock-analysis.md`,
and the `## Decisions` below.

## Principles

- Plugin rules from `AGENTS.md` hold: the skill name is a plain verb, the
  skill never loads another skill, and scripts are POSIX `sh` with `jq`,
  executable, under the skill's `scripts/`.
- Records and lessons go to the `learning` destination from `## Knowledge`.
  No path is hardcoded. With no `learning` destination the skill stops and
  points the user to `onboard`.
- Every factual claim in a source list or lesson carries a citation to a
  primary source. A gap the sources don't cover is named as a gap, not
  filled from memory.
- Answer order in a quiz comes from `scripts/shuffle.sh`, never from the
  model. Instruction-only fixes put the right answer first 33 of 33 times in
  the analysis trial.
- The `researcher` agent is a Claude Code plugin agent. The skill also works
  without it (Codex, Pi): it then does the research inline under the same
  citation and stopping rules.
- `learn`, the owner's earlier tutor (a separate repo, not on every
  machine), has no licence. Read it for ideas only; copy no text or code
  from it.
- Lessons follow a grounding order: the concrete case before the general
  rule.

## Tasks

| # | Task | Blocked by | Status |
|---|---|---|---|
| 01 | Add the `learning` Knowledge kind | | done |
| 02 | `researcher` agent: primary sources, a citation per claim, a stopping rule | 01 | done |
| 03 | `study-topic` core: mission, sources with gaps, first lesson | 01, 02 | open |
| 04 | Retrieval quiz with answers shuffled in code | 03 | open |
| 05 | Resume across sessions, misses first | 04 | open |
| 06 | Two-topic trial | 05 | open |

## Out of scope

- HTML lessons and rendered pages. Lessons are Markdown so every harness can
  read them.
- The `research` skill itself. It gets its own spec later and reuses the
  agent from 02.
- A `researcher` agent for Codex or Pi; the skill's inline fallback covers
  them.
- Spaced-repetition scheduling beyond "earlier misses first".
- Copying from the `learn` repo, which has no licence.

## Decisions

| Decision | Reason |
|---|---|
| Six issues: the `learning` kind, the agent, the core skill, the quiz, resume, a trial | Each slice is checkable on its own and fits one session |
| The shared `researcher` agent is specified here, not in its own `research` spec | The tutor needs it first; the research skill reuses it later |
| Learning records get a new `learning` Knowledge kind in `contract.md` | Study records are not cross-repo lessons and should not mix with them |
| Lessons are Markdown, not HTML | Codex and Pi read them too; the vault holds them |
| The skill is named `study-topic` | A plain verb naming the outcome, and it avoids clashing with the unlicensed `learn` repo |
