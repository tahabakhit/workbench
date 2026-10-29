#!/bin/sh
# Print a Claude Code or Codex session transcript as the conversation only:
# user prompts, assistant replies, files edited and git commands, each cut to
# a readable length.
set -eu
file=${1:?usage: digest.sh TRANSCRIPT.jsonl}
if head -n 1 "$file" | jq -e '.type == "session_meta"' >/dev/null 2>&1; then
	jq -r '
		def cut($n): if length > $n then .[0:$n] + " [...]" else . end;
		(.timestamp // "" | .[0:16]) as $ts
		| select(.type == "response_item") | .payload
		| if .type == "message" and .role == "user" then
			([.content[]? | .text // empty] | join(" ")) as $t
			| select($t != "" and ($t | startswith("<") | not) and ($t | startswith("# AGENTS.md") | not))
			| "\n[\($ts)] USER: \($t | cut(2000))"
		elif .type == "message" and .role == "assistant" then
			"[\($ts)] ASSISTANT: \([.content[]? | .text // empty] | join(" ") | cut(1500))"
		elif .type == "custom_tool_call" or .type == "function_call" then
			(.input // .arguments // "" | tostring) as $i
			| if ($i | test("apply_patch|\\*\\*\\* (Update|Add) File")) then
				"  edited \([$i | scan("\\*\\*\\* (?:Update|Add) File: ([^\\\\\\n\"]+)")[]] | unique | join(", "))"
			elif ($i | test("git (commit|push|merge)|gh pr (create|merge)")) then "  ran: \($i | cut(300))"
			else empty end
		else empty end
	' "$file"
	exit
fi
jq -r '
	def text: if type == "string" then . else map(select(.type == "text") | .text) | join(" ") end;
	def cut($n): if length > $n then .[0:$n] + " [...]" else . end;
	select(.isSidechain | not)
	| (.timestamp // "" | .[0:16]) as $ts
	| if .type == "user" then
		(.message.content | text) as $t
		| select($t != "" and ($t | startswith("<") | not))
		| "\n[\($ts)] USER: \($t | cut(2000))"
	elif .type == "assistant" then
		.message.content[]?
		| if .type == "text" then "[\($ts)] ASSISTANT: \(.text | cut(1500))"
		elif .type == "tool_use" and (.name == "Edit" or .name == "Write") then "  edited \(.input.file_path)"
		elif .type == "tool_use" and .name == "Bash" and (.input.command | test("git (commit|push|merge)|gh pr (create|merge)")) then "  ran: \(.input.command | cut(300))"
		elif .type == "tool_use" and .name == "Skill" then "  skill: \(.input.skill)"
		elif .type == "tool_use" and (.name | startswith("mcp__")) then "  tool: \(.name) \(.input.title // .input.query // "")"
		else empty end
	else empty end
' "$file"
