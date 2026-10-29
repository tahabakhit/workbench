# workbench

## Goal

A Claude Code plugin that holds one person's skills and shared contract for
working with AI coding agents in every repo, on a work Mac and a personal
Fedora PC. It replaces the skills in the agents setup repo
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
- Sources stay neutral. No skill hardcodes a tracker, chat, calendar or vault.
  Skills read the destinations in the `## Knowledge` sections.
- Hard and soft dependencies. A skill that needs a destination it can't find
  points to `onboard`. Every other skill degrades and names its fallback.
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
contract.md                      work tracking, knowledge destinations, writing rules
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
| 07 | Switch this Mac to the plugin and remove skills from the agents repo | 03, 05, 06 | open |
| 08 | Switch the Fedora PC to the plugin | 07 | open |
| 09 | `write-spec` skill: spec and issues in one pass | 02 | done |
| 10 | Five-scenario trial of `write-spec` | 09 | done |
| 11 | notebook: build the vault tool up to switch-over (own repo, own spec) | | open |
| 12 | Hooks and safety rules: decide coverage, keep them unrestrictive | 07 | open |
| 13 | Topic-tutor skill | 07 | open |

Steps 2 and 3 of the Order are issues 01 to 10. Issues 11 to 13 are steps 4
to 6. Each of those gets its own spec before work starts.

## Out of scope

- Hook and safety-rule content (issue 12) and the notebook design (issue 11).
- Permission rules in the plugin. Plugins can't ship them, so git safety
  stays in each machine's setup.
- Hooks and agents for Codex and Pi.
- An internal marketplace for colleagues.
- How notebook and workbench are installed together (still open in the
  decisions file).
- The skills dropped by the verdicts, and a router skill.
- The verdicts marked "rewrite" that no Order step covers yet: `prototype`,
  `research`, `handoff`. They are not yet specified, not rejected.
- Porting obsidian-mind commands.

Where the analysis disagrees with the decisions file, the decisions file
wins. The analysis suggests hand-bumped versions and permission rules in the
plugin; neither applies here.

## Decisions

| Decision | Reason |
|---|---|
| Skill names and merges as in the Skills table: `explain-gap`, `diagnose`, `settle`, `record` (retro plus harvest, scopes `session` and `since-last`), `improve`, `draft-questions`, `rephrase`, `onboard` (absorbs setup-repo), `plan-day`, `write-skill`, `write-spec` | Plain verbs naming the outcome; `simplify` and `plan` would clash with built-ins |
| herdr, moshi-visuals and rhel-fedora stay out of the plugin | Each is owned by its own tool or machine |
| Order step 1 is tracked in `~/dev/agents`, not here | It changes only that repo |
| `~/.agents/skills` on each machine links to this repo's skills; the agents repo stops shipping them (issues 07, 08) | One source for Claude Code, Codex and Pi |
| `scripts/check.sh` runs `claude plugin validate --strict` and tolerates only its missing-version warning | `plugin.json` has no `version` on purpose; every other warning still fails |
