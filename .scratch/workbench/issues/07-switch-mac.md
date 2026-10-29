# 07 Switch this Mac to the plugin

Status: done
Blocked by: 03, 05, 06

## What

- Claude Code gets the skills only through the plugin: remove the
  `~/.claude/skills/` symlinks for the eleven old skills (keep herdr and
  moshi-visuals).
- Codex and Pi get them from `~/.agents/skills/`: point each skill there at
  `~/dev/workbench/skills/<name>` through the chezmoi setup in `~/dev/agents`.
- Remove the old skills from `~/dev/agents/home/dot_agents/skills/` and every
  reference to their old names in that repo. Update its `AGENTS.md` and
  README. Commit there.
- In the `## Knowledge` sections of the global instructions source, key the
  `sources` lines to the new skill names (`triage` becomes `plan-day`), or
  `plan-day` finds no sources beyond the three built-in ones.

Run `git status` in both repos first.

## Done when

- `ls -la ~/.claude/skills ~/.agents/skills` shows no old names, and
  `~/.agents/skills` links to the workbench skills.
- `claude plugin details workbench` lists all ten skills.
- `rg -n 'an-metric-diff|repitch|setup-repo|writing-skills' ~/dev/agents/home`
  prints nothing.
- `rg -n 'for .triage' ~/dev/agents/home ~/.claude/CLAUDE.md` prints nothing.
