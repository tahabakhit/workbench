#!/bin/sh
# Checks that shuffle.sh keeps every option, points answer_index at the
# answer, and spreads the answer across positions over many runs.
set -u

shuffle="$(cd "$(dirname "$0")" && pwd)/shuffle.sh"
question='{"question": "What does ${v:-w} give when v is set but empty?", "answer": "w", "distractors": ["an empty string", "an error", "v"]}'
runs=30
fail=0

# One round: 30 shuffles, each checked for lost options; prints the count per
# position and succeeds when every position got 1 to 15 answers.
round() {
  counts=
  i=0
  while [ "$i" -lt "$runs" ]; do
    out=$(printf '%s\n' "$question" | sh "$shuffle") || { echo "FAIL  shuffle.sh exited non-zero"; exit 1; }
    if ! printf '%s\n' "$out" | jq -e --argjson q "$question" '
        (.options | sort) == ([$q.answer] + $q.distractors | sort)
        and .options[.answer_index] == $q.answer
        and .question == $q.question' >/dev/null; then
      echo "FAIL  run $i lost an option or misplaced answer_index: $out"
      exit 1
    fi
    counts="$counts $(printf '%s\n' "$out" | jq '.answer_index')"
    i=$((i + 1))
  done
  spread=ok
  report=
  for pos in 0 1 2 3; do
    n=$(printf '%s\n' $counts | grep -cx "$pos")
    report="$report
answer in position $pos: $n of $runs"
    [ "$n" -ge 1 ] && [ "$n" -le 15 ] || spread=bad
  done
  [ "$spread" = ok ]
}

# A fair shuffle misses the bounds in about 1 round in 250, so it gets three
# rounds; a biased one fails all three.
spread_ok=
for attempt in 1 2 3; do
  if round; then spread_ok=1; break; fi
done
if [ -n "$spread_ok" ]; then
  printf '%s\n' "$report" | sed '/^$/d; s/^/ok    /'
else
  printf '%s\n' "$report" | sed '/^$/d; s/^/FAIL  /'
  fail=1
fi

printf '%s\n' '{"question": "q", "answer": "a"}' | sh "$shuffle" >/dev/null 2>&1 &&
  { echo "FAIL  a question with no distractors was accepted"; fail=1; } ||
  echo "ok    a question with no distractors is rejected"

if [ "$fail" -eq 0 ]; then echo "shuffle tests: passed"; else echo "shuffle tests: failed"; fi
[ "$fail" -eq 0 ]
