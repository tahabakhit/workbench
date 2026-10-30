# Phase boundaries

A phase is a chunk of work inside a session: the interview, the build, the
review. It ends when the user would say "ok, that's done". The boundary
between two phases is the only place to decide what happens to the session.

Every move except continuing replaces the session as it happened with a
summary of it: less noise and more room, but the reasons behind decisions
get lossy. Pay that only when staying costs more.

Walk the questions in order; the first yes wins.

1. **Can the next phase use this context as it is?** It can when the next
   phase needs this one's reasoning word for word (interview to spec is the
   usual case) and the context has room for it. Then **continue**.
2. **Is everything here irrelevant to what comes next?** Then **clear**
   (`/clear` in Claude Code, `/new` in Codex). The old session stays
   resumable.
3. **Does the work have to travel?** To another harness, another machine or
   repo, another person, or a side task found mid-phase that should not
   derail this one. Then **`handoff`**. If nothing travels, it is not
   needed.
4. **Can the next task run without the user steering it?** A tightly scoped
   job such as a review of the diff. Then a **subagent**, and this session
   stays as it is.
5. **Otherwise compact** (`/compact`), with an instruction naming what the
   next phase needs, such as `/compact next we review the parser`.

These are judgement calls; asking them in order, at the boundary, is the
point.
