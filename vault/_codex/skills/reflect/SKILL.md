# Codex Skill: Reflect

How to run reflection and review cycles. Borrowed from Letta's sleep-time compute and ClawVault's reflection patterns, adapted for the CODEX method (step × — Multiply).

---

## Reflection Cadences

### Daily Reflection (5 min)

**When**: End of day or during agent downtime (heartbeat)
**Purpose**: Capture what happened, extract insights, flag unfinished items

**Steps**:
1. Review today's activity (messages, tasks completed, decisions made)
2. Create a daily log using `_codex/templates/daily.md`
3. Extract any **facts** that should become permanent knowledge:
   - Decisions made → create a decision record in the relevant domain
   - New information learned → create a note in the relevant domain
   - New contacts → create a person note
4. Flag **unfinished items** with `status/in-progress` tag
5. Save the daily log to the agent's memory folder or `29 Life Systems/`

### Weekly Review (15 min)

**When**: End of week (Sunday) or when prompted
**Purpose**: Consolidate the week, check goals, plan next week

**Steps**:
1. Read the last 7 daily logs
2. **Consolidate**: What patterns emerged? What was accomplished?
3. **Check goals**: Read `00 Framework/active-goals.md` (if it exists) — are we on track?
4. **Process inbox**: Run `_codex/skills/organize.md` on any items in `06 Inbox (Queue)/`
5. **Flag stale items**: Notes tagged `status/in-progress` older than 7 days → review or close
6. Create a weekly review note using `_codex/templates/review.md`
7. **Plan next week**: What are the top 3 priorities?

### Monthly Review (30 min)

**When**: First weekend of the month
**Purpose**: Bigger picture — are the systems working? What needs to change?

**Steps**:
1. Read the last 4 weekly reviews
2. **Audit each domain group**:
   - UI: Is the infrastructure healthy? Any tech debt?
   - HP: Are health/fitness/learning goals on track?
   - LE: Family & home — anything neglected?
   - GE: Business — revenue, pipeline, growth
3. **Knowledge consolidation**: Run `_codex/skills/consolidate.md`
4. **Vault health**: Run `_codex/skills/maintain.md`
5. Create a monthly review note
6. Update goals and priorities

---

## Sleep-Time Reflection (Agent-Specific)

When an agent has downtime (between conversations, during heartbeat), it can do lightweight reflection:

1. **Scan recent memory** for facts that should be promoted to the codex
2. **Check for contradictions** — does anything in recent memory conflict with codex knowledge?
3. **Extract preferences** — "Zeus said he prefers X" → update the relevant note
4. **Flag knowledge gaps** — "I was asked about X but couldn't find it" → create a stub note tagged `status/gap`

This is the Letta "sleep-time compute" pattern: use idle time to improve the knowledge base.

---

## The Three Types of Reflection Output

| Output | Where It Goes | Example |
| --- | --- | --- |
| **Insight** (new knowledge) | Codex domain folder | "Learned that NATS requires explicit stream creation" |
| **Update** (corrected knowledge) | Edit existing codex note | "VPS IP changed from X to Y" |
| **Action** (something to do) | Task system or inbox | "Need to update the deploy script" |
