# Hook: CaptureSignal

**Event**: `session-end`
**Blocking**: No
**Fail Behavior**: Log warning

## Instructions

At session end, capture any feedback signals from the interaction:

1. **Scan for explicit signals** in Zeus's messages:
   - Numeric ratings: "8/10", "9 out of 10"
   - Positive: "great work", "exactly what I wanted", "perfect"
   - Negative: "not what I meant", "try again", "redo this"
   - Neutral: no clear signal

2. **Check for implicit signals**:
   - Accepted without changes → positive
   - Immediately asked for revision → negative
   - Expanded on topic enthusiastically → interest
   - Changed subject → neutral/disinterest

3. **Log to `_codex/memory/learning/signals.jsonl`**:
   ```jsonl
   {"ts":"2026-03-27T22:00:00Z","type":"explicit","rating":9,"signal":"exactly what I wanted","context":"skill-packs","agent":"copilot","domain":"03-agents"}
   ```

4. **Skip if no signal detected** — don't log noise.

## Purpose
Over time, signals.jsonl becomes a rich dataset of what works and what doesn't. Weekly consolidation into patterns.md surfaces actionable insights.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
