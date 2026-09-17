---
id: ui-15-20260111053926
title: "Skill: Approval Management"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Skill: Approval Management"
---

# Skill: Approval Management

## Purpose

Manage the Human-in-the-Loop (HITL) approval system. Queue items requiring Zeus's decision, apply auto-approve rules for routine operations, and track approval history.

## Knowledge

### Approval Levels

| Level | Description | Auto-Approve? |
|-------|-------------|---------------|
| L3 | Read operations, routing | Yes |
| L2 | Task creation, metric updates | Yes, then notify |
| L1 | Financial, external comms | No — queue for approval |
| L0 | Deletions, legal, health emergencies | No — immediate escalation |

### Auto-Approve Rules

**Always Auto-Approve:**
- Reading any file in the vault
- Routing requests to pillar agents
- System health checks
- Generating reports and summaries

**Auto-Approve with Notification:**
- Creating Todoist tasks
- Updating memory blocks
- Logging metrics and events
- Drafting content (not publishing)

**Require Explicit Approval:**
- Financial transactions > $100
- External communications (emails, messages)
- Publishing content publicly
- Modifying credentials or access
- Any deletion operation
- Legal document changes
- Health-related recommendations

**Immediate Escalation (L0):**
- Security incidents
- Health emergencies
- Legal threats
- System failures affecting data

### Queue Management

Pending approvals are stored in the Approval Queue memory block:

```markdown
## Pending Approvals

### [ID: 001] — [Title]
- **Requested by:** [Agent]
- **Time:** [Timestamp]
- **Category:** [Financial/External/Publishing/etc]
- **Action:** [What will happen if approved]
- **Context:** [Why this is needed]
- **Options:** [Approve] [Reject] [Modify]

---
```

## Procedures

### Queue an Item for Approval

1. Validate the request requires approval (check auto-approve rules)
2. Generate unique approval ID
3. Create queue entry with full context
4. Save to Approval Queue memory block
5. Notify Zeus via preferred channel
6. Set reminder if not actioned within 24 hours

### Process Approval Response

1. Parse Zeus's decision (approve/reject/modify)
2. If **approved**: Execute the queued action
3. If **rejected**: Log rejection, notify requesting agent
4. If **modify**: Update request per instructions, re-queue if needed
5. Move item to approval history
6. Update approval statistics

### Check Approval Status

1. Query Approval Queue for pending items
2. Return count and summary to requesting agent
3. Flag items over 24 hours old

## Context

- Zeus prefers batched approval requests when possible
- Morning and evening are best times for approval reviews
- Urgent items should include clear justification
- Financial approvals need dollar amount prominently displayed
- External communications should include draft text

## Examples

**Input:** GE Agent requests to send email to investor

**Queue Entry:**
```markdown
### [ID: 047] — Send Investor Update Email
- **Requested by:** GE Agent
- **Time:** 2024-12-23 14:30
- **Category:** External Communication
- **Action:** Send quarterly update email to investor@example.com
- **Context:** Q4 investor update, includes revenue metrics and roadmap
- **Draft Preview:** "Dear [Name], I wanted to share our Q4 progress..."
- **Options:** [Approve] [Reject] [Edit Draft]
```

**Notification:**
```
📋 Approval Needed: Send Investor Update Email

GE Agent wants to send an investor email.
Category: External Communication
Amount: N/A

Reply: "approve 047" | "reject 047" | "edit 047"
```

---

**Input:** HP Agent logs workout (L2 auto-approve)

**Action:** Auto-approve, then notify

**Notification:**
```
✅ Auto-approved: Logged workout (HP Agent)
45 min strength training recorded.
```

## Integration

- **Tools:** `todoist_create`, notification system
- **Memory:** Approval Queue (read/write), System State
- **Escalation:** L0 items bypass queue, go direct to Zeus
