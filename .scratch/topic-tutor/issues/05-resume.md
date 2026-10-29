# 05 Resume across sessions, misses first

Status: done
Blocked by: 04

## What

Read `../spec.md` first.

- When `study-topic` starts on a topic that already has a folder in the
  `learning` destination, it reads `mission.md`, the lessons and `quiz.md`
  instead of starting over.
- It opens with retrieval questions on items answered wrong before, then
  writes the next lesson, numbered after the highest existing one.
- When the mission's "done" holds, it says so and offers to close the topic
  rather than writing more lessons.

## Done when

- `sh scripts/check.sh` passes.
- In a new `--plugin-dir` session on the 03 test topic, with at least one
  wrong line in `quiz.md`, the first question asked is about a missed item,
  and the next lesson is `lessons/02-<slug>.md`. Quote both.
