---
name: settle
description: Stress-test a plan, design or decision by interview, in rounds of numbered questions with a recommended answer for each.
disable-model-invocation: true
---

# Settle

Settle every open decision in a plan before work starts.

1. Map the decisions the plan depends on. A question is ready when everything
   it depends on is settled.
2. Look up every fact you can (repo, docs, tools, data). Ask the user only
   for preferences, priorities and knowledge outside your reach.
3. Ask all ready questions in one round, numbered. Give each your recommended
   answer and one line on why, so the user can reply "1 yes, 2 no: X".
4. Read the answers, mark what they settle and what they open, and ask the next
   round. Push back once on an answer that contradicts an earlier one or a
   fact you found.
5. Stop when no question is open, or when the user stops.

Done when you have listed the decisions made, each with its reason, and
anything still open. Write no files unless the user asks.
