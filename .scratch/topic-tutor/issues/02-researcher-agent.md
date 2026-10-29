# 02 `researcher` agent: primary sources, a citation per claim, a stopping rule

Status: done
Blocked by: 01

## What

Read `../spec.md` first, and the `research` row of
`.scratch/bootstrap/mattpocock-analysis.md`.

- Check the installed Claude Code's plugin docs for the agent file format
  and where plugin agents live before writing anything.
- Add `agents/researcher.md` to the plugin. It answers one research
  question with primary sources (specs, papers, official docs, source code)
  and gives a citation (URL or path) for every claim. It does the work
  itself and never delegates to another agent. It stops when every part of
  the question has a cited answer or a named gap, or after a fixed budget of
  searches it states, whichever comes first. It returns claims, citations
  and gaps, and writes nothing.
- No skill change in this issue.

## Done when

- `sh scripts/check.sh` passes, with the agent picked up by `claude plugin
  validate --strict`.
- A run of `claude -p --plugin-dir ~/dev/workbench --output-format
  stream-json` that asks the `researcher` agent one test question returns
  only cited claims plus a gaps list, and `jq` over the stream shows no
  `Agent` tool call from inside the agent. Quote the question, the answer
  and the `jq` count.
