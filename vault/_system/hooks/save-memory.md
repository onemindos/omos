# Hook: SaveMemory

**Event**: `session-end`
**Blocking**: Yes (must complete before session closes)
**Fail Behavior**: Retry once, then log error

## Instructions

Before ending a session, persist all memory updates:

1. **Work tier**: Update `_codex/memory/work/current-focus.json`
   - Mark completed tasks
   - Update active project status
   - Note any handoff items for next session

2. **Learning tier**: Ensure all signals from this session are written to `signals.jsonl`

3. **Relationship tier**: If any new preferences were discovered during this session, update `_codex/memory/relationship/interaction-notes.md`

## Failure Mode
If write fails, retry once. If still fails, output the memory updates in the response so Zeus can save them manually.
