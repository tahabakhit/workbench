# 20 Deliver the contract to Codex and Pi

Status: open

## What

Claude Code gets `contract.md` from the plugin's `SessionStart` hook. Codex
and Pi do not load the plugin, and the global instructions the dotfiles own
no longer include the contract, so their sessions run without it. Find how
each harness can load it from this repo's files without the dotfiles
including it: a native extension, a hook, or a generated instructions file
that workbench owns.

`codex plugin marketplace add` accepts this repo and `codex plugin add
workbench@workbench` installs it with all 16 skills and `hooks/hooks.json`.
Whether a live Codex session loads the skills, or runs the `SessionStart` hook
that delivers the contract, is not tested. Test that first; if the hook runs,
it may be Codex's delivery.

Trial with codex-cli 0.164.0, plugin installed in the real `~/.codex`: a new
`codex exec` session in an empty repo listed `workbench:diagnose`,
`workbench:explain-gap` and `workbench:write-skill` and none of the 13
user-invoked skills. It answered the done-when question with "NOT IN MY
INSTRUCTIONS", and its session record has no contract text. Codex keeps a
trust hash per hook under `[hooks.state]` in `config.toml`, and the plugin's
hooks have none, so the likely cause is that they are untrusted. Whether a
trusted hook runs, and whether Codex sets `CLAUDE_PLUGIN_ROOT` for it, is not
tested. The owner has to trust the hook; a bypass flag was refused.

After the owner trusted the hooks, a new `codex exec` session answered with
`~/.local/share/workbench/notebook/decisions/YYYY-MM-DD-<slug>.md`, and its
session record holds the contract. Codex delivery is the hook.

Pi: open. It gets a dedicated extension layer.

## Done when

- A new Codex session and a new Pi session each answer "where does a
  decision made outside a feature folder go?" with `decisions/` in the
  notebook.

## Out of scope

- The text of the global instructions (dotfiles).
