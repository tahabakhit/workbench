#!/bin/sh
# Reads {"question", "answer", "distractors": [...]} and prints
# {"question", "options", "answer_index"} with the options in an order taken
# from /dev/urandom, because a model placing answers puts the right one first.
set -eu

input=$(cat)
n=$(printf '%s\n' "$input" | jq -e '
  if (.question | type) == "string" and (.answer | type) == "string"
     and (.distractors | type) == "array" and (.distractors | length) > 0
  then .distractors | length + 1
  else error("need question, answer and at least one distractor") end')
keys=$(od -An -N"$((n * 4))" -tu4 /dev/urandom | tr -s ' \n' '  ')

printf '%s\n' "$input" | jq -c --arg keys "$keys" '
  . as $q
  | ($keys | [splits(" +") | select(. != "") | tonumber]) as $k
  | ([$q.answer] + $q.distractors) as $opts
  | [range(0; $opts | length) | {option: $opts[.], from: ., key: $k[.]}]
  | sort_by(.key) as $sorted
  | {question: $q.question, options: [$sorted[].option],
     answer_index: ($sorted | map(.from) | index(0))}'
