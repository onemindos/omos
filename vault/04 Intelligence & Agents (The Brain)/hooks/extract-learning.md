# Hook: ExtractLearning

**Event**: `multiply`
**Blocking**: No
**Fail Behavior**: Log warning

## Instructions

During the × Multiply (reflection) phase, extract learnings from completed work:

1. **Review what was done this session**
2. **Identify patterns**:
   - What worked well? → Log as positive pattern
   - What failed or needed revision? → Log as learning
   - What took longer than expected? → Log as efficiency insight
   - What did Zeus praise or criticize? → Log as preference signal
3. **Write insights** to `_codex/memory/learning/signals.jsonl` in JSONL format:
   ```jsonl
   {"ts":"...","type":"implicit","signal":"batch creation was efficient","context":"...","agent":"..."}
   ```
4. **Check if any insight should update `patterns.md`** — if the same pattern appears 3+ times in signals.jsonl, consolidate it.

## Purpose
This is what makes OneMind self-improving. Every session teaches the system something.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
