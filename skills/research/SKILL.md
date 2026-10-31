---
name: research
description: Use when a question needs an answer from primary sources instead of memory, such as current versions, specifications, or how an outside system behaves. Answers with a citation for every claim, names the gaps, and saves the report.
---

# Research

## 1. Question

Restate the question and split it into the parts an answer must cover. When
it is ambiguous, ask one numbered round, each question with your
recommendation.

Done when the question and its parts are written.

## 2. Gather

Answer each part under these rules: primary sources first (specs, papers,
official docs, source code), every claim cited to a source opened in this
session, gaps named rather than filled from memory, and a stop at 8 searches
and 12 page reads.

Where you see a benefit and the harness offers it, a `researcher` agent can do
this, and independent parts can go to separate agents, at most three. Doing it
in this session is equally valid.

Done when every part has cited claims or a named gap.

## 3. Check and save

Merge the reports into one, numbering sources once. Drop any claim without
a source number; a part with no claim left becomes a gap. Then add
`## Answer`: two to five sentences that answer the question from the claims,
with their source numbers.

Save it as `.scratch/<feature>/research-<slug>.md` when the question belongs
to a feature folder. Otherwise print it in full and say it was not saved.

Done when every claim in the saved or printed report carries a source
number from its `## Sources` list.
