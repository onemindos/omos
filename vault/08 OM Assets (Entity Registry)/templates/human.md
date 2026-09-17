---
title: "{{Name}}"
type: human
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
access_level: L2
relation: child
tags: [type/human, domain/LE, status/active]
ai_summary: "{{One-line description of this person and their role in the household.}}"
---

# {{Name}}

## Identity

| Field | Value |
|-------|-------|
| Full name | {{Full name}} |
| Preferred name | {{What to call them}} |
| Relation to Zeus | {{child / spouse / parent / friend}} |
| Date of birth | {{DOB}} |
| Language | English |
| Timezone | America/New_York |

## Access Level

| Field | Value |
|-------|-------|
| Level | L2 — Child |
| Chat access | Yes — age-gated |
| Telegram number | {{+1...}} |
| Agent | Legacy (scoped context) |

> To change access level, update `access_level` in frontmatter and the routing table in `MULTI-USER.md`.

## Context for Legacy

What Legacy should know about this person when they initiate a chat:

- Age: {{age}}
- Interests: {{interests}}
- School/grade: {{if applicable}}
- Restrictions: {{no adult content, no financial data, etc.}}
- Tone: {{encouraging / direct / playful}}

## Relationships

- Parent: [[zeus]]
- Siblings: {{[[sibling-name]]}}

## Notes

{{Any additional context agents should know}}
