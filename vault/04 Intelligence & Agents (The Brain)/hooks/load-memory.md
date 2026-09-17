# Hook: LoadMemory

**Event**: `session-start`
**Blocking**: No
**Fail Behavior**: Log warning, continue

## Instructions

After loading system context, load memory tiers:

1. Read `_codex/memory/work/current-focus.json` — What's active
2. Scan last 10 entries in `_codex/memory/learning/signals.jsonl` — Recent feedback
3. Read `_codex/memory/learning/patterns.md` — Consolidated lessons
4. Read `_codex/memory/relationship/interaction-notes.md` — How Zeus works

Use this to maintain continuity between sessions.

## Failure Mode
If `_codex/memory/` doesn't exist yet, proceed without it. Not all vaults have the memory tier deployed.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
