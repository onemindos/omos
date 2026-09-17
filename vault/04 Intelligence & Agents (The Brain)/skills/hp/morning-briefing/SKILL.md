---
name: morning-briefing
description: "Pulls tasks, calendar, weather, sleep metrics, and vault context from all 4 domains into a scannable morning brief under 60 seconds. Use when Zeus says 'morning briefing', 'what's on today', 'daily brief', 'good morning', or at scheduled morning trigger. Delivers via chat or Telegram."
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

# Morning Briefing — Daily Intelligence Brief

> "Start every day knowing exactly what matters."

Generates a personalized morning report by pulling from all domains and delivering
a focused, actionable briefing.

---

## When to Activate

| Signal | Examples |
|--------|----------|
| Direct request | "morning briefing", "what's on today", "daily brief" |
| Scheduled trigger | N8N workflow fires at configured wake time |
| Morning greeting | "good morning", "let's go" (in morning hours) |

---

## Data Sources

### 1. Tasks (Todoist)
```bash
# Today's tasks across all projects
curl -s "https://api.todoist.com/rest/v2/tasks?filter=today" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"

# Overdue tasks
curl -s "https://api.todoist.com/rest/v2/tasks?filter=overdue" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

### 2. Calendar
- Pull today's events from Zoho Calendar or Google Calendar API
- Flag any meetings in next 2 hours

### 3. Weather
- Location-aware (Home Assistant or weather API)
- Include if outdoor activities are planned

### 4. Vault Context
- Read `_codex/memory/work/` for active agent state
- Check for any pending reviews or approvals
- Scan recent captures in SO 24 Inbox

### 5. Health Metrics (if available)
- Oura Ring sleep score
- Previous day's activity/steps
- Any health alerts

---

## Briefing Template

```markdown
# ☀️ Morning Briefing — {weekday}, {date}

## 🎯 Top 3 Priorities
1. {highest priority from Todoist or manual override}
2. {second priority}
3. {third priority}

## 📋 Today's Tasks ({count})
### SO (Sovereign Operations)
- [ ] {task}

### HP (Holistic Performance)
- [ ] {task}

### LE (Legacy Evolution)
- [ ] {task}

### GE (Generational Empire)
- [ ] {task}

## ⚠️ Overdue ({count})
- {overdue task + age}

## 📅 Calendar
- {time} — {event}

## 🌤️ Weather
{conditions}, {temp}. {impact on plans if any}

## 💤 Sleep & Recovery
Score: {oura_score}/100 | Hours: {hours}
{recommendation based on score}

## 🔔 Pending
- {any pending reviews, approvals, captures to process}

## 💡 Daily Nudge
{One motivational/strategic insight relevant to current goals}
```

---

## Delivery Channels

| Channel | When |
|---------|------|
| Chat response | When asked directly |
| Telegram message | When triggered by N8N schedule |
| Vault daily note | Always — append to daily note in codex |

### Telegram Delivery
```bash
curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
  -d chat_id="${TELEGRAM_CHAT_ID}" \
  -d parse_mode="Markdown" \
  -d text="${briefing_markdown}"
```

---

## Vault Storage

Append briefing to daily note: `00-24 SO (Sovereign Operations)/02 Command (Dashboard & Navigation)/daily/{YYYY-MM-DD}.md`

---

## Agent Behavior

- Keep it scannable — bullet points, not paragraphs
- Flag anomalies (unusual overdue count, low sleep score, packed calendar)
- If sleep score is below 70, suggest lighter schedule
- Top 3 priorities should be curated, not just "first 3 tasks"
- Include a relevant nudge tied to current active goals or quarterly OKRs
- Total output should be readable in under 60 seconds


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
