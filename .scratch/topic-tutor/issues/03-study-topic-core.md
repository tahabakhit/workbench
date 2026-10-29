# 03 `study-topic` core: mission, sources with gaps, first lesson

Status: done
Blocked by: 01, 02

## What

Read `../spec.md` first, and the `teach` row of
`.scratch/bootstrap/mattpocock-analysis.md`. Use `workbench:write-skill`'s
guidance for wording.

- Add `skills/study-topic/SKILL.md` and `skills/study-topic/agents/openai.yaml`
  (copy the shape of another skill's `openai.yaml`).
- One session takes a topic and produces, in the `learning` destination
  under one folder per topic:
  - `mission.md`: why the owner is studying it, what "done" looks like, and
    what is out of scope, settled with the owner in a few questions;
  - `sources.md`: primary sources with what each covers, from the
    `researcher` agent, and a `## Gaps` section;
  - `lessons/01-<slug>.md`: the first lesson, concrete case before rule,
    every claim cited to `sources.md`.
- With no `learning` destination it stops and points to `onboard`.
- Without the `researcher` agent (Codex, Pi) it does the research inline
  under the same rules.
- No quiz and no resume yet; those are 04 and 05.

## Done when

- `sh scripts/check.sh` passes.
- A session started with `claude --plugin-dir ~/dev/workbench` runs
  `study-topic` on the test topic "POSIX sh parameter expansion" (or one
  the owner names instead), and the three files exist
  in the `learning` destination. Quote their paths, the `## Gaps` section
  and one cited claim from the lesson.
