# Glossary and ADR format

Propose entries at the end of a session; write them only when the user
approves. Create a file or folder only when the first entry goes in.

## GLOSSARY.md

One file at the repo root, terms grouped under subheadings when clusters
emerge:

```md
# {Context name}

{One or two sentences on what this context is.}

## Language

**Order**:
A request from a customer to buy one or more items.
_Avoid_: Purchase, transaction
```

- Pick one term per concept; list the other words under `_Avoid_`.
- One or two sentences: what the term is, not what it does.
- Only terms specific to this project; no general programming concepts and
  no implementation details.

## ADRs

`docs/adr/NNNN-slug.md`, numbered one above the highest existing number.
Offer one only when all three hold:

1. Hard to reverse: changing your mind later is costly.
2. Surprising without context: a future reader would ask "why this way?".
3. A real trade-off: there were genuine alternatives.

Otherwise the decision goes to `decisions/` in the notebook instead.

```md
# {Short title of the decision}

{One to three sentences: the context, what was decided, and why.}
```

Add "Considered options" or "Consequences" only when a rejected alternative
or a downstream effect is not obvious.
