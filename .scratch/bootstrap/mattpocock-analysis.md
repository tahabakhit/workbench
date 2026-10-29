# mattpocock-skills analysis

Source: `github.com/mattpocock/skills` at `a496f55`, read by five agents and trialled on two throwaway projects (`~/dev/ref/mp-trial-eng`, `~/dev/ref/mp-trial-prod`).

## How mattpocock-skills works

- **Packaging.** It is one Claude Code plugin, skills only, with no hooks, agents or MCP servers. `plugin.json` lists 27 skill paths, and the same repo is its own marketplace (`marketplace.json`, `source: "./"`). Skills live in `skills/<bucket>/<name>/`:
  - `engineering/` and `productivity/` ship.
  - `in-progress/` is beta and not shipped.
  - `misc/` is frozen.
  - `deprecated/` stays empty, because retired skills are deleted.
- **Two kinds of skill.**
  - About 17 are user-invoked front doors, marked `disable-model-invocation` plus `agents/openai.yaml` for Codex.
  - About 10 are model-invoked primitives (`grilling`, `tdd`, `domain-modeling`, `code-review` and others). Front doors chain to primitives with "Call the Skill tool with X". `grill-me` is literally that one line.
- **End-to-end flow.** Run `/setup-matt-pocock-skills` once per repo. Then: `grill-with-docs` interviews and writes GLOSSARY.md and ADRs. `to-spec` writes the spec with no new interview. `to-tickets` splits it into vertical-slice tickets with `Blocked by`. Then `implement`, which uses `tdd`, `code-review` and a commit. `retro` comes last. `ask-matt` is a router that tells you which skill to type next.
- **On-ramps beside the main flow.**
  - `triage`: a label state machine for issues that come in.
  - `diagnosing-bugs`: hard bugs.
  - `wayfinder`: a map of decision tickets for foggy work that spans sessions.
  - `prototype` and `handoff`: detours.
- **Repo conventions it installs:**
  - `docs/agents/issue-tracker.md`, `domain.md` and `triage-labels.md`, plus an `## Agent skills` pointer block in CLAUDE.md or AGENTS.md.
  - A root `GLOSSARY.md` (each term with `_Avoid_:` synonyms) and `docs/adr/NNNN-slug.md`. An ADR is offered only when all three gates hold.
  - The local tracker is `.scratch/<feature>/spec.md` plus `issues/NN-slug.md`, which is the same as your convention.
- **Hard and soft dependencies.** Skills whose output would be wrong without config (`to-spec`, `to-tickets`, `triage`) tell you to run setup. The rest degrade silently.
- **Governance.** SCOPE.md requires an observed failure before any change and rejects per-person tweaks. `.out-of-scope/<concept>.md` records rejected ideas. A human docs page per skill sits outside SKILL.md. Changesets bump `plugin.json`, because managed installs only update when the version changes.
- **Hidden dependencies the trials exposed:**
  - GitHub and `gh` as the default tracker.
  - Subagents: code-review degrades quietly without them.
  - A typed language with Jest-style tests ("typecheck" means nothing in sh).
  - Skill-to-skill loading, which its own docs admit is unreliable.

## Per-skill verdict

| skill | what it does | own equivalent | verdict | why |
|---|---|---|---|---|
| grilling | Asks the decision frontier in rounds, a recommended answer per question, facts looked up by subagent | grill | merge into own | Take "yes accepts the recommendation", lookups that don't block the round, and an explicit wait-for-confirmation gate. Leave out the emoji format |
| grill-me | One-line wrapper that calls grilling | grill | drop | Own grill is already the self-contained front door; the wrapper only adds a fragile skill-to-skill load |
| grill-with-docs | Router to grilling plus domain-modeling | grill | merge into own | Add an optional ending that proposes glossary terms and ADRs and writes them only on approval. In the trial the inline writes were ceremony |
| domain-modeling | Challenges terms; writes GLOSSARY.md and ADRs inline | GLOSSARY/ADR convention | merge into own | One format side file read by grill, retro and setup-repo. No model-invoked slot, no side-effect writes (ADR 0003) |
| to-spec | Turns a settled conversation into a spec; no interview | `.scratch` convention | rewrite | Keep synthesis-only and a real Out of scope list. Cut the "LONG, extremely extensive" user stories (the trial got 18 for 3 commands) and the label |
| to-tickets | Vertical-slice tickets with blocking edges, after an approval quiz | `issues/NN-slug.md` | merge into own | Fold into the rewritten spec skill: slices, prefactor first, the quiz. Plain `Status:` lines: its bold `**Status:**` breaks grep (trial) |
| implement | tdd, typecheck, code-review, commit | global "verify" rules | merge into own | Two lines in the work-tracking convention: finish at done-when, quote the output, set `Status: done`. Its order (review before commit) is broken |
| implement-spec | Parallel implementer subagents in worktrees, merged onto an integration branch | none | drop | Heavy, with documented long fix loops. Clashes with your one-worktree-per-area rule and with delegatr. Keep only the "frontier" idea, for triage |
| tdd | Red-green reference with agreed seams and anti-patterns | global expected-value rule; diagnose | merge into own | Two global lines: failing test first, mock only at boundaries. Its TypeScript examples and typecheck step didn't fit sh in the trial |
| code-review | Two parallel axes, Standards and Spec | built-in /code-review | merge into own | In the trial the Spec axis found a real bug. Make "diff vs done-when and Out of scope, in a fresh subagent" the closing check. Never name a skill code-review |
| pr | PR body: smallest visual, evidence, one-way or two-way door | repo-writing rules | merge into own | A format, not a procedure: one paragraph in instructions. Rewrite the visuals menu, which is credited to a third party |
| prototype | Throwaway HTML logic demo or UI variants that answer one question | none (moshi-visuals displays only) | rewrite | A clickable one-file demo suits an analyst who works with PMs. Short and manual; file the answer to `decisions:` |
| wayfinder | Tracker map of decision tickets for multi-session fog | brainstorm.md plus grill | merge into own | Take the vocabulary: Destination, "Not yet specified" vs Out of scope, ticket when the question is sharp. Skip the tracker machinery |
| research | Background agent, primary sources, a citation per claim | CLAUDE.md research routing | rewrite | One `researcher` agent shared with teach. Fix its nesting bug ("do the work, don't delegate") and add a stopping rule |
| diagnosing-bugs | Six-phase debugging from a failing loop | diagnose | merge into own | Add: stop and ask when no loop is possible, minimise first, show hypotheses, redact secrets, name the cause in the commit |
| triage | Moves issues through labels; agent briefs | triage (daily ranking) | keep own | Different job: theirs is maintainer inbox work. Borrow "check for an existing implementation by concept" for the spec skill |
| setup-matt-pocock-skills | Writes docs/agents files and an Agent skills block per repo | setup-repo | keep own | Own writes one AGENTS.md under 60 lines with verified commands. Trial: about 100 lines of dead template, and label roles with no done state |
| retro | Proposes fixes to the agent's environment, files nothing | retro plus improve | merge into own (improve) | Own retro keeps the filing role. Add its candidate types to improve, plus "an unguarded repo is a finding" and "mechanical gets a check" |
| codebase-design | Deep-module vocabulary for TypeScript | none | drop | Off-domain, model-invoked, burns tokens. Keep "design it twice" as one grill technique line |
| improve-codebase-architecture | HTML report of refactor candidates, then grill one | none | drop | Engineering upkeep you don't do; opens the browser. Keep "rejected for a lasting reason means offer an ADR" |
| ask-matt | Router plus a context tree for phase boundaries | none | drop | Not now; revisit as a router once the plugin's manual skills reach about 15. The phase-boundary tree becomes instruction lines when a failure shows up |
| wizard | Generates a bash script for manual setup steps; writes .env and GitHub secrets | none | drop | Conflicts with secrets in Keychain, POSIX sh, and the model-invoked cap |
| to-questionnaire | Questionnaire for the person who holds the missing knowledge | questionnaire | merge into own | Add "I don't know is welcome", a catch-all question, "why this matters" only where needed, and a coverage done-check |
| wait-what | Re-pitch in ASD-STE100 (Simplified Technical English) with glossary terms | repitch | merge into own | Allow adding a missing premise (own says "nothing new"), use GLOSSARY.md when present, add ASD-STE100 as the leading word |
| handoff | Markdown handoff to `$TMPDIR`, artifacts referenced by path | none | rewrite | Needed for moving between Claude Code, Codex and Pi and between machines. Write to `.scratch/<feature>/`, label claims verified or assumed |
| teach | Long-running tutor: mission, cited resources, HTML lessons, learning records | planned teach | merge into own | Main input for your teach: mission, sources with gaps, retrieval quizzes, records in the vault. Shuffle answers in code (instruction-only fixes put the right answer first 33 of 33 times) |
| writing-for-agents | Theory of wording: leading words, pointers, completion criteria | writing-skills | merge into own | Lift the vocabulary into a writing-skills side file. Own keeps the evidence gate and tests |
| git-guardrails (misc) | Skill that installs a PreToolUse hook blocking dangerous git | CLAUDE.md prose only | rewrite | Ship `permissions.ask` rules or a sh hook directly in the plugin. Narrower: don't block plain push |
| setup-pre-commit (misc) | Husky, lint-staged and Prettier for JS | none | drop | Node-only; gitleaks already runs on the personal PC |
| migrate-to-shoehorn (misc) | Migrates TypeScript test casts to the author's library | none | drop | Library-specific |
| scaffold-exercises (misc) | Course exercise folders | none | drop | The author's course repo only |
| chief-of-staff (in-progress) | Coordinator session with background agents and a dated wiki | vault plus om | merge into own | Into the vault tooling: on ingest update every page the fact touches, date or commit-pin claims, supersede rather than rewrite, periodic lint |
| claude-handoff (in-progress) | Handoff plus `claude --bg` launch | none | merge into own | Fold into the rewritten handoff as an optional launch step that delegatr can replace later |
| loop-me (in-progress) | Stateful grilling into workflow specs | grill | merge into own | Add trigger, checkpoint, push right and brief to GLOSSARY.md. Let grill take an optional output file |
| setup-ts-deep-modules (in-progress) | dependency-cruiser rules for TypeScript | none | drop | Keep "prove the check bites" (pass, fail, pass) for writing-skills and setup-repo |
| writing-fragments / -shape / -beats (in-progress) | Article-writing trio | none | drop | Outside the goal. "Grounding" order is useful for teach |

## Ideas to adopt beyond single skills

- **Plugin packaging.** An explicit `skills` array in `plugin.json`, plus a repo-local `marketplace.json` (`source: "./"`) so Codex reads the same manifest. `in-progress/` stays linked locally and is not shipped. Bump the version by hand, because installs update only on a version change. Run `claude plugin validate` after each change.
- **One written chain, chosen by session count.** Start with triage, then grill. A single-session job is worked in place. A multi-session job goes through the spec skill, which writes `spec.md` plus issues. Each issue is worked to its done-when, then closed with the Spec-axis check in a fresh subagent. retro and harvest come last. Put this chain in the plugin README; a router skill can come later.
- **Hard and soft dependency rule.** retro, harvest and triage need `## Knowledge` and point to first-run setup when it is missing. grill, diagnose and the spec skill degrade silently and name their fallback.
- **Self-contained skills over skill-to-skill chaining.** When a skill needs reference material, it reads a named side file at a named step. Subagent briefs say "start with a fresh context; do the work yourself, invoke no skill", which fixes their nesting bugs.
- **An "It's working if" checklist per skill**, kept in the repo next to the skill rather than in SKILL.md. It doubles as the writing-skills test cases.
- **A shared vocabulary in the plugin's GLOSSARY.md:**
  - frontier: open tasks whose blockers are done
  - tracer bullet or vertical slice
  - one-way vs two-way door
  - fog vs out of scope
  - context pointer: a link instead of pasted text
  - primary source

## What not to copy

- **Per-repo setup plumbing.** `docs/agents/*.md`, label vocabularies and git remote wiring are rejected by ADR 0003, and `setup-repo` already covers the facts.
- **Skills that edit docs as a side effect.** Inline GLOSSARY and ADR writes cause their most-reported bug; propose at the end, write on approval.
- **The engineering stack as default process.** TypeScript and Jest examples, typecheck steps, worktree fan-out and HTML architecture reports don't fit SQL, dbt, sh and docs work.
- **Their posture choices.** GitHub-first trackers with Jira refused, files dropped in `$TMPDIR` or the current directory, browser opens, and emoji question markers.
- **Distribution machinery built for a large public audience.** Changesets, release workflows, issue bots, the six-harness install wording, and bash `link-skills.sh`. The chezmoi link script already does that job.

## Gaps neither covers

- **First-run setup of `## Knowledge`** (it writes global instructions from the tools it finds) is still only a proposal, and mattpocock has no global layer at all.
- **No single home for decisions made while grilling.** Their grilling keeps decisions in the conversation; yours writes nothing unasked. A defined approved write target is needed: `spec.md`, `decisions:` or the vault.
- **No way to close the loop on work items.** Neither side sets `Status: done` in the issue and the spec table, or reconciles them with Jira or GitHub, after work finishes.
- **Ad hoc work without a feature folder.** Quick asks that never get a `.scratch` folder leave no trace for retro or harvest to file.
- **Enforcement is all prose.** There is no plugin hook or permission set for git safety, vault write validation (the brainstorm leans to one validating CLI), or a portable handoff and delegation seam toward delegatr.

## Open questions for the user

1. **ADR bar: any one gate (your AGENTS.md and retro) or all three (theirs)?**
   Recommended: all three. Fewer, sharper ADRs, and the vault's `decisions:` takes the rest.
2. **Ship one spec skill (spec plus issues in one pass), reversing "nothing the harnesses already ship: specs, tickets"?**
   Recommended: yes, and record the reversal as the plugin repo's first ADR.
3. **How should the Spec-axis review run: a closing line in the work-tracking convention, or a small manual skill?**
   Recommended: a convention line first. Promote it to a skill only if it gets skipped.
4. **Does git safety ship as plugin `permissions.ask` rules or as a hook?**
   Recommended: `ask` rules for force push, `reset --hard`, `clean -f`, `branch -D` and `checkout/restore .`. A hook only if a rule can't express it.
5. **What happens to delegatr's mattpocock files (`docs/agents/issue-tracker.md`, `triage-labels.md`) when the plugin is removed?**
   Recommended: migrate them to a `## Work tracking` section in AGENTS.md with plain `Status: open/done/wontfix`, in the same change that removes the plugin.

The `mp-trial-eng` trial repo is in `~/dev/ref/mp-trial-eng`, with its run log in `TRIAL-LOG.md` there. The `mp-trial-prod` trial repo is in `~/dev/ref/mp-trial-prod`.
