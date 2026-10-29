#!/bin/sh
# List Claude Code and interactive Codex sessions with activity at or after
# SINCE (YYYY-MM-DD or an ISO timestamp, UTC), one JSON object per line.
# Skips subagent transcripts, scripted `codex exec` runs and sessions run in
# temporary folders. `recorded` is true when the session already ran record
# or retro.
set -eu
since=${1:?usage: sessions.sh SINCE}
tmp='^/(private/)?(tmp|var)(/|$)'
if [ -d "$HOME/.claude/projects" ]; then
	find "$HOME/.claude/projects" -mindepth 2 -maxdepth 2 -name '*.jsonl' -newermt "$since" |
	while IFS= read -r f; do
		jq -n -c --arg since "$since" --arg file "$f" --arg tmp "$tmp" '
			def text: if type == "string" then . else map(select(.type == "text") | .text) | join(" ") end;
			[inputs] as $r
			| [$r[] | select(.type == "user" and (.isSidechain | not))
				| (.message.content | text) as $t
				| select($t != "" and ($t | startswith("<") | not))
				| {t: $t, cwd, branch: .gitBranch}] as $p
			| [$r[] | .timestamp // empty] as $ts
			| select(($p | length) > 0 and ($ts | max) >= $since)
			| select($p[0].cwd | test($tmp) | not)
			| {
				harness: "claude",
				file: $file,
				cwd: $p[0].cwd,
				branch: $p[0].branch,
				start: ($ts | min),
				end: ($ts | max),
				prompts: ($p | length),
				title: ([$r[] | .customTitle // .aiTitle // empty] | last),
				first: ($p[0].t | .[0:160]),
				recorded: ([$r[] | select(.type == "assistant") | .message.content[]?
					| select(.type == "tool_use" and (
						(.name == "Skill" and (.input.skill | test("^(workbench:)?(record|retro)$")))
						or (.name == "mcp__om__record_work" and .input.dry_run != true)))]
					+ [$r[] | select(.type == "user") | (.message.content | text)
						| select(test("<command-name>/(workbench:)?(record|retro)</command-name>"))]
					| length > 0)
			}' "$f" || echo "sessions.sh: skipped unreadable $f" >&2
	done
fi
if [ -d "$HOME/.codex/sessions" ]; then
	find "$HOME/.codex/sessions" -name '*.jsonl' -newermt "$since" |
	while IFS= read -r f; do
		jq -n -c --arg since "$since" --arg file "$f" --arg tmp "$tmp" '
			[inputs] as $r
			| ([$r[] | select(.type == "session_meta") | .payload] | first) as $m
			| select($m != null and $m.originator != "codex_exec")
			| [$r[] | select(.type == "response_item" and .payload.type == "message"
					and .payload.role == "user")
				| [.payload.content[]? | .text // empty] | join(" ")
				| select(. != "" and (startswith("<") | not) and (startswith("# AGENTS.md") | not))] as $p
			| [$r[] | .timestamp // empty] as $ts
			| select(($p | length) > 0 and ($ts | max) >= $since)
			| select($m.cwd | test($tmp) | not)
			| {
				harness: "codex",
				file: $file,
				cwd: $m.cwd,
				branch: $m.git.branch,
				start: ($ts | min),
				end: ($ts | max),
				prompts: ($p | length),
				title: null,
				first: ($p[0] | .[0:160]),
				recorded: ([$p[] | select(test("(^|\\s)[$/](record|retro)(\\s|$)"))] | length > 0)
			}' "$f" || echo "sessions.sh: skipped unreadable $f" >&2
	done
fi
