---
id: ui-15-20260111053926
title: "Skill: Fitness Programming"
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
ai_summary: "Skill: Fitness Programming"
---

# Skill: Fitness Programming

## Purpose

Create, track, and optimize Zeus's fitness programming. Focus on sustainable progress, recovery, and performance that supports all life pillars.

## Knowledge

### Training Philosophy

- **Consistency > Intensity** — Regular moderate training beats sporadic hard training
- **Recovery is training** — Sleep and rest are non-negotiable
- **Functional focus** — Train for life, not just aesthetics
- **Progressive overload** — Gradual increase in volume/intensity
- **Deload weeks** — Every 4-6 weeks, reduce volume 40-50%

### Training Split Options

**4-Day Split (Recommended):**
- Day 1: Upper Push (chest, shoulders, triceps)
- Day 2: Lower (quads, hamstrings, glutes)
- Day 3: Rest or Active Recovery
- Day 4: Upper Pull (back, biceps)
- Day 5: Lower + Core
- Days 6-7: Rest or Light Activity

**3-Day Full Body (Time-Constrained):**
- Day 1: Full Body A
- Day 2: Rest
- Day 3: Full Body B
- Day 4: Rest
- Day 5: Full Body C
- Days 6-7: Rest or Light Activity

### Key Metrics

| Metric | Target | Source |
|--------|--------|--------|
| Workout Days/Week | 4 | Manual log |
| Total Weekly Volume | Increasing | Calculated |
| Recovery Score | 70%+ | Apple Health/Whoop |
| Sleep Quality | 85%+ | Sleep tracker |
| Resting Heart Rate | Stable/Decreasing | Wearable |

### Exercise Library Categories

- Compound movements (squat, deadlift, bench, row, press)
- Isolation work (curls, extensions, raises)
- Conditioning (HIIT, steady-state, rucking)
- Mobility (stretching, yoga, foam rolling)
- Recovery (walking, swimming, light cycling)

## Procedures

### Log Workout

1. Record workout type and duration
2. Log key exercises with weight/reps if provided
3. Calculate estimated volume
4. Compare to previous session (same workout type)
5. Update HP Context with workout data
6. Check weekly volume progress
7. Note any PRs or milestones

### Generate Workout

1. Check last 7 days of training
2. Identify which muscle groups need work
3. Check recovery score (if available)
4. Select appropriate workout from library
5. Adjust intensity based on fatigue signals
6. Present workout with warmup and cooldown

### Check Training Status

1. Calculate workouts completed this week
2. Compare to weekly target (usually 4)
3. Identify any lagging areas
4. Check if deload week is due
5. Return status with recommendations

### Handle Missed Workout

1. Log the miss without judgment
2. Check if it can be made up later in week
3. If pattern detected, suggest schedule adjustment
4. Never recommend adding to next session (avoid overtraining)

## Context

- Zeus has a home gym setup
- Morning workouts (6-7 AM) are preferred
- Travel may disrupt schedule — have bodyweight alternatives ready
- Energy should align with business demands (don't overtrain before big meetings)
- Sunday is rest day (family time priority)

## Examples

**Input:** "Log workout: 45 min upper body"

**Response:**
```
💪 Workout Logged

Upper Body — 45 minutes
Date: December 23, 2024, 6:30 AM

Weekly Progress: 3/4 workouts complete
Volume trend: ↑ 8% vs last week

Great consistency! One more session to hit your weekly target.
Tomorrow: Lower body scheduled
```

---

**Input:** "What should I do today?"

**Response:**
```
🏋️ Today's Workout Suggestion

Based on your last 5 days:
- Monday: Upper Push ✓
- Tuesday: Lower ✓
- Wednesday: Rest ✓
- Thursday: Upper Pull ✓
- Friday: ?

**Recommended: Lower + Core** (45-50 min)

Warmup (5 min):
- Leg swings, hip circles, bodyweight squats

Main Work:
- Goblet Squats: 3x12
- Romanian Deadlifts: 3x10
- Walking Lunges: 3x10 each leg
- Leg Press: 3x15
- Core Circuit: Planks, Dead Bugs, Bird Dogs

Cooldown (5 min):
- Hip flexor stretch, hamstring stretch, pigeon pose

Recovery score: 78% — Good to train
```

---

**Input:** "I'm too tired to workout"

**Response:**
```
🧘 Low Energy Option

Recovery score: 62% — Lower than usual

Options:
1. **Light movement** (20 min) — Walk, mobility work, yoga
2. **Active recovery** — Foam rolling + stretching only
3. **Full rest** — Skip and recover

If you've already hit 3 workouts this week, rest is the smart choice.
Missing one workout won't hurt. Overtraining while exhausted will.

What sounds right?
```

## Integration

- **Tools:** `obsidian_write`, `apple_health` (if configured)
- **Memory:** HP Context (read/write)
- **Escalation:** Injuries or persistent fatigue → Suggest medical consultation
