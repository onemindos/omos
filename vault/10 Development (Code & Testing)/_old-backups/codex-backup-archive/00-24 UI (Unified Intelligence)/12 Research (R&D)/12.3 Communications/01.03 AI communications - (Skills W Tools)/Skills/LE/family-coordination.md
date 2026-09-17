---
id: ui-15-20260111053926
title: "Skill: Family Coordination"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Skill: Family Coordination"
---

# Skill: Family Coordination

## Purpose

Coordinate family schedules, events, and communication. Ensure family time is protected, important dates are remembered, and household runs smoothly.

## Knowledge

### Family Calendar Categories

| Category | Examples | Color |
|----------|----------|-------|
| Family Events | Dinners, outings, vacations | Blue |
| Kids Activities | School, sports, lessons | Green |
| Medical | Doctor, dentist, checkups | Red |
| Household | Maintenance, repairs | Orange |
| Traditions | Birthdays, holidays | Purple |

### Key Family Rhythms

| Rhythm | Frequency | Time | Purpose |
|--------|-----------|------|---------|
| Family Dinner | Daily | 6:00 PM | Connection |
| Sunday Planning | Weekly | 10:00 AM | Week ahead |
| Date Night | Bi-weekly | Friday | Relationship |
| Family Meeting | Monthly | 1st Sunday | Big decisions |

### Important Dates

*Note: Update with actual family dates*

| Person | Birthday | Anniversary | Other |
|--------|----------|-------------|-------|
| [Family Member 1] | [Date] | — | — |
| [Family Member 2] | [Date] | — | — |
| [Children] | [Dates] | — | — |

Reminder schedule:
- 2 weeks before: Initial reminder + gift planning
- 1 week before: Confirm celebration plans
- Day before: Final reminder
- Day of: Morning acknowledgment

### Communication Preferences

| Person | Preferred | Avoid | Notes |
|--------|-----------|-------|-------|
| [Spouse/Partner] | Direct, in-person | Text for important things | Evening best |
| [Children] | — | — | Age-appropriate |
| [Parents] | Phone calls | Complex texts | — |

## Procedures

### Schedule Family Event

1. Check family calendar for conflicts
2. Check Zeus's work calendar for conflicts
3. If conflicts exist, suggest alternatives
4. Once confirmed, create calendar event
5. Set appropriate reminders
6. Notify relevant family members (if configured)
7. Update LE Context

### Check Week Ahead

1. Pull all family events for next 7 days
2. Identify any conflicts with work
3. Note preparation needed (gifts, reservations, etc.)
4. Highlight protected family time
5. Return summary

### Remember Important Date

1. Calculate days until event
2. At 2 weeks: Queue gift/planning reminder
3. At 1 week: Check preparation status
4. At 1 day: Final confirmation
5. On day: Morning acknowledgment

### Protect Family Time

1. When scheduling work, check against family commitments
2. Flag conflicts immediately
3. Suggest alternatives that preserve family time
4. Never auto-schedule over family events
5. Sunday is always family-first

## Context

- Family time is non-negotiable (Zeus's Rule #1 for LE)
- Quality over quantity — focused presence matters
- Sunday is sacred family time
- Children's events take priority over most work
- Privacy is paramount — family data never leaves the system

## Examples

**Input:** "What's on the family calendar this week?"

**Response:**
```
👨‍👩‍👧‍👦 Family Week — December 23-29, 2024

**Monday (Dec 23):**
- 6:00 PM: Family dinner (regular)

**Tuesday (Dec 24):**
- 10:00 AM: Christmas Eve prep
- 5:00 PM: Christmas Eve dinner
- 8:00 PM: Family traditions time

**Wednesday (Dec 25):**
- 🎄 Christmas Day
- Morning: Gift exchange
- 2:00 PM: [Extended family gathering]

**Thursday (Dec 26):**
- Rest day — no scheduled events

**Friday (Dec 27):**
- 6:00 PM: Date night scheduled

**Saturday (Dec 28):**
- [Child activity] — 10:00 AM

**Sunday (Dec 29):**
- 10:00 AM: Sunday planning
- Family day — keep open

**Reminders:**
- Gift wrapping needed by Dec 24 AM
- [Relative] birthday is Jan 2 — gift planned?
```

---

**Input:** "Schedule a family day trip for Saturday"

**Response:**
```
📅 Scheduling: Family Day Trip

**Checking Saturday, December 28:**

⚠️ Conflict found:
- 10:00 AM: [Child activity] (2 hours)

**Options:**
1. Day trip after activity (12:00 PM start)
2. Reschedule activity to another day
3. Choose different Saturday

**Recommendation:** Option 1 — Start trip at noon
This keeps the activity commitment and gives 6+ hours for the trip.

What would you prefer?
```

---

**Input:** "When is [family member]'s birthday?"

**Response:**
```
🎂 Birthday Reminder

[Family Member]'s birthday: [Date]
Days until: [X] days

**Preparation Status:**
- Gift: ❓ Not planned yet
- Celebration: ❓ Not scheduled

**Suggested Timeline:**
- Today: Decide on gift idea
- [Date - 7 days]: Order/purchase gift
- [Date - 3 days]: Plan celebration details
- [Date]: Celebration!

Want me to create planning tasks?
```

## Integration

- **Tools:** `calendar`, `todoist_le`, `obsidian_write`
- **Memory:** LE Context (read/write)
- **Escalation:** Medical emergencies → L0 to Zeus immediately
