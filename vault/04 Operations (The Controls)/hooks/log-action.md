# Hook: LogAction

**Event**: `execute-post`
**Blocking**: No
**Fail Behavior**: Log warning

## Instructions

After completing an action, log what was done:

1. **What**: Brief description of the action taken
2. **Where**: Which file(s) or domain(s) were affected
3. **Why**: Which goal or task this served
4. **Result**: Success/failure/partial

Write to `_codex/memory/work/current-focus.json` to keep task state updated.

For significant actions (file creation, bulk edits, deploys), also note in the session's response to Zeus.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
