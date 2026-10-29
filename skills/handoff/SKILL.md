---
name: handoff
description: Write a handoff so a fresh session, in Claude Code, Codex or Pi or on another machine, can continue this work, and optionally start it in the background. Run with what the next session should do.
disable-model-invocation: true
---

# Handoff

The next session starts cold. It needs where things stand, what is proven,
and where to read, not a retelling.

## 1. Gather

The argument says what the next session does; without one, ask. Run `git
status`, `git branch --show-current` and `git log --oneline -5`, and read the
feature's `spec.md` and its open issues when a feature folder holds this
work.

Done when the next session's goal is one sentence.

## 2. Write

Write `.scratch/<feature>/handoff.md` when a feature folder holds this work,
replacing any earlier handoff there. Otherwise write
`${TMPDIR:-/tmp}/handoff-<slug>.md`, which stays on this machine.

Sections:

- `## Goal`: the sentence from step 1.
- `## State`: branch, last commit, uncommitted changes, what is done.
- `## Claims`: each fact the next session relies on, labelled `verified`
  with the command or file that shows it, or `assumed`.
- `## Next`: the next steps in order, each with its done-when.
- `## Read first`: paths or URLs of the spec, issues, files and commits.
  Point to them; never copy their content.
- `## Skills`: the slash commands the next session should run.

Leave out secrets, tokens and personal data.

Done when every claim is labelled and the file names its artifacts by path
instead of repeating them.

## 3. Pass it on

A handoff in a feature folder reaches another machine only once it is
committed and pushed: commit it, and push when the next session runs
elsewhere.

When the user asks to start the next session now in Claude Code, run
`claude --bg --name "<goal in a few words>" -- "$(cat <handoff file>)"`;
passing the file keeps the shell from expanding its text. For Codex or Pi,
give the user the command that starts it with the file as the first prompt.

Done when the user has the file's path and, if started, `claude agents`
lists the session.
