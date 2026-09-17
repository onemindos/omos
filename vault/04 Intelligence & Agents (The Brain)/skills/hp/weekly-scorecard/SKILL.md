---
name: weekly-scorecard
description: "Aggregates task completions, domain health scores, energy trends, wins, and blockers from all 4 domains into a scored weekly review note with next-week priorities. Use when Zeus says 'weekly review', 'scorecard', 'how was my week', 'week in review', or at scheduled Sunday trigger."
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

# Weekly Scorecard — Domain Health & Performance Review

> "What gets measured gets managed. What gets reviewed gets multiplied."

This is the × (Multiply) step of the CODEX method at the weekly cadence.
Aggregates the week's data into a scored report card.

---

## When to Activate

| Signal | Examples |
|--------|----------|
| Direct request | "weekly review", "scorecard", "how was my week" |
| Scheduled trigger | N8N fires Sunday evening |
| Week boundary | "wrap up the week", "week in review" |

---

## Data Collection

### 1. Task Metrics (Todoist)

```bash
# Completed this week (use filter)
curl -s "https://api.todoist.com/rest/v2/tasks?filter=completed this week" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

Calculate per domain:
- Tasks completed vs created (velocity)
- Overdue count (debt)
- Completion rate by project

### 2. Evening Review Aggregation

Read this week's daily notes from vault:
- Average energy score
- Recurring wins/patterns
- Recurring blockers
- Gratitude themes

### 3. Domain-Specific Metrics

| Domain | What to Check |
|--------|--------------|
| **SO** | Agent uptime, system health, inbox processed |
| **HP** | Workouts logged, sleep scores, energy trend |
| **LE** | Family time quality, home tasks completed |
| **GE** | Revenue actions, content published, business development |

---

## Scoring System

Each domain scored 1-10 based on:

| Score | Meaning |
|-------|---------|
| 9-10 | Exceptional — exceeded goals |
| 7-8 | On track — meeting commitments |
| 5-6 | Slipping — needs attention |
| 3-4 | Red flag — significant gaps |
| 1-2 | Crisis — immediate intervention needed |

**Overall Week Score** = weighted average (configurable weights per season)

Default weights: SO 15%, HP 35%, LE 25%, GE 25%

---

## Output Format

Save to: `00-24 SO (Sovereign Operations)/02 Command (Dashboard & Navigation)/reviews/weekly-{YYYY-Www}.md`

```markdown
---
title: "Weekly Scorecard — Week {week_num}, {year}"
type: review
status: active
created: {date}
owner: zeus
tags: [type/review, domain/so-02, cadence/weekly]
ai_summary: "Week {n} scorecard: overall {score}/10 — {one-line summary}"
---

# 📊 Weekly Scorecard — Week {week_num}
**{monday_date} → {sunday_date}**

## Overall Score: {score}/10 {trend_emoji}

| Domain | Score | Trend | Highlights |
|--------|-------|-------|------------|
| ⚡ SO | {n}/10 | {↑↓→} | {one-liner} |
| 💪 HP | {n}/10 | {↑↓→} | {one-liner} |
| 🏡 LE | {n}/10 | {↑↓→} | {one-liner} |
| 💰 GE | {n}/10 | {↑↓→} | {one-liner} |

## 📈 Task Velocity

| Metric | Count |
|--------|-------|
| Completed | {n} |
| Created | {n} |
| Carried forward | {n} |
| Overdue cleared | {n} |
| Net velocity | {+/-n} |

## 🏆 Top 3 Wins
1. {win}
2. {win}
3. {win}

## ⚠️ Top 3 Blockers
1. {blocker}
2. {blocker}
3. {blocker}

## ⚡ Energy Trend
Average: {avg}/10 | High: {day} ({score}) | Low: {day} ({score})

## 🎯 Next Week — Top 3 Priorities
1. {priority — tied to which domain}
2. {priority}
3. {priority}

## 💡 Pattern Detected
{Any recurring theme from this week's daily reviews}
```

---

## Integration Points

- **Evening Reviews**: Primary data source for qualitative scoring
- **Life Score**: Weekly scores feed into monthly PI/WI/LI calculations
- **Quarterly OKR**: Weekly trends inform quarterly adjustments
- **Morning Briefing**: Next week priorities feed into Monday's briefing

---

## Agent Behavior

- Compare to previous week's scores — always show trend direction
- If a domain drops 2+ points week-over-week, flag it prominently
- Celebrate wins before addressing gaps
- Next week priorities should be specific and actionable, not vague
- If data is missing (no evening reviews that week), note the gap itself as a signal
- Keep the full report under 2 minutes reading time


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
