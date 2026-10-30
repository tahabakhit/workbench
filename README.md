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
file copied into several skill folders differs between them, and runs the
correction flag, quiz shuffle and spec check tests.
