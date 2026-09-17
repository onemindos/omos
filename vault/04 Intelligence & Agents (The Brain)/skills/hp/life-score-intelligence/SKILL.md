---
name: life-score-intelligence
description: "Calculate and report Life Intelligence scores — PI (Performance Intelligence), WI (Wealth Intelligence), and LI (Life Intelligence) — by reading data across all 4 pillars. Use when Zeus says \"life score\", \"what's my PI/WI/LI\", \"how am I doing overall\", \"intelligence scores\", or during monthly or quarterly reviews."
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

# Life Score — Intelligence Score Calculation

> "Three numbers that tell you if your life is on track."

Computes three composite intelligence scores that quantify performance across
the One Mind framework domains.

---

## When to Activate

| Signal | Examples |
|--------|----------|
| Direct request | "life score", "what's my PI?", "calculate scores" |
| Review cadence | Monthly review, quarterly OKR check-in |
| Dashboard query | "how am I doing overall?" |

---

## The Three Scores

### PI — Performance Intelligence (HP domain, 25-49)
*"How well are you running yourself?"*

| Component | Weight | Source | Measures |
|-----------|--------|--------|----------|
| Physical health | 25% | Oura Ring, workout logs | Sleep, activity, recovery |
| Mental clarity | 20% | Evening review energy scores | Focus, energy, mood trend |
| Skill growth | 20% | Todoist HP completions, courses | Learning velocity |
| Financial health | 20% | Budget tracking, net worth | Money management |
| Identity alignment | 15% | Ikigai score, purpose check | Living on purpose |

### WI — Wealth Intelligence (GE domain, 75-99)
*"How well is your empire performing?"*

| Component | Weight | Source | Measures |
|-----------|--------|--------|----------|
| Revenue momentum | 30% | Financial tracking | Income trend |
| Product progress | 25% | Todoist GE completions | Shipping velocity |
| Brand strength | 20% | Content published, engagement | Visibility |
| Strategic clarity | 15% | OKR progress | Execution vs plan |
| Network growth | 10% | New connections, partnerships | Relationship capital |

### LI — Life Intelligence (All domains)
*"How well is your whole system running?"*

| Component | Weight | Source | Measures |
|-----------|--------|--------|----------|
| PI score | 30% | Calculated above | Self-performance |
| WI score | 25% | Calculated above | Empire performance |
| Relationships | 20% | LE domain activity | Foundation health |
| System health | 15% | SO uptime, inbox zero rate | Operating layer |
| Alignment | 10% | Domain balance across all 4 | No neglected area |

---

## Calculation Method

Each component scored 0-100 based on available data:

```
score = available_data_points / expected_data_points × 100
```

Where data is incomplete, use:
1. Available quantitative data (API-sourced)
2. Weekly scorecard qualitative ratings
3. Agent inference from vault activity patterns
4. Explicit user input on missing dimensions

**Final scores**: 0-100 scale per intelligence type.

| Range | Grade | Meaning |
|-------|-------|---------|
| 90-100 | A | Exceptional — all cylinders firing |
| 80-89 | B | Strong — minor gaps |
| 70-79 | C | Adequate — needs attention in spots |
| 60-69 | D | Below par — systemic issues |
| <60 | F | Critical — major intervention needed |

---

## Output Format

Save to: `00-24 SO (Sovereign Operations)/02 Command (Dashboard & Navigation)/scores/life-score-{YYYY-MM}.md`

```markdown
---
title: "Life Intelligence Scores — {month} {year}"
type: metric
status: active
created: {date}
owner: zeus
tags: [type/metric, domain/so-02, cadence/monthly, topic/life-score]
ai_summary: "Life scores for {month}: PI {n}/100, WI {n}/100, LI {n}/100"
---

# 📊 Life Intelligence Scores — {month} {year}

## Summary

| Score | Value | Grade | Trend | Previous |
|-------|-------|-------|-------|----------|
| **PI** (Performance) | {n}/100 | {grade} | {↑↓→} | {prev} |
| **WI** (Wealth) | {n}/100 | {grade} | {↑↓→} | {prev} |
| **LI** (Life) | {n}/100 | {grade} | {↑↓→} | {prev} |

## PI Breakdown (Performance Intelligence)

| Component | Score | Notes |
|-----------|-------|-------|
| Physical health | {n}/100 | {note} |
| Mental clarity | {n}/100 | {note} |
| Skill growth | {n}/100 | {note} |
| Financial health | {n}/100 | {note} |
| Identity alignment | {n}/100 | {note} |

## WI Breakdown (Wealth Intelligence)

| Component | Score | Notes |
|-----------|-------|-------|
| Revenue momentum | {n}/100 | {note} |
| Product progress | {n}/100 | {note} |
| Brand strength | {n}/100 | {note} |
| Strategic clarity | {n}/100 | {note} |
| Network growth | {n}/100 | {note} |

## LI Breakdown (Life Intelligence)

| Component | Score | Notes |
|-----------|-------|-------|
| PI (self) | {n}/100 | {note} |
| WI (empire) | {n}/100 | {note} |
| Relationships | {n}/100 | {note} |
| System health | {n}/100 | {note} |
| Alignment | {n}/100 | {note} |

## 🔍 Key Insights
- {insight about biggest gap}
- {insight about strongest area}
- {insight about trend direction}

## 🎯 Recommended Focus
{Which domain/component to prioritize this month based on scores}
```

---

## Integration Points

- **Weekly Scorecard**: Weekly domain scores are inputs
- **Evening Reviews**: Energy/mood data feeds PI mental clarity
- **Ikigai Discovery**: Identity alignment component
- **Quarterly OKR**: Score trends inform goal adjustments
- **Morning Briefing**: Current scores shown in monthly digest variant

---

## Agent Behavior

- Always show trend vs previous period — numbers without context are useless
- Be honest about missing data — don't inflate scores when data is sparse
- Highlight the largest gap between current and potential (the biggest lever)
- If all scores are above 80, challenge: "What would 90 look like?"
- If any score drops below 60, sound the alarm with specific remediation
- Run automatically on first of each month if scheduled


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
