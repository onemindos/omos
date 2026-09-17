---
name: contact-crm-relationships
description: "Manage personal relationship intelligence in the vault — track contacts, log interactions, monitor relationship health, and surface follow-up reminders. Use when Zeus says \"add contact\", \"log interaction\", \"who haven't I talked to\", \"relationship check\", \"CRM\", or needs to manage personal or professional network."
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Bash(python3:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Contact CRM — Relationship Intelligence Skill

> "The quality of your life equals the quality of your relationships."

Manages personal and professional relationships as structured codex notes in
the LE (Legacy Evolution) domain. This is NOT a sales CRM — it's relationship
intelligence for your foundation.

---

## When to Activate

| Signal | Examples |
|--------|----------|
| Contact management | "add contact", "update {name}", "new person" |
| Interaction logging | "just met {name}", "had coffee with {name}", "log call with {name}" |
| Follow-up queries | "who haven't I talked to in 30 days?", "overdue follow-ups" |
| Relationship review | "relationship check", "network health", "who should I reconnect with?" |
| Lookup | "what do I know about {name}?", "when did I last talk to {name}?" |

---

## Contact Note Structure

Store contacts in: `50-74 LE (Legacy Evolution)/53 Community (Network & Relationships)/contacts/`

### Contact Template

```markdown
---
title: "{Full Name}"
type: human
status: active
created: {date}
updated: {date}
owner: zeus
tags: [type/human, domain/le-53, circle/{inner|mid|outer}]
ai_summary: "{Name} — {relationship context in one line}"
---

# {Full Name}

## Identity
- **Relationship**: {friend, colleague, mentor, client, family, acquaintance}
- **Circle**: {inner (weekly), mid (monthly), outer (quarterly)}
- **Met**: {where/when first met}
- **Location**: {city/country}

## Contact
- **Phone**: {number}
- **Email**: {email}
- **Social**: {handles}

## Context
{What they do, shared interests, how you help each other}

## Interaction Log
| Date | Type | Notes |
|------|------|-------|
| {date} | {call/meeting/text/event} | {brief note} |

## Notes
{Important things to remember — their kids' names, interests, goals, etc.}

## Follow-up
- **Cadence**: {weekly/monthly/quarterly}
- **Next**: {date}
- **Topic**: {what to discuss next}
```

---

## Operations

### Add Contact
```
Input: "add contact John Smith — met at tech conference, works at Acme Corp, AI researcher"
Output: Create contact note with extracted details, set circle to outer, set quarterly cadence
```

### Log Interaction
```
Input: "had lunch with John Smith — discussed his new paper on multi-agent systems"
Output: Append to interaction log, update "updated" date, adjust "Next" follow-up
```

### Follow-up Query
```
Input: "who haven't I talked to in 30 days?"
Output: Scan all contacts, compare last interaction date, return sorted list
```

### Relationship Health Check
```
Input: "relationship check" or triggered monthly
Output: Report showing overdue contacts by circle:
  - Inner circle overdue (>7 days)
  - Mid circle overdue (>30 days)
  - Outer circle overdue (>90 days)
```

---

## Circle System

| Circle | Cadence | Who | Count Target |
|--------|---------|-----|-------------|
| **Inner** | Weekly | Close friends, immediate family, key collaborators | 5-10 people |
| **Mid** | Monthly | Good friends, mentors, active professional contacts | 15-30 people |
| **Outer** | Quarterly | Acquaintances, dormant connections, industry contacts | 50-100 people |

Moving between circles:
- **Promote**: Increased interaction frequency, deepening relationship
- **Demote**: Reduced contact, relationship cooling naturally
- **Archive**: No longer relevant or responsive (set status: inactive)

---

## Integration Points

- **Todoist**: Create follow-up tasks when cadence is overdue
- **Weekly Scorecard**: Relationship score feeds LE domain rating
- **Life Score**: Network health feeds LI calculation
- **Zoho CRM**: For business contacts, can sync with Zoho (use zoho skill)
- **Telegram**: Reminder to reach out when follow-up is overdue

---

## Agent Behavior

- When logging interactions, extract details automatically (topics discussed, action items, sentiment)
- Proactively suggest follow-ups: "You haven't talked to {name} in {n} days. They mentioned {topic} last time."
- During relationship health checks, prioritize inner circle gaps over outer circle
- Never share contact details externally — this is private relationship intelligence
- When adding new contacts, ask clarifying questions if context is sparse
- Respect that some relationships are seasonal — don't nag about contacts Zeus intentionally let drift


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
