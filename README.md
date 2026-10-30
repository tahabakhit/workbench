# workbench

A Claude Code plugin, and its own marketplace, with skills for working with
AI coding agents across repos and machines.

## Install

```sh
claude plugin marketplace add tahabakhit/workbench
claude plugin install workbench@workbench
```

Turn on auto-update for the marketplace: run `/plugin` in Claude Code, open
Marketplaces, select `workbench` and enable auto-update. New sessions then
pick up each pushed commit.

To pull a change in by hand, run `claude plugin marketplace update workbench`
and start a new session.

For Codex, run `codex plugin marketplace add tahabakhit/workbench` and
`codex plugin add workbench@workbench`, then accept the workbench hooks when
`codex` next starts.

For Pi, run `pi install git:github.com/tahabakhit/workbench`. It adds the
contract, the skills and correction flagging. A personal skill with the same
name wins over the package's, and the `researcher` agent is not included.


## Use

Run `/workbench:onboard machine` once per machine, and
`/workbench:onboard repo` once per repository to choose where its tasks live:
local markdown, GitHub, GitLab or another tracker. When unsure which skill
fits, run `/workbench:choose-skill <your situation>`.

## Check

```sh
sh scripts/check.sh
```

Validates the plugin with `claude plugin validate --strict`, fails if a side
file copied into several skill folders differs between them or `choose-skill`
misses a skill, and runs the correction flag, quiz shuffle and spec check
tests.

## Safety

The plugin ships no git safety: nothing here blocks a force push or scans for
secrets, because plugins cannot carry permission rules. Set those up in each
machine's own Claude Code settings or git hooks.
