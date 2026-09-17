---
name: OpenClaw — Standing Orders
updated: 2026-05-11
---

# Standing Orders

Standing orders grant your agent **permanent operating authority** for defined programs. Instead of prompting for every task, you define programs with clear scope, triggers, and escalation rules — and the agent executes autonomously within those boundaries.

**Official docs:** [docs.openclaw.ai/automation/standing-orders](https://docs.openclaw.ai/automation/standing-orders.md)

> "You own the weekly report. Compile it every Friday, send it, and only escalate if something looks wrong."

---

## Why Standing Orders

**Without standing orders:**
- You must prompt the agent for every task
- The agent sits idle between requests
- Routine work gets forgotten or delayed
- You become the bottleneck

**With standing orders:**
- The agent executes autonomously within defined boundaries
- Routine work happens on schedule without prompting
- You only get involved for exceptions and approvals
- The agent fills idle time productively

---

## How They Work

Standing orders are defined in agent workspace files. Put them directly in `AGENTS.md` (auto-injected every session) so the agent always has them in context. For larger configs, place them in a dedicated file like `standing-orders.md` and reference it from `AGENTS.md`.

Each program specifies:
1. **Scope** — what the agent is authorized to do
2. **Triggers** — when to execute (schedule, event, or condition)
3. **Approval gates** — what requires human sign-off before acting
4. **Escalation rules** — when to stop and ask for help

Auto-injected workspace files: `AGENTS.md`, `SOUL.md`, `TOOLS.md`, `IDENTITY.md`, `USER.md`, `HEARTBEAT.md`, `BOOTSTRAP.md`, `MEMORY.md` — arbitrary files in subdirectories are NOT auto-injected.

---

## Anatomy of a Standing Order

```markdown
## Program: Weekly Status Report

**Authority:** Compile data, generate report, deliver to stakeholders
**Trigger:** Every Friday at 4 PM (enforced via cron job)
**Approval gate:** None for standard reports. Flag anomalies for human review.
**Escalation:** If data source is unavailable or metrics look unusual (>2σ from norm)

### Execution steps

1. Pull metrics from configured sources
2. Compare to prior week and targets
3. Generate report in Reports/weekly/YYYY-MM-DD.md
4. Deliver summary via configured channel
5. Log completion to Agent/Logs/

### What NOT to do

- Do not send reports to external parties
- Do not modify source data
- Do not skip delivery if metrics look bad - report accurately
```

---

## Standing Orders + Cron Jobs

Standing orders define **what** the agent is authorized to do. Cron jobs define **when** it happens.

```
Standing Order: "You own the daily inbox triage"
    ↓
Cron Job (8 AM daily): "Execute inbox triage per standing orders"
    ↓
Agent: Reads standing orders → executes steps → reports results
```

The cron prompt should reference the standing order rather than duplicating it:

```bash
openclaw cron add \
  --name daily-inbox-triage \
  --cron "0 8 * * 1-5" \
  --tz America/New_York \
  --timeout-seconds 300 \
  --announce \
  --channel discord \
  --to "channel:CHANNEL_ID" \
  --message "Execute daily inbox triage per standing orders. Check for new alerts. Parse, categorize, and persist each item. Report summary to owner. Escalate unknowns."
```

---

## Examples

### Content & Social Media (Weekly)

```markdown
## Program: Content & Social Media

**Authority:** Draft content, schedule posts, compile engagement reports
**Approval gate:** All posts require owner review for first 30 days, then standing approval
**Trigger:** Weekly cycle (Monday review → mid-week drafts → Friday brief)

### Weekly cycle
- **Monday:** Review platform metrics and audience engagement
- **Tuesday–Thursday:** Draft social posts, create blog content
- **Friday:** Compile weekly marketing brief → deliver to owner

### Content rules
- Voice must match the brand (see SOUL.md)
- Never identify as AI in public-facing content
- Include metrics when available
```

### System Monitoring (Continuous)

```markdown
## Program: System Monitoring

**Authority:** Check system health, restart services, send alerts
**Approval gate:** Restart services automatically. Escalate if restart fails twice.
**Trigger:** Every heartbeat cycle

### Response matrix

| Condition        | Action                   | Escalate?                |
|------------------|--------------------------|--------------------------|
| Service down     | Restart automatically    | Only if restart fails 2x |
| Disk space < 10% | Alert owner              | Yes                      |
| Stale task > 24h | Remind owner             | No                       |
| Channel offline  | Log and retry next cycle | If offline > 2 hours     |
```

---

## Execute-Verify-Report Pattern

Every task in a standing order should follow this loop:

1. **Execute** — Do the actual work (don't just acknowledge)
2. **Verify** — Confirm the result is correct (file exists, message delivered, data parsed)
3. **Report** — Tell the owner what was done and what was verified

```markdown
### Execution rules

- Every task follows Execute-Verify-Report. No exceptions.
- "I'll do that" is not execution. Do it, then report.
- "Done" without verification is not acceptable. Prove it.
- If execution fails: retry once with adjusted approach.
- If still fails: report failure with diagnosis. Never silently fail.
- Never retry indefinitely — 3 attempts max, then escalate.
```

---

## Best Practices

**Do:**
- Start with narrow authority and expand as trust builds
- Define explicit approval gates for high-risk actions
- Include "What NOT to do" sections — boundaries matter as much as permissions
- Combine with cron jobs for reliable time-based execution
- Review agent logs weekly to verify standing orders are being followed

**Avoid:**
- Granting broad authority on day one ("do whatever you think is best")
- Skipping escalation rules — every program needs a "when to stop and ask" clause
- Assuming the agent will remember verbal instructions — put everything in the file
- Mixing concerns in a single program — separate programs for separate domains
- Forgetting to enforce with cron jobs — standing orders without triggers become suggestions

---

## Our Fleet Application

One Mind agents use standing orders for all recurring programs. Standing orders live in each agent's `AGENTS.md` workspace file:

```
/home/ubuntu/.openclaw/agents/{id}/workspace/AGENTS.md
```

Access via `oc://` addressing:
```bash
openclaw path read oc://legacy/AGENTS.md
```
