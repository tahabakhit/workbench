---
name: settle
description: Stress-test a plan, design or decision by interview, in rounds of numbered questions with a recommended answer for each.
disable-model-invocation: true
---

# Settle

Settle every open decision in a plan before work starts. If the user names an
output file (a spec, a workflow, a design note), keep it current as decisions
land; otherwise write nothing until step 2.

## 1. Interview

1. Name the destination: the spec, decision or change this session is
   finding its way to. It fixes the scope.
2. Map the decisions the plan depends on. A question is ready when everything
   it depends on is settled. Keep two lists beside it:
   - **Not yet specified**: in scope, but not sharp enough to ask yet.
   - **Out of scope**: past the destination. It comes back only if the
     destination changes.

   A question that is sharp but cannot be settled in this session becomes an
   open item (an issue or ticket), not a Not yet specified entry.
3. Look up every fact you can (repo, docs, tools, data). Ask the user only
   for preferences, priorities and knowledge outside your reach. Run lookups
   in a subagent where you can and don't wait on them: only the questions
   that depend on a lookup wait; ask the rest now.
4. Ask all ready questions in one round, numbered. Give each your recommended
   answer and one line on why. Word each question so "yes" accepts the
   recommendation, so the user can reply "1 yes, 2 no: X". Then wait for the
   answers.
5. Read the answers, mark what they settle and what they open, and ask the next
   round. Push back once on an answer that contradicts an earlier one or a
   fact you found.
6. Stop when no question is open, or when the user stops.

Techniques, when they fit:

- **Design it twice.** For a choice with real alternatives, sketch two or more
  very different options before recommending one.
- **Workflows.** For a recurring process, settle its trigger (an event or a
  schedule), each checkpoint where the user verifies or decides, and the
  brief each checkpoint shows: a short decision-ready summary that links to
  the full output. Push checkpoints as late as they will go, so the user is
  asked once, with everything prepared.

Done when you have listed the decisions made, each with its reason, anything
still open, and the user has confirmed that list is the shared understanding.
Apart from keeping a named output file current, act on nothing before that
confirmation.

## 2. Record on approval

Propose where the decisions go: `## Decisions` in `.scratch/<feature>/spec.md`
when a feature folder exists, otherwise the `decisions:` destination in the
`## Knowledge` sections. If neither exists, leave them in your final summary
and say so.

If the session sharpened terms or made a decision that meets the ADR bar,
read `glossary-adr-format.md` in this skill's folder and propose GLOSSARY.md
entries and ADRs in that format.

Done when the user has approved, edited or rejected each proposal, and only
the approved ones are written.
