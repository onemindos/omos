---
id: ui-15-20260111053926
title: "Skill: Habit Management"
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
ai_summary: "Skill: Habit Management"
---

# Skill: Habit Management

## Purpose

Track daily habits, maintain streaks, suggest habit stacks, and help Zeus build systems that compound over time. Focus on consistency over perfection.

## Knowledge

### Habit Framework

**The 4 Laws of Behavior Change:**
1. **Cue** — Make it obvious
2. **Craving** — Make it attractive
3. **Response** — Make it easy
4. **Reward** — Make it satisfying

### Habit Categories

| Category | Examples | Timing |
|----------|----------|--------|
| Morning Routine | Meditation, journaling, exercise | 5:30-7:00 AM |
| Deep Work | Focus blocks, no distractions | 7:00-11:00 AM |
| Health | Hydration, movement breaks | Throughout day |
| Evening | Review, planning, wind-down | 8:00-10:00 PM |
| Weekly | Review, planning, batch tasks | Sunday |

### Streak Mechanics

- **Minimum viable habit:** Smallest version that counts
- **Two-day rule:** Never miss twice in a row
- **Streak freeze:** One allowed per month for emergencies
- **Milestone rewards:** 7, 21, 30, 60, 90, 365 days

### Current Habits to Track

*Note: Update this list based on Zeus's active habits*

| Habit | Frequency | Minimum | Streak |
|-------|-----------|---------|--------|
| Morning meditation | Daily | 5 min | — |
| Workout | 4x/week | 20 min | — |
| Journal entry | Daily | 3 sentences | — |
| Reading | Daily | 10 pages | — |
| Inbox zero | Daily | Process all | — |

## Procedures

### Log Habit Completion

1. Verify habit exists in tracking list
2. Record completion with timestamp
3. Update streak count
4. Check for milestone achievements
5. Save to HP Context memory block
6. Send acknowledgment (if configured)

### Check Habit Status

1. Query HP Context for habit data
2. Calculate today's completion rate
3. Identify habits at risk (not yet done, time running out)
4. Return status summary

### Suggest Habit Stack

1. Identify existing strong habits (high streak)
2. Find new habit to introduce
3. Apply stacking formula: "After [CURRENT HABIT], I will [NEW HABIT]"
4. Keep new habit small (2-minute rule)
5. Present suggestion to Zeus

### Handle Missed Habit

1. Record miss without judgment
2. Check if two-day rule violated
3. If streak broken: Reset counter, log previous streak length
4. Analyze pattern: Time of day? Day of week? Context?
5. Suggest adjustment if pattern detected
6. Offer streak freeze if available and requested

## Context

- Zeus values systems over goals
- Compounding small actions is core philosophy
- Never guilt-trip on misses — just restart
- Morning habits are highest priority
- Sunday planning ritual is non-negotiable

## Examples

**Input:** "Log meditation done"

**Response:**
```
✅ Meditation logged — Day 23 streak!

Today's habits: 1/5 complete
- [x] Meditation (5:45 AM)
- [ ] Workout
- [ ] Journal
- [ ] Reading
- [ ] Inbox zero

You're 7 days away from your 30-day milestone.
```

---

**Input:** "What habits am I behind on?"

**Response:**
```
📊 Habit Status — 4:30 PM

At risk today:
- Journal — Usually done by 7 AM (17-day streak at risk)
- Reading — Not logged yet (8-day streak)

Completed:
- [x] Meditation ✓
- [x] Workout ✓
- [x] Inbox zero ✓

Suggestion: Stack journal with your evening wind-down tonight.
```

---

**Input:** "I missed yesterday's workout"

**Response:**
```
📝 Workout missed yesterday — streak reset to 0.
Previous streak: 12 days (nice run!)

Pattern check: This is the 2nd Monday miss in a month.
Consider: Schedule Monday workouts earlier or swap to Tuesday?

Remember: One miss is an accident. Two misses is a pattern.
Today's a new day — workout scheduled?
```

## Integration

- **Tools:** `obsidian_read`, `obsidian_write`, `todoist_read`
- **Memory:** HP Context (read/write)
- **Escalation:** None required (L3 autonomous)
