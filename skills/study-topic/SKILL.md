---
name: study-topic
description: Use when the user wants to learn a topic over several sessions. Settles a mission, gathers cited primary sources and names their gaps, writes Markdown lessons, and quizzes with shuffled answers, all in the notebook's `learning/` folder.
---

# Study topic

A tutor, not an encyclopedia. Each session moves the user one lesson closer
to their mission, and every claim traces to a source they can open.

## 1. Topic folder

The topic folder is `~/.local/share/workbench/notebook/learning/<topic-slug>/`; create it, and
never overwrite a file already in it. Commit the notebook after each step
that writes. If it already holds `mission.md`, this is a returning
session: read `mission.md`, `sources.md`, every lesson and `quiz.md`, then
go to step 6.

Done when the topic folder exists, and for a returning session its files
are read.

## 2. Mission

Ask one numbered round, each question with your recommendation: why this
topic, what the user can already do, what "done" looks like as something
they can do or explain, and what is out of scope. Write `mission.md` with
`## Why`, `## Starting point`, `## Done when` and `## Out of scope`, print it
verbatim, and change it until the user approves.

Done when the user has approved `mission.md`.

## 3. Sources

Find which primary sources cover the topic up to the mission's done-when, and
what each one covers: primary sources first (specs, papers, official docs,
source code), every claim cited to a source opened in this session, gaps named
rather than filled from memory, and a stop at 8 searches and 12 page reads.
Where the harness offers it, a `researcher` agent can do this instead.

Write `sources.md`: a numbered list of title, link, kind and what it covers
for this mission, then `## Gaps` with each part of the mission no source
covers.

Done when every line of the mission's done-when maps to a source number or
a gap.

## 4. Lesson

Take the first step toward the mission's done-when and write
`lessons/01-<slug>.md`:

- `## Goal`: one sentence on what the user can do after it.
- The body grounds before it generalises: a concrete case (an example, a
  command, a worked problem), then the rule it shows, then a second case
  that applies the rule.
- Every claim ends with `[n]`, a number from `sources.md`. A claim that needs
  a new source gets it added to `sources.md` first; a claim no source
  supports stays out.
- Readable in about ten minutes.

Print the lesson verbatim, not a summary. The lesson holds no questions or
answers; the quiz comes next.

Done when the lesson file exists, every claim in it carries a source
number, and the user has seen it in full.

## 5. Quiz

Write three to five retrieval questions on the lesson, each asking the user
to predict or recall rather than recognise, with one right answer and three
plausible wrong ones. Pass each as JSON, `{"question", "answer",
"distractors": [...]}`, to `scripts/shuffle.sh` in this skill's folder and
show the options lettered in exactly the order it returns; the script, not
you, decides where the right answer sits.

Ask one question at a time and reveal the answer, with its source number,
only after the user replies. After each reply append one line to `quiz.md`
in the topic folder, with the date from `date`:

`- YYYY-MM-DD | lesson NN | <question> | right` (or `wrong`)

Done when every question has an answer line in `quiz.md`.

## 6. Returning session

A miss is a question whose latest line in `quiz.md` says `wrong`. Before
anything new, ask each miss again as in step 5: reworded, through
`scripts/shuffle.sh`, with a new line in `quiz.md` per answer.

Then hold the mission's done-when against the lessons and the latest quiz
lines. When every item is taught, quizzed and its latest answers are
right, say so and offer to close the topic, adding `## Closed` with the date
to `mission.md`. Otherwise, or when the user declines, write the next
lesson, numbered after the highest in `lessons/`, as in step 4, and quiz it
as in step 5.

Done when every miss has a new answer line, and either the next lesson and
its quiz are done or the user has been offered to close the topic.
