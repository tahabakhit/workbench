# 20 Deliver the contract to Codex and Pi

Status: open

## What

Claude Code gets `contract.md` from the plugin's `SessionStart` hook. Codex
and Pi do not load the plugin, and the global instructions the dotfiles own
no longer include the contract, so their sessions run without it. Find how
each harness can load it from this repo's files without the dotfiles
including it: a native extension, a hook, or a generated instructions file
that workbench owns.

## Done when

- A new Codex session and a new Pi session each answer "where does a
  decision made outside a feature folder go?" with `decisions/` in the
  notebook.

## Out of scope

- The text of the global instructions (dotfiles).
