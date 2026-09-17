---
title: "Kids"
type: human
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
access_level: L2
relation: children
tags: [type/human, domain/LE, status/active]
ai_summary: "Zeus's children. Child-level access to Legacy — age-gated, homework help, schedules, safe content only. Fill in individual names and ages."
---

# Kids

> **Action required:** Fill in names, ages, Telegram numbers if any, and individual context.
> Each child can get their own file (copy `templates/human.md`) once you add their details here.

## Children Registry

| Name | Age | Telegram | Access |
|------|-----|----------|--------|
| {{Kid 1 name}} | {{age}} | {{+1... or none}} | L2 |
| {{Kid 2 name}} | {{age}} | {{+1... or none}} | L2 |

## Access Level

| Field | Value |
|-------|-------|
| Level | L2 — Child |
| Chat access | Yes — age-gated |
| Scope | Homework, schedule, Q&A, family info |
| Blocked | Adult content, financial data, SO/GE agent configs, parental private info |

## Context for Legacy

When a child messages, Legacy should:
- Use their first name
- Keep language age-appropriate (adjust per child's age)
- Be encouraging and patient — not dismissive
- Homework help: allowed (explain, guide — don't do it for them)
- If they ask about parent-private topics: "That's something for Mom/Dad to share with you"
- Emergency: always route to Zeus + wife, give home address if needed

Per-child context:
- **{{Kid 1 name}}** (age {{age}}): {{interests, school, anything relevant}}
- **{{Kid 2 name}}** (age {{age}}): {{interests, school, anything relevant}}

## Relationships

- Parents: [[zeus]], [[wife]]
- Siblings: {{each other}}
- Pets: {{[[dog-name]]}}

## Notes

{{School schedule, pickup times, extracurriculars, medical info}}
