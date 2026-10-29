---
name: study-topic
description: Tutor the user through one topic across sessions. Settle a mission, gather cited primary sources and name their gaps, and write Markdown lessons into the `learning` destination. Run with the topic to study.
disable-model-invocation: true
---

# Study topic

A tutor, not an encyclopedia. Each session moves the user one lesson closer
to their mission, and every claim traces to a source they can open.

## 1. Topic folder

Read the `learning:` entry in the `## Knowledge` sections (repo first, then
the machine's). With none, stop and point the user to `onboard`. The topic
folder is `<destination>/<topic-slug>/`; create it, and never overwrite a
file already in it.

Done when the topic folder exists.

## 2. Mission

Ask one numbered round, each question with your recommendation: why this
topic, what the user can already do, what "done" looks like as something
they can do or explain, and what is out of scope. Write `mission.md` with
`## Why`, `## Starting point`, `## Done when` and `## Out of scope`, print it
verbatim, and change it until the user approves.

Done when the user has approved `mission.md`.

## 3. Sources

Hand the `researcher` agent the mission and one question: which primary
sources cover the topic up to the mission's done-when, and what each one
covers. Where plugin agents don't exist, do that research yourself under the
agent's rules: primary sources first (specs, papers, official docs, source
code), every claim cited to a source opened in this session, gaps named
rather than filled from memory, and a stop at 8 searches and 12 page reads.

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

Print the lesson verbatim, not a summary.

Done when the lesson file exists, every claim in it carries a source
number, and the user has seen it in full.
