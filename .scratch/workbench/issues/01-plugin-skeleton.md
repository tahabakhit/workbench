# 01 Plugin skeleton with rephrase

Status: done

## What

Make this repo an installable Claude Code plugin and its own marketplace,
with one skill, as a thin end-to-end slice.

- `.claude-plugin/plugin.json`: name `workbench`, description, author; no
  `version`. Check the accepted fields with `claude plugin validate`
  (`claude plugin init` shows a scaffold) instead of guessing.
- `.claude-plugin/marketplace.json`: one plugin, `source: "./"`.
- `skills/rephrase/`: copied from
  `~/dev/agents/home/dot_agents/skills/repitch/` and renamed (folder, `name:`,
  heading, `agents/openai.yaml`). Content stays as it is; issue 06 changes it.
- `scripts/check.sh` (POSIX sh): runs `claude plugin validate --strict .`,
  then fails if any side file copied into several skill folders differs
  between them (none exist yet; the check must still run).
- `README.md`: install commands for one machine.
- On this Mac: `claude plugin marketplace add ~/dev/workbench`, then
  `claude plugin install workbench@workbench`, and turn on auto-update for
  the marketplace.

Leave `~/dev/agents` unchanged (issue 07 does the switch).

## Done when

- `sh scripts/check.sh` exits 0, quoted.
- `claude plugin details workbench` lists the `rephrase` skill.
- `rg -n '"version"' .claude-plugin/` prints nothing.
