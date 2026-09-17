---
id: overwatch-sitrep-template
title: "SITREP Template"
type: template
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/template, domain/so, status/active, topic/overwatch]
agent: Legacy
zone: SO
---

# SITREP Template — Situation Report

> Standard format for all Overwatch reports to Zeus.
> Use this template for every SITREP — daily, 6h, 1h, or on-demand.
> Copy the template below, fill in all sections, remove unused items.

---

## How to Use

- **Daily (OW-0):** One SITREP at end of day. Keep brief — GREEN domains need only one line.
- **6h (OW-1):** Four SITREPs per day. Focus on flagged domains.
- **1h (OW-2/3):** Every hour. Tight format. Lead with what changed.
- **On-demand:** When Zeus requests or when a trigger fires.
- **Hot Wash:** Immediately post-event — raw, unfiltered, fast.

---

## SITREP Format

---

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
OVERWATCH SITREP — [SERIAL NUMBER]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DATE/TIME:  [YYYY-MM-DD HH:MM]
STATE:      OW-[0/1/2/3]
FROM:       Legacy (Watchkeeper)
TO:         Zeus
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. SITUATION
[One paragraph. What is the current overall state of One Mind operations?
Lead with the most important thing. If all is green, say so clearly.]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

2. DOMAIN STATUS

SO — Operations
  Status: [🟢 GREEN / 🟡 AMBER / 🔴 RED]
  Notes: [Key items, anomalies, or all-clear]

HP — Personal
  Status: [🟢 GREEN / 🟡 AMBER / 🔴 RED]
  Notes: [Key items, anomalies, or all-clear]

LE — Home
  Status: [🟢 GREEN / 🟡 AMBER / 🔴 RED]
  Notes: [Key items, anomalies, or all-clear]

GE — Empire
  Status: [🟢 GREEN / 🟡 AMBER / 🔴 RED]
  Notes: [Key items, anomalies, or all-clear]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

3. ACTIVE FLAGS
[List all current T2+ triggers fired or anomalies under investigation.
If none, write: NONE]

  FLAG-[ID]: [Brief description] — [Status: OPEN/MONITORING/PENDING ZEUS ACTION]
  FLAG-[ID]: [Brief description] — [Status]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

4. PIR STATUS
[Status of top Priority Intelligence Requirements.
Only flag PIRs that have changed since last SITREP.]

  PIR-[ID]: [Answered/Partial/Open] — [One line update]
  PIR-[ID]: [Answered/Partial/Open] — [One line update]
  All other PIRs: [NO CHANGE / list any changes]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

5. ACTIONS TAKEN SINCE LAST SITREP
[What did Overwatch agents do since the last report?
List completed monitoring tasks, alerts sent, analysis performed.]

  - [Action taken] — [Result]
  - [Action taken] — [Result]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

6. ACTIONS PENDING / AWAITING ZEUS
[What is waiting on Zeus? Decisions needed, approvals required, questions.]

  - [Pending item] — [What Zeus needs to do / decide]
  - [Pending item] — [Deadline if applicable]

  If none: NONE

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

7. INTELLIGENCE (NEW INFORMATION)
[Any new intelligence collected since last SITREP — market, personal, operational.
Relevant to active PIRs or upcoming decisions.]

  - [Intel item] — [Source] — [Relevance]

  If none: NONE

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

8. ASSESSMENT
[Legacy's assessment of current state. Honest, direct, no filler.
Is One Mind on track? What's the main risk right now? What requires Zeus attention most?]

[2-5 sentences. BLUF format — most important thing first.]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

9. RECOMMENDED ACTIONS
[What should Zeus do next? Prioritized list.
Legacy recommends — Zeus decides.]

  1. [Action] — [Why / urgency]
  2. [Action] — [Why / urgency]
  3. [Action] — [Why / urgency]

  If no recommendations: HOLDING — AWAIT ZEUS DIRECTION

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

10. NEXT SITREP
  Time: [YYYY-MM-DD HH:MM]
  State: OW-[X]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
END SITREP [SERIAL NUMBER]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Hot Wash Format

Used immediately post-event — raw, fast, unfiltered. No formatting polish needed.

```
HOT WASH — [EVENT NAME]
Time: [YYYY-MM-DD HH:MM] (immediately after event)
Agent: Legacy

WHAT HAPPENED:
[2-3 sentences. Just the facts.]

WHAT WORKED:
- [item]

WHAT DIDN'T:
- [item]

IMMEDIATE ACTION NEEDED:
- [item]

FULL AAR: [scheduled for / pending]
```

---

## AAR Format (After Action Review)

Used within 24h of an OW-2/3 event or at sprint close.

```
AAR — [EVENT/SPRINT NAME]
Date: [YYYY-MM-DD]
Agent: Legacy + [any contributing agents]
Period Covered: [start] → [end]

━━━━━━━━━━━━

1. WHAT WAS PLANNED
[What was the mission? What was expected to happen?]

2. WHAT ACTUALLY HAPPENED
[Factual account of what occurred — no editorializing yet]

3. WHAT WENT RIGHT
[Specific things that worked as intended]
- [item]

4. WHAT WENT WRONG
[Specific things that didn't work — be honest]
- [item]

5. ROOT CAUSE ANALYSIS
[For each major deviation: why did it happen?]
- [Item]: [Root cause]

6. LESSONS LEARNED
[Specific, actionable lessons — not vague observations]
- [Lesson] → [Change to make]

7. SYSTEM UPDATES RECOMMENDED
[What in the vault, protocols, triggers, or ROE should change?]
- [Update] → [File to modify]

8. NEXT OPERATION SETUP
[How does this AAR inform the next sprint or operation?]

━━━━━━━━━━━━
Signed: Legacy / [Date]
```

---

## SITREP Serial Number Convention

Format: `SITREP-YYYYMMDD-NNN`
- YYYYMMDD = date
- NNN = sequential number for that day (001, 002, 003...)

Example: `SITREP-20260504-001`

SITREPs are logged to `watchlist.md` → SITREP Log section after delivery.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
