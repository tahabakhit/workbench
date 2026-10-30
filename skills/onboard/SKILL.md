---
name: onboard
description: Set up a machine or a repo for agent work. Machine scope writes the machine's work sources file from the tools found and creates the notebook; repo scope surveys an existing repo, proposes keep, convert or delete for its agent artifacts, and writes AGENTS.md.
disable-model-invocation: true
---

# Onboard

Scope comes from the argument: `machine` or `repo`. Without one, run
`machine` when `~/.config/workbench/sources.md` does not exist, otherwise
`repo`. With `preview`, show every proposal and write nothing.

## Machine scope

1. **Inspect.** List the MCP servers (`claude mcp list`, or the harness's
   equivalent) and the CLIs on `PATH`: `gh`, `glab`, `acli`, `jira`, `slack`,
   `gcloud`. Read `~/.config/workbench/sources.md` if it exists. Done when
   you have a table of tool, what it serves (tracker, chat, calendar,
   meeting notes, code host), and whether it answers a test call (one cheap
   read each).
2. **Ask.** One numbered round. For each of tracker, chat, calendar, meeting
   notes and code host, give your recommendation from the inventory and one
   line on why; word it so "yes" accepts it. Ask also for any limit on a
   source, such as content that must not be stored. Done when every kind has
   a source or is deliberately left out.
3. **Write.** Draft the whole file in the format of the contract's Sources
   section and print it verbatim in a code block, even when it is identical
   to the current one. Keep every existing line unless an answer changes
   it. Each line names the exact command or tool call, its filters, and
   which skills it is limited to, if any. Show a diff; write only on
   approval, to `~/.config/workbench/sources.md`. It belongs to this
   machine: never write it into a repository or an instructions file.
   Create `~/notebook` with `git init` if it is missing. Done when the
   approved file is written and `~/notebook` is a git repository.

## Repo scope

1. **Survey.** In a read-only subagent, return facts each with its source
   (path, command or URL):
   - instruction files and agent tooling: `AGENTS.md`, `CLAUDE.md`,
     `CLAUDE.local.md`, `.cursor/`, `.github/copilot-instructions.md`,
     `docs/agents/`, skill or plugin folders, hooks;
   - git remotes, and from the last 30 commits the commit style and the
     active branches;
   - open issues and PRs, only if the code host's CLI is connected;
   - specs, plans, `.scratch/`, `docs/adr/`, `GLOSSARY.md` or `CONTEXT.md`;
   - build, test and lint commands, with the result of running the safe
     ones (no deploys, no writes outside the repo, no paid calls).

   Done when every item has facts or "none found", and every command is
   marked run with its result or not run with the reason.
2. **Propose.** A table with one row per artifact found: keep, convert to
   this plugin's conventions, or delete, and a reason. Conventions:
   `AGENTS.md` under 60 lines, with `CLAUDE.md` holding only `@AGENTS.md`;
   work tracking in `.scratch/<feature>/` as the contract
   describes; glossary and ADRs in the format of `glossary-adr-format.md` in
   this skill's folder (read it before proposing a conversion); facts an
   agent can read from the code are deleted, not copied. Done when every
   survey item that is an artifact has a row.
3. **Apply on approval.** Apply only the approved rows. Draft `AGENTS.md`:
   commands first, then conventions an agent would get wrong, work tracking,
   domain docs. With an existing file, show a diff and keep what the
   user wrote unless they agree to change it. Never name a file `claude.md`
   or `agents.md` in lower case. Done when the approved changes are written
   and nothing else changed (`git status`).
4. **Prove each check bites.** Before recording a check (a lint rule, a test
   command, a hook) in `AGENTS.md`, show it passes, fails on a deliberately
   broken case, then passes once that case is reverted. Done when every
   recorded check has that pass, fail, pass, or is recorded as unproven.

## Check

Start a fresh session with a question the written file answers ("how do I
run the tests here?", "where does a new task go?") and confirm the answer comes
from it. A subagent of this session does not count: it sees the instructions
as they were when the session started. If you cannot start one, give the user
the question. Done when the answer matches the file.
