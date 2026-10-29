# 04 Retrieval quiz with answers shuffled in code

Status: open
Blocked by: 03

## What

Read `../spec.md` first.

- Add `skills/study-topic/scripts/shuffle.sh`. It reads one question as
  JSON on stdin, `{"question", "answer", "distractors": [...]}`, and prints
  `{"question", "options": [...], "answer_index"}` with the options in a
  random order taken from `/dev/urandom`, not from the model.
- Write `skills/study-topic/scripts/test-shuffle.sh` first and make
  `scripts/check.sh` run it.
- `study-topic` ends each lesson with a retrieval quiz: the model writes the
  questions, `shuffle.sh` orders the options, and each answer is appended to
  `quiz.md` in the topic folder as one line: date, lesson, question, right or
  wrong.

## Done when

- `sh scripts/check.sh` passes and prints the shuffle test results, which
  include: over 30 runs of a four-option question, the right answer lands in
  each of the four positions at least once and in no position more than 15
  times.
- In a `--plugin-dir` session, finishing lesson 01 of the 03 test topic
  shows a quiz, and `quiz.md` gains one line per answer. Quote the lines.
