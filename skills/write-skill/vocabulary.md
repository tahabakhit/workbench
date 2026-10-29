# Vocabulary

Terms for judging the wording of a skill or instruction file.

- **Context pointer**: a line in context that names material elsewhere and
  when to reach it (a skill description, an `AGENTS.md` line naming a doc).
  Its wording, not its target, decides how reliably the material gets read.
- **Branch**: a distinct case the text handles. A pointer lists one trigger
  per branch; synonyms for one branch are one trigger written twice.
- **Context load**: what always-loaded text costs on every turn, used or not.
- **Cognitive load**: what the user must remember about which skills exist
  and when to use them. Spend it only where the user's judgement matters.
- **Step** and **reference**: what the model does, in order, versus facts
  and rules it looks up. Both can sit in one file.
- **Disclosure**: moving reference out of the body into a side file behind a
  pointer. Inline what every branch needs; disclose what only some reach.
- **Co-location**: a concept's definition, rules and caveats under one
  heading, not scattered.
- **Sprawl**: a body too long even when every line is live; attention thins.
- **Completion criterion**: the "done when" that ends a step. Judge its
  clarity (can done be told from not done) and its demand (how much work it
  forces).
- **Premature completion**: ending a step early because later steps pull.
  Sharpen the criterion first; hide later steps only behind a real context
  boundary, such as a subagent.
- **Legwork**: the digging a demanding criterion forces without being
  written as its own step.
- **Leading word**: one existing word the model already knows (_tight_,
  _red_) repeated instead of a phrase, so it anchors the behaviour cheaply.
- **Negation**: a ban names the unwanted behaviour and makes it more likely.
  State the target behaviour; keep a ban only as a hard guardrail, paired
  with the positive.
- **Single source of truth**: each meaning lives in one place. Duplication
  costs upkeep and inflates its weight.
- **Cache**: text that restates what the environment already says (config,
  `--help`, layout). Keep only what a lookup cannot find: conventions,
  reasons, traps.
- **Sediment**: stale lines that pile up because adding feels safe and
  removing feels risky.
- **No-op**: a sentence the model already obeys by default. Delete the whole
  sentence; for a word too weak to beat the default, pick a stronger word.
