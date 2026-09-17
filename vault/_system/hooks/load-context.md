# Hook: LoadContext

**Event**: `session-start`
**Blocking**: Yes (must complete before agent proceeds)
**Fail Behavior**: Abort session

## Instructions

At the start of every session, load ALL files from `_codex/system/`:

1. Read `_codex/system/identity.md` — Who you serve
2. Read `_codex/system/mission.md` — The mission
3. Read `_codex/system/active-goals.md` — Current priorities
4. Read `_codex/system/beliefs.md` — Core beliefs
5. Read `_codex/system/preferences.md` — How to communicate and work
6. Read `_codex/system/models.md` — Decision-making models
7. Read `_codex/system/frames.md` — Thinking frameworks
8. Read `_codex/system/wisdom.md` — Hard-won lessons
9. Read `_codex/system/narratives.md` — Brand and personal stories
10. Read `_codex/system/predictions.md` — Future bets

Internalize this context. It informs every response.

## Failure Mode
If `_codex/system/` is missing or unreadable, alert the user and abort. Do not proceed without identity context.
