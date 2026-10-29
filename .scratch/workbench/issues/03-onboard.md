# 03 onboard skill

Status: open
Blocked by: 02

## What

`skills/onboard/`, user-invoked, one skill with two scopes. Start from
`~/dev/agents/home/dot_agents/skills/setup-repo/`, which becomes the repo
scope.

- **Machine scope**: inspect the connected tools (MCP servers, CLIs such as
  `gh`, `glab`, `acli`), ask in one numbered round with a recommendation per
  question which tracker, chat, calendar, code host and notes the user uses
  and which vault a session writes to, then write the `## Knowledge` section
  of the global instructions in the format `contract.md` defines. If the file
  is generated (it says so at the top), edit the source.
- **Repo scope**: first a read-only survey in a subagent that returns facts
  with their source: instruction files and other agent-tooling artifacts
  (skill-plugin conventions, editor rules), git remotes, recent history only
  (commit style, active branches), open issues and PRs when the code host is
  connected, existing specs, ADRs and glossary, and build and test commands
  verified by running the safe ones. Then propose per artifact: keep,
  convert to this plugin's conventions, or delete. Apply only what is
  approved. Write `AGENTS.md` under 60 lines and a repo `## Knowledge` section
  only for kinds that differ from the global one.
- Keep the `setup-ts-deep-modules` verdict line: prove a check bites (pass,
  fail, pass) before recording it.
- Copy `skills/settle/glossary-adr-format.md` unchanged into
  `skills/onboard/`. The repo scope uses it to convert an existing glossary
  or ADRs; `scripts/check.sh` fails if the copies differ.

## Done when

- `sh scripts/check.sh` exits 0.
- `.scratch/workbench/trials/onboard.md` records two runs with their output:
  machine scope on this Mac in preview mode (proposed section shown, nothing
  written) covering every kind the current global `## Knowledge` section
  has; repo scope on a work repo (it has mattpocock
  `docs/agents/` files) that proposes keep, convert or delete for each of
  them, with nothing applied.
