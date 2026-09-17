---
name: OpenClaw CLI — Commitments
updated: 2026-05-11
---

# `openclaw commitments`

Manages follow-up memories inferred from conversations — opt-in, short-lived reminders created from context.

**Official docs:** [docs.openclaw.ai/cli/commitments](https://docs.openclaw.ai/cli/commitments.md)

---

## Commands

```bash
openclaw commitments                          # Show all pending commitments
openclaw commitments --all                   # All statuses (not just pending)
openclaw commitments --agent <id>            # Filter by agent
openclaw commitments --status <status>       # Filter by status
openclaw commitments --json                  # Machine-readable output

openclaw commitments dismiss <id>            # Remove a commitment
openclaw commitments dismiss <id1> <id2>     # Dismiss multiple
```

---

## Status Values

`pending` · `sent` · `dismissed` · `snoozed` · `expired`

---

## Output Fields (text mode)

- Commitment identifier
- Current status
- Type
- Due timing
- Scope
- Suggested follow-up text

JSON mode additionally includes: storage paths, complete record data.

---

## What Are Commitments?

Short-lived follow-up memories inferred from conversation context. Example: agent notices you said "I'll review that tomorrow" and creates a pending commitment to follow up.

**Official concept docs:** [docs.openclaw.ai/concepts/commitments](https://docs.openclaw.ai/concepts/commitments.md)
