---
name: onboard
description: Use when a machine or an existing repo is being set up for agent work for the first time, or the user asks to change where work comes from. Machine scope writes the work sources file and creates the notebook; repo scope surveys the repo, proposes keep, convert or delete for its agent artifacts, and writes AGENTS.md.
---

# Onboard

Scope comes from the argument: `machine` or `repo`. Without one, run
`machine` first when `~/.local/share/workbench/sources.md` does not exist, then
`repo` when the session is in a git repository. With `preview`, show every
proposal and write nothing.

## Machine scope

1. **Inspect.** List the MCP servers (`claude mcp list`, or the harness's
   equivalent) and the CLIs on `PATH`: `gh`, `glab`, `acli`, `jira`, `slack`,
   `gcloud`. Read `~/.local/share/workbench/sources.md` if it exists. Done when
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
   it, and rekey any line that doesn't start with a contract kind. Each line names the exact command or tool call, its filters, and
   which skills it is limited to, if any. Show a diff; write only on
   approval, to `~/.local/share/workbench/sources.md`. It belongs to this
   machine: never write it into a repository or an instructions file.
   Create `~/.local/share/workbench/notebook` with `git init` if it is missing. Done when the
   approved file is written and `~/.local/share/workbench/notebook` is a git repository.

## Repo scope

1. **Survey.** Read-only, by you or a subagent, gather facts each with its source
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
2. **Ask.** One numbered round. Read `tracker-lines.md` in this skill's
   folder first; ask these, each with a recommendation worded so "yes"
   accepts it:
   1. Where do tasks for this repo live? The spec and issue files stay in
      `.scratch/<feature>/` either way. This picks who else sees them and
      where they get closed. Offer local markdown, GitHub Issues, GitLab
      Issues, or other (Jira, Linear, anything: ask for the list, create and
      close commands or tool calls).
   2. Is this repo shared with people who don't use workbench? Recommend yes
      when the last 30 commits have other authors.
   3. Do chat or code host differ from the machine's `sources.md`?

   Recommend the tracker that already holds this repo's open issues, or that
   the machine's sources name for it; with neither, local. A remote alone is
   no evidence that anyone uses its issues. Run the chosen list command once.
   Done when every question has an answer and that command has answered.
3. **Propose.** The answers become source lines in the contract's format: a
   `## Sources` section for `AGENTS.md`, or machine lines when the repo is
   shared. Any tracker but local gets a full line from `tracker-lines.md`,
   even when the machine's line already searches it. Then a table with one row per artifact found: keep, convert to
   this plugin's conventions, or delete, and a reason. Conventions:
   `AGENTS.md` under 60 lines, with `CLAUDE.md` holding only `@AGENTS.md`;
   work tracking in `.scratch/<feature>/` as the contract
   describes; glossary and ADRs in the format of `glossary-adr-format.md` in
   this skill's folder (read it before proposing a conversion); facts an
   agent can read from the code are deleted, not copied. In a shared repo,
   apply the contract's shared-repository rule: existing instruction files
   are kept and `.scratch/` goes in `.git/info/exclude`. Done when every
   survey item that is an artifact has a row.
4. **Apply on approval.** Apply only the approved rows. Draft `AGENTS.md`:
   commands first, then conventions an agent would get wrong, work tracking,
   domain docs, and the approved `## Sources` section. With an existing file, show a diff and keep what the
   user wrote unless they agree to change it. Never name a file `claude.md`
   or `agents.md` in lower case. Done when the approved changes are written
   and nothing else changed (`git status`).
5. **Prove each check bites.** Before recording a check (a lint rule, a test
   command, a hook) in `AGENTS.md`, show it passes, fails on a deliberately
   broken case, then passes once that case is reverted. Done when every
   recorded check has that pass, fail, pass, or is recorded as unproven.

## Check

Start a fresh session with a question the written file answers ("how do I
run the tests here?", "where does a new task go?") and confirm the answer comes
from it. A subagent of this session does not count: it sees the instructions
as they were when the session started. If you cannot start one, give the user
the question. Done when the answer matches the file.
