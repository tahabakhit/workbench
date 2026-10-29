# 14 The plugin delivers the contract and Knowledge itself

Status: open
Blocked by: 08

## What

- Knowledge destinations move to `~/.config/workbench/knowledge.md`, one per
  machine and in no repo. It holds the `## Knowledge` lines in the format
  the contract describes. A repo's `AGENTS.md` can still override them kind
  by kind.
- A `SessionStart` hook in the plugin adds `contract.md` and the knowledge
  file to every Claude Code session's context. When the file is missing it
  says so and points to `onboard`. Check the installed Claude Code's hook
  docs for how a plugin hook adds context before relying on it.
- `onboard` machine scope writes the knowledge file instead of a section in
  the global instructions. `record`, `plan-day`, `improve` and the contract
  read destinations from the hook's context or the file, not from the global
  instructions.
- Codex and Pi get no skills from this repo until their native extensions
  exist. They keep the contract through one include line in their global
  `AGENTS.md`, which the dotfiles own.
- On Fedora, move the current Knowledge lines from the dotfiles template into
  the knowledge file (dotfiles side: mac-dotfiles issue 11).

## Done when

- A fresh Claude Code session started outside any repo, with the global
  `AGENTS.md` holding no contract and no Knowledge section, answers "where do
  lessons go?" with the knowledge file's destination and "what line marks an
  issue's state?" with `Status:`.
- With the plugin disabled, the same session answers neither from workbench.
- `sh scripts/check.sh` passes.
- `rg -n 'global instructions' skills contract.md` shows no remaining claim
  that Knowledge lives in the global instructions.

## Out of scope

- The Pi extension and the Codex plugin.
- `settings.json` entries for the plugin (mac-dotfiles issue 04).
