---
id: ui-01-20260111053925
title: "Metrics — Pillar Scorecards & Signals"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/01]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Metrics — Pillar Scorecards & Signals"
---

# Metrics — Pillar Scorecards & Signals

> *"What gets measured gets managed. What gets scored gets improved."*

---

## Purpose

Metrics houses the **measurement framework** for the Legacy Codex. It defines how each pillar is scored, what signals to track, and how to calculate the unified Legacy Score.

---

## The Legacy Score

**Formula:** `Legacy Score = Average of Four Pillar Scores`

| Threshold | Status | Action |
|-----------|--------|--------|
| < 60 | 🔴 Red | Simplify — too much, cut scope |
| 60-79 | 🟡 Amber | Stabilize — fix foundations |
| 80+ | 🟢 Green | Scale — expand capacity |

**Target:** Maintain 80+ Legacy Score consistently.

---

## Pillar Scorecards (Weekly 0-100)

### GE (Generational Empire) Scorecard

| Metric | Weight | Measurement |
|--------|--------|-------------|
| Runway | 25% | Months of operating cash |
| Free Cash Flow | 25% | Net cash after expenses |
| Capital Efficiency | 25% | Revenue per $ invested |
| Legal Hygiene | 25% | Compliance items current |

**Calculation:**
- Runway: >12mo = 100, 6-12mo = 75, 3-6mo = 50, <3mo = 25
- FCF: Positive = 100, Break-even = 50, Negative = 25
- Efficiency: Benchmark vs. industry
- Legal: All current = 100, 1 overdue = 75, 2+ overdue = 50

---

### HP (Holistic Performance) Scorecard

| Metric | Weight | Measurement |
|--------|--------|-------------|
| Sleep Quality | 25% | Average sleep score (Oura/Whoop) |
| Recovery/Readiness | 25% | HRV-based readiness |
| Habit Streaks | 25% | Core habits maintained |
| Training Load | 25% | Weekly training consistency |

**Calculation:**
- Sleep: 85+ = 100, 70-84 = 75, 55-69 = 50, <55 = 25
- Recovery: 70+ = 100, 50-69 = 75, 33-49 = 50, <33 = 25
- Habits: 7/7 days = 100, 5-6 = 75, 3-4 = 50, <3 = 25
- Training: 5+ sessions = 100, 4 = 75, 3 = 50, <3 = 25

---

### LE (Legacy Evolution) Scorecard

| Metric | Weight | Measurement |
|--------|--------|-------------|
| Family Rituals | 25% | Weekly rituals completed |
| Homestead Progress | 25% | Projects advancing |
| Wealth Posture Delta | 25% | Net worth change |
| Sovereignty Status | 25% | Key documents current |

**Calculation:**
- Rituals: All kept = 100, 75% kept = 75, 50% = 50, <50% = 25
- Homestead: On track = 100, Minor delay = 75, Major delay = 50, Stalled = 25
- Wealth: +Growth = 100, Stable = 75, -Minor = 50, -Major = 25
- Sovereignty: All current = 100, 1 gap = 75, 2+ gaps = 50

---

### IT (Innovative Technology) Scorecard

| Metric | Weight | Measurement |
|--------|--------|-------------|
| Automation Coverage | 25% | % of repeatable tasks automated |
| MTTR | 25% | Mean time to restore |
| Backup/Restore Pass | 25% | Last drill success |
| Security Findings | 25% | Open vulnerabilities |

**Calculation:**
- Automation: >80% = 100, 60-79% = 75, 40-59% = 50, <40% = 25
- MTTR: <1hr = 100, 1-4hr = 75, 4-24hr = 50, >24hr = 25
- Backup: Pass = 100, Minor issues = 75, Major issues = 50, Fail = 25
- Security: 0 critical = 100, 1-2 = 75, 3-5 = 50, >5 = 25

---

## Signal Types

### Lead Indicators (Predictive)

Signals that predict future outcomes:

| Pillar | Lead Signals |
|--------|--------------|
| GE | Pipeline value, demo calls booked, content published |
| HP | Deep work hours, training sessions, sleep consistency |
| LE | Family time scheduled, homestead hours, savings rate |
| IT | Automation created, security scans run, backups verified |

### Lag Indicators (Results)

Signals that measure outcomes:

| Pillar | Lag Signals |
|--------|-------------|
| GE | Revenue, profit, customer count, entity compliance |
| HP | Body composition, skill certifications, habit streaks |
| LE | Net worth, property value, family milestones |
| IT | Uptime %, incidents resolved, tools deployed |

---

## Scorecard Template

```markdown
# WEEKLY SCORECARD — Week of [DATE]

## Pillar Scores
| Pillar | Score | Δ | Status |
|--------|-------|---|--------|
| GE (Generational Empire) | __/100 | __ | 🟢🟡🔴 |
| HP (Holistic Performance) | __/100 | __ | 🟢🟡🔴 |
| LE (Legacy Evolution) | __/100 | __ | 🟢🟡🔴 |
| IT (Innovative Technology) | __/100 | __ | 🟢🟡🔴 |

## LEGACY SCORE: __/100

## Key Wins
1. 
2. 
3. 

## Blockers
1. 
2. 

## Next Week Focus
1. 
2. 
3. 
```

---

## Dashboard Structure

```
LEGACY COMMAND DASHBOARD
├── Legacy Score (prominent, center)
├── Four Pillar Gauges (0-100 each)
├── Trend Chart (12-week rolling)
├── Lead Indicators (sparklines)
├── Lag Indicators (month-over-month)
└── Action Items (from Red/Amber pillars)
```

---

## Navigation

- Parent: [[Legacy Codex/00-Command/_Command]]
- Related: [[_Doctrine]], [[Legacy Codex/00-Command/Atlas/_Atlas]], [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine)/Operating-System/Operating-Rhythm]]

---

**Version:** 1.0  
**Last Updated:** 2025-12-03
