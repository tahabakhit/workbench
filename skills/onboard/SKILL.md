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
2. **Ask.** One round using the contract's question format, holding only the
   kinds with a decision to make. A kind is asked when the inventory has a
   candidate for it, or its existing line failed its test call; when the user
   asked to change one kind, ask only that one. For each, name what reads it (`tracker`:
   `plan-day`, `record`, `write-spec`; the others: `plan-day`, `record`), give
   your recommendation and one line on why; word it so "yes" accepts it. Ask
   also for any limit on a source, such as content that must not be stored.
   Do not ask about the rest: after the round, list existing lines that
   answered as kept, and kinds with no tool as left out, with "name one if
   you use it". Done when every kind has a source or is left out, and the
   list of left-out kinds has been shown.
3. **Write.** Draft the whole file in the format of the contract's Sources
   section and print it verbatim in a code block, even when it is identical
   to the current one. Keep every existing line unless an answer changes
   it, and rekey any line that doesn't start with a contract kind. Each line names the exact command or tool call, its filters, and
   which skills it is limited to, if any. Show a diff; write only on
   approval, to `~/.local/share/workbench/sources.md`. It belongs to this
   machine: never write it into a repository or an instructions file.
   Create `~/.local/share/workbench/notebook` with `git init` if it is missing; every
   later note is filed by `notebook file`. Done when the
   approved file is written and `~/.local/share/workbench/notebook` is a git repository.

## Repo scope

1. **Survey.** Read-only, by you or a subagent, gather facts each with its source
   (path, command or URL):
   - instruction files and agent tooling: `AGENTS.md`, `CLAUDE.md`,
     `CLAUDE.local.md`, `.cursor/`, `.github/copilot-instructions.md`,
     `docs/agents/`, skill or plugin folders, hooks;
   - git remotes, and from the last 30 commits the commit style and the
     active branches;
   - who can read the repo: the output of `scripts/repo-audience.sh` in this
     skill's folder (visibility, other authors, placement, and any personal
     configuration already tracked);
   - open issues and PRs, only if the code host's CLI is connected;
   - specs, plans, `.scratch/`, `docs/adr/`, `GLOSSARY.md` or `CONTEXT.md`;
   - build, test and lint commands, with the result of running the safe
     ones (no deploys, no writes outside the repo, no paid calls).

   Done when every item has facts or "none found", and every command is
   marked run with its result or not run with the reason. The survey is
   working material: later steps show the user what changes a proposal or
   needs a decision, plus failed or unrun checks with reasons, and offer the
   routine detail.
2. **Ask.** One round using the contract's question format. Read
   `tracker-lines.md` in this skill's folder first. Ask only what applies,
   each with one line on what the answer changes and a recommendation worded
   so "yes" accepts it; state the rest as findings after the round:
   1. Where do tasks for this repo live? The spec and issue files stay in
      `.scratch/<feature>/` either way. This picks who else sees them and
      where they get closed. Offer local markdown, GitHub Issues, GitLab
      Issues, or other (Jira, Linear, anything: ask for the list, create and
      close commands or tool calls). On a public or soon-public repo, say that
      `write-spec` copies issue text there for anyone to read. When the repo's
      `## Sources` or the machine's `tracker (<repo>)` line already names a
      tracker and its list command answers, keep it as a finding.
   2. With placement `ask`: will this repo go public, or be read by people
      who don't use workbench? A yes keeps personal configuration out of git.
      Recommend yes when visibility is `none` or `unknown`, or the repo has a
      licence file; otherwise no. With `shared` or `public` the repo has
      answered: state it as a finding.
   3. Does this repo use another chat than the machine's? Ask only when the
      machine has a `chat` line; a different answer gets a line for this repo.

   Without that `chat` line, chat is a finding too. Code host is never asked:
   when the remote's host is not the one in the machine's `host` line, propose
   a line for this repo. Calendar and notes belong to the machine, so they
   have no question here.

   Recommend the tracker that already holds this repo's open issues, or that
   the machine's sources name for it; with neither, local. A remote alone is
   no evidence that anyone uses its issues. Run the chosen list command once.
   Done when every asked question has an answer, every skipped one is a
   finding, and that command has answered.
3. **Propose.** The answers become source lines in the contract's format: a
   `## Sources` section for `AGENTS.md` in a repo only its owner reads now and
   later, otherwise machine lines keyed `<kind> (<repo>)`. Any tracker but
   local gets a full line from `tracker-lines.md`,
   even when the machine's line already searches it. Then a table with one row per artifact found: keep, convert to
   this plugin's conventions, or delete, and a reason. Conventions:
   `AGENTS.md` under 60 lines, with `CLAUDE.md` holding only `@AGENTS.md`;
   work tracking in `.scratch/<feature>/` as the contract
   describes; glossary and ADRs in the format of `glossary-adr-format.md` in
   this skill's folder (read it before proposing a conversion); facts an
   agent can read from the code are deleted, not copied. Apply the contract's
   readable-repository rule when others can read the repo now or later:
   `.scratch/` goes in `.git/info/exclude`. In a shared repo, existing
   instruction files are kept. In a solo public or soon-public repo,
   `AGENTS.md` can be written, for contributors: no work tracking, sources,
   notebook or home-folder paths. Anything `repo-audience.sh` reports as
   tracked gets a row: move it out of git (`git rm --cached` keeps the local
   files), or keep it when the owner works in the open on purpose. Either way
   say that git history still holds it, and that rewriting history is a
   separate decision for the owner, before the repo goes public. Done when
   every survey item that is an artifact has a row.
4. **Apply on approval.** Apply only the approved rows. Add approved machine
   lines to `~/.local/share/workbench/sources.md`, keeping its other lines.
   Draft `AGENTS.md`:
   commands first, then conventions an agent would get wrong, work tracking,
   domain docs, and the approved `## Sources` section, leaving out what step 3
   keeps out of a public repo. With an existing file, show a diff and keep what the
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
the question. When personal configuration stayed out of git, also run
`git check-ignore .scratch` and confirm it prints the path. Done when the
answer matches the file.
