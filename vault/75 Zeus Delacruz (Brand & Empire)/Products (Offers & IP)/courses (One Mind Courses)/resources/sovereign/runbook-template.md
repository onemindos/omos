# Runbook Template

> *A runbook documents how to execute a specific operational process. Use for anything you do repeatedly.*

---

```yaml
---
type: runbook
domain: [XX]
name: "[Runbook Name]"
trigger: "[What starts this process]"
frequency: [daily | weekly | monthly | on-demand | on-alert]
agent: [AgentName]
owner: Zeus
status: active
updated: YYYY-MM-DD
---
```

---

# [Runbook Name]

## Purpose

[One sentence: what this runbook does and why it exists.]

## Trigger

[What causes this runbook to run? A schedule, an alert, a specific event, a manual decision?]

## Prerequisites

Before running this runbook, confirm:

- [ ] [Prerequisite 1 — system, access, or state required]
- [ ] [Prerequisite 2]
- [ ] [Prerequisite 3]

---

## Steps

### Step 1 — [Step Name]

[What to do. Be specific enough that you could hand this to a new agent and it would execute correctly.]

```bash
# Example command if applicable
command --flag value
```

Expected output: [What success looks like]
If it fails: [What to do if this step errors]

---

### Step 2 — [Step Name]

[Instructions]

Expected output: [What success looks like]
If it fails: [What to do]

---

### Step 3 — [Step Name]

[Instructions]

Expected output: [What success looks like]
If it fails: [What to do]

---

### Step 4 — Verify

Confirm the process completed successfully:

- [ ] [Verification check 1]
- [ ] [Verification check 2]
- [ ] Log entry created at: [path or location]

---

## Rollback

If this runbook needs to be reversed:

[Describe how to undo, or "Not reversible — see incident response plan"]

---

## Notes

[Any standing context, gotchas, edge cases, or agent-specific instructions.]

---

## Change Log

| Date | Change | By |
|------|--------|----|
| YYYY-MM-DD | Initial version | Zeus |

---

*Agent: Follow each step exactly. If an unexpected state is encountered, stop and report to Zeus before continuing.*
