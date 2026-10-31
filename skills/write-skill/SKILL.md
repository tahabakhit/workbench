---
name: write-skill
description: Use when creating, editing or reviewing a skill (a SKILL.md folder) or an agent instruction file (AGENTS.md, CLAUDE.md). Covers when a skill is worth writing, how to word it so capable models are not made worse, portability between Claude Code and Codex, and how to check it triggers and helps.
---

# Write skill

A skill earns its place only by changing what the model does. Instructions the
model already follows cost context on every turn and say nothing.

## 1. Decide the home

Read `cheapest-form.md` in this skill's folder and pick the row that fits.

Done when the change has one home and nothing else restates it.

## 2. Start from evidence

Write down two or three real exchanges where the model went wrong without the
skill. They are the test cases. A skill with no observed failure behind it is
not written.

Done when each failure is a prompt plus what should have happened.

## 3. Write

Read `vocabulary.md` in this skill's folder first; it names what to look for
in each line (pointers, leading words, completion criteria, no-ops).

- Write only what closes the observed gaps. Hunt no-ops sentence by sentence.
- Concrete beats general: name the command, path, table or check.
- Fixed steps only for order-sensitive work; judgement work gets the goal and
  the traps.
- Plain words. No MUST, CRITICAL or capitals; newer models overreact to them.
- No contradictions, including with the instruction files. Codex stalls on
  them.
- Each step ends on a "done when" line the model can check.
- Description: imperative and specific ("Use this before investigating any
  bug, error ..."), naming the trigger situations. Hedges like "when the cause
  is not obvious" let the model skip the skill. Under about 400 characters,
  and no overlap with another skill's description.
- Body under about 100 lines. Move reference material to a side file that the
  body says when to read.

## 4. Portability

- Frontmatter: `name` and `description`. Nothing else, so the model can invoke
  every skill. A skill that fails step 2 of the check gets a better description.

- Name no tool from one harness without the other's equivalent.
- To call another skill, write "use the `<name>` skill" on one line at the step
  that needs it, and say what comes back. All three harnesses load it by name.
- Scripts: POSIX `sh` with `jq`, inside the skill's `scripts/` folder.

## 5. Check

1. `claude plugin validate <skills folder>` passes.
2. `scripts/trigger-check.sh <skill> yes "<prompt>"` for three prompts that
   should trigger it, and `... no "<prompt>"` for two near misses. Run it from
   a repo where the prompt makes sense; in an empty folder the model stops to
   ask instead.
3. Replay one failing case from step 2 with the skill turned off
   (`--settings '{"skillOverrides":{"<skill>":"off"}}'`) and on; keep the skill
   only if the result with it is better.
4. When the change adds a check (a script, test or hook), prove it bites:
   it passes, then fails on a deliberately broken case, then passes again
   once that case is reverted. A check that never failed is not yet a check.

Done when all four pass. For a skill worth tuning further, `claude plugin
eval` compares runs with and without it.
