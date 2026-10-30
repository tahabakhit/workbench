# Bootstrap decisions

Input for writing this repo's first spec. Treat it as the settled starting
point; everything under **Open** still needs deciding. Existing skills and
setups elsewhere are background, not a template.

## What this is

Two personal repos that together govern how one person works with AI coding
agents across every repo and machine (a work Mac and a personal Fedora PC):

- **workbench**: a Claude Code plugin with the skills, the shared contract and
  later the hooks for working in repos, tracking work, and filing work,
  decisions and lessons.
- **notebook**: a Go tool that keeps a markdown knowledge vault and serves it
  over MCP and as a CLI. It is developed in the workbench repo and ships in
  the workbench plugin.

Claude Code is the main target. Codex and Pi are kept compatible where
practical.

## Decided

### Principles

| Decision | Reason |
|---|---|
| Influence, not force. Hooks and rules mainly add context, remind, or check the result at the end. Blocking only for irreversible damage (force push to the main branch, deletes outside the repo, secrets in a commit) | Agents stay autonomous; the contract steers outcomes without limiting how models work |
| Procedures are skills, always-on rules are instruction lines, hooks and rules only for what must not be skipped | Cheapest mechanism that holds |
| Sources stay neutral: no skill hardcodes a tracker, chat or calendar. Records have one home, the notebook | The same plugin serves work and personal setups |
| Skill names are plain English verbs naming the outcome; no jargon or metaphors (`settle`, not `grill`) | Names should say what happens |
| Self-contained skills that read named side files, not skills that call other skills | Skill-to-skill loading is unreliable |

### Workflow

| Decision | Reason |
|---|---|
| One spec skill writes `spec.md` and `issues/NN-slug.md` in one pass: vertical slices, `Blocked by`, a real Out of scope list, a "done when" per issue | Slicing and an explicit Out of scope prevent confusion between sessions |
| The spec skill's done-when includes one light run each: greenfield code, brownfield repo, analytics (SQL or dbt), homelab change, non-code decision document. Pass: a fresh session starts issue 01 without asking anything | Specs must be unambiguous across all kinds of work, not only code |
| Ad hoc work records nothing during the work; the filing skills read the session transcript and git log | A folder per quick ask is ceremony that gets skipped |
| Decisions from a decision session are written only on approval: to `## Decisions` in `spec.md` when a feature folder exists, otherwise to the declared decisions destination | One home per decision, no side-effect writes |
| The session that meets an issue's done-when sets `Status: done` in the issue and the spec table | Closes the loop where the evidence is |
| Tracker updates are proposed in whichever tools the user declared (local `.scratch`, GitHub, GitLab, Jira, Linear, others), with evidence, applied on approval | The user's tools vary |
| Closing check: an instruction line to diff the result against done-when and Out of scope in a fresh subagent. Promote to a skill only if it gets skipped | Cheapest first |
| ADR only when all three hold: hard to reverse, surprising without context, a real trade-off | Fewer, sharper ADRs; other decisions go to the decisions destination |
| The mattpocock-skills per-skill verdicts in `.scratch/bootstrap/mattpocock-analysis.md` are the default for what to adopt; disputes are discussed per row | The analysis is done; revisit only disagreements |

### Onboarding and contract

| Decision | Reason |
|---|---|
| One onboarding skill with two scopes. Machine: inspects connected tools, asks which tracker, chat, calendar, code host and notes the user uses, writes the knowledge section of the global instructions. Repo: writes the repo's `AGENTS.md` and overrides | One front door; the instruction file is what every harness reads |
| Repo onboarding starts with a read-only survey of an existing repo: instruction files and other agent-tooling artifacts (skill-plugin conventions, editor rules), git remotes (code host, likely tracker), recent history (commit style, active branches), open issues and PRs when the code host is connected, existing specs, ADRs and glossary, and verified build and test commands. It then proposes per artifact: keep, convert to this plugin's conventions, or delete, and applies only what is approved | Existing work is carried over instead of duplicated or contradicted |
| The survey is bounded (recent history, not the full log) and runs in a subagent that returns facts with their source | Depth without filling the session's context |
| The plugin repo holds one `contract.md` (work tracking, knowledge destinations, writing rules). Each machine's setup includes it in the global instructions next to that machine's own policy | A plugin cannot add instruction text itself; Codex and Pi read the same file |
| A skill that needs a destination it can't find points to onboarding; skills that don't need one degrade and name their fallback | Hard versus soft dependencies |
| Records go to the notebook at `~/.agents/notebook` in one layout on every machine; only work sources are declared per machine | One home for records; sources differ by person and job |

### Plugin shape and distribution

| Decision | Reason |
|---|---|
| The repo is its own marketplace (`source: "./"`); `version` is omitted so every commit is an update; auto-update is turned on per machine | No hand-bumped versions to forget |
| Codex and Pi floor: `SKILL.md` skills from `~/.agents/skills/`. Hooks and agents are Claude Code only until a need shows up | Hooks, agents and plugin manifests differ per harness |
| Plugins cannot ship permission rules (their `settings.json` honours only `agent` and `subagentStatusLine`); git safety rules live in each machine's setup | Platform constraint |
| Plugin hooks run in every session and cannot be disabled one by one, so each needs a narrow matcher | Platform constraint |
| Colleagues install from the repo; an internal marketplace is decided when someone asks | No demand yet |

### notebook (vault tool)

| Decision | Reason |
|---|---|
| Own tool in Go, one binary serving MCP and a CLI | Language preference; MCP covers Claude Code and Codex, the CLI covers any harness and scripts |
| Notes are plain markdown in git with wikilinks, frontmatter and Mermaid. Obsidian is the viewer only; the tool never needs it running | Best browser today, replaceable without migration |
| Start with keyword search (SQLite full-text); add semantic search and rank fusion later | Smallest working core first |
| Capabilities to cover: search, recall, follow links and backlinks, file a dated work note into a project folder, file a lesson with scope and supersedes, push relevant notes at session start, validate writes | Gaps to fill; designed fresh, no existing command is ported up front |
| Developed in the workbench repo and shipped in the workbench plugin (MCP config, session-start push); skills file into it through the notebook layout in `contract.md` | They only make sense together |
| Prior art for search: `lstpsche/obsidian-mcp` (Rust, MIT). Libraries: `mark3labs/mcp-go`, SQLite FTS5, later `chromem-go` or `bleve` | No Go or Rust tool covers filing plus search together |
| The current obsidian-mind fork keeps running until notebook covers the capabilities above; then switch | No gap in daily use |

## Order

1. Close the open triage task in the existing setup repo (machine permission fixes).
2. workbench: `contract.md`, the onboarding skill, existing skills moved and renamed.
3. The spec skill and its five-scenario trial.
4. notebook, up to switch-over.
5. Hooks and safety rules: what they cover, kept unrestrictive.
6. A topic-tutor skill.

## Open

- The rename table for every skill (proposed in the spec, approved by the user).
- Which hooks and safety rules exist, and what each checks.
