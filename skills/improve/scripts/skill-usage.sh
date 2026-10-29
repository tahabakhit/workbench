#!/bin/sh
# Count how often each skill in ~/.agents/skills was used in Claude Code
# sessions from the last DAYS days (default 60), by the model or as a slash
# command. Skills with zero uses are listed too. Codex logs no skill use.
set -eu
days=${1:-60}
since=$(date -v-"${days}"d +%Y-%m-%d 2>/dev/null || date -d "-${days} days" +%Y-%m-%d)
used=$(mktemp); trap 'rm -f "$used"' EXIT
find "$HOME/.claude/projects" -name '*.jsonl' -newermt "$since" -print0 2>/dev/null |
	xargs -0 rg -o -N --no-filename \
		-e '"name":"Skill","input":\{"skill":"[^"]*"' \
		-e '<command-name>/[a-z0-9:_-]+</command-name>' 2>/dev/null |
	sed -E 's/.*"skill":"([^"]*)".*/\1/; s#<command-name>/([^<]*)</command-name>#\1#' >"$used" || true
for dir in "$HOME/.agents/skills"/*/; do
	name=$(basename "$dir")
	n=$(grep -cx "$name" "$used" || true)
	printf '%s\t%s\n' "$n" "$name"
done | sort -n
