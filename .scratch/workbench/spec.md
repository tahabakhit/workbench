# workbench

## Goal

A Claude Code plugin that holds one person's skills, shared contract and
notebook for working with AI coding agents in every repo, on any machine. It replaces the skills in the agents setup repo
(`~/dev/agents/home/dot_agents/skills/`). Codex and Pi get the same skills as
plain `SKILL.md` folders.

Settled input: `.scratch/bootstrap/decisions.md` (it wins on conflict) and the
per-skill verdicts in `.scratch/bootstrap/mattpocock-analysis.md`.

## Principles

- Influence, not force. Hooks and rules add context, remind, or check the
  result. They block only irreversible damage.
- Cheapest mechanism that holds. A procedure is a skill, an always-on rule is
  an instruction line in `contract.md`, and a hook covers only what must not
  be skipped.
- One home for records. Logs, decisions, lessons, corrections, plans, goals
  and studied topics go to the notebook at `~/.agents/notebook`, in the layout the
  contract defines, on every machine.
- Sources stay neutral. No skill hardcodes a tracker, chat or calendar;
  skills read the machine's `~/.config/workbench/sources.md` and work from
  git and `.scratch/` alone without it.
- Self-contained skills. A skill never loads another skill. Reference
  material is a side file inside the skill's own folder, read at a named step.
  When several skills need the same side file, each folder gets a copy and
  `scripts/check.sh` fails if the copies differ.
- Skill names are plain English verbs that name the outcome.
- No `version` in `plugin.json`, so every commit counts as an update.

## Skills

| Skill | From | Invoked by | Issue |
|---|---|---|---|
| `rephrase` | repitch | user | 01, 06 |
| `onboard` | setup-repo, plus new machine scope | user | 03 |
| `diagnose` | diagnose | model | 04, 06 |
| `explain-gap` | an-metric-diff | model | 04 |
| `settle` | grill | user | 04, 06 |
| `improve` | improve | user | 04, 06 |
| `draft-questions` | questionnaire | user | 04, 06 |
| `plan-day` | triage | user | 04 |
| `write-skill` | writing-skills | model | 04, 06 |
| `record` | retro and harvest, merged | user | 05 |
| `write-spec` | new | user | 09 |
| `study-topic` | teach | user | topic-tutor spec |
| `prototype` | prototype, logic branch only | user | none |
| `research` | research | user | none |
| `handoff` | handoff and claude-handoff, merged | user | none |

The plugin does not ship herdr (the herdr binary provides it through
`herdr --skill`), moshi-visuals (the Moshi app installs it) or rhel-fedora
(specific to one machine, so it belongs in the Fedora machine setup).

## Repo layout

```
.claude-plugin/plugin.json       name, description, author; no version
.claude-plugin/marketplace.json  this repo as its own marketplace, source "./"
skills/<name>/SKILL.md           one folder per skill
skills/<name>/agents/openai.yaml Codex metadata for user-invoked skills
skills/<name>/scripts/           POSIX sh helpers, executable bit set
skills/<name>/<side-file>.md     reference material the skill reads
contract.md                      work tracking, notebook layout, sources, writing rules
notebook/                        the notebook tool's Go source (issue 11)
scripts/check.sh                 plugin validation plus side-file copy check
AGENTS.md                        instructions for agents working in this repo
README.md                        install steps and the order to use the skills in
```

`contract.md` reaches the agents through each machine's setup. On the work
Mac, chezmoi templates in `~/dev/agents` build `~/.claude/CLAUDE.md` and
`~/.codex/AGENTS.md`. They include `contract.md` next to that machine's own
policy.

## Tasks

Order step 1 (machine permission fixes) is tracked in `~/dev/agents` and does
not block anything here.

| # | Task | Blocked by | Status |
|---|---|---|---|
| 01 | Plugin skeleton with `rephrase` as the first skill, installed on this Mac | | done |
| 02 | `contract.md`, included in this Mac's global instructions | 01 | done |
| 03 | `onboard` skill: machine and repo scopes | 02 | done |
| 04 | Move and rename the seven one-to-one skills | 01 | done |
| 05 | Merge retro and harvest into `record` | 02, 04 | done |
| 06 | Apply the "merge into own" verdicts to the moved skills | 04 | done |
| 07 | Switch this Mac to the plugin and remove skills from the agents repo | 03, 05, 06 | done |
| 08 | Switch the Fedora PC to the plugin | 07 | done |
| 09 | `write-spec` skill: spec and issues in one pass | 02 | done |
| 10 | Five-scenario trial of `write-spec` | 09 | done |
| 11 | notebook: build the vault tool in this repo, up to switch-over | | open |
| 12 | Hooks and safety rules: decide coverage, keep them unrestrictive | 07 | done |
| 13 | Topic-tutor skill | 07 | done |
| 14 | The plugin delivers the contract and Knowledge itself | 08 | done |
| 15 | Develop on the Fedora PC only; decide how the Mac consumes the plugin | 14 | done |
| 16 | Ship the owner's agent rules as `agent-rules.md` | | done |
| 17 | Install from GitHub on every machine | | done |
| 18 | Connect the Mac to the public repo | | open |

Steps 2 and 3 of the Order are issues 01 to 10. Issues 16 and 17 are hand-overs from the dotfiles. Issues 11 to 13 are steps 4
to 6. Each of those gets its own spec before work starts.

## Out of scope

- Hook and safety-rule content (issue 12) and the notebook design (issue 11).
- Permission rules and Bash guards in the plugin. Git safety, gitleaks
  included, is an agent rule in `agent-rules.md` only.
- Hooks and agents for Codex and Pi, apart from the Pi agents listed under Later.
- An internal marketplace for colleagues.
- How notebook and workbench are installed together (still open in the
  decisions file).
- The skills dropped by the verdicts, and a router skill.
- Porting obsidian-mind commands.

Where the analysis disagrees with the decisions file, the decisions file
wins. The analysis suggests hand-bumped versions and permission rules in the
plugin; neither applies here.

## Later

Owned by workbench, not prioritised yet. Each gets an issue when it starts.

- Statusline renderer for every harness that has a statusline. It reads a
  spec file the dotfiles ship (proposed `~/.config/statusline/statusline.json`:
  palette, lines of segments, thresholds). Start from the renderer in the
  dotfiles history: `git -C ~/dev/dotfiles log --all --
  private_dot_claude/executable_statusline.sh`.
- The six Pi agent definitions (Coding, Deep-Review, Explore, Plan, Worker,
  general-purpose). The live copies are in `~/.pi/agent/agents/`. History:
  `git -C ~/dev/dotfiles log --all -- dot_pi/agent/agents`.

## Decisions

| Decision | Reason |
|---|---|
| Skill names and merges as in the Skills table: `explain-gap`, `diagnose`, `settle`, `record` (retro plus harvest, scopes `session` and `since-last`), `improve`, `draft-questions`, `rephrase`, `onboard` (absorbs setup-repo), `plan-day`, `write-skill`, `write-spec` | Plain verbs naming the outcome; `simplify` and `plan` would clash with built-ins |
| herdr, moshi-visuals and rhel-fedora stay out of the plugin | Each is owned by its own tool or machine |
| Order step 1 is tracked in `~/dev/agents`, not here | It changes only that repo |
| The agents repo stops shipping the skills; Claude Code gets them only from the plugin (issues 07, 08) | One source for every harness |
| Codex and Pi get no skills from links in `~/.agents/skills`; native extensions come after the plugin (issue 14) | Each harness loads skills its own way |
| The plugin's `SessionStart` hook delivers the contract; the dotfiles keep only the contract include for Codex and Pi (issue 14) | Workbench does not edit the global instructions |
| `scripts/check.sh` runs `claude plugin validate --strict` and tolerates only its missing-version warning | `plugin.json` has no `version` on purpose; every other warning still fails |
| Workbench is developed, committed and pushed only from the Fedora PC (issue 15) | One place where changes are made |
| The repo is public, with history rewritten so the old trial files hold no employer details (issue 15) | A public repo lets the Mac install without the personal account's credentials |
| The Mac has no workbench checkout and installs from the `github` marketplace source when the owner gets to it (issue 18) | A public repo avoids the Mac's work-account auth problem; nothing else waits on the Mac |
| Workbench owns the text of the owner's agent rules as `agent-rules.md`; the dotfiles only include it (issue 16) | One home for agent guidance; `CLAUDE.md` keeps its standing over hook context |
| Every machine installs the plugin from the `tahabakhit/workbench` GitHub marketplace, the Fedora PC included (issue 17) | Same source everywhere; sessions run what is pushed |
| `prototype`, `research` and `handoff` ship as manual skills, and mattpocock-skills is uninstalled | They were the last verdicts marked rewrite; with them built the plugin needs no reference copy. `prototype` keeps only the one-file logic demo, since the UI-variants branch assumes a web app repo |
| Workbench owns the statusline renderer and the Pi agents; both wait under Later | Agent tooling belongs here, but neither is a priority |
| The Mac-port hook scripts and the `om` MCP server stay in the dotfiles port folder for now | Not needed yet |
| notebook is developed in this repo and ships in this plugin; workbench is the skills plus the notebook | They only make sense together, and one repo keeps the notebook layout and the skills that file into it in step |
| Records go to the notebook at `~/.agents/notebook` in one fixed layout on every machine; no per-machine or per-repo record destinations | A plugin for any machine has one home for records; per-machine destinations only bent it around an older vault setup |
| Until the notebook tool exists, skills write the notebook's markdown files directly and commit them | Plain markdown in git, so the tool takes the files over without a migration |
| Each machine keeps only its work sources in `~/.config/workbench/sources.md`, written by `onboard`; the session hook no longer loads it | Sources really differ by person and job; only the skills that gather work need them |
