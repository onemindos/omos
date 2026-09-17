---
id: so-01-intelligence-scores-20260501
domain: SO
area: "01"
type: reference
status: active
agent: Legacy
tags: [domain/so, area/01, type/reference, status/active, scores/pi-li-wi]
ai_summary: "PI/LI/WI score formulas, weights, DO DATE algorithm, and operating rhythm crons. Extracted from implementation sprint docs before deletion."
updated: 2026-05-01
---

# Intelligence Scores — Reference

> *The formulas behind OneMind's scoring system. Extracted from implementation sprint docs (Jan 2026) before folder deletion.*

---

## Score Overview

| Score | Full Name | Pillar | What It Measures |
|-------|-----------|--------|-----------------|
| **PI** | Performance Index | HP | Sleep, recovery, habits, tasks, energy |
| **LI** | Legacy Index | LE | Family time, legacy tasks, relationships, mentoring, estate |
| **WI** | Wealth Index | GE | Revenue, business tasks, investments, network, skills |
| **UI** | Unified Index | All | Simple average of PI + LI + WI |

Scale: 0–100. Default when no data: **70**.

---

## PI — Performance Index (HP)

| Component | Weight | Source |
|-----------|--------|--------|
| Habit Completion (7-day rolling) | **25%** | lifeos/habits |
| Sleep Quality | **20%** | Oura |
| Recovery/Readiness | **20%** | Oura |
| Task Completion (HP pillar, 7-day) | **20%** | lifeos/tasks |
| Energy Level | **15%** | Oura / state |

**Formula:** `PI = (habits × 0.25) + (sleep × 0.20) + (recovery × 0.20) + (tasks × 0.20) + (energy × 0.15)`

Energy mapping: `high=90`, `medium=70`, `low=50`

---

## LI — Legacy Index (LE)

| Component | Weight |
|-----------|--------|
| Family Time | **25%** |
| Legacy Tasks (LE pillar) | **25%** |
| Relationship Health | **20%** |
| Teaching / Mentoring | **15%** |
| Estate Planning Progress | **15%** |

---

## WI — Wealth Index (GE)

| Component | Weight |
|-----------|--------|
| Revenue Progress | **25%** |
| Business Tasks (GE pillar) | **25%** |
| Investment Growth | **20%** |
| Network Building | **15%** |
| Skill Development | **15%** |

---

## UI — Unified Index

`UI = (PI + LI + WI) / 3`

Balance warning fires when any single score < 60:  
→ `"{LOWEST} needs attention"`

---

## Score History & Trending

- Daily snapshots saved to `score_history` table (`user_id`, `date`, `hpi`, `lei`, `gei`, `unified`)
- Trend logic: compare latest vs oldest in window
  - `diff > +5` → **up**
  - `diff < -5` → **down**
  - otherwise → **stable**
- Anomaly detection: z-score > 2 standard deviations from 7-day mean triggers insight

---

## DO DATE Algorithm — Slot Scoring Rubric

Tasks have a **DUE DATE** (deadline) and a **DO DATE** (optimal work time). The algorithm scores each available time slot 0–100 and returns the top 3.

**Base score: 50**

| Factor | Condition | Points |
|--------|-----------|--------|
| Deadline proximity | ≤1 day until due | +20 |
| Deadline proximity | ≤3 days until due | +10 |
| Deadline proximity | >7 days until due | -5 |
| Energy match | Task=HIGH + Slot=HIGH | +20 |
| Energy match | Task=LOW + Slot=LOW | +15 |
| Energy match | Task=HIGH + Slot=LOW | -15 |
| Path balance | Pillar is underrepresented (<30%) in the day | +15 |
| Duration fit | Slot ≥ task duration | +15 |
| Duration fit | Slot ≥ 80% of task duration | +10 |
| Duration fit | Slot too short | -10 |
| Time preference | 9–11am (morning focus) | +10 |
| Time preference | 2–4pm (afternoon focus) | +5 |

**Energy prediction by hour:**
- 6–10am → HIGH
- 10am–2pm → MEDIUM
- 2–4pm → LOW (post-lunch dip)
- 4–7pm → MEDIUM
- 7pm+ → LOW

**Path balance ideal distribution:** HP 40% / LE 30% / GE 30%

---

## Operating Rhythm — Cron Schedule

| Workflow | Cron | Time |
|----------|------|------|
| Morning Briefing | `0 6 * * *` | 6:00am daily |
| Evening Review | `0 21 * * *` | 9:00pm daily |
| Weekly Planning | `0 18 * * 0` | 6:00pm Sunday |

**Morning Briefing steps:**
1. Pull Oura sleep data
2. Get all scores (PI/LI/WI/UI)
3. Get today's assets (tasks + habits + events)
4. Legacy agent generates briefing (sleep summary → top 3 priorities → habits → calendar → motivational insight)
5. Push notification to device

---

## Cross-Domain Insights Engine

**Runs on:** 30-day rolling window  
**Returns:** Top 10 insights sorted by impact score

**Correlation:** Pearson r — fires insight when `|r| > 0.5` AND `p < 0.05`  
→ Example: Sleep quality yesterday ↔ GEI today

**Anomaly detection:** Z-score > ±2 std deviations from 7-day mean  
→ `"{SCORE} dropped/spiked significantly — X std deviations from normal"`

**Trend detection:** Compare first half vs second half of 7-day window  
→ Fires when change > ±5 points


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
